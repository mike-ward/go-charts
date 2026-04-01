package render

// TextAnchor controls text alignment relative to a point.
type TextAnchor uint8

// TextAnchor constants.
const (
	TextAnchorStart TextAnchor = iota
	TextAnchorMiddle
	TextAnchorEnd
)

// TextBaseline controls vertical text alignment.
type TextBaseline uint8

// TextBaseline constants.
const (
	TextBaselineTop TextBaseline = iota
	TextBaselineMiddle
	TextBaselineBottom
)

// TODO: Text rendering requires gui.TextMeasurer integration
// for measuring text width/height before positioning labels.
