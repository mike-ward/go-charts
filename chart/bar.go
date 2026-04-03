package chart

import (
	"fmt"
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
	BarWidth   float32
	BarGap     float32
	Radius     float32 // corner radius for bars
	Horizontal bool    // draw bars left-to-right instead of bottom-to-top
	Stacked    bool    // stack series instead of grouping side-by-side
}

type barView struct {
	cfg         BarCfg
	lastVersion uint64
	yAxis       *axis.Linear
	yTicks      []axis.Tick
	hoverPx     float32
	hoverPy     float32
	hovering    bool
}

// Bar creates a bar chart view.
func Bar(cfg BarCfg) gui.View {
	cfg.applyDefaults()
	if cfg.BarGap == 0 {
		cfg.BarGap = DefaultBarGap
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &barView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (bv *barView) Draw(dc *gui.DrawContext) { bv.draw(dc) }

func (bv *barView) chartTheme() *theme.Theme { return bv.cfg.Theme }

func (bv *barView) Content() []gui.View { return nil }

func (bv *barView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &bv.cfg
	hv := loadHover(w, c.ID,
		&bv.hovering, &bv.hoverPx, &bv.hoverPy)
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version + hv,
		Clip:         true,
		OnDraw:       bv.draw,
		OnClick:      c.OnClick,
		OnHover:      bv.internalHover,
		OnMouseLeave: bv.internalMouseLeave,
	}).GenerateLayout(w)
}

func (bv *barView) internalHover(l *gui.Layout, e *gui.Event, w *gui.Window) {
	bv.hoverPx = e.MouseX - l.Shape.X
	bv.hoverPy = e.MouseY - l.Shape.Y
	bv.hovering = true
	saveHover(w, l, bv.cfg.ID, true, bv.hoverPx, bv.hoverPy)
	if bv.cfg.OnHover != nil {
		bv.cfg.OnHover(l, e, w)
	}
}

