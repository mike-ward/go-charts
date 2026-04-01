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

// BarCfg configures a bar chart.
type BarCfg struct {
	BaseCfg

	// Data
	Series []series.Category

	// Axes (optional; Y auto-created from series bounds when nil)
	YAxis *axis.Linear

	// Appearance
	BarWidth float32
	BarGap   float32
	Radius   float32 // corner radius for bars
}

type barView struct {
	cfg         BarCfg
	lastVersion uint64
	yAxis       *axis.Linear
	yTicks      []axis.Tick
}

// Bar creates a bar chart view.
func Bar(cfg BarCfg) gui.View {
	cfg.applyDefaults()
	return &barView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (bv *barView) Draw(dc *gui.DrawContext) { bv.draw(dc) }

func (bv *barView) chartTheme() *theme.Theme { return bv.cfg.Theme }

func (bv *barView) Content() []gui.View { return nil }

func (bv *barView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &bv.cfg
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  bv.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
}

func (bv *barView) draw(dc *gui.DrawContext) {
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

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	// Title.
	drawTitle(ctx, cfg.Title, th)

	// Collect all category labels from the first series.
	labels := cfg.Series[0].Values
	nCategories := len(labels)
	if nCategories == 0 {
		slog.Warn("no category data", "chart", cfg.ID)
		return
	}
	nSeries := len(cfg.Series)

	// Recompute Y axis only when version changes.
	if bv.yAxis == nil || cfg.Version != bv.lastVersion {
		if cfg.YAxis != nil {
			bv.yAxis = cfg.YAxis
		} else {
			minVal := 0.0
			maxVal := 0.0
			for _, s := range cfg.Series {
				for _, v := range s.Values {
					if !finite(v.Value) {
						continue
					}
					minVal = min(minVal, v.Value)
					maxVal = max(maxVal, v.Value)
				}
			}
			if minVal == 0 && maxVal == 0 {
				maxVal = 1
			}
			rangeVal := maxVal - minVal
			if rangeVal == 0 {
				rangeVal = 1
			}
			pad := rangeVal * 0.05
			// Pad away from zero only; bars anchor at the
			// baseline so the domain must include zero exactly.
			lo := minVal
			if lo < 0 {
				lo -= pad
			}
			hi := maxVal
			if hi > 0 {
				hi += pad
			}
			bv.yAxis = axis.NewLinear(
				axis.LinearCfg{AutoRange: true})
			bv.yAxis.SetRange(
				min(0, lo),
				max(0, hi),
			)
		}
		bv.lastVersion = cfg.Version
	}

	yAxis := bv.yAxis

	// Generate ticks.
	bv.yTicks = yAxis.Ticks(bottom, top)

	// Draw grid lines.
	for _, t := range bv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}

	// Draw axes.
	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	// Tick marks and labels on Y axis.
	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range bv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2,
			t.Label, tickStyle)
	}

	// Axis labels.
	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	// Baseline (y=0) pixel position.
	baseline := yAxis.Transform(0, bottom, top)

	// Bar layout.
	chartW := right - left
	groupWidth := chartW / float32(nCategories)

	barGap := cfg.BarGap
	if barGap == 0 {
		barGap = DefaultBarGap
	}

	barWidth := cfg.BarWidth
	if barWidth == 0 {
		usable := groupWidth - barGap*2
		if nSeries > 0 {
			barWidth = (usable - barGap*float32(nSeries-1)) /
				float32(nSeries)
		}
		barWidth = max(barWidth, 2)
	}

	// Draw bars.
	for ci := range nCategories {
		groupX := left + float32(ci)*groupWidth
		barStart := groupX + (groupWidth-
			float32(nSeries)*barWidth-
			float32(nSeries-1)*barGap)/2

		for si, s := range cfg.Series {
			if ci >= len(s.Values) {
				slog.Warn("series length mismatch",
					"chart", cfg.ID, "series", si)
				continue
			}
			v := s.Values[ci].Value
			if !finite(v) {
				continue
			}
			color := seriesColor(s.Color(), si, th.Palette)

			bx := barStart + float32(si)*(barWidth+barGap)
			by := yAxis.Transform(v, bottom, top)
			barTop := min(by, baseline)
			bh := float32(math.Abs(float64(by - baseline)))

			ctx.FilledRect(bx, barTop, barWidth, bh, color)
		}

		// Tick mark and label at center of group on X axis.
		cx := groupX + groupWidth/2
		ctx.Line(cx, bottom, cx, bottom+tickLen,
			tickColor, tickWidth)
		label := labels[ci].Label
		xts := tickStyle
		if cfg.XTickRotation != 0 {
			xts.RotationRadians = cfg.XTickRotation
			ctx.Text(cx, bottom+tickLen+2, label, xts)
		} else {
			lw := ctx.TextWidth(label, xts)
			ctx.Text(cx-lw/2, bottom+tickLen+2, label, xts)
		}
	}

	// X axis label (bar uses the first series name or a custom label
	// when a category axis is configured; skip for now since BarCfg
	// has no XAxis field — users set it via YAxis.Label()).

	// Legend.
	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
		}
	}
	drawLegend(ctx, entries, th, left, right, top, bottom,
		cfg.LegendPosition)
}
