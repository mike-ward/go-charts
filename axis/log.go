package axis

import "math"

// Log is a logarithmic axis.
type Log struct {
	title    string
	min, max float64
	base     float64
}

// LogCfg configures a logarithmic axis.
type LogCfg struct {
	Title string
	Min   float64
	Max   float64
	Base  float64
}

// NewLog creates a logarithmic axis.
func NewLog(cfg LogCfg) *Log {
	base := cfg.Base
	if base <= 0 {
		base = 10
	}
	return &Log{
		title: cfg.Title,
		min:   cfg.Min,
		max:   cfg.Max,
		base:  base,
	}
}

// Label implements Axis.
func (a *Log) Label() string { return a.title }

// Ticks implements Axis.
func (a *Log) Ticks(pixelMin, pixelMax float32) []Tick {
	// TODO: implement logarithmic tick generation
	return nil
}

// Transform implements Axis.
func (a *Log) Transform(value float64, pixelMin, pixelMax float32) float32 {
	if value <= 0 || a.min <= 0 || a.max <= a.min {
		return pixelMin
	}
	logBase := math.Log(a.base)
	logMin := math.Log(a.min) / logBase
	logMax := math.Log(a.max) / logBase
	logVal := math.Log(value) / logBase
	t := (logVal - logMin) / (logMax - logMin)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Inverse implements Axis.
func (a *Log) Inverse(pixel, pixelMin, pixelMax float32) float64 {
	if pixelMax == pixelMin || a.min <= 0 {
		return a.min
	}
	logBase := math.Log(a.base)
	logMin := math.Log(a.min) / logBase
	logMax := math.Log(a.max) / logBase
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	logVal := logMin + t*(logMax-logMin)
	return math.Pow(a.base, logVal)
}
