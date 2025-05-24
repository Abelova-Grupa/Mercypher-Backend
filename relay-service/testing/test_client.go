package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Dial the gRPC server
	conn, err := grpc.Dial("localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRelayServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Test SendMessage
	msg := &pb.Message{
		SenderId:   "77",
		ReceiverId: "55",
		Data:       "Hello from client!",
		Timestamp:  500,
	}
	status, err := client.SendMessage(ctx, msg)
	if err != nil {
		log.Fatalf("SendMessage error: %v", err)
	}
	log.Printf("SendMessage status: %v", status)

	// 3. Test GetMessages (streaming) for "user2"
	stream, err := client.GetMessages(ctx, &pb.UserId{Id: "55"})
	if err != nil {
		log.Fatalf("GetMessages error: %v", err)
	}

	log.Println("Streaming messages for user2:")
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream.Recv() error: %v", err)
		}
		log.Printf("- [%v] from %s: %s",
			m.Timestamp,
			m.GetSenderId(),
			m.GetData(),
		)
	}
}
