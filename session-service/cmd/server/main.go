package main

import (
	"log"
	"net"
	"os"

	pb "github.com/Abelova-Grupa/Mercypher/proto/session"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/server"
	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("SESSION_SERVICE_PORT")
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSessionServiceServer(grpcServer, server.NewGrpcServer())

	log.Printf("Starting gRPC server on port %v...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
