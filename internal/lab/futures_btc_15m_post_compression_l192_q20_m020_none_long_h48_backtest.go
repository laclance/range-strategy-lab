package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestName = "futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest"
	BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID         = "btc_15m_post_compression_l192_q20_m020_none_long_h48_v1"

	BTC15MPostCompressionFixedBacktestStopStatePassedNeedsReview       = "post_compression_directional_expansion_backtest_passed_needs_review"
	BTC15MPostCompressionFixedBacktestStopStateFailedNoUsableStrategy  = "post_compression_directional_expansion_backtest_failed_no_usable_strategy"
	BTC15MPostCompressionFixedBacktestStopStateFailedSourceOrResample  = "post_compression_directional_expansion_backtest_failed_source_or_resample"
	BTC15MPostCompressionFixedBacktestStopStateRejectedOptimizerContam = "post_compression_directional_expansion_backtest_rejected_optimizer_contamination"
	BTC15MPostCompressionFixedBacktestStopStateRejectedClosedReslice   = "post_compression_directional_expansion_backtest_rejected_closed_family_reslice"
	BTC15MPostCompressionFixedBacktestStopStateRejectedVetoContam      = "post_compression_directional_expansion_backtest_rejected_veto_contamination"
	BTC15MPostCompressionFixedBacktestExpectedRawCandidateRows         = 468
	BTC15MPostCompressionFixedBacktestExitReasonConcentrationLimit     = 0.95
	BTC15MPostCompressionFixedBacktestPrimarySplitTradeShareLimit      = 0.90
)

type FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig struct {
	ApprovedSourcePath         string
	ExpectedSourceRows         int
	ExpectedFirstOpenTime      string
	ExpectedLastOpenTime       string
	ExpectedGapCount           int
	ExpectedDuplicateCount     int
	ExpectedZeroVolumeCount    int
	SkipSourceFactCheck        bool
	Expected15MRows            int
	Expected15MLastOpenTime    string
	SkipCoverageCountCheck     bool
	LookbackBars               int
	CompressionPercentile      float64
	PercentileReferenceBars    int
	BreakoutATRMultiple        float64
	ATRPeriod                  int
	MaxHoldBars                int
	StopATRMultiple            float64
	TargetATRMultiple          float64
	ExpectedRawCandidateRows   int
	SkipCandidateIdentityGate  bool
	MinFullTrades              int
	MinSplitTrades             int
	FullStressProfitFactorMin  float64
	SplitStressProfitFactorMin float64
	FullMaxDrawdownLimit       float64
	SplitMaxDrawdownLimit      float64
	ForwardLabelLeakOverride   bool
	OptimizerContamination     bool
	ClosedFamilyReslice        bool
	DerivativesVetoContam      bool
}

type FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestResult struct {
	SourceRows     []BTC15MPostCompressionFixedBacktestSourceRow     `json:"source_rows"`
	CoverageRows   []BTC15MPostCompressionFixedBacktestCoverageRow   `json:"coverage_rows"`
	SignalRows     []BTC15MPostCompressionFixedBacktestSignalRow     `json:"signal_rows"`
	SkipRows       []BTC15MPostCompressionFixedBacktestSkipRow       `json:"skip_rows"`
	TradeRows      []BTC15MPostCompressionFixedBacktestTradeRow      `json:"trade_rows"`
	SummaryRows    []BTC15MPostCompressionFixedBacktestSummaryRow    `json:"summary_rows"`
	CostStressRows []BTC15MPostCompressionFixedBacktestCostStressRow `json:"cost_stress_rows"`
	Falsification  BTC15MPostCompressionFixedBacktestFalsification   `json:"falsification"`
	Trades         []Trade                                           `json:"trades"`
	StopState      string                                            `json:"stop_state"`
}

