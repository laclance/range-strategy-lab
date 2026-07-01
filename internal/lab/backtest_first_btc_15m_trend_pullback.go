package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	BacktestFirstBTC15MTrendPullbackContinuationName = "backtest_first_btc_15m_trend_pullback_continuation"
	BTC15MTrendPullbackContinuationCandidateID       = "btc_15m_trend_pullback_continuation_v1"

	BTC15MTrendPullbackContinuationStopStatePassedNeedsReview      = "btc_15m_trend_pullback_continuation_backtest_passed_needs_review"
	BTC15MTrendPullbackContinuationStopStateFailedNoUsableStrategy = "btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy"
	BTC15MTrendPullbackContinuationStopStateFailedSourceOrResample = "btc_15m_trend_pullback_continuation_backtest_failed_source_or_resample"
)

type BacktestFirstBTC15MTrendPullbackContinuationConfig struct {
	ApprovedSourcePath      string
	ExpectedSourceRows      int
	ExpectedFirstOpenTime   string
	ExpectedLastOpenTime    string
	ExpectedZeroVolumeCount int
	Expected15MRows         int
	Expected15MLastOpenTime string
	EMA20Period             int
	EMA50Period             int
	EMA200Period            int
	ATRPeriod               int
	SlopeLookbackBars       int
	PullbackWindowBars      int
	StopATRMultiple         float64
	TargetR                 float64
	MaxHoldBars             int
	MinFullTrades           int
	MinSplitTrades          int
	FullProfitFactorMin     float64
	SplitProfitFactorMin    float64
	FullMaxDrawdownLimit    float64
	SplitMaxDrawdownLimit   float64
}

type BacktestFirstBTC15MTrendPullbackContinuationResult struct {
	SourceRows    []BTC15MTrendPullbackContinuationSourceRow   `json:"source_rows"`
	CoverageRows  []BTC15MTrendPullbackContinuationCoverageRow `json:"coverage_rows"`
	SignalRows    []BTC15MTrendPullbackContinuationSignalRow   `json:"signal_rows"`
	SkipRows      []BTC15MTrendPullbackContinuationSkipRow     `json:"skip_rows"`
	TradeRows     []BTC15MTrendPullbackContinuationTradeRow    `json:"trade_rows"`
	SummaryRows   []SummaryRow                                 `json:"summary_rows"`
	Falsification BTC15MTrendPullbackContinuationFalsification `json:"falsification"`
	Trades        []Trade                                      `json:"trades"`
	StopState     string                                       `json:"stop_state"`
}

type BTC15MTrendPullbackContinuationSourceRow struct {
	BacktestName               string `json:"backtest_name"`
	CandidateID                string `json:"candidate_id"`
	Path                       string `json:"path"`
	ApprovedPath               string `json:"approved_path"`
	Venue                      string `json:"venue"`
	Product                    string `json:"product"`
	Symbol                     string `json:"symbol"`
	Interval                   string `json:"interval"`
	RowCount                   int    `json:"row_count"`
	ExpectedRowCount           int    `json:"expected_row_count"`
	FirstOpenTime              string `json:"first_open_time"`
	ExpectedFirstOpenTime      string `json:"expected_first_open_time"`
	LastOpenTime               string `json:"last_open_time"`
	ExpectedLastOpenTime       string `json:"expected_last_open_time"`
	GapCount                   int    `json:"gap_count"`
	DuplicateCount             int    `json:"duplicate_count"`
	ZeroVolumeCount            int    `json:"zero_volume_count"`
	ExpectedZeroVolumeCount    int    `json:"expected_zero_volume_count"`
	ComparisonOnly             bool   `json:"comparison_only"`
	ClosedCandleOnly           bool   `json:"closed_candle_only"`
	DerivativesVetoAsInput     bool   `json:"derivatives_veto_as_input"`
	ForwardLabelsAsSourceInput bool   `json:"forward_labels_as_source_input"`
	SourceFactsPass            bool   `json:"source_facts_pass"`
	ValidationStatus           string `json:"validation_status"`
	ValidationError            string `json:"validation_error,omitempty"`
}

type BTC15MTrendPullbackContinuationCoverageRow struct {
	BacktestName         string `json:"backtest_name"`
	CandidateID          string `json:"candidate_id"`
	Timeframe            string `json:"timeframe"`
	RowCount             int    `json:"row_count"`
	ExpectedRowCount     int    `json:"expected_row_count"`
	FirstOpenTime        string `json:"first_open_time"`
	LastOpenTime         string `json:"last_open_time"`
	ExpectedLastOpenTime string `json:"expected_last_open_time"`
	ExpectedChildBars    int    `json:"expected_child_bars"`
	MissingChildBuckets  int    `json:"missing_child_buckets"`
	ClosedCandleOnly     bool   `json:"closed_candle_only"`
	SourceResamplePass   bool   `json:"source_resample_pass"`
	ValidationStatus     string `json:"validation_status"`
	ValidationError      string `json:"validation_error,omitempty"`
}

