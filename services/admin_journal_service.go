package services

import (
	"errors"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
)

type AdminJournalService struct {
	journals *repositories.JournalRepository
	logs     *ActivityLogService
}

func NewAdminJournalService(j *repositories.JournalRepository, l *ActivityLogService) *AdminJournalService {
	return &AdminJournalService{journals: j, logs: l}
}

func (s *AdminJournalService) List(userID uint, fromStr, toStr, status string) ([]models.Journal, error) {
	f := repositories.ListJournalFilter{UserID: userID, Status: status}
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
	return s.journals.ListAll(f)
}

func (s *AdminJournalService) AddNote(adminID, journalID uint, note string) (*models.Journal, error) {
	if len(note) < 2 {
		return nil, errors.New("catatan minimal 2 karakter")
	}
	j, err := s.journals.FindByID(journalID)
	if err != nil {
		return nil, errors.New("jurnal tidak ditemukan")
	}
	j.AdminNote = note
	if err := s.journals.Update(j); err != nil {
		return nil, err
	}
	s.logs.Log(adminID, "journal_note", "Memberi catatan jurnal ID "+itoa(journalID), map[string]interface{}{"journal_id": journalID})
	return j, nil
}

type OverrideJournalInput struct {
	AdminID    uint
	JournalID  uint
	Status     string
	AdminNote  string
	Reason     string
}

func (s *AdminJournalService) Override(in OverrideJournalInput) (*models.Journal, error) {
	if len(in.Reason) < 3 {
		return nil, errors.New("alasan override wajib diisi (minimal 3 karakter)")
	}
	j, err := s.journals.FindByID(in.JournalID)
	if err != nil {
		return nil, errors.New("jurnal tidak ditemukan")
	}
	switch models.JournalStatus(in.Status) {
	case models.JournalPending, models.JournalApproved, models.JournalRejected:
		j.Status = models.JournalStatus(in.Status)
	default:
		return nil, errors.New("status tidak valid")
	}
	if in.AdminNote != "" {
		j.AdminNote = in.AdminNote
	}
	j.OverrideReason = in.Reason
	if err := s.journals.Update(j); err != nil {
		return nil, err
	}
	s.logs.Log(in.AdminID, "override_journal", "Override jurnal ID "+itoa(in.JournalID),
		map[string]interface{}{"journal_id": in.JournalID, "status": in.Status, "reason": in.Reason})
	return j, nil
}

func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}