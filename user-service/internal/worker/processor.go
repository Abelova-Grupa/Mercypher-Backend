package worker

import (
	"context"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db *gorm.DB
	emailSender *email.EmailAuth
}

func NewRedistaskProcessor(redisOpt asynq.RedisClientOpt) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			ErrorHandler:  asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error){
				log.Error().Err(err).Str("type ", task.Type()).Bytes("payload",task.Payload()).Msg("task failed to process")
			}),
			Logger: NewLogger(),
		},

	)
	return &RedisTaskProcessor {
		server: server,
		db: db.Connect(),
		emailSender: email.NewEmailAuth(),
	}
}

func (proc *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, proc.ProcessTaskSendVerifyEmail)
	return proc.server.Start(mux)
}