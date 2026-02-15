package main

import (
	//"context"

	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/grpc/server"
	worker "github.com/Abelova-Grupa/Mercypher/user-service/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// I will leave this main function as is, so if there is some need for extension we can just add another go routine
func main() {
	if err := config.LoadEnv(); err != nil {
		panic(err)
	}

	redisOpt := asynq.RedisClientOpt{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	}
	go runEmailTaskProcessor(redisOpt)
	go startUserServiceServer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func startUserServiceServer() {
	conn := db.Connect()
	port := config.GetEnv("USER_SERVICE_PORT", "")
	//Testing ci
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen to start user service")
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, server.NewGrpcServer(conn))

	log.Printf("starting user service grpc server on port %v...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to server grpc request")
	}
}

func runEmailTaskProcessor(redisOpt asynq.RedisClientOpt) {
	taskProcessor := worker.NewRedistaskProcessor(redisOpt)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
