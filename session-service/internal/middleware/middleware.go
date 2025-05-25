package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadkey = "authorization_payload"
)

func authMiddleware(tokenMaker token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		var req struct {
			TokenType token.TokenType `json:"token_type"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		tokenType := req.TokenType

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token, tokenType)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		ctx.Set(authorizationPayloadkey, payload)
		ctx.Next()
	}
}
