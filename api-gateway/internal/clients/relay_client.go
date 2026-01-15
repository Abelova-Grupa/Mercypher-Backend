package clients

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	relaypb "github.com/Abelova-Grupa/Mercypher/proto/relay"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RelayClient struct {
conn   *grpc.ClientConn
	client relaypb.RelayServiceClient
}

// NewRelayClient cretes a new client to a relay service on the given address.
//
// Note:	The situation is the same as in NewMessageClient code. Even if the
//
//	connection fails or refuses it wont be registered. Only when sending
func NewRelayClient(address string) (*RelayClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("Connection refused: nil")
	}

	client := relaypb.NewRelayServiceClient(conn)

	return &RelayClient{
		conn:   conn,
		client: client,
	}, nil

}

func (c *RelayClient) Close() error {
	return c.conn.Close()
}

// GetMessages is the only endpoint to relay service and should be used when the user
// connects to the gateway, so all accumulated undelivered messages can be emptied
// from sessions storage and user has no need for additional (refresh) requests.
//
// Method returns and array of domain ChatMessages; Parse accordingly!
func (c *RelayClient) GetMessages(userId string) ([]domain.ChatMessage, error) {

	stream, err := c.client.GetMessages(context.Background(), &relaypb.UserId{Id: userId})

	var parsedMessages []domain.ChatMessage

	if err != nil {
		log.Fatalf("failed to get message stream: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF { // Stream closed normally
			return parsedMessages, nil
		}
		if err != nil {	// Stream closed paranormally
			return parsedMessages, err

		}

		parsedMessages = append(parsedMessages, domain.ChatMessage{
			MessageId: msg.GetMessageId(),
			SenderId: msg.GetSenderId(),
			Receiver_id: msg.GetRecipientId(),
			Body: msg.GetBody(),
			Timestamp: msg.GetTimestamp(),
		})
	}

}