type BTC15MTrendPullbackContinuationSignalRow struct {
	SignalID                  string    `json:"signal_id"`
	CandidateID               string    `json:"candidate_id"`
	Timeframe                 string    `json:"timeframe"`
	Split                     string    `json:"split"`
	DecisionIndex             int       `json:"decision_index"`
	DecisionOpenTime          string    `json:"decision_open_time"`
	DecisionCloseTime         string    `json:"decision_close_time"`
	DecisionOpen              float64   `json:"decision_open"`
	DecisionHigh              float64   `json:"decision_high"`
	DecisionLow               float64   `json:"decision_low"`
	DecisionClose             float64   `json:"decision_close"`
	Side                      Direction `json:"side"`
	TimingLabel               string    `json:"timing_label"`
	EMA20                     float64   `json:"ema20"`
	EMA50                     float64   `json:"ema50"`
	EMA200                    float64   `json:"ema200"`
	PriorEMA20                float64   `json:"prior_ema20"`
	PriorEMA50                float64   `json:"prior_ema50"`
	PriorEMA200               float64   `json:"prior_ema200"`
	EMA50SlopeLookbackBars    int       `json:"ema50_slope_lookback_bars"`
	EMA50Slope                float64   `json:"ema50_slope"`
	PullbackWindowBars        int       `json:"pullback_window_bars"`
	PullbackIndex             int       `json:"pullback_index"`
	PullbackLowestLow         float64   `json:"pullback_lowest_low"`
	PullbackHighestHigh       float64   `json:"pullback_highest_high"`
	ATR14                     float64   `json:"atr14"`
	EntryIndex                int       `json:"entry_index"`
	EntryOpenTime             string    `json:"entry_open_time"`
	EntryOpen                 float64   `json:"entry_open"`
	ExpectedEntryPrice        float64   `json:"expected_entry_price"`
	StopATRMultiple           float64   `json:"stop_atr_multiple"`
	TargetR                   float64   `json:"target_r"`
	Stop                      float64   `json:"stop"`
	Target                    float64   `json:"target"`
	MaxHoldBars               int       `json:"max_hold_bars"`
	EntryGeometryValid        bool      `json:"entry_geometry_valid"`
	PrePositionCandidate      bool      `json:"pre_position_candidate"`
	Executed                  bool      `json:"executed"`
	SkippedReason             string    `json:"skipped_reason,omitempty"`
	ForwardLabelsAsInput      bool      `json:"forward_labels_as_input"`
	UsesFutureRowsForFeatures bool      `json:"uses_future_rows_for_features"`
	DerivativesVetoUsed       bool      `json:"derivatives_veto_used"`
	OptimizerSelectionUsed    bool      `json:"optimizer_selection_used"`
}

type BTC15MTrendPullbackContinuationSkipRow struct {
	CandidateID       string `json:"candidate_id"`
	Split             string `json:"split"`
	Side              string `json:"side"`
	Reason            string `json:"reason"`
	Count             int    `json:"count"`
	MissingDataPolicy string `json:"missing_data_policy"`
	ForwardFilledRows int    `json:"forward_filled_rows"`
}

type BTC15MTrendPullbackContinuationTradeRow struct {
	SignalID      string    `json:"signal_id"`
	CandidateID   string    `json:"candidate_id"`
	Timeframe     string    `json:"timeframe"`
	DecisionIndex int       `json:"decision_index"`
	EntrySplit    string    `json:"entry_split"`
	CloseSplit    string    `json:"close_split"`
	Side          Direction `json:"side"`
	EntryTime     string    `json:"entry_time"`
	ExitTime      string    `json:"exit_time"`
	OpenIndex     int       `json:"open_index"`
	CloseIndex    int       `json:"close_index"`
	EntryPrice    float64   `json:"entry_price"`
	ExitPrice     float64   `json:"exit_price"`
	Stop          float64   `json:"stop"`
	Target        float64   `json:"target"`
	Size          float64   `json:"size"`
	InitialRisk   float64   `json:"initial_risk"`
	GrossPnL      float64   `json:"gross_pnl"`
	NetPnL        float64   `json:"net_pnl"`
	Fees          float64   `json:"fees"`
	Slippage      float64   `json:"slippage"`
	GrossR        float64   `json:"gross_r"`
	NetR          float64   `json:"net_r"`
	ExitReason    string    `json:"exit_reason"`
	HoldBars      int       `json:"hold_bars"`
}

type BTC15MTrendPullbackContinuationFalsification struct {
	BacktestName                        string   `json:"backtest_name"`
	CandidateID                         string   `json:"candidate_id"`
	StopState                           string   `json:"stop_state"`
	SourceResamplePass                  bool     `json:"source_resample_pass"`
	LeakagePass                         bool     `json:"leakage_pass"`
	TradeCountPass                      bool     `json:"trade_count_pass"`
	GrossEdgePass                       bool     `json:"gross_edge_pass"`
	NetEdgePass                         bool     `json:"net_edge_pass"`
	ProfitFactorPass                    bool     `json:"profit_factor_pass"`
	DrawdownPass                        bool     `json:"drawdown_pass"`
	RobustnessPass                      bool     `json:"robustness_pass"`
	SideReportingPass                   bool     `json:"side_reporting_pass"`
	OptimizerContaminationPass          bool     `json:"optimizer_contamination_pass"`
	DerivativesVetoContaminationPass    bool     `json:"derivatives_veto_contamination_pass"`
	FullExecutedTrades                  int      `json:"full_executed_trades"`
	RequiredFullExecutedTrades          int      `json:"required_full_executed_trades"`
	MinimumPrimarySplitExecutedTrades   int      `json:"minimum_primary_split_executed_trades"`
	RequiredPrimarySplitExecutedTrades  int      `json:"required_primary_split_executed_trades"`
	FullGrossPnL                        float64  `json:"full_gross_pnl"`
	FullNetPnL                          float64  `json:"full_net_pnl"`
	FullProfitFactor                    float64  `json:"full_profit_factor"`
	FullMaxDrawdown                     float64  `json:"full_max_drawdown"`
	OOSGrossPnL                         float64  `json:"oos_gross_pnl"`
	OOSNetPnL                           float64  `json:"oos_net_pnl"`
	OOSProfitFactor                     float64  `json:"oos_profit_factor"`
	RecentGrossPnL                      float64  `json:"recent_gross_pnl"`
	RecentNetPnL                        float64  `json:"recent_net_pnl"`
	RecentProfitFactor                  float64  `json:"recent_profit_factor"`
	LongFullExecutedTrades              int      `json:"long_full_executed_trades"`
	LongFullNetPnL                      float64  `json:"long_full_net_pnl"`
	ShortFullExecutedTrades             int      `json:"short_full_executed_trades"`
	ShortFullNetPnL                     float64  `json:"short_full_net_pnl"`
	DominantPrimarySplitTradeShare      float64  `json:"dominant_primary_split_trade_share"`
	DominantPrimarySplitTradeShareLimit float64  `json:"dominant_primary_split_trade_share_limit"`
	FailureReasons                      []string `json:"failure_reasons,omitempty"`
}

