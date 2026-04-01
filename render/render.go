// Package render provides chart rendering helpers wrapping
// gui.DrawContext.
package render

import "github.com/mike-ward/go-gui/gui"

// Context wraps gui.DrawContext with chart-specific drawing
// helpers.
type Context struct {
	DC *gui.DrawContext
}

// NewContext creates a rendering context from a gui.DrawContext.
func NewContext(dc *gui.DrawContext) *Context {
	return &Context{DC: dc}
}

// Width returns the available drawing width.
func (c *Context) Width() float32 { return c.DC.Width }

// Height returns the available drawing height.
func (c *Context) Height() float32 { return c.DC.Height }

// Line draws a line segment.
func (c *Context) Line(x0, y0, x1, y1 float32, color gui.Color, width float32) {
	c.DC.Line(x0, y0, x1, y1, color, width)
}

// Polyline draws a connected series of line segments.
func (c *Context) Polyline(points []float32, color gui.Color, width float32) {
	c.DC.Polyline(points, color, width)
}

// FilledRect draws a filled rectangle.
func (c *Context) FilledRect(x, y, w, h float32, color gui.Color) {
	c.DC.FilledRect(x, y, w, h, color)
}

// Rect draws a stroked rectangle.
func (c *Context) Rect(x, y, w, h float32, color gui.Color, width float32) {
	c.DC.Rect(x, y, w, h, color, width)
}

// FilledCircle draws a filled circle.
func (c *Context) FilledCircle(cx, cy, radius float32, color gui.Color) {
	c.DC.FilledCircle(cx, cy, radius, color)
}

// Circle draws a stroked circle.
func (c *Context) Circle(cx, cy, radius float32, color gui.Color, width float32) {
	c.DC.Circle(cx, cy, radius, color, width)
}

// FilledArc draws a filled elliptical arc (pie slice).
func (c *Context) FilledArc(cx, cy, rx, ry, start, sweep float32, color gui.Color) {
	c.DC.FilledArc(cx, cy, rx, ry, start, sweep, color)
}

// Arc draws a stroked elliptical arc.
func (c *Context) Arc(cx, cy, rx, ry, start, sweep float32, color gui.Color, width float32) {
	c.DC.Arc(cx, cy, rx, ry, start, sweep, color, width)
}

// FilledPolygon draws a filled convex polygon.
func (c *Context) FilledPolygon(points []float32, color gui.Color) {
	c.DC.FilledPolygon(points, color)
}
