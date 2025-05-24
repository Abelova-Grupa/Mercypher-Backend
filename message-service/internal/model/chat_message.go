package model

import "time"

type ChatMessage struct {
	Message_id	string		`gorm:"primaryKey"`
	Sender_id	string		`gorm:"not null"`
	Receiver_id	string		`gorm:"not null"`
	Body		string		`gorm:"not null"`
	Timestamp	time.Time	`gorm:"not null"`
}