type btc15MTrendPullbackSkipKey struct {
	split  string
	side   string
	reason string
}

func DefaultBacktestFirstBTC15MTrendPullbackContinuationConfig() BacktestFirstBTC15MTrendPullbackContinuationConfig {
	return BacktestFirstBTC15MTrendPullbackContinuationConfig{
		ApprovedSourcePath:      "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedSourceRows:      573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedZeroVolumeCount: 66,
		Expected15MRows:         191328,
		Expected15MLastOpenTime: "2026-06-16T23:45:00Z",
		EMA20Period:             20,
		EMA50Period:             50,
		EMA200Period:            200,
		ATRPeriod:               14,
		SlopeLookbackBars:       16,
		PullbackWindowBars:      8,
		StopATRMultiple:         0.25,
		TargetR:                 2.0,
		MaxHoldBars:             32,
		MinFullTrades:           120,
		MinSplitTrades:          25,
		FullProfitFactorMin:     1.10,
		SplitProfitFactorMin:    1.00,
		FullMaxDrawdownLimit:    0.25,
		SplitMaxDrawdownLimit:   0.30,
	}
}

func RunBacktestFirstBTC15MTrendPullbackContinuation(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC15MTrendPullbackContinuationResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return BacktestFirstBTC15MTrendPullbackContinuationResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc15MTrendPullbackSourceRow(manifest, cfg)
	result := BacktestFirstBTC15MTrendPullbackContinuationResult{SourceRows: []BTC15MTrendPullbackContinuationSourceRow{source}}
	resampled, coverage := btc15MTrendPullbackResample(candles, cfg, source.SourceFactsPass)
	result.CoverageRows = []BTC15MTrendPullbackContinuationCoverageRow{coverage}
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = btc15MTrendPullbackFalsification(BTC15MTrendPullbackContinuationFalsification{SourceResamplePass: false, LeakagePass: true, SideReportingPass: true, OptimizerContaminationPass: true, DerivativesVetoContaminationPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}

	signals, skipCounts := btc15MTrendPullbackBuildSignals(resampled, cfg, btCfg, splits)
	signals, trades := btc15MTrendPullbackRun(resampled, signals, &skipCounts, cfg, btCfg)
	result.Trades = trades
	result.SignalRows = signals
	result.SkipRows = btc15MTrendPullbackSkipRows(skipCounts)
	result.TradeRows = btc15MTrendPullbackTradeRows(trades, signals, splits)
	result.SummaryRows = SummarizeSplits(trades, btCfg.StartBalance, splits)
	result.Falsification = btc15MTrendPullbackFalsification(btc15MTrendPullbackEvaluate(source, coverage, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}

func (cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) withDefaults() BacktestFirstBTC15MTrendPullbackContinuationConfig {
	defaults := DefaultBacktestFirstBTC15MTrendPullbackContinuationConfig()
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
	if cfg.EMA20Period == 0 {
		cfg.EMA20Period = defaults.EMA20Period
	}
	if cfg.EMA50Period == 0 {
		cfg.EMA50Period = defaults.EMA50Period
	}
	if cfg.EMA200Period == 0 {
		cfg.EMA200Period = defaults.EMA200Period
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.SlopeLookbackBars == 0 {
		cfg.SlopeLookbackBars = defaults.SlopeLookbackBars
	}
	if cfg.PullbackWindowBars == 0 {
		cfg.PullbackWindowBars = defaults.PullbackWindowBars
	}
	if cfg.StopATRMultiple == 0 {
		cfg.StopATRMultiple = defaults.StopATRMultiple
	}
	if cfg.TargetR == 0 {
		cfg.TargetR = defaults.TargetR
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
	if cfg.FullProfitFactorMin == 0 {
		cfg.FullProfitFactorMin = defaults.FullProfitFactorMin
	}
	if cfg.SplitProfitFactorMin == 0 {
		cfg.SplitProfitFactorMin = defaults.SplitProfitFactorMin
	}
	if cfg.FullMaxDrawdownLimit == 0 {
		cfg.FullMaxDrawdownLimit = defaults.FullMaxDrawdownLimit
	}
	if cfg.SplitMaxDrawdownLimit == 0 {
		cfg.SplitMaxDrawdownLimit = defaults.SplitMaxDrawdownLimit
	}
	return cfg
}

func (cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) validate() error {
	if cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("trend-pullback approved source path is required")
	}
	if cfg.EMA20Period <= 0 || cfg.EMA50Period <= 0 || cfg.EMA200Period <= 0 || cfg.ATRPeriod <= 0 {
		return fmt.Errorf("trend-pullback indicator periods must be positive")
	}
	if !(cfg.EMA20Period < cfg.EMA50Period && cfg.EMA50Period < cfg.EMA200Period) {
		return fmt.Errorf("trend-pullback EMA periods must satisfy EMA20 < EMA50 < EMA200")
	}
	if cfg.SlopeLookbackBars <= 0 || cfg.PullbackWindowBars <= 0 || cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("trend-pullback lookback and hold bars must be positive")
	}
	if cfg.StopATRMultiple <= 0 || cfg.TargetR <= 0 {
		return fmt.Errorf("trend-pullback stop and target multiples must be positive")
	}
	if cfg.MinFullTrades <= 0 || cfg.MinSplitTrades <= 0 {
		return fmt.Errorf("trend-pullback trade-count gates must be positive")
	}
	if cfg.FullProfitFactorMin <= 0 || cfg.SplitProfitFactorMin <= 0 {
		return fmt.Errorf("trend-pullback profit-factor gates must be positive")
	}
	if cfg.FullMaxDrawdownLimit <= 0 || cfg.SplitMaxDrawdownLimit <= 0 {
		return fmt.Errorf("trend-pullback drawdown limits must be positive")
	}
	return nil
}

func btc15MTrendPullbackSourceRow(manifest SourceManifest, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) BTC15MTrendPullbackContinuationSourceRow {
	row := BTC15MTrendPullbackContinuationSourceRow{
		BacktestName:               BacktestFirstBTC15MTrendPullbackContinuationName,
		CandidateID:                BTC15MTrendPullbackContinuationCandidateID,
		Path:                       manifest.Path,
		ApprovedPath:               cfg.ApprovedSourcePath,
		Venue:                      manifest.Venue,
		Product:                    manifest.Product,
		Symbol:                     manifest.Symbol,
		Interval:                   manifest.Interval,
		RowCount:                   manifest.RowCount,
		ExpectedRowCount:           cfg.ExpectedSourceRows,
		FirstOpenTime:              manifest.FirstOpenTime,
		ExpectedFirstOpenTime:      cfg.ExpectedFirstOpenTime,
		LastOpenTime:               manifest.LastOpenTime,
		ExpectedLastOpenTime:       cfg.ExpectedLastOpenTime,
		GapCount:                   manifest.GapCount,
		DuplicateCount:             manifest.DuplicateCount,
		ZeroVolumeCount:            manifest.ZeroVolumeCount,
		ExpectedZeroVolumeCount:    cfg.ExpectedZeroVolumeCount,
		ComparisonOnly:             manifest.ComparisonOnly,
		ClosedCandleOnly:           true,
		DerivativesVetoAsInput:     false,
		ForwardLabelsAsSourceInput: false,
		ValidationStatus:           "accepted",
	}
	failures := []string{}
	if manifest.ValidationStatus != "accepted" {
		failures = append(failures, "source manifest is not accepted")
	}
	if manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly {
		failures = append(failures, "source must be Binance USDT-M futures and not comparison-only")
	}
	if manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" {
		failures = append(failures, "source must be BTCUSDT 5m")
	}
	if !sameCleanPath(manifest.Path, cfg.ApprovedSourcePath) {
		failures = append(failures, fmt.Sprintf("source path %q is not approved path %q", manifest.Path, cfg.ApprovedSourcePath))
	}
	if manifest.RowCount != cfg.ExpectedSourceRows {
		failures = append(failures, fmt.Sprintf("row_count=%d expected=%d", manifest.RowCount, cfg.ExpectedSourceRows))
	}
	if manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime {
		failures = append(failures, fmt.Sprintf("first_open_time=%s expected=%s", manifest.FirstOpenTime, cfg.ExpectedFirstOpenTime))
	}
	if manifest.LastOpenTime != cfg.ExpectedLastOpenTime {
		failures = append(failures, fmt.Sprintf("last_open_time=%s expected=%s", manifest.LastOpenTime, cfg.ExpectedLastOpenTime))
	}
	if manifest.GapCount != 0 || manifest.DuplicateCount != 0 {
		failures = append(failures, "source gaps or duplicates present")
	}
	if manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
		failures = append(failures, fmt.Sprintf("zero_volume_count=%d expected=%d", manifest.ZeroVolumeCount, cfg.ExpectedZeroVolumeCount))
	}
	row.SourceFactsPass = len(failures) == 0
	if !row.SourceFactsPass {
		row.ValidationStatus = "rejected"
		row.ValidationError = strings.Join(failures, "; ")
	}
	return row
}

func btc15MTrendPullbackResample(candles []Candle, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig, sourcePass bool) ([]Candle, BTC15MTrendPullbackContinuationCoverageRow) {
	resampled, dayCoverage := btc15MPrevDayResample(candles, BacktestFirstBTC15MPreviousDayRangeReversionConfig{
		Expected15MRows:         cfg.Expected15MRows,
		Expected15MLastOpenTime: cfg.Expected15MLastOpenTime,
	}, sourcePass)
	row := BTC15MTrendPullbackContinuationCoverageRow{
		BacktestName:         BacktestFirstBTC15MTrendPullbackContinuationName,
		CandidateID:          BTC15MTrendPullbackContinuationCandidateID,
		Timeframe:            "15m",
		RowCount:             dayCoverage.RowCount,
		ExpectedRowCount:     cfg.Expected15MRows,
		FirstOpenTime:        dayCoverage.FirstOpenTime,
		LastOpenTime:         dayCoverage.LastOpenTime,
		ExpectedLastOpenTime: cfg.Expected15MLastOpenTime,
		ExpectedChildBars:    3,
		MissingChildBuckets:  dayCoverage.MissingChildBuckets,
		ClosedCandleOnly:     true,
		SourceResamplePass:   dayCoverage.SourceResamplePass,
		ValidationStatus:     dayCoverage.ValidationStatus,
		ValidationError:      dayCoverage.ValidationError,
	}
	return resampled, row
}

func btc15MTrendPullbackBuildSignals(candles []Candle, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig, btCfg BacktestConfig, splits []Split) ([]BTC15MTrendPullbackContinuationSignalRow, map[btc15MTrendPullbackSkipKey]int) {
	ema20 := btc15MTrendPullbackEMA(candles, cfg.EMA20Period)
	ema50 := btc15MTrendPullbackEMA(candles, cfg.EMA50Period)
	ema200 := btc15MTrendPullbackEMA(candles, cfg.EMA200Period)
	atr := ATR(candles, cfg.ATRPeriod)
	warmup := maxInt(cfg.EMA200Period, cfg.PullbackWindowBars+1)
	signals := []BTC15MTrendPullbackContinuationSignalRow{}
	skips := map[btc15MTrendPullbackSkipKey]int{}
	for d := 0; d < len(candles); d++ {
		split := splitNameForCloseTime(candles[d].CloseTime, splits)
		if d < warmup {
			btc15MTrendPullbackAddSkip(skips, split, "all", "missing_warmup")
			continue
		}
		if d+1 >= len(candles) {
			btc15MTrendPullbackAddSkip(skips, split, "all", "missing_entry_candle")
			continue
		}
		if !validNumber(atr[d]) || atr[d] <= 0 || !validNumber(ema20[d]) || !validNumber(ema50[d]) || !validNumber(ema200[d]) || !validNumber(ema20[d-1]) || !validNumber(ema50[d-1]) || !validNumber(ema200[d-1]) || d-cfg.SlopeLookbackBars-1 < 0 || !validNumber(ema50[d-cfg.SlopeLookbackBars-1]) {
			btc15MTrendPullbackAddSkip(skips, split, "all", "missing_indicator")
			continue
		}
		longOK, longPullbackIndex, longLow, longHigh := btc15MTrendPullbackLongSetup(candles, d, ema20, ema50, ema200, cfg)
		shortOK, shortPullbackIndex, shortLow, shortHigh := btc15MTrendPullbackShortSetup(candles, d, ema20, ema50, ema200, cfg)
		if longOK && shortOK {
			btc15MTrendPullbackAddSkip(skips, split, "all", "ambiguous_both_sides")
			continue
		}
		if longOK {
			row := btc15MTrendPullbackSignalRow(candles, d, len(signals)+1, Long, longPullbackIndex, longLow, longHigh, ema20, ema50, ema200, atr, cfg, btCfg, splits)
			if row.SkippedReason != "" {
				btc15MTrendPullbackAddSkip(skips, row.Split, string(row.Side), row.SkippedReason)
			}
			signals = append(signals, row)
		}
		if shortOK {
			row := btc15MTrendPullbackSignalRow(candles, d, len(signals)+1, Short, shortPullbackIndex, shortLow, shortHigh, ema20, ema50, ema200, atr, cfg, btCfg, splits)
			if row.SkippedReason != "" {
				btc15MTrendPullbackAddSkip(skips, row.Split, string(row.Side), row.SkippedReason)
			}
			signals = append(signals, row)
		}
	}
	return signals, skips
}

func btc15MTrendPullbackLongSetup(candles []Candle, d int, ema20, ema50, ema200 []float64, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) (bool, int, float64, float64) {
	if !(ema20[d-1] > ema50[d-1] && ema50[d-1] > ema200[d-1]) {
		return false, 0, 0, 0
	}
	if !(ema50[d-1] > ema50[d-cfg.SlopeLookbackBars-1]) {
		return false, 0, 0, 0
	}
	if !(candles[d-1].Close > ema50[d-1]) {
		return false, 0, 0, 0
	}
	pullbackIndex := -1
	lowestLow := math.Inf(1)
	highestHigh := math.Inf(-1)
	start := d - cfg.PullbackWindowBars
	for p := start; p <= d-1; p++ {
		if !validNumber(ema20[p]) || !validNumber(ema50[p]) {
			return false, 0, 0, 0
		}
		lowestLow = math.Min(lowestLow, candles[p].Low)
		highestHigh = math.Max(highestHigh, candles[p].High)
		if candles[p].Close < ema50[p] {
			return false, 0, 0, 0
		}
		if candles[p].Low <= ema20[p] && candles[p].Close >= ema50[p] {
			pullbackIndex = p
		}
	}
	if pullbackIndex < 0 {
		return false, 0, 0, 0
	}
	if !(candles[d].Close > candles[d-1].High && candles[d].Close > ema20[d] && candles[d].Close > candles[d].Open) {
		return false, 0, 0, 0
	}
	lowestLow = math.Min(lowestLow, candles[d].Low)
	highestHigh = math.Max(highestHigh, candles[d].High)
	return true, pullbackIndex, lowestLow, highestHigh
}

func btc15MTrendPullbackShortSetup(candles []Candle, d int, ema20, ema50, ema200 []float64, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) (bool, int, float64, float64) {
	if !(ema20[d-1] < ema50[d-1] && ema50[d-1] < ema200[d-1]) {
		return false, 0, 0, 0
	}
	if !(ema50[d-1] < ema50[d-cfg.SlopeLookbackBars-1]) {
		return false, 0, 0, 0
	}
	if !(candles[d-1].Close < ema50[d-1]) {
		return false, 0, 0, 0
	}
	pullbackIndex := -1
	lowestLow := math.Inf(1)
	highestHigh := math.Inf(-1)
	start := d - cfg.PullbackWindowBars
	for p := start; p <= d-1; p++ {
		if !validNumber(ema20[p]) || !validNumber(ema50[p]) {
			return false, 0, 0, 0
		}
		lowestLow = math.Min(lowestLow, candles[p].Low)
		highestHigh = math.Max(highestHigh, candles[p].High)
		if candles[p].Close > ema50[p] {
			return false, 0, 0, 0
		}
		if candles[p].High >= ema20[p] && candles[p].Close <= ema50[p] {
			pullbackIndex = p
		}
	}
	if pullbackIndex < 0 {
		return false, 0, 0, 0
	}
	if !(candles[d].Close < candles[d-1].Low && candles[d].Close < ema20[d] && candles[d].Close < candles[d].Open) {
		return false, 0, 0, 0
	}
	lowestLow = math.Min(lowestLow, candles[d].Low)
	highestHigh = math.Max(highestHigh, candles[d].High)
	return true, pullbackIndex, lowestLow, highestHigh
}

func btc15MTrendPullbackSignalRow(candles []Candle, d int, sequence int, side Direction, pullbackIndex int, pullbackLowestLow float64, pullbackHighestHigh float64, ema20, ema50, ema200, atr []float64, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig, btCfg BacktestConfig, splits []Split) BTC15MTrendPullbackContinuationSignalRow {
	decision := candles[d]
	entryIndex := d + 1
	entry := candles[entryIndex]
	expectedEntry := applySlippage(entry.Open, btCfg.SlippagePct, side, true)
	stop := pullbackLowestLow - cfg.StopATRMultiple*atr[d]
	if side == Short {
		stop = pullbackHighestHigh + cfg.StopATRMultiple*atr[d]
	}
	stopDist := math.Abs(expectedEntry - stop)
	target := expectedEntry + cfg.TargetR*stopDist
	if side == Short {
		target = expectedEntry - cfg.TargetR*stopDist
	}
	row := BTC15MTrendPullbackContinuationSignalRow{
		SignalID:                  fmt.Sprintf("%s_%06d", BTC15MTrendPullbackContinuationCandidateID, sequence),
		CandidateID:               BTC15MTrendPullbackContinuationCandidateID,
		Timeframe:                 "15m",
		Split:                     splitNameForCloseTime(decision.CloseTime, splits),
		DecisionIndex:             d,
		DecisionOpenTime:          decision.OpenTime.UTC().Format(timeLayout),
		DecisionCloseTime:         decision.CloseTime.UTC().Format(timeLayout),
		DecisionOpen:              decision.Open,
		DecisionHigh:              decision.High,
		DecisionLow:               decision.Low,
		DecisionClose:             decision.Close,
		Side:                      side,
		TimingLabel:               "next_15m_open",
		EMA20:                     ema20[d],
		EMA50:                     ema50[d],
		EMA200:                    ema200[d],
		PriorEMA20:                ema20[d-1],
		PriorEMA50:                ema50[d-1],
		PriorEMA200:               ema200[d-1],
		EMA50SlopeLookbackBars:    cfg.SlopeLookbackBars,
		EMA50Slope:                ema50[d-1] - ema50[d-cfg.SlopeLookbackBars-1],
		PullbackWindowBars:        cfg.PullbackWindowBars,
		PullbackIndex:             pullbackIndex,
		PullbackLowestLow:         pullbackLowestLow,
		PullbackHighestHigh:       pullbackHighestHigh,
		ATR14:                     atr[d],
		EntryIndex:                entryIndex,
		EntryOpenTime:             entry.OpenTime.UTC().Format(timeLayout),
		EntryOpen:                 entry.Open,
		ExpectedEntryPrice:        expectedEntry,
		StopATRMultiple:           cfg.StopATRMultiple,
		TargetR:                   cfg.TargetR,
		Stop:                      stop,
		Target:                    target,
		MaxHoldBars:               cfg.MaxHoldBars,
		ForwardLabelsAsInput:      false,
		UsesFutureRowsForFeatures: false,
		DerivativesVetoUsed:       false,
		OptimizerSelectionUsed:    false,
	}
	if entry.Open <= 0 || expectedEntry <= 0 || stop <= 0 || target <= 0 || stopDist <= 0 {
		row.SkippedReason = "non_positive_trade_price_or_risk"
		return row
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: side, Stop: stop, Target: target}, expectedEntry)
	if !row.EntryGeometryValid {
		row.SkippedReason = "invalid_entry_geometry"
		return row
	}
	row.PrePositionCandidate = true
	return row
}