type BTC15MPostCompressionFixedBacktestSourceRow struct {
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
	ExpectedGapCount           int    `json:"expected_gap_count"`
	DuplicateCount             int    `json:"duplicate_count"`
	ExpectedDuplicateCount     int    `json:"expected_duplicate_count"`
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

type BTC15MPostCompressionFixedBacktestCoverageRow struct {
	BacktestName          string `json:"backtest_name"`
	CandidateID           string `json:"candidate_id"`
	Timeframe             string `json:"timeframe"`
	IntervalMinutes       int    `json:"interval_minutes"`
	ChildBars             int    `json:"child_bars"`
	BarsPerDay            int    `json:"bars_per_day"`
	RowCount              int    `json:"row_count"`
	ExpectedRowCount      int    `json:"expected_row_count"`
	FirstOpenTime         string `json:"first_open_time"`
	LastOpenTime          string `json:"last_open_time"`
	ExpectedLastOpenTime  string `json:"expected_last_open_time"`
	FirstCloseTime        string `json:"first_close_time"`
	LastCloseTime         string `json:"last_close_time"`
	ExpectedChildBars     int    `json:"expected_child_bars"`
	CompleteBucketCount   int    `json:"complete_bucket_count"`
	PartialFinalChildBars int    `json:"partial_final_child_bars"`
	PartialFinalDropped   bool   `json:"partial_final_dropped"`
	GapCount              int    `json:"gap_count"`
	DuplicateCount        int    `json:"duplicate_count"`
	MissingChildOpenCount int    `json:"missing_child_open_count"`
	Complete              bool   `json:"complete"`
	ClosedCandleOnly      bool   `json:"closed_candle_only"`
	SourceResamplePass    bool   `json:"source_resample_pass"`
	ValidationStatus      string `json:"validation_status"`
	ValidationError       string `json:"validation_error,omitempty"`
}

type BTC15MPostCompressionFixedBacktestSignalRow struct {
	SignalID                  string    `json:"signal_id"`
	CandidateID               string    `json:"candidate_id"`
	Timeframe                 string    `json:"timeframe"`
	Split                     string    `json:"split"`
	DecisionIndex             int       `json:"decision_index"`
	DecisionOpenTime          string    `json:"decision_open_time"`
	DecisionCloseTime         string    `json:"decision_close_time"`
	DecisionClose             float64   `json:"decision_close"`
	Side                      Direction `json:"side"`
	TimingLabel               string    `json:"timing_label"`
	LookbackBars              int       `json:"lookback_bars"`
	CompressionPercentile     float64   `json:"compression_percentile"`
	PercentileReferenceBars   int       `json:"percentile_reference_bars"`
	CompressionThreshold      float64   `json:"compression_threshold"`
	RangeHigh                 float64   `json:"range_high"`
	RangeLow                  float64   `json:"range_low"`
	RangeWidthPct             float64   `json:"range_width_pct"`
	BreakoutATRMultiple       float64   `json:"breakout_atr_multiple"`
	PriorATR14                float64   `json:"prior_atr14"`
	BreakoutDistance          float64   `json:"breakout_distance"`
	VolumeMode                string    `json:"volume_mode"`
	EntryIndex                int       `json:"entry_index"`
	EntryOpenTime             string    `json:"entry_open_time"`
	EntryOpen                 float64   `json:"entry_open"`
	ExpectedEntryPrice        float64   `json:"expected_entry_price"`
	StopATRMultiple           float64   `json:"stop_atr_multiple"`
	TargetATRMultiple         float64   `json:"target_atr_multiple"`
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

type BTC15MPostCompressionFixedBacktestSkipRow struct {
	CandidateID       string `json:"candidate_id"`
	Split             string `json:"split"`
	Side              string `json:"side"`
	Reason            string `json:"reason"`
	Count             int    `json:"count"`
	MissingDataPolicy string `json:"missing_data_policy"`
	ForwardFilledRows int    `json:"forward_filled_rows"`
}

type BTC15MPostCompressionFixedBacktestTradeRow struct {
	SignalID                  string    `json:"signal_id"`
	CandidateID               string    `json:"candidate_id"`
	Timeframe                 string    `json:"timeframe"`
	DecisionIndex             int       `json:"decision_index"`
	DecisionCloseTime         string    `json:"decision_close_time"`
	EntrySplit                string    `json:"entry_split"`
	CloseSplit                string    `json:"close_split"`
	Side                      Direction `json:"side"`
	EntryTime                 string    `json:"entry_time"`
	ExitTime                  string    `json:"exit_time"`
	OpenIndex                 int       `json:"open_index"`
	CloseIndex                int       `json:"close_index"`
	EntryPrice                float64   `json:"entry_price"`
	ExitPrice                 float64   `json:"exit_price"`
	Stop                      float64   `json:"stop"`
	Target                    float64   `json:"target"`
	Size                      float64   `json:"size"`
	InitialRisk               float64   `json:"initial_risk"`
	GrossPnL                  float64   `json:"gross_pnl"`
	EngineNetPnL              float64   `json:"engine_net_pnl"`
	ExtraSlippageStressNetPnL float64   `json:"extra_slippage_stress_net_pnl"`
	Fees                      float64   `json:"fees"`
	Slippage                  float64   `json:"slippage"`
	GrossR                    float64   `json:"gross_r"`
	EngineNetR                float64   `json:"engine_net_r"`
	ExtraSlippageStressNetR   float64   `json:"extra_slippage_stress_net_r"`
	ExitReason                string    `json:"exit_reason"`
	HoldBars                  int       `json:"hold_bars"`
}

type BTC15MPostCompressionFixedBacktestSummaryRow struct {
	CandidateID               string  `json:"candidate_id"`
	Split                     string  `json:"split"`
	Side                      string  `json:"side"`
	RawCandidateRows          int     `json:"raw_candidate_rows"`
	SkippedRows               int     `json:"skipped_rows"`
	ExecutedTrades            int     `json:"executed_trades"`
	Wins                      int     `json:"wins"`
	Losses                    int     `json:"losses"`
	WinRate                   float64 `json:"win_rate"`
	GrossPnL                  float64 `json:"gross_pnl"`
	EngineNetPnL              float64 `json:"engine_net_pnl"`
	ExtraSlippageStressNetPnL float64 `json:"extra_slippage_stress_net_pnl"`
	Fees                      float64 `json:"fees"`
	Slippage                  float64 `json:"slippage"`
	StressProfitFactor        float64 `json:"stress_profit_factor"`
	GrossProfitFactor         float64 `json:"gross_profit_factor"`
	MaxDrawdown               float64 `json:"max_drawdown"`
	AvgGrossR                 float64 `json:"avg_gross_r"`
	AvgEngineNetR             float64 `json:"avg_engine_net_r"`
	AvgExtraSlippageStressR   float64 `json:"avg_extra_slippage_stress_r"`
	AvgInitialRisk            float64 `json:"avg_initial_risk"`
	AvgHoldBars               float64 `json:"avg_hold_bars"`
	StopLossExits             int     `json:"stop_loss_exits"`
	TakeProfitExits           int     `json:"take_profit_exits"`
	TimeStopExits             int     `json:"time_stop_exits"`
	ForceCloseExits           int     `json:"force_close_exits"`
	TradeCountGatePass        bool    `json:"trade_count_gate_pass"`
	GrossEdgeGatePass         bool    `json:"gross_edge_gate_pass"`
	CostedEdgeGatePass        bool    `json:"costed_edge_gate_pass"`
	ProfitFactorGatePass      bool    `json:"profit_factor_gate_pass"`
	DrawdownGatePass          bool    `json:"drawdown_gate_pass"`
	RobustnessGatePass        bool    `json:"robustness_gate_pass"`
}

type BTC15MPostCompressionFixedBacktestCostStressRow struct {
	CandidateID               string  `json:"candidate_id"`
	Split                     string  `json:"split"`
	Side                      string  `json:"side"`
	ExecutedTrades            int     `json:"executed_trades"`
	GrossPnL                  float64 `json:"gross_pnl"`
	EngineNetPnL              float64 `json:"engine_net_pnl"`
	Slippage                  float64 `json:"slippage"`
	ExtraSlippageStressNetPnL float64 `json:"extra_slippage_stress_net_pnl"`
	Fees                      float64 `json:"fees"`
	StressProfitFactor        float64 `json:"stress_profit_factor"`
	StressMaxDrawdown         float64 `json:"stress_max_drawdown"`
	PassFailUsesStressNet     bool    `json:"pass_fail_uses_stress_net"`
}

type BTC15MPostCompressionFixedBacktestFalsification struct {
	BacktestName                       string   `json:"backtest_name"`
	CandidateID                        string   `json:"candidate_id"`
	StopState                          string   `json:"stop_state"`
	SourceResamplePass                 bool     `json:"source_resample_pass"`
	CandidateIdentityPass              bool     `json:"candidate_identity_pass"`
	LeakagePass                        bool     `json:"leakage_pass"`
	TradeCountPass                     bool     `json:"trade_count_pass"`
	GrossEdgePass                      bool     `json:"gross_edge_pass"`
	CostedEdgePass                     bool     `json:"costed_edge_pass"`
	ProfitFactorPass                   bool     `json:"profit_factor_pass"`
	DrawdownPass                       bool     `json:"drawdown_pass"`
	RobustnessPass                     bool     `json:"robustness_pass"`
	OptimizerContaminationPass         bool     `json:"optimizer_contamination_pass"`
	ClosedFamilyProtectionPass         bool     `json:"closed_family_protection_pass"`
	DerivativesVetoContaminationPass   bool     `json:"derivatives_veto_contamination_pass"`
	CommonOutputsTradeCompatible       bool     `json:"common_outputs_trade_compatible"`
	ExpectedRawCandidateRows           int      `json:"expected_raw_candidate_rows"`
	FullRawCandidateRows               int      `json:"full_raw_candidate_rows"`
	FullExecutedTrades                 int      `json:"full_executed_trades"`
	RequiredFullExecutedTrades         int      `json:"required_full_executed_trades"`
	MinimumPrimarySplitExecutedTrades  int      `json:"minimum_primary_split_executed_trades"`
	RequiredPrimarySplitExecutedTrades int      `json:"required_primary_split_executed_trades"`
	FullGrossPnL                       float64  `json:"full_gross_pnl"`
	FullEngineNetPnL                   float64  `json:"full_engine_net_pnl"`
	FullExtraSlippageStressNetPnL      float64  `json:"full_extra_slippage_stress_net_pnl"`
	FullStressProfitFactor             float64  `json:"full_stress_profit_factor"`
	FullMaxDrawdown                    float64  `json:"full_max_drawdown"`
	DominantExitReasonShare            float64  `json:"dominant_exit_reason_share"`
	DominantPrimarySplitTradeShare     float64  `json:"dominant_primary_split_trade_share"`
	ExitReasonConcentrationLimit       float64  `json:"exit_reason_concentration_limit"`
	PrimarySplitTradeShareLimit        float64  `json:"primary_split_trade_share_limit"`
	FailureReasons                     []string `json:"failure_reasons,omitempty"`
}

type btcPostCompressionFixedSkipKey struct {
	split  string
	side   string
	reason string
}

func DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig() FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig {
	return FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig{
		ApprovedSourcePath:         btc15MPostCompressionSourcePath,
		ExpectedSourceRows:         btc15MPostCompressionExpectedRows,
		ExpectedFirstOpenTime:      btc15MPostCompressionExpectedFirst,
		ExpectedLastOpenTime:       btc15MPostCompressionExpectedLast,
		ExpectedGapCount:           0,
		ExpectedDuplicateCount:     0,
		ExpectedZeroVolumeCount:    btc15MPostCompressionExpectedZeroVol,
		Expected15MRows:            btc15MPostCompressionExpected15MRows,
		Expected15MLastOpenTime:    btc15MPostCompressionExpected15MLastOpen,
		LookbackBars:               192,
		CompressionPercentile:      0.20,
		PercentileReferenceBars:    1920,
		BreakoutATRMultiple:        0.20,
		ATRPeriod:                  14,
		MaxHoldBars:                48,
		StopATRMultiple:            1.0,
		TargetATRMultiple:          2.0,
		ExpectedRawCandidateRows:   BTC15MPostCompressionFixedBacktestExpectedRawCandidateRows,
		MinFullTrades:              120,
		MinSplitTrades:             25,
		FullStressProfitFactorMin:  1.20,
		SplitStressProfitFactorMin: 1.05,
		FullMaxDrawdownLimit:       0.25,
		SplitMaxDrawdownLimit:      0.30,
	}
}

func RunFuturesBTC15MPostCompressionL192Q20M020NoneLongH48Backtest(candles []Candle, manifest SourceManifest, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, btCfg BacktestConfig, splits []Split) (FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	sourceRow := btcPostCompressionFixedBacktestSourceRow(manifest, cfg)
	result := FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestResult{
		SourceRows: []BTC15MPostCompressionFixedBacktestSourceRow{sourceRow},
	}
	frame := btcPostCompression15MFrame()
	resampled, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
	coverageRow := btcPostCompressionFixedBacktestCoverageRow(coverage, cfg, sourceRow.SourceFactsPass)
	if err != nil {
		coverageRow.ValidationStatus = "rejected"
		coverageRow.ValidationError = err.Error()
		coverageRow.SourceResamplePass = false
	}
	result.CoverageRows = []BTC15MPostCompressionFixedBacktestCoverageRow{coverageRow}
	if !sourceRow.SourceFactsPass || !coverageRow.SourceResamplePass {
		result.Falsification = btcPostCompressionFixedBacktestFalsification(BTC15MPostCompressionFixedBacktestFalsification{
			SourceResamplePass:               false,
			CandidateIdentityPass:            false,
			LeakagePass:                      !cfg.ForwardLabelLeakOverride,
			OptimizerContaminationPass:       !cfg.OptimizerContamination,
			ClosedFamilyProtectionPass:       !cfg.ClosedFamilyReslice,
			DerivativesVetoContaminationPass: !cfg.DerivativesVetoContam,
			CommonOutputsTradeCompatible:     true,
			ExpectedRawCandidateRows:         cfg.ExpectedRawCandidateRows,
		})
		result.StopState = result.Falsification.StopState
		return result, nil
	}

	signals, skipCounts := btcPostCompressionFixedBacktestBuildSignals(resampled, cfg, btCfg, splits)
	signals, trades := btcPostCompressionFixedBacktestRun(resampled, signals, &skipCounts, cfg, btCfg, splits)
	tradeRows := btcPostCompressionFixedBacktestTradeRows(trades, signals, splits)
	skipRows := btcPostCompressionFixedBacktestSkipRows(skipCounts)
	summaryRows := btcPostCompressionFixedBacktestSummaryRows(signals, skipRows, tradeRows, cfg, btCfg.StartBalance, splits)
	costStressRows := btcPostCompressionFixedBacktestCostStressRows(summaryRows)
	report := btcPostCompressionFixedBacktestEvaluate(sourceRow, coverageRow, signals, summaryRows, cfg)

	result.SignalRows = signals
	result.SkipRows = skipRows
	result.Trades = trades
	result.TradeRows = tradeRows
	result.SummaryRows = summaryRows
	result.CostStressRows = costStressRows
	result.Falsification = btcPostCompressionFixedBacktestFalsification(report)
	result.StopState = result.Falsification.StopState
	return result, nil
}

func (cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig) withDefaults() FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig {
	defaults := DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
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
	if cfg.CompressionPercentile == 0 {
		cfg.CompressionPercentile = defaults.CompressionPercentile
	}
	if cfg.PercentileReferenceBars == 0 {
		cfg.PercentileReferenceBars = defaults.PercentileReferenceBars
	}
	if cfg.BreakoutATRMultiple == 0 {
		cfg.BreakoutATRMultiple = defaults.BreakoutATRMultiple
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = defaults.MaxHoldBars
	}
	if cfg.StopATRMultiple == 0 {
		cfg.StopATRMultiple = defaults.StopATRMultiple
	}
	if cfg.TargetATRMultiple == 0 {
		cfg.TargetATRMultiple = defaults.TargetATRMultiple
	}
	if cfg.ExpectedRawCandidateRows == 0 {
		cfg.ExpectedRawCandidateRows = defaults.ExpectedRawCandidateRows
	}
	if cfg.MinFullTrades == 0 {
		cfg.MinFullTrades = defaults.MinFullTrades
	}
	if cfg.MinSplitTrades == 0 {
		cfg.MinSplitTrades = defaults.MinSplitTrades
	}
	if cfg.FullStressProfitFactorMin == 0 {
		cfg.FullStressProfitFactorMin = defaults.FullStressProfitFactorMin
	}
	if cfg.SplitStressProfitFactorMin == 0 {
		cfg.SplitStressProfitFactorMin = defaults.SplitStressProfitFactorMin
	}
	if cfg.FullMaxDrawdownLimit == 0 {
		cfg.FullMaxDrawdownLimit = defaults.FullMaxDrawdownLimit
	}
	if cfg.SplitMaxDrawdownLimit == 0 {
		cfg.SplitMaxDrawdownLimit = defaults.SplitMaxDrawdownLimit
	}
	return cfg
}

func (cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig) validate() error {
	if cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("post-compression fixed backtest approved source path is required")
	}
	if cfg.LookbackBars <= 0 || cfg.PercentileReferenceBars <= 0 || cfg.ATRPeriod <= 0 || cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("post-compression fixed backtest lookback/reference/ATR/hold bars must be positive")
	}
	if cfg.CompressionPercentile <= 0 || cfg.CompressionPercentile >= 1 {
		return fmt.Errorf("post-compression fixed backtest compression percentile must be between 0 and 1")
	}
	if cfg.BreakoutATRMultiple <= 0 || cfg.StopATRMultiple <= 0 || cfg.TargetATRMultiple <= 0 {
		return fmt.Errorf("post-compression fixed backtest ATR multiples must be positive")
	}
	if cfg.ExpectedRawCandidateRows <= 0 {
		return fmt.Errorf("post-compression fixed backtest expected candidate rows must be positive")
	}
	if cfg.MinFullTrades <= 0 || cfg.MinSplitTrades <= 0 {
		return fmt.Errorf("post-compression fixed backtest trade-count gates must be positive")
	}
	if cfg.FullStressProfitFactorMin <= 0 || cfg.SplitStressProfitFactorMin <= 0 {
		return fmt.Errorf("post-compression fixed backtest profit-factor gates must be positive")
	}
	if cfg.FullMaxDrawdownLimit <= 0 || cfg.SplitMaxDrawdownLimit <= 0 {
		return fmt.Errorf("post-compression fixed backtest drawdown limits must be positive")
	}
	return nil
}

