package main

import (
	"log"
	"net"

	server "github.com/Abelova-Grupa/Mercypher/api/internal/server"
	"google.golang.org/grpc"
	pb "github.com/Abelova-Grupa/Mercypher/api/internal/grpc"
)

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGatewayServiceServer(grpcServer, &server.GatewayServer{})

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func main() {

	go startGRPCServer()

	server := server.InitServer()
	log.Fatal(server.Start(":8080"))
}
