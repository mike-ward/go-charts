package main

import (
	"github.com/mike-ward/go-charts/chart"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-gui/gui"
)

func demoBarBasic(_ *gui.Window) gui.View {
	return chart.Bar(chart.BarCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "bar-basic",
			Title:  "Sales by Region",
			Sizing: gui.FillFixed,
			Height: 350,
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
			series.NewCategory(series.CategoryCfg{
				Name:  "Q2",
				Color: gui.Hex(0xF28E2B),
				Values: []series.CategoryValue{
					{Label: "North", Value: 52},
					{Label: "South", Value: 38},
					{Label: "East", Value: 49},
					{Label: "West", Value: 55},
				},
			}),
		},
	})
}

func demoBarSingle(_ *gui.Window) gui.View {
	return chart.Bar(chart.BarCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "bar-single",
			Title:  "Monthly Rainfall (mm)",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name:  "2025",
				Color: gui.Hex(0x76B7B2),
				Values: []series.CategoryValue{
					{Label: "Jan", Value: 78},
					{Label: "Feb", Value: 63},
					{Label: "Mar", Value: 85},
					{Label: "Apr", Value: 92},
					{Label: "May", Value: 110},
					{Label: "Jun", Value: 72},
				},
			}),
		},
	})
}

func demoBarWide(_ *gui.Window) gui.View {
	return chart.Bar(chart.BarCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "bar-wide",
			Title:  "Department Headcount",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		BarWidth: 40,
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name:  "Employees",
				Color: gui.Hex(0xB07AA1),
				Values: []series.CategoryValue{
					{Label: "Eng", Value: 120},
					{Label: "Sales", Value: 85},
					{Label: "Mktg", Value: 42},
					{Label: "Ops", Value: 67},
					{Label: "HR", Value: 28},
				},
			}),
		},
	})
}

func demoBarRounded(_ *gui.Window) gui.View {
	return chart.Bar(chart.BarCfg{
		BaseCfg: chart.BaseCfg{
			ID:     "bar-rounded",
			Title:  "Product Revenue",
			Sizing: gui.FillFixed,
			Height: 350,
		},
		Radius: 4,
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name:  "Online",
				Color: gui.Hex(0x59A14F),
				Values: []series.CategoryValue{
					{Label: "Widgets", Value: 340},
					{Label: "Gadgets", Value: 280},
					{Label: "Gizmos", Value: 195},
					{Label: "Doohickeys", Value: 150},
				},
			}),
			series.NewCategory(series.CategoryCfg{
				Name:  "Retail",
				Color: gui.Hex(0xEDC948),
				Values: []series.CategoryValue{
					{Label: "Widgets", Value: 210},
					{Label: "Gadgets", Value: 320},
					{Label: "Gizmos", Value: 175},
					{Label: "Doohickeys", Value: 230},
				},
			}),
		},
	})
}
