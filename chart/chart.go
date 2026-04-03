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
		ID:           c.ID,
		Sizing:       c.Sizing,
		Width:        width,
		Height:       height,
		Version:      c.Version,
		Clip:         true,
		OnDraw:       cv.draw,
		OnClick:      c.OnClick,
		OnHover:      c.OnHover,
		OnMouseLeave: c.OnMouseLeave,
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

// hoverState persists hover information across frames via
// gui.StateMap. Chart views are typically recreated each frame
// by the view generator, so transient struct fields are lost.
// StateMap survives layout rebuilds.
type hoverState struct {
	Hovering bool
	Px, Py   float32
	Version  uint64
}

const (
	nsChartHover  = "chart-hover"
	capChartHover = 64
)

// chartHoverMap returns the persistent hover state map.
func chartHoverMap(w *gui.Window) *gui.BoundedMap[string, hoverState] {
	return gui.StateMap[string, hoverState](w, nsChartHover, capChartHover)
}

// loadHover reads persisted hover state into the view fields and
// returns the hover version for cache invalidation.
func loadHover(
	w *gui.Window, id string,
	hovering *bool, px, py *float32,
) uint64 {
	if id == "" {
		return 0
	}
	sm := chartHoverMap(w)
	hs, ok := sm.Get(id)
	if !ok {
		return 0
	}
	*hovering = hs.Hovering
	*px = hs.Px
	*py = hs.Py
	return hs.Version
}

// saveHover writes hover state to the persistent map and bumps
// l.Shape.Version so the draw-canvas cache misses this frame.
func saveHover(
	w *gui.Window, l *gui.Layout,
	id string, hovering bool, px, py float32,
) {
	if id == "" {
		return
	}
	sm := chartHoverMap(w)
	prev, _ := sm.Get(id)
	v := prev.Version + 1
	sm.Set(id, hoverState{
		Hovering: hovering,
		Px:       px,
		Py:       py,
		Version:  v,
	})
	l.Shape.Version = v
}
