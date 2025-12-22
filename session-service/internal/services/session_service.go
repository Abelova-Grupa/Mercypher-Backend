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

var (
	ErrInvalidParams = errors.New("parameters are invalid")
)

func NewSessionService(repo repository.SessionRepository, jwtMaker *token.JWTMaker) *SessionService {
	return &SessionService{repo: repo, jwtMaker: *jwtMaker}
}

func (s *SessionService) CreateToken(ctx context.Context, userID string, duration time.Duration, tokenType token.TokenType) (string, *token.Payload, error) {
	jwtMaker := token.JWTMaker{}
	token, payload, err := jwtMaker.CreateToken(userID, duration, tokenType)
	if token == "" || payload == nil || err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (s *SessionService) VerifyToken(ctx context.Context, testToken string, tokenType token.TokenType) (*token.Payload, error) {
	jwtMaker := token.JWTMaker{}
	payload, err := jwtMaker.VerifyToken(testToken, tokenType)
	if payload == nil || err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *SessionService) CreateLastSeen(ctx context.Context, lastSeenPb *pb.LastSeen) (*pb.LastSeen, error) {
	var lastSeen *models.LastSeenSession
	var err error

	if lastSeenPb.Time.AsTime().IsZero() {
		return nil, ErrInvalidParams
	}
	// Convert to DTO
	lastSeen = convertPBToLastSeen(lastSeenPb)
	lastSeen, err = s.repo.CreateLastSeen(ctx, lastSeen)

	if err != nil {
		return nil, fmt.Errorf("unable to create last seen: %v", err)
	}
	lastSeenPb = convertLastSeenToPB(lastSeen)

	return lastSeenPb, nil
}

func (s *SessionService) GetLastSeenByUserID(ctx context.Context, userID string) (*pb.LastSeen, error) {
	if userID == "" {
		return nil, ErrInvalidParams
	}
	lastSeen, err := s.repo.GetLastSeenByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get last seen by specified id: %v", err)
	}
	lastSeenPb := convertLastSeenToPB(lastSeen)
	return lastSeenPb, nil
}

func (s *SessionService) UpdateLastSeen(ctx context.Context, lastSeenPb *pb.LastSeen) (*pb.LastSeen, error) {
	if lastSeenPb.UserID == "" || lastSeenPb.Time.AsTime().IsZero() {
		return nil, ErrInvalidParams
	}
	var lastSeen *models.LastSeenSession
	var err error

	lastSeen = convertPBToLastSeen(lastSeenPb)
	lastSeen, err = s.repo.UpdateLastSeen(ctx, lastSeen)
	if err != nil {
		return nil, fmt.Errorf("unable to update last seen by specified id: %v", err)
	}

	lastSeenPb = convertLastSeenToPB(lastSeen)
	return lastSeenPb, nil
}

func (s *SessionService) DeleteLastSeen(ctx context.Context, userID *pb.UserID) error {
	if userID.UserID == "" {
		return ErrInvalidParams
	}
	err := s.repo.DeleteLastSeen(ctx, userID.UserID)
	if err != nil {
		return fmt.Errorf("unable to delete last seen object with specified id: %v", err)
	}
	return nil
}

func (s *SessionService) CreateUserLocation(ctx context.Context, userLocationPb *pb.UserLocation) (*pb.UserLocation, error) {
	if userLocationPb.APIAdress == "" {
		return nil, ErrInvalidParams
	}
	var userLocation *models.UserLocation
	var err error
	userLocation = convertPBToLocation(userLocationPb)
	userLocation, err = s.repo.CreateUserLocation(ctx, userLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to create user location object: %v", err)
	}
	userLocationPb = convertLocationToPb(userLocation)
	return userLocationPb, nil
}

func (s *SessionService) GetUserLocationByUserID(ctx context.Context, userID string) (*pb.UserLocation, error) {
	if userID == "" {
		return nil, ErrInvalidParams
	}
	userLocation, err := s.repo.GetUserLocationByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user location by specified id: %v", err)
	}

	userLocationPb := convertLocationToPb(userLocation)
	return userLocationPb, nil
}

func (s *SessionService) UpdateUserLocation(ctx context.Context, userLocationPb *pb.UserLocation) (*pb.UserLocation, error) {
	if userLocationPb.UserID == "" || userLocationPb.APIAdress == "" {
		return nil, ErrInvalidParams
	}

	var userLocation *models.UserLocation
	var err error
	userLocation = convertPBToLocation(userLocationPb)
	userLocation, err = s.repo.UpdateUserLocation(ctx, userLocation)

	if err != nil {
		return nil, fmt.Errorf("unable to update user location by specified id: %v", err)
	}

	userLocationPb = convertLocationToPb(userLocation)
	return userLocationPb, nil
}

func (s *SessionService) DeleteUserLocation(ctx context.Context, userID *pb.UserID) error {
	if userID.UserID == "" {
		return ErrInvalidParams
	}
	err := s.repo.DeleteUserLocation(ctx, userID.UserID)
	if err != nil {
		return fmt.Errorf("unable to delete location object with specified id: %v", err)
	}
	return nil
}

// Should create a session after logging in
// TODO: Change parameter *pb.Session to just userID
func (s *SessionService) CreateSession(ctx context.Context, sessionPb *pb.Session) (*pb.Session, error) {
	if sessionPb.UserID == "" {
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

func (s *SessionService) GetSessionByUserID(ctx context.Context, userID *pb.UserID) (*pb.Session, error) {
	if userID.UserID == "" {
		return nil, ErrInvalidParams
	}
	session, err := s.repo.GetSessionByUserID(ctx, userID.UserID)
	if err != nil {
		return nil, fmt.Errorf("unable to find session specified by id %v: %v", userID.UserID, err)
	}
	return convertSessionToPb(session), nil
}

func convertLastSeenToPB(lastSeen *models.LastSeenSession) *pb.LastSeen {
	return &pb.LastSeen{
		UserID:   lastSeen.UserID,
		Time: timestamppb.New(lastSeen.Time),
	}
}

func convertPBToLastSeen(lastSeenPB *pb.LastSeen) *models.LastSeenSession {
	return &models.LastSeenSession{
		UserID:   lastSeenPB.UserID,
		Time: lastSeenPB.Time.AsTime(),
	}
}

func convertPBToLocation(userLocationPB *pb.UserLocation) *models.UserLocation {
	return &models.UserLocation{
		UserID: userLocationPB.UserID,
		ApiIP:  userLocationPB.APIAdress,
	}
}

func convertLocationToPb(userLocation *models.UserLocation) *pb.UserLocation {
	return &pb.UserLocation{
		UserID:    userLocation.UserID,
		APIAdress: userLocation.ApiIP,
	}
}

func convertSessionToPb(session *models.Session) *pb.Session {
	return &pb.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		IsActive: session.IsActive,
		ConnectedAt:  timestamppb.New(session.ConnectedAt),
	}
}

func convertPbToSession(sessionPb *pb.Session) *models.Session {
	return &models.Session{
		ID:           sessionPb.ID,
		UserID:       sessionPb.UserID,
		IsActive:     sessionPb.IsActive,
		ConnectedAt:  sessionPb.ConnectedAt.AsTime(),
	}
}
