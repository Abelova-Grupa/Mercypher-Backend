package worker

import (
	"context"
	"fmt"

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

func NewRedistaskProcessor(cfg interface{}, opts ...asynq.Config) TaskProcessor {
	var server *asynq.Server
	switch v := cfg.(type){
	case asynq.RedisClientOpt:
		redisOpt := cfg.(asynq.RedisClientOpt)
		server = asynq.NewServer(
		redisOpt,
		opts[0],
	)
	case *asynq.Server:
		server = cfg.(*asynq.Server)

	default:
		log.Error().Msg(fmt.Sprintf("unable to initialize redis task processor with undefined type %v",v))
	}

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