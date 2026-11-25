package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type RecapAdminHandler struct{ svc *services.RecapAdminService }

func NewRecapAdminHandler(s *services.RecapAdminService) *RecapAdminHandler {
	return &RecapAdminHandler{svc: s}
}

func (h *RecapAdminHandler) RecapAll(c *gin.Context) {
	month, _ := parseIntQuery(c, "month")
	year, _ := parseIntQuery(c, "year")
	res, err := h.svc.RecapAll(month, year)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, res)
}