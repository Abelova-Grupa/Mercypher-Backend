package services

import (
	"context"
	"time"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
)

type SessionService struct {
	repo repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

// Think about which services should session have
func (s *SessionService) CreateToken(ctx context.Context, userID string, role string, duration time.Duration, tokenType token.TokenType) (string, *token.Payload, error) {
	jwtMaker := token.JWTMaker{}
	token, payload, err := jwtMaker.CreateToken(userID, role, duration, tokenType)
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
