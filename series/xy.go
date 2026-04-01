package series

import (
	"math"

	"github.com/mike-ward/go-gui/gui"
)

// Point represents a single (X, Y) data point.
type Point struct {
	X, Y float64
}

// XY is a series of (X, Y) data points.
type XY struct {
	name   string
	color  gui.Color
	Points []Point
}

// XYCfg configures an XY series.
type XYCfg struct {
	Name   string
	Color  gui.Color
	Points []Point
}

// NewXY creates a new XY data series.
func NewXY(cfg XYCfg) XY {
	return XY{
		name:   cfg.Name,
		color:  cfg.Color,
		Points: cfg.Points,
	}
}

// Name implements Series.
func (s XY) Name() string { return s.name }

// Len implements Series.
func (s XY) Len() int { return len(s.Points) }

// Color implements Series.
func (s XY) Color() gui.Color { return s.color }

// finite reports whether v is neither NaN nor +/-Inf.
func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// Bounds returns the min/max X and Y values. Non-finite points
// (NaN, +/-Inf) are skipped. If no finite points exist, all
// returned values are zero.
func (s XY) Bounds() (minX, maxX, minY, maxY float64) {
	// Find first finite point to seed min/max.
	i := 0
	for i < len(s.Points) {
		p := s.Points[i]
		if finite(p.X) && finite(p.Y) {
			break
		}
		i++
	}
	if i >= len(s.Points) {
		return // no finite points
	}
	minX, maxX = s.Points[i].X, s.Points[i].X
	minY, maxY = s.Points[i].Y, s.Points[i].Y
	for _, p := range s.Points[i+1:] {
		if !finite(p.X) || !finite(p.Y) {
			continue
		}
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return
}
