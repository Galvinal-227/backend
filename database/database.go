package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	level := logger.Warn
	if cfg.AppEnv == "development" && os.Getenv("DB_DEBUG") == "true" {
		level = logger.Info
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(level),
		NowFunc: func() time.Time {
			loc, e := time.LoadLocation("Asia/Jakarta")
			if e != nil {
				return time.Now()
			}
			return time.Now().In(loc)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	DB = db
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("running automigrate...")
	return db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Attendance{},
		&models.Journal{},
		&models.Leave{},
		&models.Holiday{},
		&models.Setting{},
		&models.ActivityLog{},
		&models.TokenBlacklist{},
	)
}