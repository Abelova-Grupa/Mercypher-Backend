package server

import (
	"context"
	"log"

	"time"

	messagepb "github.com/Abelova-Grupa/Mercypher/message-service/internal/grpc"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	relaypb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/proto"
)

type MessageServer struct {
	messagepb.UnimplementedMessageServiceServer
	msg_service *service.MessageService
}

func NewMessageServer(msgsvc *service.MessageService) *MessageServer {
	return &MessageServer{
		msg_service: msgsvc,
	}
}

// SendMessage implements the RPC method
func (s *MessageServer) SendMessage(ctx context.Context, msg *messagepb.ChatMessage) (*messagepb.MessageAck, error) {

	// Construct the message as defined in model (for persisting in database)
	modelMessage := model.ChatMessage{
		Message_id:  "",
		Sender_id:   msg.SenderId,
		Receiver_id: msg.RecipientId, // Gospode Boze, jednog sam nazvao receiver, a drugog recipient...
		Body:        msg.Body,
		Timestamp:   time.Unix(msg.Timestamp, 0),
	}

	s.msg_service.ProcessMessage(ctx, &modelMessage)

	ack := messagepb.MessageAck{
		MessageId: modelMessage.Message_id,
	}

	return &ack, nil
}

// Should be moved somewhere else and run somewhere else
func ConnectRelayService() {
	conn, err := grpc.Dial("localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := relaypb.NewRelayServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}

func (s *MessageServer) RelayMessage(ctx context.Context, msg *messagepb.ChatMessage) (*messagepb.RelayResponse, error) {

	return nil, nil
}

// func StartGrpcServer() {
// 	listener, err := net.Listen("tcp", ":50051")
//     if err != nil {
//         log.Fatalf("Failed to listen: %v", err)
//     }

//     grpcServer := grpc.NewServer()
//     messagepb.RegisterMessageServiceServer(grpcServer, &messageServer{})
// 	if err := grpcServer.Serve(listener); err != nil {
//         log.Fatalf("Failed to serve: %v", err)
//     }
// }
