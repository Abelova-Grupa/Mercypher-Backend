package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

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
		"username":       session.Username,
		"is_active":      session.IsActive,
		"connected_at":   session.ConnectedAt.Unix(),
		"last_seen_time": session.LastSeenTime.Unix(),
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
	if len(res) == 0 {
		return nil, fmt.Errorf("no session for user %v", username)
	}

	session, err := convertRedisHashToSession(res)
	session.Username = username
	if err != nil {
		return nil, fmt.Errorf("redis hash to struct conversion failed: %w", err)
	}

	return session, nil
}

func (s *SessionRepo) UpdateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	var res map[string]string
	var err error
	m := convertSessionToRedisHash(session)
	sessionKey := fmt.Sprintf("session:%s", session.Username)
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

	session.ConnectedAt = time.Unix(connectedAt, 0)
	session.LastSeenTime = time.Unix(last_seen_time, 0)
	session.IsActive = is_active
	return session, nil
}

func convertSessionToRedisHash(session *models.Session) (m map[string]interface{}) {
	return map[string]interface{}{
		"username":       session.Username,
		"is_active":      session.IsActive,
		"connected_at":   session.ConnectedAt.Unix(),
		"last_seen_time": session.LastSeenTime.Unix(),
	}
}
