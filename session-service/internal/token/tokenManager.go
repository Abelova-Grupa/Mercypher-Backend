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
	secretKey string
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken               string    `json:"access_token"`
	AccessTokenExpirationTime time.Time `json:"access_token_expires_at"`
}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKey {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKey)
	}
	return &JWTMaker{secretKey}, nil
}

func (jwtMaker *JWTMaker) CreateToken(userID string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(userID, role, duration, tokenType)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(jwtMaker.secretKey))
	return token, payload, err
}

func (jwtMaker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtMaker.secretKey), nil
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
	session, err := sessionRepo.GetSessionByID(ctx, refreshPayload.ID.String())
	if err != nil {
		// TODO: Make more errors for token handling
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidToken
		}
		return "", ErrInvalidToken
	}

	// Create a new Access token
	accessToken, accessPayload, err := jwtMaker.CreateToken(refreshPayload.UserID, refreshPayload.Role, time.Minute*15, tokenTypeAccessToken)
	if accessToken == "" || accessPayload == nil || err != nil {
		return "", ErrInvalidToken
	}
	session.AccessToken = accessToken
	if _, err := sessionRepo.UpdateSession(ctx, session); err != nil {
		return "", errors.New("error caused during session update")
	}
	return accessToken, nil
}
