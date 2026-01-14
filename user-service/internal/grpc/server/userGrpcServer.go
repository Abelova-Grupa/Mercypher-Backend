package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sessionClient "github.com/Abelova-Grupa/Mercypher/session-service/external/client"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	sessionpb "github.com/Abelova-Grupa/Mercypher/proto/session"
	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/service"
	"gorm.io/gorm"
)

type GrpcServer struct {
	userDB      *gorm.DB
	userRepo    repository.UserRepository
	userService service.UserService
	userpb.UnsafeUserServiceServer
	sessionClient sessionClient.GrpcClient
}

func NewGrpcServer(db *gorm.DB) *GrpcServer {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	// For now localhost is hardcocded
	// TODO: Change localhost hardcoding when ready to deploy
	grpcClient, _ := sessionClient.NewGrpcClient(fmt.Sprintf("localhost:%v",os.Getenv("SESSION_SERVICE_PORT")))
	return &GrpcServer{
		userDB:        db,
		userRepo:      repo,
		userService:   *service,
		sessionClient: *grpcClient,
	}
}

// Should only create a user not a session
func (g *GrpcServer) RegisterUser(ctx context.Context, registerRequestPb *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	res, err := g.userService.Register(ctx, registerRequestPb)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *GrpcServer) LoginUser(ctx context.Context, loginRequest *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {

	username := loginRequest.GetUsername()
	passedToken := loginRequest.GetToken()
	password := loginRequest.GetPassword()

	if passedToken != "" {
		verified, _ := g.userService.VerifyToken(ctx, &userpb.VerifyTokenRequest {Token: passedToken,})
		if verified{
			return &userpb.LoginUserResponse{Username: username, AccessToken: passedToken}, nil
		} else {
			log.Print("Token is invalid, continue with credential checking")
		}
	}

	log.Println("Checking user credentials")
	isLoggedIn, _ := g.userService.Login(ctx, username, password)
	if !isLoggedIn {
		return nil, errors.New("Authentification failed")
	}
	log.Println("Successful authentication creating session...")
	var token string
	var err error
	if token, err = g.userService.CreateToken(ctx, username, 24 * time.Hour); err != nil {
		return nil, fmt.Errorf("Failed to create auth token for user %v : %v", username, err)
	}
	_, err = g.sessionClient.Connect(ctx, &sessionpb.ConnectRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("Failed session sign in for user %v : %v ", username, err)
	}

	log.Print("Succesfull login")
	return &userpb.LoginUserResponse{Username: username, AccessToken: token}, nil

}

func (g *GrpcServer) LogoutUser(ctx context.Context, logoutRequest *userpb.LogoutUserRequest) (*emptypb.Empty, error) {
	if logoutRequest.Username == "" {
		return nil, errors.New("Invalid params for logout operation")
	}
	usernamePb := &sessionpb.DisconnectRequest{Username: logoutRequest.Username}
	if _, err := g.sessionClient.Disconnect(ctx,usernamePb); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) ValidateUserAccount(ctx context.Context, validateRequest *userpb.ValidateUserAccountRequest) (*emptypb.Empty, error) {
	if err := g.userService.ValidateAccount(ctx,validateRequest); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GrpcServer) VerifyToken(ctx context.Context, verifyTokenRequest *userpb.VerifyTokenRequest) (*wrapperspb.BoolValue, error){
	if valid, err:= g.userService.VerifyToken(ctx,verifyTokenRequest); !valid || err != nil {
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(true), nil
}
