package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	FuturesBTC15MPostCompressionDirectionalExpansionAuditName = "futures_btc_15m_post_compression_directional_expansion_zero_trade_audit"
	BTC15MPostCompressionDirectionalExpansionPremiseID        = "btc_15m_post_compression_directional_expansion_v1"

	BTC15MPostCompressionDirectionalExpansionStopStatePassedNeedsReview = "btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review"
	BTC15MPostCompressionDirectionalExpansionStopStateFailedNoPremise   = "btc_15m_post_compression_directional_expansion_zero_trade_audit_failed_no_usable_entry_premise"
	BTC15MPostCompressionDirectionalExpansionStopStateClosedReslice     = "independent_entry_premise_spec_rejected_closed_family_reslice"

	BTC15MPostCompressionSideLong  = "long"
	BTC15MPostCompressionSideShort = "short"

	BTC15MPostCompressionVolumeNone       = "none"
	BTC15MPostCompressionVolumeMedian96   = "above_prior_96_median"
	BTC15MPostCompressionVolumeP60Prior96 = "above_prior_96_p60"

	btc15MPostCompressionSourcePath          = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"
	btc15MPostCompressionExpectedRows        = 573984
	btc15MPostCompressionExpectedFirst       = "2021-01-01T00:00:00Z"
	btc15MPostCompressionExpectedLast        = "2026-06-16T23:55:00Z"
	btc15MPostCompressionExpectedZeroVol     = 66
	btc15MPostCompressionExpected15MRows     = 191328
	btc15MPostCompressionExpected15MLastOpen = "2026-06-16T23:45:00Z"
)

type FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig struct {
	ApprovedSourcePath       string
	ExpectedSourceRows       int
	ExpectedFirstOpenTime    string
	ExpectedLastOpenTime     string
	ExpectedGapCount         int
	ExpectedDuplicateCount   int
	ExpectedZeroVolumeCount  int
	SkipSourceFactCheck      bool
	Expected15MRows          int
	Expected15MLastOpenTime  string
	SkipCoverageCountCheck   bool
	CompressionLookbacks     []int
	CompressionPercentiles   []float64
	PercentileReferenceBars  int
	BreakoutATRMultiples     []float64
	VolumeModes              []string
	VolumeLookbackBars       int
	ATRPeriod                int
	HorizonsBars             []int
	MinFullDedupCandidates   int
	MinSplitDedupCandidates  int
	ClosedFamilyReslice      bool
	DerivativesVetoContam    bool
	ForwardLabelLeakOverride bool
}

type FuturesBTC15MPostCompressionDirectionalExpansionAuditResult struct {
	SourceRows       []BTC15MPostCompressionSourceRow       `json:"source_rows"`
	CoverageRows     []BTC15MPostCompressionCoverageRow     `json:"coverage_rows"`
	ParameterCells   []BTC15MPostCompressionParameterCell   `json:"parameter_cells"`
	CandidateRows    []BTC15MPostCompressionCandidateRow    `json:"candidate_rows"`
	DedupEvents      []BTC15MPostCompressionDedupEventRow   `json:"dedup_events"`
	BaselineRows     []BTC15MPostCompressionBaselineRow     `json:"baseline_rows"`
	SplitSummaryRows []BTC15MPostCompressionSplitSummaryRow `json:"split_summary_rows"`
	AdjacencyRows    []BTC15MPostCompressionAdjacencyRow    `json:"adjacency_rows"`
	MissingnessRows  []BTC15MPostCompressionMissingnessRow  `json:"missingness_rows"`
	Falsification    BTC15MPostCompressionFalsification     `json:"falsification"`
	PassingCells     int                                    `json:"passing_cells"`
	StopState        string                                 `json:"stop_state"`
}

