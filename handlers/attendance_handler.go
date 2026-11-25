package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/services"
	"pkl-attendance/backend/utils"
)

type AttendanceHandler struct{ svc *services.AttendanceService }

func NewAttendanceHandler(s *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: s}
}

type checkReq struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	var req checkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "lat dan lng wajib diisi")
		return
	}
	uid := c.GetUint("user_id")
	res, err := h.svc.CheckIn(services.CheckInInput{UserID: uid, Lat: req.Lat, Lng: req.Lng})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, res)
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	var req checkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "lat dan lng wajib diisi")
		return
	}
	uid := c.GetUint("user_id")
	res, err := h.svc.CheckOut(services.CheckInInput{UserID: uid, Lat: req.Lat, Lng: req.Lng})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, res)
}

func (h *AttendanceHandler) Today(c *gin.Context) {
	uid := c.GetUint("user_id")
	utils.Success(c, http.StatusOK, h.svc.Today(uid))
}

func (h *AttendanceHandler) History(c *gin.Context) {
	uid := c.GetUint("user_id")
	rows, err := h.svc.History(uid, c.Query("from"), c.Query("to"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, rows)
}

func (h *AttendanceHandler) Recap(c *gin.Context) {
	uid := c.GetUint("user_id")
	month, _ := parseIntQuery(c, "month")
	year, _ := parseIntQuery(c, "year")
	recap, err := h.svc.Recap(uid, month, year)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, recap)
}

func parseIntQuery(c *gin.Context, key string) (int, bool) {
	v := c.Query(key)
	if v == "" {
		return 0, false
	}
	var n int
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		n = n*10 + int(ch-'0')
	}
	return n, true
}