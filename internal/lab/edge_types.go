package lab

const (
	BacktestFirstBTC15MRangeEdgeExhaustionFadeName = "backtest_first_btc_15m_range_edge_exhaustion_fade"
	BTC15MRangeEdgeExhaustionFadeCandidateID       = "btc_15m_range_edge_exhaustion_fade_v1"

	BTC15MRangeEdgeExhaustionFadeStopStatePassedNeedsReview      = "btc_15m_range_edge_exhaustion_fade_backtest_passed_needs_review"
	BTC15MRangeEdgeExhaustionFadeStopStateFailedNoUsableStrategy = "btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy"
	BTC15MRangeEdgeExhaustionFadeStopStateFailedSourceOrResample = "btc_15m_range_edge_exhaustion_fade_backtest_failed_source_or_resample"
)

type BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig struct {
	ApprovedSourcePath      string
	ExpectedSourceRows      int
	ExpectedFirstOpenTime   string
	ExpectedLastOpenTime    string
	ExpectedZeroVolumeCount int
	Expected15MRows         int
	Expected15MLastOpenTime string
	LookbackBars            int
	ATRPeriod               int
	InteriorLowPct          float64
	InteriorHighPct         float64
	EdgeZonePct             float64
	ProgressATRMultiple     float64
	StopATRMultiple         float64
	MaxHoldBars             int
	MinFullTrades           int
	MinSplitTrades          int
	FullMaxDrawdownLimit    float64
	SplitMaxDrawdownLimit   float64
}

type BacktestFirstBTC15MRangeEdgeExhaustionFadeResult struct {
	SourceRows    []BTC15MRangeEdgeExhaustionFadeSourceRow   `json:"source_rows"`
	CoverageRows  []BTC15MRangeEdgeExhaustionFadeCoverageRow `json:"coverage_rows"`
	SignalRows    []BTC15MRangeEdgeExhaustionFadeSignalRow   `json:"signal_rows"`
	SkipRows      []BTC15MRangeEdgeExhaustionFadeSkipRow     `json:"skip_rows"`
	TradeRows     []BTC15MRangeEdgeExhaustionFadeTradeRow    `json:"trade_rows"`
	SummaryRows   []SummaryRow                               `json:"summary_rows"`
	Falsification BTC15MRangeEdgeExhaustionFadeFalsification `json:"falsification"`
	Trades        []Trade                                    `json:"trades"`
	StopState     string                                     `json:"stop_state"`
}

func DefaultBacktestFirstBTC15MRangeEdgeExhaustionFadeConfig() BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig {
	return BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig{
		ApprovedSourcePath:      "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedSourceRows:      573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedZeroVolumeCount: 66,
		Expected15MRows:         191328,
		Expected15MLastOpenTime: "2026-06-16T23:45:00Z",
		LookbackBars:            96,
		ATRPeriod:               14,
		InteriorLowPct:          0.40,
		InteriorHighPct:         0.60,
		EdgeZonePct:             0.15,
		ProgressATRMultiple:     0.35,
		StopATRMultiple:         0.25,
		MaxHoldBars:             16,
		MinFullTrades:           120,
		MinSplitTrades:          25,
		FullMaxDrawdownLimit:    0.25,
		SplitMaxDrawdownLimit:   0.30,
	}
}
