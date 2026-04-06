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

// LineCfg configures a line chart.
type LineCfg struct {
	BaseCfg

	// Data
	Series []series.XY

	// Axes (optional; auto-created from series bounds when nil)
	XAxis axis.Axis
	YAxis axis.Axis

	// Appearance
	LineWidth   float32 // 0 means default (2)
	ShowMarkers bool
	ShowArea    bool // filled area under the line

	// AutoScroll enables smooth scrolling to follow latest
	// data. Typically used with RealTimeSeries.
	AutoScroll bool
	// WindowSize is the visible X-axis range when AutoScroll
	// is enabled. Zero shows all data.
	WindowSize float64
}

type lineView struct {
	cfg         LineCfg
	lastVersion uint64
	xAxis       axis.Axis
	yAxis       axis.Axis
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	ptsBuf      []float32
	clipA       []float32 // scratch for clipConvexToRect
	clipB       []float32
	hoverPx     float32
	hoverPy     float32
	hovering    bool
	hidden      map[int]bool // legend toggle state
	lastPA      plotArea     // cached for cursor hit-testing
	lastLB      legendBounds // cached for legend click
	win         *gui.Window  // set in GenerateLayout for StateMap access
}

// Line creates a line chart view.
func Line(cfg LineCfg) gui.View {
	cfg.applyDefaults()
	if cfg.LineWidth == 0 {
		cfg.LineWidth = DefaultLineWidth
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	if cfg.ShowDataTable {
		return dataTableXY(&cfg.BaseCfg, cfg.Series)
	}
	return &lineView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (lv *lineView) Draw(dc *gui.DrawContext) { lv.draw(dc) }

func (lv *lineView) chartTheme() *theme.Theme { return lv.cfg.Theme }

func (lv *lineView) Content() []gui.View { return nil }

func (lv *lineView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &lv.cfg
	hv := loadHover(w, c.ID,
		&lv.hovering, &lv.hoverPx, &lv.hoverPy)
	var hidV uint64
	lv.hidden, hidV = loadHiddenState(w, c.ID)
	lv.lastLB = loadLegendBounds(w, c.ID)
	lv.win = w
	zv := loadZoomVersion(w, c.ID)
	av := loadAnimVersion(w, c.ID)
	tv := loadTransitionVersion(w, c.ID)
	sv := loadScrollVersion(w, c.ID)
	if c.Animate {
		startEntryAnimation(w, c.ID, c.AnimDuration)
	}
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:            c.ID,
		Sizing:        c.Sizing,
		Width:         width,
		Height:        height,
		Version:       c.Version + hv + hidV + zv + av + tv + sv,
		Clip:          true,
		OnDraw:        lv.draw,
		OnClick:       lv.internalClick,
		OnHover:       lv.internalHover,
		OnMouseMove:   lv.internalMouseMove,
		OnMouseUp:     lv.internalMouseUp,
		OnMouseLeave:  lv.internalMouseLeave,
		OnMouseScroll: lv.internalScroll,
		OnGesture:     lv.internalGesture,
	}).GenerateLayout(w)
}

func (lv *lineView) internalScroll(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if !lv.cfg.EnableZoom {
		return
	}
	handleZoomScroll(w, l, e, lv.cfg.ID, lv.lastPA, true, true)
}

func (lv *lineView) internalGesture(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if !lv.cfg.EnableZoom {
		return
	}
	handleZoomGesture(w, l, e, lv.cfg.ID, lv.lastPA, true, true)
}

func (lv *lineView) internalClick(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if lv.cfg.EnableZoom && handleDoubleClickCheck(w, l, e, lv.cfg.ID) {
		e.IsHandled = true
		return
	}
	mx := e.MouseX
	my := e.MouseY
	if idx := legendHitTest(lv.lastLB, mx, my); idx >= 0 {
		e.IsHandled = true
		l.Shape.Version = toggleHidden(w, lv.cfg.ID, idx)
		return
	}
	if lv.cfg.OnClick != nil {
		lv.cfg.OnClick(l, e, w)
	}
}

