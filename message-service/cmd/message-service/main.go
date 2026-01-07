package main

//"context"

//"time"

func main() {

	// Old code kept for later
	// 	conn := connect(getDatabaseParameters())
	// 	repo := repository.NewMessageRepository(conn)
	// 	service := service.NewMessageService(repo)
	// 	server := server.NewMessageServer(service)

	// 	listener, err := net.Listen("tcp", ":50051")
	// 	if err != nil {
	// 		log.Fatalf("Failed to listen: %v", err)
	// 	}

	// 	grpcServer := grpc.NewServer()
	// 	messagepb.RegisterMessageServiceServer(grpcServer, server)
	// 	if err := grpcServer.Serve(listener); err != nil {
	// 		log.Fatalf("Failed to serve: %v", err)
	// 	}

	// server := server.NewMessageServer()
	// var msg model.ChatMessage
	// msg.Message_id = "test_DZJLKAFSKJGDKJHFGJHLI"
	// msg.Sender_id = "testuser2"
	// msg.Receiver_id = "testuser1"
	// msg.Body = "Proba proba jen dva tri"
	// msg.Timestamp = time.Now()
	// repo.CreateMessage(context.Background(), &msg)
}
