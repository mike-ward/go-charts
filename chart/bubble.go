package chart

import (
	"log/slog"
	"math"
	"sort"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// BubbleCfg configures a bubble chart (scatter with sized markers).
type BubbleCfg struct {
	BaseCfg

	// Data
	Series []series.XYZ

	// Axes (optional; auto-created from series bounds when nil)
	XAxis *axis.Linear
	YAxis *axis.Linear

	// Appearance
	MinRadius float32     // minimum marker radius; 0 means default (4)
	MaxRadius float32     // maximum marker radius; 0 means default (30)
	Marker    MarkerShape // default marker shape (circle)
	// Markers overrides the marker shape per series. When
	// Markers[i] is set it takes precedence over Marker for
	// series i. Nil or short slices fall back to Marker.
	Markers []MarkerShape
}

type bubbleView struct {
	cfg         BubbleCfg
	lastVersion uint64
	xAxis       *axis.Linear
	yAxis       *axis.Linear
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	zMin, zMax  float64
	hoverPx     float32
	hoverPy     float32
	hovering    bool
	hidden      map[int]bool
	lastPA      plotArea
	lastLB      legendBounds
	win         *gui.Window
}

// Bubble creates a bubble chart view.
func Bubble(cfg BubbleCfg) gui.View {
	cfg.applyDefaults()
	if cfg.MinRadius == 0 {
		cfg.MinRadius = DefaultBubbleMinRadius
	}
	if cfg.MaxRadius == 0 {
		cfg.MaxRadius = DefaultBubbleMaxRadius
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &bubbleView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (bv *bubbleView) Draw(dc *gui.DrawContext) { bv.draw(dc) }

func (bv *bubbleView) chartTheme() *theme.Theme { return bv.cfg.Theme }

func (bv *bubbleView) Content() []gui.View { return nil }

func (bv *bubbleView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &bv.cfg
	hv := loadHover(w, c.ID,
		&bv.hovering, &bv.hoverPx, &bv.hoverPy)
	var hidV uint64
	bv.hidden, hidV = loadHiddenState(w, c.ID)
	bv.lastLB = loadLegendBounds(w, c.ID)
	bv.win = w
	zv := loadZoomVersion(w, c.ID)
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:            c.ID,
		Sizing:        c.Sizing,
		Width:         width,
		Height:        height,
		Version:       c.Version + hv + hidV + zv,
		Clip:          true,
		OnDraw:        bv.draw,
		OnClick:       bv.internalClick,
		OnHover:       bv.internalHover,
		OnMouseMove:   bv.internalMouseMove,
		OnMouseUp:     bv.internalMouseUp,
		OnMouseLeave:  bv.internalMouseLeave,
		OnMouseScroll: bv.internalScroll,
		OnGesture:     bv.internalGesture,
	}).GenerateLayout(w)
}

func (bv *bubbleView) internalScroll(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if !bv.cfg.EnableZoom {
		return
	}
	handleZoomScroll(w, l, e, bv.cfg.ID, bv.lastPA, true, true)
}

func (bv *bubbleView) internalGesture(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if !bv.cfg.EnableZoom {
		return
	}
	handleZoomGesture(w, l, e, bv.cfg.ID, bv.lastPA, true, true)
}

func (bv *bubbleView) internalClick(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if bv.cfg.EnableZoom && handleDoubleClickCheck(w, l, e, bv.cfg.ID) {
		e.IsHandled = true
		return
	}
	mx := e.MouseX
	my := e.MouseY
	if idx := legendHitTest(bv.lastLB, mx, my); idx >= 0 {
		e.IsHandled = true
		l.Shape.Version = toggleHidden(w, bv.cfg.ID, idx)
		return
	}
	if bv.cfg.OnClick != nil {
		bv.cfg.OnClick(l, e, w)
	}
}

func (bv *bubbleView) internalMouseMove(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if (bv.cfg.EnablePan || bv.cfg.EnableRangeSelect) &&
		handleDragHover(w, l, e, bv.cfg.ID, bv.lastPA,
			bv.cfg.EnablePan, bv.cfg.EnableRangeSelect, true, true) {
		return
	}
}

func (bv *bubbleView) internalMouseUp(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if bv.cfg.EnablePan || bv.cfg.EnableRangeSelect {
		handleDragEnd(w, l, e, bv.cfg.ID, bv.lastPA, true, true)
	}
}

func (bv *bubbleView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if isDragging(w, bv.cfg.ID) {
		return
	}
	e.IsHandled = true
	bv.hoverPx = e.MouseX - l.Shape.X
	bv.hoverPy = e.MouseY - l.Shape.Y
	bv.hovering = true
	saveHover(w, l, bv.cfg.ID, true, bv.hoverPx, bv.hoverPy)
	if legendHitTest(bv.lastLB, bv.hoverPx, bv.hoverPy) >= 0 {
		w.SetMouseCursorPointingHand()
	} else if bv.lastPA.XAxis != nil {
		_, _, _, _, ok := nearestXYZPoint(
			bv.cfg.Series, bv.lastPA, bv.hoverPx, bv.hoverPy,
			bv.cfg.MaxRadius+5)
		if ok {
			w.SetMouseCursorPointingHand()
		} else {
			w.SetMouseCursorArrow()
		}
	}
	if bv.cfg.OnHover != nil {
		bv.cfg.OnHover(l, e, w)
	}
}

func (bv *bubbleView) internalMouseLeave(l *gui.Layout, e *gui.Event, w *gui.Window) {
	e.IsHandled = true
	bv.hovering = false
	saveHover(w, l, bv.cfg.ID, false, 0, 0)
	w.SetMouseCursorArrow()
	if bv.cfg.OnMouseLeave != nil {
		bv.cfg.OnMouseLeave(l, e, w)
	}
}

// updateAxes recomputes axes from config or series bounds.
// Returns false if bounds are invalid.
func (bv *bubbleView) updateAxes() bool {
	cfg := &bv.cfg

	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	bv.zMin, bv.zMax = math.MaxFloat64, -math.MaxFloat64

	for _, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		sx0, sx1, sy0, sy1 := s.Bounds()
		minX = min(minX, sx0)
		maxX = max(maxX, sx1)
		minY = min(minY, sy0)
		maxY = max(maxY, sy1)
		sz0, sz1 := s.ZBounds()
		if finite(sz0) && finite(sz1) {
			bv.zMin = min(bv.zMin, sz0)
			bv.zMax = max(bv.zMax, sz1)
		}
	}

	hasBounds := minX <= maxX
	if hasBounds && (!finite(minX) || !finite(maxX) ||
		!finite(minY) || !finite(maxY)) {
		slog.Warn("non-finite bounds", "chart", cfg.ID)
		return false
	}

	if cfg.XAxis != nil {
		bv.xAxis = cfg.XAxis
		if hasBounds {
			bv.xAxis.SetRange(minX, maxX)
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
		bv.xAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		bv.xAxis.SetRange(minX-xRange*0.05, maxX+xRange*0.05)
	}

	if cfg.YAxis != nil {
		bv.yAxis = cfg.YAxis
		if hasBounds {
			bv.yAxis.SetRange(minY, maxY)
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
		bv.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		bv.yAxis.SetRange(minY-yRange*0.05, maxY+yRange*0.05)
	}
	bv.lastVersion = cfg.Version
	return true
}

func (bv *bubbleView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	cfg := &bv.cfg
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

	if bv.xAxis == nil || cfg.Version != bv.lastVersion {
		if !bv.updateAxes() {
			return
		}
	}

	xAxis := bv.xAxis
	yAxis := bv.yAxis

	zs := loadAndApplyZoom(bv.win, bv.cfg.ID, xAxis, yAxis, true, true)

	left = resolveLeft(ctx, th, left, bottom, top, yAxis)

	bottom = ctx.Height() - resolveBottom(ctx, th,
		maxTickLabelWidth(ctx, xAxis.Ticks(left, right), th.TickStyle),
		cfg.XTickRotation, xAxis.Label())
	bottom -= legendBottomReserve(ctx, th, cfg.LegendPosition, names, left, right)

	bv.yTicks = yAxis.Ticks(bottom, top)
	bv.xTicks = xAxis.Ticks(left, right)

	for _, t := range bv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}
	for _, t := range bv.xTicks {
		ctx.Line(t.Position, top, t.Position, bottom,
			th.GridColor, th.GridWidth)
	}

	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range bv.xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		lw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(t.Position-lw/2, bottom+tickLen+2, t.Label, tickStyle)
	}
	for _, t := range bv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2, t.Label, tickStyle)
	}

	drawXAxisLabel(ctx, xAxis.Label(), th, left, right, bottom)
	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	// Annotations.
	drawAnnotations(ctx, &cfg.Annotations, th,
		plotRect{left, right, top, bottom}, xAxis, yAxis)

	// Cache plot area for cursor hit-testing in hover callback.
	bv.lastPA = plotArea{plotRect{left, right, top, bottom}, xAxis, yAxis}

	// Hover highlight: find nearest series/point.
	hovSI := -1
	hovPI := -1
	var hovPx, hovPy float32
	snapPx := cfg.MaxRadius + 5
	if bv.hovering && xAxis != nil {
		pa := bv.lastPA
		si, pi, px, py, snapOK := nearestXYZPoint(
			cfg.Series, pa, bv.hoverPx, bv.hoverPy, snapPx)
		if snapOK {
			hovSI, hovPI, hovPx, hovPy = si, pi, px, py
		}
	}

	bv.drawBubbles(ctx, cfg, xAxis, yAxis,
		left, right, top, bottom, hovSI, hovPI, hovPx, hovPy)

	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
			Index: i,
		}
	}
	pr := plotRect{left, right, top, bottom}
	bv.lastLB = drawLegend(ctx, entries, th, pr,
		cfg.LegendPosition, bv.hidden)
	saveLegendBounds(bv.win, cfg.ID, bv.lastLB)

	drawSelectionRectIf(ctx, zs, pr, th)

	// Crosshair and tooltip.
	if bv.hovering && bv.xAxis != nil {
		drawCrosshair(ctx, th, bv.hoverPx, bv.hoverPy, pr)
		pa := plotArea{pr, xAxis, yAxis}
		drawXYZTooltip(ctx, th, cfg.Series, pa,
			bv.hoverPx, bv.hoverPy, cfg.MaxRadius+5)
	}
}

