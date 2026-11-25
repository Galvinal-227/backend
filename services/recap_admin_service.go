package services

import (
	"errors"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type RecapAdminService struct {
	attendance *AttendanceService
	users      *repositories.UserRepository
}

func NewRecapAdminService(a *AttendanceService, u *repositories.UserRepository) *RecapAdminService {
	return &RecapAdminService{attendance: a, users: u}
}

type StudentRecap struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	WorkingDays int    `json:"working_days"`
	Present     int    `json:"present"`
	Late        int    `json:"late"`
	Izin        int    `json:"izin"`
	Sakit       int    `json:"sakit"`
	Alpha       int    `json:"alpha"`
}

type RecapAllResponse struct {
	Year  int            `json:"year"`
	Month int            `json:"month"`
	Rows  []StudentRecap `json:"rows"`
}

func (s *RecapAdminService) RecapAll(month, year int) (*RecapAllResponse, error) {
	if month < 0 || month > 12 {
		return nil, errors.New("month tidak valid")
	}
	students, err := s.users.ListStudents()
	if err != nil {
		return nil, err
	}
	now := utils.JakartaNow()
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}
	out := &RecapAllResponse{Year: year, Month: month, Rows: []StudentRecap{}}
	for _, u := range students {
		recap, err := s.attendance.Recap(u.ID, month, year)
		if err != nil {
			continue
		}
		out.Rows = append(out.Rows, StudentRecap{
			UserID:      u.ID,
			Username:    u.Username,
			Name:        u.Name,
			WorkingDays: recap.WorkingDays,
			Present:     recap.Present,
			Late:        recap.Late,
			Izin:        recap.Izin,
			Sakit:       recap.Sakit,
			Alpha:       recap.Alpha,
		})
	}
	return out, nil
}

var _ = models.StatusPresent