package main

import (
	"context"
	"log"
	"net"

	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/handlers"
	pb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRelayServiceServer
}

func (s *server) SendMessage(context context.Context, message *pb.Message) (*pb.Status, error) {
	return handlers.StoreMessage(message), nil
}
func (s *server) GetMessages(userId *pb.UserId, srv grpc.ServerStreamingServer[pb.Message]) error {
	messages := handlers.GetMessagesForUserId(userId)

	// fmt.Println(messages)
	for i := range messages {
		srv.Send(messages[i])
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000") // port to be set 443 or 8080
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRelayServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
