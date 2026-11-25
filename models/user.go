package models

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleSiswa Role = "siswa"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Role      Role      `gorm:"size:16;not null;index" json:"role"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Student *Student `gorm:"foreignKey:UserID" json:"student,omitempty"`
}