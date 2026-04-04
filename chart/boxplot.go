package chart

import (
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// BoxData holds raw values for one box in a box plot.
type BoxData struct {
	Label  string
	Values []float64
	Color  gui.Color // zero = palette
}

// BoxPlotCfg configures a box plot chart.
type BoxPlotCfg struct {
	BaseCfg

	// Data holds one BoxData per box (category).
	Data []BoxData

	// YAxis overrides the auto-computed Y axis.
	YAxis *axis.Linear

	// BoxWidth is the body width in pixels. 0 = auto
	// (slot width * DefaultBoxWidthRatio).
	BoxWidth float32

	// ShowOutliers controls outlier dot rendering. nil = true.
	ShowOutliers *bool

	// OutlierRadius is the radius of outlier dots. 0 = default.
	OutlierRadius float32
}

// boxStats holds computed statistics for one box.
type boxStats struct {
	Min, Q1, Median, Q3, Max float64
	Outliers                 []float64
}

type boxplotView struct {
	cfg         BoxPlotCfg
	lastVersion uint64
	stats       []boxStats
	valid       []bool // valid[i] = computeBoxStats succeeded
	yAxis       *axis.Linear
	yTicks      []axis.Tick
	xAxis       *axis.Category
	hoverPx     float32
	hoverPy     float32
	hovering    bool
	hidden      map[int]bool
	lastLeft    float32
	lastRight   float32
	lastTop     float32
	lastBottom  float32
	lastLB      legendBounds
	win         *gui.Window
}

// BoxPlot creates a box plot chart view.
func BoxPlot(cfg BoxPlotCfg) gui.View {
	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &boxplotView{cfg: cfg}
}

// Draw renders the box plot onto dc for headless export.
func (bv *boxplotView) Draw(dc *gui.DrawContext) { bv.draw(dc) }

func (bv *boxplotView) chartTheme() *theme.Theme { return bv.cfg.Theme }

func (bv *boxplotView) Content() []gui.View { return nil }

func (bv *boxplotView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &bv.cfg
	hv := loadHover(w, c.ID, &bv.hovering, &bv.hoverPx, &bv.hoverPy)
	var hidV uint64
	bv.hidden, hidV = loadHiddenState(w, c.ID)
	bv.lastLB = loadLegendBounds(w, c.ID)
	bv.win = w
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version + hv + hidV,
		Clip:         true,
		OnDraw:       bv.draw,
		OnClick:      bv.internalClick,
		OnHover:      bv.internalHover,
		OnMouseLeave: bv.internalMouseLeave,
	}).GenerateLayout(w)
}

