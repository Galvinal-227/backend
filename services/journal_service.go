package services

import (
	"errors"
	"strings"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type JournalService struct {
	journals    *repositories.JournalRepository
	attendances *repositories.AttendanceRepository
}

func NewJournalService(j *repositories.JournalRepository, a *repositories.AttendanceRepository) *JournalService {
	return &JournalService{journals: j, attendances: a}
}

func (s *JournalService) CreateForToday(userID uint, content string) (*models.Journal, error) {
	content = strings.TrimSpace(content)
	if len(content) < 10 {
		return nil, errors.New("isi jurnal minimal 10 karakter")
	}
	today := models.DateOnlyFromTime(utils.JakartaNow())
	if _, err := s.journals.FindByUserAndDate(userID, today); err == nil {
		return nil, errors.New("jurnal hari ini sudah ada, gunakan edit")
	}
	j := &models.Journal{
		UserID:  userID,
		Date:    today,
		Content: content,
		Status:  models.JournalPending,
	}
	if err := s.journals.Create(j); err != nil {
		return nil, err
	}
	return j, nil
}

func (s *JournalService) UpdateMine(userID, journalID uint, content string) (*models.Journal, error) {
	content = strings.TrimSpace(content)
	if len(content) < 10 {
		return nil, errors.New("isi jurnal minimal 10 karakter")
	}
	j, err := s.journals.FindByID(journalID)
	if err != nil {
		return nil, errors.New("jurnal tidak ditemukan")
	}
	if j.UserID != userID {
		return nil, errors.New("anda tidak berhak mengubah jurnal ini")
	}
	today := models.DateOnlyFromTime(utils.JakartaNow())
	if j.Date.String() != today.String() {
		return nil, errors.New("hanya jurnal hari ini yang dapat diubah")
	}
	if att, _ := s.attendances.FindByUserAndDate(userID, today); att != nil && att.CheckOutAt != nil {
		return nil, errors.New("jurnal terkunci setelah absen pulang")
	}
	j.Content = content
	if err := s.journals.Update(j); err != nil {
		return nil, err
	}
	return j, nil
}

func (s *JournalService) ListMine(userID uint, fromStr, toStr string) ([]models.Journal, error) {
	now := utils.JakartaNow()
	from := models.DateOnlyFromTime(utils.FirstOfMonth(now))
	to := models.DateOnlyFromTime(utils.LastOfMonth(now))
	if fromStr != "" {
		d, err := models.ParseDateOnly(fromStr)
		if err != nil {
			return nil, errors.New("format 'from' harus YYYY-MM-DD")
		}
		from = d
	}
	if toStr != "" {
		d, err := models.ParseDateOnly(toStr)
		if err != nil {
			return nil, errors.New("format 'to' harus YYYY-MM-DD")
		}
		to = d
	}
	return s.journals.ListByUserAndRange(userID, from, to)
}