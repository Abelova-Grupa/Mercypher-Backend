package model

import "time"

type Group struct {
	ID        string	`gorm:"primaryKey;type:uuid"`
	Name      string	`gorm:"not null"`
	OwnerID   string	`gorm:"not null"`
	CreatedAt time.Time	`gorm:"not null"`
}

func NewGroup(id, name, ownerID string) *Group {
	return &Group{
		ID:        id,
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now().UTC(),
	}
}