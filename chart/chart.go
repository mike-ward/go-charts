// Package chart provides chart widgets for go-gui.
package chart

import (
	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// Cfg configures a chart widget.
type Cfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Axes
	XAxis axis.Axis
	YAxis axis.Axis

	// Data
	Series []series.Series

	// Appearance
	Theme   *theme.Theme
	Padding gui.Padding

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	// Internal
	Version uint64
}

// chartView implements gui.View for charts.
type chartView struct {
	cfg Cfg
}

// Chart creates a new chart view.
func Chart(cfg Cfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	return &chartView{cfg: cfg}
}

func (cv *chartView) Content() []gui.View { return nil }

func (cv *chartView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &cv.cfg
	width, height := c.Width, c.Height
	if width == 0 || height == 0 {
		ww, wh := w.WindowSize()
		if width == 0 {
			width = float32(ww)
		}
		if height == 0 {
			height = float32(wh)
		}
	}
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  cv.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
}

func (cv *chartView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	_ = ctx // TODO: render chart content
}
