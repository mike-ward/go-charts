// Package scale provides data-to-pixel mapping transformations.
package scale

// Scale maps data values to pixel coordinates and back.
type Scale interface {
	// Transform converts a data value to a pixel position.
	Transform(value float64, pixelMin, pixelMax float32) float32

	// Invert converts a pixel position back to a data value.
	Invert(pixel, pixelMin, pixelMax float32) float64

	// SetDomain sets the data range.
	SetDomain(min, max float64)

	// Domain returns the current data range.
	Domain() (min, max float64)
}
