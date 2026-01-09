package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	sessionClient "github.com/Abelova-Grupa/Mercypher/session-service/external/client"

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
func (g *GrpcServer) Register(ctx context.Context, user *userpb.User) (*userpb.User, error) {
	user, err := g.userService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (g *GrpcServer) Login(ctx context.Context, loginRequest *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	username := loginRequest.GetUsername()
	passedToken := loginRequest.GetAccessToken()
	password := loginRequest.GetPassword()

	if passedToken != "" {
		verified, _ := g.sessionClient.VerifyToken(ctx, &sessionpb.Token{Token: passedToken})
		if verified.Value {
			return &userpb.LoginResponse{Username: username, AccessToken: passedToken}, nil
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
	token, err := g.sessionClient.Connect(ctx, &sessionpb.Username{Name: username})
	if err != nil {
		return nil, fmt.Errorf("Failed session sign in for user %v : %v ", username, err)
	}

	log.Print("Succesfull login")
	return &userpb.LoginResponse{Username: username, AccessToken: token.Token}, nil

}

func (g *GrpcServer) Logout(ctx context.Context, logoutRequest *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	if logoutRequest.Username == "" {
		return nil, errors.New("Invalid params for logout operation")
	}

	usernamePb := &sessionpb.Username{Name: logoutRequest.Username}
	success, err := g.sessionClient.Disconnect(ctx,usernamePb)
	return &userpb.LogoutResponse{OperationSuccessfull: success.GetValue()}, err

}

func (g *GrpcServer) ValidateAccount(ctx context.Context, validateRequest *userpb.ValidateAccountRequest) (*userpb.ValidateAccountResponse, error) {
	if err := g.userService.ValidateAccount(ctx,validateRequest); err != nil {
		return &userpb.ValidateAccountResponse{Success: false}, err
	}
	return &userpb.ValidateAccountResponse{Success: true}, nil
}
