package repositories

import (
	"errors"

	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

type JournalRepository struct{ db *gorm.DB }

func NewJournalRepository(db *gorm.DB) *JournalRepository { return &JournalRepository{db: db} }

func (r *JournalRepository) Create(j *models.Journal) error { return r.db.Create(j).Error }
func (r *JournalRepository) Update(j *models.Journal) error { return r.db.Save(j).Error }

func (r *JournalRepository) FindByID(id uint) (*models.Journal, error) {
	var j models.Journal
	err := r.db.First(&j, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &j, err
}

func (r *JournalRepository) FindByUserAndDate(userID uint, d models.DateOnly) (*models.Journal, error) {
	var j models.Journal
	err := r.db.Where("user_id = ? AND date = ?", userID, d).First(&j).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &j, err
}

func (r *JournalRepository) ExistsByUserAndDate(userID uint, d models.DateOnly) (bool, error) {
	var count int64
	err := r.db.Model(&models.Journal{}).Where("user_id = ? AND date = ?", userID, d).Count(&count).Error
	return count > 0, err
}

func (r *JournalRepository) ListByUserAndRange(userID uint, from, to models.DateOnly) ([]models.Journal, error) {
	var rows []models.Journal
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).
		Order("date DESC").Find(&rows).Error
	return rows, err
}

type ListJournalFilter struct {
	UserID uint
	From   models.DateOnly
	To     models.DateOnly
	Status string
}

func (r *JournalRepository) ListAll(f ListJournalFilter) ([]models.Journal, error) {
	q := r.db.Model(&models.Journal{})
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if !f.From.Time().IsZero() {
		q = q.Where("date >= ?", f.From)
	}
	if !f.To.Time().IsZero() {
		q = q.Where("date <= ?", f.To)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	var rows []models.Journal
	err := q.Order("date DESC, user_id ASC").Find(&rows).Error
	return rows, err
}