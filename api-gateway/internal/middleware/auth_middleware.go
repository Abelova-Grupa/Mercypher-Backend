package middleware

import (
	"net/http"
	// "strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		// 	return
		// }

		// token := strings.TrimPrefix(authHeader, "Bearer ")
		// if token == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		// 	return
		// }

		tempToken := "Remove me"
		if !isTokenValid(tempToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Next()
	}
}

// Temporary stub — replace with real validation
func isTokenValid(token string) bool {

	return true // For demo purposes only
}