type BTC15MPostCompressionSourceRow struct {
	AuditName                  string `json:"audit_name"`
	PremiseID                  string `json:"premise_id"`
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

type BTC15MPostCompressionCoverageRow struct {
	AuditName             string `json:"audit_name"`
	PremiseID             string `json:"premise_id"`
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

type BTC15MPostCompressionParameterCell struct {
	CellID                   string  `json:"cell_id"`
	LookbackBars             int     `json:"lookback_bars"`
	CompressionPercentile    float64 `json:"compression_percentile"`
	PercentileReferenceBars  int     `json:"percentile_reference_bars"`
	BreakoutATRMultiple      float64 `json:"breakout_atr_multiple"`
	VolumeMode               string  `json:"volume_mode"`
	VolumeLookbackBars       int     `json:"volume_lookback_bars"`
	CandidateRows            int     `json:"candidate_rows"`
	LabelRows                int     `json:"label_rows"`
	DedupEventRows           int     `json:"dedup_event_rows"`
	LongCandidateRows        int     `json:"long_candidate_rows"`
	ShortCandidateRows       int     `json:"short_candidate_rows"`
	PassingSideHorizonCount  int     `json:"passing_side_horizon_count"`
	AdjacentPassingCellCount int     `json:"adjacent_passing_cell_count"`
	Predeclared              bool    `json:"predeclared"`
}

type BTC15MPostCompressionCandidateRow struct {
	CandidateRowID                 int     `json:"candidate_row_id"`
	BaseEventID                    int     `json:"base_event_id"`
	CellID                         string  `json:"cell_id"`
	Split                          string  `json:"split"`
	Timeframe                      string  `json:"timeframe"`
	DecisionIndex                  int     `json:"decision_index"`
	DecisionOpenTime               string  `json:"decision_open_time"`
	DecisionCloseTime              string  `json:"decision_close_time"`
	Side                           string  `json:"side"`
	TimingLabel                    string  `json:"timing_label"`
	LookbackBars                   int     `json:"lookback_bars"`
	CompressionPercentile          float64 `json:"compression_percentile"`
	CompressionThreshold           float64 `json:"compression_threshold"`
	RangeHigh                      float64 `json:"range_high"`
	RangeLow                       float64 `json:"range_low"`
	RangeWidthPct                  float64 `json:"range_width_pct"`
	BreakoutATRMultiple            float64 `json:"breakout_atr_multiple"`
	PriorATR14                     float64 `json:"prior_atr14"`
	BreakoutDistance               float64 `json:"breakout_distance"`
	ExpansionDirection             string  `json:"expansion_direction"`
	DecisionClose                  float64 `json:"decision_close"`
	DecisionVolume                 float64 `json:"decision_volume"`
	VolumeMode                     string  `json:"volume_mode"`
	VolumeThreshold                float64 `json:"volume_threshold"`
	VolumeLookbackBars             int     `json:"volume_lookback_bars"`
	HorizonBars                    int     `json:"horizon_bars"`
	LabelAnchorOpenTime            string  `json:"label_anchor_open_time"`
	LabelAnchorOpen                float64 `json:"label_anchor_open"`
	LabelWindowEndIndex            int     `json:"label_window_end_index"`
	LabelWindowEndTime             string  `json:"label_window_end_time"`
	IntendedSideForwardCloseReturn float64 `json:"intended_side_forward_close_return_bp"`
	IntendedSideFavorable          float64 `json:"intended_side_favorable_bp"`
	IntendedSideAdverse            float64 `json:"intended_side_adverse_bp"`
	FavorableGreaterThanAdverse    bool    `json:"favorable_gt_adverse"`
	FavorableMinusAdverse          float64 `json:"favorable_minus_adverse_bp"`
	ForwardLabelMetadataOnly       bool    `json:"forward_label_metadata_only"`
	ForwardLabelUsedAsFeature      bool    `json:"forward_label_used_as_feature"`
	UsesFutureRowsForFeatures      bool    `json:"uses_future_rows_for_features"`
	DerivativesVetoUsed            bool    `json:"derivatives_veto_used"`
}

type BTC15MPostCompressionDedupEventRow struct {
	DedupEventID      int    `json:"dedup_event_id"`
	Split             string `json:"split"`
	Timeframe         string `json:"timeframe"`
	DecisionIndex     int    `json:"decision_index"`
	DecisionCloseTime string `json:"decision_close_time"`
	Side              string `json:"side"`
	MatchedCellCount  int    `json:"matched_cell_count"`
	MatchedCellIDs    string `json:"matched_cell_ids"`
}

type BTC15MPostCompressionBaselineRow struct {
	Split                              string  `json:"split"`
	Side                               string  `json:"side"`
	HorizonBars                        int     `json:"horizon_bars"`
	EligibleRows                       int     `json:"eligible_rows"`
	MeanForwardCloseReturnBP           float64 `json:"mean_intended_side_forward_close_return_bp"`
	MeanFavorableBP                    float64 `json:"mean_intended_side_favorable_bp"`
	MeanAdverseBP                      float64 `json:"mean_intended_side_adverse_bp"`
	MeanFavorableMinusAdverseBP        float64 `json:"mean_favorable_minus_adverse_bp"`
	FavorableGreaterThanAdverseRate    float64 `json:"favorable_gt_adverse_rate"`
	AuditWarmupBars                    int     `json:"audit_warmup_bars"`
	ForwardLabelsAsBaselineInput       bool    `json:"forward_labels_as_baseline_input"`
	UnconditionalEligible15MBaseline   bool    `json:"unconditional_eligible_15m_baseline"`
	SharesCandidateWarmupAndLabelRules bool    `json:"shares_candidate_warmup_and_label_rules"`
}

type BTC15MPostCompressionSplitSummaryRow struct {
	CellID                         string  `json:"cell_id"`
	Split                          string  `json:"split"`
	Side                           string  `json:"side"`
	HorizonBars                    int     `json:"horizon_bars"`
	LookbackBars                   int     `json:"lookback_bars"`
	CompressionPercentile          float64 `json:"compression_percentile"`
	BreakoutATRMultiple            float64 `json:"breakout_atr_multiple"`
	VolumeMode                     string  `json:"volume_mode"`
	CandidateRows                  int     `json:"candidate_rows"`
	BaselineRows                   int     `json:"baseline_rows"`
	MeanForwardCloseReturnBP       float64 `json:"mean_intended_side_forward_close_return_bp"`
	BaselineForwardCloseReturnBP   float64 `json:"baseline_intended_side_forward_close_return_bp"`
	DeltaForwardCloseReturnBP      float64 `json:"delta_intended_side_forward_close_return_bp"`
	MeanFavorableMinusAdverseBP    float64 `json:"mean_favorable_minus_adverse_bp"`
	BaselineFavorableMinusAdverse  float64 `json:"baseline_favorable_minus_adverse_bp"`
	DeltaFavorableMinusAdverseBP   float64 `json:"delta_favorable_minus_adverse_bp"`
	FavorableGreaterThanAdverse    float64 `json:"favorable_gt_adverse_rate"`
	BaselineFavorableGTAdverseRate float64 `json:"baseline_favorable_gt_adverse_rate"`
	DeltaFavorableGTAdverseRate    float64 `json:"delta_favorable_gt_adverse_rate"`
	ForwardReturnGatePass          bool    `json:"forward_return_gate_pass"`
	FavorableMinusAdverseGatePass  bool    `json:"favorable_minus_adverse_gate_pass"`
	FavorableGTAdverseGatePass     bool    `json:"favorable_gt_adverse_gate_pass"`
	PassesBaseline                 bool    `json:"passes_baseline"`
	FullAndEveryPrimarySplitPass   bool    `json:"full_and_every_primary_split_pass"`
}

type BTC15MPostCompressionAdjacencyRow struct {
	CellID                        string `json:"cell_id"`
	Side                          string `json:"side"`
	HorizonBars                   int    `json:"horizon_bars"`
	FullAndEveryPrimarySplitPass  bool   `json:"full_and_every_primary_split_pass"`
	AdjacentPassingCellIDs        string `json:"adjacent_passing_cell_ids"`
	AdjacentPassingCellCount      int    `json:"adjacent_passing_cell_count"`
	PassesAdjacentClusterGate     bool   `json:"passes_adjacent_cluster_gate"`
	AdjacencyDefinition           string `json:"adjacency_definition"`
	OneIsolatedPassingCellFailure bool   `json:"one_isolated_passing_cell_failure"`
}

type BTC15MPostCompressionMissingnessRow struct {
	AuditName         string  `json:"audit_name"`
	PremiseID         string  `json:"premise_id"`
	Scope             string  `json:"scope"`
	Reason            string  `json:"reason"`
	Count             int     `json:"count"`
	TotalRows         int     `json:"total_rows"`
	Rate              float64 `json:"rate"`
	MissingDataPolicy string  `json:"missing_data_policy"`
	ForwardFilledRows int     `json:"forward_filled_rows"`
}

type BTC15MPostCompressionFalsification struct {
	AuditName                              string   `json:"audit_name"`
	PremiseID                              string   `json:"premise_id"`
	StopState                              string   `json:"stop_state"`
	SourceResamplePass                     bool     `json:"source_resample_pass"`
	LeakagePass                            bool     `json:"leakage_pass"`
	CandidateSizePass                      bool     `json:"candidate_size_pass"`
	SplitSizePass                          bool     `json:"split_size_pass"`
	BaselineSeparationPass                 bool     `json:"baseline_separation_pass"`
	AdjacentCellClusterPass                bool     `json:"adjacent_cell_cluster_pass"`
	SplitStabilityPass                     bool     `json:"split_stability_pass"`
	ClosedFamilyProtectionPass             bool     `json:"closed_family_protection_pass"`
	DerivativesVetoContaminationPass       bool     `json:"derivatives_veto_contamination_pass"`
	CommonOutputsZeroTrade                 bool     `json:"common_outputs_zero_trade"`
	Trades                                 int      `json:"trades"`
	FullDeduplicatedCandidateEvents        int      `json:"full_deduplicated_candidate_events"`
	MinimumPrimarySplitDedupEvents         int      `json:"minimum_primary_split_dedup_events"`
	RequiredFullDeduplicatedCandidateRows  int      `json:"required_full_deduplicated_candidate_rows"`
	RequiredSplitDeduplicatedCandidateRows int      `json:"required_split_deduplicated_candidate_rows"`
	PassingCellSideHorizons                int      `json:"passing_cell_side_horizons"`
	AdjacentPassingCellSideHorizons        int      `json:"adjacent_passing_cell_side_horizons"`
	FailureReasons                         []string `json:"failure_reasons,omitempty"`
}

type btcPostCompressionRangeFacts struct {
	lookback   int
	high       []float64
	low        []float64
	widthPct   []float64
	thresholds map[float64][]float64
}

type btcPostCompressionBaseEvent struct {
	eventID               int
	cellID                string
	split                 string
	decisionIndex         int
	side                  string
	lookback              int
	compressionPercentile float64
	compressionThreshold  float64
	rangeHigh             float64
	rangeLow              float64
	rangeWidthPct         float64
	breakoutATRMultiple   float64
	priorATR              float64
	breakoutDistance      float64
	expansionDirection    string
	volumeMode            string
	volumeThreshold       float64
}

type btcPostCompressionForwardLabel struct {
	horizonBars             int
	anchorOpenTime          string
	anchorOpen              float64
	windowEndIndex          int
	windowEndTime           string
	forwardCloseReturnBP    float64
	favorableBP             float64
	adverseBP               float64
	favorableGreaterAdverse bool
	favorableMinusAdverseBP float64
}

type btcPostCompressionMetricAccumulator struct {
	count           int
	sumForward      float64
	sumFavorable    float64
	sumAdverse      float64
	sumFavMinusAdv  float64
	favGreaterCount int
}

type btcPostCompressionSummaryKey struct {
	cellID      string
	split       string
	side        string
	horizonBars int
}

type btcPostCompressionBaselineKey struct {
	split       string
	side        string
	horizonBars int
}

type btcPostCompressionPassingKey struct {
	cellID      string
	side        string
	horizonBars int
}

type btcPostCompressionDedupKey struct {
	decisionIndex int
	side          string
}

type btcPostCompressionDedupAccumulator struct {
	split             string
	decisionCloseTime string
	cells             map[string]struct{}
}

func DefaultFuturesBTC15MPostCompressionDirectionalExpansionAuditConfig() FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig {
	return FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig{
		ApprovedSourcePath:      btc15MPostCompressionSourcePath,
		ExpectedSourceRows:      btc15MPostCompressionExpectedRows,
		ExpectedFirstOpenTime:   btc15MPostCompressionExpectedFirst,
		ExpectedLastOpenTime:    btc15MPostCompressionExpectedLast,
		ExpectedGapCount:        0,
		ExpectedDuplicateCount:  0,
		ExpectedZeroVolumeCount: btc15MPostCompressionExpectedZeroVol,
		Expected15MRows:         btc15MPostCompressionExpected15MRows,
		Expected15MLastOpenTime: btc15MPostCompressionExpected15MLastOpen,
		CompressionLookbacks:    []int{48, 96, 192},
		CompressionPercentiles:  []float64{0.20, 0.30, 0.40},
		PercentileReferenceBars: 1920,
		BreakoutATRMultiples:    []float64{0.1, 0.2, 0.3},
		VolumeModes:             []string{BTC15MPostCompressionVolumeNone, BTC15MPostCompressionVolumeMedian96, BTC15MPostCompressionVolumeP60Prior96},
		VolumeLookbackBars:      96,
		ATRPeriod:               14,
		HorizonsBars:            []int{16, 32, 48},
		MinFullDedupCandidates:  300,
		MinSplitDedupCandidates: 50,
	}
}

func RunFuturesBTC15MPostCompressionDirectionalExpansionAudit(candles []Candle, manifest SourceManifest, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, splits []Split) (FuturesBTC15MPostCompressionDirectionalExpansionAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesBTC15MPostCompressionDirectionalExpansionAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	sourceRow := btcPostCompressionSourceRow(manifest, cfg)
	result := FuturesBTC15MPostCompressionDirectionalExpansionAuditResult{
		SourceRows:     []BTC15MPostCompressionSourceRow{sourceRow},
		ParameterCells: btcPostCompressionParameterCells(cfg),
	}
	coveragePass := false
	frame := btcPostCompression15MFrame()
	resampled, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
	coverageRow := btcPostCompressionCoverageRow(coverage, cfg, sourceRow.SourceFactsPass)
	if err != nil {
		coverageRow.ValidationStatus = "rejected"
		coverageRow.ValidationError = err.Error()
		coverageRow.SourceResamplePass = false
	} else {
		coveragePass = coverageRow.SourceResamplePass
	}
	result.CoverageRows = []BTC15MPostCompressionCoverageRow{coverageRow}
	if !sourceRow.SourceFactsPass || !coveragePass {
		result.Falsification = btcPostCompressionFalsification(BTC15MPostCompressionFalsification{
			SourceResamplePass:               false,
			LeakagePass:                      !cfg.ForwardLabelLeakOverride,
			ClosedFamilyProtectionPass:       !cfg.ClosedFamilyReslice,
			DerivativesVetoContaminationPass: !cfg.DerivativesVetoContam,
			CommonOutputsZeroTrade:           true,
			Trades:                           0,
		}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}

	result = btcPostCompressionRunFrom15M(resampled, result.SourceRows, result.CoverageRows, cfg, splits)
	return result, nil
}

func btcPostCompressionRunFrom15M(candles []Candle, sourceRows []BTC15MPostCompressionSourceRow, coverageRows []BTC15MPostCompressionCoverageRow, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, splits []Split) FuturesBTC15MPostCompressionDirectionalExpansionAuditResult {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesBTC15MPostCompressionDirectionalExpansionAuditResult{
		SourceRows:     sourceRows,
		CoverageRows:   coverageRows,
		ParameterCells: btcPostCompressionParameterCells(cfg),
	}
	cellsByID := make(map[string]*BTC15MPostCompressionParameterCell, len(result.ParameterCells))
	for i := range result.ParameterCells {
		cellsByID[result.ParameterCells[i].CellID] = &result.ParameterCells[i]
	}

	lookbacks := append([]int(nil), cfg.CompressionLookbacks...)
	sort.Ints(lookbacks)
	rangeFacts := map[int]btcPostCompressionRangeFacts{}
	for _, lookback := range lookbacks {
		rangeFacts[lookback] = btcPostCompressionBuildRangeFacts(candles, lookback, cfg.CompressionPercentiles, cfg.PercentileReferenceBars)
	}
	atr := ATR(candles, cfg.ATRPeriod)
	volumeMedian := btcPostCompressionRollingPriorVolumeThreshold(candles, cfg.VolumeLookbackBars, 0.50)
	volumeP60 := btcPostCompressionRollingPriorVolumeThreshold(candles, cfg.VolumeLookbackBars, 0.60)
	auditWarmup := btcPostCompressionAuditWarmup(cfg)
	maxHorizon := btcPostCompressionMaxInt(cfg.HorizonsBars)

	missing := map[string]int{}
	dedup := map[btcPostCompressionDedupKey]*btcPostCompressionDedupAccumulator{}
	candidateAcc := map[btcPostCompressionSummaryKey]*btcPostCompressionMetricAccumulator{}
	baselineAcc := map[btcPostCompressionBaselineKey]*btcPostCompressionMetricAccumulator{}

	baseEventID := 0
	candidateRowID := 0
	for d := 0; d < len(candles); d++ {
		if d < auditWarmup {
			missing["audit_warmup"]++
			continue
		}
		if d+maxHorizon >= len(candles) {
			missing["missing_max_horizon_future"]++
		}
		for _, lookback := range lookbacks {
			facts := rangeFacts[lookback]
			if d >= len(facts.widthPct) || !validNumber(facts.widthPct[d]) || !validNumber(facts.high[d]) || !validNumber(facts.low[d]) {
				missing["missing_prior_range"]++
				continue
			}
			if d <= 0 || d-1 >= len(atr) || !validNumber(atr[d-1]) || atr[d-1] <= 0 {
				missing["missing_prior_atr"]++
				continue
			}
			for _, percentile := range cfg.CompressionPercentiles {
				thresholds := facts.thresholds[percentile]
				if d >= len(thresholds) || !validNumber(thresholds[d]) {
					missing["missing_compression_reference"]++
					continue
				}
				if facts.widthPct[d] > thresholds[d] {
					continue
				}
				for _, multiple := range cfg.BreakoutATRMultiples {
					up := candles[d].Close >= facts.high[d]+multiple*atr[d-1]
					down := candles[d].Close <= facts.low[d]-multiple*atr[d-1]
					if !up && !down {
						continue
					}
					for _, volumeMode := range cfg.VolumeModes {
						volumeThreshold, ok := btcPostCompressionVolumePass(candles[d].Volume, volumeMode, volumeMedian, volumeP60, d)
						if !ok {
							if volumeMode != BTC15MPostCompressionVolumeNone {
								missing["missing_volume_reference"]++
							}
							continue
						}
						if up {
							baseEventID++
							event := btcPostCompressionBaseEventFor(candles, d, baseEventID, lookback, percentile, thresholds[d], facts, multiple, atr[d-1], volumeMode, volumeThreshold, BTC15MPostCompressionSideLong, splits)
							candidateRowID = btcPostCompressionAddEvent(candles, event, cfg, &result, cellsByID, dedup, candidateAcc, missing, candidateRowID)
						}
						if down {
							baseEventID++
							event := btcPostCompressionBaseEventFor(candles, d, baseEventID, lookback, percentile, thresholds[d], facts, multiple, atr[d-1], volumeMode, volumeThreshold, BTC15MPostCompressionSideShort, splits)
							candidateRowID = btcPostCompressionAddEvent(candles, event, cfg, &result, cellsByID, dedup, candidateAcc, missing, candidateRowID)
						}
					}
				}
			}
		}
	}

	for d := auditWarmup; d < len(candles); d++ {
		for _, horizon := range cfg.HorizonsBars {
			for _, side := range []string{BTC15MPostCompressionSideLong, BTC15MPostCompressionSideShort} {
				label, ok := btcPostCompressionForwardLabelFor(candles, d, side, horizon)
				if !ok {
					continue
				}
				split := splitNameForCloseTime(candles[d].CloseTime, splits)
				for _, splitName := range btcPostCompressionSplitNames(split) {
					key := btcPostCompressionBaselineKey{split: splitName, side: side, horizonBars: horizon}
					btcPostCompressionEnsureBaselineAcc(baselineAcc, key).add(label)
				}
			}
		}
	}

	result.DedupEvents = btcPostCompressionDedupRows(dedup)
	for _, acc := range dedup {
		for cellID := range acc.cells {
			if cell := cellsByID[cellID]; cell != nil {
				cell.DedupEventRows++
			}
		}
	}
	result.BaselineRows = btcPostCompressionBaselineRows(baselineAcc, cfg, splits)
	summaryRows, passing := btcPostCompressionSplitSummaryRows(result.ParameterCells, candidateAcc, baselineAcc, cfg, splits)
	result.SplitSummaryRows = summaryRows
	result.AdjacencyRows = btcPostCompressionAdjacencyRows(result.ParameterCells, passing, cfg)
	result.MissingnessRows = btcPostCompressionMissingnessRows(missing, len(candles))

	for key := range passing {
		if passing[key] {
			result.PassingCells++
			if cell := cellsByID[key.cellID]; cell != nil {
				cell.PassingSideHorizonCount++
			}
		}
	}
	adjacentPassing := 0
	for _, row := range result.AdjacencyRows {
		if row.PassesAdjacentClusterGate {
			adjacentPassing++
			if cell := cellsByID[row.CellID]; cell != nil {
				cell.AdjacentPassingCellCount++
			}
		}
	}
	for i := range result.ParameterCells {
		key := result.ParameterCells[i].CellID
		if cell := cellsByID[key]; cell != nil {
			result.ParameterCells[i] = *cell
		}
	}

	fullDedup, minSplitDedup := btcPostCompressionDedupCounts(result.DedupEvents, splits)
	report := BTC15MPostCompressionFalsification{
		SourceResamplePass:                     btcPostCompressionSourceResamplePass(sourceRows, coverageRows),
		LeakagePass:                            !cfg.ForwardLabelLeakOverride,
		CandidateSizePass:                      fullDedup >= cfg.MinFullDedupCandidates,
		SplitSizePass:                          minSplitDedup >= cfg.MinSplitDedupCandidates,
		BaselineSeparationPass:                 result.PassingCells > 0,
		AdjacentCellClusterPass:                adjacentPassing > 0,
		SplitStabilityPass:                     result.PassingCells > 0,
		ClosedFamilyProtectionPass:             !cfg.ClosedFamilyReslice,
		DerivativesVetoContaminationPass:       !cfg.DerivativesVetoContam,
		CommonOutputsZeroTrade:                 true,
		Trades:                                 0,
		FullDeduplicatedCandidateEvents:        fullDedup,
		MinimumPrimarySplitDedupEvents:         minSplitDedup,
		RequiredFullDeduplicatedCandidateRows:  cfg.MinFullDedupCandidates,
		RequiredSplitDeduplicatedCandidateRows: cfg.MinSplitDedupCandidates,
		PassingCellSideHorizons:                result.PassingCells,
		AdjacentPassingCellSideHorizons:        adjacentPassing,
	}
	result.Falsification = btcPostCompressionFalsification(report, cfg)
	result.StopState = result.Falsification.StopState
	return result
}

func (cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) withDefaults() FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig {
	defaults := DefaultFuturesBTC15MPostCompressionDirectionalExpansionAuditConfig()
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
	if len(cfg.CompressionLookbacks) == 0 {
		cfg.CompressionLookbacks = append([]int(nil), defaults.CompressionLookbacks...)
	}
	if len(cfg.CompressionPercentiles) == 0 {
		cfg.CompressionPercentiles = append([]float64(nil), defaults.CompressionPercentiles...)
	}
	if cfg.PercentileReferenceBars == 0 {
		cfg.PercentileReferenceBars = defaults.PercentileReferenceBars
	}
	if len(cfg.BreakoutATRMultiples) == 0 {
		cfg.BreakoutATRMultiples = append([]float64(nil), defaults.BreakoutATRMultiples...)
	}
	if len(cfg.VolumeModes) == 0 {
		cfg.VolumeModes = append([]string(nil), defaults.VolumeModes...)
	}
	if cfg.VolumeLookbackBars == 0 {
		cfg.VolumeLookbackBars = defaults.VolumeLookbackBars
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.MinFullDedupCandidates == 0 {
		cfg.MinFullDedupCandidates = defaults.MinFullDedupCandidates
	}
	if cfg.MinSplitDedupCandidates == 0 {
		cfg.MinSplitDedupCandidates = defaults.MinSplitDedupCandidates
	}
	return cfg
}

func (cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) validate() error {
	if cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("post-compression audit approved source path is required")
	}
	if cfg.PercentileReferenceBars <= 0 {
		return fmt.Errorf("post-compression audit percentile reference bars must be positive")
	}
	if cfg.VolumeLookbackBars <= 0 {
		return fmt.Errorf("post-compression audit volume lookback bars must be positive")
	}
	if cfg.ATRPeriod <= 0 {
		return fmt.Errorf("post-compression audit ATR period must be positive")
	}
	if cfg.MinFullDedupCandidates <= 0 {
		return fmt.Errorf("post-compression audit full dedup candidate minimum must be positive")
	}
	if cfg.MinSplitDedupCandidates <= 0 {
		return fmt.Errorf("post-compression audit split dedup candidate minimum must be positive")
	}
	for _, lookback := range cfg.CompressionLookbacks {
		if lookback <= 0 {
			return fmt.Errorf("post-compression audit lookbacks must be positive")
		}
	}
	for _, percentile := range cfg.CompressionPercentiles {
		if percentile <= 0 || percentile >= 1 {
			return fmt.Errorf("post-compression audit compression percentiles must be between 0 and 1")
		}
	}
	for _, multiple := range cfg.BreakoutATRMultiples {
		if multiple <= 0 {
			return fmt.Errorf("post-compression audit breakout ATR multiples must be positive")
		}
	}
	for _, mode := range cfg.VolumeModes {
		switch mode {
		case BTC15MPostCompressionVolumeNone, BTC15MPostCompressionVolumeMedian96, BTC15MPostCompressionVolumeP60Prior96:
		default:
			return fmt.Errorf("post-compression audit unknown volume mode %q", mode)
		}
	}
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("post-compression audit horizons must be positive")
		}
	}
	return nil
}

func btcPostCompressionSourceRow(manifest SourceManifest, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) BTC15MPostCompressionSourceRow {
	row := BTC15MPostCompressionSourceRow{
		AuditName:                  FuturesBTC15MPostCompressionDirectionalExpansionAuditName,
		PremiseID:                  BTC15MPostCompressionDirectionalExpansionPremiseID,
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

func btcPostCompression15MFrame() rangeDiscoveryFrameDef {
	for _, frame := range rangeDiscoveryFrameDefs() {
		if frame.timeframe == RangeDiscoveryTimeframe15m {
			return frame
		}
	}
	return rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe15m, interval: 15 * time.Minute, childBars: 3, barsPerDay: 96}
}

func btcPostCompressionCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, sourcePass bool) BTC15MPostCompressionCoverageRow {
	row := BTC15MPostCompressionCoverageRow{
		AuditName:             FuturesBTC15MPostCompressionDirectionalExpansionAuditName,
		PremiseID:             BTC15MPostCompressionDirectionalExpansionPremiseID,
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

func btcPostCompressionParameterCells(cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) []BTC15MPostCompressionParameterCell {
	cfg = cfg.withDefaults()
	lookbacks := append([]int(nil), cfg.CompressionLookbacks...)
	sort.Ints(lookbacks)
	percentiles := append([]float64(nil), cfg.CompressionPercentiles...)
	sort.Float64s(percentiles)
	multiples := append([]float64(nil), cfg.BreakoutATRMultiples...)
	sort.Float64s(multiples)
	modes := btcPostCompressionSortedVolumeModes(cfg.VolumeModes)
	rows := []BTC15MPostCompressionParameterCell{}
	for _, lookback := range lookbacks {
		for _, percentile := range percentiles {
			for _, multiple := range multiples {
				for _, mode := range modes {
					rows = append(rows, BTC15MPostCompressionParameterCell{
						CellID:                  btcPostCompressionCellID(lookback, percentile, multiple, mode),
						LookbackBars:            lookback,
						CompressionPercentile:   percentile,
						PercentileReferenceBars: cfg.PercentileReferenceBars,
						BreakoutATRMultiple:     multiple,
						VolumeMode:              mode,
						VolumeLookbackBars:      cfg.VolumeLookbackBars,
						Predeclared:             true,
					})
				}
			}
		}
	}
	return rows
}

func btcPostCompressionBuildRangeFacts(candles []Candle, lookback int, percentiles []float64, referenceBars int) btcPostCompressionRangeFacts {
	high := nanSlice(len(candles))
	low := nanSlice(len(candles))
	width := nanSlice(len(candles))
	for d := lookback; d < len(candles); d++ {
		rangeHigh := candles[d-lookback].High
		rangeLow := candles[d-lookback].Low
		for i := d - lookback + 1; i <= d-1; i++ {
			if candles[i].High > rangeHigh {
				rangeHigh = candles[i].High
			}
			if candles[i].Low < rangeLow {
				rangeLow = candles[i].Low
			}
		}
		if candles[d-1].Close <= 0 || rangeHigh < rangeLow {
			continue
		}
		high[d] = rangeHigh
		low[d] = rangeLow
		width[d] = (rangeHigh - rangeLow) / candles[d-1].Close
	}
	thresholds := make(map[float64][]float64, len(percentiles))
	for _, percentile := range percentiles {
		thresholds[percentile] = btcPostCompressionRollingPriorValidPercentile(width, referenceBars, percentile)
	}
	return btcPostCompressionRangeFacts{lookback: lookback, high: high, low: low, widthPct: width, thresholds: thresholds}
}

func btcPostCompressionRollingPriorValidPercentile(values []float64, lookback int, percentile float64) []float64 {
	out := nanSlice(len(values))
	if lookback <= 0 || percentile <= 0 || percentile >= 1 {
		return out
	}
	queue := make([]float64, 0, lookback)
	sortedValues := make([]float64, 0, lookback)
	for i, value := range values {
		if len(queue) == lookback {
			out[i] = percentileFromSorted(sortedValues, percentile)
		}
		if validNumber(value) {
			queue = append(queue, value)
			sortedValues = insertSorted(sortedValues, value)
			if len(queue) > lookback {
				old := queue[0]
				queue = queue[1:]
				sortedValues = removeSorted(sortedValues, old)
			}
		}
	}
	return out
}

func btcPostCompressionRollingPriorVolumeThreshold(candles []Candle, lookback int, percentile float64) []float64 {
	values := nanSlice(len(candles))
	for i, candle := range candles {
		values[i] = candle.Volume
	}
	if percentile == 0.50 {
		return btcPostCompressionRollingPriorMedian(values, lookback)
	}
	return btcPostCompressionRollingPriorValidPercentile(values, lookback, percentile)
}

func btcPostCompressionRollingPriorMedian(values []float64, lookback int) []float64 {
	out := nanSlice(len(values))
	if lookback <= 0 {
		return out
	}
	queue := make([]float64, 0, lookback)
	sortedValues := make([]float64, 0, lookback)
	for i, value := range values {
		if len(queue) == lookback {
			out[i] = medianFloat(sortedValues)
		}
		if validNumber(value) {
			queue = append(queue, value)
			sortedValues = insertSorted(sortedValues, value)
			if len(queue) > lookback {
				old := queue[0]
				queue = queue[1:]
				sortedValues = removeSorted(sortedValues, old)
			}
		}
	}
	return out
}

func btcPostCompressionVolumePass(volume float64, mode string, medians, p60 []float64, index int) (float64, bool) {
	switch mode {
	case BTC15MPostCompressionVolumeNone:
		return 0, true
	case BTC15MPostCompressionVolumeMedian96:
		if index < 0 || index >= len(medians) || !validNumber(medians[index]) {
			return math.NaN(), false
		}
		return medians[index], volume > medians[index]
	case BTC15MPostCompressionVolumeP60Prior96:
		if index < 0 || index >= len(p60) || !validNumber(p60[index]) {
			return math.NaN(), false
		}
		return p60[index], volume > p60[index]
	default:
		return math.NaN(), false
	}
}

func btcPostCompressionBaseEventFor(candles []Candle, d int, eventID int, lookback int, percentile, threshold float64, facts btcPostCompressionRangeFacts, multiple, atr float64, volumeMode string, volumeThreshold float64, side string, splits []Split) btcPostCompressionBaseEvent {
	expansion := "upside"
	distance := candles[d].Close - facts.high[d]
	if side == BTC15MPostCompressionSideShort {
		expansion = "downside"
		distance = facts.low[d] - candles[d].Close
	}
	return btcPostCompressionBaseEvent{
		eventID:               eventID,
		cellID:                btcPostCompressionCellID(lookback, percentile, multiple, volumeMode),
		split:                 splitNameForCloseTime(candles[d].CloseTime, splits),
		decisionIndex:         d,
		side:                  side,
		lookback:              lookback,
		compressionPercentile: percentile,
		compressionThreshold:  threshold,
		rangeHigh:             facts.high[d],
		rangeLow:              facts.low[d],
		rangeWidthPct:         facts.widthPct[d],
		breakoutATRMultiple:   multiple,
		priorATR:              atr,
		breakoutDistance:      distance,
		expansionDirection:    expansion,
		volumeMode:            volumeMode,
		volumeThreshold:       volumeThreshold,
	}
}

func btcPostCompressionAddEvent(candles []Candle, event btcPostCompressionBaseEvent, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, result *FuturesBTC15MPostCompressionDirectionalExpansionAuditResult, cells map[string]*BTC15MPostCompressionParameterCell, dedup map[btcPostCompressionDedupKey]*btcPostCompressionDedupAccumulator, candidateAcc map[btcPostCompressionSummaryKey]*btcPostCompressionMetricAccumulator, missing map[string]int, candidateRowID int) int {
	cell := cells[event.cellID]
	if cell != nil {
		cell.CandidateRows++
		if event.side == BTC15MPostCompressionSideLong {
			cell.LongCandidateRows++
		} else {
			cell.ShortCandidateRows++
		}
	}
	dedupKey := btcPostCompressionDedupKey{decisionIndex: event.decisionIndex, side: event.side}
	if dedup[dedupKey] == nil {
		dedup[dedupKey] = &btcPostCompressionDedupAccumulator{
			split:             event.split,
			decisionCloseTime: candles[event.decisionIndex].CloseTime.UTC().Format(timeLayout),
			cells:             map[string]struct{}{},
		}
	}
	dedup[dedupKey].cells[event.cellID] = struct{}{}

	for _, horizon := range cfg.HorizonsBars {
		label, ok := btcPostCompressionForwardLabelFor(candles, event.decisionIndex, event.side, horizon)
		if !ok {
			missing["missing_forward_label"]++
			continue
		}
		candidateRowID++
		row := btcPostCompressionCandidateRow(candles, event, label, candidateRowID, cfg)
		result.CandidateRows = append(result.CandidateRows, row)
		if cell != nil {
			cell.LabelRows++
		}
		for _, splitName := range btcPostCompressionSplitNames(event.split) {
			key := btcPostCompressionSummaryKey{cellID: event.cellID, split: splitName, side: event.side, horizonBars: horizon}
			btcPostCompressionEnsureCandidateAcc(candidateAcc, key).add(label)
		}
	}
	return candidateRowID
}

func btcPostCompressionCandidateRow(candles []Candle, event btcPostCompressionBaseEvent, label btcPostCompressionForwardLabel, rowID int, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) BTC15MPostCompressionCandidateRow {
	decision := candles[event.decisionIndex]
	return BTC15MPostCompressionCandidateRow{
		CandidateRowID:                 rowID,
		BaseEventID:                    event.eventID,
		CellID:                         event.cellID,
		Split:                          event.split,
		Timeframe:                      RangeDiscoveryTimeframe15m,
		DecisionIndex:                  event.decisionIndex,
		DecisionOpenTime:               decision.OpenTime.UTC().Format(timeLayout),
		DecisionCloseTime:              decision.CloseTime.UTC().Format(timeLayout),
		Side:                           event.side,
		TimingLabel:                    "next_15m_open",
		LookbackBars:                   event.lookback,
		CompressionPercentile:          event.compressionPercentile,
		CompressionThreshold:           event.compressionThreshold,
		RangeHigh:                      event.rangeHigh,
		RangeLow:                       event.rangeLow,
		RangeWidthPct:                  event.rangeWidthPct,
		BreakoutATRMultiple:            event.breakoutATRMultiple,
		PriorATR14:                     event.priorATR,
		BreakoutDistance:               event.breakoutDistance,
		ExpansionDirection:             event.expansionDirection,
		DecisionClose:                  decision.Close,
		DecisionVolume:                 decision.Volume,
		VolumeMode:                     event.volumeMode,
		VolumeThreshold:                event.volumeThreshold,
		VolumeLookbackBars:             cfg.VolumeLookbackBars,
		HorizonBars:                    label.horizonBars,
		LabelAnchorOpenTime:            label.anchorOpenTime,
		LabelAnchorOpen:                label.anchorOpen,
		LabelWindowEndIndex:            label.windowEndIndex,
		LabelWindowEndTime:             label.windowEndTime,
		IntendedSideForwardCloseReturn: label.forwardCloseReturnBP,
		IntendedSideFavorable:          label.favorableBP,
		IntendedSideAdverse:            label.adverseBP,
		FavorableGreaterThanAdverse:    label.favorableGreaterAdverse,
		FavorableMinusAdverse:          label.favorableMinusAdverseBP,
		ForwardLabelMetadataOnly:       true,
		ForwardLabelUsedAsFeature:      false,
		UsesFutureRowsForFeatures:      false,
		DerivativesVetoUsed:            false,
	}
}

func btcPostCompressionForwardLabelFor(candles []Candle, decisionIndex int, side string, horizon int) (btcPostCompressionForwardLabel, bool) {
	if decisionIndex < 0 || horizon <= 0 || decisionIndex+1 >= len(candles) || decisionIndex+horizon >= len(candles) {
		return btcPostCompressionForwardLabel{}, false
	}
	anchor := candles[decisionIndex+1]
	if anchor.Open <= 0 {
		return btcPostCompressionForwardLabel{}, false
	}
	end := decisionIndex + horizon
	maxHigh := candles[decisionIndex+1].High
	minLow := candles[decisionIndex+1].Low
	for i := decisionIndex + 1; i <= end; i++ {
		if candles[i].High > maxHigh {
			maxHigh = candles[i].High
		}
		if candles[i].Low < minLow {
			minLow = candles[i].Low
		}
	}
	endClose := candles[end].Close
	forward := 0.0
	favorable := 0.0
	adverse := 0.0
	switch side {
	case BTC15MPostCompressionSideLong:
		forward = (endClose - anchor.Open) / anchor.Open * 10000
		favorable = (maxHigh - anchor.Open) / anchor.Open * 10000
		adverse = (anchor.Open - minLow) / anchor.Open * 10000
	case BTC15MPostCompressionSideShort:
		forward = (anchor.Open - endClose) / anchor.Open * 10000
		favorable = (anchor.Open - minLow) / anchor.Open * 10000
		adverse = (maxHigh - anchor.Open) / anchor.Open * 10000
	default:
		return btcPostCompressionForwardLabel{}, false
	}
	if favorable < 0 {
		favorable = 0
	}
	if adverse < 0 {
		adverse = 0
	}
	return btcPostCompressionForwardLabel{
		horizonBars:             horizon,
		anchorOpenTime:          anchor.OpenTime.UTC().Format(timeLayout),
		anchorOpen:              anchor.Open,
		windowEndIndex:          end,
		windowEndTime:           candles[end].CloseTime.UTC().Format(timeLayout),
		forwardCloseReturnBP:    forward,
		favorableBP:             favorable,
		adverseBP:               adverse,
		favorableGreaterAdverse: favorable > adverse,
		favorableMinusAdverseBP: favorable - adverse,
	}, true
}

func btcPostCompressionDedupRows(dedup map[btcPostCompressionDedupKey]*btcPostCompressionDedupAccumulator) []BTC15MPostCompressionDedupEventRow {
	keys := make([]btcPostCompressionDedupKey, 0, len(dedup))
	for key := range dedup {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].decisionIndex != keys[j].decisionIndex {
			return keys[i].decisionIndex < keys[j].decisionIndex
		}
		return btcPostCompressionSideSortKey(keys[i].side) < btcPostCompressionSideSortKey(keys[j].side)
	})
	rows := make([]BTC15MPostCompressionDedupEventRow, 0, len(keys))
	for i, key := range keys {
		acc := dedup[key]
		cellIDs := make([]string, 0, len(acc.cells))
		for cellID := range acc.cells {
			cellIDs = append(cellIDs, cellID)
		}
		sort.Strings(cellIDs)
		rows = append(rows, BTC15MPostCompressionDedupEventRow{
			DedupEventID:      i + 1,
			Split:             acc.split,
			Timeframe:         RangeDiscoveryTimeframe15m,
			DecisionIndex:     key.decisionIndex,
			DecisionCloseTime: acc.decisionCloseTime,
			Side:              key.side,
			MatchedCellCount:  len(cellIDs),
			MatchedCellIDs:    strings.Join(cellIDs, ";"),
		})
	}
	return rows
}

func btcPostCompressionBaselineRows(accs map[btcPostCompressionBaselineKey]*btcPostCompressionMetricAccumulator, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, splits []Split) []BTC15MPostCompressionBaselineRow {
	keys := make([]btcPostCompressionBaselineKey, 0, len(accs))
	for key := range accs {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if splitSortKey(keys[i].split) != splitSortKey(keys[j].split) {
			return splitSortKey(keys[i].split) < splitSortKey(keys[j].split)
		}
		if btcPostCompressionSideSortKey(keys[i].side) != btcPostCompressionSideSortKey(keys[j].side) {
			return btcPostCompressionSideSortKey(keys[i].side) < btcPostCompressionSideSortKey(keys[j].side)
		}
		return keys[i].horizonBars < keys[j].horizonBars
	})
	rows := []BTC15MPostCompressionBaselineRow{}
	for _, key := range keys {
		acc := accs[key]
		rows = append(rows, BTC15MPostCompressionBaselineRow{
			Split:                              key.split,
			Side:                               key.side,
			HorizonBars:                        key.horizonBars,
			EligibleRows:                       acc.count,
			MeanForwardCloseReturnBP:           acc.meanForward(),
			MeanFavorableBP:                    acc.meanFavorable(),
			MeanAdverseBP:                      acc.meanAdverse(),
			MeanFavorableMinusAdverseBP:        acc.meanFavMinusAdv(),
			FavorableGreaterThanAdverseRate:    acc.favorableGreaterRate(),
			AuditWarmupBars:                    btcPostCompressionAuditWarmup(cfg),
			ForwardLabelsAsBaselineInput:       false,
			UnconditionalEligible15MBaseline:   true,
			SharesCandidateWarmupAndLabelRules: true,
		})
	}
	return rows
}

