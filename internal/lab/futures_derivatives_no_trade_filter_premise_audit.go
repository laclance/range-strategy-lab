package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesDerivativesNoTradeFilterPremiseAuditName = "futures_derivatives_no_trade_filter_premise_audit"

	DerivativesNoTradeFilterPremiseStopStateSourceGap                  = "derivatives_context_no_trade_filter_premise_audit_source_gap"
	DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak    = "derivatives_context_no_trade_filter_premise_audit_rejected_future_label_leak"
	DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue = "derivatives_context_no_trade_filter_premise_audit_rejected_closed_family_rescue"
	DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry      = "derivatives_context_no_trade_filter_premise_audit_rejected_rotation_entry_rescue"
	DerivativesNoTradeFilterPremiseStopStateFailedNoUsableFilter       = "derivatives_context_no_trade_filter_premise_audit_failed_no_usable_filter"
	DerivativesNoTradeFilterPremiseStopStatePassedNeedsSpec            = "derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec"

	derivativesNoTradeFilterID = "btc_15m_basis_discount_no_trade_veto_v1"
)

type FuturesDerivativesNoTradeFilterPremiseAuditConfig struct {
	SourceAuditConfig              FuturesDerivativesContextSourceAuditConfig
	BTCSource                      FuturesBTCRegimeETHSOLContextSourceConfig
	StateConfig                    FuturesRangeStateConstructionLoopAuditConfig
	Timeframe                      string
	HorizonBars                    int
	MinBasisContextCoverage        float64
	MinFullCandidateRows           int
	MinSplitCandidateRows          int
	MinFullToxicRate               float64
	MinWorstSplitToxicRate         float64
	MinFullToxicImprovement        float64
	MinUnionToxicRateFull          float64
	MinUnionToxicRateEverySplit    float64
	RequireExactMetricReproduction bool
	RescueClosedFamily             bool
	RescueRotationEntry            bool
}

type FuturesDerivativesNoTradeFilterPremiseAuditResult struct {
	SourceRows            []FuturesDerivativesNoTradeFilterPremiseSourceRow         `json:"source_rows"`
	CoverageRows          []FuturesDerivativesNoTradeFilterPremiseCoverageRow       `json:"coverage_rows"`
	FilterDefinitionRows  []FuturesDerivativesNoTradeFilterPremiseDefinitionRow     `json:"filter_definition_rows"`
	ExactCandidateRows    []FuturesDerivativesNoTradeFilterPremiseExactCandidateRow `json:"exact_candidate_rows"`
	CanonicalUnionRows    []FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow `json:"canonical_union_rows"`
	OverlapRows           []FuturesDerivativesNoTradeFilterPremiseOverlapRow        `json:"overlap_rows"`
	VetoCandidateRows     []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow  `json:"veto_candidate_rows"`
	CollateralRows        []FuturesDerivativesNoTradeFilterPremiseCollateralRow     `json:"collateral_rows"`
	MissingnessRows       []FuturesDerivativesNoTradeFilterPremiseMissingnessRow    `json:"missingness_rows"`
	SummaryRows           []FuturesDerivativesNoTradeFilterPremiseSummaryRow        `json:"summary_rows"`
	ExactCandidatesPassed int                                                       `json:"exact_candidates_passed"`
	CanonicalUnionPassed  bool                                                      `json:"canonical_union_passed"`
	StopState             string                                                    `json:"stop_state"`
}

type FuturesDerivativesNoTradeFilterPremiseSourceRow struct {
	AuditName                  string `json:"audit_name"`
	Role                       string `json:"role"`
	Symbol                     string `json:"symbol"`
	SourceFamily               string `json:"source_family"`
	ProductScope               string `json:"product_scope"`
	NativeInterval             string `json:"native_interval"`
	DurablePath                string `json:"durable_path"`
	Required                   bool   `json:"required"`
	RowCount                   int    `json:"row_count"`
	ExpectedRowCount           int    `json:"expected_row_count"`
	FirstOpenTime              string `json:"first_open_time"`
	LastOpenTime               string `json:"last_open_time"`
	GapCount                   int    `json:"gap_count"`
	DuplicateCount             int    `json:"duplicate_count"`
	ZeroVolumeCount            int    `json:"zero_volume_count,omitempty"`
	FileSHA256                 string `json:"file_sha256,omitempty"`
	FinalityRule               string `json:"finality_rule"`
	ValidationStatus           string `json:"validation_status"`
	ValidationError            string `json:"validation_error,omitempty"`
	ForwardLabelsAsSourceInput bool   `json:"forward_labels_as_source_input"`
}

type FuturesDerivativesNoTradeFilterPremiseCoverageRow struct {
	AuditName             string  `json:"audit_name"`
	Scope                 string  `json:"scope"`
	Symbol                string  `json:"symbol"`
	Timeframe             string  `json:"timeframe,omitempty"`
	SourceFamily          string  `json:"source_family,omitempty"`
	Required              bool    `json:"required"`
	RowCount              int     `json:"row_count"`
	StateRows             int     `json:"state_rows"`
	FeatureRows           int     `json:"feature_rows"`
	LabelRows             int     `json:"label_rows"`
	VetoCandidateRows     int     `json:"veto_candidate_rows"`
	MissingContextRows    int     `json:"missing_context_rows"`
	LagCoveragePct        float64 `json:"lag_coverage_pct"`
	ContextCoveragePct    float64 `json:"context_coverage_pct"`
	RequiredCoverageFloor float64 `json:"required_coverage_floor"`
	CoverageGatePass      bool    `json:"coverage_gate_pass"`
	ClosedCandleOnly      bool    `json:"closed_candle_only"`
	ForwardLabelsAsInputs bool    `json:"forward_labels_as_inputs"`
	NoFillPolicy          bool    `json:"no_fill_policy"`
	ValidationStatus      string  `json:"validation_status"`
	ValidationError       string  `json:"validation_error,omitempty"`
}

type FuturesDerivativesNoTradeFilterPremiseDefinitionRow struct {
	CandidateID                  string  `json:"candidate_id"`
	SelectedContextRank          int     `json:"selected_context_rank"`
	FilterScope                  string  `json:"filter_scope"`
	Symbol                       string  `json:"symbol"`
	Timeframe                    string  `json:"timeframe"`
	HorizonBars                  int     `json:"horizon_bars"`
	RouteCandidate               string  `json:"route_candidate"`
	LocalRangeBucketID           string  `json:"local_range_bucket_id"`
	DerivativesBucketType        string  `json:"derivatives_bucket_type"`
	DerivativesBucketID          string  `json:"derivatives_bucket_id"`
	BasisLevelBucket             string  `json:"basis_level_bucket"`
	BasisChangeBucket            string  `json:"basis_change_bucket,omitempty"`
	PremiumLevelBucket           string  `json:"premium_level_bucket,omitempty"`
	ExpectedRows                 int     `json:"expected_rows"`
	ExpectedWeakestSplitRows     int     `json:"expected_weakest_split_rows"`
	ExpectedFullToxicRate        float64 `json:"expected_full_toxic_rate"`
	ExpectedWorstSplitToxicRate  float64 `json:"expected_worst_split_toxic_rate"`
	ExpectedFullToxicImprovement float64 `json:"expected_full_toxic_improvement"`
	NestedInsideCandidateID      string  `json:"nested_inside_candidate_id,omitempty"`
	CorroboratorBound            bool    `json:"corroborator_bound"`
	RotationEntryAllowed         bool    `json:"rotation_entry_allowed"`
	ClosedFamilyRescueAllowed    bool    `json:"closed_family_rescue_allowed"`
	ForwardLabelsAsFilterInput   bool    `json:"forward_labels_as_filter_input"`
}

type FuturesDerivativesNoTradeFilterPremiseExactCandidateRow struct {
	CandidateID                    string  `json:"candidate_id"`
	Split                          string  `json:"split"`
	Symbol                         string  `json:"symbol"`
	Timeframe                      string  `json:"timeframe"`
	HorizonBars                    int     `json:"horizon_bars"`
	RouteCandidate                 string  `json:"route_candidate"`
	LocalRangeBucketID             string  `json:"local_range_bucket_id"`
	DerivativesBucketType          string  `json:"derivatives_bucket_type"`
	DerivativesBucketID            string  `json:"derivatives_bucket_id"`
	CandidateCount                 int     `json:"candidate_count"`
	NoTradeToxicCount              int     `json:"no_trade_toxic_count"`
	RotationUsefulBlockedCount     int     `json:"rotation_useful_blocked_count"`
	ContinuationUsefulBlockedCount int     `json:"continuation_useful_blocked_count"`
	DiagnosticOnlyCount            int     `json:"diagnostic_only_count"`
	ToxicRate                      float64 `json:"toxic_rate"`
	DominantForwardLabel           string  `json:"dominant_forward_label"`
	DominantForwardLabelRate       float64 `json:"dominant_forward_label_rate"`
	BaselineCandidateCount         int     `json:"baseline_candidate_count"`
	BaselineToxicRate              float64 `json:"baseline_toxic_rate"`
	FullRows                       int     `json:"full_rows"`
	WeakestSplitRows               int     `json:"weakest_split_rows"`
	MaxSplitContributionRate       float64 `json:"max_split_contribution_rate"`
	FullToxicRate                  float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate            float64 `json:"worst_split_toxic_rate"`
	MinSplitToxicRate              float64 `json:"min_split_toxic_rate"`
	FullToxicImprovement           float64 `json:"full_toxic_improvement"`
	MinSplitToxicImprovement       float64 `json:"min_split_toxic_improvement"`
	ExpectedRows                   int     `json:"expected_rows"`
	ExpectedWeakestSplitRows       int     `json:"expected_weakest_split_rows"`
	ExpectedFullToxicRate          float64 `json:"expected_full_toxic_rate"`
	ExpectedWorstSplitToxicRate    float64 `json:"expected_worst_split_toxic_rate"`
	ExpectedFullToxicImprovement   float64 `json:"expected_full_toxic_improvement"`
	ExactMetricReproductionPass    bool    `json:"exact_metric_reproduction_pass"`
	ReviewableCountGatePass        bool    `json:"reviewable_count_gate_pass"`
	ToxicRateGatePass              bool    `json:"toxic_rate_gate_pass"`
	ImprovementGatePass            bool    `json:"improvement_gate_pass"`
	MissingnessGatePass            bool    `json:"missingness_gate_pass"`
	FutureLeakProtectionPass       bool    `json:"future_leak_protection_pass"`
	ClosedFamilyProtectionPass     bool    `json:"closed_family_protection_pass"`
	RotationEntryProtectionPass    bool    `json:"rotation_entry_protection_pass"`
	PassesGate                     bool    `json:"passes_gate"`
	FailureReason                  string  `json:"failure_reason,omitempty"`
}

type FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow struct {
	FilterID                          string  `json:"filter_id"`
	Split                             string  `json:"split"`
	Symbol                            string  `json:"symbol"`
	Timeframe                         string  `json:"timeframe"`
	HorizonBars                       int     `json:"horizon_bars"`
	MatchedExactCandidates            string  `json:"matched_exact_candidates"`
	SumExactCandidateRows             int     `json:"sum_exact_candidate_rows"`
	DeduplicatedRows                  int     `json:"deduplicated_rows"`
	OverlapRows                       int     `json:"overlap_rows"`
	NestedTrendDownPremiumOverlapRows int     `json:"nested_trend_down_premium_overlap_rows"`
	NoTradeToxicCount                 int     `json:"no_trade_toxic_count"`
	RotationUsefulBlockedCount        int     `json:"rotation_useful_blocked_count"`
	ContinuationUsefulBlockedCount    int     `json:"continuation_useful_blocked_count"`
	DiagnosticOnlyCount               int     `json:"diagnostic_only_count"`
	ToxicRate                         float64 `json:"toxic_rate"`
	FullRows                          int     `json:"full_rows"`
	WeakestSplitRows                  int     `json:"weakest_split_rows"`
	MaxSplitContributionRate          float64 `json:"max_split_contribution_rate"`
	FullToxicRate                     float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate               float64 `json:"worst_split_toxic_rate"`
	MinSplitToxicRate                 float64 `json:"min_split_toxic_rate"`
	LocalOnlyBaselineRows             int     `json:"local_only_baseline_rows"`
	LocalOnlyBaselineToxicRate        float64 `json:"local_only_baseline_toxic_rate"`
	FullToxicImprovement              float64 `json:"full_toxic_improvement"`
	ToxicDominatedEverySplit          bool    `json:"toxic_dominated_every_split"`
	DoubleCountingProtectionPass      bool    `json:"double_counting_protection_pass"`
	CollateralReported                bool    `json:"collateral_reported"`
	MissingnessGatePass               bool    `json:"missingness_gate_pass"`
	FutureLeakProtectionPass          bool    `json:"future_leak_protection_pass"`
	ClosedFamilyProtectionPass        bool    `json:"closed_family_protection_pass"`
	RotationEntryProtectionPass       bool    `json:"rotation_entry_protection_pass"`
	PassesGate                        bool    `json:"passes_gate"`
	FailureReason                     string  `json:"failure_reason,omitempty"`
}

type FuturesDerivativesNoTradeFilterPremiseOverlapRow struct {
	Split                string `json:"split"`
	CandidateID          string `json:"candidate_id"`
	OtherCandidateID     string `json:"other_candidate_id"`
	OverlapRows          int    `json:"overlap_rows"`
	CandidateRows        int    `json:"candidate_rows"`
	OtherCandidateRows   int    `json:"other_candidate_rows"`
	NestedOverlap        bool   `json:"nested_overlap"`
	DoubleCountedInUnion bool   `json:"double_counted_in_union"`
	OverlapPolicy        string `json:"overlap_policy"`
}

type FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow struct {
	VetoRowID                     int     `json:"veto_row_id"`
	Symbol                        string  `json:"symbol"`
	Split                         string  `json:"split"`
	Timeframe                     string  `json:"timeframe"`
	HorizonBars                   int     `json:"horizon_bars"`
	LocalStateRowID               int     `json:"local_state_row_id"`
	FeatureRowID                  int     `json:"feature_row_id"`
	DecisionCloseTime             string  `json:"decision_close_time"`
	SourceCloseTime               string  `json:"source_close_time"`
	LocalRangeBucketID            string  `json:"local_range_bucket_id"`
	BasisLevelBucket              string  `json:"basis_level_bucket"`
	BasisChangeBucket             string  `json:"basis_change_bucket"`
	PremiumLevelBucket            string  `json:"premium_level_bucket"`
	BasisBPS                      float64 `json:"basis_bps"`
	PremiumBPS                    float64 `json:"premium_bps"`
	MatchedCandidateIDs           string  `json:"matched_candidate_ids"`
	CandidateMatchCount           int     `json:"candidate_match_count"`
	NestedTrendDownPremiumOverlap bool    `json:"nested_trend_down_premium_overlap"`
	ForwardLabel                  string  `json:"forward_label"`
	NoTradeToxic                  bool    `json:"no_trade_toxic"`
	RotationUsefulBlocked         bool    `json:"rotation_useful_blocked"`
	ContinuationUsefulBlocked     bool    `json:"continuation_useful_blocked"`
	DiagnosticOnly                bool    `json:"diagnostic_only"`
	ForwardLabelMetadataOnly      bool    `json:"forward_label_metadata_only"`
	ForwardLabelUsedAsFeature     bool    `json:"forward_label_used_as_feature"`
	ClosedCandleOnly              bool    `json:"closed_candle_only"`
	UsesFutureRows                bool    `json:"uses_future_rows"`
}

type FuturesDerivativesNoTradeFilterPremiseCollateralRow struct {
	Split                         string  `json:"split"`
	ForwardLabel                  string  `json:"forward_label"`
	BlockedRows                   int     `json:"blocked_rows"`
	ShareOfBlockedRows            float64 `json:"share_of_blocked_rows"`
	NoTradeToxicRows              int     `json:"no_trade_toxic_rows"`
	RotationUsefulBlockedRows     int     `json:"rotation_useful_blocked_rows"`
	ContinuationUsefulBlockedRows int     `json:"continuation_useful_blocked_rows"`
	DiagnosticOnlyRows            int     `json:"diagnostic_only_rows"`
	CollateralDamageReported      bool    `json:"collateral_damage_reported"`
}

type FuturesDerivativesNoTradeFilterPremiseMissingnessRow struct {
	AuditName         string  `json:"audit_name"`
	Scope             string  `json:"scope"`
	Symbol            string  `json:"symbol"`
	Timeframe         string  `json:"timeframe,omitempty"`
	SourceFamily      string  `json:"source_family,omitempty"`
	Reason            string  `json:"reason"`
	Count             int     `json:"count"`
	TotalRows         int     `json:"total_rows"`
	Rate              float64 `json:"rate"`
	ForwardFilledRows int     `json:"forward_filled_rows"`
	MissingDataPolicy string  `json:"missing_data_policy"`
	CoverageFloor     float64 `json:"coverage_floor"`
	CoverageGatePass  bool    `json:"coverage_gate_pass"`
}

type FuturesDerivativesNoTradeFilterPremiseSummaryRow struct {
	Split                              string  `json:"split"`
	FilterID                           string  `json:"filter_id"`
	Symbol                             string  `json:"symbol"`
	Timeframe                          string  `json:"timeframe"`
	HorizonBars                        int     `json:"horizon_bars"`
	SourceRows                         int     `json:"source_rows"`
	CoverageRows                       int     `json:"coverage_rows"`
	FilterDefinitionRows               int     `json:"filter_definition_rows"`
	ExactCandidateRows                 int     `json:"exact_candidate_rows"`
	ExactCandidatesPassed              int     `json:"exact_candidates_passed"`
	VetoCandidateRows                  int     `json:"veto_candidate_rows"`
	CanonicalDeduplicatedRows          int     `json:"canonical_deduplicated_rows"`
	CanonicalOverlapRows               int     `json:"canonical_overlap_rows"`
	NestedTrendDownPremiumOverlapRows  int     `json:"nested_trend_down_premium_overlap_rows"`
	CanonicalNoTradeToxicRows          int     `json:"canonical_no_trade_toxic_rows"`
	CanonicalToxicRate                 float64 `json:"canonical_toxic_rate"`
	CanonicalMinSplitToxicRate         float64 `json:"canonical_min_split_toxic_rate"`
	CanonicalRotationUsefulBlocked     int     `json:"canonical_rotation_useful_blocked"`
	CanonicalContinuationUsefulBlocked int     `json:"canonical_continuation_useful_blocked"`
	CollateralRows                     int     `json:"collateral_rows"`
	MissingnessRows                    int     `json:"missingness_rows"`
	RequiredCoverageFloor              float64 `json:"required_coverage_floor"`
	MinLagCoveragePct                  float64 `json:"min_lag_coverage_pct"`
	SourceScopePass                    bool    `json:"source_scope_pass"`
	CoveragePass                       bool    `json:"coverage_pass"`
	ExactRowsReproduced                bool    `json:"exact_rows_reproduced"`
	CanonicalUnionPass                 bool    `json:"canonical_union_pass"`
	CommonOutputsZeroTrade             bool    `json:"common_outputs_zero_trade"`
	ForwardLabelsAsInputs              bool    `json:"forward_labels_as_inputs"`
	ClosedFamilyProtectionPass         bool    `json:"closed_family_protection_pass"`
	RotationEntryProtectionPass        bool    `json:"rotation_entry_protection_pass"`
	Trades                             int     `json:"trades"`
	StopState                          string  `json:"stop_state"`
}

