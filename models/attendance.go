package models

import "time"

type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusLate    AttendanceStatus = "late"
	StatusAlpha   AttendanceStatus = "alpha"
	StatusIzin    AttendanceStatus = "izin"
	StatusSakit   AttendanceStatus = "sakit"
)

type Attendance struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserID uint      `gorm:"uniqueIndex:idx_att_user_date;not null" json:"user_id"`
	Date   DateOnly  `gorm:"uniqueIndex:idx_att_user_date;type:date;not null" json:"date"`

	CheckInAt       *time.Time `json:"check_in_at"`
	CheckInLat      *float64   `json:"check_in_lat"`
	CheckInLng      *float64   `json:"check_in_lng"`
	CheckInAddress  string     `gorm:"size:512" json:"check_in_address"`

	CheckOutAt      *time.Time `json:"check_out_at"`
	CheckOutLat     *float64   `json:"check_out_lat"`
	CheckOutLng     *float64   `json:"check_out_lng"`
	CheckOutAddress string     `gorm:"size:512" json:"check_out_address"`

	Status AttendanceStatus `gorm:"size:16;not null" json:"status"`

	CheckInOutsideGeofence  bool `gorm:"not null;default:false" json:"check_in_outside_geofence"`
	CheckOutOutsideGeofence bool `gorm:"not null;default:false" json:"check_out_outside_geofence"`

	Note      string    `gorm:"size:512" json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}