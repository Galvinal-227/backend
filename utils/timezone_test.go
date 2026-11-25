package utils

import (
	"testing"
	"time"
)

func TestParseHHMM(t *testing.T) {
	h, m, err := ParseHHMM("07:30")
	if err != nil || h != 7 || m != 30 {
		t.Fatalf("unexpected: %d:%d %v", h, m, err)
	}
	if _, _, err := ParseHHMM("25:00"); err == nil {
		t.Fatal("expected error")
	}
}

func TestIsWorkingDayDefault(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	monday := time.Date(2026, 11, 9, 8, 0, 0, 0, loc)
	saturday := time.Date(2026, 11, 14, 8, 0, 0, 0, loc)
	days := []int{1, 2, 3, 4, 5}
	if !IsWorkingDay(monday, days) {
		t.Fatal("Monday should be working day")
	}
	if IsWorkingDay(saturday, days) {
		t.Fatal("Saturday should not be working day")
	}
}

func TestDetermineAttendanceStatus(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	cases := []struct {
		name  string
		hour, min int
		want  string
	}{
		{"07:00 present", 7, 0, "present"},
		{"07:30 present", 7, 30, "present"},
		{"07:45 present", 7, 45, "present"},
		{"07:46 late", 7, 46, "late"},
		{"08:00 late", 8, 0, "late"},
	}
	for _, tc := range cases {
		now := time.Date(2026, 11, 9, tc.hour, tc.min, 0, 0, loc)
		if got := DetermineAttendanceStatus(now, "07:30", "15"); got != tc.want {
			t.Fatalf("%s: want %s got %s", tc.name, tc.want, got)
		}
	}
}