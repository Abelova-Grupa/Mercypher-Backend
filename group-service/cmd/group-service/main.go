package main

import (
	"log"
	"net"
	"os"

	"github.com/Abelova-Grupa/Mercypher/group-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/repository"
	"github.com/Abelova-Grupa/Mercypher/group-service/internal/server"
	grouppb "github.com/Abelova-Grupa/Mercypher/proto/group"


	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	dbURL := getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/mercypher?sslmode=disable")
	port := getEnv("PORT", "50056")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "group_service.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&model.Group{},
		&model.GroupMember{},
	)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	groupRepo := repository.NewGroupRepository(db)

	groupServer := server.NewGroupServer(groupRepo)

	grpcServer := grpc.NewServer()
	grouppb.RegisterGroupServiceServer(grpcServer, groupServer)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Group service running on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}