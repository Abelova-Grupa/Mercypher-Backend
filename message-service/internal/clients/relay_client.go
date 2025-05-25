package relay_client

import (
	"context"
	"log"

	relaypb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Should be moved somewhere else and run somewhere else
func ConnectRelayService() relaypb.RelayServiceClient {
	conn, err := grpc.NewClient("localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return relaypb.NewRelayServiceClient(conn)
}

func modelMessage2RelayMessage(msg *model.ChatMessage) relaypb.ChatMessage {
	return relaypb.ChatMessage{
		MessageId:   msg.Message_id,
		SenderId:    msg.Sender_id,
		RecipientId: msg.Receiver_id,
		Timestamp:   12345, // to be fixed
		Body:        msg.Body,
	}
}

func RelayMessage(ctx context.Context, msg *model.ChatMessage) error {
	// To be continued
	conn := ConnectRelayService()
	relayMessage := modelMessage2RelayMessage(msg)
	_, err := conn.SendMessage(ctx, &relayMessage)
	return err
}