func (bv *barView) internalMouseLeave(l *gui.Layout, e *gui.Event, w *gui.Window) {
	bv.hovering = false
	saveHover(w, l, bv.cfg.ID, false, 0, 0)
	if bv.cfg.OnMouseLeave != nil {
		bv.cfg.OnMouseLeave(l, e, w)
	}
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

	// Recompute value axis only when version changes.
	if bv.yAxis == nil || cfg.Version != bv.lastVersion {
		if cfg.YAxis != nil {
			bv.yAxis = cfg.YAxis
		} else {
			minVal := 0.0
			maxVal := 0.0
			if cfg.Stacked {
				// Range is the max cumulative sum per category.
				posSums := make([]float64, nCategories)
				negSums := make([]float64, nCategories)
				for _, s := range cfg.Series {
					for ci, v := range s.Values {
						if !finite(v.Value) {
							continue
						}
						if v.Value >= 0 {
							posSums[ci] += v.Value
							maxVal = max(maxVal, posSums[ci])
						} else {
							negSums[ci] += v.Value
							minVal = min(minVal, negSums[ci])
						}
					}
				}
			} else {
				for _, s := range cfg.Series {
					for _, v := range s.Values {
						if !finite(v.Value) {
							continue
						}
						minVal = min(minVal, v.Value)
						maxVal = max(maxVal, v.Value)
					}
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
			lo := minVal
			if lo < 0 {
				lo -= pad
			}
			hi := maxVal
			if hi > 0 {
				hi += pad
			}
			bv.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
			bv.yAxis.SetRange(min(0, lo), max(0, hi))
		}
		bv.lastVersion = cfg.Version
	}

	if cfg.Horizontal {
		bv.drawHorizontal(ctx, cfg, th, nCategories, nSeries, labels,
			left, right, top, bottom)
	} else {
		bv.drawVertical(ctx, cfg, th, nCategories, nSeries, labels,
			left, right, top, bottom)
	}

	// Crosshair and tooltip.
	if bv.hovering {
		drawCrosshair(ctx, th, bv.hoverPx, bv.hoverPy,
			left, right, top, bottom)
		bv.tooltipBar(ctx, left, right, top, bottom, th)
	}
}

func (bv *barView) drawVertical(
	ctx *render.Context, cfg *BarCfg, th *theme.Theme,
	nCategories, nSeries int, labels []series.CategoryValue,
	left, right, top, bottom float32,
) {
	hovCI, hovSI, hovOK := -1, -1, false
	if bv.hovering {
		hovCI, hovSI, hovOK = bv.hoveredBarVertical(
			bv.hoverPx, bv.hoverPy, left, right, top, bottom)
	}

	yAxis := bv.yAxis
	bv.yTicks = yAxis.Ticks(bottom, top)

	for _, t := range bv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}

	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range bv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2, t.Label, tickStyle)
	}

	drawYAxisLabel(ctx, yAxis.Label(), th, top, bottom)

	baseline := yAxis.Transform(0, bottom, top)
	chartW := right - left
	groupWidth := chartW / float32(nCategories)
	barGap := cfg.BarGap

	if cfg.Stacked {
		barWidth := groupWidth - barGap*2
		barWidth = max(barWidth, 2)

		for ci := range nCategories {
			groupX := left + float32(ci)*groupWidth
			bx := groupX + barGap

			posOff := 0.0
			negOff := 0.0
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
				if hovOK && (ci != hovCI || si != hovSI) {
					color = dimColor(color, HoverDimAlpha)
				}

				var segTop, segBot float32
				if v >= 0 {
					segBot = yAxis.Transform(posOff, bottom, top)
					segTop = yAxis.Transform(posOff+v, bottom, top)
					posOff += v
				} else {
					segTop = yAxis.Transform(negOff, bottom, top)
					segBot = yAxis.Transform(negOff+v, bottom, top)
					negOff += v
				}
				bh := float32(math.Abs(float64(segBot - segTop)))
				by := min(segTop, segBot)
				if cfg.Radius > 0 {
					ctx.FilledRoundedRect(bx, by, barWidth, bh, cfg.Radius, color)
				} else {
					ctx.FilledRect(bx, by, barWidth, bh, color)
				}
			}

			cx := groupX + groupWidth/2
			ctx.Line(cx, bottom, cx, bottom+tickLen, tickColor, tickWidth)
			lw := ctx.TextWidth(labels[ci].Label, tickStyle)
			ctx.Text(cx-lw/2, bottom+tickLen+2, labels[ci].Label, tickStyle)
		}
	} else {
		barWidth := cfg.BarWidth
		if barWidth == 0 {
			usable := groupWidth - barGap*2
			if nSeries > 0 {
				barWidth = (usable - barGap*float32(nSeries-1)) /
					float32(nSeries)
			}
			barWidth = max(barWidth, 2)
		}

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
				if hovOK && (ci != hovCI || si != hovSI) {
					color = dimColor(color, HoverDimAlpha)
				}

				bx := barStart + float32(si)*(barWidth+barGap)
				by := yAxis.Transform(v, bottom, top)
				barTop := min(by, baseline)
				bh := float32(math.Abs(float64(by - baseline)))

				if cfg.Radius > 0 {
					ctx.FilledRoundedRect(bx, barTop, barWidth, bh, cfg.Radius, color)
				} else {
					ctx.FilledRect(bx, barTop, barWidth, bh, color)
				}
			}

			cx := groupX + groupWidth/2
			ctx.Line(cx, bottom, cx, bottom+tickLen, tickColor, tickWidth)
			xts := tickStyle
			if cfg.XTickRotation != 0 {
				xts.RotationRadians = cfg.XTickRotation
				ctx.Text(cx, bottom+tickLen+2, labels[ci].Label, xts)
			} else {
				lw := ctx.TextWidth(labels[ci].Label, xts)
				ctx.Text(cx-lw/2, bottom+tickLen+2, labels[ci].Label, xts)
			}
		}
	}

	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
		}
	}
	drawLegend(ctx, entries, th, left, right, top, bottom, cfg.LegendPosition)
}

