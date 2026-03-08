package servers

import (
	"context"
	"log"
	"time"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
)

type KafkaConsumer struct {
	gwChan chan *domain.ChatMessage
	reader *kafka.Reader
}

type AzureBusConsumer struct {
	gwChan chan *domain.ChatMessage
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, outChan chan *domain.ChatMessage) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		reader: reader,
		gwChan: outChan,
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

func (c *KafkaConsumer) StartLiveForwarder(ctx context.Context) {

	log.Println("Started reading from Kafka.")

	defer c.Close()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			break
		}

		var protoMsg pb.ChatMessage
		var msg domain.ChatMessage

		if err := proto.Unmarshal(m.Value, &protoMsg); err != nil {
			log.Printf("Kafka read error: %v", err)
		}

		msg = domain.ChatMessage{
			MessageId:   protoMsg.Id,
			SenderId:    protoMsg.SenderId,
			Receiver_id: protoMsg.RecieverId,
			Body:        protoMsg.Body,
		}

		log.Printf("Message received from Kafka: %s", msg.Body)

		c.gwChan <- &msg

	}
}

type BusReceiver struct {
	GwChan chan *domain.ChatMessage
	Receiver *azservicebus.Receiver
}

func NewBusReceiver(receiver *azservicebus.Receiver, gwChan chan *domain.ChatMessage) *BusReceiver {
	return &BusReceiver{
		Receiver: receiver,
		GwChan: gwChan,
	}
}

func (b *BusReceiver) Complete(ctx context.Context, msg *azservicebus.ReceivedMessage) error {
	return b.Receiver.CompleteMessage(ctx, msg, nil)
}

func (b *BusReceiver) Start(ctx context.Context) {
	log.Println("Started reading from Azure Service Bus for persistence.")

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
			var msg domain.ChatMessage
			msg = domain.ChatMessage{
			MessageId:   protoMsg.Id,
			SenderId:    protoMsg.SenderId,
			Receiver_id: protoMsg.RecieverId,
			Body:        protoMsg.Body,
			}

			log.Printf("Message received from Azure Service bus: %s", msg.Body)

			b.GwChan <- &msg

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
