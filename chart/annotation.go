package chart

import (
	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// AnnotationAxis selects which axis an annotation is associated with.
type AnnotationAxis uint8

const (
	// AnnotationX draws a vertical line or X-range region.
	AnnotationX AnnotationAxis = iota
	// AnnotationY draws a horizontal line or Y-range region.
	AnnotationY
)

// LineAnnotation draws a reference line spanning the plot area at
// a data-coordinate value. Horizontal for AnnotationY, vertical
// for AnnotationX.
type LineAnnotation struct {
	Axis    AnnotationAxis
	Value   float64
	Color   gui.Color // zero → theme GridColor
	Width   float32   // zero → DefaultAnnotationLineWidth
	DashLen float32   // zero → solid line
	GapLen  float32
	Label   string
}

// TextAnnotation draws a text label at a data-coordinate position.
type TextAnnotation struct {
	X, Y     float64
	Text     string
	Color    gui.Color // zero → theme TickStyle color
	FontSize float32   // zero → theme TickStyle size
}

// RegionAnnotation draws a shaded rectangle between two
// data-coordinate values on one axis, spanning the full extent
// of the other axis.
type RegionAnnotation struct {
	Axis  AnnotationAxis
	Min   float64
	Max   float64
	Color gui.Color // zero → semi-transparent gray
	Label string
}

// Annotations groups all annotation types for a chart.
// Zero value is empty and results in no drawing.
type Annotations struct {
	Lines   []LineAnnotation
	Texts   []TextAnnotation
	Regions []RegionAnnotation
}

// empty reports whether there are no annotations to draw.
func (a *Annotations) empty() bool {
	return len(a.Lines) == 0 &&
		len(a.Texts) == 0 &&
		len(a.Regions) == 0
}

// drawAnnotations renders annotations in the plot area. Labels
// for line annotations render just outside the plot boundary.
// Either axis may be nil; annotations referencing a nil axis
// are skipped.
func drawAnnotations(
	ctx *render.Context, ann *Annotations, th *theme.Theme,
	pr plotRect, xAxis, yAxis axis.Axis,
) {
	if ann.empty() {
		return
	}
	for i := range ann.Regions {
		drawRegionAnnotation(ctx, &ann.Regions[i], th, pr, xAxis, yAxis)
	}
	for i := range ann.Lines {
		drawLineAnnotation(ctx, &ann.Lines[i], th, pr, xAxis, yAxis)
	}
	for i := range ann.Texts {
		drawTextAnnotation(ctx, &ann.Texts[i], th, pr, xAxis, yAxis)
	}
}

func drawRegionAnnotation(
	ctx *render.Context, r *RegionAnnotation, th *theme.Theme,
	pr plotRect, xAxis, yAxis axis.Axis,
) {
	color := r.Color
	if !color.IsSet() {
		color = gui.RGBA(128, 128, 128, 30)
	}

	if !finite(r.Min) || !finite(r.Max) {
		return
	}

	var x, y, w, h float32
	switch r.Axis {
	case AnnotationX:
		if xAxis == nil {
			return
		}
		x0 := xAxis.Transform(r.Min, pr.Left, pr.Right)
		x1 := xAxis.Transform(r.Max, pr.Left, pr.Right)
		xLo := min(x0, x1)
		xHi := max(x0, x1)
		xLo = clampF(xLo, pr.Left, pr.Right)
		xHi = clampF(xHi, pr.Left, pr.Right)
		x, y, w, h = xLo, pr.Top, xHi-xLo, pr.Bottom-pr.Top

	case AnnotationY:
		if yAxis == nil {
			return
		}
		y0 := yAxis.Transform(r.Min, pr.Bottom, pr.Top)
		y1 := yAxis.Transform(r.Max, pr.Bottom, pr.Top)
		// Y pixel coords are inverted (top < bottom).
		yLo := min(y0, y1)
		yHi := max(y0, y1)
		yLo = clampF(yLo, pr.Top, pr.Bottom)
		yHi = clampF(yHi, pr.Top, pr.Bottom)
		x, y, w, h = pr.Left, yLo, pr.Right-pr.Left, yHi-yLo
	}

	if w <= 0 || h <= 0 {
		return
	}
	ctx.FilledRect(x, y, w, h, color)

	if r.Label != "" {
		style := th.TickStyle
		tw := ctx.TextWidth(r.Label, style)
		fh := ctx.FontHeight(style)
		lx := x + (w-tw)/2
		ly := y + (h-fh)/2
		ctx.Text(lx, ly, r.Label, style)
	}
}

func drawLineAnnotation(
	ctx *render.Context, la *LineAnnotation, th *theme.Theme,
	pr plotRect, xAxis, yAxis axis.Axis,
) {
	color := la.Color
	if !color.IsSet() {
		color = th.GridColor
	}
	width := la.Width
	if width <= 0 {
		width = DefaultAnnotationLineWidth
	}

	if !finite(la.Value) {
		return
	}

	switch la.Axis {
	case AnnotationX:
		if xAxis == nil {
			return
		}
		px := xAxis.Transform(la.Value, pr.Left, pr.Right)
		if px < pr.Left || px > pr.Right {
			return
		}
		if la.DashLen > 0 && la.GapLen > 0 {
			ctx.DashedLine(px, pr.Top, px, pr.Bottom,
				color, width, la.DashLen, la.GapLen)
		} else {
			ctx.Line(px, pr.Top, px, pr.Bottom, color, width)
		}
		if la.Label != "" {
			style := th.TickStyle
			ctx.Text(px+4, pr.Top+2, la.Label, style)
		}

	case AnnotationY:
		if yAxis == nil {
			return
		}
		py := yAxis.Transform(la.Value, pr.Bottom, pr.Top)
		if py < pr.Top || py > pr.Bottom {
			return
		}
		if la.DashLen > 0 && la.GapLen > 0 {
			ctx.DashedLine(pr.Left, py, pr.Right, py,
				color, width, la.DashLen, la.GapLen)
		} else {
			ctx.Line(pr.Left, py, pr.Right, py, color, width)
		}
		if la.Label != "" {
			style := th.TickStyle
			tw := ctx.TextWidth(la.Label, style)
			fh := ctx.FontHeight(style)
			ctx.Text(pr.Right-tw-4, py-fh-2, la.Label, style)
		}
	}
}

func drawTextAnnotation(
	ctx *render.Context, ta *TextAnnotation, th *theme.Theme,
	pr plotRect, xAxis, yAxis axis.Axis,
) {
	if xAxis == nil || yAxis == nil {
		return
	}
	if !finite(ta.X) || !finite(ta.Y) {
		return
	}
	px := xAxis.Transform(ta.X, pr.Left, pr.Right)
	py := yAxis.Transform(ta.Y, pr.Bottom, pr.Top)
	if px < pr.Left || px > pr.Right || py < pr.Top || py > pr.Bottom {
		return
	}
	style := th.TickStyle
	if ta.Color.IsSet() {
		style.Color = ta.Color
	}
	if ta.FontSize > 0 {
		style.Size = ta.FontSize
	}
	ctx.Text(px, py, ta.Text, style)
}

// clampF restricts v to [lo, hi].
func clampF(v, lo, hi float32) float32 {
	return max(lo, min(v, hi))
}