func btcPostCompressionSplitSummaryRows(cells []BTC15MPostCompressionParameterCell, candidateAcc map[btcPostCompressionSummaryKey]*btcPostCompressionMetricAccumulator, baselineAcc map[btcPostCompressionBaselineKey]*btcPostCompressionMetricAccumulator, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig, splits []Split) ([]BTC15MPostCompressionSplitSummaryRow, map[btcPostCompressionPassingKey]bool) {
	splitNames := btcPostCompressionAllSplitNames(splits)
	rows := []BTC15MPostCompressionSplitSummaryRow{}
	rowByKey := map[btcPostCompressionSummaryKey]BTC15MPostCompressionSplitSummaryRow{}
	for _, cell := range cells {
		for _, side := range []string{BTC15MPostCompressionSideLong, BTC15MPostCompressionSideShort} {
			for _, horizon := range cfg.HorizonsBars {
				for _, split := range splitNames {
					key := btcPostCompressionSummaryKey{cellID: cell.CellID, split: split, side: side, horizonBars: horizon}
					cand := candidateAcc[key]
					base := baselineAcc[btcPostCompressionBaselineKey{split: split, side: side, horizonBars: horizon}]
					row := BTC15MPostCompressionSplitSummaryRow{
						CellID:                cell.CellID,
						Split:                 split,
						Side:                  side,
						HorizonBars:           horizon,
						LookbackBars:          cell.LookbackBars,
						CompressionPercentile: cell.CompressionPercentile,
						BreakoutATRMultiple:   cell.BreakoutATRMultiple,
						VolumeMode:            cell.VolumeMode,
					}
					if cand != nil {
						row.CandidateRows = cand.count
						row.MeanForwardCloseReturnBP = cand.meanForward()
						row.MeanFavorableMinusAdverseBP = cand.meanFavMinusAdv()
						row.FavorableGreaterThanAdverse = cand.favorableGreaterRate()
					}
					if base != nil {
						row.BaselineRows = base.count
						row.BaselineForwardCloseReturnBP = base.meanForward()
						row.BaselineFavorableMinusAdverse = base.meanFavMinusAdv()
						row.BaselineFavorableGTAdverseRate = base.favorableGreaterRate()
					}
					row.DeltaForwardCloseReturnBP = row.MeanForwardCloseReturnBP - row.BaselineForwardCloseReturnBP
					row.DeltaFavorableMinusAdverseBP = row.MeanFavorableMinusAdverseBP - row.BaselineFavorableMinusAdverse
					row.DeltaFavorableGTAdverseRate = row.FavorableGreaterThanAdverse - row.BaselineFavorableGTAdverseRate
					row.ForwardReturnGatePass = row.CandidateRows > 0 && row.BaselineRows > 0 && row.MeanForwardCloseReturnBP > row.BaselineForwardCloseReturnBP
					row.FavorableMinusAdverseGatePass = row.CandidateRows > 0 && row.BaselineRows > 0 && row.MeanFavorableMinusAdverseBP > row.BaselineFavorableMinusAdverse
					row.FavorableGTAdverseGatePass = row.CandidateRows > 0 && row.BaselineRows > 0 && row.FavorableGreaterThanAdverse > row.BaselineFavorableGTAdverseRate
					row.PassesBaseline = row.ForwardReturnGatePass && row.FavorableMinusAdverseGatePass && row.FavorableGTAdverseGatePass
					rowByKey[key] = row
					rows = append(rows, row)
				}
			}
		}
	}

	passing := map[btcPostCompressionPassingKey]bool{}
	for _, cell := range cells {
		for _, side := range []string{BTC15MPostCompressionSideLong, BTC15MPostCompressionSideShort} {
			for _, horizon := range cfg.HorizonsBars {
				pass := true
				for _, split := range splitNames {
					key := btcPostCompressionSummaryKey{cellID: cell.CellID, split: split, side: side, horizonBars: horizon}
					if !rowByKey[key].PassesBaseline {
						pass = false
						break
					}
				}
				if pass {
					passing[btcPostCompressionPassingKey{cellID: cell.CellID, side: side, horizonBars: horizon}] = true
				}
			}
		}
	}
	for i := range rows {
		passKey := btcPostCompressionPassingKey{cellID: rows[i].CellID, side: rows[i].Side, horizonBars: rows[i].HorizonBars}
		rows[i].FullAndEveryPrimarySplitPass = passing[passKey]
	}
	sort.Slice(rows, func(i, j int) bool {
		return btcPostCompressionLessSplitSummary(rows[i], rows[j])
	})
	return rows, passing
}