func btc15MTrendPullbackRun(candles []Candle, signals []BTC15MTrendPullbackContinuationSignalRow, skips *map[btc15MTrendPullbackSkipKey]int, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig, btCfg BacktestConfig) ([]BTC15MTrendPullbackContinuationSignalRow, []Trade) {
	byDecision := map[int]int{}
	for i, signal := range signals {
		if signal.PrePositionCandidate {
			byDecision[signal.DecisionIndex] = i
		}
	}
	balance := btCfg.StartBalance
	var trades []Trade
	var pos *Position
	var pending *Position
	for i, bar := range candles {
		if pending != nil && pending.OpenIndex == i {
			pos = pending
			pending = nil
		}
		if pos != nil && i >= pos.OpenIndex {
			if tr, ok := maybeExit(*pos, bar, i, btCfg); ok {
				trades = append(trades, tr)
				balance += tr.NetPnL
				pos = nil
			}
		}
		signalIdx, ok := byDecision[i]
		if !ok {
			continue
		}
		signal := signals[signalIdx]
		if signal.SkippedReason != "" {
			continue
		}
		if pos != nil || pending != nil {
			signals[signalIdx].SkippedReason = "open_position_already_active"
			btc15MTrendPullbackAddSkip(*skips, signal.Split, string(signal.Side), "open_position_already_active")
			continue
		}
		stopDist := math.Abs(signal.ExpectedEntryPrice - signal.Stop)
		size := positionSize(balance, signal.ExpectedEntryPrice, stopDist, btCfg)
		if size <= 0 {
			signals[signalIdx].SkippedReason = "non_positive_position_size"
			btc15MTrendPullbackAddSkip(*skips, signal.Split, string(signal.Side), "non_positive_position_size")
			continue
		}
		entryFee := signal.ExpectedEntryPrice * size * btCfg.FeePct
		entrySlip := math.Abs(signal.ExpectedEntryPrice-signal.EntryOpen) * size
		pending = &Position{
			Side:        signal.Side,
			EntryPrice:  signal.ExpectedEntryPrice,
			Stop:        signal.Stop,
			Target:      signal.Target,
			Size:        size,
			OpenIndex:   signal.EntryIndex,
			OpenTime:    signal.EntryOpenTime,
			EntryFee:    entryFee,
			EntrySlip:   entrySlip,
			MaxHoldBars: cfg.MaxHoldBars,
			Reason:      signal.SignalID,
		}
	}
	if pos != nil && len(candles) > 0 {
		lastIdx := len(candles) - 1
		last := candles[lastIdx]
		tr := closePosition(*pos, last.Close, last.Close, "force_close", lastIdx, last.CloseTime.Format(timeLayout), btCfg)
		trades = append(trades, tr)
	}
	executed := map[string]bool{}
	for _, trade := range trades {
		executed[trade.Signal] = true
	}
	for i := range signals {
		if executed[signals[i].SignalID] {
			signals[i].Executed = true
		}
	}
	return signals, trades
}

