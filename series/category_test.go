package series

import "testing"

func TestCategoryFromMap(t *testing.T) {
	m := map[string]float64{"Banana": 3, "Apple": 5, "Cherry": 1}
	s := CategoryFromMap("fruit", m)
	if s.Name() != "fruit" {
		t.Errorf("Name = %q, want %q", s.Name(), "fruit")
	}
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
	// Sorted alphabetically.
	if s.Values[0].Label != "Apple" {
		t.Errorf("Values[0].Label = %q, want %q", s.Values[0].Label, "Apple")
	}
	if s.Values[1].Label != "Banana" {
		t.Errorf("Values[1].Label = %q, want %q", s.Values[1].Label, "Banana")
	}
	if s.Values[2].Value != 1 {
		t.Errorf("Values[2].Value = %v, want 1", s.Values[2].Value)
	}
}

func TestCategoryFromMapEmpty(t *testing.T) {
	s := CategoryFromMap("empty", nil)
	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}
