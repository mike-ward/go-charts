package scale

import (
	"math"
	"testing"
)

func TestLinearMapEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		min, max float64
		value    float64
		pMin     float32
		pMax     float32
		want     float32
	}{
		{"NaN value", 0, 100, math.NaN(), 0, 500, 0},
		{"+Inf value", 0, 100, math.Inf(1), 0, 500, 0},
		{"-Inf value", 0, 100, math.Inf(-1), 0, 500, 0},
		{"zero range", 50, 50, 50, 0, 500, 0},
		{"very large domain overflow", -1e308, 1e308, 0, 0, 1000, 0},
		{"very small domain", 0, 1e-300, 5e-301, 0, 1000, 500},
		{"inverted domain", 100, 0, 50, 0, 500, 250},
		{"NaN min", math.NaN(), 100, 50, 0, 500, 0},
		{"NaN max", 0, math.NaN(), 50, 0, 500, 0},
		{"equal pixel range", 0, 100, 50, 200, 200, 200},
		{"negative pixels", 0, 100, 50, -100, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			got := s.Map(tt.value, tt.pMin, tt.pMax)
			if math.IsNaN(float64(got)) {
				t.Errorf("Map returned NaN")
				return
			}
			if math.Abs(float64(got-tt.want)) > 0.5 {
				t.Errorf("Map(%v) = %v, want %v",
					tt.value, got, tt.want)
			}
		})
	}
}

func TestLinearInvertEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		min, max float64
		pixel    float32
		pMin     float32
		pMax     float32
		want     float64
	}{
		{"equal pixel range", 0, 100, 200, 200, 200, 0},
		{"zero pixel range at origin", 0, 100, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			got := s.Invert(tt.pixel, tt.pMin, tt.pMax)
			if math.IsNaN(got) {
				t.Errorf("Invert returned NaN")
				return
			}
			if math.Abs(got-tt.want) > 0.5 {
				t.Errorf("Invert(%v) = %v, want %v",
					tt.pixel, got, tt.want)
			}
		})
	}
}