func btc15MTrendPullbackTradeRows(trades []Trade, signals []BTC15MTrendPullbackContinuationSignalRow, splits []Split) []BTC15MTrendPullbackContinuationTradeRow {
	signalByID := map[string]BTC15MTrendPullbackContinuationSignalRow{}
	for _, signal := range signals {
		signalByID[signal.SignalID] = signal
	}
	rows := make([]BTC15MTrendPullbackContinuationTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := signalByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := BTC15MTrendPullbackContinuationTradeRow{
			SignalID:      signal.SignalID,
			CandidateID:   signal.CandidateID,
			Timeframe:     signal.Timeframe,
			DecisionIndex: signal.DecisionIndex,
			EntrySplit:    splitNameForCloseTime(entryTime, splits),
			CloseSplit:    splitNameForCloseTime(exitTime, splits),
			Side:          trade.Side,
			EntryTime:     trade.EntryTime,
			ExitTime:      trade.ExitTime,
			OpenIndex:     trade.OpenIndex,
			CloseIndex:    trade.CloseIndex,
			EntryPrice:    trade.EntryPrice,
			ExitPrice:     trade.ExitPrice,
			Stop:          trade.Stop,
			Target:        trade.Target,
			Size:          trade.Size,
			InitialRisk:   initialRisk,
			GrossPnL:      trade.GrossPnL,
			NetPnL:        trade.NetPnL,
			Fees:          trade.Fees,
			Slippage:      trade.Slippage,
			ExitReason:    trade.Reason,
			HoldBars:      trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.NetR = trade.NetPnL / initialRisk
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].EntryTime != rows[j].EntryTime {
			return rows[i].EntryTime < rows[j].EntryTime
		}
		return rows[i].SignalID < rows[j].SignalID
	})
	return rows
}

