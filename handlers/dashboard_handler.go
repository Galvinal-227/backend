package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type DashboardHandler struct {
	attendance *services.AttendanceService
	journals   *services.JournalService
	leaves     *services.LeaveService
	settings   *services.SettingService
	users      *repositories.UserRepository
	students   *repositories.StudentRepository
}

func NewDashboardHandler(
	a *services.AttendanceService,
	j *services.JournalService,
	l *services.LeaveService,
	s *services.SettingService,
	u *repositories.UserRepository,
	st *repositories.StudentRepository,
) *DashboardHandler {
	return &DashboardHandler{a, j, l, s, u, st}
}

func (h *DashboardHandler) SiswaDashboard(c *gin.Context) {
	uid := c.GetUint("user_id")
	student, _ := h.students.FindByUserID(uid)
	today := h.attendance.Today(uid)
	recap, _ := h.attendance.Recap(uid, 0, 0)
	settings, _ := h.settings.GetAll()

	utils.Success(c, http.StatusOK, gin.H{
		"student":    student,
		"today":      today,
		"recap":      recap,
		"settings":   settings,
	})
}