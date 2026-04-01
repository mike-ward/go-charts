package chart

import (
	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// BarCfg configures a bar chart.
type BarCfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Data
	Series []series.Category

	// Appearance
	Theme    *theme.Theme
	BarWidth float32
	BarGap   float32
	Radius   float32 // corner radius for bars

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	Version uint64
}

type barView struct {
	cfg BarCfg
}

// Bar creates a bar chart view.
func Bar(cfg BarCfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	return &barView{cfg: cfg}
}

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
		return
	}

	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	if right <= left || bottom <= top {
		return
	}

	// Collect all category labels from the first series.
	labels := cfg.Series[0].Values
	nCategories := len(labels)
	if nCategories == 0 {
		return
	}
	nSeries := len(cfg.Series)

	// Find value range across all series.
	maxVal := 0.0
	for _, s := range cfg.Series {
		for _, v := range s.Values {
			maxVal = max(maxVal, v.Value)
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	// Y-axis for value scaling (0 to maxVal with 5% headroom).
	yAxis := axis.NewLinear(axis.LinearCfg{AutoRange: true})
	yAxis.SetRange(0, maxVal*1.05)

	// Draw grid lines.
	yTicks := yAxis.Ticks(bottom, top)
	for _, t := range yTicks {
		ctx.Line(left, t.Position, right, t.Position,
			th.GridColor, th.GridWidth)
	}

	// Draw axes.
	ctx.Line(left, bottom, right, bottom, th.AxisColor, th.AxisWidth)
	ctx.Line(left, top, left, bottom, th.AxisColor, th.AxisWidth)

	// Tick marks on Y axis.
	const tickLen float32 = 5
	for _, t := range yTicks {
		ctx.Line(left-tickLen, t.Position, left, t.Position,
			th.AxisColor, th.AxisWidth)
	}

	// Bar layout.
	chartW := right - left
	groupWidth := chartW / float32(nCategories)

	barGap := cfg.BarGap
	if barGap == 0 {
		barGap = 4
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
				continue
			}
			v := s.Values[ci].Value
			color := s.Color()
			if color == (gui.Color{}) && si < len(th.Palette) {
				color = th.Palette[si]
			}

			bx := barStart + float32(si)*(barWidth+barGap)
			by := yAxis.Transform(v, bottom, top)
			bh := bottom - by

			ctx.FilledRect(bx, by, barWidth, bh, color)
		}

		// Tick mark at center of group on X axis.
		cx := groupX + groupWidth/2
		ctx.Line(cx, bottom, cx, bottom+tickLen,
			th.AxisColor, th.AxisWidth)
	}
}
