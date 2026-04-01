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
