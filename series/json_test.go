package series

import (
	"strings"
	"testing"
	"time"
)

func TestXYFromJSON(t *testing.T) {
	data := `[{"x": 1, "y": 10}, {"x": 2, "y": 20}]`
	s, err := XYFromJSON(strings.NewReader(data), "test", "x", "y")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name() != "test" {
		t.Errorf("name = %q, want %q", s.Name(), "test")
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Points[0].X != 1 || s.Points[0].Y != 10 {
		t.Errorf("point[0] = %v, want (1, 10)", s.Points[0])
	}
}

func TestXYFromJSON_IntValues(t *testing.T) {
	data := `[{"x": 1, "y": 2}]`
	s, err := XYFromJSON(strings.NewReader(data), "int", "x", "y")
	if err != nil {
		t.Fatal(err)
	}
	if s.Points[0].X != 1 || s.Points[0].Y != 2 {
		t.Errorf("point[0] = %v, want (1, 2)", s.Points[0])
	}
}

func TestXYFromJSON_EmptyArray(t *testing.T) {
	data := `[]`
	s, err := XYFromJSON(strings.NewReader(data), "empty", "x", "y")
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 0 {
		t.Fatalf("len = %d, want 0", s.Len())
	}
}

func TestXYFromJSON_MissingField(t *testing.T) {
	data := `[{"x": 1}]`
	_, err := XYFromJSON(strings.NewReader(data), "err", "x", "y")
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestXYFromJSON_MalformedInput(t *testing.T) {
	data := `not json`
	_, err := XYFromJSON(strings.NewReader(data), "err", "x", "y")
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestXYMultiFromJSON(t *testing.T) {
	data := `[
		{"x": 1, "a": 10, "b": 100},
		{"x": 2, "a": 20, "b": 200}
	]`
	ss, err := XYMultiFromJSON(strings.NewReader(data),
		"x", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 2 {
		t.Fatalf("len = %d, want 2", len(ss))
	}
	if ss[0].Name() != "a" || ss[1].Name() != "b" {
		t.Errorf("names = %q, %q", ss[0].Name(), ss[1].Name())
	}
	if ss[1].Points[1].Y != 200 {
		t.Errorf("ss[1].Points[1].Y = %v, want 200",
			ss[1].Points[1].Y)
	}
}

func TestCategoryFromJSON(t *testing.T) {
	data := `[
		{"region": "North", "sales": 45},
		{"region": "South", "sales": 32}
	]`
	s, err := CategoryFromJSON(strings.NewReader(data),
		"regions", "region", "sales")
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("len = %d, want 2", s.Len())
	}
	if s.Values[0].Label != "North" || s.Values[0].Value != 45 {
		t.Errorf("values[0] = %v", s.Values[0])
	}
}

func TestXYZFromJSON(t *testing.T) {
	data := `[{"x": 1, "y": 2, "z": 3}]`
	s, err := XYZFromJSON(strings.NewReader(data),
		"bubble", "x", "y", "z")
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 1 {
		t.Fatalf("len = %d, want 1", s.Len())
	}
	if s.Points[0].Z != 3 {
		t.Errorf("z = %v, want 3", s.Points[0].Z)
	}
}

func TestOHLCFromJSON(t *testing.T) {
	data := `[
		{"time": "2024-01-01", "open": 100, "high": 110,
		 "low": 95, "close": 105, "vol": 1000}
	]`
	s, err := OHLCFromJSON(strings.NewReader(data), "stock",
		OHLCJSONCfg{
			TimeLayout:  "2006-01-02",
			VolumeField: "vol",
		})
	if err != nil {
		t.Fatal(err)
	}
	if s.Len() != 1 {
		t.Fatalf("len = %d, want 1", s.Len())
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !s.Points[0].Time.Equal(expected) {
		t.Errorf("time = %v, want %v",
			s.Points[0].Time, expected)
	}
	if s.Points[0].Open != 100 {
		t.Errorf("open = %v, want 100", s.Points[0].Open)
	}
	if s.Points[0].Volume != 1000 {
		t.Errorf("volume = %v, want 1000", s.Points[0].Volume)
	}
}

func TestOHLCFromJSON_DefaultFields(t *testing.T) {
	data := `[
		{"time": "2024-01-01T00:00:00Z", "open": 100,
		 "high": 110, "low": 95, "close": 105}
	]`
	s, err := OHLCFromJSON(strings.NewReader(data), "def",
		OHLCJSONCfg{})
	if err != nil {
		t.Fatal(err)
	}
	if s.Points[0].Close != 105 {
		t.Errorf("close = %v, want 105", s.Points[0].Close)
	}
}

func TestGridFromJSON(t *testing.T) {
	data := `{
		"rows": ["Alice", "Bob"],
		"cols": ["Mon", "Tue"],
		"values": [[1, 2], [3, 4]]
	}`
	g, err := GridFromJSON(strings.NewReader(data), "grid")
	if err != nil {
		t.Fatal(err)
	}
	if g.NumRows() != 2 || g.NumCols() != 2 {
		t.Fatalf("dims = %dx%d, want 2x2",
			g.NumRows(), g.NumCols())
	}
	if g.At(1, 1) != 4 {
		t.Errorf("At(1,1) = %v, want 4", g.At(1, 1))
	}
}

func TestGridFromJSON_DimensionMismatch(t *testing.T) {
	data := `{
		"rows": ["Alice"],
		"cols": ["Mon", "Tue"],
		"values": [[1, 2], [3, 4]]
	}`
	_, err := GridFromJSON(strings.NewReader(data), "err")
	if err == nil {
		t.Fatal("expected error for dimension mismatch")
	}
}

func TestTreeNodeFromJSON(t *testing.T) {
	data := `{
		"label": "root",
		"value": 0,
		"children": [
			{"label": "A", "value": 10, "children": []},
			{"label": "B", "value": 0, "children": [
				{"label": "B1", "value": 5, "children": []},
				{"label": "B2", "value": 3, "children": []}
			]}
		]
	}`
	n, err := TreeNodeFromJSON(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if n.Label != "root" {
		t.Errorf("label = %q, want root", n.Label)
	}
	if len(n.Children) != 2 {
		t.Fatalf("children = %d, want 2", len(n.Children))
	}
	if n.Children[0].Label != "A" || n.Children[0].Value != 10 {
		t.Errorf("child[0] = %v", n.Children[0])
	}
	if n.TotalValue() != 18 {
		t.Errorf("totalValue = %v, want 18", n.TotalValue())
	}
	if len(n.Children[1].Children) != 2 {
		t.Errorf("B children = %d, want 2",
			len(n.Children[1].Children))
	}
}

func TestTreeNodeFromJSON_Leaf(t *testing.T) {
	data := `{"label": "leaf", "value": 42}`
	n, err := TreeNodeFromJSON(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if !n.IsLeaf() {
		t.Error("expected leaf node")
	}
	if n.Value != 42 {
		t.Errorf("value = %v, want 42", n.Value)
	}
}

func TestXYFromJSON_WrongType(t *testing.T) {
	data := `[{"x": "not a number", "y": 10}]`
	_, err := XYFromJSON(strings.NewReader(data), "err", "x", "y")
	if err == nil {
		t.Fatal("expected error for string in number field")
	}
}

func TestCategoryFromJSON_WrongLabelType(t *testing.T) {
	data := `[{"label": 123, "value": 10}]`
	_, err := CategoryFromJSON(strings.NewReader(data),
		"err", "label", "value")
	if err == nil {
		t.Fatal("expected error for number in string field")
	}
}

func TestOHLCFromJSON_BadTime(t *testing.T) {
	data := `[{"time": "not-a-date", "open": 100,
		"high": 110, "low": 95, "close": 105}]`
	_, err := OHLCFromJSON(strings.NewReader(data), "err",
		OHLCJSONCfg{TimeLayout: "2006-01-02"})
	if err == nil {
		t.Fatal("expected error for bad time string")
	}
}

func TestXYMultiFromJSON_EmptyFields(t *testing.T) {
	data := `[{"x": 1}]`
	ss, err := XYMultiFromJSON(strings.NewReader(data),
		"x", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 0 {
		t.Fatalf("len = %d, want 0", len(ss))
	}
}

func TestGridFromJSON_EmptyObject(t *testing.T) {
	data := `{}`
	g, err := GridFromJSON(strings.NewReader(data), "empty")
	if err != nil {
		t.Fatal(err)
	}
	if g.NumRows() != 0 || g.NumCols() != 0 {
		t.Errorf("dims = %dx%d, want 0x0",
			g.NumRows(), g.NumCols())
	}
}

func TestTreeNodeFromJSON_Malformed(t *testing.T) {
	data := `{not json`
	_, err := TreeNodeFromJSON(strings.NewReader(data))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestJSON_NilReader(t *testing.T) {
	_, err := XYFromJSON(nil, "x", "x", "y")
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
	_, err = GridFromJSON(nil, "g")
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
	_, err = TreeNodeFromJSON(nil)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}
