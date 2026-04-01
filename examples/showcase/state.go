package main

import "strings"

const (
	groupAll     = "all"
	groupTypes   = "types"
	groupLine    = "line"
	groupBar     = "bar"
	groupPie     = "pie"
	groupArea    = "area"
	groupScatter = "scatter"
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
	{groupArea, "Area"},
	{groupScatter, "Scatter"},
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
	{ID: "bar_wide", Label: "Wide Bars", Group: groupBar, Summary: "Department headcount with custom bar width.", Tags: []string{"bar", "wide", "headcount"}},
	{ID: "bar_rounded", Label: "Rounded Bars", Group: groupBar, Summary: "Product category revenue with rounded corners.", Tags: []string{"bar", "rounded", "radius"}},

	// Pie (stubs)
	{ID: "pie_basic", Label: "Basic Pie", Group: groupPie, Summary: "Browser market share distribution.", Tags: []string{"pie", "share", "percent"}},
	{ID: "pie_donut", Label: "Donut Chart", Group: groupPie, Summary: "Budget allocation with inner radius.", Tags: []string{"pie", "donut", "budget"}},

	// Area (stubs)
	{ID: "area_basic", Label: "Basic Area", Group: groupArea, Summary: "User signups over time.", Tags: []string{"area", "signups", "time"}},
	{ID: "area_stacked", Label: "Stacked Area", Group: groupArea, Summary: "Revenue breakdown by product line.", Tags: []string{"area", "stacked", "revenue"}},

	// Scatter (stubs)
	{ID: "scatter_basic", Label: "Basic Scatter", Group: groupScatter, Summary: "Height versus weight correlation.", Tags: []string{"scatter", "correlation", "points"}},
	{ID: "scatter_markers", Label: "Marker Shapes", Group: groupScatter, Summary: "Wind speed versus temperature with different marker shapes.", Tags: []string{"scatter", "markers", "shapes"}},
}

func init() {
	for i := range demoEntries {
		demoEntries[i].idLower = strings.ToLower(demoEntries[i].ID)
		demoEntries[i].labelLower = strings.ToLower(demoEntries[i].Label)
	}
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

func preferredComponentForGroup(_ string, entries []DemoEntry) string {
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
