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

// AreaCfg configures an area chart.
type AreaCfg struct {
	BaseCfg

	// Data
	Series []series.XY

	// Axes (optional; auto-created from series bounds when nil)
	XAxis *axis.Linear
	YAxis *axis.Linear

	// Appearance
	Stacked   bool
	LineWidth float32 // 0 means default (2)
	Opacity   float32 // fill opacity 0-1; 0 means default (0.3)
}

type areaView struct {
	cfg        AreaCfg
	lastVersion uint64
	xAxis       *axis.Linear
	yAxis       *axis.Linear
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	ptsBuf      []float32
	prevPtsBuf  []float32
	hoverPx     float32
	hoverPy     float32
	hovering    bool
}

// Area creates an area chart view.
func Area(cfg AreaCfg) gui.View {
	cfg.applyDefaults()
	if cfg.LineWidth == 0 {
		cfg.LineWidth = DefaultLineWidth
	}
	if cfg.Opacity == 0 {
		cfg.Opacity = DefaultAreaOpacity
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &areaView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (av *areaView) Draw(dc *gui.DrawContext) { av.draw(dc) }

func (av *areaView) chartTheme() *theme.Theme { return av.cfg.Theme }

func (av *areaView) Content() []gui.View { return nil }

func (av *areaView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &av.cfg
	hv := loadHover(w, c.ID,
		&av.hovering, &av.hoverPx, &av.hoverPy)
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version + hv,
		Clip:         true,
		OnDraw:       av.draw,
		OnClick:      c.OnClick,
		OnHover:      av.internalHover,
		OnMouseLeave: av.internalMouseLeave,
	}).GenerateLayout(w)
}

func (av *areaView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
	av.hoverPx = e.MouseX - l.Shape.X
	av.hoverPy = e.MouseY - l.Shape.Y
	av.hovering = true
	saveHover(w, l, av.cfg.ID, true, av.hoverPx, av.hoverPy)
	if av.cfg.OnHover != nil {
		av.cfg.OnHover(l, e, w)
	}
}

func (av *areaView) internalMouseLeave(l *gui.Layout, e *gui.Event, w *gui.Window) {
	av.hovering = false
	saveHover(w, l, av.cfg.ID, false, 0, 0)
	if av.cfg.OnMouseLeave != nil {
		av.cfg.OnMouseLeave(l, e, w)
	}
}

// updateAxes recomputes axes from config or series bounds.
// Returns false if bounds are invalid (empty or non-finite).
func (av *areaView) updateAxes() bool {
	cfg := &av.cfg

	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		sx0, sx1, sy0, sy1 := s.Bounds()
		minX = min(minX, sx0)
		maxX = max(maxX, sx1)
		if !cfg.Stacked {
			minY = min(minY, sy0)
			maxY = max(maxY, sy1)
		}
	}

	// For stacked mode, Y max is the maximum cumulative sum per index.
	if cfg.Stacked {
		refLen := 0
		for _, s := range cfg.Series {
			if s.Len() > 0 {
				refLen = s.Len()
				break
			}
		}
		if refLen > 0 {
			sums := make([]float64, refLen)
			for _, s := range cfg.Series {
				n := min(s.Len(), refLen)
				for j := range n {
					sums[j] += s.Points[j].Y
					maxY = max(maxY, sums[j])
				}
			}
			minY = 0
		}
	}

	hasBounds := minX <= maxX
	if hasBounds && (!finite(minX) || !finite(maxX) ||
		!finite(minY) || !finite(maxY)) {
		slog.Warn("non-finite bounds", "chart", cfg.ID)
		return false
	}

	if cfg.XAxis != nil {
		av.xAxis = cfg.XAxis
		if hasBounds {
			av.xAxis.SetRange(minX, maxX)
		}
	} else {
		if !hasBounds {
			slog.Warn("all series empty", "chart", cfg.ID)
			return false
		}
		av.xAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		av.xAxis.SetRange(minX, maxX)
	}

	if cfg.YAxis != nil {
		av.yAxis = cfg.YAxis
		if hasBounds {
			av.yAxis.SetRange(minY, maxY)
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
		if !cfg.Stacked {
			minY -= yRange * 0.05
		}
		maxY += yRange * 0.05
		av.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		av.yAxis.SetRange(minY, maxY)
	}
	av.lastVersion = cfg.Version
	return true
}

