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

func TestOverrideDomainPreventsTickExpansion(t *testing.T) {
	t.Parallel()
	a := NewLinear(LinearCfg{AutoRange: true})
	a.SetRange(1.3, 4.7) // not "nice" boundaries
	a.SetOverrideDomain(true)
	_ = a.Ticks(0, 400) // should not expand domain

	dMin, dMax := a.Domain()
	if dMin != 1.3 || dMax != 4.7 {
		t.Errorf("domain changed to [%g, %g], want [1.3, 4.7]",
			dMin, dMax)
	}
}

func TestOverrideDomainFalseAllowsExpansion(t *testing.T) {
	t.Parallel()
	a := NewLinear(LinearCfg{AutoRange: true})
	a.SetRange(1.3, 4.7)
	// overrideDomain defaults to false — Ticks should expand.
	_ = a.Ticks(0, 400)

	dMin, dMax := a.Domain()
	// Nice ticks for [1.3, 4.7] will expand to something like [1, 5].
	if dMin == 1.3 && dMax == 4.7 {
		t.Error("domain was not expanded by Ticks with AutoRange")
	}
	if dMin > 1.3 || dMax < 4.7 {
		t.Errorf("domain [%g, %g] does not contain [1.3, 4.7]",
			dMin, dMax)
	}
}

func TestGenerateNiceTicksPrecision(t *testing.T) {
	t.Parallel()
	// Verify ticks are cleanly snapped for various spacings,
	// including those with float64 representation issues.
	tests := []struct {
		min, max float64
		ticks    int
	}{
		{0, 1, 10},   // spacing ~ 0.1
		{0, 2, 10},   // spacing ~ 0.2
		{0, 3, 10},   // spacing ~ 0.3
		{0, 0.5, 10}, // spacing ~ 0.05
		{0, 50, 10},  // spacing ~ 5 (integer)
	}
	for _, tt := range tests {
		ticks := GenerateNiceTicks(tt.min, tt.max, tt.ticks)
		if len(ticks) < 2 {
			t.Errorf("[%g,%g]: got %d ticks, want >= 2",
				tt.min, tt.max, len(ticks))
			continue
		}
		spacing := ticks[1] - ticks[0]
		for i := 1; i < len(ticks); i++ {
			gap := ticks[i] - ticks[i-1]
			if math.Abs(gap-spacing) > 1e-9 {
				t.Errorf("[%g,%g]: uneven gap at tick %d: %g vs %g",
					tt.min, tt.max, i, gap, spacing)
			}
		}
	}
}

func TestDomainMethod(t *testing.T) {
	t.Parallel()
	a := NewLinear(LinearCfg{Min: 10, Max: 20})
	lo, hi := a.Domain()
	if lo != 10 || hi != 20 {
		t.Errorf("Domain() = [%g, %g], want [10, 20]", lo, hi)
	}
}
