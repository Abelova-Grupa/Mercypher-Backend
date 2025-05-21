package api

import (
	"log"

	"github.com/Abelova-Grupa/Mercypher/api/internal/handlers"
	"github.com/Abelova-Grupa/Mercypher/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func InitServer() *Server {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &Server{}
	router := gin.Default()

	router.POST("/login", handlers.HandleLogin)
	router.POST("/register", handlers.HandleRegister)

	router.GET("/logout", handlers.HandleLogout)
	router.GET("/user", handlers.HandleSearchUser)
	router.GET("/ws", middleware.AuthMiddleware() ,handlers.HandleWebSocket)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	log.Println("Server started on: ", address)	
	return server.router.Run(address)
}
