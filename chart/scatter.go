package chart

import (
	"log/slog"
	"math"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// MarkerShape controls the shape of scatter plot markers.
type MarkerShape uint8

// MarkerShape constants.
const (
	MarkerCircle MarkerShape = iota
	MarkerSquare
	MarkerTriangle
	MarkerDiamond
	MarkerCross
)

// ScatterCfg configures a scatter plot.
type ScatterCfg struct {
	BaseCfg

	// Data
	Series []series.XY

	// Axes (optional; auto-created from series bounds when nil)
	XAxis *axis.Linear
	YAxis *axis.Linear

	// Appearance
	MarkerSize float32 // 0 means default (6)
	Marker     MarkerShape
}

type scatterView struct {
	cfg         ScatterCfg
	lastVersion uint64
	xAxis       *axis.Linear
	yAxis       *axis.Linear
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	hoverPx     float32
	hoverPy     float32
	hovering    bool
	hidden      map[int]bool // legend toggle state
	lastPA      plotArea     // cached for cursor hit-testing
	lastLB      legendBounds // cached for legend click
	win         *gui.Window
}

// Scatter creates a scatter plot view.
func Scatter(cfg ScatterCfg) gui.View {
	cfg.applyDefaults()
	if cfg.MarkerSize == 0 {
		cfg.MarkerSize = DefaultMarkerSize
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &scatterView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (sv *scatterView) Draw(dc *gui.DrawContext) { sv.draw(dc) }

func (sv *scatterView) chartTheme() *theme.Theme { return sv.cfg.Theme }

func (sv *scatterView) Content() []gui.View { return nil }

func (sv *scatterView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &sv.cfg
	hv := loadHover(w, c.ID,
		&sv.hovering, &sv.hoverPx, &sv.hoverPy)
	var hidV uint64
	sv.hidden, hidV = loadHiddenState(w, c.ID)
	sv.lastLB = loadLegendBounds(w, c.ID)
	sv.win = w
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version + hv + hidV,
		Clip:         true,
		OnDraw:       sv.draw,
		OnClick:      sv.internalClick,
		OnHover:      sv.internalHover,
		OnMouseLeave: sv.internalMouseLeave,
	}).GenerateLayout(w)
}

func (sv *scatterView) internalClick(l *gui.Layout, e *gui.Event, w *gui.Window) {
	mx := e.MouseX
	my := e.MouseY
	if idx := legendHitTest(sv.lastLB, mx, my); idx >= 0 {
		e.IsHandled = true
		l.Shape.Version = toggleHidden(w, sv.cfg.ID, idx)
		return
	}
	if sv.cfg.OnClick != nil {
		sv.cfg.OnClick(l, e, w)
	}
}

func (sv *scatterView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
	e.IsHandled = true
	sv.hoverPx = e.MouseX - l.Shape.X
	sv.hoverPy = e.MouseY - l.Shape.Y
	sv.hovering = true
	saveHover(w, l, sv.cfg.ID, true, sv.hoverPx, sv.hoverPy)
	if legendHitTest(sv.lastLB, sv.hoverPx, sv.hoverPy) >= 0 {
		w.SetMouseCursorPointingHand()
	} else if sv.lastPA.XAxis != nil {
		_, _, _, _, ok := nearestXYPoint(
			sv.cfg.Series, sv.lastPA, sv.hoverPx, sv.hoverPy, 20)
		if ok {
			w.SetMouseCursorPointingHand()
		} else {
			w.SetMouseCursorArrow()
		}
	}
	if sv.cfg.OnHover != nil {
		sv.cfg.OnHover(l, e, w)
	}
}

func (sv *scatterView) internalMouseLeave(l *gui.Layout, e *gui.Event, w *gui.Window) {
	e.IsHandled = true
	sv.hovering = false
	saveHover(w, l, sv.cfg.ID, false, 0, 0)
	w.SetMouseCursorArrow()
	if sv.cfg.OnMouseLeave != nil {
		sv.cfg.OnMouseLeave(l, e, w)
	}
}

// updateAxes recomputes axes from config or series bounds.
// Returns false if bounds are invalid.
func (sv *scatterView) updateAxes() bool {
	cfg := &sv.cfg

	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		sx0, sx1, sy0, sy1 := s.Bounds()
		minX = min(minX, sx0)
		maxX = max(maxX, sx1)
		minY = min(minY, sy0)
		maxY = max(maxY, sy1)
	}

	hasBounds := minX <= maxX
	if hasBounds && (!finite(minX) || !finite(maxX) ||
		!finite(minY) || !finite(maxY)) {
		slog.Warn("non-finite bounds", "chart", cfg.ID)
		return false
	}

	if cfg.XAxis != nil {
		sv.xAxis = cfg.XAxis
		if hasBounds {
			sv.xAxis.SetRange(minX, maxX)
		}
	} else {
		if !hasBounds {
			slog.Warn("all series empty", "chart", cfg.ID)
			return false
		}
		xRange := maxX - minX
		if xRange == 0 {
			xRange = 1
		}
		sv.xAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		sv.xAxis.SetRange(minX-xRange*0.05, maxX+xRange*0.05)
	}

	if cfg.YAxis != nil {
		sv.yAxis = cfg.YAxis
		if hasBounds {
			sv.yAxis.SetRange(minY, maxY)
		}
	} else {
		if !hasBounds {
			slog.Warn("all series empty", "chart", cfg.ID)
			return false
		}
		yRange := maxY - minY
		if yRange == 0 {
			yRange = 1
		}
		sv.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		sv.yAxis.SetRange(minY-yRange*0.05, maxY+yRange*0.05)
	}
	sv.lastVersion = cfg.Version
	return true
}

