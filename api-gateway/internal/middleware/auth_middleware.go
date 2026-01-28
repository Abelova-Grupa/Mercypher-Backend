package middleware

import (
	"net/http"

	// "strings"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/clients"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sc *clients.UserClient) gin.HandlerFunc {
	// return func(c *gin.Context) {
	// 	authHeader := c.GetHeader("Authorization")
	// 	if authHeader == "" {
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
	// 		return
	// 	}

	// 	parts := strings.SplitN(authHeader, " ", 2)
	// 	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
	// 		return
	// 	}

	// 	token := parts[1]

	// 	isValid, _ := sc.VerifyToken(token)

	// 	if !isValid {
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
	// 		return
	// 	}

	// 	c.Next()
	// }
	return func(c *gin.Context) {
        token, err := c.Cookie("access_token")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Couldn't read token"})
            return
        }

        valid, err := sc.VerifyToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
            return
        }
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		userid, err := sc.DecodeToken(token)
		if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Couldn't extract user data from token"})
            return
        }

        c.Set("userID", userid)
        c.Next()
    }
}

