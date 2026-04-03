package chart

import (
	"math"
	"strings"

	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// legendEntry describes one item in the chart legend.
type legendEntry struct {
	Name  string
	Color gui.Color
	Index int // original series index (for toggle mapping)
}

// legendBounds holds the pixel bounds of legend entries for
// click hit-testing.
type legendBounds struct {
	// EntryRects maps original series index to the entry's pixel
	// rect within the canvas.
	EntryRects []legendEntryRect
}

type legendEntryRect struct {
	Index               int // original series index
	X, Y, Width, Height float32
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
// override. Hidden entries are drawn dimmed with a strikethrough.
// Returns bounds for click hit-testing.
func drawLegend(
	ctx *render.Context,
	entries []legendEntry,
	th *theme.Theme,
	left, right, top, bottom float32,
	posOverride *theme.LegendPosition,
	hidden map[int]bool,
) legendBounds {
	// Filter to entries with names.
	named := make([]legendEntry, 0, len(entries))
	for _, e := range entries {
		if e.Name != "" {
			named = append(named, e)
		}
	}
	if len(named) == 0 {
		return legendBounds{}
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
	lb := legendBounds{
		EntryRects: make([]legendEntryRect, len(named)),
	}
	for i, e := range named {
		ey := by + padding +
			float32(i)*(rowH+rowGap)

		lb.EntryRects[i] = legendEntryRect{
			Index:  e.Index,
			X:      bx,
			Y:      ey,
			Width:  boxW,
			Height: rowH,
		}

		isHidden := hidden[e.Index]
		color := e.Color
		textStyle := style
		if isHidden {
			color = dimColor(color, HoverDimAlpha)
			textStyle.Color = gui.RGBA(
				style.Color.R, style.Color.G,
				style.Color.B, HoverDimAlpha)
		}

		// Color swatch.
		sx := bx + padding
		sy := ey + (rowH-swatchSize)/2
		ctx.FilledRoundedRect(sx, sy,
			swatchSize, swatchSize, 2, color)

		// Strikethrough for hidden entries.
		if isHidden {
			mid := sy + swatchSize/2
			ctx.Line(sx-1, mid, sx+swatchSize+1, mid,
				gui.RGBA(200, 200, 200, 180), 1.5)
		}

		// Label.
		tx := sx + swatchSize + itemGap
		ty := ey + (rowH-fh)/2
		ctx.Text(tx, ty, e.Name, textStyle)
	}

	return lb
}

// legendHitTest returns the original series index of the legend
// entry under (mx, my), or -1 if none.
func legendHitTest(lb legendBounds, mx, my float32) int {
	for _, r := range lb.EntryRects {
		if mx >= r.X && mx <= r.X+r.Width &&
			my >= r.Y && my <= r.Y+r.Height {
			return r.Index
		}
	}
	return -1
}

// drawTooltip draws a tooltip directly on the canvas near (tx, ty).
// label may contain '\n'-separated lines.
func drawTooltip(
	ctx *render.Context, tx, ty float32, label string, th *theme.Theme,
) {
	lines := strings.Split(label, "\n")
	textStyle := th.TickStyle
	textStyle.Color = gui.Hex(0xEEEEEE)

	fh := ctx.FontHeight(textStyle)
	const padding = float32(6)
	const lineGap = float32(2)

	// Measure box dimensions.
	maxW := float32(0)
	for _, ln := range lines {
		w := ctx.TextWidth(ln, textStyle)
		maxW = max(maxW, w)
	}
	boxW := maxW + padding*2
	boxH := float32(len(lines))*fh +
		float32(len(lines)-1)*lineGap + padding*2

	// Position above-right of (tx, ty), clamped to canvas.
	bx := tx + 8
	by := ty - 8 - boxH
	if bx+boxW > ctx.Width() {
		bx = ctx.Width() - boxW
	}
	if bx < 0 {
		bx = 0
	}
	if by < 0 {
		by = 0
	}
	if by+boxH > ctx.Height() {
		by = ctx.Height() - boxH
	}

	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4,
		gui.RGBA(20, 20, 20, 220))

	for i, ln := range lines {
		lx := bx + padding
		ly := by + padding + float32(i)*(fh+lineGap)
		ctx.Text(lx, ly, ln, textStyle)
	}
}

// drawCrosshair draws dashed vertical and horizontal tracking lines at
// (mx, my). Does nothing when the cursor is outside the plot area.
func drawCrosshair(
	ctx *render.Context, th *theme.Theme,
	mx, my, left, right, top, bottom float32,
) {
	if mx < left || mx > right || my < top || my > bottom {
		return
	}
	cs := th.Crosshair
	color := cs.Color
	if !color.IsSet() {
		color = gui.RGBA(128, 128, 128, 160)
	}
	width := cs.Width
	if width == 0 {
		width = 1
	}
	dashLen := cs.DashLen
	if dashLen == 0 {
		dashLen = 6
	}
	gapLen := cs.GapLen
	if gapLen == 0 {
		gapLen = 4
	}
	ctx.DashedLine(mx, top, mx, bottom, color, width, dashLen, gapLen)
	ctx.DashedLine(left, my, right, my, color, width, dashLen, gapLen)
}
