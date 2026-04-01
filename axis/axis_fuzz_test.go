package axis

import (
	"math"
	"testing"
)

func FuzzCategoryTransform(f *testing.F) {
	f.Add(0.0, float32(0), float32(500))
	f.Add(5.0, float32(0), float32(0))
	f.Add(-1.0, float32(-100), float32(100))
	f.Add(1e10, float32(0), float32(1000))

	cats := []string{"A", "B", "C", "D", "E"}

	f.Fuzz(func(t *testing.T, value float64, pMin, pMax float32) {
		a := NewCategory(CategoryCfg{Categories: cats})

		got := a.Transform(value, pMin, pMax)
		if math.IsNaN(float64(got)) {
			t.Errorf("Transform NaN for value=%v pMin=%v pMax=%v",
				value, pMin, pMax)
		}
		if math.IsInf(float64(got), 0) {
			t.Errorf("Transform Inf for value=%v pMin=%v pMax=%v",
				value, pMin, pMax)
		}

		inv := a.Inverse(got, pMin, pMax)
		if math.IsNaN(inv) {
			t.Errorf("Inverse NaN for pixel=%v pMin=%v pMax=%v",
				got, pMin, pMax)
		}
		if math.IsInf(inv, 0) {
			t.Errorf("Inverse Inf for pixel=%v pMin=%v pMax=%v",
				got, pMin, pMax)
		}
	})
}
