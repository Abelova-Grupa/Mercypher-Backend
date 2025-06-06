package client

import (
	pb "github.com/Abelova-Grupa/Mercypher/user-service/external/proto"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	pb.UserServiceClient
}

func NewGrpcClient(address string) (*GrpcClient, error) {
	// TODO: Use credentials
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &GrpcClient{client}, nil
}