func (bv *boxplotView) internalClick(
	l *gui.Layout, e *gui.Event, w *gui.Window,
) {
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

func (bv *boxplotView) internalHover(
	l *gui.Layout, e *gui.Event, w *gui.Window,
) {
	e.IsHandled = true
	bv.hoverPx = e.MouseX - l.Shape.X
	bv.hoverPy = e.MouseY - l.Shape.Y
	bv.hovering = true
	saveHover(w, l, bv.cfg.ID, true, bv.hoverPx, bv.hoverPy)
	if legendHitTest(bv.lastLB, bv.hoverPx, bv.hoverPy) >= 0 {
		w.SetMouseCursorPointingHand()
	} else if bv.hoverPx >= bv.lastLeft &&
		bv.hoverPx <= bv.lastRight &&
		bv.hoverPy >= bv.lastTop &&
		bv.hoverPy <= bv.lastBottom {
		w.SetMouseCursorPointingHand()
	} else {
		w.SetMouseCursorArrow()
	}
	if bv.cfg.OnHover != nil {
		bv.cfg.OnHover(l, e, w)
	}
}

func (bv *boxplotView) internalMouseLeave(
	l *gui.Layout, e *gui.Event, w *gui.Window,
) {
	e.IsHandled = true
	bv.hovering = false
	saveHover(w, l, bv.cfg.ID, false, 0, 0)
	w.SetMouseCursorArrow()
	if bv.cfg.OnMouseLeave != nil {
		bv.cfg.OnMouseLeave(l, e, w)
	}
}

func (bv *boxplotView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	cfg := &bv.cfg
	th := cfg.Theme

	if len(cfg.Data) == 0 {
		slog.Warn("no data", "chart", cfg.ID)
		return
	}

	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	names := make([]string, len(cfg.Data))
	for i, d := range cfg.Data {
		names[i] = d.Label
	}
	right -= legendRightReserve(ctx, th, cfg.LegendPosition, names)
	top += legendTopReserve(ctx, th, cfg.LegendPosition, names, left, right)

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	drawTitle(ctx, cfg.Title, th)

	// Rebuild axes and stats when version changes.
	if bv.yAxis == nil || cfg.Version != bv.lastVersion {
		bv.buildAxesAndStats(cfg, th)
		bv.lastVersion = cfg.Version
	}

	left = resolveLeft(ctx, th, left, bottom, top, bv.yAxis)

	bottom = ctx.Height() - resolveBottom(ctx, th,
		maxTickLabelWidth(ctx, bv.xAxis.Ticks(left, right), th.TickStyle),
		cfg.XTickRotation, bv.xAxis.Label())
	bottom -= legendBottomReserve(ctx, th, cfg.LegendPosition, names, left, right)

	bv.lastLeft = left
	bv.lastRight = right
	bv.lastTop = top
	bv.lastBottom = bottom

	bv.yTicks = bv.yAxis.Ticks(bottom, top)

	// Grid lines.
	for _, t := range bv.yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}

	// Axis lines.
	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	// Y tick marks and labels.
	tickLen, tickWidth, tickColor := resolvedTickMark(th)
	tickStyle := th.TickStyle
	fh := ctx.FontHeight(tickStyle)
	for _, t := range bv.yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			tickColor, tickWidth)
		tw := ctx.TextWidth(t.Label, tickStyle)
		ctx.Text(left-tickLen-tw-2, t.Position-fh/2, t.Label, tickStyle)
	}
	drawYAxisLabel(ctx, bv.yAxis.Label(), th, top, bottom)

	// Annotations.
	drawAnnotations(ctx, &cfg.Annotations, th,
		plotRect{left, right, top, bottom}, bv.xAxis, bv.yAxis)

	// Determine hovered box index.
	hovI := -1
	if bv.hovering {
		hovI = bv.hoveredBox(bv.hoverPx, left, right)
	}

	nCat := len(cfg.Data)
	slotW := (right - left) / float32(nCat)
	boxW := cfg.BoxWidth
	if boxW <= 0 {
		boxW = slotW * DefaultBoxWidthRatio
	}
	boxW = max(boxW, 1)

	showOutliers := cfg.ShowOutliers == nil || *cfg.ShowOutliers
	outlierR := cfg.OutlierRadius
	if outlierR <= 0 {
		outlierR = DefaultOutlierRadius
	}
	outlierR = min(outlierR, boxW/2)

	// Draw boxes.
	for i, d := range cfg.Data {
		if bv.hidden[i] || !bv.valid[i] {
			continue
		}
		st := bv.stats[i]

		color := seriesColor(d.Color, i, th.Palette)
		if hovI >= 0 && hovI != i {
			color = dimColor(color, HoverDimAlpha)
		}

		cx := bv.xAxis.Transform(float64(i), left, right)

		q1Px := bv.yAxis.Transform(st.Q1, bottom, top)
		q3Px := bv.yAxis.Transform(st.Q3, bottom, top)
		medPx := bv.yAxis.Transform(st.Median, bottom, top)
		minPx := bv.yAxis.Transform(st.Min, bottom, top)
		maxPx := bv.yAxis.Transform(st.Max, bottom, top)

		// Box body: Q1 to Q3.
		bodyTop := min(q1Px, q3Px)
		bodyH := float32(math.Abs(float64(q1Px - q3Px)))
		bodyH = max(bodyH, 1)
		ctx.FilledRect(cx-boxW/2, bodyTop, boxW, bodyH, color)
		ctx.Rect(cx-boxW/2, bodyTop, boxW, bodyH,
			th.AxisColor, th.AxisWidth)

		// Median line.
		ctx.Line(cx-boxW/2, medPx, cx+boxW/2, medPx,
			th.AxisColor, 2.5)

		// Whiskers: vertical lines.
		ctx.Line(cx, q3Px, cx, maxPx, color, 1.5)
		ctx.Line(cx, q1Px, cx, minPx, color, 1.5)

		// Whisker caps.
		capW := boxW / 4
		ctx.Line(cx-capW, maxPx, cx+capW, maxPx, color, 1.5)
		ctx.Line(cx-capW, minPx, cx+capW, minPx, color, 1.5)

		// Outliers.
		if showOutliers {
			for _, v := range st.Outliers {
				oy := bv.yAxis.Transform(v, bottom, top)
				ctx.FilledCircle(cx, oy, outlierR, color)
			}
		}
	}

	// X tick marks and labels.
	xts := tickStyle
	if cfg.XTickRotation != 0 {
		xts.RotationRadians = cfg.XTickRotation
	}
	xTicks := bv.xAxis.Ticks(left, right)
	for _, t := range xTicks {
		ctx.Line(t.Position, bottom, t.Position, bottom+tickLen,
			tickColor, tickWidth)
		if cfg.XTickRotation != 0 {
			ctx.Text(t.Position, bottom+tickLen+2, t.Label, xts)
		} else {
			lw := ctx.TextWidth(t.Label, xts)
			ctx.Text(t.Position-lw/2, bottom+tickLen+2, t.Label, xts)
		}
	}

	// Legend.
	entries := make([]legendEntry, 0, len(cfg.Data))
	for i, d := range cfg.Data {
		entries = append(entries, legendEntry{
			Name:  d.Label,
			Color: seriesColor(d.Color, i, th.Palette),
			Index: i,
		})
	}
	pr := plotRect{left, right, top, bottom}
	bv.lastLB = drawLegend(ctx, entries, th, pr,
		cfg.LegendPosition, bv.hidden)
	saveLegendBounds(bv.win, cfg.ID, bv.lastLB)

	if bv.hovering {
		drawCrosshair(ctx, th, bv.hoverPx, bv.hoverPy, pr)
		bv.tooltipBoxPlot(ctx, pr, th)
	}
}

