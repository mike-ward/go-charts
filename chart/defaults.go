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
)
