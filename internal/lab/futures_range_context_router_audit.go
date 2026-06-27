package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeContextRouterAuditName = "futures_range_context_router_audit"

	RangeContextRouterStopStateSourceGap               = "range_context_router_source_gap"
	RangeContextRouterStopStateFailedNoActionableRoute = "range_context_router_failed_no_actionable_route"
	RangeContextRouterStopStatePassedNoTradeFilterOnly = "range_context_router_passed_no_trade_filter_only"
	RangeContextRouterStopStatePassedNeedsRotationSpec = "range_context_router_passed_needs_rotation_premise_spec"
	RangeContextRouterStopStatePassedNeedsContinuation = "range_context_router_passed_needs_continuation_premise_spec"
	RangeContextRouterStopStateRejectedClosedReslice   = "range_context_router_rejected_closed_family_reslice"
	RangeContextRouterLabelNoTrade                     = "no_trade"
	RangeContextRouterLabelTradableRotation            = "tradable_rotation"
	RangeContextRouterLabelTrendContinuation           = "trend_continuation"
	RangeContextRouterLabelDiagnosticOnly              = "diagnostic_only"
	rangeContextRouterRuleSourcePassedStateRanking     = "range_state_construction_loop_passed_ranking"
	rangeContextRouterRuleStatusActive                 = "active"
	rangeContextRouterRuleStatusRejectedClosedFamily   = "rejected_closed_family_reslice"
	rangeContextRouterRuleStatusRejectedFutureLeak     = "rejected_future_label_feature"
	rangeContextRouterRowReasonNoRuleMatch             = "no_passing_rule_match"
	rangeContextRouterRowReasonMatchedRule             = "matched_passing_rule"
	rangeContextRouterRowReasonConflictingRules        = "conflicting_route_rules"
	rangeContextRouterRowReasonStateIneligible         = "state_row_ineligible"
	rangeContextRouterRowReasonMissingStateRollup      = "missing_state_rollup"
)

type FuturesRangeContextRouterAuditConfig struct {
	StateAuditConfig             FuturesRangeStateConstructionLoopAuditConfig
	MinFullRouterRows            int
	MinSplitRouterRows           int
	MinNoTradeFullRouterRows     int
	MinNoTradeSplitRouterRows    int
	MaxSplitContributionRate     float64
	MinPositiveExpectedRateFull  float64
	MinPositiveExpectedRateSplit float64
	MaxPositiveAdverseRateFull   float64
	MaxPositiveAdverseRateSplit  float64
	MinPositiveMarginFull        float64
	MinPositiveMarginSplit       float64
	MinNoTradeToxicRateFull      float64
	MinNoTradeToxicRateSplit     float64
}

type FuturesRangeContextRouterAuditResult struct {
	SourceRows     []FuturesRangeContextRouterSourceRow   `json:"source_rows"`
	CoverageRows   []FuturesRangeContextRouterCoverageRow `json:"coverage_rows"`
	RuleRows       []FuturesRangeContextRouterRuleRow     `json:"rule_rows"`
	Rows           []FuturesRangeContextRouterRow         `json:"rows"`
	CohortRows     []FuturesRangeContextRouterCohortRow   `json:"cohort_rows"`
	RankingRows    []FuturesRangeContextRouterRankingRow  `json:"ranking_rows"`
	SummaryRows    []FuturesRangeContextRouterSummaryRow  `json:"summary_rows"`
	SkipRows       []FuturesRangeContextRouterSkipRow     `json:"skip_rows"`
	PassingCohorts int                                    `json:"passing_cohorts"`
	StopState      string                                 `json:"stop_state"`
	StateStopState string                                 `json:"state_stop_state"`
}

type FuturesRangeContextRouterSourceRow struct {
	AuditName               string `json:"audit_name"`
	StateAuditName          string `json:"state_audit_name"`
	StateAuditStopState     string `json:"state_audit_stop_state"`
	SourceResamplePass      bool   `json:"source_resample_pass"`
	RouterUsesForwardLabels bool   `json:"router_uses_forward_labels"`
	FuturesRangeStateConstructionLoopSourceRow
}

type FuturesRangeContextRouterCoverageRow struct {
	AuditName          string `json:"audit_name"`
	StateRows          int    `json:"state_rows"`
	RouterRows         int    `json:"router_rows"`
	SourceResamplePass bool   `json:"source_resample_pass"`
	FuturesRangeStateConstructionLoopCoverageRow
}

type FuturesRangeContextRouterRuleRow struct {
	RuleID                     string  `json:"rule_id"`
	RuleSource                 string  `json:"rule_source"`
	RuleStatus                 string  `json:"rule_status"`
	RouterLabel                string  `json:"router_label"`
	SourceRouteCandidate       string  `json:"source_route_candidate"`
	SourceStateRank            int     `json:"source_state_rank"`
	SourceStateCohortID        string  `json:"source_state_cohort_id"`
	RollupType                 string  `json:"rollup_type"`
	RollupID                   string  `json:"rollup_id"`
	Timeframe                  string  `json:"timeframe"`
	HorizonBars                int     `json:"horizon_bars"`
	FullPeriodRows             int     `json:"full_period_rows"`
	WeakestSplitRows           int     `json:"weakest_split_rows"`
	FullExpectedRouteRate      float64 `json:"full_expected_route_rate"`
	WeakestSplitExpectedRate   float64 `json:"weakest_split_expected_rate"`
	FullAdverseRouteRate       float64 `json:"full_adverse_route_rate"`
	WorstSplitAdverseRate      float64 `json:"worst_split_adverse_rate"`
	FullExpectedMinusAdverse   float64 `json:"full_expected_minus_adverse_margin"`
	WeakestSplitMargin         float64 `json:"weakest_split_margin"`
	ClosedFamilyReslice        bool    `json:"closed_family_reslice"`
	FutureLeakProtectionPass   bool    `json:"future_leak_protection_pass"`
	ForwardLabelsAsRouterInput bool    `json:"forward_labels_as_router_input"`
	FeatureInputs              string  `json:"feature_inputs"`
	FailureReason              string  `json:"failure_reason,omitempty"`
}

