package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type JournalAdminHandler struct{ svc *services.AdminJournalService }

func NewJournalAdminHandler(s *services.AdminJournalService) *JournalAdminHandler {
	return &JournalAdminHandler{svc: s}
}

func (h *JournalAdminHandler) List(c *gin.Context) {
	uid := uint(0)
	if q := c.Query("user_id"); q != "" {
		n, _ := strconv.ParseUint(q, 10, 32)
		uid = uint(n)
	}
	rows, err := h.svc.List(uid, c.Query("from"), c.Query("to"), c.Query("status"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}

type noteReq struct {
	Note string `json:"note" binding:"required"`
}

func (h *JournalAdminHandler) AddNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req noteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "catatan wajib diisi")
		return
	}
	j, err := h.svc.AddNote(c.GetUint("user_id"), uint(id), req.Note)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, j)
}

type overrideJournalReq struct {
	Status    string `json:"status" binding:"required"`
	AdminNote string `json:"admin_note"`
	Reason    string `json:"reason" binding:"required"`
}

func (h *JournalAdminHandler) Override(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req overrideJournalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}
	j, err := h.svc.Override(services.OverrideJournalInput{
		AdminID: c.GetUint("user_id"), JournalID: uint(id),
		Status: req.Status, AdminNote: req.AdminNote, Reason: req.Reason,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, j)
}