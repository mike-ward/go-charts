package axis

import (
	"math"
	"testing"
)

func TestLogTicksBase10(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 1000})
	ticks := a.Ticks(0, 600)
	if len(ticks) == 0 {
		t.Fatal("expected ticks, got none")
	}

	// Collect major tick values.
	var majors []float64
	for _, tk := range ticks {
		if !tk.Minor {
			majors = append(majors, tk.Value)
		}
	}
	want := []float64{1, 10, 100, 1000}
	if len(majors) != len(want) {
		t.Fatalf("major ticks = %v, want %v", majors, want)
	}
	for i, v := range majors {
		if math.Abs(v-want[i]) > 1e-9 {
			t.Errorf("major[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestLogTicksMinorBase10(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 100})
	ticks := a.Ticks(0, 600)

	var minors int
	for _, tk := range ticks {
		if tk.Minor {
			minors++
		}
	}
	// Between 1-10: 2,3,4,5,6,7,8,9 = 8 minor ticks
	// Between 10-100: 20,30,...,90 = 8 minor ticks
	// Total = 16
	if minors != 16 {
		t.Errorf("minor tick count = %d, want 16", minors)
	}
}

func TestLogTicksBase2(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 16, Base: 2})
	ticks := a.Ticks(0, 600)

	var majors []float64
	for _, tk := range ticks {
		if !tk.Minor {
			majors = append(majors, tk.Value)
		}
	}
	want := []float64{1, 2, 4, 8, 16}
	if len(majors) != len(want) {
		t.Fatalf("base-2 majors = %v, want %v", majors, want)
	}
	for i, v := range majors {
		if math.Abs(v-want[i]) > 1e-9 {
			t.Errorf("major[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestLogTicksNonPositiveDomain(t *testing.T) {
	for _, tc := range []struct {
		name     string
		min, max float64
	}{
		{"zero min", 0, 100},
		{"negative min", -1, 100},
		{"negative both", -100, -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := NewLog(LogCfg{Min: tc.min, Max: tc.max})
			ticks := a.Ticks(0, 600)
			if ticks != nil {
				t.Errorf("expected nil, got %d ticks", len(ticks))
			}
		})
	}
}

func TestLogTicksDegenerateRange(t *testing.T) {
	a := NewLog(LogCfg{Min: 10, Max: 10})
	ticks := a.Ticks(0, 600)
	if ticks != nil {
		t.Errorf("min==max: expected nil, got %d ticks", len(ticks))
	}
}

func TestLogTicksInvertedRange(t *testing.T) {
	a := NewLog(LogCfg{Min: 100, Max: 1})
	ticks := a.Ticks(0, 600)
	if ticks != nil {
		t.Errorf("inverted: expected nil, got %d ticks", len(ticks))
	}
}

func TestLogTicksVeryWideRange(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 1e100})
	ticks := a.Ticks(0, 600)
	if len(ticks) > 500 {
		t.Errorf("tick count = %d, want <= 500", len(ticks))
	}
	if len(ticks) == 0 {
		t.Error("expected some ticks for wide range")
	}
}

func TestLogTicksCustomFormat(t *testing.T) {
	a := NewLog(LogCfg{
		Min: 1, Max: 1000,
		TickFormat: func(v float64) string { return "custom" },
	})
	ticks := a.Ticks(0, 600)
	for _, tk := range ticks {
		if tk.Minor {
			if tk.Label != "" {
				t.Errorf("minor label = %q, want empty", tk.Label)
			}
			continue
		}
		if tk.Label != "custom" {
			t.Errorf("major label = %q, want \"custom\"", tk.Label)
		}
	}
}

func TestLogTicksMonotonicPositions(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 10000})
	ticks := a.Ticks(0, 800)
	for i := 1; i < len(ticks); i++ {
		if ticks[i].Position <= ticks[i-1].Position {
			t.Errorf("non-monotonic at %d: %.1f <= %.1f",
				i, ticks[i].Position, ticks[i-1].Position)
		}
	}
}

func TestLogTicksOverrideDomain(t *testing.T) {
	a := NewLog(LogCfg{Min: 5, Max: 500})
	a.SetOverrideDomain(true)
	ticks := a.Ticks(0, 600)
	for _, tk := range ticks {
		if tk.Value < 5*(1-1e-9) || tk.Value > 500*(1+1e-9) {
			t.Errorf("tick %v outside domain [5, 500]", tk.Value)
		}
	}
}

func TestLogSetRangeDomain(t *testing.T) {
	a := NewLog(LogCfg{Min: 1, Max: 100})
	a.SetRange(10, 10000)
	lo, hi := a.Domain()
	if lo != 10 || hi != 10000 {
		t.Errorf("Domain() = (%v, %v), want (10, 10000)", lo, hi)
	}
}