type FuturesRangeContextRouterRow struct {
	RouterRowID                int    `json:"router_row_id"`
	StateRowID                 int    `json:"state_row_id"`
	Timestamp                  string `json:"timestamp"`
	Timeframe                  string `json:"timeframe"`
	Split                      string `json:"split"`
	RangeEpisodeID             int    `json:"range_episode_id"`
	RangeStartIndex            int    `json:"range_start_index"`
	DecisionIndex              int    `json:"decision_index"`
	RangeStartTime             string `json:"range_start_time"`
	DecisionCloseTime          string `json:"decision_close_time"`
	StateID                    string `json:"state_id"`
	GeometryVolID              string `json:"geometry_vol_id"`
	GeometryTrendID            string `json:"geometry_trend_id"`
	GeometryImpulseID          string `json:"geometry_impulse_id"`
	GeometryParticipationID    string `json:"geometry_participation_id"`
	GeometryVolTrendID         string `json:"geometry_vol_trend_id"`
	GeometryVolTrendImpulseID  string `json:"geometry_vol_trend_impulse_id"`
	AllFamiliesID              string `json:"all_families_id"`
	RouterLabel                string `json:"router_label"`
	RouterReason               string `json:"router_reason"`
	Actionable                 bool   `json:"actionable"`
	MatchedRuleCount           int    `json:"matched_rule_count"`
	MatchedRuleIDs             string `json:"matched_rule_ids,omitempty"`
	MatchedRouterLabels        string `json:"matched_router_labels,omitempty"`
	MatchedRollupTypes         string `json:"matched_rollup_types,omitempty"`
	ConflictingRuleMatch       bool   `json:"conflicting_rule_match"`
	MissingRuleMatch           bool   `json:"missing_rule_match"`
	ClosedCandleOnly           bool   `json:"closed_candle_only"`
	ForwardLabelsAsRouterInput bool   `json:"forward_labels_as_router_input"`
	ForwardLabelColumnsPresent bool   `json:"forward_label_columns_present"`
}

type FuturesRangeContextRouterCohortRow struct {
	CohortID                       string  `json:"cohort_id"`
	RouterLabel                    string  `json:"router_label"`
	Split                          string  `json:"split"`
	Timeframe                      string  `json:"timeframe"`
	HorizonBars                    int     `json:"horizon_bars"`
	CandidateCount                 int     `json:"candidate_count"`
	ExpectedRouteHitCount          int     `json:"expected_route_hit_count"`
	AdverseRouteHitCount           int     `json:"adverse_route_hit_count"`
	RotationUsefulCount            int     `json:"rotation_useful_count"`
	RotationToxicCount             int     `json:"rotation_toxic_count"`
	ContinuationUsefulCount        int     `json:"continuation_useful_count"`
	ContinuationToxicCount         int     `json:"continuation_toxic_count"`
	NoTradeToxicCount              int     `json:"no_trade_toxic_count"`
	ContainedRotationCount         int     `json:"contained_rotation_count"`
	CleanExpansionUpCount          int     `json:"clean_expansion_up_count"`
	CleanExpansionDownCount        int     `json:"clean_expansion_down_count"`
	FalseBreakReentryUpCount       int     `json:"false_break_reentry_up_count"`
	FalseBreakReentryDownCount     int     `json:"false_break_reentry_down_count"`
	BoundaryChopCount              int     `json:"boundary_chop_count"`
	DriftThroughUpCount            int     `json:"drift_through_up_count"`
	DriftThroughDownCount          int     `json:"drift_through_down_count"`
	LowWidthNoiseCount             int     `json:"low_width_noise_count"`
	NoResolutionCount              int     `json:"no_resolution_count"`
	ExpectedRouteHitRate           float64 `json:"expected_route_hit_rate"`
	AdverseRouteHitRate            float64 `json:"adverse_route_hit_rate"`
	ExpectedMinusAdverseMargin     float64 `json:"expected_minus_adverse_margin"`
	DominantForwardLabel           string  `json:"dominant_forward_label"`
	DominantForwardLabelRate       float64 `json:"dominant_forward_label_rate"`
	FullPeriodRows                 int     `json:"full_period_rows"`
	WeakestSplitRows               int     `json:"weakest_split_rows"`
	MaxSplitContributionRate       float64 `json:"max_split_contribution_rate"`
	FullExpectedRouteHitRate       float64 `json:"full_expected_route_hit_rate"`
	WeakestSplitExpectedHitRate    float64 `json:"weakest_split_expected_hit_rate"`
	FullAdverseRouteHitRate        float64 `json:"full_adverse_route_hit_rate"`
	WorstSplitAdverseHitRate       float64 `json:"worst_split_adverse_hit_rate"`
	FullExpectedMinusAdverseMargin float64 `json:"full_expected_minus_adverse_margin"`
	WeakestSplitMargin             float64 `json:"weakest_split_margin"`
	ReviewableCountGatePass        bool    `json:"reviewable_count_gate_pass"`
	SplitStabilityGatePass         bool    `json:"split_stability_gate_pass"`
	SplitContributionGatePass      bool    `json:"split_contribution_gate_pass"`
	RouteRateGatePass              bool    `json:"route_rate_gate_pass"`
	ClosedFamilyProtectionPass     bool    `json:"closed_family_protection_pass"`
	FutureLeakProtectionPass       bool    `json:"future_leak_protection_pass"`
	PassesReviewGate               bool    `json:"passes_review_gate"`
	FailureReason                  string  `json:"failure_reason,omitempty"`
}

type FuturesRangeContextRouterRankingRow struct {
	Rank                           int     `json:"rank"`
	CohortID                       string  `json:"cohort_id"`
	RouterLabel                    string  `json:"router_label"`
	Timeframe                      string  `json:"timeframe"`
	HorizonBars                    int     `json:"horizon_bars"`
	PassesGate                     bool    `json:"passes_gate"`
	RankScore                      float64 `json:"rank_score"`
	FullPeriodRows                 int     `json:"full_period_rows"`
	WeakestSplitRows               int     `json:"weakest_split_rows"`
	MaxSplitContributionRate       float64 `json:"max_split_contribution_rate"`
	FullExpectedRouteHitRate       float64 `json:"full_expected_route_hit_rate"`
	WeakestSplitExpectedHitRate    float64 `json:"weakest_split_expected_hit_rate"`
	FullAdverseRouteHitRate        float64 `json:"full_adverse_route_hit_rate"`
	WorstSplitAdverseHitRate       float64 `json:"worst_split_adverse_hit_rate"`
	FullExpectedMinusAdverseMargin float64 `json:"full_expected_minus_adverse_margin"`
	WeakestSplitMargin             float64 `json:"weakest_split_margin"`
	DominantForwardLabel           string  `json:"dominant_forward_label"`
	DominantForwardLabelRate       float64 `json:"dominant_forward_label_rate"`
	FailureReason                  string  `json:"failure_reason,omitempty"`
}

type FuturesRangeContextRouterSummaryRow struct {
	Split                 string `json:"split"`
	Timeframe             string `json:"timeframe"`
	HorizonBars           int    `json:"horizon_bars"`
	SourceResamplePass    bool   `json:"source_resample_pass"`
	StateAuditStopState   string `json:"state_audit_stop_state"`
	RuleRows              int    `json:"rule_rows"`
	RouterRows            int    `json:"router_rows"`
	NoTradeRows           int    `json:"no_trade_rows"`
	TradableRotationRows  int    `json:"tradable_rotation_rows"`
	TrendContinuationRows int    `json:"trend_continuation_rows"`
	DiagnosticOnlyRows    int    `json:"diagnostic_only_rows"`
	ConflictRows          int    `json:"conflict_rows"`
	MissingRuleRows       int    `json:"missing_rule_rows"`
	CohortRows            int    `json:"cohort_rows"`
	RankingRows           int    `json:"ranking_rows"`
	PassingCohorts        int    `json:"passing_cohorts"`
	SkipRows              int    `json:"skip_rows"`
	StopState             string `json:"stop_state"`
}

