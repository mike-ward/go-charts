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

// resolvedTickMark returns the tick mark visual properties,
// falling back to axis defaults from the theme.
func resolvedTickMark(
	th *theme.Theme,
) (length, width float32, color gui.Color) {
	tms := th.TickMark
	length = tms.Length
	if length == 0 {
		length = DefaultTickLength
	}
	color = tms.Color
	if !color.IsSet() {
		color = th.AxisColor
	}
	width = tms.Width
	if width == 0 {
		width = th.AxisWidth
	}
	return
}

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

// drawLegend renders the legend in the plot area. Position is
// determined by the theme LegendStyle with an optional per-chart
// override. Skipped when no entries have names.
func drawLegend(
	ctx *render.Context,
	entries []legendEntry,
	th *theme.Theme,
	left, right, top, bottom float32,
	posOverride *theme.LegendPosition,
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

	ls := th.Legend

	// Resolve style fields with defaults.
	style := ls.TextStyle
	if !style.Color.IsSet() && style.Size == 0 {
		style = th.LabelStyle
	}
	bgColor := ls.Background
	if !bgColor.IsSet() {
		bgColor = gui.RGBA(0, 0, 0, 120)
	}
	swatchSize := ls.SwatchSize
	if swatchSize == 0 {
		swatchSize = 12
	}
	padding := ls.Padding
	if padding == 0 {
		padding = 6
	}
	itemGap := ls.ItemGap
	if itemGap == 0 {
		itemGap = 4
	}
	rowGap := ls.RowGap
	if rowGap == 0 {
		rowGap = 2
	}

	fh := ctx.FontHeight(style)
	rowH := max(fh, swatchSize)

	// Measure widest entry.
	maxW := float32(0)
	for _, e := range named {
		w := ctx.TextWidth(e.Name, style)
		maxW = max(maxW, w)
	}

	boxW := padding*2 + swatchSize + itemGap + maxW
	boxH := padding*2 +
		float32(len(named))*rowH +
		float32(len(named)-1)*rowGap

	// Determine position.
	pos := ls.Position
	if posOverride != nil {
		pos = *posOverride
	}
	var bx, by float32
	switch pos {
	case theme.LegendTopLeft:
		bx = left + 4
		by = top + 4
	case theme.LegendBottomRight:
		bx = right - boxW - 4
		by = bottom - boxH - 4
	case theme.LegendBottomLeft:
		bx = left + 4
		by = bottom - boxH - 4
	default: // LegendTopRight
		bx = right - boxW - 4
		by = top + 4
	}

	// Background.
	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4, bgColor)

	// Entries.
	for i, e := range named {
		ey := by + padding +
			float32(i)*(rowH+rowGap)
		// Color swatch.
		sx := bx + padding
		sy := ey + (rowH-swatchSize)/2
		ctx.FilledRoundedRect(sx, sy,
			swatchSize, swatchSize, 2, e.Color)
		// Label.
		tx := sx + swatchSize + itemGap
		ty := ey + (rowH-fh)/2
		ctx.Text(tx, ty, e.Name, style)
	}
}
