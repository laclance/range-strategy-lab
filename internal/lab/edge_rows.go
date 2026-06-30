package lab

type BTC15MRangeEdgeExhaustionFadeSourceRow struct {
	BacktestName               string `json:"backtest_name"`
	CandidateID                string `json:"candidate_id"`
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

type BTC15MRangeEdgeExhaustionFadeCoverageRow struct {
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

type BTC15MRangeEdgeExhaustionFadeSignalRow struct {
	SignalID                  string    `json:"signal_id"`
	CandidateID               string    `json:"candidate_id"`
	Split                     string    `json:"split"`
	DecisionIndex             int       `json:"decision_index"`
	DecisionCloseTime         string    `json:"decision_close_time"`
	DecisionClose             float64   `json:"decision_close"`
	Side                      Direction `json:"side"`
	TimingLabel               string    `json:"timing_label"`
	LookbackBars              int       `json:"lookback_bars"`
	RangeHigh                 float64   `json:"range_high"`
	RangeLow                  float64   `json:"range_low"`
	RangeMidpoint             float64   `json:"range_midpoint"`
	RangeWidth                float64   `json:"range_width"`
	StartClosePosition        float64   `json:"start_close_position"`
	FinalProgressATR          float64   `json:"final_progress_atr"`
	PriorATR14                float64   `json:"prior_atr14"`
	Stop                      float64   `json:"stop"`
	Target                    float64   `json:"target"`
	MaxHoldBars               int       `json:"max_hold_bars"`
	Executed                  bool      `json:"executed"`
	ForwardLabelsAsInput      bool      `json:"forward_labels_as_input"`
	UsesFutureRowsForFeatures bool      `json:"uses_future_rows_for_features"`
	DerivativesVetoUsed       bool      `json:"derivatives_veto_used"`
	OptimizerSelectionUsed    bool      `json:"optimizer_selection_used"`
}

type BTC15MRangeEdgeExhaustionFadeSkipRow struct {
	CandidateID       string `json:"candidate_id"`
	Split             string `json:"split"`
	Reason            string `json:"reason"`
	Count             int    `json:"count"`
	MissingDataPolicy string `json:"missing_data_policy"`
	ForwardFilledRows int    `json:"forward_filled_rows"`
}

type BTC15MRangeEdgeExhaustionFadeTradeRow struct {
	SignalID   string    `json:"signal_id"`
	CloseSplit string    `json:"close_split"`
	Side       Direction `json:"side"`
	EntryTime  string    `json:"entry_time"`
	ExitTime   string    `json:"exit_time"`
	GrossPnL   float64   `json:"gross_pnl"`
	NetPnL     float64   `json:"net_pnl"`
	Fees       float64   `json:"fees"`
	Slippage   float64   `json:"slippage"`
	ExitReason string    `json:"exit_reason"`
	HoldBars   int       `json:"hold_bars"`
}

type BTC15MRangeEdgeExhaustionFadeFalsification struct {
	BacktestName                      string   `json:"backtest_name"`
	CandidateID                       string   `json:"candidate_id"`
	StopState                         string   `json:"stop_state"`
	SourceResamplePass                bool     `json:"source_resample_pass"`
	LeakagePass                       bool     `json:"leakage_pass"`
	TradeCountPass                    bool     `json:"trade_count_pass"`
	GrossEdgePass                     bool     `json:"gross_edge_pass"`
	NetEdgePass                       bool     `json:"net_edge_pass"`
	DrawdownPass                      bool     `json:"drawdown_pass"`
	RobustnessPass                    bool     `json:"robustness_pass"`
	SideReportingPass                 bool     `json:"side_reporting_pass"`
	FullExecutedTrades                int      `json:"full_executed_trades"`
	RequiredFullExecutedTrades        int      `json:"required_full_executed_trades"`
	MinimumPrimarySplitExecutedTrades int      `json:"minimum_primary_split_executed_trades"`
	RequiredPrimarySplitTrades        int      `json:"required_primary_split_trades"`
	FullGrossPnL                      float64  `json:"full_gross_pnl"`
	FullNetPnL                        float64  `json:"full_net_pnl"`
	FullProfitFactor                  float64  `json:"full_profit_factor"`
	FullMaxDrawdown                   float64  `json:"full_max_drawdown"`
	DominantPrimarySplitTradeShare    float64  `json:"dominant_primary_split_trade_share"`
	FailureReasons                    []string `json:"failure_reasons,omitempty"`
}
