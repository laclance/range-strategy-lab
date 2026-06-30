package lab

type btc15MEdgeRange struct {
	high  float64
	low   float64
	mid   float64
	width float64
}

func (cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig) withDefaults() BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig {
	defaults := DefaultBacktestFirstBTC15MRangeEdgeExhaustionFadeConfig()
	if cfg.ApprovedSourcePath == "" {
		cfg.ApprovedSourcePath = defaults.ApprovedSourcePath
	}
	if cfg.ExpectedSourceRows == 0 {
		cfg.ExpectedSourceRows = defaults.ExpectedSourceRows
	}
	if cfg.ExpectedFirstOpenTime == "" {
		cfg.ExpectedFirstOpenTime = defaults.ExpectedFirstOpenTime
	}
	if cfg.ExpectedLastOpenTime == "" {
		cfg.ExpectedLastOpenTime = defaults.ExpectedLastOpenTime
	}
	if cfg.ExpectedZeroVolumeCount == 0 {
		cfg.ExpectedZeroVolumeCount = defaults.ExpectedZeroVolumeCount
	}
	if cfg.Expected15MRows == 0 {
		cfg.Expected15MRows = defaults.Expected15MRows
	}
	if cfg.Expected15MLastOpenTime == "" {
		cfg.Expected15MLastOpenTime = defaults.Expected15MLastOpenTime
	}
	if cfg.LookbackBars == 0 {
		cfg.LookbackBars = defaults.LookbackBars
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.InteriorLowPct == 0 {
		cfg.InteriorLowPct = defaults.InteriorLowPct
	}
	if cfg.InteriorHighPct == 0 {
		cfg.InteriorHighPct = defaults.InteriorHighPct
	}
	if cfg.EdgeZonePct == 0 {
		cfg.EdgeZonePct = defaults.EdgeZonePct
	}
	if cfg.ProgressATRMultiple == 0 {
		cfg.ProgressATRMultiple = defaults.ProgressATRMultiple
	}
	if cfg.StopATRMultiple == 0 {
		cfg.StopATRMultiple = defaults.StopATRMultiple
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = defaults.MaxHoldBars
	}
	if cfg.MinFullTrades == 0 {
		cfg.MinFullTrades = defaults.MinFullTrades
	}
	if cfg.MinSplitTrades == 0 {
		cfg.MinSplitTrades = defaults.MinSplitTrades
	}
	if cfg.FullMaxDrawdownLimit == 0 {
		cfg.FullMaxDrawdownLimit = defaults.FullMaxDrawdownLimit
	}
	if cfg.SplitMaxDrawdownLimit == 0 {
		cfg.SplitMaxDrawdownLimit = defaults.SplitMaxDrawdownLimit
	}
	return cfg
}