func (lv *lineView) internalMouseMove(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if (lv.cfg.EnablePan || lv.cfg.EnableRangeSelect) &&
		handleDragHover(w, l, e, lv.cfg.ID, lv.lastPA,
			lv.cfg.EnablePan, lv.cfg.EnableRangeSelect, true, true) {
		return
	}
}

func (lv *lineView) internalMouseUp(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if lv.cfg.EnablePan || lv.cfg.EnableRangeSelect {
		handleDragEnd(w, l, e, lv.cfg.ID, lv.lastPA, true, true)
	}
}

func (lv *lineView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
	if isDragging(w, lv.cfg.ID) {
		lv.hovering = false
		saveHover(w, l, lv.cfg.ID, false, 0, 0)
		return
	}
	e.IsHandled = true
	lv.hoverPx = e.MouseX - l.Shape.X
	lv.hoverPy = e.MouseY - l.Shape.Y
	lv.hovering = true
	saveHover(w, l, lv.cfg.ID, true, lv.hoverPx, lv.hoverPy)
	if legendHitTest(lv.lastLB, lv.hoverPx, lv.hoverPy) >= 0 {
		w.SetMouseCursorPointingHand()
	} else if lv.lastPA.XAxis != nil {
		_, _, _, _, ok := nearestXYPoint(
			lv.cfg.Series, lv.lastPA, lv.hoverPx, lv.hoverPy, 20)
		if ok {
			w.SetMouseCursorPointingHand()
		} else {
			w.SetMouseCursorArrow()
		}
	}
	if lv.cfg.OnHover != nil {
		lv.cfg.OnHover(l, e, w)
	}
}

func (lv *lineView) internalMouseLeave(l *gui.Layout, e *gui.Event, w *gui.Window) {
	e.IsHandled = true
	lv.hovering = false
	saveHover(w, l, lv.cfg.ID, false, 0, 0)
	w.SetMouseCursorArrow()
	if lv.cfg.OnMouseLeave != nil {
		lv.cfg.OnMouseLeave(l, e, w)
	}
}

// cacheTransitionData stores current Y values and axis bounds
// for future transition animations. Skips while a transition
// is active to preserve old values.
func (lv *lineView) cacheTransitionData() {
	cfg := &lv.cfg
	if cfg.AnimateTransitions &&
		!transitionActive(lv.win, cfg.ID) {
		saveTransitionData(lv.win, cfg.ID,
			snapshotYValues(cfg.Series))
		if lv.xAxis != nil && lv.yAxis != nil {
			xMin, xMax := lv.xAxis.Domain()
			yMin, yMax := lv.yAxis.Domain()
			saveTransitionBounds(lv.win, cfg.ID,
				xMin, xMax, yMin, yMax)
		}
	}
}

// maybeStartTransition starts a transition animation if
// AnimateTransitions is enabled and cfg.Version actually
// changed. Uses LastCfgVer in StateMap to detect real data
// changes (lineView is recreated each frame in immediate-mode
// so struct fields cannot track version across frames).
func (lv *lineView) maybeStartTransition() {
	cfg := &lv.cfg
	if !cfg.AnimateTransitions {
		return
	}
	sm := chartTransitionMap(lv.win)
	ts, _ := sm.Get(cfg.ID)
	if ts.Active || cfg.Version == ts.LastCfgVer {
		return
	}
	ts.LastCfgVer = cfg.Version
	sm.Set(cfg.ID, ts)
	startTransition(lv.win, cfg.ID, cfg.AnimDuration)
}

