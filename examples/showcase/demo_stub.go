package main

import (
	"github.com/mike-ward/go-charts/chart"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-gui/gui"
)

func demoStubPie(w *gui.Window) gui.View {
	return demoWithCode(w, "pie-basic", chart.Pie(chart.PieCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "pie-basic",
			Title:  "Browser Market Share",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		ShowLabels:  true,
		ShowPercent: true,
		Slices: []chart.PieSlice{
			{Label: "Chrome", Value: 65, Color: gui.Hex(0x4E79A7)},
			{Label: "Safari", Value: 18, Color: gui.Hex(0xF28E2B)},
			{Label: "Firefox", Value: 8, Color: gui.Hex(0xE15759)},
			{Label: "Edge", Value: 5, Color: gui.Hex(0x76B7B2)},
			{Label: "Other", Value: 4, Color: gui.Hex(0x59A14F)},
		},
	}), `chart.Pie(chart.PieCfg{
    BaseCfg: chart.BaseCfg{
        Title: "Browser Market Share",
    },
    ShowLabels:  true,
    ShowPercent: true,
    Slices: []chart.PieSlice{
        {Label: "Chrome", Value: 65, Color: gui.Hex(0x4E79A7)},
        {Label: "Safari", Value: 18, Color: gui.Hex(0xF28E2B)},
        {Label: "Firefox", Value: 8, Color: gui.Hex(0xE15759)},
        {Label: "Edge", Value: 5, Color: gui.Hex(0x76B7B2)},
        {Label: "Other", Value: 4, Color: gui.Hex(0x59A14F)},
    },
})`)
}

func demoStubDonut(w *gui.Window) gui.View {
	return demoWithCode(w, "pie-donut", chart.Pie(chart.PieCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "pie-donut",
			Title:  "Budget Allocation",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		InnerRadius: 60,
		ShowLabels:  true,
		Slices: []chart.PieSlice{
			{Label: "R&D", Value: 35, Color: gui.Hex(0x4E79A7)},
			{Label: "Marketing", Value: 25, Color: gui.Hex(0xF28E2B)},
			{Label: "Operations", Value: 22, Color: gui.Hex(0xE15759)},
			{Label: "Admin", Value: 18, Color: gui.Hex(0x76B7B2)},
		},
	}), `chart.Pie(chart.PieCfg{
    BaseCfg: chart.BaseCfg{
        Title: "Budget Allocation",
    },
    InnerRadius: 60,
    ShowLabels:  true,
    Slices: []chart.PieSlice{
        {Label: "R&D", Value: 35, Color: gui.Hex(0x4E79A7)},
        {Label: "Marketing", Value: 25, Color: gui.Hex(0xF28E2B)},
        {Label: "Operations", Value: 22, Color: gui.Hex(0xE15759)},
        {Label: "Admin", Value: 18, Color: gui.Hex(0x76B7B2)},
    },
})`)
}

func demoStubArea(w *gui.Window) gui.View {
	return demoWithCode(w, "area-basic", chart.Area(chart.AreaCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "area-basic",
			Title:  "User Signups",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "Organic",
				Color: gui.Hex(0x4E79A7),
				Points: []series.Point{
					{X: 1, Y: 120}, {X: 2, Y: 180}, {X: 3, Y: 250},
					{X: 4, Y: 310}, {X: 5, Y: 400}, {X: 6, Y: 480},
				},
			}),
			series.NewXY(series.XYCfg{
				Name:  "Paid",
				Color: gui.Hex(0xF28E2B),
				Points: []series.Point{
					{X: 1, Y: 80}, {X: 2, Y: 110}, {X: 3, Y: 160},
					{X: 4, Y: 200}, {X: 5, Y: 260}, {X: 6, Y: 300},
				},
			}),
		},
	}), `chart.Area(chart.AreaCfg{
    BaseCfg: chart.BaseCfg{
        Title: "User Signups",
    },
    Series: []series.XY{
        series.NewXY(series.XYCfg{
            Name:   "Organic",
            Color:  gui.Hex(0x4E79A7),
            Points: []series.Point{
                {X: 1, Y: 120}, {X: 2, Y: 180}, {X: 3, Y: 250},
                {X: 4, Y: 310}, {X: 5, Y: 400}, {X: 6, Y: 480},
            },
        }),
        series.NewXY(series.XYCfg{
            Name:   "Paid",
            Color:  gui.Hex(0xF28E2B),
            Points: []series.Point{ ... },
        }),
    },
})`)
}

func demoStubAreaStacked(w *gui.Window) gui.View {
	return demoWithCode(w, "area-stacked", chart.Area(chart.AreaCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "area-stacked",
			Title:  "Revenue by Product",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		Stacked: true,
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "Widgets",
				Color: gui.Hex(0x4E79A7),
				Points: []series.Point{
					{X: 1, Y: 50}, {X: 2, Y: 55}, {X: 3, Y: 62},
					{X: 4, Y: 58}, {X: 5, Y: 70}, {X: 6, Y: 75},
				},
			}),
			series.NewXY(series.XYCfg{
				Name:  "Gadgets",
				Color: gui.Hex(0xF28E2B),
				Points: []series.Point{
					{X: 1, Y: 30}, {X: 2, Y: 38}, {X: 3, Y: 42},
					{X: 4, Y: 45}, {X: 5, Y: 50}, {X: 6, Y: 55},
				},
			}),
			series.NewXY(series.XYCfg{
				Name:  "Gizmos",
				Color: gui.Hex(0xE15759),
				Points: []series.Point{
					{X: 1, Y: 20}, {X: 2, Y: 22}, {X: 3, Y: 25},
					{X: 4, Y: 28}, {X: 5, Y: 32}, {X: 6, Y: 35},
				},
			}),
		},
	}), `chart.Area(chart.AreaCfg{
    BaseCfg: chart.BaseCfg{
        Title: "Revenue by Product",
    },
    Stacked: true,
    Series: []series.XY{
        series.NewXY(series.XYCfg{
            Name:   "Widgets",
            Color:  gui.Hex(0x4E79A7),
            Points: []series.Point{
                {X: 1, Y: 50}, {X: 2, Y: 55}, ...
            },
        }),
        series.NewXY(series.XYCfg{
            Name:   "Gadgets",
            Color:  gui.Hex(0xF28E2B),
            Points: []series.Point{ ... },
        }),
        series.NewXY(series.XYCfg{
            Name:   "Gizmos",
            Color:  gui.Hex(0xE15759),
            Points: []series.Point{ ... },
        }),
    },
})`)
}

func demoStubScatter(w *gui.Window) gui.View {
	return demoWithCode(w, "scatter-basic", chart.Scatter(chart.ScatterCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "scatter-basic",
			Title:  "Height vs Weight",
			Sizing: gui.FillFixed,
			Height: 350,
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

func demoStubScatterMarkers(w *gui.Window) gui.View {
	return demoWithCode(w, "scatter-markers", chart.Scatter(chart.ScatterCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "scatter-markers",
			Title:  "Wind Speed vs Temperature",
			Sizing: gui.FillFixed,
			Height: 350,
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
