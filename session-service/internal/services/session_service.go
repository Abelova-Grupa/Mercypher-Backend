package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/proto/session"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionService struct {
	repo     repository.SessionRepository
	jwtMaker token.JWTMaker
}

type CreateSessionInput struct {
	Username string
	ConnectedAt time.Time
}

type CreateSessionResponse struct {
	Username string
	IsActive bool 
	ConnectedAt time.Time
}

var (
	ErrInvalidParams = errors.New("parameters are invalid")
)

func NewSessionService(repo repository.SessionRepository, jwtMaker *token.JWTMaker) *SessionService {
	return &SessionService{repo: repo, jwtMaker: *jwtMaker}
}

// Should create a session after logging in
func (s *SessionService) CreateSession(ctx context.Context, username string, connectedAt time.Time) (*pb.Session, error) {
	if username == ""  || connectedAt.IsZero(){
		return nil, ErrInvalidParams
	}

	session := &models.Session{
		Username: username,
		IsActive: true,
		ConnectedAt: connectedAt,
	}
	
	createdSession, err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("unable to create session for specified user: %v", err)
	}

	return convertSessionToPb(createdSession), nil
}

func (s *SessionService) GetSessionByUsername(ctx context.Context, username string) (*pb.Session, error) {
	if username == "" {
		return nil, ErrInvalidParams
	}
	session, err := s.repo.GetSessionByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("unable to find session specified by username %v: %v", username, err)
	}
	return convertSessionToPb(session), nil
}

// This method is used to connect only authenticated users
func (s *SessionService) Connect(ctx context.Context, username string) (error) {
	if username == "" {
		return ErrInvalidParams
	}
	
	session, _ := s.repo.GetSessionByUsername(ctx,username)
	var err error
	if session == nil {
		_, err = s.repo.CreateSession(ctx,&models.Session{Username: username, IsActive: true, ConnectedAt: time.Now()})
	}else{
		session.IsActive = true
		session.ConnectedAt = time.Now()
		_, err = s.repo.UpdateSession(ctx,session)
	}
	if err != nil {
		return fmt.Errorf("Failed to connect user with username %v: %v", username, err)
	}

	return nil
}

func (s *SessionService) Disconnect(ctx context.Context, username string) (error) {
	if username == "" {
		return ErrInvalidParams
	}

	session, err := s.repo.GetSessionByUsername(ctx,username)
	if session == nil || err != nil {
		return fmt.Errorf("Session for user with specified username %v doesn't exist: %v",username, err)
	}

	session.IsActive = false
	session.LastSeenTime = time.Now()
	_, err = s.repo.UpdateSession(ctx,session)
	if err != nil {
		return fmt.Errorf("User %v didn't properly disconnect: %v",username, err)
	}
	return nil
}

// MAPPERS
func convertSessionToPb(session *models.Session) *pb.Session {
	return &pb.Session{
		ID:          session.ID,
		Username:      session.Username,
		IsActive:    session.IsActive,
		ConnectedAt: timestamppb.New(session.ConnectedAt),
	}
}

