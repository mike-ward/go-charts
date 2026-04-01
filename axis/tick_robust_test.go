package axis

import (
	"math"
	"testing"
)

func TestNiceNumberEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		round bool
		want  float64
		isNaN bool // expect NaN result
	}{
		{"negative", -7.5, true, -10, false},
		{"very small", 1e-300, true, 1e-300, false},
		{"very large", 1e300, true, 1e300, false},
		{"zero", 0, true, 0, false},
		{"NaN", math.NaN(), true, 0, true},
		{"+Inf", math.Inf(1), true, 0, true},
		{"-Inf", math.Inf(-1), true, 0, true},
		{"negative small", -1e-300, false, -1e-300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NiceNumber(tt.value, tt.round)
			if tt.isNaN {
				// NaN/Inf input may produce NaN/Inf; just verify
				// no panic. GenerateNiceTicks handles these.
				return
			}
			if math.IsNaN(got) {
				t.Errorf("NiceNumber(%v, %v) = NaN", tt.value, tt.round)
				return
			}
			// Verify same sign and within order of magnitude.
			if tt.want != 0 && math.Abs(got/tt.want-1) > 1 {
				t.Errorf("NiceNumber(%v, %v) = %v, want ~%v",
					tt.value, tt.round, got, tt.want)
			}
		})
	}
}

func TestGenerateNiceTicksEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		min, max float64
		maxTicks int
		wantNil  bool
		maxLen   int // max expected result length; 0 = don't check
	}{
		{"NaN min", math.NaN(), 100, 8, false, 1},
		{"NaN max", 0, math.NaN(), 8, false, 1},
		{"both NaN", math.NaN(), math.NaN(), 8, true, 0},
		{"+Inf min", math.Inf(1), 100, 8, false, 1},
		{"-Inf max", 0, math.Inf(-1), 8, false, 1},
		{"both Inf", math.Inf(1), math.Inf(-1), 8, true, 0},
		{"min == max", 50, 50, 8, false, 1},
		{"inverted range", 100, 0, 8, false, 1},
		{"very large range", -1e300, 1e300, 8, false, 500},
		{"very small range", 0, 1e-300, 8, false, 500},
		{"maxTicks 0", 0, 100, 0, false, 500},
		{"maxTicks 1", 0, 100, 1, false, 500},
		{"maxTicks negative", 0, 100, -5, false, 500},
		{"normal", 0, 100, 5, false, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticks := GenerateNiceTicks(tt.min, tt.max, tt.maxTicks)
			if tt.wantNil {
				if ticks != nil {
					t.Errorf("expected nil, got %v", ticks)
				}
				return
			}
			if ticks == nil {
				t.Fatal("unexpected nil")
			}
			if tt.maxLen > 0 && len(ticks) > tt.maxLen {
				t.Errorf("got %d ticks, max expected %d",
					len(ticks), tt.maxLen)
			}
			// Verify all tick values are finite.
			for i, v := range ticks {
				if math.IsNaN(v) || math.IsInf(v, 0) {
					t.Errorf("tick[%d] = %v (non-finite)", i, v)
				}
			}
		})
	}
}
