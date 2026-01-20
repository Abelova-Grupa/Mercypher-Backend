package kafka

import (
	"context"
	"log"
	"time"

	gp "github.com/Abelova-Grupa/Mercypher/proto/gateway"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func StartLiveForwarder(ctx context.Context, kafkaBrokers []string, gatewayAddr string) {
	// 1. Setup gRPC Connection to Gateway
	conn, err := grpc.NewClient(gatewayAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Gateway: %v", err)
	}
	defer conn.Close()

	client := gp.NewGatewayServiceClient(conn)

	// Create a stream used for communication
	stream, err := client.Stream(ctx)
	if err != nil {
		log.Fatalf("Failed to open Gateway stream: %v", err)
	}

	// 2. Setup Kafka Reader with a specific GroupID
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		Topic:    TopicName,
		GroupID:  "gateway-live-forwarder", // Separate from DB worker
		MinBytes: 10e3,                     // ~10KB
		MaxBytes: 10e6,                     // ~10MB
	})
	defer r.Close()

	log.Printf("Live Forwarder started. Forwarding to %s", gatewayAddr)

	for {
		// ReadMessage automatically commits offsets after reading
		// (use FetchMessage/CommitMessages for more control later)
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			break
		}

		var chatMsg pb.ChatMessage
		if err := proto.Unmarshal(m.Value, &chatMsg); err != nil {
			log.Printf("Failed to unmarshal kafka message: %v", err)
			continue
		}

		// 3. Map Message Service proto to Gateway proto
		req := &gp.GatewayRequest{
			Payload: &gp.GatewayRequest_ChatMessage{
				ChatMessage: &gp.ChatMessage{
					MessageId:   chatMsg.Id,
					SenderId:    chatMsg.SenderId,
					RecipientId: chatMsg.RecieverId,
					Timestamp:   chatMsg.Timestamp,
					Body:        chatMsg.Body,
				},
			},
		}

		// 4. Push to the long-lived stream
		if err := stream.Send(req); err != nil {
			log.Printf("Failed to send to Gateway stream: %v. Attempting to reconnect...", err)

			// Simple reconnection logic
			time.Sleep(2 * time.Second)
			newStream, err := client.Stream(ctx)
			if err == nil {
				stream = newStream
				_ = stream.Send(req) // GIve up after 1st try
			}
		}
	}
}
