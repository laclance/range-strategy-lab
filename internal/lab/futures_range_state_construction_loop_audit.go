package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	FuturesRangeStateConstructionLoopAuditName = "futures_range_state_construction_loop_audit"

	RangeStateConstructionLoopStopStateSourceGap               = "range_state_construction_loop_source_gap"
	RangeStateConstructionLoopStopStateNoEligibleStates        = "range_state_construction_loop_no_eligible_states"
	RangeStateConstructionLoopStopStateFailedNoUsableState     = "range_state_construction_loop_audit_failed_no_usable_state"
	RangeStateConstructionLoopStopStatePassedNoTradeFilterOnly = "range_state_construction_loop_audit_passed_no_trade_filter_only"
	RangeStateConstructionLoopStopStatePassedNeedsRouterSpec   = "range_state_construction_loop_audit_passed_needs_router_spec"
	RangeStateConstructionLoopStopStatePassedNeedsStrategySpec = "range_state_construction_loop_audit_passed_needs_strategy_premise_spec"
	RangeStateConstructionLoopStopStateRejectedClosedReslice   = "range_state_construction_loop_rejected_closed_family_reslice"
	RangeStateConstructionLoopRouteRotation                    = "tradable_rotation_candidate"
	RangeStateConstructionLoopRouteContinuation                = "trend_continuation_candidate"
	RangeStateConstructionLoopRouteNoTradeToxic                = "no_trade_toxic"
	RangeStateConstructionLoopRouteDiagnosticOnly              = "diagnostic_only"
	RangeStateConstructionLoopRollupGeometryVol                = "geometry+vol"
	RangeStateConstructionLoopRollupGeometryTrend              = "geometry+trend"
	RangeStateConstructionLoopRollupGeometryImpulse            = "geometry+impulse"
	RangeStateConstructionLoopRollupGeometryParticipation      = "geometry+participation"
	RangeStateConstructionLoopRollupGeometryVolTrend           = "geometry+vol+trend"
	RangeStateConstructionLoopRollupGeometryVolTrendImpulse    = "geometry+vol+trend+impulse"
	RangeStateConstructionLoopRollupAllFamilies                = "all_families"
	rangeStateConstructionLoopSourcePath                       = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"
	rangeStateConstructionLoopExpectedRows                     = 573984
	rangeStateConstructionLoopExpectedFirst                    = "2021-01-01T00:00:00Z"
	rangeStateConstructionLoopExpectedLast                     = "2026-06-16T23:55:00Z"
	rangeStateConstructionLoopExpectedZeroVol                  = 66
	rangeStateConstructionLoopExpected15MRows                  = 191328
	rangeStateConstructionLoopExpected15MLast                  = "2026-06-16T23:45:00Z"
	rangeStateConstructionLoopExpected1HRows                   = 47832
	rangeStateConstructionLoopExpected1HLast                   = "2026-06-16T23:00:00Z"
	rangeStateConstructionLoopExpected4HRows                   = 11958
	rangeStateConstructionLoopExpected4HLast                   = "2026-06-16T20:00:00Z"
	rangeStateConstructionLoopDetectorProfileID                = "p30_c12_bollinger_on_adx_off"
)

type FuturesRangeStateConstructionLoopAuditConfig struct {
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
	ShortWindowBars              int
	MediumWindowBars             int
	FeatureLookbackBars          int
	LongLookbackDays             int
	ATRPeriod                    int
	ADXPeriod                    int
	Horizons15M                  []int
	Horizons1H                   []int
	Horizons4H                   []int
	ReentryWindowBars            int
	CleanExpansionThreshold      float64
	DriftThreshold               float64
	BoundaryChopTransitions      int
	LowWidthPctThreshold         float64
	LowWidthATRThreshold         float64
	WideWidthPctThreshold        float64
	WideWidthATRThreshold        float64
	AbnormalTRATRThreshold       float64
	LargeBodyToRangeThreshold    float64
	MinFullCohortCount           int
	MinSplitCohortCount          int
	MaxSplitContributionRate     float64
	MaxDominantRouteLabelRate    float64
	MinPositiveUsefulRateFull    float64
	MinPositiveUsefulRateSplit   float64
	MaxPositiveToxicRateFull     float64
	MaxPositiveToxicRateSplit    float64
	MinPositiveMarginFull        float64
	MinPositiveMarginSplit       float64
	MinToxicRateFull             float64
	MinToxicRateSplit            float64
	MinNoTradeFullCohortCount    int
	MinNoTradeSplitCohortCount   int
}

type FuturesRangeStateConstructionLoopAuditResult struct {
	SourceRows        []FuturesRangeStateConstructionLoopSourceRow        `json:"source_rows"`
	CoverageRows      []FuturesRangeStateConstructionLoopCoverageRow      `json:"coverage_rows"`
	FeatureWindowRows []FuturesRangeStateConstructionLoopFeatureWindowRow `json:"feature_window_rows"`
	StateRows         []FuturesRangeStateConstructionLoopStateRow         `json:"state_rows"`
	LabelRows         []FuturesRangeStateConstructionLoopLabelRow         `json:"label_rows"`
	CohortRows        []FuturesRangeStateConstructionLoopCohortRow        `json:"cohort_rows"`
	RankingRows       []FuturesRangeStateConstructionLoopRankingRow       `json:"ranking_rows"`
	SummaryRows       []FuturesRangeStateConstructionLoopSummaryRow       `json:"summary_rows"`
	SkipRows          []FuturesRangeStateConstructionLoopSkipRow          `json:"skip_rows"`
	PassingCohorts    int                                                 `json:"passing_cohorts"`
	StopState         string                                              `json:"stop_state"`
}

