package chart

import "strconv"

// Validate checks LineCfg for invalid or contradictory settings.
// Returns nil when valid.
func (c *LineCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Series) == 0 {
		errs = append(errs, "no series data")
	}
	if c.LineWidth < 0 {
		errs = append(errs, "negative LineWidth")
	}
	return buildError("chart.Line", errs)
}

// Validate checks BarCfg for invalid or contradictory settings.
// Returns nil when valid.
func (c *BarCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Series) == 0 {
		errs = append(errs, "no series data")
	}
	if c.BarWidth < 0 {
		errs = append(errs, "negative BarWidth")
	}
	if c.BarGap < 0 {
		errs = append(errs, "negative BarGap")
	}
	if c.Radius < 0 {
		errs = append(errs, "negative Radius")
	}
	if len(c.Series) > 1 {
		n := len(c.Series[0].Values)
		for i, s := range c.Series[1:] {
			if len(s.Values) != n {
				errs = append(errs,
					"series length mismatch: series 0 has "+
						strconv.Itoa(n)+" values, series "+
						strconv.Itoa(i+1)+" has "+
						strconv.Itoa(len(s.Values)))
				break
			}
		}
	}
	return buildError("chart.Bar", errs)
}

// Validate checks AreaCfg for invalid or contradictory settings.
// Returns nil when valid.
func (c *AreaCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Series) == 0 {
		errs = append(errs, "no series data")
	}
	if c.LineWidth < 0 {
		errs = append(errs, "negative LineWidth")
	}
	if c.Opacity < 0 || c.Opacity > 1 {
		errs = append(errs, "Opacity out of range [0,1]")
	}
	return buildError("chart.Area", errs)
}

// Validate checks ScatterCfg for invalid or contradictory
// settings. Returns nil when valid.
func (c *ScatterCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Series) == 0 {
		errs = append(errs, "no series data")
	}
	if c.MarkerSize < 0 {
		errs = append(errs, "negative MarkerSize")
	}
	return buildError("chart.Scatter", errs)
}

// Validate checks PieCfg for invalid or contradictory settings.
// Returns nil when valid.
func (c *PieCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Slices) == 0 {
		errs = append(errs, "no slice data")
	}
	if c.InnerRadius < 0 {
		errs = append(errs, "negative InnerRadius")
	}
	return buildError("chart.Pie", errs)
}

// Validate checks HistogramCfg for invalid settings.
// Returns nil when valid.
func (c *HistogramCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if c.Bins < 0 {
		errs = append(errs, "negative Bins")
	}
	if c.Radius < 0 {
		errs = append(errs, "negative Radius")
	}
	if len(c.BinEdges) == 1 {
		errs = append(errs, "BinEdges must have 0 or 2+ entries")
	}
	return buildError("chart.Histogram", errs)
}

// Validate checks WaterfallCfg for invalid settings.
// Returns nil when valid.
func (c *WaterfallCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Values) == 0 {
		errs = append(errs, "no values")
	}
	if c.BarWidth < 0 {
		errs = append(errs, "negative BarWidth")
	}
	if c.Radius < 0 {
		errs = append(errs, "negative Radius")
	}
	return buildError("chart.Waterfall", errs)
}

// Validate checks BoxPlotCfg for invalid settings.
// Returns nil when valid.
func (c *BoxPlotCfg) Validate() error {
	if err := c.BaseCfg.Validate(); err != nil {
		return err
	}
	var errs []string
	if len(c.Data) == 0 {
		errs = append(errs, "no data")
	}
	if c.BoxWidth < 0 {
		errs = append(errs, "negative BoxWidth")
	}
	if c.OutlierRadius < 0 {
		errs = append(errs, "negative OutlierRadius")
	}
	return buildError("chart.BoxPlot", errs)
}
