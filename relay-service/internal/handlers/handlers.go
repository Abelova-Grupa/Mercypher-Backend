package handlers

import (
	pb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/repository"
)

func StoreMessage(message *pb.ChatMessage) (status *pb.Status) {
	repository.SaveMessage(message)
	return &pb.Status{Status: 0}
}

func GetMessagesForUserId(id *pb.UserId) []*pb.ChatMessage {
	var result []*pb.ChatMessage
	result, _ = repository.GetMessages(id)
	return result
}
