package chart

import (
	"math"
	"testing"
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