func btcPostCompressionFixedBacktestSourceRow(manifest SourceManifest, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig) BTC15MPostCompressionFixedBacktestSourceRow {
	row := BTC15MPostCompressionFixedBacktestSourceRow{
		BacktestName:               FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestName,
		CandidateID:                BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
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
		ExpectedGapCount:           cfg.ExpectedGapCount,
		DuplicateCount:             manifest.DuplicateCount,
		ExpectedDuplicateCount:     cfg.ExpectedDuplicateCount,
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
	if manifest.Symbol != RangeUniverseSymbolBTCUSDT || manifest.Interval != "5m" {
		failures = append(failures, "source must be BTCUSDT 5m")
	}
	if !sameCleanPath(manifest.Path, cfg.ApprovedSourcePath) {
		failures = append(failures, fmt.Sprintf("source path %q is not approved path %q", manifest.Path, cfg.ApprovedSourcePath))
	}
	if !cfg.SkipSourceFactCheck {
		if manifest.RowCount != cfg.ExpectedSourceRows {
			failures = append(failures, fmt.Sprintf("row_count=%d expected=%d", manifest.RowCount, cfg.ExpectedSourceRows))
		}
		if manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime {
			failures = append(failures, fmt.Sprintf("first_open_time=%s expected=%s", manifest.FirstOpenTime, cfg.ExpectedFirstOpenTime))
		}
		if manifest.LastOpenTime != cfg.ExpectedLastOpenTime {
			failures = append(failures, fmt.Sprintf("last_open_time=%s expected=%s", manifest.LastOpenTime, cfg.ExpectedLastOpenTime))
		}
		if manifest.GapCount != cfg.ExpectedGapCount {
			failures = append(failures, fmt.Sprintf("gap_count=%d expected=%d", manifest.GapCount, cfg.ExpectedGapCount))
		}
		if manifest.DuplicateCount != cfg.ExpectedDuplicateCount {
			failures = append(failures, fmt.Sprintf("duplicate_count=%d expected=%d", manifest.DuplicateCount, cfg.ExpectedDuplicateCount))
		}
		if manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
			failures = append(failures, fmt.Sprintf("zero_volume_count=%d expected=%d", manifest.ZeroVolumeCount, cfg.ExpectedZeroVolumeCount))
		}
	}
	row.SourceFactsPass = len(failures) == 0
	if !row.SourceFactsPass {
		row.ValidationStatus = "rejected"
		row.ValidationError = strings.Join(failures, "; ")
	}
	return row
}

func btcPostCompressionFixedBacktestCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, sourcePass bool) BTC15MPostCompressionFixedBacktestCoverageRow {
	row := BTC15MPostCompressionFixedBacktestCoverageRow{
		BacktestName:          FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestName,
		CandidateID:           BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
		Timeframe:             base.Timeframe,
		IntervalMinutes:       base.IntervalMinutes,
		ChildBars:             base.ChildBars,
		BarsPerDay:            base.BarsPerDay,
		RowCount:              base.RowCount,
		ExpectedRowCount:      cfg.Expected15MRows,
		FirstOpenTime:         base.FirstOpenTime,
		LastOpenTime:          base.LastOpenTime,
		ExpectedLastOpenTime:  cfg.Expected15MLastOpenTime,
		FirstCloseTime:        base.FirstCloseTime,
		LastCloseTime:         base.LastCloseTime,
		ExpectedChildBars:     base.ExpectedChildBars,
		CompleteBucketCount:   base.CompleteBucketCount,
		PartialFinalChildBars: base.PartialFinalChildBars,
		PartialFinalDropped:   base.PartialFinalDropped,
		GapCount:              base.GapCount,
		DuplicateCount:        base.DuplicateCount,
		MissingChildOpenCount: base.MissingChildOpenCount,
		Complete:              base.Complete,
		ClosedCandleOnly:      true,
		ValidationStatus:      base.ValidationStatus,
		ValidationError:       base.ValidationError,
	}
	pass := sourcePass && base.ValidationStatus == "accepted" && base.Complete && base.Timeframe == RangeDiscoveryTimeframe15m && base.ChildBars == 3
	if !cfg.SkipCoverageCountCheck {
		pass = pass && base.RowCount == cfg.Expected15MRows && base.LastOpenTime == cfg.Expected15MLastOpenTime
	}
	row.SourceResamplePass = pass
	if !pass && row.ValidationStatus == "accepted" {
		row.ValidationStatus = "rejected"
		reasons := []string{}
		if !sourcePass {
			reasons = append(reasons, "source facts failed")
		}
		if !base.Complete || base.Timeframe != RangeDiscoveryTimeframe15m || base.ChildBars != 3 {
			reasons = append(reasons, "15m resample coverage failed")
		}
		if !cfg.SkipCoverageCountCheck && base.RowCount != cfg.Expected15MRows {
			reasons = append(reasons, fmt.Sprintf("15m row_count=%d expected=%d", base.RowCount, cfg.Expected15MRows))
		}
		if !cfg.SkipCoverageCountCheck && base.LastOpenTime != cfg.Expected15MLastOpenTime {
			reasons = append(reasons, fmt.Sprintf("15m last_open_time=%s expected=%s", base.LastOpenTime, cfg.Expected15MLastOpenTime))
		}
		row.ValidationError = strings.Join(reasons, "; ")
	}
	return row
}

