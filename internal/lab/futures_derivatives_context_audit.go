package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	FuturesDerivativesContextAuditName = "futures_derivatives_context_audit"

	DerivativesContextAuditStopStateSourceGap                  = "derivatives_context_zero_trade_context_audit_source_gap"
	DerivativesContextAuditStopStateRejectedFutureLabelLeak    = "derivatives_context_zero_trade_context_audit_rejected_future_label_leak"
	DerivativesContextAuditStopStateRejectedClosedFamilyRescue = "derivatives_context_zero_trade_context_audit_rejected_closed_family_rescue"
	DerivativesContextAuditStopStateFailedNoUsableContext      = "derivatives_context_zero_trade_context_audit_failed_no_usable_context"
	DerivativesContextAuditStopStatePassedNeedsSpec            = "derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec"

	DerivativesContextComparisonLocalOnly      = "local_only_baseline"
	DerivativesContextComparisonLocalPlusDeriv = "derivatives_context_plus_local"

	DerivativesContextBucketBasisLevel        = "basis_level"
	DerivativesContextBucketBasisLevelChange  = "basis_level_plus_change"
	DerivativesContextBucketBasisLevelPremium = "basis_level_plus_premium"
)

type FuturesDerivativesContextAuditConfig struct {
	SourceAuditConfig                  FuturesDerivativesContextSourceAuditConfig
	LocalSources                       []FuturesBTCRegimeETHSOLContextSourceConfig
	StateConfig                        FuturesRangeStateConstructionLoopAuditConfig
	MinBasisContextCoverage            float64
	BasisChangeLookbackIntervals       int
	BasisVolatilityLookbackIntervals   int
	MinFullCohortRows                  int
	MinSplitCohortRows                 int
	MaxSplitContributionRate           float64
	MinUsefulRateFull                  float64
	MinUsefulRateSplit                 float64
	MaxToxicRateFull                   float64
	MaxToxicRateSplit                  float64
	MinUsefulMinusToxicMarginFull      float64
	MinUsefulMinusToxicMarginSplit     float64
	MinToxicRateFull                   float64
	MinToxicRateSplit                  float64
	MinContextUsefulImprovementFull    float64
	MinContextUsefulImprovementSplit   float64
	MinContextMarginImprovementFull    float64
	MinContextMarginImprovementSplit   float64
	MinContextToxicImprovementFull     float64
	MinContextToxicImprovementSplit    float64
	MaxBucketShareOfLocalBaselineFull  float64
	MaxBucketShareOfLocalBaselineSplit float64
	RescueClosedFamily                 bool
}

type FuturesDerivativesContextAuditResult struct {
	SourceRows       []FuturesDerivativesContextAuditSourceRow       `json:"source_rows"`
	CoverageRows     []FuturesDerivativesContextAuditCoverageRow     `json:"coverage_rows"`
	BasisFeatureRows []FuturesDerivativesContextAuditBasisFeatureRow `json:"basis_feature_rows"`
	LocalStateRows   []FuturesDerivativesContextAuditLocalStateRow   `json:"local_state_rows"`
	LabelRows        []FuturesDerivativesContextAuditLabelRow        `json:"label_rows"`
	CohortRows       []FuturesDerivativesContextAuditCohortRow       `json:"cohort_rows"`
	RankingRows      []FuturesDerivativesContextAuditRankingRow      `json:"ranking_rows"`
	MissingnessRows  []FuturesDerivativesContextAuditMissingnessRow  `json:"missingness_rows"`
	SummaryRows      []FuturesDerivativesContextAuditSummaryRow      `json:"summary_rows"`
	PassingCohorts   int                                             `json:"passing_cohorts"`
	StopState        string                                          `json:"stop_state"`
}

