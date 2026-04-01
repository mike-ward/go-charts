package axis

import (
	"math"
	"testing"
)

func FuzzNiceNumber(f *testing.F) {
	f.Add(0.0, true)
	f.Add(7.5, false)
	f.Add(-1e300, true)
	f.Add(math.NaN(), true)
	f.Add(math.Inf(1), false)
	f.Add(math.Inf(-1), true)
	f.Add(5e-324, true) // smallest positive float64

	f.Fuzz(func(t *testing.T, value float64, round bool) {
		// Must not panic.
		_ = NiceNumber(value, round)
	})
}

func FuzzGenerateNiceTicks(f *testing.F) {
	f.Add(0.0, 100.0, 8)
	f.Add(0.0, 0.0, 8)
	f.Add(100.0, 0.0, 8)
	f.Add(math.NaN(), math.NaN(), 8)
	f.Add(math.Inf(1), math.Inf(-1), 8)
	f.Add(-1e308, 1e308, 2)
	f.Add(0.0, 5e-324, 100)

	f.Fuzz(func(t *testing.T, dataMin, dataMax float64, maxTicks int) {
		ticks := GenerateNiceTicks(dataMin, dataMax, maxTicks)

		// Must terminate (we got here) and not exceed cap.
		if len(ticks) > 500 {
			t.Errorf("too many ticks: %d", len(ticks))
		}

		// All tick values must be finite (if any returned).
		for i, v := range ticks {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				t.Errorf("tick[%d] non-finite: %v", i, v)
			}
		}
	})
}