type derivativesNoTradeFilterPremiseLabelWithFeature struct {
	label   FuturesDerivativesContextAuditLabelRow
	feature FuturesDerivativesContextAuditBasisFeatureRow
}

type derivativesNoTradeFilterPremiseMetricAccumulator struct {
	row    FuturesDerivativesNoTradeFilterPremiseExactCandidateRow
	labels map[string]int
}

type derivativesNoTradeFilterPremiseUnionAccumulator struct {
	row    FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow
	labels map[string]int
}

func DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig() FuturesDerivativesNoTradeFilterPremiseAuditConfig {
	sourceCfg := derivativesNoTradeFilterPremiseSourceConfig()
	stateCfg := DefaultFuturesRangeStateConstructionLoopAuditConfig()
	return FuturesDerivativesNoTradeFilterPremiseAuditConfig{
		SourceAuditConfig:              sourceCfg,
		BTCSource:                      derivativesNoTradeFilterPremiseBTCSource(),
		StateConfig:                    stateCfg,
		Timeframe:                      RangeDiscoveryTimeframe15m,
		HorizonBars:                    48,
		MinBasisContextCoverage:        0.994472,
		MinFullCandidateRows:           300,
		MinSplitCandidateRows:          60,
		MinFullToxicRate:               0.65,
		MinWorstSplitToxicRate:         0.69,
		MinFullToxicImprovement:        0.04,
		MinUnionToxicRateFull:          0.50,
		MinUnionToxicRateEverySplit:    0.50,
		RequireExactMetricReproduction: true,
	}
}

func RunFuturesDerivativesNoTradeFilterPremiseAudit(cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) (FuturesDerivativesNoTradeFilterPremiseAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesDerivativesNoTradeFilterPremiseAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesDerivativesNoTradeFilterPremiseAuditResult{}
	result.FilterDefinitionRows = derivativesNoTradeFilterPremiseDefinitions(cfg)
	if cfg.RescueClosedFamily {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}
	if cfg.RescueRotationEntry {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	sourceAudit, err := RunFuturesDerivativesContextSourceAudit(cfg.SourceAuditConfig, splits)
	if err != nil {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateSourceGap
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}
	result.SourceRows = derivativesNoTradeFilterPremiseSourceRows(sourceAudit)
	result.CoverageRows = append(result.CoverageRows, derivativesNoTradeFilterPremiseSourceCoverageRows(sourceAudit, cfg)...)
	result.MissingnessRows = append(result.MissingnessRows, derivativesNoTradeFilterPremiseSourceMissingnessRows(sourceAudit, cfg)...)
	sourcePass := derivativesNoTradeFilterPremiseSourcePass(sourceAudit, cfg)
	if sourceAudit.StopState != DerivativesContextSourceAuditStopStatePassedNeedsContextBrief || !sourcePass {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateSourceGap
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	btcCfg := FuturesBTCRegimeETHSOLContextAuditConfig{Sources: []FuturesBTCRegimeETHSOLContextSourceConfig{cfg.BTCSource}, StateConfig: cfg.StateConfig}
	data := btcRegimeETHSOLContextLoadSymbol(cfg.BTCSource, btcCfg, splits)
	result.SourceRows = append(result.SourceRows, derivativesNoTradeFilterPremiseAnchorSourceRow(data.source))
	if !data.sourceOK {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateSourceGap
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}
	built, err := btcRegimeETHSOLContextBuildSymbolStates(data, btcCfg, splits)
	if err != nil {
		return result, err
	}
	for _, cov := range built.coverage {
		result.CoverageRows = append(result.CoverageRows, derivativesNoTradeFilterPremiseLocalCoverageRow(cov, cfg))
	}
	if !built.coverageOK {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateSourceGap
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	streams := derivativesContextLoadStreams(FuturesDerivativesContextAuditConfig{SourceAuditConfig: cfg.SourceAuditConfig})
	featureByState := map[int]FuturesDerivativesContextAuditBasisFeatureRow{}
	missingContext := 0
	stateRows := 0
	for _, state := range built.states {
		if state.Timeframe != cfg.Timeframe {
			continue
		}
		stateRows++
		feature, ok, _ := derivativesContextFeatureForState(len(featureByState)+1, RangeUniverseSymbolBTCUSDT, state, streams, FuturesDerivativesContextAuditConfig{
			SourceAuditConfig:                cfg.SourceAuditConfig,
			BasisChangeLookbackIntervals:     12,
			BasisVolatilityLookbackIntervals: 36,
		}.withDefaults())
		if !ok {
			missingContext++
			continue
		}
		featureByState[state.StateRowID] = feature
	}
	labels := []derivativesNoTradeFilterPremiseLabelWithFeature{}
	for _, label := range built.labels {
		if label.Timeframe != cfg.Timeframe || label.HorizonBars != cfg.HorizonBars {
			continue
		}
		feature, ok := featureByState[label.StateRowID]
		if !ok {
			continue
		}
		row := derivativesContextLabelRow(len(labels)+1, RangeUniverseSymbolBTCUSDT, feature, label)
		labels = append(labels, derivativesNoTradeFilterPremiseLabelWithFeature{label: row, feature: feature})
	}
	contextCoverage := safeDiv(float64(len(featureByState)), float64(stateRows))
	result.CoverageRows = append(result.CoverageRows, FuturesDerivativesNoTradeFilterPremiseCoverageRow{
		AuditName:             FuturesDerivativesNoTradeFilterPremiseAuditName,
		Scope:                 "lagged_basis_context",
		Symbol:                RangeUniverseSymbolBTCUSDT,
		Timeframe:             cfg.Timeframe,
		SourceFamily:          "mark_minus_index_basis",
		Required:              true,
		StateRows:             stateRows,
		FeatureRows:           len(featureByState),
		LabelRows:             len(labels),
		MissingContextRows:    missingContext,
		LagCoveragePct:        derivativesNoTradeFilterPremiseMinLagCoverage(sourceAudit),
		ContextCoveragePct:    contextCoverage,
		RequiredCoverageFloor: cfg.MinBasisContextCoverage,
		CoverageGatePass:      derivativesNoTradeFilterPremiseMinLagCoverage(sourceAudit)+1e-6 >= cfg.MinBasisContextCoverage,
		ClosedCandleOnly:      true,
		ForwardLabelsAsInputs: false,
		NoFillPolicy:          true,
		ValidationStatus:      "accepted",
	})
	result.MissingnessRows = append(result.MissingnessRows, FuturesDerivativesNoTradeFilterPremiseMissingnessRow{
		AuditName:         FuturesDerivativesNoTradeFilterPremiseAuditName,
		Scope:             "derived_context",
		Symbol:            RangeUniverseSymbolBTCUSDT,
		Timeframe:         cfg.Timeframe,
		SourceFamily:      "mark_minus_index_basis",
		Reason:            "missing_required_lagged_basis_context",
		Count:             missingContext,
		TotalRows:         stateRows,
		Rate:              safeDiv(float64(missingContext), float64(stateRows)),
		ForwardFilledRows: 0,
		MissingDataPolicy: derivativesMissingPolicy,
		CoverageFloor:     cfg.MinBasisContextCoverage,
		CoverageGatePass:  derivativesNoTradeFilterPremiseMinLagCoverage(sourceAudit)+1e-6 >= cfg.MinBasisContextCoverage,
	})
	if !derivativesNoTradeFilterPremiseCoveragePass(result.CoverageRows) {
		result.StopState = DerivativesNoTradeFilterPremiseStopStateSourceGap
		result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
		return result, nil
	}

	derivativesNoTradeFilterPremiseEvaluate(&result, labels, cfg, splits)
	result.StopState = FuturesDerivativesNoTradeFilterPremiseAuditStopState(result)
	result.SummaryRows = derivativesNoTradeFilterPremiseSummaryRows(result, cfg, splits)
	return result, nil
}

func FuturesDerivativesNoTradeFilterPremiseAuditStopState(result FuturesDerivativesNoTradeFilterPremiseAuditResult) string {
	if result.StopState == DerivativesNoTradeFilterPremiseStopStateSourceGap ||
		result.StopState == DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak ||
		result.StopState == DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue ||
		result.StopState == DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry {
		return result.StopState
	}
	if !derivativesNoTradeFilterPremiseCoveragePass(result.CoverageRows) {
		return DerivativesNoTradeFilterPremiseStopStateSourceGap
	}
	for _, row := range result.ExactCandidateRows {
		if !row.FutureLeakProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak
		}
		if !row.ClosedFamilyProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue
		}
		if !row.RotationEntryProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry
		}
	}
	for _, row := range result.VetoCandidateRows {
		if row.ForwardLabelUsedAsFeature {
			return DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak
		}
		if row.UsesFutureRows {
			return DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak
		}
	}
	for _, row := range result.CanonicalUnionRows {
		if !row.FutureLeakProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak
		}
		if !row.ClosedFamilyProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue
		}
		if !row.RotationEntryProtectionPass {
			return DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry
		}
	}
	if result.ExactCandidatesPassed == len(result.FilterDefinitionRows) && result.CanonicalUnionPassed {
		return DerivativesNoTradeFilterPremiseStopStatePassedNeedsSpec
	}
	return DerivativesNoTradeFilterPremiseStopStateFailedNoUsableFilter
}

