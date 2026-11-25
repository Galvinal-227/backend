package models

import "time"

type ActivityLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Action      string    `gorm:"size:64;not null;index" json:"action"`
	Description string    `gorm:"size:512" json:"description"`
	Metadata    string    `gorm:"type:text" json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}