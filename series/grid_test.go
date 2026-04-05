package series

import (
	"math"
	"testing"
)

func TestNewGrid(t *testing.T) {
	g, err := NewGrid(GridCfg{
		Name:   "test",
		Rows:   []string{"A", "B"},
		Cols:   []string{"X", "Y", "Z"},
		Values: [][]float64{{1, 2, 3}, {4, 5, 6}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Name() != "test" {
		t.Errorf("Name() = %q, want %q", g.Name(), "test")
	}
	if g.Len() != 6 {
		t.Errorf("Len() = %d, want 6", g.Len())
	}
	if g.NumRows() != 2 || g.NumCols() != 3 {
		t.Errorf("dims = %dx%d, want 2x3", g.NumRows(), g.NumCols())
	}
	if v := g.At(1, 2); v != 6 {
		t.Errorf("At(1,2) = %v, want 6", v)
	}
}

func TestNewGridRowMismatch(t *testing.T) {
	_, err := NewGrid(GridCfg{
		Rows:   []string{"A"},
		Cols:   []string{"X"},
		Values: [][]float64{{1}, {2}},
	})
	if err == nil {
		t.Error("expected error for row count mismatch")
	}
}

func TestNewGridColMismatch(t *testing.T) {
	_, err := NewGrid(GridCfg{
		Rows:   []string{"A", "B"},
		Cols:   []string{"X", "Y"},
		Values: [][]float64{{1, 2}, {3}},
	})
	if err == nil {
		t.Error("expected error for col count mismatch")
	}
}

func TestGridBounds(t *testing.T) {
	g, _ := NewGrid(GridCfg{
		Rows:   []string{"A", "B"},
		Cols:   []string{"X", "Y"},
		Values: [][]float64{{3, -1}, {7, 2}},
	})
	vMin, vMax := g.Bounds()
	if vMin != -1 || vMax != 7 {
		t.Errorf("Bounds() = (%v, %v), want (-1, 7)", vMin, vMax)
	}
}

func TestGridBoundsSkipsNaN(t *testing.T) {
	nan := math.NaN()
	g, _ := NewGrid(GridCfg{
		Rows:   []string{"A", "B"},
		Cols:   []string{"X", "Y"},
		Values: [][]float64{{nan, 5}, {3, nan}},
	})
	vMin, vMax := g.Bounds()
	if vMin != 3 || vMax != 5 {
		t.Errorf("Bounds() = (%v, %v), want (3, 5)", vMin, vMax)
	}
}

func TestGridBoundsAllNaN(t *testing.T) {
	nan := math.NaN()
	g, _ := NewGrid(GridCfg{
		Rows:   []string{"A"},
		Cols:   []string{"X"},
		Values: [][]float64{{nan}},
	})
	vMin, vMax := g.Bounds()
	if vMin != 0 || vMax != 0 {
		t.Errorf("Bounds() = (%v, %v), want (0, 0)", vMin, vMax)
	}
}

func TestGridBoundsEmpty(t *testing.T) {
	g := Grid{}
	vMin, vMax := g.Bounds()
	if vMin != 0 || vMax != 0 {
		t.Errorf("Bounds() = (%v, %v), want (0, 0)", vMin, vMax)
	}
}

func TestGridFromFunc(t *testing.T) {
	rows := []string{"R0", "R1"}
	cols := []string{"C0", "C1", "C2"}
	g := GridFromFunc("fn", rows, cols, func(r, c int) float64 {
		return float64(r*10 + c)
	})
	if g.At(1, 2) != 12 {
		t.Errorf("At(1,2) = %v, want 12", g.At(1, 2))
	}
	if g.NumRows() != 2 || g.NumCols() != 3 {
		t.Errorf("dims = %dx%d, want 2x3", g.NumRows(), g.NumCols())
	}
}

func TestGridIsNaN(t *testing.T) {
	g, _ := NewGrid(GridCfg{
		Rows:   []string{"A"},
		Cols:   []string{"X", "Y"},
		Values: [][]float64{{math.NaN(), 5}},
	})
	if !g.IsNaN(0, 0) {
		t.Error("IsNaN(0,0) = false, want true")
	}
	if g.IsNaN(0, 1) {
		t.Error("IsNaN(0,1) = true, want false")
	}
}

func TestGridString(t *testing.T) {
	g, _ := NewGrid(GridCfg{
		Name:   "m",
		Rows:   []string{"A", "B"},
		Cols:   []string{"X"},
		Values: [][]float64{{1}, {2}},
	})
	got := g.String()
	want := `Grid{"m", 2x1}`
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
