package server

import (
	"context"
	"log"
	"net"

	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/config"
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

	for i := range messages {
		srv.Send(messages[i])
	}

	return nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":"+config.GetEnv("RELAY_SERVICE_PORT", "9001"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRelayServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
