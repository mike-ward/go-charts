package main

import (
	"sort"

	"github.com/mike-ward/go-gui/gui"
)

func catalogPanel(w *gui.Window) gui.View {
	t := gui.CurrentTheme()
	app := gui.State[ShowcaseApp](w)
	entries := filteredEntries(app)

	switch {
	case len(entries) == 0:
		app.SelectedComponent = ""
		w.ScrollVerticalTo(scrollDetail, 0)
		w.ScrollHorizontalTo(scrollDetail, 0)
	case !hasEntry(entries, app.SelectedComponent):
		app.SelectedComponent = preferredComponentForGroup(
			app.SelectedGroup, entries)
		w.ScrollVerticalTo(scrollDetail, 0)
		w.ScrollHorizontalTo(scrollDetail, 0)
	}

	return gui.Column(gui.ContainerCfg{
		Width:   catalogWidth,
		Sizing:  gui.FixedFill,
		Color:   t.ColorPanel,
		Padding: gui.SomeP(12, 12, 12, 12),
		Spacing: gui.SomeF(8),
		Content: []gui.View{
			gui.Text(gui.TextCfg{
				Text:      "Chart Catalog",
				TextStyle: t.B3,
			}),
			searchInput(app),
			groupPicker(app),
			line(),
			gui.Column(gui.ContainerCfg{
				IDScroll: scrollCatalog,
				Sizing:   gui.FillFill,
				Padding: gui.Some(gui.Padding{
					Right: t.ScrollbarStyle.Size + 4,
				}),
				Spacing:       gui.SomeF(2),
				ScrollbarCfgY: &gui.ScrollbarCfg{GapEdge: 3},
				Content:       catalogRows(entries, app),
			}),
		},
	})
}

func searchInput(app *ShowcaseApp) gui.View {
	return gui.Input(gui.InputCfg{
		ID:          "showcase-nav-search",
		IDFocus:     focusSearch,
		Sizing:      gui.FillFit,
		Text:        app.NavQuery,
		Placeholder: "Search charts...",
		OnTextChanged: func(_ *gui.Layout, text string, w *gui.Window) {
			gui.State[ShowcaseApp](w).NavQuery = text
		},
	})
}

func groupPicker(app *ShowcaseApp) gui.View {
	items := make([]gui.View, len(demoGroups))
	for i, g := range demoGroups {
		items[i] = groupPickerItem(g.Label, g.Key, app)
	}
	return gui.Wrap(gui.ContainerCfg{
		Sizing:  gui.FillFit,
		Padding: gui.NoPadding,
		Spacing: gui.SomeF(3),
		Content: items,
	})
}

func groupPickerItem(label, key string, app *ShowcaseApp) gui.View {
	selected := app.SelectedGroup == key
	color := gui.CurrentTheme().ColorBackground
	if selected {
		color = gui.CurrentTheme().ColorActive
	}

	return gui.Button(gui.ButtonCfg{
		ID:          "grp-" + key,
		Color:       color,
		ColorBorder: color,
		Radius:      gui.SomeF(3),
		Padding:     gui.SomeP(3, 6, 3, 6),
		Content: []gui.View{
			gui.Text(gui.TextCfg{
				Text:      label,
				TextStyle: gui.CurrentTheme().N5,
			}),
		},
		OnClick: func(_ *gui.Layout, e *gui.Event, w *gui.Window) {
			sa := gui.State[ShowcaseApp](w)
			sa.SelectedGroup = key
			sa.NavQuery = ""
			entries := filteredEntries(sa)
			sa.SelectedComponent = preferredComponentForGroup(
				key, entries)
			w.ScrollVerticalTo(scrollCatalog, 0)
			w.ScrollVerticalTo(scrollDetail, 0)
			w.ScrollHorizontalTo(scrollDetail, 0)
			e.IsHandled = true
		},
	})
}

func catalogRows(entries []DemoEntry, app *ShowcaseApp) []gui.View {
	if len(entries) == 0 {
		return []gui.View{
			gui.Text(gui.TextCfg{
				Text:      "No matching charts",
				TextStyle: gui.CurrentTheme().N4,
			}),
		}
	}

	rows := make([]gui.View, 0, len(entries)+len(demoGroups)*2)
	for _, group := range demoGroups {
		if group.Key == groupAll {
			continue
		}
		groupEntries := make([]DemoEntry, 0, 8)
		for _, entry := range entries {
			if entry.Group == group.Key {
				groupEntries = append(groupEntries, entry)
			}
		}
		if len(groupEntries) == 0 {
			continue
		}
		sort.SliceStable(groupEntries, func(i, j int) bool {
			return entrySortBefore(groupEntries[i], groupEntries[j])
		})
		if len(rows) > 0 {
			rows = append(rows, gui.Row(gui.ContainerCfg{
				Height:  6,
				Sizing:  gui.FillFixed,
				Padding: gui.NoPadding,
			}))
		}
		rows = append(rows, gui.Text(gui.TextCfg{
			Text:      group.Label,
			TextStyle: gui.CurrentTheme().B5,
		}))
		for _, entry := range groupEntries {
			rows = append(rows, catalogRow(entry, app))
		}
	}
	return rows
}

func catalogRow(entry DemoEntry, app *ShowcaseApp) gui.View {
	selected := app.SelectedComponent == entry.ID
	color := gui.ColorTransparent
	if selected {
		color = gui.CurrentTheme().ColorActive
	}

	return gui.Button(gui.ButtonCfg{
		ID:               "cat-" + entry.ID,
		Sizing:           gui.FillFit,
		Color:            color,
		ColorHover:       gui.CurrentTheme().MenubarStyle.ColorSelect,
		ColorClick:       gui.CurrentTheme().ColorActive,
		ColorFocus:       color,
		ColorBorder:      gui.ColorTransparent,
		ColorBorderFocus: gui.ColorTransparent,
		Radius:           gui.SomeF(4),
		Padding:          gui.SomeP(3, 6, 3, 6),
		HAlign:           gui.Some(gui.HAlignLeft),
		Content: []gui.View{
			gui.Text(gui.TextCfg{
				Text:      entry.Label,
				TextStyle: gui.CurrentTheme().N4,
			}),
		},
		OnClick: func(_ *gui.Layout, e *gui.Event, w *gui.Window) {
			sa := gui.State[ShowcaseApp](w)
			sa.SelectedComponent = entry.ID
			w.ScrollVerticalTo(scrollDetail, 0)
			w.ScrollHorizontalTo(scrollDetail, 0)
			e.IsHandled = true
		},
	})
}
