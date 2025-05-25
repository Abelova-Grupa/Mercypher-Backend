package servers

import (
	"io"
	"log"

	pb "github.com/Abelova-Grupa/Mercypher/api/internal/grpc"
)

type GrpcServer struct {
	pb.UnimplementedGatewayServiceServer
}

func (s *GrpcServer) Stream(stream pb.GatewayService_StreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Stream recv error: %v", err)
			return err
		}

		switch payload := req.Payload.(type) {

		case *pb.GatewayRequest_ChatMessage:
			msg := payload.ChatMessage
			log.Printf("Chat message: %s -> %s: %s", msg.SenderId, msg.RecipientId, msg.Body)

			// TODO: Forward to the correct routine

			stream.Send(&pb.GatewayResponse{
				Status: "ok",
				Body:   "chat message forwarded",
			})

		case *pb.GatewayRequest_MessageStatus:
			status := payload.MessageStatus
			log.Printf("Status update: %s marked %s as %s", status.RecipientId, status.MessageId, status.Status)

			// TODO: Forward to the correct routine

			stream.Send(&pb.GatewayResponse{
				Status: "ok",
				Body:   "status update forwarded",
			})

		default:
			log.Println("Unknown payload")
			stream.Send(&pb.GatewayResponse{
				Status: "error",
				Body: "unknown payload",
			})
		}
	}
}
