package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AdminSeed struct {
	Username string
	Password string
	Name     string
}

type Config struct {
	AppEnv    string
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPassword string
	DBName    string
	JWTSecret string
	JWTExpiry time.Duration
	AdminSeed AdminSeed
}

var C *Config

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "pkl_attendance"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		JWTExpiry:  time.Duration(getEnvInt("JWT_EXPIRY_HOURS", 720)) * time.Hour,
		AdminSeed: AdminSeed{
			Username: getEnv("ADMIN_SEED_USERNAME", "admin"),
			Password: getEnv("ADMIN_SEED_PASSWORD", ""),
			Name:     getEnv("ADMIN_SEED_NAME", "Administrator"),
		},
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET minimal 32 karakter")
	}

	C = cfg
	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}