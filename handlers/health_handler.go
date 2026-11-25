package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/database"
	"pkl-attendance/backend/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "ok"
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}
	} else {
		dbStatus = "uninitialized"
	}
	utils.Success(c, http.StatusOK, gin.H{
		"status": "ok",
		"db":     dbStatus,
	})
}