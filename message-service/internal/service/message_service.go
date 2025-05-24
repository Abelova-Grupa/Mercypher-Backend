package service

import (
	"context"
	"log"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// ProcessMessage assigns a unique ID to the message, saves it to the database,
// and returns a pointer to the updated message.
func (s *MessageService) ProcessMessage(ctx context.Context, msg *model.ChatMessage) (*model.ChatMessage, error) {

	// Generate new message id
	msg.Message_id = uuid.New().String()
	log.Println(msg)

	// Store the message in database
	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		log.Println("Error saving the file: ", err)
		return nil, err
	}

	// DELEGATE TO SEPARATE UNIT

	// Send ack to sender (with id)
	// Wait for ack from sender

	// Check whether the recipient is online (get api location)

	// (a) Send message to online user
	// (a) Wait for ack

	// (b) Send message to relay service
	// (b) Wait for ack

	return msg, nil

}
