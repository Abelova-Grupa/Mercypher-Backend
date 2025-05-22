package repository

import (
	"context"
	"errors"
	"service-session/internal/models"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, ID string) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error

	CreateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) error
	GetLastSeenByUserID(ctx context.Context, userID string) (*models.LastSeenSession, error)
	UpdateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) error

	CreateUserLocation(tx context.Context, userLocation *models.UserLocation) error
	GetUserLocationByUserID(tx context.Context, userID string) (*models.UserLocation, error)
	UpdateUserLocation(tx context.Context, userLocation *models.UserLocation) error
}

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &SessionRepo{DB: db}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session *models.Session) error {
	return s.DB.WithContext(ctx).Create(session).Error
}

// Should return payloadID from refreshToken
func (s *SessionRepo) GetSessionByID(ctx context.Context, ID string) (*models.Session, error) {
	var session models.Session
	result := s.DB.WithContext(ctx).Where("id = ?", ID).First(&session)
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

func (s *SessionRepo) UpdateSession(ctx context.Context, session *models.Session) error {
	return s.DB.WithContext(ctx).Save(session).Error
}

func (s *SessionRepo) CreateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) error {
	return s.DB.WithContext(ctx).Create(lastSeen).Error
}

func (s *SessionRepo) GetLastSeenByUserID(ctx context.Context, userID string) (*models.LastSeenSession, error) {
	var lastSeen models.LastSeenSession
	result := s.DB.WithContext(ctx).Where("id = ?", userID).First(&lastSeen)
	return &lastSeen, result.Error
}

func (s *SessionRepo) UpdateLastSeen(ctx context.Context, lastSeen *models.LastSeenSession) error {
	return s.DB.WithContext(ctx).Save(lastSeen).Error
}

func (s *SessionRepo) CreateUserLocation(ctx context.Context, userLocation *models.UserLocation) error {
	return s.DB.WithContext(ctx).Create(userLocation).Error
}

func (s *SessionRepo) GetUserLocationByUserID(ctx context.Context, userID string) (*models.UserLocation, error) {
	var userLocation models.UserLocation
	results := s.DB.WithContext(ctx).Where("id = ?", userID).First(&userLocation)
	return &userLocation, results.Error
}

func (s *SessionRepo) UpdateUserLocation(ctx context.Context, userLocation *models.UserLocation) error {
	return s.DB.WithContext(ctx).Save(userLocation).Error
}
