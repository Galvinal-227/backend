package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type StudentAdminHandler struct{ svc *services.StudentManagementService }

func NewStudentAdminHandler(s *services.StudentManagementService) *StudentAdminHandler {
	return &StudentAdminHandler{svc: s}
}

type createStudentReq struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Name         string `json:"name" binding:"required"`
	NIS          string `json:"nis"`
	Class        string `json:"class"`
	Major        string `json:"major"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	PKLStartDate string `json:"pkl_start_date"`
	PKLEndDate   string `json:"pkl_end_date"`
}

func (h *StudentAdminHandler) Create(c *gin.Context) {
	var req createStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak lengkap: "+err.Error())
		return
	}
	u, err := h.svc.Create(c.GetUint("user_id"), services.CreateStudentInput{
		Username: req.Username, Password: req.Password, Name: req.Name,
		NIS: req.NIS, Class: req.Class, Major: req.Major, Phone: req.Phone, Address: req.Address,
		PKLStartDate: req.PKLStartDate, PKLEndDate: req.PKLEndDate,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, u)
}

type updateStudentReq struct {
	Name         string `json:"name"`
	NIS          string `json:"nis"`
	Class        string `json:"class"`
	Major        string `json:"major"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	PKLStartDate string `json:"pkl_start_date"`
	PKLEndDate   string `json:"pkl_end_date"`
}

func (h *StudentAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req updateStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "input tidak valid")
		return
	}
	u, err := h.svc.Update(c.GetUint("user_id"), services.UpdateStudentInput{
		UserID: uint(id), Name: req.Name, NIS: req.NIS, Class: req.Class, Major: req.Major,
		Phone: req.Phone, Address: req.Address, PKLStartDate: req.PKLStartDate, PKLEndDate: req.PKLEndDate,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, u)
}

type resetPasswordReq struct {
	Password string `json:"password" binding:"required"`
}

func (h *StudentAdminHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "password wajib diisi")
		return
	}
	if err := h.svc.ResetPassword(c.GetUint("user_id"), uint(id), req.Password); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "password direset"})
}

type toggleActiveReq struct {
	Active bool `json:"active"`
}

func (h *StudentAdminHandler) ToggleActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	var req toggleActiveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "aktif wajib diisi")
		return
	}
	if err := h.svc.ToggleActive(c.GetUint("user_id"), uint(id), req.Active); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "status diperbarui"})
}

func (h *StudentAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}
	if err := h.svc.Delete(c.GetUint("user_id"), uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "siswa berhasil dihapus"})
}

func (h *StudentAdminHandler) List(c *gin.Context) {
	rows, err := h.svc.List()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}