func (av *areaView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	cfg := &av.cfg
	th := cfg.Theme

	if len(cfg.Series) == 0 {
		slog.Warn("no series data", "chart", cfg.ID)
		return
	}

	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	drawTitle(ctx, cfg.Title, th)

	if av.xAxis == nil || cfg.Version != av.lastVersion {
		if !av.updateAxes() {
			return
		}
	}

	xAxis := av.xAxis
	yAxis := av.yAxis

	av.yTicks = yAxis.Ticks(bottom, top)
	av.xTicks = xAxis.Ticks(left, right)

	for _, t := range av.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}
	for _, t := range av.xTicks {
		ctx.Line(t.Position, top, t.Position, bottom,
			th.GridColor, th.GridWidth)
	}

	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range av.xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		lw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(t.Position-lw/2, bottom+tickLen+2, t.Label, tickStyle)
	}
	for _, t := range av.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2, t.Label, tickStyle)
	}

	drawXAxisLabel(ctx, xAxis.Label(), th, left, right, bottom)
	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	alpha := uint8(cfg.Opacity * 255)

	if cfg.Stacked {
		av.drawStacked(ctx, cfg, xAxis, yAxis, left, right, top, bottom, alpha)
	} else {
		av.drawOverlapping(ctx, cfg, xAxis, yAxis, left, right, top, bottom, alpha)
	}

	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
		}
	}
	drawLegend(ctx, entries, th, left, right, top, bottom, cfg.LegendPosition)

	// Tooltip.
	if av.hovering && av.xAxis != nil {
		pa := plotArea{left, right, top, bottom, xAxis, yAxis}
		drawXYTooltip(ctx, th, cfg.Series, pa,
			av.hoverPx, av.hoverPy)
	}
}

func (av *areaView) drawOverlapping(
	ctx *render.Context, cfg *AreaCfg,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	alpha uint8,
) {
	for i, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		color := seriesColor(s.Color(), i, cfg.Theme.Palette)

		needed := s.Len() * 2
		if cap(av.ptsBuf) < needed {
			av.ptsBuf = make([]float32, 0, needed)
		}
		pts := av.ptsBuf[:0]
		for _, p := range s.Points {
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(p.Y, bottom, top)
			pts = append(pts, px, py)
		}
		av.ptsBuf = pts

		if len(pts) >= 4 {
			fill := gui.RGBA(color.R, color.G, color.B, alpha)
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
	}
}

func (av *areaView) drawStacked(
	ctx *render.Context, cfg *AreaCfg,
	xAxis, yAxis *axis.Linear,
	left, right, top, bottom float32,
	alpha uint8,
) {
	// Find reference point count from first non-empty series.
	refLen := 0
	for _, s := range cfg.Series {
		if s.Len() > 0 {
			refLen = s.Len()
			break
		}
	}
	if refLen == 0 {
		return
	}

	cumY := make([]float64, refLen)

	// prevPts holds the pixel coords of the previous series' top edge.
	// Initialize to baseline Y using first series' X positions.
	needed := refLen * 2
	if cap(av.prevPtsBuf) < needed {
		av.prevPtsBuf = make([]float32, needed)
	}
	prev := av.prevPtsBuf[:needed]
	for j, p := range cfg.Series[0].Points {
		if j >= refLen {
			break
		}
		prev[j*2] = xAxis.Transform(p.X, left, right)
		prev[j*2+1] = bottom
	}

	for i, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		color := seriesColor(s.Color(), i, cfg.Theme.Palette)
		fill := gui.RGBA(color.R, color.G, color.B, alpha)

		n := min(s.Len(), refLen)
		// cur is allocated per-series so prev can safely reference the
		// previous iteration's slice.
		cur := make([]float32, n*2)
		for j := range n {
			p := s.Points[j]
			cumY[j] += p.Y
			cur[j*2] = xAxis.Transform(p.X, left, right)
			cur[j*2+1] = yAxis.Transform(cumY[j], bottom, top)
		}

		// Fill quad between cur top edge and prev top edge (or baseline).
		if len(cur) >= 4 {
			var quad [8]float32
			for k := 0; k < len(cur)-2; k += 2 {
				quad[0] = cur[k]
				quad[1] = cur[k+1]
				quad[2] = cur[k+2]
				quad[3] = cur[k+3]
				quad[4] = prev[k+2]
				quad[5] = prev[k+3]
				quad[6] = prev[k]
				quad[7] = prev[k+1]
				ctx.FilledPolygon(quad[:], fill)
			}
		}
		ctx.Polyline(cur, color, cfg.LineWidth)
		copy(prev, cur)
	}
}