func btcPostCompressionAdjacencyRows(cells []BTC15MPostCompressionParameterCell, passing map[btcPostCompressionPassingKey]bool, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) []BTC15MPostCompressionAdjacencyRow {
	rows := []BTC15MPostCompressionAdjacencyRow{}
	for _, cell := range cells {
		for _, side := range []string{BTC15MPostCompressionSideLong, BTC15MPostCompressionSideShort} {
			for _, horizon := range cfg.HorizonsBars {
				key := btcPostCompressionPassingKey{cellID: cell.CellID, side: side, horizonBars: horizon}
				adjacent := []string{}
				if passing[key] {
					for _, other := range cells {
						otherKey := btcPostCompressionPassingKey{cellID: other.CellID, side: side, horizonBars: horizon}
						if other.CellID == cell.CellID || !passing[otherKey] {
							continue
						}
						if btcPostCompressionAdjacentCells(cell, other, cfg) {
							adjacent = append(adjacent, other.CellID)
						}
					}
				}
				sort.Strings(adjacent)
				rows = append(rows, BTC15MPostCompressionAdjacencyRow{
					CellID:                        cell.CellID,
					Side:                          side,
					HorizonBars:                   horizon,
					FullAndEveryPrimarySplitPass:  passing[key],
					AdjacentPassingCellIDs:        strings.Join(adjacent, ";"),
					AdjacentPassingCellCount:      len(adjacent),
					PassesAdjacentClusterGate:     passing[key] && len(adjacent) > 0,
					AdjacencyDefinition:           "one ordered step different in exactly one parameter dimension; all other dimensions equal",
					OneIsolatedPassingCellFailure: passing[key] && len(adjacent) == 0,
				})
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CellID != rows[j].CellID {
			return rows[i].CellID < rows[j].CellID
		}
		if btcPostCompressionSideSortKey(rows[i].Side) != btcPostCompressionSideSortKey(rows[j].Side) {
			return btcPostCompressionSideSortKey(rows[i].Side) < btcPostCompressionSideSortKey(rows[j].Side)
		}
		return rows[i].HorizonBars < rows[j].HorizonBars
	})
	return rows
}

func btcPostCompressionAdjacentCells(a, b BTC15MPostCompressionParameterCell, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) bool {
	diff := 0
	if btcPostCompressionOrderedStepInt(a.LookbackBars, b.LookbackBars, cfg.CompressionLookbacks) {
		diff++
	} else if a.LookbackBars != b.LookbackBars {
		return false
	}
	if btcPostCompressionOrderedStepFloat(a.CompressionPercentile, b.CompressionPercentile, cfg.CompressionPercentiles) {
		diff++
	} else if a.CompressionPercentile != b.CompressionPercentile {
		return false
	}
	if btcPostCompressionOrderedStepFloat(a.BreakoutATRMultiple, b.BreakoutATRMultiple, cfg.BreakoutATRMultiples) {
		diff++
	} else if a.BreakoutATRMultiple != b.BreakoutATRMultiple {
		return false
	}
	if btcPostCompressionOrderedStepString(a.VolumeMode, b.VolumeMode, btcPostCompressionSortedVolumeModes(cfg.VolumeModes)) {
		diff++
	} else if a.VolumeMode != b.VolumeMode {
		return false
	}
	return diff == 1
}

func btcPostCompressionMissingnessRows(missing map[string]int, totalRows int) []BTC15MPostCompressionMissingnessRow {
	reasons := make([]string, 0, len(missing))
	for reason := range missing {
		reasons = append(reasons, reason)
	}
	sort.Strings(reasons)
	rows := []BTC15MPostCompressionMissingnessRow{}
	for _, reason := range reasons {
		count := missing[reason]
		rate := 0.0
		if totalRows > 0 {
			rate = float64(count) / float64(totalRows)
		}
		rows = append(rows, BTC15MPostCompressionMissingnessRow{
			AuditName:         FuturesBTC15MPostCompressionDirectionalExpansionAuditName,
			PremiseID:         BTC15MPostCompressionDirectionalExpansionPremiseID,
			Scope:             RangeDiscoveryTimeframe15m,
			Reason:            reason,
			Count:             count,
			TotalRows:         totalRows,
			Rate:              rate,
			MissingDataPolicy: "skip rows; no fill/interpolation/nearest-future",
			ForwardFilledRows: 0,
		})
	}
	return rows
}

func btcPostCompressionFalsification(report BTC15MPostCompressionFalsification, cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) BTC15MPostCompressionFalsification {
	report.AuditName = FuturesBTC15MPostCompressionDirectionalExpansionAuditName
	report.PremiseID = BTC15MPostCompressionDirectionalExpansionPremiseID
	failures := []string{}
	if !report.SourceResamplePass {
		failures = append(failures, "source_or_resample_failed")
	}
	if !report.LeakagePass {
		failures = append(failures, "future_label_leak")
	}
	if !report.CandidateSizePass {
		failures = append(failures, "full_deduplicated_candidate_count_below_minimum")
	}
	if !report.SplitSizePass {
		failures = append(failures, "primary_split_deduplicated_candidate_count_below_minimum")
	}
	if !report.BaselineSeparationPass {
		failures = append(failures, "no_cell_cluster_separates_beyond_baseline")
	}
	if !report.AdjacentCellClusterPass {
		failures = append(failures, "only_isolated_or_no_passing_parameter_cells")
	}
	if !report.SplitStabilityPass {
		failures = append(failures, "full_sample_separation_not_supported_in_every_primary_split")
	}
	if !report.ClosedFamilyProtectionPass {
		failures = append(failures, "closed_family_reslice")
	}
	if !report.DerivativesVetoContaminationPass {
		failures = append(failures, "derivatives_veto_contamination")
	}
	report.FailureReasons = failures
	switch {
	case !report.ClosedFamilyProtectionPass:
		report.StopState = BTC15MPostCompressionDirectionalExpansionStopStateClosedReslice
	case report.SourceResamplePass &&
		report.LeakagePass &&
		report.CandidateSizePass &&
		report.SplitSizePass &&
		report.BaselineSeparationPass &&
		report.AdjacentCellClusterPass &&
		report.SplitStabilityPass &&
		report.ClosedFamilyProtectionPass &&
		report.DerivativesVetoContaminationPass:
		report.StopState = BTC15MPostCompressionDirectionalExpansionStopStatePassedNeedsReview
	default:
		report.StopState = BTC15MPostCompressionDirectionalExpansionStopStateFailedNoPremise
	}
	if report.RequiredFullDeduplicatedCandidateRows == 0 {
		report.RequiredFullDeduplicatedCandidateRows = cfg.MinFullDedupCandidates
	}
	if report.RequiredSplitDeduplicatedCandidateRows == 0 {
		report.RequiredSplitDeduplicatedCandidateRows = cfg.MinSplitDedupCandidates
	}
	return report
}

func (acc *btcPostCompressionMetricAccumulator) add(label btcPostCompressionForwardLabel) {
	acc.count++
	acc.sumForward += label.forwardCloseReturnBP
	acc.sumFavorable += label.favorableBP
	acc.sumAdverse += label.adverseBP
	acc.sumFavMinusAdv += label.favorableMinusAdverseBP
	if label.favorableGreaterAdverse {
		acc.favGreaterCount++
	}
}

func (acc btcPostCompressionMetricAccumulator) meanForward() float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sumForward / float64(acc.count)
}

