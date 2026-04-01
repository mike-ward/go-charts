package theme

import "github.com/mike-ward/go-gui/gui"

// WithAlpha returns a copy of the color with the given alpha
// (0.0 = fully transparent, 1.0 = fully opaque).
func WithAlpha(c gui.Color, alpha float64) gui.Color {
	a := uint8(max(0, min(255, alpha*255)))
	return gui.RGBA(c.R, c.G, c.B, a)
}

// Lighten returns a lighter version of the color. Amount is
// 0.0 (unchanged) to 1.0 (white).
func Lighten(c gui.Color, amount float64) gui.Color {
	f := max(0, min(1, amount))
	r := c.R + uint8(float64(255-c.R)*f)
	g := c.G + uint8(float64(255-c.G)*f)
	b := c.B + uint8(float64(255-c.B)*f)
	return gui.RGBA(r, g, b, c.A)
}

// Darken returns a darker version of the color. Amount is
// 0.0 (unchanged) to 1.0 (black).
func Darken(c gui.Color, amount float64) gui.Color {
	f := 1 - max(0, min(1, amount))
	r := uint8(float64(c.R) * f)
	g := uint8(float64(c.G) * f)
	b := uint8(float64(c.B) * f)
	return gui.RGBA(r, g, b, c.A)
}
