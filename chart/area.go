package chart

import (
	"log/slog"

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
	cfg AreaCfg
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
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  av.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
}

func (av *areaView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	_ = ctx // TODO: render area series
}
