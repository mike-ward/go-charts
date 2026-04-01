package series

import (
	"math"
	"testing"
)

func TestBoundsNaNInf(t *testing.T) {
	nan := math.NaN()
	inf := math.Inf(1)
	ninf := math.Inf(-1)

	tests := []struct {
		name               string
		points             []Point
		wantMinX, wantMaxX float64
		wantMinY, wantMaxY float64
		wantZero           bool // all returned values should be 0
	}{
		{
			name:     "all NaN",
			points:   []Point{{nan, nan}, {nan, nan}},
			wantZero: true,
		},
		{
			name:     "all +Inf",
			points:   []Point{{inf, inf}, {inf, inf}},
			wantZero: true,
		},
		{
			name:     "all -Inf",
			points:   []Point{{ninf, ninf}},
			wantZero: true,
		},
		{
			name:     "mixed Inf",
			points:   []Point{{inf, ninf}, {ninf, inf}},
			wantZero: true,
		},
		{
			name:     "NaN X valid Y",
			points:   []Point{{nan, 5}},
			wantZero: true,
		},
		{
			name:     "valid X NaN Y",
			points:   []Point{{5, nan}},
			wantZero: true,
		},
		{
			name:     "valid among NaN",
			points:   []Point{{nan, nan}, {1, 2}, {nan, nan}, {3, 4}},
			wantMinX: 1, wantMaxX: 3,
			wantMinY: 2, wantMaxY: 4,
		},
		{
			name:     "valid among Inf",
			points:   []Point{{inf, inf}, {10, 20}, {ninf, ninf}, {30, 40}},
			wantMinX: 10, wantMaxX: 30,
			wantMinY: 20, wantMaxY: 40,
		},
		{
			name:     "single valid point",
			points:   []Point{{7, 11}},
			wantMinX: 7, wantMaxX: 7,
			wantMinY: 11, wantMaxY: 11,
		},
		{
			name:     "two identical points",
			points:   []Point{{5, 5}, {5, 5}},
			wantMinX: 5, wantMaxX: 5,
			wantMinY: 5, wantMaxY: 5,
		},
		{
			name:     "very large values",
			points:   []Point{{1e308, -1e308}, {-1e308, 1e308}},
			wantMinX: -1e308, wantMaxX: 1e308,
			wantMinY: -1e308, wantMaxY: 1e308,
		},
		{
			name:     "very small values",
			points:   []Point{{1e-308, 2e-308}, {3e-308, 4e-308}},
			wantMinX: 1e-308, wantMaxX: 3e-308,
			wantMinY: 2e-308, wantMaxY: 4e-308,
		},
		{
			name:     "negative values only",
			points:   []Point{{-10, -20}, {-5, -30}},
			wantMinX: -10, wantMaxX: -5,
			wantMinY: -30, wantMaxY: -20,
		},
		{
			name:     "empty",
			points:   nil,
			wantZero: true,
		},
		{
			name:     "first valid last NaN",
			points:   []Point{{1, 2}, {nan, nan}},
			wantMinX: 1, wantMaxX: 1,
			wantMinY: 2, wantMaxY: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewXY(XYCfg{Points: tt.points})
			minX, maxX, minY, maxY := s.Bounds()

			if tt.wantZero {
				if minX != 0 || maxX != 0 || minY != 0 || maxY != 0 {
					t.Errorf("got (%v,%v,%v,%v), want all zeros",
						minX, maxX, minY, maxY)
				}
				return
			}

			if minX != tt.wantMinX || maxX != tt.wantMaxX ||
				minY != tt.wantMinY || maxY != tt.wantMaxY {
				t.Errorf("got (%v,%v,%v,%v), want (%v,%v,%v,%v)",
					minX, maxX, minY, maxY,
					tt.wantMinX, tt.wantMaxX,
					tt.wantMinY, tt.wantMaxY)
			}
		})
	}
}
