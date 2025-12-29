package client

import (
	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	userpb.UserServiceClient
}

func NewGrpcClient(address string) (*GrpcClient, error) {
	// TODO: Use credentials
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := userpb.NewUserServiceClient(conn)
	return &GrpcClient{client}, nil
}
