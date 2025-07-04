package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	sessionClient "github.com/Abelova-Grupa/Mercypher/session-service/external/client"
	sessionpb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"

	pb "github.com/Abelova-Grupa/Mercypher/user-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/service"
	"gorm.io/gorm"
)

type GrpcServer struct {
	userDB      *gorm.DB
	userRepo    repository.UserRepository
	userService service.UserService
	pb.UnsafeUserServiceServer
	sessionClient sessionClient.GrpcClient
}

func NewGrpcServer(db *gorm.DB) *GrpcServer {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	grpcClient, _ := sessionClient.NewGrpcClient("localhost:50055")
	return &GrpcServer{
		userDB:        db,
		userRepo:      repo,
		userService:   *service,
		sessionClient: *grpcClient,
	}
}

// Should only create a user not a session
func (g *GrpcServer) Register(ctx context.Context, user *pb.User) (*pb.User, error) {
	user, err := g.userService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Note to future maintainers: Assume that user (and Gateway) doesn't know its ID for it is provided
//
//	at the time of successful registration and/or login. Therefore, asking for ID on login forwards
//	the nil value to other services creating hard to find errors.
//
//	Also, username is an unique key!
func (g *GrpcServer) Login(ctx context.Context, loginRequest *pb.LoginRequest) (*pb.LoginResponse, error) {

	//Check if the session exists with token and userID
	userID := &sessionpb.UserID{
		UserID: loginRequest.GetUserID(),
	}

	sessionPb, _ := g.sessionClient.GetSessionByUserID(ctx, userID)

	// Retrieves access token, already in session
	if sessionPb != nil {

		log.Println("Refreshing session for user ", loginRequest.GetUsername())
		return &pb.LoginResponse{
			UserID:      sessionPb.GetUserID(),
			Username:    loginRequest.Username,
			AccessToken: sessionPb.GetAccessToken(),
		}, nil
	} else {

		log.Print("Validating session for user ", loginRequest.GetUsername())

		// User logging in first time: Check username and password
		isLoggedIn, err := g.userService.Login(ctx, loginRequest.GetUsername(), loginRequest.GetPassword())
		if err != nil {
			log.Println("...AUTHORIZATION FAILED!")
			return nil, err
		}
		if isLoggedIn {

			log.Println("...OK")
			// Get the id from the database for user doesn't need to supply it in the request!
			user, err := g.userRepo.GetUserByUsername(context.Background(), loginRequest.Username)
			if err != nil {
				fmt.Print(err)
				return nil, err
			}

			createdSessionPb, err := g.sessionClient.CreateSession(ctx, &sessionpb.Session{UserID: user.ID})
			if err != nil {
				fmt.Print(err)
				return nil, err
			}
			return &pb.LoginResponse{
				UserID:      user.ID,
				Username:    loginRequest.GetUsername(),
				AccessToken: createdSessionPb.AccessToken,
			}, nil
		} else {
			log.Println("...Invalid credentials.")
			return nil, errors.New("invalid credentials")
		}
	}

}
