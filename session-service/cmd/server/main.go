package main

import (
	"log"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/handlers"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/routes"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/services"
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
