package chart

import (
	"math"

	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// legendEntry describes one item in the chart legend.
type legendEntry struct {
	Name  string
	Color gui.Color
}

const (
	legendSwatchW float32 = 12
	legendSwatchH float32 = 12
	legendItemGap float32 = 4
	legendPadding float32 = 6
	legendRowGap  float32 = 2
)

// drawTitle renders the chart title centered above the plot area.
func drawTitle(
	ctx *render.Context, title string, th *theme.Theme,
) {
	if title == "" {
		return
	}
	style := th.TitleStyle
	tw := ctx.TextWidth(title, style)
	x := (ctx.Width() - tw) / 2
	fh := ctx.FontHeight(style)
	y := (th.PaddingTop - fh) / 2
	ctx.Text(x, y, title, style)
}

// drawXAxisLabel renders the X axis title centered below the
// tick labels. bottom is the Y coordinate of the X axis line.
func drawXAxisLabel(
	ctx *render.Context, label string, th *theme.Theme,
	left, right, bottom float32,
) {
	if label == "" {
		return
	}
	style := th.LabelStyle
	tw := ctx.TextWidth(label, style)
	tickFh := ctx.FontHeight(th.TickStyle)
	// Position below tick marks (5px) + tick labels + gap.
	x := (left + right - tw) / 2
	y := bottom + 5 + tickFh + 6
	ctx.Text(x, y, label, style)
}

// drawYAxisLabel renders the Y axis title rotated 90° CCW,
// centered vertically along the left edge of the plot area.
func drawYAxisLabel(
	ctx *render.Context, label string, th *theme.Theme,
	top, bottom float32,
) {
	if label == "" {
		return
	}
	style := th.LabelStyle
	style.RotationRadians = -math.Pi / 2
	fh := ctx.FontHeight(style)
	tw := ctx.TextWidth(label, style)
	// After -90° rotation, the text's visual width becomes its
	// height and vice versa. Position so the label is centered
	// vertically in the plot area, offset from the left edge.
	x := fh / 2
	y := (top + bottom + tw) / 2
	ctx.Text(x, y, label, style)
}

// drawLegend renders a legend in the top-right corner of the
// plot area. Skipped when no entries have names.
func drawLegend(
	ctx *render.Context,
	entries []legendEntry,
	th *theme.Theme,
	right, top float32,
) {
	// Filter to entries with names.
	named := make([]legendEntry, 0, len(entries))
	for _, e := range entries {
		if e.Name != "" {
			named = append(named, e)
		}
	}
	if len(named) == 0 {
		return
	}

	style := th.LabelStyle
	fh := ctx.FontHeight(style)
	rowH := max(fh, legendSwatchH)

	// Measure widest entry.
	maxW := float32(0)
	for _, e := range named {
		w := ctx.TextWidth(e.Name, style)
		maxW = max(maxW, w)
	}

	boxW := legendPadding*2 + legendSwatchW + legendItemGap + maxW
	boxH := legendPadding*2 +
		float32(len(named))*rowH +
		float32(len(named)-1)*legendRowGap

	bx := right - boxW - 4
	by := top + 4

	// Background.
	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4,
		gui.RGBA(0, 0, 0, 120))

	// Entries.
	for i, e := range named {
		ey := by + legendPadding +
			float32(i)*(rowH+legendRowGap)
		// Color swatch.
		sx := bx + legendPadding
		sy := ey + (rowH-legendSwatchH)/2
		ctx.FilledRoundedRect(sx, sy,
			legendSwatchW, legendSwatchH, 2, e.Color)
		// Label.
		tx := sx + legendSwatchW + legendItemGap
		ty := ey + (rowH-fh)/2
		ctx.Text(tx, ty, e.Name, style)
	}
}
