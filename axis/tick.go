package axis

import "math"

// NiceNumber computes a "nice" number approximately equal to the
// given value. If round is true, it rounds; otherwise it takes
// the ceiling.
func NiceNumber(value float64, round bool) float64 {
	if value == 0 {
		return 0
	}
	exp := math.Floor(math.Log10(math.Abs(value)))
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
	return nice * math.Pow(10, exp)
}

// GenerateNiceTicks generates evenly-spaced tick values for the
// given data range, targeting approximately maxTicks ticks.
func GenerateNiceTicks(dataMin, dataMax float64, maxTicks int) []float64 {
	if maxTicks < 2 {
		maxTicks = 2
	}
	rangeVal := NiceNumber(dataMax-dataMin, false)
	spacing := NiceNumber(rangeVal/float64(maxTicks-1), true)
	if spacing == 0 {
		return []float64{dataMin}
	}

	niceMin := math.Floor(dataMin/spacing) * spacing
	niceMax := math.Ceil(dataMax/spacing) * spacing

	var ticks []float64
	for v := niceMin; v <= niceMax+spacing*0.5; v += spacing {
		ticks = append(ticks, v)
	}
	return ticks
}