func (acc btcPostCompressionMetricAccumulator) meanFavorable() float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sumFavorable / float64(acc.count)
}

func (acc btcPostCompressionMetricAccumulator) meanAdverse() float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sumAdverse / float64(acc.count)
}

func (acc btcPostCompressionMetricAccumulator) meanFavMinusAdv() float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sumFavMinusAdv / float64(acc.count)
}

func (acc btcPostCompressionMetricAccumulator) favorableGreaterRate() float64 {
	if acc.count == 0 {
		return 0
	}
	return float64(acc.favGreaterCount) / float64(acc.count)
}

func btcPostCompressionEnsureCandidateAcc(accs map[btcPostCompressionSummaryKey]*btcPostCompressionMetricAccumulator, key btcPostCompressionSummaryKey) *btcPostCompressionMetricAccumulator {
	if accs[key] == nil {
		accs[key] = &btcPostCompressionMetricAccumulator{}
	}
	return accs[key]
}

func btcPostCompressionEnsureBaselineAcc(accs map[btcPostCompressionBaselineKey]*btcPostCompressionMetricAccumulator, key btcPostCompressionBaselineKey) *btcPostCompressionMetricAccumulator {
	if accs[key] == nil {
		accs[key] = &btcPostCompressionMetricAccumulator{}
	}
	return accs[key]
}

