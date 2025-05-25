package servers

import (
	"log"

	"github.com/Abelova-Grupa/Mercypher/api/internal/handlers"
	"github.com/Abelova-Grupa/Mercypher/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
}

func NewHttpServer() *HttpServer {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &HttpServer{}
	router := gin.Default()

	router.POST("/login", handlers.HandleLogin)
	router.POST("/register", handlers.HandleRegister)

	router.GET("/logout", handlers.HandleLogout)
	router.GET("/user", handlers.HandleSearchUser)
	router.GET("/ws", middleware.AuthMiddleware() ,handlers.HandleWebSocket)

	server.router = router
	go server.Start(":8080")
	return server
}

func (server *HttpServer) Start(address string) error {
	log.Println("Server started on: ", address)	
	return server.router.Run(address)
}