// buildAxesAndStats computes box statistics and rebuilds axes.
func (bv *boxplotView) buildAxesAndStats(
	cfg *BoxPlotCfg, _ *theme.Theme,
) {
	nCat := len(cfg.Data)
	bv.stats = make([]boxStats, nCat)
	bv.valid = make([]bool, nCat)

	globalMin := math.MaxFloat64
	globalMax := -math.MaxFloat64

	for i, d := range cfg.Data {
		if bv.hidden[i] {
			continue
		}
		st, ok := computeBoxStats(d.Values)
		if !ok {
			continue
		}
		bv.stats[i] = st
		bv.valid[i] = true

		lo := st.Min
		hi := st.Max
		for _, v := range st.Outliers {
			lo = min(lo, v)
			hi = max(hi, v)
		}
		globalMin = min(globalMin, lo)
		globalMax = max(globalMax, hi)
	}

	if globalMin == math.MaxFloat64 {
		globalMin = 0
		globalMax = 1
	}
	if globalMin == globalMax {
		globalMax = globalMin + 1
	}

	// Y axis.
	if cfg.YAxis != nil {
		bv.yAxis = cfg.YAxis
	} else {
		span := globalMax - globalMin
		pad := span * 0.05
		bv.yAxis = axis.NewLinear(axis.LinearCfg{AutoRange: true})
		bv.yAxis.SetRange(globalMin-pad, globalMax+pad)
	}

	// X axis: category labels.
	labels := make([]string, nCat)
	for i, d := range cfg.Data {
		labels[i] = d.Label
	}
	bv.xAxis = axis.NewCategory(axis.CategoryCfg{Categories: labels})
}

// hoveredBox returns the index of the box under pixel mx,
// or -1 when outside the plot area.
func (bv *boxplotView) hoveredBox(mx, left, right float32) int {
	if mx < left || mx > right {
		return -1
	}
	nCat := len(bv.cfg.Data)
	if nCat == 0 {
		return -1
	}
	slotW := (right - left) / float32(nCat)
	slot := int((mx - left) / slotW)
	return max(0, min(slot, nCat-1))
}

