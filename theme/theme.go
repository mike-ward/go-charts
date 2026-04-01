// Package theme provides theming for charts.
package theme

import (
	"github.com/mike-ward/go-gui/gui"
)

// Default padding values for chart themes.
const (
	DefaultPaddingTop    float32 = 40
	DefaultPaddingRight  float32 = 40
	DefaultPaddingBottom float32 = 60
	DefaultPaddingLeft   float32 = 60
)

// TickMarkStyle controls axis tick mark appearance. Zero values
// fall back to the axis line style (AxisColor / AxisWidth / 5px).
type TickMarkStyle struct {
	Length float32   // 0 → default (5)
	Color  gui.Color // zero → AxisColor
	Width  float32   // 0 → AxisWidth
}

// LegendPosition selects where the legend box is placed within
// the plot area.
type LegendPosition uint8

// Legend position constants.
const (
	LegendTopRight LegendPosition = iota
	LegendTopLeft
	LegendBottomRight
	LegendBottomLeft
)

// LegendStyle controls legend appearance. Zero values preserve
// the original defaults.
type LegendStyle struct {
	Position   LegendPosition
	TextStyle  gui.TextStyle // zero → Theme.LabelStyle
	Background gui.Color     // zero → RGBA(0,0,0,120)
	SwatchSize float32       // zero → 12
	Padding    float32       // zero → 6
	ItemGap    float32       // zero → 4
	RowGap     float32       // zero → 2
}

// Theme defines the visual style for charts.
type Theme struct {
	// Background
	Background gui.Color

	// Text
	TitleStyle gui.TextStyle
	LabelStyle gui.TextStyle
	TickStyle  gui.TextStyle

	// Axes
	AxisColor gui.Color
	AxisWidth float32
	GridColor gui.Color
	GridWidth float32

	// Tick marks
	TickMark TickMarkStyle

	// Legend
	Legend LegendStyle

	// Series palette
	Palette []gui.Color

	// Spacing
	PaddingTop    float32
	PaddingRight  float32
	PaddingBottom float32
	PaddingLeft   float32
}

var globalDefault *Theme

// SetDefault sets the global default chart theme. Passing nil
// reverts to the auto-generated theme from gui.CurrentTheme().
// Call once at startup before creating charts.
func SetDefault(t *Theme) { globalDefault = t }

// Default returns the global default theme if set via
// SetDefault, otherwise creates one from gui.CurrentTheme().
// When no global default is set, allocates on each call;
// callers rendering at interactive frame rates should cache
// the returned value or call SetDefault once at startup.
func Default() *Theme {
	if globalDefault != nil {
		return globalDefault
	}
	t := gui.CurrentTheme()
	return &Theme{
		Background:    t.ColorBackground,
		TitleStyle:    t.B1,
		LabelStyle:    t.TextStyleDef,
		TickStyle:     t.TextStyleDef,
		AxisColor:     t.ColorBorder,
		AxisWidth:     1,
		GridColor:     gui.RGBA(128, 128, 128, 40),
		GridWidth:     0.5,
		Palette:       DefaultPalette(),
		PaddingTop:    DefaultPaddingTop,
		PaddingRight:  DefaultPaddingRight,
		PaddingBottom: DefaultPaddingBottom,
		PaddingLeft:   DefaultPaddingLeft,
	}
}
