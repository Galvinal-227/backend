package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type HolidayAdminHandler struct{ svc *services.HolidayService }

func NewHolidayAdminHandler(s *services.HolidayService) *HolidayAdminHandler {
	return &HolidayAdminHandler{svc: s}
}

func (h *HolidayAdminHandler) List(c *gin.Context) {
	rows, err := h.svc.List()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}

type holidayReq struct {
	Date string `json:"date" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *HolidayAdminHandler) Create(c *gin.Context) {
	var req holidayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak lengkap")
		return
	}
	hl, err := h.svc.Create(c.GetUint("user_id"), req.Date, req.Name)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, hl)
}

func (h *HolidayAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req holidayReq
	_ = c.ShouldBindJSON(&req)
	hl, err := h.svc.Update(c.GetUint("user_id"), uint(id), req.Date, req.Name)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, hl)
}

func (h *HolidayAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	if err := h.svc.Delete(c.GetUint("user_id"), uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "libur dihapus"})
}