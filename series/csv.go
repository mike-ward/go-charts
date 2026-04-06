package series

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

var errNilReader = errors.New("nil io.Reader")

// CSVCfg controls CSV parsing behavior. Zero value gives sensible
// defaults (comma delimiter, no header, no comment).
type CSVCfg struct {
	Delimiter  rune // field delimiter; 0 defaults to ','
	Comment    rune // comment character; 0 means none
	HasHeader  bool // first row is column headers
	TrimSpace  bool // trim leading/trailing whitespace from fields
	SkipErrors bool // skip rows with parse errors instead of failing
}

// Col identifies a CSV column by index or header name.
type Col struct {
	name string
	idx  int // 1-based when set via ColIdx (0 means use name)
	neg  bool
}

// ColIdx returns a Col that refers to a 0-based column index.
// Negative indices produce an error on resolve.
func ColIdx(i int) Col {
	if i < 0 {
		return Col{neg: true}
	}
	return Col{idx: i + 1}
}

// ColName returns a Col that refers to a column by header name.
func ColName(s string) Col { return Col{name: s} }

// resolve returns the 0-based column index. When Col uses a name,
// headers must contain that name.
func (c Col) resolve(headers []string) (int, error) {
	if c.neg {
		return 0, errors.New("negative column index")
	}
	if c.idx > 0 {
		return c.idx - 1, nil
	}
	for i, h := range headers {
		if h == c.name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %q not found in headers %v",
		c.name, headers)
}

// OHLCCSVCfg configures OHLC CSV parsing.
type OHLCCSVCfg struct {
	CSVCfg
	TimeCol    Col
	OpenCol    Col
	HighCol    Col
	LowCol     Col
	CloseCol   Col
	VolumeCol  Col    // zero value means no volume column
	TimeLayout string // Go time layout; empty = time.RFC3339
}

func newCSVReader(r io.Reader, cfg CSVCfg) *csv.Reader {
	cr := csv.NewReader(r)
	if cfg.Delimiter != 0 {
		cr.Comma = cfg.Delimiter
	}
	if cfg.Comment != 0 {
		cr.Comment = cfg.Comment
	}
	cr.ReuseRecord = true
	cr.FieldsPerRecord = -1
	return cr
}

func readHeader(cr *csv.Reader, cfg CSVCfg) ([]string, error) {
	if !cfg.HasHeader {
		return nil, nil
	}
	row, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}
	h := make([]string, len(row))
	for i, v := range row {
		if cfg.TrimSpace {
			v = strings.TrimSpace(v)
		}
		h[i] = v
	}
	return h, nil
}

func csvParseFloat(s string, row, col int) (float64, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("row %d, col %d: %w", row, col, err)
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("row %d, col %d: "+
			"non-finite value %q", row, col, s)
	}
	return v, nil
}

func csvField(record []string, col int, row int,
	trimSpace bool,
) (string, error) {
	if col < 0 || col >= len(record) {
		return "", fmt.Errorf("row %d: column index %d out of range "+
			"(record has %d fields)", row, col, len(record))
	}
	s := record[col]
	if trimSpace {
		s = strings.TrimSpace(s)
	}
	return s, nil
}

// csvReadFloat reads a field and parses it as float64.
func csvReadFloat(
	record []string, col, row int, trimSpace bool,
) (float64, error) {
	s, err := csvField(record, col, row, trimSpace)
	if err != nil {
		return 0, err
	}
	return csvParseFloat(s, row, col)
}

// csvReadString reads a string field.
func csvReadString(
	record []string, col, row int, trimSpace bool,
) (string, error) {
	return csvField(record, col, row, trimSpace)
}

