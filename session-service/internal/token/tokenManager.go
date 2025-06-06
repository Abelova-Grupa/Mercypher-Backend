package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const minSecretKey = 32

type JWTMaker struct {
	secretAccessKey  string
	secretRefreshKey string
}

func NewJWTMaker(secretAccessKey string, secretRefreshKey string) (*JWTMaker, error) {
	if len(secretAccessKey) < minSecretKey || len(secretRefreshKey) < minSecretKey {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKey)
	}
	return &JWTMaker{secretAccessKey, secretRefreshKey}, nil
}

func (jwtMaker *JWTMaker) CreateToken(userID string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	var err error
	payload, err := NewPayload(userID, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	var token string
	if tokenType == TokenTypeAccessToken {
		token, err = jwtToken.SignedString([]byte(jwtMaker.secretAccessKey))
	} else {
		token, err = jwtToken.SignedString([]byte(jwtMaker.secretRefreshKey))
	}
	return token, payload, err
}

func (jwtMaker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		if tokenType == TokenTypeAccessToken {
			return []byte(jwtMaker.secretAccessKey), nil
		} else {
			return []byte(jwtMaker.secretRefreshKey), nil
		}
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (jwtMaker *JWTMaker) RefreshToken(ctx context.Context, refreshToken string, tokenType TokenType) (string, error) {
	// Validating refresh token
	refreshPayload, err := jwtMaker.VerifyToken(refreshToken, tokenType)
	if err != nil {
		return "", ErrInvalidToken
	}
	if refreshPayload.ExpiresAt.Before(time.Now()) {
		return "", ErrExpiredToken
	}

	// Checking if the token exists in database
	sessionRepo := &repository.SessionRepo{DB: &gorm.DB{}}
	session, err := sessionRepo.GetSessionByUserID(ctx, refreshPayload.UserID)
	if err != nil {
		// TODO: Make more errors for token handling
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidToken
		}
		return "", ErrInvalidToken
	}

	// Create a new Access token
	accessToken, accessPayload, err := jwtMaker.CreateToken(refreshPayload.UserID, time.Minute*15, TokenTypeAccessToken)
	if accessToken == "" || accessPayload == nil || err != nil {
		return "", ErrInvalidToken
	}
	session.AccessToken = accessToken
	if _, err := sessionRepo.UpdateSession(ctx, session); err != nil {
		return "", errors.New("error caused during session update")
	}
	return accessToken, nil
}
