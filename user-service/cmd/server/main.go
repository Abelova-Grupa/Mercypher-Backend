package main

import (
	//"context"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/grpc/server"
)

// I will leave this main function as is, so if there is some need for extension we can just add another go routine
func main() {
	go startUserServiceServer()
	// go startSessionServiceClient()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func startUserServiceServer() {
	conn := db.Connect()
	port := config.GetEnv("USER_SERVICE_PORT", "")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, server.NewGrpcServer(conn))

	log.Printf("starting user service grpc server on port %v...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