func (sv *scatterView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	cfg := &sv.cfg
	th := cfg.Theme

	if len(cfg.Series) == 0 {
		slog.Warn("no series data", "chart", cfg.ID)
		return
	}

	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	names := make([]string, len(cfg.Series))
	for i, s := range cfg.Series {
		names[i] = s.Name()
	}
	right -= legendRightReserve(ctx, th, cfg.LegendPosition, names)
	top += legendTopReserve(ctx, th, cfg.LegendPosition, names, left, right)

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	drawTitle(ctx, cfg.Title, th)

	if sv.xAxis == nil || cfg.Version != sv.lastVersion {
		if !sv.updateAxes() {
			return
		}
	}

	xAxis := sv.xAxis
	yAxis := sv.yAxis

	left = resolveLeft(ctx, th, left, bottom, top, yAxis)

	bottom = ctx.Height() - resolveBottom(ctx, th,
		maxTickLabelWidth(ctx, xAxis.Ticks(left, right), th.TickStyle),
		cfg.XTickRotation, xAxis.Label())
	bottom -= legendBottomReserve(ctx, th, cfg.LegendPosition, names, left, right)

	sv.yTicks = yAxis.Ticks(bottom, top)
	sv.xTicks = xAxis.Ticks(left, right)

	for _, t := range sv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}
	for _, t := range sv.xTicks {
		ctx.Line(t.Position, top, t.Position, bottom,
			th.GridColor, th.GridWidth)
	}

	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range sv.xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		lw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(t.Position-lw/2, bottom+tickLen+2, t.Label, tickStyle)
	}
	for _, t := range sv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2, t.Label, tickStyle)
	}

	drawXAxisLabel(ctx, xAxis.Label(), th, left, right, bottom)
	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	// Cache plot area for cursor hit-testing in hover callback.
	sv.lastPA = plotArea{left, right, top, bottom, xAxis, yAxis}

	// Hover highlight: find nearest series/point.
	hovSI := -1
	var hovPx, hovPy float32
	if sv.hovering && xAxis != nil {
		pa := sv.lastPA
		si, _, px, py, snapOK := nearestXYPoint(
			cfg.Series, pa, sv.hoverPx, sv.hoverPy, 20)
		if snapOK {
			hovSI, hovPx, hovPy = si, px, py
		}
	}

	for i, s := range cfg.Series {
		if s.Len() == 0 || sv.hidden[i] {
			continue
		}
		color := seriesColor(s.Color(), i, th.Palette)
		if hovSI >= 0 && i != hovSI {
			color = dimColor(color, HoverDimAlpha)
		}
		for _, p := range s.Points {
			if !finite(p.X) || !finite(p.Y) {
				continue
			}
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(p.Y, bottom, top)
			drawMarker(ctx, px, py, cfg.MarkerSize, cfg.Marker, color)
		}
	}

	// Enlarged marker on hovered series/point.
	if hovSI >= 0 && !sv.hidden[hovSI] {
		hc := seriesColor(cfg.Series[hovSI].Color(), hovSI, th.Palette)
		drawMarker(ctx, hovPx, hovPy, cfg.MarkerSize*2, cfg.Marker, hc)
	}

	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
			Index: i,
		}
	}
	sv.lastLB = drawLegend(ctx, entries, th, left, right, top, bottom,
		cfg.LegendPosition, sv.hidden)
	saveLegendBounds(sv.win, cfg.ID, sv.lastLB)

	// Crosshair and tooltip.
	if sv.hovering && sv.xAxis != nil {
		drawCrosshair(ctx, th, sv.hoverPx, sv.hoverPy,
			left, right, top, bottom)
		pa := plotArea{left, right, top, bottom, xAxis, yAxis}
		drawXYTooltip(ctx, th, cfg.Series, pa,
			sv.hoverPx, sv.hoverPy)
	}
}

// drawMarker renders a single marker at (cx, cy) with the given size and shape.
func drawMarker(
	ctx *render.Context, cx, cy, size float32,
	shape MarkerShape, color gui.Color,
) {
	h := size / 2
	switch shape {
	case MarkerSquare:
		ctx.FilledRect(cx-h, cy-h, size, size, color)
	case MarkerTriangle:
		// Equilateral triangle pointing up.
		pts := [6]float32{
			cx, cy - h,
			cx + h, cy + h,
			cx - h, cy + h,
		}
		ctx.FilledPolygon(pts[:], color)
	case MarkerDiamond:
		pts := [8]float32{
			cx, cy - h,
			cx + h, cy,
			cx, cy + h,
			cx - h, cy,
		}
		ctx.FilledPolygon(pts[:], color)
	case MarkerCross:
		w := size / 4
		ctx.Line(cx-h, cy, cx+h, cy, color, w)
		ctx.Line(cx, cy-h, cx, cy+h, color, w)
	default: // MarkerCircle
		ctx.FilledCircle(cx, cy, h, color)
	}
}
