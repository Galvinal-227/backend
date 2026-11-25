package repositories

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pkl-attendance/backend/models"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Revoke(jti string, expiresAt time.Time, reason string) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.TokenBlacklist{
		JTI:       jti,
		ExpiresAt: expiresAt,
		Reason:    reason,
	}).Error
}

func (r *TokenRepository) IsRevoked(jti string) bool {
	var count int64
	r.db.Model(&models.TokenBlacklist{}).Where("jti = ?", jti).Count(&count)
	return count > 0
}

func (r *TokenRepository) CleanupExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.TokenBlacklist{}).Error
}