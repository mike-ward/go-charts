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

// timeToSeconds converts a time.Time to seconds as float64,
// avoiding int64 overflow from UnixNano.
func timeToSeconds(t time.Time) float64 {
	return float64(t.Unix()) + float64(t.Nanosecond())/1e9
}

// Transform implements Axis. Value is expected as seconds
// (float64), matching the output of timeToSeconds.
func (a *Time) Transform(value float64, pixelMin, pixelMax float32) float32 {
	minSec := timeToSeconds(a.min)
	rangeSec := timeToSeconds(a.max) - minSec
	if rangeSec == 0 {
		return pixelMin
	}
	t := (value - minSec) / rangeSec
	return pixelMin + float32(t)*(pixelMax-pixelMin)
}

// Inverse implements Axis. Returns seconds as float64.
func (a *Time) Inverse(pixel, pixelMin, pixelMax float32) float64 {
	minSec := timeToSeconds(a.min)
	rangeSec := timeToSeconds(a.max) - minSec
	if pixelMax == pixelMin || rangeSec == 0 {
		return minSec
	}
	t := float64(pixel-pixelMin) / float64(pixelMax-pixelMin)
	return minSec + t*rangeSec
}