// hoveredBarVertical returns the (catIdx, seriesIdx) of the bar under
// (mx, my) for vertical bar charts, or ok=false when none.
func (bv *barView) hoveredBarVertical(
	mx, my, left, right, top, bottom float32,
) (ci, si int, ok bool) {
	cfg := &bv.cfg
	if len(cfg.Series) == 0 || bv.yAxis == nil {
		return
	}
	yAxis := bv.yAxis
	labels := cfg.Series[0].Values
	nCat := len(labels)
	if nCat == 0 {
		return
	}
	nSer := len(cfg.Series)
	barGap := cfg.BarGap
	chartW := right - left
	groupW := chartW / float32(nCat)
	baseline := yAxis.Transform(0, bottom, top)

	if cfg.Stacked {
		barW := max(groupW-barGap*2, 2)
		for c := range nCat {
			bx := left + float32(c)*groupW + barGap
			if mx < bx || mx > bx+barW {
				continue
			}
			posOff, negOff := 0.0, 0.0
			for s, ser := range cfg.Series {
				if c >= len(ser.Values) {
					continue
				}
				v := ser.Values[c].Value
				if !finite(v) {
					continue
				}
				var segTop, segBot float32
				if v >= 0 {
					segBot = yAxis.Transform(posOff, bottom, top)
					segTop = yAxis.Transform(posOff+v, bottom, top)
					posOff += v
				} else {
					segTop = yAxis.Transform(negOff, bottom, top)
					segBot = yAxis.Transform(negOff+v, bottom, top)
					negOff += v
				}
				by := min(segTop, segBot)
				bh := float32(math.Abs(float64(segBot - segTop)))
				if my >= by && my <= by+bh {
					return c, s, true
				}
			}
		}
	} else {
		barW := cfg.BarWidth
		if barW == 0 {
			usable := groupW - barGap*2
			if nSer > 0 {
				barW = (usable - barGap*float32(nSer-1)) / float32(nSer)
			}
			barW = max(barW, 2)
		}
		for c := range nCat {
			groupX := left + float32(c)*groupW
			barStart := groupX + (groupW-
				float32(nSer)*barW-
				float32(nSer-1)*barGap)/2
			for s, ser := range cfg.Series {
				if c >= len(ser.Values) {
					continue
				}
				v := ser.Values[c].Value
				if !finite(v) {
					continue
				}
				bx := barStart + float32(s)*(barW+barGap)
				if mx < bx || mx > bx+barW {
					continue
				}
				by := yAxis.Transform(v, bottom, top)
				barTop := min(by, baseline)
				bh := float32(math.Abs(float64(by - baseline)))
				if my >= barTop && my <= barTop+bh {
					return c, s, true
				}
			}
		}
	}
	return
}

// hoveredBarHorizontal returns the (catIdx, seriesIdx) of the bar under
// (mx, my) for horizontal bar charts, or ok=false when none.
func (bv *barView) hoveredBarHorizontal(
	mx, my, left, right, top, bottom float32,
) (ci, si int, ok bool) {
	cfg := &bv.cfg
	if len(cfg.Series) == 0 || bv.yAxis == nil {
		return
	}
	xAxis := bv.yAxis // horizontal mode reuses yAxis for the value axis
	labels := cfg.Series[0].Values
	nCat := len(labels)
	if nCat == 0 {
		return
	}
	nSer := len(cfg.Series)
	barGap := cfg.BarGap
	chartH := bottom - top
	groupH := chartH / float32(nCat)
	baseline := xAxis.Transform(0, left, right)

	if cfg.Stacked {
		barH := max(groupH-barGap*2, 2)
		for c := range nCat {
			by := top + float32(c)*groupH + barGap
			if my < by || my > by+barH {
				continue
			}
			posOff, negOff := 0.0, 0.0
			for s, ser := range cfg.Series {
				if c >= len(ser.Values) {
					continue
				}
				v := ser.Values[c].Value
				if !finite(v) {
					continue
				}
				var segL, segR float32
				if v >= 0 {
					segL = xAxis.Transform(posOff, left, right)
					segR = xAxis.Transform(posOff+v, left, right)
					posOff += v
				} else {
					segR = xAxis.Transform(negOff, left, right)
					segL = xAxis.Transform(negOff+v, left, right)
					negOff += v
				}
				bx := min(segL, segR)
				bw := float32(math.Abs(float64(segR - segL)))
				if mx >= bx && mx <= bx+bw {
					return c, s, true
				}
			}
		}
	} else {
		barH := cfg.BarWidth
		if barH == 0 {
			usable := groupH - barGap*2
			if nSer > 0 {
				barH = (usable - barGap*float32(nSer-1)) / float32(nSer)
			}
			barH = max(barH, 2)
		}
		for c := range nCat {
			groupY := top + float32(c)*groupH
			barStart := groupY + (groupH-
				float32(nSer)*barH-
				float32(nSer-1)*barGap)/2
			for s, ser := range cfg.Series {
				if c >= len(ser.Values) {
					continue
				}
				v := ser.Values[c].Value
				if !finite(v) {
					continue
				}
				by := barStart + float32(s)*(barH+barGap)
				if my < by || my > by+barH {
					continue
				}
				bx := xAxis.Transform(v, left, right)
				barLeft := min(bx, baseline)
				bw := float32(math.Abs(float64(bx - baseline)))
				if mx >= barLeft && mx <= barLeft+bw {
					return c, s, true
				}
			}
		}
	}
	return
}

