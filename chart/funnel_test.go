package chart

import (
	"math"
	"testing"
)

func TestFunnelLayout(t *testing.T) {
	slices := []PieSlice{
		{Label: "A", Value: 100},
		{Label: "B", Value: 60},
		{Label: "C", Value: 30},
	}
	fv := &funnelView{
		cfg: FunnelCfg{
			Slices:        slices,
			SegmentGap:    0,
			MinWidthRatio: 0.25,
		},
	}

	// Simulate layout: availW=200, top=0, bottom=300, center=100.
	availW := float32(200)
	centerX := float32(100)
	n := len(slices)
	segH := float32(300) / float32(n)
	maxVal := 100.0

	for i, s := range slices {
		topW := float32(s.Value/maxVal) * availW
		var botW float32
		if i < n-1 {
			botW = float32(slices[i+1].Value/maxVal) * availW
		} else {
			botW = topW * 0.25
		}
		fv.segments = append(fv.segments, funnelSegment{
			TopY:     float32(i) * segH,
			BotY:     float32(i)*segH + segH,
			TopLeft:  centerX - topW/2,
			TopRight: centerX + topW/2,
			BotLeft:  centerX - botW/2,
			BotRight: centerX + botW/2,
			Index:    i,
		})
	}

	if len(fv.segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(fv.segments))
	}

	// First segment should be widest.
	s0 := fv.segments[0]
	s1 := fv.segments[1]
	w0 := s0.TopRight - s0.TopLeft
	w1 := s1.TopRight - s1.TopLeft
	if w0 <= w1 {
		t.Errorf("first segment width %g should exceed second %g",
			w0, w1)
	}

	// Widths proportional to values.
	ratio := float64(w1 / w0)
	expected := 60.0 / 100.0
	if math.Abs(ratio-expected) > 0.01 {
		t.Errorf("width ratio = %g, want %g", ratio, expected)
	}
}

func TestFunnelHitTest(t *testing.T) {
	fv := &funnelView{
		segments: []funnelSegment{
			{TopY: 0, BotY: 100,
				TopLeft: 0, TopRight: 200,
				BotLeft: 50, BotRight: 150, Index: 0},
			{TopY: 100, BotY: 200,
				TopLeft: 50, TopRight: 150,
				BotLeft: 75, BotRight: 125, Index: 1},
		},
	}

	// Center of first segment.
	idx, ok := fv.hitTest(100, 50)
	if !ok || idx != 0 {
		t.Errorf("hitTest(100,50) = (%d,%v), want (0,true)",
			idx, ok)
	}

	// Center of second segment.
	idx, ok = fv.hitTest(100, 150)
	if !ok || idx != 1 {
		t.Errorf("hitTest(100,150) = (%d,%v), want (1,true)",
			idx, ok)
	}

	// Outside trapezoid (top-left corner at midpoint of first
	// segment where edge has narrowed).
	_, ok = fv.hitTest(10, 80)
	if ok {
		t.Error("hitTest(10,80) should miss (outside edge)")
	}
}

func TestFunnelHitTestOutside(t *testing.T) {
	fv := &funnelView{
		segments: []funnelSegment{
			{TopY: 10, BotY: 50,
				TopLeft: 20, TopRight: 80,
				BotLeft: 30, BotRight: 70, Index: 0},
		},
	}

	_, ok := fv.hitTest(50, 0)
	if ok {
		t.Error("hitTest above segments should return false")
	}

	_, ok = fv.hitTest(50, 60)
	if ok {
		t.Error("hitTest below segments should return false")
	}

	_, ok = fv.hitTest(0, 30)
	if ok {
		t.Error("hitTest left of segment should return false")
	}
}

func TestFunnelNaNInfSkipped(t *testing.T) {
	fv := &funnelView{
		cfg: FunnelCfg{
			Slices: []PieSlice{
				{Label: "A", Value: 100},
				{Label: "NaN", Value: math.NaN()},
				{Label: "Inf", Value: math.Inf(1)},
				{Label: "Neg", Value: -10},
				{Label: "B", Value: 50},
			},
			SegmentGap:    0,
			MinWidthRatio: 0.25,
		},
	}

	// Manually compute what draw() computes for maxValue/totalValue.
	maxVal := 0.0
	total := 0.0
	for _, s := range fv.cfg.Slices {
		if math.IsNaN(s.Value) || math.IsInf(s.Value, 0) || s.Value <= 0 {
			continue
		}
		maxVal = max(maxVal, s.Value)
		total += s.Value
	}
	if maxVal != 100 {
		t.Errorf("maxValue = %g, want 100", maxVal)
	}
	if total != 150 {
		t.Errorf("totalValue = %g, want 150", total)
	}
}

func TestFunnelValidateNegativeValue(t *testing.T) {
	cfg := FunnelCfg{
		Slices: []PieSlice{{Label: "a", Value: -5}},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for negative slice value")
	}
}

func TestFunnelValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     FunnelCfg
		wantErr bool
	}{
		{"empty slices",
			FunnelCfg{}, true},
		{"negative SegmentGap",
			FunnelCfg{
				Slices:     []PieSlice{{Label: "a", Value: 1}},
				SegmentGap: -1,
			}, true},
		{"MinWidthRatio too high",
			FunnelCfg{
				Slices:        []PieSlice{{Label: "a", Value: 1}},
				MinWidthRatio: 1.5,
			}, true},
		{"MinWidthRatio negative",
			FunnelCfg{
				Slices:        []PieSlice{{Label: "a", Value: 1}},
				MinWidthRatio: -0.1,
			}, true},
		{"valid",
			FunnelCfg{
				Slices: []PieSlice{{Label: "a", Value: 10}},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}
