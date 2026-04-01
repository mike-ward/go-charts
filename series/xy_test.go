package series

import (
	"testing"

	"github.com/mike-ward/go-gui/gui"
)

func TestXYBounds(t *testing.T) {
	s := NewXY(XYCfg{
		Name:  "test",
		Color: gui.Blue,
		Points: []Point{
			{X: 1, Y: 10},
			{X: 5, Y: 2},
			{X: 3, Y: 8},
		},
	})
	minX, maxX, minY, maxY := s.Bounds()
	if minX != 1 || maxX != 5 || minY != 2 || maxY != 10 {
		t.Errorf("Bounds = (%v,%v,%v,%v), want (1,5,2,10)",
			minX, maxX, minY, maxY)
	}
}

func TestXYEmpty(t *testing.T) {
	s := NewXY(XYCfg{Name: "empty"})
	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}
