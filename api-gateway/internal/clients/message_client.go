package clients

import (
	"context"
	"errors"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	messagepb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MessageClient struct {
	conn   *grpc.ClientConn
	client messagepb.MessageServiceClient
}

// TODO: Fix GRPC Client connection connecting to non-existent service madness.

// NewMessageClient cretes a new client to a message service on the given address.
//
// Note:	For some inhumane, ungodly and barbaric reason, grpc.NewClient does not verify the
//
//	connection immediately. It returns a *ClientConn regardless of whether
//	the server exists — errors only show up when you the connection is used.
//	For development purposes, this will work, yet I will be looking for
//	a soulution and implement it asap.
func NewMessageClient(address string) (*MessageClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("Connection refused: nil")
	}

	client := messagepb.NewMessageServiceClient(conn)

	return &MessageClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *MessageClient) Close() error {
	return c.conn.Close()
}

// SendMessage accepts a domain message struct, parses it to grpc format
// and sends it to the message service.
//
// Note: Watch out for errors for they might be associated with a bad conn.
func (c *MessageClient) SendMessage(msg domain.ChatMessage) error {
	var grpcMsg = &messagepb.ChatMessage{
		SenderId:   msg.SenderId,
		RecieverId: msg.Receiver_id,
		Body:       msg.Body,
		Timestamp:  msg.Timestamp,
	}
	_, err := c.client.SendMessage(context.Background(), grpcMsg)

	return err
}

func (c *MessageClient) GetMessages(participant1 string, participant2 string, limit int64) ([]domain.ChatMessage, error) {
	messages, _ := c.client.GetMessages(context.Background(), &messagepb.MessageRange{
		Participant1: participant1,
		Participant2: participant2,
		Limit: limit,
	})

	domainMessages := make([]domain.ChatMessage, 0)

	for _, msg := range messages.Messages {
		domainMessages = append(domainMessages, domain.ChatMessage{
			MessageId: msg.Id,
			SenderId: msg.SenderId,
			Receiver_id: msg.RecieverId,
			Body: msg.Body,
		})
	}

	return domainMessages, nil
}

// TODO: Implement status
