package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	GetSessionByUsername(ctx context.Context, username string) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error)
}

type SessionRepo struct {
	RDB *redis.Client
}

func NewSessionRepository(redis_cli *redis.Client) *SessionRepo {
	return &SessionRepo{RDB: redis_cli}
}

func (s *SessionRepo) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", session.Username)
	m := map[string]interface{}{
		"is_active":      session.IsActive,
		"connected_at":   session.ConnectedAt,
		"last_seen_time": session.LastSeenTime,
	}
	err := s.RDB.HSet(ctx, sessionKey, m).Err()
	if err != nil {
		return nil, fmt.Errorf("unable to store a new session in redis cache: %w", err)
	}

	return session, nil
}

func (s *SessionRepo) GetSessionByUsername(ctx context.Context, username string) (*models.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", username)
	res, err := s.RDB.HGetAll(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, err
	}

	session, err := convertRedisHashToSession(res)
	if err != nil {
		return nil, fmt.Errorf("redis hash to struct conversion failed: %w", err)
	}

	return session, nil
}

func (s *SessionRepo) UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	var res map[string]string
	var err error

	sessionKey := fmt.Sprintf("session:%s", session.Username)
	m := map[string]interface{}{
		"is_active":      session.IsActive,
		"connected_at":   session.ConnectedAt,
		"last_seen_time": session.LastSeenTime,
	}
	err = s.RDB.HSet(ctx, sessionKey, m).Err()
	if err != nil {
		return nil, fmt.Errorf("unable to store a new session in redis cache: %w", err)
	}

	res, err = s.RDB.HGetAll(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, err
	}

	session, err = convertRedisHashToSession(res)
	if err != nil {
		return nil, fmt.Errorf("redis hash to struct conversion failed: %w", err)
	}

	return session, nil
}

func convertRedisHashToSession(m map[string]string) (*models.Session, error) {
	session := &models.Session{Username: m["username"]}
	connectedAt, err := strconv.ParseInt(m["connected_at"], 10, 64)
	if err != nil {
		return nil, err
	}
	last_seen_time, err := strconv.ParseInt(m["last_seen_time"], 10, 64)
	if err != nil {
		return nil, err
	}
	is_active, err := strconv.ParseBool(m["is_active"])
	if err != nil {
		return nil, err
	}

	session.ConnectedAt = connectedAt
	session.LastSeenTime = last_seen_time
	session.IsActive = is_active
	return session, nil
}
