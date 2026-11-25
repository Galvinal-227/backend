package repositories

import (
	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

type ActivityLogRepository struct{ db *gorm.DB }

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(a *models.ActivityLog) error { return r.db.Create(a).Error }