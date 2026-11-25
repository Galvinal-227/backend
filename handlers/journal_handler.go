package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type JournalHandler struct{ svc *services.JournalService }

func NewJournalHandler(s *services.JournalService) *JournalHandler {
	return &JournalHandler{svc: s}
}

type journalReq struct {
	Content string `json:"content" binding:"required"`
}

func (h *JournalHandler) Create(c *gin.Context) {
	var req journalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "content wajib diisi")
		return
	}
	j, err := h.svc.CreateForToday(c.GetUint("user_id"), req.Content)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, j)
}

func (h *JournalHandler) UpdateMine(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req journalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "content wajib diisi")
		return
	}
	j, err := h.svc.UpdateMine(c.GetUint("user_id"), uint(id), req.Content)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, j)
}

func (h *JournalHandler) ListMine(c *gin.Context) {
	rows, err := h.svc.ListMine(c.GetUint("user_id"), c.Query("from"), c.Query("to"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}