package utils

import "math"

const earthRadiusMeters = 6371000.0

// DistanceMeters menghitung jarak antara dua koordinat dengan formula Haversine.
func DistanceMeters(lat1, lng1, lat2, lng2 float64) float64 {
	const rad = math.Pi / 180.0
	phi1 := lat1 * rad
	phi2 := lat2 * rad
	dphi := (lat2 - lat1) * rad
	dlam := (lng2 - lng1) * rad
	a := math.Sin(dphi/2)*math.Sin(dphi/2) +
		math.Cos(phi1)*math.Cos(phi2)*math.Sin(dlam/2)*math.Sin(dlam/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}