func btc15MTrendPullbackEvaluate(source BTC15MTrendPullbackContinuationSourceRow, coverage BTC15MTrendPullbackContinuationCoverageRow, rows []SummaryRow, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) BTC15MTrendPullbackContinuationFalsification {
	byKey := map[string]SummaryRow{}
	for _, row := range rows {
		byKey[row.Split+"|"+row.Side] = row
	}
	full := byKey[fullSplitName+"|all"]
	stress := byKey["2021_2022_stress|all"]
	oos := byKey["2023_2024_oos|all"]
	recent := byKey["2025_2026_recent|all"]
	longFull := byKey[fullSplitName+"|"+string(Long)]
	shortFull := byKey[fullSplitName+"|"+string(Short)]
	minPrimary := minInt(stress.TotalTrades, minInt(oos.TotalTrades, recent.TotalTrades))
	dominantShare := 0.0
	if full.TotalTrades > 0 {
		dominantShare = float64(maxInt(stress.TotalTrades, maxInt(oos.TotalTrades, recent.TotalTrades))) / float64(full.TotalTrades)
	}
	tradeCountPass := full.TotalTrades >= cfg.MinFullTrades && minPrimary >= cfg.MinSplitTrades
	grossEdgePass := full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0
	netEdgePass := full.NetPnL > 0 && oos.NetPnL >= 0 && recent.NetPnL >= 0
	profitFactorPass := full.ProfitFactor >= cfg.FullProfitFactorMin && oos.ProfitFactor >= cfg.SplitProfitFactorMin && recent.ProfitFactor >= cfg.SplitProfitFactorMin
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit && stress.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && oos.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && recent.MaxDrawdown <= cfg.SplitMaxDrawdownLimit
	sideReportingPass := true
	for _, split := range append(DefaultSplits(), Split{Name: fullSplitName}) {
		_, hasAll := byKey[split.Name+"|all"]
		_, hasLong := byKey[split.Name+"|"+string(Long)]
		_, hasShort := byKey[split.Name+"|"+string(Short)]
		if !hasAll || !hasLong || !hasShort {
			sideReportingPass = false
		}
	}
	return BTC15MTrendPullbackContinuationFalsification{
		SourceResamplePass:                  source.SourceFactsPass && coverage.SourceResamplePass,
		LeakagePass:                         true,
		TradeCountPass:                      tradeCountPass,
		GrossEdgePass:                       grossEdgePass,
		NetEdgePass:                         netEdgePass,
		ProfitFactorPass:                    profitFactorPass,
		DrawdownPass:                        drawdownPass,
		RobustnessPass:                      tradeCountPass && dominantShare <= 0.90,
		SideReportingPass:                   sideReportingPass,
		OptimizerContaminationPass:          true,
		DerivativesVetoContaminationPass:    true,
		FullExecutedTrades:                  full.TotalTrades,
		RequiredFullExecutedTrades:          cfg.MinFullTrades,
		MinimumPrimarySplitExecutedTrades:   minPrimary,
		RequiredPrimarySplitExecutedTrades:  cfg.MinSplitTrades,
		FullGrossPnL:                        full.GrossPnL,
		FullNetPnL:                          full.NetPnL,
		FullProfitFactor:                    full.ProfitFactor,
		FullMaxDrawdown:                     full.MaxDrawdown,
		OOSGrossPnL:                         oos.GrossPnL,
		OOSNetPnL:                           oos.NetPnL,
		OOSProfitFactor:                     oos.ProfitFactor,
		RecentGrossPnL:                      recent.GrossPnL,
		RecentNetPnL:                        recent.NetPnL,
		RecentProfitFactor:                  recent.ProfitFactor,
		LongFullExecutedTrades:              longFull.TotalTrades,
		LongFullNetPnL:                      longFull.NetPnL,
		ShortFullExecutedTrades:             shortFull.TotalTrades,
		ShortFullNetPnL:                     shortFull.NetPnL,
		DominantPrimarySplitTradeShare:      dominantShare,
		DominantPrimarySplitTradeShareLimit: 0.90,
	}
}

