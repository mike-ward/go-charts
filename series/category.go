package series

import "github.com/mike-ward/go-gui/gui"

// CategoryValue represents a labeled value.
type CategoryValue struct {
	Label string
	Value float64
}

// Category is a series of labeled categorical data.
type Category struct {
	name   string
	color  gui.Color
	Values []CategoryValue
}

// CategoryCfg configures a category series.
type CategoryCfg struct {
	Name   string
	Color  gui.Color
	Values []CategoryValue
}

// NewCategory creates a new category data series.
func NewCategory(cfg CategoryCfg) Category {
	return Category{
		name:   cfg.Name,
		color:  cfg.Color,
		Values: cfg.Values,
	}
}

// Name implements Series.
func (s Category) Name() string { return s.name }

// Len implements Series.
func (s Category) Len() int { return len(s.Values) }

// Color implements Series.
func (s Category) Color() gui.Color { return s.color }
