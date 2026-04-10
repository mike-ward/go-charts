package axis

import (
	"testing"
)

func TestCategoryTicksCount(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		wantLen    int
	}{
		{"empty", nil, 0},
		{"single", []string{"A"}, 1},
		{"multiple", []string{"A", "B", "C", "D"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewCategory(CategoryCfg{Categories: tt.categories})
			ticks := a.Ticks(0, 400)
			if len(ticks) != tt.wantLen {
				t.Errorf("Ticks len = %d, want %d", len(ticks), tt.wantLen)
			}
		})
	}
}

func TestCategoryTickLabels(t *testing.T) {
	cats := []string{"Alpha", "Beta", "Gamma"}
	a := NewCategory(CategoryCfg{Categories: cats})
	ticks := a.Ticks(0, 300)
	for i, tk := range ticks {
		if tk.Label != cats[i] {
			t.Errorf("ticks[%d].Label = %q, want %q", i, tk.Label, cats[i])
		}
	}
}

func TestCategoryTickPositionsCentered(t *testing.T) {
	// With 2 categories over [0, 200], each slot = 100px.
	// Centers should be at 50 and 150.
	a := NewCategory(CategoryCfg{Categories: []string{"X", "Y"}})
	ticks := a.Ticks(0, 200)
	if len(ticks) != 2 {
		t.Fatalf("expected 2 ticks, got %d", len(ticks))
	}
	if ticks[0].Position < 40 || ticks[0].Position > 60 {
		t.Errorf("ticks[0].Position = %v, want ~50", ticks[0].Position)
	}
	if ticks[1].Position < 140 || ticks[1].Position > 160 {
		t.Errorf("ticks[1].Position = %v, want ~150", ticks[1].Position)
	}
}

func TestCategoryDomain(t *testing.T) {
	tests := []struct {
		name    string
		cats    []string
		wantMin float64
		wantMax float64
	}{
		{"empty", nil, 0, 0},
		{"one", []string{"A"}, 0, 0},
		{"three", []string{"A", "B", "C"}, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewCategory(CategoryCfg{Categories: tt.cats})
			min, max := a.Domain()
			if min != tt.wantMin || max != tt.wantMax {
				t.Errorf("Domain() = (%v, %v), want (%v, %v)",
					min, max, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestCategoryTransformInvert(t *testing.T) {
	// Transform maps integer index i to the center of slot i.
	// Invert(center of slot i) returns i + 0.5 (pixel ÷ step, no center
	// adjustment). Verify that Invert falls within the expected slot.
	tests := []struct {
		name       string
		categories []string
		index      float64
		pixMin     float32
		pixMax     float32
	}{
		{"first", []string{"A", "B", "C"}, 0, 0, 300},
		{"last", []string{"A", "B", "C"}, 2, 0, 300},
		{"middle", []string{"A", "B", "C"}, 1, 0, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewCategory(CategoryCfg{Categories: tt.categories})
			px := a.Transform(tt.index, tt.pixMin, tt.pixMax)
			got := a.Invert(px, tt.pixMin, tt.pixMax)
			// Invert of the slot center falls within [index, index+1).
			if got < tt.index || got >= tt.index+1 {
				t.Errorf("Invert of slot center %v: got %v, want in [%v, %v)",
					px, got, tt.index, tt.index+1)
			}
		})
	}
}

func TestCategoryEmptyTransformInvert(t *testing.T) {
	a := NewCategory(CategoryCfg{})
	if px := a.Transform(0, 0, 100); px != 0 {
		t.Errorf("empty Transform = %v, want 0", px)
	}
	if v := a.Invert(50, 0, 100); v != 0 {
		t.Errorf("empty Invert = %v, want 0", v)
	}
}
