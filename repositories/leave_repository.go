package repositories

import (
	"errors"

	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

type LeaveRepository struct{ db *gorm.DB }

func NewLeaveRepository(db *gorm.DB) *LeaveRepository { return &LeaveRepository{db: db} }

func (r *LeaveRepository) Create(l *models.Leave) error { return r.db.Create(l).Error }
func (r *LeaveRepository) Update(l *models.Leave) error { return r.db.Save(l).Error }

func (r *LeaveRepository) FindByID(id uint) (*models.Leave, error) {
	var l models.Leave
	err := r.db.First(&l, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &l, err
}

func (r *LeaveRepository) ListByUser(userID uint) ([]models.Leave, error) {
	var rows []models.Leave
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *LeaveRepository) OverlapExists(userID uint, start, end models.DateOnly) (bool, error) {
	var count int64
	err := r.db.Model(&models.Leave{}).
		Where("user_id = ? AND status IN ?", userID, []models.LeaveStatus{models.LeavePending, models.LeaveApproved}).
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?)",
			end, start, start, start).
		Count(&count).Error
	return count > 0, err
}

func (r *LeaveRepository) FindApprovedByUserAndDate(userID uint, d models.DateOnly) (*models.Leave, error) {
	var l models.Leave
	err := r.db.Where("user_id = ? AND status = ? AND start_date <= ? AND end_date >= ?",
		userID, models.LeaveApproved, d, d).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &l, err
}

type ListLeaveFilter struct {
	UserID uint
	Status string
}

func (r *LeaveRepository) ListAll(f ListLeaveFilter) ([]models.Leave, error) {
	q := r.db.Model(&models.Leave{})
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	var rows []models.Leave
	err := q.Order("created_at DESC").Find(&rows).Error
	return rows, err
}