package axis

// Category is a categorical (discrete) axis.
type Category struct {
	title      string
	categories []string
}

// CategoryCfg configures a category axis.
type CategoryCfg struct {
	Title      string
	Categories []string
}

// NewCategory creates a category axis.
func NewCategory(cfg CategoryCfg) *Category {
	return &Category{
		title:      cfg.Title,
		categories: cfg.Categories,
	}
}

// Label implements Axis.
func (a *Category) Label() string { return a.title }

// SetRange implements Axis. No-op for categorical axes.
func (a *Category) SetRange(_, _ float64) {}

// Domain implements Axis. Returns [0, len(categories)-1].
func (a *Category) Domain() (float64, float64) {
	n := len(a.categories)
	if n == 0 {
		return 0, 0
	}
	return 0, float64(n - 1)
}

// SetOverrideDomain implements Axis. No-op for categorical axes.
func (a *Category) SetOverrideDomain(_ bool) {}

// Ticks implements Axis.
func (a *Category) Ticks(pixelMin, pixelMax float32) []Tick {
	n := len(a.categories)
	if n == 0 {
		return nil
	}
	span := pixelMax - pixelMin
	step := span / float32(n)
	ticks := make([]Tick, n)
	for i, label := range a.categories {
		ticks[i] = Tick{
			Value:    float64(i),
			Label:    label,
			Position: pixelMin + step*float32(i) + step/2,
		}
	}
	return ticks
}

// Transform implements Axis.
func (a *Category) Transform(value float64, pixelMin, pixelMax float32) float32 {
	n := len(a.categories)
	if n == 0 {
		return pixelMin
	}
	span := pixelMax - pixelMin
	step := span / float32(n)
	pos := pixelMin + step*float32(value) + step/2
	return max(pixelMin, min(pixelMax, pos))
}

// Invert implements Axis.
func (a *Category) Invert(pixel, pixelMin, pixelMax float32) float64 {
	n := len(a.categories)
	if n == 0 {
		return 0
	}
	span := pixelMax - pixelMin
	step := span / float32(n)
	if step == 0 {
		return 0
	}
	v := float64((pixel - pixelMin) / step)
	return max(0, min(float64(n-1), v))
}