func btcPostCompressionDedupCounts(rows []BTC15MPostCompressionDedupEventRow, splits []Split) (int, int) {
	full := len(rows)
	primaryCounts := map[string]int{}
	for _, split := range btcPostCompressionPrimarySplitNames(splits) {
		primaryCounts[split] = 0
	}
	for _, row := range rows {
		if _, ok := primaryCounts[row.Split]; ok {
			primaryCounts[row.Split]++
		}
	}
	if len(primaryCounts) == 0 {
		return full, full
	}
	minSplit := int(^uint(0) >> 1)
	for _, count := range primaryCounts {
		if count < minSplit {
			minSplit = count
		}
	}
	if minSplit == int(^uint(0)>>1) {
		minSplit = 0
	}
	return full, minSplit
}

func btcPostCompressionSourceResamplePass(sources []BTC15MPostCompressionSourceRow, coverage []BTC15MPostCompressionCoverageRow) bool {
	if len(sources) == 0 || len(coverage) == 0 {
		return false
	}
	for _, row := range sources {
		if !row.SourceFactsPass || row.ValidationStatus != "accepted" || row.DerivativesVetoAsInput || row.ForwardLabelsAsSourceInput {
			return false
		}
	}
	for _, row := range coverage {
		if !row.SourceResamplePass || row.ValidationStatus != "accepted" || !row.Complete {
			return false
		}
	}
	return true
}

