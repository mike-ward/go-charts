package axis

import "github.com/mike-ward/go-charts/scale"

// Log is a logarithmic axis.
type Log struct {
	title      string
	sc         *scale.Log
	tickFormat TickFormat
}

// LogCfg configures a logarithmic axis.
type LogCfg struct {
	Title      string
	Min        float64
	Max        float64
	Base       float64
	TickFormat TickFormat
}

// NewLog creates a logarithmic axis. Base defaults to 10 via
// scale.NewLog when <= 0.
func NewLog(cfg LogCfg) *Log {
	return &Log{
		title:      cfg.Title,
		sc:         scale.NewLog(cfg.Min, cfg.Max, cfg.Base),
		tickFormat: cfg.TickFormat,
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
	return a.sc.Map(value, pixelMin, pixelMax)
}

// Inverse implements Axis.
func (a *Log) Inverse(pixel, pixelMin, pixelMax float32) float64 {
	return a.sc.Invert(pixel, pixelMin, pixelMax)
}
