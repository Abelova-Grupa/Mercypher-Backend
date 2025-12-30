package handlers

import (
	"log"

	pb "github.com/Abelova-Grupa/Mercypher/proto/relay"
	"github.com/Abelova-Grupa/Mercypher/relay-service/internal/repository"
)

func StoreMessage(message *pb.ChatMessage) (status *pb.Status) {
	repository.SaveMessage(message)
	// testing block
	log.Printf("	Stored:\n 	%v\n	\n", message)
	//
	return &pb.Status{Status: 0}
}

func GetMessagesForUserId(id *pb.UserId) []*pb.ChatMessage {
	var result []*pb.ChatMessage
	result, _ = repository.GetMessages(id)
	return result
}
