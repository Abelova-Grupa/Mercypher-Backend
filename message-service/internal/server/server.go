package server

import (
	"context"

	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/service"
	messagepb "github.com/Abelova-Grupa/Mercypher/proto/message"
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
