package repositories

import (
	"errors"

	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

type HolidayRepository struct{ db *gorm.DB }

func NewHolidayRepository(db *gorm.DB) *HolidayRepository { return &HolidayRepository{db: db} }

func (r *HolidayRepository) Create(h *models.Holiday) error { return r.db.Create(h).Error }
func (r *HolidayRepository) Update(h *models.Holiday) error { return r.db.Save(h).Error }
func (r *HolidayRepository) Delete(id uint) error           { return r.db.Delete(&models.Holiday{}, id).Error }

func (r *HolidayRepository) FindByID(id uint) (*models.Holiday, error) {
	var h models.Holiday
	err := r.db.First(&h, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &h, err
}

func (r *HolidayRepository) ExistsOnDate(d models.DateOnly) (bool, error) {
	var count int64
	err := r.db.Model(&models.Holiday{}).Where("date = ?", d).Count(&count).Error
	return count > 0, err
}

func (r *HolidayRepository) ListInRange(from, to models.DateOnly) ([]models.Holiday, error) {
	var rows []models.Holiday
	err := r.db.Where("date BETWEEN ? AND ?", from, to).Order("date ASC").Find(&rows).Error
	return rows, err
}

func (r *HolidayRepository) ListAll() ([]models.Holiday, error) {
	var rows []models.Holiday
	err := r.db.Order("date ASC").Find(&rows).Error
	return rows, err
}