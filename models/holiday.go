package models

import "time"

type Holiday struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Date      DateOnly  `gorm:"type:date;uniqueIndex;not null" json:"date"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}