package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/Abelova-Grupa/Mercypher/session-service/internal/db"
	pb "github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/pb"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Loading env variables from cloud
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	var creds credentials.TransportCredentials
	var err error
	// These variables are only store on the cloud
	if tlsCert == "" || tlsKey == "" {
		creds, err = credentials.NewServerTLSFromFile("../../internal/certs/server.crt", "../../internal/certs/server.key")
		if err != nil {
			log.Fatalf("Failed to load TLS keys: %v", err)
		}
	} else {
		// Creating a certificate : key pair
		cert, err := tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey))
		if err != nil {
			log.Fatalf("Failed to generate x509 pair: %v", err)
		}
		// Creating tls configuration based on certificate pair
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		// Creating credentials
		creds = credentials.NewTLS(tlsConfig)
	}

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterSessionServiceServer(grpcServer, server.NewGrpcServer(db.Connect(db.GetDBUrl())))

	log.Println("Starting gRPC server on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
