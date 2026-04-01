package axis

import (
	"fmt"
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

func TestLinearTickFormat(t *testing.T) {
	a := NewLinear(LinearCfg{
		Min: 0, Max: 100,
		TickFormat: func(v float64) string {
			return fmt.Sprintf("$%.0f", v)
		},
	})
	ticks := a.Ticks(0, 500)
	if len(ticks) == 0 {
		t.Fatal("expected ticks")
	}
	for _, tk := range ticks {
		if tk.Label[0] != '$' {
			t.Errorf("tick label %q missing $ prefix", tk.Label)
		}
	}
}

func TestLinearTickFormatNil(t *testing.T) {
	a := NewLinear(LinearCfg{Min: 0, Max: 10})
	ticks := a.Ticks(0, 100)
	if len(ticks) == 0 {
		t.Fatal("expected ticks")
	}
	// Default format should not be empty.
	for _, tk := range ticks {
		if tk.Label == "" {
			t.Errorf("tick label should not be empty")
		}
	}
}