// tooltipBar hit-tests the cursor against actual bar rectangles
// and draws a tooltip only when the cursor is over a bar.
func (bv *barView) tooltipBar(
	ctx *render.Context,
	left, right, top, bottom float32,
	th *theme.Theme,
) {
	cfg := &bv.cfg
	if len(cfg.Series) == 0 || bv.yAxis == nil {
		return
	}
	mx := bv.hoverPx
	my := bv.hoverPy

	labels := cfg.Series[0].Values
	nCat := len(labels)
	if nCat == 0 {
		return
	}
	nSer := len(cfg.Series)
	barGap := cfg.BarGap

	if cfg.Horizontal {
		bv.tooltipHorizontal(ctx, mx, my,
			left, right, top, bottom, th,
			labels, nCat, nSer, barGap)
	} else {
		bv.tooltipVertical(ctx, mx, my,
			left, right, top, bottom, th,
			labels, nCat, nSer, barGap)
	}
}

func (bv *barView) tooltipVertical(
	ctx *render.Context, mx, my,
	left, right, top, bottom float32,
	th *theme.Theme,
	labels []series.CategoryValue,
	nCat, nSer int, barGap float32,
) {
	cfg := &bv.cfg
	yAxis := bv.yAxis
	baseline := yAxis.Transform(0, bottom, top)
	chartW := right - left
	groupW := chartW / float32(nCat)

	if cfg.Stacked {
		barW := max(groupW-barGap*2, 2)
		for ci := range nCat {
			bx := left + float32(ci)*groupW + barGap
			if mx < bx || mx > bx+barW {
				continue
			}
			posOff, negOff := 0.0, 0.0
			for _, s := range cfg.Series {
				if ci >= len(s.Values) {
					continue
				}
				v := s.Values[ci].Value
				if !finite(v) {
					continue
				}
				var segTop, segBot float32
				if v >= 0 {
					segBot = yAxis.Transform(posOff, bottom, top)
					segTop = yAxis.Transform(posOff+v, bottom, top)
					posOff += v
				} else {
					segTop = yAxis.Transform(negOff, bottom, top)
					segBot = yAxis.Transform(negOff+v, bottom, top)
					negOff += v
				}
				by := min(segTop, segBot)
				bh := float32(math.Abs(float64(segBot - segTop)))
				if my >= by && my <= by+bh {
					emitBarTooltip(ctx, mx, my, th,
						s.Name(), labels[ci].Label, v)
					return
				}
			}
		}
	} else {
		barW := cfg.BarWidth
		if barW == 0 {
			usable := groupW - barGap*2
			if nSer > 0 {
				barW = (usable - barGap*float32(nSer-1)) /
					float32(nSer)
			}
			barW = max(barW, 2)
		}
		for ci := range nCat {
			groupX := left + float32(ci)*groupW
			barStart := groupX + (groupW-
				float32(nSer)*barW-
				float32(nSer-1)*barGap)/2
			for si, s := range cfg.Series {
				if ci >= len(s.Values) {
					continue
				}
				v := s.Values[ci].Value
				if !finite(v) {
					continue
				}
				bx := barStart + float32(si)*(barW+barGap)
				if mx < bx || mx > bx+barW {
					continue
				}
				by := yAxis.Transform(v, bottom, top)
				barTop := min(by, baseline)
				bh := float32(math.Abs(float64(by - baseline)))
				if my >= barTop && my <= barTop+bh {
					emitBarTooltip(ctx, mx, my, th,
						s.Name(), labels[ci].Label, v)
					return
				}
			}
		}
	}
}

