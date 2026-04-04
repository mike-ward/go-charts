package chart

import (
	"math"
	"strings"

	"github.com/mike-ward/go-charts/axis"
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

// resolveLeft returns left expanded so that Y-axis tick labels and
// the rotated Y-axis title don't overlap. Call after the Y axis
// domain is set. Returns left unchanged when yAxis has no label.
func resolveLeft(
	ctx *render.Context, th *theme.Theme,
	left, bottom, top float32,
	yAxis *axis.Linear,
) float32 {
	if yAxis == nil || yAxis.Label() == "" {
		return left
	}
	ticks := yAxis.Ticks(bottom, top)
	maxTickW := float32(0)
	for _, t := range ticks {
		maxTickW = max(maxTickW, ctx.TextWidth(t.Label, th.TickStyle))
	}
	tickLen, _, _ := resolvedTickMark(th)
	labelFH := ctx.FontHeight(th.LabelStyle)
	needed := labelFH*2 + tickLen + maxTickW + 4
	return max(left, needed)
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

// maxTickLabelWidth returns the widest tick label in pixels.
func maxTickLabelWidth(
	ctx *render.Context, ticks []axis.Tick, style gui.TextStyle,
) float32 {
	maxW := float32(0)
	for _, t := range ticks {
		maxW = max(maxW, ctx.TextWidth(t.Label, style))
	}
	return maxW
}

// resolveBottom computes the bottom padding needed for X-axis
// decorations: tick marks, tick labels (accounting for rotation),
// and the optional axis title. Returns the pixel distance from
// the canvas bottom to the X-axis line.
func resolveBottom(
	ctx *render.Context, th *theme.Theme,
	maxTickLabelWidth, xTickRotation float32,
	xAxisLabel string,
) float32 {
	tickLen, _, _ := resolvedTickMark(th)
	tickFH := ctx.FontHeight(th.TickStyle)

	var tickLabelH float32
	if xTickRotation != 0 {
		sinA := float32(math.Abs(math.Sin(float64(xTickRotation))))
		cosA := float32(math.Abs(math.Cos(float64(xTickRotation))))
		tickLabelH = sinA*maxTickLabelWidth + cosA*tickFH
	} else {
		tickLabelH = tickFH
	}

	// tick marks + gap + labels + margin
	needed := tickLen + 2 + tickLabelH + 6
	if xAxisLabel != "" {
		labelFH := ctx.FontHeight(th.LabelStyle)
		needed += 6 + labelFH
	}
	return needed
}

// legendRightReserve returns the horizontal space to subtract
// from the right edge of the plot area when LegendRight is
// active. Returns 0 for all other positions.
func legendRightReserve(
	ctx *render.Context, th *theme.Theme,
	posOverride *theme.LegendPosition,
	names []string,
) float32 {
	pos := th.Legend.Position
	if posOverride != nil {
		pos = *posOverride
	}
	if pos != theme.LegendRight {
		return 0
	}

	ls := th.Legend
	style := ls.TextStyle
	if !style.Color.IsSet() && style.Size == 0 {
		style = th.LabelStyle
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

	maxW := float32(0)
	any := false
	for _, n := range names {
		if n == "" {
			continue
		}
		any = true
		w := ctx.TextWidth(n, style)
		maxW = max(maxW, w)
	}
	if !any {
		return 0
	}
	// box width + gap between plot edge and legend
	return padding*2 + swatchSize + itemGap + maxW + 8
}

// legendTopReserve returns the vertical space to add to the top
// edge of the plot area when LegendTop is active. Returns 0 for
// all other positions.
func legendTopReserve(
	ctx *render.Context, th *theme.Theme,
	posOverride *theme.LegendPosition,
	names []string,
	left, right float32,
) float32 {
	pos := th.Legend.Position
	if posOverride != nil {
		pos = *posOverride
	}
	if pos != theme.LegendTop {
		return 0
	}

	ls := th.Legend
	style := ls.TextStyle
	if !style.Color.IsSet() && style.Size == 0 {
		style = th.LabelStyle
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

	// Count rows needed for the horizontal layout.
	const interItemGap = float32(12)
	availW := right - left
	nRows := 1
	rowW := float32(0)
	any := false
	for _, n := range names {
		if n == "" {
			continue
		}
		any = true
		tw := ctx.TextWidth(n, style)
		w := swatchSize + itemGap + tw
		addition := w
		if rowW > 0 {
			addition += interItemGap
		}
		if rowW > 0 && rowW+addition > availW {
			nRows++
			rowW = w
		} else {
			rowW += addition
		}
	}
	if !any {
		return 0
	}
	// box height + gap below legend
	return padding*2 +
		float32(nRows)*rowH +
		float32(max(nRows-1, 0))*rowGap + 8
}

// legendBottomReserve returns the vertical space to add to the
// bottom edge of the plot area when LegendBottom is active.
// Returns 0 for all other positions.
func legendBottomReserve(
	ctx *render.Context, th *theme.Theme,
	posOverride *theme.LegendPosition,
	names []string,
	left, right float32,
) float32 {
	pos := th.Legend.Position
	if posOverride != nil {
		pos = *posOverride
	}
	if pos != theme.LegendBottom {
		return 0
	}

	ls := th.Legend
	style := ls.TextStyle
	if !style.Color.IsSet() && style.Size == 0 {
		style = th.LabelStyle
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

	const interItemGap = float32(12)
	availW := right - left
	nRows := 1
	rowW := float32(0)
	any := false
	for _, n := range names {
		if n == "" {
			continue
		}
		any = true
		tw := ctx.TextWidth(n, style)
		w := swatchSize + itemGap + tw
		addition := w
		if rowW > 0 {
			addition += interItemGap
		}
		if rowW > 0 && rowW+addition > availW {
			nRows++
			rowW = w
		} else {
			rowW += addition
		}
	}
	if !any {
		return 0
	}
	return padding*2 +
		float32(nRows)*rowH +
		float32(max(nRows-1, 0))*rowGap
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

	// Determine position.
	pos := ls.Position
	if posOverride != nil {
		pos = *posOverride
	}

	if pos == theme.LegendNone {
		return legendBounds{}
	}

	fh := ctx.FontHeight(style)
	rowH := max(fh, swatchSize)

	if pos == theme.LegendBottom {
		return drawLegendBottom(ctx, named, hidden,
			style, bgColor, swatchSize, padding,
			itemGap, rowGap, fh, rowH, left, right)
	}

	if pos == theme.LegendRight {
		return drawLegendRight(ctx, named, hidden,
			style, bgColor, swatchSize, padding,
			itemGap, rowGap, fh, rowH, right, top)
	}

	if pos == theme.LegendTop {
		return drawLegendTop(ctx, named, hidden,
			style, bgColor, swatchSize, padding,
			itemGap, rowGap, fh, rowH, left, right, top)
	}

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

// drawLegendBottom renders a horizontal legend centered below the
// plot area. Items wrap to additional rows when they exceed the
// available width.
func drawLegendBottom(
	ctx *render.Context,
	entries []legendEntry,
	hidden map[int]bool,
	style gui.TextStyle,
	bgColor gui.Color,
	swatchSize, padding, itemGap, rowGap,
	fh, rowH, left, right float32,
) legendBounds {
	const interItemGap = float32(12)

	// Measure each item's width.
	type itemInfo struct {
		width float32
	}
	items := make([]itemInfo, len(entries))
	for i, e := range entries {
		tw := ctx.TextWidth(e.Name, style)
		items[i] = itemInfo{width: swatchSize + itemGap + tw}
	}

	// Layout rows, wrapping when the row exceeds available width.
	availW := right - left
	type layoutRow struct {
		start, end int // index range [start, end)
		width      float32
	}
	var rows []layoutRow
	rowStart := 0
	rowW := float32(0)
	for i, item := range items {
		addition := item.width
		if i > rowStart {
			addition += interItemGap
		}
		if i > rowStart && rowW+addition > availW {
			rows = append(rows, layoutRow{rowStart, i, rowW})
			rowStart = i
			rowW = item.width
		} else {
			rowW += addition
		}
	}
	if rowStart < len(items) {
		rows = append(rows, layoutRow{rowStart, len(items), rowW})
	}

	// Compute box dimensions.
	maxRowW := float32(0)
	for _, r := range rows {
		maxRowW = max(maxRowW, r.width)
	}
	boxW := maxRowW + padding*2
	boxH := padding*2 +
		float32(len(rows))*rowH +
		float32(max(len(rows)-1, 0))*rowGap

	// Position centered horizontally, at the bottom of the canvas.
	bx := (left + right - boxW) / 2
	by := ctx.Height() - boxH

	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4, bgColor)

	// Draw entries row by row.
	lb := legendBounds{
		EntryRects: make([]legendEntryRect, len(entries)),
	}
	for ri, r := range rows {
		ey := by + padding + float32(ri)*(rowH+rowGap)
		// Center this row within the box.
		x := bx + padding + (maxRowW-r.width)/2
		for i := r.start; i < r.end; i++ {
			if i > r.start {
				x += interItemGap
			}
			e := entries[i]
			lb.EntryRects[i] = legendEntryRect{
				Index:  e.Index,
				X:      x,
				Y:      ey,
				Width:  items[i].width,
				Height: rowH,
			}

			isHidden := hidden[e.Index]
			color := e.Color
			ts := style
			if isHidden {
				color = dimColor(color, HoverDimAlpha)
				ts.Color = gui.RGBA(
					style.Color.R, style.Color.G,
					style.Color.B, HoverDimAlpha)
			}

			// Color swatch.
			sy := ey + (rowH-swatchSize)/2
			ctx.FilledRoundedRect(x, sy,
				swatchSize, swatchSize, 2, color)

			if isHidden {
				mid := sy + swatchSize/2
				ctx.Line(x-1, mid, x+swatchSize+1, mid,
					gui.RGBA(200, 200, 200, 180), 1.5)
			}

			// Label.
			tx := x + swatchSize + itemGap
			ty := ey + (rowH-fh)/2
			ctx.Text(tx, ty, e.Name, ts)

			x += items[i].width
		}
	}

	return lb
}

// drawLegendTop renders a horizontal legend centered between the
// title and the plot area. Items wrap to additional rows when
// they exceed the available width. top is the adjusted plot-area
// top (after legendTopReserve has been applied).
func drawLegendTop(
	ctx *render.Context,
	entries []legendEntry,
	hidden map[int]bool,
	style gui.TextStyle,
	bgColor gui.Color,
	swatchSize, padding, itemGap, rowGap,
	fh, rowH, left, right, top float32,
) legendBounds {
	const interItemGap = float32(12)

	type itemInfo struct{ width float32 }
	items := make([]itemInfo, len(entries))
	for i, e := range entries {
		tw := ctx.TextWidth(e.Name, style)
		items[i] = itemInfo{width: swatchSize + itemGap + tw}
	}

	availW := right - left
	type layoutRow struct {
		start, end int
		width      float32
	}
	var rows []layoutRow
	rowStart := 0
	rowW := float32(0)
	for i, item := range items {
		addition := item.width
		if i > rowStart {
			addition += interItemGap
		}
		if i > rowStart && rowW+addition > availW {
			rows = append(rows, layoutRow{rowStart, i, rowW})
			rowStart = i
			rowW = item.width
		} else {
			rowW += addition
		}
	}
	if rowStart < len(items) {
		rows = append(rows, layoutRow{rowStart, len(items), rowW})
	}

	maxRowW := float32(0)
	for _, r := range rows {
		maxRowW = max(maxRowW, r.width)
	}
	boxW := maxRowW + padding*2
	boxH := padding*2 +
		float32(len(rows))*rowH +
		float32(max(len(rows)-1, 0))*rowGap

	// Position centered horizontally, just above the plot area.
	bx := (left + right - boxW) / 2
	by := top - boxH - 8

	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4, bgColor)

	lb := legendBounds{
		EntryRects: make([]legendEntryRect, len(entries)),
	}
	for ri, r := range rows {
		ey := by + padding + float32(ri)*(rowH+rowGap)
		x := bx + padding + (maxRowW-r.width)/2
		for i := r.start; i < r.end; i++ {
			if i > r.start {
				x += interItemGap
			}
			e := entries[i]
			lb.EntryRects[i] = legendEntryRect{
				Index:  e.Index,
				X:      x,
				Y:      ey,
				Width:  items[i].width,
				Height: rowH,
			}

			isHidden := hidden[e.Index]
			color := e.Color
			ts := style
			if isHidden {
				color = dimColor(color, HoverDimAlpha)
				ts.Color = gui.RGBA(
					style.Color.R, style.Color.G,
					style.Color.B, HoverDimAlpha)
			}

			sy := ey + (rowH-swatchSize)/2
			ctx.FilledRoundedRect(x, sy,
				swatchSize, swatchSize, 2, color)

			if isHidden {
				mid := sy + swatchSize/2
				ctx.Line(x-1, mid, x+swatchSize+1, mid,
					gui.RGBA(200, 200, 200, 180), 1.5)
			}

			tx := x + swatchSize + itemGap
			ty := ey + (rowH-fh)/2
			ctx.Text(tx, ty, e.Name, ts)

			x += items[i].width
		}
	}

	return lb
}

// drawLegendRight renders a vertical legend outside the plot area,
// to the right and top-aligned.
func drawLegendRight(
	ctx *render.Context,
	entries []legendEntry,
	hidden map[int]bool,
	style gui.TextStyle,
	bgColor gui.Color,
	swatchSize, padding, itemGap, rowGap,
	fh, rowH, right, top float32,
) legendBounds {
	// Measure widest entry.
	maxW := float32(0)
	for _, e := range entries {
		w := ctx.TextWidth(e.Name, style)
		maxW = max(maxW, w)
	}

	boxW := padding*2 + swatchSize + itemGap + maxW
	boxH := padding*2 +
		float32(len(entries))*rowH +
		float32(len(entries)-1)*rowGap

	// Position to the right of the plot area, top-aligned.
	bx := right + 8
	by := top

	// Clamp to canvas width.
	if bx+boxW > ctx.Width() {
		bx = ctx.Width() - boxW
	}

	ctx.FilledRoundedRect(bx, by, boxW, boxH, 4, bgColor)

	lb := legendBounds{
		EntryRects: make([]legendEntryRect, len(entries)),
	}
	for i, e := range entries {
		ey := by + padding + float32(i)*(rowH+rowGap)

		lb.EntryRects[i] = legendEntryRect{
			Index:  e.Index,
			X:      bx,
			Y:      ey,
			Width:  boxW,
			Height: rowH,
		}

		isHidden := hidden[e.Index]
		color := e.Color
		ts := style
		if isHidden {
			color = dimColor(color, HoverDimAlpha)
			ts.Color = gui.RGBA(
				style.Color.R, style.Color.G,
				style.Color.B, HoverDimAlpha)
		}

		sx := bx + padding
		sy := ey + (rowH-swatchSize)/2
		ctx.FilledRoundedRect(sx, sy,
			swatchSize, swatchSize, 2, color)

		if isHidden {
			mid := sy + swatchSize/2
			ctx.Line(sx-1, mid, sx+swatchSize+1, mid,
				gui.RGBA(200, 200, 200, 180), 1.5)
		}

		tx := sx + swatchSize + itemGap
		ty := ey + (rowH-fh)/2
		ctx.Text(tx, ty, e.Name, ts)
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
