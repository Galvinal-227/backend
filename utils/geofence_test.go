package utils

import (
	"math"
	"testing"
)

func TestDistanceMetersSamePoint(t *testing.T) {
	if d := DistanceMeters(-6.2, 106.8, -6.2, 106.8); d != 0 {
		t.Fatalf("same point should be 0, got %v", d)
	}
}

func TestDistanceMetersOneDegreeLat(t *testing.T) {
	d := DistanceMeters(0, 0, 1, 0)
	expected := 111195.0
	if math.Abs(d-expected) > 500 {
		t.Fatalf("expected ~%v, got %v", expected, d)
	}
}

func TestDistanceMetersApprox100m(t *testing.T) {
	// 0.0009 deg lat ≈ 100m
	d := DistanceMeters(0, 0, 0.0009, 0)
	if d < 90 || d > 110 {
		t.Fatalf("expected ~100m, got %v", d)
	}
}