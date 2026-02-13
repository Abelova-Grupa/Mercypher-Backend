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

// range za sada ide normalno 1, 2, 3 -> [1, 3] gde je 0 najstarija
// BITNO: posto je getmessages  spora funkcija, ctx sa timeout-om moze da presece pre odgovora
// moguce da zavisi od broja poruka koje se citaju ali meni sa 5 sekundi pukne ali 20 je okej
func (c *MessageClient) GetMessages(ctx context.Context, from int64, to int64) ([]*pb.ChatMessage, error) {
	req := &pb.MessageRange{
		From: from,
		To:   to,
	}

	resp, err := c.service.GetMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}