package eventbus

import (
	"context"
	"fmt"

	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"google.golang.org/protobuf/proto"
)

type BusSender struct {
	BusName string
	Sender  *azservicebus.Sender
}

func NewBusSender(BusName string) *BusSender {
	return &BusSender{
		BusName: BusName,
	}
}

func (b *BusSender) CreateSender(azureCli *azservicebus.Client) error {
	sender, err := azureCli.NewSender(b.BusName, &azservicebus.NewSenderOptions{})
	if err != nil {
		return err
	}
	b.Sender = sender
	return nil
}

func (b *BusSender) SendMessage(msg *pb.ChatMessage) (string, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto: %w", err)
	}

	sbMessage := &azservicebus.Message{
		Body: data,
	}

	if err = b.Sender.SendMessage(context.TODO(), sbMessage, nil); err != nil {
		return "", fmt.Errorf("unable to write to azure service bus")
	}
	return msg.Id, nil
}

func (b *BusSender) Close(ctx context.Context) error {
	return b.Sender.Close(ctx)
}

// func (b *BusSender) SendMessageBatch(msgs []*pb.ChatMessage) ([]string, error)
