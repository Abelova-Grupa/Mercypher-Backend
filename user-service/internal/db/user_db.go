package db

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


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

	var sslMode string
	if os.Getenv("ENVIRONMENT") == "azure" {
		sslMode = "sslmode=require"
	} else{
		sslMode = "sslmode=disable"
	}

	migrateUrl := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(user,password),
		Host: fmt.Sprintf("%s:%s",host,port),
		Path: "/" + dbname,
		RawQuery: sslMode,
	}

	var m *migrate.Migrate
	for i := 0; i < 10; i++ {
		m, err = migrate.New("file://internal/migrations", migrateUrl.String())
		if err == nil {
			break
    }
		log.Info().Msg("DB not ready, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
}
	
	if err != nil {
        log.Fatal().Err(err).Msg("Failed to initialize migration engine")
    }

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal().Err(err).Msg("Migrations failed.")
    }

	log.Info().Msg("Migrations applied successfully!")

	db, err := gorm.Open(postgres.Open(migrateUrl.String()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "user_service.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	return db
}
