package server

import (
	"context"

	sessionpb "github.com/Abelova-Grupa/Mercypher/proto/session"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type grpcServer struct {
	sessionDB      *gorm.DB
	sessionRepo    repository.SessionRepository
	sessionService services.SessionService
	sessionpb.UnsafeSessionServiceServer
}

// Change to pointer needed structs
func NewGrpcServer(db *gorm.DB) *grpcServer {
	repo := repository.NewSessionRepository(db)
	jwtMaker, _ := token.NewJWTMaker(uuid.NewString())
	service := services.NewSessionService(repo, jwtMaker)
	return &grpcServer{
		sessionDB:      db,
		sessionRepo:    repo,
		sessionService: *service,
	}
}


func (s *grpcServer) Connect(ctx context.Context, username *sessionpb.Username) (*sessionpb.Token, error) {
	token,connected, err := s.sessionService.Connect(ctx,username.Name)
	if !connected {
		return nil, err
	}
	return token, nil
}

func (s *grpcServer) Disconnect(ctx context.Context, credentials *sessionpb.ConnectionCredentials) (*wrapperspb.BoolValue, error) {
	panic("Unimplemented")
}

func (s *grpcServer) VerifyToken(ctx context.Context, token *sessionpb.Token) (*wrapperspb.BoolValue, error) {
	verified, err := s.sessionService.VerifyToken(ctx,token)
	return wrapperspb.Bool(verified), err
}
