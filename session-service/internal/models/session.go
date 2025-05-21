package models

type Session struct {
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	AccessToken  string `gorm:"not null"`
}
