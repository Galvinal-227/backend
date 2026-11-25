package services

import (
	"errors"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type AdminLeaveService struct {
	leaves      *repositories.LeaveRepository
	attendances *repositories.AttendanceRepository
	holidays    *repositories.HolidayRepository
	settings    *SettingService
	logs        *ActivityLogService
}

func NewAdminLeaveService(
	l *repositories.LeaveRepository,
	a *repositories.AttendanceRepository,
	h *repositories.HolidayRepository,
	s *SettingService,
	log *ActivityLogService,
) *AdminLeaveService {
	return &AdminLeaveService{leaves: l, attendances: a, holidays: h, settings: s, logs: log}
}

func (s *AdminLeaveService) List(userID uint, status string) ([]models.Leave, error) {
	return s.leaves.ListAll(repositories.ListLeaveFilter{UserID: userID, Status: status})
}

type DecideLeaveInput struct {
	AdminID uint
	LeaveID uint
	Approve bool
	Note    string
}

func (s *AdminLeaveService) Decide(in DecideLeaveInput) (*models.Leave, error) {
	l, err := s.leaves.FindByID(in.LeaveID)
	if err != nil {
		return nil, errors.New("pengajuan tidak ditemukan")
	}
	if l.Status != models.LeavePending {
		return nil, errors.New("hanya pengajuan pending yang dapat diputuskan")
	}
	if in.Approve {
		l.Status = models.LeaveApproved
	} else {
		l.Status = models.LeaveRejected
	}
	l.AdminNote = in.Note
	if err := s.leaves.Update(l); err != nil {
		return nil, err
	}

	if in.Approve {
		s.markAttendance(l)
	}

	action := "reject_leave"
	if in.Approve {
		action = "approve_leave"
	}
	s.logs.Log(in.AdminID, action, "Keputusan izin/sakit ID "+itoa(l.ID),
		map[string]interface{}{"leave_id": l.ID, "user_id": l.UserID, "status": l.Status})
	return l, nil
}

func (s *AdminLeaveService) markAttendance(l *models.Leave) {
	settings, err := s.settings.GetAll()
	if err != nil {
		return
	}
	workDays := utils.ParseWorkDays(settings["work_days"])
	holidaySet := map[string]bool{}
	if hols, err := s.holidays.ListInRange(l.StartDate, l.EndDate); err == nil {
		for _, h := range hols {
			holidaySet[h.Date.String()] = true
		}
	}
	statusMap := map[models.LeaveType]models.AttendanceStatus{
		models.LeaveIzin:  models.StatusIzin,
		models.LeaveSakit: models.StatusSakit,
	}
	attStatus := statusMap[l.Type]

	for d := l.StartDate.Time(); !d.After(l.EndDate.Time()); d = d.AddDate(0, 0, 1) {
		if !utils.IsWorkingDay(d, workDays) {
			continue
		}
		do := models.DateOnlyFromTime(d)
		if holidaySet[do.String()] {
			continue
		}
		rec := &models.Attendance{
			UserID: l.UserID,
			Date:   do,
			Status: attStatus,
			Note:   "Auto dari keputusan izin/sakit ID " + itoa(l.ID),
		}
		_, _ = s.attendances.UpsertIfAbsent(rec)
	}
}