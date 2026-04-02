package scale

import (
	"math"
	"testing"
)

func TestLogTransform(t *testing.T) {
	s := NewLog(1, 1000, 10)
	got := s.Transform(100, 0, 300)
	want := float32(200)
	if math.Abs(float64(got-want)) > 0.01 {
		t.Errorf("Transform(100) = %v, want %v", got, want)
	}
}

func TestLogInvert(t *testing.T) {
	s := NewLog(1, 1000, 10)
	got := s.Invert(200, 0, 300)
	if math.Abs(got-100) > 0.5 {
		t.Errorf("Invert(200) = %v, want 100", got)
	}
}

func TestLogInvalidInput(t *testing.T) {
	s := NewLog(1, 1000, 10)
	got := s.Transform(-1, 0, 300)
	if got != 0 {
		t.Errorf("Transform(-1) = %v, want 0", got)
	}
}
