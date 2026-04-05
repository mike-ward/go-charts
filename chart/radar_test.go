package chart

import (
	"math"
	"testing"
)

func TestRadarCfgValidate(t *testing.T) {
	threeAxes := []RadarAxis{
		{Label: "A", Max: 100},
		{Label: "B", Max: 100},
		{Label: "C", Max: 100},
	}
	oneSeries := []RadarSeries{
		{Name: "S1", Values: []float64{50, 60, 70}},
	}

	tests := []struct {
		name    string
		cfg     RadarCfg
		wantErr bool
	}{
		{
			name: "valid defaults",
			cfg: RadarCfg{
				Axes:   threeAxes,
				Series: oneSeries,
			},
			wantErr: false,
		},
		{
			name: "fewer than 3 axes",
			cfg: RadarCfg{
				Axes: []RadarAxis{
					{Label: "A", Max: 10},
					{Label: "B", Max: 10},
				},
				Series: []RadarSeries{
					{Name: "S", Values: []float64{1, 2}},
				},
			},
			wantErr: true,
		},
		{
			name: "series length mismatch",
			cfg: RadarCfg{
				Axes: threeAxes,
				Series: []RadarSeries{
					{Name: "S", Values: []float64{1, 2}},
				},
			},
			wantErr: true,
		},
		{
			name: "axis min >= max",
			cfg: RadarCfg{
				Axes: []RadarAxis{
					{Label: "A", Min: 10, Max: 5},
					{Label: "B", Max: 100},
					{Label: "C", Max: 100},
				},
				Series: oneSeries,
			},
			wantErr: true,
		},
		{
			name: "valid explicit min/max",
			cfg: RadarCfg{
				Axes: []RadarAxis{
					{Label: "A", Min: -10, Max: 100},
					{Label: "B", Min: 0, Max: 50},
					{Label: "C", Min: 0, Max: 200},
				},
				Series: []RadarSeries{
					{Name: "S", Values: []float64{50, 25, 100}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.applyRadarDefaults()
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestRadarConstructor(t *testing.T) {
	v := Radar(RadarCfg{
		BaseCfg: BaseCfg{ID: "r1"},
		Axes: []RadarAxis{
			{Label: "A", Max: 100},
			{Label: "B", Max: 100},
			{Label: "C", Max: 100},
		},
		Series: []RadarSeries{
			{Name: "S", Values: []float64{50, 60, 70}},
		},
	})
	if v == nil {
		t.Fatal("Radar returned nil")
	}
	rv, ok := v.(*radarView)
	if !ok {
		t.Fatal("Radar did not return *radarView")
	}
	if !rv.cfg.ShowGrid {
		t.Error("default ShowGrid not set")
	}
	if !rv.cfg.ShowArea {
		t.Error("default ShowArea not set")
	}
	if !rv.cfg.ShowMarkers {
		t.Error("default ShowMarkers not set")
	}
	if rv.cfg.GridLevels != 5 {
		t.Errorf("default GridLevels = %d, want 5", rv.cfg.GridLevels)
	}
	if rv.cfg.StartAngle != float32(-math.Pi/2) {
		t.Errorf("default StartAngle = %g, want %g",
			rv.cfg.StartAngle, float32(-math.Pi/2))
	}
}

func TestRadarAutoMax(t *testing.T) {
	v := Radar(RadarCfg{
		BaseCfg: BaseCfg{ID: "r2"},
		Axes: []RadarAxis{
			{Label: "A"},
			{Label: "B"},
			{Label: "C"},
		},
		Series: []RadarSeries{
			{Name: "S1", Values: []float64{30, 80, 50}},
			{Name: "S2", Values: []float64{60, 40, 90}},
		},
	})
	rv := v.(*radarView)
	// Auto-max should pick the max across all series.
	wantMax := []float64{60, 80, 90}
	for i, want := range wantMax {
		if rv.cfg.Axes[i].Max != want {
			t.Errorf("axis %d Max = %g, want %g",
				i, rv.cfg.Axes[i].Max, want)
		}
	}
}

func TestRadarPolygonGridDefaults(t *testing.T) {
	// PolygonGrid: true must not prevent auto-enabling visual
	// features (ShowGrid, ShowArea, ShowMarkers).
	v := Radar(RadarCfg{
		BaseCfg:     BaseCfg{ID: "rpg"},
		PolygonGrid: true,
		Axes: []RadarAxis{
			{Label: "A", Max: 10},
			{Label: "B", Max: 10},
			{Label: "C", Max: 10},
		},
		Series: []RadarSeries{
			{Name: "S", Values: []float64{5, 6, 7}},
		},
	})
	rv := v.(*radarView)
	if !rv.cfg.ShowGrid {
		t.Error("ShowGrid should default to true with PolygonGrid")
	}
	if !rv.cfg.ShowArea {
		t.Error("ShowArea should default to true with PolygonGrid")
	}
	if !rv.cfg.ShowMarkers {
		t.Error("ShowMarkers should default to true with PolygonGrid")
	}
}

func TestRadarHoveredAxisIndex(t *testing.T) {
	// Build a radarView with cached geometry so hoveredAxisIndex
	// can compute angles from the center.
	cfg := RadarCfg{
		BaseCfg: BaseCfg{ID: "hai"},
		Axes: []RadarAxis{
			{Label: "A", Max: 100},
			{Label: "B", Max: 100},
			{Label: "C", Max: 100},
			{Label: "D", Max: 100},
			{Label: "E", Max: 100},
		},
		Series: []RadarSeries{
			{Name: "S", Values: []float64{50, 50, 50, 50, 50}},
		},
	}
	cfg.applyRadarDefaults()
	rv := &radarView{cfg: cfg, cx: 200, cy: 200, radius: 100}

	nAxes := len(cfg.Axes)
	for i := range nAxes {
		// Place the cursor on the spoke for axis i at radius/2.
		a := radarAxisAngle(cfg.StartAngle, i, nAxes)
		mx := rv.cx + 50*float32(math.Cos(float64(a)))
		my := rv.cy + 50*float32(math.Sin(float64(a)))
		got := rv.hoveredAxisIndex(mx, my)
		if got != i {
			t.Errorf("axis %d: hoveredAxisIndex at spoke = %d, want %d",
				i, got, i)
		}
	}
}

func TestRadarAxisAngle(t *testing.T) {
	tests := []struct {
		name       string
		startAngle float32
		i, nAxes   int
		want       float32
	}{
		{"first axis top", -math.Pi / 2, 0, 5, -math.Pi / 2},
		{"second of 4", 0, 1, 4, math.Pi / 2},
		{"third of 3", 0, 2, 3, 4 * math.Pi / 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := radarAxisAngle(tt.startAngle, tt.i, tt.nAxes)
			if math.Abs(float64(got-tt.want)) > 1e-5 {
				t.Errorf("radarAxisAngle(%g, %d, %d) = %g, want %g",
					tt.startAngle, tt.i, tt.nAxes, got, tt.want)
			}
		})
	}
}

func TestRadarAutoMaxInfIgnored(t *testing.T) {
	v := Radar(RadarCfg{
		BaseCfg: BaseCfg{ID: "rinf"},
		Axes: []RadarAxis{
			{Label: "A"},
			{Label: "B"},
			{Label: "C"},
		},
		Series: []RadarSeries{
			{Name: "S1", Values: []float64{
				30, math.Inf(1), 50,
			}},
			{Name: "S2", Values: []float64{
				60, 40, math.NaN(),
			}},
		},
	})
	rv := v.(*radarView)
	// Inf and NaN should be ignored in auto-max.
	wantMax := []float64{60, 40, 50}
	for i, want := range wantMax {
		if rv.cfg.Axes[i].Max != want {
			t.Errorf("axis %d Max = %g, want %g",
				i, rv.cfg.Axes[i].Max, want)
		}
	}
}

func TestRadarNormalize(t *testing.T) {
	tests := []struct {
		value, axisMin, axisMax float64
		want                    float64
	}{
		{50, 0, 100, 0.5},
		{0, 0, 100, 0},
		{100, 0, 100, 1},
		{-10, 0, 100, 0},          // clamped below
		{200, 0, 100, 1},          // clamped above
		{50, 50, 50, 0},           // degenerate
		{75, 50, 100, 0.5},        // offset range
		{math.NaN(), 0, 100, 0},   // NaN value
		{50, math.NaN(), 100, 0},  // NaN axisMin
		{50, 0, math.Inf(1), 0},   // Inf axisMax
		{math.Inf(-1), 0, 100, 0}, // -Inf value
	}
	for _, tt := range tests {
		got := radarNormalize(tt.value, tt.axisMin, tt.axisMax)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("radarNormalize(%g, %g, %g) = %g, want %g",
				tt.value, tt.axisMin, tt.axisMax, got, tt.want)
		}
	}
}
