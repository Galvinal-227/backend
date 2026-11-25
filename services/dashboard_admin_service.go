package services

import (
	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type DashboardAdminService struct {
	attendances *repositories.AttendanceRepository
	journals    *repositories.JournalRepository
	leaves      *repositories.LeaveRepository
	users       *repositories.UserRepository
}

func NewDashboardAdminService(
	a *repositories.AttendanceRepository,
	j *repositories.JournalRepository,
	l *repositories.LeaveRepository,
	u *repositories.UserRepository,
) *DashboardAdminService {
	return &DashboardAdminService{a, j, l, u}
}

func (s *DashboardAdminService) Summary() (map[string]interface{}, error) {
	now := utils.JakartaNow()
	today := models.DateOnlyFromTime(now)
	from := models.DateOnlyFromTime(utils.FirstOfMonth(now))
	to := models.DateOnlyFromTime(utils.LastOfMonth(now))

	attendancesToday, _ := s.attendances.ListAll(repositories.ListAttendanceFilter{From: today, To: today})
	journalsToday, _ := s.journals.ListAll(repositories.ListJournalFilter{From: today, To: today})
	pendingLeaves, _ := s.leaves.ListAll(repositories.ListLeaveFilter{Status: "pending"})
	monthAttendances, _ := s.attendances.ListAll(repositories.ListAttendanceFilter{From: from, To: to})
	students, _ := s.users.ListStudents()

	in := 0
	out := 0
	for _, r := range attendancesToday {
		if r.CheckInAt != nil {
			in++
		}
		if r.CheckOutAt != nil {
			out++
		}
	}

	return map[string]interface{}{
		"total_students":        len(students),
		"check_in_today":        in,
		"check_out_today":       out,
		"journals_today":        len(journalsToday),
		"pending_leaves":        len(pendingLeaves),
		"attendance_this_month": len(monthAttendances),
	}, nil
}