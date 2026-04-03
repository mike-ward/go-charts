package fmath

import (
	"math"
	"testing"
)

func TestFinite(t *testing.T) {
	tests := []struct {
		name string
		v    float64
		want bool
	}{
		{"zero", 0, true},
		{"positive", 1.5, true},
		{"negative", -1.5, true},
		{"NaN", math.NaN(), false},
		{"PosInf", math.Inf(1), false},
		{"NegInf", math.Inf(-1), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Finite(tc.v); got != tc.want {
				t.Errorf("Finite(%v) = %v, want %v", tc.v, got, tc.want)
			}
		})
	}
}
