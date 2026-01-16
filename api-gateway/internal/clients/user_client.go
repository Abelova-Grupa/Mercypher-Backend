package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserClient cretes a new client to a user service on the given address.
//
// Note:	The situation is the same as in NewMessageClient code. Even if the
//
//	connection fails or refuses it wont be registered. Only when sending
//	messages to an unexisting address will the error be thrown.
func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("Connection refused: nil")
	}

	client := userpb.NewUserServiceClient(conn)

	return &UserClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

// Register method returns ID of the created user.
func (c *UserClient) Register(user domain.User, password string) (string, error) {
	response, err := c.client.RegisterUser(context.Background(),
		&userpb.RegisterUserRequest{
			Username:  user.Username,
			Email:     user.Email,
			Password:  password,
			CreatedAt: timestamppb.Now(),
		})
	fmt.Print(response)

	if err != nil {
		return "", err
	}

	return response.Username, nil
}

// Login method returns access token of the logged user
func (c *UserClient) Login(user domain.User, password string, accessToken string) (string, error) {
	response, err := c.client.LoginUser(context.Background(),
		&userpb.LoginUserRequest{
			Username:    user.Username, // Redundant?
			Password:    password,
			Token: accessToken,
		})

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	return response.AccessToken, nil
}

func (c *UserClient) VerifyToken(token string) (bool, error) {
	resp, err := c.client.VerifyToken(context.Background(), &userpb.VerifyTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, err
	} else {
		return resp.Value, nil
	}
}
