package chart

import "testing"

func TestAnimProgressDefaultsToOne(t *testing.T) {
	// nil window returns 1.0 (fully visible).
	if got := animProgress(nil, "test"); got != 1 {
		t.Errorf("animProgress(nil, id) = %v, want 1", got)
	}
}

func TestAnimProgressEmptyID(t *testing.T) {
	if got := animProgress(nil, ""); got != 1 {
		t.Errorf("animProgress(nil, \"\") = %v, want 1", got)
	}
}

func TestLoadAnimVersionNilWindow(t *testing.T) {
	if got := loadAnimVersion(nil, "test"); got != 0 {
		t.Errorf("loadAnimVersion(nil, id) = %v, want 0", got)
	}
}

func TestLoadAnimVersionEmptyID(t *testing.T) {
	if got := loadAnimVersion(nil, ""); got != 0 {
		t.Errorf("loadAnimVersion(nil, \"\") = %v, want 0", got)
	}
}

func TestTransitionProgressDefaultsToOne(t *testing.T) {
	if got := transitionProgress(nil, "test"); got != 1 {
		t.Errorf("transitionProgress(nil, id) = %v, want 1", got)
	}
}

func TestTransitionActiveNilWindow(t *testing.T) {
	if got := transitionActive(nil, "test"); got {
		t.Error("transitionActive(nil, id) = true, want false")
	}
}

func TestLoadTransitionVersionNilWindow(t *testing.T) {
	if got := loadTransitionVersion(nil, "test"); got != 0 {
		t.Errorf("loadTransitionVersion(nil, id) = %v, want 0", got)
	}
}

func TestStartEntryAnimationNilWindow(t *testing.T) {
	// Should not panic.
	startEntryAnimation(nil, "test", 0)
}

func TestStartTransitionNilWindow(t *testing.T) {
	// Should not panic.
	startTransition(nil, "test", 0)
}

func TestResetEntryAnimationNilWindow(t *testing.T) {
	// Should not panic.
	ResetEntryAnimation(nil, "test")
}
