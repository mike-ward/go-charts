package chart

import "testing"

func TestShouldReduceDetailNilWindow(t *testing.T) {
	if shouldReduceDetail(nil) {
		t.Error("shouldReduceDetail(nil) = true, want false")
	}
}

func TestAvgFrameMsNilWindow(t *testing.T) {
	if got := avgFrameMs(nil); got != 0 {
		t.Errorf("avgFrameMs(nil) = %v, want 0", got)
	}
}

func TestUpdateFPSTrackerNilWindow(t *testing.T) {
	// Should not panic.
	updateFPSTracker(nil)
}
