package main

import (
	"log"
	"net"

	server "github.com/Abelova-Grupa/Mercypher/api/internal/servers"
	"google.golang.org/grpc"
	pb "github.com/Abelova-Grupa/Mercypher/api/internal/grpc"
)

// func startHTTPServer() {
// 	httpServer := server.NewServer()
// 	log.Fatal(httpServer.Start(":8080"))
// }

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGatewayServiceServer(grpcServer, &server.GrpcServer{})

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func main() {

	if httpServer := server.NewHttpServer(); httpServer == nil {
		log.Fatal("Couldn't start http server")
	}
	go startGRPCServer()

	for {
		
	}

}