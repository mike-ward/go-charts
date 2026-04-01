package chart

import (
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// PieCfg configures a pie or donut chart.
type PieCfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Data
	Labels []string
	Values []float64
	Colors []gui.Color

	// Appearance
	Theme       *theme.Theme
	InnerRadius float32 // >0 makes it a donut chart
	StartAngle  float32 // in radians
	ShowLabels  bool
	ShowPercent bool

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	Version uint64
}

type pieView struct {
	cfg PieCfg
}

// Pie creates a pie or donut chart view.
func Pie(cfg PieCfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	return &pieView{cfg: cfg}
}

func (pv *pieView) Content() []gui.View { return nil }

func (pv *pieView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &pv.cfg
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  pv.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
}

func (pv *pieView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	_ = ctx // TODO: render pie segments
}
