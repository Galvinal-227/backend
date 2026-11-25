package repositories

import (
	"errors"

	"gorm.io/gorm"

	"pkl-attendance/backend/models"
)

var ErrNotFound = errors.New("not found")

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) Create(u *models.User) error { return r.db.Create(u).Error }
func (r *UserRepository) Update(u *models.User) error { return r.db.Save(u).Error }

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ListStudents() ([]models.User, error) {
	var rows []models.User
	err := r.db.Where("role = ?", models.RoleSiswa).Preload("Student").Order("id ASC").Find(&rows).Error
	return rows, err
}

// DeleteWithAssociations menghapus user dan SEMUA data terkait (absensi, jurnal, izin, log, profil siswa)
func (r *UserRepository) DeleteWithAssociations(userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Hapus data anak (relasi)
		if err := tx.Where("user_id = ?", userID).Delete(&models.Attendance{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Journal{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Leave{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.ActivityLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&models.Student{}).Error; err != nil {
			return err
		}
		
		// 2. Hapus data induk (User) paling akhir
		return tx.Delete(&models.User{}, userID).Error
	})
}