func (cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) withDefaults() FuturesDerivativesNoTradeFilterPremiseAuditConfig {
	defaults := DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
	if len(cfg.SourceAuditConfig.DerivativeSources) == 0 {
		cfg.SourceAuditConfig = defaults.SourceAuditConfig
	}
	if len(cfg.BTCSource.Symbol) == 0 {
		cfg.BTCSource = defaults.BTCSource
	}
	if len(cfg.StateConfig.Timeframes) == 0 {
		cfg.StateConfig = defaults.StateConfig
	}
	if cfg.Timeframe == "" {
		cfg.Timeframe = defaults.Timeframe
	}
	if cfg.HorizonBars == 0 {
		cfg.HorizonBars = defaults.HorizonBars
	}
	if cfg.MinBasisContextCoverage == 0 {
		cfg.MinBasisContextCoverage = defaults.MinBasisContextCoverage
	}
	if cfg.MinFullCandidateRows == 0 {
		cfg.MinFullCandidateRows = defaults.MinFullCandidateRows
	}
	if cfg.MinSplitCandidateRows == 0 {
		cfg.MinSplitCandidateRows = defaults.MinSplitCandidateRows
	}
	if cfg.MinFullToxicRate == 0 {
		cfg.MinFullToxicRate = defaults.MinFullToxicRate
	}
	if cfg.MinWorstSplitToxicRate == 0 {
		cfg.MinWorstSplitToxicRate = defaults.MinWorstSplitToxicRate
	}
	if cfg.MinFullToxicImprovement == 0 {
		cfg.MinFullToxicImprovement = defaults.MinFullToxicImprovement
	}
	if cfg.MinUnionToxicRateFull == 0 {
		cfg.MinUnionToxicRateFull = defaults.MinUnionToxicRateFull
	}
	if cfg.MinUnionToxicRateEverySplit == 0 {
		cfg.MinUnionToxicRateEverySplit = defaults.MinUnionToxicRateEverySplit
	}
	return cfg
}

func (cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) validate() error {
	if cfg.Timeframe != RangeDiscoveryTimeframe15m || cfg.HorizonBars != 48 {
		return fmt.Errorf("derivatives no-trade filter premise audit requires BTCUSDT 15m h48 scope")
	}
	if strings.ToUpper(strings.TrimSpace(cfg.BTCSource.Symbol)) != RangeUniverseSymbolBTCUSDT {
		return fmt.Errorf("derivatives no-trade filter premise audit requires BTCUSDT candle source")
	}
	for _, source := range cfg.SourceAuditConfig.DerivativeSources {
		if strings.ToUpper(strings.TrimSpace(source.Symbol)) != RangeUniverseSymbolBTCUSDT {
			return fmt.Errorf("derivatives no-trade filter premise audit forbids non-BTCUSDT derivative source %q", source.Symbol)
		}
	}
	for _, anchor := range cfg.SourceAuditConfig.Anchors {
		if strings.ToUpper(strings.TrimSpace(anchor.Symbol)) != RangeUniverseSymbolBTCUSDT {
			return fmt.Errorf("derivatives no-trade filter premise audit forbids non-BTCUSDT anchor %q", anchor.Symbol)
		}
	}
	if cfg.MinBasisContextCoverage <= 0 || cfg.MinBasisContextCoverage > 1 {
		return fmt.Errorf("basis context coverage floor must be in (0,1]")
	}
	return nil
}

func derivativesNoTradeFilterPremiseSourceConfig() FuturesDerivativesContextSourceAuditConfig {
	cfg := DefaultFuturesDerivativesContextSourceAuditConfig()
	filtered := make([]FuturesDerivativesContextSourceFileConfig, 0, 3)
	for _, source := range cfg.DerivativeSources {
		if strings.ToUpper(strings.TrimSpace(source.Symbol)) != RangeUniverseSymbolBTCUSDT {
			continue
		}
		source.Required = true
		if source.SourceFamily == "premium_index_klines" {
			source.AllowNonPositive = true
		}
		filtered = append(filtered, source)
	}
	cfg.DerivativeSources = filtered
	cfg.Anchors = []FuturesRangeUniverseSourceConfig{{
		Symbol:       RangeUniverseSymbolBTCUSDT,
		Path:         "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ApprovedPath: "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
	}}
	return cfg
}

func derivativesNoTradeFilterPremiseBTCSource() FuturesBTCRegimeETHSOLContextSourceConfig {
	base := DefaultFuturesBTCRegimeETHSOLContextAuditConfig().withDefaults()
	for _, source := range base.Sources {
		if strings.ToUpper(strings.TrimSpace(source.Symbol)) == RangeUniverseSymbolBTCUSDT {
			return source
		}
	}
	return FuturesBTCRegimeETHSOLContextSourceConfig{
		Symbol:                  RangeUniverseSymbolBTCUSDT,
		Path:                    "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ApprovedPath:            "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		ExpectedRowCount:        573984,
		ExpectedFirstOpenTime:   "2021-01-01T00:00:00Z",
		ExpectedLastOpenTime:    "2026-06-16T23:55:00Z",
		ExpectedGapCount:        0,
		ExpectedDuplicateCount:  0,
		ExpectedZeroVolumeCount: 66,
	}
}

func derivativesNoTradeFilterPremiseDefinitions(cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) []FuturesDerivativesNoTradeFilterPremiseDefinitionRow {
	local := func(trend string) string {
		return strings.Join([]string{cfg.Timeframe, "geometry_midline_balanced", "vol_compressed", trend, "impulse_none"}, "::")
	}
	rows := []FuturesDerivativesNoTradeFilterPremiseDefinitionRow{
		{
			CandidateID: "exact_1_trend_down_basis_discount_premium_discount", SelectedContextRank: 1, FilterScope: "exact_context_row",
			Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic,
			LocalRangeBucketID: local("trend_down_pressure"), DerivativesBucketType: DerivativesContextBucketBasisLevelPremium,
			DerivativesBucketID: "basis=basis_discount_small::premium=premium_discount_small", BasisLevelBucket: "basis_discount_small", PremiumLevelBucket: "premium_discount_small",
			ExpectedRows: 515, ExpectedWeakestSplitRows: 110, ExpectedFullToxicRate: 0.732039, ExpectedWorstSplitToxicRate: 0.800000, ExpectedFullToxicImprovement: 0.049738,
			NestedInsideCandidateID: "exact_2_trend_down_basis_discount", CorroboratorBound: true,
		},
		{
			CandidateID: "exact_2_trend_down_basis_discount", SelectedContextRank: 2, FilterScope: "exact_context_row",
			Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic,
			LocalRangeBucketID: local("trend_down_pressure"), DerivativesBucketType: DerivativesContextBucketBasisLevel,
			DerivativesBucketID: "basis=basis_discount_small", BasisLevelBucket: "basis_discount_small",
			ExpectedRows: 622, ExpectedWeakestSplitRows: 142, ExpectedFullToxicRate: 0.729904, ExpectedWorstSplitToxicRate: 0.802817, ExpectedFullToxicImprovement: 0.047603,
		},
		{
			CandidateID: "exact_3_trend_flat_basis_discount_change_flat", SelectedContextRank: 3, FilterScope: "exact_context_row",
			Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic,
			LocalRangeBucketID: local("trend_flat"), DerivativesBucketType: DerivativesContextBucketBasisLevelChange,
			DerivativesBucketID: "basis=basis_discount_small::change=basis_change_flat", BasisLevelBucket: "basis_discount_small", BasisChangeBucket: "basis_change_flat",
			ExpectedRows: 356, ExpectedWeakestSplitRows: 62, ExpectedFullToxicRate: 0.662921, ExpectedWorstSplitToxicRate: 0.699387, ExpectedFullToxicImprovement: 0.045126,
			CorroboratorBound: true,
		},
		{
			CandidateID: "exact_4_trend_up_basis_discount", SelectedContextRank: 4, FilterScope: "exact_context_row",
			Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic,
			LocalRangeBucketID: local("trend_up_pressure"), DerivativesBucketType: DerivativesContextBucketBasisLevel,
			DerivativesBucketID: "basis=basis_discount_small", BasisLevelBucket: "basis_discount_small",
			ExpectedRows: 613, ExpectedWeakestSplitRows: 124, ExpectedFullToxicRate: 0.654160, ExpectedWorstSplitToxicRate: 0.759358, ExpectedFullToxicImprovement: 0.050840,
		},
		{
			CandidateID: "exact_5_trend_flat_basis_discount_premium_discount", SelectedContextRank: 5, FilterScope: "exact_context_row",
			Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic,
			LocalRangeBucketID: local("trend_flat"), DerivativesBucketType: DerivativesContextBucketBasisLevelPremium,
			DerivativesBucketID: "basis=basis_discount_small::premium=premium_discount_small", BasisLevelBucket: "basis_discount_small", PremiumLevelBucket: "premium_discount_small",
			ExpectedRows: 538, ExpectedWeakestSplitRows: 115, ExpectedFullToxicRate: 0.659851, ExpectedWorstSplitToxicRate: 0.719212, ExpectedFullToxicImprovement: 0.042056,
			CorroboratorBound: true,
		},
	}
	for i := range rows {
		rows[i].RotationEntryAllowed = false
		rows[i].ClosedFamilyRescueAllowed = false
		rows[i].ForwardLabelsAsFilterInput = false
	}
	return rows
}

