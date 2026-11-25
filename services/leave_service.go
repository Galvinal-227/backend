package services

import (
	"errors"
	"strings"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
)

type LeaveService struct {
	leaves *repositories.LeaveRepository
}

func NewLeaveService(l *repositories.LeaveRepository) *LeaveService {
	return &LeaveService{leaves: l}
}

type CreateLeaveInput struct {
	UserID    uint
	Type      string
	StartDate string
	EndDate   string
	Reason    string
}

func (s *LeaveService) Create(in CreateLeaveInput) (*models.Leave, error) {
	t := models.LeaveType(in.Type)
	if t != models.LeaveIzin && t != models.LeaveSakit {
		return nil, errors.New("tipe harus 'izin' atau 'sakit'")
	}
	start, err := models.ParseDateOnly(in.StartDate)
	if err != nil {
		return nil, errors.New("start_date harus YYYY-MM-DD")
	}
	end, err := models.ParseDateOnly(in.EndDate)
	if err != nil {
		return nil, errors.New("end_date harus YYYY-MM-DD")
	}
	if end.Time().Before(start.Time()) {
		return nil, errors.New("end_date tidak boleh sebelum start_date")
	}
	if strings.TrimSpace(in.Reason) == "" {
		return nil, errors.New("alasan wajib diisi")
	}
	overlap, _ := s.leaves.OverlapExists(in.UserID, start, end)
	if overlap {
		return nil, errors.New("sudah ada pengajuan izin/sakit pada rentang tanggal tersebut")
	}
	l := &models.Leave{
		UserID:    in.UserID,
		Type:      t,
		StartDate: start,
		EndDate:   end,
		Reason:    strings.TrimSpace(in.Reason),
		Status:    models.LeavePending,
	}
	if err := s.leaves.Create(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *LeaveService) ListMine(userID uint) ([]models.Leave, error) {
	return s.leaves.ListByUser(userID)
}