// applyAutoScroll overrides X-axis domain to follow latest data
// when auto-scroll is enabled and zoom is not active.
func applyAutoScroll(
	w *gui.Window, id string, autoScroll bool,
	windowSize float64, zoomed bool,
	ss []series.XY, xAxis axis.Axis,
) {
	if !autoScroll || windowSize <= 0 || zoomed {
		return
	}
	_, dataXMax, _, _ := seriesBoundsXY(ss)
	updateAutoScroll(w, id, dataXMax, windowSize)
	if xMax, ok := scrollXMax(w, id); ok {
		xAxis.SetRange(xMax-windowSize, xMax)
		xAxis.SetOverrideDomain(true)
	}
}

// updateAxes recomputes axes from config or series bounds.
// Returns false if bounds are invalid (empty or non-finite).
func (lv *lineView) updateAxes() bool {
	cfg := &lv.cfg

	// Always compute bounds from series data so explicit
	// AutoRange axes can be sized correctly.
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
		lv.xAxis = cfg.XAxis
		if hasBounds {
			lv.xAxis.SetRange(minX, maxX)
		}
	} else {
		if !hasBounds {
			slog.Warn("all series empty", "chart", cfg.ID)
			return false
		}
		lv.xAxis = axis.NewLinear(
			axis.LinearCfg{AutoRange: true})
		lv.xAxis.SetRange(minX, maxX)
	}

	if cfg.YAxis != nil {
		lv.yAxis = cfg.YAxis
		if hasBounds {
			lv.yAxis.SetRange(minY, maxY)
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
		minY -= yRange * 0.05
		maxY += yRange * 0.05
		lv.yAxis = axis.NewLinear(
			axis.LinearCfg{AutoRange: true})
		lv.yAxis.SetRange(minY, maxY)
	}
	lv.lastVersion = cfg.Version
	return true
}

