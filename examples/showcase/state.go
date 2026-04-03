package main

import (
	"sort"
	"strings"
)

const (
	groupAll         = "all"
	groupTypes       = "types"
	groupLine        = "line"
	groupBar         = "bar"
	groupPie         = "pie"
	groupGauge       = "gauge"
	groupArea        = "area"
	groupScatter     = "scatter"
	groupCandlestick = "candlestick"
	groupStyles      = "styles"
)

// ShowcaseApp holds all state for the charts showcase.
type ShowcaseApp struct {
	NavQuery          string
	SelectedGroup     string
	SelectedComponent string
}

func newShowcaseApp() *ShowcaseApp {
	return &ShowcaseApp{
		SelectedGroup:     groupAll,
		SelectedComponent: "line_basic",
	}
}

// DemoEntry describes one demo in the catalog.
type DemoEntry struct {
	ID      string
	Label   string
	Group   string
	Summary string
	Tags    []string

	idLower    string
	labelLower string
}

// DemoGroup describes a catalog group in the left pane.
type DemoGroup struct {
	Key   string
	Label string
}

var demoGroups = []DemoGroup{
	{groupAll, "All"},
	{groupTypes, "Types"},
	{groupLine, "Line"},
	{groupBar, "Bar"},
	{groupPie, "Pie"},
	{groupGauge, "Gauge"},
	{groupArea, "Area"},
	{groupScatter, "Scatter"},
	{groupCandlestick, "Candlestick"},
	{groupStyles, "Styles"},
}

var demoEntries = []DemoEntry{
	// Types
	{ID: "type_basecfg", Label: "BaseCfg", Group: groupTypes, Summary: "Common configuration fields embedded by all chart types.", Tags: []string{"base", "config", "id", "title", "sizing", "theme", "version"}},
	{ID: "type_series_xy", Label: "series.XY", Group: groupTypes, Summary: "XY data series for Line, Area, and Scatter charts.", Tags: []string{"series", "xy", "point", "data", "slices"}},
	{ID: "type_series_cat", Label: "series.Category", Group: groupTypes, Summary: "Categorical data series for Bar charts.", Tags: []string{"series", "category", "label", "value", "map"}},
	{ID: "type_theme", Label: "theme.Theme", Group: groupTypes, Summary: "Visual style: colors, fonts, padding, palette.", Tags: []string{"theme", "palette", "color", "style", "tableau", "pastel", "vivid"}},
	{ID: "type_axis", Label: "axis.Linear", Group: groupTypes, Summary: "Linear numeric axis with auto-tick generation.", Tags: []string{"axis", "linear", "tick", "range", "format", "auto"}},

	// Line
	{ID: "line_basic", Label: "Basic Line", Group: groupLine, Summary: "Monthly revenue comparison across two years.", Tags: []string{"line", "multi", "revenue"}},
	{ID: "line_markers", Label: "Line with Markers", Group: groupLine, Summary: "Temperature readings with visible data points.", Tags: []string{"line", "markers", "temperature"}},
	{ID: "line_area", Label: "Line with Area Fill", Group: groupLine, Summary: "Website traffic with shaded area under the curve.", Tags: []string{"line", "area", "fill", "traffic"}},
	{ID: "line_multi", Label: "Multi-Series Line", Group: groupLine, Summary: "Three series with custom line width.", Tags: []string{"line", "multi", "series"}},

	// Bar
	{ID: "bar_basic", Label: "Basic Bar", Group: groupBar, Summary: "Regional sales comparison across quarters.", Tags: []string{"bar", "grouped", "sales"}},
	{ID: "bar_single", Label: "Single Series Bar", Group: groupBar, Summary: "Monthly rainfall totals.", Tags: []string{"bar", "single", "rainfall"}},
	{ID: "bar_wide", Label: "Fixed Width Bars", Group: groupBar, Summary: "Department headcount with fixed bar width.", Tags: []string{"bar", "fixed", "width", "headcount"}},
	{ID: "bar_rounded", Label: "Rounded Bars", Group: groupBar, Summary: "Product category revenue with rounded corners.", Tags: []string{"bar", "rounded", "radius"}},
	{ID: "bar_horizontal", Label: "Horizontal Bar", Group: groupBar, Summary: "Survey results with bars drawn left-to-right.", Tags: []string{"bar", "horizontal", "survey"}},
	{ID: "bar_stacked", Label: "Stacked Bar", Group: groupBar, Summary: "Traffic by channel stacked per month.", Tags: []string{"bar", "stacked", "channel", "traffic"}},

	// Pie
	{ID: "pie_basic", Label: "Basic Pie", Group: groupPie, Summary: "Browser market share distribution.", Tags: []string{"pie", "share", "percent"}},
	{ID: "pie_donut", Label: "Donut Chart", Group: groupPie, Summary: "Budget allocation with inner radius.", Tags: []string{"pie", "donut", "budget"}},

	// Gauge
	{ID: "gauge_basic", Label: "Gauge with Zones", Group: groupGauge, Summary: "CPU usage gauge with colored warning zones.", Tags: []string{"gauge", "arc", "zone", "value"}},
	{ID: "gauge_simple", Label: "Simple Gauge", Group: groupGauge, Summary: "Completion percentage with no zones.", Tags: []string{"gauge", "simple", "percent"}},

	// Area
	{ID: "area_basic", Label: "Basic Area", Group: groupArea, Summary: "User signups over time.", Tags: []string{"area", "signups", "time"}},
	{ID: "area_stacked", Label: "Stacked Area", Group: groupArea, Summary: "Revenue breakdown by product line.", Tags: []string{"area", "stacked", "revenue"}},

	// Scatter
	{ID: "scatter_basic", Label: "Basic Scatter", Group: groupScatter, Summary: "Height versus weight correlation.", Tags: []string{"scatter", "correlation", "points"}},
	{ID: "scatter_markers", Label: "Marker Shapes", Group: groupScatter, Summary: "Wind speed versus temperature with different marker shapes.", Tags: []string{"scatter", "markers", "shapes"}},

	// Candlestick
	{ID: "candlestick_basic", Label: "Basic Candlestick", Group: groupCandlestick, Summary: "AAPL price action with up/down candle colors.", Tags: []string{"candlestick", "ohlc", "price", "stock", "finance"}},

	// Styles
	{ID: "style_palette", Label: "Palette Swap", Group: groupStyles, Summary: "Same data rendered with Tableau 10, Pastel, and Vivid palettes.", Tags: []string{"style", "palette", "theme", "color", "tableau", "pastel", "vivid"}},
	{ID: "style_tick_marks", Label: "Tick Marks", Group: groupStyles, Summary: "Custom tick mark length, color, and width.", Tags: []string{"style", "tick", "marks", "axis"}},
	{ID: "style_legend_pos", Label: "Legend Positions", Group: groupStyles, Summary: "Legend placed in each of the four corners.", Tags: []string{"style", "legend", "position", "corner"}},
	{ID: "style_legend_cfg", Label: "Legend Styling", Group: groupStyles, Summary: "Custom swatch size, padding, gaps, and background.", Tags: []string{"style", "legend", "swatch", "padding", "gap"}},
	{ID: "style_rotation", Label: "Rotated Labels", Group: groupStyles, Summary: "X-axis tick labels rotated for long category names.", Tags: []string{"style", "rotation", "labels", "tick", "angle"}},
	{ID: "style_padding", Label: "Custom Padding", Group: groupStyles, Summary: "Tight versus spacious chart padding.", Tags: []string{"style", "padding", "spacing", "inset"}},
	{ID: "style_kitchen", Label: "Kitchen Sink", Group: groupStyles, Summary: "All style knobs combined on a single chart.", Tags: []string{"style", "combined", "kitchen", "sink", "all"}},
}

