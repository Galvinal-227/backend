package services

import (
	"encoding/json"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
)

type ActivityLogService struct{ repo *repositories.ActivityLogRepository }

func NewActivityLogService(r *repositories.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{repo: r}
}

func (s *ActivityLogService) Log(userID uint, action, description string, meta map[string]interface{}) {
	metaBytes, _ := json.Marshal(meta)
	_ = s.repo.Create(&models.ActivityLog{
		UserID:      userID,
		Action:      action,
		Description: description,
		Metadata:    string(metaBytes),
	})
}