func (lv *lineView) draw(dc *gui.DrawContext) {
	updateFPSTracker(lv.win)
	ctx := render.NewContext(dc)
	cfg := &lv.cfg
	th := cfg.Theme

	if len(cfg.Series) == 0 {
		slog.Warn("no series data", "chart", cfg.ID)
		return
	}

	// Chart area inset by theme padding.
	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	// Reserve space for outside legends.
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

	// Title.
	drawTitle(ctx, cfg.Title, th)

	// Recompute axes only when version changes.
	if lv.xAxis == nil || cfg.Version != lv.lastVersion {
		lv.maybeStartTransition()
		if !lv.updateAxes() {
			return
		}
	}

	xAxis := lv.xAxis
	yAxis := lv.yAxis

	// Animated scale transition: lerp axis domain from old to
	// new bounds so the grid moves with the data.
	tp := transitionProgress(lv.win, cfg.ID)
	if tp < 1 {
		oMinX, oMaxX, oMinY, oMaxY, ok :=
			loadTransitionBounds(lv.win, cfg.ID)
		if ok {
			nMinX, nMaxX := xAxis.Domain()
			nMinY, nMaxY := yAxis.Domain()
			lerpAxisRange(xAxis, tp, oMinX, oMaxX, nMinX, nMaxX)
			lerpAxisRange(yAxis, tp, oMinY, oMaxY, nMinY, nMaxY)
		}
	}

	zs := loadAndApplyZoom(lv.win, lv.cfg.ID, xAxis, yAxis, true, true)
	applyAutoScroll(lv.win, cfg.ID, cfg.AutoScroll,
		cfg.WindowSize, zs.Zoomed, cfg.Series, xAxis)

	left = resolveLeft(ctx, th, left, bottom, top, yAxis)

	// Resolve bottom from actual X-axis content.
	bottom = ctx.Height() - resolveBottom(ctx, th,
		maxTickLabelWidth(ctx, xAxis.Ticks(left, right), th.TickStyle),
		cfg.XTickRotation, xAxis.Label())
	bottom -= legendBottomReserve(ctx, th, cfg.LegendPosition, names, left, right)

	// Generate ticks.
	lv.yTicks = yAxis.Ticks(bottom, top)
	lv.xTicks = xAxis.Ticks(left, right)

	// FPS-aware: skip grid when animating under load.
	reduceFPS := animProgress(lv.win, cfg.ID) < 1 &&
		shouldReduceDetail(lv.win)

	// Draw grid lines.
	if !reduceFPS {
		for _, t := range lv.yTicks {
			ctx.Line(left, t.Position, right, t.Position,
				th.GridColor, th.GridWidth)
		}
		for _, t := range lv.xTicks {
			ctx.Line(t.Position, top, t.Position, bottom,
				th.GridColor, th.GridWidth)
		}
	}

	// Draw axes.
	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth) // X
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)     // Y

	// Draw tick marks and labels on axes.
	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range lv.xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		xts := tickStyle
		if cfg.XTickRotation != 0 {
			xts.RotationRadians = cfg.XTickRotation
			ctx.Text(t.Position, bottom+tickLen+2,
				t.Label, xts)
		} else {
			lw := ctx.TextWidth(t.Label, xts)
			ctx.Text(t.Position-lw/2, bottom+tickLen+2,
				t.Label, xts)
		}
	}
	for _, t := range lv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2,
			t.Label, tickStyle)
	}

	// Axis labels.
	drawXAxisLabel(ctx, xAxis.Label(), th, left, right, bottom)
	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	// Annotations.
	drawAnnotations(ctx, &cfg.Annotations, th,
		plotRect{left, right, top, bottom}, xAxis, yAxis)

	// Cache plot area for cursor hit-testing in hover callback.
	lv.lastPA = plotArea{plotRect{left, right, top, bottom}, xAxis, yAxis}

	// Hover highlight: find nearest series/point.
	hovSI, hovPx, hovPy := lv.hoverHighlight(cfg, xAxis)

	progress := animProgress(lv.win, cfg.ID)
	var oldYs [][]float64
	if tp < 1 {
		oldYs, _ = loadTransitionData(lv.win, cfg.ID)
	}
	lv.drawSeries(ctx, cfg, th, xAxis, yAxis,
		left, right, top, bottom, hovSI, progress, tp, oldYs)

	// Enlarged point marker on hovered series (only if inside plot).
	if hovSI >= 0 && !lv.hidden[hovSI] &&
		hovPx >= left && hovPx <= right &&
		hovPy >= top && hovPy <= bottom {
		hc := seriesColor(cfg.Series[hovSI].Color(), hovSI, th.Palette)
		ctx.FilledCircle(hovPx, hovPy, cfg.LineWidth*4, hc)
	}

	// Legend.
	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
			Index: i,
		}
	}
	pr := plotRect{left, right, top, bottom}
	lv.lastLB = drawLegend(ctx, entries, th, pr,
		cfg.LegendPosition, lv.hidden)
	saveLegendBounds(lv.win, cfg.ID, lv.lastLB)

	drawSelectionRectIf(ctx, zs, pr, th)

	// Crosshair and tooltip (skip during FPS reduction).
	if lv.hovering && lv.xAxis != nil && !reduceFPS {
		drawCrosshair(ctx, th, lv.hoverPx, lv.hoverPy, pr)
		pa := lv.lastPA
		drawXYTooltip(ctx, th, cfg.Series, pa,
			lv.hoverPx, lv.hoverPy)
	}

	lv.cacheTransitionData()
}

// hoverHighlight returns the hovered series index and pixel
// coordinates, or -1 if nothing is hovered.
func (lv *lineView) hoverHighlight(
	cfg *LineCfg, xAxis axis.Axis,
) (int, float32, float32) {
	if !lv.hovering || xAxis == nil {
		return -1, 0, 0
	}
	pa := lv.lastPA
	si, _, px, py, ok := nearestXYPoint(
		cfg.Series, pa, lv.hoverPx, lv.hoverPy, 20)
	if !ok {
		return -1, 0, 0
	}
	return si, px, py
}

