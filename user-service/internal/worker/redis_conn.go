package worker

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/hibiken/asynq"
	entraid "github.com/redis/go-redis-entraid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type TaskAsynq interface {
	RunTaskProcessor()
	NewTaskProcessor() TaskProcessor
	NewTaskDistributor() TaskDistributor
}

type AzureTaskAsynq struct {
	cli redis.UniversalClient
}

func (p *AzureTaskAsynq) RunTaskProcessor() {
	taskProcessor := p.NewTaskProcessor()
	log.Info().Msg("start task processor on azure")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func (p *AzureTaskAsynq) NewTaskProcessor() TaskProcessor{
	provider, err := entraid.NewDefaultAzureCredentialsProvider(entraid.DefaultAzureCredentialsProviderOptions{})
	if err != nil {
		log.Error().Msg("unable to create default azure credentials provider")
	}
	p.cli = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string {os.Getenv("AZURE_REDIS_CACHE_URL") + ":" + os.Getenv("AZURE_REDIS_CACHE_PORT_TLS")},
		StreamingCredentialsProvider: provider,
		Username: os.Getenv("USER_OBJECT_ID"),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

	conf := asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type ", task.Type()).Bytes("payload", task.Payload()).Msg("task failed to process")
		}),
		Logger: NewLogger(),
	}

	asynqCli := asynq.NewServerFromRedisClient(p.cli, conf)
	return NewRedistaskProcessor(asynqCli,conf)
}

func (p *AzureTaskAsynq) NewTaskDistributor() TaskDistributor {
	provider, err := entraid.NewDefaultAzureCredentialsProvider(entraid.DefaultAzureCredentialsProviderOptions{})
	if err != nil {
		log.Error().Msg("unable to create default azure credentials provider")
	}
	p.cli = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string {os.Getenv("AZURE_REDIS_CACHE_URL") + ":" + os.Getenv("AZURE_REDIS_CACHE_PORT_TLS")},
		StreamingCredentialsProvider: provider,
		Username: os.Getenv("USER_OBJECT_ID"),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

	asynqCli := asynq.NewClientFromRedisClient(p.cli)
	return NewRedisTaskDistributor(asynqCli)
}

type LocalTaskAsynq struct {
	cli asynq.RedisClientOpt
}

func (p *LocalTaskAsynq) RunTaskProcessor() {
	taskProcessor := p.NewTaskProcessor()
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func (p *LocalTaskAsynq) NewTaskProcessor() TaskProcessor {
	redisOpt := asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	}

	conf := asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type ", task.Type()).Bytes("payload", task.Payload()).Msg("task failed to process")
		}),
		Logger: NewLogger(),
	}
	return NewRedistaskProcessor(redisOpt,conf)
}

func (p *LocalTaskAsynq) NewTaskDistributor() TaskDistributor {
	redisOpt := asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	}
	return NewRedisTaskDistributor(redisOpt)
}
