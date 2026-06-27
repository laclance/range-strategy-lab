package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	FuturesRangeContextTriageAuditName = "futures_range_context_triage_audit"

	RangeContextTriageStopStateSourceGap                = "range_context_triage_source_gap"
	RangeContextTriageStopStateNoRangeEpisodes          = "range_context_triage_no_range_episodes"
	RangeContextTriageStopStateNoUsableCohorts          = "range_context_triage_no_usable_cohorts"
	RangeContextTriageStopStateFailedNoStrategyPremise  = "range_context_triage_failed_no_strategy_premise"
	RangeContextTriageStopStateReadyForStrategySpec     = "range_context_triage_ready_for_strategy_spec"
	RangeContextTriageStopStateRejectedClosedReslice    = "range_context_triage_rejected_closed_family_reslice"
	RangeContextTriageQualityTooNarrowNoise             = "too_narrow_noise"
	RangeContextTriageQualityNarrowOrderly              = "narrow_orderly"
	RangeContextTriageQualityBalancedOrderly            = "balanced_orderly"
	RangeContextTriageQualityWideVolatile               = "wide_volatile"
	RangeContextTriageQualityChoppy                     = "choppy"
	RangeContextTriageQualityUnknown                    = "unknown"
	RangeContextTriageLabelContainedRotation            = "contained_rotation"
	RangeContextTriageLabelCleanExpansionUp             = "clean_expansion_up"
	RangeContextTriageLabelCleanExpansionDown           = "clean_expansion_down"
	RangeContextTriageLabelFalseBreakReentryUp          = "false_break_reentry_up"
	RangeContextTriageLabelFalseBreakReentryDown        = "false_break_reentry_down"
	RangeContextTriageLabelBoundaryChop                 = "boundary_chop"
	RangeContextTriageLabelDriftThroughUp               = "drift_through_up"
	RangeContextTriageLabelDriftThroughDown             = "drift_through_down"
	RangeContextTriageLabelLowWidthNoise                = "low_width_noise"
	RangeContextTriageLabelNoResolution                 = "no_resolution"
	RangeContextTriageLabelMissingFuture                = "missing_future"
	RangeContextTriageSessionAsia                       = "asia_utc_00_07"
	RangeContextTriageSessionEurope                     = "europe_utc_08_12"
	RangeContextTriageSessionUSOverlap                  = "us_overlap_utc_13_16"
	RangeContextTriageSessionUSLate                     = "us_late_utc_17_23"
	RangeContextTriageCohortAll                         = "all"
	RangeContextTriageCohortQualityBucket               = "quality_bucket"
	RangeContextTriageCohortMatureSession               = "mature_session"
	RangeContextTriageCohortPrimaryContextLabel         = "primary_context_label"
	RangeContextTriageCohortQualityBucketMatureSession  = "quality_bucket_mature_session"
	rangeContextTriageSourcePath                        = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"
	rangeContextTriageExpectedRows                      = 573984
	rangeContextTriageExpectedFirst                     = "2021-01-01T00:00:00Z"
	rangeContextTriageExpectedLast                      = "2026-06-16T23:55:00Z"
	rangeContextTriageExpectedZeroVol                   = 66
	rangeContextTriageExpected15MRows                   = 191328
	rangeContextTriageExpected15MLast                   = "2026-06-16T23:45:00Z"
	rangeContextTriageExpected1HRows                    = 47832
	rangeContextTriageExpected1HLast                    = "2026-06-16T23:00:00Z"
	rangeContextTriageExpected4HRows                    = 11958
	rangeContextTriageExpected4HLast                    = "2026-06-16T20:00:00Z"
	rangeContextTriageDetectorProfileID                 = "p30_c12_bollinger_on_adx_off"
	rangeContextTriageCohortRankableDecisionContextNote = "outcome_label_cohort_not_decision_context"
)

type FuturesRangeContextTriageAuditConfig struct {
	SourcePath                   string
	ApprovedSourcePath           string
	ExpectedSourceRows           int
	ExpectedFirstOpenTime        string
	ExpectedLastOpenTime         string
	ExpectedGapCount             int
	ExpectedDuplicateCount       int
	ExpectedZeroVolumeCount      int
	SkipSourceFactCheck          bool
	Timeframes                   []string
	Expected15MRows              int
	Expected15MLastOpenTime      string
	Expected1HRows               int
	Expected1HLastOpenTime       string
	Expected4HRows               int
	Expected4HLastOpenTime       string
	SkipCoverageCountCheck       bool
	DetectorLookbackDays         int
	DetectorLookbackBarsOverride int
	DetectorPercentile           float64
	DetectorMinConsecutiveBars   int
	HorizonsBars                 []int
	QuickFailureHorizonBars      int
	ReentryWindowBars            int
	CleanExpansionThreshold      float64
	DriftThreshold               float64
	BoundaryChopTransitions      int
	LowWidthPctThreshold         float64
	LowWidthATRThreshold         float64
	WideWidthPctThreshold        float64
	WideWidthATRThreshold        float64
	ChoppyMidCrossThreshold      int
	ChoppyBoundaryTouchThreshold int
	MinFullCohortCount           int
	MinSplitCohortCount          int
	MinSessionSplitCohortCount   int
	MinUsableContextRateFull     float64
	MinUsableContextRateSplit    float64
	MaxToxicContextRateFull      float64
	MaxToxicContextRateSplit     float64
	MaxMissingFutureRateFull     float64
	MaxDominantToxicRate         float64
}

type FuturesRangeContextTriageAuditResult struct {
	SourceRows      []FuturesRangeContextTriageSourceRow      `json:"source_rows"`
	CoverageRows    []FuturesRangeContextTriageCoverageRow    `json:"coverage_rows"`
	EpisodeRows     []FuturesRangeContextTriageEpisodeRow     `json:"episode_rows"`
	QualityRows     []FuturesRangeContextTriageQualityRow     `json:"quality_rows"`
	SessionRows     []FuturesRangeContextTriageSessionRow     `json:"session_rows"`
	FailureModeRows []FuturesRangeContextTriageFailureModeRow `json:"failure_mode_rows"`
	CohortRows      []FuturesRangeContextTriageCohortRow      `json:"cohort_rows"`
	RankingRows     []FuturesRangeContextTriageRankingRow     `json:"ranking_rows"`
	SummaryRows     []FuturesRangeContextTriageSummaryRow     `json:"summary_rows"`
	PassingCohorts  int                                       `json:"passing_cohorts"`
	StopState       string                                    `json:"stop_state"`
}

type FuturesRangeContextTriageSourceRow struct {
	Path                    string `json:"path"`
	ApprovedPath            string `json:"approved_path"`
	Venue                   string `json:"venue"`
	Product                 string `json:"product"`
	Symbol                  string `json:"symbol"`
	Interval                string `json:"interval"`
	RowCount                int    `json:"row_count"`
	ExpectedRowCount        int    `json:"expected_row_count"`
	FirstOpenTime           string `json:"first_open_time"`
	ExpectedFirstOpenTime   string `json:"expected_first_open_time"`
	LastOpenTime            string `json:"last_open_time"`
	ExpectedLastOpenTime    string `json:"expected_last_open_time"`
	GapCount                int    `json:"gap_count"`
	ExpectedGapCount        int    `json:"expected_gap_count"`
	DuplicateCount          int    `json:"duplicate_count"`
	ExpectedDuplicateCount  int    `json:"expected_duplicate_count"`
	ZeroVolumeCount         int    `json:"zero_volume_count"`
	ExpectedZeroVolumeCount int    `json:"expected_zero_volume_count"`
	ComparisonOnly          bool   `json:"comparison_only"`
	SourceFactsPass         bool   `json:"source_facts_pass"`
	ValidationStatus        string `json:"validation_status"`
	ValidationError         string `json:"validation_error,omitempty"`
}