func btcPostCompressionAuditWarmup(cfg FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig) int {
	return btcPostCompressionMaxInt(cfg.CompressionLookbacks) + cfg.PercentileReferenceBars
}

func btcPostCompressionMaxInt(values []int) int {
	maxValue := 0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func btcPostCompressionCellID(lookback int, percentile float64, multiple float64, volumeMode string) string {
	return fmt.Sprintf("l%d_q%02d_m%03d_%s", lookback, int(math.Round(percentile*100)), int(math.Round(multiple*100)), volumeMode)
}

func btcPostCompressionSortedVolumeModes(modes []string) []string {
	order := map[string]int{
		BTC15MPostCompressionVolumeNone:       0,
		BTC15MPostCompressionVolumeMedian96:   1,
		BTC15MPostCompressionVolumeP60Prior96: 2,
	}
	out := append([]string(nil), modes...)
	sort.Slice(out, func(i, j int) bool {
		if order[out[i]] != order[out[j]] {
			return order[out[i]] < order[out[j]]
		}
		return out[i] < out[j]
	})
	return out
}

func btcPostCompressionSideSortKey(side string) int {
	switch side {
	case BTC15MPostCompressionSideLong:
		return 0
	case BTC15MPostCompressionSideShort:
		return 1
	default:
		return 99
	}
}

func btcPostCompressionSplitNames(split string) []string {
	if split == "" || split == fullSplitName {
		return []string{fullSplitName}
	}
	return []string{split, fullSplitName}
}

func btcPostCompressionAllSplitNames(splits []Split) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, split := range splits {
		if seen[split.Name] {
			continue
		}
		seen[split.Name] = true
		out = append(out, split.Name)
	}
	sort.Slice(out, func(i, j int) bool {
		return splitSortKey(out[i]) < splitSortKey(out[j])
	})
	return out
}

