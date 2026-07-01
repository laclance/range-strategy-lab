package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	BacktestFirstBTC15MSessionOpeningRangeExpansionName = "backtest_first_btc_15m_session_opening_range_expansion"
	BTC15MSessionOpeningRangeExpansionCandidateID       = "btc_15m_session_opening_range_expansion_v1"

	BTC15MSessionOpeningRangeExpansionStopStatePassedNeedsReview      = "btc_15m_session_opening_range_expansion_backtest_passed_needs_review"
	BTC15MSessionOpeningRangeExpansionStopStateFailedNoUsableStrategy = "btc_15m_session_opening_range_expansion_backtest_failed_no_usable_strategy"
	BTC15MSessionOpeningRangeExpansionStopStateFailedSourceOrResample = "btc_15m_session_opening_range_expansion_backtest_failed_source_or_resample"
)

type BacktestFirstBTC15MSessionOpeningRangeExpansionConfig struct {
	ApprovedSourcePath      string
	ExpectedSourceRows      int
	ExpectedFirstOpenTime   string
	ExpectedLastOpenTime    string
	ExpectedGapCount        int
	ExpectedDuplicateCount  int
	ExpectedZeroVolumeCount int
	Expected15MRows         int
	Expected15MLastOpenTime string
	ATRPeriod               int
	AcceptanceATRMultiple   float64
	SessionAnchorHour       int
	SessionAnchorMinute     int
	OpeningRangeBars        int
	ExpansionStartHour      int
	ExpansionStartMinute    int
	ExpansionEndHour        int
	ExpansionEndMinute      int
	TargetR                 float64
	MaxHoldBars             int
	MinFullTrades           int
	MinSplitTrades          int
	FullProfitFactorMin     float64
	SplitProfitFactorMin    float64
	FullMaxDrawdownLimit    float64
	SplitMaxDrawdownLimit   float64
	SkipSourceFactCheck     bool
	SkipCoverageCountCheck  bool
}

type BacktestFirstBTC15MSessionOpeningRangeExpansionResult struct {
	SourceRows       []BTC15MSessionOpeningRangeExpansionSourceRow       `json:"source_rows"`
	CoverageRows     []BTC15MSessionOpeningRangeExpansionCoverageRow     `json:"coverage_rows"`
	SessionRangeRows []BTC15MSessionOpeningRangeExpansionSessionRangeRow `json:"session_range_rows"`
	SignalRows       []BTC15MSessionOpeningRangeExpansionSignalRow       `json:"signal_rows"`
	SkipRows         []BTC15MSessionOpeningRangeExpansionSkipRow         `json:"skip_rows"`
	TradeRows        []BTC15MSessionOpeningRangeExpansionTradeRow        `json:"trade_rows"`
	SummaryRows      []SummaryRow                                        `json:"summary_rows"`
	Falsification    BTC15MSessionOpeningRangeExpansionFalsification     `json:"falsification"`
	Trades           []Trade                                             `json:"trades"`
	StopState        string                                              `json:"stop_state"`
}

