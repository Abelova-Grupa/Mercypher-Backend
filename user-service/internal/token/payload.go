package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	// ID - meant to generate each time a person logs in 
	ID string  `json:"id"`
	// ServerID - meant to be static part of payload
	ServerID   string `json:"server_id"`
	//userid
	Username    string    `json:"user_id"`
	Revoked bool `json:"revoked"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID:= os.Getenv("USER_SERVICE_UUID")
	if tokenID == "" {
		return nil, fmt.Errorf("Unable to generate uuid for token payload")
	}

	payload := &Payload{
		ID: uuid.NewString(),
		ServerID:        tokenID,
		Username:    username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return errors.New("token has expired")
	}
	return nil
}

// Needs to implement methods from Claims interface

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.ExpiresAt,
	}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
