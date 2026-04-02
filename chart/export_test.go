package chart

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/mike-ward/go-charts/axis"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-gui/gui"
)

// nonChartView is a minimal gui.View that does not implement
// Drawer, used to test ExportPNG rejection.
type nonChartView struct{}

func (nonChartView) Content() []gui.View                     { return nil }
func (nonChartView) GenerateLayout(_ *gui.Window) gui.Layout { return gui.Layout{} }

func TestExportPNG_Line(t *testing.T) {
	v := Line(LineCfg{
		BaseCfg: BaseCfg{
			ID:    "test-line",
			Width: 400, Height: 300,
		},
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:   "s1",
				Color:  gui.Hex(0x4E79A7),
				Points: []series.Point{{X: 0, Y: 1}, {X: 1, Y: 3}, {X: 2, Y: 2}},
			}),
		},
	})

	path := filepath.Join(t.TempDir(), "line.png")
	if err := ExportPNG(v, 400, 300, path); err != nil {
		t.Fatal(err)
	}
	assertValidPNG(t, path, 400, 300)
}

func TestExportPNG_Bar(t *testing.T) {
	v := Bar(BarCfg{
		BaseCfg: BaseCfg{
			ID:    "test-bar",
			Width: 400, Height: 300,
		},
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name:  "s1",
				Color: gui.Hex(0xE15759),
				Values: []series.CategoryValue{
					{Label: "A", Value: 10},
					{Label: "B", Value: 20},
					{Label: "C", Value: 15},
				},
			}),
		},
	})

	path := filepath.Join(t.TempDir(), "bar.png")
	if err := ExportPNG(v, 400, 300, path); err != nil {
		t.Fatal(err)
	}
	assertValidPNG(t, path, 400, 300)
}

func TestExportPNG_LineWithAxes(t *testing.T) {
	xa := axis.NewLinear(axis.LinearCfg{AutoRange: true})
	xa.SetRange(0, 10)
	ya := axis.NewLinear(axis.LinearCfg{AutoRange: true})
	ya.SetRange(0, 100)

	v := Line(LineCfg{
		BaseCfg: BaseCfg{
			ID:    "test-axes",
			Width: 800, Height: 600,
		},
		XAxis: xa,
		YAxis: ya,
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:   "data",
				Color:  gui.Hex(0x59A14F),
				Points: []series.Point{{X: 1, Y: 20}, {X: 5, Y: 80}, {X: 9, Y: 40}},
			}),
		},
		ShowMarkers: true,
	})

	path := filepath.Join(t.TempDir(), "axes.png")
	if err := ExportPNG(v, 800, 600, path); err != nil {
		t.Fatal(err)
	}
	assertValidPNG(t, path, 800, 600)
}

func TestExportPNG_LineAreaMonotonicX(t *testing.T) {
	// Monotonic X — area fill should render correctly.
	v := Line(LineCfg{
		BaseCfg: BaseCfg{
			ID:    "test-area-monotonic",
			Width: 400, Height: 300,
		},
		ShowArea: true,
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "s1",
				Color: gui.Hex(0x59A14F),
				Points: []series.Point{
					{X: 1, Y: 10}, {X: 2, Y: 30}, {X: 3, Y: 20},
					{X: 4, Y: 50}, {X: 5, Y: 40}, {X: 6, Y: 60},
				},
			}),
		},
	})

	path := filepath.Join(t.TempDir(), "area_monotonic.png")
	if err := ExportPNG(v, 400, 300, path); err != nil {
		t.Fatal(err)
	}
	assertValidPNG(t, path, 400, 300)
}

func TestExportPNG_LineAreaNonMonotonicX(t *testing.T) {
	// Non-monotonic X — fill trapezoids follow source point order,
	// same as the rendered polyline. Export must succeed.
	v := Line(LineCfg{
		BaseCfg: BaseCfg{
			ID:    "test-area-nonmonotonic",
			Width: 400, Height: 300,
		},
		ShowArea: true,
		Series: []series.XY{
			series.NewXY(series.XYCfg{
				Name:  "s1",
				Color: gui.Hex(0x59A14F),
				Points: []series.Point{
					{X: 0, Y: 0}, {X: 10, Y: 10},
					{X: 5, Y: 5}, {X: 15, Y: 15},
				},
			}),
		},
	})

	path := filepath.Join(t.TempDir(), "area_nonmonotonic.png")
	if err := ExportPNG(v, 400, 300, path); err != nil {
		t.Fatal(err)
	}
	assertValidPNG(t, path, 400, 300)
}

func TestExportPNG_RejectsNonChartView(t *testing.T) {
	v := nonChartView{}
	err := ExportPNG(v, 100, 100, "/dev/null")
	if err == nil {
		t.Fatal("expected error for non-chart view")
	}
}

func TestExportPNG_RejectsZeroDimensions(t *testing.T) {
	v := Line(LineCfg{
		BaseCfg: BaseCfg{Width: 100, Height: 100},
	})
	err := ExportPNG(v, 0, 100, "/dev/null")
	if err == nil {
		t.Fatal("expected error for zero width")
	}
}

func assertValidPNG(t *testing.T, path string, wantW, wantH int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatal("invalid PNG:", err)
	}
	b := img.Bounds()
	if b.Dx() != wantW || b.Dy() != wantH {
		t.Fatalf("dimensions: got %dx%d, want %dx%d",
			b.Dx(), b.Dy(), wantW, wantH)
	}
}
