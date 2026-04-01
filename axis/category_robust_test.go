package axis

import (
	"math"
	"testing"
)

func TestCategoryEmpty(t *testing.T) {
	a := NewCategory(CategoryCfg{})

	ticks := a.Ticks(0, 500)
	if ticks != nil {
		t.Errorf("Ticks on empty = %v, want nil", ticks)
	}

	got := a.Transform(0, 0, 500)
	if got != 0 {
		t.Errorf("Transform on empty = %v, want 0", got)
	}

	inv := a.Inverse(250, 0, 500)
	if inv != 0 {
		t.Errorf("Inverse on empty = %v, want 0", inv)
	}
}

func TestCategorySingle(t *testing.T) {
	a := NewCategory(CategoryCfg{Categories: []string{"A"}})

	ticks := a.Ticks(0, 500)
	if len(ticks) != 1 {
		t.Fatalf("Ticks len = %d, want 1", len(ticks))
	}
	if ticks[0].Label != "A" {
		t.Errorf("Ticks[0].Label = %q, want A", ticks[0].Label)
	}
}

func TestCategoryZeroPixelRange(t *testing.T) {
	a := NewCategory(CategoryCfg{Categories: []string{"A", "B"}})

	// pixelMin == pixelMax → step == 0
	got := a.Transform(0, 100, 100)
	if math.IsNaN(float64(got)) || math.IsInf(float64(got), 0) {
		t.Errorf("Transform with zero pixel range = %v", got)
	}

	inv := a.Inverse(100, 100, 100)
	if math.IsNaN(inv) || math.IsInf(inv, 0) {
		t.Errorf("Inverse with zero pixel range = %v", inv)
	}
}

func TestCategoryOutOfRange(t *testing.T) {
	a := NewCategory(CategoryCfg{
		Categories: []string{"A", "B", "C"},
	})

	// Negative value clamped.
	got := a.Transform(-1, 0, 300)
	if got < 0 {
		t.Errorf("Transform(-1) = %v, expected >= 0", got)
	}

	// Value beyond n clamped.
	got = a.Transform(10, 0, 300)
	if got > 300 {
		t.Errorf("Transform(10) = %v, expected <= 300", got)
	}
}

func TestCategoryLargeCount(t *testing.T) {
	cats := make([]string, 1000)
	for i := range cats {
		cats[i] = "x"
	}
	a := NewCategory(CategoryCfg{Categories: cats})
	ticks := a.Ticks(0, 1000)
	if len(ticks) != 1000 {
		t.Errorf("Ticks len = %d, want 1000", len(ticks))
	}
}
