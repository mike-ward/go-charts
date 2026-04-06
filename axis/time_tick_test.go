package axis

import (
	"math"
	"testing"
	"time"
)

func TestTimeTicksHours(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(24 * time.Hour)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) < 4 || len(ticks) > 25 {
		t.Errorf("24h range: got %d ticks, want 4-25", len(ticks))
	}
}

func TestTimeTicksDays(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(0, 0, 30)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) < 4 || len(ticks) > 32 {
		t.Errorf("30d range: got %d ticks, want 4-32", len(ticks))
	}
}

func TestTimeTicksMonths(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(2, 0, 0)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) < 4 {
		t.Errorf("2y range: got %d ticks, want >= 4", len(ticks))
	}
	// Verify ticks are on month boundaries (in UTC).
	for _, tk := range ticks {
		tm := secondsToTime(tk.Value).UTC()
		if tm.Day() != 1 || tm.Hour() != 0 {
			t.Errorf("tick not on month start: %v", tm)
		}
	}
}

func TestTimeTicksYears(t *testing.T) {
	t0 := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(50, 0, 0)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) < 3 || len(ticks) > 52 {
		t.Errorf("50y range: got %d ticks, want 3-52", len(ticks))
	}
}

func TestTimeTicksSeconds(t *testing.T) {
	t0 := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	t1 := t0.Add(2 * time.Minute)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) < 4 {
		t.Errorf("2m range: got %d ticks, want >= 4", len(ticks))
	}
}

func TestTimeTicksZeroRange(t *testing.T) {
	now := time.Now()
	a := NewTime(TimeCfg{Min: now, Max: now})
	ticks := a.Ticks(0, 800)
	if ticks != nil {
		t.Errorf("zero range: expected nil, got %d ticks", len(ticks))
	}
}

func TestTimeTicksInvertedRange(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if ticks != nil {
		t.Errorf("inverted: expected nil, got %d ticks", len(ticks))
	}
}

func TestTimeTicksAlignment(t *testing.T) {
	// Daily ticks should align to midnight.
	t0 := time.Date(2025, 3, 5, 14, 30, 0, 0, time.UTC)
	t1 := t0.AddDate(0, 0, 14)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	for _, tk := range ticks {
		tm := secondsToTime(tk.Value).UTC()
		if tm.Hour() != 0 || tm.Minute() != 0 || tm.Second() != 0 {
			t.Errorf("tick not at midnight: %v", tm)
		}
	}
}

func TestTimeTicksWideRange(t *testing.T) {
	t0 := time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	if len(ticks) == 0 {
		t.Error("wide range: expected ticks")
	}
	for _, tk := range ticks {
		if math.IsNaN(tk.Value) || math.IsInf(tk.Value, 0) {
			t.Errorf("non-finite tick value: %v", tk.Value)
		}
		if math.IsNaN(float64(tk.Position)) ||
			math.IsInf(float64(tk.Position), 0) {
			t.Errorf("non-finite position: %v", tk.Position)
		}
	}
}

func TestTimeTicksCustomFormat(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(1, 0, 0)
	a := NewTime(TimeCfg{
		Min:        t0,
		Max:        t1,
		TickFormat: func(float64) string { return "X" },
	})
	ticks := a.Ticks(0, 800)
	for _, tk := range ticks {
		if tk.Label != "X" {
			t.Errorf("label = %q, want \"X\"", tk.Label)
		}
	}
}

func TestTimeTicksExplicitFormat(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(1, 0, 0)
	a := NewTime(TimeCfg{
		Min:    t0,
		Max:    t1,
		Format: "01/02",
	})
	ticks := a.Ticks(0, 800)
	if len(ticks) == 0 {
		t.Fatal("expected ticks")
	}
	// All ticks should use the explicit format "01/02".
	for _, tk := range ticks {
		tm := secondsToTime(tk.Value).In(t0.Location())
		want := tm.Format("01/02")
		if tk.Label != want {
			t.Errorf("label = %q, want %q", tk.Label, want)
		}
	}
}

func TestTimeTicksMonotonicPositions(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(0, 6, 0)
	a := NewTime(TimeCfg{Min: t0, Max: t1})
	ticks := a.Ticks(0, 800)
	for i := 1; i < len(ticks); i++ {
		if ticks[i].Position <= ticks[i-1].Position {
			t.Errorf("non-monotonic at %d: %.1f <= %.1f",
				i, ticks[i].Position, ticks[i-1].Position)
		}
	}
}

func TestTimeSetRangeDomain(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.AddDate(1, 0, 0)
	a := NewTime(TimeCfg{Min: t0, Max: t1})

	newMin := timeToSeconds(t0.AddDate(0, 3, 0))
	newMax := timeToSeconds(t0.AddDate(0, 9, 0))
	a.SetRange(newMin, newMax)

	lo, hi := a.Domain()
	if math.Abs(lo-newMin) > 1 || math.Abs(hi-newMax) > 1 {
		t.Errorf("Domain() = (%v, %v), want (%v, %v)",
			lo, hi, newMin, newMax)
	}
}
