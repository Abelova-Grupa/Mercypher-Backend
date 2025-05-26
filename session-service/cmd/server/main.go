package main

import (
	"log"
	"net"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/db"
	pb "github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/pb"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("../../internal/certs/server.crt", "../../internal/certs/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterSessionServiceServer(grpcServer, server.NewGrpcServer(db.Connect(db.GetDBUrl())))

	log.Println("Starting gRPC server on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
