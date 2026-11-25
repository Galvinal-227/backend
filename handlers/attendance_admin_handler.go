package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type AttendanceAdminHandler struct{ svc *services.AdminAttendanceService }

func NewAttendanceAdminHandler(s *services.AdminAttendanceService) *AttendanceAdminHandler {
	return &AttendanceAdminHandler{svc: s}
}

func (h *AttendanceAdminHandler) List(c *gin.Context) {
	uid := uint(0)
	if q := c.Query("user_id"); q != "" {
		n, err := strconv.ParseUint(q, 10, 32)
		if err == nil {
			uid = uint(n)
		}
	}
	rows, err := h.svc.List(uid, c.Query("from"), c.Query("to"), c.Query("status"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}

type correctAttendanceReq struct {
	Reason     string     `json:"reason" binding:"required"`
	CheckInAt  *time.Time `json:"check_in_at"`
	CheckOutAt *time.Time `json:"check_out_at"`
	Status     string     `json:"status"`
	Note       string     `json:"note"`
}

func (h *AttendanceAdminHandler) Correct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req correctAttendanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}
	a, err := h.svc.Correct(services.CorrectAttendanceInput{
		AdminID: c.GetUint("user_id"), AttendanceID: uint(id),
		Reason: req.Reason, CheckInAt: req.CheckInAt, CheckOutAt: req.CheckOutAt,
		Status: req.Status, Note: req.Note,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, a)
}