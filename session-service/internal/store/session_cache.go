package store

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/config"
	entraid "github.com/redis/go-redis-entraid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidEnvVars = errors.New("invalid env variables for redis client")
)

func NewSessionCache(ctx context.Context) *redis.Client {
	var redisCli *redis.Client
	err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	if os.Getenv("ENVIRONMENT") == "azure" {
		redisCli = NewSessionCacheAzure(ctx)
	}else{
		redisCli = NewSessionCacheLocal(ctx)
	}
	return redisCli
}

func NewSessionCacheLocal(ctx context.Context) *redis.Client {
	
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
func NewSessionCacheAzure(ctx context.Context) *redis.Client {
	redisHost := net.JoinHostPort(
		os.Getenv("AZURE_REDIS_CACHE_URL"),
		os.Getenv("AZURE_REDIS_CACHE_PORT_TLS"),
	)
	provider, err := entraid.NewDefaultAzureCredentialsProvider(entraid.DefaultAzureCredentialsProviderOptions{})

	if err != nil {
		log.Error().Msg("unable to connect to azure cache for redis")
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
		StreamingCredentialsProvider: provider,
		Username: os.Getenv("USER_OBJECT_ID"),
		TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12, 
    },
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Error().Msg("Could not ping azure cache for redis")
	}
	log.Info().Msg("successfully connected to azure cache for redis")
	return client
}