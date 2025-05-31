package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	GetSessionByID(ctx context.Context, ID string) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error)

	CreateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) (*models.LastSeenSession, error)
	GetLastSeenByUserID(ctx context.Context, userID string) (*models.LastSeenSession, error)
	UpdateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) (*models.LastSeenSession, error)
	DeleteLastSeen(ctx context.Context, userID string) error

	CreateUserLocation(tx context.Context, userLocation *models.UserLocation) (*models.UserLocation, error)
	GetUserLocationByUserID(tx context.Context, userID string) (*models.UserLocation, error)
	UpdateUserLocation(tx context.Context, userLocation *models.UserLocation) (*models.UserLocation, error)
	DeleteUserLocation(ctx context.Context, userID string) error
}

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepo {
	return &SessionRepo{DB: db}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	err := s.DB.WithContext(ctx).Create(session).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store a new session in db: %v", err)
	}
	return session, nil
}

// Should return payloadID from refreshToken
func (s *SessionRepo) GetSessionByID(ctx context.Context, ID string) (*models.Session, error) {
	var session models.Session
	result := s.DB.WithContext(ctx).Where("user_id = ?", ID).First(&session)
	return &session, result.Error
}

func (s *SessionRepo) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	result := s.DB.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &session, result.Error
}

func (s *SessionRepo) UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	err := s.DB.WithContext(ctx).Save(session).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store an updated session in db: %v", err)
	}
	return session, nil
}

func (s *SessionRepo) CreateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) (*models.LastSeenSession, error) {
	err := s.DB.WithContext(ctx).Create(lastSeen).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store a new last seen object to db: %v", err)
	}
	return lastSeen, nil
}

func (s *SessionRepo) GetLastSeenByUserID(ctx context.Context, userID string) (*models.LastSeenSession, error) {
	var lastSeen models.LastSeenSession
	result := s.DB.WithContext(ctx).Where("user_id = ?", userID).First(&lastSeen)

	return &lastSeen, result.Error
}

func (s *SessionRepo) UpdateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) (*models.LastSeenSession, error) {
	err := s.DB.WithContext(ctx).Save(lastSeen).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store an updated last seen object to db: %v", err)
	}
	return lastSeen, nil
}

func (s *SessionRepo) DeleteLastSeen(ctx context.Context, userID string) error {
	err := s.DB.Delete(&models.LastSeenSession{}, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) CreateUserLocation(ctx context.Context, userLocation *models.UserLocation) (*models.UserLocation, error) {
	err := s.DB.WithContext(ctx).Create(userLocation).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store new user location in db: %v", err)
	}
	return userLocation, nil
}

func (s *SessionRepo) GetUserLocationByUserID(ctx context.Context, userID string) (*models.UserLocation, error) {
	var userLocation models.UserLocation
	results := s.DB.WithContext(ctx).Where("user_id = ?", userID).First(&userLocation)
	return &userLocation, results.Error
}

func (s *SessionRepo) UpdateUserLocation(ctx context.Context, userLocation *models.UserLocation) (*models.UserLocation, error) {
	err := s.DB.WithContext(ctx).Save(userLocation).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store updated user location in db: %v", err)
	}
	return userLocation, nil
}

func (s *SessionRepo) DeleteUserLocation(ctx context.Context, userId string) error {
	err := s.DB.Delete(&models.UserLocation{}, userId).Error
	if err != nil {
		return err
	}
	return nil
}
