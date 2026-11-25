package services

import (
	"errors"
	"strconv"
	"strings"
)

type SettingsAdminService struct {
	settings *SettingService
	logs     *ActivityLogService
}

func NewSettingsAdminService(s *SettingService, l *ActivityLogService) *SettingsAdminService {
	return &SettingsAdminService{settings: s, logs: l}
}

func (s *SettingsAdminService) GetAll() (map[string]string, error) { return s.settings.GetAll() }

func (s *SettingsAdminService) Update(adminID uint, updates map[string]string) (map[string]string, error) {
	allowed := map[string]bool{
		"work_days":         true,
		"check_in_time":     true,
		"tolerance_minutes": true,
		"geofence_lat":      true,
		"geofence_lng":      true,
		"geofence_radius":   true,
		"pkl_location_name": true,
	}
	validated := map[string]string{}
	for k, v := range updates {
		if !allowed[k] {
			continue
		}
		v = strings.TrimSpace(v)
		switch k {
		case "work_days":
			parts := strings.Split(v, ",")
			for _, p := range parts {
				n, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil || n < 0 || n > 6 {
					return nil, errors.New("work_days harus berupa angka 0-6 dipisah koma")
				}
			}
		case "check_in_time":
			if _, _, err := parseHourMinute(v); err != nil {
				return nil, errors.New("check_in_time harus format HH:MM")
			}
		case "tolerance_minutes":
			if n, err := strconv.Atoi(v); err != nil || n < 0 || n > 120 {
				return nil, errors.New("tolerance_minutes harus angka 0-120")
			}
		case "geofence_lat":
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, errors.New("geofence_lat tidak valid")
			}
		case "geofence_lng":
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, errors.New("geofence_lng tidak valid")
			}
		case "geofence_radius":
			if n, err := strconv.ParseFloat(v, 64); err != nil || n < 10 || n > 10000 {
				return nil, errors.New("geofence_radius harus 10-10000 (meter)")
			}
		case "pkl_location_name":
			if v == "" {
				return nil, errors.New("pkl_location_name tidak boleh kosong")
			}
		}
		validated[k] = v
	}
	if len(validated) == 0 {
		return nil, errors.New("tidak ada setting yang diperbarui")
	}
	for k, v := range validated {
		if err := s.settings.Set(k, v); err != nil {
			return nil, err
		}
	}
	s.logs.Log(adminID, "update_settings", "Memperbarui pengaturan aplikasi",
		map[string]interface{}{"keys": keysOf(validated)})
	return s.settings.GetAll()
}

func parseHourMinute(s string) (int, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid")
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, 0, errors.New("invalid")
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, 0, errors.New("invalid")
	}
	return h, m, nil
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}