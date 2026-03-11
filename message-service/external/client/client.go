package client

import (
	"context"
	"fmt"

	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MessageClient struct {
	conn    *grpc.ClientConn
	service pb.MessageServiceClient
}

// address string tipa "localhost:50051"
func NewMessageClient(address string) (*MessageClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client connection: %w", err)
	}

	client := pb.NewMessageServiceClient(conn)

	return &MessageClient{
		conn:    conn,
		service: client,
	}, nil
}

func (c *MessageClient) Close() error {
	return c.conn.Close()
}

func (c *MessageClient) SendMessageParts(ctx context.Context, senderId string, receiverId string, body string) (*pb.MessageAck, error) {
	req := &pb.ChatMessage{
		SenderId:   senderId,
		RecieverId: receiverId,
		Body:       body,
	}
	return c.service.SendMessage(ctx, req)
}
func (c *MessageClient) SendMessageWhole(ctx context.Context, req *pb.ChatMessage) (*pb.MessageAck, error) {
	return c.service.SendMessage(ctx, req)
}

func (c *MessageClient) GetMessages(ctx context.Context, p1, p2 string, lastSeen int64, limit int64) ([]*pb.ChatMessage, error) {
	req := &pb.MessageRange{
		Participant1: p1,
		Participant2: p2,
		LastSeen:     lastSeen,
		Limit:        limit,
	}

	resp, err := c.service.GetMessages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("rpc error: %w", err)
	}

	return resp.Messages, nil
}
