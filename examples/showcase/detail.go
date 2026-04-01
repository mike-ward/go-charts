package main

import "github.com/mike-ward/go-gui/gui"

func detailPanel(w *gui.Window, entries []DemoEntry) gui.View {
	t := gui.CurrentTheme()
	app := gui.State[ShowcaseApp](w)

	if len(entries) == 0 {
		return gui.Column(gui.ContainerCfg{
			IDScroll:   scrollDetail,
			Sizing:     gui.FillFill,
			SizeBorder: gui.NoBorder,
			Padding:    gui.Some(detailPadding()),
			ScrollbarCfgY: &gui.ScrollbarCfg{
				GapEdge: 4,
			},
			Content: []gui.View{
				gui.Text(gui.TextCfg{
					Text:      "No chart matches filter",
					TextStyle: t.B2,
				}),
			},
		})
	}

	entry := selectedEntry(entries, app.SelectedComponent)

	content := componentDemo(w, entry.ID)

	return gui.Column(gui.ContainerCfg{
		IDScroll:   scrollDetail,
		Sizing:     gui.FillFill,
		Color:      t.ColorBackground,
		SizeBorder: gui.NoBorder,
		Padding:    gui.Some(detailPadding()),
		Spacing:    gui.Some(t.SpacingLarge),
		ScrollbarCfgY: &gui.ScrollbarCfg{
			GapEdge: 4,
		},
		Content: []gui.View{
			viewTitleBar(entry),
			gui.Text(gui.TextCfg{
				Text:      entry.Summary,
				TextStyle: t.N3,
				Mode:      gui.TextModeWrap,
			}),
			content,
		},
	})
}

func detailPadding() gui.Padding {
	base := gui.CurrentTheme().PaddingLarge
	base.Right += gui.CurrentTheme().ScrollbarStyle.Size + 4
	return base
}

func viewTitleBar(entry DemoEntry) gui.View {
	return gui.Column(gui.ContainerCfg{
		Sizing:     gui.FillFit,
		Spacing:    gui.NoSpacing,
		Padding:    gui.NoPadding,
		SizeBorder: gui.NoBorder,
		Content: []gui.View{
			gui.Text(gui.TextCfg{
				Text:      entry.Label,
				TextStyle: gui.CurrentTheme().B1,
			}),
			line(),
		},
	})
}

var componentDemos = map[string]func(*gui.Window) gui.View{
	"type_basecfg":    demoBaseCfg,
	"type_series_xy":  demoSeriesXY,
	"type_series_cat": demoSeriesCategory,
	"type_theme":      demoTheme,
	"type_axis":       demoAxisLinear,
	"line_basic":      demoLineBasic,
	"line_markers":    demoLineMarkers,
	"line_area":       demoLineArea,
	"line_multi":      demoLineMulti,
	"bar_basic":       demoBarBasic,
	"bar_single":      demoBarSingle,
	"bar_wide":        demoBarWide,
	"bar_rounded":     demoBarRounded,
	"pie_basic":       demoPie,
	"pie_donut":       demoDonut,
	"area_basic":      demoArea,
	"area_stacked":    demoAreaStacked,
	"scatter_basic":   demoScatter,
	"scatter_markers": demoScatterMarkers,
}

func componentDemo(w *gui.Window, id string) gui.View {
	if fn, ok := componentDemos[id]; ok {
		return fn(w)
	}
	return demoPlaceholder("Demo: " + id)
}

func demoPlaceholder(text string) gui.View {
	t := gui.CurrentTheme()
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFit,
		Color:   t.ColorPanel,
		Padding: gui.SomeP(24, 24, 24, 24),
		Radius:  gui.SomeF(8),
		Content: []gui.View{
			gui.Text(gui.TextCfg{
				Text:      text,
				TextStyle: t.N3,
				Mode:      gui.TextModeWrap,
			}),
		},
	})
}
