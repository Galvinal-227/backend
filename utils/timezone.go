package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var jakartaLoc *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	jakartaLoc = loc
}

func JakartaLoc() *time.Location { return jakartaLoc }
func JakartaNow() time.Time       { return time.Now().In(jakartaLoc) }

// ParseWorkDays mengubah "1,2,3,4,5" menjadi []int.
func ParseWorkDays(csv string) []int {
	parts := strings.Split(csv, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err == nil && n >= 0 && n <= 6 {
			out = append(out, n)
		}
	}
	return out
}

// IsWorkingDay: time.Weekday() Sunday=0..Saturday=6, sama dengan konvensi kita.
func IsWorkingDay(date time.Time, workDays []int) bool {
	wd := int(date.In(jakartaLoc).Weekday())
	for _, d := range workDays {
		if d == wd {
			return true
		}
	}
	return false
}

// ParseHHMM mengembalikan jam dan menit dari "07:30".
func ParseHHMM(s string) (int, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s", s)
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("time out of range: %s", s)
	}
	return h, m, nil
}

// DetermineAttendanceStatus mengembalikan "present" atau "late"
// berdasarkan jam masuk + toleransi. Default: 07:30 / 15 menit.
func DetermineAttendanceStatus(now time.Time, checkInTimeStr string, toleranceStr string) string {
	h, m, err := ParseHHMM(checkInTimeStr)
	if err != nil {
		h, m = 7, 30
	}
	tol, err := strconv.Atoi(toleranceStr)
	if err != nil || tol < 0 {
		tol = 15
	}
	jt := now.In(jakartaLoc)
	limit := time.Date(jt.Year(), jt.Month(), jt.Day(), h, m, 0, 0, jakartaLoc).
		Add(time.Duration(tol) * time.Minute)
	if jt.After(limit) {
		return "late"
	}
	return "present"
}

// FirstOfMonth mengembalikan tanggal 1 dari bulan berjalan di zona Jakarta.
func FirstOfMonth(t time.Time) time.Time {
	jt := t.In(jakartaLoc)
	return time.Date(jt.Year(), jt.Month(), 1, 0, 0, 0, 0, jakartaLoc)
}

// LastOfMonth mengembalikan tanggal terakhir bulan berjalan.
func LastOfMonth(t time.Time) time.Time {
	return FirstOfMonth(t).AddDate(0, 1, -1)
}