package main

import (
	//"context"

	"log"
	"net"

	"google.golang.org/grpc"

	userpb "github.com/Abelova-Grupa/Mercypher/proto/user"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/config"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/user-service/internal/grpc/server"
)

// TODO: Move to config?
// func getDatabaseParameters() string {
// 	config.LoadEnv()

// 	user := config.GetEnv("DB_USER", "postgres")
// 	pass := config.GetEnv("DB_PASSWORD", "")
// 	host := config.GetEnv("DB_HOST", "127.0.0.1")
// 	port := config.GetEnv("DB_PORT", "5432")
// 	name := config.GetEnv("DB_NAME", "users")
// 	ssl := config.GetEnv("DB_SSLMODE", "disable")
// 	tz := config.GetEnv("DB_TIMEZONE", "UTC")

// 	return fmt.Sprintf(
// 		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
// 		user, pass, host, port, name, ssl, tz,
// 	)
// }

func main() {
	conn := db.Connect()

	port := config.GetEnv("USER_SERVICE_PORT", "")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, server.NewGrpcServer(conn))

	log.Printf("starting user service grpc server on port %v", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Setup the router and start routing
	// router := routes.SetupRouter(userHandler)
	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatal("Failed to start server:", err)
	// }

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
