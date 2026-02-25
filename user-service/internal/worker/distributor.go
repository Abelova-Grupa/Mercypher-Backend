package worker

import (
	"context"
	"fmt"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *email.EmailPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}


func NewRedisTaskDistributor(cfg interface{}) TaskDistributor {
	var client *asynq.Client
	switch v := cfg.(type) {
	case asynq.RedisClientOpt:
		redisOpt := cfg.(asynq.RedisClientOpt)
		client = asynq.NewClient(redisOpt)
	case *asynq.Client:
		client = cfg.(*asynq.Client)
	default:
		log.Error().Msg(fmt.Sprintf("unable to initialize task distributor, non supported type %v", v))
		return nil
	}

	return &RedisTaskDistributor{
		client: client,
	}
}