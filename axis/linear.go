package axis

import "fmt"

// Linear is a linear numeric axis with auto-tick generation.
type Linear struct {
	title     string
	min, max  float64
	autoRange bool
}

// LinearCfg configures a linear axis.
type LinearCfg struct {
	Title     string
	Min       float64
	Max       float64
	AutoRange bool
}

// NewLinear creates a linear axis.
func NewLinear(cfg LinearCfg) *Linear {
	return &Linear{
		title:     cfg.Title,
		min:       cfg.Min,
		max:       cfg.Max,
		autoRange: cfg.AutoRange,
	}
}

// SetRange updates the axis data range.
func (a *Linear) SetRange(min, max float64) {
	a.min = min
	a.max = max
}

// Label implements Axis.
func (a *Linear) Label() string { return a.title }

// Ticks implements Axis.
func (a *Linear) Ticks(pixelMin, pixelMax float32) []Tick {
	values := GenerateNiceTicks(a.min, a.max, 8)
	ticks := make([]Tick, 0, len(values))
	for _, v := range values {
		ticks = append(ticks, Tick{
			Value:    v,
			Label:    formatTickValue(v),
			Position: a.Transform(v, pixelMin, pixelMax),
		})
	}
	return ticks
}

func formatTickValue(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%g", v)
}

// Transform implements Axis.
func (a *Linear) Transform(value float64, pixelMin, pixelMax float32) float32 {
	if a.max == a.min {
		return pixelMin
	}
	t := (value - a.min) / (a.max - a.min)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Inverse implements Axis.
func (a *Linear) Inverse(pixel, pixelMin, pixelMax float32) float64 {
	if pixelMax == pixelMin {
		return a.min
	}
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	return a.min + t*(a.max-a.min)
}
