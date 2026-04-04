package chart

import (
	"math"
	"testing"
)

func TestComputeBoxStats_Empty(t *testing.T) {
	_, ok := computeBoxStats(nil)
	if ok {
		t.Fatal("expected ok=false for nil input")
	}
	_, ok = computeBoxStats([]float64{})
	if ok {
		t.Fatal("expected ok=false for empty input")
	}
}

func TestComputeBoxStats_SingleValue(t *testing.T) {
	st, ok := computeBoxStats([]float64{42})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Q1 != 42 || st.Median != 42 || st.Q3 != 42 {
		t.Errorf("single value: Q1=%g Median=%g Q3=%g; want all 42",
			st.Q1, st.Median, st.Q3)
	}
	if len(st.Outliers) != 0 {
		t.Errorf("single value: %d outliers; want 0", len(st.Outliers))
	}
}

func TestComputeBoxStats_OddCount(t *testing.T) {
	// [1, 2, 3, 4, 5]: median=3, Q1=1.5, Q3=4.5
	st, ok := computeBoxStats([]float64{5, 3, 1, 4, 2})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Median != 3 {
		t.Errorf("median=%g; want 3", st.Median)
	}
	if st.Q1 != 1.5 {
		t.Errorf("Q1=%g; want 1.5", st.Q1)
	}
	if st.Q3 != 4.5 {
		t.Errorf("Q3=%g; want 4.5", st.Q3)
	}
}

func TestComputeBoxStats_EvenCount(t *testing.T) {
	// [1, 2, 3, 4]: median=2.5, Q1=1.5, Q3=3.5
	st, ok := computeBoxStats([]float64{4, 2, 1, 3})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Median != 2.5 {
		t.Errorf("median=%g; want 2.5", st.Median)
	}
	if st.Q1 != 1.5 {
		t.Errorf("Q1=%g; want 1.5", st.Q1)
	}
	if st.Q3 != 3.5 {
		t.Errorf("Q3=%g; want 3.5", st.Q3)
	}
}

func TestComputeBoxStats_WithOutliers(t *testing.T) {
	// Data with clear outliers outside 1.5*IQR fences.
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 50}
	st, ok := computeBoxStats(data)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(st.Outliers) == 0 {
		t.Fatal("expected at least one outlier")
	}
	// 50 should be an outlier.
	found := false
	for _, v := range st.Outliers {
		if v == 50 {
			found = true
		}
	}
	if !found {
		t.Errorf("50 not found in outliers: %v", st.Outliers)
	}
	// Whisker max should be less than 50.
	if st.Max >= 50 {
		t.Errorf("whisker max=%g; should be < 50", st.Max)
	}
}

func TestComputeBoxStats_NaNInfIgnored(t *testing.T) {
	data := []float64{1, math.NaN(), 2, math.Inf(1), 3, math.Inf(-1)}
	st, ok := computeBoxStats(data)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Median != 2 {
		t.Errorf("median=%g; want 2", st.Median)
	}
}

func TestComputeBoxStats_AllIdentical(t *testing.T) {
	st, ok := computeBoxStats([]float64{5, 5, 5, 5})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Q1 != 5 || st.Median != 5 || st.Q3 != 5 {
		t.Errorf("Q1=%g Median=%g Q3=%g; want all 5",
			st.Q1, st.Median, st.Q3)
	}
	if len(st.Outliers) != 0 {
		t.Errorf("%d outliers; want 0", len(st.Outliers))
	}
}

func TestComputeBoxStats_TwoValues(t *testing.T) {
	st, ok := computeBoxStats([]float64{10, 20})
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Median != 15 {
		t.Errorf("median=%g; want 15", st.Median)
	}
	if st.Q1 != 10 {
		t.Errorf("Q1=%g; want 10", st.Q1)
	}
	if st.Q3 != 20 {
		t.Errorf("Q3=%g; want 20", st.Q3)
	}
}

func TestComputeBoxStats_AllNonFinite(t *testing.T) {
	_, ok := computeBoxStats([]float64{
		math.NaN(), math.Inf(1), math.Inf(-1),
	})
	if ok {
		t.Fatal("expected ok=false for all non-finite input")
	}
}

func TestComputeBoxStats_WhiskerBounds(t *testing.T) {
	// [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 200]
	// Q1=3, Q3=9, IQR=6, fences=[-6, 18]
	// Whisker min=1 (>=−6), whisker max=10 (<=18)
	// Outliers: 100, 200
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 200}
	st, ok := computeBoxStats(data)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if st.Min != 1 {
		t.Errorf("whisker min=%g; want 1", st.Min)
	}
	if st.Max != 10 {
		t.Errorf("whisker max=%g; want 10", st.Max)
	}
	if len(st.Outliers) != 2 {
		t.Fatalf("outliers=%d; want 2", len(st.Outliers))
	}
}

func TestComputeBoxStats_OutliersBothSides(t *testing.T) {
	// Data with outliers on both ends.
	data := []float64{-100, -50, 10, 20, 30, 40, 50, 60, 70, 150, 200}
	st, ok := computeBoxStats(data)
	if !ok {
		t.Fatal("expected ok=true")
	}
	// Outliers must be outside the whisker range.
	hasLow, hasHigh := false, false
	for _, v := range st.Outliers {
		if v >= st.Min && v <= st.Max {
			t.Errorf("outlier %g inside whisker range [%g, %g]",
				v, st.Min, st.Max)
		}
		if v < st.Min {
			hasLow = true
		}
		if v > st.Max {
			hasHigh = true
		}
	}
	if !hasLow {
		t.Error("expected low-side outlier")
	}
	if !hasHigh {
		t.Error("expected high-side outlier")
	}
}

func TestBoxPlotValidate_NoData(t *testing.T) {
	cfg := BoxPlotCfg{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for no data")
	}
}

func TestBoxPlotValidate_NegativeBoxWidth(t *testing.T) {
	cfg := BoxPlotCfg{
		Data:     []BoxData{{Label: "a", Values: []float64{1}}},
		BoxWidth: -1,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative BoxWidth")
	}
}

func TestBoxPlotValidate_NegativeOutlierRadius(t *testing.T) {
	cfg := BoxPlotCfg{
		Data:          []BoxData{{Label: "a", Values: []float64{1}}},
		OutlierRadius: -1,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative OutlierRadius")
	}
}

func TestBoxPlotValidate_Valid(t *testing.T) {
	cfg := BoxPlotCfg{
		Data: []BoxData{{Label: "a", Values: []float64{1, 2, 3}}},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
