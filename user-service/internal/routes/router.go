package routes

import (
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {

	// Biraj koji ces, samo nemoj oba
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
	}

	return r
}
