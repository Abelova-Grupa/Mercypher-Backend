package server

import (
	"context"

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
	return &GrpcServer{
		userDB:      db,
		userRepo:    repo,
		userService: *service,
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

func (g *GrpcServer) Login(ctx context.Context, loginRequest *pb.LoginRequest) (*pb.LoginResponse, error) {
	//Check if the session exists with token and userID

	userID := &sessionpb.UserID{
		UserID: loginRequest.GetUserID(),
	}
	sessionPb, _ := g.sessionClient.GetSessionByUserID(ctx, userID)
	// Retrieves access token, already in session
	if sessionPb != nil {
		return &pb.LoginResponse{
			UserID:      sessionPb.GetUserID(),
			Username:    loginRequest.Username,
			AccessToken: sessionPb.GetAccessToken(),
		}, nil
	} else {
		// Check username and password
		isLoggedIn, err := g.userService.Login(ctx, loginRequest.GetUsername(), loginRequest.GetPassword())
		if err != nil {
			return nil, err
		}
		if isLoggedIn {
			createdSessionPb, err := g.sessionClient.CreateSession(ctx, &sessionpb.Session{UserID: loginRequest.GetUserID()})
			if err != nil {
				return nil, err
			}
			return &pb.LoginResponse{
				UserID:      loginRequest.GetUserID(),
				Username:    loginRequest.GetUsername(),
				AccessToken: createdSessionPb.AccessToken,
			}, nil
		} else {
			return nil, err
		}
	}

}
