// Package fmath provides shared floating-point helpers.
package fmath

import "math"

// Finite reports whether v is neither NaN nor +/-Inf.
func Finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
