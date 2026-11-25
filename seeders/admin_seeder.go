package seeders

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/models"
	"pkl-attendance/backend/utils"
)

var defaultSettings = map[string]string{
	"work_days":         "1,2,3,4,5",
	"check_in_time":     "07:30",
	"tolerance_minutes": "15",
	"geofence_lat":      "-6.200000",
	"geofence_lng":      "106.816666",
	"geofence_radius":   "100",
	"pkl_location_name": "Lokasi PKL",
}

func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	if cfg.AdminSeed.Password == "" {
		return errors.New("ADMIN_SEED_PASSWORD wajib diisi untuk seeding admin")
	}
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		log.Println("admin sudah ada, seed admin dilewati")
	} else {
		hash, err := utils.HashPassword(cfg.AdminSeed.Password)
		if err != nil {
			return err
		}
		admin := &models.User{
			Username: cfg.AdminSeed.Username,
			Password: hash,
			Name:     cfg.AdminSeed.Name,
			Role:     models.RoleAdmin,
			IsActive: true,
		}
		if err := db.Create(admin).Error; err != nil {
			return err
		}
		log.Printf("admin berhasil dibuat: %s\n", admin.Username)
	}

	for k, v := range defaultSettings {
		var existing models.Setting
		err := db.Where("setting_key = ?", k).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&models.Setting{Key: k, Value: v}).Error; err != nil {
				return err
			}
			log.Printf("setting default dibuat: %s=%s\n", k, v)
		}
	}
	return nil
}