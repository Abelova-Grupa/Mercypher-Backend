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

type Contact struct {
	Username1  string `gorm:"primaryKey"`
	Username2  string `gorm:"primaryKey"`
	FirstUser  User   `gorm:"foreignKey:Username1;contraint:OnDelete:CASCADE"`
	SecondUser User   `gorm:"foreignKey:Username2;contraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
}
