package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xiaowuzai/simplebank/api"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/gapi"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/util"
)

//go:embed doc/swagger/*
var content embed.FS

func main() {
	// log.Println(util.NewPasetoSymmetricKey())

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	fmt.Printf("config %v\n", config)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC server: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	// grpcMux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	// embed fs
	fs := http.FileServer(http.FS(content))
	httpMux.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener err: ", err)
	}
	defer listener.Close()

	log.Printf("start gateway server at %s", listener.Addr().String())
	if err = http.Serve(listener, httpMux); err != nil {
		log.Fatal("cannot start gateway server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gRPC server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener err: ", err)
	}
	defer listener.Close()

	log.Printf("start gRPC server at %s", listener.Addr().String())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start grpc server: ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot new server: ", err)
	}

	log.Printf("start HTTP server at %s", config.HTTPServerAddress)
	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