type FuturesDerivativesContextAuditSourceRow struct {
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

type FuturesDerivativesContextAuditCoverageRow struct {
	AuditName               string  `json:"audit_name"`
	Scope                   string  `json:"scope"`
	Symbol                  string  `json:"symbol"`
	Timeframe               string  `json:"timeframe,omitempty"`
	SourceFamily            string  `json:"source_family,omitempty"`
	Required                bool    `json:"required"`
	RowCount                int     `json:"row_count"`
	StateRows               int     `json:"state_rows"`
	LabelRows               int     `json:"label_rows"`
	BasisFeatureRows        int     `json:"basis_feature_rows"`
	MissingBasisRows        int     `json:"missing_basis_rows"`
	MissingBasisChangeRows  int     `json:"missing_basis_change_rows"`
	MissingPremiumRows      int     `json:"missing_premium_rows"`
	LagCoveragePct          float64 `json:"lag_coverage_pct"`
	BasisContextCoveragePct float64 `json:"basis_context_coverage_pct"`
	RequiredCoverageFloor   float64 `json:"required_coverage_floor"`
	CoverageGatePass        bool    `json:"coverage_gate_pass"`
	ClosedCandleOnly        bool    `json:"closed_candle_only"`
	ForwardLabelsAsInputs   bool    `json:"forward_labels_as_inputs"`
	ValidationStatus        string  `json:"validation_status"`
	ValidationError         string  `json:"validation_error,omitempty"`
}

type FuturesDerivativesContextAuditBasisFeatureRow struct {
	FeatureRowID                 int     `json:"feature_row_id"`
	Symbol                       string  `json:"symbol"`
	LocalStateRowID              int     `json:"local_state_row_id"`
	Timestamp                    string  `json:"timestamp"`
	Timeframe                    string  `json:"timeframe"`
	Split                        string  `json:"split"`
	DecisionIndex                int     `json:"decision_index"`
	DecisionCloseTime            string  `json:"decision_close_time"`
	SourceOpenTime               string  `json:"source_open_time"`
	SourceCloseTime              string  `json:"source_close_time"`
	SourceLagIntervals           int     `json:"source_lag_intervals"`
	AntiLookaheadRule            string  `json:"anti_lookahead_rule"`
	MarkClose                    float64 `json:"mark_close"`
	IndexClose                   float64 `json:"index_close"`
	BasisRaw                     float64 `json:"basis_raw"`
	BasisBPS                     float64 `json:"basis_bps"`
	BasisSignBucket              string  `json:"basis_sign_bucket"`
	BasisLevelBucket             string  `json:"basis_level_bucket"`
	BasisChangeLookbackIntervals int     `json:"basis_change_lookback_intervals"`
	BasisChangeBPS               float64 `json:"basis_change_bps"`
	BasisChangeBucket            string  `json:"basis_change_bucket"`
	BasisChangePresent           bool    `json:"basis_change_present"`
	BasisVolatilityBPS           float64 `json:"basis_volatility_bps"`
	BasisVolatilityBucket        string  `json:"basis_volatility_bucket"`
	PremiumClose                 float64 `json:"premium_close"`
	PremiumBPS                   float64 `json:"premium_bps"`
	PremiumLevelBucket           string  `json:"premium_level_bucket"`
	PremiumPresent               bool    `json:"premium_present"`
	LocalRangeBucketID           string  `json:"local_range_bucket_id"`
	ContextStateID               string  `json:"context_state_id"`
	ClosedCandleOnly             bool    `json:"closed_candle_only"`
	UsesFutureRows               bool    `json:"uses_future_rows"`
	ForwardLabelsAsFeatureInput  bool    `json:"forward_labels_as_feature_input"`
}

type FuturesDerivativesContextAuditLocalStateRow struct {
	LocalStateRowID             int    `json:"local_state_row_id"`
	Symbol                      string `json:"symbol"`
	Timestamp                   string `json:"timestamp"`
	Timeframe                   string `json:"timeframe"`
	Split                       string `json:"split"`
	DecisionIndex               int    `json:"decision_index"`
	DecisionCloseTime           string `json:"decision_close_time"`
	RangeEpisodeID              int    `json:"range_episode_id"`
	LocalRangeBucketID          string `json:"local_range_bucket_id"`
	StateID                     string `json:"state_id"`
	GeometryBucket              string `json:"geometry_bucket"`
	VolBucket                   string `json:"vol_bucket"`
	TrendBucket                 string `json:"trend_bucket"`
	ImpulseBucket               string `json:"impulse_bucket"`
	ParticipationBucket         string `json:"participation_bucket"`
	ContextStateID              string `json:"context_state_id,omitempty"`
	BasisLevelBucket            string `json:"basis_level_bucket,omitempty"`
	BasisChangeBucket           string `json:"basis_change_bucket,omitempty"`
	PremiumLevelBucket          string `json:"premium_level_bucket,omitempty"`
	DerivativesContextAvailable bool   `json:"derivatives_context_available"`
	MissingContextReason        string `json:"missing_context_reason,omitempty"`
	ClosedCandleOnly            bool   `json:"closed_candle_only"`
	AuthorityCandidateOnly      bool   `json:"authority_candidate_only"`
	ForwardLabelsAsStateInput   bool   `json:"forward_labels_as_state_input"`
	ClosedFamilyRescue          bool   `json:"closed_family_rescue"`
}

type FuturesDerivativesContextAuditLabelRow struct {
	LabelRowID                int    `json:"label_row_id"`
	Symbol                    string `json:"symbol"`
	LocalStateRowID           int    `json:"local_state_row_id"`
	FeatureRowID              int    `json:"feature_row_id"`
	Timestamp                 string `json:"timestamp"`
	Timeframe                 string `json:"timeframe"`
	Split                     string `json:"split"`
	HorizonBars               int    `json:"horizon_bars"`
	LocalRangeBucketID        string `json:"local_range_bucket_id"`
	BasisLevelBucket          string `json:"basis_level_bucket"`
	BasisChangeBucket         string `json:"basis_change_bucket"`
	BasisChangePresent        bool   `json:"basis_change_present"`
	PremiumLevelBucket        string `json:"premium_level_bucket"`
	PremiumPresent            bool   `json:"premium_present"`
	ContextStateID            string `json:"context_state_id"`
	ForwardLabel              string `json:"forward_label"`
	RotationUseful            bool   `json:"rotation_useful"`
	RotationToxic             bool   `json:"rotation_toxic"`
	ContinuationUseful        bool   `json:"continuation_useful"`
	ContinuationToxic         bool   `json:"continuation_toxic"`
	NoTradeToxic              bool   `json:"no_trade_toxic"`
	DiagnosticOnly            bool   `json:"diagnostic_only"`
	LabelWindowStartIndex     int    `json:"label_window_start_index"`
	LabelWindowEndIndex       int    `json:"label_window_end_index"`
	LabelWindowStartTime      string `json:"label_window_start_time"`
	LabelWindowEndTime        string `json:"label_window_end_time"`
	ForwardLabelMetadataOnly  bool   `json:"forward_label_metadata_only"`
	ForwardLabelUsedAsFeature bool   `json:"forward_label_used_as_feature"`
}

type FuturesDerivativesContextAuditCohortRow struct {
	CohortID                           string  `json:"cohort_id"`
	ComparisonType                     string  `json:"comparison_type"`
	Symbol                             string  `json:"symbol"`
	Split                              string  `json:"split"`
	Timeframe                          string  `json:"timeframe"`
	HorizonBars                        int     `json:"horizon_bars"`
	RouteCandidate                     string  `json:"route_candidate"`
	LocalRangeBucketID                 string  `json:"local_range_bucket_id"`
	DerivativesBucketType              string  `json:"derivatives_bucket_type,omitempty"`
	DerivativesBucketID                string  `json:"derivatives_bucket_id,omitempty"`
	BasisLevelBucket                   string  `json:"basis_level_bucket,omitempty"`
	BasisChangeBucket                  string  `json:"basis_change_bucket,omitempty"`
	PremiumLevelBucket                 string  `json:"premium_level_bucket,omitempty"`
	CandidateCount                     int     `json:"candidate_count"`
	UsefulCount                        int     `json:"useful_count"`
	ToxicCount                         int     `json:"toxic_count"`
	RotationUsefulCount                int     `json:"rotation_useful_count"`
	RotationToxicCount                 int     `json:"rotation_toxic_count"`
	ContinuationUsefulCount            int     `json:"continuation_useful_count"`
	ContinuationToxicCount             int     `json:"continuation_toxic_count"`
	NoTradeToxicCount                  int     `json:"no_trade_toxic_count"`
	UsefulRate                         float64 `json:"useful_rate"`
	ToxicRate                          float64 `json:"toxic_rate"`
	UsefulMinusToxicMargin             float64 `json:"useful_minus_toxic_margin"`
	DominantForwardLabel               string  `json:"dominant_forward_label"`
	DominantForwardLabelRate           float64 `json:"dominant_forward_label_rate"`
	FullPeriodRows                     int     `json:"full_period_rows"`
	WeakestSplitRows                   int     `json:"weakest_split_rows"`
	MaxSplitContributionRate           float64 `json:"max_split_contribution_rate"`
	FullUsefulRate                     float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate             float64 `json:"weakest_split_useful_rate"`
	FullToxicRate                      float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate                float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin         float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin                 float64 `json:"weakest_split_margin"`
	BaselineFullRows                   int     `json:"baseline_full_rows"`
	BaselineFullUsefulRate             float64 `json:"baseline_full_useful_rate"`
	BaselineWeakestSplitUsefulRate     float64 `json:"baseline_weakest_split_useful_rate"`
	BaselineFullToxicRate              float64 `json:"baseline_full_toxic_rate"`
	BaselineWorstSplitToxicRate        float64 `json:"baseline_worst_split_toxic_rate"`
	BaselineFullMargin                 float64 `json:"baseline_full_margin"`
	BaselineWeakestSplitMargin         float64 `json:"baseline_weakest_split_margin"`
	FullUsefulImprovement              float64 `json:"full_useful_improvement"`
	WeakestSplitUsefulImprovement      float64 `json:"weakest_split_useful_improvement"`
	FullToxicImprovement               float64 `json:"full_toxic_improvement"`
	WeakestSplitToxicImprovement       float64 `json:"weakest_split_toxic_improvement"`
	FullMarginImprovement              float64 `json:"full_margin_improvement"`
	WeakestSplitMarginImprovement      float64 `json:"weakest_split_margin_improvement"`
	BucketShareOfLocalBaselineFull     float64 `json:"bucket_share_of_local_baseline_full"`
	MaxBucketShareOfLocalBaselineSplit float64 `json:"max_bucket_share_of_local_baseline_split"`
	ReviewableCountGatePass            bool    `json:"reviewable_count_gate_pass"`
	SplitStabilityGatePass             bool    `json:"split_stability_gate_pass"`
	SplitContributionGatePass          bool    `json:"split_contribution_gate_pass"`
	RouteRateGatePass                  bool    `json:"route_rate_gate_pass"`
	DerivativesImprovementGatePass     bool    `json:"derivatives_improvement_gate_pass"`
	OrthogonalityGatePass              bool    `json:"orthogonality_gate_pass"`
	MissingnessGatePass                bool    `json:"missingness_gate_pass"`
	PremiumOnlyGatePass                bool    `json:"premium_only_gate_pass"`
	ClosedFamilyProtectionPass         bool    `json:"closed_family_protection_pass"`
	FutureLeakProtectionPass           bool    `json:"future_leak_protection_pass"`
	PassesReviewGate                   bool    `json:"passes_review_gate"`
	FailureReason                      string  `json:"failure_reason,omitempty"`
}

type FuturesDerivativesContextAuditRankingRow struct {
	Rank                               int     `json:"rank"`
	CohortID                           string  `json:"cohort_id"`
	Symbol                             string  `json:"symbol"`
	Timeframe                          string  `json:"timeframe"`
	HorizonBars                        int     `json:"horizon_bars"`
	RouteCandidate                     string  `json:"route_candidate"`
	LocalRangeBucketID                 string  `json:"local_range_bucket_id"`
	DerivativesBucketType              string  `json:"derivatives_bucket_type"`
	DerivativesBucketID                string  `json:"derivatives_bucket_id"`
	PassesGate                         bool    `json:"passes_gate"`
	RankScore                          float64 `json:"rank_score"`
	FullPeriodRows                     int     `json:"full_period_rows"`
	WeakestSplitRows                   int     `json:"weakest_split_rows"`
	MaxSplitContributionRate           float64 `json:"max_split_contribution_rate"`
	FullUsefulRate                     float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate             float64 `json:"weakest_split_useful_rate"`
	FullToxicRate                      float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate                float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin         float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin                 float64 `json:"weakest_split_margin"`
	FullUsefulImprovement              float64 `json:"full_useful_improvement"`
	WeakestSplitUsefulImprovement      float64 `json:"weakest_split_useful_improvement"`
	FullToxicImprovement               float64 `json:"full_toxic_improvement"`
	WeakestSplitToxicImprovement       float64 `json:"weakest_split_toxic_improvement"`
	FullMarginImprovement              float64 `json:"full_margin_improvement"`
	WeakestSplitMarginImprovement      float64 `json:"weakest_split_margin_improvement"`
	BucketShareOfLocalBaselineFull     float64 `json:"bucket_share_of_local_baseline_full"`
	MaxBucketShareOfLocalBaselineSplit float64 `json:"max_bucket_share_of_local_baseline_split"`
	OrthogonalityGatePass              bool    `json:"orthogonality_gate_pass"`
	DerivativesImprovementGatePass     bool    `json:"derivatives_improvement_gate_pass"`
	DominantForwardLabel               string  `json:"dominant_forward_label"`
	DominantForwardLabelRate           float64 `json:"dominant_forward_label_rate"`
	FutureLeakProtectionPass           bool    `json:"future_leak_protection_pass"`
	ClosedFamilyProtectionPass         bool    `json:"closed_family_protection_pass"`
	FailureReason                      string  `json:"failure_reason,omitempty"`
}

type FuturesDerivativesContextAuditMissingnessRow struct {
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

type FuturesDerivativesContextAuditSummaryRow struct {
	Split                    string  `json:"split"`
	Symbol                   string  `json:"symbol"`
	Timeframe                string  `json:"timeframe"`
	HorizonBars              int     `json:"horizon_bars"`
	SourceRows               int     `json:"source_rows"`
	CoverageRows             int     `json:"coverage_rows"`
	BasisFeatureRows         int     `json:"basis_feature_rows"`
	LocalStateRows           int     `json:"local_state_rows"`
	LabelRows                int     `json:"label_rows"`
	CohortRows               int     `json:"cohort_rows"`
	RankingRows              int     `json:"ranking_rows"`
	PassingCohorts           int     `json:"passing_cohorts"`
	MissingnessRows          int     `json:"missingness_rows"`
	RequiredCoverageFloor    float64 `json:"required_coverage_floor"`
	MinBasisCoveragePct      float64 `json:"min_basis_coverage_pct"`
	SourceScopePass          bool    `json:"source_scope_pass"`
	CoveragePass             bool    `json:"coverage_pass"`
	CommonOutputsZeroTrade   bool    `json:"common_outputs_zero_trade"`
	AuthorityCandidateOnly   bool    `json:"authority_candidate_only"`
	ForwardLabelsAsInputs    bool    `json:"forward_labels_as_inputs"`
	OrthogonalityGateApplied bool    `json:"orthogonality_gate_applied"`
	Trades                   int     `json:"trades"`
	StopState                string  `json:"stop_state"`
}

type derivativesContextSymbolData struct {
	symbol     string
	candles    []Candle
	source     FuturesBTCRegimeETHSOLContextSourceRow
	coverage   []FuturesBTCRegimeETHSOLContextCoverageRow
	states     []FuturesRangeStateConstructionLoopStateRow
	labels     []FuturesRangeStateConstructionLoopLabelRow
	sourceOK   bool
	coverageOK bool
}

type derivativesContextFeatureKey struct {
	symbol     string
	timeframe  string
	stateRowID int
}

type derivativesContextCohortKey struct {
	cohortID      string
	comparison    string
	symbol        string
	split         string
	timeframe     string
	horizonBars   int
	route         string
	localBucket   string
	bucketType    string
	derivBucketID string
	basisLevel    string
	basisChange   string
	premium       string
}

type derivativesContextCohortAccumulator struct {
	row    FuturesDerivativesContextAuditCohortRow
	labels map[string]int
}

func DefaultFuturesDerivativesContextAuditConfig() FuturesDerivativesContextAuditConfig {
	btcCfg := DefaultFuturesBTCRegimeETHSOLContextAuditConfig()
	return FuturesDerivativesContextAuditConfig{
		SourceAuditConfig:                  DefaultFuturesDerivativesContextSourceAuditConfig(),
		LocalSources:                       append([]FuturesBTCRegimeETHSOLContextSourceConfig(nil), btcCfg.Sources...),
		StateConfig:                        DefaultFuturesRangeStateConstructionLoopAuditConfig(),
		MinBasisContextCoverage:            0.994472,
		BasisChangeLookbackIntervals:       12,
		BasisVolatilityLookbackIntervals:   36,
		MinFullCohortRows:                  300,
		MinSplitCohortRows:                 60,
		MaxSplitContributionRate:           0.60,
		MinUsefulRateFull:                  0.56,
		MinUsefulRateSplit:                 0.50,
		MaxToxicRateFull:                   0.44,
		MaxToxicRateSplit:                  0.50,
		MinUsefulMinusToxicMarginFull:      0.10,
		MinUsefulMinusToxicMarginSplit:     0.03,
		MinToxicRateFull:                   0.58,
		MinToxicRateSplit:                  0.52,
		MinContextUsefulImprovementFull:    0.04,
		MinContextUsefulImprovementSplit:   0.01,
		MinContextMarginImprovementFull:    0.04,
		MinContextMarginImprovementSplit:   0.01,
		MinContextToxicImprovementFull:     0.04,
		MinContextToxicImprovementSplit:    0.01,
		MaxBucketShareOfLocalBaselineFull:  0.90,
		MaxBucketShareOfLocalBaselineSplit: 0.95,
	}
}

func RunFuturesDerivativesContextAudit(cfg FuturesDerivativesContextAuditConfig, splits []Split) (FuturesDerivativesContextAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesDerivativesContextAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesDerivativesContextAuditResult{}
	if cfg.RescueClosedFamily {
		result.StopState = DerivativesContextAuditStopStateRejectedClosedFamilyRescue
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	sourceAudit, err := RunFuturesDerivativesContextSourceAudit(cfg.SourceAuditConfig, splits)
	if err != nil {
		result.StopState = DerivativesContextAuditStopStateSourceGap
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}
	result.SourceRows = derivativesContextSourceRowsFromSourceAudit(sourceAudit)
	for _, row := range sourceAudit.MissingnessRows {
		result.MissingnessRows = append(result.MissingnessRows, FuturesDerivativesContextAuditMissingnessRow{
			AuditName:         FuturesDerivativesContextAuditName,
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
	if sourceAudit.StopState != DerivativesContextSourceAuditStopStatePassedNeedsContextBrief ||
		!derivativesContextSourceAuditPass(sourceAudit, cfg) {
		result.StopState = DerivativesContextAuditStopStateSourceGap
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	dataBySymbol := map[string]derivativesContextSymbolData{}
	sourceOK := true
	btcCfg := FuturesBTCRegimeETHSOLContextAuditConfig{Sources: cfg.LocalSources, StateConfig: cfg.StateConfig}
	for _, source := range cfg.LocalSources {
		data := btcRegimeETHSOLContextLoadSymbol(source, btcCfg, splits)
		result.SourceRows = append(result.SourceRows, derivativesContextSourceRowFromAnchor(data.source))
		dataBySymbol[data.symbol] = derivativesContextSymbolData{
			symbol: data.symbol, candles: data.candles, source: data.source,
			sourceOK: data.sourceOK, coverageOK: data.coverageOK,
		}
		if !data.sourceOK {
			sourceOK = false
		}
	}
	if !sourceOK || !derivativesContextRequiredSymbolsPresent(dataBySymbol) {
		result.StopState = DerivativesContextAuditStopStateSourceGap
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	coverageOK := true
	for _, symbol := range derivativesContextSymbols() {
		base := btcRegimeContextSymbolData{
			symbol: dataBySymbol[symbol].symbol, candles: dataBySymbol[symbol].candles,
			source: dataBySymbol[symbol].source, sourceOK: dataBySymbol[symbol].sourceOK,
			skips: rangeStateSkipAccumulator{},
		}
		built, err := btcRegimeETHSOLContextBuildSymbolStates(base, btcCfg, splits)
		if err != nil {
			return result, err
		}
		dataBySymbol[symbol] = derivativesContextSymbolData{
			symbol: built.symbol, candles: built.candles, source: built.source,
			coverage: built.coverage, states: built.states, labels: built.labels,
			sourceOK: built.sourceOK, coverageOK: built.coverageOK,
		}
		if !built.coverageOK {
			coverageOK = false
		}
	}
	if !coverageOK {
		result.StopState = DerivativesContextAuditStopStateSourceGap
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	streams := derivativesContextLoadStreams(cfg)
	featureByState := map[derivativesContextFeatureKey]FuturesDerivativesContextAuditBasisFeatureRow{}
	missingBySymbolTimeframe := map[string]map[string]int{}
	for _, symbol := range derivativesContextSymbols() {
		data := dataBySymbol[symbol]
		for _, cov := range data.coverage {
			result.CoverageRows = append(result.CoverageRows, derivativesContextCoverageRowFromLocalCoverage(symbol, cov, cfg))
		}
		for _, state := range data.states {
			localRow := derivativesContextLocalStateRow(symbol, state)
			feature, ok, reason := derivativesContextFeatureForState(len(result.BasisFeatureRows)+1, symbol, state, streams, cfg)
			if ok {
				localRow.DerivativesContextAvailable = true
				localRow.ContextStateID = feature.ContextStateID
				localRow.BasisLevelBucket = feature.BasisLevelBucket
				localRow.BasisChangeBucket = feature.BasisChangeBucket
				localRow.PremiumLevelBucket = feature.PremiumLevelBucket
				result.BasisFeatureRows = append(result.BasisFeatureRows, feature)
				featureByState[derivativesContextFeatureKey{symbol: symbol, timeframe: state.Timeframe, stateRowID: state.StateRowID}] = feature
			} else {
				localRow.MissingContextReason = reason
				if missingBySymbolTimeframe[symbol] == nil {
					missingBySymbolTimeframe[symbol] = map[string]int{}
				}
				missingBySymbolTimeframe[symbol][state.Timeframe]++
			}
			result.LocalStateRows = append(result.LocalStateRows, localRow)
		}
	}

	labelID := 0
	for _, symbol := range derivativesContextSymbols() {
		for _, label := range dataBySymbol[symbol].labels {
			feature, ok := featureByState[derivativesContextFeatureKey{symbol: symbol, timeframe: label.Timeframe, stateRowID: label.StateRowID}]
			if !ok {
				continue
			}
			labelID++
			result.LabelRows = append(result.LabelRows, derivativesContextLabelRow(labelID, symbol, feature, label))
		}
	}

	result.CoverageRows = append(result.CoverageRows, derivativesContextFeatureCoverageRows(result, sourceAudit, cfg)...)
	for _, row := range derivativesContextFeatureMissingnessRows(result, missingBySymbolTimeframe, cfg) {
		result.MissingnessRows = append(result.MissingnessRows, row)
	}
	if !derivativesContextCoveragePass(result.CoverageRows) {
		result.StopState = DerivativesContextAuditStopStateSourceGap
		result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	result.CohortRows = derivativesContextCohortRows(result.LabelRows, cfg, splits, true)
	result.RankingRows = derivativesContextRankingRows(result.CohortRows)
	result.PassingCohorts = derivativesContextPassingRankingCount(result.RankingRows)
	result.StopState = FuturesDerivativesContextAuditStopState(result)
	result.SummaryRows = derivativesContextSummaryRows(result, cfg, splits)
	return result, nil
}

func FuturesDerivativesContextAuditStopState(result FuturesDerivativesContextAuditResult) string {
	if result.StopState == DerivativesContextAuditStopStateSourceGap ||
		result.StopState == DerivativesContextAuditStopStateRejectedFutureLabelLeak ||
		result.StopState == DerivativesContextAuditStopStateRejectedClosedFamilyRescue {
		return result.StopState
	}
	if !derivativesContextCoveragePass(result.CoverageRows) {
		return DerivativesContextAuditStopStateSourceGap
	}
	for _, row := range result.LocalStateRows {
		if row.ClosedFamilyRescue {
			return DerivativesContextAuditStopStateRejectedClosedFamilyRescue
		}
		if row.ForwardLabelsAsStateInput {
			return DerivativesContextAuditStopStateRejectedFutureLabelLeak
		}
	}
	for _, row := range result.LabelRows {
		if row.ForwardLabelUsedAsFeature {
			return DerivativesContextAuditStopStateRejectedFutureLabelLeak
		}
	}
	allCollinear := len(result.RankingRows) > 0
	for _, row := range result.RankingRows {
		if !row.FutureLeakProtectionPass {
			return DerivativesContextAuditStopStateRejectedFutureLabelLeak
		}
		if !row.ClosedFamilyProtectionPass {
			return DerivativesContextAuditStopStateRejectedClosedFamilyRescue
		}
		if row.PassesGate {
			return DerivativesContextAuditStopStatePassedNeedsSpec
		}
		if !strings.Contains(row.FailureReason, "basis_bucket_collinear_with_local_state") {
			allCollinear = false
		}
	}
	if allCollinear {
		return DerivativesContextAuditStopStateRejectedClosedFamilyRescue
	}
	return DerivativesContextAuditStopStateFailedNoUsableContext
}

func (cfg FuturesDerivativesContextAuditConfig) withDefaults() FuturesDerivativesContextAuditConfig {
	defaults := DefaultFuturesDerivativesContextAuditConfig()
	if len(cfg.SourceAuditConfig.DerivativeSources) == 0 {
		cfg.SourceAuditConfig = defaults.SourceAuditConfig
	}
	if len(cfg.LocalSources) == 0 {
		cfg.LocalSources = append([]FuturesBTCRegimeETHSOLContextSourceConfig(nil), defaults.LocalSources...)
	}
	if len(cfg.StateConfig.Timeframes) == 0 {
		cfg.StateConfig = defaults.StateConfig
	}
	btcDefaults := FuturesBTCRegimeETHSOLContextAuditConfig{
		Sources:     cfg.LocalSources,
		StateConfig: cfg.StateConfig,
	}.withDefaults()
	cfg.LocalSources = btcDefaults.Sources
	cfg.StateConfig = btcDefaults.StateConfig
	if cfg.MinBasisContextCoverage == 0 {
		cfg.MinBasisContextCoverage = defaults.MinBasisContextCoverage
	}
	if cfg.BasisChangeLookbackIntervals == 0 {
		cfg.BasisChangeLookbackIntervals = defaults.BasisChangeLookbackIntervals
	}
	if cfg.BasisVolatilityLookbackIntervals == 0 {
		cfg.BasisVolatilityLookbackIntervals = defaults.BasisVolatilityLookbackIntervals
	}
	if cfg.MinFullCohortRows == 0 {
		cfg.MinFullCohortRows = defaults.MinFullCohortRows
	}
	if cfg.MinSplitCohortRows == 0 {
		cfg.MinSplitCohortRows = defaults.MinSplitCohortRows
	}
	if cfg.MaxSplitContributionRate == 0 {
		cfg.MaxSplitContributionRate = defaults.MaxSplitContributionRate
	}
	if cfg.MinUsefulRateFull == 0 {
		cfg.MinUsefulRateFull = defaults.MinUsefulRateFull
	}
	if cfg.MinUsefulRateSplit == 0 {
		cfg.MinUsefulRateSplit = defaults.MinUsefulRateSplit
	}
	if cfg.MaxToxicRateFull == 0 {
		cfg.MaxToxicRateFull = defaults.MaxToxicRateFull
	}
	if cfg.MaxToxicRateSplit == 0 {
		cfg.MaxToxicRateSplit = defaults.MaxToxicRateSplit
	}
	if cfg.MinUsefulMinusToxicMarginFull == 0 {
		cfg.MinUsefulMinusToxicMarginFull = defaults.MinUsefulMinusToxicMarginFull
	}
	if cfg.MinUsefulMinusToxicMarginSplit == 0 {
		cfg.MinUsefulMinusToxicMarginSplit = defaults.MinUsefulMinusToxicMarginSplit
	}
	if cfg.MinToxicRateFull == 0 {
		cfg.MinToxicRateFull = defaults.MinToxicRateFull
	}
	if cfg.MinToxicRateSplit == 0 {
		cfg.MinToxicRateSplit = defaults.MinToxicRateSplit
	}
	if cfg.MinContextUsefulImprovementFull == 0 {
		cfg.MinContextUsefulImprovementFull = defaults.MinContextUsefulImprovementFull
	}
	if cfg.MinContextUsefulImprovementSplit == 0 {
		cfg.MinContextUsefulImprovementSplit = defaults.MinContextUsefulImprovementSplit
	}
	if cfg.MinContextMarginImprovementFull == 0 {
		cfg.MinContextMarginImprovementFull = defaults.MinContextMarginImprovementFull
	}
	if cfg.MinContextMarginImprovementSplit == 0 {
		cfg.MinContextMarginImprovementSplit = defaults.MinContextMarginImprovementSplit
	}
	if cfg.MinContextToxicImprovementFull == 0 {
		cfg.MinContextToxicImprovementFull = defaults.MinContextToxicImprovementFull
	}
	if cfg.MinContextToxicImprovementSplit == 0 {
		cfg.MinContextToxicImprovementSplit = defaults.MinContextToxicImprovementSplit
	}
	if cfg.MaxBucketShareOfLocalBaselineFull == 0 {
		cfg.MaxBucketShareOfLocalBaselineFull = defaults.MaxBucketShareOfLocalBaselineFull
	}
	if cfg.MaxBucketShareOfLocalBaselineSplit == 0 {
		cfg.MaxBucketShareOfLocalBaselineSplit = defaults.MaxBucketShareOfLocalBaselineSplit
	}
	return cfg
}

func (cfg FuturesDerivativesContextAuditConfig) validate() error {
	if err := cfg.StateConfig.withDefaults().validate(); err != nil {
		return err
	}
	if len(cfg.LocalSources) != 3 {
		return fmt.Errorf("derivatives context audit requires exactly BTCUSDT, ETHUSDT, and SOLUSDT candle anchors")
	}
	seen := map[string]bool{}
	for _, source := range cfg.LocalSources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if symbol != RangeUniverseSymbolBTCUSDT && symbol != RangeUniverseSymbolETHUSDT && symbol != RangeUniverseSymbolSOLUSDT {
			return fmt.Errorf("derivatives context audit unsupported local symbol %q", source.Symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("derivatives context audit duplicate local symbol %s", symbol)
		}
		seen[symbol] = true
	}
	for _, symbol := range derivativesContextSymbols() {
		if !seen[symbol] {
			return fmt.Errorf("derivatives context audit missing local symbol %s", symbol)
		}
	}
	if cfg.MinBasisContextCoverage <= 0 || cfg.MinBasisContextCoverage > 1 {
		return fmt.Errorf("derivatives context audit basis coverage floor must be in (0,1]")
	}
	if cfg.MinFullCohortRows <= 0 || cfg.MinSplitCohortRows <= 0 {
		return fmt.Errorf("derivatives context audit cohort count gates must be positive")
	}
	if cfg.MaxSplitContributionRate <= 0 || cfg.MaxSplitContributionRate > 1 {
		return fmt.Errorf("derivatives context audit max split contribution must be in (0,1]")
	}
	return nil
}

func derivativesContextSourceRowsFromSourceAudit(sourceAudit FuturesDerivativesContextSourceAuditResult) []FuturesDerivativesContextAuditSourceRow {
	rows := make([]FuturesDerivativesContextAuditSourceRow, 0, len(sourceAudit.SourceRows))
	for _, row := range sourceAudit.SourceRows {
		rows = append(rows, FuturesDerivativesContextAuditSourceRow{
			AuditName:        FuturesDerivativesContextAuditName,
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

func derivativesContextSourceRowFromAnchor(source FuturesBTCRegimeETHSOLContextSourceRow) FuturesDerivativesContextAuditSourceRow {
	return FuturesDerivativesContextAuditSourceRow{
		AuditName:                  FuturesDerivativesContextAuditName,
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

func derivativesContextSourceAuditPass(sourceAudit FuturesDerivativesContextSourceAuditResult, cfg FuturesDerivativesContextAuditConfig) bool {
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
	if required != 6 || accepted != 6 {
		return false
	}
	for _, row := range sourceAudit.TimestampAlignmentRows {
		if row.Required && (!row.MeetsMinCoverage || row.LagCoveragePct < cfg.SourceAuditConfig.MinAlignedCoverage) {
			return false
		}
		if row.UsesFutureRows || !row.ExactClosedIntervalJoin {
			return false
		}
	}
	return true
}

func derivativesContextLoadStreams(cfg FuturesDerivativesContextAuditConfig) map[string]*derivStreamData {
	streams := map[string]*derivStreamData{}
	eraStartMs := cfg.SourceAuditConfig.EraStartMs
	eraEndMs := cfg.SourceAuditConfig.EraEndMs
	if eraStartMs == 0 {
		eraStartMs = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	}
	if eraEndMs == 0 {
		eraEndMs = time.Date(2026, 6, 16, 23, 55, 0, 0, time.UTC).UnixMilli()
	}
	for _, source := range cfg.SourceAuditConfig.DerivativeSources {
		streams[derivativesContextStreamKey(source.Symbol, source.SourceFamily)] = loadDerivativesStream(source, eraStartMs, eraEndMs)
	}
	return streams
}

func derivativesContextFeatureForState(id int, symbol string, state FuturesRangeStateConstructionLoopStateRow, streams map[string]*derivStreamData, cfg FuturesDerivativesContextAuditConfig) (FuturesDerivativesContextAuditBasisFeatureRow, bool, string) {
	decisionClose, err := time.Parse(timeLayout, state.DecisionCloseTime)
	if err != nil {
		return FuturesDerivativesContextAuditBasisFeatureRow{}, false, "invalid_decision_close_time"
	}
	lagIntervals := cfg.SourceAuditConfig.ConservativeLagIntervals
	if lagIntervals <= 0 {
		lagIntervals = 1
	}
	sourceOpenMs := derivativesContextLaggedSourceOpen(decisionClose.UnixMilli(), lagIntervals)
	sourceCloseMs := sourceOpenMs + derivativesIntervalMs - 1
	mark := streams[derivativesContextStreamKey(symbol, "mark_price_klines")]
	index := streams[derivativesContextStreamKey(symbol, "index_price_klines")]
	if mark == nil || index == nil {
		return FuturesDerivativesContextAuditBasisFeatureRow{}, false, "missing_required_basis_stream"
	}
	markClose, markOK := mark.openClose[sourceOpenMs]
	indexClose, indexOK := index.openClose[sourceOpenMs]
	if !markOK || !indexOK || indexClose == 0 {
		return FuturesDerivativesContextAuditBasisFeatureRow{}, false, "missing_required_lagged_basis_context"
	}
	basisRaw := markClose - indexClose
	basisBPS := basisRaw / indexClose * 10000
	lookbackOpen := sourceOpenMs - int64(cfg.BasisChangeLookbackIntervals)*derivativesIntervalMs
	basisChangePresent := false
	basisChangeBPS := 0.0
	if prevMark, ok := mark.openClose[lookbackOpen]; ok {
		if prevIndex, ok := index.openClose[lookbackOpen]; ok && prevIndex != 0 {
			prevBasisBPS := (prevMark - prevIndex) / prevIndex * 10000
			basisChangeBPS = basisBPS - prevBasisBPS
			basisChangePresent = true
		}
	}
	volBPS, volBucket := derivativesContextBasisVolatility(mark, index, sourceOpenMs, cfg.BasisVolatilityLookbackIntervals)
	premiumClose := 0.0
	premiumBPS := 0.0
	premiumPresent := false
	if premium := streams[derivativesContextStreamKey(symbol, "premium_index_klines")]; premium != nil {
		if v, ok := premium.openClose[sourceOpenMs]; ok {
			premiumClose = v
			premiumBPS = v * 10000
			premiumPresent = true
		}
	}
	localBucket := btcRegimeETHSOLContextLocalBucketID(state)
	changeBucket := derivativesContextBasisChangeBucket(basisChangeBPS, basisChangePresent)
	premiumBucket := derivativesContextPremiumBucket(premiumBPS, premiumPresent)
	contextID := fmt.Sprintf("derivatives_context_v1::%s::%s::%s::basis=%s::change=%s::premium=%s",
		symbol, state.Timeframe, localBucket, derivativesContextBasisLevelBucket(basisBPS), changeBucket, premiumBucket)
	return FuturesDerivativesContextAuditBasisFeatureRow{
		FeatureRowID:                 id,
		Symbol:                       symbol,
		LocalStateRowID:              state.StateRowID,
		Timestamp:                    state.Timestamp,
		Timeframe:                    state.Timeframe,
		Split:                        state.Split,
		DecisionIndex:                state.DecisionIndex,
		DecisionCloseTime:            state.DecisionCloseTime,
		SourceOpenTime:               derivativesFmtMs(sourceOpenMs),
		SourceCloseTime:              derivativesFmtMs(sourceCloseMs),
		SourceLagIntervals:           lagIntervals,
		AntiLookaheadRule:            "source_close_time + 5m <= decision_candle_close_time",
		MarkClose:                    markClose,
		IndexClose:                   indexClose,
		BasisRaw:                     basisRaw,
		BasisBPS:                     basisBPS,
		BasisSignBucket:              derivativesContextBasisSignBucket(basisBPS),
		BasisLevelBucket:             derivativesContextBasisLevelBucket(basisBPS),
		BasisChangeLookbackIntervals: cfg.BasisChangeLookbackIntervals,
		BasisChangeBPS:               basisChangeBPS,
		BasisChangeBucket:            changeBucket,
		BasisChangePresent:           basisChangePresent,
		BasisVolatilityBPS:           volBPS,
		BasisVolatilityBucket:        volBucket,
		PremiumClose:                 premiumClose,
		PremiumBPS:                   premiumBPS,
		PremiumLevelBucket:           premiumBucket,
		PremiumPresent:               premiumPresent,
		LocalRangeBucketID:           localBucket,
		ContextStateID:               contextID,
		ClosedCandleOnly:             true,
		UsesFutureRows:               false,
		ForwardLabelsAsFeatureInput:  false,
	}, true, ""
}

func derivativesContextLocalStateRow(symbol string, state FuturesRangeStateConstructionLoopStateRow) FuturesDerivativesContextAuditLocalStateRow {
	return FuturesDerivativesContextAuditLocalStateRow{
		LocalStateRowID:           state.StateRowID,
		Symbol:                    symbol,
		Timestamp:                 state.Timestamp,
		Timeframe:                 state.Timeframe,
		Split:                     state.Split,
		DecisionIndex:             state.DecisionIndex,
		DecisionCloseTime:         state.DecisionCloseTime,
		RangeEpisodeID:            state.RangeEpisodeID,
		LocalRangeBucketID:        btcRegimeETHSOLContextLocalBucketID(state),
		StateID:                   state.StateID,
		GeometryBucket:            state.GeometryBucket,
		VolBucket:                 state.VolBucket,
		TrendBucket:               state.TrendBucket,
		ImpulseBucket:             state.ImpulseBucket,
		ParticipationBucket:       state.ParticipationBucket,
		ClosedCandleOnly:          true,
		AuthorityCandidateOnly:    true,
		ForwardLabelsAsStateInput: false,
		ClosedFamilyRescue:        false,
	}
}

func derivativesContextLabelRow(id int, symbol string, feature FuturesDerivativesContextAuditBasisFeatureRow, label FuturesRangeStateConstructionLoopLabelRow) FuturesDerivativesContextAuditLabelRow {
	return FuturesDerivativesContextAuditLabelRow{
		LabelRowID:                id,
		Symbol:                    symbol,
		LocalStateRowID:           feature.LocalStateRowID,
		FeatureRowID:              feature.FeatureRowID,
		Timestamp:                 label.Timestamp,
		Timeframe:                 label.Timeframe,
		Split:                     label.Split,
		HorizonBars:               label.HorizonBars,
		LocalRangeBucketID:        feature.LocalRangeBucketID,
		BasisLevelBucket:          feature.BasisLevelBucket,
		BasisChangeBucket:         feature.BasisChangeBucket,
		BasisChangePresent:        feature.BasisChangePresent,
		PremiumLevelBucket:        feature.PremiumLevelBucket,
		PremiumPresent:            feature.PremiumPresent,
		ContextStateID:            feature.ContextStateID,
		ForwardLabel:              label.ForwardLabel,
		RotationUseful:            label.RotationUseful,
		RotationToxic:             label.RotationToxic,
		ContinuationUseful:        label.ContinuationUseful,
		ContinuationToxic:         label.ContinuationToxic,
		NoTradeToxic:              label.NoTradeToxic,
		DiagnosticOnly:            label.DiagnosticOnly,
		LabelWindowStartIndex:     label.LabelWindowStartIndex,
		LabelWindowEndIndex:       label.LabelWindowEndIndex,
		LabelWindowStartTime:      label.LabelWindowStartTime,
		LabelWindowEndTime:        label.LabelWindowEndTime,
		ForwardLabelMetadataOnly:  true,
		ForwardLabelUsedAsFeature: false,
	}
}

func derivativesContextCohortRows(labels []FuturesDerivativesContextAuditLabelRow, cfg FuturesDerivativesContextAuditConfig, splits []Split, sourcePass bool) []FuturesDerivativesContextAuditCohortRow {
	acc := map[derivativesContextCohortKey]*derivativesContextCohortAccumulator{}
	for _, label := range labels {
		for _, split := range rangeDiscoverySplitCombos(label.Split) {
			for _, route := range []string{RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation, RangeStateConstructionLoopRouteNoTradeToxic} {
				for _, key := range derivativesContextCohortKeys(label, split, route) {
					a := acc[key]
					if a == nil {
						a = &derivativesContextCohortAccumulator{labels: map[string]int{}}
						a.row = FuturesDerivativesContextAuditCohortRow{
							CohortID:                   key.cohortID,
							ComparisonType:             key.comparison,
							Symbol:                     key.symbol,
							Split:                      key.split,
							Timeframe:                  key.timeframe,
							HorizonBars:                key.horizonBars,
							RouteCandidate:             key.route,
							LocalRangeBucketID:         key.localBucket,
							DerivativesBucketType:      key.bucketType,
							DerivativesBucketID:        key.derivBucketID,
							BasisLevelBucket:           key.basisLevel,
							BasisChangeBucket:          key.basisChange,
							PremiumLevelBucket:         key.premium,
							ClosedFamilyProtectionPass: true,
							FutureLeakProtectionPass:   true,
							MissingnessGatePass:        true,
							PremiumOnlyGatePass:        key.bucketType != "premium_only",
						}
						if !sourcePass {
							a.row.FailureReason = "source_or_coverage_gap"
						}
						acc[key] = a
					}
					a.add(label)
				}
			}
		}
	}
	rows := make([]FuturesDerivativesContextAuditCohortRow, 0, len(acc))
	for _, a := range acc {
		rows = append(rows, a.finalRow())
	}
	derivativesContextMarkCohortGates(rows, cfg, splits, sourcePass)
	sort.Slice(rows, func(i, j int) bool { return derivativesContextLessCohort(rows[i], rows[j]) })
	return rows
}

func (acc *derivativesContextCohortAccumulator) add(label FuturesDerivativesContextAuditLabelRow) {
	acc.row.CandidateCount++
	acc.labels[label.ForwardLabel]++
	if derivativesContextRouteUseful(acc.row.RouteCandidate, label) {
		acc.row.UsefulCount++
	}
	if derivativesContextRouteToxic(acc.row.RouteCandidate, label) {
		acc.row.ToxicCount++
	}
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
}

func (acc *derivativesContextCohortAccumulator) finalRow() FuturesDerivativesContextAuditCohortRow {
	row := acc.row
	if row.CandidateCount > 0 {
		row.UsefulRate = float64(row.UsefulCount) / float64(row.CandidateCount)
		row.ToxicRate = float64(row.ToxicCount) / float64(row.CandidateCount)
		row.UsefulMinusToxicMargin = row.UsefulRate - row.ToxicRate
	}
	row.DominantForwardLabel, row.DominantForwardLabelRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, false)
	return row
}

func derivativesContextCohortKeys(label FuturesDerivativesContextAuditLabelRow, split, route string) []derivativesContextCohortKey {
	keys := []derivativesContextCohortKey{
		derivativesContextCohortKeyFor(label, split, route, DerivativesContextComparisonLocalOnly, "", "", "", ""),
		derivativesContextCohortKeyFor(label, split, route, DerivativesContextComparisonLocalPlusDeriv, DerivativesContextBucketBasisLevel, label.BasisLevelBucket, "", ""),
	}
	if label.BasisChangePresent {
		keys = append(keys, derivativesContextCohortKeyFor(label, split, route, DerivativesContextComparisonLocalPlusDeriv, DerivativesContextBucketBasisLevelChange, label.BasisLevelBucket, label.BasisChangeBucket, ""))
	}
	if label.PremiumPresent {
		keys = append(keys, derivativesContextCohortKeyFor(label, split, route, DerivativesContextComparisonLocalPlusDeriv, DerivativesContextBucketBasisLevelPremium, label.BasisLevelBucket, "", label.PremiumLevelBucket))
	}
	return keys
}

func derivativesContextCohortKeyFor(label FuturesDerivativesContextAuditLabelRow, split, route, comparison, bucketType, basisLevel, basisChange, premium string) derivativesContextCohortKey {
	bucketID := ""
	parts := []string{"derivatives_context_v1", comparison, label.Symbol, label.Timeframe, fmt.Sprintf("h%d", label.HorizonBars), route, "local=" + label.LocalRangeBucketID}
	if comparison == DerivativesContextComparisonLocalPlusDeriv {
		bucketParts := []string{"basis=" + basisLevel}
		if basisChange != "" {
			bucketParts = append(bucketParts, "change="+basisChange)
		}
		if premium != "" {
			bucketParts = append(bucketParts, "premium="+premium)
		}
		bucketID = strings.Join(bucketParts, "::")
		parts = append(parts, "bucket_type="+bucketType, bucketID)
	}
	return derivativesContextCohortKey{
		cohortID:      strings.Join(parts, "|"),
		comparison:    comparison,
		symbol:        label.Symbol,
		split:         split,
		timeframe:     label.Timeframe,
		horizonBars:   label.HorizonBars,
		route:         route,
		localBucket:   label.LocalRangeBucketID,
		bucketType:    bucketType,
		derivBucketID: bucketID,
		basisLevel:    basisLevel,
		basisChange:   basisChange,
		premium:       premium,
	}
}

func derivativesContextMarkCohortGates(rows []FuturesDerivativesContextAuditCohortRow, cfg FuturesDerivativesContextAuditConfig, splits []Split, sourcePass bool) {
	byIDSplit := map[string]map[string]*FuturesDerivativesContextAuditCohortRow{}
	localBaseline := map[string]map[string]*FuturesDerivativesContextAuditCohortRow{}
	for i := range rows {
		if byIDSplit[rows[i].CohortID] == nil {
			byIDSplit[rows[i].CohortID] = map[string]*FuturesDerivativesContextAuditCohortRow{}
		}
		byIDSplit[rows[i].CohortID][rows[i].Split] = &rows[i]
		if rows[i].ComparisonType == DerivativesContextComparisonLocalOnly {
			key := derivativesContextBaselineKey(rows[i])
			if localBaseline[key] == nil {
				localBaseline[key] = map[string]*FuturesDerivativesContextAuditCohortRow{}
			}
			localBaseline[key][rows[i].Split] = &rows[i]
		}
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	for i := range rows {
		cohortRows := byIDSplit[rows[i].CohortID]
		full := cohortRows[fullSplitName]
		reasons := []string{}
		if full != nil {
			rows[i].FullPeriodRows = full.CandidateCount
			rows[i].FullUsefulRate = full.UsefulRate
			rows[i].FullToxicRate = full.ToxicRate
			rows[i].FullUsefulMinusToxicMargin = full.UsefulMinusToxicMargin
		}
		rows[i].WeakestSplitRows = math.MaxInt
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
			if row.CandidateCount < cfg.MinSplitCohortRows {
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
			if rows[i].FullPeriodRows > 0 {
				maxContribution = math.Max(maxContribution, float64(row.CandidateCount)/float64(rows[i].FullPeriodRows))
			}
		}
		if rows[i].WeakestSplitRows == math.MaxInt {
			rows[i].WeakestSplitRows = 0
		}
		if math.IsInf(rows[i].WeakestSplitUsefulRate, 1) {
			rows[i].WeakestSplitUsefulRate = 0
		}
		if math.IsInf(rows[i].WeakestSplitMargin, 1) {
			rows[i].WeakestSplitMargin = 0
		}
		rows[i].MaxSplitContributionRate = maxContribution
		rows[i].ReviewableCountGatePass = rows[i].FullPeriodRows >= cfg.MinFullCohortRows && splitCountsPass
		rows[i].SplitStabilityGatePass = splitRatesPresent && len(periodSplits) > 0
		rows[i].SplitContributionGatePass = maxContribution <= cfg.MaxSplitContributionRate || len(periodSplits) == 0
		rows[i].RouteRateGatePass = derivativesContextRouteRateGate(rows[i], cfg)
		rows[i].ClosedFamilyProtectionPass = true
		rows[i].FutureLeakProtectionPass = !strings.Contains(rows[i].CohortID, "forward_label")
		rows[i].MissingnessGatePass = !strings.Contains(rows[i].DerivativesBucketID, "missing")
		rows[i].PremiumOnlyGatePass = rows[i].DerivativesBucketType != "premium_only"

		if rows[i].ComparisonType == DerivativesContextComparisonLocalPlusDeriv {
			derivativesContextFillBaselineComparison(&rows[i], cohortRows, localBaseline[derivativesContextBaselineKey(rows[i])], periodSplits)
			rows[i].DerivativesImprovementGatePass = derivativesContextImprovementGate(rows[i], cfg)
			rows[i].OrthogonalityGatePass = rows[i].DerivativesImprovementGatePass &&
				rows[i].BucketShareOfLocalBaselineFull <= cfg.MaxBucketShareOfLocalBaselineFull &&
				rows[i].MaxBucketShareOfLocalBaselineSplit <= cfg.MaxBucketShareOfLocalBaselineSplit
		}

		if !sourcePass {
			reasons = append(reasons, "source_or_coverage_gap")
		}
		if rows[i].ComparisonType != DerivativesContextComparisonLocalPlusDeriv {
			reasons = append(reasons, "local_only_baseline_not_context_candidate")
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
		if !rows[i].RouteRateGatePass {
			reasons = append(reasons, "route_rate_gate_failed")
		}
		if rows[i].ComparisonType == DerivativesContextComparisonLocalPlusDeriv {
			if !rows[i].DerivativesImprovementGatePass {
				reasons = append(reasons, "derivatives_context_improvement_gate_failed")
			}
			if rows[i].BucketShareOfLocalBaselineFull > cfg.MaxBucketShareOfLocalBaselineFull ||
				rows[i].MaxBucketShareOfLocalBaselineSplit > cfg.MaxBucketShareOfLocalBaselineSplit {
				reasons = append(reasons, "basis_bucket_collinear_with_local_state")
			}
			if !rows[i].OrthogonalityGatePass {
				reasons = append(reasons, "orthogonality_gate_failed")
			}
		}
		if !rows[i].MissingnessGatePass {
			reasons = append(reasons, "missing_context_bucket")
		}
		if !rows[i].PremiumOnlyGatePass {
			reasons = append(reasons, "premium_only_context_forbidden")
		}
		if !rows[i].FutureLeakProtectionPass {
			reasons = append(reasons, "future_label_as_feature")
		}
		rows[i].PassesReviewGate = sourcePass &&
			rows[i].ComparisonType == DerivativesContextComparisonLocalPlusDeriv &&
			rows[i].ReviewableCountGatePass &&
			rows[i].SplitStabilityGatePass &&
			rows[i].SplitContributionGatePass &&
			rows[i].RouteRateGatePass &&
			rows[i].DerivativesImprovementGatePass &&
			rows[i].OrthogonalityGatePass &&
			rows[i].MissingnessGatePass &&
			rows[i].PremiumOnlyGatePass &&
			rows[i].ClosedFamilyProtectionPass &&
			rows[i].FutureLeakProtectionPass &&
			rows[i].Split == fullSplitName
		rows[i].FailureReason = uniqueJoinedReasons(reasons)
	}
}

func derivativesContextFillBaselineComparison(row *FuturesDerivativesContextAuditCohortRow, derivBySplit map[string]*FuturesDerivativesContextAuditCohortRow, baselineBySplit map[string]*FuturesDerivativesContextAuditCohortRow, periodSplits []Split) {
	if baselineBySplit == nil {
		return
	}
	full := baselineBySplit[fullSplitName]
	if full != nil {
		row.BaselineFullRows = full.CandidateCount
		row.BaselineFullUsefulRate = full.UsefulRate
		row.BaselineFullToxicRate = full.ToxicRate
		row.BaselineFullMargin = full.UsefulMinusToxicMargin
		row.FullUsefulImprovement = row.FullUsefulRate - full.UsefulRate
		row.FullToxicImprovement = row.FullToxicRate - full.ToxicRate
		row.FullMarginImprovement = row.FullUsefulMinusToxicMargin - full.UsefulMinusToxicMargin
		if full.CandidateCount > 0 {
			row.BucketShareOfLocalBaselineFull = float64(row.FullPeriodRows) / float64(full.CandidateCount)
		}
	}
	row.BaselineWeakestSplitUsefulRate = math.Inf(1)
	row.BaselineWeakestSplitMargin = math.Inf(1)
	for _, split := range periodSplits {
		base := baselineBySplit[split.Name]
		if base == nil {
			row.BaselineWeakestSplitUsefulRate = 0
			row.BaselineWeakestSplitMargin = 0
			continue
		}
		deriv := derivBySplit[split.Name]
		if base.UsefulRate < row.BaselineWeakestSplitUsefulRate {
			row.BaselineWeakestSplitUsefulRate = base.UsefulRate
		}
		if base.ToxicRate > row.BaselineWorstSplitToxicRate {
			row.BaselineWorstSplitToxicRate = base.ToxicRate
		}
		if base.UsefulMinusToxicMargin < row.BaselineWeakestSplitMargin {
			row.BaselineWeakestSplitMargin = base.UsefulMinusToxicMargin
		}
		if deriv != nil && base.CandidateCount > 0 {
			row.MaxBucketShareOfLocalBaselineSplit = math.Max(row.MaxBucketShareOfLocalBaselineSplit, safeDiv(float64(deriv.CandidateCount), float64(base.CandidateCount)))
		}
	}
	if math.IsInf(row.BaselineWeakestSplitUsefulRate, 1) {
		row.BaselineWeakestSplitUsefulRate = 0
	}
	if math.IsInf(row.BaselineWeakestSplitMargin, 1) {
		row.BaselineWeakestSplitMargin = 0
	}
	row.WeakestSplitUsefulImprovement = row.WeakestSplitUsefulRate - row.BaselineWeakestSplitUsefulRate
	row.WeakestSplitToxicImprovement = row.WorstSplitToxicRate - row.BaselineWorstSplitToxicRate
	row.WeakestSplitMarginImprovement = row.WeakestSplitMargin - row.BaselineWeakestSplitMargin
}

func derivativesContextRankingRows(cohorts []FuturesDerivativesContextAuditCohortRow) []FuturesDerivativesContextAuditRankingRow {
	rows := []FuturesDerivativesContextAuditRankingRow{}
	for _, cohort := range cohorts {
		if cohort.Split != fullSplitName || cohort.ComparisonType != DerivativesContextComparisonLocalPlusDeriv {
			continue
		}
		row := FuturesDerivativesContextAuditRankingRow{
			CohortID:                           cohort.CohortID,
			Symbol:                             cohort.Symbol,
			Timeframe:                          cohort.Timeframe,
			HorizonBars:                        cohort.HorizonBars,
			RouteCandidate:                     cohort.RouteCandidate,
			LocalRangeBucketID:                 cohort.LocalRangeBucketID,
			DerivativesBucketType:              cohort.DerivativesBucketType,
			DerivativesBucketID:                cohort.DerivativesBucketID,
			PassesGate:                         cohort.PassesReviewGate,
			FullPeriodRows:                     cohort.FullPeriodRows,
			WeakestSplitRows:                   cohort.WeakestSplitRows,
			MaxSplitContributionRate:           cohort.MaxSplitContributionRate,
			FullUsefulRate:                     cohort.FullUsefulRate,
			WeakestSplitUsefulRate:             cohort.WeakestSplitUsefulRate,
			FullToxicRate:                      cohort.FullToxicRate,
			WorstSplitToxicRate:                cohort.WorstSplitToxicRate,
			FullUsefulMinusToxicMargin:         cohort.FullUsefulMinusToxicMargin,
			WeakestSplitMargin:                 cohort.WeakestSplitMargin,
			FullUsefulImprovement:              cohort.FullUsefulImprovement,
			WeakestSplitUsefulImprovement:      cohort.WeakestSplitUsefulImprovement,
			FullToxicImprovement:               cohort.FullToxicImprovement,
			WeakestSplitToxicImprovement:       cohort.WeakestSplitToxicImprovement,
			FullMarginImprovement:              cohort.FullMarginImprovement,
			WeakestSplitMarginImprovement:      cohort.WeakestSplitMarginImprovement,
			BucketShareOfLocalBaselineFull:     cohort.BucketShareOfLocalBaselineFull,
			MaxBucketShareOfLocalBaselineSplit: cohort.MaxBucketShareOfLocalBaselineSplit,
			OrthogonalityGatePass:              cohort.OrthogonalityGatePass,
			DerivativesImprovementGatePass:     cohort.DerivativesImprovementGatePass,
			DominantForwardLabel:               cohort.DominantForwardLabel,
			DominantForwardLabelRate:           cohort.DominantForwardLabelRate,
			FutureLeakProtectionPass:           cohort.FutureLeakProtectionPass,
			ClosedFamilyProtectionPass:         cohort.ClosedFamilyProtectionPass,
			FailureReason:                      cohort.FailureReason,
		}
		row.RankScore = derivativesContextRankScore(row)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		if rows[i].Symbol != rows[j].Symbol {
			return rangeUniverseSymbolSortKey(rows[i].Symbol) < rangeUniverseSymbolSortKey(rows[j].Symbol)
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
		if rows[i].DerivativesBucketType != rows[j].DerivativesBucketType {
			return rows[i].DerivativesBucketType < rows[j].DerivativesBucketType
		}
		return rows[i].CohortID < rows[j].CohortID
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func derivativesContextSummaryRows(result FuturesDerivativesContextAuditResult, cfg FuturesDerivativesContextAuditConfig, splits []Split) []FuturesDerivativesContextAuditSummaryRow {
	coveragePass := derivativesContextCoveragePass(result.CoverageRows)
	minCov := derivativesContextMinCoverage(result.CoverageRows)
	rows := []FuturesDerivativesContextAuditSummaryRow{{
		Split:                    fullSplitName,
		Symbol:                   "all",
		Timeframe:                "all",
		HorizonBars:              0,
		SourceRows:               len(result.SourceRows),
		CoverageRows:             len(result.CoverageRows),
		BasisFeatureRows:         len(result.BasisFeatureRows),
		LocalStateRows:           len(result.LocalStateRows),
		LabelRows:                len(result.LabelRows),
		CohortRows:               len(result.CohortRows),
		RankingRows:              len(result.RankingRows),
		PassingCohorts:           result.PassingCohorts,
		MissingnessRows:          len(result.MissingnessRows),
		RequiredCoverageFloor:    cfg.MinBasisContextCoverage,
		MinBasisCoveragePct:      minCov,
		SourceScopePass:          coveragePass,
		CoveragePass:             coveragePass,
		CommonOutputsZeroTrade:   true,
		AuthorityCandidateOnly:   true,
		ForwardLabelsAsInputs:    false,
		OrthogonalityGateApplied: true,
		Trades:                   0,
		StopState:                result.StopState,
	}}
	for _, split := range splits {
		for _, symbol := range derivativesContextSymbols() {
			for _, timeframe := range cfg.StateConfig.withDefaults().Timeframes {
				for _, horizon := range rangeStateConstructionLoopHorizons(timeframe, cfg.StateConfig.withDefaults()) {
					row := FuturesDerivativesContextAuditSummaryRow{
						Split:                    split.Name,
						Symbol:                   symbol,
						Timeframe:                timeframe,
						HorizonBars:              horizon,
						SourceRows:               len(result.SourceRows),
						CoverageRows:             len(result.CoverageRows),
						MissingnessRows:          len(result.MissingnessRows),
						RequiredCoverageFloor:    cfg.MinBasisContextCoverage,
						MinBasisCoveragePct:      minCov,
						SourceScopePass:          coveragePass,
						CoveragePass:             coveragePass,
						CommonOutputsZeroTrade:   true,
						AuthorityCandidateOnly:   true,
						ForwardLabelsAsInputs:    false,
						OrthogonalityGateApplied: true,
						Trades:                   0,
						StopState:                result.StopState,
					}
					for _, local := range result.LocalStateRows {
						if local.Symbol == symbol && local.Timeframe == timeframe && rangeStateRowInSplit(local.DecisionCloseTime, split) {
							row.LocalStateRows++
						}
					}
					for _, feature := range result.BasisFeatureRows {
						if feature.Symbol == symbol && feature.Timeframe == timeframe && rangeStateRowInSplit(feature.DecisionCloseTime, split) {
							row.BasisFeatureRows++
						}
					}
					for _, label := range result.LabelRows {
						if label.Symbol == symbol && label.Timeframe == timeframe && label.HorizonBars == horizon && rangeStateRowInSplit(label.Timestamp, split) {
							row.LabelRows++
						}
					}
					for _, cohort := range result.CohortRows {
						if cohort.Symbol == symbol && cohort.Timeframe == timeframe && cohort.HorizonBars == horizon && cohort.Split == split.Name {
							row.CohortRows++
							if cohort.PassesReviewGate {
								row.PassingCohorts++
							}
						}
					}
					for _, ranking := range result.RankingRows {
						if ranking.Symbol == symbol && ranking.Timeframe == timeframe && ranking.HorizonBars == horizon {
							row.RankingRows++
						}
					}
					rows = append(rows, row)
				}
			}
		}
	}
	return rows
}

func derivativesContextCoverageRowFromLocalCoverage(symbol string, base FuturesBTCRegimeETHSOLContextCoverageRow, cfg FuturesDerivativesContextAuditConfig) FuturesDerivativesContextAuditCoverageRow {
	return FuturesDerivativesContextAuditCoverageRow{
		AuditName:             FuturesDerivativesContextAuditName,
		Scope:                 "local_range_state_resample",
		Symbol:                symbol,
		Timeframe:             base.Timeframe,
		SourceFamily:          "trade_klines",
		Required:              true,
		RowCount:              base.RowCount,
		RequiredCoverageFloor: cfg.MinBasisContextCoverage,
		CoverageGatePass:      base.CoverageFactsPass && base.ValidationStatus == "accepted" && base.Complete,
		ClosedCandleOnly:      true,
		ForwardLabelsAsInputs: false,
		ValidationStatus:      base.ValidationStatus,
		ValidationError:       base.ValidationError,
	}
}

func derivativesContextFeatureCoverageRows(result FuturesDerivativesContextAuditResult, sourceAudit FuturesDerivativesContextSourceAuditResult, cfg FuturesDerivativesContextAuditConfig) []FuturesDerivativesContextAuditCoverageRow {
	rows := []FuturesDerivativesContextAuditCoverageRow{}
	for _, symbol := range derivativesContextSymbols() {
		for _, timeframe := range cfg.StateConfig.withDefaults().Timeframes {
			states, features, labels, missingBasis, missingChange, missingPremium := 0, 0, 0, 0, 0, 0
			for _, row := range result.LocalStateRows {
				if row.Symbol == symbol && row.Timeframe == timeframe {
					states++
					if !row.DerivativesContextAvailable {
						missingBasis++
					}
				}
			}
			for _, row := range result.BasisFeatureRows {
				if row.Symbol == symbol && row.Timeframe == timeframe {
					features++
					if !row.BasisChangePresent {
						missingChange++
					}
					if !row.PremiumPresent {
						missingPremium++
					}
				}
			}
			for _, row := range result.LabelRows {
				if row.Symbol == symbol && row.Timeframe == timeframe {
					labels++
				}
			}
			coverage := safeDiv(float64(features), float64(states))
			lagCoverage := derivativesContextMinRequiredLagCoverage(sourceAudit, symbol)
			rows = append(rows, FuturesDerivativesContextAuditCoverageRow{
				AuditName:               FuturesDerivativesContextAuditName,
				Scope:                   "lagged_basis_context",
				Symbol:                  symbol,
				Timeframe:               timeframe,
				SourceFamily:            "mark_minus_index_basis",
				Required:                true,
				StateRows:               states,
				LabelRows:               labels,
				BasisFeatureRows:        features,
				MissingBasisRows:        missingBasis,
				MissingBasisChangeRows:  missingChange,
				MissingPremiumRows:      missingPremium,
				LagCoveragePct:          lagCoverage,
				BasisContextCoveragePct: coverage,
				RequiredCoverageFloor:   cfg.MinBasisContextCoverage,
				CoverageGatePass:        lagCoverage+1e-6 >= cfg.MinBasisContextCoverage,
				ClosedCandleOnly:        true,
				ForwardLabelsAsInputs:   false,
				ValidationStatus:        "accepted",
			})
		}
	}
	return rows
}

func derivativesContextFeatureMissingnessRows(result FuturesDerivativesContextAuditResult, missing map[string]map[string]int, cfg FuturesDerivativesContextAuditConfig) []FuturesDerivativesContextAuditMissingnessRow {
	rows := []FuturesDerivativesContextAuditMissingnessRow{}
	coveragePassBySymbolTimeframe := map[string]bool{}
	for _, row := range result.CoverageRows {
		if row.Scope == "lagged_basis_context" {
			coveragePassBySymbolTimeframe[row.Symbol+"|"+row.Timeframe] = row.CoverageGatePass
		}
	}
	for _, symbol := range derivativesContextSymbols() {
		for _, timeframe := range cfg.StateConfig.withDefaults().Timeframes {
			total := 0
			missingBasis := missing[symbol][timeframe]
			missingChange := 0
			missingPremium := 0
			for _, feature := range result.BasisFeatureRows {
				if feature.Symbol != symbol || feature.Timeframe != timeframe {
					continue
				}
				total++
				if !feature.BasisChangePresent {
					missingChange++
				}
				if !feature.PremiumPresent {
					missingPremium++
				}
			}
			stateTotal := total + missingBasis
			for _, item := range []struct {
				reason string
				count  int
				denom  int
			}{
				{"missing_required_lagged_basis_context", missingBasis, stateTotal},
				{"missing_basis_change_lookback", missingChange, total},
				{"missing_optional_premium_context", missingPremium, total},
			} {
				coveragePass := true
				if item.reason == "missing_required_lagged_basis_context" {
					coveragePass = coveragePassBySymbolTimeframe[symbol+"|"+timeframe]
				}
				rows = append(rows, FuturesDerivativesContextAuditMissingnessRow{
					AuditName:         FuturesDerivativesContextAuditName,
					Scope:             "derived_context",
					Symbol:            symbol,
					Timeframe:         timeframe,
					SourceFamily:      "mark_minus_index_basis",
					Reason:            item.reason,
					Count:             item.count,
					TotalRows:         item.denom,
					Rate:              safeDiv(float64(item.count), float64(item.denom)),
					ForwardFilledRows: 0,
					MissingDataPolicy: derivativesMissingPolicy,
					CoverageFloor:     cfg.MinBasisContextCoverage,
					CoverageGatePass:  coveragePass,
				})
			}
		}
	}
	return rows
}

func derivativesContextCoveragePass(rows []FuturesDerivativesContextAuditCoverageRow) bool {
	if len(rows) == 0 {
		return false
	}
	hasBasis := false
	for _, row := range rows {
		if row.ForwardLabelsAsInputs || !row.CoverageGatePass {
			return false
		}
		if row.Scope == "lagged_basis_context" {
			hasBasis = true
		}
	}
	return hasBasis
}

func derivativesContextMinCoverage(rows []FuturesDerivativesContextAuditCoverageRow) float64 {
	minCov := math.Inf(1)
	for _, row := range rows {
		if row.Scope != "lagged_basis_context" {
			continue
		}
		if row.BasisContextCoveragePct < minCov {
			minCov = row.BasisContextCoveragePct
		}
	}
	if math.IsInf(minCov, 1) {
		return 0
	}
	return minCov
}

func derivativesContextMinRequiredLagCoverage(sourceAudit FuturesDerivativesContextSourceAuditResult, symbol string) float64 {
	minCov := math.Inf(1)
	for _, row := range sourceAudit.TimestampAlignmentRows {
		if row.Symbol == symbol && row.Required && row.LagCoveragePct < minCov {
			minCov = row.LagCoveragePct
		}
	}
	if math.IsInf(minCov, 1) {
		return 0
	}
	return minCov
}

func sourceAuditLagAlignedRows(sourceAudit FuturesDerivativesContextSourceAuditResult, symbol, family string) int {
	for _, row := range sourceAudit.TimestampAlignmentRows {
		if row.Symbol == symbol && row.SourceFamily == family {
			return row.LagAlignedRows
		}
	}
	return 0
}

func derivativesContextBasisVolatility(mark, index *derivStreamData, sourceOpenMs int64, lookback int) (float64, string) {
	if lookback <= 1 {
		return 0, "basis_vol_unknown"
	}
	values := []float64{}
	for i := lookback - 1; i >= 0; i-- {
		open := sourceOpenMs - int64(i)*derivativesIntervalMs
		m, mok := mark.openClose[open]
		idx, iok := index.openClose[open]
		if !mok || !iok || idx == 0 {
			continue
		}
		values = append(values, (m-idx)/idx*10000)
	}
	if len(values) < maxInt(3, lookback/2) {
		return 0, "basis_vol_missing"
	}
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))
	var sum float64
	for _, v := range values {
		d := v - mean
		sum += d * d
	}
	vol := math.Sqrt(sum / float64(len(values)))
	return vol, derivativesContextBasisVolBucket(vol)
}

func derivativesContextLaggedSourceOpen(decisionCloseMs int64, lagIntervals int) int64 {
	return decisionCloseMs - int64(lagIntervals)*derivativesIntervalMs - (derivativesIntervalMs - 1000)
}

func derivativesContextStreamKey(symbol, family string) string {
	return strings.ToUpper(strings.TrimSpace(symbol)) + "|" + family
}

func derivativesContextBasisSignBucket(bps float64) string {
	switch {
	case bps > 0.25:
		return "perp_premium"
	case bps < -0.25:
		return "perp_discount"
	default:
		return "basis_flat"
	}
}

func derivativesContextBasisLevelBucket(bps float64) string {
	switch {
	case bps <= -5:
		return "basis_discount_wide"
	case bps < -1:
		return "basis_discount_small"
	case bps <= 1:
		return "basis_flat"
	case bps < 5:
		return "basis_premium_small"
	default:
		return "basis_premium_wide"
	}
}

func derivativesContextBasisChangeBucket(change float64, present bool) string {
	if !present {
		return "basis_change_missing"
	}
	switch {
	case change <= -3:
		return "basis_contracting_fast"
	case change < -0.75:
		return "basis_contracting"
	case change <= 0.75:
		return "basis_change_flat"
	case change < 3:
		return "basis_expanding"
	default:
		return "basis_expanding_fast"
	}
}

func derivativesContextBasisVolBucket(vol float64) string {
	switch {
	case !validNumber(vol) || vol == 0:
		return "basis_vol_unknown"
	case vol < 0.75:
		return "basis_vol_low"
	case vol < 2:
		return "basis_vol_mid"
	default:
		return "basis_vol_high"
	}
}

func derivativesContextPremiumBucket(bps float64, present bool) string {
	if !present {
		return "premium_missing"
	}
	switch {
	case bps <= -5:
		return "premium_discount_wide"
	case bps < -1:
		return "premium_discount_small"
	case bps <= 1:
		return "premium_flat"
	case bps < 5:
		return "premium_positive_small"
	default:
		return "premium_positive_wide"
	}
}

func derivativesContextRouteUseful(route string, label FuturesDerivativesContextAuditLabelRow) bool {
	switch route {
	case RangeStateConstructionLoopRouteRotation:
		return label.RotationUseful
	case RangeStateConstructionLoopRouteContinuation:
		return label.ContinuationUseful
	default:
		return false
	}
}

func derivativesContextRouteToxic(route string, label FuturesDerivativesContextAuditLabelRow) bool {
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

func derivativesContextRouteRateGate(row FuturesDerivativesContextAuditCohortRow, cfg FuturesDerivativesContextAuditConfig) bool {
	switch row.RouteCandidate {
	case RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation:
		return row.FullUsefulRate >= cfg.MinUsefulRateFull &&
			row.WeakestSplitUsefulRate >= cfg.MinUsefulRateSplit &&
			row.FullToxicRate <= cfg.MaxToxicRateFull &&
			row.WorstSplitToxicRate <= cfg.MaxToxicRateSplit &&
			row.FullUsefulMinusToxicMargin >= cfg.MinUsefulMinusToxicMarginFull &&
			row.WeakestSplitMargin >= cfg.MinUsefulMinusToxicMarginSplit
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return row.FullToxicRate >= cfg.MinToxicRateFull &&
			row.WorstSplitToxicRate >= cfg.MinToxicRateSplit
	default:
		return false
	}
}

func derivativesContextImprovementGate(row FuturesDerivativesContextAuditCohortRow, cfg FuturesDerivativesContextAuditConfig) bool {
	if row.BaselineFullRows <= 0 {
		return false
	}
	switch row.RouteCandidate {
	case RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation:
		return row.FullUsefulImprovement >= cfg.MinContextUsefulImprovementFull &&
			row.WeakestSplitUsefulImprovement >= cfg.MinContextUsefulImprovementSplit &&
			row.FullMarginImprovement >= cfg.MinContextMarginImprovementFull &&
			row.WeakestSplitMarginImprovement >= cfg.MinContextMarginImprovementSplit
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return row.FullToxicImprovement >= cfg.MinContextToxicImprovementFull &&
			row.WeakestSplitToxicImprovement >= cfg.MinContextToxicImprovementSplit
	default:
		return false
	}
}

func derivativesContextRankScore(row FuturesDerivativesContextAuditRankingRow) float64 {
	switch row.RouteCandidate {
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return row.FullToxicRate + row.WorstSplitToxicRate + row.FullToxicImprovement + row.WeakestSplitToxicImprovement - row.BucketShareOfLocalBaselineFull
	default:
		return row.FullUsefulMinusToxicMargin + row.WeakestSplitMargin + row.FullMarginImprovement + row.WeakestSplitMarginImprovement - row.BucketShareOfLocalBaselineFull
	}
}

func derivativesContextPassingRankingCount(rows []FuturesDerivativesContextAuditRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func derivativesContextBaselineKey(row FuturesDerivativesContextAuditCohortRow) string {
	return strings.Join([]string{row.Symbol, row.Timeframe, fmt.Sprintf("%d", row.HorizonBars), row.RouteCandidate, row.LocalRangeBucketID}, "|")
}

func derivativesContextRequiredSymbolsPresent(data map[string]derivativesContextSymbolData) bool {
	for _, symbol := range derivativesContextSymbols() {
		if _, ok := data[symbol]; !ok {
			return false
		}
	}
	return true
}

func derivativesContextSymbols() []string {
	return []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}
}

func derivativesContextLessCohort(a, b FuturesDerivativesContextAuditCohortRow) bool {
	if a.Symbol != b.Symbol {
		return rangeUniverseSymbolSortKey(a.Symbol) < rangeUniverseSymbolSortKey(b.Symbol)
	}
	if a.Timeframe != b.Timeframe {
		return rangeContextTriageTimeframeSortKey(a.Timeframe) < rangeContextTriageTimeframeSortKey(b.Timeframe)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.RouteCandidate != b.RouteCandidate {
		return rangeStateRouteSortKey(a.RouteCandidate) < rangeStateRouteSortKey(b.RouteCandidate)
	}
	if a.ComparisonType != b.ComparisonType {
		return a.ComparisonType < b.ComparisonType
	}
	if a.DerivativesBucketType != b.DerivativesBucketType {
		return a.DerivativesBucketType < b.DerivativesBucketType
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	return a.CohortID < b.CohortID
}

func safeDiv(n, d float64) float64 {
	if d == 0 {
		return 0
	}
	return n / d
}
