package lab

const (
	RangeOptimizationWorkbenchName = "range_optimization_workbench_v1"

	RangeOptimizationWorkbenchStopStateImplementationAddedNeedsLocalRun = "range_optimization_workbench_implementation_added_needs_local_run"
	RangeOptimizationWorkbenchStopStateFailedNoCandidate               = "range_optimization_workbench_failed_no_candidate"
	RangeOptimizationWorkbenchStopStateCandidateSelectedNeedsValidation = "range_optimization_workbench_candidate_selected_needs_fixed_validation"
	RangeOptimizationWorkbenchStopStateRejectedOverfitRisk             = "range_optimization_workbench_rejected_overfit_risk"
	RangeOptimizationWorkbenchStopStateFailedSourceOrResample          = "range_optimization_workbench_failed_source_or_resample"
)

type RangeOptimizationWorkbenchConfig struct {
	ApprovedSourcePath      string
	ExpectedSourceRows      int
	ExpectedFirstOpenTime   string
	ExpectedLastOpenTime    string
	ExpectedZeroVolumeCount int
	Expected15MRows         int
	Expected15MLastOpenTime string
	MaxTrials               int
}

type RangeOptimizationWorkbenchResult struct {
	RunID             string                                  `json:"run_id"`
	SourceRows        []RangeOptimizationWorkbenchSourceRow   `json:"source_rows"`
	CoverageRows      []RangeOptimizationWorkbenchCoverageRow `json:"coverage_rows"`
	GridRows          []RangeOptimizationWorkbenchTrialSpec   `json:"grid_rows"`
	TrialRows         []RangeOptimizationWorkbenchTrialRow    `json:"trial_rows"`
	TrialSummaryRows  []RangeOptimizationWorkbenchSummaryRow  `json:"trial_summary_rows"`
	TopCandidates     []RangeOptimizationWorkbenchCandidate   `json:"top_candidates"`
	RejectedCandidates []RangeOptimizationWorkbenchRejected    `json:"rejected_candidates"`
	RobustnessSummary RangeOptimizationWorkbenchRobustness    `json:"robustness_summary"`
	Falsification     RangeOptimizationWorkbenchFalsification `json:"falsification"`
	StopState         string                                  `json:"stop_state"`
}

