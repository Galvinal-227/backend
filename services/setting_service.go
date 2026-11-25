package services

import (
	"pkl-attendance/backend/repositories"
)

var settingDefaults = map[string]string{
	"work_days":         "1,2,3,4,5",
	"check_in_time":     "07:30",
	"tolerance_minutes": "15",
	"geofence_lat":      "-6.200000",
	"geofence_lng":      "106.816666",
	"geofence_radius":   "100",
	"pkl_location_name": "Lokasi PKL",
}

type SettingService struct{ repo *repositories.SettingRepository }

func NewSettingService(r *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: r}
}

func (s *SettingService) GetAll() (map[string]string, error) {
	out := make(map[string]string, len(settingDefaults))
	for k, v := range settingDefaults {
		out[k] = v
	}
	rows, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.Value != "" {
			out[r.Key] = r.Value
		}
	}
	return out, nil
}

func (s *SettingService) Set(key, value string) error {
	return s.repo.Upsert(key, value)
}