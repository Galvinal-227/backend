package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type geocodeEntry struct {
	address string
	exp     time.Time
}

type GeocodingService struct {
	client   *http.Client
	mu       sync.Mutex
	cache    map[string]*geocodeEntry
	lastCall time.Time
}

func NewGeocodingService() *GeocodingService {
	return &GeocodingService{
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  make(map[string]*geocodeEntry),
	}
}

func cacheKey(lat, lng float64) string {
	return fmt.Sprintf("%.4f:%.4f", lat, lng)
}

// ReverseGeocode memakai Nominatim. Gagal geocode bukan error fatal.
func (s *GeocodingService) ReverseGeocode(lat, lng float64) string {
	key := cacheKey(lat, lng)
	s.mu.Lock()
	if e, ok := s.cache[key]; ok && time.Now().Before(e.exp) {
		s.mu.Unlock()
		return e.address
	}
	// Rate limit 1 req/detik terhadap Nominatim
	if delta := time.Since(s.lastCall); delta < time.Second {
		time.Sleep(time.Second - delta)
	}
	s.lastCall = time.Now()
	s.mu.Unlock()

	address := s.fetch(lat, lng)
	s.mu.Lock()
	s.cache[key] = &geocodeEntry{address: address, exp: time.Now().Add(24 * time.Hour)}
	s.mu.Unlock()
	return address
}

func (s *GeocodingService) fetch(lat, lng float64) string {
	q := url.Values{}
	q.Set("format", "json")
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lng))
	q.Set("zoom", "18")
	q.Set("addressdetails", "0")

	req, err := http.NewRequest("GET", "https://nominatim.openstreetmap.org/reverse?"+q.Encode(), nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "PKL-Attendance/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var out struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return ""
	}
	return out.DisplayName
}