func derivativesNoTradeFilterPremiseEvaluate(result *FuturesDerivativesNoTradeFilterPremiseAuditResult, labels []derivativesNoTradeFilterPremiseLabelWithFeature, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) {
	defs := result.FilterDefinitionRows
	exactAcc := map[string]*derivativesNoTradeFilterPremiseMetricAccumulator{}
	baselineAcc := map[string]*derivativesNoTradeFilterPremiseMetricAccumulator{}
	unionByKey := map[string]derivativesNoTradeFilterPremiseLabelWithFeature{}
	unionMatchesByKey := map[string][]string{}
	for _, item := range labels {
		label := item.label
		for _, split := range rangeDiscoverySplitCombos(label.Split) {
			baseKey := label.LocalRangeBucketID + "|" + split
			acc := baselineAcc[baseKey]
			if acc == nil {
				acc = &derivativesNoTradeFilterPremiseMetricAccumulator{labels: map[string]int{}}
				baselineAcc[baseKey] = acc
			}
			derivativesNoTradeFilterPremiseAddLabel(acc, label)
		}
		matches := derivativesNoTradeFilterPremiseMatchingDefinitions(label, defs)
		if len(matches) == 0 {
			continue
		}
		matchIDs := make([]string, 0, len(matches))
		for _, def := range matches {
			matchIDs = append(matchIDs, def.CandidateID)
			for _, split := range rangeDiscoverySplitCombos(label.Split) {
				key := def.CandidateID + "|" + split
				acc := exactAcc[key]
				if acc == nil {
					acc = &derivativesNoTradeFilterPremiseMetricAccumulator{labels: map[string]int{}}
					acc.row = FuturesDerivativesNoTradeFilterPremiseExactCandidateRow{
						CandidateID: def.CandidateID, Split: split, Symbol: def.Symbol, Timeframe: def.Timeframe, HorizonBars: def.HorizonBars,
						RouteCandidate: def.RouteCandidate, LocalRangeBucketID: def.LocalRangeBucketID, DerivativesBucketType: def.DerivativesBucketType,
						DerivativesBucketID: def.DerivativesBucketID, ExpectedRows: def.ExpectedRows, ExpectedWeakestSplitRows: def.ExpectedWeakestSplitRows,
						ExpectedFullToxicRate: def.ExpectedFullToxicRate, ExpectedWorstSplitToxicRate: def.ExpectedWorstSplitToxicRate,
						ExpectedFullToxicImprovement: def.ExpectedFullToxicImprovement,
					}
					exactAcc[key] = acc
				}
				derivativesNoTradeFilterPremiseAddLabel(acc, label)
			}
		}
		sort.Strings(matchIDs)
		unionKey := fmt.Sprintf("%s|%s|%d", label.Symbol, label.Timeframe, label.LocalStateRowID)
		unionByKey[unionKey] = item
		unionMatchesByKey[unionKey] = matchIDs
	}

	result.ExactCandidateRows = derivativesNoTradeFilterPremiseExactRows(exactAcc, baselineAcc, defs, cfg, splits)
	for _, row := range result.ExactCandidateRows {
		if row.Split == fullSplitName && row.PassesGate {
			result.ExactCandidatesPassed++
		}
	}
	result.VetoCandidateRows = derivativesNoTradeFilterPremiseVetoRows(unionByKey, unionMatchesByKey)
	result.OverlapRows = derivativesNoTradeFilterPremiseOverlapRows(result.VetoCandidateRows, result.ExactCandidateRows, defs, splits)
	result.CanonicalUnionRows = derivativesNoTradeFilterPremiseUnionRows(result.VetoCandidateRows, labels, result.ExactCandidateRows, cfg, splits)
	for _, row := range result.CanonicalUnionRows {
		if row.Split == fullSplitName && row.PassesGate {
			result.CanonicalUnionPassed = true
		}
	}
	result.CollateralRows = derivativesNoTradeFilterPremiseCollateralRows(result.VetoCandidateRows, splits)
	for i := range result.CanonicalUnionRows {
		result.CanonicalUnionRows[i].CollateralReported = len(result.CollateralRows) > 0
		if !result.CanonicalUnionRows[i].CollateralReported {
			result.CanonicalUnionRows[i].PassesGate = false
			result.CanonicalUnionRows[i].FailureReason = uniqueJoinedReasons([]string{result.CanonicalUnionRows[i].FailureReason, "collateral_damage_not_reported"})
		}
	}
}

func derivativesNoTradeFilterPremiseAddLabel(acc *derivativesNoTradeFilterPremiseMetricAccumulator, label FuturesDerivativesContextAuditLabelRow) {
	acc.row.CandidateCount++
	acc.labels[label.ForwardLabel]++
	if label.NoTradeToxic {
		acc.row.NoTradeToxicCount++
	}
	if label.RotationUseful {
		acc.row.RotationUsefulBlockedCount++
	}
	if label.ContinuationUseful {
		acc.row.ContinuationUsefulBlockedCount++
	}
	if label.DiagnosticOnly {
		acc.row.DiagnosticOnlyCount++
	}
}

func derivativesNoTradeFilterPremiseMatchingDefinitions(label FuturesDerivativesContextAuditLabelRow, defs []FuturesDerivativesNoTradeFilterPremiseDefinitionRow) []FuturesDerivativesNoTradeFilterPremiseDefinitionRow {
	matches := []FuturesDerivativesNoTradeFilterPremiseDefinitionRow{}
	for _, def := range defs {
		if label.Symbol != def.Symbol || label.Timeframe != def.Timeframe || label.HorizonBars != def.HorizonBars || label.LocalRangeBucketID != def.LocalRangeBucketID {
			continue
		}
		if label.BasisLevelBucket != def.BasisLevelBucket {
			continue
		}
		if def.BasisChangeBucket != "" && label.BasisChangeBucket != def.BasisChangeBucket {
			continue
		}
		if def.PremiumLevelBucket != "" && label.PremiumLevelBucket != def.PremiumLevelBucket {
			continue
		}
		matches = append(matches, def)
	}
	return matches
}

func derivativesNoTradeFilterPremiseExactRows(exactAcc, baselineAcc map[string]*derivativesNoTradeFilterPremiseMetricAccumulator, defs []FuturesDerivativesNoTradeFilterPremiseDefinitionRow, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) []FuturesDerivativesNoTradeFilterPremiseExactCandidateRow {
	rows := []FuturesDerivativesNoTradeFilterPremiseExactCandidateRow{}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	byCandidateSplit := map[string]map[string]FuturesDerivativesNoTradeFilterPremiseExactCandidateRow{}
	for _, def := range defs {
		for _, split := range append([]string{fullSplitName}, derivativesNoTradeFilterPremiseSplitNames(periodSplits)...) {
			key := def.CandidateID + "|" + split
			row := FuturesDerivativesNoTradeFilterPremiseExactCandidateRow{
				CandidateID: def.CandidateID, Split: split, Symbol: def.Symbol, Timeframe: def.Timeframe, HorizonBars: def.HorizonBars,
				RouteCandidate: def.RouteCandidate, LocalRangeBucketID: def.LocalRangeBucketID, DerivativesBucketType: def.DerivativesBucketType,
				DerivativesBucketID: def.DerivativesBucketID, ExpectedRows: def.ExpectedRows, ExpectedWeakestSplitRows: def.ExpectedWeakestSplitRows,
				ExpectedFullToxicRate: def.ExpectedFullToxicRate, ExpectedWorstSplitToxicRate: def.ExpectedWorstSplitToxicRate,
				ExpectedFullToxicImprovement: def.ExpectedFullToxicImprovement,
				FutureLeakProtectionPass:     true, ClosedFamilyProtectionPass: true, RotationEntryProtectionPass: true, MissingnessGatePass: true,
			}
			if acc := exactAcc[key]; acc != nil {
				row = acc.row
				row.ToxicRate = safeDiv(float64(row.NoTradeToxicCount), float64(row.CandidateCount))
				row.DominantForwardLabel, row.DominantForwardLabelRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, false)
			}
			if base := baselineAcc[def.LocalRangeBucketID+"|"+split]; base != nil {
				row.BaselineCandidateCount = base.row.CandidateCount
				row.BaselineToxicRate = safeDiv(float64(base.row.NoTradeToxicCount), float64(base.row.CandidateCount))
			}
			if byCandidateSplit[def.CandidateID] == nil {
				byCandidateSplit[def.CandidateID] = map[string]FuturesDerivativesNoTradeFilterPremiseExactCandidateRow{}
			}
			byCandidateSplit[def.CandidateID][split] = row
		}
	}
	for _, def := range defs {
		candidateRows := byCandidateSplit[def.CandidateID]
		full := candidateRows[fullSplitName]
		full.FullRows = full.CandidateCount
		full.FullToxicRate = full.ToxicRate
		full.FullToxicImprovement = full.ToxicRate - full.BaselineToxicRate
		full.WeakestSplitRows = math.MaxInt
		full.MinSplitToxicRate = math.Inf(1)
		full.MinSplitToxicImprovement = math.Inf(1)
		reasons := []string{}
		for _, split := range periodSplits {
			row := candidateRows[split.Name]
			if row.CandidateCount < full.WeakestSplitRows {
				full.WeakestSplitRows = row.CandidateCount
			}
			if row.ToxicRate > full.WorstSplitToxicRate {
				full.WorstSplitToxicRate = row.ToxicRate
			}
			if row.ToxicRate < full.MinSplitToxicRate {
				full.MinSplitToxicRate = row.ToxicRate
			}
			improvement := row.ToxicRate - row.BaselineToxicRate
			if improvement < full.MinSplitToxicImprovement {
				full.MinSplitToxicImprovement = improvement
			}
			if full.CandidateCount > 0 {
				full.MaxSplitContributionRate = math.Max(full.MaxSplitContributionRate, safeDiv(float64(row.CandidateCount), float64(full.CandidateCount)))
			}
		}
		if full.WeakestSplitRows == math.MaxInt {
			full.WeakestSplitRows = 0
		}
		if math.IsInf(full.MinSplitToxicRate, 1) {
			full.MinSplitToxicRate = 0
		}
		if math.IsInf(full.MinSplitToxicImprovement, 1) {
			full.MinSplitToxicImprovement = 0
		}
		full.ExactMetricReproductionPass = !cfg.RequireExactMetricReproduction || derivativesNoTradeFilterPremiseExactMatches(full)
		full.ReviewableCountGatePass = full.FullRows >= cfg.MinFullCandidateRows && full.WeakestSplitRows >= cfg.MinSplitCandidateRows
		full.ToxicRateGatePass = full.FullToxicRate >= cfg.MinFullToxicRate && full.WorstSplitToxicRate >= cfg.MinWorstSplitToxicRate
		full.ImprovementGatePass = full.FullToxicImprovement >= cfg.MinFullToxicImprovement
		full.MissingnessGatePass = !strings.Contains(full.DerivativesBucketID, "missing")
		full.FutureLeakProtectionPass = true
		full.ClosedFamilyProtectionPass = true
		full.RotationEntryProtectionPass = true
		if !full.ExactMetricReproductionPass {
			reasons = append(reasons, "exact_context_row_metrics_not_reproduced")
		}
		if !full.ReviewableCountGatePass {
			reasons = append(reasons, "candidate_count_gate_failed")
		}
		if !full.ToxicRateGatePass {
			reasons = append(reasons, "toxic_rate_gate_failed")
		}
		if !full.ImprovementGatePass {
			reasons = append(reasons, "local_only_toxic_improvement_gate_failed")
		}
		if !full.MissingnessGatePass {
			reasons = append(reasons, "missing_context_bucket")
		}
		full.PassesGate = full.ExactMetricReproductionPass && full.ReviewableCountGatePass && full.ToxicRateGatePass && full.ImprovementGatePass && full.MissingnessGatePass && full.FutureLeakProtectionPass && full.ClosedFamilyProtectionPass && full.RotationEntryProtectionPass
		full.FailureReason = uniqueJoinedReasons(reasons)
		rows = append(rows, full)
		for _, split := range periodSplits {
			row := candidateRows[split.Name]
			row.FullRows = full.FullRows
			row.WeakestSplitRows = full.WeakestSplitRows
			row.MaxSplitContributionRate = full.MaxSplitContributionRate
			row.FullToxicRate = full.FullToxicRate
			row.WorstSplitToxicRate = full.WorstSplitToxicRate
			row.MinSplitToxicRate = full.MinSplitToxicRate
			row.FullToxicImprovement = full.FullToxicImprovement
			row.MinSplitToxicImprovement = full.MinSplitToxicImprovement
			row.ExactMetricReproductionPass = full.ExactMetricReproductionPass
			row.ReviewableCountGatePass = full.ReviewableCountGatePass
			row.ToxicRateGatePass = full.ToxicRateGatePass
			row.ImprovementGatePass = full.ImprovementGatePass
			row.MissingnessGatePass = full.MissingnessGatePass
			row.FutureLeakProtectionPass = full.FutureLeakProtectionPass
			row.ClosedFamilyProtectionPass = full.ClosedFamilyProtectionPass
			row.RotationEntryProtectionPass = full.RotationEntryProtectionPass
			row.PassesGate = full.PassesGate
			row.FailureReason = full.FailureReason
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CandidateID != rows[j].CandidateID {
			return rows[i].CandidateID < rows[j].CandidateID
		}
		return splitSortKey(rows[i].Split) < splitSortKey(rows[j].Split)
	})
	return rows
}

