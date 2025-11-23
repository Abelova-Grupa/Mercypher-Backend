package routes

import (
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(sessionHandler *handlers.SessionHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	api := r.Group("/session")
	api.GET("/verify", sessionHandler.VerifyToken)
	api.POST("/token", sessionHandler.CreateToken)

	return r
}
