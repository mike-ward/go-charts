package axis

import (
	"math"
	"testing"
	"time"
)

func TestTimeZeroValues(t *testing.T) {
	a := NewTime(TimeCfg{})
	got := a.Transform(0, 0, 500)
	if math.IsNaN(float64(got)) || math.IsInf(float64(got), 0) {
		t.Errorf("Transform with zero times = %v", got)
	}
}

func TestTimeMinEqualsMax(t *testing.T) {
	now := time.Now()
	a := NewTime(TimeCfg{Min: now, Max: now})
	got := a.Transform(timeToSeconds(now), 0, 500)
	if got != 0 {
		t.Errorf("Transform(min==max) = %v, want pixelMin 0", got)
	}
}

func TestTimeInverted(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewTime(TimeCfg{Min: t2, Max: t1}) // inverted

	got := a.Transform(timeToSeconds(t1), 0, 500)
	if math.IsNaN(float64(got)) || math.IsInf(float64(got), 0) {
		t.Errorf("Transform inverted = %v (non-finite)", got)
	}
}

func TestTimeWideRange(t *testing.T) {
	// Year 1000 to 3000 would overflow UnixNano. Now safe.
	t1 := time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewTime(TimeCfg{Min: t1, Max: t2})

	mid := timeToSeconds(t1) + (timeToSeconds(t2)-timeToSeconds(t1))/2
	got := a.Transform(mid, 0, 1000)
	if math.IsNaN(float64(got)) || math.IsInf(float64(got), 0) {
		t.Errorf("Transform wide range = %v (non-finite)", got)
	}
	if math.Abs(float64(got)-500) > 1 {
		t.Errorf("Transform midpoint = %v, want ~500", got)
	}
}

func TestTimeInverseZeroPixelRange(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewTime(TimeCfg{Min: t1, Max: t2})

	got := a.Inverse(100, 100, 100)
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Errorf("Inverse zero pixel range = %v (non-finite)", got)
	}
}

func TestTimeInverseMinEqualsMax(t *testing.T) {
	now := time.Now()
	a := NewTime(TimeCfg{Min: now, Max: now})
	got := a.Inverse(250, 0, 500)
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Errorf("Inverse min==max = %v (non-finite)", got)
	}
}
