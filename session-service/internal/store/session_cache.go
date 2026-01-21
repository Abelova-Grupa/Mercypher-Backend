package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidEnvVars = errors.New("invalid env variables for redis client")
)

func NewSessionCache(ctx context.Context) *redis.Client {
	err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	redisUser := config.GetEnv("REDIS_USER", "")
	redisPass := config.GetEnv("REDIS_PASSWORD", "")
	redisHost := config.GetEnv("REDIS_HOST", "")
	redisDB := config.GetEnv("REDIS_DB", "")

	if redisUser == "" || redisPass == "" || redisHost == "" || redisDB == "" {
		panic(ErrInvalidEnvVars)
	}

	redis_str := fmt.Sprintf("redis://%s:%s@%s/%s", redisUser, redisPass, redisHost, redisDB)
	log.Print(redis_str)
	opt, err := redis.ParseURL(redis_str)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("successfuly connected to session cache")
	rdb := redis.NewClient(opt)
	err = rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
	return rdb
}
