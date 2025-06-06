package server

import (
	"context"
	"errors"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type grpcServer struct {
	sessionDB      *gorm.DB
	sessionRepo    repository.SessionRepository
	sessionService services.SessionService
	pb.UnsafeSessionServiceServer
}

// Change to pointer needed structs
func NewGrpcServer(db *gorm.DB) *grpcServer {
	repo := repository.NewSessionRepository(db)
	jwtMaker, _ := token.NewJWTMaker(uuid.NewString(), uuid.NewString())
	service := services.NewSessionService(repo, jwtMaker)
	return &grpcServer{
		sessionDB:      db,
		sessionRepo:    repo,
		sessionService: *service,
	}
}

func (s *grpcServer) CreateUserLocation(ctx context.Context, userLocation *pb.UserLocation) (*pb.UserLocation, error) {
	userLocation, err := s.sessionService.CreateUserLocation(ctx, userLocation)
	if err != nil {
		return nil, err
	}
	return userLocation, err
}

func (s *grpcServer) GetUserLocation(ctx context.Context, userID *pb.UserID) (*pb.UserLocation, error) {
	userLocation, err := s.sessionService.GetUserLocationByUserID(ctx, userID.UserID)
	if err != nil {
		return nil, err
	}
	return userLocation, nil
}

func (s *grpcServer) UpdateUserLocation(ctx context.Context, userLoc *pb.UserLocation) (*pb.UserLocation, error) {
	// If the userID doesnt exist it will create a new UserLocation, otherwise it will update existing UserLocation
	userLocation, err := s.sessionService.UpdateUserLocation(ctx, userLoc)
	if err != nil {
		return nil, errors.New("unable to update user location")
	}
	return userLocation, nil
}

func (s *grpcServer) DeleteUserLocation(ctx context.Context, userID *pb.UserID) (*emptypb.Empty, error) {
	err := s.sessionService.DeleteUserLocation(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *grpcServer) CreateLastSeen(ctx context.Context, lastSeen *pb.LastSeen) (*pb.LastSeen, error) {
	lastSeen, err := s.sessionService.CreateLastSeen(ctx, lastSeen)
	if err != nil {
		return nil, err
	}
	return lastSeen, nil
}

func (s *grpcServer) GetLastSeen(ctx context.Context, userID *pb.UserID) (*pb.LastSeen, error) {
	lastSeen, err := s.sessionService.GetLastSeenByUserID(ctx, userID.UserID)
	if err != nil {
		return nil, errors.New("unable to retreive last seen info")
	}
	return lastSeen, nil
}

func (s *grpcServer) UpdateLastSeen(ctx context.Context, lastSeen *pb.LastSeen) (*pb.LastSeen, error) {
	lastSeen, err := s.sessionService.UpdateLastSeen(ctx, lastSeen)
	if err != nil {
		return nil, errors.New("unable to update last seen info")
	}

	return lastSeen, nil
}

func (s *grpcServer) DeleteLastSeen(ctx context.Context, userID *pb.UserID) (*emptypb.Empty, error) {
	err := s.sessionService.DeleteLastSeen(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *grpcServer) CreateToken(ctx context.Context, userID *pb.UserID) (*pb.Token, error) {
	token, _, err := s.sessionService.CreateToken(ctx, userID.UserID, time.Minute*15, 1)
	if err != nil {
		return nil, err
	}
	return &pb.Token{
		Token:     token,
		TokenType: "access-token",
	}, nil
}

func (s *grpcServer) VerifyToken(ctx context.Context, token *pb.Token) (*pb.VerifiedToken, error) {
	//Convert pb.Token.TokenType into entity token.TokenType
	_, err := s.sessionService.VerifyToken(ctx, token.Token, 1)
	if err != nil {
		return &pb.VerifiedToken{IsValid: false}, err
	}
	return &pb.VerifiedToken{IsValid: true}, err
}

func (s *grpcServer) RefreshToken(ctx context.Context, token *pb.Token) (*pb.Token, error) {
	newToken, err := s.sessionService.RefreshToken(ctx, token.Token, 1)
	if err != nil {
		return nil, err
	}
	return &pb.Token{Token: newToken, TokenType: "access_token"}, nil
}

func (s *grpcServer) CreateSession(ctx context.Context, sessionPb *pb.Session) (*pb.Session, error) {
	newSession, err := s.sessionService.CreateSession(ctx, sessionPb)
	if err != nil {
		return nil, err
	}
	return newSession, nil
}

func (s *grpcServer) GetSessionByUserID(ctx context.Context, userID *pb.UserID) (*pb.Session, error) {
	session, err := s.sessionService.GetSessionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return session, nil
}
