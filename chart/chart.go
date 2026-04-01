// Package chart provides chart widgets for go-gui.
package chart

import (
	"log/slog"
	"math"

	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// Cfg configures a chart widget.
type Cfg struct {
	BaseCfg

	// Data
	Series []series.Series

	// Appearance
	Padding gui.Padding
}

// chartView implements gui.View for charts.
type chartView struct {
	cfg Cfg
}

// Chart creates a new chart view.
func Chart(cfg Cfg) gui.View {
	cfg.applyDefaults()
	return &chartView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (cv *chartView) Draw(dc *gui.DrawContext) { cv.draw(dc) }

func (cv *chartView) chartTheme() *theme.Theme { return cv.cfg.Theme }

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

// finite reports whether v is neither NaN nor +/-Inf.
func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
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
