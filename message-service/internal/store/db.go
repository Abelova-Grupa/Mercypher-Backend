package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/rs/zerolog/log"
)

func NewMessageDB(ctx context.Context) (*sql.DB, error) {
	config.LoadEnv()
	if os.Getenv("ENVIRONMENT") == "azure" {
		return NewPostgresAzure(ctx)
	}
	return NewPostgresLocal()
}

func NewPostgresLocal() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL") // e.g., postgres://user:pass@localhost:5432/dbname?sslmode=disable
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info().Msg("successfully connected to local postgres")
	return db, nil
}

func NewPostgresAzure(ctx context.Context) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Get Token from Azure Entra ID
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure credentials: %w", err)
	}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://ossrdbms-aad.database.windows.net/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get azure auth token: %w", err)
	}

	// Password is the OAuth2 token
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, token.Token, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool limits for production
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Info().Msg("successfully connected to azure postgres via entra id")
	return db, nil
}
