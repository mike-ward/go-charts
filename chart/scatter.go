package chart

import (
	"github.com/mike-ward/go-charts/axis"
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
	BaseCfg

	// Data
	Series []series.XY

	// Axes (optional; auto-created from series bounds when nil)
	XAxis *axis.Linear
	YAxis *axis.Linear

	// Appearance
	MarkerSize float32 // 0 means default (6)
	Marker     MarkerShape
}

type scatterView struct {
	cfg ScatterCfg
}

// Scatter creates a scatter plot view.
func Scatter(cfg ScatterCfg) gui.View {
	cfg.applyDefaults()
	if cfg.MarkerSize == 0 {
		cfg.MarkerSize = DefaultMarkerSize
	}
	return &scatterView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (sv *scatterView) Draw(dc *gui.DrawContext) { sv.draw(dc) }

func (sv *scatterView) chartTheme() *theme.Theme { return sv.cfg.Theme }

func (sv *scatterView) Content() []gui.View { return nil }

func (sv *scatterView) GenerateLayout(w *gui.Window) gui.Layout {
	c := &sv.cfg
	width, height := resolveSize(c.Width, c.Height, w)
	return gui.DrawCanvas(gui.DrawCanvasCfg{
		ID:      c.ID,
		Sizing:  c.Sizing,
		Width:   width,
		Height:  height,
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
