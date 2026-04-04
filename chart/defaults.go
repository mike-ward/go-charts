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

	DefaultGaugeArcAngle   float32 = 3 * 3.14159265 / 2 // 270°
	DefaultGaugeInnerRatio float32 = 0.7

	// DefaultCandleWidthRatio is the fraction of the slot width used
	// for the candle body when CandleWidth is 0.
	DefaultCandleWidthRatio float32 = 0.6

	// DefaultBoxWidthRatio is the fraction of the slot width used
	// for the box body when BoxWidth is 0.
	DefaultBoxWidthRatio float32 = 0.5

	// DefaultOutlierRadius is the radius of outlier dots when
	// OutlierRadius is 0.
	DefaultOutlierRadius float32 = 3

	// DefaultWaterfallWidthRatio is the fraction of the slot
	// width used for bars when BarWidth is 0.
	DefaultWaterfallWidthRatio float32 = 0.6

	// DefaultConnectorWidth is the line width for waterfall
	// connector lines between bars.
	DefaultConnectorWidth float32 = 1

	// DefaultAnnotationLineWidth is the line width for
	// annotation reference lines when Width is 0.
	DefaultAnnotationLineWidth float32 = 1.5

	// DefaultZoomFactor is the scale factor per scroll tick.
	// >1 zooms in; 1/factor zooms out.
	DefaultZoomFactor float64 = 1.15

	// DefaultMinDragPx is the minimum pixel distance to
	// distinguish a drag from a click.
	DefaultMinDragPx float32 = 4

	// DefaultMinZoomRange is the smallest allowed domain span
	// to prevent degenerate transforms.
	DefaultMinZoomRange float64 = 1e-12

	// zoomDoubleClickFrames is the frame-count threshold for
	// detecting a mouse double-click (~400ms at 60fps).
	zoomDoubleClickFrames uint64 = 24
)