func derivativesNoTradeFilterPremiseExactMatches(row FuturesDerivativesNoTradeFilterPremiseExactCandidateRow) bool {
	return row.FullRows == row.ExpectedRows &&
		row.WeakestSplitRows == row.ExpectedWeakestSplitRows &&
		math.Abs(row.FullToxicRate-row.ExpectedFullToxicRate) <= 0.000001 &&
		math.Abs(row.WorstSplitToxicRate-row.ExpectedWorstSplitToxicRate) <= 0.000001 &&
		math.Abs(row.FullToxicImprovement-row.ExpectedFullToxicImprovement) <= 0.000001
}

func derivativesNoTradeFilterPremiseVetoRows(items map[string]derivativesNoTradeFilterPremiseLabelWithFeature, matches map[string][]string) []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	rows := make([]FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow, 0, len(keys))
	for _, key := range keys {
		item := items[key]
		label := item.label
		feature := item.feature
		matchIDs := append([]string(nil), matches[key]...)
		sort.Strings(matchIDs)
		rows = append(rows, FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow{
			VetoRowID: len(rows) + 1, Symbol: label.Symbol, Split: label.Split, Timeframe: label.Timeframe, HorizonBars: label.HorizonBars,
			LocalStateRowID: label.LocalStateRowID, FeatureRowID: label.FeatureRowID, DecisionCloseTime: label.Timestamp, SourceCloseTime: feature.SourceCloseTime,
			LocalRangeBucketID: label.LocalRangeBucketID, BasisLevelBucket: label.BasisLevelBucket, BasisChangeBucket: label.BasisChangeBucket,
			PremiumLevelBucket: label.PremiumLevelBucket, BasisBPS: feature.BasisBPS, PremiumBPS: feature.PremiumBPS,
			MatchedCandidateIDs: strings.Join(matchIDs, ";"), CandidateMatchCount: len(matchIDs),
			NestedTrendDownPremiumOverlap: derivativesNoTradeFilterPremiseHasNestedTrendDownOverlap(matchIDs),
			ForwardLabel:                  label.ForwardLabel, NoTradeToxic: label.NoTradeToxic, RotationUsefulBlocked: label.RotationUseful,
			ContinuationUsefulBlocked: label.ContinuationUseful, DiagnosticOnly: label.DiagnosticOnly, ForwardLabelMetadataOnly: true,
			ForwardLabelUsedAsFeature: label.ForwardLabelUsedAsFeature, ClosedCandleOnly: feature.ClosedCandleOnly, UsesFutureRows: feature.UsesFutureRows,
		})
	}
	return rows
}

func derivativesNoTradeFilterPremiseOverlapRows(vetoRows []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow, exactRows []FuturesDerivativesNoTradeFilterPremiseExactCandidateRow, defs []FuturesDerivativesNoTradeFilterPremiseDefinitionRow, splits []Split) []FuturesDerivativesNoTradeFilterPremiseOverlapRow {
	candidateCounts := map[string]map[string]int{}
	for _, row := range exactRows {
		if candidateCounts[row.CandidateID] == nil {
			candidateCounts[row.CandidateID] = map[string]int{}
		}
		candidateCounts[row.CandidateID][row.Split] = row.CandidateCount
	}
	defIDs := make([]string, 0, len(defs))
	for _, def := range defs {
		defIDs = append(defIDs, def.CandidateID)
	}
	splitNames := append([]string{fullSplitName}, derivativesNoTradeFilterPremiseSplitNames(rangeDiscoveryPeriodSplits(splits))...)
	overlap := map[string]map[string]int{}
	for _, veto := range vetoRows {
		for _, split := range rangeDiscoverySplitCombos(veto.Split) {
			ids := strings.Split(veto.MatchedCandidateIDs, ";")
			for i := 0; i < len(ids); i++ {
				for j := i + 1; j < len(ids); j++ {
					key := split + "|" + ids[i] + "|" + ids[j]
					if overlap[key] == nil {
						overlap[key] = map[string]int{}
					}
					overlap[key]["count"]++
				}
			}
		}
	}
	rows := []FuturesDerivativesNoTradeFilterPremiseOverlapRow{}
	for _, split := range splitNames {
		for i := 0; i < len(defIDs); i++ {
			for j := i + 1; j < len(defIDs); j++ {
				count := overlap[split+"|"+defIDs[i]+"|"+defIDs[j]]["count"]
				nested := (defIDs[i] == "exact_1_trend_down_basis_discount_premium_discount" && defIDs[j] == "exact_2_trend_down_basis_discount") ||
					(defIDs[j] == "exact_1_trend_down_basis_discount_premium_discount" && defIDs[i] == "exact_2_trend_down_basis_discount")
				rows = append(rows, FuturesDerivativesNoTradeFilterPremiseOverlapRow{
					Split: split, CandidateID: defIDs[i], OtherCandidateID: defIDs[j], OverlapRows: count,
					CandidateRows: candidateCounts[defIDs[i]][split], OtherCandidateRows: candidateCounts[defIDs[j]][split],
					NestedOverlap: nested && count > 0, DoubleCountedInUnion: false,
					OverlapPolicy: "deduplicated_by_symbol_timeframe_local_state_and_horizon",
				})
			}
		}
	}
	return rows
}

