package models

import "time"

//	The session component shouldn't be deleted only changed
//
// Only valid reason for the session record to be deleted is if the user deletes
// his account
type Session struct {
	//Session ID should be payloadID of refreshToken
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"not null;unique"`
	IsActive 	 bool `gorm:"not null"`
	ConnectedAt  time.Time `gorm:"not null"`
}

type LastSeenSession struct {
	UserID   string `gorm:"primaryKey;foreignKey:UserID;referenced:UserID"`
	Time time.Time  `gorm:"not null"`
}

type UserLocation struct {
	UserID string `gorm:"primaryKey;foreignKey:UserID;referenced:UserID"`
	ApiIP  string `gorm:"not null"`
}
