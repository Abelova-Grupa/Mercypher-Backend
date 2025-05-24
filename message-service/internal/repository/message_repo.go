package repository

import (
	"context"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *model.ChatMessage) error
	GetMessageById(ctx context.Context, id string) (*model.ChatMessage, error)
	UpdateMessage(ctx context.Context, msg *model.ChatMessage) error }

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepo{db: db}
}

func (r *messageRepo) CreateMessage(ctx context.Context, msg *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *messageRepo) GetMessageById(ctx context.Context, id string) (*model.ChatMessage, error) {
	var message model.ChatMessage
	result := r.db.WithContext(ctx).Where("message_id = ?", id).First(&message)
	return &message, result.Error
}

func (r *messageRepo) UpdateMessage(ctx context.Context, msg *model.ChatMessage) error {
	return r.db.WithContext(ctx).Save(msg).Error
}


