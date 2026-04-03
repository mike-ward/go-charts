package chart

import (
	"fmt"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
)

// nearestXYPoint finds the series/point index and pixel position of
// the data point closest to (mx, my) within snapPx pixels. Returns
// ok=false when no point is within the threshold.
func nearestXYPoint(
	serieses []series.XY,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	mx, my float32,
	snapPx float32,
) (si, pi int, px, py float32, ok bool) {
	best := snapPx * snapPx
	for i, s := range serieses {
		for j, p := range s.Points {
			ppx := xAxis.Transform(p.X, left, right)
			ppy := yAxis.Transform(p.Y, bottom, top)
			dx := ppx - mx
			dy := ppy - my
			d2 := dx*dx + dy*dy
			if d2 < best {
				best = d2
				si, pi, px, py = i, j, ppx, ppy
				ok = true
			}
		}
	}
	return
}

// drawXYTooltip draws a tooltip for the nearest XY data point.
// Shared by line, scatter, and area charts.
func drawXYTooltip(
	ctx *render.Context, th *theme.Theme,
	serieses []series.XY,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	mx, my float32,
) {
	si, pi, px, py, ok := nearestXYPoint(
		serieses, xAxis, yAxis,
		left, right, top, bottom,
		mx, my, 20)
	if !ok {
		return
	}
	s := serieses[si]
	p := s.Points[pi]
	var label string
	if s.Name() != "" {
		label = fmt.Sprintf(
			"%s\nX: %g\nY: %g", s.Name(), p.X, p.Y)
	} else {
		label = fmt.Sprintf("X: %g\nY: %g", p.X, p.Y)
	}
	drawTooltip(ctx, px, py, label, th)
}
