package client

import (
	pb "github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/pb"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	pb.SessionServiceClient
}

func NewGrpcClient(address string) (*GrpcClient, error) {
	// TODO: Use credentials
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewSessionServiceClient(conn)
	return &GrpcClient{client}, nil
}
