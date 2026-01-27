package db

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	// "github.com/Abelova-Grupa/Mercypher/user-service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetDBUrl() string {
	err := config.LoadEnv()
	// TODO: Remove these lines because Railway isn't used no more
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

	host := os.Getenv("POSTGRES_HOST")
	if host == "" || host == "localhost" {
		host = "localhost"
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	log.Info().Msg(dsn)

	// Before gorm starts, run migrations

	migrateUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
        user, password, host, port, dbname)

	m, err := migrate.New("file://internal/migrations", migrateUrl)
	
	if err != nil {
        log.Fatal().Err(err).Msg("Failed to initialize migration engine")
    }

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal().Err(err).Msg("Migrations failed.")
    }

    log.Info().Msg("Migrations applied successfully!")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "user_service.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// db.Exec("CREATE SCHEMA IF NOT EXISTS user_service")

	// err = db.AutoMigrate(&models.User{})
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to migrate database")
	// }

	return db
}
