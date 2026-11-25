package models

import "time"

type Student struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	NIS          string    `gorm:"size:32;index" json:"nis"`
	Class        string    `gorm:"size:64" json:"class"`
	Major        string    `gorm:"size:128" json:"major"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Address      string    `gorm:"size:255" json:"address"`
	PKLStartDate *DateOnly `gorm:"type:date" json:"pkl_start_date"`
	PKLEndDate   *DateOnly `gorm:"type:date" json:"pkl_end_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}