package main

import (
	//"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/handlers"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/routes"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/service"
)

// TODO: Move to config?
func getDatabaseParameters() string {
    config.LoadEnv()

    user := config.GetEnv("DB_USER",     "postgres")
    pass := config.GetEnv("DB_PASSWORD", "")
    host := config.GetEnv("DB_HOST",     "127.0.0.1")
    port := config.GetEnv("DB_PORT",     "5432")
    name := config.GetEnv("DB_NAME",     "users")
    ssl  := config.GetEnv("DB_SSLMODE",  "disable")
    tz   := config.GetEnv("DB_TIMEZONE", "UTC")   

    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
        user, pass, host, port, name, ssl, tz,
    )
}

// TODO: Move to config?
func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("Attempting to connect to the users database...")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Connected to the users database.")
	}

	// Auto-migrate 
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("auto-migration failed:", err)
	}

	return db;
}

func main() {
	
	db := connect(getDatabaseParameters())
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Setup the router and start routing
	router := routes.SetupRouter(userHandler)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// TESTING
	// test_user := models.User{
	// 	ID: "1",
	// 	Username: "jezdimir1",
	// 	Email: "jezdimir.bekrija1@gmail.com",
	// 	PasswordHash: "RodjaRaicevic123",
	// }

	// test_user2, _ := userRepo.GetUserByID(context.Background(), "0")
	// log.Println(*test_user2)





	//userRepo.CreateUser(context.Background(), &test_user)
}
