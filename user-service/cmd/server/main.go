package main

import (
	//"context"
	"log"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	//"github.com/Abelova-Grupa/Mercypher/user-service/internal/repository"
)

// TODO: Move to config?
func getDatabaseParameters() string {
	config.LoadEnv()

	user := config.GetEnv("DB_USER", "root")
	pass := config.GetEnv("DB_PASSWORD", "")
	host := config.GetEnv("DB_HOST", "127.0.0.1")
	port := config.GetEnv("DB_PORT", "3306")
	name := config.GetEnv("DB_NAME", "users")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
}

// TODO: Move to config?
func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	log.Println("Attempting to connect to the users database...")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Connected to the users database.")
	}

	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")

	// Auto-migrate 
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("auto-migration failed:", err)
	}

	return db;
}

func main() {
	
	//db := connect(getDatabaseParameters())

	
	// userRepo := repository.NewUserRepository(db)

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