type FuturesRangeStateConstructionLoopSourceRow struct {
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

type FuturesRangeStateConstructionLoopCoverageRow struct {
	ExpectedRowCount      int    `json:"expected_row_count"`
	ExpectedFirstOpenTime string `json:"expected_first_open_time"`
	ExpectedLastOpenTime  string `json:"expected_last_open_time"`
	CoverageFactsPass     bool   `json:"coverage_facts_pass"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesRangeStateConstructionLoopFeatureWindowRow struct {
	Timeframe                  string `json:"timeframe"`
	BarsPerDay                 int    `json:"bars_per_day"`
	DetectorProfileID          string `json:"detector_profile_id"`
	DetectorLookbackDays       int    `json:"detector_lookback_days"`
	DetectorLookbackBars       int    `json:"detector_lookback_bars"`
	DetectorMinConsecutiveBars int    `json:"detector_min_consecutive_bars"`
	ShortWindowBars            int    `json:"short_window_bars"`
	MediumWindowBars           int    `json:"medium_window_bars"`
	FeatureLookbackBars        int    `json:"feature_lookback_bars"`
	LongLookbackBars           int    `json:"long_lookback_bars"`
	ATRPeriod                  int    `json:"atr_period"`
	ADXPeriod                  int    `json:"adx_period"`
	HorizonsBars               string `json:"horizons_bars"`
	MaxHorizonBars             int    `json:"max_horizon_bars"`
	ClosedCandleOnly           bool   `json:"closed_candle_only"`
	ForwardLabelsAsFeatures    bool   `json:"forward_labels_as_features"`
}

type FuturesRangeStateConstructionLoopStateRow struct {
	StateRowID                       int     `json:"state_row_id"`
	Timestamp                        string  `json:"timestamp"`
	Timeframe                        string  `json:"timeframe"`
	Split                            string  `json:"split"`
	RangeEpisodeID                   int     `json:"range_episode_id"`
	RangeStartIndex                  int     `json:"range_start_index"`
	DecisionIndex                    int     `json:"decision_index"`
	RangeStartTime                   string  `json:"range_start_time"`
	DecisionCloseTime                string  `json:"decision_close_time"`
	RangeHigh                        float64 `json:"range_high"`
	RangeLow                         float64 `json:"range_low"`
	RangeMid                         float64 `json:"range_mid"`
	RangeWidth                       float64 `json:"range_width"`
	RangeAgeBars                     int     `json:"range_age_bars"`
	RangeWidthPct                    float64 `json:"range_width_pct"`
	RangeWidthATRRatio               float64 `json:"range_width_atr_ratio"`
	CloseLocationPct                 float64 `json:"close_location_pct"`
	DistanceToLowPct                 float64 `json:"distance_to_low_pct"`
	DistanceToHighPct                float64 `json:"distance_to_high_pct"`
	DistanceToMidPct                 float64 `json:"distance_to_mid_pct"`
	BoundaryTouchCountLookback       int     `json:"boundary_touch_count_lookback"`
	MidlineCrossCountLookback        int     `json:"midline_cross_count_lookback"`
	CloseLocationEntropyLookback     float64 `json:"close_location_entropy_lookback"`
	RangeDutyCycleLookback           float64 `json:"range_duty_cycle_lookback"`
	ATRShort                         float64 `json:"atr_short"`
	ATRPercentileLong                float64 `json:"atr_percentile_long"`
	RealizedVolPercentileLong        float64 `json:"realized_vol_percentile_long"`
	ATRToRangeWidth                  float64 `json:"atr_to_range_width"`
	TrueRangeExpansionRatio          float64 `json:"true_range_expansion_ratio"`
	AbnormalRangeBarCountLookback    int     `json:"abnormal_range_bar_count_lookback"`
	ReturnShortPct                   float64 `json:"return_short_pct"`
	ReturnMediumPct                  float64 `json:"return_medium_pct"`
	MovingAverageSlopePct            float64 `json:"moving_average_slope_pct"`
	HigherTimeframeDirectionProxyPct float64 `json:"higher_timeframe_direction_proxy_pct"`
	ADX                              float64 `json:"adx"`
	DistanceFromMAPct                float64 `json:"distance_from_ma_pct"`
	LastAbnormalCandleSide           string  `json:"last_abnormal_candle_side"`
	BarsSinceLastAbnormalCandle      int     `json:"bars_since_last_abnormal_candle"`
	LargeBodyCandleCountLookback     int     `json:"large_body_candle_count_lookback"`
	LargeRangeCandleCountLookback    int     `json:"large_range_candle_count_lookback"`
	ImpulseContinuationPressure      float64 `json:"impulse_continuation_pressure"`
	ImpulseExhaustionProxy           float64 `json:"impulse_exhaustion_proxy"`
	VolumePercentileLong             float64 `json:"volume_percentile_long"`
	VolumeChangeRatio                float64 `json:"volume_change_ratio"`
	VolumePerRangeWidthProxy         float64 `json:"volume_per_range_width_proxy"`
	CandleSpreadProxy                float64 `json:"candle_spread_proxy"`
	WickBodyStructureSummary         float64 `json:"wick_body_structure_summary"`
	ZeroVolumeRowFlag                bool    `json:"zero_volume_row_flag"`
	GeometryBucket                   string  `json:"geometry_bucket"`
	VolBucket                        string  `json:"vol_bucket"`
	TrendBucket                      string  `json:"trend_bucket"`
	ImpulseBucket                    string  `json:"impulse_bucket"`
	ParticipationBucket              string  `json:"participation_bucket"`
	StateID                          string  `json:"state_id"`
	GeometryVolID                    string  `json:"geometry_vol_id"`
	GeometryTrendID                  string  `json:"geometry_trend_id"`
	GeometryImpulseID                string  `json:"geometry_impulse_id"`
	GeometryParticipationID          string  `json:"geometry_participation_id"`
	GeometryVolTrendID               string  `json:"geometry_vol_trend_id"`
	GeometryVolTrendImpulseID        string  `json:"geometry_vol_trend_impulse_id"`
	AllFamiliesID                    string  `json:"all_families_id"`
	Eligible                         bool    `json:"eligible"`
	SkippedReason                    string  `json:"skipped_reason,omitempty"`
}

type FuturesRangeStateConstructionLoopLabelRow struct {
	StateRowID                   int     `json:"state_row_id"`
	StateID                      string  `json:"state_id"`
	Timestamp                    string  `json:"timestamp"`
	Timeframe                    string  `json:"timeframe"`
	Split                        string  `json:"split"`
	RangeEpisodeID               int     `json:"range_episode_id"`
	HorizonBars                  int     `json:"horizon_bars"`
	ForwardLabel                 string  `json:"forward_label"`
	RotationUseful               bool    `json:"rotation_useful"`
	RotationToxic                bool    `json:"rotation_toxic"`
	ContinuationUseful           bool    `json:"continuation_useful"`
	ContinuationToxic            bool    `json:"continuation_toxic"`
	NoTradeToxic                 bool    `json:"no_trade_toxic"`
	DiagnosticOnly               bool    `json:"diagnostic_only"`
	LabelWindowStartIndex        int     `json:"label_window_start_index"`
	LabelWindowEndIndex          int     `json:"label_window_end_index"`
	LabelWindowStartTime         string  `json:"label_window_start_time"`
	LabelWindowEndTime           string  `json:"label_window_end_time"`
	FirstOutsideCloseSide        string  `json:"first_outside_close_side"`
	BarsToFirstOutsideClose      int     `json:"bars_to_first_outside_close"`
	ReentryWithinWindow          bool    `json:"reentry_within_window"`
	MaxExcursionAboveRangeWidths float64 `json:"max_excursion_above_range_widths"`
	MaxExcursionBelowRangeWidths float64 `json:"max_excursion_below_range_widths"`
	FinalClosePosition           float64 `json:"final_close_position"`
	InsideCloseRate              float64 `json:"inside_close_rate"`
	MidpointCrossCount           int     `json:"midpoint_cross_count"`
	OutsideStateTransitionCount  int     `json:"outside_state_transition_count"`
	GeometryVolID                string  `json:"geometry_vol_id"`
	GeometryTrendID              string  `json:"geometry_trend_id"`
	GeometryImpulseID            string  `json:"geometry_impulse_id"`
	GeometryParticipationID      string  `json:"geometry_participation_id"`
	GeometryVolTrendID           string  `json:"geometry_vol_trend_id"`
	GeometryVolTrendImpulseID    string  `json:"geometry_vol_trend_impulse_id"`
	AllFamiliesID                string  `json:"all_families_id"`
	FutureWindowMetadataOnly     bool    `json:"future_window_metadata_only"`
	ForwardLabelUsedAsFeature    bool    `json:"forward_label_used_as_feature"`
}

type FuturesRangeStateConstructionLoopCohortRow struct {
	CohortID                   string  `json:"cohort_id"`
	RouteCandidate             string  `json:"route_candidate"`
	RollupType                 string  `json:"rollup_type"`
	RollupID                   string  `json:"rollup_id"`
	Split                      string  `json:"split"`
	Timeframe                  string  `json:"timeframe"`
	HorizonBars                int     `json:"horizon_bars"`
	CandidateCount             int     `json:"candidate_count"`
	UsefulCount                int     `json:"useful_count"`
	ToxicCount                 int     `json:"toxic_count"`
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
	UsefulRate                 float64 `json:"useful_rate"`
	ToxicRate                  float64 `json:"toxic_rate"`
	UsefulMinusToxicMargin     float64 `json:"useful_minus_toxic_margin"`
	DominantForwardLabel       string  `json:"dominant_forward_label"`
	DominantForwardLabelRate   float64 `json:"dominant_forward_label_rate"`
	FullPeriodRows             int     `json:"full_period_rows"`
	WeakestSplitRows           int     `json:"weakest_split_rows"`
	MaxSplitContributionRate   float64 `json:"max_split_contribution_rate"`
	FullUsefulRate             float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate     float64 `json:"weakest_split_useful_rate"`
	FullToxicRate              float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate        float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin         float64 `json:"weakest_split_margin"`
	ReviewableCountGatePass    bool    `json:"reviewable_count_gate_pass"`
	SplitStabilityGatePass     bool    `json:"split_stability_gate_pass"`
	SplitContributionGatePass  bool    `json:"split_contribution_gate_pass"`
	DominantLabelGatePass      bool    `json:"dominant_label_gate_pass"`
	FeatureKnownGatePass       bool    `json:"feature_known_gate_pass"`
	RouteRateGatePass          bool    `json:"route_rate_gate_pass"`
	ClosedFamilyReslice        bool    `json:"closed_family_reslice"`
	ClosedFamilyProtectionPass bool    `json:"closed_family_protection_pass"`
	FutureLeakProtectionPass   bool    `json:"future_leak_protection_pass"`
	PassesReviewGate           bool    `json:"passes_review_gate"`
	FailureReason              string  `json:"failure_reason,omitempty"`
}

type FuturesRangeStateConstructionLoopRankingRow struct {
	Rank                       int     `json:"rank"`
	CohortID                   string  `json:"cohort_id"`
	RouteCandidate             string  `json:"route_candidate"`
	RollupType                 string  `json:"rollup_type"`
	RollupID                   string  `json:"rollup_id"`
	Timeframe                  string  `json:"timeframe"`
	HorizonBars                int     `json:"horizon_bars"`
	PassesGate                 bool    `json:"passes_gate"`
	RankScore                  float64 `json:"rank_score"`
	FullPeriodRows             int     `json:"full_period_rows"`
	WeakestSplitRows           int     `json:"weakest_split_rows"`
	MaxSplitContributionRate   float64 `json:"max_split_contribution_rate"`
	FullUsefulRate             float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate     float64 `json:"weakest_split_useful_rate"`
	FullToxicRate              float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate        float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin         float64 `json:"weakest_split_margin"`
	DominantForwardLabel       string  `json:"dominant_forward_label"`
	DominantForwardLabelRate   float64 `json:"dominant_forward_label_rate"`
	ClosedFamilyReslice        bool    `json:"closed_family_reslice"`
	FutureLeakProtectionPass   bool    `json:"future_leak_protection_pass"`
	FailureReason              string  `json:"failure_reason,omitempty"`
}

type FuturesRangeStateConstructionLoopSummaryRow struct {
	Split              string `json:"split"`
	Timeframe          string `json:"timeframe"`
	HorizonBars        int    `json:"horizon_bars"`
	SourceResamplePass bool   `json:"source_resample_pass"`
	StateRows          int    `json:"state_rows"`
	LabelRows          int    `json:"label_rows"`
	CohortRows         int    `json:"cohort_rows"`
	RankingRows        int    `json:"ranking_rows"`
	SkipRows           int    `json:"skip_rows"`
	PassingCohorts     int    `json:"passing_cohorts"`
	StopState          string `json:"stop_state"`
}

type FuturesRangeStateConstructionLoopSkipRow struct {
	Timeframe string `json:"timeframe"`
	Split     string `json:"split"`
	Reason    string `json:"reason"`
	Count     int    `json:"count"`
}

type rangeStateFrameData struct {
	frame      rangeDiscoveryFrameDef
	candles    []Candle
	coverage   FuturesRangeStateConstructionLoopCoverageRow
	metrics    rangeStateMetrics
	horizons   []int
	maxHorizon int
}

type rangeStateMetrics struct {
	normalizedATR         []float64
	atr                   []float64
	adx                   []float64
	realizedVol           []float64
	atrPercentile         []float64
	realizedVolPercentile []float64
	volumePercentile      []float64
	trueRange             []float64
}

type rangeStateSkipAccumulator map[rangeStateSkipKey]int

type rangeStateSkipKey struct {
	timeframe string
	split     string
	reason    string
}

type rangeStateCohortKey struct {
	cohortID       string
	routeCandidate string
	rollupType     string
	rollupID       string
	split          string
	timeframe      string
	horizonBars    int
}

type rangeStateCohortAccumulator struct {
	row    FuturesRangeStateConstructionLoopCohortRow
	labels map[string]int
}

func DefaultFuturesRangeStateConstructionLoopAuditConfig() FuturesRangeStateConstructionLoopAuditConfig {
	return FuturesRangeStateConstructionLoopAuditConfig{
		SourcePath:                 rangeStateConstructionLoopSourcePath,
		ApprovedSourcePath:         rangeStateConstructionLoopSourcePath,
		ExpectedSourceRows:         rangeStateConstructionLoopExpectedRows,
		ExpectedFirstOpenTime:      rangeStateConstructionLoopExpectedFirst,
		ExpectedLastOpenTime:       rangeStateConstructionLoopExpectedLast,
		ExpectedGapCount:           0,
		ExpectedDuplicateCount:     0,
		ExpectedZeroVolumeCount:    rangeStateConstructionLoopExpectedZeroVol,
		Timeframes:                 []string{RangeDiscoveryTimeframe15m, RangeDiscoveryTimeframe1h, RangeDiscoveryTimeframe4h},
		Expected15MRows:            rangeStateConstructionLoopExpected15MRows,
		Expected15MLastOpenTime:    rangeStateConstructionLoopExpected15MLast,
		Expected1HRows:             rangeStateConstructionLoopExpected1HRows,
		Expected1HLastOpenTime:     rangeStateConstructionLoopExpected1HLast,
		Expected4HRows:             rangeStateConstructionLoopExpected4HRows,
		Expected4HLastOpenTime:     rangeStateConstructionLoopExpected4HLast,
		DetectorLookbackDays:       20,
		DetectorPercentile:         0.30,
		DetectorMinConsecutiveBars: 12,
		ShortWindowBars:            12,
		MediumWindowBars:           48,
		FeatureLookbackBars:        48,
		LongLookbackDays:           20,
		ATRPeriod:                  14,
		ADXPeriod:                  14,
		Horizons15M:                []int{12, 24, 48},
		Horizons1H:                 []int{12, 24, 48},
		Horizons4H:                 []int{6, 12, 24},
		ReentryWindowBars:          6,
		CleanExpansionThreshold:    0.75,
		DriftThreshold:             0.50,
		BoundaryChopTransitions:    3,
		LowWidthPctThreshold:       0.0015,
		LowWidthATRThreshold:       0.75,
		WideWidthPctThreshold:      0.0150,
		WideWidthATRThreshold:      4.00,
		AbnormalTRATRThreshold:     1.80,
		LargeBodyToRangeThreshold:  0.60,
		MinFullCohortCount:         150,
		MinSplitCohortCount:        30,
		MaxSplitContributionRate:   0.60,
		MaxDominantRouteLabelRate:  0.80,
		MinPositiveUsefulRateFull:  0.58,
		MinPositiveUsefulRateSplit: 0.52,
		MaxPositiveToxicRateFull:   0.42,
		MaxPositiveToxicRateSplit:  0.48,
		MinPositiveMarginFull:      0.12,
		MinPositiveMarginSplit:     0.04,
		MinToxicRateFull:           0.58,
		MinToxicRateSplit:          0.52,
		MinNoTradeFullCohortCount:  200,
		MinNoTradeSplitCohortCount: 40,
	}
}

func RunFuturesRangeStateConstructionLoopAudit(candles []Candle, manifest SourceManifest, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split) (FuturesRangeStateConstructionLoopAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeStateConstructionLoopAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesRangeStateConstructionLoopAuditResult{}
	sourceRow := rangeStateConstructionLoopSourceRow(candles, manifest, cfg)
	result.SourceRows = append(result.SourceRows, sourceRow)
	if sourceRow.ValidationStatus != "accepted" || !sourceRow.SourceFactsPass {
		result.StopState = RangeStateConstructionLoopStopStateSourceGap
		result.SummaryRows = rangeStateConstructionLoopSummaryRows(result, cfg, splits)
		return result, nil
	}

	frameDataByTimeframe := map[string]rangeStateFrameData{}
	for _, timeframe := range cfg.Timeframes {
		frame, ok := rangeContextTriageFrameDef(timeframe)
		if !ok {
			return result, fmt.Errorf("range state construction loop missing frame definition for %s", timeframe)
		}
		frameCandles, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
		coverageRow := rangeStateConstructionLoopCoverageRow(coverage, cfg)
		result.CoverageRows = append(result.CoverageRows, coverageRow)
		if err != nil || !coverageRow.CoverageFactsPass || coverageRow.ValidationStatus != "accepted" || !coverageRow.Complete {
			result.StopState = RangeStateConstructionLoopStopStateSourceGap
			result.SummaryRows = rangeStateConstructionLoopSummaryRows(result, cfg, splits)
			return result, nil
		}
		horizons := rangeStateConstructionLoopHorizons(timeframe, cfg)
		data := rangeStateFrameData{
			frame:      frame,
			candles:    frameCandles,
			coverage:   coverageRow,
			metrics:    rangeStateConstructionLoopMetrics(frameCandles, frame, cfg),
			horizons:   horizons,
			maxHorizon: maxIntInSlice(horizons),
		}
		frameDataByTimeframe[timeframe] = data
		result.FeatureWindowRows = append(result.FeatureWindowRows, rangeStateConstructionLoopFeatureWindowRow(frame, cfg, horizons))
	}

	skips := rangeStateSkipAccumulator{}
	for _, timeframe := range cfg.Timeframes {
		data := frameDataByTimeframe[timeframe]
		stateRows, err := rangeStateConstructionLoopStateRows(data, frameDataByTimeframe, cfg, splits, skips, len(result.StateRows))
		if err != nil {
			return result, err
		}
		result.StateRows = append(result.StateRows, stateRows...)
	}
	result.SkipRows = rangeStateConstructionLoopSkipRows(skips)
	result.LabelRows = rangeStateConstructionLoopLabelRows(result.StateRows, frameDataByTimeframe, cfg)
	result.CohortRows = rangeStateConstructionLoopCohortRows(result.LabelRows, result.SourceRows, result.CoverageRows, cfg, splits)
	result.RankingRows = rangeStateConstructionLoopRankingRows(result.CohortRows, splits)
	result.PassingCohorts = rangeStateConstructionLoopPassingRankingCount(result.RankingRows)
	result.StopState = FuturesRangeStateConstructionLoopAuditStopState(result)
	result.SummaryRows = rangeStateConstructionLoopSummaryRows(result, cfg, splits)
	return result, nil
}

func FuturesRangeStateConstructionLoopAuditStopState(result FuturesRangeStateConstructionLoopAuditResult) string {
	if result.StopState == RangeStateConstructionLoopStopStateSourceGap || result.StopState == RangeStateConstructionLoopStopStateRejectedClosedReslice {
		return result.StopState
	}
	if !rangeStateConstructionLoopSourceResamplePass(result.SourceRows, result.CoverageRows) {
		return RangeStateConstructionLoopStopStateSourceGap
	}
	for _, row := range result.StateRows {
		if row.SkippedReason == "closed_family_reslice" {
			return RangeStateConstructionLoopStopStateRejectedClosedReslice
		}
	}
	if len(result.StateRows) == 0 {
		return RangeStateConstructionLoopStopStateNoEligibleStates
	}
	positivePasses := 0
	toxicPasses := 0
	for _, row := range result.RankingRows {
		if !row.PassesGate {
			continue
		}
		switch row.RouteCandidate {
		case RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation:
			positivePasses++
		case RangeStateConstructionLoopRouteNoTradeToxic:
			toxicPasses++
		}
	}
	switch {
	case positivePasses == 0 && toxicPasses == 0:
		return RangeStateConstructionLoopStopStateFailedNoUsableState
	case positivePasses == 0 && toxicPasses > 0:
		return RangeStateConstructionLoopStopStatePassedNoTradeFilterOnly
	case positivePasses == 1 && toxicPasses == 0:
		return RangeStateConstructionLoopStopStatePassedNeedsStrategySpec
	default:
		return RangeStateConstructionLoopStopStatePassedNeedsRouterSpec
	}
}

func (cfg FuturesRangeStateConstructionLoopAuditConfig) withDefaults() FuturesRangeStateConstructionLoopAuditConfig {
	defaults := DefaultFuturesRangeStateConstructionLoopAuditConfig()
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
	if cfg.ShortWindowBars == 0 {
		cfg.ShortWindowBars = defaults.ShortWindowBars
	}
	if cfg.MediumWindowBars == 0 {
		cfg.MediumWindowBars = defaults.MediumWindowBars
	}
	if cfg.FeatureLookbackBars == 0 {
		cfg.FeatureLookbackBars = defaults.FeatureLookbackBars
	}
	if cfg.LongLookbackDays == 0 {
		cfg.LongLookbackDays = defaults.LongLookbackDays
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.ADXPeriod == 0 {
		cfg.ADXPeriod = defaults.ADXPeriod
	}
	if len(cfg.Horizons15M) == 0 {
		cfg.Horizons15M = append([]int(nil), defaults.Horizons15M...)
	}
	if len(cfg.Horizons1H) == 0 {
		cfg.Horizons1H = append([]int(nil), defaults.Horizons1H...)
	}
	if len(cfg.Horizons4H) == 0 {
		cfg.Horizons4H = append([]int(nil), defaults.Horizons4H...)
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
	if cfg.AbnormalTRATRThreshold == 0 {
		cfg.AbnormalTRATRThreshold = defaults.AbnormalTRATRThreshold
	}
	if cfg.LargeBodyToRangeThreshold == 0 {
		cfg.LargeBodyToRangeThreshold = defaults.LargeBodyToRangeThreshold
	}
	if cfg.MinFullCohortCount == 0 {
		cfg.MinFullCohortCount = defaults.MinFullCohortCount
	}
	if cfg.MinSplitCohortCount == 0 {
		cfg.MinSplitCohortCount = defaults.MinSplitCohortCount
	}
	if cfg.MaxSplitContributionRate == 0 {
		cfg.MaxSplitContributionRate = defaults.MaxSplitContributionRate
	}
	if cfg.MaxDominantRouteLabelRate == 0 {
		cfg.MaxDominantRouteLabelRate = defaults.MaxDominantRouteLabelRate
	}
	if cfg.MinPositiveUsefulRateFull == 0 {
		cfg.MinPositiveUsefulRateFull = defaults.MinPositiveUsefulRateFull
	}
	if cfg.MinPositiveUsefulRateSplit == 0 {
		cfg.MinPositiveUsefulRateSplit = defaults.MinPositiveUsefulRateSplit
	}
	if cfg.MaxPositiveToxicRateFull == 0 {
		cfg.MaxPositiveToxicRateFull = defaults.MaxPositiveToxicRateFull
	}
	if cfg.MaxPositiveToxicRateSplit == 0 {
		cfg.MaxPositiveToxicRateSplit = defaults.MaxPositiveToxicRateSplit
	}
	if cfg.MinPositiveMarginFull == 0 {
		cfg.MinPositiveMarginFull = defaults.MinPositiveMarginFull
	}
	if cfg.MinPositiveMarginSplit == 0 {
		cfg.MinPositiveMarginSplit = defaults.MinPositiveMarginSplit
	}
	if cfg.MinToxicRateFull == 0 {
		cfg.MinToxicRateFull = defaults.MinToxicRateFull
	}
	if cfg.MinToxicRateSplit == 0 {
		cfg.MinToxicRateSplit = defaults.MinToxicRateSplit
	}
	if cfg.MinNoTradeFullCohortCount == 0 {
		cfg.MinNoTradeFullCohortCount = defaults.MinNoTradeFullCohortCount
	}
	if cfg.MinNoTradeSplitCohortCount == 0 {
		cfg.MinNoTradeSplitCohortCount = defaults.MinNoTradeSplitCohortCount
	}
	return cfg
}

func (cfg FuturesRangeStateConstructionLoopAuditConfig) validate() error {
	if cfg.SourcePath == "" || cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("range state construction loop source paths must not be empty")
	}
	for _, timeframe := range cfg.Timeframes {
		if timeframe != RangeDiscoveryTimeframe15m && timeframe != RangeDiscoveryTimeframe1h && timeframe != RangeDiscoveryTimeframe4h {
			return fmt.Errorf("range state construction loop unsupported timeframe %q", timeframe)
		}
		if len(rangeStateConstructionLoopHorizons(timeframe, cfg)) == 0 {
			return fmt.Errorf("range state construction loop %s horizons must not be empty", timeframe)
		}
	}
	if cfg.DetectorLookbackDays <= 0 && cfg.DetectorLookbackBarsOverride <= 0 {
		return fmt.Errorf("range state construction loop detector lookback must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("range state construction loop detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 || cfg.ShortWindowBars <= 0 || cfg.MediumWindowBars <= 0 || cfg.FeatureLookbackBars <= 0 || cfg.LongLookbackDays <= 0 || cfg.ATRPeriod <= 0 || cfg.ADXPeriod <= 0 {
		return fmt.Errorf("range state construction loop windows must be positive")
	}
	if cfg.MediumWindowBars < cfg.ShortWindowBars {
		return fmt.Errorf("range state construction loop medium window must be >= short window")
	}
	for _, horizon := range append(append(append([]int{}, cfg.Horizons15M...), cfg.Horizons1H...), cfg.Horizons4H...) {
		if horizon <= 0 {
			return fmt.Errorf("range state construction loop horizon bars must be positive")
		}
	}
	if cfg.ReentryWindowBars <= 0 || cfg.BoundaryChopTransitions <= 0 || cfg.CleanExpansionThreshold <= 0 || cfg.DriftThreshold <= 0 {
		return fmt.Errorf("range state construction loop label thresholds must be positive")
	}
	if cfg.MinFullCohortCount <= 0 || cfg.MinSplitCohortCount <= 0 || cfg.MinNoTradeFullCohortCount <= 0 || cfg.MinNoTradeSplitCohortCount <= 0 {
		return fmt.Errorf("range state construction loop count gates must be positive")
	}
	return nil
}

func rangeStateConstructionLoopSourceRow(parent []Candle, manifest SourceManifest, cfg FuturesRangeStateConstructionLoopAuditConfig) FuturesRangeStateConstructionLoopSourceRow {
	row := FuturesRangeStateConstructionLoopSourceRow{
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
		reject(fmt.Sprintf("range state construction loop requires BTCUSDT Binance USDT-M futures 5m comparison_only=false; got product=%q symbol=%q interval=%q comparison_only=%t", manifest.Product, manifest.Symbol, manifest.Interval, manifest.ComparisonOnly))
		return row
	}
	if !sameCleanPath(manifest.Path, cfg.ApprovedSourcePath) {
		reject(fmt.Sprintf("range state construction loop source path %q is not approved path %q", manifest.Path, cfg.ApprovedSourcePath))
		return row
	}
	if cfg.SkipSourceFactCheck {
		return row
	}
	switch {
	case len(parent) != cfg.ExpectedSourceRows || manifest.RowCount != cfg.ExpectedSourceRows:
		reject(fmt.Sprintf("range state construction loop source rows=%d manifest_rows=%d expected=%d", len(parent), manifest.RowCount, cfg.ExpectedSourceRows))
	case manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime:
		reject(fmt.Sprintf("range state construction loop source first_open_time=%s expected=%s", manifest.FirstOpenTime, cfg.ExpectedFirstOpenTime))
	case manifest.LastOpenTime != cfg.ExpectedLastOpenTime:
		reject(fmt.Sprintf("range state construction loop source last_open_time=%s expected=%s", manifest.LastOpenTime, cfg.ExpectedLastOpenTime))
	case manifest.GapCount != cfg.ExpectedGapCount:
		reject(fmt.Sprintf("range state construction loop source gap_count=%d expected=%d", manifest.GapCount, cfg.ExpectedGapCount))
	case manifest.DuplicateCount != cfg.ExpectedDuplicateCount:
		reject(fmt.Sprintf("range state construction loop source duplicate_count=%d expected=%d", manifest.DuplicateCount, cfg.ExpectedDuplicateCount))
	case manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount:
		reject(fmt.Sprintf("range state construction loop source zero_volume_count=%d expected=%d", manifest.ZeroVolumeCount, cfg.ExpectedZeroVolumeCount))
	}
	return row
}

func rangeStateConstructionLoopCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesRangeStateConstructionLoopAuditConfig) FuturesRangeStateConstructionLoopCoverageRow {
	row := FuturesRangeStateConstructionLoopCoverageRow{
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
		row.ValidationError = fmt.Sprintf("unsupported range state construction loop timeframe %q", base.Timeframe)
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

func rangeStateConstructionLoopFeatureWindowRow(frame rangeDiscoveryFrameDef, cfg FuturesRangeStateConstructionLoopAuditConfig, horizons []int) FuturesRangeStateConstructionLoopFeatureWindowRow {
	detectorLookbackBars := cfg.DetectorLookbackDays * frame.barsPerDay
	if cfg.DetectorLookbackBarsOverride > 0 {
		detectorLookbackBars = cfg.DetectorLookbackBarsOverride
	}
	return FuturesRangeStateConstructionLoopFeatureWindowRow{
		Timeframe:                  frame.timeframe,
		BarsPerDay:                 frame.barsPerDay,
		DetectorProfileID:          rangeStateConstructionLoopDetectorProfileID,
		DetectorLookbackDays:       cfg.DetectorLookbackDays,
		DetectorLookbackBars:       detectorLookbackBars,
		DetectorMinConsecutiveBars: cfg.DetectorMinConsecutiveBars,
		ShortWindowBars:            cfg.ShortWindowBars,
		MediumWindowBars:           cfg.MediumWindowBars,
		FeatureLookbackBars:        cfg.FeatureLookbackBars,
		LongLookbackBars:           cfg.LongLookbackDays * frame.barsPerDay,
		ATRPeriod:                  cfg.ATRPeriod,
		ADXPeriod:                  cfg.ADXPeriod,
		HorizonsBars:               rangeStateConstructionLoopIntList(horizons),
		MaxHorizonBars:             maxIntInSlice(horizons),
		ClosedCandleOnly:           true,
		ForwardLabelsAsFeatures:    false,
	}
}

func rangeStateConstructionLoopMetrics(candles []Candle, frame rangeDiscoveryFrameDef, cfg FuturesRangeStateConstructionLoopAuditConfig) rangeStateMetrics {
	normalizedATR := NormalizedATR(candles, cfg.ATRPeriod)
	atr := ATR(candles, cfg.ATRPeriod)
	realizedVol := rangeStateRollingRealizedVol(candles, cfg.ShortWindowBars)
	volume := make([]float64, len(candles))
	tr := make([]float64, len(candles))
	for i := range candles {
		volume[i] = candles[i].Volume
		tr[i] = trueRangeAt(candles, i)
	}
	longLookback := cfg.LongLookbackDays * frame.barsPerDay
	return rangeStateMetrics{
		normalizedATR:         normalizedATR,
		atr:                   atr,
		adx:                   ADX(candles, cfg.ADXPeriod),
		realizedVol:           realizedVol,
		atrPercentile:         rangeStateRollingPriorPercentRank(normalizedATR, longLookback),
		realizedVolPercentile: rangeStateRollingPriorPercentRank(realizedVol, longLookback),
		volumePercentile:      rangeStateRollingPriorPercentRank(volume, longLookback),
		trueRange:             tr,
	}
}

func rangeStateConstructionLoopDetectorConfig(cfg FuturesRangeStateConstructionLoopAuditConfig, barsPerDay int) RangeDetectorConfig {
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

func rangeStateConstructionLoopStateRows(data rangeStateFrameData, allFrames map[string]rangeStateFrameData, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split, skips rangeStateSkipAccumulator, offset int) ([]FuturesRangeStateConstructionLoopStateRow, error) {
	classifications, err := (CompressionRangeDetector{Config: rangeStateConstructionLoopDetectorConfig(cfg, data.frame.barsPerDay)}).Classify(data.candles)
	if err != nil {
		return nil, err
	}
	rows := []FuturesRangeStateConstructionLoopStateRow{}
	episodeID := 0
	for i := 0; i < len(classifications); {
		if !classifications[i].RawActive {
			i++
			continue
		}
		episodeID++
		start := i
		high := data.candles[i].High
		low := data.candles[i].Low
		for i < len(classifications) && classifications[i].RawActive {
			candle := data.candles[i]
			if candle.High > high {
				high = candle.High
			}
			if candle.Low < low {
				low = candle.Low
			}
			split := splitNameForCloseTime(candle.CloseTime, splits)
			switch {
			case !classifications[i].Active:
				skips.add(data.frame.timeframe, split, "not_mature_active")
			case i+data.maxHorizon >= len(data.candles):
				skips.add(data.frame.timeframe, split, "missing_future")
			case i-cfg.MediumWindowBars+1 < 0 || i-cfg.FeatureLookbackBars+1 < 0:
				skips.add(data.frame.timeframe, split, "insufficient_feature_window")
			default:
				row, ok, reason := rangeStateConstructionLoopStateRow(data, allFrames, cfg, splits, offset+len(rows)+1, episodeID, start, i, high, low)
				if !ok {
					skips.add(data.frame.timeframe, split, reason)
				} else {
					rows = append(rows, row)
				}
			}
			i++
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		return rows[i].DecisionIndex < rows[j].DecisionIndex
	})
	return rows, nil
}

func rangeStateConstructionLoopStateRow(data rangeStateFrameData, allFrames map[string]rangeStateFrameData, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split, rowID, episodeID, start, decision int, high, low float64) (FuturesRangeStateConstructionLoopStateRow, bool, string) {
	candles := data.candles
	candle := candles[decision]
	width := high - low
	if width <= 0 || !finitePositive(width) || candle.Close <= 0 {
		return FuturesRangeStateConstructionLoopStateRow{}, false, "non_positive_width_or_price"
	}
	atrShort := rangeStateAvg(data.metrics.normalizedATR, decision-cfg.ShortWindowBars+1, decision)
	if !validNumber(atrShort) || atrShort <= 0 {
		return FuturesRangeStateConstructionLoopStateRow{}, false, "missing_atr"
	}
	atrPct := rangeStateValueAt(data.metrics.atrPercentile, decision)
	realizedPct := rangeStateValueAt(data.metrics.realizedVolPercentile, decision)
	volumePct := rangeStateValueAt(data.metrics.volumePercentile, decision)
	adx := rangeStateValueAt(data.metrics.adx, decision)
	if !validNumber(atrPct) || !validNumber(realizedPct) || !validNumber(volumePct) || !validNumber(adx) {
		return FuturesRangeStateConstructionLoopStateRow{}, false, "missing_percentile_or_adx"
	}
	widthPct := width / candle.Close
	if widthPct <= 0 || !validNumber(widthPct) {
		return FuturesRangeStateConstructionLoopStateRow{}, false, "invalid_width_pct"
	}
	mid := (high + low) / 2
	lookbackStart := maxInt(start, decision-cfg.FeatureLookbackBars+1)
	closeLocation := (candle.Close - low) / width
	maNow := rangeStateMeanClose(candles, decision-cfg.ShortWindowBars+1, decision)
	maPrev := rangeStateMeanClose(candles, decision-cfg.ShortWindowBars*2+1, decision-cfg.ShortWindowBars)
	mediumMA := rangeStateMeanClose(candles, decision-cfg.MediumWindowBars+1, decision)
	if !validNumber(maNow) || !validNumber(maPrev) || !validNumber(mediumMA) || mediumMA <= 0 {
		return FuturesRangeStateConstructionLoopStateRow{}, false, "missing_moving_average"
	}
	retShort := rangeStateReturn(candles, decision, cfg.ShortWindowBars)
	retMedium := rangeStateReturn(candles, decision, cfg.MediumWindowBars)
	maSlope := 0.0
	if maPrev > 0 {
		maSlope = (maNow - maPrev) / maPrev
	}
	avgTRShort := rangeStateAvg(data.metrics.trueRange, decision-cfg.ShortWindowBars+1, decision)
	avgTRPrior := rangeStateAvg(data.metrics.trueRange, decision-cfg.ShortWindowBars*2+1, decision-cfg.ShortWindowBars)
	trExpansion := 0.0
	if avgTRPrior > 0 {
		trExpansion = avgTRShort / avgTRPrior
	}
	lastSide, barsSinceAbnormal, largeBodies, largeRanges, continuation, exhaustion := rangeStateImpulseFeatures(candles, data.metrics.atr, data.metrics.trueRange, decision, lookbackStart, cfg)
	volChange := rangeStateVolumeChangeRatio(candles, decision, cfg.ShortWindowBars)
	volPerWidth := 0.0
	if widthPct > 0 {
		volPerWidth = rangeStateAvgVolume(candles, decision-cfg.ShortWindowBars+1, decision) / widthPct
	}
	body := math.Abs(candle.Close - candle.Open)
	upperWick := candle.High - math.Max(candle.Open, candle.Close)
	lowerWick := math.Min(candle.Open, candle.Close) - candle.Low
	wickBody := 0.0
	if body > 0 {
		wickBody = (upperWick + lowerWick) / body
	} else if candle.High > candle.Low {
		wickBody = 999
	}
	row := FuturesRangeStateConstructionLoopStateRow{
		StateRowID:                       rowID,
		Timestamp:                        candle.CloseTime.UTC().Format(timeLayout),
		Timeframe:                        data.frame.timeframe,
		Split:                            splitNameForCloseTime(candle.CloseTime, splits),
		RangeEpisodeID:                   episodeID,
		RangeStartIndex:                  start,
		DecisionIndex:                    decision,
		RangeStartTime:                   candles[start].CloseTime.UTC().Format(timeLayout),
		DecisionCloseTime:                candle.CloseTime.UTC().Format(timeLayout),
		RangeHigh:                        high,
		RangeLow:                         low,
		RangeMid:                         mid,
		RangeWidth:                       width,
		RangeAgeBars:                     decision - start + 1,
		RangeWidthPct:                    widthPct,
		RangeWidthATRRatio:               widthPct / atrShort,
		CloseLocationPct:                 closeLocation,
		DistanceToLowPct:                 (candle.Close - low) / candle.Close,
		DistanceToHighPct:                (high - candle.Close) / candle.Close,
		DistanceToMidPct:                 math.Abs(candle.Close-mid) / candle.Close,
		BoundaryTouchCountLookback:       rangeContextTriageBoundaryTouchCount(candles, lookbackStart, decision, low, high),
		MidlineCrossCountLookback:        rangeContextTriageMidCrossCount(candles, lookbackStart, decision, mid),
		CloseLocationEntropyLookback:     rangeStateCloseLocationEntropy(candles, lookbackStart, decision, low, high),
		RangeDutyCycleLookback:           rangeStateRangeDutyCycle(candles, lookbackStart, decision, low, high),
		ATRShort:                         atrShort,
		ATRPercentileLong:                atrPct,
		RealizedVolPercentileLong:        realizedPct,
		ATRToRangeWidth:                  atrShort / widthPct,
		TrueRangeExpansionRatio:          trExpansion,
		AbnormalRangeBarCountLookback:    largeRanges,
		ReturnShortPct:                   retShort,
		ReturnMediumPct:                  retMedium,
		MovingAverageSlopePct:            maSlope,
		HigherTimeframeDirectionProxyPct: rangeStateHigherTimeframeProxy(data.frame.timeframe, candle.CloseTime, allFrames, cfg),
		ADX:                              adx,
		DistanceFromMAPct:                (candle.Close - mediumMA) / candle.Close,
		LastAbnormalCandleSide:           lastSide,
		BarsSinceLastAbnormalCandle:      barsSinceAbnormal,
		LargeBodyCandleCountLookback:     largeBodies,
		LargeRangeCandleCountLookback:    largeRanges,
		ImpulseContinuationPressure:      continuation,
		ImpulseExhaustionProxy:           exhaustion,
		VolumePercentileLong:             volumePct,
		VolumeChangeRatio:                volChange,
		VolumePerRangeWidthProxy:         volPerWidth,
		CandleSpreadProxy:                (candle.High - candle.Low) / candle.Close,
		WickBodyStructureSummary:         wickBody,
		ZeroVolumeRowFlag:                candle.Volume == 0,
		Eligible:                         true,
	}
	row.GeometryBucket = rangeStateGeometryBucket(row, cfg)
	row.VolBucket = rangeStateVolBucket(row)
	row.TrendBucket = rangeStateTrendBucket(row)
	row.ImpulseBucket = rangeStateImpulseBucket(row, cfg)
	row.ParticipationBucket = rangeStateParticipationBucket(row)
	row.StateID = fmt.Sprintf("range_state_v1::%s::%s::%s::%s::%s::%s", row.Timeframe, row.GeometryBucket, row.VolBucket, row.TrendBucket, row.ImpulseBucket, row.ParticipationBucket)
	row.GeometryVolID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.VolBucket}, "::")
	row.GeometryTrendID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.TrendBucket}, "::")
	row.GeometryImpulseID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.ImpulseBucket}, "::")
	row.GeometryParticipationID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.ParticipationBucket}, "::")
	row.GeometryVolTrendID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.VolBucket, row.TrendBucket}, "::")
	row.GeometryVolTrendImpulseID = strings.Join([]string{row.Timeframe, row.GeometryBucket, row.VolBucket, row.TrendBucket, row.ImpulseBucket}, "::")
	row.AllFamiliesID = row.StateID
	return row, true, ""
}

func rangeStateConstructionLoopLabelRows(states []FuturesRangeStateConstructionLoopStateRow, frames map[string]rangeStateFrameData, cfg FuturesRangeStateConstructionLoopAuditConfig) []FuturesRangeStateConstructionLoopLabelRow {
	rows := []FuturesRangeStateConstructionLoopLabelRow{}
	for _, state := range states {
		data := frames[state.Timeframe]
		for _, horizon := range data.horizons {
			rows = append(rows, rangeStateConstructionLoopLabelRow(data.candles, state, horizon, cfg))
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].StateRowID != rows[j].StateRowID {
			return rows[i].StateRowID < rows[j].StateRowID
		}
		return rows[i].HorizonBars < rows[j].HorizonBars
	})
	return rows
}

func rangeStateConstructionLoopLabelRow(candles []Candle, state FuturesRangeStateConstructionLoopStateRow, horizon int, cfg FuturesRangeStateConstructionLoopAuditConfig) FuturesRangeStateConstructionLoopLabelRow {
	row := FuturesRangeStateConstructionLoopLabelRow{
		StateRowID:                state.StateRowID,
		StateID:                   state.StateID,
		Timestamp:                 state.Timestamp,
		Timeframe:                 state.Timeframe,
		Split:                     state.Split,
		RangeEpisodeID:            state.RangeEpisodeID,
		HorizonBars:               horizon,
		ForwardLabel:              RangeContextTriageLabelNoResolution,
		LabelWindowStartIndex:     state.DecisionIndex + 1,
		LabelWindowEndIndex:       state.DecisionIndex + horizon,
		FirstOutsideCloseSide:     "none",
		BarsToFirstOutsideClose:   -1,
		GeometryVolID:             state.GeometryVolID,
		GeometryTrendID:           state.GeometryTrendID,
		GeometryImpulseID:         state.GeometryImpulseID,
		GeometryParticipationID:   state.GeometryParticipationID,
		GeometryVolTrendID:        state.GeometryVolTrendID,
		GeometryVolTrendImpulseID: state.GeometryVolTrendImpulseID,
		AllFamiliesID:             state.AllFamiliesID,
		FutureWindowMetadataOnly:  true,
		ForwardLabelUsedAsFeature: false,
	}
	if row.LabelWindowStartIndex < len(candles) {
		row.LabelWindowStartTime = candles[row.LabelWindowStartIndex].CloseTime.UTC().Format(timeLayout)
	}
	if row.LabelWindowEndIndex < len(candles) {
		row.LabelWindowEndTime = candles[row.LabelWindowEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	if horizon <= 0 || row.LabelWindowEndIndex >= len(candles) || state.RangeWidth <= 0 {
		row.ForwardLabel = RangeContextTriageLabelMissingFuture
		return rangeStateFinalizeLabelRow(row)
	}
	insideCloses := 0
	prevMidSide := 0
	prevOutsideState := ""
	outsideTransitionsStarted := false
	firstOutsideIndex := -1
	reentryDeadline := -1
	upperHalfSeen := false
	lowerHalfSeen := false
	for i := row.LabelWindowStartIndex; i <= row.LabelWindowEndIndex; i++ {
		candle := candles[i]
		delay := i - state.DecisionIndex
		if candle.High > state.RangeHigh {
			row.MaxExcursionAboveRangeWidths = math.Max(row.MaxExcursionAboveRangeWidths, (candle.High-state.RangeHigh)/state.RangeWidth)
		}
		if candle.Low < state.RangeLow {
			row.MaxExcursionBelowRangeWidths = math.Max(row.MaxExcursionBelowRangeWidths, (state.RangeLow-candle.Low)/state.RangeWidth)
		}
		inside := candle.Close >= state.RangeLow && candle.Close <= state.RangeHigh
		if inside {
			insideCloses++
			if candle.Close >= state.RangeMid {
				upperHalfSeen = true
			}
			if candle.Close <= state.RangeMid {
				lowerHalfSeen = true
			}
		}
		midSide := 0
		if candle.Close > state.RangeMid {
			midSide = 1
		} else if candle.Close < state.RangeMid {
			midSide = -1
		}
		if prevMidSide != 0 && midSide != 0 && midSide != prevMidSide {
			row.MidpointCrossCount++
		}
		if midSide != 0 {
			prevMidSide = midSide
		}
		outsideState := "inside"
		if candle.Close > state.RangeHigh {
			outsideState = RangeDiscoverySideUp
		} else if candle.Close < state.RangeLow {
			outsideState = RangeDiscoverySideDown
		}
		if firstOutsideIndex < 0 && outsideState != "inside" {
			firstOutsideIndex = i
			row.BarsToFirstOutsideClose = delay
			row.FirstOutsideCloseSide = outsideState
			reentryDeadline = i + cfg.ReentryWindowBars
		}
		if firstOutsideIndex >= 0 {
			if outsideTransitionsStarted && outsideState != prevOutsideState {
				row.OutsideStateTransitionCount++
			}
			prevOutsideState = outsideState
			outsideTransitionsStarted = true
			if !row.ReentryWithinWindow && i > firstOutsideIndex && i <= reentryDeadline && inside {
				row.ReentryWithinWindow = true
			}
		}
	}
	count := row.LabelWindowEndIndex - row.LabelWindowStartIndex + 1
	if count > 0 {
		row.InsideCloseRate = float64(insideCloses) / float64(count)
	}
	row.FinalClosePosition = (candles[row.LabelWindowEndIndex].Close - state.RangeLow) / state.RangeWidth
	row.ForwardLabel = rangeStateForwardLabel(row, state, upperHalfSeen, lowerHalfSeen, cfg)
	return rangeStateFinalizeLabelRow(row)
}

func rangeStateForwardLabel(row FuturesRangeStateConstructionLoopLabelRow, state FuturesRangeStateConstructionLoopStateRow, upperHalfSeen, lowerHalfSeen bool, cfg FuturesRangeStateConstructionLoopAuditConfig) string {
	switch {
	case state.RangeWidthPct < cfg.LowWidthPctThreshold || state.RangeWidthATRRatio < cfg.LowWidthATRThreshold:
		return RangeContextTriageLabelLowWidthNoise
	case row.FirstOutsideCloseSide == "none" && row.MidpointCrossCount >= 2 && upperHalfSeen && lowerHalfSeen:
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

func rangeStateFinalizeLabelRow(row FuturesRangeStateConstructionLoopLabelRow) FuturesRangeStateConstructionLoopLabelRow {
	row.RotationUseful = row.ForwardLabel == RangeContextTriageLabelContainedRotation || row.ForwardLabel == RangeContextTriageLabelFalseBreakReentryUp || row.ForwardLabel == RangeContextTriageLabelFalseBreakReentryDown
	row.RotationToxic = row.ForwardLabel == RangeContextTriageLabelBoundaryChop || row.ForwardLabel == RangeContextTriageLabelCleanExpansionUp || row.ForwardLabel == RangeContextTriageLabelCleanExpansionDown
	row.ContinuationUseful = row.ForwardLabel == RangeContextTriageLabelCleanExpansionUp || row.ForwardLabel == RangeContextTriageLabelCleanExpansionDown || row.ForwardLabel == RangeContextTriageLabelDriftThroughUp || row.ForwardLabel == RangeContextTriageLabelDriftThroughDown
	row.ContinuationToxic = row.ForwardLabel == RangeContextTriageLabelFalseBreakReentryUp || row.ForwardLabel == RangeContextTriageLabelFalseBreakReentryDown || row.ForwardLabel == RangeContextTriageLabelBoundaryChop
	row.NoTradeToxic = row.ForwardLabel == RangeContextTriageLabelBoundaryChop || row.ForwardLabel == RangeContextTriageLabelLowWidthNoise || row.ForwardLabel == RangeContextTriageLabelNoResolution
	row.DiagnosticOnly = !row.RotationUseful && !row.RotationToxic && !row.ContinuationUseful && !row.ContinuationToxic && !row.NoTradeToxic
	return row
}

func rangeStateConstructionLoopCohortRows(labels []FuturesRangeStateConstructionLoopLabelRow, sources []FuturesRangeStateConstructionLoopSourceRow, coverage []FuturesRangeStateConstructionLoopCoverageRow, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split) []FuturesRangeStateConstructionLoopCohortRow {
	acc := map[rangeStateCohortKey]*rangeStateCohortAccumulator{}
	sourcePass := rangeStateConstructionLoopSourceResamplePass(sources, coverage)
	for _, label := range labels {
		for _, split := range rangeDiscoverySplitCombos(label.Split) {
			for _, rollup := range rangeStateRollups(label) {
				for _, route := range []string{RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation, RangeStateConstructionLoopRouteNoTradeToxic, RangeStateConstructionLoopRouteDiagnosticOnly} {
					key := rangeStateCohortKeyFor(split, label.Timeframe, label.HorizonBars, route, rollup.rollupType, rollup.rollupID)
					a := acc[key]
					if a == nil {
						a = &rangeStateCohortAccumulator{labels: map[string]int{}}
						a.row = FuturesRangeStateConstructionLoopCohortRow{
							CohortID:                   key.cohortID,
							RouteCandidate:             key.routeCandidate,
							RollupType:                 key.rollupType,
							RollupID:                   key.rollupID,
							Split:                      key.split,
							Timeframe:                  key.timeframe,
							HorizonBars:                key.horizonBars,
							ClosedFamilyProtectionPass: true,
							FutureLeakProtectionPass:   true,
						}
						if !sourcePass {
							a.row.FailureReason = "source_or_resample_gap"
						}
						acc[key] = a
					}
					a.add(label)
				}
			}
		}
	}
	rows := make([]FuturesRangeStateConstructionLoopCohortRow, 0, len(acc))
	for _, a := range acc {
		rows = append(rows, a.finalRow())
	}
	rangeStateMarkCohortGates(rows, cfg, splits, sourcePass)
	sort.Slice(rows, func(i, j int) bool {
		return rangeStateLessCohort(rows[i], rows[j])
	})
	return rows
}

func (acc *rangeStateCohortAccumulator) add(label FuturesRangeStateConstructionLoopLabelRow) {
	acc.row.CandidateCount++
	acc.labels[label.ForwardLabel]++
	if rangeStateRouteUseful(acc.row.RouteCandidate, label) {
		acc.row.UsefulCount++
	}
	if rangeStateRouteToxic(acc.row.RouteCandidate, label) {
		acc.row.ToxicCount++
	}
	switch label.ForwardLabel {
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

func (acc *rangeStateCohortAccumulator) finalRow() FuturesRangeStateConstructionLoopCohortRow {
	row := acc.row
	if row.CandidateCount > 0 {
		row.UsefulRate = float64(row.UsefulCount) / float64(row.CandidateCount)
		row.ToxicRate = float64(row.ToxicCount) / float64(row.CandidateCount)
		row.UsefulMinusToxicMargin = row.UsefulRate - row.ToxicRate
	}
	row.DominantForwardLabel, row.DominantForwardLabelRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, false)
	return row
}

func rangeStateMarkCohortGates(rows []FuturesRangeStateConstructionLoopCohortRow, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split, sourcePass bool) {
	byIDSplit := map[string]map[string]*FuturesRangeStateConstructionLoopCohortRow{}
	for i := range rows {
		if byIDSplit[rows[i].CohortID] == nil {
			byIDSplit[rows[i].CohortID] = map[string]*FuturesRangeStateConstructionLoopCohortRow{}
		}
		byIDSplit[rows[i].CohortID][rows[i].Split] = &rows[i]
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	for i := range rows {
		cohortRows := byIDSplit[rows[i].CohortID]
		full := cohortRows[fullSplitName]
		reasons := []string{}
		fullRows := 0
		if full != nil {
			fullRows = full.CandidateCount
			rows[i].FullUsefulRate = full.UsefulRate
			rows[i].FullToxicRate = full.ToxicRate
			rows[i].FullUsefulMinusToxicMargin = full.UsefulMinusToxicMargin
		}
		rows[i].FullPeriodRows = fullRows
		rows[i].WeakestSplitRows = int(^uint(0) >> 1)
		rows[i].WeakestSplitUsefulRate = math.Inf(1)
		rows[i].WeakestSplitMargin = math.Inf(1)
		splitCountsPass := true
		splitRatesPresent := true
		maxContribution := 0.0
		for _, split := range periodSplits {
			row := cohortRows[split.Name]
			if row == nil {
				splitCountsPass = false
				splitRatesPresent = false
				rows[i].WeakestSplitRows = 0
				rows[i].WeakestSplitUsefulRate = 0
				rows[i].WeakestSplitMargin = 0
				continue
			}
			if row.CandidateCount < cfg.MinSplitCohortCount {
				splitCountsPass = false
			}
			if rows[i].RouteCandidate == RangeStateConstructionLoopRouteNoTradeToxic && row.CandidateCount < cfg.MinNoTradeSplitCohortCount {
				splitCountsPass = false
			}
			if row.CandidateCount < rows[i].WeakestSplitRows {
				rows[i].WeakestSplitRows = row.CandidateCount
			}
			if row.UsefulRate < rows[i].WeakestSplitUsefulRate {
				rows[i].WeakestSplitUsefulRate = row.UsefulRate
			}
			if row.ToxicRate > rows[i].WorstSplitToxicRate {
				rows[i].WorstSplitToxicRate = row.ToxicRate
			}
			if row.UsefulMinusToxicMargin < rows[i].WeakestSplitMargin {
				rows[i].WeakestSplitMargin = row.UsefulMinusToxicMargin
			}
			if fullRows > 0 {
				maxContribution = math.Max(maxContribution, float64(row.CandidateCount)/float64(fullRows))
			}
		}
		if rows[i].WeakestSplitRows == int(^uint(0)>>1) {
			rows[i].WeakestSplitRows = 0
		}
		if math.IsInf(rows[i].WeakestSplitUsefulRate, 1) {
			rows[i].WeakestSplitUsefulRate = 0
		}
		if math.IsInf(rows[i].WeakestSplitMargin, 1) {
			rows[i].WeakestSplitMargin = 0
		}
		rows[i].MaxSplitContributionRate = maxContribution
		fullCountGate := fullRows >= cfg.MinFullCohortCount
		if rows[i].RouteCandidate == RangeStateConstructionLoopRouteNoTradeToxic {
			fullCountGate = fullRows >= cfg.MinNoTradeFullCohortCount
		}
		rows[i].ReviewableCountGatePass = fullCountGate && splitCountsPass
		rows[i].SplitStabilityGatePass = splitRatesPresent && len(periodSplits) > 0
		rows[i].SplitContributionGatePass = maxContribution <= cfg.MaxSplitContributionRate || len(periodSplits) == 0
		rows[i].DominantLabelGatePass = rows[i].RouteCandidate == RangeStateConstructionLoopRouteNoTradeToxic || (full != nil && full.DominantForwardLabelRate <= cfg.MaxDominantRouteLabelRate)
		rows[i].FeatureKnownGatePass = !strings.Contains(rows[i].RollupID, "unknown")
		rows[i].ClosedFamilyProtectionPass = !rows[i].ClosedFamilyReslice
		rows[i].FutureLeakProtectionPass = rows[i].RollupType != "forward_label"
		switch rows[i].RouteCandidate {
		case RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation:
			rows[i].RouteRateGatePass = full != nil &&
				rows[i].FullUsefulRate >= cfg.MinPositiveUsefulRateFull &&
				rows[i].WeakestSplitUsefulRate >= cfg.MinPositiveUsefulRateSplit &&
				rows[i].FullToxicRate <= cfg.MaxPositiveToxicRateFull &&
				rows[i].WorstSplitToxicRate <= cfg.MaxPositiveToxicRateSplit &&
				rows[i].FullUsefulMinusToxicMargin >= cfg.MinPositiveMarginFull &&
				rows[i].WeakestSplitMargin >= cfg.MinPositiveMarginSplit
		case RangeStateConstructionLoopRouteNoTradeToxic:
			rows[i].RouteRateGatePass = full != nil &&
				rows[i].FullToxicRate >= cfg.MinToxicRateFull &&
				rows[i].WorstSplitToxicRate >= cfg.MinToxicRateSplit
		default:
			rows[i].RouteRateGatePass = false
		}
		if !sourcePass {
			reasons = append(reasons, "source_or_resample_gap")
		}
		if !rows[i].ReviewableCountGatePass {
			reasons = append(reasons, "inadequate_cohort_count")
		}
		if !rows[i].SplitStabilityGatePass {
			reasons = append(reasons, "missing_period_split")
		}
		if !rows[i].SplitContributionGatePass {
			reasons = append(reasons, "single_split_contribution_above_gate")
		}
		if !rows[i].DominantLabelGatePass {
			reasons = append(reasons, "single_forward_label_above_gate")
		}
		if !rows[i].FeatureKnownGatePass {
			reasons = append(reasons, "unknown_feature_bucket")
		}
		if !rows[i].RouteRateGatePass {
			reasons = append(reasons, "route_rate_gate_failed")
		}
		if rows[i].ClosedFamilyReslice {
			reasons = append(reasons, "closed_family_reslice")
		}
		if !rows[i].FutureLeakProtectionPass {
			reasons = append(reasons, "future_label_as_feature")
		}
		rows[i].PassesReviewGate = sourcePass &&
			rows[i].ReviewableCountGatePass &&
			rows[i].SplitStabilityGatePass &&
			rows[i].SplitContributionGatePass &&
			rows[i].DominantLabelGatePass &&
			rows[i].FeatureKnownGatePass &&
			rows[i].RouteRateGatePass &&
			rows[i].ClosedFamilyProtectionPass &&
			rows[i].FutureLeakProtectionPass &&
			rows[i].Split == fullSplitName
		rows[i].FailureReason = uniqueJoinedReasons(reasons)
	}
}

func rangeStateConstructionLoopRankingRows(cohorts []FuturesRangeStateConstructionLoopCohortRow, splits []Split) []FuturesRangeStateConstructionLoopRankingRow {
	rows := []FuturesRangeStateConstructionLoopRankingRow{}
	for _, cohort := range cohorts {
		if cohort.Split != fullSplitName || cohort.RouteCandidate == RangeStateConstructionLoopRouteDiagnosticOnly {
			continue
		}
		rank := FuturesRangeStateConstructionLoopRankingRow{
			CohortID:                   cohort.CohortID,
			RouteCandidate:             cohort.RouteCandidate,
			RollupType:                 cohort.RollupType,
			RollupID:                   cohort.RollupID,
			Timeframe:                  cohort.Timeframe,
			HorizonBars:                cohort.HorizonBars,
			PassesGate:                 cohort.PassesReviewGate,
			FullPeriodRows:             cohort.FullPeriodRows,
			WeakestSplitRows:           cohort.WeakestSplitRows,
			MaxSplitContributionRate:   cohort.MaxSplitContributionRate,
			FullUsefulRate:             cohort.FullUsefulRate,
			WeakestSplitUsefulRate:     cohort.WeakestSplitUsefulRate,
			FullToxicRate:              cohort.FullToxicRate,
			WorstSplitToxicRate:        cohort.WorstSplitToxicRate,
			FullUsefulMinusToxicMargin: cohort.FullUsefulMinusToxicMargin,
			WeakestSplitMargin:         cohort.WeakestSplitMargin,
			DominantForwardLabel:       cohort.DominantForwardLabel,
			DominantForwardLabelRate:   cohort.DominantForwardLabelRate,
			ClosedFamilyReslice:        cohort.ClosedFamilyReslice,
			FutureLeakProtectionPass:   cohort.FutureLeakProtectionPass,
			FailureReason:              cohort.FailureReason,
		}
		rank.RankScore = rangeStateRankScore(rank)
		rows = append(rows, rank)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		if rows[i].RouteCandidate != rows[j].RouteCandidate {
			return rangeStateRouteSortKey(rows[i].RouteCandidate) < rangeStateRouteSortKey(rows[j].RouteCandidate)
		}
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].HorizonBars != rows[j].HorizonBars {
			return rows[i].HorizonBars < rows[j].HorizonBars
		}
		return rows[i].CohortID < rows[j].CohortID
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func rangeStateConstructionLoopSummaryRows(result FuturesRangeStateConstructionLoopAuditResult, cfg FuturesRangeStateConstructionLoopAuditConfig, splits []Split) []FuturesRangeStateConstructionLoopSummaryRow {
	sourcePass := rangeStateConstructionLoopSourceResamplePass(result.SourceRows, result.CoverageRows)
	rows := []FuturesRangeStateConstructionLoopSummaryRow{{
		Split:              fullSplitName,
		Timeframe:          "all",
		HorizonBars:        0,
		SourceResamplePass: sourcePass,
		StateRows:          len(result.StateRows),
		LabelRows:          len(result.LabelRows),
		CohortRows:         len(result.CohortRows),
		RankingRows:        len(result.RankingRows),
		SkipRows:           rangeStateSkipTotal(result.SkipRows),
		PassingCohorts:     result.PassingCohorts,
		StopState:          result.StopState,
	}}
	for _, split := range splits {
		for _, timeframe := range cfg.Timeframes {
			for _, horizon := range rangeStateConstructionLoopHorizons(timeframe, cfg) {
				row := FuturesRangeStateConstructionLoopSummaryRow{
					Split:              split.Name,
					Timeframe:          timeframe,
					HorizonBars:        horizon,
					SourceResamplePass: sourcePass,
					StopState:          result.StopState,
				}
				for _, state := range result.StateRows {
					if state.Timeframe == timeframe && rangeStateRowInSplit(state.DecisionCloseTime, split) {
						row.StateRows++
					}
				}
				for _, label := range result.LabelRows {
					if label.Timeframe == timeframe && label.HorizonBars == horizon && rangeStateRowInSplit(label.Timestamp, split) {
						row.LabelRows++
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
				for _, skip := range result.SkipRows {
					if skip.Timeframe == timeframe && skip.Split == split.Name {
						row.SkipRows += skip.Count
					}
				}
				rows = append(rows, row)
			}
		}
	}
	return rows
}

func rangeStateConstructionLoopSourceResamplePass(sources []FuturesRangeStateConstructionLoopSourceRow, coverage []FuturesRangeStateConstructionLoopCoverageRow) bool {
	if len(sources) == 0 {
		return false
	}
	for _, row := range sources {
		if row.ValidationStatus != "accepted" || !row.SourceFactsPass {
			return false
		}
	}
	if len(coverage) == 0 {
		return false
	}
	for _, row := range coverage {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass || !row.Complete {
			return false
		}
	}
	return true
}

func rangeStateConstructionLoopHorizons(timeframe string, cfg FuturesRangeStateConstructionLoopAuditConfig) []int {
	switch timeframe {
	case RangeDiscoveryTimeframe15m:
		return cfg.Horizons15M
	case RangeDiscoveryTimeframe1h:
		return cfg.Horizons1H
	case RangeDiscoveryTimeframe4h:
		return cfg.Horizons4H
	default:
		return nil
	}
}

func rangeStateConstructionLoopSkipRows(skips rangeStateSkipAccumulator) []FuturesRangeStateConstructionLoopSkipRow {
	rows := make([]FuturesRangeStateConstructionLoopSkipRow, 0, len(skips))
	for key, count := range skips {
		rows = append(rows, FuturesRangeStateConstructionLoopSkipRow{
			Timeframe: key.timeframe,
			Split:     key.split,
			Reason:    key.reason,
			Count:     count,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].Split != rows[j].Split {
			return splitSortKey(rows[i].Split) < splitSortKey(rows[j].Split)
		}
		return rows[i].Reason < rows[j].Reason
	})
	return rows
}

func (skips rangeStateSkipAccumulator) add(timeframe, split, reason string) {
	if split == "" {
		split = fullSplitName
	}
	skips[rangeStateSkipKey{timeframe: timeframe, split: split, reason: reason}]++
	if split != fullSplitName {
		skips[rangeStateSkipKey{timeframe: timeframe, split: fullSplitName, reason: reason}]++
	}
}

func rangeStateSkipTotal(rows []FuturesRangeStateConstructionLoopSkipRow) int {
	total := 0
	for _, row := range rows {
		if row.Split == fullSplitName {
			total += row.Count
		}
	}
	return total
}

func rangeStateRollups(label FuturesRangeStateConstructionLoopLabelRow) []struct {
	rollupType string
	rollupID   string
} {
	return []struct {
		rollupType string
		rollupID   string
	}{
		{RangeStateConstructionLoopRollupGeometryVol, label.GeometryVolID},
		{RangeStateConstructionLoopRollupGeometryTrend, label.GeometryTrendID},
		{RangeStateConstructionLoopRollupGeometryImpulse, label.GeometryImpulseID},
		{RangeStateConstructionLoopRollupGeometryParticipation, label.GeometryParticipationID},
		{RangeStateConstructionLoopRollupGeometryVolTrend, label.GeometryVolTrendID},
		{RangeStateConstructionLoopRollupGeometryVolTrendImpulse, label.GeometryVolTrendImpulseID},
		{RangeStateConstructionLoopRollupAllFamilies, label.AllFamiliesID},
	}
}

func rangeStateCohortKeyFor(split, timeframe string, horizon int, route, rollupType, rollupID string) rangeStateCohortKey {
	cohortID := strings.Join([]string{"range_state_v1", timeframe, fmt.Sprintf("h%d", horizon), route, rollupType, rollupID}, "|")
	return rangeStateCohortKey{
		cohortID:       cohortID,
		routeCandidate: route,
		rollupType:     rollupType,
		rollupID:       rollupID,
		split:          split,
		timeframe:      timeframe,
		horizonBars:    horizon,
	}
}

func rangeStateRouteUseful(route string, label FuturesRangeStateConstructionLoopLabelRow) bool {
	switch route {
	case RangeStateConstructionLoopRouteRotation:
		return label.RotationUseful
	case RangeStateConstructionLoopRouteContinuation:
		return label.ContinuationUseful
	default:
		return false
	}
}

func rangeStateRouteToxic(route string, label FuturesRangeStateConstructionLoopLabelRow) bool {
	switch route {
	case RangeStateConstructionLoopRouteRotation:
		return label.RotationToxic
	case RangeStateConstructionLoopRouteContinuation:
		return label.ContinuationToxic
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return label.NoTradeToxic
	default:
		return false
	}
}

func rangeStateRankScore(row FuturesRangeStateConstructionLoopRankingRow) float64 {
	switch row.RouteCandidate {
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return row.FullToxicRate + row.WorstSplitToxicRate - row.MaxSplitContributionRate + math.Log1p(float64(row.FullPeriodRows))/20
	default:
		return row.FullUsefulMinusToxicMargin + row.WeakestSplitMargin + row.FullUsefulRate - row.WorstSplitToxicRate + math.Log1p(float64(row.WeakestSplitRows))/20
	}
}

func rangeStateGeometryBucket(row FuturesRangeStateConstructionLoopStateRow, cfg FuturesRangeStateConstructionLoopAuditConfig) string {
	switch {
	case row.BoundaryTouchCountLookback >= 8:
		return "geometry_boundary_crowded"
	case row.MidlineCrossCountLookback >= 4 && row.CloseLocationEntropyLookback >= 0.80:
		return "geometry_midline_balanced"
	case row.RangeWidthPct >= cfg.WideWidthPctThreshold || row.RangeWidthATRRatio >= cfg.WideWidthATRThreshold:
		return "geometry_wide_volatile"
	case row.RangeWidthPct < 0.0030 && row.RangeWidthATRRatio < 1.50 && row.RangeDutyCycleLookback >= 0.85:
		return "geometry_narrow_orderly"
	default:
		return "geometry_balanced_orderly"
	}
}

func rangeStateVolBucket(row FuturesRangeStateConstructionLoopStateRow) string {
	switch {
	case row.ATRPercentileLong >= 0.90 || row.RealizedVolPercentileLong >= 0.90 || row.AbnormalRangeBarCountLookback >= 6:
		return "vol_extreme"
	case row.TrueRangeExpansionRatio >= 1.35 || row.ATRPercentileLong >= 0.70 || row.RealizedVolPercentileLong >= 0.70:
		return "vol_expanding"
	case row.ATRPercentileLong <= 0.25 && row.RealizedVolPercentileLong <= 0.35:
		return "vol_compressed"
	default:
		return "vol_normal"
	}
}

func rangeStateTrendBucket(row FuturesRangeStateConstructionLoopStateRow) string {
	trendScore := row.ReturnMediumPct + row.MovingAverageSlopePct + row.HigherTimeframeDirectionProxyPct/2
	switch {
	case trendScore >= 0.015 && row.ADX >= 25:
		return "trend_strong_up"
	case trendScore <= -0.015 && row.ADX >= 25:
		return "trend_strong_down"
	case trendScore >= 0.003:
		return "trend_up_pressure"
	case trendScore <= -0.003:
		return "trend_down_pressure"
	default:
		return "trend_flat"
	}
}

func rangeStateImpulseBucket(row FuturesRangeStateConstructionLoopStateRow, cfg FuturesRangeStateConstructionLoopAuditConfig) string {
	switch {
	case row.LargeRangeCandleCountLookback >= 3:
		return "impulse_clustered"
	case row.LastAbnormalCandleSide == "none":
		return "impulse_none"
	case row.BarsSinceLastAbnormalCandle > cfg.MediumWindowBars/2:
		return "impulse_stale"
	case row.LastAbnormalCandleSide == RangeDiscoverySideUp:
		return "impulse_up_recent"
	case row.LastAbnormalCandleSide == RangeDiscoverySideDown:
		return "impulse_down_recent"
	default:
		return "impulse_none"
	}
}

func rangeStateParticipationBucket(row FuturesRangeStateConstructionLoopStateRow) string {
	switch {
	case row.ZeroVolumeRowFlag || row.VolumePercentileLong <= 0.20:
		return "participation_low"
	case row.VolumePercentileLong >= 0.85 && (row.CandleSpreadProxy >= 0.01 || row.WickBodyStructureSummary >= 4):
		return "participation_dislocated"
	case row.VolumePercentileLong >= 0.75 || row.VolumeChangeRatio >= 1.50:
		return "participation_high"
	default:
		return "participation_normal"
	}
}

func rangeStateRollingRealizedVol(candles []Candle, window int) []float64 {
	out := nanSlice(len(candles))
	if window <= 1 {
		return out
	}
	returns := nanSlice(len(candles))
	for i := 1; i < len(candles); i++ {
		if candles[i-1].Close > 0 {
			returns[i] = (candles[i].Close - candles[i-1].Close) / candles[i-1].Close
		}
	}
	for i := window; i < len(candles); i++ {
		sum := 0.0
		count := 0
		for j := i - window + 1; j <= i; j++ {
			if validNumber(returns[j]) {
				sum += returns[j]
				count++
			}
		}
		if count == 0 {
			continue
		}
		mean := sum / float64(count)
		variance := 0.0
		for j := i - window + 1; j <= i; j++ {
			if validNumber(returns[j]) {
				diff := returns[j] - mean
				variance += diff * diff
			}
		}
		out[i] = math.Sqrt(variance / float64(count))
	}
	return out
}

func rangeStateRollingPriorPercentRank(values []float64, lookback int) []float64 {
	out := nanSlice(len(values))
	if lookback <= 0 {
		return out
	}
	window := make([]float64, 0, lookback)
	for i, value := range values {
		if i >= lookback && len(window) > 0 && validNumber(value) {
			idx := sort.SearchFloat64s(window, value)
			out[i] = float64(idx) / float64(len(window))
		}
		if validNumber(value) {
			window = insertSorted(window, value)
		}
		if i >= lookback && validNumber(values[i-lookback]) {
			window = removeSorted(window, values[i-lookback])
		}
	}
	return out
}

func rangeStateAvg(values []float64, start, end int) float64 {
	if start < 0 || end >= len(values) || end < start {
		return math.NaN()
	}
	sum := 0.0
	count := 0
	for i := start; i <= end; i++ {
		if validNumber(values[i]) {
			sum += values[i]
			count++
		}
	}
	if count == 0 {
		return math.NaN()
	}
	return sum / float64(count)
}

func rangeStateValueAt(values []float64, index int) float64 {
	if index < 0 || index >= len(values) {
		return math.NaN()
	}
	return values[index]
}

func rangeStateReturn(candles []Candle, index, bars int) float64 {
	if bars <= 0 || index-bars < 0 || index >= len(candles) || candles[index-bars].Close <= 0 {
		return 0
	}
	return (candles[index].Close - candles[index-bars].Close) / candles[index-bars].Close
}

func rangeStateMeanClose(candles []Candle, start, end int) float64 {
	if start < 0 || end >= len(candles) || end < start {
		return math.NaN()
	}
	sum := 0.0
	for i := start; i <= end; i++ {
		sum += candles[i].Close
	}
	return sum / float64(end-start+1)
}

func rangeStateAvgVolume(candles []Candle, start, end int) float64 {
	if start < 0 || end >= len(candles) || end < start {
		return 0
	}
	sum := 0.0
	for i := start; i <= end; i++ {
		sum += candles[i].Volume
	}
	return sum / float64(end-start+1)
}

func rangeStateVolumeChangeRatio(candles []Candle, index, window int) float64 {
	now := rangeStateAvgVolume(candles, index-window+1, index)
	prior := rangeStateAvgVolume(candles, index-window*2+1, index-window)
	if prior <= 0 {
		return 0
	}
	return now / prior
}

func rangeStateImpulseFeatures(candles []Candle, atr []float64, tr []float64, decision, start int, cfg FuturesRangeStateConstructionLoopAuditConfig) (string, int, int, int, float64, float64) {
	lastSide := "none"
	lastIndex := -1
	lastClose := 0.0
	largeBodies := 0
	largeRanges := 0
	for i := start; i <= decision; i++ {
		if i >= len(candles) || i >= len(atr) || i >= len(tr) {
			continue
		}
		body := math.Abs(candles[i].Close - candles[i].Open)
		if tr[i] > 0 && body/tr[i] >= cfg.LargeBodyToRangeThreshold {
			largeBodies++
		}
		if validNumber(atr[i]) && atr[i] > 0 && tr[i]/atr[i] >= cfg.AbnormalTRATRThreshold {
			largeRanges++
			lastIndex = i
			lastClose = candles[i].Close
			if candles[i].Close >= candles[i].Open {
				lastSide = RangeDiscoverySideUp
			} else {
				lastSide = RangeDiscoverySideDown
			}
		}
	}
	if lastIndex < 0 || lastClose <= 0 {
		return lastSide, 0, largeBodies, largeRanges, 0, 0
	}
	signedMove := (candles[decision].Close - lastClose) / lastClose
	if lastSide == RangeDiscoverySideDown {
		signedMove = -signedMove
	}
	continuation := math.Max(0, signedMove)
	exhaustion := math.Max(0, -signedMove)
	return lastSide, decision - lastIndex, largeBodies, largeRanges, continuation, exhaustion
}

func rangeStateHigherTimeframeProxy(timeframe string, decisionClose time.Time, frames map[string]rangeStateFrameData, cfg FuturesRangeStateConstructionLoopAuditConfig) float64 {
	target := timeframe
	switch timeframe {
	case RangeDiscoveryTimeframe15m:
		target = RangeDiscoveryTimeframe1h
	case RangeDiscoveryTimeframe1h:
		target = RangeDiscoveryTimeframe4h
	}
	data, ok := frames[target]
	if !ok || len(data.candles) == 0 {
		return 0
	}
	idx := sort.Search(len(data.candles), func(i int) bool {
		return !data.candles[i].CloseTime.Before(decisionClose)
	})
	if idx >= len(data.candles) {
		idx = len(data.candles) - 1
	}
	if data.candles[idx].CloseTime.After(decisionClose) && idx > 0 {
		idx--
	}
	return rangeStateReturn(data.candles, idx, cfg.ShortWindowBars)
}

func rangeStateCloseLocationEntropy(candles []Candle, start, end int, low, high float64) float64 {
	if end < start || high <= low {
		return 0
	}
	counts := [4]int{}
	total := 0
	width := high - low
	for i := start; i <= end && i < len(candles); i++ {
		loc := (candles[i].Close - low) / width
		bucket := int(math.Floor(loc * 4))
		if bucket < 0 {
			bucket = 0
		}
		if bucket > 3 {
			bucket = 3
		}
		counts[bucket]++
		total++
	}
	if total == 0 {
		return 0
	}
	entropy := 0.0
	for _, count := range counts {
		if count == 0 {
			continue
		}
		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}
	return entropy / 2
}

func rangeStateRangeDutyCycle(candles []Candle, start, end int, low, high float64) float64 {
	if end < start {
		return 0
	}
	inside := 0
	total := 0
	for i := start; i <= end && i < len(candles); i++ {
		total++
		if candles[i].Close >= low && candles[i].Close <= high {
			inside++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(inside) / float64(total)
}

func rangeStateRowInSplit(timestamp string, split Split) bool {
	parsed, err := time.Parse(timeLayout, timestamp)
	if err != nil {
		return false
	}
	return split.Contains(parsed)
}

func rangeStateConstructionLoopPassingRankingCount(rows []FuturesRangeStateConstructionLoopRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func rangeStateLessCohort(a, b FuturesRangeStateConstructionLoopCohortRow) bool {
	if a.Split != b.Split {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.RouteCandidate != b.RouteCandidate {
		return rangeStateRouteSortKey(a.RouteCandidate) < rangeStateRouteSortKey(b.RouteCandidate)
	}
	if a.Timeframe != b.Timeframe {
		return rangeContextTriageTimeframeSortKey(a.Timeframe) < rangeContextTriageTimeframeSortKey(b.Timeframe)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.RollupType != b.RollupType {
		return rangeStateRollupSortKey(a.RollupType) < rangeStateRollupSortKey(b.RollupType)
	}
	return a.RollupID < b.RollupID
}

func rangeStateRouteSortKey(route string) int {
	switch route {
	case RangeStateConstructionLoopRouteRotation:
		return 0
	case RangeStateConstructionLoopRouteContinuation:
		return 1
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return 2
	case RangeStateConstructionLoopRouteDiagnosticOnly:
		return 3
	default:
		return 99
	}
}

func rangeStateRollupSortKey(rollup string) int {
	switch rollup {
	case RangeStateConstructionLoopRollupGeometryVol:
		return 0
	case RangeStateConstructionLoopRollupGeometryTrend:
		return 1
	case RangeStateConstructionLoopRollupGeometryImpulse:
		return 2
	case RangeStateConstructionLoopRollupGeometryParticipation:
		return 3
	case RangeStateConstructionLoopRollupGeometryVolTrend:
		return 4
	case RangeStateConstructionLoopRollupGeometryVolTrendImpulse:
		return 5
	case RangeStateConstructionLoopRollupAllFamilies:
		return 6
	default:
		return 99
	}
}

func rangeStateConstructionLoopIntList(values []int) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, fmt.Sprint(value))
	}
	return strings.Join(parts, ";")
}
