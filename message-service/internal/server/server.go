package server

import (
	"context"
	"log"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/kafka"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MessageServer handles incoming gRPC requests implementing protobuf
type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	brokers []string
}

func NewMessageServer(brokers []string) *MessageServer {
	return &MessageServer{
		brokers: brokers,
	}
}

func (s *MessageServer) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.MessageAck, error) {
	// bare minimum checks
	if req.Body == "" || req.RecieverId == "" || req.SenderId == "" {
		return nil, status.Error(codes.InvalidArgument, "body and recipient_id and sender_id are required")
	}

	req.Id = uuid.New().String()
	// when is timestamp added? here maybe?
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}

	log.Printf("Queueing message from %s to %s", req.SenderId, req.RecieverId)

	generatedID, err := kafka.PublishMessage(ctx, s.brokers, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to queue: %v", err)
	}

	return &pb.MessageAck{
		MessageId: generatedID,
	}, nil
}

func (s *MessageServer) GetMessages(ctx context.Context, req *pb.MessageRange) (*pb.MessageList, error) {
	log.Printf("Fetching messages from %d to %d", req.From, req.To)

	msgs, err := kafka.FetchMessages(ctx, s.brokers, req.From, req.To)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch: %v", err)
	}
	return &pb.MessageList{Messages: msgs}, nil
}