// XYFromCSV reads one XY series from CSV. The X and Y columns are
// identified by Col (index or header name).
func XYFromCSV(
	r io.Reader, name string, xCol, yCol Col, cfg CSVCfg,
) (XY, error) {
	if r == nil {
		return XY{}, fmt.Errorf("series.XYFromCSV: %w", errNilReader)
	}
	cr := newCSVReader(r, cfg)
	headers, err := readHeader(cr, cfg)
	if err != nil {
		return XY{}, fmt.Errorf("series.XYFromCSV: %w", err)
	}
	xi, err := xCol.resolve(headers)
	if err != nil {
		return XY{}, fmt.Errorf("series.XYFromCSV: xCol: %w", err)
	}
	yi, err := yCol.resolve(headers)
	if err != nil {
		return XY{}, fmt.Errorf("series.XYFromCSV: yCol: %w", err)
	}
	pts := make([]Point, 0, 64)
	row := 1
	if cfg.HasHeader {
		row = 2
	}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XY{}, fmt.Errorf("series.XYFromCSV: row %d: %w",
				row, err)
		}
		x, err := csvReadFloat(record, xi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XY{}, fmt.Errorf("series.XYFromCSV: %w", err)
		}
		y, err := csvReadFloat(record, yi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XY{}, fmt.Errorf("series.XYFromCSV: %w", err)
		}
		pts = append(pts, Point{X: x, Y: y})
		row++
	}
	return XY{name: name, Points: pts}, nil
}

// XYMultiFromCSV reads multiple XY series from CSV, sharing a
// single X column. Each yCols entry produces one series. When
// names is nil and HasHeader is true, series names come from
// the header row.
func XYMultiFromCSV(
	r io.Reader, xCol Col, yCols []Col,
	names []string, cfg CSVCfg,
) ([]XY, error) {
	if r == nil {
		return nil, fmt.Errorf(
			"series.XYMultiFromCSV: %w", errNilReader)
	}
	cr := newCSVReader(r, cfg)
	headers, err := readHeader(cr, cfg)
	if err != nil {
		return nil, fmt.Errorf("series.XYMultiFromCSV: %w", err)
	}
	xi, err := xCol.resolve(headers)
	if err != nil {
		return nil, fmt.Errorf("series.XYMultiFromCSV: xCol: %w", err)
	}
	yis := make([]int, len(yCols))
	for i, yc := range yCols {
		yis[i], err = yc.resolve(headers)
		if err != nil {
			return nil, fmt.Errorf(
				"series.XYMultiFromCSV: yCols[%d]: %w", i, err)
		}
	}
	// Determine series names.
	sNames := make([]string, len(yCols))
	if names != nil {
		copy(sNames, names)
	} else if cfg.HasHeader {
		for i, yi := range yis {
			if yi < len(headers) {
				sNames[i] = headers[yi]
			}
		}
	}
	allPts := make([][]Point, len(yCols))
	for i := range allPts {
		allPts[i] = make([]Point, 0, 64)
	}
	row := 1
	if cfg.HasHeader {
		row = 2
	}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return nil, fmt.Errorf(
				"series.XYMultiFromCSV: row %d: %w", row, err)
		}
		x, err := csvReadFloat(record, xi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return nil, fmt.Errorf("series.XYMultiFromCSV: %w", err)
		}
		// Parse all Y values before appending to keep series
		// lengths aligned when SkipErrors is true.
		yVals := make([]float64, len(yis))
		skip := false
		for si, yi := range yis {
			yVals[si], err = csvReadFloat(
				record, yi, row, cfg.TrimSpace)
			if err != nil {
				if cfg.SkipErrors {
					skip = true
					break
				}
				return nil, fmt.Errorf(
					"series.XYMultiFromCSV: %w", err)
			}
		}
		if skip {
			row++
			continue
		}
		for si, y := range yVals {
			allPts[si] = append(allPts[si],
				Point{X: x, Y: y})
		}
		row++
	}
	result := make([]XY, len(yCols))
	for i := range result {
		result[i] = XY{name: sNames[i], Points: allPts[i]}
	}
	return result, nil
}

