package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

	if redisHost == "" || redisDB == "" {
		panic(ErrInvalidEnvVars)
	}

	if redisUser == "" || redisPass == "" {
		dbNum, _ := strconv.Atoi(redisDB)
		return redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: "",
			DB:       dbNum,
		})

	}
	redis_str := fmt.Sprintf("redis://%s:%s@%s/%s", redisUser, redisPass, redisHost, redisDB)
	log.Print(redis_str)

	opt, err := redis.ParseURL(redis_str)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	log.Info().Msg("successfuly connected to session cache")

	return rdb
}
