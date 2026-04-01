package scale

// Linear is a linear data-to-pixel scale.
type Linear struct {
	min, max float64
}

// NewLinear creates a linear scale.
func NewLinear(min, max float64) *Linear {
	return &Linear{min: min, max: max}
}

// SetDomain implements Scale.
func (s *Linear) SetDomain(min, max float64) {
	s.min = min
	s.max = max
}

// Domain implements Scale.
func (s *Linear) Domain() (float64, float64) {
	return s.min, s.max
}

// Map implements Scale.
func (s *Linear) Map(value float64, pixelMin, pixelMax float32) float32 {
	if s.max == s.min {
		return pixelMin
	}
	t := (value - s.min) / (s.max - s.min)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Invert implements Scale.
func (s *Linear) Invert(pixel, pixelMin, pixelMax float32) float64 {
	if pixelMax == pixelMin {
		return s.min
	}
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	return s.min + t*(s.max-s.min)
}