func derivativesNoTradeFilterPremiseUnionRows(vetoRows []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow, allLabels []derivativesNoTradeFilterPremiseLabelWithFeature, exactRows []FuturesDerivativesNoTradeFilterPremiseExactCandidateRow, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) []FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow {
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	accs := map[string]*derivativesNoTradeFilterPremiseUnionAccumulator{}
	add := func(split string, veto FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow) {
		acc := accs[split]
		if acc == nil {
			acc = &derivativesNoTradeFilterPremiseUnionAccumulator{labels: map[string]int{}}
			acc.row = FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{
				FilterID: derivativesNoTradeFilterID, Split: split, Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars,
				MatchedExactCandidates:   derivativesNoTradeFilterPremiseDefinitionIDList(),
				FutureLeakProtectionPass: true, ClosedFamilyProtectionPass: true, RotationEntryProtectionPass: true, MissingnessGatePass: true,
			}
			accs[split] = acc
		}
		acc.row.DeduplicatedRows++
		acc.row.OverlapRows += maxInt(0, veto.CandidateMatchCount-1)
		if veto.NestedTrendDownPremiumOverlap {
			acc.row.NestedTrendDownPremiumOverlapRows++
		}
		if veto.NoTradeToxic {
			acc.row.NoTradeToxicCount++
		}
		if veto.RotationUsefulBlocked {
			acc.row.RotationUsefulBlockedCount++
		}
		if veto.ContinuationUsefulBlocked {
			acc.row.ContinuationUsefulBlockedCount++
		}
		if veto.DiagnosticOnly {
			acc.row.DiagnosticOnlyCount++
		}
		acc.labels[veto.ForwardLabel]++
	}
	for _, veto := range vetoRows {
		for _, split := range rangeDiscoverySplitCombos(veto.Split) {
			add(split, veto)
		}
	}
	exactBySplit := map[string]int{}
	for _, row := range exactRows {
		exactBySplit[row.Split] += row.CandidateCount
	}
	baselineBySplit := derivativesNoTradeFilterPremiseUnionBaseline(allLabels, cfg, splits)
	rows := []FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{}
	splitNames := append([]string{fullSplitName}, derivativesNoTradeFilterPremiseSplitNames(periodSplits)...)
	for _, split := range splitNames {
		acc := accs[split]
		row := FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{
			FilterID: derivativesNoTradeFilterID, Split: split, Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars,
			MatchedExactCandidates:   derivativesNoTradeFilterPremiseDefinitionIDList(),
			FutureLeakProtectionPass: true, ClosedFamilyProtectionPass: true, RotationEntryProtectionPass: true, MissingnessGatePass: true,
		}
		if acc != nil {
			row = acc.row
		}
		row.SumExactCandidateRows = exactBySplit[split]
		row.ToxicRate = safeDiv(float64(row.NoTradeToxicCount), float64(row.DeduplicatedRows))
		base := baselineBySplit[split]
		row.LocalOnlyBaselineRows = base.CandidateCount
		row.LocalOnlyBaselineToxicRate = safeDiv(float64(base.NoTradeToxicCount), float64(base.CandidateCount))
		rows = append(rows, row)
	}
	bySplit := map[string]*FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{}
	for i := range rows {
		bySplit[rows[i].Split] = &rows[i]
	}
	full := bySplit[fullSplitName]
	if full != nil {
		full.FullRows = full.DeduplicatedRows
		full.FullToxicRate = full.ToxicRate
		full.FullToxicImprovement = full.ToxicRate - full.LocalOnlyBaselineToxicRate
		full.WeakestSplitRows = math.MaxInt
		full.MinSplitToxicRate = math.Inf(1)
		reasons := []string{}
		for _, split := range periodSplits {
			row := bySplit[split.Name]
			if row == nil {
				full.WeakestSplitRows = 0
				full.MinSplitToxicRate = 0
				continue
			}
			if row.DeduplicatedRows < full.WeakestSplitRows {
				full.WeakestSplitRows = row.DeduplicatedRows
			}
			if row.ToxicRate > full.WorstSplitToxicRate {
				full.WorstSplitToxicRate = row.ToxicRate
			}
			if row.ToxicRate < full.MinSplitToxicRate {
				full.MinSplitToxicRate = row.ToxicRate
			}
			if full.FullRows > 0 {
				full.MaxSplitContributionRate = math.Max(full.MaxSplitContributionRate, safeDiv(float64(row.DeduplicatedRows), float64(full.FullRows)))
			}
		}
		if full.WeakestSplitRows == math.MaxInt {
			full.WeakestSplitRows = 0
		}
		if math.IsInf(full.MinSplitToxicRate, 1) {
			full.MinSplitToxicRate = 0
		}
		full.ToxicDominatedEverySplit = full.FullToxicRate >= cfg.MinUnionToxicRateFull && full.MinSplitToxicRate >= cfg.MinUnionToxicRateEverySplit
		full.DoubleCountingProtectionPass = full.SumExactCandidateRows >= full.DeduplicatedRows && full.OverlapRows == full.SumExactCandidateRows-full.DeduplicatedRows
		if !full.ToxicDominatedEverySplit {
			reasons = append(reasons, "canonical_union_not_toxic_dominated_every_split")
		}
		if !full.DoubleCountingProtectionPass {
			reasons = append(reasons, "canonical_union_overlap_accounting_failed")
		}
		full.PassesGate = full.ToxicDominatedEverySplit && full.DoubleCountingProtectionPass && full.FutureLeakProtectionPass && full.ClosedFamilyProtectionPass && full.RotationEntryProtectionPass && full.MissingnessGatePass
		full.FailureReason = uniqueJoinedReasons(reasons)
		for i := range rows {
			rows[i].FullRows = full.FullRows
			rows[i].WeakestSplitRows = full.WeakestSplitRows
			rows[i].MaxSplitContributionRate = full.MaxSplitContributionRate
			rows[i].FullToxicRate = full.FullToxicRate
			rows[i].WorstSplitToxicRate = full.WorstSplitToxicRate
			rows[i].MinSplitToxicRate = full.MinSplitToxicRate
			rows[i].FullToxicImprovement = full.FullToxicImprovement
			rows[i].ToxicDominatedEverySplit = full.ToxicDominatedEverySplit
			rows[i].DoubleCountingProtectionPass = full.DoubleCountingProtectionPass
			rows[i].MissingnessGatePass = full.MissingnessGatePass
			rows[i].FutureLeakProtectionPass = full.FutureLeakProtectionPass
			rows[i].ClosedFamilyProtectionPass = full.ClosedFamilyProtectionPass
			rows[i].RotationEntryProtectionPass = full.RotationEntryProtectionPass
			rows[i].PassesGate = full.PassesGate
			rows[i].FailureReason = full.FailureReason
		}
	}
	sort.Slice(rows, func(i, j int) bool { return splitSortKey(rows[i].Split) < splitSortKey(rows[j].Split) })
	return rows
}

type derivativesNoTradeFilterPremiseBaselineCounts struct {
	CandidateCount    int
	NoTradeToxicCount int
}

func derivativesNoTradeFilterPremiseUnionBaseline(labels []derivativesNoTradeFilterPremiseLabelWithFeature, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) map[string]derivativesNoTradeFilterPremiseBaselineCounts {
	selectedLocal := map[string]bool{}
	for _, def := range derivativesNoTradeFilterPremiseDefinitions(cfg) {
		selectedLocal[def.LocalRangeBucketID] = true
	}
	out := map[string]derivativesNoTradeFilterPremiseBaselineCounts{}
	for _, item := range labels {
		label := item.label
		if !selectedLocal[label.LocalRangeBucketID] {
			continue
		}
		for _, split := range rangeDiscoverySplitCombos(label.Split) {
			row := out[split]
			row.CandidateCount++
			if label.NoTradeToxic {
				row.NoTradeToxicCount++
			}
			out[split] = row
		}
	}
	return out
}

func derivativesNoTradeFilterPremiseCollateralRows(vetoRows []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow, splits []Split) []FuturesDerivativesNoTradeFilterPremiseCollateralRow {
	type acc struct {
		row FuturesDerivativesNoTradeFilterPremiseCollateralRow
	}
	accs := map[string]*acc{}
	for _, veto := range vetoRows {
		for _, split := range rangeDiscoverySplitCombos(veto.Split) {
			key := split + "|" + veto.ForwardLabel
			a := accs[key]
			if a == nil {
				a = &acc{row: FuturesDerivativesNoTradeFilterPremiseCollateralRow{Split: split, ForwardLabel: veto.ForwardLabel, CollateralDamageReported: true}}
				accs[key] = a
			}
			a.row.BlockedRows++
			if veto.NoTradeToxic {
				a.row.NoTradeToxicRows++
			}
			if veto.RotationUsefulBlocked {
				a.row.RotationUsefulBlockedRows++
			}
			if veto.ContinuationUsefulBlocked {
				a.row.ContinuationUsefulBlockedRows++
			}
			if veto.DiagnosticOnly {
				a.row.DiagnosticOnlyRows++
			}
		}
	}
	totalBySplit := map[string]int{}
	for _, a := range accs {
		totalBySplit[a.row.Split] += a.row.BlockedRows
	}
	rows := make([]FuturesDerivativesNoTradeFilterPremiseCollateralRow, 0, len(accs))
	for _, a := range accs {
		a.row.ShareOfBlockedRows = safeDiv(float64(a.row.BlockedRows), float64(totalBySplit[a.row.Split]))
		rows = append(rows, a.row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if splitSortKey(rows[i].Split) != splitSortKey(rows[j].Split) {
			return splitSortKey(rows[i].Split) < splitSortKey(rows[j].Split)
		}
		if rows[i].BlockedRows != rows[j].BlockedRows {
			return rows[i].BlockedRows > rows[j].BlockedRows
		}
		return rows[i].ForwardLabel < rows[j].ForwardLabel
	})
	return rows
}

func derivativesNoTradeFilterPremiseSourceRows(sourceAudit FuturesDerivativesContextSourceAuditResult) []FuturesDerivativesNoTradeFilterPremiseSourceRow {
	rows := []FuturesDerivativesNoTradeFilterPremiseSourceRow{}
	for _, row := range sourceAudit.SourceRows {
		rows = append(rows, FuturesDerivativesNoTradeFilterPremiseSourceRow{
			AuditName:        FuturesDerivativesNoTradeFilterPremiseAuditName,
			Role:             "derivatives_context_source",
			Symbol:           row.Symbol,
			SourceFamily:     row.SourceFamily,
			ProductScope:     row.ProductScope,
			NativeInterval:   row.NativeInterval,
			DurablePath:      row.DurablePath,
			Required:         row.Required,
			RowCount:         row.RowCount,
			ExpectedRowCount: row.ExpectedRowCount,
			FirstOpenTime:    row.FirstOpenTime,
			LastOpenTime:     row.LastOpenTime,
			GapCount:         row.GapCount,
			DuplicateCount:   row.DuplicateCount,
			FileSHA256:       row.FileSHA256,
			FinalityRule:     row.FinalityRule,
			ValidationStatus: row.ValidationStatus,
			ValidationError:  row.ValidationError,
		})
	}
	return rows
}

func derivativesNoTradeFilterPremiseAnchorSourceRow(source FuturesBTCRegimeETHSOLContextSourceRow) FuturesDerivativesNoTradeFilterPremiseSourceRow {
	return FuturesDerivativesNoTradeFilterPremiseSourceRow{
		AuditName:                  FuturesDerivativesNoTradeFilterPremiseAuditName,
		Role:                       "candle_anchor_local_range_state",
		Symbol:                     source.Symbol,
		SourceFamily:               "trade_klines",
		ProductScope:               source.Product,
		NativeInterval:             source.Interval,
		DurablePath:                source.Path,
		Required:                   true,
		RowCount:                   source.RowCount,
		ExpectedRowCount:           source.ExpectedRowCount,
		FirstOpenTime:              source.FirstOpenTime,
		LastOpenTime:               source.LastOpenTime,
		GapCount:                   source.GapCount,
		DuplicateCount:             source.DuplicateCount,
		ZeroVolumeCount:            source.ZeroVolumeCount,
		FinalityRule:               "confirmed closed candle local range-state anchor",
		ValidationStatus:           source.ValidationStatus,
		ValidationError:            source.ValidationError,
		ForwardLabelsAsSourceInput: false,
	}
}

