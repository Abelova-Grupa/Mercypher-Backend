package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createFn   func(ctx context.Context, s *models.Session) (*models.Session, error)
	getFn      func(ctx context.Context, username string) (*models.Session, error)
	updateFn   func(ctx context.Context, s *models.Session) (*models.Session, error)
}

func (m *mockRepo) CreateSession(ctx context.Context, s *models.Session) (*models.Session, error) {
	return m.createFn(ctx, s)
}

func (m *mockRepo) GetSessionByUsername(ctx context.Context, username string) (*models.Session, error) {
	return m.getFn(ctx, username)
}

func (m *mockRepo) UpdateSession(ctx context.Context, s *models.Session) (*models.Session, error) {
	return m.updateFn(ctx, s)
}

func TestCreateSession_Success(t *testing.T) {
	now := time.Now()

	repo := &mockRepo{
		createFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
	}

	service := NewSessionService(repo)

	resp, err := service.CreateSession(context.Background(), "alice", now)

	assert.NoError(t, err)
	assert.Equal(t, "alice", resp.Username)
	assert.True(t, resp.IsActive)
}

func TestCreateSession_InvalidParams(t *testing.T) {
	service := NewSessionService(nil)

	_, err := service.CreateSession(context.Background(), "", time.Now())

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidParams, err)
}

func TestGetSessionByUsername_Success(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, username string) (*models.Session, error) {
			return &models.Session{
				Username:    username,
				IsActive:    true,
				ConnectedAt: time.Now(),
			}, nil
		},
	}

	service := NewSessionService(repo)

	resp, err := service.GetSessionByUsername(context.Background(), "alice")

	assert.NoError(t, err)
	assert.Equal(t, "alice", resp.Username)
}

func TestConnect_CreateNewSession(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, username string) (*models.Session, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
		updateFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
	}

	service := NewSessionService(repo)

	err := service.Connect(context.Background(), "alice")

	assert.NoError(t, err)
}

func TestConnect_UpdateExistingSession(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, username string) (*models.Session, error) {
			return &models.Session{
				Username: username,
				IsActive: false,
			}, nil
		},
		createFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
		updateFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
	}

	service := NewSessionService(repo)

	err := service.Connect(context.Background(), "alice")

	assert.NoError(t, err)
}

func TestDisconnect_Success(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, username string) (*models.Session, error) {
			return &models.Session{
				Username: username,
				IsActive: true,
			}, nil
		},
		updateFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
		createFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
	}

	service := NewSessionService(repo)

	err := service.Disconnect(context.Background(), "alice")

	assert.NoError(t, err)
}

func TestDisconnect_NotFound(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, username string) (*models.Session, error) {
			return nil, errors.New("not found")
		},
		updateFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
		createFn: func(ctx context.Context, s *models.Session) (*models.Session, error) {
			return s, nil
		},
	}

	service := NewSessionService(repo)

	err := service.Disconnect(context.Background(), "alice")

	assert.Error(t, err)
}