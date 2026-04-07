package series

import (
	"math"
	"testing"
)

func TestErrorXYBounds(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "test",
		Points: []ErrorPoint{
			{X: 1, Y: 10, YErr: ErrorBar{Low: 2, High: 3}},
			{X: 5, Y: 20, YErr: ErrorBar{Low: 1, High: 4}, XErr: ErrorBar{Low: 0.5, High: 0.5}},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	// Point 1: X=1 (no XErr), Point 2: X=5 XErr={0.5,0.5} → [4.5,5.5]
	if minX != 1 {
		t.Errorf("minX = %v, want 1", minX)
	}
	if maxX != 5.5 {
		t.Errorf("maxX = %v, want 5.5", maxX)
	}
	// Point 1: Y=10 YErr={2,3} → [8,13], Point 2: Y=20 YErr={1,4} → [19,24]
	if minY != 8 {
		t.Errorf("minY = %v, want 8", minY)
	}
	if maxY != 24 {
		t.Errorf("maxY = %v, want 24", maxY)
	}
}

func TestErrorXYBoundsZeroError(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "no-err",
		Points: []ErrorPoint{
			{X: 1, Y: 10},
			{X: 5, Y: 20},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 1 || maxX != 5 || minY != 10 || maxY != 20 {
		t.Errorf("bounds = (%v,%v,%v,%v), want (1,5,10,20)",
			minX, maxX, minY, maxY)
	}
}

func TestErrorXYBoundsNonFinite(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "nan",
		Points: []ErrorPoint{
			{X: math.NaN(), Y: 10},
			{X: 1, Y: math.Inf(1)},
			{X: 2, Y: 20, YErr: ErrorBar{Low: 1, High: 1}},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 2 || maxX != 2 || minY != 19 || maxY != 21 {
		t.Errorf("bounds = (%v,%v,%v,%v), want (2,2,19,21)",
			minX, maxX, minY, maxY)
	}
}

func TestErrorXYBoundsEmpty(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{Name: "empty"})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 0 || maxX != 0 || minY != 0 || maxY != 0 {
		t.Errorf("bounds = (%v,%v,%v,%v), want all zero",
			minX, maxX, minY, maxY)
	}
}

func TestErrorXYBoundsNegativeError(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "neg-err",
		Points: []ErrorPoint{
			{X: 5, Y: 10, YErr: ErrorBar{Low: -3, High: 2}},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	// Negative error clamped to 0: minY = 10-0 = 10
	if minY != 10 {
		t.Errorf("minY = %v, want 10 (negative error clamped)", minY)
	}
	if maxY != 12 {
		t.Errorf("maxY = %v, want 12", maxY)
	}
	if minX != 5 || maxX != 5 {
		t.Errorf("X bounds = (%v,%v), want (5,5)", minX, maxX)
	}
}

func TestErrorXYInterface(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "test",
		Points: []ErrorPoint{
			{X: 1, Y: 2},
			{X: 3, Y: 4},
		},
	})
	if s.Name() != "test" {
		t.Errorf("Name = %q, want %q", s.Name(), "test")
	}
	if s.Len() != 2 {
		t.Errorf("Len = %d, want 2", s.Len())
	}
	str := s.String()
	if str != `ErrorXY{"test", 2 points}` {
		t.Errorf("String = %q", str)
	}
}

func TestSymmetric(t *testing.T) {
	e := Symmetric(1.5)
	if e.Low != 1.5 || e.High != 1.5 {
		t.Errorf("Symmetric(1.5) = %+v, want {1.5, 1.5}", e)
	}
}

func TestErrorXYBoundsNonFiniteErrorBar(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Name: "nan-err",
		Points: []ErrorPoint{
			{X: 1, Y: 10, YErr: ErrorBar{
				Low: math.NaN(), High: math.Inf(1),
			}, XErr: ErrorBar{
				Low: math.Inf(-1), High: math.NaN(),
			}},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	// Non-finite errors clamp to 0 → bounds collapse to point.
	if minX != 1 || maxX != 1 || minY != 10 || maxY != 10 {
		t.Errorf("bounds = (%v,%v,%v,%v), want (1,1,10,10)",
			minX, maxX, minY, maxY)
	}
}

func TestErrorXYBoundsSinglePoint(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Points: []ErrorPoint{
			{X: 5, Y: 10, YErr: ErrorBar{Low: 1, High: 2}},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 5 || maxX != 5 || minY != 9 || maxY != 12 {
		t.Errorf("bounds = (%v,%v,%v,%v), want (5,5,9,12)",
			minX, maxX, minY, maxY)
	}
}

func TestErrorXYBoundsMixedErrors(t *testing.T) {
	s := NewErrorXY(ErrorXYCfg{
		Points: []ErrorPoint{
			{X: 1, Y: 10}, // no err
			{X: 2, Y: 20, YErr: ErrorBar{Low: 5, High: 5}}, // with err
			{X: 3, Y: 15}, // no err
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 1 || maxX != 3 || minY != 10 || maxY != 25 {
		t.Errorf("bounds = (%v,%v,%v,%v), want (1,3,10,25)",
			minX, maxX, minY, maxY)
	}
}
