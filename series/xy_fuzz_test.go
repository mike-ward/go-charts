package series

import (
	"math"
	"testing"
)

func FuzzXYBounds(f *testing.F) {
	f.Add(0.0, 0.0, 1.0, 1.0)
	f.Add(math.NaN(), math.NaN(), math.NaN(), math.NaN())
	f.Add(math.Inf(1), math.Inf(-1), 0.0, 0.0)
	f.Add(1e308, -1e308, 1e-308, -1e-308)
	f.Add(0.0, 0.0, 0.0, 0.0)

	f.Fuzz(func(t *testing.T, x1, y1, x2, y2 float64) {
		s := NewXY(XYCfg{
			Points: []Point{{x1, y1}, {x2, y2}},
		})
		minX, maxX, minY, maxY := s.Bounds()

		// Bounds must be finite or all zero (no finite input).
		allZero := minX == 0 && maxX == 0 && minY == 0 && maxY == 0
		if allZero {
			return
		}
		if math.IsNaN(minX) || math.IsNaN(maxX) ||
			math.IsNaN(minY) || math.IsNaN(maxY) {
			t.Errorf("NaN in bounds: (%v,%v,%v,%v)",
				minX, maxX, minY, maxY)
		}
		if math.IsInf(minX, 0) || math.IsInf(maxX, 0) ||
			math.IsInf(minY, 0) || math.IsInf(maxY, 0) {
			t.Errorf("Inf in bounds: (%v,%v,%v,%v)",
				minX, maxX, minY, maxY)
		}
		if minX > maxX {
			t.Errorf("minX %v > maxX %v", minX, maxX)
		}
		if minY > maxY {
			t.Errorf("minY %v > maxY %v", minY, maxY)
		}
	})
}