func derivativesNoTradeFilterPremiseSourceCoverageRows(sourceAudit FuturesDerivativesContextSourceAuditResult, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) []FuturesDerivativesNoTradeFilterPremiseCoverageRow {
	rows := []FuturesDerivativesNoTradeFilterPremiseCoverageRow{}
	for _, row := range sourceAudit.TimestampAlignmentRows {
		rows = append(rows, FuturesDerivativesNoTradeFilterPremiseCoverageRow{
			AuditName:             FuturesDerivativesNoTradeFilterPremiseAuditName,
			Scope:                 "source_lag_alignment",
			Symbol:                row.Symbol,
			SourceFamily:          row.SourceFamily,
			Required:              row.Required,
			RowCount:              row.AnchorDecisionCandles,
			LagCoveragePct:        row.LagCoveragePct,
			RequiredCoverageFloor: cfg.MinBasisContextCoverage,
			CoverageGatePass:      row.LagCoveragePct+1e-6 >= cfg.MinBasisContextCoverage && !row.UsesFutureRows && row.ExactClosedIntervalJoin,
			ClosedCandleOnly:      true,
			ForwardLabelsAsInputs: false,
			NoFillPolicy:          true,
			ValidationStatus:      "accepted",
		})
	}
	return rows
}

func derivativesNoTradeFilterPremiseLocalCoverageRow(base FuturesBTCRegimeETHSOLContextCoverageRow, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) FuturesDerivativesNoTradeFilterPremiseCoverageRow {
	return FuturesDerivativesNoTradeFilterPremiseCoverageRow{
		AuditName:             FuturesDerivativesNoTradeFilterPremiseAuditName,
		Scope:                 "local_range_state_resample",
		Symbol:                RangeUniverseSymbolBTCUSDT,
		Timeframe:             base.Timeframe,
		SourceFamily:          "trade_klines",
		Required:              true,
		RowCount:              base.RowCount,
		RequiredCoverageFloor: cfg.MinBasisContextCoverage,
		CoverageGatePass:      base.CoverageFactsPass && base.ValidationStatus == "accepted" && base.Complete,
		ClosedCandleOnly:      true,
		ForwardLabelsAsInputs: false,
		NoFillPolicy:          true,
		ValidationStatus:      base.ValidationStatus,
		ValidationError:       base.ValidationError,
	}
}

func derivativesNoTradeFilterPremiseSourceMissingnessRows(sourceAudit FuturesDerivativesContextSourceAuditResult, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) []FuturesDerivativesNoTradeFilterPremiseMissingnessRow {
	rows := []FuturesDerivativesNoTradeFilterPremiseMissingnessRow{}
	for _, row := range sourceAudit.MissingnessRows {
		rows = append(rows, FuturesDerivativesNoTradeFilterPremiseMissingnessRow{
			AuditName:         FuturesDerivativesNoTradeFilterPremiseAuditName,
			Scope:             "source_stream",
			Symbol:            row.Symbol,
			SourceFamily:      row.SourceFamily,
			Reason:            "source_audit_recorded_missingness",
			Count:             row.LagMissingContextRows,
			TotalRows:         row.LagMissingContextRows + sourceAuditLagAlignedRows(sourceAudit, row.Symbol, row.SourceFamily),
			Rate:              safeDiv(float64(row.LagMissingContextRows), float64(row.LagMissingContextRows+sourceAuditLagAlignedRows(sourceAudit, row.Symbol, row.SourceFamily))),
			ForwardFilledRows: row.ForwardFilledRows,
			MissingDataPolicy: row.MissingDataPolicy,
			CoverageFloor:     cfg.MinBasisContextCoverage,
			CoverageGatePass:  true,
		})
	}
	return rows
}

func derivativesNoTradeFilterPremiseSourcePass(sourceAudit FuturesDerivativesContextSourceAuditResult, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig) bool {
	if sourceAudit.StopState != DerivativesContextSourceAuditStopStatePassedNeedsContextBrief {
		return false
	}
	required := 0
	accepted := 0
	for _, row := range sourceAudit.SourceRows {
		if !row.Required {
			continue
		}
		required++
		if row.ValidationStatus != "rejected" {
			accepted++
		}
	}
	if required != 3 || accepted != 3 {
		return false
	}
	for _, row := range sourceAudit.TimestampAlignmentRows {
		if row.Required && (!row.MeetsMinCoverage || row.LagCoveragePct+1e-6 < cfg.MinBasisContextCoverage) {
			return false
		}
		if row.UsesFutureRows || !row.ExactClosedIntervalJoin {
			return false
		}
	}
	return true
}

func derivativesNoTradeFilterPremiseCoveragePass(rows []FuturesDerivativesNoTradeFilterPremiseCoverageRow) bool {
	if len(rows) == 0 {
		return false
	}
	hasSource, hasLocal, hasContext := false, false, false
	for _, row := range rows {
		if row.ForwardLabelsAsInputs || !row.CoverageGatePass || !row.NoFillPolicy {
			return false
		}
		switch row.Scope {
		case "source_lag_alignment":
			hasSource = true
		case "local_range_state_resample":
			hasLocal = true
		case "lagged_basis_context":
			hasContext = true
		}
	}
	return hasSource && hasLocal && hasContext
}

func derivativesNoTradeFilterPremiseMinLagCoverage(sourceAudit FuturesDerivativesContextSourceAuditResult) float64 {
	minCov := math.Inf(1)
	for _, row := range sourceAudit.TimestampAlignmentRows {
		if row.Required && row.LagCoveragePct < minCov {
			minCov = row.LagCoveragePct
		}
	}
	if math.IsInf(minCov, 1) {
		return 0
	}
	return minCov
}

func derivativesNoTradeFilterPremiseSummaryRows(result FuturesDerivativesNoTradeFilterPremiseAuditResult, cfg FuturesDerivativesNoTradeFilterPremiseAuditConfig, splits []Split) []FuturesDerivativesNoTradeFilterPremiseSummaryRow {
	coveragePass := derivativesNoTradeFilterPremiseCoveragePass(result.CoverageRows)
	minLag := math.Inf(1)
	for _, row := range result.CoverageRows {
		if row.Scope == "source_lag_alignment" && row.Required && row.LagCoveragePct < minLag {
			minLag = row.LagCoveragePct
		}
	}
	if math.IsInf(minLag, 1) {
		minLag = 0
	}
	canonBySplit := map[string]FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{}
	for _, row := range result.CanonicalUnionRows {
		canonBySplit[row.Split] = row
	}
	rows := []FuturesDerivativesNoTradeFilterPremiseSummaryRow{}
	for _, split := range append([]Split{{Name: fullSplitName}}, rangeDiscoveryPeriodSplits(splits)...) {
		canon := canonBySplit[split.Name]
		rows = append(rows, FuturesDerivativesNoTradeFilterPremiseSummaryRow{
			Split: split.Name, FilterID: derivativesNoTradeFilterID, Symbol: RangeUniverseSymbolBTCUSDT, Timeframe: cfg.Timeframe, HorizonBars: cfg.HorizonBars,
			SourceRows: len(result.SourceRows), CoverageRows: len(result.CoverageRows), FilterDefinitionRows: len(result.FilterDefinitionRows),
			ExactCandidateRows: len(result.ExactCandidateRows), ExactCandidatesPassed: result.ExactCandidatesPassed,
			VetoCandidateRows: len(result.VetoCandidateRows), CanonicalDeduplicatedRows: canon.DeduplicatedRows,
			CanonicalOverlapRows: canon.OverlapRows, NestedTrendDownPremiumOverlapRows: canon.NestedTrendDownPremiumOverlapRows,
			CanonicalNoTradeToxicRows: canon.NoTradeToxicCount, CanonicalToxicRate: canon.ToxicRate,
			CanonicalMinSplitToxicRate: canon.MinSplitToxicRate, CanonicalRotationUsefulBlocked: canon.RotationUsefulBlockedCount,
			CanonicalContinuationUsefulBlocked: canon.ContinuationUsefulBlockedCount, CollateralRows: len(result.CollateralRows),
			MissingnessRows: len(result.MissingnessRows), RequiredCoverageFloor: cfg.MinBasisContextCoverage, MinLagCoveragePct: minLag,
			SourceScopePass: coveragePass, CoveragePass: coveragePass,
			ExactRowsReproduced: result.ExactCandidatesPassed == len(result.FilterDefinitionRows), CanonicalUnionPass: result.CanonicalUnionPassed,
			CommonOutputsZeroTrade: true, ForwardLabelsAsInputs: false, ClosedFamilyProtectionPass: true, RotationEntryProtectionPass: true,
			Trades: 0, StopState: result.StopState,
		})
	}
	return rows
}

func derivativesNoTradeFilterPremiseSplitNames(splits []Split) []string {
	out := make([]string, 0, len(splits))
	for _, split := range splits {
		out = append(out, split.Name)
	}
	return out
}

func derivativesNoTradeFilterPremiseDefinitionIDList() string {
	defs := derivativesNoTradeFilterPremiseDefinitions(DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig())
	ids := make([]string, 0, len(defs))
	for _, def := range defs {
		ids = append(ids, def.CandidateID)
	}
	sort.Strings(ids)
	return strings.Join(ids, ";")
}

func derivativesNoTradeFilterPremiseHasNestedTrendDownOverlap(ids []string) bool {
	seen := map[string]bool{}
	for _, id := range ids {
		seen[id] = true
	}
	return seen["exact_1_trend_down_basis_discount_premium_discount"] && seen["exact_2_trend_down_basis_discount"]
}
