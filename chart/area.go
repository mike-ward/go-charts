package chart

import (
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// AreaCfg configures an area chart.
type AreaCfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Data
	Series []series.XY

	// Appearance
	Theme     *theme.Theme
	Stacked   bool
	LineWidth float32
	Opacity   float32 // fill opacity (0-1)

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	Version uint64
}

type areaView struct {
	cfg AreaCfg
}

// Area creates an area chart view.
func Area(cfg AreaCfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	if cfg.LineWidth == 0 {
		cfg.LineWidth = 2
	}
	if cfg.Opacity == 0 {
		cfg.Opacity = 0.3
	}
	return &areaView{cfg: cfg}
}

func (av *areaView) Content() []gui.View { return nil }

func (av *areaView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &av.cfg
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   c.Width,
		Height:  c.Height,
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