type BTC15MSessionOpeningRangeExpansionSourceRow struct {
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

type BTC15MSessionOpeningRangeExpansionCoverageRow struct {
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

type BTC15MSessionOpeningRangeExpansionSessionRangeRow struct {
	CandidateID             string  `json:"candidate_id"`
	Timeframe               string  `json:"timeframe"`
	SessionDateUTC          string  `json:"session_date_utc"`
	SessionAnchorUTC        string  `json:"session_anchor_utc"`
	SessionAnchorOpenTime   string  `json:"session_anchor_open_time"`
	OpeningRangeStartTime   string  `json:"opening_range_start_time"`
	OpeningRangeEndTime     string  `json:"opening_range_end_time"`
	OpeningRangeBars        int     `json:"opening_range_bars"`
	OpeningRangeHigh        float64 `json:"opening_range_high"`
	OpeningRangeLow         float64 `json:"opening_range_low"`
	OpeningRangeWidth       float64 `json:"opening_range_width"`
	ExpansionWindowStartUTC string  `json:"expansion_window_start_utc"`
	ExpansionWindowEndUTC   string  `json:"expansion_window_end_utc"`
	ATRPeriod               int     `json:"atr_period"`
	AcceptanceATRMultiple   float64 `json:"acceptance_atr_multiple"`
	RangeReady              bool    `json:"range_ready"`
	RangeWidthPositive      bool    `json:"range_width_positive"`
	ClosedCandleOnly        bool    `json:"closed_candle_only"`
	DSTShifted              bool    `json:"dst_shifted"`
	AlternateAnchorCompared bool    `json:"alternate_anchor_compared"`
	SkippedReason           string  `json:"skipped_reason,omitempty"`
}

type BTC15MSessionOpeningRangeExpansionSignalRow struct {
	SignalID                  string    `json:"signal_id"`
	CandidateID               string    `json:"candidate_id"`
	Timeframe                 string    `json:"timeframe"`
	Split                     string    `json:"split"`
	SessionDateUTC            string    `json:"session_date_utc"`
	DecisionIndex             int       `json:"decision_index"`
	DecisionOpenTime          string    `json:"decision_open_time"`
	DecisionCloseTime         string    `json:"decision_close_time"`
	DecisionOpen              float64   `json:"decision_open"`
	DecisionHigh              float64   `json:"decision_high"`
	DecisionLow               float64   `json:"decision_low"`
	DecisionClose             float64   `json:"decision_close"`
	Side                      Direction `json:"side"`
	TimingLabel               string    `json:"timing_label"`
	OpeningRangeHigh          float64   `json:"opening_range_high"`
	OpeningRangeLow           float64   `json:"opening_range_low"`
	OpeningRangeWidth         float64   `json:"opening_range_width"`
	ATR14                     float64   `json:"atr14"`
	AcceptanceATRMultiple     float64   `json:"acceptance_atr_multiple"`
	AcceptanceBuffer          float64   `json:"acceptance_buffer"`
	LongTriggerPrice          float64   `json:"long_trigger_price"`
	ShortTriggerPrice         float64   `json:"short_trigger_price"`
	EntryIndex                int       `json:"entry_index"`
	EntryOpenTime             string    `json:"entry_open_time"`
	EntryOpen                 float64   `json:"entry_open"`
	ExpectedEntryPrice        float64   `json:"expected_entry_price"`
	Stop                      float64   `json:"stop"`
	Target                    float64   `json:"target"`
	TargetR                   float64   `json:"target_r"`
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

type BTC15MSessionOpeningRangeExpansionSkipRow struct {
	CandidateID       string `json:"candidate_id"`
	Split             string `json:"split"`
	Side              string `json:"side"`
	Reason            string `json:"reason"`
	Count             int    `json:"count"`
	MissingDataPolicy string `json:"missing_data_policy"`
	ForwardFilledRows int    `json:"forward_filled_rows"`
}

type BTC15MSessionOpeningRangeExpansionTradeRow struct {
	SignalID      string    `json:"signal_id"`
	CandidateID   string    `json:"candidate_id"`
	Timeframe     string    `json:"timeframe"`
	SessionDate   string    `json:"session_date_utc"`
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

type BTC15MSessionOpeningRangeExpansionFalsification struct {
	BacktestName                       string   `json:"backtest_name"`
	CandidateID                        string   `json:"candidate_id"`
	StopState                          string   `json:"stop_state"`
	SourceResamplePass                 bool     `json:"source_resample_pass"`
	FixedSessionSpecPass               bool     `json:"fixed_session_spec_pass"`
	LeakagePass                        bool     `json:"leakage_pass"`
	TradeCountPass                     bool     `json:"trade_count_pass"`
	GrossEdgePass                      bool     `json:"gross_edge_pass"`
	NetEdgePass                        bool     `json:"net_edge_pass"`
	ProfitFactorPass                   bool     `json:"profit_factor_pass"`
	DrawdownPass                       bool     `json:"drawdown_pass"`
	SideReportingPass                  bool     `json:"side_reporting_pass"`
	CombinedBaselineSelectionPass      bool     `json:"combined_baseline_selection_pass"`
	OptimizerContaminationPass         bool     `json:"optimizer_contamination_pass"`
	DerivativesVetoContaminationPass   bool     `json:"derivatives_veto_contamination_pass"`
	FullExecutedTrades                 int      `json:"full_executed_trades"`
	RequiredFullExecutedTrades         int      `json:"required_full_executed_trades"`
	MinimumPrimarySplitExecutedTrades  int      `json:"minimum_primary_split_executed_trades"`
	RequiredPrimarySplitExecutedTrades int      `json:"required_primary_split_executed_trades"`
	FullGrossPnL                       float64  `json:"full_gross_pnl"`
	FullNetPnL                         float64  `json:"full_net_pnl"`
	FullProfitFactor                   float64  `json:"full_profit_factor"`
	FullMaxDrawdown                    float64  `json:"full_max_drawdown"`
	StressGrossPnL                     float64  `json:"stress_gross_pnl"`
	StressNetPnL                       float64  `json:"stress_net_pnl"`
	StressProfitFactor                 float64  `json:"stress_profit_factor"`
	StressMaxDrawdown                  float64  `json:"stress_max_drawdown"`
	OOSGrossPnL                        float64  `json:"oos_gross_pnl"`
	OOSNetPnL                          float64  `json:"oos_net_pnl"`
	OOSProfitFactor                    float64  `json:"oos_profit_factor"`
	OOSMaxDrawdown                     float64  `json:"oos_max_drawdown"`
	RecentGrossPnL                     float64  `json:"recent_gross_pnl"`
	RecentNetPnL                       float64  `json:"recent_net_pnl"`
	RecentProfitFactor                 float64  `json:"recent_profit_factor"`
	RecentMaxDrawdown                  float64  `json:"recent_max_drawdown"`
	LongFullExecutedTrades             int      `json:"long_full_executed_trades"`
	LongFullGrossPnL                   float64  `json:"long_full_gross_pnl"`
	LongFullNetPnL                     float64  `json:"long_full_net_pnl"`
	LongFullProfitFactor               float64  `json:"long_full_profit_factor"`
	ShortFullExecutedTrades            int      `json:"short_full_executed_trades"`
	ShortFullGrossPnL                  float64  `json:"short_full_gross_pnl"`
	ShortFullNetPnL                    float64  `json:"short_full_net_pnl"`
	ShortFullProfitFactor              float64  `json:"short_full_profit_factor"`
	FailureReasons                     []string `json:"failure_reasons,omitempty"`
}

type btc15MSessionOpeningRangeSkipKey struct {
	split  string
	side   string
	reason string
}

func DefaultBacktestFirstBTC15MSessionOpeningRangeExpansionConfig() BacktestFirstBTC15MSessionOpeningRangeExpansionConfig {
	return BacktestFirstBTC15MSessionOpeningRangeExpansionConfig{
		ApprovedSourcePath:      "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedSourceRows:      573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedGapCount:        0,
		ExpectedDuplicateCount:  0,
		ExpectedZeroVolumeCount: 66,
		Expected15MRows:         191328,
		Expected15MLastOpenTime: "2026-06-16T23:45:00Z",
		ATRPeriod:               14,
		AcceptanceATRMultiple:   0.10,
		SessionAnchorHour:       13,
		SessionAnchorMinute:     30,
		OpeningRangeBars:        4,
		ExpansionStartHour:      14,
		ExpansionStartMinute:    30,
		ExpansionEndHour:        17,
		ExpansionEndMinute:      30,
		TargetR:                 1.5,
		MaxHoldBars:             24,
		MinFullTrades:           120,
		MinSplitTrades:          25,
		FullProfitFactorMin:     1.10,
		SplitProfitFactorMin:    1.00,
		FullMaxDrawdownLimit:    0.25,
		SplitMaxDrawdownLimit:   0.30,
	}
}

func RunBacktestFirstBTC15MSessionOpeningRangeExpansion(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC15MSessionOpeningRangeExpansionResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return BacktestFirstBTC15MSessionOpeningRangeExpansionResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc15MSessionOpeningRangeSourceRow(manifest, cfg)
	result := BacktestFirstBTC15MSessionOpeningRangeExpansionResult{SourceRows: []BTC15MSessionOpeningRangeExpansionSourceRow{source}}
	resampled, coverage := btc15MSessionOpeningRangeResample(candles, cfg, source.SourceFactsPass)
	result.CoverageRows = []BTC15MSessionOpeningRangeExpansionCoverageRow{coverage}
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = btc15MSessionOpeningRangeFalsification(BTC15MSessionOpeningRangeExpansionFalsification{SourceResamplePass: false, FixedSessionSpecPass: true, LeakagePass: true, SideReportingPass: true, CombinedBaselineSelectionPass: true, OptimizerContaminationPass: true, DerivativesVetoContaminationPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}

	sessionRows, signals, skipCounts := btc15MSessionOpeningRangeBuildSignals(resampled, cfg, btCfg, splits)
	signals, trades := btc15MSessionOpeningRangeRun(resampled, signals, &skipCounts, cfg, btCfg)
	result.SessionRangeRows = sessionRows
	result.Trades = trades
	result.SignalRows = signals
	result.SkipRows = btc15MSessionOpeningRangeSkipRows(skipCounts)
	result.TradeRows = btc15MSessionOpeningRangeTradeRows(trades, signals, splits)
	result.SummaryRows = SummarizeSplits(trades, btCfg.StartBalance, splits)
	result.Falsification = btc15MSessionOpeningRangeFalsification(btc15MSessionOpeningRangeEvaluate(source, coverage, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}

func (cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig) withDefaults() BacktestFirstBTC15MSessionOpeningRangeExpansionConfig {
	defaults := DefaultBacktestFirstBTC15MSessionOpeningRangeExpansionConfig()
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
	if cfg.ExpectedZeroVolumeCount == 0 && !cfg.SkipSourceFactCheck {
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
	if cfg.AcceptanceATRMultiple == 0 {
		cfg.AcceptanceATRMultiple = defaults.AcceptanceATRMultiple
	}
	if cfg.SessionAnchorHour == 0 && cfg.SessionAnchorMinute == 0 {
		cfg.SessionAnchorHour = defaults.SessionAnchorHour
		cfg.SessionAnchorMinute = defaults.SessionAnchorMinute
	}
	if cfg.OpeningRangeBars == 0 {
		cfg.OpeningRangeBars = defaults.OpeningRangeBars
	}
	if cfg.ExpansionStartHour == 0 && cfg.ExpansionStartMinute == 0 {
		cfg.ExpansionStartHour = defaults.ExpansionStartHour
		cfg.ExpansionStartMinute = defaults.ExpansionStartMinute
	}
	if cfg.ExpansionEndHour == 0 && cfg.ExpansionEndMinute == 0 {
		cfg.ExpansionEndHour = defaults.ExpansionEndHour
		cfg.ExpansionEndMinute = defaults.ExpansionEndMinute
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

func (cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig) validate() error {
	if cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("session opening-range expansion approved source path is required")
	}
	if cfg.ATRPeriod <= 0 {
		return fmt.Errorf("session opening-range expansion ATR period must be positive")
	}
	if cfg.AcceptanceATRMultiple <= 0 || cfg.TargetR <= 0 {
		return fmt.Errorf("session opening-range expansion acceptance and target multiples must be positive")
	}
	if cfg.OpeningRangeBars != 4 {
		return fmt.Errorf("session opening-range expansion opening-range bars must stay fixed at 4")
	}
	if cfg.SessionAnchorHour != 13 || cfg.SessionAnchorMinute != 30 {
		return fmt.Errorf("session opening-range expansion session anchor must stay fixed at 13:30 UTC")
	}
	if cfg.ExpansionStartHour != 14 || cfg.ExpansionStartMinute != 30 || cfg.ExpansionEndHour != 17 || cfg.ExpansionEndMinute != 30 {
		return fmt.Errorf("session opening-range expansion window must stay fixed at [14:30, 17:30) UTC")
	}
	if cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("session opening-range expansion max hold bars must be positive")
	}
	if cfg.MinFullTrades <= 0 || cfg.MinSplitTrades <= 0 {
		return fmt.Errorf("session opening-range expansion trade-count gates must be positive")
	}
	if cfg.FullProfitFactorMin <= 0 || cfg.SplitProfitFactorMin <= 0 {
		return fmt.Errorf("session opening-range expansion profit-factor gates must be positive")
	}
	if cfg.FullMaxDrawdownLimit <= 0 || cfg.SplitMaxDrawdownLimit <= 0 {
		return fmt.Errorf("session opening-range expansion drawdown limits must be positive")
	}
	return nil
}

func btc15MSessionOpeningRangeSourceRow(manifest SourceManifest, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig) BTC15MSessionOpeningRangeExpansionSourceRow {
	row := BTC15MSessionOpeningRangeExpansionSourceRow{
		BacktestName:               BacktestFirstBTC15MSessionOpeningRangeExpansionName,
		CandidateID:                BTC15MSessionOpeningRangeExpansionCandidateID,
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
	if manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" {
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

func btc15MSessionOpeningRangeResample(candles []Candle, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, sourcePass bool) ([]Candle, BTC15MSessionOpeningRangeExpansionCoverageRow) {
	resampled, base := btc15MPrevDayResample(candles, BacktestFirstBTC15MPreviousDayRangeReversionConfig{
		Expected15MRows:         cfg.Expected15MRows,
		Expected15MLastOpenTime: cfg.Expected15MLastOpenTime,
	}, sourcePass)
	pass := base.SourceResamplePass
	if cfg.SkipCoverageCountCheck {
		pass = sourcePass && base.MissingChildBuckets == 0
	}
	row := BTC15MSessionOpeningRangeExpansionCoverageRow{
		BacktestName:         BacktestFirstBTC15MSessionOpeningRangeExpansionName,
		CandidateID:          BTC15MSessionOpeningRangeExpansionCandidateID,
		Timeframe:            "15m",
		RowCount:             base.RowCount,
		ExpectedRowCount:     cfg.Expected15MRows,
		FirstOpenTime:        base.FirstOpenTime,
		LastOpenTime:         base.LastOpenTime,
		ExpectedLastOpenTime: cfg.Expected15MLastOpenTime,
		ExpectedChildBars:    3,
		MissingChildBuckets:  base.MissingChildBuckets,
		ClosedCandleOnly:     true,
		SourceResamplePass:   pass,
		ValidationStatus:     "accepted",
	}
	if !row.SourceResamplePass {
		row.ValidationStatus = "rejected"
		row.ValidationError = base.ValidationError
	}
	return resampled, row
}

func btc15MSessionOpeningRangeBuildSignals(candles []Candle, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, btCfg BacktestConfig, splits []Split) ([]BTC15MSessionOpeningRangeExpansionSessionRangeRow, []BTC15MSessionOpeningRangeExpansionSignalRow, map[btc15MSessionOpeningRangeSkipKey]int) {
	atr := ATR(candles, cfg.ATRPeriod)
	indexByOpen := make(map[time.Time]int, len(candles))
	dateSeen := map[string]bool{}
	dates := []time.Time{}
	for i, candle := range candles {
		open := candle.OpenTime.UTC()
		indexByOpen[open] = i
		date := time.Date(open.Year(), open.Month(), open.Day(), 0, 0, 0, 0, time.UTC)
		key := date.Format("2006-01-02")
		if !dateSeen[key] {
			dateSeen[key] = true
			dates = append(dates, date)
		}
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

	sessionRows := []BTC15MSessionOpeningRangeExpansionSessionRangeRow{}
	signals := []BTC15MSessionOpeningRangeExpansionSignalRow{}
	skips := map[btc15MSessionOpeningRangeSkipKey]int{}
	for _, date := range dates {
		row, signal, ok := btc15MSessionOpeningRangeEvaluateDate(candles, atr, indexByOpen, date, len(signals)+1, cfg, btCfg, splits, skips)
		sessionRows = append(sessionRows, row)
		if ok {
			signals = append(signals, signal)
		}
	}
	return sessionRows, signals, skips
}

func btc15MSessionOpeningRangeEvaluateDate(candles []Candle, atr []float64, indexByOpen map[time.Time]int, date time.Time, sequence int, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, btCfg BacktestConfig, splits []Split, skips map[btc15MSessionOpeningRangeSkipKey]int) (BTC15MSessionOpeningRangeExpansionSessionRangeRow, BTC15MSessionOpeningRangeExpansionSignalRow, bool) {
	anchor := time.Date(date.Year(), date.Month(), date.Day(), cfg.SessionAnchorHour, cfg.SessionAnchorMinute, 0, 0, time.UTC)
	expansionStart := time.Date(date.Year(), date.Month(), date.Day(), cfg.ExpansionStartHour, cfg.ExpansionStartMinute, 0, 0, time.UTC)
	expansionEnd := time.Date(date.Year(), date.Month(), date.Day(), cfg.ExpansionEndHour, cfg.ExpansionEndMinute, 0, 0, time.UTC)
	split := splitNameForCloseTime(anchor.Add(time.Hour), splits)
	row := BTC15MSessionOpeningRangeExpansionSessionRangeRow{
		CandidateID:             BTC15MSessionOpeningRangeExpansionCandidateID,
		Timeframe:               "15m",
		SessionDateUTC:          date.Format("2006-01-02"),
		SessionAnchorUTC:        "13:30:00Z",
		SessionAnchorOpenTime:   anchor.Format(timeLayout),
		OpeningRangeStartTime:   anchor.Format(timeLayout),
		ExpansionWindowStartUTC: expansionStart.Format(timeLayout),
		ExpansionWindowEndUTC:   expansionEnd.Format(timeLayout),
		OpeningRangeBars:        cfg.OpeningRangeBars,
		ATRPeriod:               cfg.ATRPeriod,
		AcceptanceATRMultiple:   cfg.AcceptanceATRMultiple,
		ClosedCandleOnly:        true,
		DSTShifted:              false,
		AlternateAnchorCompared: false,
	}
	anchorIdx, ok := indexByOpen[anchor]
	if !ok {
		row.SkippedReason = "missing_session_anchor"
		btc15MSessionOpeningRangeAddSkip(skips, split, "all", row.SkippedReason)
		return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
	}
	if anchorIdx+cfg.OpeningRangeBars > len(candles) {
		row.SkippedReason = "missing_opening_range_bars"
		btc15MSessionOpeningRangeAddSkip(skips, split, "all", row.SkippedReason)
		return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
	}
	rangeHigh := math.Inf(-1)
	rangeLow := math.Inf(1)
	for offset := 0; offset < cfg.OpeningRangeBars; offset++ {
		idx := anchorIdx + offset
		expectedOpen := anchor.Add(time.Duration(offset*15) * time.Minute)
		if idx >= len(candles) || !candles[idx].OpenTime.UTC().Equal(expectedOpen) {
			row.SkippedReason = "missing_opening_range_bars"
			btc15MSessionOpeningRangeAddSkip(skips, split, "all", row.SkippedReason)
			return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
		}
		rangeHigh = math.Max(rangeHigh, candles[idx].High)
		rangeLow = math.Min(rangeLow, candles[idx].Low)
		row.OpeningRangeEndTime = candles[idx].CloseTime.UTC().Format(timeLayout)
	}
	row.OpeningRangeHigh = rangeHigh
	row.OpeningRangeLow = rangeLow
	row.OpeningRangeWidth = rangeHigh - rangeLow
	row.RangeReady = true
	row.RangeWidthPositive = row.OpeningRangeWidth > 0
	if !row.RangeWidthPositive {
		row.SkippedReason = "non_positive_opening_range_width"
		btc15MSessionOpeningRangeAddSkip(skips, split, "all", row.SkippedReason)
		return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
	}

	foundExpansionWindowBar := false
	for d := anchorIdx + cfg.OpeningRangeBars; d < len(candles); d++ {
		open := candles[d].OpenTime.UTC()
		if !sameUTCDate(open, date) {
			break
		}
		if open.Before(expansionStart) {
			continue
		}
		if !open.Before(expansionEnd) {
			break
		}
		foundExpansionWindowBar = true
		decisionSplit := splitNameForCloseTime(candles[d].CloseTime, splits)
		if d+1 >= len(candles) {
			btc15MSessionOpeningRangeAddSkip(skips, decisionSplit, "all", "missing_entry_candle")
			return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
		}
		if !validNumber(atr[d]) || atr[d] <= 0 {
			btc15MSessionOpeningRangeAddSkip(skips, decisionSplit, "all", "missing_atr")
			continue
		}
		buffer := cfg.AcceptanceATRMultiple * atr[d]
		longOK := candles[d].Close >= rangeHigh+buffer
		shortOK := candles[d].Close <= rangeLow-buffer
		if longOK && shortOK {
			btc15MSessionOpeningRangeAddSkip(skips, decisionSplit, "all", "ambiguous_both_sides")
			return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
		}
		if !longOK && !shortOK {
			continue
		}
		side := Long
		if shortOK {
			side = Short
		}
		signal := btc15MSessionOpeningRangeSignalRow(candles, d, sequence, side, date, rangeHigh, rangeLow, atr[d], cfg, btCfg, splits)
		if signal.SkippedReason != "" {
			btc15MSessionOpeningRangeAddSkip(skips, signal.Split, string(signal.Side), signal.SkippedReason)
		}
		return row, signal, true
	}
	if !foundExpansionWindowBar {
		btc15MSessionOpeningRangeAddSkip(skips, split, "all", "missing_expansion_window")
		return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
	}
	btc15MSessionOpeningRangeAddSkip(skips, split, "all", "no_expansion_acceptance")
	return row, BTC15MSessionOpeningRangeExpansionSignalRow{}, false
}

func btc15MSessionOpeningRangeSignalRow(candles []Candle, d int, sequence int, side Direction, date time.Time, rangeHigh, rangeLow, atrValue float64, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, btCfg BacktestConfig, splits []Split) BTC15MSessionOpeningRangeExpansionSignalRow {
	decision := candles[d]
	entryIndex := d + 1
	entry := candles[entryIndex]
	expectedEntry := applySlippage(entry.Open, btCfg.SlippagePct, side, true)
	buffer := cfg.AcceptanceATRMultiple * atrValue
	stop := rangeLow - buffer
	if side == Short {
		stop = rangeHigh + buffer
	}
	stopDist := math.Abs(expectedEntry - stop)
	target := expectedEntry + cfg.TargetR*stopDist
	if side == Short {
		target = expectedEntry - cfg.TargetR*stopDist
	}
	row := BTC15MSessionOpeningRangeExpansionSignalRow{
		SignalID:                  fmt.Sprintf("%s_%06d", BTC15MSessionOpeningRangeExpansionCandidateID, sequence),
		CandidateID:               BTC15MSessionOpeningRangeExpansionCandidateID,
		Timeframe:                 "15m",
		Split:                     splitNameForCloseTime(decision.CloseTime, splits),
		SessionDateUTC:            date.Format("2006-01-02"),
		DecisionIndex:             d,
		DecisionOpenTime:          decision.OpenTime.UTC().Format(timeLayout),
		DecisionCloseTime:         decision.CloseTime.UTC().Format(timeLayout),
		DecisionOpen:              decision.Open,
		DecisionHigh:              decision.High,
		DecisionLow:               decision.Low,
		DecisionClose:             decision.Close,
		Side:                      side,
		TimingLabel:               "next_15m_open",
		OpeningRangeHigh:          rangeHigh,
		OpeningRangeLow:           rangeLow,
		OpeningRangeWidth:         rangeHigh - rangeLow,
		ATR14:                     atrValue,
		AcceptanceATRMultiple:     cfg.AcceptanceATRMultiple,
		AcceptanceBuffer:          buffer,
		LongTriggerPrice:          rangeHigh + buffer,
		ShortTriggerPrice:         rangeLow - buffer,
		EntryIndex:                entryIndex,
		EntryOpenTime:             entry.OpenTime.UTC().Format(timeLayout),
		EntryOpen:                 entry.Open,
		ExpectedEntryPrice:        expectedEntry,
		Stop:                      stop,
		Target:                    target,
		TargetR:                   cfg.TargetR,
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

func btc15MSessionOpeningRangeRun(candles []Candle, signals []BTC15MSessionOpeningRangeExpansionSignalRow, skips *map[btc15MSessionOpeningRangeSkipKey]int, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig, btCfg BacktestConfig) ([]BTC15MSessionOpeningRangeExpansionSignalRow, []Trade) {
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
			btc15MSessionOpeningRangeAddSkip(*skips, signal.Split, string(signal.Side), "open_position_already_active")
			continue
		}
		stopDist := math.Abs(signal.ExpectedEntryPrice - signal.Stop)
		size := positionSize(balance, signal.ExpectedEntryPrice, stopDist, btCfg)
		if size <= 0 {
			signals[signalIdx].SkippedReason = "non_positive_position_size"
			btc15MSessionOpeningRangeAddSkip(*skips, signal.Split, string(signal.Side), "non_positive_position_size")
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

func btc15MSessionOpeningRangeTradeRows(trades []Trade, signals []BTC15MSessionOpeningRangeExpansionSignalRow, splits []Split) []BTC15MSessionOpeningRangeExpansionTradeRow {
	signalByID := map[string]BTC15MSessionOpeningRangeExpansionSignalRow{}
	for _, signal := range signals {
		signalByID[signal.SignalID] = signal
	}
	rows := make([]BTC15MSessionOpeningRangeExpansionTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := signalByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := BTC15MSessionOpeningRangeExpansionTradeRow{
			SignalID:      signal.SignalID,
			CandidateID:   signal.CandidateID,
			Timeframe:     signal.Timeframe,
			SessionDate:   signal.SessionDateUTC,
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

func btc15MSessionOpeningRangeEvaluate(source BTC15MSessionOpeningRangeExpansionSourceRow, coverage BTC15MSessionOpeningRangeExpansionCoverageRow, rows []SummaryRow, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig) BTC15MSessionOpeningRangeExpansionFalsification {
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
	tradeCountPass := full.TotalTrades >= cfg.MinFullTrades && minPrimary >= cfg.MinSplitTrades
	grossEdgePass := full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0
	netEdgePass := full.NetPnL > 0 && oos.NetPnL >= 0 && recent.NetPnL >= 0
	profitFactorPass := full.ProfitFactor >= cfg.FullProfitFactorMin && oos.ProfitFactor >= cfg.SplitProfitFactorMin && recent.ProfitFactor >= cfg.SplitProfitFactorMin
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit && stress.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && oos.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && recent.MaxDrawdown <= cfg.SplitMaxDrawdownLimit
	sideReportingPass := true
	for _, split := range DefaultSplits() {
		_, hasAll := byKey[split.Name+"|all"]
		_, hasLong := byKey[split.Name+"|"+string(Long)]
		_, hasShort := byKey[split.Name+"|"+string(Short)]
		if !hasAll || !hasLong || !hasShort {
			sideReportingPass = false
		}
	}
	return BTC15MSessionOpeningRangeExpansionFalsification{
		SourceResamplePass:                 source.SourceFactsPass && coverage.SourceResamplePass,
		FixedSessionSpecPass:               true,
		LeakagePass:                        true,
		TradeCountPass:                     tradeCountPass,
		GrossEdgePass:                      grossEdgePass,
		NetEdgePass:                        netEdgePass,
		ProfitFactorPass:                   profitFactorPass,
		DrawdownPass:                       drawdownPass,
		SideReportingPass:                  sideReportingPass,
		CombinedBaselineSelectionPass:      true,
		OptimizerContaminationPass:         true,
		DerivativesVetoContaminationPass:   true,
		FullExecutedTrades:                 full.TotalTrades,
		RequiredFullExecutedTrades:         cfg.MinFullTrades,
		MinimumPrimarySplitExecutedTrades:  minPrimary,
		RequiredPrimarySplitExecutedTrades: cfg.MinSplitTrades,
		FullGrossPnL:                       full.GrossPnL,
		FullNetPnL:                         full.NetPnL,
		FullProfitFactor:                   full.ProfitFactor,
		FullMaxDrawdown:                    full.MaxDrawdown,
		StressGrossPnL:                     stress.GrossPnL,
		StressNetPnL:                       stress.NetPnL,
		StressProfitFactor:                 stress.ProfitFactor,
		StressMaxDrawdown:                  stress.MaxDrawdown,
		OOSGrossPnL:                        oos.GrossPnL,
		OOSNetPnL:                          oos.NetPnL,
		OOSProfitFactor:                    oos.ProfitFactor,
		OOSMaxDrawdown:                     oos.MaxDrawdown,
		RecentGrossPnL:                     recent.GrossPnL,
		RecentNetPnL:                       recent.NetPnL,
		RecentProfitFactor:                 recent.ProfitFactor,
		RecentMaxDrawdown:                  recent.MaxDrawdown,
		LongFullExecutedTrades:             longFull.TotalTrades,
		LongFullGrossPnL:                   longFull.GrossPnL,
		LongFullNetPnL:                     longFull.NetPnL,
		LongFullProfitFactor:               longFull.ProfitFactor,
		ShortFullExecutedTrades:            shortFull.TotalTrades,
		ShortFullGrossPnL:                  shortFull.GrossPnL,
		ShortFullNetPnL:                    shortFull.NetPnL,
		ShortFullProfitFactor:              shortFull.ProfitFactor,
	}
}

func btc15MSessionOpeningRangeFalsification(report BTC15MSessionOpeningRangeExpansionFalsification, cfg BacktestFirstBTC15MSessionOpeningRangeExpansionConfig) BTC15MSessionOpeningRangeExpansionFalsification {
	report.BacktestName = BacktestFirstBTC15MSessionOpeningRangeExpansionName
	report.CandidateID = BTC15MSessionOpeningRangeExpansionCandidateID
	report.RequiredFullExecutedTrades = cfg.MinFullTrades
	report.RequiredPrimarySplitExecutedTrades = cfg.MinSplitTrades
	failures := []string{}
	if !report.SourceResamplePass {
		failures = append(failures, "source_or_resample_failed")
	}
	if !report.FixedSessionSpecPass {
		failures = append(failures, "fixed_session_spec_failed")
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
	if !report.SideReportingPass {
		failures = append(failures, "side_reporting_failed")
	}
	if !report.CombinedBaselineSelectionPass {
		failures = append(failures, "post_result_side_selection")
	}
	if !report.OptimizerContaminationPass {
		failures = append(failures, "optimizer_contamination")
	}
	if !report.DerivativesVetoContaminationPass {
		failures = append(failures, "derivatives_veto_contamination")
	}
	report.FailureReasons = failures
	switch {
	case !report.SourceResamplePass || !report.FixedSessionSpecPass || !report.LeakagePass:
		report.StopState = BTC15MSessionOpeningRangeExpansionStopStateFailedSourceOrResample
	case report.TradeCountPass &&
		report.GrossEdgePass &&
		report.NetEdgePass &&
		report.ProfitFactorPass &&
		report.DrawdownPass &&
		report.SideReportingPass &&
		report.CombinedBaselineSelectionPass &&
		report.OptimizerContaminationPass &&
		report.DerivativesVetoContaminationPass:
		report.StopState = BTC15MSessionOpeningRangeExpansionStopStatePassedNeedsReview
	default:
		report.StopState = BTC15MSessionOpeningRangeExpansionStopStateFailedNoUsableStrategy
	}
	return report
}

func btc15MSessionOpeningRangeAddSkip(skips map[btc15MSessionOpeningRangeSkipKey]int, split, side, reason string) {
	if split == "" {
		split = fullSplitName
	}
	if side == "" {
		side = "all"
	}
	skips[btc15MSessionOpeningRangeSkipKey{split: split, side: side, reason: reason}]++
	if split != fullSplitName {
		skips[btc15MSessionOpeningRangeSkipKey{split: fullSplitName, side: side, reason: reason}]++
	}
}

func btc15MSessionOpeningRangeSkipRows(skips map[btc15MSessionOpeningRangeSkipKey]int) []BTC15MSessionOpeningRangeExpansionSkipRow {
	keys := make([]btc15MSessionOpeningRangeSkipKey, 0, len(skips))
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
	rows := make([]BTC15MSessionOpeningRangeExpansionSkipRow, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, BTC15MSessionOpeningRangeExpansionSkipRow{
			CandidateID:       BTC15MSessionOpeningRangeExpansionCandidateID,
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

func sameUTCDate(t, date time.Time) bool {
	t = t.UTC()
	date = date.UTC()
	return t.Year() == date.Year() && t.Month() == date.Month() && t.Day() == date.Day()
}