func btcPostCompressionPrimarySplitNames(splits []Split) []string {
	names := []string{}
	for _, split := range splits {
		if split.Name != fullSplitName {
			names = append(names, split.Name)
		}
	}
	return names
}

func btcPostCompressionOrderedStepInt(a, b int, order []int) bool {
	ordered := append([]int(nil), order...)
	sort.Ints(ordered)
	ai, bi := -1, -1
	for i, value := range ordered {
		if value == a {
			ai = i
		}
		if value == b {
			bi = i
		}
	}
	return ai >= 0 && bi >= 0 && int(math.Abs(float64(ai-bi))) == 1
}

func btcPostCompressionOrderedStepFloat(a, b float64, order []float64) bool {
	ordered := append([]float64(nil), order...)
	sort.Float64s(ordered)
	ai, bi := -1, -1
	for i, value := range ordered {
		if value == a {
			ai = i
		}
		if value == b {
			bi = i
		}
	}
	return ai >= 0 && bi >= 0 && int(math.Abs(float64(ai-bi))) == 1
}

func btcPostCompressionOrderedStepString(a, b string, order []string) bool {
	ai, bi := -1, -1
	for i, value := range order {
		if value == a {
			ai = i
		}
		if value == b {
			bi = i
		}
	}
	return ai >= 0 && bi >= 0 && int(math.Abs(float64(ai-bi))) == 1
}

func btcPostCompressionLessSplitSummary(a, b BTC15MPostCompressionSplitSummaryRow) bool {
	if a.CellID != b.CellID {
		return a.CellID < b.CellID
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if btcPostCompressionSideSortKey(a.Side) != btcPostCompressionSideSortKey(b.Side) {
		return btcPostCompressionSideSortKey(a.Side) < btcPostCompressionSideSortKey(b.Side)
	}
	return a.HorizonBars < b.HorizonBars
}
