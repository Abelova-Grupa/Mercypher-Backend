package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SessionService struct {
	repo repository.SessionRepository
}

var (
	ErrInvalidParams = errors.New("parameters are invalid")
)

// TODO implement parameter conversion in service layer

func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

// Think about which services should session have
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

func (s *SessionService) RefreshToken(ctx context.Context, refreshToken string, tokenType token.TokenType) (string, error) {
	jwtMaker := token.JWTMaker{}
	newToken, err := jwtMaker.RefreshToken(ctx, refreshToken, tokenType)
	if newToken == "" || err != nil {
		return "", err
	}
	return newToken, nil
}

func (s *SessionService) CreateLastSeen(ctx context.Context, lastSeenPb *pb.LastSeen) (*pb.LastSeen, error) {
	var lastSeen *models.LastSeenSession
	var err error

	if lastSeenPb.LastSeen.AsTime().IsZero() {
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
	if lastSeenPb.UserID == "" || lastSeenPb.LastSeen.AsTime().IsZero() {
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

func convertLastSeenToPB(lastSeen *models.LastSeenSession) *pb.LastSeen {
	return &pb.LastSeen{
		UserID:   lastSeen.UserID,
		LastSeen: timestamppb.New(time.Unix(lastSeen.LastSeen, 0)),
	}
}

func convertPBToLastSeen(lastSeenPB *pb.LastSeen) *models.LastSeenSession {
	return &models.LastSeenSession{
		UserID:   lastSeenPB.UserID,
		LastSeen: lastSeenPB.LastSeen.AsTime().Unix(),
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
