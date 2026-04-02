package chart

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
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
	if err := cfg.Validate(); err != nil {
		slog.Warn("invalid config", "error", err)
	}
	return &pieView{cfg: cfg}
}

// Draw renders the chart onto dc for headless export.
func (pv *pieView) Draw(dc *gui.DrawContext) { pv.draw(dc) }

func (pv *pieView) chartTheme() *theme.Theme { return pv.cfg.Theme }

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
	cfg := &pv.cfg
	th := cfg.Theme

	if len(cfg.Slices) == 0 {
		slog.Warn("no slice data", "chart", cfg.ID)
		return
	}

	left := th.PaddingLeft
	right := ctx.Width() - th.PaddingRight
	top := th.PaddingTop
	bottom := ctx.Height() - th.PaddingBottom

	if right <= left || bottom <= top {
		slog.Warn("plot area too small", "chart", cfg.ID)
		return
	}

	drawTitle(ctx, cfg.Title, th)

	// Sum all positive slice values.
	total := 0.0
	for _, s := range cfg.Slices {
		if s.Value > 0 {
			total += s.Value
		}
	}
	if total == 0 {
		slog.Warn("all slice values zero or negative", "chart", cfg.ID)
		return
	}

	// Outer radius — leave 15% margin for labels.
	plotW := right - left
	plotH := bottom - top
	outerR := min(plotW, plotH) / 2 * 0.85
	cx := (left + right) / 2
	cy := (top + bottom) / 2

	// Draw slices.
	angle := cfg.StartAngle
	for i, s := range cfg.Slices {
		if s.Value <= 0 {
			continue
		}
		sweep := float32(s.Value/total) * (2 * math.Pi)
		color := s.Color
		if !color.IsSet() {
			color = seriesColor(gui.Color{}, i, th.Palette)
		}
		ctx.FilledArc(cx, cy, outerR, outerR, angle, sweep, color)
		angle += sweep
	}

	// Donut hole: overdraw center with background color.
	// Note: this only looks correct over a solid background.
	if cfg.InnerRadius > 0 {
		ctx.FilledCircle(cx, cy, cfg.InnerRadius, th.Background)
	}

	// Labels at midpoint of each slice arc.
	if cfg.ShowLabels || cfg.ShowPercent {
		style := th.TickStyle
		fh := ctx.FontHeight(style)
		labelR := outerR * 0.7
		if cfg.InnerRadius > 0 {
			// Place labels between inner radius and outer radius.
			labelR = (cfg.InnerRadius + outerR) / 2
		}

		angle = cfg.StartAngle
		for _, s := range cfg.Slices {
			if s.Value <= 0 {
				continue
			}
			sweep := float32(s.Value/total) * (2 * math.Pi)
			mid := angle + sweep/2
			lx := cx + labelR*float32(math.Cos(float64(mid)))
			ly := cy + labelR*float32(math.Sin(float64(mid)))

			text := ""
			if cfg.ShowLabels {
				text = s.Label
			}
			if cfg.ShowPercent {
				pct := fmt.Sprintf("%.1f%%", s.Value/total*100)
				if text != "" {
					text = text + " " + pct
				} else {
					text = pct
				}
			}
			if text == "" {
				angle += sweep
				continue
			}

			tw := ctx.TextWidth(text, style)
			ctx.Text(lx-tw/2, ly-fh/2, text, style)
			angle += sweep
		}
	}

	// Legend.
	entries := make([]legendEntry, len(cfg.Slices))
	for i, s := range cfg.Slices {
		color := s.Color
		if !color.IsSet() {
			color = seriesColor(gui.Color{}, i, th.Palette)
		}
		entries[i] = legendEntry{Name: s.Label, Color: color}
	}
	drawLegend(ctx, entries, th, left, right, top, bottom, cfg.LegendPosition)
}
