package servers

import (
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/Abelova-Grupa/Mercypher/api/external/grpc"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedGatewayServiceServer
	wg          *sync.WaitGroup
	grpcServer  *grpc.Server
}

// Constructor
func NewGrpcServer(wg *sync.WaitGroup) *GrpcServer {
	return &GrpcServer{
		wg:         wg,
		grpcServer: grpc.NewServer(),
	}
}

// Start method
func (s *GrpcServer) Start(addr string) {
	s.wg.Add(1)
	defer s.wg.Done()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}

	pb.RegisterGatewayServiceServer(s.grpcServer, s)

	log.Println("gRPC server thread running on: ", addr)

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
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
