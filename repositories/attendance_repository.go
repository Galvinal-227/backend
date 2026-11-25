package repositories

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pkl-attendance/backend/models"
)

type AttendanceRepository struct{ db *gorm.DB }

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository { return &AttendanceRepository{db: db} }

func (r *AttendanceRepository) Create(a *models.Attendance) error { return r.db.Create(a).Error }
func (r *AttendanceRepository) Update(a *models.Attendance) error { return r.db.Save(a).Error }

func (r *AttendanceRepository) FindByID(id uint) (*models.Attendance, error) {
	var a models.Attendance
	err := r.db.First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &a, err
}

func (r *AttendanceRepository) FindByUserAndDate(userID uint, d models.DateOnly) (*models.Attendance, error) {
	var a models.Attendance
	err := r.db.Where("user_id = ? AND date = ?", userID, d).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &a, err
}

func (r *AttendanceRepository) ListByUserAndRange(userID uint, from, to models.DateOnly) ([]models.Attendance, error) {
	var rows []models.Attendance
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).
		Order("date DESC").Find(&rows).Error
	return rows, err
}

type ListAttendanceFilter struct {
	UserID uint
	From   models.DateOnly
	To     models.DateOnly
	Status string
}

func (r *AttendanceRepository) ListAll(f ListAttendanceFilter) ([]models.Attendance, error) {
	q := r.db.Model(&models.Attendance{})
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
	var rows []models.Attendance
	err := q.Order("date DESC, user_id ASC").Find(&rows).Error
	return rows, err
}

// UpsertIfAbsent membuat record attendance hanya jika belum ada pada tanggal tersebut.
func (r *AttendanceRepository) UpsertIfAbsent(a *models.Attendance) (created bool, err error) {
	var existing models.Attendance
	err = r.db.Where("user_id = ? AND date = ?", a.UserID, a.Date).First(&existing).Error
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	err = r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(a).Error
	return err == nil, err
}