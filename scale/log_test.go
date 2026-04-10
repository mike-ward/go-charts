package scale

import (
	"math"
	"testing"
)

func TestLogTransform(t *testing.T) {
	tests := []struct {
		name           string
		min, max, base float64
		value          float64
		pixMin, pixMax float32
		want           float32
		wantPixMin     bool // true if result must == pixMin (degenerate)
	}{
		{"mid-decade", 1, 1000, 10, 100, 0, 300, 200, false},
		{"at min", 1, 1000, 10, 1, 0, 300, 0, false},
		{"at max", 1, 1000, 10, 1000, 0, 300, 300, false},
		{"non-positive value", 1, 1000, 10, -1, 0, 300, 0, true},
		{"zero value", 1, 1000, 10, 0, 0, 300, 0, true},
		{"non-positive min", -1, 1000, 10, 100, 0, 300, 0, true},
		{"NaN value", 1, 1000, 10, math.NaN(), 0, 300, 0, true},
		{"base 2", 1, 8, 2, 4, 0, 300, 200, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(tt.min, tt.max, tt.base)
			got := s.Transform(tt.value, tt.pixMin, tt.pixMax)
			if tt.wantPixMin {
				if got != tt.pixMin {
					t.Errorf("Transform(%v) = %v, want pixMin (%v)",
						tt.value, got, tt.pixMin)
				}
				return
			}
			if math.Abs(float64(got-tt.want)) > 0.5 {
				t.Errorf("Transform(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLogInvert(t *testing.T) {
	tests := []struct {
		name           string
		min, max, base float64
		pixel          float32
		pixMin, pixMax float32
		want           float64
		wantMin        bool // true if result must == domain min (degenerate)
	}{
		{"mid-decade", 1, 1000, 10, 200, 0, 300, 100, false},
		{"at pixMin", 1, 1000, 10, 0, 0, 300, 1, false},
		{"at pixMax", 1, 1000, 10, 300, 0, 300, 1000, false},
		{"zero pixel range", 1, 1000, 10, 150, 200, 200, 1, true},
		{"non-positive domain min", -1, 1000, 10, 150, 0, 300, -1, true},
		{"base 2", 1, 8, 2, 200, 0, 300, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(tt.min, tt.max, tt.base)
			got := s.Invert(tt.pixel, tt.pixMin, tt.pixMax)
			if tt.wantMin {
				if got != tt.min {
					t.Errorf("Invert(%v) = %v, want domain min (%v)",
						tt.pixel, got, tt.min)
				}
				return
			}
			if math.Abs(got-tt.want)/tt.want > 0.01 {
				t.Errorf("Invert(%v) = %v, want %v", tt.pixel, got, tt.want)
			}
		})
	}
}

func TestLogRoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		min, max, base float64
		value          float64
	}{
		{"decade midpoint base10", 1, 1000, 10, 100},
		{"arbitrary base2", 1, 64, 2, 8},
		{"near min", 1, 1000, 10, 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(tt.min, tt.max, tt.base)
			px := s.Transform(tt.value, 0, 1000)
			got := s.Invert(px, 0, 1000)
			if math.Abs(got-tt.value)/tt.value > 1e-6 {
				t.Errorf("round-trip(%v): got %v", tt.value, got)
			}
		})
	}
}

func TestLogBaseDefaults(t *testing.T) {
	tests := []struct {
		name string
		base float64
	}{
		{"zero base", 0},
		{"negative base", -5},
		{"base 1", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLog(1, 1000, tt.base)
			if s.Base() != 10 {
				t.Errorf("base %v: got %v, want 10 (default)", tt.base, s.Base())
			}
		})
	}
}
