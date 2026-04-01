package chart

import (
	"errors"
	"strings"

	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

// BaseCfg contains fields common to all chart configuration
// structs. Embed it to inherit ID, Title, sizing, theme,
// interaction callbacks, and version tracking.
type BaseCfg struct {
	ID     string
	Title  string
	Sizing gui.Sizing
	Width  float32
	Height float32

	Theme   *theme.Theme
	OnClick func(*gui.Layout, *gui.Event, *gui.Window)
	OnHover func(*gui.Layout, *gui.Event, *gui.Window)

	// XTickRotation rotates X-axis tick labels (radians).
	// 0 = horizontal.
	XTickRotation float32

	// LegendPosition overrides the theme legend position for
	// this chart. nil = use theme default.
	LegendPosition *theme.LegendPosition

	Version uint64
}

// applyDefaults sets sensible zero-value defaults.
func (b *BaseCfg) applyDefaults() {
	if b.Sizing == (gui.Sizing{}) {
		b.Sizing = gui.FillFill
	}
	if b.Theme == nil {
		b.Theme = theme.Default()
	}
}

// Validate checks BaseCfg for invalid settings. Returns nil
// when valid.
func (b *BaseCfg) Validate() error {
	var errs []string
	if b.Width < 0 {
		errs = append(errs, "negative Width")
	}
	if b.Height < 0 {
		errs = append(errs, "negative Height")
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New("chart: " + strings.Join(errs, "; "))
}

// buildError joins error strings with a prefix, returning nil
// when empty.
func buildError(prefix string, errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.New(prefix + ": " + strings.Join(errs, "; "))
}
