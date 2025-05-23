package main

import (
	"log"
	"service-session/internal/db"
	"service-session/internal/handlers"
	"service-session/internal/repository"
	"service-session/internal/routes"
	"service-session/internal/services"
)

func main() {

	db := db.Connect(db.GetDBUrl())
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := services.NewSessionService(sessionRepo)
	sessionHandler := handlers.NewSessionHandler(sessionService)

	router := routes.SetupRouter(sessionHandler)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
