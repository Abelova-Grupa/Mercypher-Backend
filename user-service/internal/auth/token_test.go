package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	// Example for testing
	maker, err := NewJWTMaker("Dk4cLr7zvUeFYRAxmPlgwXqJ3uEZntBS")
	require.NoError(t, err)

	username := "Cole"
	role := "admin"
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration, tokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker("Dk4cLr7zvUeFYRAxmPlgwXqJ3uEZntBS")
	require.NoError(t, err)

	token, payload, err := maker.CreateToken("Cole", "admin", -time.Minute, tokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, 1)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload("Cole", "admin", time.Minute, tokenTypeAccessToken)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker("Dk4cLr7zvUeFYRAxmPlgwXqJ3uEZntBS")
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token, tokenTypeAccessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestJWTWrongTokenType(t *testing.T) {
	maker, err := NewJWTMaker("Dk4cLr7zvUeFYRAxmPlgwXqJ3uEZntBS")
	require.NoError(t, err)

	token, payload, err := maker.CreateToken("Cole", "admin", time.Minute, tokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, anotherTokenType)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}
