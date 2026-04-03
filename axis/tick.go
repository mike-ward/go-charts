package axis

import (
	"math"

	"github.com/mike-ward/go-charts/internal/fmath"
)

// NiceNumber computes a "nice" number approximately equal to the
// given value. If round is true, it rounds; otherwise it takes
// the ceiling.
func NiceNumber(value float64, round bool) float64 {
	if value == 0 {
		return 0
	}
	sign := 1.0
	if value < 0 {
		sign = -1
		value = -value
	}
	exp := math.Floor(math.Log10(value))
	frac := value / math.Pow(10, exp)

	var nice float64
	if round {
		switch {
		case frac < 1.5:
			nice = 1
		case frac < 3:
			nice = 2
		case frac < 7:
			nice = 5
		default:
			nice = 10
		}
	} else {
		switch {
		case frac <= 1:
			nice = 1
		case frac <= 2:
			nice = 2
		case frac <= 5:
			nice = 5
		default:
			nice = 10
		}
	}
	return sign * nice * math.Pow(10, exp)
}

// GenerateNiceTicks generates evenly-spaced tick values for the
// given data range, targeting approximately maxTicks ticks.
// Non-finite or degenerate inputs produce a safe fallback.
func GenerateNiceTicks(dataMin, dataMax float64, maxTicks int) []float64 {
	if maxTicks < 2 {
		maxTicks = 2
	}

	// Guard: non-finite inputs.
	if !fmath.Finite(dataMin) || !fmath.Finite(dataMax) {
		if fmath.Finite(dataMin) {
			return []float64{dataMin}
		}
		if fmath.Finite(dataMax) {
			return []float64{dataMax}
		}
		return nil
	}

	rangeVal := NiceNumber(dataMax-dataMin, false)
	if !fmath.Finite(rangeVal) {
		return []float64{dataMin}
	}

	spacing := NiceNumber(rangeVal/float64(maxTicks-1), true)
	if spacing <= 0 || !fmath.Finite(spacing) {
		return []float64{dataMin}
	}

	niceMin := math.Floor(dataMin/spacing) * spacing
	niceMax := math.Ceil(dataMax/spacing) * spacing
	if !fmath.Finite(niceMin) || !fmath.Finite(niceMax) {
		return []float64{dataMin}
	}

	const maxTickCount = 500
	estCap := (niceMax - niceMin) / spacing
	capVal := maxTickCount
	if estCap >= 0 && estCap < float64(maxTickCount) {
		capVal = int(estCap) + 2
	}
	ticks := make([]float64, 0, capVal)
	for v := niceMin; v <= niceMax+spacing*0.5; v += spacing {
		ticks = append(ticks, v)
		if len(ticks) >= maxTickCount {
			break
		}
		// Guard: float64 stall (v+spacing == v).
		if v+spacing == v {
			break
		}
	}
	return ticks
}
