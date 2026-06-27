package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeRouterRotationPremiseAuditName = "futures_range_router_rotation_premise_audit"

	RangeRouterRotationPremiseStopStateSourceRouterGap          = "range_router_rotation_premise_audit_source_router_gap"
	RangeRouterRotationPremiseStopStateRejectedClosedReslice    = "range_router_rotation_premise_audit_rejected_closed_family_reslice"
	RangeRouterRotationPremiseStopStateNoEligibleEvents         = "range_router_rotation_premise_audit_no_eligible_events"
	RangeRouterRotationPremiseStopStateFailedNoPremise          = "range_router_rotation_premise_audit_failed_no_premise"
	RangeRouterRotationPremiseStopStatePassedNeedsTriggerAudit  = "range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit"
	RangeRouterRotationPremiseStopStateRejectedStrategyBacktest = "range_router_rotation_premise_audit_rejected_as_strategy_backtest_request"
	RangeRouterRotationPremiseName                              = "router_gated_boundary_reclaim_rotation_v1"
	RangeRouterRotationPremiseCohortID                          = "range_router_rotation_premise_v1|15m|h24|router_gated_boundary_reclaim_rotation"
	RangeRouterRotationPremiseRequiredRouterCohortID            = "range_context_router_v1|15m|h24|tradable_rotation"
	rangeRouterRotationPremiseRulePrefix                        = "range_context_router_v1|tradable_rotation|15m|h24|"
	RangeRouterRotationPremiseSideAll                           = "all"
	RangeRouterRotationPremiseSideLower                         = "lower_boundary_reclaim"
	RangeRouterRotationPremiseSideUpper                         = "upper_boundary_reclaim"
	RangeRouterRotationPremiseOutcomeMidline                    = "midline_rotation_first"
	RangeRouterRotationPremiseOutcomeOppositeInnerQuartile      = "opposite_inner_quartile_first"
	RangeRouterRotationPremiseOutcomeBoundaryFailure            = "boundary_failure_first"
	RangeRouterRotationPremiseOutcomeCleanExpansionAgainst      = "clean_expansion_against_rotation"
	RangeRouterRotationPremiseOutcomeBoundaryChopNoRotation     = "boundary_chop_no_rotation"
	RangeRouterRotationPremiseOutcomeNoResolution               = "no_resolution"
	RangeRouterRotationPremiseOutcomeMissingFuture              = "missing_future"
	rangeRouterRotationPremiseSkipDuplicateContext              = "duplicate_router_context_row"
	rangeRouterRotationPremiseSkipIneligibleRouter              = "ineligible_router_context"
	rangeRouterRotationPremiseSkipMissingState                  = "missing_state_row"
	rangeRouterRotationPremiseSkipNonPositiveWidth              = "non_positive_frozen_range_width"
	rangeRouterRotationPremiseSkipNoBoundaryReclaim             = "no_boundary_reclaim_event"
	rangeRouterRotationPremiseSkipMissingFuture                 = "missing_future"
	rangeRouterRotationPremiseSkipRouterDependency              = "router_dependency_gap"
	rangeRouterRotationPremiseSkipClosedFamilyReslice           = "closed_family_reslice"
)

type FuturesRangeRouterRotationPremiseAuditConfig struct {
	RouterAuditConfig                 FuturesRangeContextRouterAuditConfig
	RequiredRouterStopState           string
	Timeframe                         string
	HorizonBars                       int
	EventSearchBars                   int
	BoundaryZoneWidthFraction         float64
	CleanExpansionWidthFraction       float64
	MinContextSegmentsFull            int
	MinValidEventsFull                int
	MinValidEventsPerSplit            int
	MinValidEventsPerSide             int
	MaxSplitContributionRate          float64
	RouterFullExpectedHitRate         float64
	RouterWeakestSplitExpectedRate    float64
	MinFullMidlineImprovement         float64
	MinWeakestSplitMidlineImprovement float64
	MaxFullHardAdverseRate            float64
	MaxWorstSplitHardAdverseRate      float64
	MaxFullChopNoResolutionRate       float64
	MaxDominantStateIDRate            float64
	AllowWeakSideDiagnostic           bool
	SkipCoverageCountCheck            bool
	RejectAsStrategyBacktestRequest   bool
}

type FuturesRangeRouterRotationPremiseAuditResult struct {
	SourceRows      []FuturesRangeRouterRotationPremiseSourceRow           `json:"source_rows"`
	CoverageRows    []FuturesRangeRouterRotationPremiseCoverageRow         `json:"coverage_rows"`
	DependencyRows  []FuturesRangeRouterRotationPremiseRouterDependencyRow `json:"router_dependency_rows"`
	SegmentRows     []FuturesRangeRouterRotationPremiseContextSegmentRow   `json:"context_segment_rows"`
	EventRows       []FuturesRangeRouterRotationPremiseEventRow            `json:"event_rows"`
	OutcomeRows     []FuturesRangeRouterRotationPremiseOutcomeRow          `json:"outcome_rows"`
	CohortRows      []FuturesRangeRouterRotationPremiseCohortRow           `json:"cohort_rows"`
	RankingRows     []FuturesRangeRouterRotationPremiseRankingRow          `json:"ranking_rows"`
	SummaryRows     []FuturesRangeRouterRotationPremiseSummaryRow          `json:"summary_rows"`
	SkipRows        []FuturesRangeRouterRotationPremiseSkipRow             `json:"skip_rows"`
	PassingCohorts  int                                                    `json:"passing_cohorts"`
	StopState       string                                                 `json:"stop_state"`
	RouterStopState string                                                 `json:"router_stop_state"`
}