type FuturesRangeContextTriageCoverageRow struct {
	ExpectedRowCount      int    `json:"expected_row_count"`
	ExpectedFirstOpenTime string `json:"expected_first_open_time"`
	ExpectedLastOpenTime  string `json:"expected_last_open_time"`
	CoverageFactsPass     bool   `json:"coverage_facts_pass"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesRangeContextTriageEpisodeRow struct {
	EpisodeID                   int     `json:"episode_id"`
	Timeframe                   string  `json:"timeframe"`
	Split                       string  `json:"split"`
	DetectorProfileID           string  `json:"detector_profile_id"`
	StartIndex                  int     `json:"start_index"`
	MatureIndex                 int     `json:"mature_index"`
	RawEndIndex                 int     `json:"raw_end_index"`
	StartCloseTime              string  `json:"start_close_time"`
	MatureCloseTime             string  `json:"mature_close_time"`
	RawEndCloseTime             string  `json:"raw_end_close_time"`
	High                        float64 `json:"high"`
	Low                         float64 `json:"low"`
	Mid                         float64 `json:"mid"`
	UpperQuartile               float64 `json:"upper_quartile"`
	LowerQuartile               float64 `json:"lower_quartile"`
	Width                       float64 `json:"width"`
	WidthPct                    float64 `json:"width_pct"`
	WidthToATRRatio             float64 `json:"width_to_atr_ratio"`
	DurationBarsToMaturity      int     `json:"duration_bars_to_maturity"`
	ActiveBarsToMaturity        int     `json:"active_bars_to_maturity"`
	MatureClosePosition         float64 `json:"mature_close_position"`
	PreMatureMidCrossCount      int     `json:"pre_mature_mid_cross_count"`
	PreMatureBoundaryTouchCount int     `json:"pre_mature_boundary_touch_count"`
	PreMatureCloseInsideRate    float64 `json:"pre_mature_close_inside_rate"`
	PreMatureWickOvershootCount int     `json:"pre_mature_wick_overshoot_count"`
	QualityBucket               string  `json:"quality_bucket"`
	MatureSession               string  `json:"mature_session"`
	Eligible                    bool    `json:"eligible"`
	SkippedReason               string  `json:"skipped_reason,omitempty"`
}

type FuturesRangeContextTriageQualityRow struct {
	EpisodeID                   int     `json:"episode_id"`
	Timeframe                   string  `json:"timeframe"`
	Split                       string  `json:"split"`
	MatureCloseTime             string  `json:"mature_close_time"`
	WidthPct                    float64 `json:"width_pct"`
	WidthToATRRatio             float64 `json:"width_to_atr_ratio"`
	DurationBarsToMaturity      int     `json:"duration_bars_to_maturity"`
	ActiveBarsToMaturity        int     `json:"active_bars_to_maturity"`
	MatureClosePosition         float64 `json:"mature_close_position"`
	PreMatureMidCrossCount      int     `json:"pre_mature_mid_cross_count"`
	PreMatureBoundaryTouchCount int     `json:"pre_mature_boundary_touch_count"`
	PreMatureCloseInsideRate    float64 `json:"pre_mature_close_inside_rate"`
	PreMatureWickOvershootCount int     `json:"pre_mature_wick_overshoot_count"`
	QualityBucket               string  `json:"quality_bucket"`
	Eligible                    bool    `json:"eligible"`
	SkippedReason               string  `json:"skipped_reason,omitempty"`
}

type FuturesRangeContextTriageSessionRow struct {
	Split                   string  `json:"split"`
	Timeframe               string  `json:"timeframe"`
	HorizonBars             int     `json:"horizon_bars"`
	MatureSession           string  `json:"mature_session"`
	FirstOutsideSession     string  `json:"first_outside_session,omitempty"`
	CandidateCount          int     `json:"candidate_count"`
	UsableContextCount      int     `json:"usable_context_count"`
	ToxicContextCount       int     `json:"toxic_context_count"`
	MissingFutureCount      int     `json:"missing_future_count"`
	UsableContextRate       float64 `json:"usable_context_rate"`
	ToxicContextRate        float64 `json:"toxic_context_rate"`
	MissingFutureRate       float64 `json:"missing_future_rate"`
	SessionEdgeRate         float64 `json:"session_edge_rate"`
	SessionEdgeVsWorstPass  bool    `json:"session_edge_vs_worst_pass"`
	DominantContextLabel    string  `json:"dominant_context_label,omitempty"`
	DominantContextRate     float64 `json:"dominant_context_rate"`
	ContainedRotationCount  int     `json:"contained_rotation_count"`
	CleanExpansionUpCount   int     `json:"clean_expansion_up_count"`
	CleanExpansionDownCount int     `json:"clean_expansion_down_count"`
}

type FuturesRangeContextTriageFailureModeRow struct {
	EpisodeID                         int     `json:"episode_id"`
	Timeframe                         string  `json:"timeframe"`
	Split                             string  `json:"split"`
	HorizonBars                       int     `json:"horizon_bars"`
	StartIndex                        int     `json:"start_index"`
	MatureIndex                       int     `json:"mature_index"`
	RawEndIndex                       int     `json:"raw_end_index"`
	MatureCloseTime                   string  `json:"mature_close_time"`
	MatureSession                     string  `json:"mature_session"`
	QualityBucket                     string  `json:"quality_bucket"`
	RangeHigh                         float64 `json:"range_high"`
	RangeLow                          float64 `json:"range_low"`
	RangeMid                          float64 `json:"range_mid"`
	RangeWidth                        float64 `json:"range_width"`
	RangeWidthPct                     float64 `json:"range_width_pct"`
	WidthToATRRatio                   float64 `json:"width_to_atr_ratio"`
	LabelWindowStartIndex             int     `json:"label_window_start_index"`
	LabelWindowEndIndex               int     `json:"label_window_end_index"`
	LabelWindowStartTime              string  `json:"label_window_start_time"`
	LabelWindowEndTime                string  `json:"label_window_end_time"`
	FirstOutsideCloseSide             string  `json:"first_outside_close_side"`
	BarsToFirstOutsideClose           int     `json:"bars_to_first_outside_close"`
	FirstOutsideSession               string  `json:"first_outside_session,omitempty"`
	ReentryWithinWindow               bool    `json:"reentry_within_window"`
	MaxExcursionAboveRangeWidths      float64 `json:"max_excursion_above_range_widths"`
	MaxExcursionBelowRangeWidths      float64 `json:"max_excursion_below_range_widths"`
	FinalClosePosition                float64 `json:"final_close_position"`
	InsideCloseRate                   float64 `json:"inside_close_rate"`
	MidpointCrossCount                int     `json:"midpoint_cross_count"`
	OutsideStateTransitionCount       int     `json:"outside_state_transition_count"`
	UpperHalfInsideCloseSeen          bool    `json:"upper_half_inside_close_seen"`
	LowerHalfInsideCloseSeen          bool    `json:"lower_half_inside_close_seen"`
	MissingFuture                     bool    `json:"missing_future"`
	QuickFailure                      bool    `json:"quick_failure"`
	PrimaryContextLabel               string  `json:"primary_context_label"`
	ConstructiveContext               bool    `json:"constructive_context"`
	ToxicContext                      bool    `json:"toxic_context"`
	SkippedReason                     string  `json:"skipped_reason,omitempty"`
	ParentWidthNormalizedFavorableMax float64 `json:"parent_width_normalized_favorable_max"`
	ParentWidthNormalizedAdverseMax   float64 `json:"parent_width_normalized_adverse_max"`
}

type FuturesRangeContextTriageCohortRow struct {
	CohortID                   string  `json:"cohort_id"`
	CohortType                 string  `json:"cohort_type"`
	Split                      string  `json:"split"`
	Timeframe                  string  `json:"timeframe"`
	HorizonBars                int     `json:"horizon_bars"`
	QualityBucket              string  `json:"quality_bucket,omitempty"`
	MatureSession              string  `json:"mature_session,omitempty"`
	PrimaryContextLabel        string  `json:"primary_context_label,omitempty"`
	CandidateCount             int     `json:"candidate_count"`
	UsableContextCount         int     `json:"usable_context_count"`
	ToxicContextCount          int     `json:"toxic_context_count"`
	MissingFutureCount         int     `json:"missing_future_count"`
	ContainedRotationCount     int     `json:"contained_rotation_count"`
	CleanExpansionUpCount      int     `json:"clean_expansion_up_count"`
	CleanExpansionDownCount    int     `json:"clean_expansion_down_count"`
	FalseBreakReentryUpCount   int     `json:"false_break_reentry_up_count"`
	FalseBreakReentryDownCount int     `json:"false_break_reentry_down_count"`
	BoundaryChopCount          int     `json:"boundary_chop_count"`
	DriftThroughUpCount        int     `json:"drift_through_up_count"`
	DriftThroughDownCount      int     `json:"drift_through_down_count"`
	LowWidthNoiseCount         int     `json:"low_width_noise_count"`
	NoResolutionCount          int     `json:"no_resolution_count"`
	UsableContextRate          float64 `json:"usable_context_rate"`
	ToxicContextRate           float64 `json:"toxic_context_rate"`
	MissingFutureRate          float64 `json:"missing_future_rate"`
	DominantToxicLabel         string  `json:"dominant_toxic_label,omitempty"`
	DominantToxicRate          float64 `json:"dominant_toxic_rate"`
	MinFullCohortCount         int     `json:"min_full_cohort_count"`
	MinSplitCohortCount        int     `json:"min_split_cohort_count"`
	MinSessionSplitCohortCount int     `json:"min_session_split_cohort_count"`
	SourceResamplePass         bool    `json:"source_resample_pass"`
	CountGatePass              bool    `json:"count_gate_pass"`
	UsableGatePass             bool    `json:"usable_gate_pass"`
	ToxicGatePass              bool    `json:"toxic_gate_pass"`
	MissingFutureGatePass      bool    `json:"missing_future_gate_pass"`
	DominantToxicGatePass      bool    `json:"dominant_toxic_gate_pass"`
	ClosedFamilyReslice        bool    `json:"closed_family_reslice"`
	RankableDecisionContext    bool    `json:"rankable_decision_context"`
	PassesReviewGate           bool    `json:"passes_review_gate"`
	FailureReason              string  `json:"failure_reason,omitempty"`
}

type FuturesRangeContextTriageRankingRow struct {
	Rank                    int     `json:"rank"`
	CohortID                string  `json:"cohort_id"`
	CohortType              string  `json:"cohort_type"`
	Timeframe               string  `json:"timeframe"`
	HorizonBars             int     `json:"horizon_bars"`
	QualityBucket           string  `json:"quality_bucket,omitempty"`
	MatureSession           string  `json:"mature_session,omitempty"`
	PassesGate              bool    `json:"passes_gate"`
	RankScore               float64 `json:"rank_score"`
	FullCandidateCount      int     `json:"full_candidate_count"`
	FullUsableContextRate   float64 `json:"full_usable_context_rate"`
	WeakestSplitUsableRate  float64 `json:"weakest_split_usable_rate"`
	FullToxicContextRate    float64 `json:"full_toxic_context_rate"`
	WorstSplitToxicRate     float64 `json:"worst_split_toxic_rate"`
	FullMissingFutureRate   float64 `json:"full_missing_future_rate"`
	SessionEdgeRate         float64 `json:"session_edge_rate"`
	PeriodSplitsRequired    int     `json:"period_splits_required"`
	PeriodSplitsPassing     int     `json:"period_splits_passing"`
	RankableDecisionContext bool    `json:"rankable_decision_context"`
	FailureReason           string  `json:"failure_reason,omitempty"`
}

type FuturesRangeContextTriageSummaryRow struct {
	Split              string `json:"split"`
	Timeframe          string `json:"timeframe"`
	HorizonBars        int    `json:"horizon_bars"`
	SourceResamplePass bool   `json:"source_resample_pass"`
	EpisodeRows        int    `json:"episode_rows"`
	EligibleEpisodes   int    `json:"eligible_episodes"`
	FailureModeRows    int    `json:"failure_mode_rows"`
	CohortRows         int    `json:"cohort_rows"`
	RankingRows        int    `json:"ranking_rows"`
	PassingCohorts     int    `json:"passing_cohorts"`
	StopState          string `json:"stop_state"`
}

type rangeContextTriageCohortKey struct {
	cohortID            string
	cohortType          string
	split               string
	timeframe           string
	horizonBars         int
	qualityBucket       string
	matureSession       string
	primaryContextLabel string
}

type rangeContextTriageCohortAccumulator struct {
	key    rangeContextTriageCohortKey
	labels map[string]int
	row    FuturesRangeContextTriageCohortRow
}

func DefaultFuturesRangeContextTriageAuditConfig() FuturesRangeContextTriageAuditConfig {
	return FuturesRangeContextTriageAuditConfig{
		SourcePath:                   rangeContextTriageSourcePath,
		ApprovedSourcePath:           rangeContextTriageSourcePath,
		ExpectedSourceRows:           rangeContextTriageExpectedRows,
		ExpectedFirstOpenTime:        rangeContextTriageExpectedFirst,
		ExpectedLastOpenTime:         rangeContextTriageExpectedLast,
		ExpectedGapCount:             0,
		ExpectedDuplicateCount:       0,
		ExpectedZeroVolumeCount:      rangeContextTriageExpectedZeroVol,
		Timeframes:                   []string{RangeDiscoveryTimeframe15m, RangeDiscoveryTimeframe1h, RangeDiscoveryTimeframe4h},
		Expected15MRows:              rangeContextTriageExpected15MRows,
		Expected15MLastOpenTime:      rangeContextTriageExpected15MLast,
		Expected1HRows:               rangeContextTriageExpected1HRows,
		Expected1HLastOpenTime:       rangeContextTriageExpected1HLast,
		Expected4HRows:               rangeContextTriageExpected4HRows,
		Expected4HLastOpenTime:       rangeContextTriageExpected4HLast,
		DetectorLookbackDays:         20,
		DetectorPercentile:           0.30,
		DetectorMinConsecutiveBars:   12,
		HorizonsBars:                 []int{12, 24, 48},
		QuickFailureHorizonBars:      6,
		ReentryWindowBars:            6,
		CleanExpansionThreshold:      0.75,
		DriftThreshold:               0.50,
		BoundaryChopTransitions:      3,
		LowWidthPctThreshold:         0.0015,
		LowWidthATRThreshold:         0.75,
		WideWidthPctThreshold:        0.0150,
		WideWidthATRThreshold:        4.00,
		ChoppyMidCrossThreshold:      3,
		ChoppyBoundaryTouchThreshold: 4,
		MinFullCohortCount:           300,
		MinSplitCohortCount:          50,
		MinSessionSplitCohortCount:   30,
		MinUsableContextRateFull:     0.55,
		MinUsableContextRateSplit:    0.45,
		MaxToxicContextRateFull:      0.45,
		MaxToxicContextRateSplit:     0.55,
		MaxMissingFutureRateFull:     0.02,
		MaxDominantToxicRate:         0.50,
	}
}

func RunFuturesRangeContextTriageAudit(candles []Candle, manifest SourceManifest, cfg FuturesRangeContextTriageAuditConfig, splits []Split) (FuturesRangeContextTriageAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeContextTriageAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesRangeContextTriageAuditResult{}
	sourceRow := rangeContextTriageSourceRow(candles, manifest, cfg)
	result.SourceRows = append(result.SourceRows, sourceRow)
	if sourceRow.ValidationStatus != "accepted" || !sourceRow.SourceFactsPass {
		result.StopState = RangeContextTriageStopStateSourceGap
		result.SummaryRows = rangeContextTriageSummaryRows(result, cfg, splits)
		return result, nil
	}

	for _, timeframe := range cfg.Timeframes {
		frame, ok := rangeContextTriageFrameDef(timeframe)
		if !ok {
			return result, fmt.Errorf("range context triage missing frame definition for %s", timeframe)
		}
		frameCandles, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
		coverageRow := rangeContextTriageCoverageRow(coverage, cfg)
		result.CoverageRows = append(result.CoverageRows, coverageRow)
		if err != nil || !coverageRow.CoverageFactsPass || coverageRow.ValidationStatus != "accepted" || !coverageRow.Complete {
			result.StopState = RangeContextTriageStopStateSourceGap
			result.SummaryRows = rangeContextTriageSummaryRows(result, cfg, splits)
			return result, nil
		}
		classifications, err := (CompressionRangeDetector{Config: rangeContextTriageDetectorConfig(cfg, frame.barsPerDay)}).Classify(frameCandles)
		if err != nil {
			return result, err
		}
		normalizedATR := NormalizedATR(frameCandles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
		episodes := rangeContextTriageEpisodeRows(frameCandles, classifications, normalizedATR, timeframe, cfg, splits, len(result.EpisodeRows))
		result.EpisodeRows = append(result.EpisodeRows, episodes...)
	}
	result.QualityRows = rangeContextTriageQualityRows(result.EpisodeRows)
	result.FailureModeRows = rangeContextTriageFailureRows(result.EpisodeRows, candlesByTriageTimeframe(candles, result.CoverageRows, cfg), cfg)
	result.SessionRows = rangeContextTriageSessionRows(result.FailureModeRows)
	result.CohortRows = rangeContextTriageCohortRows(result.FailureModeRows, result.SourceRows, result.CoverageRows, cfg, splits)
	result.RankingRows = rangeContextTriageRankingRows(result.CohortRows, cfg, splits)
	result.PassingCohorts = rangeContextTriagePassingRankingCount(result.RankingRows)
	result.StopState = FuturesRangeContextTriageAuditStopState(result)
	result.SummaryRows = rangeContextTriageSummaryRows(result, cfg, splits)
	for i := range result.CohortRows {
		if result.CohortRows[i].PassesReviewGate && result.StopState != RangeContextTriageStopStateReadyForStrategySpec {
			result.CohortRows[i].FailureReason = appendReason(result.CohortRows[i].FailureReason, "stop_state_not_ready")
		}
	}
	return result, nil
}

func FuturesRangeContextTriageAuditStopState(result FuturesRangeContextTriageAuditResult) string {
	if result.StopState == RangeContextTriageStopStateSourceGap || result.StopState == RangeContextTriageStopStateRejectedClosedReslice {
		return result.StopState
	}
	for _, row := range result.SourceRows {
		if row.ValidationStatus != "accepted" || !row.SourceFactsPass {
			return RangeContextTriageStopStateSourceGap
		}
	}
	for _, row := range result.CoverageRows {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass || !row.Complete {
			return RangeContextTriageStopStateSourceGap
		}
	}
	eligible := 0
	for _, row := range result.EpisodeRows {
		if row.Eligible {
			eligible++
		}
		if row.SkippedReason == "closed_family_reslice" {
			return RangeContextTriageStopStateRejectedClosedReslice
		}
	}
	if eligible == 0 {
		return RangeContextTriageStopStateNoRangeEpisodes
	}
	if len(result.CohortRows) == 0 || len(result.RankingRows) == 0 {
		return RangeContextTriageStopStateNoUsableCohorts
	}
	for _, row := range result.RankingRows {
		if row.PassesGate {
			return RangeContextTriageStopStateReadyForStrategySpec
		}
	}
	return RangeContextTriageStopStateFailedNoStrategyPremise
}

func (cfg FuturesRangeContextTriageAuditConfig) withDefaults() FuturesRangeContextTriageAuditConfig {
	defaults := DefaultFuturesRangeContextTriageAuditConfig()
	if cfg.SourcePath == "" {
		cfg.SourcePath = defaults.SourcePath
	}
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
	if len(cfg.Timeframes) == 0 {
		cfg.Timeframes = append([]string(nil), defaults.Timeframes...)
	}
	if cfg.Expected15MRows == 0 {
		cfg.Expected15MRows = defaults.Expected15MRows
	}
	if cfg.Expected15MLastOpenTime == "" {
		cfg.Expected15MLastOpenTime = defaults.Expected15MLastOpenTime
	}
	if cfg.Expected1HRows == 0 {
		cfg.Expected1HRows = defaults.Expected1HRows
	}
	if cfg.Expected1HLastOpenTime == "" {
		cfg.Expected1HLastOpenTime = defaults.Expected1HLastOpenTime
	}
	if cfg.Expected4HRows == 0 {
		cfg.Expected4HRows = defaults.Expected4HRows
	}
	if cfg.Expected4HLastOpenTime == "" {
		cfg.Expected4HLastOpenTime = defaults.Expected4HLastOpenTime
	}
	if cfg.DetectorLookbackDays == 0 {
		cfg.DetectorLookbackDays = defaults.DetectorLookbackDays
	}
	if cfg.DetectorPercentile == 0 {
		cfg.DetectorPercentile = defaults.DetectorPercentile
	}
	if cfg.DetectorMinConsecutiveBars == 0 {
		cfg.DetectorMinConsecutiveBars = defaults.DetectorMinConsecutiveBars
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickFailureHorizonBars == 0 {
		cfg.QuickFailureHorizonBars = defaults.QuickFailureHorizonBars
	}
	if cfg.ReentryWindowBars == 0 {
		cfg.ReentryWindowBars = defaults.ReentryWindowBars
	}
	if cfg.CleanExpansionThreshold == 0 {
		cfg.CleanExpansionThreshold = defaults.CleanExpansionThreshold
	}
	if cfg.DriftThreshold == 0 {
		cfg.DriftThreshold = defaults.DriftThreshold
	}
	if cfg.BoundaryChopTransitions == 0 {
		cfg.BoundaryChopTransitions = defaults.BoundaryChopTransitions
	}
	if cfg.LowWidthPctThreshold == 0 {
		cfg.LowWidthPctThreshold = defaults.LowWidthPctThreshold
	}
	if cfg.LowWidthATRThreshold == 0 {
		cfg.LowWidthATRThreshold = defaults.LowWidthATRThreshold
	}
	if cfg.WideWidthPctThreshold == 0 {
		cfg.WideWidthPctThreshold = defaults.WideWidthPctThreshold
	}
	if cfg.WideWidthATRThreshold == 0 {
		cfg.WideWidthATRThreshold = defaults.WideWidthATRThreshold
	}
	if cfg.ChoppyMidCrossThreshold == 0 {
		cfg.ChoppyMidCrossThreshold = defaults.ChoppyMidCrossThreshold
	}
	if cfg.ChoppyBoundaryTouchThreshold == 0 {
		cfg.ChoppyBoundaryTouchThreshold = defaults.ChoppyBoundaryTouchThreshold
	}
	if cfg.MinFullCohortCount == 0 {
		cfg.MinFullCohortCount = defaults.MinFullCohortCount
	}
	if cfg.MinSplitCohortCount == 0 {
		cfg.MinSplitCohortCount = defaults.MinSplitCohortCount
	}
	if cfg.MinSessionSplitCohortCount == 0 {
		cfg.MinSessionSplitCohortCount = defaults.MinSessionSplitCohortCount
	}
	if cfg.MinUsableContextRateFull == 0 {
		cfg.MinUsableContextRateFull = defaults.MinUsableContextRateFull
	}
	if cfg.MinUsableContextRateSplit == 0 {
		cfg.MinUsableContextRateSplit = defaults.MinUsableContextRateSplit
	}
	if cfg.MaxToxicContextRateFull == 0 {
		cfg.MaxToxicContextRateFull = defaults.MaxToxicContextRateFull
	}
	if cfg.MaxToxicContextRateSplit == 0 {
		cfg.MaxToxicContextRateSplit = defaults.MaxToxicContextRateSplit
	}
	if cfg.MaxMissingFutureRateFull == 0 {
		cfg.MaxMissingFutureRateFull = defaults.MaxMissingFutureRateFull
	}
	if cfg.MaxDominantToxicRate == 0 {
		cfg.MaxDominantToxicRate = defaults.MaxDominantToxicRate
	}
	return cfg
}

func (cfg FuturesRangeContextTriageAuditConfig) validate() error {
	if cfg.SourcePath == "" || cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("range context triage source paths must not be empty")
	}
	for _, timeframe := range cfg.Timeframes {
		if timeframe != RangeDiscoveryTimeframe15m && timeframe != RangeDiscoveryTimeframe1h && timeframe != RangeDiscoveryTimeframe4h {
			return fmt.Errorf("range context triage unsupported timeframe %q", timeframe)
		}
	}
	if cfg.DetectorLookbackDays <= 0 && cfg.DetectorLookbackBarsOverride <= 0 {
		return fmt.Errorf("range context triage detector lookback must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("range context triage detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("range context triage detector min consecutive bars must be positive")
	}
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("range context triage horizon bars must be positive")
		}
	}
	if cfg.QuickFailureHorizonBars <= 0 || cfg.ReentryWindowBars <= 0 || cfg.BoundaryChopTransitions <= 0 {
		return fmt.Errorf("range context triage horizon thresholds must be positive")
	}
	if cfg.CleanExpansionThreshold <= 0 || cfg.DriftThreshold <= 0 {
		return fmt.Errorf("range context triage expansion thresholds must be positive")
	}
	if cfg.LowWidthPctThreshold <= 0 || cfg.LowWidthATRThreshold <= 0 || cfg.WideWidthPctThreshold <= 0 || cfg.WideWidthATRThreshold <= 0 {
		return fmt.Errorf("range context triage quality thresholds must be positive")
	}
	if cfg.MinFullCohortCount <= 0 || cfg.MinSplitCohortCount <= 0 || cfg.MinSessionSplitCohortCount <= 0 {
		return fmt.Errorf("range context triage count gates must be positive")
	}
	if cfg.MinUsableContextRateFull <= 0 || cfg.MinUsableContextRateSplit <= 0 || cfg.MaxToxicContextRateFull <= 0 || cfg.MaxToxicContextRateSplit <= 0 || cfg.MaxMissingFutureRateFull <= 0 || cfg.MaxDominantToxicRate <= 0 {
		return fmt.Errorf("range context triage rate gates must be positive")
	}
	return nil
}

func rangeContextTriageSourceRow(parent []Candle, manifest SourceManifest, cfg FuturesRangeContextTriageAuditConfig) FuturesRangeContextTriageSourceRow {
	row := FuturesRangeContextTriageSourceRow{
		Path:                    manifest.Path,
		ApprovedPath:            cfg.ApprovedSourcePath,
		Venue:                   manifest.Venue,
		Product:                 manifest.Product,
		Symbol:                  manifest.Symbol,
		Interval:                manifest.Interval,
		RowCount:                manifest.RowCount,
		ExpectedRowCount:        cfg.ExpectedSourceRows,
		FirstOpenTime:           manifest.FirstOpenTime,
		ExpectedFirstOpenTime:   cfg.ExpectedFirstOpenTime,
		LastOpenTime:            manifest.LastOpenTime,
		ExpectedLastOpenTime:    cfg.ExpectedLastOpenTime,
		GapCount:                manifest.GapCount,
		ExpectedGapCount:        cfg.ExpectedGapCount,
		DuplicateCount:          manifest.DuplicateCount,
		ExpectedDuplicateCount:  cfg.ExpectedDuplicateCount,
		ZeroVolumeCount:         manifest.ZeroVolumeCount,
		ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount,
		ComparisonOnly:          manifest.ComparisonOnly,
		SourceFactsPass:         true,
		ValidationStatus:        "accepted",
	}
	reject := func(reason string) {
		row.SourceFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = reason
	}
	if manifest.ValidationStatus != "accepted" {
		reject(manifest.ValidationError)
		return row
	}
	if manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly || manifest.Symbol != RangeUniverseSymbolBTCUSDT || manifest.Interval != "5m" {
		reject(fmt.Sprintf("range context triage requires BTCUSDT Binance USDT-M futures 5m comparison_only=false; got product=%q symbol=%q interval=%q comparison_only=%t", manifest.Product, manifest.Symbol, manifest.Interval, manifest.ComparisonOnly))
		return row
	}
	if !sameCleanPath(manifest.Path, cfg.ApprovedSourcePath) {
		reject(fmt.Sprintf("range context triage source path %q is not approved path %q", manifest.Path, cfg.ApprovedSourcePath))
		return row
	}
	if cfg.SkipSourceFactCheck {
		return row
	}
	switch {
	case len(parent) != cfg.ExpectedSourceRows || manifest.RowCount != cfg.ExpectedSourceRows:
		reject(fmt.Sprintf("range context triage source rows=%d manifest_rows=%d expected=%d", len(parent), manifest.RowCount, cfg.ExpectedSourceRows))
	case manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime:
		reject(fmt.Sprintf("range context triage source first_open_time=%s expected=%s", manifest.FirstOpenTime, cfg.ExpectedFirstOpenTime))
	case manifest.LastOpenTime != cfg.ExpectedLastOpenTime:
		reject(fmt.Sprintf("range context triage source last_open_time=%s expected=%s", manifest.LastOpenTime, cfg.ExpectedLastOpenTime))
	case manifest.GapCount != cfg.ExpectedGapCount:
		reject(fmt.Sprintf("range context triage source gap_count=%d expected=%d", manifest.GapCount, cfg.ExpectedGapCount))
	case manifest.DuplicateCount != cfg.ExpectedDuplicateCount:
		reject(fmt.Sprintf("range context triage source duplicate_count=%d expected=%d", manifest.DuplicateCount, cfg.ExpectedDuplicateCount))
	case manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount:
		reject(fmt.Sprintf("range context triage source zero_volume_count=%d expected=%d", manifest.ZeroVolumeCount, cfg.ExpectedZeroVolumeCount))
	}
	return row
}

func rangeContextTriageCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesRangeContextTriageAuditConfig) FuturesRangeContextTriageCoverageRow {
	row := FuturesRangeContextTriageCoverageRow{
		FuturesRangeDiscoveryCoverageRow: base,
		ExpectedFirstOpenTime:            cfg.ExpectedFirstOpenTime,
		CoverageFactsPass:                true,
	}
	switch base.Timeframe {
	case RangeDiscoveryTimeframe15m:
		row.ExpectedRowCount = cfg.Expected15MRows
		row.ExpectedLastOpenTime = cfg.Expected15MLastOpenTime
	case RangeDiscoveryTimeframe1h:
		row.ExpectedRowCount = cfg.Expected1HRows
		row.ExpectedLastOpenTime = cfg.Expected1HLastOpenTime
	case RangeDiscoveryTimeframe4h:
		row.ExpectedRowCount = cfg.Expected4HRows
		row.ExpectedLastOpenTime = cfg.Expected4HLastOpenTime
	default:
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("unsupported range context triage timeframe %q", base.Timeframe)
		return row
	}
	if base.ValidationStatus != "accepted" || !base.Complete {
		row.CoverageFactsPass = false
		return row
	}
	if !cfg.SkipCoverageCountCheck && (row.RowCount != row.ExpectedRowCount || row.FirstOpenTime != row.ExpectedFirstOpenTime || row.LastOpenTime != row.ExpectedLastOpenTime) {
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("%s coverage row_count=%d first=%s last=%s expected row_count=%d first=%s last=%s", base.Timeframe, row.RowCount, row.FirstOpenTime, row.LastOpenTime, row.ExpectedRowCount, row.ExpectedFirstOpenTime, row.ExpectedLastOpenTime)
	}
	if row.GapCount != 0 || row.DuplicateCount != 0 || row.MissingChildOpenCount != 0 {
		row.CoverageFactsPass = false
		if row.ValidationError == "" {
			row.ValidationStatus = "rejected"
			row.ValidationError = fmt.Sprintf("%s coverage gap=%d duplicate=%d missing_child_open=%d", base.Timeframe, row.GapCount, row.DuplicateCount, row.MissingChildOpenCount)
		}
	}
	return row
}

func rangeContextTriageFrameDef(timeframe string) (rangeDiscoveryFrameDef, bool) {
	for _, frame := range rangeDiscoveryFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, true
		}
	}
	return rangeDiscoveryFrameDef{}, false
}

func rangeContextTriageDetectorConfig(cfg FuturesRangeContextTriageAuditConfig, barsPerDay int) RangeDetectorConfig {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = barsPerDay
	detectorCfg.LookbackDays = cfg.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DetectorMinConsecutiveBars
	detectorCfg.UseBollinger = true
	detectorCfg.UseADX = false
	if cfg.DetectorLookbackBarsOverride > 0 {
		detectorCfg.LookbackDays = 1
		detectorCfg.BarsPerDay = cfg.DetectorLookbackBarsOverride
	}
	return detectorCfg
}

func rangeContextTriageEpisodeRows(candles []Candle, classifications []RangeClassification, normalizedATR []float64, timeframe string, cfg FuturesRangeContextTriageAuditConfig, splits []Split, offset int) []FuturesRangeContextTriageEpisodeRow {
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	rows := []FuturesRangeContextTriageEpisodeRow{}
	for i := 0; i < len(classifications); {
		if !classifications[i].RawActive {
			i++
			continue
		}
		start := i
		firstActive := -1
		rawHigh := candles[i].High
		rawLow := candles[i].Low
		activeBars := 0
		for i < len(classifications) && classifications[i].RawActive {
			if candles[i].High > rawHigh {
				rawHigh = candles[i].High
			}
			if candles[i].Low < rawLow {
				rawLow = candles[i].Low
			}
			if classifications[i].Active {
				activeBars++
				if firstActive == -1 {
					firstActive = i
				}
			}
			i++
		}
		rawEnd := i - 1
		row := FuturesRangeContextTriageEpisodeRow{
			EpisodeID:            offset + len(rows) + 1,
			Timeframe:            timeframe,
			DetectorProfileID:    rangeContextTriageDetectorProfileID,
			StartIndex:           start,
			MatureIndex:          firstActive,
			RawEndIndex:          rawEnd,
			High:                 rawHigh,
			Low:                  rawLow,
			ActiveBarsToMaturity: activeBars,
			SkippedReason:        "",
		}
		if start >= 0 && start < len(candles) {
			row.StartCloseTime = candles[start].CloseTime.UTC().Format(timeLayout)
		}
		if rawEnd >= 0 && rawEnd < len(candles) {
			row.RawEndCloseTime = candles[rawEnd].CloseTime.UTC().Format(timeLayout)
		}
		if firstActive < 0 {
			row.SkippedReason = "no_mature_active_bar"
			rows = append(rows, row)
			continue
		}
		row.MatureCloseTime = candles[firstActive].CloseTime.UTC().Format(timeLayout)
		row.Split = splitNameForCloseTime(candles[firstActive].CloseTime, splits)
		row.MatureSession = RangeContextTriageUTCSession(candles[firstActive].CloseTime)
		freezeHigh := candles[start].High
		freezeLow := candles[start].Low
		atrSum := 0.0
		atrCount := 0
		for j := start; j <= firstActive; j++ {
			if candles[j].High > freezeHigh {
				freezeHigh = candles[j].High
			}
			if candles[j].Low < freezeLow {
				freezeLow = candles[j].Low
			}
			if j < len(normalizedATR) && validNumber(normalizedATR[j]) {
				atrSum += normalizedATR[j]
				atrCount++
			}
		}
		row.High = freezeHigh
		row.Low = freezeLow
		row.Width = freezeHigh - freezeLow
		row.Mid = (freezeHigh + freezeLow) / 2
		row.LowerQuartile = freezeLow + row.Width*0.25
		row.UpperQuartile = freezeLow + row.Width*0.75
		row.DurationBarsToMaturity = firstActive - start + 1
		row.ActiveBarsToMaturity = 0
		for j := start; j <= firstActive; j++ {
			if classifications[j].Active {
				row.ActiveBarsToMaturity++
			}
		}
		if row.Width <= 0 || !finitePositive(row.Width) {
			row.SkippedReason = "non_positive_width"
			rows = append(rows, row)
			continue
		}
		if candles[firstActive].Close <= 0 || !finitePositive(candles[firstActive].Close) {
			row.SkippedReason = "non_positive_price"
			rows = append(rows, row)
			continue
		}
		row.WidthPct = row.Width / candles[firstActive].Close
		if atrCount == 0 || atrSum <= 0 {
			row.SkippedReason = "missing_atr"
			rows = append(rows, row)
			continue
		}
		avgATR := atrSum / float64(atrCount)
		row.WidthToATRRatio = row.WidthPct / avgATR
		row.MatureClosePosition = (candles[firstActive].Close - row.Low) / row.Width
		row.PreMatureMidCrossCount = rangeContextTriageMidCrossCount(candles, start, firstActive, row.Mid)
		row.PreMatureBoundaryTouchCount = rangeContextTriageBoundaryTouchCount(candles, start, firstActive, row.Low, row.High)
		row.PreMatureCloseInsideRate = rangeContextTriageInsideCloseRate(candles, start, firstActive, row.Low, row.High)
		row.PreMatureWickOvershootCount = rangeContextTriageWickOvershootCount(candles, start, firstActive, row.Low, row.High)
		row.QualityBucket = rangeContextTriageQualityBucket(row, cfg)
		if firstActive+1 >= len(candles) {
			row.SkippedReason = "missing_future"
			rows = append(rows, row)
			continue
		}
		row.Eligible = true
		rows = append(rows, row)
	}
	return rows
}

func rangeContextTriageQualityRows(episodes []FuturesRangeContextTriageEpisodeRow) []FuturesRangeContextTriageQualityRow {
	rows := make([]FuturesRangeContextTriageQualityRow, 0, len(episodes))
	for _, episode := range episodes {
		rows = append(rows, FuturesRangeContextTriageQualityRow{
			EpisodeID:                   episode.EpisodeID,
			Timeframe:                   episode.Timeframe,
			Split:                       episode.Split,
			MatureCloseTime:             episode.MatureCloseTime,
			WidthPct:                    episode.WidthPct,
			WidthToATRRatio:             episode.WidthToATRRatio,
			DurationBarsToMaturity:      episode.DurationBarsToMaturity,
			ActiveBarsToMaturity:        episode.ActiveBarsToMaturity,
			MatureClosePosition:         episode.MatureClosePosition,
			PreMatureMidCrossCount:      episode.PreMatureMidCrossCount,
			PreMatureBoundaryTouchCount: episode.PreMatureBoundaryTouchCount,
			PreMatureCloseInsideRate:    episode.PreMatureCloseInsideRate,
			PreMatureWickOvershootCount: episode.PreMatureWickOvershootCount,
			QualityBucket:               episode.QualityBucket,
			Eligible:                    episode.Eligible,
			SkippedReason:               episode.SkippedReason,
		})
	}
	return rows
}

func candlesByTriageTimeframe(parent []Candle, coverage []FuturesRangeContextTriageCoverageRow, cfg FuturesRangeContextTriageAuditConfig) map[string][]Candle {
	out := map[string][]Candle{}
	for _, timeframe := range cfg.Timeframes {
		frame, ok := rangeContextTriageFrameDef(timeframe)
		if !ok {
			continue
		}
		candles, _, err := resampleRangeDiscoveryFrame(parent, frame)
		if err == nil {
			out[timeframe] = candles
		}
	}
	return out
}

func rangeContextTriageFailureRows(episodes []FuturesRangeContextTriageEpisodeRow, candlesByTimeframe map[string][]Candle, cfg FuturesRangeContextTriageAuditConfig) []FuturesRangeContextTriageFailureModeRow {
	rows := []FuturesRangeContextTriageFailureModeRow{}
	for _, episode := range episodes {
		if !episode.Eligible {
			continue
		}
		candles := candlesByTimeframe[episode.Timeframe]
		for _, horizon := range cfg.HorizonsBars {
			rows = append(rows, rangeContextTriageFailureRow(candles, episode, horizon, cfg))
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].EpisodeID != rows[j].EpisodeID {
			return rows[i].EpisodeID < rows[j].EpisodeID
		}
		return rows[i].HorizonBars < rows[j].HorizonBars
	})
	return rows
}

func rangeContextTriageFailureRow(candles []Candle, episode FuturesRangeContextTriageEpisodeRow, horizon int, cfg FuturesRangeContextTriageAuditConfig) FuturesRangeContextTriageFailureModeRow {
	row := FuturesRangeContextTriageFailureModeRow{
		EpisodeID:               episode.EpisodeID,
		Timeframe:               episode.Timeframe,
		Split:                   episode.Split,
		HorizonBars:             horizon,
		StartIndex:              episode.StartIndex,
		MatureIndex:             episode.MatureIndex,
		RawEndIndex:             episode.RawEndIndex,
		MatureCloseTime:         episode.MatureCloseTime,
		MatureSession:           episode.MatureSession,
		QualityBucket:           episode.QualityBucket,
		RangeHigh:               episode.High,
		RangeLow:                episode.Low,
		RangeMid:                episode.Mid,
		RangeWidth:              episode.Width,
		RangeWidthPct:           episode.WidthPct,
		WidthToATRRatio:         episode.WidthToATRRatio,
		LabelWindowStartIndex:   -1,
		LabelWindowEndIndex:     -1,
		FirstOutsideCloseSide:   "none",
		BarsToFirstOutsideClose: -1,
		PrimaryContextLabel:     RangeContextTriageLabelNoResolution,
	}
	start := episode.MatureIndex + 1
	end := episode.MatureIndex + horizon
	if start < len(candles) {
		row.LabelWindowStartIndex = start
		row.LabelWindowStartTime = candles[start].CloseTime.UTC().Format(timeLayout)
		row.LabelWindowEndIndex = minInt(end, len(candles)-1)
		row.LabelWindowEndTime = candles[row.LabelWindowEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	if horizon <= 0 || end >= len(candles) || episode.Width <= 0 {
		row.MissingFuture = true
		row.PrimaryContextLabel = RangeContextTriageLabelMissingFuture
		return rangeContextTriageFinalizeFailureRow(row)
	}
	row.LabelWindowEndIndex = end
	row.LabelWindowEndTime = candles[end].CloseTime.UTC().Format(timeLayout)
	insideCloses := 0
	prevMidSide := 0
	prevOutsideState := ""
	outsideTransitionsStarted := false
	firstOutsideIndex := -1
	reentryDeadline := -1
	for i := start; i <= end; i++ {
		candle := candles[i]
		delay := i - episode.MatureIndex
		if candle.High > episode.High {
			row.MaxExcursionAboveRangeWidths = math.Max(row.MaxExcursionAboveRangeWidths, (candle.High-episode.High)/episode.Width)
		}
		if candle.Low < episode.Low {
			row.MaxExcursionBelowRangeWidths = math.Max(row.MaxExcursionBelowRangeWidths, (episode.Low-candle.Low)/episode.Width)
		}
		inside := candle.Close >= episode.Low && candle.Close <= episode.High
		if inside {
			insideCloses++
			if candle.Close >= episode.Mid {
				row.UpperHalfInsideCloseSeen = true
			}
			if candle.Close <= episode.Mid {
				row.LowerHalfInsideCloseSeen = true
			}
		}
		midSide := 0
		if candle.Close > episode.Mid {
			midSide = 1
		} else if candle.Close < episode.Mid {
			midSide = -1
		}
		if prevMidSide != 0 && midSide != 0 && midSide != prevMidSide {
			row.MidpointCrossCount++
		}
		if midSide != 0 {
			prevMidSide = midSide
		}
		state := "inside"
		if candle.Close > episode.High {
			state = RangeDiscoverySideUp
		} else if candle.Close < episode.Low {
			state = RangeDiscoverySideDown
		}
		if firstOutsideIndex < 0 && state != "inside" {
			firstOutsideIndex = i
			row.BarsToFirstOutsideClose = delay
			row.FirstOutsideCloseSide = state
			row.FirstOutsideSession = RangeContextTriageUTCSession(candle.CloseTime)
			reentryDeadline = i + cfg.ReentryWindowBars
			if delay <= cfg.QuickFailureHorizonBars {
				row.QuickFailure = true
			}
		}
		if firstOutsideIndex >= 0 {
			if outsideTransitionsStarted && state != prevOutsideState {
				row.OutsideStateTransitionCount++
			}
			prevOutsideState = state
			outsideTransitionsStarted = true
			if !row.ReentryWithinWindow && i > firstOutsideIndex && i <= reentryDeadline && inside {
				row.ReentryWithinWindow = true
			}
		}
	}
	count := end - start + 1
	if count > 0 {
		row.InsideCloseRate = float64(insideCloses) / float64(count)
	}
	row.FinalClosePosition = (candles[end].Close - episode.Low) / episode.Width
	row.ParentWidthNormalizedFavorableMax = math.Max(row.MaxExcursionAboveRangeWidths, row.MaxExcursionBelowRangeWidths)
	row.ParentWidthNormalizedAdverseMax = math.Min(row.MaxExcursionAboveRangeWidths, row.MaxExcursionBelowRangeWidths)
	row.PrimaryContextLabel = rangeContextTriagePrimaryLabel(row, cfg)
	return rangeContextTriageFinalizeFailureRow(row)
}

func rangeContextTriagePrimaryLabel(row FuturesRangeContextTriageFailureModeRow, cfg FuturesRangeContextTriageAuditConfig) string {
	switch {
	case row.MissingFuture:
		return RangeContextTriageLabelMissingFuture
	case row.QualityBucket == RangeContextTriageQualityTooNarrowNoise:
		return RangeContextTriageLabelLowWidthNoise
	case row.FirstOutsideCloseSide == "none" && row.MidpointCrossCount >= 2 && row.UpperHalfInsideCloseSeen && row.LowerHalfInsideCloseSeen:
		return RangeContextTriageLabelContainedRotation
	case row.OutsideStateTransitionCount >= cfg.BoundaryChopTransitions:
		return RangeContextTriageLabelBoundaryChop
	case row.FirstOutsideCloseSide == RangeDiscoverySideUp && row.ReentryWithinWindow:
		return RangeContextTriageLabelFalseBreakReentryUp
	case row.FirstOutsideCloseSide == RangeDiscoverySideDown && row.ReentryWithinWindow:
		return RangeContextTriageLabelFalseBreakReentryDown
	case row.FirstOutsideCloseSide == RangeDiscoverySideUp && row.MaxExcursionAboveRangeWidths >= cfg.CleanExpansionThreshold:
		return RangeContextTriageLabelCleanExpansionUp
	case row.FirstOutsideCloseSide == RangeDiscoverySideDown && row.MaxExcursionBelowRangeWidths >= cfg.CleanExpansionThreshold:
		return RangeContextTriageLabelCleanExpansionDown
	case row.FinalClosePosition > 1 && row.FinalClosePosition-1 < cfg.DriftThreshold:
		return RangeContextTriageLabelDriftThroughUp
	case row.FinalClosePosition < 0 && -row.FinalClosePosition < cfg.DriftThreshold:
		return RangeContextTriageLabelDriftThroughDown
	default:
		return RangeContextTriageLabelNoResolution
	}
}

func rangeContextTriageFinalizeFailureRow(row FuturesRangeContextTriageFailureModeRow) FuturesRangeContextTriageFailureModeRow {
	row.ConstructiveContext = rangeContextTriageConstructiveLabel(row.PrimaryContextLabel)
	row.ToxicContext = rangeContextTriageToxicLabel(row.PrimaryContextLabel)
	return row
}

func rangeContextTriageSessionRows(failures []FuturesRangeContextTriageFailureModeRow) []FuturesRangeContextTriageSessionRow {
	type key struct {
		split               string
		timeframe           string
		horizon             int
		matureSession       string
		firstOutsideSession string
	}
	acc := map[key]*FuturesRangeContextTriageSessionRow{}
	labelCounts := map[key]map[string]int{}
	for _, row := range failures {
		for _, split := range rangeDiscoverySplitCombos(row.Split) {
			k := key{split: split, timeframe: row.Timeframe, horizon: row.HorizonBars, matureSession: row.MatureSession, firstOutsideSession: row.FirstOutsideSession}
			session := acc[k]
			if session == nil {
				session = &FuturesRangeContextTriageSessionRow{
					Split:               split,
					Timeframe:           row.Timeframe,
					HorizonBars:         row.HorizonBars,
					MatureSession:       row.MatureSession,
					FirstOutsideSession: row.FirstOutsideSession,
				}
				acc[k] = session
				labelCounts[k] = map[string]int{}
			}
			session.CandidateCount++
			labelCounts[k][row.PrimaryContextLabel]++
			if row.ConstructiveContext {
				session.UsableContextCount++
			}
			if row.ToxicContext {
				session.ToxicContextCount++
			}
			if row.MissingFuture {
				session.MissingFutureCount++
			}
			switch row.PrimaryContextLabel {
			case RangeContextTriageLabelContainedRotation:
				session.ContainedRotationCount++
			case RangeContextTriageLabelCleanExpansionUp:
				session.CleanExpansionUpCount++
			case RangeContextTriageLabelCleanExpansionDown:
				session.CleanExpansionDownCount++
			}
		}
	}
	rows := make([]FuturesRangeContextTriageSessionRow, 0, len(acc))
	for k, row := range acc {
		if row.CandidateCount > 0 {
			denom := float64(row.CandidateCount)
			row.UsableContextRate = float64(row.UsableContextCount) / denom
			row.ToxicContextRate = float64(row.ToxicContextCount) / denom
			row.MissingFutureRate = float64(row.MissingFutureCount) / denom
		}
		row.DominantContextLabel, row.DominantContextRate = rangeContextTriageDominantLabel(labelCounts[k], row.CandidateCount, false)
		rows = append(rows, *row)
	}
	rangeContextTriageMarkSessionEdges(rows)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Split != rows[j].Split {
			return rangeContextTriageSplitSortKey(rows[i].Split) < rangeContextTriageSplitSortKey(rows[j].Split)
		}
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].HorizonBars != rows[j].HorizonBars {
			return rows[i].HorizonBars < rows[j].HorizonBars
		}
		if rows[i].MatureSession != rows[j].MatureSession {
			return rows[i].MatureSession < rows[j].MatureSession
		}
		return rows[i].FirstOutsideSession < rows[j].FirstOutsideSession
	})
	return rows
}

func rangeContextTriageMarkSessionEdges(rows []FuturesRangeContextTriageSessionRow) {
	type groupKey struct {
		split     string
		timeframe string
		horizon   int
	}
	minRate := map[groupKey]float64{}
	for _, row := range rows {
		if row.CandidateCount == 0 || row.FirstOutsideSession != "" {
			continue
		}
		k := groupKey{split: row.Split, timeframe: row.Timeframe, horizon: row.HorizonBars}
		if _, ok := minRate[k]; !ok || row.UsableContextRate < minRate[k] {
			minRate[k] = row.UsableContextRate
		}
	}
	for i := range rows {
		if rows[i].FirstOutsideSession != "" {
			continue
		}
		k := groupKey{split: rows[i].Split, timeframe: rows[i].Timeframe, horizon: rows[i].HorizonBars}
		rows[i].SessionEdgeRate = rows[i].UsableContextRate - minRate[k]
		rows[i].SessionEdgeVsWorstPass = rows[i].SessionEdgeRate >= 0.15
	}
}

func rangeContextTriageCohortRows(failures []FuturesRangeContextTriageFailureModeRow, sources []FuturesRangeContextTriageSourceRow, coverage []FuturesRangeContextTriageCoverageRow, cfg FuturesRangeContextTriageAuditConfig, splits []Split) []FuturesRangeContextTriageCohortRow {
	acc := map[rangeContextTriageCohortKey]*rangeContextTriageCohortAccumulator{}
	for _, row := range failures {
		for _, split := range rangeDiscoverySplitCombos(row.Split) {
			keys := []rangeContextTriageCohortKey{
				rangeContextTriageCohortKeyFor(split, row.Timeframe, row.HorizonBars, RangeContextTriageCohortAll, "", "", ""),
				rangeContextTriageCohortKeyFor(split, row.Timeframe, row.HorizonBars, RangeContextTriageCohortQualityBucket, row.QualityBucket, "", ""),
				rangeContextTriageCohortKeyFor(split, row.Timeframe, row.HorizonBars, RangeContextTriageCohortMatureSession, "", row.MatureSession, ""),
				rangeContextTriageCohortKeyFor(split, row.Timeframe, row.HorizonBars, RangeContextTriageCohortPrimaryContextLabel, "", "", row.PrimaryContextLabel),
				rangeContextTriageCohortKeyFor(split, row.Timeframe, row.HorizonBars, RangeContextTriageCohortQualityBucketMatureSession, row.QualityBucket, row.MatureSession, ""),
			}
			for _, key := range keys {
				a := acc[key]
				if a == nil {
					a = &rangeContextTriageCohortAccumulator{key: key, labels: map[string]int{}}
					a.row = FuturesRangeContextTriageCohortRow{
						CohortID:                   key.cohortID,
						CohortType:                 key.cohortType,
						Split:                      key.split,
						Timeframe:                  key.timeframe,
						HorizonBars:                key.horizonBars,
						QualityBucket:              key.qualityBucket,
						MatureSession:              key.matureSession,
						PrimaryContextLabel:        key.primaryContextLabel,
						MinFullCohortCount:         cfg.MinFullCohortCount,
						MinSplitCohortCount:        cfg.MinSplitCohortCount,
						MinSessionSplitCohortCount: cfg.MinSessionSplitCohortCount,
						SourceResamplePass:         rangeContextTriageSourceResamplePass(sources, coverage),
						RankableDecisionContext:    key.cohortType != RangeContextTriageCohortPrimaryContextLabel,
					}
					acc[key] = a
				}
				a.add(row)
			}
		}
	}
	rows := make([]FuturesRangeContextTriageCohortRow, 0, len(acc))
	for _, a := range acc {
		row := a.finalRow()
		rows = append(rows, row)
	}
	rangeContextTriageMarkCohortGates(rows, cfg, splits, rangeContextTriageSourceResamplePass(sources, coverage))
	sort.Slice(rows, func(i, j int) bool {
		return rangeContextTriageLessCohort(rows[i], rows[j])
	})
	return rows
}

func (acc *rangeContextTriageCohortAccumulator) add(row FuturesRangeContextTriageFailureModeRow) {
	acc.row.CandidateCount++
	acc.labels[row.PrimaryContextLabel]++
	if row.ConstructiveContext {
		acc.row.UsableContextCount++
	}
	if row.ToxicContext {
		acc.row.ToxicContextCount++
	}
	if row.MissingFuture {
		acc.row.MissingFutureCount++
	}
	switch row.PrimaryContextLabel {
	case RangeContextTriageLabelContainedRotation:
		acc.row.ContainedRotationCount++
	case RangeContextTriageLabelCleanExpansionUp:
		acc.row.CleanExpansionUpCount++
	case RangeContextTriageLabelCleanExpansionDown:
		acc.row.CleanExpansionDownCount++
	case RangeContextTriageLabelFalseBreakReentryUp:
		acc.row.FalseBreakReentryUpCount++
	case RangeContextTriageLabelFalseBreakReentryDown:
		acc.row.FalseBreakReentryDownCount++
	case RangeContextTriageLabelBoundaryChop:
		acc.row.BoundaryChopCount++
	case RangeContextTriageLabelDriftThroughUp:
		acc.row.DriftThroughUpCount++
	case RangeContextTriageLabelDriftThroughDown:
		acc.row.DriftThroughDownCount++
	case RangeContextTriageLabelLowWidthNoise:
		acc.row.LowWidthNoiseCount++
	case RangeContextTriageLabelNoResolution:
		acc.row.NoResolutionCount++
	}
}

func (acc *rangeContextTriageCohortAccumulator) finalRow() FuturesRangeContextTriageCohortRow {
	row := acc.row
	if row.CandidateCount > 0 {
		denom := float64(row.CandidateCount)
		row.UsableContextRate = float64(row.UsableContextCount) / denom
		row.ToxicContextRate = float64(row.ToxicContextCount) / denom
		row.MissingFutureRate = float64(row.MissingFutureCount) / denom
	}
	row.DominantToxicLabel, row.DominantToxicRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, true)
	return row
}

func rangeContextTriageMarkCohortGates(rows []FuturesRangeContextTriageCohortRow, cfg FuturesRangeContextTriageAuditConfig, splits []Split, sourcePass bool) {
	byIDSplit := map[string]map[string]*FuturesRangeContextTriageCohortRow{}
	for i := range rows {
		if byIDSplit[rows[i].CohortID] == nil {
			byIDSplit[rows[i].CohortID] = map[string]*FuturesRangeContextTriageCohortRow{}
		}
		byIDSplit[rows[i].CohortID][rows[i].Split] = &rows[i]
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	for i := range rows {
		cohortRows := byIDSplit[rows[i].CohortID]
		full := cohortRows[fullSplitName]
		reasons := []string{}
		fullCountPass := full != nil && full.CandidateCount >= cfg.MinFullCohortCount
		splitCountPass := true
		usablePass := full != nil && full.UsableContextRate >= cfg.MinUsableContextRateFull
		toxicPass := full != nil && full.ToxicContextRate <= cfg.MaxToxicContextRateFull
		missingPass := full != nil && full.MissingFutureRate <= cfg.MaxMissingFutureRateFull
		dominantPass := true
		requiredSplitCount := cfg.MinSplitCohortCount
		if rows[i].MatureSession != "" {
			requiredSplitCount = cfg.MinSessionSplitCohortCount
		}
		periodPresent := 0
		for _, split := range periodSplits {
			row := cohortRows[split.Name]
			if row == nil {
				splitCountPass = false
				usablePass = false
				toxicPass = false
				continue
			}
			periodPresent++
			if row.CandidateCount < requiredSplitCount {
				splitCountPass = false
			}
			if row.UsableContextRate < cfg.MinUsableContextRateSplit {
				usablePass = false
			}
			if row.ToxicContextRate > cfg.MaxToxicContextRateSplit {
				toxicPass = false
			}
		}
		if periodPresent != len(periodSplits) {
			splitCountPass = false
		}
		for _, label := range rangeContextTriageToxicLabels() {
			allDominant := len(periodSplits) > 0
			for _, split := range periodSplits {
				row := cohortRows[split.Name]
				if row == nil || row.DominantToxicLabel != label || row.DominantToxicRate <= cfg.MaxDominantToxicRate {
					allDominant = false
					break
				}
			}
			if allDominant {
				dominantPass = false
				break
			}
		}
		if !sourcePass {
			reasons = append(reasons, "source_or_resample_gap")
		}
		if !fullCountPass {
			reasons = append(reasons, "inadequate_full_cohort_count")
		}
		if !splitCountPass {
			reasons = append(reasons, "inadequate_split_cohort_count")
		}
		if !usablePass {
			reasons = append(reasons, "usable_context_rate_below_gate")
		}
		if !toxicPass {
			reasons = append(reasons, "toxic_context_rate_above_gate")
		}
		if !missingPass {
			reasons = append(reasons, "missing_future_rate_above_gate")
		}
		if !dominantPass {
			reasons = append(reasons, "single_toxic_label_dominates_all_splits")
		}
		if rows[i].ClosedFamilyReslice {
			reasons = append(reasons, "closed_family_reslice")
		}
		if !rows[i].RankableDecisionContext {
			reasons = append(reasons, rangeContextTriageCohortRankableDecisionContextNote)
		}
		rows[i].SourceResamplePass = sourcePass
		rows[i].CountGatePass = fullCountPass && splitCountPass
		if rows[i].Split != fullSplitName {
			rows[i].CountGatePass = rows[i].CandidateCount >= requiredSplitCount
		}
		rows[i].UsableGatePass = usablePass
		rows[i].ToxicGatePass = toxicPass
		rows[i].MissingFutureGatePass = missingPass
		rows[i].DominantToxicGatePass = dominantPass
		rows[i].PassesReviewGate = sourcePass && fullCountPass && splitCountPass && usablePass && toxicPass && missingPass && dominantPass && !rows[i].ClosedFamilyReslice && rows[i].RankableDecisionContext && rows[i].Split == fullSplitName
		rows[i].FailureReason = strings.Join(reasons, ";")
	}
}

func rangeContextTriageRankingRows(cohorts []FuturesRangeContextTriageCohortRow, cfg FuturesRangeContextTriageAuditConfig, splits []Split) []FuturesRangeContextTriageRankingRow {
	byIDSplit := map[string]map[string]FuturesRangeContextTriageCohortRow{}
	for _, row := range cohorts {
		if byIDSplit[row.CohortID] == nil {
			byIDSplit[row.CohortID] = map[string]FuturesRangeContextTriageCohortRow{}
		}
		byIDSplit[row.CohortID][row.Split] = row
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	sessionEdges := rangeContextTriageSessionEdgeByTimeframeHorizon(cohorts, cfg)
	rows := []FuturesRangeContextTriageRankingRow{}
	for cohortID, rowsBySplit := range byIDSplit {
		full, ok := rowsBySplit[fullSplitName]
		if !ok || !full.RankableDecisionContext {
			continue
		}
		rank := FuturesRangeContextTriageRankingRow{
			CohortID:                cohortID,
			CohortType:              full.CohortType,
			Timeframe:               full.Timeframe,
			HorizonBars:             full.HorizonBars,
			QualityBucket:           full.QualityBucket,
			MatureSession:           full.MatureSession,
			PassesGate:              full.PassesReviewGate,
			FullCandidateCount:      full.CandidateCount,
			FullUsableContextRate:   full.UsableContextRate,
			FullToxicContextRate:    full.ToxicContextRate,
			FullMissingFutureRate:   full.MissingFutureRate,
			PeriodSplitsRequired:    len(periodSplits),
			RankableDecisionContext: full.RankableDecisionContext,
			FailureReason:           full.FailureReason,
		}
		rank.WeakestSplitUsableRate = math.Inf(1)
		for _, split := range periodSplits {
			row := rowsBySplit[split.Name]
			if row.CandidateCount > 0 {
				rank.PeriodSplitsPassing++
			}
			if row.UsableContextRate < rank.WeakestSplitUsableRate {
				rank.WeakestSplitUsableRate = row.UsableContextRate
			}
			if row.ToxicContextRate > rank.WorstSplitToxicRate {
				rank.WorstSplitToxicRate = row.ToxicContextRate
			}
		}
		if math.IsInf(rank.WeakestSplitUsableRate, 1) {
			rank.WeakestSplitUsableRate = 0
		}
		rank.SessionEdgeRate = sessionEdges[rank.Timeframe+"|"+fmt.Sprint(rank.HorizonBars)]
		if full.MatureSession != "" {
			rank.SessionEdgeRate = math.Max(0, full.UsableContextRate-rangeContextTriageWorstSessionRate(cohorts, full.Timeframe, full.HorizonBars))
		}
		rank.RankScore = rangeContextTriageRankScore(rank)
		rows = append(rows, rank)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rangeContextTriageLessRanking(rows[i], rows[j])
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func rangeContextTriageSummaryRows(result FuturesRangeContextTriageAuditResult, cfg FuturesRangeContextTriageAuditConfig, splits []Split) []FuturesRangeContextTriageSummaryRow {
	sourcePass := rangeContextTriageSourceResamplePass(result.SourceRows, result.CoverageRows)
	rows := []FuturesRangeContextTriageSummaryRow{{
		Split:              fullSplitName,
		Timeframe:          "all",
		HorizonBars:        0,
		SourceResamplePass: sourcePass,
		EpisodeRows:        len(result.EpisodeRows),
		EligibleEpisodes:   rangeContextTriageEligibleEpisodeCount(result.EpisodeRows),
		FailureModeRows:    len(result.FailureModeRows),
		CohortRows:         len(result.CohortRows),
		RankingRows:        len(result.RankingRows),
		PassingCohorts:     result.PassingCohorts,
		StopState:          result.StopState,
	}}
	for _, split := range splits {
		for _, timeframe := range cfg.Timeframes {
			for _, horizon := range cfg.HorizonsBars {
				row := FuturesRangeContextTriageSummaryRow{
					Split:              split.Name,
					Timeframe:          timeframe,
					HorizonBars:        horizon,
					SourceResamplePass: sourcePass,
					StopState:          result.StopState,
				}
				for _, episode := range result.EpisodeRows {
					if episode.Timeframe == timeframe && rangeContextTriageRowInSplit(episode.MatureCloseTime, split) {
						row.EpisodeRows++
						if episode.Eligible {
							row.EligibleEpisodes++
						}
					}
				}
				for _, failure := range result.FailureModeRows {
					if failure.Timeframe == timeframe && failure.HorizonBars == horizon && rangeContextTriageRowInSplit(failure.MatureCloseTime, split) {
						row.FailureModeRows++
					}
				}
				for _, cohort := range result.CohortRows {
					if cohort.Split == split.Name && cohort.Timeframe == timeframe && cohort.HorizonBars == horizon {
						row.CohortRows++
						if cohort.PassesReviewGate {
							row.PassingCohorts++
						}
					}
				}
				for _, ranking := range result.RankingRows {
					if ranking.Timeframe == timeframe && ranking.HorizonBars == horizon {
						row.RankingRows++
					}
				}
				rows = append(rows, row)
			}
		}
	}
	return rows
}

func RangeContextTriageUTCSession(t time.Time) string {
	hour := t.UTC().Hour()
	switch {
	case hour <= 7:
		return RangeContextTriageSessionAsia
	case hour <= 12:
		return RangeContextTriageSessionEurope
	case hour <= 16:
		return RangeContextTriageSessionUSOverlap
	default:
		return RangeContextTriageSessionUSLate
	}
}

func rangeContextTriageQualityBucket(row FuturesRangeContextTriageEpisodeRow, cfg FuturesRangeContextTriageAuditConfig) string {
	switch {
	case row.WidthPct < cfg.LowWidthPctThreshold || row.WidthToATRRatio < cfg.LowWidthATRThreshold:
		return RangeContextTriageQualityTooNarrowNoise
	case row.WidthPct >= cfg.WideWidthPctThreshold || row.WidthToATRRatio >= cfg.WideWidthATRThreshold:
		return RangeContextTriageQualityWideVolatile
	case row.PreMatureMidCrossCount >= cfg.ChoppyMidCrossThreshold && row.PreMatureBoundaryTouchCount >= cfg.ChoppyBoundaryTouchThreshold:
		return RangeContextTriageQualityChoppy
	case row.WidthPct < 0.0030 && row.WidthToATRRatio < 1.50 && row.PreMatureCloseInsideRate >= 0.90:
		return RangeContextTriageQualityNarrowOrderly
	case row.PreMatureCloseInsideRate >= 0.90:
		return RangeContextTriageQualityBalancedOrderly
	default:
		return RangeContextTriageQualityUnknown
	}
}

func rangeContextTriageMidCrossCount(candles []Candle, start, end int, mid float64) int {
	count := 0
	prev := 0
	for i := start; i <= end && i < len(candles); i++ {
		side := 0
		if candles[i].Close > mid {
			side = 1
		} else if candles[i].Close < mid {
			side = -1
		}
		if prev != 0 && side != 0 && side != prev {
			count++
		}
		if side != 0 {
			prev = side
		}
	}
	return count
}

func rangeContextTriageBoundaryTouchCount(candles []Candle, start, end int, low, high float64) int {
	count := 0
	for i := start; i <= end && i < len(candles); i++ {
		if candles[i].Low <= low || candles[i].High >= high {
			count++
		}
	}
	return count
}

func rangeContextTriageInsideCloseRate(candles []Candle, start, end int, low, high float64) float64 {
	if end < start {
		return 0
	}
	count := 0
	total := 0
	for i := start; i <= end && i < len(candles); i++ {
		total++
		if candles[i].Close >= low && candles[i].Close <= high {
			count++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total)
}

func rangeContextTriageWickOvershootCount(candles []Candle, start, end int, low, high float64) int {
	count := 0
	for i := start; i <= end && i < len(candles); i++ {
		if candles[i].High > high || candles[i].Low < low {
			count++
		}
	}
	return count
}

func rangeContextTriageCohortKeyFor(split, timeframe string, horizon int, cohortType, quality, session, label string) rangeContextTriageCohortKey {
	parts := []string{"range_context", timeframe, fmt.Sprintf("h%d", horizon), cohortType}
	if quality != "" {
		parts = append(parts, quality)
	}
	if session != "" {
		parts = append(parts, session)
	}
	if label != "" {
		parts = append(parts, label)
	}
	return rangeContextTriageCohortKey{
		cohortID:            strings.Join(parts, "_"),
		cohortType:          cohortType,
		split:               split,
		timeframe:           timeframe,
		horizonBars:         horizon,
		qualityBucket:       quality,
		matureSession:       session,
		primaryContextLabel: label,
	}
}

func rangeContextTriageSourceResamplePass(sources []FuturesRangeContextTriageSourceRow, coverage []FuturesRangeContextTriageCoverageRow) bool {
	if len(sources) == 0 {
		return false
	}
	for _, row := range sources {
		if row.ValidationStatus != "accepted" || !row.SourceFactsPass {
			return false
		}
	}
	for _, row := range coverage {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass || !row.Complete {
			return false
		}
	}
	return true
}

func rangeContextTriageConstructiveLabel(label string) bool {
	return label == RangeContextTriageLabelContainedRotation || label == RangeContextTriageLabelCleanExpansionUp || label == RangeContextTriageLabelCleanExpansionDown
}

func rangeContextTriageToxicLabel(label string) bool {
	switch label {
	case RangeContextTriageLabelFalseBreakReentryUp, RangeContextTriageLabelFalseBreakReentryDown, RangeContextTriageLabelBoundaryChop, RangeContextTriageLabelDriftThroughUp, RangeContextTriageLabelDriftThroughDown, RangeContextTriageLabelLowWidthNoise, RangeContextTriageLabelNoResolution:
		return true
	default:
		return false
	}
}

func rangeContextTriageToxicLabels() []string {
	return []string{
		RangeContextTriageLabelFalseBreakReentryUp,
		RangeContextTriageLabelFalseBreakReentryDown,
		RangeContextTriageLabelBoundaryChop,
		RangeContextTriageLabelDriftThroughUp,
		RangeContextTriageLabelDriftThroughDown,
		RangeContextTriageLabelLowWidthNoise,
		RangeContextTriageLabelNoResolution,
	}
}

func rangeContextTriageDominantLabel(labels map[string]int, total int, toxicOnly bool) (string, float64) {
	if total == 0 {
		return "", 0
	}
	bestLabel := ""
	bestCount := 0
	for label, count := range labels {
		if toxicOnly && !rangeContextTriageToxicLabel(label) {
			continue
		}
		if count > bestCount || (count == bestCount && label < bestLabel) {
			bestLabel = label
			bestCount = count
		}
	}
	if bestCount == 0 {
		return "", 0
	}
	return bestLabel, float64(bestCount) / float64(total)
}

func rangeContextTriageSessionEdgeByTimeframeHorizon(cohorts []FuturesRangeContextTriageCohortRow, cfg FuturesRangeContextTriageAuditConfig) map[string]float64 {
	type minMax struct{ min, max float64 }
	values := map[string]minMax{}
	for _, row := range cohorts {
		if row.Split != fullSplitName || row.CohortType != RangeContextTriageCohortMatureSession || row.CandidateCount < cfg.MinSessionSplitCohortCount {
			continue
		}
		key := row.Timeframe + "|" + fmt.Sprint(row.HorizonBars)
		mm, ok := values[key]
		if !ok {
			values[key] = minMax{min: row.UsableContextRate, max: row.UsableContextRate}
			continue
		}
		mm.min = math.Min(mm.min, row.UsableContextRate)
		mm.max = math.Max(mm.max, row.UsableContextRate)
		values[key] = mm
	}
	out := map[string]float64{}
	for key, mm := range values {
		out[key] = math.Max(0, mm.max-mm.min)
	}
	return out
}

func rangeContextTriageWorstSessionRate(cohorts []FuturesRangeContextTriageCohortRow, timeframe string, horizon int) float64 {
	worst := math.Inf(1)
	for _, row := range cohorts {
		if row.Split == fullSplitName && row.CohortType == RangeContextTriageCohortMatureSession && row.Timeframe == timeframe && row.HorizonBars == horizon {
			worst = math.Min(worst, row.UsableContextRate)
		}
	}
	if math.IsInf(worst, 1) {
		return 0
	}
	return worst
}

func rangeContextTriageRankScore(row FuturesRangeContextTriageRankingRow) float64 {
	countBonus := 0.0
	if row.FullCandidateCount > 0 {
		countBonus = math.Min(0.20, math.Log10(float64(row.FullCandidateCount))/20)
	}
	return row.FullUsableContextRate + row.WeakestSplitUsableRate - row.FullToxicContextRate - row.WorstSplitToxicRate + math.Min(0.20, row.SessionEdgeRate) + countBonus
}

func rangeContextTriageLessRanking(a, b FuturesRangeContextTriageRankingRow) bool {
	if a.PassesGate != b.PassesGate {
		return a.PassesGate
	}
	if !rangeContextTriageFloatEqual(a.RankScore, b.RankScore) {
		return a.RankScore > b.RankScore
	}
	if !rangeContextTriageFloatEqual(a.WeakestSplitUsableRate, b.WeakestSplitUsableRate) {
		return a.WeakestSplitUsableRate > b.WeakestSplitUsableRate
	}
	if !rangeContextTriageFloatEqual(a.WorstSplitToxicRate, b.WorstSplitToxicRate) {
		return a.WorstSplitToxicRate < b.WorstSplitToxicRate
	}
	if a.FullCandidateCount != b.FullCandidateCount {
		return a.FullCandidateCount > b.FullCandidateCount
	}
	if rangeContextTriageTimeframeSortKey(a.Timeframe) != rangeContextTriageTimeframeSortKey(b.Timeframe) {
		return rangeContextTriageTimeframeSortKey(a.Timeframe) < rangeContextTriageTimeframeSortKey(b.Timeframe)
	}
	if rangeContextTriageCohortComplexity(a.CohortType) != rangeContextTriageCohortComplexity(b.CohortType) {
		return rangeContextTriageCohortComplexity(a.CohortType) < rangeContextTriageCohortComplexity(b.CohortType)
	}
	return a.CohortID < b.CohortID
}

func rangeContextTriageLessCohort(a, b FuturesRangeContextTriageCohortRow) bool {
	if a.Split != b.Split {
		return rangeContextTriageSplitSortKey(a.Split) < rangeContextTriageSplitSortKey(b.Split)
	}
	if a.Timeframe != b.Timeframe {
		return rangeContextTriageTimeframeSortKey(a.Timeframe) < rangeContextTriageTimeframeSortKey(b.Timeframe)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if rangeContextTriageCohortComplexity(a.CohortType) != rangeContextTriageCohortComplexity(b.CohortType) {
		return rangeContextTriageCohortComplexity(a.CohortType) < rangeContextTriageCohortComplexity(b.CohortType)
	}
	return a.CohortID < b.CohortID
}

func rangeContextTriageCohortComplexity(cohortType string) int {
	switch cohortType {
	case RangeContextTriageCohortAll:
		return 0
	case RangeContextTriageCohortQualityBucket:
		return 1
	case RangeContextTriageCohortMatureSession:
		return 2
	case RangeContextTriageCohortQualityBucketMatureSession:
		return 3
	case RangeContextTriageCohortPrimaryContextLabel:
		return 4
	default:
		return 9
	}
}

func rangeContextTriageTimeframeSortKey(timeframe string) int {
	switch timeframe {
	case RangeDiscoveryTimeframe1h:
		return 0
	case RangeDiscoveryTimeframe15m:
		return 1
	case RangeDiscoveryTimeframe4h:
		return 2
	default:
		return 9
	}
}

func rangeContextTriageSplitSortKey(split string) int {
	switch split {
	case "2021_2022_stress":
		return 0
	case "2023_2024_oos":
		return 1
	case "2025_2026_recent":
		return 2
	case fullSplitName:
		return 3
	default:
		return 9
	}
}

func rangeContextTriageEligibleEpisodeCount(rows []FuturesRangeContextTriageEpisodeRow) int {
	count := 0
	for _, row := range rows {
		if row.Eligible {
			count++
		}
	}
	return count
}

func rangeContextTriagePassingRankingCount(rows []FuturesRangeContextTriageRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func rangeContextTriageFloatEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func rangeContextTriageRowInSplit(rawTime string, split Split) bool {
	t, err := parseTime(rawTime)
	if err != nil {
		return false
	}
	return split.Contains(t)
}

func appendReason(existing string, reason string) string {
	if existing == "" {
		return reason
	}
	return existing + ";" + reason
}
