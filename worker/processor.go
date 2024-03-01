package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/mail"
	"github.com/xiaowuzai/simplebank/util"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

// 任务处理接口
type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Shutdown()
}

// 使用 redis 实现任务处理接口
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
	config util.Config
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	mailer mail.EmailSender,
	config util.Config,
) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("types", task.Type()).
					Str("payload", string(task.Payload())).
					Msg("task processing failed")
			}),
			Logger: &Logger{},
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
		config: config,
	}
}

// 启动任务处理服务
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	// 注册 handler
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

// 关闭服务
func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
