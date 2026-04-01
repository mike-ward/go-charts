package series

import "github.com/mike-ward/go-gui/gui"

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

// Bounds returns the min/max X and Y values.
func (s XY) Bounds() (minX, maxX, minY, maxY float64) {
	if len(s.Points) == 0 {
		return
	}
	minX, maxX = s.Points[0].X, s.Points[0].X
	minY, maxY = s.Points[0].Y, s.Points[0].Y
	for _, p := range s.Points[1:] {
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
