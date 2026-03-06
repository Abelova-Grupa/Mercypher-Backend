package eventbus

import (
	"context"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type BusReceiver struct {
	Receiver *azservicebus.Receiver
	Repo     *repository.MessageRepo
}

func NewBusReceiver(repo *repository.MessageRepo, receiver *azservicebus.Receiver) *BusReceiver {
	return &BusReceiver{
		Receiver: receiver,
		Repo:     repo,
	}
}

func (b *BusReceiver) Complete(ctx context.Context, msg *azservicebus.ReceivedMessage) error {
	return b.Receiver.CompleteMessage(ctx, msg, nil)
}

func (b *BusReceiver) Start(ctx context.Context) {
	log.Info().Msg("Started reading from Azure Service Bus for persistence.")

	defer func() {
		if err := b.Receiver.Close(ctx); err != nil {
			log.Printf("Error closing Service Bus receiver: %v", err)
		}
	}()

	for {
		messages, err := b.Receiver.ReceiveMessages(ctx, 32, nil)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Service Bus read error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, message := range messages {
			var protoMsg pb.ChatMessage
			if err := proto.Unmarshal(message.Body, &protoMsg); err != nil {
				log.Printf("Unmarshal error: %v", err)
				// Sending malformed messages do dead letter queue
				_ = b.Receiver.DeadLetterMessage(ctx, message, nil)
				continue
			}

			msg := &model.ChatMessage{
				Message_id:  protoMsg.Id,
				Sender_id:   protoMsg.SenderId,
				Receiver_id: protoMsg.RecieverId,
				Body:        protoMsg.Body,
				Timestamp:   time.Unix(protoMsg.Timestamp, 0),
			}

			if err := b.Repo.SaveMessage(ctx, msg); err != nil {
				log.Printf("Failed to save message to DB: %v", err)
				continue
			}

			if err := b.Complete(ctx, message); err != nil {
				log.Printf("Failed to complete message: %v", err)
			}

		}

	}

}

func (b *BusReceiver) Close(ctx context.Context) error {
	return b.Receiver.Close(ctx)
}

func CreateQueueReceiver(azureCli *azservicebus.Client, queueName string, opt *azservicebus.ReceiverOptions) (*azservicebus.Receiver, error) {
	return azureCli.NewReceiverForQueue(queueName, opt)
}

func CreateTopicReceiver(azureCli *azservicebus.Client, topicName string, subsName string, opt *azservicebus.ReceiverOptions) (*azservicebus.Receiver, error) {
	return azureCli.NewReceiverForSubscription(topicName, subsName, opt)
}
