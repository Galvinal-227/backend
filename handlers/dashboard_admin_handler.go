package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type DashboardAdminHandler struct{ svc *services.DashboardAdminService }

func NewDashboardAdminHandler(s *services.DashboardAdminService) *DashboardAdminHandler {
	return &DashboardAdminHandler{svc: s}
}

func (h *DashboardAdminHandler) Summary(c *gin.Context) {
	v, err := h.svc.Summary()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, v)
}