package kafka

import (
	"context"
	"log"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	repo   *repository.MessageRepo
}

func NewKafkaConsumer(repo *repository.MessageRepo, brokers []string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    "chat-messages-v1",
		GroupID:  "message-service-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		reader: reader,
		repo:   repo,
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Started reading from Kafka for persistence.")
	defer c.Close()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var protoMsg pb.ChatMessage
		if err := proto.Unmarshal(m.Value, &protoMsg); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		// Mapping to the model.ChatMessage used by your repo
		msg := &model.ChatMessage{
			Message_id:  protoMsg.Id,
			Sender_id:   protoMsg.SenderId,
			Receiver_id: protoMsg.RecieverId, // Matching proto field naming
			Body:        protoMsg.Body,
			Timestamp:   time.Unix(protoMsg.Timestamp, 0),
		}

		// Direct call to the repo instead of a channel
		if err := c.repo.SaveMessage(ctx, msg); err != nil {
			log.Printf("Failed to save message to DB: %v", err)
			continue
		}

		log.Printf("Message persisted: %s", msg.Message_id)
	}
}
