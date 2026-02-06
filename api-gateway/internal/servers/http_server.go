package servers

import (
	// "encoding/json"

	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/clients"
	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/middleware"
	"github.com/Abelova-Grupa/Mercypher/api-gateway/internal/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HttpServer serves incoming http requests (excluding websocket) such as login and
// register.
//
// Note to self:	It could be more optimal to remove register and unregister channels,
//
//	and to define envelope messages for that purpose. Something that
//	should be tested in the future.
type HttpServer struct {
	router     *gin.Engine               // HTTP Servers internal gin router
	wg         *sync.WaitGroup           // Wait group that holds for HTTP server routine
	gwIn       chan *domain.Envelope     // Channel for sending envelopes to gateway
	gwOut      chan *domain.Envelope     // Channel for receiving envelopes from gateway
	register   chan *websocket.Websocket // Channel for registering new user in gateway
	unregister chan *websocket.Websocket // Channel for unregistering user from gateway

	userClient    *clients.UserClient    // Temporary solution for handling login requests
	sessionClient *clients.SessionClient // Temporary solution for handling token validation
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ContactRequest struct {
	Contact string `json:"contact"`
	Nickname string `json:"nickname"`
}

func (s *HttpServer) handleLogin(ctx *gin.Context) {

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	//Quick fix
	token, err := s.userClient.Login(domain.User{Username: req.Username}, req.Password, req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist"})
		return
	}

	ctx.SetCookie(
		"access_token",
		token,
		9000,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token, // Delete this after testing
	})
}

func (s *HttpServer) handleRegister(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, err := s.userClient.Register(domain.User{Username: req.Username, Email: req.Email}, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Couldn't register user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      id,
	})
}

func (s *HttpServer) handleLogout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Byeee",
	})
}

func (s *HttpServer) handleMe(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": userID})
}

func (s *HttpServer) handleCreateContact(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	username := fmt.Sprint(userID)

	if err := s.userClient.CreateContact(username, req.Contact, req.Nickname); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create contact."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Contact saved."})
	}
}

func (s *HttpServer) handleDeleteContact(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	username := fmt.Sprint(userID)

	if err := s.userClient.DeleteContact(username, req.Contact); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete contact."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Contact deleted."})
	}
}

func (s *HttpServer) handleGetcontacts(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username := fmt.Sprint(userID)

	if resp, err := s.userClient.GetContacts(username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve contacts."})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Success.", "contacts":resp})
	}



}

func (s *HttpServer) handleWebSocket(ctx *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ws := websocket.NewWebsocket(conn, domain.User{
		UserId:   fmt.Sprint(userID), // TODO: Remove id for its the same as username
		Username: fmt.Sprint(userID),
		Email:    "", // Nil here because (for now) we are only iterested in username
	}, s.unregister, s.gwIn)

	//TODO: Register this ws in gateway.
	s.register <- ws

	// Handle this client in a new goroutine    // HttpOnly
	go ws.HandleClient()
}

func (s *HttpServer) setupRoutes() {

	// HTTP POST request routes
	//
	// Body of a login request should contain an username/email with password.
	// Body of a register request should contain a full user.
	//
	// Check README.md (for api gateway) for more detailed info about format.
	s.router.POST("/login", s.handleLogin)
	s.router.POST("/register", s.handleRegister)
	s.router.POST("/createContact", middleware.AuthMiddleware(s.userClient), s.handleCreateContact)
	s.router.POST("/deleteContact", middleware.AuthMiddleware(s.userClient), s.handleDeleteContact)

	// HTTP GET requset routes.
	//
	// Websocket route (/ws) must contain a valid token issued by login request.
	s.router.GET("/logout", s.handleLogout)
	s.router.GET("/ws", middleware.AuthMiddleware(s.userClient), s.handleWebSocket)
	s.router.GET("/me", middleware.AuthMiddleware(s.userClient), s.handleMe)
	s.router.GET("/contacts", middleware.AuthMiddleware(s.userClient), s.handleGetcontacts)
}

func NewHttpServer(wg *sync.WaitGroup, gwIn chan *domain.Envelope, gwOut chan *domain.Envelope, reg chan *websocket.Websocket, unreg chan *websocket.Websocket) *HttpServer {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &HttpServer{}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Clients to other serivces
	server.userClient, _ = clients.NewUserClient("localhost:50054")
	server.sessionClient, _ = clients.NewSessionClient("localhost:50055")

	// Server parameters
	server.wg = wg

	server.router = router
	server.setupRoutes()

	server.gwIn = gwIn
	server.gwOut = gwOut

	server.register = reg
	server.unregister = unreg

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
