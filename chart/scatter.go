package chart

import (
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// MarkerShape controls the shape of scatter plot markers.
type MarkerShape uint8

// MarkerShape constants.
const (
	MarkerCircle MarkerShape = iota
	MarkerSquare
	MarkerTriangle
	MarkerDiamond
	MarkerCross
)

// ScatterCfg configures a scatter plot.
type ScatterCfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	// Data
	Series []series.XY

	// Appearance
	Theme      *theme.Theme
	MarkerSize float32
	Marker     MarkerShape

	// Interaction
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	Version uint64
}

type scatterView struct {
	cfg ScatterCfg
}

// Scatter creates a scatter plot view.
func Scatter(cfg ScatterCfg) gui.View {
	if cfg.Sizing == (gui.Sizing{}) {
		cfg.Sizing = gui.FillFill
	}
	if cfg.MarkerSize == 0 {
		cfg.MarkerSize = 6
	}
	if cfg.Theme == nil {
		cfg.Theme = theme.Default()
	}
	return &scatterView{cfg: cfg}
}

func (sv *scatterView) Content() []gui.View { return nil }

func (sv *scatterView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &sv.cfg
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   c.Width,
		Height:  c.Height,
		Version: c.Version,
		Clip:    true,
		OnDraw:  sv.draw,
		OnClick: c.OnClick,
		OnHover: c.OnHover,
	}).GenerateLayout(w)
}

func (sv *scatterView) draw(dc *gui.DrawContext) {
	ctx := render.NewContext(dc)
	_ = ctx // TODO: render scatter markers
}
