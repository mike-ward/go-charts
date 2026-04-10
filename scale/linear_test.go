package scale

import (
	"math"
	"testing"
)

func TestLinearTransform(t *testing.T) {
	tests := []struct {
		name           string
		min, max       float64
		value          float64
		pixMin, pixMax float32
		want           float32
	}{
		{"midpoint", 0, 100, 50, 0, 500, 250},
		{"at min", 0, 100, 0, 0, 500, 0},
		{"at max", 0, 100, 100, 0, 500, 500},
		{"negative domain", -100, 100, 0, 0, 400, 200},
		{"pixel offset", 0, 100, 25, 100, 300, 150},
		{"value below min", 0, 100, -10, 0, 100, -10},
		{"value above max", 0, 100, 110, 0, 100, 110},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			got := s.Transform(tt.value, tt.pixMin, tt.pixMax)
			if math.Abs(float64(got-tt.want)) > 0.01 {
				t.Errorf("Transform(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLinearTransformDegenerate(t *testing.T) {
	tests := []struct {
		name     string
		min, max float64
		value    float64
	}{
		{"zero domain", 50, 50, 50},
		{"NaN value", 0, 100, math.NaN()},
		{"Inf value", 0, 100, math.Inf(1)},
		{"NaN domain min", math.NaN(), 100, 50},
		{"NaN domain max", 0, math.NaN(), 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			got := s.Transform(tt.value, 0, 500)
			// Non-finite inputs must return pixelMin (0), not NaN or Inf.
			if math.IsNaN(float64(got)) || math.IsInf(float64(got), 0) {
				t.Errorf("Transform returned non-finite %v", got)
			}
		})
	}
}

func TestLinearInvert(t *testing.T) {
	tests := []struct {
		name           string
		min, max       float64
		pixel          float32
		pixMin, pixMax float32
		want           float64
	}{
		{"midpoint", 0, 100, 250, 0, 500, 50},
		{"at pixMin", 0, 100, 0, 0, 500, 0},
		{"at pixMax", 0, 100, 500, 0, 500, 100},
		{"negative domain", -100, 100, 200, 0, 400, 0},
		{"pixel offset", 0, 100, 150, 100, 300, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			got := s.Invert(tt.pixel, tt.pixMin, tt.pixMax)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("Invert(%v) = %v, want %v", tt.pixel, got, tt.want)
			}
		})
	}
}

func TestLinearInvertZeroPixelRange(t *testing.T) {
	s := NewLinear(0, 100)
	got := s.Invert(50, 200, 200)
	if got != s.min {
		t.Errorf("Invert with zero pixel range = %v, want %v (domain min)", got, s.min)
	}
}

func TestLinearRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		min, max float64
		value    float64
	}{
		{"zero midpoint", -100, 100, 0},
		{"quarter point", 0, 400, 100},
		{"float value", 0, 1, 0.333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLinear(tt.min, tt.max)
			px := s.Transform(tt.value, 0, 1000)
			got := s.Invert(px, 0, 1000)
			if math.Abs(got-tt.value) > 1e-9 {
				t.Errorf("round-trip(%v): got %v", tt.value, got)
			}
		})
	}
}
