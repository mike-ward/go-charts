package chart

import (
	"math"
	"testing"
)

func TestBubbleRadius(t *testing.T) {
	tests := []struct {
		name          string
		z, zMin, zMax float64
		minR, maxR    float32
		want          float32
	}{
		{
			name: "minimum Z returns minR",
			z:    0, zMin: 0, zMax: 100,
			minR: 4, maxR: 30,
			want: 4,
		},
		{
			name: "maximum Z returns maxR",
			z:    100, zMin: 0, zMax: 100,
			minR: 4, maxR: 30,
			want: 30,
		},
		{
			name: "midpoint Z returns sqrt-scaled radius",
			z:    50, zMin: 0, zMax: 100,
			minR: 4, maxR: 30,
			want: 4 + 26*float32(math.Sqrt(0.5)),
		},
		{
			name: "equal zMin zMax returns midpoint radius",
			z:    10, zMin: 10, zMax: 10,
			minR: 4, maxR: 30,
			want: 17,
		},
		{
			name: "zMax less than zMin returns midpoint radius",
			z:    5, zMin: 10, zMax: 5,
			minR: 4, maxR: 30,
			want: 17,
		},
		{
			name: "quarter Z",
			z:    25, zMin: 0, zMax: 100,
			minR: 4, maxR: 30,
			want: 4 + 26*float32(math.Sqrt(0.25)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bubbleRadius(tt.z, tt.zMin, tt.zMax, tt.minR, tt.maxR)
			if diff := got - tt.want; diff > 0.001 || diff < -0.001 {
				t.Errorf("bubbleRadius(%g, %g, %g, %g, %g) = %g, want %g",
					tt.z, tt.zMin, tt.zMax, tt.minR, tt.maxR, got, tt.want)
			}
		})
	}
}

func TestBubbleRadiusNegativeZ(t *testing.T) {
	// Z below zMin should clamp to minR, not produce NaN
	// via math.Sqrt(negative).
	got := bubbleRadius(-10, 0, 100, 4, 30)
	if got != 4 {
		t.Errorf("bubbleRadius(-10, 0, 100) = %g, want 4 (minR)", got)
	}
}

func TestBubbleRadiusAboveMax(t *testing.T) {
	// Z above zMax should clamp to maxR.
	got := bubbleRadius(200, 0, 100, 4, 30)
	if got != 30 {
		t.Errorf("bubbleRadius(200, 0, 100) = %g, want 30 (maxR)", got)
	}
}

func TestBubbleMarkerShape(t *testing.T) {
	cfg := &BubbleCfg{
		Marker:  MarkerCircle,
		Markers: []MarkerShape{MarkerSquare, MarkerDiamond},
	}

	if got := bubbleMarkerShape(cfg, 0); got != MarkerSquare {
		t.Errorf("index 0: got %v, want MarkerSquare", got)
	}
	if got := bubbleMarkerShape(cfg, 1); got != MarkerDiamond {
		t.Errorf("index 1: got %v, want MarkerDiamond", got)
	}
	// Falls back to Marker for out-of-range index.
	if got := bubbleMarkerShape(cfg, 2); got != MarkerCircle {
		t.Errorf("index 2: got %v, want MarkerCircle (fallback)", got)
	}
}