// tooltipBoxPlot draws a tooltip when the cursor is over a box
// (horizontally within the box body, vertically within whisker
// range).
func (bv *boxplotView) tooltipBoxPlot(
	ctx *render.Context, pr plotRect, th *theme.Theme,
) {
	if bv.yAxis == nil {
		return
	}
	left, right, top, bottom := pr.Left, pr.Right, pr.Top, pr.Bottom
	cfg := &bv.cfg
	mx := bv.hoverPx
	my := bv.hoverPy
	if mx < left || mx > right || my < top || my > bottom {
		return
	}

	i := bv.hoveredBox(mx, left, right)
	if i < 0 || i >= len(bv.stats) || bv.hidden[i] || !bv.valid[i] {
		return
	}

	st := bv.stats[i]
	cx := bv.xAxis.Transform(float64(i), left, right)

	// Resolve box width (must match draw loop).
	nCat := len(cfg.Data)
	slotW := (right - left) / float32(nCat)
	boxW := cfg.BoxWidth
	if boxW <= 0 {
		boxW = slotW * DefaultBoxWidthRatio
	}
	boxW = max(boxW, 1)

	// Horizontal: cursor must be within the box body.
	if mx < cx-boxW/2 || mx > cx+boxW/2 {
		return
	}

	// Vertical: cursor must be within the whisker range
	// (Min to Max, excluding outliers which are small dots).
	maxPx := bv.yAxis.Transform(st.Max, bottom, top)
	minPx := bv.yAxis.Transform(st.Min, bottom, top)
	if my < maxPx || my > minPx {
		return
	}

	q3Px := bv.yAxis.Transform(st.Q3, bottom, top)

	label := cfg.Data[i].Label
	nOut := len(st.Outliers)
	text := fmt.Sprintf("%s\nMax:    %.4g\nQ3:     %.4g\n"+
		"Median: %.4g\nQ1:     %.4g\nMin:    %.4g",
		label, st.Max, st.Q3, st.Median, st.Q1, st.Min)
	if nOut > 0 {
		text += fmt.Sprintf("\nOutliers: %d", nOut)
	}

	drawTooltip(ctx, cx, q3Px, text, th)
}

// computeBoxStats calculates quartiles, whiskers, and outliers
// from raw values. NaN/Inf values are silently filtered. Returns
// ok=false when fewer than 1 finite value remains.
func computeBoxStats(raw []float64) (boxStats, bool) {
	clean := make([]float64, 0, len(raw))
	for _, v := range raw {
		if finite(v) {
			clean = append(clean, v)
		}
	}
	if len(clean) == 0 {
		return boxStats{}, false
	}

	slices.Sort(clean)
	n := len(clean)

	med := median(clean)
	var q1, q3 float64
	if n == 1 {
		q1 = clean[0]
		q3 = clean[0]
	} else {
		mid := n / 2
		q1 = median(clean[:mid])
		if n%2 == 0 {
			q3 = median(clean[mid:])
		} else {
			q3 = median(clean[mid+1:])
		}
	}

	iqr := q3 - q1
	lowerFence := q1 - 1.5*iqr
	upperFence := q3 + 1.5*iqr

	// Whisker endpoints: furthest data within fences.
	wMin := clean[0]
	wMax := clean[n-1]
	for _, v := range clean {
		if v >= lowerFence {
			wMin = v
			break
		}
	}
	for i := n - 1; i >= 0; i-- {
		if clean[i] <= upperFence {
			wMax = clean[i]
			break
		}
	}

	// Collect outliers.
	var outliers []float64
	for _, v := range clean {
		if v < lowerFence || v > upperFence {
			outliers = append(outliers, v)
		}
	}

	return boxStats{
		Min:      wMin,
		Q1:       q1,
		Median:   med,
		Q3:       q3,
		Max:      wMax,
		Outliers: outliers,
	}, true
}

// median returns the median of a sorted, non-empty slice.
func median(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}
