package series

import (
	"math"
	"testing"
)

func TestNewXYZ(t *testing.T) {
	s := NewXYZ(XYZCfg{
		Name:   "test",
		Points: []XYZPoint{{1, 2, 3}, {4, 5, 6}},
	})
	if s.Name() != "test" {
		t.Errorf("Name() = %q, want %q", s.Name(), "test")
	}
	if s.Len() != 2 {
		t.Errorf("Len() = %d, want 2", s.Len())
	}
}

func TestXYZFromSlices(t *testing.T) {
	s, err := XYZFromSlices("s",
		[]float64{1, 2}, []float64{3, 4}, []float64{5, 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Errorf("Len() = %d, want 2", s.Len())
	}
	if s.Points[0].Z != 5 {
		t.Errorf("Points[0].Z = %v, want 5", s.Points[0].Z)
	}
}

func TestXYZFromSlicesMismatch(t *testing.T) {
	_, err := XYZFromSlices("s",
		[]float64{1, 2}, []float64{3}, []float64{5, 6})
	if err == nil {
		t.Error("expected error for mismatched slice lengths")
	}
}

func TestXYZBounds(t *testing.T) {
	s := NewXYZ(XYZCfg{
		Points: []XYZPoint{
			{1, 10, 100},
			{5, 2, 200},
			{3, 8, 50},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 1 || maxX != 5 {
		t.Errorf("X bounds = (%v, %v), want (1, 5)", minX, maxX)
	}
	if minY != 2 || maxY != 10 {
		t.Errorf("Y bounds = (%v, %v), want (2, 10)", minY, maxY)
	}
}

func TestXYZBoundsSkipsNonFinite(t *testing.T) {
	s := NewXYZ(XYZCfg{
		Points: []XYZPoint{
			{math.NaN(), 1, 1},
			{2, 3, 10},
			{4, math.Inf(1), 20},
			{6, 7, 30},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 2 || maxX != 6 {
		t.Errorf("X bounds = (%v, %v), want (2, 6)", minX, maxX)
	}
	if minY != 3 || maxY != 7 {
		t.Errorf("Y bounds = (%v, %v), want (3, 7)", minY, maxY)
	}
}

func TestXYZBoundsEmpty(t *testing.T) {
	s := NewXYZ(XYZCfg{})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 0 || maxX != 0 || minY != 0 || maxY != 0 {
		t.Errorf("empty bounds = (%v,%v,%v,%v), want zeros",
			minX, maxX, minY, maxY)
	}
}

func TestXYZZBounds(t *testing.T) {
	s := NewXYZ(XYZCfg{
		Points: []XYZPoint{
			{1, 2, 50},
			{3, 4, 10},
			{5, 6, 200},
		},
	})
	minZ, maxZ := s.ZBounds()
	if minZ != 10 || maxZ != 200 {
		t.Errorf("Z bounds = (%v, %v), want (10, 200)", minZ, maxZ)
	}
}

func TestXYZZBoundsSkipsNonFinite(t *testing.T) {
	s := NewXYZ(XYZCfg{
		Points: []XYZPoint{
			{1, 2, math.NaN()},
			{3, 4, 5},
			{5, 6, math.Inf(1)},
			{7, 8, 15},
		},
	})
	minZ, maxZ := s.ZBounds()
	if minZ != 5 || maxZ != 15 {
		t.Errorf("Z bounds = (%v, %v), want (5, 15)", minZ, maxZ)
	}
}

func TestXYZZBoundsEmpty(t *testing.T) {
	s := NewXYZ(XYZCfg{})
	minZ, maxZ := s.ZBounds()
	if minZ != 0 || maxZ != 0 {
		t.Errorf("empty Z bounds = (%v, %v), want zeros", minZ, maxZ)
	}
}

func TestXYZString(t *testing.T) {
	s := NewXYZ(XYZCfg{Name: "test", Points: []XYZPoint{{1, 2, 3}}})
	got := s.String()
	if got != `XYZ{"test", 1 points}` {
		t.Errorf("String() = %q", got)
	}
}

func TestXYZPointString(t *testing.T) {
	p := XYZPoint{1.5, 2.5, 3.5}
	got := p.String()
	if got != "(1.5, 2.5, 3.5)" {
		t.Errorf("String() = %q", got)
	}
}