// drawSeries renders each visible series as polylines with
// optional area fill and markers.
func (lv *lineView) drawSeries(
	ctx *render.Context, cfg *LineCfg, th *theme.Theme,
	xAxis, yAxis axis.Axis,
	left, right, top, bottom float32,
	hovSI int, progress float32,
	tp float32, oldYs [][]float64,
) {
	for i, s := range cfg.Series {
		n := s.Len()
		if n == 0 || lv.hidden[i] {
			continue
		}
		// Resolve old Y values for this series if transitioning.
		var serOldY []float64
		if tp < 1 && i < len(oldYs) {
			serOldY = oldYs[i]
		}
		color := seriesColor(s.Color(), i, th.Palette)
		if hovSI >= 0 && i != hovSI {
			color = dimColor(color, HoverDimAlpha)
		}

		// Build polyline points (flat x,y pairs), reusing buffer.
		needed := n * 2
		if cap(lv.ptsBuf) < needed {
			lv.ptsBuf = make([]float32, 0, needed)
		}
		pts := lv.ptsBuf[:0]
		for j, p := range s.Points[:n] {
			if !finite(p.X) || !finite(p.Y) {
				continue
			}
			y := p.Y
			// Transition: interpolate from old Y to new Y.
			if serOldY != nil && j < len(serOldY) {
				y = lerpFloat64(serOldY[j], p.Y, float64(tp))
			}
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(y, bottom, top)
			// Entry animation: lerp Y from baseline toward
			// actual value for smooth grow-from-zero effect.
			if progress < 1 {
				py = bottom + (py-bottom)*progress
			}
			pts = append(pts, px, py)
		}
		lv.ptsBuf = pts

		// Clip polyline to plot rect for correct boundary
		// intersections.
		clipped := clipPolylineToRect(pts, left, right, top, bottom)

		// Filled area under the line. Clip each quad to the plot
		// rect using Sutherland-Hodgman so fill stays correct when
		// zoomed beyond visible data range.
		if cfg.ShowArea && len(pts) >= 4 {
			fill := gui.RGBA(color.R, color.G, color.B, 40)
			var quad [8]float32
			for k := 0; k < len(pts)-2; k += 2 {
				qy0 := min(pts[k+1], bottom)
				qy1 := min(pts[k+3], bottom)
				if qy0 == bottom && qy1 == bottom {
					continue
				}
				quad[0] = pts[k]
				quad[1] = qy0
				quad[2] = pts[k+2]
				quad[3] = qy1
				quad[4] = pts[k+2]
				quad[5] = bottom
				quad[6] = pts[k]
				quad[7] = bottom
				var clippedQ []float32
				clippedQ, lv.clipA, lv.clipB = clipConvexToRect(
					quad[:], left, right, top, bottom,
					lv.clipA, lv.clipB)
				if clippedQ != nil {
					ctx.FilledPolygon(clippedQ, fill)
				}
			}
		}

		ctx.Polyline(clipped, color, cfg.LineWidth)

		if cfg.ShowMarkers {
			drawLineMarkers(ctx, s.Points[:n], serOldY, tp,
				xAxis, yAxis, left, right, top, bottom,
				cfg.LineWidth, color, progress)
		}
	}
}

// drawLineMarkers renders point markers with transition
// interpolation for the line chart.
func drawLineMarkers(
	ctx *render.Context, pts []series.Point,
	oldY []float64, tp float32,
	xAxis, yAxis axis.Axis,
	left, right, top, bottom, lineWidth float32,
	color gui.Color, progress float32,
) {
	for j, p := range pts {
		if !finite(p.X) || !finite(p.Y) {
			continue
		}
		y := p.Y
		if oldY != nil && j < len(oldY) {
			y = lerpFloat64(oldY[j], p.Y, float64(tp))
		}
		px := xAxis.Transform(p.X, left, right)
		py := yAxis.Transform(y, bottom, top)
		if progress < 1 {
			py = bottom + (py-bottom)*progress
		}
		if px < left || px > right || py < top || py > bottom {
			continue
		}
		ctx.FilledCircle(px, py, lineWidth*2, color)
	}
}
