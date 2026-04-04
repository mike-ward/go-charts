package chart

import (
	"math"
	"testing"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

func TestZoomAroundCursorCenter(t *testing.T) {
	t.Parallel()
	// Cursor at center of [0, 100], zoom in 2x.
	lo, hi := zoomAroundCursor(50, 0, 100, 2)
	if math.Abs(lo-25) > 1e-9 || math.Abs(hi-75) > 1e-9 {
		t.Errorf("got [%g, %g], want [25, 75]", lo, hi)
	}
}

func TestZoomAroundCursorEdge(t *testing.T) {
	t.Parallel()
	// Cursor at left edge.
	lo, hi := zoomAroundCursor(0, 0, 100, 2)
	if math.Abs(lo-0) > 1e-9 || math.Abs(hi-50) > 1e-9 {
		t.Errorf("got [%g, %g], want [0, 50]", lo, hi)
	}
	// Cursor at right edge.
	lo, hi = zoomAroundCursor(100, 0, 100, 2)
	if math.Abs(lo-50) > 1e-9 || math.Abs(hi-100) > 1e-9 {
		t.Errorf("got [%g, %g], want [50, 100]", lo, hi)
	}
}

func TestZoomAroundCursorZoomOut(t *testing.T) {
	t.Parallel()
	// Factor < 1 → zoom out (double range).
	lo, hi := zoomAroundCursor(50, 25, 75, 0.5)
	if math.Abs(lo-0) > 1e-9 || math.Abs(hi-100) > 1e-9 {
		t.Errorf("got [%g, %g], want [0, 100]", lo, hi)
	}
}

func TestZoomAroundCursorZeroDomain(t *testing.T) {
	t.Parallel()
	lo, hi := zoomAroundCursor(5, 5, 5, 2)
	if lo != 5 || hi != 5 {
		t.Errorf("got [%g, %g], want [5, 5]", lo, hi)
	}
}

func TestZoomAroundCursorZeroFactor(t *testing.T) {
	t.Parallel()
	lo, hi := zoomAroundCursor(50, 0, 100, 0)
	if lo != 0 || hi != 100 {
		t.Errorf("got [%g, %g], want [0, 100]", lo, hi)
	}
}

func TestPanDomain(t *testing.T) {
	t.Parallel()
	// 10px drag on 100px span → 10% shift.
	lo, hi := panDomain(0, 100, 10, 100)
	// float32→float64 conversion introduces ~1e-7 error.
	if math.Abs(lo-(-10)) > 1e-5 || math.Abs(hi-90) > 1e-5 {
		t.Errorf("got [%g, %g], want [-10, 90]", lo, hi)
	}
}

func TestPanDomainZeroSpan(t *testing.T) {
	t.Parallel()
	lo, hi := panDomain(0, 100, 10, 0)
	if lo != 0 || hi != 100 {
		t.Errorf("got [%g, %g], want [0, 100]", lo, hi)
	}
}

func TestClampZoomRangeMinSpan(t *testing.T) {
	t.Parallel()
	lo, hi := clampZoomRange(50, 50+1e-15, 0, 100)
	span := hi - lo
	// Allow small floating-point rounding (< 1% of min range).
	if span < DefaultMinZoomRange*0.99 {
		t.Errorf("span %g < min %g", span, DefaultMinZoomRange)
	}
}

func TestClampZoomRangeMaxExtent(t *testing.T) {
	t.Parallel()
	lo, hi := clampZoomRange(-10, 110, 0, 100)
	if lo < 0 || hi > 100 {
		t.Errorf("got [%g, %g], want within [0, 100]", lo, hi)
	}
}

func TestClampZoomRangeShiftRight(t *testing.T) {
	t.Parallel()
	// Panned left beyond origin → shift right.
	lo, hi := clampZoomRange(-5, 45, 0, 100)
	if math.Abs(lo-0) > 1e-9 || math.Abs(hi-50) > 1e-9 {
		t.Errorf("got [%g, %g], want [0, 50]", lo, hi)
	}
}

func TestClampZoomRangeShiftLeft(t *testing.T) {
	t.Parallel()
	// Panned right beyond max → shift left.
	lo, hi := clampZoomRange(60, 110, 0, 100)
	if math.Abs(lo-50) > 1e-9 || math.Abs(hi-100) > 1e-9 {
		t.Errorf("got [%g, %g], want [50, 100]", lo, hi)
	}
}

func TestClampZoomRangeMinSpanAtBoundary(t *testing.T) {
	t.Parallel()
	// Near origMax: min-span expansion must not exceed origMax.
	lo, hi := clampZoomRange(100-1e-15, 100, 0, 100)
	if hi > 100 {
		t.Errorf("dMax %g exceeds origMax 100", hi)
	}
	if lo < 0 {
		t.Errorf("dMin %g below origMin 0", lo)
	}
	// Near origMin: min-span expansion must not go below origMin.
	lo, hi = clampZoomRange(0, 0+1e-15, 0, 100)
	if lo < 0 {
		t.Errorf("dMin %g below origMin 0", lo)
	}
	if hi > 100 {
		t.Errorf("dMax %g exceeds origMax 100", hi)
	}
}

// --- clampRectToPlot ---

func TestClampRectToPlotInside(t *testing.T) {
	t.Parallel()
	x, y, w, h, vis := clampRectToPlot(20, 20, 60, 40, 0, 100, 0, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if x != 20 || y != 20 || w != 60 || h != 40 {
		t.Errorf("got (%g,%g,%g,%g), want (20,20,60,40)", x, y, w, h)
	}
}

func TestClampRectToPlotClipsTop(t *testing.T) {
	t.Parallel()
	// Rect extends above plot top (top=10).
	x, y, w, h, vis := clampRectToPlot(20, 5, 30, 20, 0, 100, 10, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y != 10 || h != 15 {
		t.Errorf("got y=%g h=%g, want y=10 h=15", y, h)
	}
	if x != 20 || w != 30 {
		t.Errorf("x/w changed: got (%g,%g)", x, w)
	}
}

func TestClampRectToPlotClipsBottom(t *testing.T) {
	t.Parallel()
	// Rect extends below plot bottom (bottom=80).
	_, y, _, h, vis := clampRectToPlot(20, 70, 30, 20, 0, 100, 0, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y != 70 || h != 10 {
		t.Errorf("got y=%g h=%g, want y=70 h=10", y, h)
	}
}

func TestClampRectToPlotClipsLeft(t *testing.T) {
	t.Parallel()
	x, _, w, _, vis := clampRectToPlot(-5, 20, 20, 10, 0, 100, 0, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if x != 0 || w != 15 {
		t.Errorf("got x=%g w=%g, want x=0 w=15", x, w)
	}
}

func TestClampRectToPlotClipsRight(t *testing.T) {
	t.Parallel()
	x, _, w, _, vis := clampRectToPlot(90, 20, 20, 10, 0, 100, 0, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if x != 90 || w != 10 {
		t.Errorf("got x=%g w=%g, want x=90 w=10", x, w)
	}
}

func TestClampRectToPlotFullyOutside(t *testing.T) {
	t.Parallel()
	_, _, _, _, vis := clampRectToPlot(110, 20, 20, 10, 0, 100, 0, 80)
	if vis {
		t.Error("expected not visible (right of plot)")
	}
	_, _, _, _, vis = clampRectToPlot(20, 90, 20, 10, 0, 100, 0, 80)
	if vis {
		t.Error("expected not visible (below plot)")
	}
}

// --- clampVerticalLine ---

func TestClampVerticalLineInside(t *testing.T) {
	t.Parallel()
	y0, y1, vis := clampVerticalLine(20, 60, 10, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y0 != 20 || y1 != 60 {
		t.Errorf("got [%g,%g], want [20,60]", y0, y1)
	}
}

func TestClampVerticalLineSwapped(t *testing.T) {
	t.Parallel()
	// y0 > y1 should be normalized.
	y0, y1, vis := clampVerticalLine(60, 20, 10, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y0 != 20 || y1 != 60 {
		t.Errorf("got [%g,%g], want [20,60]", y0, y1)
	}
}

func TestClampVerticalLineClipsTop(t *testing.T) {
	t.Parallel()
	y0, y1, vis := clampVerticalLine(5, 50, 10, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y0 != 10 || y1 != 50 {
		t.Errorf("got [%g,%g], want [10,50]", y0, y1)
	}
}

func TestClampVerticalLineClipsBottom(t *testing.T) {
	t.Parallel()
	y0, y1, vis := clampVerticalLine(30, 90, 10, 80)
	if !vis {
		t.Fatal("expected visible")
	}
	if y0 != 30 || y1 != 80 {
		t.Errorf("got [%g,%g], want [30,80]", y0, y1)
	}
}

func TestClampVerticalLineFullyOutside(t *testing.T) {
	t.Parallel()
	_, _, vis := clampVerticalLine(85, 95, 10, 80)
	if vis {
		t.Error("expected not visible (below plot)")
	}
	_, _, vis = clampVerticalLine(2, 8, 10, 80)
	if vis {
		t.Error("expected not visible (above plot)")
	}
}

// --- insidePlot ---

func TestInsidePlot(t *testing.T) {
	t.Parallel()
	if !insidePlot(50, 50, 0, 100, 0, 100) {
		t.Error("center should be inside")
	}
	if !insidePlot(0, 0, 0, 100, 0, 100) {
		t.Error("top-left corner should be inside (inclusive)")
	}
	if !insidePlot(100, 100, 0, 100, 0, 100) {
		t.Error("bottom-right corner should be inside (inclusive)")
	}
	if insidePlot(-1, 50, 0, 100, 0, 100) {
		t.Error("left of plot should be outside")
	}
	if insidePlot(50, 101, 0, 100, 0, 100) {
		t.Error("below plot should be outside")
	}
}

// --- selectionColors ---

func TestSelectionColorsDefaults(t *testing.T) {
	t.Parallel()
	th := &theme.Theme{}
	fill, border := selectionColors(th)
	if fill != gui.RGBA(70, 130, 220, 30) {
		t.Errorf("fill = %v, want default", fill)
	}
	if border != gui.RGBA(70, 130, 220, 180) {
		t.Errorf("border = %v, want default", border)
	}
}

func TestSelectionColorsCustom(t *testing.T) {
	t.Parallel()
	custom := gui.RGBA(255, 0, 0, 50)
	th := &theme.Theme{SelectionFill: custom}
	fill, border := selectionColors(th)
	if fill != custom {
		t.Errorf("fill = %v, want custom %v", fill, custom)
	}
	// Border should still be default.
	if border != gui.RGBA(70, 130, 220, 180) {
		t.Errorf("border = %v, want default", border)
	}
}

// --- ensureOrigBounds ---

func TestEnsureOrigBoundsInitsDomain(t *testing.T) {
	t.Parallel()
	xAxis := axis.NewLinear(axis.LinearCfg{Min: 10, Max: 90})
	yAxis := axis.NewLinear(axis.LinearCfg{Min: -5, Max: 50})
	zs := zoomState{}

	ensureOrigBounds(&zs, xAxis, yAxis, true, true)

	if !zs.OrigStored {
		t.Fatal("OrigStored should be true")
	}
	if zs.OrigXMin != 10 || zs.OrigXMax != 90 {
		t.Errorf("OrigX = [%g,%g], want [10,90]",
			zs.OrigXMin, zs.OrigXMax)
	}
	if zs.OrigYMin != -5 || zs.OrigYMax != 50 {
		t.Errorf("OrigY = [%g,%g], want [-5,50]",
			zs.OrigYMin, zs.OrigYMax)
	}
	// When not zoomed, zoom domain should match orig.
	if zs.XMin != 10 || zs.XMax != 90 {
		t.Errorf("XMin/XMax = [%g,%g], want [10,90]",
			zs.XMin, zs.XMax)
	}
	if zs.YMin != -5 || zs.YMax != 50 {
		t.Errorf("YMin/YMax = [%g,%g], want [-5,50]",
			zs.YMin, zs.YMax)
	}
}

func TestEnsureOrigBoundsPreservesZoomed(t *testing.T) {
	t.Parallel()
	xAxis := axis.NewLinear(axis.LinearCfg{Min: 0, Max: 100})
	zs := zoomState{
		Zoomed: true,
		XMin:   20, XMax: 80,
	}

	ensureOrigBounds(&zs, xAxis, nil, true, false)

	if zs.OrigXMin != 0 || zs.OrigXMax != 100 {
		t.Errorf("OrigX = [%g,%g], want [0,100]",
			zs.OrigXMin, zs.OrigXMax)
	}
	// Zoomed domain should NOT be overwritten.
	if zs.XMin != 20 || zs.XMax != 80 {
		t.Errorf("XMin/XMax = [%g,%g], want [20,80]",
			zs.XMin, zs.XMax)
	}
}

func TestEnsureOrigBoundsNoop(t *testing.T) {
	t.Parallel()
	xAxis := axis.NewLinear(axis.LinearCfg{Min: 0, Max: 100})
	zs := zoomState{
		OrigStored: true,
		OrigXMin:   5, OrigXMax: 95,
	}

	ensureOrigBounds(&zs, xAxis, nil, true, false)

	// Already stored — should not overwrite.
	if zs.OrigXMin != 5 || zs.OrigXMax != 95 {
		t.Errorf("OrigX changed to [%g,%g], want [5,95]",
			zs.OrigXMin, zs.OrigXMax)
	}
}

func TestEnsureOrigBoundsNilAxes(t *testing.T) {
	t.Parallel()
	zs := zoomState{}
	ensureOrigBounds(&zs, nil, nil, true, true)

	if !zs.OrigStored {
		t.Fatal("OrigStored should be true even with nil axes")
	}
	// Domains stay zero.
	if zs.OrigXMin != 0 || zs.OrigXMax != 0 {
		t.Errorf("OrigX = [%g,%g], want [0,0]",
			zs.OrigXMin, zs.OrigXMax)
	}
}

func TestEnsureOrigBoundsPartialAxes(t *testing.T) {
	t.Parallel()
	yAxis := axis.NewLinear(axis.LinearCfg{Min: -10, Max: 10})
	zs := zoomState{}

	// zoomX=true but xAxis=nil, zoomY=true with valid yAxis.
	ensureOrigBounds(&zs, nil, yAxis, true, true)

	if zs.OrigXMin != 0 || zs.OrigXMax != 0 {
		t.Errorf("OrigX = [%g,%g], want [0,0] (nil axis)",
			zs.OrigXMin, zs.OrigXMax)
	}
	if zs.OrigYMin != -10 || zs.OrigYMax != 10 {
		t.Errorf("OrigY = [%g,%g], want [-10,10]",
			zs.OrigYMin, zs.OrigYMax)
	}
}
