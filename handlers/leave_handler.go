package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type LeaveHandler struct{ svc *services.LeaveService }

func NewLeaveHandler(s *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{svc: s}
}

type leaveReq struct {
	Type      string `json:"type" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
}

func (h *LeaveHandler) Create(c *gin.Context) {
	var req leaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak lengkap")
		return
	}
	l, err := h.svc.Create(services.CreateLeaveInput{
		UserID:    c.GetUint("user_id"),
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Reason:    req.Reason,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, l)
}

func (h *LeaveHandler) ListMine(c *gin.Context) {
	rows, err := h.svc.ListMine(c.GetUint("user_id"))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}