// CategoryFromCSV reads a Category series from CSV with a label
// column and a value column.
func CategoryFromCSV(
	r io.Reader, name string, labelCol, valueCol Col, cfg CSVCfg,
) (Category, error) {
	if r == nil {
		return Category{}, fmt.Errorf(
			"series.CategoryFromCSV: %w", errNilReader)
	}
	cr := newCSVReader(r, cfg)
	headers, err := readHeader(cr, cfg)
	if err != nil {
		return Category{}, fmt.Errorf(
			"series.CategoryFromCSV: %w", err)
	}
	li, err := labelCol.resolve(headers)
	if err != nil {
		return Category{}, fmt.Errorf(
			"series.CategoryFromCSV: labelCol: %w", err)
	}
	vi, err := valueCol.resolve(headers)
	if err != nil {
		return Category{}, fmt.Errorf(
			"series.CategoryFromCSV: valueCol: %w", err)
	}
	vals := make([]CategoryValue, 0, 64)
	row := 1
	if cfg.HasHeader {
		row = 2
	}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return Category{}, fmt.Errorf(
				"series.CategoryFromCSV: row %d: %w", row, err)
		}
		label, err := csvReadString(record, li, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return Category{}, fmt.Errorf(
				"series.CategoryFromCSV: %w", err)
		}
		v, err := csvReadFloat(record, vi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return Category{}, fmt.Errorf(
				"series.CategoryFromCSV: %w", err)
		}
		vals = append(vals, CategoryValue{Label: label, Value: v})
		row++
	}
	return Category{name: name, Values: vals}, nil
}

// XYZFromCSV reads an XYZ series from CSV.
func XYZFromCSV(
	r io.Reader, name string, xCol, yCol, zCol Col, cfg CSVCfg,
) (XYZ, error) {
	if r == nil {
		return XYZ{}, fmt.Errorf(
			"series.XYZFromCSV: %w", errNilReader)
	}
	cr := newCSVReader(r, cfg)
	headers, err := readHeader(cr, cfg)
	if err != nil {
		return XYZ{}, fmt.Errorf("series.XYZFromCSV: %w", err)
	}
	xi, err := xCol.resolve(headers)
	if err != nil {
		return XYZ{}, fmt.Errorf("series.XYZFromCSV: xCol: %w", err)
	}
	yi, err := yCol.resolve(headers)
	if err != nil {
		return XYZ{}, fmt.Errorf("series.XYZFromCSV: yCol: %w", err)
	}
	zi, err := zCol.resolve(headers)
	if err != nil {
		return XYZ{}, fmt.Errorf("series.XYZFromCSV: zCol: %w", err)
	}
	pts := make([]XYZPoint, 0, 64)
	row := 1
	if cfg.HasHeader {
		row = 2
	}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XYZ{}, fmt.Errorf(
				"series.XYZFromCSV: row %d: %w", row, err)
		}
		x, err := csvReadFloat(record, xi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XYZ{}, fmt.Errorf("series.XYZFromCSV: %w", err)
		}
		y, err := csvReadFloat(record, yi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XYZ{}, fmt.Errorf("series.XYZFromCSV: %w", err)
		}
		z, err := csvReadFloat(record, zi, row, cfg.TrimSpace)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return XYZ{}, fmt.Errorf("series.XYZFromCSV: %w", err)
		}
		pts = append(pts, XYZPoint{X: x, Y: y, Z: z})
		row++
	}
	return XYZ{name: name, Points: pts}, nil
}

// ohlcColIdx holds resolved column indices for OHLC CSV parsing.
type ohlcColIdx struct {
	time, open, high, low, close, volume int
	hasVolume                            bool
	trimSpace                            bool
	layout                               string
}

func (o ohlcColIdx) parseRow(
	record []string, row int,
) (OHLC, error) {
	ts, err := csvReadString(record, o.time, row, o.trimSpace)
	if err != nil {
		return OHLC{}, err
	}
	t, err := time.Parse(o.layout, ts)
	if err != nil {
		return OHLC{}, fmt.Errorf("row %d, time: %w", row, err)
	}
	op, err := csvReadFloat(record, o.open, row, o.trimSpace)
	if err != nil {
		return OHLC{}, err
	}
	hi, err := csvReadFloat(record, o.high, row, o.trimSpace)
	if err != nil {
		return OHLC{}, err
	}
	lo, err := csvReadFloat(record, o.low, row, o.trimSpace)
	if err != nil {
		return OHLC{}, err
	}
	cl, err := csvReadFloat(record, o.close, row, o.trimSpace)
	if err != nil {
		return OHLC{}, err
	}
	vol := 0.0
	if o.hasVolume {
		vol, err = csvReadFloat(record, o.volume, row, o.trimSpace)
		if err != nil {
			return OHLC{}, err
		}
	}
	return OHLC{
		Time: t, Open: op, High: hi, Low: lo, Close: cl,
		Volume: vol,
	}, nil
}

