package axis

import "time"

// Time is a time-based axis.
type Time struct {
	title    string
	min, max time.Time
	format   string
}

// TimeCfg configures a time axis.
type TimeCfg struct {
	Title  string
	Min    time.Time
	Max    time.Time
	Format string // Go time format string
}

// NewTime creates a time axis.
func NewTime(cfg TimeCfg) *Time {
	format := cfg.Format
	if format == "" {
		format = "2006-01-02"
	}
	return &Time{
		title:  cfg.Title,
		min:    cfg.Min,
		max:    cfg.Max,
		format: format,
	}
}

// Label implements Axis.
func (a *Time) Label() string { return a.title }

// Ticks implements Axis.
func (a *Time) Ticks(pixelMin, pixelMax float32) []Tick {
	// TODO: implement time-aware tick generation
	return nil
}

// Transform implements Axis.
func (a *Time) Transform(value float64, pixelMin, pixelMax float32) float32 {
	minNano := float64(a.min.UnixNano())
	maxNano := float64(a.max.UnixNano())
	if maxNano == minNano {
		return pixelMin
	}
	t := (value - minNano) / (maxNano - minNano)
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Inverse implements Axis.
func (a *Time) Inverse(pixel, pixelMin, pixelMax float32) float64 {
	minNano := float64(a.min.UnixNano())
	maxNano := float64(a.max.UnixNano())
	if pixelMax == pixelMin {
		return minNano
	}
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	return minNano + t*(maxNano-minNano)
}
