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

	// Loading grpc server
	tlsPort := loadGrpcServerPort()
	// creds := loadTransportCredentials()

	listener, err := net.Listen("tcp", tlsPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpcServer := grpc.NewServer(grpc.Creds(creds))
	grpcServer := grpc.NewServer()
	pb.RegisterSessionServiceServer(grpcServer, server.NewGrpcServer())

	log.Printf("Starting gRPC server on port %v...", tlsPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadGrpcServerPort() string {
	tlsPort := ":" + os.Getenv("SESSION_SERVICE_PORT")
	if tlsPort == ":" {
		tlsPort = ":50055"
	}
	return tlsPort
}
