package series

import (
	"strings"
	"testing"
	"time"
)

func TestXYFromCSV_HeaderNames(t *testing.T) {
	data := "x,y\n1,10\n2,20\n3,30\n"
	s, err := XYFromCSV(strings.NewReader(data), "test",
		ColName("x"), ColName("y"), CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name() != "test" {
		t.Errorf("name = %q, want %q", s.Name(), "test")
	}
	if s.Len() != 3 {
		t.Fatalf("len = %d, want 3", s.Len())
	}
	if s.Points[0].X != 1 || s.Points[0].Y != 10 {
		t.Errorf("point[0] = %v, want (1, 10)", s.Points[0])
	}
	if s.Points[2].X != 3 || s.Points[2].Y != 30 {
		t.Errorf("point[2] = %v, want (3, 30)", s.Points[2])
	}
}

func TestXYFromCSV_ColIdx(t *testing.T) {
	data := "1,10\n2,20\n"
	s, err := XYFromCSV(strings.NewReader(data), "idx",
		ColIdx(0), ColIdx(1), CSVCfg{})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[1].X != 2 || s.Points[1].Y != 20 {
		t.Errorf("point[1] = %v, want (2, 20)", s.Points[1])
	}
}

func TestXYFromCSV_CustomDelimiter(t *testing.T) {
	data := "x\ty\n1\t10\n2\t20\n"
	s, err := XYFromCSV(strings.NewReader(data), "tab",
		ColName("x"), ColName("y"),
		CSVCfg{HasHeader: true, Delimiter: '\t'})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
}

func TestXYFromCSV_TrimSpace(t *testing.T) {
	data := " x , y \n 1 , 10 \n 2 , 20 \n"
	s, err := XYFromCSV(strings.NewReader(data), "trim",
		ColName("x"), ColName("y"),
		CSVCfg{HasHeader: true, TrimSpace: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[0].X != 1 || s.Points[0].Y != 10 {
		t.Errorf("point[0] = %v, want (1, 10)", s.Points[0])
	}
}

func TestXYFromCSV_SkipErrors(t *testing.T) {
	data := "x,y\n1,10\nbad,20\n3,30\n"
	s, err := XYFromCSV(strings.NewReader(data), "skip",
		ColName("x"), ColName("y"),
		CSVCfg{HasHeader: true, SkipErrors: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[0].X != 1 || s.Points[1].X != 3 {
		t.Errorf("points = %v, want x=1,3", s.Points)
	}
}

func TestXYFromCSV_ErrorMissingCol(t *testing.T) {
	data := "a,b\n1,2\n"
	_, err := XYFromCSV(strings.NewReader(data), "err",
		ColName("x"), ColName("y"), CSVCfg{HasHeader: true})
	if err == nil {
		t.Fatal("expected error for missing column")
	}
}

func TestXYFromCSV_Comment(t *testing.T) {
	data := "# comment\n1,10\n2,20\n"
	s, err := XYFromCSV(strings.NewReader(data), "cmt",
		ColIdx(0), ColIdx(1), CSVCfg{Comment: '#'})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
}

func TestXYFromCSV_Empty(t *testing.T) {
	data := "x,y\n"
	s, err := XYFromCSV(strings.NewReader(data), "empty",
		ColName("x"), ColName("y"), CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 0 {
		t.Fatalf("len = %d, want 0", s.Len())
	}
}

func TestXYMultiFromCSV(t *testing.T) {
	data := "x,a,b\n1,10,100\n2,20,200\n"
	ss, err := XYMultiFromCSV(strings.NewReader(data),
		ColName("x"),
		[]Col{ColName("a"), ColName("b")},
		nil,
		CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 2 {
		t.Fatalf("len = %d, want 2", len(ss))
	}
	if ss[0].Name() != "a" || ss[1].Name() != "b" {
		t.Errorf("names = %q, %q, want a, b",
			ss[0].Name(), ss[1].Name())
	}
	if ss[0].Len() != 2 || ss[1].Len() != 2 {
		t.Fatalf("lens = %d, %d, want 2, 2",
			ss[0].Len(), ss[1].Len())
	}
	if ss[1].Points[1].Y != 200 {
		t.Errorf("ss[1].Points[1].Y = %v, want 200",
			ss[1].Points[1].Y)
	}
}

func TestXYMultiFromCSV_SkipErrorsAligned(t *testing.T) {
	// Row 2 has "bad" in the second Y column. With SkipErrors,
	// the entire row must be skipped to keep series aligned.
	data := "x,a,b\n1,10,100\n2,20,bad\n3,30,300\n"
	ss, err := XYMultiFromCSV(strings.NewReader(data),
		ColName("x"),
		[]Col{ColName("a"), ColName("b")},
		nil,
		CSVCfg{HasHeader: true, SkipErrors: true})
	if err != nil {
		t.Fatal(err)
	}
	if ss[0].Len() != ss[1].Len() {
		t.Fatalf("series lengths differ: %d vs %d",
			ss[0].Len(), ss[1].Len())
	}
	if ss[0].Len() != 2 {
		t.Fatalf("len = %d, want 2", ss[0].Len())
	}
	if ss[0].Points[1].X != 3 {
		t.Errorf("ss[0].Points[1].X = %v, want 3",
			ss[0].Points[1].X)
	}
}

func TestXYMultiFromCSV_ExplicitNames(t *testing.T) {
	data := "1,10,100\n2,20,200\n"
	ss, err := XYMultiFromCSV(strings.NewReader(data),
		ColIdx(0),
		[]Col{ColIdx(1), ColIdx(2)},
		[]string{"Series A", "Series B"},
		CSVCfg{})
	if err != nil {
		t.Fatal(err)
	}
	if ss[0].Name() != "Series A" {
		t.Errorf("name = %q, want %q", ss[0].Name(), "Series A")
	}
}

func TestCategoryFromCSV(t *testing.T) {
	data := "label,value\nNorth,45\nSouth,32\n"
	s, err := CategoryFromCSV(strings.NewReader(data), "regions",
		ColName("label"), ColName("value"),
		CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name() != "regions" {
		t.Errorf("name = %q, want %q", s.Name(), "regions")
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Values[0].Label != "North" || s.Values[0].Value != 45 {
		t.Errorf("values[0] = %v, want North:45", s.Values[0])
	}
}

func TestXYZFromCSV(t *testing.T) {
	data := "x,y,z\n1,2,3\n4,5,6\n"
	s, err := XYZFromCSV(strings.NewReader(data), "bubble",
		ColName("x"), ColName("y"), ColName("z"),
		CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[0].Z != 3 {
		t.Errorf("points[0].Z = %v, want 3", s.Points[0].Z)
	}
}

func TestOHLCFromCSV(t *testing.T) {
	data := "time,open,high,low,close,volume\n" +
		"2024-01-01,100,110,95,105,1000\n" +
		"2024-01-02,105,115,100,112,1200\n"
	s, err := OHLCFromCSV(strings.NewReader(data), "stock",
		OHLCCSVCfg{
			CSVCfg:     CSVCfg{HasHeader: true},
			TimeCol:    ColName("time"),
			OpenCol:    ColName("open"),
			HighCol:    ColName("high"),
			LowCol:     ColName("low"),
			CloseCol:   ColName("close"),
			VolumeCol:  ColName("volume"),
			TimeLayout: "2006-01-02",
		})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[0].Open != 100 {
		t.Errorf("points[0].Open = %v, want 100",
			s.Points[0].Open)
	}
	if s.Points[1].Volume != 1200 {
		t.Errorf("points[1].Volume = %v, want 1200",
			s.Points[1].Volume)
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !s.Points[0].Time.Equal(expected) {
		t.Errorf("points[0].Time = %v, want %v",
			s.Points[0].Time, expected)
	}
}

func TestOHLCFromCSV_NoVolume(t *testing.T) {
	data := "time,open,high,low,close\n" +
		"2024-01-01,100,110,95,105\n"
	s, err := OHLCFromCSV(strings.NewReader(data), "noVol",
		OHLCCSVCfg{
			CSVCfg:     CSVCfg{HasHeader: true},
			TimeCol:    ColName("time"),
			OpenCol:    ColName("open"),
			HighCol:    ColName("high"),
			LowCol:     ColName("low"),
			CloseCol:   ColName("close"),
			TimeLayout: "2006-01-02",
		})
	if err != nil {
		t.Fatal(err)
	}
	if s.Points[0].Volume != 0 {
		t.Errorf("volume = %v, want 0", s.Points[0].Volume)
	}
}

func TestGridFromCSV(t *testing.T) {
	data := ",Mon,Tue,Wed\nAlice,1,2,3\nBob,4,5,6\n"
	g, err := GridFromCSV(strings.NewReader(data), "grid",
		CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if g.NumRows() != 2 || g.NumCols() != 3 {
		t.Fatalf("dims = %dx%d, want 2x3",
			g.NumRows(), g.NumCols())
	}
	if g.Rows()[0] != "Alice" || g.Rows()[1] != "Bob" {
		t.Errorf("rows = %v", g.Rows())
	}
	if g.Cols()[0] != "Mon" {
		t.Errorf("cols[0] = %q, want Mon", g.Cols()[0])
	}
	if g.At(1, 2) != 6 {
		t.Errorf("At(1,2) = %v, want 6", g.At(1, 2))
	}
}

func TestGridFromCSV_NoHeader(t *testing.T) {
	data := "Alice,1,2\n"
	_, err := GridFromCSV(strings.NewReader(data), "err",
		CSVCfg{HasHeader: false})
	if err == nil {
		t.Fatal("expected error when HasHeader is false")
	}
}

func TestXYFromCSV_ColIdxOutOfRange(t *testing.T) {
	data := "1,2\n"
	_, err := XYFromCSV(strings.NewReader(data), "x",
		ColIdx(0), ColIdx(99), CSVCfg{})
	if err == nil {
		t.Fatal("expected error for out-of-range ColIdx")
	}
}

func TestXYFromCSV_SingleRow(t *testing.T) {
	data := "x,y\n42,99\n"
	s, err := XYFromCSV(strings.NewReader(data), "one",
		ColName("x"), ColName("y"), CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 1 {
		t.Fatalf("len = %d, want 1", s.Len())
	}
	if s.Points[0].X != 42 || s.Points[0].Y != 99 {
		t.Errorf("point = %v, want (42, 99)", s.Points[0])
	}
}

func TestXYFromCSV_ZeroCol(t *testing.T) {
	data := "1,2\n"
	_, err := XYFromCSV(strings.NewReader(data), "x",
		Col{}, ColIdx(0), CSVCfg{})
	if err == nil {
		t.Fatal("expected error for zero-value Col")
	}
}

func TestXYMultiFromCSV_EmptyYCols(t *testing.T) {
	data := "x\n1\n2\n"
	ss, err := XYMultiFromCSV(strings.NewReader(data),
		ColName("x"), nil, nil, CSVCfg{HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 0 {
		t.Fatalf("len = %d, want 0", len(ss))
	}
}

func TestCategoryFromCSV_SkipErrors(t *testing.T) {
	data := "label,value\nA,10\nB,bad\nC,30\n"
	s, err := CategoryFromCSV(strings.NewReader(data), "cat",
		ColName("label"), ColName("value"),
		CSVCfg{HasHeader: true, SkipErrors: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Values[1].Label != "C" {
		t.Errorf("values[1].Label = %q, want C", s.Values[1].Label)
	}
}

func TestOHLCFromCSV_DefaultTimeLayout(t *testing.T) {
	data := "time,open,high,low,close\n" +
		"2024-01-01T00:00:00Z,100,110,95,105\n"
	s, err := OHLCFromCSV(strings.NewReader(data), "rfc",
		OHLCCSVCfg{
			CSVCfg:   CSVCfg{HasHeader: true},
			TimeCol:  ColName("time"),
			OpenCol:  ColName("open"),
			HighCol:  ColName("high"),
			LowCol:   ColName("low"),
			CloseCol: ColName("close"),
		})
	if err != nil {
		t.Fatal(err)
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !s.Points[0].Time.Equal(expected) {
		t.Errorf("time = %v, want %v",
			s.Points[0].Time, expected)
	}
}

func TestOHLCFromCSV_SkipBadTime(t *testing.T) {
	data := "time,open,high,low,close\n" +
		"2024-01-01,100,110,95,105\n" +
		"not-a-date,100,110,95,105\n" +
		"2024-01-03,100,110,95,105\n"
	s, err := OHLCFromCSV(strings.NewReader(data), "skip",
		OHLCCSVCfg{
			CSVCfg:     CSVCfg{HasHeader: true, SkipErrors: true},
			TimeCol:    ColName("time"),
			OpenCol:    ColName("open"),
			HighCol:    ColName("high"),
			LowCol:     ColName("low"),
			CloseCol:   ColName("close"),
			TimeLayout: "2006-01-02",
		})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
}

func TestGridFromCSV_SingleColumn(t *testing.T) {
	data := "only\nAlice\n"
	_, err := GridFromCSV(strings.NewReader(data), "err",
		CSVCfg{HasHeader: true})
	if err == nil {
		t.Fatal("expected error for single-column header")
	}
}

func TestGridFromCSV_ShortRow(t *testing.T) {
	data := ",Mon,Tue\nAlice,1\n"
	_, err := GridFromCSV(strings.NewReader(data), "err",
		CSVCfg{HasHeader: true})
	if err == nil {
		t.Fatal("expected error for short row")
	}
}

func TestCSV_NilReader(t *testing.T) {
	_, err := XYFromCSV(nil, "x", ColIdx(0), ColIdx(1), CSVCfg{})
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestCSV_NegativeColIdx(t *testing.T) {
	data := "1,2\n"
	_, err := XYFromCSV(strings.NewReader(data), "x",
		ColIdx(-1), ColIdx(0), CSVCfg{})
	if err == nil {
		t.Fatal("expected error for negative ColIdx")
	}
}

func TestCSV_NaNInfRejected(t *testing.T) {
	for _, val := range []string{"NaN", "Inf", "-Inf", "+Inf"} {
		data := "1," + val + "\n"
		_, err := XYFromCSV(strings.NewReader(data), "x",
			ColIdx(0), ColIdx(1), CSVCfg{})
		if err == nil {
			t.Errorf("expected error for %q", val)
		}
	}
}

func TestCSV_NaNInfSkipped(t *testing.T) {
	data := "1,10\n2,NaN\n3,30\n"
	s, err := XYFromCSV(strings.NewReader(data), "x",
		ColIdx(0), ColIdx(1),
		CSVCfg{SkipErrors: true})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
}
