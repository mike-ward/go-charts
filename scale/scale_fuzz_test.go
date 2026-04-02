package scale

import (
	"math"
	"testing"
)

func FuzzLinearTransform(f *testing.F) {
	f.Add(0.0, 100.0, 50.0, float32(0), float32(500))
	f.Add(0.0, 0.0, 0.0, float32(0), float32(0))
	f.Add(math.NaN(), math.NaN(), math.NaN(), float32(0), float32(100))
	f.Add(-1e308, 1e308, 0.0, float32(-1000), float32(1000))

	f.Fuzz(func(t *testing.T, dMin, dMax, value float64, pMin, pMax float32) {
		s := NewLinear(dMin, dMax)
		got := s.Transform(value, pMin, pMax)
		if math.IsNaN(float64(got)) {
			t.Errorf("Transform returned NaN for value=%v domain=[%v,%v]",
				value, dMin, dMax)
		}
	})
}

func FuzzLogTransform(f *testing.F) {
	f.Add(1.0, 1000.0, 10.0, 100.0, float32(0), float32(300))
	f.Add(0.0, 0.0, 0.0, 0.0, float32(0), float32(0))
	f.Add(1.0, 1.0, 1.0, 1.0, float32(0), float32(100))
	f.Add(math.NaN(), math.NaN(), math.NaN(), math.NaN(), float32(0), float32(100))

	f.Fuzz(func(t *testing.T, dMin, dMax, base, value float64, pMin, pMax float32) {
		s := NewLog(dMin, dMax, base)
		got := s.Transform(value, pMin, pMax)
		if math.IsNaN(float64(got)) {
			t.Errorf("Transform returned NaN for value=%v domain=[%v,%v] base=%v",
				value, dMin, dMax, base)
		}
		if math.IsInf(float64(got), 0) {
			t.Errorf("Transform returned Inf for value=%v domain=[%v,%v] base=%v",
				value, dMin, dMax, base)
		}
	})
}
