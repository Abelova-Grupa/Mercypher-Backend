package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	GetSessionByUsername(ctx context.Context, username string) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error)

}

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepo {
	return &SessionRepo{DB: db}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	if session.Username == "" {
		return nil, errors.New("Username cannot be empty during session creation")
	}

	if session.ConnectedAt.IsZero() {
		session.ConnectedAt = time.Now()
	}

	err := s.DB.WithContext(ctx).Create(session).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store a new session in db: %v", err)
	}
	return session, nil
}

func (s *SessionRepo) GetSessionByUsername(ctx context.Context, username string) (*models.Session, error) {
	var session models.Session
	result := s.DB.WithContext(ctx).Where("username = ?",username).First(&session)
	if errors.Is(result.Error,gorm.ErrRecordNotFound){
		return nil, result.Error
	}
	return &session, nil
}

func (s *SessionRepo) UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	err := s.DB.WithContext(ctx).Save(session).Error
	if err != nil {
		return nil, fmt.Errorf("unable to store an updated session in db: %v", err)
	}
	return session, nil
}
