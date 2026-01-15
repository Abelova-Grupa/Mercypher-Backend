package relay_client

import (
	"context"
	"log"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	relaypb "github.com/Abelova-Grupa/Mercypher/proto/relay"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var relayService relaypb.RelayServiceClient

func modelMessage2RelayMessage(msg *model.ChatMessage) relaypb.ChatMessage {
	return relaypb.ChatMessage{
		MessageId:   msg.Message_id,
		SenderId:    msg.Sender_id,
		RecipientId: msg.Receiver_id,
		Timestamp:   msg.Timestamp.Unix(), // to be fixed
		Body:        msg.Body,
	}
}

func RelayMessage(ctx context.Context, msg *model.ChatMessage) error {
	// To be continued
	relayMessage := modelMessage2RelayMessage(msg)
	_, err := relayService.SendMessage(ctx, &relayMessage)
	return err
}

func StartClient(channel chan func()) {
	// Connecting to relay service
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	relayService = relaypb.NewRelayServiceClient(conn)

	// waiting for command
	for function := range channel {
		function()
	}
}