func (bv *barView) tooltipHorizontal(
	ctx *render.Context, mx, my,
	left, right, top, bottom float32,
	th *theme.Theme,
	labels []series.CategoryValue,
	nCat, nSer int, barGap float32,
) {
	cfg := &bv.cfg
	xAxis := bv.yAxis
	baseline := xAxis.Transform(0, left, right)
	chartH := bottom - top
	groupH := chartH / float32(nCat)

	if cfg.Stacked {
		barH := max(groupH-barGap*2, 2)
		for ci := range nCat {
			by := top + float32(ci)*groupH + barGap
			if my < by || my > by+barH {
				continue
			}
			posOff, negOff := 0.0, 0.0
			for _, s := range cfg.Series {
				if ci >= len(s.Values) {
					continue
				}
				v := s.Values[ci].Value
				if !finite(v) {
					continue
				}
				var segL, segR float32
				if v >= 0 {
					segL = xAxis.Transform(posOff, left, right)
					segR = xAxis.Transform(posOff+v, left, right)
					posOff += v
				} else {
					segR = xAxis.Transform(negOff, left, right)
					segL = xAxis.Transform(negOff+v, left, right)
					negOff += v
				}
				bx := min(segL, segR)
				bw := float32(math.Abs(float64(segR - segL)))
				if mx >= bx && mx <= bx+bw {
					emitBarTooltip(ctx, mx, my, th,
						s.Name(), labels[ci].Label, v)
					return
				}
			}
		}
	} else {
		barH := cfg.BarWidth
		if barH == 0 {
			usable := groupH - barGap*2
			if nSer > 0 {
				barH = (usable - barGap*float32(nSer-1)) /
					float32(nSer)
			}
			barH = max(barH, 2)
		}
		for ci := range nCat {
			groupY := top + float32(ci)*groupH
			barStart := groupY + (groupH-
				float32(nSer)*barH-
				float32(nSer-1)*barGap)/2
			for si, s := range cfg.Series {
				if ci >= len(s.Values) {
					continue
				}
				v := s.Values[ci].Value
				if !finite(v) {
					continue
				}
				by := barStart + float32(si)*(barH+barGap)
				if my < by || my > by+barH {
					continue
				}
				bx := xAxis.Transform(v, left, right)
				barLeft := min(bx, baseline)
				bw := float32(math.Abs(float64(bx - baseline)))
				if mx >= barLeft && mx <= barLeft+bw {
					emitBarTooltip(ctx, mx, my, th,
						s.Name(), labels[ci].Label, v)
					return
				}
			}
		}
	}
}

// emitBarTooltip formats and draws a single bar's tooltip.
func emitBarTooltip(
	ctx *render.Context, mx, my float32,
	th *theme.Theme,
	serName, catLabel string, v float64,
) {
	var label string
	if serName != "" {
		label = fmt.Sprintf("%s / %s: %g", serName, catLabel, v)
	} else {
		label = fmt.Sprintf("%s: %g", catLabel, v)
	}
	drawTooltip(ctx, mx, my, label, th)
}

