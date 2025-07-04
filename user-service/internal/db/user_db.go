package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetDBUrl() string {
	err := config.LoadEnv()
	// If LoadEnv returns an error there is no .env file and this is run on railway
	if err != nil {
		return os.Getenv("USER_LOCAL_DB_URL")
	}
	return config.GetEnv("USER_LOCAL_DB_URL", "")
}

func Connect() *gorm.DB {
	err := config.LoadEnv()
	if err != nil {
		return nil
	}

	var host string
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		host = "localhost"
	} else {
		host = "user-db"
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "user_service.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	db.Exec("CREATE SCHEMA IF NOT EXISTS user_service")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	return db
}
