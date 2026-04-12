package series

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"time"
)

// maxPrealloc caps pre-allocated slice capacity to avoid huge
// allocations from a single decoded JSON array.
const maxPrealloc = 1 << 16 // 65536

// OHLCJSONCfg maps JSON object fields to OHLC values.
// Empty field names default to standard names (time, open,
// high, low, close).
type OHLCJSONCfg struct {
	TimeField   string // default "time"
	OpenField   string // default "open"
	HighField   string // default "high"
	LowField    string // default "low"
	CloseField  string // default "close"
	VolumeField string // empty means no volume
	TimeLayout  string // Go time layout; empty = time.RFC3339
}

func (c *OHLCJSONCfg) applyDefaults() {
	c.TimeField = cmp.Or(c.TimeField, "time")
	c.OpenField = cmp.Or(c.OpenField, "open")
	c.HighField = cmp.Or(c.HighField, "high")
	c.LowField = cmp.Or(c.LowField, "low")
	c.CloseField = cmp.Or(c.CloseField, "close")
	c.TimeLayout = cmp.Or(c.TimeLayout, time.RFC3339)
}

func jsonFloat(rec map[string]any, field string, row int,
) (float64, error) {
	v, ok := rec[field]
	if !ok {
		return 0, fmt.Errorf("row %d: missing field %q", row, field)
	}
	n, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("row %d: field %q: expected number, "+
			"got %T", row, field, v)
	}
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("row %d: field %q: "+
			"non-finite value", row, field)
	}
	return n, nil
}

func jsonString(rec map[string]any, field string, row int,
) (string, error) {
	v, ok := rec[field]
	if !ok {
		return "", fmt.Errorf(
			"row %d: missing field %q", row, field)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("row %d: field %q: expected string, "+
			"got %T", row, field, v)
	}
	return s, nil
}

func decodeRecords(r io.Reader) ([]map[string]any, error) {
	if r == nil {
		return nil, errNilReader
	}
	var records []map[string]any
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}

func capPrealloc(n int) int {
	return min(n, maxPrealloc)
}

// XYFromJSON reads an XY series from a JSON array of objects.
//
// Example input:
//
//	[{"x": 1, "y": 10}, {"x": 2, "y": 20}]
func XYFromJSON(
	r io.Reader, name, xField, yField string,
) (XY, error) {
	records, err := decodeRecords(r)
	if err != nil {
		return XY{}, fmt.Errorf("series.XYFromJSON: %w", err)
	}
	pts := make([]Point, 0, capPrealloc(len(records)))
	for i, rec := range records {
		x, err := jsonFloat(rec, xField, i)
		if err != nil {
			return XY{}, fmt.Errorf("series.XYFromJSON: %w", err)
		}
		y, err := jsonFloat(rec, yField, i)
		if err != nil {
			return XY{}, fmt.Errorf("series.XYFromJSON: %w", err)
		}
		pts = append(pts, Point{X: x, Y: y})
	}
	return XY{name: name, Points: pts}, nil
}

// XYMultiFromJSON reads multiple XY series from a JSON array of
// objects sharing a single X field. Each yField produces one
// series named after the field.
func XYMultiFromJSON(
	r io.Reader, xField string, yFields []string,
) ([]XY, error) {
	records, err := decodeRecords(r)
	if err != nil {
		return nil, fmt.Errorf(
			"series.XYMultiFromJSON: %w", err)
	}
	allPts := make([][]Point, len(yFields))
	for i := range allPts {
		allPts[i] = make([]Point, 0, capPrealloc(len(records)))
	}
	for i, rec := range records {
		x, err := jsonFloat(rec, xField, i)
		if err != nil {
			return nil, fmt.Errorf(
				"series.XYMultiFromJSON: %w", err)
		}
		for si, yf := range yFields {
			y, err := jsonFloat(rec, yf, i)
			if err != nil {
				return nil, fmt.Errorf(
					"series.XYMultiFromJSON: %w", err)
			}
			allPts[si] = append(allPts[si], Point{X: x, Y: y})
		}
	}
	result := make([]XY, len(yFields))
	for i, yf := range yFields {
		result[i] = XY{name: yf, Points: allPts[i]}
	}
	return result, nil
}

// CategoryFromJSON reads a Category series from a JSON array of
// objects with a label field and a value field.
//
// Example input:
//
//	[{"label": "North", "value": 45}, {"label": "South", "value": 32}]
func CategoryFromJSON(
	r io.Reader, name, labelField, valueField string,
) (Category, error) {
	records, err := decodeRecords(r)
	if err != nil {
		return Category{}, fmt.Errorf(
			"series.CategoryFromJSON: %w", err)
	}
	vals := make([]CategoryValue, 0, capPrealloc(len(records)))
	for i, rec := range records {
		label, err := jsonString(rec, labelField, i)
		if err != nil {
			return Category{}, fmt.Errorf(
				"series.CategoryFromJSON: %w", err)
		}
		v, err := jsonFloat(rec, valueField, i)
		if err != nil {
			return Category{}, fmt.Errorf(
				"series.CategoryFromJSON: %w", err)
		}
		vals = append(vals, CategoryValue{Label: label, Value: v})
	}
	return Category{name: name, Values: vals}, nil
}

