package db

import (
	"log"
	"os"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBUrl() string {
	err := config.LoadEnv()
	// If LoadEnv returns an error there is no .env file and this is run on railway
	if err != nil {
		return os.Getenv("USER_RAILWAY_DB_URL")
	}
	return config.GetEnv("USER_RAILWAY_DB_URL", "")
}

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return db
}
