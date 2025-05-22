package models

import (
	"time"
)

//	The session component shouldn't be deleted only changed
//
// Only valid reason for the session record to be deleted is if the user deletes
// his account
type Session struct {
	//Session ID should be payloadID of refreshToken
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"not null;unique"`
	RefreshToken string `gorm:"not null;unique"`
	AccessToken  string `gorm:"not null;unique"`
}

type LastSeenSession struct {
	UserID   string    `gorm:"primaryKey;foreignKey:UserID;referenced:UserID"`
	LastSeen time.Time `gorm:"not null"`
}

type UserLocation struct {
	UserID string `gorm:"primaryKey;foreignKey:UserID;referenced:UserID"`
	ApiIP  string `gorm:"not null"`
}
