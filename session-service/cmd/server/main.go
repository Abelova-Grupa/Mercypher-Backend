package main

import (
	"log"
	"net"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/db"
	pb "github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/pb"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/server"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterSessionServiceServer(s, server.NewGrpcServer(db.Connect(db.GetDBUrl())))

	log.Println("Starting gRPC server on port 50052...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
