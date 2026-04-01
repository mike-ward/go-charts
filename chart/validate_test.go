package chart

import (
	"testing"

	"github.com/mike-ward/go-charts/series"
)

func TestLineCfgValidateOK(t *testing.T) {
	cfg := LineCfg{
		Series: []series.XY{series.XYFromYValues("s", []float64{1})},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLineCfgValidateEmpty(t *testing.T) {
	cfg := LineCfg{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty series")
	}
}

func TestLineCfgValidateNegativeWidth(t *testing.T) {
	cfg := LineCfg{
		BaseCfg: BaseCfg{Width: -1},
		Series:  []series.XY{series.XYFromYValues("s", []float64{1})},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative Width")
	}
}

func TestBarCfgValidateEmpty(t *testing.T) {
	cfg := BarCfg{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty series")
	}
}

func TestBarCfgValidateNegativeBarWidth(t *testing.T) {
	cfg := BarCfg{
		Series:   []series.Category{series.CategoryFromMap("s", map[string]float64{"a": 1})},
		BarWidth: -1,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative BarWidth")
	}
}

func TestPieCfgValidateOK(t *testing.T) {
	cfg := PieCfg{
		Slices: []PieSlice{{Label: "a", Value: 1}},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPieCfgValidateEmpty(t *testing.T) {
	cfg := PieCfg{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty slices")
	}
}