func btc15MTrendPullbackFalsification(report BTC15MTrendPullbackContinuationFalsification, cfg BacktestFirstBTC15MTrendPullbackContinuationConfig) BTC15MTrendPullbackContinuationFalsification {
	report.BacktestName = BacktestFirstBTC15MTrendPullbackContinuationName
	report.CandidateID = BTC15MTrendPullbackContinuationCandidateID
	report.RequiredFullExecutedTrades = cfg.MinFullTrades
	report.RequiredPrimarySplitExecutedTrades = cfg.MinSplitTrades
	if report.DominantPrimarySplitTradeShareLimit == 0 {
		report.DominantPrimarySplitTradeShareLimit = 0.90
	}
	failures := []string{}
	if !report.SourceResamplePass {
		failures = append(failures, "source_or_resample_failed")
	}
	if !report.LeakagePass {
		failures = append(failures, "leakage")
	}
	if !report.TradeCountPass {
		failures = append(failures, "insufficient_executed_trades")
	}
	if !report.GrossEdgePass {
		failures = append(failures, "gross_edge_gate_failed")
	}
	if !report.NetEdgePass {
		failures = append(failures, "net_edge_gate_failed")
	}
	if !report.ProfitFactorPass {
		failures = append(failures, "profit_factor_gate_failed")
	}
	if !report.DrawdownPass {
		failures = append(failures, "drawdown_gate_failed")
	}
	if !report.RobustnessPass {
		failures = append(failures, "robustness_gate_failed")
	}
	if !report.SideReportingPass {
		failures = append(failures, "side_reporting_failed")
	}
	if !report.OptimizerContaminationPass {
		failures = append(failures, "optimizer_contamination")
	}
	if !report.DerivativesVetoContaminationPass {
		failures = append(failures, "derivatives_veto_contamination")
	}
	report.FailureReasons = failures
	switch {
	case !report.SourceResamplePass || !report.LeakagePass:
		report.StopState = BTC15MTrendPullbackContinuationStopStateFailedSourceOrResample
	case report.TradeCountPass &&
		report.GrossEdgePass &&
		report.NetEdgePass &&
		report.ProfitFactorPass &&
		report.DrawdownPass &&
		report.RobustnessPass &&
		report.SideReportingPass &&
		report.OptimizerContaminationPass &&
		report.DerivativesVetoContaminationPass:
		report.StopState = BTC15MTrendPullbackContinuationStopStatePassedNeedsReview
	default:
		report.StopState = BTC15MTrendPullbackContinuationStopStateFailedNoUsableStrategy
	}
	return report
}

