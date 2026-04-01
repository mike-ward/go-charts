package chart

import (
	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-gui/gui"
)

// PieSlice represents a single slice of a pie chart.
type PieSlice struct {
	Label string
	Value float64
	Color gui.Color
}

// PieCfg configures a pie or donut chart.
type PieCfg struct {
	BaseCfg

	// Data
	Slices []PieSlice

	// Appearance
	InnerRadius float32 // >0 makes it a donut chart
	StartAngle  float32 // in radians
	ShowLabels  bool
	ShowPercent bool
}

type pieView struct {
	cfg PieCfg
}

// Pie creates a pie or donut chart view.
func Pie(cfg PieCfg) gui.View {
	cfg.applyDefaults()
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
