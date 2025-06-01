package servers

import (
	"log"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router 	*gin.Engine
	wg 		*sync.WaitGroup
	gwIn	chan *domain.Envelope
	gwOut	chan *domain.Envelope
}

func (s *HttpServer) handleLogin(ctx *gin.Context) {
	
}

func (s *HttpServer) handleRegister(ctx *gin.Context) {

}

func (s *HttpServer) handleLogout(ctx *gin.Context) {

}

func handleWebSocket(ctx *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	//ws := NewWe
	//conn, err := websocket.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	//if err != nil {
	//	log.Println("Upgrade error:", err)
	//	return
	//}

	// Handle this client in a new goroutine
	//go websocket.HandleClient(conn)
}

func (s *HttpServer) setupRoutes() {
	
	// HTTP POST request routes.out 
	//
	// Body of a login request should contain an username/email with password.
	// Body of a register request should contain a full user.
	s.router.POST("/login", s.handleLogin)
	s.router.POST("/register", s.handleRegister)

	// HTTP GET requset routes.
	//
	// Websocket route (/ws) must contain a valid token issued by login request.
	s.router.GET("/logout", s.handleLogout)
	s.router.GET("/ws", middleware.AuthMiddleware(), handleWebSocket)


} 

func NewHttpServer(wg *sync.WaitGroup, gwIn chan *domain.Envelope, gwOut chan *domain.Envelope) *HttpServer {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &HttpServer{}
	router := gin.Default()

	// Server parameters
	server.wg = wg
	
	server.router = router
	server.setupRoutes()

	server.gwIn = gwIn
	server.gwOut = gwOut

	return server
}

func (server *HttpServer) Start(address string) {
	server.wg.Add(1)
		
	// HTTP Server must run in its own routine for it has to work concurrently with
	// a gRPC server and main gateway router.
	go func() {
		defer server.wg.Done()
		
		log.Println("HTTP server thread started on: ", address)	
		if err := server.router.Run(address); err != nil {
			log.Fatal("Something went wrong while starting http server.")
		}
	}()	
}
