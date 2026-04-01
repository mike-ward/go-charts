package main

import (
	"github.com/mike-ward/go-charts/chart"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-gui/gui"
	"github.com/mike-ward/go-gui/gui/backend"
)

type App struct{}

func main() {
	gui.SetTheme(gui.ThemeDarkBordered)

	w := gui.NewWindow(gui.WindowCfg{
		State:  &App{},
		Title:  "Bar Chart",
		Width:  800,
		Height: 600,
		OnInit: func(w *gui.Window) {
			w.UpdateView(view)
		},
	})

	backend.Run(w)
}

func view(w *gui.Window) gui.View {
	return chart.Bar(chart.BarCfg{
		BaseCfg: chart.BaseCfg{
			ID:    "bar-demo",
			Title: "Sales by Region",
		},
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name:  "Q1",
				Color: gui.Hex(0x4E79A7),
				Values: []series.CategoryValue{
					{Label: "North", Value: 45},
					{Label: "South", Value: 32},
					{Label: "East", Value: 58},
					{Label: "West", Value: 41},
				},
			}),
		},
	})
}
