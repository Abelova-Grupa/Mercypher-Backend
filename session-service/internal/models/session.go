package models

import "time"

//	The session component shouldn't be deleted only changed
// Only valid reason for the session record to be deleted is if the user deletes his account
type Session struct {
	ID           string `gorm:"primaryKey"`
	Username       string `gorm:"primaryKey"`
	IsActive 	 bool `gorm:"not null"`
	ConnectedAt  time.Time `gorm:"not null"`
	LastSeenTime time.Time `gorm:"index"`
}

type Token struct {
	Text string
}