func (bv *barView) drawHorizontal(
	ctx *render.Context, cfg *BarCfg, th *theme.Theme,
	nCategories, nSeries int, labels []series.CategoryValue,
	left, right, top, bottom float32,
) {
	hovCI, hovSI, hovOK := -1, -1, false
	if bv.hovering {
		hovCI, hovSI, hovOK = bv.hoveredBarHorizontal(
			bv.hoverPx, bv.hoverPy, left, right, top, bottom)
	}

	// In horizontal mode bv.yAxis is the value axis mapped to the X
	// direction: Transform(v, left, right) → X pixel.
	xAxis := bv.yAxis
	bv.yTicks = xAxis.Ticks(left, right)

	// Vertical grid lines (parallel to bars).
	for _, t := range bv.yTicks {
		ctx.Line(t.Position, top, t.Position, bottom,
			th.GridColor, th.GridWidth)
	}

	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)

	// Value ticks along the bottom (X axis).
	for _, t := range bv.yTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		lw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(t.Position-lw/2, bottom+tickLen+2, t.Label, tickStyle)
	}

	drawXAxisLabel(ctx, xAxis.Label(), th, left, right, bottom)

	baseline := xAxis.Transform(0, left, right)
	chartH := bottom - top
	groupHeight := chartH / float32(nCategories)
	barGap := cfg.BarGap

	if cfg.Stacked {
		barHeight := groupHeight - barGap*2
		barHeight = max(barHeight, 2)

		for ci := range nCategories {
			groupY := top + float32(ci)*groupHeight
			by := groupY + barGap

			posOff := 0.0
			negOff := 0.0
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
				if hovOK && (ci != hovCI || si != hovSI) {
					color = dimColor(color, HoverDimAlpha)
				}

				var segLeft, segRight float32
				if v >= 0 {
					segLeft = xAxis.Transform(posOff, left, right)
					segRight = xAxis.Transform(posOff+v, left, right)
					posOff += v
				} else {
					segRight = xAxis.Transform(negOff, left, right)
					segLeft = xAxis.Transform(negOff+v, left, right)
					negOff += v
				}
				bw := float32(math.Abs(float64(segRight - segLeft)))
				bx := min(segLeft, segRight)
				if cfg.Radius > 0 {
					ctx.FilledRoundedRect(bx, by, bw, barHeight, cfg.Radius, color)
				} else {
					ctx.FilledRect(bx, by, bw, barHeight, color)
				}
			}

			cy := groupY + groupHeight/2
			ctx.Line(left-tickLen, cy, left, cy, tickColor, tickWidth)
			lw := ctx.TextWidth(labels[ci].Label, tickStyle)
			ctx.Text(left-tickLen-lw-2, cy-fh/2, labels[ci].Label, tickStyle)
		}
	} else {
		barHeight := cfg.BarWidth // reuse BarWidth as bar thickness in H mode
		if barHeight == 0 {
			usable := groupHeight - barGap*2
			if nSeries > 0 {
				barHeight = (usable - barGap*float32(nSeries-1)) /
					float32(nSeries)
			}
			barHeight = max(barHeight, 2)
		}

		for ci := range nCategories {
			groupY := top + float32(ci)*groupHeight
			barStart := groupY + (groupHeight-
				float32(nSeries)*barHeight-
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
				if hovOK && (ci != hovCI || si != hovSI) {
					color = dimColor(color, HoverDimAlpha)
				}

				bx := xAxis.Transform(v, left, right)
				barLeft := min(bx, baseline)
				bw := float32(math.Abs(float64(bx - baseline)))
				by := barStart + float32(si)*(barHeight+barGap)

				if cfg.Radius > 0 {
					ctx.FilledRoundedRect(barLeft, by, bw, barHeight, cfg.Radius, color)
				} else {
					ctx.FilledRect(barLeft, by, bw, barHeight, color)
				}
			}

			// Category tick and label on the left.
			cy := groupY + groupHeight/2
			ctx.Line(left-tickLen, cy, left, cy, tickColor, tickWidth)
			lw := ctx.TextWidth(labels[ci].Label, tickStyle)
			ctx.Text(left-tickLen-lw-2, cy-fh/2, labels[ci].Label, tickStyle)
		}
	}

	entries := make([]legendEntry, len(cfg.Series))
	for i, s := range cfg.Series {
		entries[i] = legendEntry{
			Name:  s.Name(),
			Color: seriesColor(s.Color(), i, th.Palette),
		}
	}
	drawLegend(ctx, entries, th, left, right, top, bottom, cfg.LegendPosition)
}