type FuturesRangeContextRouterSkipRow struct {
	Timeframe string `json:"timeframe"`
	Split     string `json:"split"`
	Reason    string `json:"reason"`
	Count     int    `json:"count"`
}

type rangeContextRouterCohortKey struct {
	cohortID    string
	routerLabel string
	split       string
	timeframe   string
	horizonBars int
}

type rangeContextRouterCohortAccumulator struct {
	row    FuturesRangeContextRouterCohortRow
	labels map[string]int
}

func DefaultFuturesRangeContextRouterAuditConfig() FuturesRangeContextRouterAuditConfig {
	stateCfg := DefaultFuturesRangeStateConstructionLoopAuditConfig()
	return FuturesRangeContextRouterAuditConfig{
		StateAuditConfig:             stateCfg,
		MinFullRouterRows:            stateCfg.MinFullCohortCount,
		MinSplitRouterRows:           stateCfg.MinSplitCohortCount,
		MinNoTradeFullRouterRows:     stateCfg.MinNoTradeFullCohortCount,
		MinNoTradeSplitRouterRows:    stateCfg.MinNoTradeSplitCohortCount,
		MaxSplitContributionRate:     stateCfg.MaxSplitContributionRate,
		MinPositiveExpectedRateFull:  stateCfg.MinPositiveUsefulRateFull,
		MinPositiveExpectedRateSplit: stateCfg.MinPositiveUsefulRateSplit,
		MaxPositiveAdverseRateFull:   stateCfg.MaxPositiveToxicRateFull,
		MaxPositiveAdverseRateSplit:  stateCfg.MaxPositiveToxicRateSplit,
		MinPositiveMarginFull:        stateCfg.MinPositiveMarginFull,
		MinPositiveMarginSplit:       stateCfg.MinPositiveMarginSplit,
		MinNoTradeToxicRateFull:      stateCfg.MinToxicRateFull,
		MinNoTradeToxicRateSplit:     stateCfg.MinToxicRateSplit,
	}
}

func RunFuturesRangeContextRouterAudit(candles []Candle, manifest SourceManifest, cfg FuturesRangeContextRouterAuditConfig, splits []Split) (FuturesRangeContextRouterAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeContextRouterAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	stateResult, err := RunFuturesRangeStateConstructionLoopAudit(candles, manifest, cfg.StateAuditConfig, splits)
	if err != nil {
		return FuturesRangeContextRouterAuditResult{}, err
	}

	result := FuturesRangeContextRouterAuditResult{StateStopState: stateResult.StopState}
	sourcePass := rangeStateConstructionLoopSourceResamplePass(stateResult.SourceRows, stateResult.CoverageRows)
	result.SourceRows = rangeContextRouterSourceRows(stateResult.SourceRows, stateResult.StopState, sourcePass)
	result.CoverageRows = rangeContextRouterCoverageRows(stateResult.CoverageRows, stateResult.StateRows, nil, sourcePass)
	if !sourcePass || stateResult.StopState == RangeStateConstructionLoopStopStateSourceGap {
		result.StopState = RangeContextRouterStopStateSourceGap
		result.SummaryRows = rangeContextRouterSummaryRows(result, cfg, splits)
		return result, nil
	}
	if stateResult.StopState == RangeStateConstructionLoopStopStateRejectedClosedReslice {
		result.StopState = RangeContextRouterStopStateRejectedClosedReslice
		result.SummaryRows = rangeContextRouterSummaryRows(result, cfg, splits)
		return result, nil
	}

	result.RuleRows = rangeContextRouterRuleRows(stateResult.RankingRows)
	result.Rows = rangeContextRouterRows(stateResult.StateRows, result.RuleRows)
	result.CoverageRows = rangeContextRouterCoverageRows(stateResult.CoverageRows, stateResult.StateRows, result.Rows, sourcePass)
	result.SkipRows = rangeContextRouterSkipRows(stateResult.SkipRows, result.Rows)
	result.CohortRows = rangeContextRouterCohortRows(result.Rows, stateResult.LabelRows, cfg, splits, sourcePass)
	result.RankingRows = rangeContextRouterRankingRows(result.CohortRows)
	result.PassingCohorts = rangeContextRouterPassingRankingCount(result.RankingRows)
	result.StopState = FuturesRangeContextRouterAuditStopState(result)
	result.SummaryRows = rangeContextRouterSummaryRows(result, cfg, splits)
	return result, nil
}

func FuturesRangeContextRouterAuditStopState(result FuturesRangeContextRouterAuditResult) string {
	if result.StopState == RangeContextRouterStopStateSourceGap || result.StopState == RangeContextRouterStopStateRejectedClosedReslice {
		return result.StopState
	}
	if !rangeContextRouterSourcePass(result.SourceRows, result.CoverageRows) {
		return RangeContextRouterStopStateSourceGap
	}
	for _, rule := range result.RuleRows {
		if rule.ClosedFamilyReslice || rule.RuleStatus == rangeContextRouterRuleStatusRejectedClosedFamily {
			return RangeContextRouterStopStateRejectedClosedReslice
		}
	}
	noTradePasses := 0
	rotationPasses := 0
	continuationPasses := 0
	for _, row := range result.RankingRows {
		if !row.PassesGate {
			continue
		}
		switch row.RouterLabel {
		case RangeContextRouterLabelNoTrade:
			noTradePasses++
		case RangeContextRouterLabelTradableRotation:
			rotationPasses++
		case RangeContextRouterLabelTrendContinuation:
			continuationPasses++
		}
	}
	switch {
	case rotationPasses == 0 && continuationPasses == 0 && noTradePasses == 0:
		return RangeContextRouterStopStateFailedNoActionableRoute
	case rotationPasses == 0 && continuationPasses == 0 && noTradePasses > 0:
		return RangeContextRouterStopStatePassedNoTradeFilterOnly
	case continuationPasses > 0 && rotationPasses == 0:
		return RangeContextRouterStopStatePassedNeedsContinuation
	default:
		return RangeContextRouterStopStatePassedNeedsRotationSpec
	}
}