func init() {
	for i := range demoEntries {
		demoEntries[i].idLower = strings.ToLower(demoEntries[i].ID)
		demoEntries[i].labelLower = strings.ToLower(demoEntries[i].Label)
	}
	sort.SliceStable(demoEntries, func(i, j int) bool {
		return entrySortBefore(demoEntries[i], demoEntries[j])
	})
}

func entryMatchesQuery(entry DemoEntry, query string) bool {
	if query == "" {
		return true
	}
	q := strings.ToLower(query)
	if strings.Contains(entry.idLower, q) ||
		strings.Contains(entry.labelLower, q) ||
		strings.Contains(entry.Group, q) {
		return true
	}
	for _, tag := range entry.Tags {
		if strings.Contains(tag, q) {
			return true
		}
	}
	return strings.Contains(strings.ToLower(entry.Summary), q)
}

func filteredEntries(app *ShowcaseApp) []DemoEntry {
	out := make([]DemoEntry, 0, len(demoEntries))
	for _, entry := range demoEntries {
		if app.SelectedGroup != groupAll && entry.Group != app.SelectedGroup {
			continue
		}
		if !entryMatchesQuery(entry, app.NavQuery) {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func hasEntry(entries []DemoEntry, selected string) bool {
	for _, entry := range entries {
		if entry.ID == selected {
			return true
		}
	}
	return false
}

func selectedEntry(entries []DemoEntry, selected string) DemoEntry {
	for _, entry := range entries {
		if entry.ID == selected {
			return entry
		}
	}
	if len(entries) == 0 {
		return DemoEntry{}
	}
	return entries[0]
}

func preferredComponentForGroup(entries []DemoEntry) string {
	if len(entries) == 0 {
		return ""
	}
	best := entries[0]
	for _, entry := range entries[1:] {
		if entrySortBefore(entry, best) {
			best = entry
		}
	}
	return best.ID
}

func entrySortBefore(a, b DemoEntry) bool {
	if a.labelLower != b.labelLower {
		return a.labelLower < b.labelLower
	}
	return a.idLower < b.idLower
}
