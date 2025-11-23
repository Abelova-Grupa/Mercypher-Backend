package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKey = 32

type JWTMaker struct {
	secretAccessKey  string
}

func NewJWTMaker(secretAccessKey string, secretRefreshKey string) (*JWTMaker, error) {
	if len(secretAccessKey) < minSecretKey || len(secretRefreshKey) < minSecretKey {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKey)
	}
	return &JWTMaker{secretAccessKey}, nil
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
	} 
	return token, payload, err
}

func (jwtMaker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		} else {
			return []byte(jwtMaker.secretAccessKey), nil
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


