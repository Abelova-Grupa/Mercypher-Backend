package db

import (
	"log"
	"service-session/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBUrl() string {
	config.LoadEnv()
	return config.GetEnv("SESSION_RAILWAY_DB_URL", "")
}

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return db
}