type RangeOptimizationWorkbenchSourceRow struct {
	BacktestName               string `json:"backtest_name"`
	RunID                      string `json:"run_id"`
	Path                       string `json:"path"`
	ApprovedPath               string `json:"approved_path"`
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

type RangeOptimizationWorkbenchCoverageRow struct {
	BacktestName         string `json:"backtest_name"`
	RunID                string `json:"run_id"`
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

type RangeOptimizationWorkbenchTrialSpec struct {
	TrialID               string  `json:"trial_id"`
	FamilyID              string  `json:"family_id"`
	Timeframe             string  `json:"timeframe"`
	EntryArchetype        string  `json:"entry_archetype"`
	RangeLookback         int     `json:"range_lookback"`
	EdgeZonePct           float64 `json:"edge_zone_pct"`
	InteriorLowPct        float64 `json:"interior_low_pct"`
	InteriorHighPct       float64 `json:"interior_high_pct"`
	VWAPDistanceRangePct  float64 `json:"vwap_distance_range_pct"`
	ProgressATRMultiple   float64 `json:"progress_atr_multiple"`
	MinRangeATR           float64 `json:"min_range_atr"`
	MaxRangeATR           float64 `json:"max_range_atr"`
	ATRPercentileMin      string  `json:"atr_percentile_min"`
	ATRPercentileMax      string  `json:"atr_percentile_max"`
	TrendMode             string  `json:"trend_mode"`
	EMALength             int     `json:"ema_length"`
	ImpulseATRMin         string  `json:"impulse_atr_min"`
	ImpulseATRMax         string  `json:"impulse_atr_max"`
	VolumeMode            string  `json:"volume_mode"`
	VolumeLookback        int     `json:"volume_lookback"`
	TargetMode            string  `json:"target_mode"`
	StopMode              string  `json:"stop_mode"`
	TimeStopBars          int     `json:"time_stop_bars"`
}

type RangeOptimizationWorkbenchTrialRow struct {
	RangeOptimizationWorkbenchTrialSpec
	RunID                         string  `json:"run_id"`
	SourceRef                     string  `json:"source_ref"`
	FullTrades                    int     `json:"full_trades"`
	StressTrades                  int     `json:"stress_trades"`
	OOSTrades                     int     `json:"oos_trades"`
	RecentTrades                  int     `json:"recent_trades"`
	FullGrossPnL                  float64 `json:"full_gross_pnl"`
	FullNetPnL                    float64 `json:"full_net_pnl"`
	FullProfitFactor              float64 `json:"full_profit_factor"`
	FullMaxDrawdown               float64 `json:"full_max_drawdown"`
	StressGrossPnL                float64 `json:"stress_gross_pnl"`
	StressNetPnL                  float64 `json:"stress_net_pnl"`
	OOSGrossPnL                   float64 `json:"oos_gross_pnl"`
	OOSNetPnL                     float64 `json:"oos_net_pnl"`
	RecentGrossPnL                float64 `json:"recent_gross_pnl"`
	RecentNetPnL                  float64 `json:"recent_net_pnl"`
	LongTrades                    int     `json:"long_trades"`
	ShortTrades                   int     `json:"short_trades"`
	LongNetPnL                    float64 `json:"long_net_pnl"`
	ShortNetPnL                   float64 `json:"short_net_pnl"`
	MinimumPrimarySplitTrades     int     `json:"minimum_primary_split_trades"`
	DominantPrimarySplitTradeShare float64 `json:"dominant_primary_split_trade_share"`
	RobustnessScore               float64 `json:"robustness_score"`
	FailureReasons                string  `json:"failure_reasons"`
	SelectedForLocking            bool    `json:"selected_for_locking"`
}

type RangeOptimizationWorkbenchSummaryRow struct {
	RunID       string  `json:"run_id"`
	TrialID     string  `json:"trial_id"`
	Split       string  `json:"split"`
	Side        string  `json:"side"`
	TotalTrades int     `json:"total_trades"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
	GrossPnL    float64 `json:"gross_pnl"`
	NetPnL      float64 `json:"net_pnl"`
	TotalCosts  float64 `json:"total_costs"`
	ProfitFactor float64 `json:"profit_factor"`
	GrossProfitFactor float64 `json:"gross_profit_factor"`
	MaxDrawdown float64 `json:"max_drawdown"`
	Expectancy  float64 `json:"expectancy"`
	AvgHoldBars float64 `json:"avg_hold_bars"`
}

type RangeOptimizationWorkbenchCandidate struct {
	RunID           string  `json:"run_id"`
	Rank            int     `json:"rank"`
	TrialID         string  `json:"trial_id"`
	CandidateID     string  `json:"candidate_id"`
	FamilyID        string  `json:"family_id"`
	Timeframe       string  `json:"timeframe"`
	EntryArchetype  string  `json:"entry_archetype"`
	RobustnessScore float64 `json:"robustness_score"`
	FullTrades      int     `json:"full_trades"`
	FullGrossPnL    float64 `json:"full_gross_pnl"`
	FullNetPnL      float64 `json:"full_net_pnl"`
	FullProfitFactor float64 `json:"full_profit_factor"`
	FullMaxDrawdown float64 `json:"full_max_drawdown"`
	SelectionReason string  `json:"selection_reason"`
}

type RangeOptimizationWorkbenchRejected struct {
	RunID          string `json:"run_id"`
	TrialID        string `json:"trial_id"`
	FamilyID       string `json:"family_id"`
	Timeframe      string `json:"timeframe"`
	EntryArchetype string `json:"entry_archetype"`
	FailureReasons string `json:"failure_reasons"`
}

type RangeOptimizationWorkbenchRobustness struct {
	RunID              string `json:"run_id"`
	TotalTrials        int    `json:"total_trials"`
	MaxTrials          int    `json:"max_trials"`
	PassingCandidates  int    `json:"passing_candidates"`
	RejectedCandidates int    `json:"rejected_candidates"`
	SelectedTrialID     string `json:"selected_trial_id,omitempty"`
	SelectedCandidateID string `json:"selected_candidate_id,omitempty"`
	StopState           string `json:"stop_state"`
}

type RangeOptimizationWorkbenchFalsification struct {
	BacktestName        string   `json:"backtest_name"`
	RunID               string   `json:"run_id"`
	StopState           string   `json:"stop_state"`
	SourceResamplePass  bool     `json:"source_resample_pass"`
	TotalTrials         int      `json:"total_trials"`
	MaxTrials           int      `json:"max_trials"`
	PassingCandidates   int      `json:"passing_candidates"`
	SelectedTrialID      string   `json:"selected_trial_id,omitempty"`
	SelectedCandidateID  string   `json:"selected_candidate_id,omitempty"`
	FailureReasons      []string `json:"failure_reasons,omitempty"`
}

func DefaultRangeOptimizationWorkbenchConfig() RangeOptimizationWorkbenchConfig {
	return RangeOptimizationWorkbenchConfig{
		ApprovedSourcePath:      "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedSourceRows:      573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedZeroVolumeCount: 66,
		Expected15MRows:         191328,
		Expected15MLastOpenTime: "2026-06-16T23:45:00Z",
		MaxTrials:               2500,
	}
}
