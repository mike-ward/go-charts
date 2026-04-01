package scale

import (
	"math"
	"testing"
)

func TestLinearMap(t *testing.T) {
	s := NewLinear(0, 100)
	got := s.Map(50, 0, 500)
	if math.Abs(float64(got)-250) > 0.01 {
		t.Errorf("Map(50) = %v, want 250", got)
	}
}

func TestLinearInvert(t *testing.T) {
	s := NewLinear(0, 100)
	got := s.Invert(250, 0, 500)
	if math.Abs(got-50) > 0.01 {
		t.Errorf("Invert(250) = %v, want 50", got)
	}
}

func TestLinearZeroRange(t *testing.T) {
	s := NewLinear(50, 50)
	got := s.Map(50, 0, 500)
	if got != 0 {
		t.Errorf("Map with zero range = %v, want 0", got)
	}
}