func (cfg FuturesRangeContextRouterAuditConfig) withDefaults() FuturesRangeContextRouterAuditConfig {
	defaults := DefaultFuturesRangeContextRouterAuditConfig()
	cfg.StateAuditConfig = cfg.StateAuditConfig.withDefaults()
	if cfg.MinFullRouterRows == 0 {
		cfg.MinFullRouterRows = defaults.MinFullRouterRows
	}
	if cfg.MinSplitRouterRows == 0 {
		cfg.MinSplitRouterRows = defaults.MinSplitRouterRows
	}
	if cfg.MinNoTradeFullRouterRows == 0 {
		cfg.MinNoTradeFullRouterRows = defaults.MinNoTradeFullRouterRows
	}
	if cfg.MinNoTradeSplitRouterRows == 0 {
		cfg.MinNoTradeSplitRouterRows = defaults.MinNoTradeSplitRouterRows
	}
	if cfg.MaxSplitContributionRate == 0 {
		cfg.MaxSplitContributionRate = defaults.MaxSplitContributionRate
	}
	if cfg.MinPositiveExpectedRateFull == 0 {
		cfg.MinPositiveExpectedRateFull = defaults.MinPositiveExpectedRateFull
	}
	if cfg.MinPositiveExpectedRateSplit == 0 {
		cfg.MinPositiveExpectedRateSplit = defaults.MinPositiveExpectedRateSplit
	}
	if cfg.MaxPositiveAdverseRateFull == 0 {
		cfg.MaxPositiveAdverseRateFull = defaults.MaxPositiveAdverseRateFull
	}
	if cfg.MaxPositiveAdverseRateSplit == 0 {
		cfg.MaxPositiveAdverseRateSplit = defaults.MaxPositiveAdverseRateSplit
	}
	if cfg.MinPositiveMarginFull == 0 {
		cfg.MinPositiveMarginFull = defaults.MinPositiveMarginFull
	}
	if cfg.MinPositiveMarginSplit == 0 {
		cfg.MinPositiveMarginSplit = defaults.MinPositiveMarginSplit
	}
	if cfg.MinNoTradeToxicRateFull == 0 {
		cfg.MinNoTradeToxicRateFull = defaults.MinNoTradeToxicRateFull
	}
	if cfg.MinNoTradeToxicRateSplit == 0 {
		cfg.MinNoTradeToxicRateSplit = defaults.MinNoTradeToxicRateSplit
	}
	return cfg
}

func (cfg FuturesRangeContextRouterAuditConfig) validate() error {
	if err := cfg.StateAuditConfig.validate(); err != nil {
		return err
	}
	if cfg.MinFullRouterRows <= 0 || cfg.MinSplitRouterRows <= 0 || cfg.MinNoTradeFullRouterRows <= 0 || cfg.MinNoTradeSplitRouterRows <= 0 {
		return fmt.Errorf("range context router count gates must be positive")
	}
	if cfg.MaxSplitContributionRate <= 0 || cfg.MaxSplitContributionRate > 1 {
		return fmt.Errorf("range context router split contribution gate must be in (0,1]")
	}
	return nil
}

func rangeContextRouterSourceRows(sourceRows []FuturesRangeStateConstructionLoopSourceRow, stateStop string, sourcePass bool) []FuturesRangeContextRouterSourceRow {
	rows := make([]FuturesRangeContextRouterSourceRow, 0, len(sourceRows))
	for _, source := range sourceRows {
		rows = append(rows, FuturesRangeContextRouterSourceRow{
			AuditName:               FuturesRangeContextRouterAuditName,
			StateAuditName:          FuturesRangeStateConstructionLoopAuditName,
			StateAuditStopState:     stateStop,
			SourceResamplePass:      sourcePass,
			RouterUsesForwardLabels: false,
			FuturesRangeStateConstructionLoopSourceRow: source,
		})
	}
	return rows
}

