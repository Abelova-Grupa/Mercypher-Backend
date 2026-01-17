package models

import "time"

type User struct {
	Username     string `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	Validated    bool `gorm:"not null"`
	// TODO: Think if this code should be hashed, if so how would you do it,
	// will this require a constant signature
	AuthCode string
}
