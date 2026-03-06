package store

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewMessageDB(ctx context.Context) (*gorm.DB, error) {
	config.LoadEnv()
	env := os.Getenv("ENVIRONMENT")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	var sslMode string
	if env == "azure" {
		sslMode = "sslmode=require"
	} else {
		sslMode = "sslmode=disable"
	}

	dbUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     "/" + dbName,
		RawQuery: sslMode,
	}

	if err := migrateDB(dbUrl); err != nil {
		return nil, err
	}

	db, err := connectDB(dbUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(migrateUrl *url.URL) error {
	fmt.Println(migrateUrl.String())
	var m *migrate.Migrate
	var err error
	for i := 0; i < 10; i++ {
		m, err = migrate.New("file://internal/migrations", migrateUrl.String())
		if err == nil {
			break
		}
		log.Info().Msg("DB not ready, retrying in 2 seconds...")
		log.Info().Err(err).Msgf("Attempt %d: DB not ready, retrying...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize migration engine")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Migrations failed.")
	}

	log.Info().Msg("Migrations applied successfully!")
	return nil
}

func connectDB(dbUrl *url.URL) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl.String()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "message_service.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to message database")
		return nil, err
	} else {
		log.Info().Msg("successfully connected to the message database")
	}

	return db, nil
}