func rangeContextRouterCoverageRows(coverageRows []FuturesRangeStateConstructionLoopCoverageRow, stateRows []FuturesRangeStateConstructionLoopStateRow, routerRows []FuturesRangeContextRouterRow, sourcePass bool) []FuturesRangeContextRouterCoverageRow {
	rows := make([]FuturesRangeContextRouterCoverageRow, 0, len(coverageRows))
	for _, coverage := range coverageRows {
		row := FuturesRangeContextRouterCoverageRow{
			AuditName:          FuturesRangeContextRouterAuditName,
			SourceResamplePass: sourcePass,
			FuturesRangeStateConstructionLoopCoverageRow: coverage,
		}
		for _, state := range stateRows {
			if state.Timeframe == coverage.Timeframe {
				row.StateRows++
			}
		}
		for _, router := range routerRows {
			if router.Timeframe == coverage.Timeframe {
				row.RouterRows++
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func rangeContextRouterRuleRows(rankings []FuturesRangeStateConstructionLoopRankingRow) []FuturesRangeContextRouterRuleRow {
	rows := []FuturesRangeContextRouterRuleRow{}
	for _, ranking := range rankings {
		if !ranking.PassesGate {
			continue
		}
		routerLabel := rangeContextRouterLabelFromRoute(ranking.RouteCandidate)
		if routerLabel == "" || routerLabel == RangeContextRouterLabelDiagnosticOnly {
			continue
		}
		status := rangeContextRouterRuleStatusActive
		reason := ""
		if ranking.ClosedFamilyReslice {
			status = rangeContextRouterRuleStatusRejectedClosedFamily
			reason = "closed_family_reslice"
		}
		if !ranking.FutureLeakProtectionPass {
			status = rangeContextRouterRuleStatusRejectedFutureLeak
			reason = uniqueJoinedReasons([]string{reason, "future_label_as_feature"})
		}
		row := FuturesRangeContextRouterRuleRow{
			RuleSource:                 rangeContextRouterRuleSourcePassedStateRanking,
			RuleStatus:                 status,
			RouterLabel:                routerLabel,
			SourceRouteCandidate:       ranking.RouteCandidate,
			SourceStateRank:            ranking.Rank,
			SourceStateCohortID:        ranking.CohortID,
			RollupType:                 ranking.RollupType,
			RollupID:                   ranking.RollupID,
			Timeframe:                  ranking.Timeframe,
			HorizonBars:                ranking.HorizonBars,
			FullPeriodRows:             ranking.FullPeriodRows,
			WeakestSplitRows:           ranking.WeakestSplitRows,
			FullExpectedRouteRate:      ranking.FullUsefulRate,
			WeakestSplitExpectedRate:   ranking.WeakestSplitUsefulRate,
			FullAdverseRouteRate:       ranking.FullToxicRate,
			WorstSplitAdverseRate:      ranking.WorstSplitToxicRate,
			FullExpectedMinusAdverse:   ranking.FullUsefulMinusToxicMargin,
			WeakestSplitMargin:         ranking.WeakestSplitMargin,
			ClosedFamilyReslice:        ranking.ClosedFamilyReslice,
			FutureLeakProtectionPass:   ranking.FutureLeakProtectionPass,
			ForwardLabelsAsRouterInput: false,
			FeatureInputs:              ranking.RollupType,
			FailureReason:              reason,
		}
		if routerLabel == RangeContextRouterLabelNoTrade {
			row.FullExpectedRouteRate = ranking.FullToxicRate
			row.WeakestSplitExpectedRate = ranking.WorstSplitToxicRate
			row.FullAdverseRouteRate = ranking.FullUsefulRate
			row.WorstSplitAdverseRate = ranking.WeakestSplitUsefulRate
			row.FullExpectedMinusAdverse = ranking.FullToxicRate - ranking.FullUsefulRate
			row.WeakestSplitMargin = ranking.WorstSplitToxicRate - ranking.WeakestSplitUsefulRate
		}
		row.RuleID = rangeContextRouterRuleID(row)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RouterLabel != rows[j].RouterLabel {
			return rangeContextRouterLabelSortKey(rows[i].RouterLabel) < rangeContextRouterLabelSortKey(rows[j].RouterLabel)
		}
		if rows[i].SourceStateRank != rows[j].SourceStateRank {
			return rows[i].SourceStateRank < rows[j].SourceStateRank
		}
		return rows[i].RuleID < rows[j].RuleID
	})
	return rows
}

func rangeContextRouterRows(states []FuturesRangeStateConstructionLoopStateRow, rules []FuturesRangeContextRouterRuleRow) []FuturesRangeContextRouterRow {
	rows := make([]FuturesRangeContextRouterRow, 0, len(states))
	activeRules := make([]FuturesRangeContextRouterRuleRow, 0, len(rules))
	for _, rule := range rules {
		if rule.RuleStatus == rangeContextRouterRuleStatusActive {
			activeRules = append(activeRules, rule)
		}
	}
	for _, state := range states {
		row := FuturesRangeContextRouterRow{
			RouterRowID:                len(rows) + 1,
			StateRowID:                 state.StateRowID,
			Timestamp:                  state.Timestamp,
			Timeframe:                  state.Timeframe,
			Split:                      state.Split,
			RangeEpisodeID:             state.RangeEpisodeID,
			RangeStartIndex:            state.RangeStartIndex,
			DecisionIndex:              state.DecisionIndex,
			RangeStartTime:             state.RangeStartTime,
			DecisionCloseTime:          state.DecisionCloseTime,
			StateID:                    state.StateID,
			GeometryVolID:              state.GeometryVolID,
			GeometryTrendID:            state.GeometryTrendID,
			GeometryImpulseID:          state.GeometryImpulseID,
			GeometryParticipationID:    state.GeometryParticipationID,
			GeometryVolTrendID:         state.GeometryVolTrendID,
			GeometryVolTrendImpulseID:  state.GeometryVolTrendImpulseID,
			AllFamiliesID:              state.AllFamiliesID,
			RouterLabel:                RangeContextRouterLabelDiagnosticOnly,
			RouterReason:               rangeContextRouterRowReasonNoRuleMatch,
			MissingRuleMatch:           true,
			ClosedCandleOnly:           true,
			ForwardLabelsAsRouterInput: false,
			ForwardLabelColumnsPresent: false,
		}
		if !state.Eligible {
			row.RouterReason = rangeContextRouterRowReasonStateIneligible
			rows = append(rows, row)
			continue
		}
		matched := rangeContextRouterMatchedRules(state, activeRules)
		if len(matched) == 0 {
			if rangeContextRouterMissingStateRollup(state) {
				row.RouterReason = rangeContextRouterRowReasonMissingStateRollup
			}
			rows = append(rows, row)
			continue
		}
		labelSet := map[string]struct{}{}
		ruleIDs := make([]string, 0, len(matched))
		rollupTypes := make([]string, 0, len(matched))
		for _, rule := range matched {
			labelSet[rule.RouterLabel] = struct{}{}
			ruleIDs = append(ruleIDs, rule.RuleID)
			rollupTypes = append(rollupTypes, rule.RollupType)
		}
		labels := sortedRouterLabels(labelSet)
		row.MatchedRuleCount = len(matched)
		row.MatchedRuleIDs = strings.Join(uniqueStrings(ruleIDs), ";")
		row.MatchedRollupTypes = strings.Join(uniqueStrings(rollupTypes), ";")
		row.MatchedRouterLabels = strings.Join(labels, ";")
		row.MissingRuleMatch = false
		if len(labels) > 1 {
			row.RouterLabel = RangeContextRouterLabelDiagnosticOnly
			row.RouterReason = rangeContextRouterRowReasonConflictingRules
			row.ConflictingRuleMatch = true
			rows = append(rows, row)
			continue
		}
		row.RouterLabel = labels[0]
		row.RouterReason = rangeContextRouterRowReasonMatchedRule
		row.Actionable = row.RouterLabel != RangeContextRouterLabelDiagnosticOnly
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		return rows[i].DecisionIndex < rows[j].DecisionIndex
	})
	return rows
}

func rangeContextRouterMatchedRules(state FuturesRangeStateConstructionLoopStateRow, rules []FuturesRangeContextRouterRuleRow) []FuturesRangeContextRouterRuleRow {
	matched := []FuturesRangeContextRouterRuleRow{}
	for _, rule := range rules {
		if rule.Timeframe != state.Timeframe {
			continue
		}
		if rangeContextRouterStateRollupID(state, rule.RollupType) == rule.RollupID {
			matched = append(matched, rule)
		}
	}
	return matched
}

func rangeContextRouterCohortRows(rows []FuturesRangeContextRouterRow, labels []FuturesRangeStateConstructionLoopLabelRow, cfg FuturesRangeContextRouterAuditConfig, splits []Split, sourcePass bool) []FuturesRangeContextRouterCohortRow {
	routerByStateID := map[int]FuturesRangeContextRouterRow{}
	for _, row := range rows {
		routerByStateID[row.StateRowID] = row
	}
	accs := map[rangeContextRouterCohortKey]*rangeContextRouterCohortAccumulator{}
	for _, label := range labels {
		router, ok := routerByStateID[label.StateRowID]
		if !ok {
			continue
		}
		for _, splitName := range rangeContextRouterSplitNames(router.Split) {
			key := rangeContextRouterCohortKeyFor(splitName, router.Timeframe, label.HorizonBars, router.RouterLabel)
			acc := accs[key]
			if acc == nil {
				acc = &rangeContextRouterCohortAccumulator{
					row: FuturesRangeContextRouterCohortRow{
						CohortID:    key.cohortID,
						RouterLabel: key.routerLabel,
						Split:       key.split,
						Timeframe:   key.timeframe,
						HorizonBars: key.horizonBars,
					},
					labels: map[string]int{},
				}
				accs[key] = acc
			}
			acc.add(label)
		}
	}
	cohorts := make([]FuturesRangeContextRouterCohortRow, 0, len(accs))
	for _, acc := range accs {
		cohorts = append(cohorts, acc.finalRow())
	}
	rangeContextRouterMarkCohortGates(cohorts, cfg, splits, sourcePass)
	sort.Slice(cohorts, func(i, j int) bool {
		if cohorts[i].RouterLabel != cohorts[j].RouterLabel {
			return rangeContextRouterLabelSortKey(cohorts[i].RouterLabel) < rangeContextRouterLabelSortKey(cohorts[j].RouterLabel)
		}
		if cohorts[i].Timeframe != cohorts[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(cohorts[i].Timeframe) < rangeContextTriageTimeframeSortKey(cohorts[j].Timeframe)
		}
		if cohorts[i].HorizonBars != cohorts[j].HorizonBars {
			return cohorts[i].HorizonBars < cohorts[j].HorizonBars
		}
		if cohorts[i].Split != cohorts[j].Split {
			return splitSortKey(cohorts[i].Split) < splitSortKey(cohorts[j].Split)
		}
		return cohorts[i].CohortID < cohorts[j].CohortID
	})
	return cohorts
}

func (acc *rangeContextRouterCohortAccumulator) add(label FuturesRangeStateConstructionLoopLabelRow) {
	acc.row.CandidateCount++
	acc.labels[label.ForwardLabel]++
	if label.RotationUseful {
		acc.row.RotationUsefulCount++
	}
	if label.RotationToxic {
		acc.row.RotationToxicCount++
	}
	if label.ContinuationUseful {
		acc.row.ContinuationUsefulCount++
	}
	if label.ContinuationToxic {
		acc.row.ContinuationToxicCount++
	}
	if label.NoTradeToxic {
		acc.row.NoTradeToxicCount++
	}
	if rangeContextRouterExpectedHit(acc.row.RouterLabel, label) {
		acc.row.ExpectedRouteHitCount++
	}
	if rangeContextRouterAdverseHit(acc.row.RouterLabel, label) {
		acc.row.AdverseRouteHitCount++
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

func (acc *rangeContextRouterCohortAccumulator) finalRow() FuturesRangeContextRouterCohortRow {
	row := acc.row
	if row.CandidateCount > 0 {
		row.ExpectedRouteHitRate = float64(row.ExpectedRouteHitCount) / float64(row.CandidateCount)
		row.AdverseRouteHitRate = float64(row.AdverseRouteHitCount) / float64(row.CandidateCount)
		row.ExpectedMinusAdverseMargin = row.ExpectedRouteHitRate - row.AdverseRouteHitRate
	}
	row.DominantForwardLabel, row.DominantForwardLabelRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, false)
	return row
}

func rangeContextRouterMarkCohortGates(rows []FuturesRangeContextRouterCohortRow, cfg FuturesRangeContextRouterAuditConfig, splits []Split, sourcePass bool) {
	byIDSplit := map[string]map[string]*FuturesRangeContextRouterCohortRow{}
	for i := range rows {
		if byIDSplit[rows[i].CohortID] == nil {
			byIDSplit[rows[i].CohortID] = map[string]*FuturesRangeContextRouterCohortRow{}
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
			rows[i].FullExpectedRouteHitRate = full.ExpectedRouteHitRate
			rows[i].FullAdverseRouteHitRate = full.AdverseRouteHitRate
			rows[i].FullExpectedMinusAdverseMargin = full.ExpectedMinusAdverseMargin
		}
		rows[i].FullPeriodRows = fullRows
		rows[i].WeakestSplitRows = int(^uint(0) >> 1)
		rows[i].WeakestSplitExpectedHitRate = math.Inf(1)
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
				rows[i].WeakestSplitExpectedHitRate = 0
				rows[i].WeakestSplitMargin = 0
				continue
			}
			minSplitRows := cfg.MinSplitRouterRows
			if rows[i].RouterLabel == RangeContextRouterLabelNoTrade {
				minSplitRows = cfg.MinNoTradeSplitRouterRows
			}
			if row.CandidateCount < minSplitRows {
				splitCountsPass = false
			}
			if row.CandidateCount < rows[i].WeakestSplitRows {
				rows[i].WeakestSplitRows = row.CandidateCount
			}
			if row.ExpectedRouteHitRate < rows[i].WeakestSplitExpectedHitRate {
				rows[i].WeakestSplitExpectedHitRate = row.ExpectedRouteHitRate
			}
			if row.AdverseRouteHitRate > rows[i].WorstSplitAdverseHitRate {
				rows[i].WorstSplitAdverseHitRate = row.AdverseRouteHitRate
			}
			if row.ExpectedMinusAdverseMargin < rows[i].WeakestSplitMargin {
				rows[i].WeakestSplitMargin = row.ExpectedMinusAdverseMargin
			}
			if fullRows > 0 {
				maxContribution = math.Max(maxContribution, float64(row.CandidateCount)/float64(fullRows))
			}
		}
		if rows[i].WeakestSplitRows == int(^uint(0)>>1) {
			rows[i].WeakestSplitRows = 0
		}
		if math.IsInf(rows[i].WeakestSplitExpectedHitRate, 1) {
			rows[i].WeakestSplitExpectedHitRate = 0
		}
		if math.IsInf(rows[i].WeakestSplitMargin, 1) {
			rows[i].WeakestSplitMargin = 0
		}
		rows[i].MaxSplitContributionRate = maxContribution
		minFullRows := cfg.MinFullRouterRows
		if rows[i].RouterLabel == RangeContextRouterLabelNoTrade {
			minFullRows = cfg.MinNoTradeFullRouterRows
		}
		rows[i].ReviewableCountGatePass = fullRows >= minFullRows && splitCountsPass
		rows[i].SplitStabilityGatePass = splitRatesPresent && len(periodSplits) > 0
		rows[i].SplitContributionGatePass = maxContribution <= cfg.MaxSplitContributionRate || len(periodSplits) == 0
		rows[i].ClosedFamilyProtectionPass = true
		rows[i].FutureLeakProtectionPass = true
		switch rows[i].RouterLabel {
		case RangeContextRouterLabelTradableRotation, RangeContextRouterLabelTrendContinuation:
			rows[i].RouteRateGatePass = full != nil &&
				rows[i].FullExpectedRouteHitRate >= cfg.MinPositiveExpectedRateFull &&
				rows[i].WeakestSplitExpectedHitRate >= cfg.MinPositiveExpectedRateSplit &&
				rows[i].FullAdverseRouteHitRate <= cfg.MaxPositiveAdverseRateFull &&
				rows[i].WorstSplitAdverseHitRate <= cfg.MaxPositiveAdverseRateSplit &&
				rows[i].FullExpectedMinusAdverseMargin >= cfg.MinPositiveMarginFull &&
				rows[i].WeakestSplitMargin >= cfg.MinPositiveMarginSplit
		case RangeContextRouterLabelNoTrade:
			rows[i].RouteRateGatePass = full != nil &&
				rows[i].FullExpectedRouteHitRate >= cfg.MinNoTradeToxicRateFull &&
				rows[i].WeakestSplitExpectedHitRate >= cfg.MinNoTradeToxicRateSplit
		default:
			rows[i].RouteRateGatePass = false
		}
		if !sourcePass {
			reasons = append(reasons, "source_or_resample_gap")
		}
		if !rows[i].ReviewableCountGatePass {
			reasons = append(reasons, "inadequate_router_count")
		}
		if !rows[i].SplitStabilityGatePass {
			reasons = append(reasons, "missing_period_split")
		}
		if !rows[i].SplitContributionGatePass {
			reasons = append(reasons, "single_split_contribution_above_gate")
		}
		if !rows[i].RouteRateGatePass {
			reasons = append(reasons, "route_rate_gate_failed")
		}
		rows[i].PassesReviewGate = sourcePass &&
			rows[i].ReviewableCountGatePass &&
			rows[i].SplitStabilityGatePass &&
			rows[i].SplitContributionGatePass &&
			rows[i].RouteRateGatePass &&
			rows[i].ClosedFamilyProtectionPass &&
			rows[i].FutureLeakProtectionPass &&
			rows[i].Split == fullSplitName &&
			rows[i].RouterLabel != RangeContextRouterLabelDiagnosticOnly
		rows[i].FailureReason = uniqueJoinedReasons(reasons)
	}
}

func rangeContextRouterRankingRows(cohorts []FuturesRangeContextRouterCohortRow) []FuturesRangeContextRouterRankingRow {
	rows := []FuturesRangeContextRouterRankingRow{}
	for _, cohort := range cohorts {
		if cohort.Split != fullSplitName || cohort.RouterLabel == RangeContextRouterLabelDiagnosticOnly {
			continue
		}
		row := FuturesRangeContextRouterRankingRow{
			CohortID:                       cohort.CohortID,
			RouterLabel:                    cohort.RouterLabel,
			Timeframe:                      cohort.Timeframe,
			HorizonBars:                    cohort.HorizonBars,
			PassesGate:                     cohort.PassesReviewGate,
			FullPeriodRows:                 cohort.FullPeriodRows,
			WeakestSplitRows:               cohort.WeakestSplitRows,
			MaxSplitContributionRate:       cohort.MaxSplitContributionRate,
			FullExpectedRouteHitRate:       cohort.FullExpectedRouteHitRate,
			WeakestSplitExpectedHitRate:    cohort.WeakestSplitExpectedHitRate,
			FullAdverseRouteHitRate:        cohort.FullAdverseRouteHitRate,
			WorstSplitAdverseHitRate:       cohort.WorstSplitAdverseHitRate,
			FullExpectedMinusAdverseMargin: cohort.FullExpectedMinusAdverseMargin,
			WeakestSplitMargin:             cohort.WeakestSplitMargin,
			DominantForwardLabel:           cohort.DominantForwardLabel,
			DominantForwardLabelRate:       cohort.DominantForwardLabelRate,
			FailureReason:                  cohort.FailureReason,
		}
		row.RankScore = rangeContextRouterRankScore(row)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		if rows[i].RouterLabel != rows[j].RouterLabel {
			return rangeContextRouterLabelSortKey(rows[i].RouterLabel) < rangeContextRouterLabelSortKey(rows[j].RouterLabel)
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

func rangeContextRouterSummaryRows(result FuturesRangeContextRouterAuditResult, cfg FuturesRangeContextRouterAuditConfig, splits []Split) []FuturesRangeContextRouterSummaryRow {
	sourcePass := rangeContextRouterSourcePass(result.SourceRows, result.CoverageRows)
	rows := []FuturesRangeContextRouterSummaryRow{rangeContextRouterSummaryRowFor(result, fullSplitName, "all", 0, sourcePass)}
	for _, split := range splits {
		for _, timeframe := range cfg.StateAuditConfig.Timeframes {
			for _, horizon := range rangeStateConstructionLoopHorizons(timeframe, cfg.StateAuditConfig) {
				rows = append(rows, rangeContextRouterSummaryRowFor(result, split.Name, timeframe, horizon, sourcePass))
			}
		}
	}
	return rows
}

func rangeContextRouterSummaryRowFor(result FuturesRangeContextRouterAuditResult, split, timeframe string, horizon int, sourcePass bool) FuturesRangeContextRouterSummaryRow {
	row := FuturesRangeContextRouterSummaryRow{
		Split:               split,
		Timeframe:           timeframe,
		HorizonBars:         horizon,
		SourceResamplePass:  sourcePass,
		StateAuditStopState: result.StateStopState,
		RuleRows:            len(result.RuleRows),
		StopState:           result.StopState,
	}
	for _, router := range result.Rows {
		if !rangeContextRouterSummaryRowMatches(router, split, timeframe) {
			continue
		}
		row.RouterRows++
		switch router.RouterLabel {
		case RangeContextRouterLabelNoTrade:
			row.NoTradeRows++
		case RangeContextRouterLabelTradableRotation:
			row.TradableRotationRows++
		case RangeContextRouterLabelTrendContinuation:
			row.TrendContinuationRows++
		default:
			row.DiagnosticOnlyRows++
		}
		if router.ConflictingRuleMatch {
			row.ConflictRows++
		}
		if router.MissingRuleMatch {
			row.MissingRuleRows++
		}
	}
	for _, cohort := range result.CohortRows {
		if cohort.Split == split && (timeframe == "all" || cohort.Timeframe == timeframe) && (horizon == 0 || cohort.HorizonBars == horizon) {
			row.CohortRows++
			if cohort.PassesReviewGate {
				row.PassingCohorts++
			}
		}
	}
	for _, ranking := range result.RankingRows {
		if (timeframe == "all" || ranking.Timeframe == timeframe) && (horizon == 0 || ranking.HorizonBars == horizon) {
			row.RankingRows++
		}
	}
	for _, skip := range result.SkipRows {
		if skip.Split == split && (timeframe == "all" || skip.Timeframe == timeframe) {
			row.SkipRows += skip.Count
		}
	}
	return row
}

func rangeContextRouterSkipRows(stateSkips []FuturesRangeStateConstructionLoopSkipRow, rows []FuturesRangeContextRouterRow) []FuturesRangeContextRouterSkipRow {
	skips := map[rangeStateSkipKey]int{}
	for _, skip := range stateSkips {
		skips[rangeStateSkipKey{timeframe: skip.Timeframe, split: skip.Split, reason: "state_audit_" + skip.Reason}] += skip.Count
	}
	for _, row := range rows {
		if row.RouterLabel != RangeContextRouterLabelDiagnosticOnly || row.RouterReason == "" {
			continue
		}
		rangeContextRouterAddSkip(skips, row.Timeframe, row.Split, row.RouterReason)
	}
	out := make([]FuturesRangeContextRouterSkipRow, 0, len(skips))
	for key, count := range skips {
		out = append(out, FuturesRangeContextRouterSkipRow{
			Timeframe: key.timeframe,
			Split:     key.split,
			Reason:    key.reason,
			Count:     count,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Timeframe != out[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(out[i].Timeframe) < rangeContextTriageTimeframeSortKey(out[j].Timeframe)
		}
		if out[i].Split != out[j].Split {
			return splitSortKey(out[i].Split) < splitSortKey(out[j].Split)
		}
		return out[i].Reason < out[j].Reason
	})
	return out
}

func rangeContextRouterAddSkip(skips map[rangeStateSkipKey]int, timeframe, split, reason string) {
	if split == "" {
		split = fullSplitName
	}
	skips[rangeStateSkipKey{timeframe: timeframe, split: split, reason: reason}]++
	if split != fullSplitName {
		skips[rangeStateSkipKey{timeframe: timeframe, split: fullSplitName, reason: reason}]++
	}
}

func rangeContextRouterSourcePass(sources []FuturesRangeContextRouterSourceRow, coverage []FuturesRangeContextRouterCoverageRow) bool {
	if len(sources) == 0 || len(coverage) == 0 {
		return false
	}
	for _, source := range sources {
		if !source.SourceResamplePass || source.ValidationStatus != "accepted" || !source.SourceFactsPass || source.RouterUsesForwardLabels {
			return false
		}
	}
	for _, row := range coverage {
		if !row.SourceResamplePass || row.ValidationStatus != "accepted" || !row.CoverageFactsPass || !row.Complete {
			return false
		}
	}
	return true
}

func rangeContextRouterLabelFromRoute(route string) string {
	switch route {
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return RangeContextRouterLabelNoTrade
	case RangeStateConstructionLoopRouteRotation:
		return RangeContextRouterLabelTradableRotation
	case RangeStateConstructionLoopRouteContinuation:
		return RangeContextRouterLabelTrendContinuation
	case RangeStateConstructionLoopRouteDiagnosticOnly:
		return RangeContextRouterLabelDiagnosticOnly
	default:
		return ""
	}
}

func rangeContextRouterRuleID(row FuturesRangeContextRouterRuleRow) string {
	parts := []string{"range_context_router_v1", row.RouterLabel, row.Timeframe, fmt.Sprintf("h%d", row.HorizonBars), row.RollupType, row.RollupID}
	return strings.Join(parts, "|")
}

func rangeContextRouterStateRollupID(state FuturesRangeStateConstructionLoopStateRow, rollupType string) string {
	switch rollupType {
	case RangeStateConstructionLoopRollupGeometryVol:
		return state.GeometryVolID
	case RangeStateConstructionLoopRollupGeometryTrend:
		return state.GeometryTrendID
	case RangeStateConstructionLoopRollupGeometryImpulse:
		return state.GeometryImpulseID
	case RangeStateConstructionLoopRollupGeometryParticipation:
		return state.GeometryParticipationID
	case RangeStateConstructionLoopRollupGeometryVolTrend:
		return state.GeometryVolTrendID
	case RangeStateConstructionLoopRollupGeometryVolTrendImpulse:
		return state.GeometryVolTrendImpulseID
	case RangeStateConstructionLoopRollupAllFamilies:
		return state.AllFamiliesID
	default:
		return ""
	}
}

func rangeContextRouterMissingStateRollup(state FuturesRangeStateConstructionLoopStateRow) bool {
	return state.GeometryVolID == "" ||
		state.GeometryTrendID == "" ||
		state.GeometryImpulseID == "" ||
		state.GeometryParticipationID == "" ||
		state.GeometryVolTrendID == "" ||
		state.GeometryVolTrendImpulseID == "" ||
		state.AllFamiliesID == ""
}

func rangeContextRouterExpectedHit(routerLabel string, label FuturesRangeStateConstructionLoopLabelRow) bool {
	switch routerLabel {
	case RangeContextRouterLabelNoTrade:
		return label.NoTradeToxic
	case RangeContextRouterLabelTradableRotation:
		return label.RotationUseful
	case RangeContextRouterLabelTrendContinuation:
		return label.ContinuationUseful
	default:
		return false
	}
}

func rangeContextRouterAdverseHit(routerLabel string, label FuturesRangeStateConstructionLoopLabelRow) bool {
	switch routerLabel {
	case RangeContextRouterLabelNoTrade:
		return label.RotationUseful || label.ContinuationUseful
	case RangeContextRouterLabelTradableRotation:
		return label.RotationToxic
	case RangeContextRouterLabelTrendContinuation:
		return label.ContinuationToxic
	default:
		return false
	}
}

func rangeContextRouterCohortKeyFor(split, timeframe string, horizon int, routerLabel string) rangeContextRouterCohortKey {
	cohortID := strings.Join([]string{"range_context_router_v1", timeframe, fmt.Sprintf("h%d", horizon), routerLabel}, "|")
	return rangeContextRouterCohortKey{
		cohortID:    cohortID,
		routerLabel: routerLabel,
		split:       split,
		timeframe:   timeframe,
		horizonBars: horizon,
	}
}

func rangeContextRouterSplitNames(split string) []string {
	if split == "" || split == fullSplitName {
		return []string{fullSplitName}
	}
	return []string{split, fullSplitName}
}

func rangeContextRouterSummaryRowMatches(row FuturesRangeContextRouterRow, split, timeframe string) bool {
	if timeframe != "all" && row.Timeframe != timeframe {
		return false
	}
	if split == fullSplitName {
		return true
	}
	return row.Split == split
}

func rangeContextRouterRankScore(row FuturesRangeContextRouterRankingRow) float64 {
	score := row.FullExpectedRouteHitRate + row.WeakestSplitExpectedHitRate - row.FullAdverseRouteHitRate - row.WorstSplitAdverseHitRate
	score += math.Min(float64(row.FullPeriodRows), 1000) / 1000
	score += math.Min(float64(row.WeakestSplitRows), 250) / 1000
	score -= row.MaxSplitContributionRate / 10
	return score
}

func rangeContextRouterPassingRankingCount(rows []FuturesRangeContextRouterRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func rangeContextRouterLabelSortKey(label string) int {
	switch label {
	case RangeContextRouterLabelNoTrade:
		return 0
	case RangeContextRouterLabelTradableRotation:
		return 1
	case RangeContextRouterLabelTrendContinuation:
		return 2
	case RangeContextRouterLabelDiagnosticOnly:
		return 3
	default:
		return 99
	}
}

func sortedRouterLabels(labelSet map[string]struct{}) []string {
	labels := make([]string, 0, len(labelSet))
	for label := range labelSet {
		labels = append(labels, label)
	}
	sort.Slice(labels, func(i, j int) bool {
		return rangeContextRouterLabelSortKey(labels[i]) < rangeContextRouterLabelSortKey(labels[j])
	})
	return labels
}
