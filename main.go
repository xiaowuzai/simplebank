package main

import (
	"context"
	"embed"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xiaowuzai/simplebank/api"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/gapi"
	"github.com/xiaowuzai/simplebank/mail"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/util"
	"github.com/xiaowuzai/simplebank/worker"
)

//go:embed doc/swagger/*
var content embed.FS

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config: %s", err)
	}

	// 开发环境日志输出
	if config.Environment == "local" || config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	dbPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s", err)
	}
	defer dbPool.Close()

	// 升级数据库
	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(dbPool)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store, mailer)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("connot start server")
	}
}

// 升级数据库
func runDBMigration(migrationUrl, dbSource string) {
	m, err := migrate.New(
		migrationUrl,
		dbSource)
	if err != nil {
		log.Fatal().Msgf("create new migrate instance error: %s", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Msgf("migration db error: %s", err)
	}

	log.Info().Msg("db migrated successfully")
}

// 运行异步任务处理器
func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	mailer mail.EmailSender,
) {

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer, config)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	distributor worker.TaskDistributor,
) {

	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().Msgf("cannot create gRPC server: %s", err)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msgf("cannot register handler server: %s", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// embed fs
	fs := http.FileServer(http.FS(content))
	httpMux.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	httpServer := &http.Server{
		Handler: gapi.HttpLogger(httpMux),
		Addr:    config.HTTPServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msgf("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done() // 从通道中读出信号，表示已经关闭

		log.Info().Msg("graceful shutdown HTTP server")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msgf("failed to shutdown HTTP gateway server")
			return err
		}
		log.Info().Msg("HTTP server is stopped")

		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	distributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, store, distributor)
	if err != nil {
		log.Fatal().Msgf("cannot create gRPC server: %s", err)
	}

	// 处理中间件
	grpcLogger := grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(gapi.GrpcLogger))
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	waitGroup.Go(func() error {
		listener, err := net.Listen("tcp", config.GRPCServerAddress)
		if err != nil {
			log.Fatal().Msgf("cannot create listener err: %s", err)
		}
		defer listener.Close()

		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
		if err = grpcServer.Serve(listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("cannot start grpc server")
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done() // 从通道中读出信号，表示已经关闭

		log.Info().Msg("graceful shutdown gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot new server: %s", err)
	}

	log.Printf("start HTTP server at %s", config.HTTPServerAddress)
	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal().Msgf("cannot start server: %s", err)
	}
}
