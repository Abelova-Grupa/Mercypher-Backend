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