// XYZFromJSON reads an XYZ series from a JSON array of objects.
//
// Example input:
//
//	[{"x": 1, "y": 2, "z": 3}]
func XYZFromJSON(
	r io.Reader, name, xField, yField, zField string,
) (XYZ, error) {
	records, err := decodeRecords(r)
	if err != nil {
		return XYZ{}, fmt.Errorf("series.XYZFromJSON: %w", err)
	}
	pts := make([]XYZPoint, 0, capPrealloc(len(records)))
	for i, rec := range records {
		x, err := jsonFloat(rec, xField, i)
		if err != nil {
			return XYZ{}, fmt.Errorf(
				"series.XYZFromJSON: %w", err)
		}
		y, err := jsonFloat(rec, yField, i)
		if err != nil {
			return XYZ{}, fmt.Errorf(
				"series.XYZFromJSON: %w", err)
		}
		z, err := jsonFloat(rec, zField, i)
		if err != nil {
			return XYZ{}, fmt.Errorf(
				"series.XYZFromJSON: %w", err)
		}
		pts = append(pts, XYZPoint{X: x, Y: y, Z: z})
	}
	return XYZ{name: name, Points: pts}, nil
}

// OHLCFromJSON reads an OHLCSeries from a JSON array of objects.
func OHLCFromJSON(
	r io.Reader, name string, cfg OHLCJSONCfg,
) (OHLCSeries, error) {
	cfg.applyDefaults()
	records, err := decodeRecords(r)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromJSON: %w", err)
	}
	pts := make([]OHLC, 0, capPrealloc(len(records)))
	for i, rec := range records {
		ts, err := jsonString(rec, cfg.TimeField, i)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: %w", err)
		}
		t, err := time.Parse(cfg.TimeLayout, ts)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: row %d, time: %w", i, err)
		}
		o, err := jsonFloat(rec, cfg.OpenField, i)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: %w", err)
		}
		h, err := jsonFloat(rec, cfg.HighField, i)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: %w", err)
		}
		l, err := jsonFloat(rec, cfg.LowField, i)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: %w", err)
		}
		c, err := jsonFloat(rec, cfg.CloseField, i)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromJSON: %w", err)
		}
		vol := 0.0
		if cfg.VolumeField != "" {
			vol, err = jsonFloat(rec, cfg.VolumeField, i)
			if err != nil {
				return OHLCSeries{}, fmt.Errorf(
					"series.OHLCFromJSON: %w", err)
			}
		}
		pts = append(pts, OHLC{
			Time: t, Open: o, High: h, Low: l, Close: c,
			Volume: vol,
		})
	}
	return OHLCSeries{name: name, Points: pts}, nil
}

// gridJSON is the expected JSON structure for GridFromJSON.
type gridJSON struct {
	Rows   []string    `json:"rows"`
	Cols   []string    `json:"cols"`
	Values [][]float64 `json:"values"`
}

// GridFromJSON reads a Grid from a JSON object with structure:
//
//	{"rows": [...], "cols": [...], "values": [[...], ...]}
func GridFromJSON(r io.Reader, name string) (Grid, error) {
	if r == nil {
		return Grid{}, fmt.Errorf(
			"series.GridFromJSON: %w", errNilReader)
	}
	var g gridJSON
	if err := json.NewDecoder(r).Decode(&g); err != nil {
		return Grid{}, fmt.Errorf("series.GridFromJSON: %w", err)
	}
	return NewGrid(GridCfg{
		Name:   name,
		Rows:   g.Rows,
		Cols:   g.Cols,
		Values: g.Values,
	})
}

// treeNodeJSON is the expected JSON structure for
// TreeNodeFromJSON.
type treeNodeJSON struct {
	Label    string         `json:"label"`
	Value    float64        `json:"value"`
	Children []treeNodeJSON `json:"children"`
}

func (n treeNodeJSON) toTreeNode() TreeNode {
	children := make([]TreeNode, len(n.Children))
	for i := range n.Children {
		children[i] = n.Children[i].toTreeNode()
	}
	return TreeNode{
		Label:    n.Label,
		Value:    n.Value,
		Children: children,
	}
}

// TreeNodeFromJSON reads a TreeNode hierarchy from JSON.
//
// Expected shape:
//
//	{"label": "root", "value": 0, "children": [...]}
func TreeNodeFromJSON(r io.Reader) (TreeNode, error) {
	if r == nil {
		return TreeNode{}, fmt.Errorf(
			"series.TreeNodeFromJSON: %w", errNilReader)
	}
	var n treeNodeJSON
	if err := json.NewDecoder(r).Decode(&n); err != nil {
		return TreeNode{}, fmt.Errorf(
			"series.TreeNodeFromJSON: %w", err)
	}
	return n.toTreeNode(), nil
}
