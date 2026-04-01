// Package chart provides chart widgets for go-gui.
package chart

import (
	"log/slog"

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
	width, height := resolveSize(c.Width, c.Height, w)
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
	slog.Debug("chart.Chart has no renderer; use Line, Bar, etc.",
		"chart", cv.cfg.ID)
}

// seriesColor returns the explicit color if set, otherwise wraps
// into the palette. Falls back to visible gray if palette is empty.
func seriesColor(
	color gui.Color, index int, palette []gui.Color,
) gui.Color {
	if color != (gui.Color{}) {
		return color
	}
	if len(palette) == 0 {
		return gui.Hex(0x808080)
	}
	return palette[index%len(palette)]
}

// resolveSize returns width/height, falling back to window
// dimensions when either is zero.
func resolveSize(width, height float32, w *gui.Window) (float32, float32) {
	if width == 0 || height == 0 {
		ww, wh := w.WindowSize()
		if width == 0 {
			width = float32(ww)
		}
		if height == 0 {
			height = float32(wh)
		}
	}
	return width, height
}
