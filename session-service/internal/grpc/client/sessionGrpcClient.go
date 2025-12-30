package client

import (
	sessionpb "github.com/Abelova-Grupa/Mercypher/proto/session"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	sessionpb.SessionServiceClient
}

func NewGrpcClient(address string) (*GrpcClient, error) {
	// TODO: Use credentials
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := sessionpb.NewSessionServiceClient(conn)
	return &GrpcClient{client}, nil
}
