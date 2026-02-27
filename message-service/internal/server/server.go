package server

import (
	"context"
	"log"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/kafka"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MessageServer handles incoming gRPC requests implementing protobuf
type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	brokers []string
	repo    repository.MessageRepository
}

func NewMessageServer(brokers []string, repo repository.MessageRepository) *MessageServer {
	return &MessageServer{
		brokers: brokers,
		repo:    repo,
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
	lastSeen := time.Unix(req.LastSeen, 0)
	if req.Limit < 1 {
		req.Limit = 20
	}

	messages, err := s.repo.GetChatHistory(ctx, req.Participant1, req.Participant2, lastSeen, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch history: %v", err)
	}

	// 4. Map repository models to Protobuf models
	var pbMessages []*pb.ChatMessage
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.ChatMessage{
			Id:         m.Message_id,
			SenderId:   m.Sender_id,
			RecieverId: m.Receiver_id,
			Body:       m.Body,
			Timestamp:  m.Timestamp.Unix(),
		})
	}

	return &pb.MessageList{
		Messages: pbMessages,
	}, nil
}
