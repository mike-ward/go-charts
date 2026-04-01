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
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Data
	Series []series.XY

	// Appearance
	Theme       *theme.Theme
	LineWidth   float32 // 0 means default (2)
	ShowMarkers bool
	ShowArea    bool // filled area under the line

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	Version uint64
}

type lineView struct {
	cfg         LineCfg
	lastVersion uint64
	xAxis       *axis.Linear
	yAxis       *axis.Linear
	xTicks      []axis.Tick
	yTicks      []axis.Tick
	ptsBuf      []float32
}

// Line creates a line chart view.
func Line(cfg LineCfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.LineWidth == 0 {
		cfg.LineWidth = 2
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	return &lineView{cfg: cfg}
}

func (lv *lineView) Content() []gui.View { return nil }

func (lv *lineView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &lv.cfg
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  lv.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
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

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	// Recompute axes only when version changes.
	if lv.xAxis == nil || cfg.Version != lv.lastVersion {
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
		if minX > maxX {
			// All series empty.
			slog.Warn("all series empty", "chart", cfg.ID)
			return
		}

		// Defense-in-depth: reject non-finite bounds that
		// slipped through series filtering.
		if !finite(minX) || !finite(maxX) ||
			!finite(minY) || !finite(maxY) {
			slog.Warn("non-finite bounds", "chart", cfg.ID)
			return
		}

		yRange := maxY - minY
		if yRange == 0 {
			yRange = 1
		}
		minY -= yRange * 0.05
		maxY += yRange * 0.05

		lv.xAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		lv.xAxis.SetRange(minX, maxX)
		lv.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		lv.yAxis.SetRange(minY, maxY)
		lv.lastVersion = cfg.Version
	}

	xAxis := lv.xAxis
	yAxis := lv.yAxis

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

	// Draw tick marks on axes.
	const tickLen float32 = 5
	for _, t := range lv.xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			th.AxisColor, th.AxisWidth)
	}
	for _, t := range lv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			th.AxisColor, th.AxisWidth)
	}

	// Draw each series.
	for i, s := range cfg.Series {
		if s.Len() == 0 {
			continue
		}
		color := seriesColor(s.Color(), i, th.Palette)

		// Build polyline points (flat x,y pairs), reusing buffer.
		needed := s.Len() * 2
		if cap(lv.ptsBuf) < needed {
			lv.ptsBuf = make([]float32, 0, needed)
		}
		pts := lv.ptsBuf[:0]
		for _, p := range s.Points {
			px := xAxis.Transform(p.X, left, right)
			py := yAxis.Transform(p.Y, bottom, top)
			pts = append(pts, px, py)
		}
		lv.ptsBuf = pts

		// Filled area under the line.
		if cfg.ShowArea && len(pts) >= 4 {
			area := make([]float32, len(pts), len(pts)+4)
			copy(area, pts)
			area = append(area, pts[len(pts)-2], bottom)
			area = append(area, pts[0], bottom)
			fill := gui.RGBA(color.R, color.G, color.B, 40)
			ctx.FilledPolygon(area, fill)
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
