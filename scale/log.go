package scale

import "math"

// Log is a logarithmic data-to-pixel scale.
type Log struct {
	min, max float64
	base     float64
}

// NewLog creates a logarithmic scale.
func NewLog(min, max, base float64) *Log {
	if base <= 0 {
		base = 10
	}
	return &Log{min: min, max: max, base: base}
}

// SetDomain implements Scale.
func (s *Log) SetDomain(min, max float64) {
	s.min = min
	s.max = max
}

// Domain implements Scale.
func (s *Log) Domain() (float64, float64) {
	return s.min, s.max
}

// Map implements Scale.
func (s *Log) Map(value float64, pixelMin, pixelMax float32) float32 {
	if value <= 0 || s.min <= 0 || s.max <= s.min {
		return pixelMin
	}
	logBase := math.Log(s.base)
	logMin := math.Log(s.min) / logBase
	logMax := math.Log(s.max) / logBase
	logVal := math.Log(value) / logBase
	t := (logVal - logMin) / (logMax - logMin)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Invert implements Scale.
func (s *Log) Invert(pixel, pixelMin, pixelMax float32) float64 {
	if pixelMax == pixelMin || s.min <= 0 {
		return s.min
	}
	logBase := math.Log(s.base)
	logMin := math.Log(s.min) / logBase
	logMax := math.Log(s.max) / logBase
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	logVal := logMin + t*(logMax-logMin)
	return math.Pow(s.base, logVal)
}
