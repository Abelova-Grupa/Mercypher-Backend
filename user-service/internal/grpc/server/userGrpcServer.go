package server

import (
	"context"

	sessionClient "github.com/Abelova-Grupa/Mercypher/session-service/external/client"

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

func (g *GrpcServer) Register(ctx context.Context, user *pb.User) (*pb.User, error) {
	return nil, nil
}

func (g *GrpcServer) Login(ctx context.Context, userCredentials *pb.LoginRequest) (*pb.LoginResponse, error) {
	//Check if the session exists with token and userID

	return nil, nil
}
