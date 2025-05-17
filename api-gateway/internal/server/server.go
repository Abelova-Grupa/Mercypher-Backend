package api

import (
	"github.com/Abelova-Grupa/Mercypher/api/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func initServer() *Server {
	server := &Server{}
	router := gin.Default()

	//Here goes route handling
	router.GET("/ws", handlers.HandleWebSocket)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponses(err error) gin.H {
	return gin.H{"error": err}
}
