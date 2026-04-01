// Package theme provides theming for charts.
package theme

import "github.com/mike-ward/go-gui/gui"

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

	// Series palette
	Palette []gui.Color

	// Spacing
	PaddingTop    float32
	PaddingRight  float32
	PaddingBottom float32
	PaddingLeft   float32
}

// Default returns the default chart theme that inherits from
// the current gui theme.
func Default() *Theme {
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
		PaddingTop:    40,
		PaddingRight:  20,
		PaddingBottom: 40,
		PaddingLeft:   60,
	}
}
