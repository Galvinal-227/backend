package models

import "time"

type JournalStatus string

const (
	JournalPending  JournalStatus = "pending"
	JournalApproved JournalStatus = "approved"
	JournalRejected JournalStatus = "rejected"
)

type Journal struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	UserID         uint          `gorm:"uniqueIndex:idx_journal_user_date;not null" json:"user_id"`
	Date           DateOnly      `gorm:"uniqueIndex:idx_journal_user_date;type:date;not null" json:"date"`
	Content        string        `gorm:"type:text;not null" json:"content"`
	Status         JournalStatus `gorm:"size:16;not null;default:pending" json:"status"`
	AdminNote      string        `gorm:"size:512" json:"admin_note"`
	OverrideReason string        `gorm:"size:512" json:"override_reason"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}