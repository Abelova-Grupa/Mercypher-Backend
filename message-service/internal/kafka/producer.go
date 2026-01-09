package kafka

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const TopicName = "chat_messages_v1"

// PublishMessage writes to kafka
func PublishMessage(ctx context.Context, brokers []string, msg *pb.ChatMessage) (string, error) {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: TopicName,
	}
	defer writer.Close()

	// converting struct to bytes
	data, err := proto.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto: %w", err)
	}

	// writting the actual data
	err = writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

	if err != nil {
		return "", fmt.Errorf("failed to write to kafka: %w", err)
	}

	return msg.Id, nil
}

// Reading kafka topic within the range (currently all chats :D)
func FetchMessages(ctx context.Context, brokers []string, from, to int64) ([]*pb.ChatMessage, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   TopicName,
		// Partition: 0, // In single topic dev, we usually have 1 partition
		// MinBytes:  1,
		// MaxBytes:  1024,
	})
	defer reader.Close()

	// moving the offset to the begging of given range
	if err := reader.SetOffset(from); err != nil {
		return nil, err
	}

	var messages []*pb.ChatMessage
	count := to - from
	if count <= 0 {
		return nil, nil
	}

	for i := int64(0); i < count; i++ {
		readCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		m, err := reader.ReadMessage(readCtx)
		cancel()

		if err != nil {
			break // break if we reached the end before range
		}

		var chatMsg pb.ChatMessage
		if err := proto.Unmarshal(m.Value, &chatMsg); err == nil {
			messages = append(messages, &chatMsg)
		}
	}

	return messages, nil
}
