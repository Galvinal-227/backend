package services

import (
	"errors"
	"strings"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
)

type HolidayService struct {
	holidays *repositories.HolidayRepository
	logs     *ActivityLogService
}

func NewHolidayService(h *repositories.HolidayRepository, l *ActivityLogService) *HolidayService {
	return &HolidayService{holidays: h, logs: l}
}

func (s *HolidayService) List() ([]models.Holiday, error) { return s.holidays.ListAll() }

func (s *HolidayService) Create(adminID uint, date, name string) (*models.Holiday, error) {
	d, err := models.ParseDateOnly(date)
	if err != nil {
		return nil, errors.New("tanggal harus YYYY-MM-DD")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("nama libur wajib diisi")
	}
	if exists, _ := s.holidays.ExistsOnDate(d); exists {
		return nil, errors.New("tanggal sudah terdaftar sebagai libur")
	}
	h := &models.Holiday{Date: d, Name: name}
	if err := s.holidays.Create(h); err != nil {
		return nil, err
	}
	s.logs.Log(adminID, "create_holiday", "Menambah libur "+d.String()+" "+name,
		map[string]interface{}{"holiday_id": h.ID})
	return h, nil
}

func (s *HolidayService) Update(adminID, id uint, date, name string) (*models.Holiday, error) {
	h, err := s.holidays.FindByID(id)
	if err != nil {
		return nil, errors.New("libur tidak ditemukan")
	}
	if date != "" {
		d, err := models.ParseDateOnly(date)
		if err != nil {
			return nil, errors.New("tanggal harus YYYY-MM-DD")
		}
		if d.String() != h.Date.String() {
			if exists, _ := s.holidays.ExistsOnDate(d); exists {
				return nil, errors.New("tanggal sudah terdaftar sebagai libur lain")
			}
		}
		h.Date = d
	}
	if name = strings.TrimSpace(name); name != "" {
		h.Name = name
	}
	if err := s.holidays.Update(h); err != nil {
		return nil, err
	}
	s.logs.Log(adminID, "update_holiday", "Mengubah libur ID "+itoa(id),
		map[string]interface{}{"holiday_id": id})
	return h, nil
}

func (s *HolidayService) Delete(adminID, id uint) error {
	if _, err := s.holidays.FindByID(id); err != nil {
		return errors.New("libur tidak ditemukan")
	}
	if err := s.holidays.Delete(id); err != nil {
		return err
	}
	s.logs.Log(adminID, "delete_holiday", "Menghapus libur ID "+itoa(id),
		map[string]interface{}{"holiday_id": id})
	return nil
}