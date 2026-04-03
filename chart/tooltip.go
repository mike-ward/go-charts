package chart

import (
	"fmt"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
)

// plotArea describes the pixel bounds and axes of the chart's data
// region. Passed to tooltip helpers to avoid long parameter lists.
type plotArea struct {
	Left, Right, Top, Bottom float32
	XAxis, YAxis             *axis.Linear
}

// nearestXYPoint finds the series/point index and pixel position
// of the data point closest to (mx, my) within snapPx pixels.
// Returns ok=false when no point is within the threshold.
func nearestXYPoint(
	serieses []series.XY, pa plotArea,
	mx, my, snapPx float32,
) (si, pi int, px, py float32, ok bool) {
	best := snapPx * snapPx
	for i, s := range serieses {
		for j, p := range s.Points {
			ppx := pa.XAxis.Transform(p.X, pa.Left, pa.Right)
			ppy := pa.YAxis.Transform(p.Y, pa.Bottom, pa.Top)
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
	serieses []series.XY, pa plotArea,
	mx, my float32,
) {
	si, pi, px, py, ok := nearestXYPoint(
		serieses, pa, mx, my, 20)
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
