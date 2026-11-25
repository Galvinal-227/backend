package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pkl-attendance/backend/models"
)

type SettingRepository struct{ db *gorm.DB }

func NewSettingRepository(db *gorm.DB) *SettingRepository { return &SettingRepository{db: db} }

func (r *SettingRepository) All() ([]models.Setting, error) {
	var rows []models.Setting
	err := r.db.Find(&rows).Error
	return rows, err
}

func (r *SettingRepository) Upsert(key, value string) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&models.Setting{Key: key, Value: value}).Error
}

func (r *SettingRepository) Get(key string) (string, error) {
	var s models.Setting
	err := r.db.Where("setting_key = ?", key).First(&s).Error
	if err != nil {
		return "", err
	}
	return s.Value, nil
}