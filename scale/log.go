package scale

import (
	"log/slog"
	"math"
)

// Log is a logarithmic data-to-pixel scale.
type Log struct {
	min, max float64
	base     float64
}

// NewLog creates a logarithmic scale. Base defaults to 10 when
// <= 0 or == 1 (log(1) == 0 causes division by zero).
func NewLog(min, max, base float64) *Log {
	if base <= 0 || base == 1 {
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

// Transform implements Scale. Non-finite values and non-positive
// domain/value return pixelMin.
func (s *Log) Transform(value float64, pixelMin, pixelMax float32) float32 {
	if value <= 0 || s.min <= 0 || s.max <= s.min ||
		!finiteF64(value) {
		if value <= 0 {
			slog.Debug("log scale: non-positive value",
				"value", value)
		}
		if s.min <= 0 {
			slog.Debug("log scale: non-positive domain min",
				"min", s.min)
		}
		return pixelMin
	}
	logBase := math.Log(s.base)
	logMin := math.Log(s.min) / logBase
	logMax := math.Log(s.max) / logBase
	if logMax == logMin {
		return pixelMin
	}
	logVal := math.Log(value) / logBase
	t := (logVal - logMin) / (logMax - logMin)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Invert implements Scale.
func (s *Log) Invert(pixel, pixelMin, pixelMax float32) float64 {
	if pixelMax == pixelMin || s.min <= 0 {
		if s.min <= 0 {
			slog.Debug("log scale: non-positive domain min in Invert",
				"min", s.min)
		}
		return s.min
	}
	logBase := math.Log(s.base)
	logMin := math.Log(s.min) / logBase
	logMax := math.Log(s.max) / logBase
	if logMax == logMin {
		slog.Debug("log scale: degenerate domain in Invert",
			"min", s.min, "max", s.max)
		return s.min
	}
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	logVal := logMin + t*(logMax-logMin)
	return math.Pow(s.base, logVal)
}
