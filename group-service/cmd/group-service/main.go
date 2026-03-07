package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/Abelova-Grupa/Mercypher/group-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/server"
	grouppb "github.com/Abelova-Grupa/Mercypher/proto/group"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnv()
	port := config.GetEnv("PORT", "50056")

	db, err := NewMessageDB(context.Background())

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	err = db.AutoMigrate(
		&model.Group{},
		&model.GroupMember{},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate")
	}

	groupRepo := repository.NewGroupRepository(db)

	groupServer := server.NewGroupServer(groupRepo)

	grpcServer := grpc.NewServer()
	grouppb.RegisterGroupServiceServer(grpcServer, groupServer)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	log.Printf("Group service running on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
}

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
	// log.Info().Str("db_url", dbUrl.String()).Msg("constructed database URL")

	if err := migrateDB(dbUrl); err != nil {
		return nil, err
	}

	db, err := connectDB(dbUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB(dbUrl *url.URL) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl.String()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "group_service.",
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

func migrateDB(migrateUrl *url.URL) error {
	fmt.Println(migrateUrl.String())
	var m *migrate.Migrate
	var err error
	for i := 0; i < 10; i++ {
		m, err = migrate.New("file://./internal/migrations", migrateUrl.String())
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
