package servers

import (
	"context"
	"log"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
)

type KafkaConsumer struct {
	gwChan		chan *domain.ChatMessage
	reader		*kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, outChan chan *domain.ChatMessage) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
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
			MessageId: protoMsg.Id,
			SenderId: protoMsg.SenderId,
			Receiver_id: protoMsg.RecieverId,
			Body: protoMsg.Body,
		}
		
		log.Printf("Message received from Kafka: %s", msg.Body)

		c.gwChan <- &msg
		
	}
}