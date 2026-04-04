package main

import (
	"github.com/mike-ward/go-charts/chart"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-gui/gui"
)

func demoScatter(w *gui.Window) gui.View {
	return demoWithCode(w, "scatter-basic", chart.Scatter(chart.ScatterCfg{
		BaseCfg: chart.BaseCfg{
			ID:             "scatter-basic",
			Title:          "Height vs Weight",
			Sizing:         gui.FillFixed,
			Height:         350,
			LegendPosition: &posBottom,
		},
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "Subjects",
				Color: gui.Hex(0x4E79A7),
				Points: []series.Point{
					{X: 155, Y: 52}, {X: 160, Y: 58}, {X: 162, Y: 55},
					{X: 165, Y: 62}, {X: 167, Y: 60}, {X: 168, Y: 65},
					{X: 170, Y: 68}, {X: 172, Y: 70}, {X: 173, Y: 66},
					{X: 175, Y: 75}, {X: 176, Y: 72}, {X: 178, Y: 78},
					{X: 180, Y: 80}, {X: 181, Y: 76}, {X: 183, Y: 85},
					{X: 185, Y: 82}, {X: 187, Y: 88}, {X: 188, Y: 90},
					{X: 190, Y: 92}, {X: 193, Y: 95},
				},
			}),
		},
	}), `chart.Scatter(chart.ScatterCfg{
    BaseCfg: chart.BaseCfg{
        Title: "Height vs Weight",
    },
    Series: []series.XY{
        series.NewXY(series.XYCfg{
            Name:   "Subjects",
            Color:  gui.Hex(0x4E79A7),
            Points: []series.Point{
                {X: 155, Y: 52}, {X: 160, Y: 58},
                {X: 165, Y: 62}, {X: 170, Y: 68},
                // ... 20 points total
            },
        }),
    },
})`)
}

func demoScatterMarkers(w *gui.Window) gui.View {
	return demoWithCode(w, "scatter-markers", chart.Scatter(chart.ScatterCfg{
		BaseCfg: chart.BaseCfg{
			ID:             "scatter-markers",
			Title:          "Wind Speed vs Temperature",
			Sizing:         gui.FillFixed,
			Height:         350,
			LegendPosition: &posBottom,
		},
		Marker: chart.MarkerSquare,
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "Coastal",
				Color: gui.Hex(0x4E79A7),
				Points: []series.Point{
					{X: 5, Y: 22}, {X: 8, Y: 20}, {X: 12, Y: 18},
					{X: 15, Y: 16}, {X: 18, Y: 14}, {X: 22, Y: 12},
					{X: 25, Y: 10}, {X: 28, Y: 8},
				},
			}),
			series.NewXY(series.XYCfg{
				Name:  "Inland",
				Color: gui.Hex(0xE15759),
				Points: []series.Point{
					{X: 3, Y: 28}, {X: 6, Y: 25}, {X: 10, Y: 22},
					{X: 14, Y: 19}, {X: 17, Y: 16}, {X: 20, Y: 13},
					{X: 24, Y: 10}, {X: 27, Y: 7},
				},
			}),
		},
	}), `chart.Scatter(chart.ScatterCfg{
    BaseCfg: chart.BaseCfg{
        Title: "Wind Speed vs Temperature",
    },
    Marker: chart.MarkerSquare,
    Series: []series.XY{
        series.NewXY(series.XYCfg{
            Name:   "Coastal",
            Color:  gui.Hex(0x4E79A7),
            Points: []series.Point{
                {X: 5, Y: 22}, {X: 8, Y: 20}, {X: 12, Y: 18},
                {X: 15, Y: 16}, {X: 18, Y: 14}, {X: 22, Y: 12},
                {X: 25, Y: 10}, {X: 28, Y: 8},
            },
        }),
        series.NewXY(series.XYCfg{
            Name:   "Inland",
            Color:  gui.Hex(0xE15759),
            Points: []series.Point{ ... },
        }),
    },
})`)
}
