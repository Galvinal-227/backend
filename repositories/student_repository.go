package repositories

import (
	"errors"

	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

type StudentRepository struct{ db *gorm.DB }

func NewStudentRepository(db *gorm.DB) *StudentRepository { return &StudentRepository{db: db} }

func (r *StudentRepository) FindByUserID(userID uint) (*models.Student, error) {
	var s models.Student
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *StudentRepository) FindByID(id uint) (*models.Student, error) {
	var s models.Student
	err := r.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}

func (r *StudentRepository) Create(s *models.Student) error { return r.db.Create(s).Error }
func (r *StudentRepository) Update(s *models.Student) error { return r.db.Save(s).Error }