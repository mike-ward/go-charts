package scale

import (
	"math"
	"testing"
)

func TestLogTransformEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		min, max     float64
		base         float64
		value        float64
		pMin, pMax   float32
		wantPixelMin bool // expect pixelMin returned
	}{
		{"base 1 defaults to 10", 1, 1000, 1, 100, 0, 300, false},
		{"base 0 defaults to 10", 1, 1000, 0, 100, 0, 300, false},
		{"base -5 defaults to 10", 1, 1000, -5, 100, 0, 300, false},
		{"min == max positive", 10, 10, 10, 10, 0, 300, true},
		{"NaN value", 1, 1000, 10, math.NaN(), 0, 300, true},
		{"+Inf value", 1, 1000, 10, math.Inf(1), 0, 300, true},
		{"-Inf value", 1, 1000, 10, math.Inf(-1), 0, 300, true},
		{"zero value", 1, 1000, 10, 0, 0, 300, true},
		{"negative value", 1, 1000, 10, -5, 0, 300, true},
		{"min > max", 1000, 1, 10, 100, 0, 300, true},
		{"min == 0", 0, 1000, 10, 100, 0, 300, true},
		{"min < 0", -1, 1000, 10, 100, 0, 300, true},
		{"very close min/max", 1.0, 1.0 + 1e-15, 10, 1.0, 0, 300, false},
		{"very large domain", 1e100, 1e200, 10, 1e150, 0, 300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(tt.min, tt.max, tt.base)
			got := s.Transform(tt.value, tt.pMin, tt.pMax)
			if math.IsNaN(float64(got)) {
				t.Fatalf("Transform returned NaN")
			}
			if math.IsInf(float64(got), 0) {
				t.Fatalf("Transform returned Inf")
			}
			if tt.wantPixelMin && got != tt.pMin {
				t.Errorf("Transform = %v, want pixelMin %v", got, tt.pMin)
			}
		})
	}
}

func TestLogInvertEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		min, max   float64
		base       float64
		pixel      float32
		pMin, pMax float32
		wantMin    bool // expect s.min returned
	}{
		{"equal pixel range", 1, 1000, 10, 100, 100, 100, true},
		{"min <= 0", 0, 1000, 10, 150, 0, 300, true},
		{"min == max", 10, 10, 10, 150, 0, 300, true},
		{"base 1 safe", 1, 1000, 1, 150, 0, 300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(tt.min, tt.max, tt.base)
			got := s.Invert(tt.pixel, tt.pMin, tt.pMax)
			if math.IsNaN(got) {
				t.Fatalf("Invert returned NaN")
			}
			if math.IsInf(got, 0) {
				t.Fatalf("Invert returned Inf")
			}
			if tt.wantMin && got != s.min {
				t.Errorf("Invert = %v, want s.min %v", got, s.min)
			}
		})
	}
}
