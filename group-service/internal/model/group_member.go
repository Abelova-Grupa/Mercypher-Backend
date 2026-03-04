package model

import "time"

type GroupMember struct {
	GroupID  string		`gorm:"primaryKey;type:uuid"`
	UserID   string		`gorm:"primaryKey"`
	Role     GroupRole	`gorm:"not null"`
	JoinedAt time.Time	`gorm:"not null"`
}

func NewOwnerMember(groupID, userID string) *GroupMember {
	return &GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     RoleOwner,
		JoinedAt: time.Now().UTC(),
	}
}