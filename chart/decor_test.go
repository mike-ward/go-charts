package chart

import (
	"testing"

	"github.com/mike-ward/go-charts/render"
	"github.com/mike-ward/go-charts/series"
	"github.com/mike-ward/go-charts/theme"
	"github.com/mike-ward/go-gui/gui"
)

func testCtx(w, h float32) (*render.Context, *gui.DrawContext) {
	dc := &gui.DrawContext{Width: w, Height: h}
	return render.NewContext(dc), dc
}

func testTheme() *theme.Theme {
	return &theme.Theme{
		TitleStyle:    gui.TextStyle{Size: 14, Color: gui.White},
		LabelStyle:    gui.TextStyle{Size: 11, Color: gui.White},
		TickStyle:     gui.TextStyle{Size: 11, Color: gui.White},
		AxisColor:     gui.White,
		AxisWidth:     1,
		GridColor:     gui.RGBA(128, 128, 128, 40),
		GridWidth:     0.5,
		Palette:       theme.DefaultPalette(),
		PaddingTop:    40,
		PaddingRight:  20,
		PaddingBottom: 40,
		PaddingLeft:   60,
	}
}

// --- drawTitle ---

func TestDrawTitleRendersText(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	drawTitle(ctx, "My Chart", th)
	if len(dc.Texts()) != 1 {
		t.Fatalf("texts = %d, want 1", len(dc.Texts()))
	}
	if dc.Texts()[0].Text != "My Chart" {
		t.Errorf("text = %q, want %q",
			dc.Texts()[0].Text, "My Chart")
	}
}

func TestDrawTitleEmptySkipped(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	drawTitle(ctx, "", th)
	if len(dc.Texts()) != 0 {
		t.Errorf("empty title should produce no text, got %d",
			len(dc.Texts()))
	}
}

// --- drawLegend ---

func TestDrawLegendRendersEntries(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	entries := []legendEntry{
		{Name: "Series A", Color: gui.Hex(0xFF0000)},
		{Name: "Series B", Color: gui.Hex(0x00FF00)},
	}
	drawLegend(ctx, entries, th, 380, 40)
	// Background rect + 2 swatches = 3 rounded rects → batches.
	if len(dc.Batches()) == 0 {
		t.Error("expected batches for legend background/swatches")
	}
	// 2 text labels.
	if len(dc.Texts()) != 2 {
		t.Errorf("texts = %d, want 2", len(dc.Texts()))
	}
}

func TestDrawLegendSkipsUnnamedEntries(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	entries := []legendEntry{
		{Name: "", Color: gui.Hex(0xFF0000)},
		{Name: "Visible", Color: gui.Hex(0x00FF00)},
	}
	drawLegend(ctx, entries, th, 380, 40)
	// Only 1 named entry → 1 text label.
	if len(dc.Texts()) != 1 {
		t.Errorf("texts = %d, want 1", len(dc.Texts()))
	}
}

func TestDrawLegendAllUnnamedSkipped(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	entries := []legendEntry{
		{Name: "", Color: gui.Hex(0xFF0000)},
		{Name: "", Color: gui.Hex(0x00FF00)},
	}
	drawLegend(ctx, entries, th, 380, 40)
	if len(dc.Texts()) != 0 {
		t.Errorf("all-unnamed should produce no text, got %d",
			len(dc.Texts()))
	}
	if len(dc.Batches()) != 0 {
		t.Errorf("all-unnamed should produce no batches, got %d",
			len(dc.Batches()))
	}
}

func TestDrawLegendEmptyEntries(t *testing.T) {
	ctx, dc := testCtx(400, 300)
	th := testTheme()
	drawLegend(ctx, nil, th, 380, 40)
	if len(dc.Texts()) != 0 || len(dc.Batches()) != 0 {
		t.Error("nil entries should produce no output")
	}
}

// --- Integration: bar chart title + legend ---

func TestBarDrawTitleAndLegend(t *testing.T) {
	bv := Bar(BarCfg{
		BaseCfg: BaseCfg{
			ID:    "test-bar",
			Title: "Bar Title",
		},
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Name: "S1",
				Values: []series.CategoryValue{
					{Label: "A", Value: 10},
					{Label: "B", Value: 20},
				},
			}),
			series.NewCategory(series.CategoryCfg{
				Name: "S2",
				Values: []series.CategoryValue{
					{Label: "A", Value: 15},
					{Label: "B", Value: 25},
				},
			}),
		},
	}).(*barView)

	dc := gui.DrawContext{Width: 400, Height: 300}
	bv.draw(&dc)

	// Should have title text + Y tick labels + X tick labels +
	// legend labels.
	hasTitle := false
	hasLegendS1 := false
	hasLegendS2 := false
	for _, txt := range dc.Texts() {
		switch txt.Text {
		case "Bar Title":
			hasTitle = true
		case "S1":
			hasLegendS1 = true
		case "S2":
			hasLegendS2 = true
		}
	}
	if !hasTitle {
		t.Error("title not rendered")
	}
	if !hasLegendS1 {
		t.Error("legend entry S1 not rendered")
	}
	if !hasLegendS2 {
		t.Error("legend entry S2 not rendered")
	}
}

// --- Integration: line chart title + legend ---

func TestLineDrawTitleAndLegend(t *testing.T) {
	lv := Line(LineCfg{
		BaseCfg: BaseCfg{
			ID:    "test-line",
			Title: "Line Title",
		},
		Series: []series.XY{
			series.XYFromYValues("Alpha", []float64{1, 2, 3}),
			series.XYFromYValues("Beta", []float64{3, 2, 1}),
		},
	}).(*lineView)

	dc := gui.DrawContext{Width: 400, Height: 300}
	lv.draw(&dc)

	hasTitle := false
	hasAlpha := false
	hasBeta := false
	for _, txt := range dc.Texts() {
		switch txt.Text {
		case "Line Title":
			hasTitle = true
		case "Alpha":
			hasAlpha = true
		case "Beta":
			hasBeta = true
		}
	}
	if !hasTitle {
		t.Error("title not rendered")
	}
	if !hasAlpha {
		t.Error("legend entry Alpha not rendered")
	}
	if !hasBeta {
		t.Error("legend entry Beta not rendered")
	}
}

// --- No title/legend when absent ---

func TestBarDrawNoTitleNoLegend(t *testing.T) {
	bv := Bar(BarCfg{
		BaseCfg: BaseCfg{ID: "test-bar-bare"},
		Series: []series.Category{
			series.NewCategory(series.CategoryCfg{
				Values: []series.CategoryValue{
					{Label: "A", Value: 10},
				},
			}),
		},
	}).(*barView)

	dc := gui.DrawContext{Width: 400, Height: 300}
	bv.draw(&dc)

	for _, txt := range dc.Texts() {
		// Should only have tick labels, no title or legend text.
		if txt.Text == "" {
			t.Error("empty text rendered")
		}
	}
	// No legend entries since series has no name.
	for _, txt := range dc.Texts() {
		if txt.Text == "test-bar-bare" {
			t.Error("ID should not appear as title")
		}
	}
}
