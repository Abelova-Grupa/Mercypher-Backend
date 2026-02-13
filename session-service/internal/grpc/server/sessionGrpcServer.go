package server

import (
	"context"

	sessionpb "github.com/Abelova-Grupa/Mercypher/proto/session"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/store"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	sessionCache   *redis.Client
	sessionRepo    repository.SessionRepository
	sessionService services.SessionService
	sessionpb.UnsafeSessionServiceServer
}

// Change to pointer needed structs
func NewGrpcServer() *grpcServer {
	ctx := context.Background()
	rdb := store.NewSessionCache(ctx)
	repo := repository.NewSessionRepository(rdb)
	jwtMaker, _ := token.NewJWTMaker(uuid.NewString())
	service := services.NewSessionService(repo, jwtMaker)
	return &grpcServer{
		sessionCache:   rdb,
		sessionRepo:    repo,
		sessionService: *service,
	}
}

func (s *grpcServer) Connect(ctx context.Context, connectRequest *sessionpb.ConnectRequest) (*emptypb.Empty, error) {
	if connectRequest == nil || connectRequest.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments for session connection")
	}
	if err := s.sessionService.Connect(ctx, connectRequest.Username); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *grpcServer) Disconnect(ctx context.Context, disconnectRequest *sessionpb.DisconnectRequest) (*emptypb.Empty, error) {
	if disconnectRequest == nil || disconnectRequest.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments for session disconnection")
	}
	if err := s.sessionService.Disconnect(ctx, disconnectRequest.Username); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
