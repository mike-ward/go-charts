# Go-Charts

[![Go 1.26+](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![License: PolyForm NC 1.0](https://img.shields.io/badge/License-PolyForm%20NC%201.0-blue)](LICENSE)

Professional charting library for Go, built on
[go-gui](https://github.com/mike-ward/go-gui). Immediate-mode rendering
via `DrawCanvas` — no virtual DOM, no diffing, just fast composable charts.

## Status

Early development. Package structure and APIs are scaffolded. Rendering
implementations are in progress.

## Quick Start

```go
package main

import (
    "github.com/mike-ward/go-charts/chart"
    "github.com/mike-ward/go-charts/series"
    "github.com/mike-ward/go-gui/gui"
)

func main() {
    w := gui.NewWindow(gui.WindowCfg{
        Title: "Line Chart",
        Width: 800, Height: 600,
        View: func(w *gui.Window) gui.View {
            return chart.Line(chart.LineCfg{
                ID: "demo",
                Series: []series.XY{
                    series.NewXY(series.XYCfg{
                        Name:  "Sales",
                        Color: gui.Hex(0x4E79A7),
                        Points: []series.Point{
                            {X: 1, Y: 10},
                            {X: 2, Y: 25},
                            {X: 3, Y: 18},
                            {X: 4, Y: 32},
                        },
                    }),
                },
            })
        },
    })
    w.Run()
}
```

## Packages

| Package  | Description                                    |
|----------|------------------------------------------------|
| `chart`  | Chart widgets: Line, Bar, Area, Scatter, Pie   |
| `axis`   | Axis types: Linear, Log, Time, Category        |
| `series` | Data series: XY, Category, OHLC                |
| `scale`  | Data-to-pixel mapping: Linear, Log             |
| `render` | DrawContext adapter for chart rendering         |
| `theme`  | Theming and color palettes                     |

## Chart Types

| Type    | Function         | Status |
|---------|------------------|--------|
| Line    | `chart.Line()`   | Stub   |
| Bar     | `chart.Bar()`    | Stub   |
| Area    | `chart.Area()`   | Stub   |
| Scatter | `chart.Scatter()`| Stub   |
| Pie     | `chart.Pie()`    | Stub   |

See [doc/ROADMAP.md](doc/ROADMAP.md) for planned chart types and features.

## Requirements

- Go 1.26+
- [go-gui](https://github.com/mike-ward/go-gui) (sibling directory)
- SDL2 (for running examples)

## Build

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

## License

[PolyForm Noncommercial License 1.0.0](LICENSE)
