package theme

import (
	"testing"

	"github.com/mike-ward/go-gui/gui"
)

func TestWithAlpha(t *testing.T) {
	c := gui.Hex(0xFF0000)
	got := WithAlpha(c, 0.5)
	if got.A != 127 {
		t.Errorf("A = %d, want 127", got.A)
	}
	if got.R != 255 || got.G != 0 || got.B != 0 {
		t.Errorf("RGB changed: %v", got)
	}
}

func TestWithAlphaClamped(t *testing.T) {
	c := gui.Hex(0x000000)
	if got := WithAlpha(c, -1); got.A != 0 {
		t.Errorf("A = %d, want 0 for negative", got.A)
	}
	if got := WithAlpha(c, 2); got.A != 255 {
		t.Errorf("A = %d, want 255 for >1", got.A)
	}
}

func TestLighten(t *testing.T) {
	c := gui.Hex(0x000000)
	got := Lighten(c, 1.0)
	if got.R != 255 || got.G != 255 || got.B != 255 {
		t.Errorf("Lighten(black, 1.0) = %v, want white", got)
	}
}

func TestLightenZero(t *testing.T) {
	c := gui.Hex(0x804020)
	got := Lighten(c, 0)
	if got.R != c.R || got.G != c.G || got.B != c.B {
		t.Errorf("Lighten(c, 0) changed color: %v", got)
	}
}

func TestDarken(t *testing.T) {
	c := gui.Hex(0xFFFFFF)
	got := Darken(c, 1.0)
	if got.R != 0 || got.G != 0 || got.B != 0 {
		t.Errorf("Darken(white, 1.0) = %v, want black", got)
	}
}

func TestDarkenZero(t *testing.T) {
	c := gui.Hex(0x804020)
	got := Darken(c, 0)
	if got.R != c.R || got.G != c.G || got.B != c.B {
		t.Errorf("Darken(c, 0) changed color: %v", got)
	}
}
