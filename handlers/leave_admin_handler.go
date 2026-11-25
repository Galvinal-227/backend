package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type LeaveAdminHandler struct{ svc *services.AdminLeaveService }

func NewLeaveAdminHandler(s *services.AdminLeaveService) *LeaveAdminHandler {
	return &LeaveAdminHandler{svc: s}
}

func (h *LeaveAdminHandler) List(c *gin.Context) {
	uid := uint(0)
	if q := c.Query("user_id"); q != "" {
		n, _ := strconv.ParseUint(q, 10, 32)
		uid = uint(n)
	}
	rows, err := h.svc.List(uid, c.Query("status"))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}

type decideLeaveReq struct {
	Approve bool   `json:"approve"`
	Note    string `json:"note"`
}

func (h *LeaveAdminHandler) Decide(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req decideLeaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	l, err := h.svc.Decide(services.DecideLeaveInput{
		AdminID: c.GetUint("user_id"), LeaveID: uint(id), Approve: req.Approve, Note: req.Note,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, l)
}