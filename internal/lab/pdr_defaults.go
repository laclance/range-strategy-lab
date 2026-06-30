package lab

type btc15MDayRange struct {
	key      string
	high     float64
	low      float64
	mid      float64
	width    float64
	count    int
	complete bool
}

func (cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) withDefaults() BacktestFirstBTC15MPreviousDayRangeReversionConfig {
	defaults := DefaultBacktestFirstBTC15MPreviousDayRangeReversionConfig()
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
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.OuterZonePct == 0 {
		cfg.OuterZonePct = defaults.OuterZonePct
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
