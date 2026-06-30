package lab

const (
	BacktestFirstBTC15MPreviousDayRangeReversionName = "backtest_first_btc_15m_previous_day_range_reversion"
	BTC15MPreviousDayRangeReversionCandidateID       = "btc_15m_previous_day_range_reversion_v1"

	BTC15MPreviousDayRangeReversionStopStatePassedNeedsReview      = "btc_15m_previous_day_range_reversion_backtest_passed_needs_review"
	BTC15MPreviousDayRangeReversionStopStateFailedNoUsableStrategy = "btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy"
	BTC15MPreviousDayRangeReversionStopStateFailedSourceOrResample = "btc_15m_previous_day_range_reversion_backtest_failed_source_or_resample"
)

type BacktestFirstBTC15MPreviousDayRangeReversionConfig struct {
	ApprovedSourcePath      string
	ExpectedSourceRows      int
	ExpectedFirstOpenTime   string
	ExpectedLastOpenTime    string
	ExpectedZeroVolumeCount int
	Expected15MRows         int
	Expected15MLastOpenTime string
	ATRPeriod               int
	OuterZonePct            float64
	StopATRMultiple         float64
	MaxHoldBars             int
	MinFullTrades           int
	MinSplitTrades          int
	FullMaxDrawdownLimit    float64
	SplitMaxDrawdownLimit   float64
}

type BacktestFirstBTC15MPreviousDayRangeReversionResult struct {
	SourceRows    []BTC15MPreviousDayRangeReversionSourceRow   `json:"source_rows"`
	CoverageRows  []BTC15MPreviousDayRangeReversionCoverageRow `json:"coverage_rows"`
	SignalRows    []BTC15MPreviousDayRangeReversionSignalRow   `json:"signal_rows"`
	SkipRows      []BTC15MPreviousDayRangeReversionSkipRow     `json:"skip_rows"`
	TradeRows     []BTC15MPreviousDayRangeReversionTradeRow    `json:"trade_rows"`
	SummaryRows   []SummaryRow                                 `json:"summary_rows"`
	Falsification BTC15MPreviousDayRangeReversionFalsification `json:"falsification"`
	Trades        []Trade                                      `json:"trades"`
	StopState     string                                       `json:"stop_state"`
}

func DefaultBacktestFirstBTC15MPreviousDayRangeReversionConfig() BacktestFirstBTC15MPreviousDayRangeReversionConfig {
	return BacktestFirstBTC15MPreviousDayRangeReversionConfig{
		ApprovedSourcePath:      "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedSourceRows:      573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedZeroVolumeCount: 66,
		Expected15MRows:         191328,
		Expected15MLastOpenTime: "2026-06-16T23:45:00Z",
		ATRPeriod:               14,
		OuterZonePct:            0.10,
		StopATRMultiple:         0.25,
		MaxHoldBars:             24,
		MinFullTrades:           120,
		MinSplitTrades:          25,
		FullMaxDrawdownLimit:    0.25,
		SplitMaxDrawdownLimit:   0.30,
	}
}