func btc15MTrendPullbackEMA(candles []Candle, period int) []float64 {
	out := nanSlice(len(candles))
	if period <= 0 || len(candles) < period {
		return out
	}
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += candles[i].Close
	}
	out[period-1] = sum / float64(period)
	alpha := 2.0 / float64(period+1)
	for i := period; i < len(candles); i++ {
		out[i] = alpha*candles[i].Close + (1-alpha)*out[i-1]
	}
	return out
}

func btc15MTrendPullbackAddSkip(skips map[btc15MTrendPullbackSkipKey]int, split, side, reason string) {
	if split == "" {
		split = fullSplitName
	}
	if side == "" {
		side = "all"
	}
	skips[btc15MTrendPullbackSkipKey{split: split, side: side, reason: reason}]++
	if split != fullSplitName {
		skips[btc15MTrendPullbackSkipKey{split: fullSplitName, side: side, reason: reason}]++
	}
}

func btc15MTrendPullbackSkipRows(skips map[btc15MTrendPullbackSkipKey]int) []BTC15MTrendPullbackContinuationSkipRow {
	keys := make([]btc15MTrendPullbackSkipKey, 0, len(skips))
	for key := range skips {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if splitSortKey(keys[i].split) != splitSortKey(keys[j].split) {
			return splitSortKey(keys[i].split) < splitSortKey(keys[j].split)
		}
		if keys[i].side != keys[j].side {
			return keys[i].side < keys[j].side
		}
		return keys[i].reason < keys[j].reason
	})
	rows := make([]BTC15MTrendPullbackContinuationSkipRow, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, BTC15MTrendPullbackContinuationSkipRow{
			CandidateID:       BTC15MTrendPullbackContinuationCandidateID,
			Split:             key.split,
			Side:              key.side,
			Reason:            key.reason,
			Count:             skips[key],
			MissingDataPolicy: "skip rows; no fill/interpolation/nearest-future",
			ForwardFilledRows: 0,
		})
	}
	return rows
}
