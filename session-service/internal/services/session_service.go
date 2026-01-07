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
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionService struct {
	repo     repository.SessionRepository
	jwtMaker token.JWTMaker
}

var (
	// Number of minutes in a day
	sessionDuration = 1440 * time.Minute
)

var (
	ErrInvalidParams = errors.New("parameters are invalid")
)

func NewSessionService(repo repository.SessionRepository, jwtMaker *token.JWTMaker) *SessionService {
	return &SessionService{repo: repo, jwtMaker: *jwtMaker}
}

// TOKEN AND SESSION

func (s *SessionService) CreateToken(ctx context.Context, username string, duration time.Duration) (string,error) {
	jwtMaker := token.JWTMaker{}
	token, _, err := jwtMaker.CreateToken(username, duration)
	if token == "" || err != nil {
		return "", err
	}

	return token, nil
}

func (s *SessionService) VerifyToken(ctx context.Context, tokenPb *pb.Token) (bool, error) {
	jwtMaker := token.JWTMaker{}
	payload, err := jwtMaker.VerifyToken(tokenPb.Token)
	if payload == nil || err != nil {
		return false, err
	}
	return true, nil
}



// Should create a session after logging in
// TODO: Change parameter *pb.Session to just userID
func (s *SessionService) CreateSession(ctx context.Context, sessionPb *pb.Session) (*pb.Session, error) {
	if sessionPb.Username == "" {
		return nil, ErrInvalidParams
	}

	session := convertPbToSession(sessionPb)
	session.IsActive = true
	session.ConnectedAt = sessionPb.ConnectedAt.AsTime()

	createdSession, err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("unable to create session for specified user: %v", err)
	}

	return convertSessionToPb(createdSession), nil
}

func (s *SessionService) GetSessionByUsername(ctx context.Context, usernamePb *pb.Username) (*pb.Session, error) {
	if usernamePb.Name == "" {
		return nil, ErrInvalidParams
	}
	session, err := s.repo.GetSessionByUsername(ctx, usernamePb.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to find session specified by username %v: %v", usernamePb.Name, err)
	}
	return convertSessionToPb(session), nil
}

// This method is used to connect only authenticated users
func (s *SessionService) Connect(ctx context.Context, username string) (*pb.Token,bool,error) {
	if username == "" {
		return nil, false,ErrInvalidParams
	}
	
	var token string

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
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
			err = fmt.Errorf("Failed to connect user with username %v: %v", username, err)
		}
		return err
	})

	group.Go(func() error {
		var err error
		token, err = s.CreateToken(ctx,username, time.Duration(sessionDuration))
		if err != nil {
			err = fmt.Errorf("Failed to create a token for user %v : %v", username, err)
		}
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, false, err
	}

	return &pb.Token{Token: token}, true, nil
}

func (s *SessionService) Disconnect(ctx context.Context, usernamePb *pb.Username) (bool, error) {
	if usernamePb == nil || usernamePb.Name == "" {
		return false, ErrInvalidParams
	}

	session, err := s.repo.GetSessionByUsername(ctx,usernamePb.Name)
	if session == nil || err != nil {
		return false, fmt.Errorf("Session for user with specified username %v doesn't exist: %v",usernamePb.Name, err)
	}

	session.IsActive = false
	session.LastSeenTime = time.Now()
	_, err = s.repo.UpdateSession(ctx,session)
	if err != nil {
		return false, fmt.Errorf("User %v didn't properly disconnect: %v",usernamePb.Name, err)
	}
	return true, nil
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

func convertPbToSession(sessionPb *pb.Session) *models.Session {
	return &models.Session{
		ID:          sessionPb.ID,
		Username:      sessionPb.Username,
		IsActive:    sessionPb.IsActive,
		ConnectedAt: sessionPb.ConnectedAt.AsTime(),
	}
}