// bubbleRadius maps a Z value to a pixel radius using sqrt scaling
// so that marker area is proportional to Z.
func bubbleRadius(z, zMin, zMax float64, minR, maxR float32) float32 {
	if zMax <= zMin {
		return (minR + maxR) / 2
	}
	frac := (z - zMin) / (zMax - zMin)
	frac = max(0, min(1, frac))
	t := math.Sqrt(frac)
	return minR + (maxR-minR)*float32(t)
}

// bubbleMarkerShape returns the marker shape for series index i,
// using Markers[i] when available, otherwise the default Marker.
func bubbleMarkerShape(cfg *BubbleCfg, i int) MarkerShape {
	if i < len(cfg.Markers) {
		return cfg.Markers[i]
	}
	return cfg.Marker
}

// bubbleDrawItem holds the data needed to draw a single bubble,
// used to sort by Z descending so smaller bubbles draw on top.
type bubbleDrawItem struct {
	px, py float32
	radius float32
	color  gui.Color
	shape  MarkerShape
}

// drawBubbles renders bubble markers sorted by Z descending and
// the hovered highlight, skipping points outside the plot area.
func (bv *bubbleView) drawBubbles(
	ctx *render.Context, cfg *BubbleCfg,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	hovSI, hovPI int, hovPx, hovPy float32,
) {
	th := cfg.Theme
	zMin, zMax := bv.zMin, bv.zMax

	// Pre-compute capacity for the items slice.
	totalPts := 0
	for i, s := range cfg.Series {
		if s.Len() > 0 && !bv.hidden[i] {
			totalPts += s.Len()
		}
	}

	// Collect all visible bubbles for Z-sorted drawing.
	items := make([]bubbleDrawItem, 0, totalPts)
	for i, s := range cfg.Series {
		if s.Len() == 0 || bv.hidden[i] {
			continue
		}
		color := seriesColor(s.Color(), i, th.Palette)
		if hovSI >= 0 && i != hovSI {
			color = dimColor(color, HoverDimAlpha)
		}
		shape := bubbleMarkerShape(cfg, i)
		for _, p := range s.Points {
			if !finite(p.X) || !finite(p.Y) || !finite(p.Z) {
				continue
			}
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(p.Y, bottom, top)
			if px < left || px > right || py < top || py > bottom {
				continue
			}
			r := bubbleRadius(p.Z, zMin, zMax,
				cfg.MinRadius, cfg.MaxRadius)
			items = append(items, bubbleDrawItem{px, py, r, color, shape})
		}
	}

	// Sort Z descending (largest first) so smaller bubbles
	// draw on top.
	sort.Slice(items, func(i, j int) bool {
		return items[i].radius > items[j].radius
	})

	for _, it := range items {
		drawMarker(ctx, it.px, it.py, it.radius*2,
			it.shape, it.color)
	}

	// Hover highlight: draw hovered point at 1.3x radius.
	if hovSI >= 0 && hovPI >= 0 && !bv.hidden[hovSI] &&
		hovPx >= left && hovPx <= right &&
		hovPy >= top && hovPy <= bottom {
		hc := seriesColor(cfg.Series[hovSI].Color(), hovSI, th.Palette)
		hovShape := bubbleMarkerShape(cfg, hovSI)
		z := cfg.Series[hovSI].Points[hovPI].Z
		r := bubbleRadius(z, zMin, zMax,
			cfg.MinRadius, cfg.MaxRadius)
		drawMarker(ctx, hovPx, hovPy, r*2*1.3,
			hovShape, hc)
	}
}
