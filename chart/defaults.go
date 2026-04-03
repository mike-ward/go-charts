package chart

// Default values for chart configuration. Applied when the
// corresponding Cfg field is zero-valued.
const (
	DefaultLineWidth   float32 = 2
	DefaultBarGap      float32 = 4
	DefaultMarkerSize  float32 = 6
	DefaultAreaOpacity float32 = 0.3
	DefaultTickLength  float32 = 5

	HoverDimAlpha    uint8   = 60 // alpha for non-hovered elements (~24%)
	HoverExplodeDist float32 = 8  // px offset for exploded pie segment
)
