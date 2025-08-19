package models

import "time"

type RefreshToken struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"not null"`
	Token  string    `gorm:"uniqueIndex"`
	Expiry time.Time `gorm:"not null"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
