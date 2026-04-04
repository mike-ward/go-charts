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
	XAxis *axis.Linear
	YAxis *axis.Linear

	// Appearance
	LineWidth   float32 // 0 means default (2)
	ShowMarkers bool
	ShowArea    bool // filled area under the line
}

type lineView struct {
	cfg         LineCfg
	lastVersion uint64
	xAxis       *axis.Linear
	yAxis       *axis.Linear
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	ptsBuf      []float32
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
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version + hv + hidV,
		Clip:         true,
		OnDraw:       lv.draw,
		OnClick:      lv.internalClick,
		OnHover:      lv.internalHover,
		OnMouseLeave: lv.internalMouseLeave,
	}).GenerateLayout(w)
}

func (lv *lineView) internalClick(l *gui.Layout, e *gui.Event, w *gui.Window) {
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

func (lv *lineView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
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
		if !lv.updateAxes() {
			return
		}
	}

	xAxis := lv.xAxis
	yAxis := lv.yAxis

	left = resolveLeft(ctx, th, left, bottom, top, yAxis)

	// Resolve bottom from actual X-axis content.
	bottom = ctx.Height() - resolveBottom(ctx, th,
		maxTickLabelWidth(ctx, xAxis.Ticks(left, right), th.TickStyle),
		cfg.XTickRotation, xAxis.Label())
	bottom -= legendBottomReserve(ctx, th, cfg.LegendPosition, names, left, right)

	// Generate ticks.
	lv.yTicks = yAxis.Ticks(bottom, top)
	lv.xTicks = xAxis.Ticks(left, right)

	// Draw grid lines.
	for _, t := range lv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}
	for _, t := range lv.xTicks {
		ctx.Line(t.Position, top, t.Position, bottom,
			th.GridColor, th.GridWidth)
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

	// Cache plot area for cursor hit-testing in hover callback.
	lv.lastPA = plotArea{plotRect{left, right, top, bottom}, xAxis, yAxis}

	// Hover highlight: find nearest series/point.
	hovSI := -1
	var hovPx, hovPy float32
	if lv.hovering && xAxis != nil {
		pa := lv.lastPA
		si, _, px, py, snapOK := nearestXYPoint(
			cfg.Series, pa, lv.hoverPx, lv.hoverPy, 20)
		if snapOK {
			hovSI, hovPx, hovPy = si, px, py
		}
	}

	lv.drawSeries(ctx, cfg, th, xAxis, yAxis,
		left, right, top, bottom, hovSI)

	// Enlarged point marker on hovered series.
	if hovSI >= 0 && !lv.hidden[hovSI] {
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

	// Crosshair and tooltip.
	if lv.hovering && lv.xAxis != nil {
		drawCrosshair(ctx, th, lv.hoverPx, lv.hoverPy, pr)
		pa := lv.lastPA
		drawXYTooltip(ctx, th, cfg.Series, pa,
			lv.hoverPx, lv.hoverPy)
	}
}

// drawSeries renders each visible series as polylines with
// optional area fill and markers.
func (lv *lineView) drawSeries(
	ctx *render.Context, cfg *LineCfg, th *theme.Theme,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	hovSI int,
) {
	for i, s := range cfg.Series {
		if s.Len() == 0 || lv.hidden[i] {
			continue
		}
		color := seriesColor(s.Color(), i, th.Palette)
		if hovSI >= 0 && i != hovSI {
			color = dimColor(color, HoverDimAlpha)
		}

		// Build polyline points (flat x,y pairs), reusing buffer.
		needed := s.Len() * 2
		if cap(lv.ptsBuf) < needed {
			lv.ptsBuf = make([]float32, 0, needed)
		}
		pts := lv.ptsBuf[:0]
		for _, p := range s.Points {
			if !finite(p.X) || !finite(p.Y) {
				continue
			}
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(p.Y, bottom, top)
			pts = append(pts, px, py)
		}
		lv.ptsBuf = pts

		// Filled area under the line.
		if cfg.ShowArea && len(pts) >= 4 {
			fill := gui.RGBA(color.R, color.G, color.B, 40)
			var quad [8]float32
			for k := 0; k < len(pts)-2; k += 2 {
				quad[0] = pts[k]
				quad[1] = pts[k+1]
				quad[2] = pts[k+2]
				quad[3] = pts[k+3]
				quad[4] = pts[k+2]
				quad[5] = bottom
				quad[6] = pts[k]
				quad[7] = bottom
				ctx.FilledPolygon(quad[:], fill)
			}
		}

		ctx.Polyline(pts, color, cfg.LineWidth)

		// Markers at each data point.
		if cfg.ShowMarkers {
			for j := 0; j < len(pts); j += 2 {
				ctx.FilledCircle(pts[j], pts[j+1], cfg.LineWidth*2, color)
			}
		}
	}
}
