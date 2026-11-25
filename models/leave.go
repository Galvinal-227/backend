package models

import "time"

type LeaveType string
type LeaveStatus string

const (
	LeaveIzin  LeaveType = "izin"
	LeaveSakit LeaveType = "sakit"

	LeavePending  LeaveStatus = "pending"
	LeaveApproved LeaveStatus = "approved"
	LeaveRejected LeaveStatus = "rejected"
)

type Leave struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"index;not null" json:"user_id"`
	Type      LeaveType   `gorm:"size:16;not null" json:"type"`
	StartDate DateOnly    `gorm:"type:date;not null" json:"start_date"`
	EndDate   DateOnly    `gorm:"type:date;not null" json:"end_date"`
	Reason    string      `gorm:"size:512;not null" json:"reason"`
	Status    LeaveStatus `gorm:"size:16;not null;default:pending" json:"status"`
	AdminNote string      `gorm:"size:512" json:"admin_note"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}