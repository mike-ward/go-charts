# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Commands

```
go test ./...                        # run all tests
go test ./chart/... -run TestFoo     # run single test
go vet ./...                         # static analysis
golangci-lint run ./...              # full lint
go build ./...                       # build all packages
```

## Architecture

Professional charting library built on go-gui. Charts render via
`gui.DrawCanvas` using immediate-mode `OnDraw(*DrawContext)` callbacks.
Retained tessellation cache skips re-render when `Version` unchanged.

```
chart.Line(LineCfg{...}) → gui.View
  → GenerateLayout() wraps gui.DrawCanvas
  → OnDraw(*DrawContext) renders axes, series, legend
  → DrawContext primitives: Line, Polyline, FilledRect, FilledArc, ...
```

### Packages

- `chart/`  — chart widget views (Line, Bar, Area, Scatter, Pie);
  each implements `gui.View` via `DrawCanvas`
- `axis/`   — Axis interface + Linear, Log, Time, Category axes;
  tick generation (nice-number algorithm)
- `series/` — data series types: XY, Category, OHLC
- `scale/`  — Scale interface: data-to-pixel mapping (Linear, Log)
- `render/` — DrawContext adapter with chart-specific helpers
- `theme/`  — chart Theme type inheriting `gui.CurrentTheme()`;
  palettes (Tableau 10, Pastel, Vivid)

### Key Types

- `chart.*Cfg` — config structs for each chart type (zero-initializable)
- `axis.Axis` — interface: `Label()`, `Ticks()`, `Transform()`, `Inverse()`
- `series.Series` — interface: `Name()`, `Len()`, `Color()`
- `series.XY` — `[]Point` with `Bounds()` method
- `scale.Scale` — interface: `Map()`, `Invert()`, `SetDomain()`, `Domain()`
- `theme.Theme` — colors, text styles, palette, padding
- `render.Context` — wraps `*gui.DrawContext`

### Dependencies

- `github.com/mike-ward/go-gui` — GUI framework (local replace `../go-gui`)
- No other external dependencies

### Pattern Notes

- All chart types follow go-gui `*Cfg` struct convention
- Charts implement `gui.View` (`Content() []View`, `GenerateLayout(*Window) Layout`)
- Event callbacks: `func(*gui.Layout, *gui.Event, *gui.Window)`
- Default sizing: `gui.FillFill`
- `gui.Hex(0xRRGGBB)` — 3-byte RGB, alpha defaults to 255
- `gui.RGBA(r, g, b, a)` — explicit alpha

## Coding Conventions

- **No variable shadowing.** Use `=` for existing variables, not `:=`.
- **Clean lint and format.** `golangci-lint run ./...` and `gofmt` must
  pass with zero issues.
- Comments wrap at 90 columns when practical.
- Performance improvements should favor reducing heap allocations.
