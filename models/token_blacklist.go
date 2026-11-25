package models

import "time"

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JTI       string    `gorm:"size:64;uniqueIndex;not null" json:"jti"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	Reason    string    `gorm:"size:64" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}