func btcPostCompressionFixedBacktestBuildSignals(candles []Candle, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, btCfg BacktestConfig, splits []Split) ([]BTC15MPostCompressionFixedBacktestSignalRow, map[btcPostCompressionFixedSkipKey]int) {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	facts := btcPostCompressionBuildRangeFacts(candles, cfg.LookbackBars, []float64{cfg.CompressionPercentile}, cfg.PercentileReferenceBars)
	thresholds := facts.thresholds[cfg.CompressionPercentile]
	atr := ATR(candles, cfg.ATRPeriod)
	warmup := cfg.LookbackBars + cfg.PercentileReferenceBars
	skips := map[btcPostCompressionFixedSkipKey]int{}
	signals := []BTC15MPostCompressionFixedBacktestSignalRow{}
	for d := 0; d < len(candles); d++ {
		split := splitNameForCloseTime(candles[d].CloseTime, splits)
		if d < warmup {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_warmup")
			continue
		}
		if d >= len(facts.widthPct) || !validNumber(facts.widthPct[d]) || !validNumber(facts.high[d]) || !validNumber(facts.low[d]) {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_prior_range")
			continue
		}
		if d >= len(thresholds) || !validNumber(thresholds[d]) {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_compression_reference")
			continue
		}
		if d <= 0 || d-1 >= len(atr) || !validNumber(atr[d-1]) || atr[d-1] <= 0 {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_prior_atr")
			continue
		}
		if facts.widthPct[d] > thresholds[d] {
			continue
		}
		if candles[d].Close < facts.high[d]+cfg.BreakoutATRMultiple*atr[d-1] {
			continue
		}
		entryIndex := d + 1
		if entryIndex >= len(candles) {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_entry_candle")
			continue
		}
		if entryIndex+cfg.MaxHoldBars >= len(candles) {
			btcPostCompressionFixedAddSkip(skips, split, string(Long), "missing_future_path")
			continue
		}
		row := btcPostCompressionFixedSignalRow(candles, d, len(signals)+1, facts, thresholds[d], atr[d-1], cfg, btCfg, splits)
		if row.SkippedReason != "" {
			btcPostCompressionFixedAddSkip(skips, row.Split, string(row.Side), row.SkippedReason)
		}
		signals = append(signals, row)
	}
	return signals, skips
}

func btcPostCompressionFixedSignalRow(candles []Candle, d int, sequence int, facts btcPostCompressionRangeFacts, threshold, priorATR float64, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, btCfg BacktestConfig, splits []Split) BTC15MPostCompressionFixedBacktestSignalRow {
	decision := candles[d]
	entryIndex := d + 1
	entry := candles[entryIndex]
	expectedEntry := applySlippage(entry.Open, btCfg.SlippagePct, Long, true)
	stop := expectedEntry - cfg.StopATRMultiple*priorATR
	target := expectedEntry + cfg.TargetATRMultiple*priorATR
	row := BTC15MPostCompressionFixedBacktestSignalRow{
		SignalID:                  fmt.Sprintf("%s_%06d", BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID, sequence),
		CandidateID:               BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
		Timeframe:                 RangeDiscoveryTimeframe15m,
		Split:                     splitNameForCloseTime(decision.CloseTime, splits),
		DecisionIndex:             d,
		DecisionOpenTime:          decision.OpenTime.UTC().Format(timeLayout),
		DecisionCloseTime:         decision.CloseTime.UTC().Format(timeLayout),
		DecisionClose:             decision.Close,
		Side:                      Long,
		TimingLabel:               "next_15m_open",
		LookbackBars:              cfg.LookbackBars,
		CompressionPercentile:     cfg.CompressionPercentile,
		PercentileReferenceBars:   cfg.PercentileReferenceBars,
		CompressionThreshold:      threshold,
		RangeHigh:                 facts.high[d],
		RangeLow:                  facts.low[d],
		RangeWidthPct:             facts.widthPct[d],
		BreakoutATRMultiple:       cfg.BreakoutATRMultiple,
		PriorATR14:                priorATR,
		BreakoutDistance:          decision.Close - facts.high[d],
		VolumeMode:                BTC15MPostCompressionVolumeNone,
		EntryIndex:                entryIndex,
		EntryOpenTime:             entry.OpenTime.UTC().Format(timeLayout),
		EntryOpen:                 entry.Open,
		ExpectedEntryPrice:        expectedEntry,
		StopATRMultiple:           cfg.StopATRMultiple,
		TargetATRMultiple:         cfg.TargetATRMultiple,
		Stop:                      stop,
		Target:                    target,
		MaxHoldBars:               cfg.MaxHoldBars,
		ForwardLabelsAsInput:      false,
		UsesFutureRowsForFeatures: false,
		DerivativesVetoUsed:       false,
		OptimizerSelectionUsed:    false,
	}
	if entry.Open <= 0 || expectedEntry <= 0 || stop <= 0 || target <= 0 {
		row.SkippedReason = "non_positive_trade_price"
		return row
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: Long, Stop: stop, Target: target}, expectedEntry)
	if !row.EntryGeometryValid {
		row.SkippedReason = "invalid_entry_geometry"
		return row
	}
	row.PrePositionCandidate = true
	return row
}

