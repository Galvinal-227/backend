package services

import (
	"errors"
	"strconv"
	"time"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type AttendanceService struct {
	attendances *repositories.AttendanceRepository
	students    *repositories.StudentRepository
	journals    *repositories.JournalRepository
	leaves      *repositories.LeaveRepository
	holidays    *repositories.HolidayRepository
	settings    *SettingService
	geocoder    *GeocodingService
}

func NewAttendanceService(
	attendances *repositories.AttendanceRepository,
	students *repositories.StudentRepository,
	journals *repositories.JournalRepository,
	leaves *repositories.LeaveRepository,
	holidays *repositories.HolidayRepository,
	settings *SettingService,
	geocoder *GeocodingService,
) *AttendanceService {
	return &AttendanceService{attendances, students, journals, leaves, holidays, settings, geocoder}
}

type CheckInInput struct {
	UserID uint
	Lat    float64
	Lng    float64
}

type CheckResult struct {
	Attendance *models.Attendance `json:"attendance"`
	Distance   float64            `json:"distance_meters"`
	Address    string             `json:"address"`
	Message    string             `json:"message"`
}

func (s *AttendanceService) CheckIn(in CheckInInput) (*CheckResult, error) {
	now := utils.JakartaNow()
	today := models.DateOnlyFromTime(now)

	student, err := s.students.FindByUserID(in.UserID)
	if err != nil {
		return nil, errors.New("profil siswa tidak ditemukan")
	}
	if student.PKLStartDate != nil && today.Time().Before(student.PKLStartDate.Time()) {
		return nil, errors.New("belum memasuki periode PKL")
	}
	if student.PKLEndDate != nil && today.Time().After(student.PKLEndDate.Time()) {
		return nil, errors.New("periode PKL telah berakhir")
	}

	settings, err := s.settings.GetAll()
	if err != nil {
		return nil, err
	}
	workDays := utils.ParseWorkDays(settings["work_days"])
	if !utils.IsWorkingDay(now, workDays) {
		return nil, errors.New("hari ini bukan hari kerja")
	}
	if isHoliday, _ := s.holidays.ExistsOnDate(today); isHoliday {
		return nil, errors.New("hari ini libur nasional/khusus")
	}

	if existing, err := s.attendances.FindByUserAndDate(in.UserID, today); err == nil && existing.CheckInAt != nil {
		return nil, errors.New("anda sudah absen masuk hari ini")
	}

	geoLat, _ := strconv.ParseFloat(settings["geofence_lat"], 64)
	geoLng, _ := strconv.ParseFloat(settings["geofence_lng"], 64)
	geoRad, _ := strconv.ParseFloat(settings["geofence_radius"], 64)
	distance := utils.DistanceMeters(in.Lat, in.Lng, geoLat, geoLng)
	outside := distance > geoRad

	statusStr := utils.DetermineAttendanceStatus(now, settings["check_in_time"], settings["tolerance_minutes"])
	address := s.geocoder.ReverseGeocode(in.Lat, in.Lng)

	att := &models.Attendance{
		UserID:                  in.UserID,
		Date:                    today,
		CheckInAt:               &now,
		CheckInLat:              &in.Lat,
		CheckInLng:              &in.Lng,
		CheckInAddress:          address,
		Status:                  models.AttendanceStatus(statusStr),
		CheckInOutsideGeofence: outside,
	}
	if err := s.attendances.Create(att); err != nil {
		return nil, err
	}

	msg := "Absen masuk berhasil dicatat"
	if outside {
		msg = "Absen masuk berhasil, tetapi Anda berada di luar radius geofence (" +
			strconv.FormatFloat(distance, 'f', 0, 64) + " m dari lokasi PKL)"
	}
	if statusStr == "late" {
		msg += " — status: terlambat"
	}
	return &CheckResult{Attendance: att, Distance: distance, Address: address, Message: msg}, nil
}

func (s *AttendanceService) CheckOut(in CheckInInput) (*CheckResult, error) {
	now := utils.JakartaNow()
	today := models.DateOnlyFromTime(now)

	existing, err := s.attendances.FindByUserAndDate(in.UserID, today)
	if err != nil || existing.CheckInAt == nil {
		return nil, errors.New("anda belum absen masuk hari ini")
	}
	if existing.CheckOutAt != nil {
		return nil, errors.New("anda sudah absen pulang hari ini")
	}

	hasJournal, _ := s.journals.ExistsByUserAndDate(in.UserID, today)
	if !hasJournal {
		return nil, errors.New("jurnal harian wajib diisi sebelum absen pulang")
	}

	settings, _ := s.settings.GetAll()
	geoLat, _ := strconv.ParseFloat(settings["geofence_lat"], 64)
	geoLng, _ := strconv.ParseFloat(settings["geofence_lng"], 64)
	geoRad, _ := strconv.ParseFloat(settings["geofence_radius"], 64)
	distance := utils.DistanceMeters(in.Lat, in.Lng, geoLat, geoLng)
	outside := distance > geoRad

	address := s.geocoder.ReverseGeocode(in.Lat, in.Lng)

	existing.CheckOutAt = &now
	existing.CheckOutLat = &in.Lat
	existing.CheckOutLng = &in.Lng
	existing.CheckOutAddress = address
	existing.CheckOutOutsideGeofence = outside
	if err := s.attendances.Update(existing); err != nil {
		return nil, err
	}

	msg := "Absen pulang berhasil dicatat"
	if outside {
		msg = "Absen pulang berhasil, tetapi Anda berada di luar radius geofence (" +
			strconv.FormatFloat(distance, 'f', 0, 64) + " m dari lokasi PKL)"
	}
	return &CheckResult{Attendance: existing, Distance: distance, Address: address, Message: msg}, nil
}

func (s *AttendanceService) Today(userID uint) *models.Attendance {
	today := models.DateOnlyFromTime(utils.JakartaNow())
	att, err := s.attendances.FindByUserAndDate(userID, today)
	if err != nil {
		return nil
	}
	return att
}

func (s *AttendanceService) History(userID uint, fromStr, toStr string) ([]models.Attendance, error) {
	now := utils.JakartaNow()
	from := models.DateOnlyFromTime(utils.FirstOfMonth(now))
	to := models.DateOnlyFromTime(utils.LastOfMonth(now))
	if fromStr != "" {
		d, err := models.ParseDateOnly(fromStr)
		if err != nil {
			return nil, errors.New("format 'from' harus YYYY-MM-DD")
		}
		from = d
	}
	if toStr != "" {
		d, err := models.ParseDateOnly(toStr)
		if err != nil {
			return nil, errors.New("format 'to' harus YYYY-MM-DD")
		}
		to = d
	}
	return s.attendances.ListByUserAndRange(userID, from, to)
}

type RecapSummary struct {
	Year        int                 `json:"year"`
	Month       int                 `json:"month"`
	WorkingDays int                 `json:"working_days"`
	Present     int                 `json:"present"`
	Late        int                 `json:"late"`
	Izin        int                 `json:"izin"`
	Sakit       int                 `json:"sakit"`
	Alpha       int                 `json:"alpha"`
	Records     []models.Attendance `json:"records"`
	Holidays    []models.Holiday    `json:"holidays"`
}

func (s *AttendanceService) Recap(userID uint, month, year int) (*RecapSummary, error) {
	now := utils.JakartaNow()
	if month < 1 || month > 12 {
		month = int(now.Month())
	}
	if year <= 0 {
		year = now.Year()
	}
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, utils.JakartaLoc())
	last := first.AddDate(0, 1, -1)

	firstDO := models.DateOnlyFromTime(first)
	lastDO := models.DateOnlyFromTime(last)

	rows, err := s.attendances.ListByUserAndRange(userID, firstDO, lastDO)
	if err != nil {
		return nil, err
	}
	holidays, _ := s.holidays.ListInRange(firstDO, lastDO)

	settings, _ := s.settings.GetAll()
	workDays := utils.ParseWorkDays(settings["work_days"])
	holidaySet := make(map[string]bool)
	for _, h := range holidays {
		holidaySet[h.Date.String()] = true
	}

	recap := &RecapSummary{Year: year, Month: month, Records: rows, Holidays: holidays}
	for d := first; !d.After(last); d = d.AddDate(0, 0, 1) {
		if !utils.IsWorkingDay(d, workDays) {
			continue
		}
		if holidaySet[models.DateOnlyFromTime(d).String()] {
			continue
		}
		recap.WorkingDays++
	}

	for _, r := range rows {
		switch r.Status {
		case models.StatusPresent:
			recap.Present++
		case models.StatusLate:
			recap.Late++
		case models.StatusIzin:
			recap.Izin++
		case models.StatusSakit:
			recap.Sakit++
		case models.StatusAlpha:
			recap.Alpha++
		}
	}
	filled := recap.Present + recap.Late + recap.Izin + recap.Sakit + recap.Alpha
	if recap.WorkingDays > filled {
		recap.Alpha += recap.WorkingDays - filled
	}
	return recap, nil
}