package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type SettingsAdminHandler struct{ svc *services.SettingsAdminService }

func NewSettingsAdminHandler(s *services.SettingsAdminService) *SettingsAdminHandler {
	return &SettingsAdminHandler{svc: s}
}

func (h *SettingsAdminHandler) Get(c *gin.Context) {
	v, err := h.svc.GetAll()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, v)
}

func (h *SettingsAdminHandler) Update(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "body harus berupa object key-value string")
		return
	}
	v, err := h.svc.Update(c.GetUint("user_id"), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, v)
}