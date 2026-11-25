package services

import (
	"errors"
	"time"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
)

type AdminAttendanceService struct {
	attendances *repositories.AttendanceRepository
	logs        *ActivityLogService
}

func NewAdminAttendanceService(a *repositories.AttendanceRepository, l *ActivityLogService) *AdminAttendanceService {
	return &AdminAttendanceService{attendances: a, logs: l}
}

func (s *AdminAttendanceService) List(userID uint, fromStr, toStr, status string) ([]models.Attendance, error) {
	f := repositories.ListAttendanceFilter{UserID: userID, Status: status}
	if fromStr != "" {
		d, err := models.ParseDateOnly(fromStr)
		if err != nil {
			return nil, errors.New("format 'from' harus YYYY-MM-DD")
		}
		f.From = d
	}
	if toStr != "" {
		d, err := models.ParseDateOnly(toStr)
		if err != nil {
			return nil, errors.New("format 'to' harus YYYY-MM-DD")
		}
		f.To = d
	}
	return s.attendances.ListAll(f)
}

type CorrectAttendanceInput struct {
	AdminID      uint
	AttendanceID uint
	Reason       string
	CheckInAt    *time.Time
	CheckOutAt   *time.Time
	Status       string
	Note         string
}

func (s *AdminAttendanceService) Correct(in CorrectAttendanceInput) (*models.Attendance, error) {
	if len(in.Reason) < 3 {
		return nil, errors.New("alasan koreksi wajib diisi (minimal 3 karakter)")
	}
	a, err := s.attendances.FindByID(in.AttendanceID)
	if err != nil {
		return nil, errors.New("absensi tidak ditemukan")
	}
	if in.Status != "" {
		switch models.AttendanceStatus(in.Status) {
		case models.StatusPresent, models.StatusLate, models.StatusAlpha, models.StatusIzin, models.StatusSakit:
			a.Status = models.AttendanceStatus(in.Status)
		default:
			return nil, errors.New("status tidak valid")
		}
	}
	if in.CheckInAt != nil {
		a.CheckInAt = in.CheckInAt
	}
	if in.CheckOutAt != nil {
		a.CheckOutAt = in.CheckOutAt
	}
	if in.Note != "" {
		a.Note = in.Note
	}
	if err := s.attendances.Update(a); err != nil {
		return nil, err
	}
	s.logs.Log(in.AdminID, "correct_attendance",
		"Koreksi absensi tanggal "+a.Date.String(),
		map[string]interface{}{
			"attendance_id": a.ID,
			"reason":        in.Reason,
			"status":        a.Status,
			"user_id":       a.UserID,
		})
	return a, nil
}