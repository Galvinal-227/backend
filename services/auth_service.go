package services

import (
"errors"
"time"

"pkl-attendance/backend/config"
"pkl-attendance/backend/models"
"pkl-attendance/backend/repositories"
"pkl-attendance/backend/utils"
)

type AuthService struct {
cfg    *config.Config
users  *repositories.UserRepository
tokens *repositories.TokenRepository
logs   *ActivityLogService
}

func NewAuthService(cfg *config.Config, users *repositories.UserRepository, tokens *repositories.TokenRepository, logs *ActivityLogService) *AuthService {
return &AuthService{cfg: cfg, users: users, tokens: tokens, logs: logs}
}

type LoginResult struct {
Token     string       `json:"token"`
ExpiresAt int64        `json:"expires_at"`
User      *models.User `json:"user"`
}

func (s *AuthService) Login(username, password string) (*LoginResult, error) {
user, err := s.users.FindByUsername(username)
if err != nil {
return nil, errors.New("invalid credentials")
}
if !user.IsActive {
return nil, errors.New("account is disabled")
}
if !utils.CheckPassword(user.Password, password) {
return nil, errors.New("invalid credentials")
}
token, _, exp, err := utils.GenerateToken(
s.cfg.JWTSecret, user.ID, user.Username, string(user.Role), s.cfg.JWTExpiry,
)
if err != nil {
return nil, err
}
return &LoginResult{
Token:     token,
ExpiresAt: exp.Unix(),
User:      user,
}, nil
}

func (s *AuthService) Logout(jti string, exp int64) error {
expiresAt := time.Unix(exp, 0)
return s.tokens.Revoke(jti, expiresAt, "logout")
}

func (s *AuthService) Me(userID uint) (*models.User, error) {
return s.users.FindByID(userID)
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
if len(newPassword) < 6 {
return errors.New("password baru minimal 6 karakter")
}
if oldPassword == newPassword {
return errors.New("password baru harus berbeda dari password lama")
}
user, err := s.users.FindByID(userID)
if err != nil {
return errors.New("user tidak ditemukan")
}
if !utils.CheckPassword(user.Password, oldPassword) {
return errors.New("password lama tidak sesuai")
}
hash, err := utils.HashPassword(newPassword)
if err != nil {
return err
}
user.Password = hash
if err := s.users.Update(user); err != nil {
return err
}
if s.logs != nil {
s.logs.Log(userID, "change_password", "User mengganti password sendiri",
map[string]interface{}{"user_id": userID})
}
return nil
}