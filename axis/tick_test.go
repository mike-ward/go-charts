package axis

import (
	"math"
	"testing"
)

func TestNiceNumber(t *testing.T) {
	tests := []struct {
		value float64
		round bool
		want  float64
	}{
		{0.7, true, 0.5},
		{3.5, true, 5},
		{7.5, true, 10},
		{12, true, 10},
		{0, true, 0},
	}
	for _, tt := range tests {
		got := NiceNumber(tt.value, tt.round)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("NiceNumber(%v, %v) = %v, want %v",
				tt.value, tt.round, got, tt.want)
		}
	}
}

func TestGenerateNiceTicks(t *testing.T) {
	ticks := GenerateNiceTicks(0, 100, 5)
	if len(ticks) == 0 {
		t.Fatal("expected ticks, got none")
	}
	if ticks[0] > 0 {
		t.Errorf("first tick %v should be <= 0", ticks[0])
	}
	last := ticks[len(ticks)-1]
	if last < 100 {
		t.Errorf("last tick %v should be >= 100", last)
	}
}
