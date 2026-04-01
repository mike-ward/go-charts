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
		Title:  "Line Chart",
		Width:  800,
		Height: 600,
		OnInit: func(w *gui.Window) {
			w.UpdateView(view)
		},
	})

	backend.Run(w)
}

func view(w *gui.Window) gui.View {
	return chart.Line(chart.LineCfg{
		BaseCfg: chart.BaseCfg{
			ID:    "line-demo",
			Title: "Monthly Revenue",
		},
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "2025",
				Color: gui.Hex(0x4E79A7),
				Points: []series.Point{
					{X: 1, Y: 12},
					{X: 2, Y: 19},
					{X: 3, Y: 15},
					{X: 4, Y: 28},
					{X: 5, Y: 24},
					{X: 6, Y: 31},
				},
			}),
			series.NewXY(series.XYCfg{
				Name:  "2024",
				Color: gui.Hex(0xF28E2B),
				Points: []series.Point{
					{X: 1, Y: 8},
					{X: 2, Y: 14},
					{X: 3, Y: 11},
					{X: 4, Y: 22},
					{X: 5, Y: 19},
					{X: 6, Y: 26},
				},
			}),
		},
	})
}