type FuturesRangeRouterRotationPremiseSourceRow struct {
	AuditName               string `json:"audit_name"`
	PremiseName             string `json:"premise_name"`
	RouterAuditName         string `json:"router_audit_name"`
	RouterStopState         string `json:"router_stop_state"`
	RequiredRouterStopState string `json:"required_router_stop_state"`
	SourceResamplePass      bool   `json:"source_resample_pass"`
	RouterDependencyPass    bool   `json:"router_dependency_pass"`
	Closed15MResamplePass   bool   `json:"closed_15m_resample_pass"`
	RouterUsesForwardLabels bool   `json:"router_uses_forward_labels"`
	ForwardLabelsAsInputs   bool   `json:"forward_labels_as_inputs"`
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

type FuturesRangeRouterRotationPremiseCoverageRow struct {
	AuditName             string `json:"audit_name"`
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
	CoverageFactsPass     bool   `json:"coverage_facts_pass"`
	ValidationStatus      string `json:"validation_status"`
	ValidationError       string `json:"validation_error,omitempty"`
	RouterRows            int    `json:"router_rows"`
	ContextSegments       int    `json:"context_segments"`
	EventRows             int    `json:"event_rows"`
	OutcomeRows           int    `json:"outcome_rows"`
	SourceResamplePass    bool   `json:"source_resample_pass"`
}

type FuturesRangeRouterRotationPremiseRouterDependencyRow struct {
	AuditName                       string  `json:"audit_name"`
	RouterAuditName                 string  `json:"router_audit_name"`
	RequiredRouterStopState         string  `json:"required_router_stop_state"`
	ActualRouterStopState           string  `json:"actual_router_stop_state"`
	DependencyPass                  bool    `json:"dependency_pass"`
	RouterRuleRows                  int     `json:"router_rule_rows"`
	RouterRows                      int     `json:"router_rows"`
	RouterCohortRows                int     `json:"router_cohort_rows"`
	RouterRankingRows               int     `json:"router_ranking_rows"`
	RouterPassingCohorts            int     `json:"router_passing_cohorts"`
	RequiredCohortID                string  `json:"required_cohort_id"`
	RequiredCohortPassesGate        bool    `json:"required_cohort_passes_gate"`
	RequiredCohortFullRows          int     `json:"required_cohort_full_rows"`
	RequiredCohortWeakestSplitRows  int     `json:"required_cohort_weakest_split_rows"`
	RequiredCohortFullExpectedRate  float64 `json:"required_cohort_full_expected_rate"`
	RequiredCohortWeakestSplitRate  float64 `json:"required_cohort_weakest_split_rate"`
	RequiredCohortFullAdverseRate   float64 `json:"required_cohort_full_adverse_rate"`
	RequiredCohortWorstSplitAdverse float64 `json:"required_cohort_worst_split_adverse_rate"`
	RequiredCohortDominantLabel     string  `json:"required_cohort_dominant_forward_label"`
	RequiredCohortDominantLabelRate float64 `json:"required_cohort_dominant_forward_label_rate"`
	ForwardLabelsAsRouterInput      bool    `json:"forward_labels_as_router_input"`
	ForwardLabelColumnsPresent      bool    `json:"forward_label_columns_present"`
	FailureReason                   string  `json:"failure_reason,omitempty"`
}

type FuturesRangeRouterRotationPremiseContextSegmentRow struct {
	SegmentID                         int     `json:"segment_id"`
	Timeframe                         string  `json:"timeframe"`
	HorizonBars                       int     `json:"horizon_bars"`
	Split                             string  `json:"split"`
	RangeEpisodeID                    int     `json:"range_episode_id"`
	SegmentStartRouterRowID           int     `json:"segment_start_router_row_id"`
	SegmentEndRouterRowID             int     `json:"segment_end_router_row_id"`
	SegmentStartStateRowID            int     `json:"segment_start_state_row_id"`
	SegmentEndStateRowID              int     `json:"segment_end_state_row_id"`
	StartDecisionIndex                int     `json:"start_decision_index"`
	EndDecisionIndex                  int     `json:"end_decision_index"`
	SegmentStartCloseTime             string  `json:"segment_start_close_time"`
	SegmentEndCloseTime               string  `json:"segment_end_close_time"`
	RangeStartIndex                   int     `json:"range_start_index"`
	RangeStartTime                    string  `json:"range_start_time"`
	FrozenRangeHigh                   float64 `json:"frozen_range_high"`
	FrozenRangeLow                    float64 `json:"frozen_range_low"`
	FrozenRangeMid                    float64 `json:"frozen_range_mid"`
	FrozenUpperQuartile               float64 `json:"frozen_upper_quartile"`
	FrozenLowerQuartile               float64 `json:"frozen_lower_quartile"`
	FrozenRangeWidth                  float64 `json:"frozen_range_width"`
	BoundaryZoneWidth                 float64 `json:"boundary_zone_width"`
	FrozenBoundsKnownThroughIndex     int     `json:"frozen_bounds_known_through_index"`
	FrozenBoundsKnownThroughCloseTime string  `json:"frozen_bounds_known_through_close_time"`
	RouterLabel                       string  `json:"router_label"`
	MatchedRuleCount                  int     `json:"matched_rule_count"`
	MatchedRuleIDs                    string  `json:"matched_rule_ids"`
	StateID                           string  `json:"state_id"`
	AllFamiliesID                     string  `json:"all_families_id"`
	ContextRowCount                   int     `json:"context_row_count"`
	DuplicateRouterRowsSkipped        int     `json:"duplicate_router_rows_skipped"`
	EventSearchStartIndex             int     `json:"event_search_start_index"`
	EventSearchEndIndex               int     `json:"event_search_end_index"`
	EventSearchStartTime              string  `json:"event_search_start_time"`
	EventSearchEndTime                string  `json:"event_search_end_time"`
	EventFound                        bool    `json:"event_found"`
	EventID                           int     `json:"event_id,omitempty"`
	ClosedCandleOnly                  bool    `json:"closed_candle_only"`
	ForwardLabelsAsContextInput       bool    `json:"forward_labels_as_context_input"`
	ForwardLabelsAsEventInput         bool    `json:"forward_labels_as_event_input"`
	Eligible                          bool    `json:"eligible"`
	SkipReason                        string  `json:"skip_reason,omitempty"`
}

type FuturesRangeRouterRotationPremiseEventRow struct {
	EventID                      int     `json:"event_id"`
	SegmentID                    int     `json:"segment_id"`
	Timeframe                    string  `json:"timeframe"`
	HorizonBars                  int     `json:"horizon_bars"`
	Split                        string  `json:"split"`
	Side                         string  `json:"side"`
	RangeEpisodeID               int     `json:"range_episode_id"`
	ContextStartStateRowID       int     `json:"context_start_state_row_id"`
	ContextStartDecisionIndex    int     `json:"context_start_decision_index"`
	EventIndex                   int     `json:"event_index"`
	EventOpenTime                string  `json:"event_open_time"`
	EventCloseTime               string  `json:"event_close_time"`
	EventDelayBars               int     `json:"event_delay_bars"`
	FrozenRangeHigh              float64 `json:"frozen_range_high"`
	FrozenRangeLow               float64 `json:"frozen_range_low"`
	FrozenRangeMid               float64 `json:"frozen_range_mid"`
	FrozenUpperQuartile          float64 `json:"frozen_upper_quartile"`
	FrozenLowerQuartile          float64 `json:"frozen_lower_quartile"`
	FrozenRangeWidth             float64 `json:"frozen_range_width"`
	BoundaryZoneWidth            float64 `json:"boundary_zone_width"`
	EventOpen                    float64 `json:"event_open"`
	EventHigh                    float64 `json:"event_high"`
	EventLow                     float64 `json:"event_low"`
	EventClose                   float64 `json:"event_close"`
	EventClosePosition           float64 `json:"event_close_position"`
	EventInputWindowEndIndex     int     `json:"event_input_window_end_index"`
	EventInputWindowEndCloseTime string  `json:"event_input_window_end_close_time"`
	ClosedCandleOnly             bool    `json:"closed_candle_only"`
	ForwardLabelsAsEventInput    bool    `json:"forward_labels_as_event_input"`
	EventValid                   bool    `json:"event_valid"`
	RejectReason                 string  `json:"reject_reason,omitempty"`
	StateID                      string  `json:"state_id"`
	AllFamiliesID                string  `json:"all_families_id"`
}

type FuturesRangeRouterRotationPremiseOutcomeRow struct {
	OutcomeID                     int    `json:"outcome_id"`
	EventID                       int    `json:"event_id"`
	SegmentID                     int    `json:"segment_id"`
	Timeframe                     string `json:"timeframe"`
	HorizonBars                   int    `json:"horizon_bars"`
	Split                         string `json:"split"`
	Side                          string `json:"side"`
	LabelWindowStartIndex         int    `json:"label_window_start_index"`
	LabelWindowEndIndex           int    `json:"label_window_end_index"`
	LabelWindowStartTime          string `json:"label_window_start_time"`
	LabelWindowEndTime            string `json:"label_window_end_time"`
	PrimaryOutcome                string `json:"primary_outcome"`
	MidlineRotationFirst          bool   `json:"midline_rotation_first"`
	OppositeInnerQuartileFirst    bool   `json:"opposite_inner_quartile_first"`
	BoundaryFailureFirst          bool   `json:"boundary_failure_first"`
	CleanExpansionAgainstRotation bool   `json:"clean_expansion_against_rotation"`
	BoundaryChopNoRotation        bool   `json:"boundary_chop_no_rotation"`
	NoResolution                  bool   `json:"no_resolution"`
	HardAdverse                   bool   `json:"hard_adverse"`
	BarsToOutcome                 int    `json:"bars_to_outcome"`
	SameBoundaryProbeCount        int    `json:"same_boundary_probe_count"`
	ConsecutiveOutsideCloseCount  int    `json:"consecutive_outside_close_count"`
	OutcomeComplete               bool   `json:"outcome_complete"`
	MissingFuture                 bool   `json:"missing_future"`
	FutureLabelsUsedAsInput       bool   `json:"future_labels_used_as_input"`
	StateID                       string `json:"state_id"`
	AllFamiliesID                 string `json:"all_families_id"`
}

type FuturesRangeRouterRotationPremiseCohortRow struct {
	CohortID                           string  `json:"cohort_id"`
	Split                              string  `json:"split"`
	Side                               string  `json:"side"`
	Timeframe                          string  `json:"timeframe"`
	HorizonBars                        int     `json:"horizon_bars"`
	ContextSegmentCount                int     `json:"context_segment_count"`
	EventCount                         int     `json:"event_count"`
	CompleteOutcomeCount               int     `json:"complete_outcome_count"`
	MissingFutureCount                 int     `json:"missing_future_count"`
	MidlineRotationCount               int     `json:"midline_rotation_count"`
	OppositeInnerQuartileCount         int     `json:"opposite_inner_quartile_count"`
	BoundaryFailureCount               int     `json:"boundary_failure_count"`
	CleanExpansionAgainstRotationCount int     `json:"clean_expansion_against_rotation_count"`
	BoundaryChopNoRotationCount        int     `json:"boundary_chop_no_rotation_count"`
	NoResolutionCount                  int     `json:"no_resolution_count"`
	HardAdverseCount                   int     `json:"hard_adverse_count"`
	ChopNoResolutionCount              int     `json:"chop_no_resolution_count"`
	MidlineRotationRate                float64 `json:"midline_rotation_rate"`
	OppositeInnerQuartileRate          float64 `json:"opposite_inner_quartile_rate"`
	HardAdverseRate                    float64 `json:"hard_adverse_rate"`
	ChopNoResolutionRate               float64 `json:"chop_no_resolution_rate"`
	DominantOutcome                    string  `json:"dominant_outcome"`
	DominantOutcomeRate                float64 `json:"dominant_outcome_rate"`
	DominantStateID                    string  `json:"dominant_state_id"`
	DominantStateIDRate                float64 `json:"dominant_state_id_rate"`
	FullContextSegments                int     `json:"full_context_segments"`
	FullValidEvents                    int     `json:"full_valid_events"`
	WeakestSplitEvents                 int     `json:"weakest_split_events"`
	LowerSideFullEvents                int     `json:"lower_side_full_events"`
	UpperSideFullEvents                int     `json:"upper_side_full_events"`
	WeakerSideDiagnostic               bool    `json:"weaker_side_diagnostic"`
	SideSymmetricLaterPremiseAllowed   bool    `json:"side_symmetric_later_premise_allowed"`
	MaxSplitContributionRate           float64 `json:"max_split_contribution_rate"`
	FullMidlineRotationRate            float64 `json:"full_midline_rotation_rate"`
	WeakestSplitMidlineRotationRate    float64 `json:"weakest_split_midline_rotation_rate"`
	FullHardAdverseRate                float64 `json:"full_hard_adverse_rate"`
	WorstSplitHardAdverseRate          float64 `json:"worst_split_hard_adverse_rate"`
	FullChopNoResolutionRate           float64 `json:"full_chop_no_resolution_rate"`
	SourceRouterGatePass               bool    `json:"source_router_gate_pass"`
	ContextCountGatePass               bool    `json:"context_count_gate_pass"`
	EventCountGatePass                 bool    `json:"event_count_gate_pass"`
	SplitCountGatePass                 bool    `json:"split_count_gate_pass"`
	SideCountGatePass                  bool    `json:"side_count_gate_pass"`
	SplitContributionGatePass          bool    `json:"split_contribution_gate_pass"`
	BehaviorGatePass                   bool    `json:"behavior_gate_pass"`
	StateRollupGatePass                bool    `json:"state_rollup_gate_pass"`
	ClosedFamilyProtectionPass         bool    `json:"closed_family_protection_pass"`
	FutureLeakProtectionPass           bool    `json:"future_leak_protection_pass"`
	PassesReviewGate                   bool    `json:"passes_review_gate"`
	FailureReason                      string  `json:"failure_reason,omitempty"`
}

type FuturesRangeRouterRotationPremiseRankingRow struct {
	Rank                             int     `json:"rank"`
	CohortID                         string  `json:"cohort_id"`
	Side                             string  `json:"side"`
	Timeframe                        string  `json:"timeframe"`
	HorizonBars                      int     `json:"horizon_bars"`
	PassesGate                       bool    `json:"passes_gate"`
	RankScore                        float64 `json:"rank_score"`
	FullContextSegments              int     `json:"full_context_segments"`
	FullValidEvents                  int     `json:"full_valid_events"`
	WeakestSplitEvents               int     `json:"weakest_split_events"`
	LowerSideFullEvents              int     `json:"lower_side_full_events"`
	UpperSideFullEvents              int     `json:"upper_side_full_events"`
	MaxSplitContributionRate         float64 `json:"max_split_contribution_rate"`
	FullMidlineRotationRate          float64 `json:"full_midline_rotation_rate"`
	WeakestSplitMidlineRotationRate  float64 `json:"weakest_split_midline_rotation_rate"`
	FullHardAdverseRate              float64 `json:"full_hard_adverse_rate"`
	WorstSplitHardAdverseRate        float64 `json:"worst_split_hard_adverse_rate"`
	FullChopNoResolutionRate         float64 `json:"full_chop_no_resolution_rate"`
	DominantOutcome                  string  `json:"dominant_outcome"`
	DominantOutcomeRate              float64 `json:"dominant_outcome_rate"`
	DominantStateID                  string  `json:"dominant_state_id"`
	DominantStateIDRate              float64 `json:"dominant_state_id_rate"`
	WeakerSideDiagnostic             bool    `json:"weaker_side_diagnostic"`
	SideSymmetricLaterPremiseAllowed bool    `json:"side_symmetric_later_premise_allowed"`
	FailureReason                    string  `json:"failure_reason,omitempty"`
}

type FuturesRangeRouterRotationPremiseSummaryRow struct {
	Split                    string `json:"split"`
	Timeframe                string `json:"timeframe"`
	HorizonBars              int    `json:"horizon_bars"`
	SourceResamplePass       bool   `json:"source_resample_pass"`
	RouterDependencyPass     bool   `json:"router_dependency_pass"`
	RouterStopState          string `json:"router_stop_state"`
	ContextSegments          int    `json:"context_segments"`
	Events                   int    `json:"events"`
	CompleteOutcomes         int    `json:"complete_outcomes"`
	MissingFutureOutcomes    int    `json:"missing_future_outcomes"`
	LowerEvents              int    `json:"lower_events"`
	UpperEvents              int    `json:"upper_events"`
	MidlineRotationOutcomes  int    `json:"midline_rotation_outcomes"`
	HardAdverseOutcomes      int    `json:"hard_adverse_outcomes"`
	ChopNoResolutionOutcomes int    `json:"chop_no_resolution_outcomes"`
	CohortRows               int    `json:"cohort_rows"`
	RankingRows              int    `json:"ranking_rows"`
	PassingCohorts           int    `json:"passing_cohorts"`
	SkipRows                 int    `json:"skip_rows"`
	StopState                string `json:"stop_state"`
}

type FuturesRangeRouterRotationPremiseSkipRow struct {
	Timeframe string `json:"timeframe"`
	Split     string `json:"split"`
	Reason    string `json:"reason"`
	Count     int    `json:"count"`
}

type rangeRouterRotationPremiseCohortKey struct {
	split       string
	side        string
	timeframe   string
	horizonBars int
}

type rangeRouterRotationPremiseCohortAccumulator struct {
	row      FuturesRangeRouterRotationPremiseCohortRow
	outcomes map[string]int
	stateIDs map[string]int
}

func DefaultFuturesRangeRouterRotationPremiseAuditConfig() FuturesRangeRouterRotationPremiseAuditConfig {
	return FuturesRangeRouterRotationPremiseAuditConfig{
		RouterAuditConfig:                 DefaultFuturesRangeContextRouterAuditConfig(),
		RequiredRouterStopState:           RangeContextRouterStopStatePassedNeedsRotationSpec,
		Timeframe:                         RangeDiscoveryTimeframe15m,
		HorizonBars:                       24,
		EventSearchBars:                   6,
		BoundaryZoneWidthFraction:         0.20,
		CleanExpansionWidthFraction:       0.15,
		MinContextSegmentsFull:            250,
		MinValidEventsFull:                150,
		MinValidEventsPerSplit:            40,
		MinValidEventsPerSide:             40,
		MaxSplitContributionRate:          0.45,
		RouterFullExpectedHitRate:         0.599776,
		RouterWeakestSplitExpectedRate:    0.554667,
		MinFullMidlineImprovement:         0.03,
		MinWeakestSplitMidlineImprovement: 0.02,
		MaxFullHardAdverseRate:            0.22,
		MaxWorstSplitHardAdverseRate:      0.28,
		MaxFullChopNoResolutionRate:       0.30,
		MaxDominantStateIDRate:            0.80,
		AllowWeakSideDiagnostic:           true,
	}
}

func RunFuturesRangeRouterRotationPremiseAudit(candles []Candle, manifest SourceManifest, cfg FuturesRangeRouterRotationPremiseAuditConfig, splits []Split) (FuturesRangeRouterRotationPremiseAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeRouterRotationPremiseAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesRangeRouterRotationPremiseAuditResult{}
	if cfg.RejectAsStrategyBacktestRequest {
		result.StopState = RangeRouterRotationPremiseStopStateRejectedStrategyBacktest
		result.SummaryRows = rangeRouterRotationPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	routerResult, err := RunFuturesRangeContextRouterAudit(candles, manifest, cfg.RouterAuditConfig, splits)
	if err != nil {
		return result, err
	}
	result.RouterStopState = routerResult.StopState

	frame, ok := rangeContextTriageFrameDef(RangeDiscoveryTimeframe15m)
	if !ok {
		return result, fmt.Errorf("range router rotation premise missing frame definition for 15m")
	}
	frameCandles, coverage, resampleErr := resampleRangeDiscoveryFrame(candles, frame)
	coverageRow := rangeRouterRotationPremiseCoverageRow(coverage, cfg)
	resamplePass := resampleErr == nil && coverageRow.CoverageFactsPass && coverageRow.ValidationStatus == "accepted" && coverageRow.Complete
	dependency := rangeRouterRotationPremiseDependencyRow(routerResult, cfg)
	parentResamplePass := rangeContextRouterSourcePass(routerResult.SourceRows, routerResult.CoverageRows) && resamplePass
	sourcePass := parentResamplePass && dependency.DependencyPass
	result.SourceRows = rangeRouterRotationPremiseSourceRows(routerResult.SourceRows, routerResult.StopState, cfg, parentResamplePass, resamplePass, dependency.DependencyPass)
	result.CoverageRows = []FuturesRangeRouterRotationPremiseCoverageRow{coverageRow}
	result.DependencyRows = []FuturesRangeRouterRotationPremiseRouterDependencyRow{dependency}

	if routerResult.StopState == RangeContextRouterStopStateRejectedClosedReslice {
		rangeRouterRotationPremiseAddSkipRows(&result, cfg.Timeframe, fullSplitName, rangeRouterRotationPremiseSkipClosedFamilyReslice, 1)
		result.StopState = RangeRouterRotationPremiseStopStateRejectedClosedReslice
		result.SummaryRows = rangeRouterRotationPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}
	if !sourcePass || routerResult.StopState != cfg.RequiredRouterStopState {
		rangeRouterRotationPremiseAddSkipRows(&result, cfg.Timeframe, fullSplitName, rangeRouterRotationPremiseSkipRouterDependency, 1)
		result.StopState = RangeRouterRotationPremiseStopStateSourceRouterGap
		result.SummaryRows = rangeRouterRotationPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	stateByID := map[int]FuturesRangeStateConstructionLoopStateRow{}
	for _, state := range routerResult.StateRows {
		stateByID[state.StateRowID] = state
	}
	skipAcc := map[rangeStateSkipKey]int{}
	result.SegmentRows, result.EventRows, result.OutcomeRows = rangeRouterRotationPremiseBuildSegmentsEvents(routerResult.Rows, stateByID, frameCandles, cfg, splits, skipAcc)
	result.SkipRows = rangeRouterRotationPremiseSkipRows(skipAcc)
	result.CohortRows = rangeRouterRotationPremiseCohortRows(result.SegmentRows, result.OutcomeRows, cfg, splits, sourcePass)
	result.RankingRows = rangeRouterRotationPremiseRankingRows(result.CohortRows)
	result.PassingCohorts = rangeRouterRotationPremisePassingRankingCount(result.RankingRows)
	result.StopState = FuturesRangeRouterRotationPremiseAuditStopState(result)
	result.SummaryRows = rangeRouterRotationPremiseSummaryRows(result, cfg, splits)
	for i := range result.CoverageRows {
		result.CoverageRows[i].RouterRows = len(routerResult.Rows)
		result.CoverageRows[i].ContextSegments = len(result.SegmentRows)
		result.CoverageRows[i].EventRows = len(result.EventRows)
		result.CoverageRows[i].OutcomeRows = len(result.OutcomeRows)
		result.CoverageRows[i].SourceResamplePass = rangeRouterRotationPremiseSourcePass(result.SourceRows, result.CoverageRows, result.DependencyRows)
	}
	return result, nil
}

func FuturesRangeRouterRotationPremiseAuditStopState(result FuturesRangeRouterRotationPremiseAuditResult) string {
	if result.StopState == RangeRouterRotationPremiseStopStateRejectedStrategyBacktest {
		return result.StopState
	}
	if result.RouterStopState == RangeContextRouterStopStateRejectedClosedReslice {
		return RangeRouterRotationPremiseStopStateRejectedClosedReslice
	}
	if !rangeRouterRotationPremiseSourcePass(result.SourceRows, result.CoverageRows, result.DependencyRows) {
		return RangeRouterRotationPremiseStopStateSourceRouterGap
	}
	complete := 0
	for _, outcome := range result.OutcomeRows {
		if outcome.OutcomeComplete {
			complete++
		}
	}
	if complete == 0 {
		return RangeRouterRotationPremiseStopStateNoEligibleEvents
	}
	for _, ranking := range result.RankingRows {
		if ranking.PassesGate {
			return RangeRouterRotationPremiseStopStatePassedNeedsTriggerAudit
		}
	}
	return RangeRouterRotationPremiseStopStateFailedNoPremise
}

func (cfg FuturesRangeRouterRotationPremiseAuditConfig) withDefaults() FuturesRangeRouterRotationPremiseAuditConfig {
	defaults := DefaultFuturesRangeRouterRotationPremiseAuditConfig()
	cfg.RouterAuditConfig = cfg.RouterAuditConfig.withDefaults()
	if cfg.RequiredRouterStopState == "" {
		cfg.RequiredRouterStopState = defaults.RequiredRouterStopState
	}
	if cfg.Timeframe == "" {
		cfg.Timeframe = defaults.Timeframe
	}
	if cfg.HorizonBars == 0 {
		cfg.HorizonBars = defaults.HorizonBars
	}
	if cfg.EventSearchBars == 0 {
		cfg.EventSearchBars = defaults.EventSearchBars
	}
	if cfg.BoundaryZoneWidthFraction == 0 {
		cfg.BoundaryZoneWidthFraction = defaults.BoundaryZoneWidthFraction
	}
	if cfg.CleanExpansionWidthFraction == 0 {
		cfg.CleanExpansionWidthFraction = defaults.CleanExpansionWidthFraction
	}
	if cfg.MinContextSegmentsFull == 0 {
		cfg.MinContextSegmentsFull = defaults.MinContextSegmentsFull
	}
	if cfg.MinValidEventsFull == 0 {
		cfg.MinValidEventsFull = defaults.MinValidEventsFull
	}
	if cfg.MinValidEventsPerSplit == 0 {
		cfg.MinValidEventsPerSplit = defaults.MinValidEventsPerSplit
	}
	if cfg.MinValidEventsPerSide == 0 {
		cfg.MinValidEventsPerSide = defaults.MinValidEventsPerSide
	}
	if cfg.MaxSplitContributionRate == 0 {
		cfg.MaxSplitContributionRate = defaults.MaxSplitContributionRate
	}
	if cfg.RouterFullExpectedHitRate == 0 {
		cfg.RouterFullExpectedHitRate = defaults.RouterFullExpectedHitRate
	}
	if cfg.RouterWeakestSplitExpectedRate == 0 {
		cfg.RouterWeakestSplitExpectedRate = defaults.RouterWeakestSplitExpectedRate
	}
	if cfg.MinFullMidlineImprovement == 0 {
		cfg.MinFullMidlineImprovement = defaults.MinFullMidlineImprovement
	}
	if cfg.MinWeakestSplitMidlineImprovement == 0 {
		cfg.MinWeakestSplitMidlineImprovement = defaults.MinWeakestSplitMidlineImprovement
	}
	if cfg.MaxFullHardAdverseRate == 0 {
		cfg.MaxFullHardAdverseRate = defaults.MaxFullHardAdverseRate
	}
	if cfg.MaxWorstSplitHardAdverseRate == 0 {
		cfg.MaxWorstSplitHardAdverseRate = defaults.MaxWorstSplitHardAdverseRate
	}
	if cfg.MaxFullChopNoResolutionRate == 0 {
		cfg.MaxFullChopNoResolutionRate = defaults.MaxFullChopNoResolutionRate
	}
	if cfg.MaxDominantStateIDRate == 0 {
		cfg.MaxDominantStateIDRate = defaults.MaxDominantStateIDRate
	}
	if !cfg.AllowWeakSideDiagnostic {
		cfg.AllowWeakSideDiagnostic = defaults.AllowWeakSideDiagnostic
	}
	return cfg
}

func (cfg FuturesRangeRouterRotationPremiseAuditConfig) validate() error {
	if err := cfg.RouterAuditConfig.validate(); err != nil {
		return err
	}
	if cfg.Timeframe != RangeDiscoveryTimeframe15m {
		return fmt.Errorf("range router rotation premise audit supports only 15m timeframe")
	}
	if cfg.HorizonBars != 24 {
		return fmt.Errorf("range router rotation premise audit supports only h24")
	}
	if cfg.EventSearchBars <= 0 || cfg.BoundaryZoneWidthFraction <= 0 || cfg.CleanExpansionWidthFraction <= 0 {
		return fmt.Errorf("range router rotation premise event windows and thresholds must be positive")
	}
	if cfg.MaxSplitContributionRate <= 0 || cfg.MaxSplitContributionRate > 1 || cfg.MaxDominantStateIDRate <= 0 || cfg.MaxDominantStateIDRate > 1 {
		return fmt.Errorf("range router rotation premise contribution gates must be in (0,1]")
	}
	return nil
}

func rangeRouterRotationPremiseSourceRows(routerSources []FuturesRangeContextRouterSourceRow, routerStop string, cfg FuturesRangeRouterRotationPremiseAuditConfig, sourcePass, resamplePass, dependencyPass bool) []FuturesRangeRouterRotationPremiseSourceRow {
	rows := make([]FuturesRangeRouterRotationPremiseSourceRow, 0, len(routerSources))
	for _, source := range routerSources {
		rows = append(rows, FuturesRangeRouterRotationPremiseSourceRow{
			AuditName:               FuturesRangeRouterRotationPremiseAuditName,
			PremiseName:             RangeRouterRotationPremiseName,
			RouterAuditName:         FuturesRangeContextRouterAuditName,
			RouterStopState:         routerStop,
			RequiredRouterStopState: cfg.RequiredRouterStopState,
			SourceResamplePass:      sourcePass,
			RouterDependencyPass:    dependencyPass,
			Closed15MResamplePass:   resamplePass,
			RouterUsesForwardLabels: source.RouterUsesForwardLabels,
			ForwardLabelsAsInputs:   false,
			Path:                    source.Path,
			ApprovedPath:            source.ApprovedPath,
			Venue:                   source.Venue,
			Product:                 source.Product,
			Symbol:                  source.Symbol,
			Interval:                source.Interval,
			RowCount:                source.RowCount,
			ExpectedRowCount:        source.ExpectedRowCount,
			FirstOpenTime:           source.FirstOpenTime,
			ExpectedFirstOpenTime:   source.ExpectedFirstOpenTime,
			LastOpenTime:            source.LastOpenTime,
			ExpectedLastOpenTime:    source.ExpectedLastOpenTime,
			GapCount:                source.GapCount,
			ExpectedGapCount:        source.ExpectedGapCount,
			DuplicateCount:          source.DuplicateCount,
			ExpectedDuplicateCount:  source.ExpectedDuplicateCount,
			ZeroVolumeCount:         source.ZeroVolumeCount,
			ExpectedZeroVolumeCount: source.ExpectedZeroVolumeCount,
			ComparisonOnly:          source.ComparisonOnly,
			SourceFactsPass:         source.SourceFactsPass,
			ValidationStatus:        source.ValidationStatus,
			ValidationError:         source.ValidationError,
		})
	}
	return rows
}

func rangeRouterRotationPremiseCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesRangeRouterRotationPremiseAuditConfig) FuturesRangeRouterRotationPremiseCoverageRow {
	row := FuturesRangeRouterRotationPremiseCoverageRow{
		AuditName:             FuturesRangeRouterRotationPremiseAuditName,
		Timeframe:             base.Timeframe,
		IntervalMinutes:       base.IntervalMinutes,
		ChildBars:             base.ChildBars,
		BarsPerDay:            base.BarsPerDay,
		RowCount:              base.RowCount,
		ExpectedRowCount:      rangeStateConstructionLoopExpected15MRows,
		FirstOpenTime:         base.FirstOpenTime,
		LastOpenTime:          base.LastOpenTime,
		ExpectedLastOpenTime:  rangeStateConstructionLoopExpected15MLast,
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
		ValidationStatus:      base.ValidationStatus,
		ValidationError:       base.ValidationError,
		SourceResamplePass:    base.ValidationStatus == "accepted" && base.Complete,
	}
	row.CoverageFactsPass = row.Timeframe == RangeDiscoveryTimeframe15m &&
		row.ValidationStatus == "accepted" &&
		row.Complete &&
		row.GapCount == 0 &&
		row.DuplicateCount == 0 &&
		row.MissingChildOpenCount == 0
	if !cfg.SkipCoverageCountCheck {
		row.CoverageFactsPass = row.CoverageFactsPass &&
			row.RowCount == row.ExpectedRowCount &&
			row.LastOpenTime == row.ExpectedLastOpenTime
	}
	return row
}

func rangeRouterRotationPremiseDependencyRow(routerResult FuturesRangeContextRouterAuditResult, cfg FuturesRangeRouterRotationPremiseAuditConfig) FuturesRangeRouterRotationPremiseRouterDependencyRow {
	row := FuturesRangeRouterRotationPremiseRouterDependencyRow{
		AuditName:               FuturesRangeRouterRotationPremiseAuditName,
		RouterAuditName:         FuturesRangeContextRouterAuditName,
		RequiredRouterStopState: cfg.RequiredRouterStopState,
		ActualRouterStopState:   routerResult.StopState,
		RouterRuleRows:          len(routerResult.RuleRows),
		RouterRows:              len(routerResult.Rows),
		RouterCohortRows:        len(routerResult.CohortRows),
		RouterRankingRows:       len(routerResult.RankingRows),
		RouterPassingCohorts:    routerResult.PassingCohorts,
		RequiredCohortID:        RangeRouterRotationPremiseRequiredRouterCohortID,
	}
	for _, ranking := range routerResult.RankingRows {
		if ranking.CohortID != RangeRouterRotationPremiseRequiredRouterCohortID {
			continue
		}
		row.RequiredCohortPassesGate = ranking.PassesGate
		row.RequiredCohortFullRows = ranking.FullPeriodRows
		row.RequiredCohortWeakestSplitRows = ranking.WeakestSplitRows
		row.RequiredCohortFullExpectedRate = ranking.FullExpectedRouteHitRate
		row.RequiredCohortWeakestSplitRate = ranking.WeakestSplitExpectedHitRate
		row.RequiredCohortFullAdverseRate = ranking.FullAdverseRouteHitRate
		row.RequiredCohortWorstSplitAdverse = ranking.WorstSplitAdverseHitRate
		row.RequiredCohortDominantLabel = ranking.DominantForwardLabel
		row.RequiredCohortDominantLabelRate = ranking.DominantForwardLabelRate
	}
	for _, router := range routerResult.Rows {
		if router.ForwardLabelsAsRouterInput {
			row.ForwardLabelsAsRouterInput = true
		}
		if router.ForwardLabelColumnsPresent {
			row.ForwardLabelColumnsPresent = true
		}
	}
	reasons := []string{}
	if routerResult.StopState != cfg.RequiredRouterStopState {
		reasons = append(reasons, "router_stop_state_mismatch")
	}
	if !row.RequiredCohortPassesGate {
		reasons = append(reasons, "required_rotation_router_cohort_missing_or_failed")
	}
	if row.ForwardLabelsAsRouterInput || row.ForwardLabelColumnsPresent {
		reasons = append(reasons, "router_forward_label_leak")
	}
	row.DependencyPass = len(reasons) == 0
	row.FailureReason = uniqueJoinedReasons(reasons)
	return row
}

func rangeRouterRotationPremiseBuildSegmentsEvents(routerRows []FuturesRangeContextRouterRow, stateByID map[int]FuturesRangeStateConstructionLoopStateRow, candles []Candle, cfg FuturesRangeRouterRotationPremiseAuditConfig, splits []Split, skipAcc map[rangeStateSkipKey]int) ([]FuturesRangeRouterRotationPremiseContextSegmentRow, []FuturesRangeRouterRotationPremiseEventRow, []FuturesRangeRouterRotationPremiseOutcomeRow) {
	rows := append([]FuturesRangeContextRouterRow(nil), routerRows...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RangeEpisodeID != rows[j].RangeEpisodeID {
			return rows[i].RangeEpisodeID < rows[j].RangeEpisodeID
		}
		if rows[i].DecisionIndex != rows[j].DecisionIndex {
			return rows[i].DecisionIndex < rows[j].DecisionIndex
		}
		return rows[i].RouterRowID < rows[j].RouterRowID
	})

	segments := []FuturesRangeRouterRotationPremiseContextSegmentRow{}
	events := []FuturesRangeRouterRotationPremiseEventRow{}
	outcomes := []FuturesRangeRouterRotationPremiseOutcomeRow{}
	var active *FuturesRangeRouterRotationPremiseContextSegmentRow
	var activeStartState FuturesRangeStateConstructionLoopStateRow
	finalize := func() {
		if active == nil {
			return
		}
		segment := *active
		event, found := rangeRouterRotationPremiseFindEvent(segment, candles)
		if found {
			event.EventID = len(events) + 1
			segment.EventFound = true
			segment.EventID = event.EventID
			events = append(events, event)
			outcomes = append(outcomes, rangeRouterRotationPremiseOutcomeRow(candles, event, len(outcomes)+1, cfg))
		} else {
			segment.SkipReason = rangeRouterRotationPremiseSkipNoBoundaryReclaim
			rangeRouterRotationPremiseAddSkip(skipAcc, segment.Timeframe, segment.Split, rangeRouterRotationPremiseSkipNoBoundaryReclaim)
		}
		segments = append(segments, segment)
		active = nil
		activeStartState = FuturesRangeStateConstructionLoopStateRow{}
	}

	for _, router := range rows {
		if !rangeRouterRotationPremiseEligibleRouterRow(router, cfg) {
			if router.Timeframe == cfg.Timeframe {
				rangeRouterRotationPremiseAddSkip(skipAcc, router.Timeframe, router.Split, rangeRouterRotationPremiseSkipIneligibleRouter)
			}
			continue
		}
		state, ok := stateByID[router.StateRowID]
		if !ok {
			rangeRouterRotationPremiseAddSkip(skipAcc, router.Timeframe, router.Split, rangeRouterRotationPremiseSkipMissingState)
			continue
		}
		if state.RangeWidth <= 0 {
			rangeRouterRotationPremiseAddSkip(skipAcc, router.Timeframe, router.Split, rangeRouterRotationPremiseSkipNonPositiveWidth)
			continue
		}
		if active != nil &&
			router.RangeEpisodeID == active.RangeEpisodeID &&
			router.DecisionIndex == active.EndDecisionIndex+1 {
			active.SegmentEndRouterRowID = router.RouterRowID
			active.SegmentEndStateRowID = router.StateRowID
			active.EndDecisionIndex = router.DecisionIndex
			active.SegmentEndCloseTime = router.DecisionCloseTime
			active.ContextRowCount++
			active.DuplicateRouterRowsSkipped++
			rangeRouterRotationPremiseAddSkip(skipAcc, router.Timeframe, router.Split, rangeRouterRotationPremiseSkipDuplicateContext)
			continue
		}
		finalize()
		activeStartState = state
		segment := rangeRouterRotationPremiseNewSegment(len(segments)+1, router, activeStartState, candles, cfg)
		active = &segment
	}
	finalize()
	return segments, events, outcomes
}

func rangeRouterRotationPremiseNewSegment(segmentID int, router FuturesRangeContextRouterRow, state FuturesRangeStateConstructionLoopStateRow, candles []Candle, cfg FuturesRangeRouterRotationPremiseAuditConfig) FuturesRangeRouterRotationPremiseContextSegmentRow {
	width := state.RangeWidth
	lowerQ := state.RangeLow + width*0.25
	upperQ := state.RangeLow + width*0.75
	searchStart := router.DecisionIndex + 1
	searchEnd := minInt(router.DecisionIndex+cfg.EventSearchBars, len(candles)-1)
	row := FuturesRangeRouterRotationPremiseContextSegmentRow{
		SegmentID:                         segmentID,
		Timeframe:                         cfg.Timeframe,
		HorizonBars:                       cfg.HorizonBars,
		Split:                             router.Split,
		RangeEpisodeID:                    router.RangeEpisodeID,
		SegmentStartRouterRowID:           router.RouterRowID,
		SegmentEndRouterRowID:             router.RouterRowID,
		SegmentStartStateRowID:            router.StateRowID,
		SegmentEndStateRowID:              router.StateRowID,
		StartDecisionIndex:                router.DecisionIndex,
		EndDecisionIndex:                  router.DecisionIndex,
		SegmentStartCloseTime:             router.DecisionCloseTime,
		SegmentEndCloseTime:               router.DecisionCloseTime,
		RangeStartIndex:                   state.RangeStartIndex,
		RangeStartTime:                    state.RangeStartTime,
		FrozenRangeHigh:                   state.RangeHigh,
		FrozenRangeLow:                    state.RangeLow,
		FrozenRangeMid:                    state.RangeMid,
		FrozenUpperQuartile:               upperQ,
		FrozenLowerQuartile:               lowerQ,
		FrozenRangeWidth:                  width,
		BoundaryZoneWidth:                 width * cfg.BoundaryZoneWidthFraction,
		FrozenBoundsKnownThroughIndex:     router.DecisionIndex,
		FrozenBoundsKnownThroughCloseTime: router.DecisionCloseTime,
		RouterLabel:                       router.RouterLabel,
		MatchedRuleCount:                  router.MatchedRuleCount,
		MatchedRuleIDs:                    router.MatchedRuleIDs,
		StateID:                           router.StateID,
		AllFamiliesID:                     state.AllFamiliesID,
		ContextRowCount:                   1,
		EventSearchStartIndex:             searchStart,
		EventSearchEndIndex:               searchEnd,
		ClosedCandleOnly:                  router.ClosedCandleOnly,
		ForwardLabelsAsContextInput:       false,
		ForwardLabelsAsEventInput:         false,
		Eligible:                          true,
	}
	if searchStart >= 0 && searchStart < len(candles) {
		row.EventSearchStartTime = candles[searchStart].CloseTime.UTC().Format(timeLayout)
	}
	if searchEnd >= 0 && searchEnd < len(candles) && searchEnd >= searchStart {
		row.EventSearchEndTime = candles[searchEnd].CloseTime.UTC().Format(timeLayout)
	}
	return row
}

func rangeRouterRotationPremiseFindEvent(segment FuturesRangeRouterRotationPremiseContextSegmentRow, candles []Candle) (FuturesRangeRouterRotationPremiseEventRow, bool) {
	if segment.EventSearchStartIndex >= len(candles) || segment.EventSearchEndIndex < segment.EventSearchStartIndex || segment.FrozenRangeWidth <= 0 {
		return FuturesRangeRouterRotationPremiseEventRow{}, false
	}
	for idx := segment.EventSearchStartIndex; idx <= segment.EventSearchEndIndex && idx < len(candles); idx++ {
		candle := candles[idx]
		if event, ok := rangeRouterRotationPremiseEventForCandle(segment, candle, idx, RangeRouterRotationPremiseSideLower); ok {
			return event, true
		}
		if event, ok := rangeRouterRotationPremiseEventForCandle(segment, candle, idx, RangeRouterRotationPremiseSideUpper); ok {
			return event, true
		}
	}
	return FuturesRangeRouterRotationPremiseEventRow{}, false
}

func rangeRouterRotationPremiseEventForCandle(segment FuturesRangeRouterRotationPremiseContextSegmentRow, candle Candle, eventIndex int, side string) (FuturesRangeRouterRotationPremiseEventRow, bool) {
	candleRange := candle.High - candle.Low
	if candleRange <= 0 || !validNumber(candleRange) {
		return FuturesRangeRouterRotationPremiseEventRow{}, false
	}
	if candle.Close < segment.FrozenRangeLow || candle.Close > segment.FrozenRangeHigh {
		return FuturesRangeRouterRotationPremiseEventRow{}, false
	}
	closePosition := (candle.Close - candle.Low) / candleRange
	lowerZoneTop := segment.FrozenRangeLow + segment.BoundaryZoneWidth
	upperZoneBottom := segment.FrozenRangeHigh - segment.BoundaryZoneWidth
	switch side {
	case RangeRouterRotationPremiseSideLower:
		if !(candle.Low <= lowerZoneTop &&
			candle.Close > lowerZoneTop &&
			candle.Close < segment.FrozenRangeMid &&
			candle.High < segment.FrozenRangeMid &&
			candle.Close > candle.Open &&
			closePosition >= 0.60) {
			return FuturesRangeRouterRotationPremiseEventRow{}, false
		}
	case RangeRouterRotationPremiseSideUpper:
		if !(candle.High >= upperZoneBottom &&
			candle.Close < upperZoneBottom &&
			candle.Close > segment.FrozenRangeMid &&
			candle.Low > segment.FrozenRangeMid &&
			candle.Close < candle.Open &&
			closePosition <= 0.40) {
			return FuturesRangeRouterRotationPremiseEventRow{}, false
		}
	default:
		return FuturesRangeRouterRotationPremiseEventRow{}, false
	}
	return FuturesRangeRouterRotationPremiseEventRow{
		SegmentID:                    segment.SegmentID,
		Timeframe:                    segment.Timeframe,
		HorizonBars:                  segment.HorizonBars,
		Split:                        splitNameForCloseTime(candle.CloseTime, DefaultSplits()),
		Side:                         side,
		RangeEpisodeID:               segment.RangeEpisodeID,
		ContextStartStateRowID:       segment.SegmentStartStateRowID,
		ContextStartDecisionIndex:    segment.StartDecisionIndex,
		EventIndex:                   eventIndex,
		EventOpenTime:                candle.OpenTime.UTC().Format(timeLayout),
		EventCloseTime:               candle.CloseTime.UTC().Format(timeLayout),
		EventDelayBars:               eventIndex - segment.StartDecisionIndex,
		FrozenRangeHigh:              segment.FrozenRangeHigh,
		FrozenRangeLow:               segment.FrozenRangeLow,
		FrozenRangeMid:               segment.FrozenRangeMid,
		FrozenUpperQuartile:          segment.FrozenUpperQuartile,
		FrozenLowerQuartile:          segment.FrozenLowerQuartile,
		FrozenRangeWidth:             segment.FrozenRangeWidth,
		BoundaryZoneWidth:            segment.BoundaryZoneWidth,
		EventOpen:                    candle.Open,
		EventHigh:                    candle.High,
		EventLow:                     candle.Low,
		EventClose:                   candle.Close,
		EventClosePosition:           closePosition,
		EventInputWindowEndIndex:     eventIndex,
		EventInputWindowEndCloseTime: candle.CloseTime.UTC().Format(timeLayout),
		ClosedCandleOnly:             true,
		ForwardLabelsAsEventInput:    false,
		EventValid:                   true,
		StateID:                      segment.StateID,
		AllFamiliesID:                segment.AllFamiliesID,
	}, true
}

func rangeRouterRotationPremiseOutcomeRow(candles []Candle, event FuturesRangeRouterRotationPremiseEventRow, outcomeID int, cfg FuturesRangeRouterRotationPremiseAuditConfig) FuturesRangeRouterRotationPremiseOutcomeRow {
	row := FuturesRangeRouterRotationPremiseOutcomeRow{
		OutcomeID:               outcomeID,
		EventID:                 event.EventID,
		SegmentID:               event.SegmentID,
		Timeframe:               event.Timeframe,
		HorizonBars:             cfg.HorizonBars,
		Split:                   event.Split,
		Side:                    event.Side,
		LabelWindowStartIndex:   event.EventIndex + 1,
		LabelWindowEndIndex:     event.EventIndex + cfg.HorizonBars,
		PrimaryOutcome:          RangeRouterRotationPremiseOutcomeNoResolution,
		NoResolution:            true,
		BarsToOutcome:           -1,
		FutureLabelsUsedAsInput: false,
		StateID:                 event.StateID,
		AllFamiliesID:           event.AllFamiliesID,
	}
	if row.LabelWindowStartIndex < len(candles) {
		row.LabelWindowStartTime = candles[row.LabelWindowStartIndex].CloseTime.UTC().Format(timeLayout)
	}
	if row.LabelWindowEndIndex < len(candles) {
		row.LabelWindowEndTime = candles[row.LabelWindowEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	if row.LabelWindowEndIndex >= len(candles) || row.LabelWindowStartIndex >= len(candles) {
		row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeMissingFuture
		row.NoResolution = false
		row.OutcomeComplete = false
		row.MissingFuture = true
		return row
	}

	inf := int(^uint(0) >> 1)
	midIdx, oppIdx, boundaryIdx, cleanIdx, chopIdx := inf, inf, inf, inf, inf
	consecutiveOutside := 0
	probes := 0
	for idx := row.LabelWindowStartIndex; idx <= row.LabelWindowEndIndex; idx++ {
		candle := candles[idx]
		if rangeRouterRotationPremiseMidlineHit(event, candle) && midIdx == inf {
			midIdx = idx
		}
		if rangeRouterRotationPremiseOppositeQuartileHit(event, candle) && oppIdx == inf {
			oppIdx = idx
		}
		outsideClose := rangeRouterRotationPremiseOutsideCloseAgainst(event, candle)
		if outsideClose {
			consecutiveOutside++
			if boundaryIdx == inf {
				boundaryIdx = idx
			}
		} else {
			consecutiveOutside = 0
		}
		if (rangeRouterRotationPremiseCleanExpansionAgainst(event, candle, cfg) || consecutiveOutside >= 2) && cleanIdx == inf {
			cleanIdx = idx
		}
		if rangeRouterRotationPremiseSameBoundaryProbe(event, candle) {
			probes++
			if probes >= 3 && chopIdx == inf {
				chopIdx = idx
			}
		}
		if consecutiveOutside > row.ConsecutiveOutsideCloseCount {
			row.ConsecutiveOutsideCloseCount = consecutiveOutside
		}
		row.SameBoundaryProbeCount = probes
	}
	hardIdx := minInt(boundaryIdx, cleanIdx)
	row.OutcomeComplete = true
	switch {
	case hardIdx != inf && hardIdx <= midIdx:
		if cleanIdx == hardIdx {
			row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeCleanExpansionAgainst
			row.CleanExpansionAgainstRotation = true
		} else {
			row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeBoundaryFailure
			row.BoundaryFailureFirst = true
		}
		row.NoResolution = false
		row.HardAdverse = true
		row.BarsToOutcome = hardIdx - event.EventIndex
	case oppIdx != inf:
		row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeOppositeInnerQuartile
		row.MidlineRotationFirst = true
		row.OppositeInnerQuartileFirst = true
		row.NoResolution = false
		row.BarsToOutcome = oppIdx - event.EventIndex
	case midIdx != inf:
		row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeMidline
		row.MidlineRotationFirst = true
		row.NoResolution = false
		row.BarsToOutcome = midIdx - event.EventIndex
	case chopIdx != inf:
		row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeBoundaryChopNoRotation
		row.BoundaryChopNoRotation = true
		row.NoResolution = false
		row.BarsToOutcome = chopIdx - event.EventIndex
	default:
		row.PrimaryOutcome = RangeRouterRotationPremiseOutcomeNoResolution
		row.NoResolution = true
		row.BarsToOutcome = -1
	}
	return row
}

func rangeRouterRotationPremiseMidlineHit(event FuturesRangeRouterRotationPremiseEventRow, candle Candle) bool {
	if event.Side == RangeRouterRotationPremiseSideLower {
		return candle.High >= event.FrozenRangeMid
	}
	return candle.Low <= event.FrozenRangeMid
}

func rangeRouterRotationPremiseOppositeQuartileHit(event FuturesRangeRouterRotationPremiseEventRow, candle Candle) bool {
	if event.Side == RangeRouterRotationPremiseSideLower {
		return candle.High >= event.FrozenUpperQuartile
	}
	return candle.Low <= event.FrozenLowerQuartile
}

func rangeRouterRotationPremiseOutsideCloseAgainst(event FuturesRangeRouterRotationPremiseEventRow, candle Candle) bool {
	if event.Side == RangeRouterRotationPremiseSideLower {
		return candle.Close < event.FrozenRangeLow
	}
	return candle.Close > event.FrozenRangeHigh
}

func rangeRouterRotationPremiseCleanExpansionAgainst(event FuturesRangeRouterRotationPremiseEventRow, candle Candle, cfg FuturesRangeRouterRotationPremiseAuditConfig) bool {
	if event.Side == RangeRouterRotationPremiseSideLower {
		return candle.Close <= event.FrozenRangeLow-event.FrozenRangeWidth*cfg.CleanExpansionWidthFraction
	}
	return candle.Close >= event.FrozenRangeHigh+event.FrozenRangeWidth*cfg.CleanExpansionWidthFraction
}

func rangeRouterRotationPremiseSameBoundaryProbe(event FuturesRangeRouterRotationPremiseEventRow, candle Candle) bool {
	if event.Side == RangeRouterRotationPremiseSideLower {
		return candle.Low <= event.FrozenRangeLow+event.BoundaryZoneWidth
	}
	return candle.High >= event.FrozenRangeHigh-event.BoundaryZoneWidth
}

func rangeRouterRotationPremiseEligibleRouterRow(row FuturesRangeContextRouterRow, cfg FuturesRangeRouterRotationPremiseAuditConfig) bool {
	return row.RouterLabel == RangeContextRouterLabelTradableRotation &&
		row.Timeframe == cfg.Timeframe &&
		strings.Contains(row.MatchedRuleIDs, rangeRouterRotationPremiseRulePrefix) &&
		!row.ConflictingRuleMatch &&
		!row.MissingRuleMatch &&
		row.ClosedCandleOnly &&
		!row.ForwardLabelsAsRouterInput &&
		!row.ForwardLabelColumnsPresent
}

func rangeRouterRotationPremiseCohortRows(segments []FuturesRangeRouterRotationPremiseContextSegmentRow, outcomes []FuturesRangeRouterRotationPremiseOutcomeRow, cfg FuturesRangeRouterRotationPremiseAuditConfig, splits []Split, sourcePass bool) []FuturesRangeRouterRotationPremiseCohortRow {
	accs := map[rangeRouterRotationPremiseCohortKey]*rangeRouterRotationPremiseCohortAccumulator{}
	ensure := func(split, side string) *rangeRouterRotationPremiseCohortAccumulator {
		key := rangeRouterRotationPremiseCohortKey{split: split, side: side, timeframe: cfg.Timeframe, horizonBars: cfg.HorizonBars}
		acc := accs[key]
		if acc == nil {
			acc = &rangeRouterRotationPremiseCohortAccumulator{
				row: FuturesRangeRouterRotationPremiseCohortRow{
					CohortID:    rangeRouterRotationPremiseCohortID(key),
					Split:       split,
					Side:        side,
					Timeframe:   cfg.Timeframe,
					HorizonBars: cfg.HorizonBars,
				},
				outcomes: map[string]int{},
				stateIDs: map[string]int{},
			}
			accs[key] = acc
		}
		return acc
	}
	for _, segment := range segments {
		for _, split := range rangeContextRouterSplitNames(segment.Split) {
			ensure(split, RangeRouterRotationPremiseSideAll).row.ContextSegmentCount++
		}
	}
	for _, outcome := range outcomes {
		for _, split := range rangeContextRouterSplitNames(outcome.Split) {
			ensure(split, RangeRouterRotationPremiseSideAll).add(outcome)
			ensure(split, outcome.Side).add(outcome)
		}
	}
	rows := make([]FuturesRangeRouterRotationPremiseCohortRow, 0, len(accs))
	for _, acc := range accs {
		rows = append(rows, acc.finalRow())
	}
	rangeRouterRotationPremiseMarkCohortGates(rows, cfg, splits, sourcePass)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Split != rows[j].Split {
			return splitSortKey(rows[i].Split) < splitSortKey(rows[j].Split)
		}
		if rows[i].Side != rows[j].Side {
			return rangeRouterRotationPremiseSideSortKey(rows[i].Side) < rangeRouterRotationPremiseSideSortKey(rows[j].Side)
		}
		return rows[i].CohortID < rows[j].CohortID
	})
	return rows
}

func (acc *rangeRouterRotationPremiseCohortAccumulator) add(outcome FuturesRangeRouterRotationPremiseOutcomeRow) {
	acc.row.EventCount++
	if outcome.MissingFuture {
		acc.row.MissingFutureCount++
		return
	}
	if !outcome.OutcomeComplete {
		return
	}
	acc.row.CompleteOutcomeCount++
	acc.outcomes[outcome.PrimaryOutcome]++
	acc.stateIDs[outcome.AllFamiliesID]++
	if outcome.MidlineRotationFirst {
		acc.row.MidlineRotationCount++
	}
	if outcome.OppositeInnerQuartileFirst {
		acc.row.OppositeInnerQuartileCount++
	}
	if outcome.BoundaryFailureFirst {
		acc.row.BoundaryFailureCount++
	}
	if outcome.CleanExpansionAgainstRotation {
		acc.row.CleanExpansionAgainstRotationCount++
	}
	if outcome.BoundaryChopNoRotation {
		acc.row.BoundaryChopNoRotationCount++
	}
	if outcome.NoResolution {
		acc.row.NoResolutionCount++
	}
	if outcome.HardAdverse {
		acc.row.HardAdverseCount++
	}
	if outcome.BoundaryChopNoRotation || outcome.NoResolution {
		acc.row.ChopNoResolutionCount++
	}
}

func (acc *rangeRouterRotationPremiseCohortAccumulator) finalRow() FuturesRangeRouterRotationPremiseCohortRow {
	row := acc.row
	if row.CompleteOutcomeCount > 0 {
		row.MidlineRotationRate = float64(row.MidlineRotationCount) / float64(row.CompleteOutcomeCount)
		row.OppositeInnerQuartileRate = float64(row.OppositeInnerQuartileCount) / float64(row.CompleteOutcomeCount)
		row.HardAdverseRate = float64(row.HardAdverseCount) / float64(row.CompleteOutcomeCount)
		row.ChopNoResolutionRate = float64(row.ChopNoResolutionCount) / float64(row.CompleteOutcomeCount)
	}
	row.DominantOutcome, row.DominantOutcomeRate = rangeContextTriageDominantLabel(acc.outcomes, row.CompleteOutcomeCount, false)
	row.DominantStateID, row.DominantStateIDRate = rangeContextTriageDominantLabel(acc.stateIDs, row.CompleteOutcomeCount, false)
	return row
}

func rangeRouterRotationPremiseMarkCohortGates(rows []FuturesRangeRouterRotationPremiseCohortRow, cfg FuturesRangeRouterRotationPremiseAuditConfig, splits []Split, sourcePass bool) {
	bySplitSide := map[string]*FuturesRangeRouterRotationPremiseCohortRow{}
	for i := range rows {
		bySplitSide[rows[i].Split+"|"+rows[i].Side] = &rows[i]
	}
	full := bySplitSide[fullSplitName+"|"+RangeRouterRotationPremiseSideAll]
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	fullContextSegments := 0
	fullValidEvents := 0
	fullMidlineRate := 0.0
	fullHardAdverseRate := 0.0
	fullChopNoResolutionRate := 0.0
	fullDominantStateIDRate := 0.0
	if full != nil {
		fullContextSegments = full.ContextSegmentCount
		fullValidEvents = full.CompleteOutcomeCount
		fullMidlineRate = full.MidlineRotationRate
		fullHardAdverseRate = full.HardAdverseRate
		fullChopNoResolutionRate = full.ChopNoResolutionRate
		fullDominantStateIDRate = full.DominantStateIDRate
	}
	weakestSplitEvents := int(^uint(0) >> 1)
	weakestSplitMidline := math.Inf(1)
	worstSplitHardAdverse := 0.0
	maxContribution := 0.0
	splitCountPass := true
	for _, split := range periodSplits {
		row := bySplitSide[split.Name+"|"+RangeRouterRotationPremiseSideAll]
		if row == nil {
			splitCountPass = false
			weakestSplitEvents = 0
			weakestSplitMidline = 0
			continue
		}
		if row.CompleteOutcomeCount < cfg.MinValidEventsPerSplit {
			splitCountPass = false
		}
		if row.CompleteOutcomeCount < weakestSplitEvents {
			weakestSplitEvents = row.CompleteOutcomeCount
		}
		if row.MidlineRotationRate < weakestSplitMidline {
			weakestSplitMidline = row.MidlineRotationRate
		}
		if row.HardAdverseRate > worstSplitHardAdverse {
			worstSplitHardAdverse = row.HardAdverseRate
		}
		if fullValidEvents > 0 {
			maxContribution = math.Max(maxContribution, float64(row.CompleteOutcomeCount)/float64(fullValidEvents))
		}
	}
	if weakestSplitEvents == int(^uint(0)>>1) {
		weakestSplitEvents = 0
	}
	if math.IsInf(weakestSplitMidline, 1) {
		weakestSplitMidline = 0
	}
	lowerEvents, upperEvents := 0, 0
	if lower := bySplitSide[fullSplitName+"|"+RangeRouterRotationPremiseSideLower]; lower != nil {
		lowerEvents = lower.CompleteOutcomeCount
	}
	if upper := bySplitSide[fullSplitName+"|"+RangeRouterRotationPremiseSideUpper]; upper != nil {
		upperEvents = upper.CompleteOutcomeCount
	}
	weakerDiagnostic := cfg.AllowWeakSideDiagnostic && fullValidEvents >= cfg.MinValidEventsFull && (lowerEvents >= cfg.MinValidEventsPerSide || upperEvents >= cfg.MinValidEventsPerSide)
	sideCountPass := (lowerEvents >= cfg.MinValidEventsPerSide && upperEvents >= cfg.MinValidEventsPerSide) || weakerDiagnostic

	for i := range rows {
		reasons := []string{}
		rows[i].FullContextSegments = fullContextSegments
		rows[i].FullValidEvents = fullValidEvents
		rows[i].WeakestSplitEvents = weakestSplitEvents
		rows[i].LowerSideFullEvents = lowerEvents
		rows[i].UpperSideFullEvents = upperEvents
		rows[i].WeakerSideDiagnostic = weakerDiagnostic && !(lowerEvents >= cfg.MinValidEventsPerSide && upperEvents >= cfg.MinValidEventsPerSide)
		rows[i].SideSymmetricLaterPremiseAllowed = lowerEvents >= cfg.MinValidEventsPerSide && upperEvents >= cfg.MinValidEventsPerSide
		rows[i].MaxSplitContributionRate = maxContribution
		rows[i].FullMidlineRotationRate = fullMidlineRate
		rows[i].WeakestSplitMidlineRotationRate = weakestSplitMidline
		rows[i].FullHardAdverseRate = fullHardAdverseRate
		rows[i].WorstSplitHardAdverseRate = worstSplitHardAdverse
		rows[i].FullChopNoResolutionRate = fullChopNoResolutionRate
		rows[i].SourceRouterGatePass = sourcePass
		rows[i].ContextCountGatePass = fullContextSegments >= cfg.MinContextSegmentsFull
		rows[i].EventCountGatePass = fullValidEvents >= cfg.MinValidEventsFull
		rows[i].SplitCountGatePass = splitCountPass
		rows[i].SideCountGatePass = sideCountPass
		rows[i].SplitContributionGatePass = maxContribution <= cfg.MaxSplitContributionRate || len(periodSplits) == 0
		rows[i].BehaviorGatePass = fullMidlineRate >= cfg.RouterFullExpectedHitRate+cfg.MinFullMidlineImprovement &&
			weakestSplitMidline >= cfg.RouterWeakestSplitExpectedRate+cfg.MinWeakestSplitMidlineImprovement &&
			fullHardAdverseRate <= cfg.MaxFullHardAdverseRate &&
			worstSplitHardAdverse <= cfg.MaxWorstSplitHardAdverseRate &&
			fullChopNoResolutionRate <= cfg.MaxFullChopNoResolutionRate
		rows[i].StateRollupGatePass = fullDominantStateIDRate > 0 && fullDominantStateIDRate <= cfg.MaxDominantStateIDRate
		rows[i].ClosedFamilyProtectionPass = true
		rows[i].FutureLeakProtectionPass = true
		if !rows[i].SourceRouterGatePass {
			reasons = append(reasons, "source_router_gate_failed")
		}
		if !rows[i].ContextCountGatePass {
			reasons = append(reasons, "inadequate_context_segment_count")
		}
		if !rows[i].EventCountGatePass {
			reasons = append(reasons, "inadequate_event_count")
		}
		if !rows[i].SplitCountGatePass {
			reasons = append(reasons, "inadequate_split_event_count")
		}
		if !rows[i].SideCountGatePass {
			reasons = append(reasons, "inadequate_side_event_count")
		}
		if !rows[i].SplitContributionGatePass {
			reasons = append(reasons, "single_split_contribution_above_gate")
		}
		if !rows[i].BehaviorGatePass {
			reasons = append(reasons, "behavior_gate_failed")
		}
		if !rows[i].StateRollupGatePass {
			reasons = append(reasons, "single_state_rollup_carry")
		}
		rows[i].PassesReviewGate = rows[i].Split == fullSplitName &&
			rows[i].Side == RangeRouterRotationPremiseSideAll &&
			rows[i].SourceRouterGatePass &&
			rows[i].ContextCountGatePass &&
			rows[i].EventCountGatePass &&
			rows[i].SplitCountGatePass &&
			rows[i].SideCountGatePass &&
			rows[i].SplitContributionGatePass &&
			rows[i].BehaviorGatePass &&
			rows[i].StateRollupGatePass &&
			rows[i].ClosedFamilyProtectionPass &&
			rows[i].FutureLeakProtectionPass
		rows[i].FailureReason = uniqueJoinedReasons(reasons)
	}
}

func rangeRouterRotationPremiseRankingRows(cohorts []FuturesRangeRouterRotationPremiseCohortRow) []FuturesRangeRouterRotationPremiseRankingRow {
	rows := []FuturesRangeRouterRotationPremiseRankingRow{}
	for _, cohort := range cohorts {
		if cohort.Split != fullSplitName {
			continue
		}
		row := FuturesRangeRouterRotationPremiseRankingRow{
			CohortID:                         cohort.CohortID,
			Side:                             cohort.Side,
			Timeframe:                        cohort.Timeframe,
			HorizonBars:                      cohort.HorizonBars,
			PassesGate:                       cohort.PassesReviewGate,
			FullContextSegments:              cohort.FullContextSegments,
			FullValidEvents:                  cohort.FullValidEvents,
			WeakestSplitEvents:               cohort.WeakestSplitEvents,
			LowerSideFullEvents:              cohort.LowerSideFullEvents,
			UpperSideFullEvents:              cohort.UpperSideFullEvents,
			MaxSplitContributionRate:         cohort.MaxSplitContributionRate,
			FullMidlineRotationRate:          cohort.FullMidlineRotationRate,
			WeakestSplitMidlineRotationRate:  cohort.WeakestSplitMidlineRotationRate,
			FullHardAdverseRate:              cohort.FullHardAdverseRate,
			WorstSplitHardAdverseRate:        cohort.WorstSplitHardAdverseRate,
			FullChopNoResolutionRate:         cohort.FullChopNoResolutionRate,
			DominantOutcome:                  cohort.DominantOutcome,
			DominantOutcomeRate:              cohort.DominantOutcomeRate,
			DominantStateID:                  cohort.DominantStateID,
			DominantStateIDRate:              cohort.DominantStateIDRate,
			WeakerSideDiagnostic:             cohort.WeakerSideDiagnostic,
			SideSymmetricLaterPremiseAllowed: cohort.SideSymmetricLaterPremiseAllowed,
			FailureReason:                    cohort.FailureReason,
		}
		row.RankScore = rangeRouterRotationPremiseRankScore(row)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		return rangeRouterRotationPremiseSideSortKey(rows[i].Side) < rangeRouterRotationPremiseSideSortKey(rows[j].Side)
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func rangeRouterRotationPremiseSummaryRows(result FuturesRangeRouterRotationPremiseAuditResult, cfg FuturesRangeRouterRotationPremiseAuditConfig, splits []Split) []FuturesRangeRouterRotationPremiseSummaryRow {
	sourcePass := rangeRouterRotationPremiseSourcePass(result.SourceRows, result.CoverageRows, result.DependencyRows)
	rows := []FuturesRangeRouterRotationPremiseSummaryRow{rangeRouterRotationPremiseSummaryRowFor(result, fullSplitName, cfg.Timeframe, cfg.HorizonBars, sourcePass)}
	for _, split := range splits {
		if split.Name == fullSplitName {
			continue
		}
		rows = append(rows, rangeRouterRotationPremiseSummaryRowFor(result, split.Name, cfg.Timeframe, cfg.HorizonBars, sourcePass))
	}
	return rows
}

func rangeRouterRotationPremiseSummaryRowFor(result FuturesRangeRouterRotationPremiseAuditResult, split, timeframe string, horizon int, sourcePass bool) FuturesRangeRouterRotationPremiseSummaryRow {
	row := FuturesRangeRouterRotationPremiseSummaryRow{
		Split:                split,
		Timeframe:            timeframe,
		HorizonBars:          horizon,
		SourceResamplePass:   sourcePass,
		RouterDependencyPass: len(result.DependencyRows) > 0 && result.DependencyRows[0].DependencyPass,
		RouterStopState:      result.RouterStopState,
		StopState:            result.StopState,
	}
	for _, segment := range result.SegmentRows {
		if split == fullSplitName || segment.Split == split {
			row.ContextSegments++
		}
	}
	for _, event := range result.EventRows {
		if split != fullSplitName && event.Split != split {
			continue
		}
		row.Events++
		if event.Side == RangeRouterRotationPremiseSideLower {
			row.LowerEvents++
		}
		if event.Side == RangeRouterRotationPremiseSideUpper {
			row.UpperEvents++
		}
	}
	for _, outcome := range result.OutcomeRows {
		if split != fullSplitName && outcome.Split != split {
			continue
		}
		if outcome.OutcomeComplete {
			row.CompleteOutcomes++
		}
		if outcome.MissingFuture {
			row.MissingFutureOutcomes++
		}
		if outcome.MidlineRotationFirst {
			row.MidlineRotationOutcomes++
		}
		if outcome.HardAdverse {
			row.HardAdverseOutcomes++
		}
		if outcome.BoundaryChopNoRotation || outcome.NoResolution {
			row.ChopNoResolutionOutcomes++
		}
	}
	for _, cohort := range result.CohortRows {
		if cohort.Split == split {
			row.CohortRows++
		}
	}
	for _, ranking := range result.RankingRows {
		if split == fullSplitName && ranking.Timeframe == timeframe && ranking.HorizonBars == horizon {
			row.RankingRows++
		}
	}
	for _, ranking := range result.RankingRows {
		if ranking.PassesGate {
			row.PassingCohorts++
		}
	}
	for _, skip := range result.SkipRows {
		if skip.Split == split {
			row.SkipRows += skip.Count
		}
	}
	return row
}

func rangeRouterRotationPremiseSourcePass(sources []FuturesRangeRouterRotationPremiseSourceRow, coverage []FuturesRangeRouterRotationPremiseCoverageRow, dependency []FuturesRangeRouterRotationPremiseRouterDependencyRow) bool {
	if len(sources) == 0 || len(coverage) == 0 || len(dependency) == 0 {
		return false
	}
	for _, source := range sources {
		if !source.SourceResamplePass ||
			!source.SourceFactsPass ||
			source.ValidationStatus != "accepted" ||
			source.ComparisonOnly ||
			source.Product != "Binance USDT-M futures" ||
			source.Symbol != RangeUniverseSymbolBTCUSDT ||
			source.Interval != "5m" ||
			!sameCleanPath(source.Path, source.ApprovedPath) ||
			source.RowCount != source.ExpectedRowCount ||
			source.FirstOpenTime != source.ExpectedFirstOpenTime ||
			source.LastOpenTime != source.ExpectedLastOpenTime ||
			source.GapCount != source.ExpectedGapCount ||
			source.DuplicateCount != source.ExpectedDuplicateCount ||
			source.ZeroVolumeCount != source.ExpectedZeroVolumeCount ||
			source.RouterUsesForwardLabels ||
			source.ForwardLabelsAsInputs {
			return false
		}
	}
	for _, row := range coverage {
		if !row.SourceResamplePass || !row.CoverageFactsPass || row.ValidationStatus != "accepted" || !row.Complete {
			return false
		}
	}
	for _, row := range dependency {
		if !row.DependencyPass || row.ForwardLabelsAsRouterInput || row.ForwardLabelColumnsPresent {
			return false
		}
	}
	return true
}

func rangeRouterRotationPremiseSkipRows(skips map[rangeStateSkipKey]int) []FuturesRangeRouterRotationPremiseSkipRow {
	rows := make([]FuturesRangeRouterRotationPremiseSkipRow, 0, len(skips))
	for key, count := range skips {
		rows = append(rows, FuturesRangeRouterRotationPremiseSkipRow{
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

func rangeRouterRotationPremiseAddSkip(skips map[rangeStateSkipKey]int, timeframe, split, reason string) {
	if split == "" {
		split = fullSplitName
	}
	skips[rangeStateSkipKey{timeframe: timeframe, split: split, reason: reason}]++
	if split != fullSplitName {
		skips[rangeStateSkipKey{timeframe: timeframe, split: fullSplitName, reason: reason}]++
	}
}

func rangeRouterRotationPremiseAddSkipRows(result *FuturesRangeRouterRotationPremiseAuditResult, timeframe, split, reason string, count int) {
	if split == "" {
		split = fullSplitName
	}
	result.SkipRows = append(result.SkipRows, FuturesRangeRouterRotationPremiseSkipRow{Timeframe: timeframe, Split: split, Reason: reason, Count: count})
}

func rangeRouterRotationPremiseCohortID(key rangeRouterRotationPremiseCohortKey) string {
	return strings.Join([]string{RangeRouterRotationPremiseCohortID, key.split, key.side}, "|")
}

func rangeRouterRotationPremiseRankScore(row FuturesRangeRouterRotationPremiseRankingRow) float64 {
	score := row.FullMidlineRotationRate + row.WeakestSplitMidlineRotationRate - row.FullHardAdverseRate - row.WorstSplitHardAdverseRate - row.FullChopNoResolutionRate
	score += math.Min(float64(row.FullValidEvents), 1000) / 1000
	score += math.Min(float64(row.WeakestSplitEvents), 250) / 1000
	score -= row.MaxSplitContributionRate / 10
	score -= row.DominantStateIDRate / 10
	return score
}

func rangeRouterRotationPremisePassingRankingCount(rows []FuturesRangeRouterRotationPremiseRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func rangeRouterRotationPremiseSideSortKey(side string) int {
	switch side {
	case RangeRouterRotationPremiseSideAll:
		return 0
	case RangeRouterRotationPremiseSideLower:
		return 1
	case RangeRouterRotationPremiseSideUpper:
		return 2
	default:
		return 9
	}
}
