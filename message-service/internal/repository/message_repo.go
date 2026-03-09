package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg *model.ChatMessage) error
	GetChatHistory(ctx context.Context, p1, p2 string, lastSeen time.Time, limit int) ([]model.ChatMessage, error)
	GetGroupHistory(ctx context.Context, participant2 string, lastSeen time.Time, limit int) ([]model.ChatMessage, error)
}

type MessageRepo struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepo {
	return &MessageRepo{DB: db}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, msg *model.ChatMessage) error {
	// query := `INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
	//           VALUES ($1, $2, $3, $4, $5)`

	// _, err := r.DB.ExecContext(ctx, query, msg.Message_id, msg.Sender_id, msg.Receiver_id, msg.Body, time.Now())
	// if err != nil {
	// 	return fmt.Errorf("failed to insert message: %w", err)
	// }
	// return nil
	// Postavljamo trenutno vreme ako nije već postavljeno
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	if err := r.DB.WithContext(ctx).Create(msg).Error; err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

// when getting 1st set of messages in front, request with lastSeen = now
func (r *MessageRepo) GetChatHistory(ctx context.Context, p1, p2 string, lastSeen time.Time, limit int) ([]model.ChatMessage, error) {
	// query := `
	//     SELECT id, sender_id, receiver_id, content, created_at
	//     FROM messages
	//     WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
	//       AND created_at < $3
	//     ORDER BY created_at DESC
	//     LIMIT $4`

	// rows, err := r.DB.QueryContext(ctx, query, p1, p2, lastSeen, limit)
	// if err != nil {
	// 	return nil, fmt.Errorf("querying history with cursor: %w", err)
	// }
	// defer rows.Close()

	// var messages []model.ChatMessage
	// for rows.Next() {
	// 	var m model.ChatMessage
	// 	if err := rows.Scan(&m.Message_id, &m.Sender_id, &m.Receiver_id, &m.Body, &m.Timestamp); err != nil {
	// 		return nil, err
	// 	}
	// 	messages = append(messages, m)
	// }
	// return messages, nil
	var messages []model.ChatMessage

	err := r.DB.WithContext(ctx).
		Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))", p1, p2, p2, p1).
		Where("timestamp < ?", lastSeen).
		Order("timestamp DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching chat history: %w", err)
	}

	return messages, nil
}

func (r *MessageRepo) GetGroupHistory(ctx context.Context, participant2 string, lastSeen time.Time, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage

	// In a group, we only care about messages SENT TO the group ID
	err := r.DB.WithContext(ctx).
		Where("receiver_id = ?", participant2).
		Where("timestamp < ?", lastSeen).
		Order("timestamp DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching group history: %w", err)
	}

	return messages, nil
}
