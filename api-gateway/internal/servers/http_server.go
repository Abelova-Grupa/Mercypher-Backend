package servers

import (
	"log"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api/internal/handlers"
	"github.com/Abelova-Grupa/Mercypher/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
	wg *sync.WaitGroup
	// TODO: Implement incoming and outgoing channels

}

func NewHttpServer(wg *sync.WaitGroup) *HttpServer {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &HttpServer{}
	router := gin.Default()

	// HTTP POST request routes. 
	//
	// Body of a login request should contain an username/email with password.
	// Body of a register request should contain a full user.
	router.POST("/login", handlers.HandleLogin)
	router.POST("/register", handlers.HandleRegister)

	// HTTP GET requset routes.
	//
	// Websocket route (/ws) must contain a valid token issued by login request.
	router.GET("/logout", handlers.HandleLogout)
	router.GET("/ws", middleware.AuthMiddleware(), handlers.HandleWebSocket)

	server.wg = wg
	server.router = router
	return server
}

func (server *HttpServer) Start(address string) {
	defer server.wg.Done()	

	// HTTP Server must run in its own routine for it has to work concurrently with
	// a gRPC server and main gateway router.
	go func() {
		log.Println("HTTP server thread started on ", address)	
		if err := server.router.Run(address); err != nil {
			log.Fatal("Something went wrong while starting http server.")
		}
	}()	
}