func btcPostCompressionFixedBacktestRun(candles []Candle, signals []BTC15MPostCompressionFixedBacktestSignalRow, skips *map[btcPostCompressionFixedSkipKey]int, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, btCfg BacktestConfig, splits []Split) ([]BTC15MPostCompressionFixedBacktestSignalRow, []Trade) {
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
			btcPostCompressionFixedAddSkip(*skips, signal.Split, string(signal.Side), "open_position_already_active")
			continue
		}
		stopDist := math.Abs(signal.ExpectedEntryPrice - signal.Stop)
		size := positionSize(balance, signal.ExpectedEntryPrice, stopDist, btCfg)
		if size <= 0 {
			signals[signalIdx].SkippedReason = "non_positive_position_size"
			btcPostCompressionFixedAddSkip(*skips, signal.Split, string(signal.Side), "non_positive_position_size")
			continue
		}
		entryFee := signal.ExpectedEntryPrice * size * btCfg.FeePct
		entrySlip := math.Abs(signal.ExpectedEntryPrice-signal.EntryOpen) * size
		pending = &Position{
			Side:        Long,
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

func btcPostCompressionFixedBacktestTradeRows(trades []Trade, signals []BTC15MPostCompressionFixedBacktestSignalRow, splits []Split) []BTC15MPostCompressionFixedBacktestTradeRow {
	signalByID := map[string]BTC15MPostCompressionFixedBacktestSignalRow{}
	for _, signal := range signals {
		signalByID[signal.SignalID] = signal
	}
	rows := make([]BTC15MPostCompressionFixedBacktestTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := signalByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		stressNet := trade.NetPnL - trade.Slippage
		row := BTC15MPostCompressionFixedBacktestTradeRow{
			SignalID:                  signal.SignalID,
			CandidateID:               signal.CandidateID,
			Timeframe:                 signal.Timeframe,
			DecisionIndex:             signal.DecisionIndex,
			DecisionCloseTime:         signal.DecisionCloseTime,
			EntrySplit:                splitNameForCloseTime(entryTime, splits),
			CloseSplit:                splitNameForCloseTime(exitTime, splits),
			Side:                      trade.Side,
			EntryTime:                 trade.EntryTime,
			ExitTime:                  trade.ExitTime,
			OpenIndex:                 trade.OpenIndex,
			CloseIndex:                trade.CloseIndex,
			EntryPrice:                trade.EntryPrice,
			ExitPrice:                 trade.ExitPrice,
			Stop:                      trade.Stop,
			Target:                    trade.Target,
			Size:                      trade.Size,
			InitialRisk:               initialRisk,
			GrossPnL:                  trade.GrossPnL,
			EngineNetPnL:              trade.NetPnL,
			ExtraSlippageStressNetPnL: stressNet,
			Fees:                      trade.Fees,
			Slippage:                  trade.Slippage,
			ExitReason:                trade.Reason,
			HoldBars:                  trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.EngineNetR = trade.NetPnL / initialRisk
			row.ExtraSlippageStressNetR = stressNet / initialRisk
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

func btcPostCompressionFixedBacktestSummaryRows(signals []BTC15MPostCompressionFixedBacktestSignalRow, skips []BTC15MPostCompressionFixedBacktestSkipRow, trades []BTC15MPostCompressionFixedBacktestTradeRow, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, startBalance float64, splits []Split) []BTC15MPostCompressionFixedBacktestSummaryRow {
	rows := make([]BTC15MPostCompressionFixedBacktestSummaryRow, 0, len(splits)*3)
	for _, split := range splits {
		for _, side := range []string{"all", string(Long), string(Short)} {
			filteredTrades := btcPostCompressionFixedFilterTrades(trades, split, side)
			row := btcPostCompressionFixedSummarizeTrades(filteredTrades, startBalance)
			row.CandidateID = BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID
			row.Split = split.Name
			row.Side = side
			row.RawCandidateRows = btcPostCompressionFixedSignalCount(signals, split, side)
			row.SkippedRows = btcPostCompressionFixedSkipCount(skips, split, side)
			rows = append(rows, row)
		}
	}
	btcPostCompressionFixedMarkSummaryGates(rows, cfg, splits)
	return rows
}

func btcPostCompressionFixedBacktestCostStressRows(summaryRows []BTC15MPostCompressionFixedBacktestSummaryRow) []BTC15MPostCompressionFixedBacktestCostStressRow {
	rows := make([]BTC15MPostCompressionFixedBacktestCostStressRow, 0, len(summaryRows))
	for _, row := range summaryRows {
		rows = append(rows, BTC15MPostCompressionFixedBacktestCostStressRow{
			CandidateID:               row.CandidateID,
			Split:                     row.Split,
			Side:                      row.Side,
			ExecutedTrades:            row.ExecutedTrades,
			GrossPnL:                  row.GrossPnL,
			EngineNetPnL:              row.EngineNetPnL,
			Slippage:                  row.Slippage,
			ExtraSlippageStressNetPnL: row.ExtraSlippageStressNetPnL,
			Fees:                      row.Fees,
			StressProfitFactor:        row.StressProfitFactor,
			StressMaxDrawdown:         row.MaxDrawdown,
			PassFailUsesStressNet:     true,
		})
	}
	return rows
}

func btcPostCompressionFixedSummarizeTrades(trades []BTC15MPostCompressionFixedBacktestTradeRow, startBalance float64) BTC15MPostCompressionFixedBacktestSummaryRow {
	row := BTC15MPostCompressionFixedBacktestSummaryRow{ExecutedTrades: len(trades)}
	balance := startBalance
	equity := []float64{startBalance}
	stressProfit, stressLoss := 0.0, 0.0
	grossProfit, grossLoss := 0.0, 0.0
	holdBars := 0
	for _, trade := range trades {
		row.GrossPnL += trade.GrossPnL
		row.EngineNetPnL += trade.EngineNetPnL
		row.ExtraSlippageStressNetPnL += trade.ExtraSlippageStressNetPnL
		row.Fees += trade.Fees
		row.Slippage += trade.Slippage
		row.AvgGrossR += trade.GrossR
		row.AvgEngineNetR += trade.EngineNetR
		row.AvgExtraSlippageStressR += trade.ExtraSlippageStressNetR
		row.AvgInitialRisk += trade.InitialRisk
		holdBars += trade.HoldBars
		switch trade.ExitReason {
		case "stop_loss":
			row.StopLossExits++
		case "take_profit":
			row.TakeProfitExits++
		case "time_stop":
			row.TimeStopExits++
		case "force_close":
			row.ForceCloseExits++
		}
		if trade.ExtraSlippageStressNetPnL > 0 {
			row.Wins++
			stressProfit += trade.ExtraSlippageStressNetPnL
		} else if trade.ExtraSlippageStressNetPnL < 0 {
			row.Losses++
			stressLoss += -trade.ExtraSlippageStressNetPnL
		}
		if trade.GrossPnL > 0 {
			grossProfit += trade.GrossPnL
		} else if trade.GrossPnL < 0 {
			grossLoss += -trade.GrossPnL
		}
		balance += trade.ExtraSlippageStressNetPnL
		equity = append(equity, balance)
	}
	if row.ExecutedTrades > 0 {
		count := float64(row.ExecutedTrades)
		row.WinRate = float64(row.Wins) / count
		row.AvgGrossR /= count
		row.AvgEngineNetR /= count
		row.AvgExtraSlippageStressR /= count
		row.AvgInitialRisk /= count
		row.AvgHoldBars = float64(holdBars) / count
	}
	if stressLoss > 0 {
		row.StressProfitFactor = stressProfit / stressLoss
	} else if stressProfit > 0 {
		row.StressProfitFactor = 999.99
	}
	if grossLoss > 0 {
		row.GrossProfitFactor = grossProfit / grossLoss
	} else if grossProfit > 0 {
		row.GrossProfitFactor = 999.99
	}
	row.MaxDrawdown = MaxDrawdown(equity)
	return row
}

func btcPostCompressionFixedBacktestEvaluate(source BTC15MPostCompressionFixedBacktestSourceRow, coverage BTC15MPostCompressionFixedBacktestCoverageRow, signals []BTC15MPostCompressionFixedBacktestSignalRow, rows []BTC15MPostCompressionFixedBacktestSummaryRow, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig) BTC15MPostCompressionFixedBacktestFalsification {
	byKey := btcPostCompressionFixedSummaryByKey(rows)
	full := byKey[fullSplitName+"|"+string(Long)]
	oos := byKey["2023_2024_oos|"+string(Long)]
	recent := byKey["2025_2026_recent|"+string(Long)]
	fullRaw := btcPostCompressionFixedRawCandidateCount(signals)
	minPrimaryTrades := btcPostCompressionFixedMinPrimaryTrades(byKey, cfg)
	dominantExitShare := btcPostCompressionFixedDominantExitShare(full)
	dominantSplitShare := btcPostCompressionFixedDominantPrimarySplitTradeShare(byKey)
	sourcePass := source.SourceFactsPass && coverage.SourceResamplePass && source.ValidationStatus == "accepted" && coverage.ValidationStatus == "accepted"
	candidatePass := cfg.SkipCandidateIdentityGate || fullRaw == cfg.ExpectedRawCandidateRows
	tradeCountPass := full.ExecutedTrades >= cfg.MinFullTrades && minPrimaryTrades >= cfg.MinSplitTrades
	grossPass := full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0
	costedPass := full.ExtraSlippageStressNetPnL > 0 && oos.ExtraSlippageStressNetPnL >= 0 && recent.ExtraSlippageStressNetPnL >= 0
	pfPass := full.StressProfitFactor >= cfg.FullStressProfitFactorMin && oos.StressProfitFactor >= cfg.SplitStressProfitFactorMin && recent.StressProfitFactor >= cfg.SplitStressProfitFactorMin
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit
	for _, split := range btcPostCompressionPrimarySplitNames(DefaultSplits()) {
		row := byKey[split+"|"+string(Long)]
		if row.MaxDrawdown > cfg.SplitMaxDrawdownLimit {
			drawdownPass = false
		}
	}
	robustnessPass := tradeCountPass &&
		dominantExitShare <= BTC15MPostCompressionFixedBacktestExitReasonConcentrationLimit &&
		dominantSplitShare <= BTC15MPostCompressionFixedBacktestPrimarySplitTradeShareLimit
	return BTC15MPostCompressionFixedBacktestFalsification{
		SourceResamplePass:                 sourcePass,
		CandidateIdentityPass:              candidatePass,
		LeakagePass:                        !cfg.ForwardLabelLeakOverride,
		TradeCountPass:                     tradeCountPass,
		GrossEdgePass:                      grossPass,
		CostedEdgePass:                     costedPass,
		ProfitFactorPass:                   pfPass,
		DrawdownPass:                       drawdownPass,
		RobustnessPass:                     robustnessPass,
		OptimizerContaminationPass:         !cfg.OptimizerContamination,
		ClosedFamilyProtectionPass:         !cfg.ClosedFamilyReslice,
		DerivativesVetoContaminationPass:   !cfg.DerivativesVetoContam,
		CommonOutputsTradeCompatible:       true,
		ExpectedRawCandidateRows:           cfg.ExpectedRawCandidateRows,
		FullRawCandidateRows:               fullRaw,
		FullExecutedTrades:                 full.ExecutedTrades,
		RequiredFullExecutedTrades:         cfg.MinFullTrades,
		MinimumPrimarySplitExecutedTrades:  minPrimaryTrades,
		RequiredPrimarySplitExecutedTrades: cfg.MinSplitTrades,
		FullGrossPnL:                       full.GrossPnL,
		FullEngineNetPnL:                   full.EngineNetPnL,
		FullExtraSlippageStressNetPnL:      full.ExtraSlippageStressNetPnL,
		FullStressProfitFactor:             full.StressProfitFactor,
		FullMaxDrawdown:                    full.MaxDrawdown,
		DominantExitReasonShare:            dominantExitShare,
		DominantPrimarySplitTradeShare:     dominantSplitShare,
		ExitReasonConcentrationLimit:       BTC15MPostCompressionFixedBacktestExitReasonConcentrationLimit,
		PrimarySplitTradeShareLimit:        BTC15MPostCompressionFixedBacktestPrimarySplitTradeShareLimit,
	}
}

func btcPostCompressionFixedBacktestFalsification(report BTC15MPostCompressionFixedBacktestFalsification) BTC15MPostCompressionFixedBacktestFalsification {
	report.BacktestName = FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestName
	report.CandidateID = BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID
	failures := []string{}
	if !report.SourceResamplePass {
		failures = append(failures, "source_or_resample_failed")
	}
	if !report.CandidateIdentityPass {
		failures = append(failures, "representative_candidate_count_mismatch")
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
	if !report.CostedEdgePass {
		failures = append(failures, "extra_slippage_stress_edge_gate_failed")
	}
	if !report.ProfitFactorPass {
		failures = append(failures, "stress_profit_factor_gate_failed")
	}
	if !report.DrawdownPass {
		failures = append(failures, "drawdown_gate_failed")
	}
	if !report.RobustnessPass {
		failures = append(failures, "robustness_gate_failed")
	}
	if !report.OptimizerContaminationPass {
		failures = append(failures, "optimizer_contamination")
	}
	if !report.ClosedFamilyProtectionPass {
		failures = append(failures, "closed_family_reslice")
	}
	if !report.DerivativesVetoContaminationPass {
		failures = append(failures, "derivatives_veto_contamination")
	}
	report.FailureReasons = failures
	switch {
	case !report.SourceResamplePass || !report.CandidateIdentityPass || !report.LeakagePass:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStateFailedSourceOrResample
	case !report.OptimizerContaminationPass:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStateRejectedOptimizerContam
	case !report.ClosedFamilyProtectionPass:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStateRejectedClosedReslice
	case !report.DerivativesVetoContaminationPass:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStateRejectedVetoContam
	case report.TradeCountPass &&
		report.GrossEdgePass &&
		report.CostedEdgePass &&
		report.ProfitFactorPass &&
		report.DrawdownPass &&
		report.RobustnessPass:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStatePassedNeedsReview
	default:
		report.StopState = BTC15MPostCompressionFixedBacktestStopStateFailedNoUsableStrategy
	}
	return report
}

func btcPostCompressionFixedMarkSummaryGates(rows []BTC15MPostCompressionFixedBacktestSummaryRow, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig, splits []Split) {
	byKey := btcPostCompressionFixedSummaryByKey(rows)
	full := byKey[fullSplitName+"|"+string(Long)]
	oos := byKey["2023_2024_oos|"+string(Long)]
	recent := byKey["2025_2026_recent|"+string(Long)]
	minPrimaryTrades := btcPostCompressionFixedMinPrimaryTrades(byKey, cfg)
	for i := range rows {
		if rows[i].Side != string(Long) && rows[i].Side != "all" {
			continue
		}
		rows[i].TradeCountGatePass = full.ExecutedTrades >= cfg.MinFullTrades && minPrimaryTrades >= cfg.MinSplitTrades
		rows[i].GrossEdgeGatePass = full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0
		rows[i].CostedEdgeGatePass = full.ExtraSlippageStressNetPnL > 0 && oos.ExtraSlippageStressNetPnL >= 0 && recent.ExtraSlippageStressNetPnL >= 0
		rows[i].ProfitFactorGatePass = full.StressProfitFactor >= cfg.FullStressProfitFactorMin && oos.StressProfitFactor >= cfg.SplitStressProfitFactorMin && recent.StressProfitFactor >= cfg.SplitStressProfitFactorMin
		rows[i].DrawdownGatePass = full.MaxDrawdown <= cfg.FullMaxDrawdownLimit
		for _, split := range btcPostCompressionPrimarySplitNames(splits) {
			if byKey[split+"|"+string(Long)].MaxDrawdown > cfg.SplitMaxDrawdownLimit {
				rows[i].DrawdownGatePass = false
			}
		}
		rows[i].RobustnessGatePass = rows[i].TradeCountGatePass &&
			btcPostCompressionFixedDominantExitShare(full) <= BTC15MPostCompressionFixedBacktestExitReasonConcentrationLimit &&
			btcPostCompressionFixedDominantPrimarySplitTradeShare(byKey) <= BTC15MPostCompressionFixedBacktestPrimarySplitTradeShareLimit
	}
}

func btcPostCompressionFixedFilterTrades(trades []BTC15MPostCompressionFixedBacktestTradeRow, split Split, side string) []BTC15MPostCompressionFixedBacktestTradeRow {
	out := make([]BTC15MPostCompressionFixedBacktestTradeRow, 0, len(trades))
	for _, trade := range trades {
		exitTime, err := parseTime(trade.ExitTime)
		if err != nil || !split.Contains(exitTime) {
			continue
		}
		if side != "all" && string(trade.Side) != side {
			continue
		}
		out = append(out, trade)
	}
	return out
}

func btcPostCompressionFixedSignalCount(signals []BTC15MPostCompressionFixedBacktestSignalRow, split Split, side string) int {
	count := 0
	for _, signal := range signals {
		if !signal.PrePositionCandidate {
			continue
		}
		decisionTime, err := parseTime(signal.DecisionCloseTime)
		if err != nil || !split.Contains(decisionTime) {
			continue
		}
		if side != "all" && string(signal.Side) != side {
			continue
		}
		count++
	}
	return count
}

func btcPostCompressionFixedSkipCount(skips []BTC15MPostCompressionFixedBacktestSkipRow, split Split, side string) int {
	count := 0
	for _, skip := range skips {
		if skip.Split != split.Name && split.Name != fullSplitName {
			continue
		}
		if side != "all" && skip.Side != side {
			continue
		}
		count += skip.Count
	}
	return count
}

func btcPostCompressionFixedRawCandidateCount(signals []BTC15MPostCompressionFixedBacktestSignalRow) int {
	count := 0
	for _, signal := range signals {
		if signal.PrePositionCandidate {
			count++
		}
	}
	return count
}

func btcPostCompressionFixedMinPrimaryTrades(byKey map[string]BTC15MPostCompressionFixedBacktestSummaryRow, cfg FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig) int {
	minTrades := int(^uint(0) >> 1)
	for _, split := range btcPostCompressionPrimarySplitNames(DefaultSplits()) {
		row := byKey[split+"|"+string(Long)]
		if row.ExecutedTrades < minTrades {
			minTrades = row.ExecutedTrades
		}
	}
	if minTrades == int(^uint(0)>>1) {
		return 0
	}
	return minTrades
}

func btcPostCompressionFixedDominantExitShare(row BTC15MPostCompressionFixedBacktestSummaryRow) float64 {
	if row.ExecutedTrades == 0 {
		return 1
	}
	maxCount := row.StopLossExits
	for _, count := range []int{row.TakeProfitExits, row.TimeStopExits, row.ForceCloseExits} {
		if count > maxCount {
			maxCount = count
		}
	}
	return float64(maxCount) / float64(row.ExecutedTrades)
}

func btcPostCompressionFixedDominantPrimarySplitTradeShare(byKey map[string]BTC15MPostCompressionFixedBacktestSummaryRow) float64 {
	full := byKey[fullSplitName+"|"+string(Long)]
	if full.ExecutedTrades == 0 {
		return 1
	}
	maxTrades := 0
	for _, split := range btcPostCompressionPrimarySplitNames(DefaultSplits()) {
		trades := byKey[split+"|"+string(Long)].ExecutedTrades
		if trades > maxTrades {
			maxTrades = trades
		}
	}
	return float64(maxTrades) / float64(full.ExecutedTrades)
}

func btcPostCompressionFixedSummaryByKey(rows []BTC15MPostCompressionFixedBacktestSummaryRow) map[string]BTC15MPostCompressionFixedBacktestSummaryRow {
	out := map[string]BTC15MPostCompressionFixedBacktestSummaryRow{}
	for _, row := range rows {
		out[row.Split+"|"+row.Side] = row
	}
	return out
}

func btcPostCompressionFixedAddSkip(skips map[btcPostCompressionFixedSkipKey]int, split, side, reason string) {
	if split == "" {
		split = fullSplitName
	}
	if side == "" {
		side = string(Long)
	}
	skips[btcPostCompressionFixedSkipKey{split: split, side: side, reason: reason}]++
	if split != fullSplitName {
		skips[btcPostCompressionFixedSkipKey{split: fullSplitName, side: side, reason: reason}]++
	}
}

func btcPostCompressionFixedBacktestSkipRows(skips map[btcPostCompressionFixedSkipKey]int) []BTC15MPostCompressionFixedBacktestSkipRow {
	keys := make([]btcPostCompressionFixedSkipKey, 0, len(skips))
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
	rows := make([]BTC15MPostCompressionFixedBacktestSkipRow, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, BTC15MPostCompressionFixedBacktestSkipRow{
			CandidateID:       BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
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