// OHLCFromCSV reads an OHLCSeries from CSV.
func OHLCFromCSV(
	r io.Reader, name string, cfg OHLCCSVCfg,
) (OHLCSeries, error) {
	if r == nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: %w", errNilReader)
	}
	layout := cfg.TimeLayout
	if layout == "" {
		layout = time.RFC3339
	}
	cr := newCSVReader(r, cfg.CSVCfg)
	headers, err := readHeader(cr, cfg.CSVCfg)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: %w", err)
	}
	idx := ohlcColIdx{
		trimSpace: cfg.TrimSpace,
		layout:    layout,
	}
	idx.time, err = cfg.TimeCol.resolve(headers)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: TimeCol: %w", err)
	}
	idx.open, err = cfg.OpenCol.resolve(headers)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: OpenCol: %w", err)
	}
	idx.high, err = cfg.HighCol.resolve(headers)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: HighCol: %w", err)
	}
	idx.low, err = cfg.LowCol.resolve(headers)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: LowCol: %w", err)
	}
	idx.close, err = cfg.CloseCol.resolve(headers)
	if err != nil {
		return OHLCSeries{}, fmt.Errorf(
			"series.OHLCFromCSV: CloseCol: %w", err)
	}
	idx.hasVolume = cfg.VolumeCol.idx > 0 || cfg.VolumeCol.name != ""
	if idx.hasVolume {
		idx.volume, err = cfg.VolumeCol.resolve(headers)
		if err != nil {
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromCSV: VolumeCol: %w", err)
		}
	}
	pts := make([]OHLC, 0, 64)
	row := 1
	if cfg.HasHeader {
		row = 2
	}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromCSV: row %d: %w", row, err)
		}
		pt, err := idx.parseRow(record, row)
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return OHLCSeries{}, fmt.Errorf(
				"series.OHLCFromCSV: %w", err)
		}
		pts = append(pts, pt)
		row++
	}
	return OHLCSeries{name: name, Points: pts}, nil
}

// GridFromCSV reads a Grid series from CSV. The first column is
// treated as row labels; the header row provides column labels.
// HasHeader must be true. All non-label cells must be numeric.
func GridFromCSV(
	r io.Reader, name string, cfg CSVCfg,
) (Grid, error) {
	if r == nil {
		return Grid{}, fmt.Errorf(
			"series.GridFromCSV: %w", errNilReader)
	}
	if !cfg.HasHeader {
		return Grid{}, fmt.Errorf(
			"series.GridFromCSV: HasHeader must be true")
	}
	cr := newCSVReader(r, cfg)
	headers, err := readHeader(cr, cfg)
	if err != nil {
		return Grid{}, fmt.Errorf("series.GridFromCSV: %w", err)
	}
	if len(headers) < 2 {
		return Grid{}, fmt.Errorf(
			"series.GridFromCSV: need at least 2 columns, got %d",
			len(headers))
	}
	cols := headers[1:]
	var rows []string
	var values [][]float64
	row := 2
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if cfg.SkipErrors {
				row++
				continue
			}
			return Grid{}, fmt.Errorf(
				"series.GridFromCSV: row %d: %w", row, err)
		}
		if len(record) < len(headers) {
			if cfg.SkipErrors {
				row++
				continue
			}
			return Grid{}, fmt.Errorf(
				"series.GridFromCSV: row %d: expected %d fields, "+
					"got %d", row, len(headers), len(record))
		}
		label := record[0]
		if cfg.TrimSpace {
			label = strings.TrimSpace(label)
		}
		rowVals := make([]float64, len(cols))
		skip := false
		for i := range cols {
			s := record[i+1]
			if cfg.TrimSpace {
				s = strings.TrimSpace(s)
			}
			rowVals[i], err = csvParseFloat(s, row, i+1)
			if err != nil {
				if cfg.SkipErrors {
					skip = true
					break
				}
				return Grid{}, fmt.Errorf(
					"series.GridFromCSV: %w", err)
			}
		}
		if skip {
			row++
			continue
		}
		rows = append(rows, label)
		values = append(values, rowVals)
		row++
	}
	return Grid{
		name:   name,
		rows:   rows,
		cols:   cols,
		values: values,
	}, nil
}
