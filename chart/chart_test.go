package chart

import (
	"testing"

	"github.com/mike-ward/go-gui/gui"
)

func TestChartCreatesView(t *testing.T) {
	v := Chart(Cfg{ID: "test-chart"})
	if v == nil {
		t.Fatal("Chart() returned nil")
	}
}

func TestChartDefaultSizing(t *testing.T) {
	cv := Chart(Cfg{ID: "test"}).(*chartView)
	if cv.cfg.Sizing != gui.FillFill {
		t.Errorf("expected GrowGrow, got %v", cv.cfg.Sizing)
	}
}
