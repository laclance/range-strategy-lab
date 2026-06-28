package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesBTCRegimeETHSOLContextAuditName = "futures_btc_regime_eth_sol_context_audit"

	BTCRegimeETHSOLContextStopStateSourceGap             = "btc_regime_eth_sol_context_zero_trade_audit_source_gap"
	BTCRegimeETHSOLContextStopStateRejectedClosedReslice = "btc_regime_eth_sol_context_zero_trade_audit_rejected_closed_family_reslice"
	BTCRegimeETHSOLContextStopStateRejectedFutureLeak    = "btc_regime_eth_sol_context_zero_trade_audit_rejected_future_label_leak"
	BTCRegimeETHSOLContextStopStateFailedNoUsableContext = "btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context"
	BTCRegimeETHSOLContextStopStatePassedNeedsSpec       = "btc_regime_eth_sol_context_zero_trade_audit_passed_needs_strategy_premise_spec"

	BTCRegimeETHSOLContextRoleBTCMarketRegime = "btc_market_regime_context_diagnostic_only"
	BTCRegimeETHSOLContextRoleAuthorityLocal  = "eth_sol_local_range_context_authority_candidate"

	BTCRegimeETHSOLContextComparisonBTCPlusLocal = "btc_regime_plus_local"
	BTCRegimeETHSOLContextComparisonLocalOnly    = "local_only_baseline"

	btcRegimeETHSOLContextExpectedRows        = 573984
	btcRegimeETHSOLContextExpectedFirst       = "2021-01-01T00:00:00Z"
	btcRegimeETHSOLContextExpectedLast        = "2026-06-16T23:55:00Z"
	btcRegimeETHSOLContextExpectedBTCZeroVol  = 66
	btcRegimeETHSOLContextExpectedETHZeroVol  = 47
	btcRegimeETHSOLContextExpectedSOLZeroVol  = 47
	btcRegimeETHSOLContextExpectedSOLPhysical = 1
)

type FuturesBTCRegimeETHSOLContextSourceConfig struct {
	Symbol                       string
	Path                         string
	ApprovedPath                 string
	ExpectedRowCount             int
	ExpectedFirstOpenTime        string
	ExpectedLastOpenTime         string
	ExpectedGapCount             int
	ExpectedDuplicateCount       int
	ExpectedZeroVolumeCount      int
	ExpectedPhysicalNonMonotonic int
	ExpectedSortedForValidation  bool
	SkipSourceFactCheck          bool
	SkipSplitEligibilityCheck    bool
}

type FuturesBTCRegimeETHSOLContextAuditConfig struct {
	Sources                          []FuturesBTCRegimeETHSOLContextSourceConfig
	StateConfig                      FuturesRangeStateConstructionLoopAuditConfig
	MinFullCohortRows                int
	MinSplitCohortRows               int
	MaxSplitContributionRate         float64
	MinUsefulRateFull                float64
	MinUsefulRateSplit               float64
	MaxToxicRateFull                 float64
	MaxToxicRateSplit                float64
	MinUsefulMinusToxicMarginFull    float64
	MinUsefulMinusToxicMarginSplit   float64
	MinToxicRateFull                 float64
	MinToxicRateSplit                float64
	MinContextUsefulImprovementFull  float64
	MinContextUsefulImprovementSplit float64
	MinContextMarginImprovementFull  float64
	MinContextMarginImprovementSplit float64
	MinContextToxicImprovementFull   float64
	MinContextToxicImprovementSplit  float64
}

type FuturesBTCRegimeETHSOLContextAuditResult struct {
	SourceRows           []FuturesBTCRegimeETHSOLContextSourceRow           `json:"source_rows"`
	CoverageRows         []FuturesBTCRegimeETHSOLContextCoverageRow         `json:"coverage_rows"`
	BTCStateRows         []FuturesBTCRegimeETHSOLContextBTCStateRow         `json:"btc_state_rows"`
	LocalStateRows       []FuturesBTCRegimeETHSOLContextLocalStateRow       `json:"local_state_rows"`
	RelativeStrengthRows []FuturesBTCRegimeETHSOLContextRelativeStrengthRow `json:"relative_strength_rows"`
	LabelRows            []FuturesBTCRegimeETHSOLContextLabelRow            `json:"label_rows"`
	CohortRows           []FuturesBTCRegimeETHSOLContextCohortRow           `json:"cohort_rows"`
	RankingRows          []FuturesBTCRegimeETHSOLContextRankingRow          `json:"ranking_rows"`
	SummaryRows          []FuturesBTCRegimeETHSOLContextSummaryRow          `json:"summary_rows"`
	PassingCohorts       int                                                `json:"passing_cohorts"`
	StopState            string                                             `json:"stop_state"`
}

type FuturesBTCRegimeETHSOLContextSourceRow struct {
	AuditName                    string `json:"audit_name"`
	Role                         string `json:"role"`
	ApprovedLocalScopeOnly       bool   `json:"approved_local_scope_only"`
	ExpectedRowCount             int    `json:"expected_row_count"`
	ExpectedFirstOpenTime        string `json:"expected_first_open_time"`
	ExpectedLastOpenTime         string `json:"expected_last_open_time"`
	ExpectedGapCount             int    `json:"expected_gap_count"`
	ExpectedDuplicateCount       int    `json:"expected_duplicate_count"`
	ExpectedZeroVolumeCount      int    `json:"expected_zero_volume_count"`
	ExpectedPhysicalNonMonotonic int    `json:"expected_physical_non_monotonic_count"`
	ExpectedSortedForValidation  bool   `json:"expected_sorted_for_validation"`
	SourceFactsPass              bool   `json:"source_facts_pass"`
	ForwardLabelsAsSourceInput   bool   `json:"forward_labels_as_source_input"`
	FuturesRangeUniverseSourceRow
}

type FuturesBTCRegimeETHSOLContextCoverageRow struct {
	AuditName             string `json:"audit_name"`
	Symbol                string `json:"symbol"`
	Role                  string `json:"role"`
	ExpectedRowCount      int    `json:"expected_row_count"`
	ExpectedFirstOpenTime string `json:"expected_first_open_time"`
	ExpectedLastOpenTime  string `json:"expected_last_open_time"`
	CoverageFactsPass     bool   `json:"coverage_facts_pass"`
	StateRows             int    `json:"state_rows"`
	LabelRows             int    `json:"label_rows"`
	ContextMatchedRows    int    `json:"context_matched_rows"`
	MissingBTCContextRows int    `json:"missing_btc_context_rows"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesBTCRegimeETHSOLContextBTCStateRow struct {
	BTCStateRowID             int     `json:"btc_state_row_id"`
	Timestamp                 string  `json:"timestamp"`
	Timeframe                 string  `json:"timeframe"`
	Split                     string  `json:"split"`
	DecisionIndex             int     `json:"decision_index"`
	DecisionCloseTime         string  `json:"decision_close_time"`
	RangeEpisodeID            int     `json:"range_episode_id"`
	RangeAgeBars              int     `json:"range_age_bars"`
	CloseLocationPct          float64 `json:"close_location_pct"`
	ReturnShortPct            float64 `json:"return_short_pct"`
	ReturnMediumPct           float64 `json:"return_medium_pct"`
	RealizedVolPercentileLong float64 `json:"realized_vol_percentile_long"`
	VolumePercentileLong      float64 `json:"volume_percentile_long"`
	StateID                   string  `json:"state_id"`
	BTCRegimeID               string  `json:"btc_regime_id"`
	GeometryBucket            string  `json:"geometry_bucket"`
	VolBucket                 string  `json:"vol_bucket"`
	TrendBucket               string  `json:"trend_bucket"`
	ImpulseBucket             string  `json:"impulse_bucket"`
	ParticipationBucket       string  `json:"participation_bucket"`
	CloseLocationBucket       string  `json:"close_location_bucket"`
	RangeAgeBucket            string  `json:"range_age_bucket"`
	ClosedCandleOnly          bool    `json:"closed_candle_only"`
	DiagnosticOnlyAuthority   bool    `json:"diagnostic_only_authority"`
	ForwardLabelsAsStateInput bool    `json:"forward_labels_as_state_input"`
}

type FuturesBTCRegimeETHSOLContextLocalStateRow struct {
	LocalStateRowID           int     `json:"local_state_row_id"`
	Symbol                    string  `json:"symbol"`
	Timestamp                 string  `json:"timestamp"`
	Timeframe                 string  `json:"timeframe"`
	Split                     string  `json:"split"`
	DecisionIndex             int     `json:"decision_index"`
	DecisionCloseTime         string  `json:"decision_close_time"`
	RangeEpisodeID            int     `json:"range_episode_id"`
	RangeAgeBars              int     `json:"range_age_bars"`
	CloseLocationPct          float64 `json:"close_location_pct"`
	ReturnShortPct            float64 `json:"return_short_pct"`
	ReturnMediumPct           float64 `json:"return_medium_pct"`
	StateID                   string  `json:"state_id"`
	LocalRangeBucketID        string  `json:"local_range_bucket_id"`
	BTCStateRowID             int     `json:"btc_state_row_id"`
	BTCRegimeID               string  `json:"btc_regime_id"`
	RelativeStrengthBucket    string  `json:"relative_strength_bucket"`
	ContextStateID            string  `json:"context_state_id"`
	GeometryBucket            string  `json:"geometry_bucket"`
	VolBucket                 string  `json:"vol_bucket"`
	TrendBucket               string  `json:"trend_bucket"`
	ImpulseBucket             string  `json:"impulse_bucket"`
	ParticipationBucket       string  `json:"participation_bucket"`
	ClosedCandleOnly          bool    `json:"closed_candle_only"`
	AuthorityCandidateOnly    bool    `json:"authority_candidate_only"`
	ForwardLabelsAsStateInput bool    `json:"forward_labels_as_state_input"`
	ClosedFamilyReslice       bool    `json:"closed_family_reslice"`
}

type FuturesBTCRegimeETHSOLContextRelativeStrengthRow struct {
	RelativeStrengthRowID   int     `json:"relative_strength_row_id"`
	Symbol                  string  `json:"symbol"`
	Timestamp               string  `json:"timestamp"`
	Timeframe               string  `json:"timeframe"`
	DecisionIndex           int     `json:"decision_index"`
	LocalStateRowID         int     `json:"local_state_row_id"`
	BTCStateRowID           int     `json:"btc_state_row_id"`
	LocalReturnShortPct     float64 `json:"local_return_short_pct"`
	BTCReturnShortPct       float64 `json:"btc_return_short_pct"`
	RelativeReturnShortPct  float64 `json:"relative_return_short_pct"`
	LocalReturnMediumPct    float64 `json:"local_return_medium_pct"`
	BTCReturnMediumPct      float64 `json:"btc_return_medium_pct"`
	RelativeReturnMediumPct float64 `json:"relative_return_medium_pct"`
	RelativeStrengthBucket  string  `json:"relative_strength_bucket"`
	ClosedCandleOnly        bool    `json:"closed_candle_only"`
	ForwardLabelsAsRSInput  bool    `json:"forward_labels_as_relative_strength_input"`
}

type FuturesBTCRegimeETHSOLContextLabelRow struct {
	LabelRowID                int     `json:"label_row_id"`
	Symbol                    string  `json:"symbol"`
	LocalStateRowID           int     `json:"local_state_row_id"`
	BTCStateRowID             int     `json:"btc_state_row_id"`
	Timestamp                 string  `json:"timestamp"`
	Timeframe                 string  `json:"timeframe"`
	Split                     string  `json:"split"`
	HorizonBars               int     `json:"horizon_bars"`
	LocalRangeBucketID        string  `json:"local_range_bucket_id"`
	BTCRegimeID               string  `json:"btc_regime_id"`
	RelativeStrengthBucket    string  `json:"relative_strength_bucket"`
	ContextStateID            string  `json:"context_state_id"`
	ForwardLabel              string  `json:"forward_label"`
	RotationUseful            bool    `json:"rotation_useful"`
	RotationToxic             bool    `json:"rotation_toxic"`
	ContinuationUseful        bool    `json:"continuation_useful"`
	ContinuationToxic         bool    `json:"continuation_toxic"`
	NoTradeToxic              bool    `json:"no_trade_toxic"`
	DiagnosticOnly            bool    `json:"diagnostic_only"`
	LabelWindowStartIndex     int     `json:"label_window_start_index"`
	LabelWindowEndIndex       int     `json:"label_window_end_index"`
	LabelWindowStartTime      string  `json:"label_window_start_time"`
	LabelWindowEndTime        string  `json:"label_window_end_time"`
	InsideCloseRate           float64 `json:"inside_close_rate"`
	MaxExcursionAboveWidths   float64 `json:"max_excursion_above_range_widths"`
	MaxExcursionBelowWidths   float64 `json:"max_excursion_below_range_widths"`
	ForwardLabelMetadataOnly  bool    `json:"forward_label_metadata_only"`
	ForwardLabelUsedAsFeature bool    `json:"forward_label_used_as_feature"`
}

type FuturesBTCRegimeETHSOLContextCohortRow struct {
	CohortID                       string  `json:"cohort_id"`
	ComparisonType                 string  `json:"comparison_type"`
	Symbol                         string  `json:"symbol"`
	Split                          string  `json:"split"`
	Timeframe                      string  `json:"timeframe"`
	HorizonBars                    int     `json:"horizon_bars"`
	RouteCandidate                 string  `json:"route_candidate"`
	LocalRangeBucketID             string  `json:"local_range_bucket_id"`
	BTCRegimeID                    string  `json:"btc_regime_id,omitempty"`
	RelativeStrengthBucket         string  `json:"relative_strength_bucket,omitempty"`
	CandidateCount                 int     `json:"candidate_count"`
	UsefulCount                    int     `json:"useful_count"`
	ToxicCount                     int     `json:"toxic_count"`
	RotationUsefulCount            int     `json:"rotation_useful_count"`
	RotationToxicCount             int     `json:"rotation_toxic_count"`
	ContinuationUsefulCount        int     `json:"continuation_useful_count"`
	ContinuationToxicCount         int     `json:"continuation_toxic_count"`
	NoTradeToxicCount              int     `json:"no_trade_toxic_count"`
	UsefulRate                     float64 `json:"useful_rate"`
	ToxicRate                      float64 `json:"toxic_rate"`
	UsefulMinusToxicMargin         float64 `json:"useful_minus_toxic_margin"`
	DominantForwardLabel           string  `json:"dominant_forward_label"`
	DominantForwardLabelRate       float64 `json:"dominant_forward_label_rate"`
	FullPeriodRows                 int     `json:"full_period_rows"`
	WeakestSplitRows               int     `json:"weakest_split_rows"`
	MaxSplitContributionRate       float64 `json:"max_split_contribution_rate"`
	FullUsefulRate                 float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate         float64 `json:"weakest_split_useful_rate"`
	FullToxicRate                  float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate            float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin     float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin             float64 `json:"weakest_split_margin"`
	BaselineFullRows               int     `json:"baseline_full_rows"`
	BaselineFullUsefulRate         float64 `json:"baseline_full_useful_rate"`
	BaselineWeakestSplitUsefulRate float64 `json:"baseline_weakest_split_useful_rate"`
	BaselineFullToxicRate          float64 `json:"baseline_full_toxic_rate"`
	BaselineWorstSplitToxicRate    float64 `json:"baseline_worst_split_toxic_rate"`
	BaselineFullMargin             float64 `json:"baseline_full_margin"`
	BaselineWeakestSplitMargin     float64 `json:"baseline_weakest_split_margin"`
	FullUsefulImprovement          float64 `json:"full_useful_improvement"`
	WeakestSplitUsefulImprovement  float64 `json:"weakest_split_useful_improvement"`
	FullToxicImprovement           float64 `json:"full_toxic_improvement"`
	WeakestSplitToxicImprovement   float64 `json:"weakest_split_toxic_improvement"`
	FullMarginImprovement          float64 `json:"full_margin_improvement"`
	WeakestSplitMarginImprovement  float64 `json:"weakest_split_margin_improvement"`
	ReviewableCountGatePass        bool    `json:"reviewable_count_gate_pass"`
	SplitStabilityGatePass         bool    `json:"split_stability_gate_pass"`
	SplitContributionGatePass      bool    `json:"split_contribution_gate_pass"`
	BTCContextImprovementGatePass  bool    `json:"btc_context_improvement_gate_pass"`
	RouteRateGatePass              bool    `json:"route_rate_gate_pass"`
	ClosedFamilyProtectionPass     bool    `json:"closed_family_protection_pass"`
	FutureLeakProtectionPass       bool    `json:"future_leak_protection_pass"`
	PassesReviewGate               bool    `json:"passes_review_gate"`
	FailureReason                  string  `json:"failure_reason,omitempty"`
}

type FuturesBTCRegimeETHSOLContextRankingRow struct {
	Rank                          int     `json:"rank"`
	CohortID                      string  `json:"cohort_id"`
	Symbol                        string  `json:"symbol"`
	Timeframe                     string  `json:"timeframe"`
	HorizonBars                   int     `json:"horizon_bars"`
	RouteCandidate                string  `json:"route_candidate"`
	LocalRangeBucketID            string  `json:"local_range_bucket_id"`
	BTCRegimeID                   string  `json:"btc_regime_id"`
	RelativeStrengthBucket        string  `json:"relative_strength_bucket"`
	PassesGate                    bool    `json:"passes_gate"`
	RankScore                     float64 `json:"rank_score"`
	FullPeriodRows                int     `json:"full_period_rows"`
	WeakestSplitRows              int     `json:"weakest_split_rows"`
	MaxSplitContributionRate      float64 `json:"max_split_contribution_rate"`
	FullUsefulRate                float64 `json:"full_useful_rate"`
	WeakestSplitUsefulRate        float64 `json:"weakest_split_useful_rate"`
	FullToxicRate                 float64 `json:"full_toxic_rate"`
	WorstSplitToxicRate           float64 `json:"worst_split_toxic_rate"`
	FullUsefulMinusToxicMargin    float64 `json:"full_useful_minus_toxic_margin"`
	WeakestSplitMargin            float64 `json:"weakest_split_margin"`
	FullUsefulImprovement         float64 `json:"full_useful_improvement"`
	WeakestSplitUsefulImprovement float64 `json:"weakest_split_useful_improvement"`
	FullToxicImprovement          float64 `json:"full_toxic_improvement"`
	WeakestSplitToxicImprovement  float64 `json:"weakest_split_toxic_improvement"`
	FullMarginImprovement         float64 `json:"full_margin_improvement"`
	WeakestSplitMarginImprovement float64 `json:"weakest_split_margin_improvement"`
	DominantForwardLabel          string  `json:"dominant_forward_label"`
	DominantForwardLabelRate      float64 `json:"dominant_forward_label_rate"`
	FutureLeakProtectionPass      bool    `json:"future_leak_protection_pass"`
	ClosedFamilyProtectionPass    bool    `json:"closed_family_protection_pass"`
	FailureReason                 string  `json:"failure_reason,omitempty"`
}

type FuturesBTCRegimeETHSOLContextSummaryRow struct {
	Split                        string `json:"split"`
	Symbol                       string `json:"symbol"`
	Timeframe                    string `json:"timeframe"`
	HorizonBars                  int    `json:"horizon_bars"`
	SourceRows                   int    `json:"source_rows"`
	CoverageRows                 int    `json:"coverage_rows"`
	BTCStateRows                 int    `json:"btc_state_rows"`
	LocalStateRows               int    `json:"local_state_rows"`
	RelativeStrengthRows         int    `json:"relative_strength_rows"`
	LabelRows                    int    `json:"label_rows"`
	CohortRows                   int    `json:"cohort_rows"`
	RankingRows                  int    `json:"ranking_rows"`
	PassingCohorts               int    `json:"passing_cohorts"`
	SourceScopePass              bool   `json:"source_scope_pass"`
	CoveragePass                 bool   `json:"coverage_pass"`
	CommonOutputsZeroTrade       bool   `json:"common_outputs_zero_trade"`
	BTCContextDiagnosticOnly     bool   `json:"btc_context_diagnostic_only"`
	ETHSOLAuthorityCandidateOnly bool   `json:"eth_sol_authority_candidate_only"`
	ForwardLabelsAsInputs        bool   `json:"forward_labels_as_inputs"`
	StopState                    string `json:"stop_state"`
}

type btcRegimeContextSymbolData struct {
	symbol     string
	candles    []Candle
	source     FuturesBTCRegimeETHSOLContextSourceRow
	coverage   []FuturesBTCRegimeETHSOLContextCoverageRow
	frames     map[string]rangeStateFrameData
	states     []FuturesRangeStateConstructionLoopStateRow
	labels     []FuturesRangeStateConstructionLoopLabelRow
	skips      rangeStateSkipAccumulator
	sourceOK   bool
	coverageOK bool
}

type btcRegimeContextStateKey struct {
	symbol    string
	timeframe string
	timestamp string
}

type btcRegimeContextCohortKey struct {
	cohortID       string
	comparisonType string
	symbol         string
	split          string
	timeframe      string
	horizonBars    int
	routeCandidate string
	localBucketID  string
	btcRegimeID    string
	rsBucket       string
}

type btcRegimeContextCohortAccumulator struct {
	row    FuturesBTCRegimeETHSOLContextCohortRow
	labels map[string]int
}

func DefaultFuturesBTCRegimeETHSOLContextAuditConfig() FuturesBTCRegimeETHSOLContextAuditConfig {
	return FuturesBTCRegimeETHSOLContextAuditConfig{
		Sources: []FuturesBTCRegimeETHSOLContextSourceConfig{
			{
				Symbol:                  RangeUniverseSymbolBTCUSDT,
				Path:                    "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
				ApprovedPath:            "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
				ExpectedZeroVolumeCount: btcRegimeETHSOLContextExpectedBTCZeroVol,
			},
			{
				Symbol:                  RangeUniverseSymbolETHUSDT,
				Path:                    "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv",
				ApprovedPath:            "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv",
				ExpectedZeroVolumeCount: btcRegimeETHSOLContextExpectedETHZeroVol,
			},
			{
				Symbol:                       RangeUniverseSymbolSOLUSDT,
				Path:                         "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv",
				ApprovedPath:                 "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv",
				ExpectedZeroVolumeCount:      btcRegimeETHSOLContextExpectedSOLZeroVol,
				ExpectedPhysicalNonMonotonic: btcRegimeETHSOLContextExpectedSOLPhysical,
				ExpectedSortedForValidation:  true,
			},
		},
		StateConfig:                      DefaultFuturesRangeStateConstructionLoopAuditConfig(),
		MinFullCohortRows:                300,
		MinSplitCohortRows:               60,
		MaxSplitContributionRate:         0.60,
		MinUsefulRateFull:                0.56,
		MinUsefulRateSplit:               0.50,
		MaxToxicRateFull:                 0.44,
		MaxToxicRateSplit:                0.50,
		MinUsefulMinusToxicMarginFull:    0.10,
		MinUsefulMinusToxicMarginSplit:   0.03,
		MinToxicRateFull:                 0.58,
		MinToxicRateSplit:                0.52,
		MinContextUsefulImprovementFull:  0.04,
		MinContextUsefulImprovementSplit: 0.01,
		MinContextMarginImprovementFull:  0.04,
		MinContextMarginImprovementSplit: 0.01,
		MinContextToxicImprovementFull:   0.04,
		MinContextToxicImprovementSplit:  0.01,
	}
}

func RunFuturesBTCRegimeETHSOLContextAudit(cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split) (FuturesBTCRegimeETHSOLContextAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesBTCRegimeETHSOLContextAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	result := FuturesBTCRegimeETHSOLContextAuditResult{}
	dataBySymbol := map[string]btcRegimeContextSymbolData{}
	sourceOK := true
	for _, source := range cfg.Sources {
		data := btcRegimeETHSOLContextLoadSymbol(source, cfg, splits)
		result.SourceRows = append(result.SourceRows, data.source)
		dataBySymbol[data.symbol] = data
		if !data.sourceOK {
			sourceOK = false
		}
	}
	if !sourceOK || !btcRegimeETHSOLContextRequiredSourcesPresent(dataBySymbol) {
		result.StopState = BTCRegimeETHSOLContextStopStateSourceGap
		result.SummaryRows = btcRegimeETHSOLContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	coverageOK := true
	for symbol := range dataBySymbol {
		data := dataBySymbol[symbol]
		built, err := btcRegimeETHSOLContextBuildSymbolStates(data, cfg, splits)
		if err != nil {
			return result, err
		}
		dataBySymbol[symbol] = built
		result.CoverageRows = append(result.CoverageRows, built.coverage...)
		if !built.coverageOK {
			coverageOK = false
		}
	}
	sort.Slice(result.CoverageRows, func(i, j int) bool {
		if result.CoverageRows[i].Symbol != result.CoverageRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.CoverageRows[i].Symbol) < rangeUniverseSymbolSortKey(result.CoverageRows[j].Symbol)
		}
		return rangeContextTriageTimeframeSortKey(result.CoverageRows[i].Timeframe) < rangeContextTriageTimeframeSortKey(result.CoverageRows[j].Timeframe)
	})
	if !coverageOK {
		result.StopState = BTCRegimeETHSOLContextStopStateSourceGap
		result.SummaryRows = btcRegimeETHSOLContextSummaryRows(result, cfg, splits)
		return result, nil
	}

	btcData := dataBySymbol[RangeUniverseSymbolBTCUSDT]
	btcRows := btcRegimeETHSOLContextBTCStateRows(btcData.states)
	result.BTCStateRows = btcRows
	btcByKey := map[btcRegimeContextStateKey]FuturesBTCRegimeETHSOLContextBTCStateRow{}
	for _, row := range btcRows {
		btcByKey[btcRegimeContextStateKey{symbol: RangeUniverseSymbolBTCUSDT, timeframe: row.Timeframe, timestamp: row.DecisionCloseTime}] = row
	}

	localStateByKey := map[btcRegimeContextStateKey]FuturesBTCRegimeETHSOLContextLocalStateRow{}
	localStateByID := map[string]FuturesBTCRegimeETHSOLContextLocalStateRow{}
	rsID := 0
	for _, symbol := range []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		data := dataBySymbol[symbol]
		for _, state := range data.states {
			key := btcRegimeContextStateKey{symbol: RangeUniverseSymbolBTCUSDT, timeframe: state.Timeframe, timestamp: state.DecisionCloseTime}
			btc, ok := btcByKey[key]
			if !ok {
				continue
			}
			rsID++
			rs := btcRegimeETHSOLContextRelativeStrengthRow(rsID, symbol, state, btc)
			local := btcRegimeETHSOLContextLocalStateRow(symbol, state, btc, rs.RelativeStrengthBucket)
			result.LocalStateRows = append(result.LocalStateRows, local)
			result.RelativeStrengthRows = append(result.RelativeStrengthRows, rs)
			localStateByKey[btcRegimeContextStateKey{symbol: symbol, timeframe: local.Timeframe, timestamp: local.DecisionCloseTime}] = local
			localStateByID[btcRegimeETHSOLContextStateIDKey(symbol, local.LocalStateRowID)] = local
		}
	}
	sort.Slice(result.LocalStateRows, func(i, j int) bool {
		if result.LocalStateRows[i].Symbol != result.LocalStateRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.LocalStateRows[i].Symbol) < rangeUniverseSymbolSortKey(result.LocalStateRows[j].Symbol)
		}
		if result.LocalStateRows[i].Timeframe != result.LocalStateRows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(result.LocalStateRows[i].Timeframe) < rangeContextTriageTimeframeSortKey(result.LocalStateRows[j].Timeframe)
		}
		return result.LocalStateRows[i].DecisionIndex < result.LocalStateRows[j].DecisionIndex
	})

	labelID := 0
	for _, symbol := range []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		data := dataBySymbol[symbol]
		for _, label := range data.labels {
			local, ok := localStateByID[btcRegimeETHSOLContextStateIDKey(symbol, label.StateRowID)]
			if !ok {
				continue
			}
			labelID++
			result.LabelRows = append(result.LabelRows, btcRegimeETHSOLContextLabelRow(labelID, symbol, local, label))
		}
	}

	result.CohortRows = btcRegimeETHSOLContextCohortRows(result.LabelRows, cfg, splits, btcRegimeETHSOLContextSourceCoveragePass(result.SourceRows, result.CoverageRows))
	result.RankingRows = btcRegimeETHSOLContextRankingRows(result.CohortRows)
	result.PassingCohorts = btcRegimeETHSOLContextPassingRankingCount(result.RankingRows)
	result.StopState = FuturesBTCRegimeETHSOLContextAuditStopState(result)
	btcRegimeETHSOLContextFillCoverageCounts(&result, dataBySymbol)
	result.SummaryRows = btcRegimeETHSOLContextSummaryRows(result, cfg, splits)
	return result, nil
}

func FuturesBTCRegimeETHSOLContextAuditStopState(result FuturesBTCRegimeETHSOLContextAuditResult) string {
	if result.StopState == BTCRegimeETHSOLContextStopStateSourceGap ||
		result.StopState == BTCRegimeETHSOLContextStopStateRejectedClosedReslice ||
		result.StopState == BTCRegimeETHSOLContextStopStateRejectedFutureLeak {
		return result.StopState
	}
	if !btcRegimeETHSOLContextSourceCoveragePass(result.SourceRows, result.CoverageRows) {
		return BTCRegimeETHSOLContextStopStateSourceGap
	}
	for _, row := range result.LocalStateRows {
		if row.ClosedFamilyReslice {
			return BTCRegimeETHSOLContextStopStateRejectedClosedReslice
		}
		if row.ForwardLabelsAsStateInput {
			return BTCRegimeETHSOLContextStopStateRejectedFutureLeak
		}
	}
	for _, row := range result.LabelRows {
		if row.ForwardLabelUsedAsFeature {
			return BTCRegimeETHSOLContextStopStateRejectedFutureLeak
		}
	}
	for _, row := range result.RankingRows {
		if !row.FutureLeakProtectionPass {
			return BTCRegimeETHSOLContextStopStateRejectedFutureLeak
		}
		if !row.ClosedFamilyProtectionPass {
			return BTCRegimeETHSOLContextStopStateRejectedClosedReslice
		}
		if row.PassesGate {
			return BTCRegimeETHSOLContextStopStatePassedNeedsSpec
		}
	}
	return BTCRegimeETHSOLContextStopStateFailedNoUsableContext
}

func (cfg FuturesBTCRegimeETHSOLContextAuditConfig) withDefaults() FuturesBTCRegimeETHSOLContextAuditConfig {
	defaults := DefaultFuturesBTCRegimeETHSOLContextAuditConfig()
	if len(cfg.Sources) == 0 {
		cfg.Sources = append([]FuturesBTCRegimeETHSOLContextSourceConfig(nil), defaults.Sources...)
	}
	for i := range cfg.Sources {
		if cfg.Sources[i].ExpectedRowCount == 0 {
			cfg.Sources[i].ExpectedRowCount = btcRegimeETHSOLContextExpectedRows
		}
		if cfg.Sources[i].ExpectedFirstOpenTime == "" {
			cfg.Sources[i].ExpectedFirstOpenTime = btcRegimeETHSOLContextExpectedFirst
		}
		if cfg.Sources[i].ExpectedLastOpenTime == "" {
			cfg.Sources[i].ExpectedLastOpenTime = btcRegimeETHSOLContextExpectedLast
		}
	}
	if len(cfg.StateConfig.Timeframes) == 0 {
		cfg.StateConfig = defaults.StateConfig
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
	return cfg
}

func (cfg FuturesBTCRegimeETHSOLContextAuditConfig) validate() error {
	if len(cfg.Sources) != 3 {
		return fmt.Errorf("BTC regime ETH/SOL context audit requires exactly BTCUSDT, ETHUSDT, and SOLUSDT sources")
	}
	seen := map[string]bool{}
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if symbol != RangeUniverseSymbolBTCUSDT && symbol != RangeUniverseSymbolETHUSDT && symbol != RangeUniverseSymbolSOLUSDT {
			return fmt.Errorf("BTC regime ETH/SOL context audit unsupported source symbol %q", source.Symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("BTC regime ETH/SOL context audit duplicate source symbol %s", symbol)
		}
		seen[symbol] = true
		if source.Path == "" || source.ApprovedPath == "" {
			return fmt.Errorf("BTC regime ETH/SOL context audit source paths must not be empty")
		}
	}
	for _, symbol := range []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		if !seen[symbol] {
			return fmt.Errorf("BTC regime ETH/SOL context audit missing required %s source", symbol)
		}
	}
	if err := cfg.StateConfig.withDefaults().validate(); err != nil {
		return err
	}
	if cfg.MinFullCohortRows <= 0 || cfg.MinSplitCohortRows <= 0 {
		return fmt.Errorf("BTC regime ETH/SOL context audit count gates must be positive")
	}
	if cfg.MaxSplitContributionRate <= 0 || cfg.MaxSplitContributionRate > 1 {
		return fmt.Errorf("BTC regime ETH/SOL context audit max split contribution must be in (0,1]")
	}
	return nil
}

func btcRegimeETHSOLContextLoadSymbol(source FuturesBTCRegimeETHSOLContextSourceConfig, cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split) btcRegimeContextSymbolData {
	source.Symbol = strings.ToUpper(strings.TrimSpace(source.Symbol))
	row := FuturesBTCRegimeETHSOLContextSourceRow{
		AuditName:                    FuturesBTCRegimeETHSOLContextAuditName,
		Role:                         btcRegimeETHSOLContextRole(source.Symbol),
		ApprovedLocalScopeOnly:       true,
		ExpectedRowCount:             source.ExpectedRowCount,
		ExpectedFirstOpenTime:        source.ExpectedFirstOpenTime,
		ExpectedLastOpenTime:         source.ExpectedLastOpenTime,
		ExpectedGapCount:             source.ExpectedGapCount,
		ExpectedDuplicateCount:       source.ExpectedDuplicateCount,
		ExpectedZeroVolumeCount:      source.ExpectedZeroVolumeCount,
		ExpectedPhysicalNonMonotonic: source.ExpectedPhysicalNonMonotonic,
		ExpectedSortedForValidation:  source.ExpectedSortedForValidation,
		SourceFactsPass:              true,
		ForwardLabelsAsSourceInput:   false,
	}
	universeSource := FuturesRangeUniverseSourceConfig{
		Symbol:                    source.Symbol,
		Path:                      source.Path,
		ApprovedPath:              source.ApprovedPath,
		SkipSplitEligibilityCheck: source.SkipSplitEligibilityCheck,
	}
	candles, universeRow, err := LoadFuturesRangeUniverseSource(universeSource, splits)
	row.FuturesRangeUniverseSourceRow = universeRow
	data := btcRegimeContextSymbolData{symbol: source.Symbol, candles: candles, source: row, sourceOK: true, skips: rangeStateSkipAccumulator{}}
	reject := func(reason string) btcRegimeContextSymbolData {
		row.SourceFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = reason
		data.source = row
		data.sourceOK = false
		return data
	}
	if err != nil {
		return reject(err.Error())
	}
	if universeRow.Product != "Binance USDT-M futures" || universeRow.ComparisonOnly || universeRow.Interval != "5m" {
		return reject(fmt.Sprintf("%s requires Binance USDT-M futures 5m comparison_only=false; got product=%q interval=%q comparison_only=%t", source.Symbol, universeRow.Product, universeRow.Interval, universeRow.ComparisonOnly))
	}
	if !source.SkipSourceFactCheck {
		switch {
		case universeRow.RowCount != source.ExpectedRowCount:
			return reject(fmt.Sprintf("%s rows=%d expected=%d", source.Symbol, universeRow.RowCount, source.ExpectedRowCount))
		case universeRow.FirstOpenTime != source.ExpectedFirstOpenTime:
			return reject(fmt.Sprintf("%s first_open_time=%s expected=%s", source.Symbol, universeRow.FirstOpenTime, source.ExpectedFirstOpenTime))
		case universeRow.LastOpenTime != source.ExpectedLastOpenTime:
			return reject(fmt.Sprintf("%s last_open_time=%s expected=%s", source.Symbol, universeRow.LastOpenTime, source.ExpectedLastOpenTime))
		case universeRow.GapCount != source.ExpectedGapCount:
			return reject(fmt.Sprintf("%s gap_count=%d expected=%d", source.Symbol, universeRow.GapCount, source.ExpectedGapCount))
		case universeRow.DuplicateCount != source.ExpectedDuplicateCount:
			return reject(fmt.Sprintf("%s duplicate_count=%d expected=%d", source.Symbol, universeRow.DuplicateCount, source.ExpectedDuplicateCount))
		case universeRow.ZeroVolumeCount != source.ExpectedZeroVolumeCount:
			return reject(fmt.Sprintf("%s zero_volume_count=%d expected=%d", source.Symbol, universeRow.ZeroVolumeCount, source.ExpectedZeroVolumeCount))
		case universeRow.PhysicalNonMonotonicCount != source.ExpectedPhysicalNonMonotonic:
			return reject(fmt.Sprintf("%s physical_non_monotonic_count=%d expected=%d", source.Symbol, universeRow.PhysicalNonMonotonicCount, source.ExpectedPhysicalNonMonotonic))
		case universeRow.SortedForValidation != source.ExpectedSortedForValidation:
			return reject(fmt.Sprintf("%s sorted_for_validation=%t expected=%t", source.Symbol, universeRow.SortedForValidation, source.ExpectedSortedForValidation))
		}
	}
	data.source = row
	return data
}

func btcRegimeETHSOLContextBuildSymbolStates(data btcRegimeContextSymbolData, cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split) (btcRegimeContextSymbolData, error) {
	data.coverageOK = true
	stateCfg := cfg.StateConfig.withDefaults()
	frameDataByTimeframe := map[string]rangeStateFrameData{}
	for _, timeframe := range stateCfg.Timeframes {
		frame, ok := rangeContextTriageFrameDef(timeframe)
		if !ok {
			return data, fmt.Errorf("BTC regime ETH/SOL context missing frame definition for %s", timeframe)
		}
		frameCandles, coverage, err := resampleRangeDiscoveryFrame(data.candles, frame)
		coverageRow := btcRegimeETHSOLContextCoverageRow(data.symbol, coverage, stateCfg, cfg)
		data.coverage = append(data.coverage, coverageRow)
		if err != nil || !coverageRow.CoverageFactsPass || coverageRow.ValidationStatus != "accepted" || !coverageRow.Complete {
			data.coverageOK = false
			continue
		}
		horizons := rangeStateConstructionLoopHorizons(timeframe, stateCfg)
		frameDataByTimeframe[timeframe] = rangeStateFrameData{
			frame:      frame,
			candles:    frameCandles,
			coverage:   FuturesRangeStateConstructionLoopCoverageRow{FuturesRangeDiscoveryCoverageRow: coverage},
			metrics:    rangeStateConstructionLoopMetrics(frameCandles, frame, stateCfg),
			horizons:   horizons,
			maxHorizon: maxIntInSlice(horizons),
		}
	}
	if !data.coverageOK {
		return data, nil
	}
	for _, timeframe := range stateCfg.Timeframes {
		stateRows, err := rangeStateConstructionLoopStateRows(frameDataByTimeframe[timeframe], frameDataByTimeframe, stateCfg, splits, data.skips, len(data.states))
		if err != nil {
			return data, err
		}
		data.states = append(data.states, stateRows...)
	}
	data.labels = rangeStateConstructionLoopLabelRows(data.states, frameDataByTimeframe, stateCfg)
	return data, nil
}

func btcRegimeETHSOLContextCoverageRow(symbol string, base FuturesRangeDiscoveryCoverageRow, stateCfg FuturesRangeStateConstructionLoopAuditConfig, cfg FuturesBTCRegimeETHSOLContextAuditConfig) FuturesBTCRegimeETHSOLContextCoverageRow {
	stateCfg = stateCfg.withDefaults()
	row := FuturesBTCRegimeETHSOLContextCoverageRow{
		AuditName:                        FuturesBTCRegimeETHSOLContextAuditName,
		Symbol:                           symbol,
		Role:                             btcRegimeETHSOLContextRole(symbol),
		ExpectedFirstOpenTime:            stateCfg.ExpectedFirstOpenTime,
		CoverageFactsPass:                true,
		FuturesRangeDiscoveryCoverageRow: base,
	}
	switch base.Timeframe {
	case RangeDiscoveryTimeframe15m:
		row.ExpectedRowCount = stateCfg.Expected15MRows
		row.ExpectedLastOpenTime = stateCfg.Expected15MLastOpenTime
	case RangeDiscoveryTimeframe1h:
		row.ExpectedRowCount = stateCfg.Expected1HRows
		row.ExpectedLastOpenTime = stateCfg.Expected1HLastOpenTime
	case RangeDiscoveryTimeframe4h:
		row.ExpectedRowCount = stateCfg.Expected4HRows
		row.ExpectedLastOpenTime = stateCfg.Expected4HLastOpenTime
	default:
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("unsupported BTC regime ETH/SOL context timeframe %q", base.Timeframe)
		return row
	}
	if base.ValidationStatus != "accepted" || !base.Complete {
		row.CoverageFactsPass = false
		return row
	}
	if !stateCfg.SkipCoverageCountCheck && (row.RowCount != row.ExpectedRowCount || row.FirstOpenTime != row.ExpectedFirstOpenTime || row.LastOpenTime != row.ExpectedLastOpenTime) {
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("%s %s coverage row_count=%d first=%s last=%s expected row_count=%d first=%s last=%s", symbol, base.Timeframe, row.RowCount, row.FirstOpenTime, row.LastOpenTime, row.ExpectedRowCount, row.ExpectedFirstOpenTime, row.ExpectedLastOpenTime)
	}
	if row.GapCount != 0 || row.DuplicateCount != 0 || row.MissingChildOpenCount != 0 {
		row.CoverageFactsPass = false
		if row.ValidationError == "" {
			row.ValidationStatus = "rejected"
			row.ValidationError = fmt.Sprintf("%s %s coverage gap=%d duplicate=%d missing_child_open=%d", symbol, base.Timeframe, row.GapCount, row.DuplicateCount, row.MissingChildOpenCount)
		}
	}
	return row
}

func btcRegimeETHSOLContextBTCStateRows(states []FuturesRangeStateConstructionLoopStateRow) []FuturesBTCRegimeETHSOLContextBTCStateRow {
	rows := make([]FuturesBTCRegimeETHSOLContextBTCStateRow, 0, len(states))
	for _, state := range states {
		rows = append(rows, FuturesBTCRegimeETHSOLContextBTCStateRow{
			BTCStateRowID:             state.StateRowID,
			Timestamp:                 state.Timestamp,
			Timeframe:                 state.Timeframe,
			Split:                     state.Split,
			DecisionIndex:             state.DecisionIndex,
			DecisionCloseTime:         state.DecisionCloseTime,
			RangeEpisodeID:            state.RangeEpisodeID,
			RangeAgeBars:              state.RangeAgeBars,
			CloseLocationPct:          state.CloseLocationPct,
			ReturnShortPct:            state.ReturnShortPct,
			ReturnMediumPct:           state.ReturnMediumPct,
			RealizedVolPercentileLong: state.RealizedVolPercentileLong,
			VolumePercentileLong:      state.VolumePercentileLong,
			StateID:                   state.StateID,
			BTCRegimeID:               btcRegimeETHSOLContextBTCRegimeID(state),
			GeometryBucket:            state.GeometryBucket,
			VolBucket:                 state.VolBucket,
			TrendBucket:               state.TrendBucket,
			ImpulseBucket:             state.ImpulseBucket,
			ParticipationBucket:       state.ParticipationBucket,
			CloseLocationBucket:       btcRegimeETHSOLContextCloseLocationBucket(state.CloseLocationPct),
			RangeAgeBucket:            btcRegimeETHSOLContextRangeAgeBucket(state.RangeAgeBars),
			ClosedCandleOnly:          true,
			DiagnosticOnlyAuthority:   true,
			ForwardLabelsAsStateInput: false,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Timeframe != rows[j].Timeframe {
			return rangeContextTriageTimeframeSortKey(rows[i].Timeframe) < rangeContextTriageTimeframeSortKey(rows[j].Timeframe)
		}
		return rows[i].DecisionIndex < rows[j].DecisionIndex
	})
	return rows
}

func btcRegimeETHSOLContextLocalStateRow(symbol string, state FuturesRangeStateConstructionLoopStateRow, btc FuturesBTCRegimeETHSOLContextBTCStateRow, rsBucket string) FuturesBTCRegimeETHSOLContextLocalStateRow {
	localBucket := btcRegimeETHSOLContextLocalBucketID(state)
	contextStateID := fmt.Sprintf("btc_regime_eth_sol_v1::%s::%s::%s::%s::%s", symbol, state.Timeframe, btc.BTCRegimeID, localBucket, rsBucket)
	return FuturesBTCRegimeETHSOLContextLocalStateRow{
		LocalStateRowID:           state.StateRowID,
		Symbol:                    symbol,
		Timestamp:                 state.Timestamp,
		Timeframe:                 state.Timeframe,
		Split:                     state.Split,
		DecisionIndex:             state.DecisionIndex,
		DecisionCloseTime:         state.DecisionCloseTime,
		RangeEpisodeID:            state.RangeEpisodeID,
		RangeAgeBars:              state.RangeAgeBars,
		CloseLocationPct:          state.CloseLocationPct,
		ReturnShortPct:            state.ReturnShortPct,
		ReturnMediumPct:           state.ReturnMediumPct,
		StateID:                   state.StateID,
		LocalRangeBucketID:        localBucket,
		BTCStateRowID:             btc.BTCStateRowID,
		BTCRegimeID:               btc.BTCRegimeID,
		RelativeStrengthBucket:    rsBucket,
		ContextStateID:            contextStateID,
		GeometryBucket:            state.GeometryBucket,
		VolBucket:                 state.VolBucket,
		TrendBucket:               state.TrendBucket,
		ImpulseBucket:             state.ImpulseBucket,
		ParticipationBucket:       state.ParticipationBucket,
		ClosedCandleOnly:          true,
		AuthorityCandidateOnly:    true,
		ForwardLabelsAsStateInput: false,
		ClosedFamilyReslice:       false,
	}
}

func btcRegimeETHSOLContextRelativeStrengthRow(id int, symbol string, local FuturesRangeStateConstructionLoopStateRow, btc FuturesBTCRegimeETHSOLContextBTCStateRow) FuturesBTCRegimeETHSOLContextRelativeStrengthRow {
	short := local.ReturnShortPct - btc.ReturnShortPct
	medium := local.ReturnMediumPct - btc.ReturnMediumPct
	return FuturesBTCRegimeETHSOLContextRelativeStrengthRow{
		RelativeStrengthRowID:   id,
		Symbol:                  symbol,
		Timestamp:               local.Timestamp,
		Timeframe:               local.Timeframe,
		DecisionIndex:           local.DecisionIndex,
		LocalStateRowID:         local.StateRowID,
		BTCStateRowID:           btc.BTCStateRowID,
		LocalReturnShortPct:     local.ReturnShortPct,
		BTCReturnShortPct:       btc.ReturnShortPct,
		RelativeReturnShortPct:  short,
		LocalReturnMediumPct:    local.ReturnMediumPct,
		BTCReturnMediumPct:      btc.ReturnMediumPct,
		RelativeReturnMediumPct: medium,
		RelativeStrengthBucket:  btcRegimeETHSOLContextRelativeStrengthBucket(short, medium),
		ClosedCandleOnly:        true,
		ForwardLabelsAsRSInput:  false,
	}
}

func btcRegimeETHSOLContextLabelRow(id int, symbol string, local FuturesBTCRegimeETHSOLContextLocalStateRow, label FuturesRangeStateConstructionLoopLabelRow) FuturesBTCRegimeETHSOLContextLabelRow {
	return FuturesBTCRegimeETHSOLContextLabelRow{
		LabelRowID:                id,
		Symbol:                    symbol,
		LocalStateRowID:           local.LocalStateRowID,
		BTCStateRowID:             local.BTCStateRowID,
		Timestamp:                 label.Timestamp,
		Timeframe:                 label.Timeframe,
		Split:                     label.Split,
		HorizonBars:               label.HorizonBars,
		LocalRangeBucketID:        local.LocalRangeBucketID,
		BTCRegimeID:               local.BTCRegimeID,
		RelativeStrengthBucket:    local.RelativeStrengthBucket,
		ContextStateID:            local.ContextStateID,
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
		InsideCloseRate:           label.InsideCloseRate,
		MaxExcursionAboveWidths:   label.MaxExcursionAboveRangeWidths,
		MaxExcursionBelowWidths:   label.MaxExcursionBelowRangeWidths,
		ForwardLabelMetadataOnly:  true,
		ForwardLabelUsedAsFeature: false,
	}
}

func btcRegimeETHSOLContextCohortRows(labels []FuturesBTCRegimeETHSOLContextLabelRow, cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split, sourcePass bool) []FuturesBTCRegimeETHSOLContextCohortRow {
	acc := map[btcRegimeContextCohortKey]*btcRegimeContextCohortAccumulator{}
	for _, label := range labels {
		for _, split := range rangeDiscoverySplitCombos(label.Split) {
			for _, route := range []string{RangeStateConstructionLoopRouteRotation, RangeStateConstructionLoopRouteContinuation, RangeStateConstructionLoopRouteNoTradeToxic} {
				for _, comparison := range []string{BTCRegimeETHSOLContextComparisonLocalOnly, BTCRegimeETHSOLContextComparisonBTCPlusLocal} {
					key := btcRegimeETHSOLContextCohortKeyFor(label, split, route, comparison)
					a := acc[key]
					if a == nil {
						a = &btcRegimeContextCohortAccumulator{labels: map[string]int{}}
						a.row = FuturesBTCRegimeETHSOLContextCohortRow{
							CohortID:                   key.cohortID,
							ComparisonType:             key.comparisonType,
							Symbol:                     key.symbol,
							Split:                      key.split,
							Timeframe:                  key.timeframe,
							HorizonBars:                key.horizonBars,
							RouteCandidate:             key.routeCandidate,
							LocalRangeBucketID:         key.localBucketID,
							BTCRegimeID:                key.btcRegimeID,
							RelativeStrengthBucket:     key.rsBucket,
							ClosedFamilyProtectionPass: true,
							FutureLeakProtectionPass:   true,
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
	rows := make([]FuturesBTCRegimeETHSOLContextCohortRow, 0, len(acc))
	for _, a := range acc {
		rows = append(rows, a.finalRow())
	}
	btcRegimeETHSOLContextMarkCohortGates(rows, cfg, splits, sourcePass)
	sort.Slice(rows, func(i, j int) bool {
		return btcRegimeETHSOLContextLessCohort(rows[i], rows[j])
	})
	return rows
}

func (acc *btcRegimeContextCohortAccumulator) add(label FuturesBTCRegimeETHSOLContextLabelRow) {
	acc.row.CandidateCount++
	acc.labels[label.ForwardLabel]++
	if btcRegimeETHSOLContextRouteUseful(acc.row.RouteCandidate, label) {
		acc.row.UsefulCount++
	}
	if btcRegimeETHSOLContextRouteToxic(acc.row.RouteCandidate, label) {
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

func (acc *btcRegimeContextCohortAccumulator) finalRow() FuturesBTCRegimeETHSOLContextCohortRow {
	row := acc.row
	if row.CandidateCount > 0 {
		row.UsefulRate = float64(row.UsefulCount) / float64(row.CandidateCount)
		row.ToxicRate = float64(row.ToxicCount) / float64(row.CandidateCount)
		row.UsefulMinusToxicMargin = row.UsefulRate - row.ToxicRate
	}
	row.DominantForwardLabel, row.DominantForwardLabelRate = rangeContextTriageDominantLabel(acc.labels, row.CandidateCount, false)
	return row
}

func btcRegimeETHSOLContextMarkCohortGates(rows []FuturesBTCRegimeETHSOLContextCohortRow, cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split, sourcePass bool) {
	byIDSplit := map[string]map[string]*FuturesBTCRegimeETHSOLContextCohortRow{}
	localBaseline := map[string]map[string]*FuturesBTCRegimeETHSOLContextCohortRow{}
	for i := range rows {
		if byIDSplit[rows[i].CohortID] == nil {
			byIDSplit[rows[i].CohortID] = map[string]*FuturesBTCRegimeETHSOLContextCohortRow{}
		}
		byIDSplit[rows[i].CohortID][rows[i].Split] = &rows[i]
		if rows[i].ComparisonType == BTCRegimeETHSOLContextComparisonLocalOnly {
			baseKey := btcRegimeETHSOLContextBaselineKey(rows[i])
			if localBaseline[baseKey] == nil {
				localBaseline[baseKey] = map[string]*FuturesBTCRegimeETHSOLContextCohortRow{}
			}
			localBaseline[baseKey][rows[i].Split] = &rows[i]
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
		rows[i].ClosedFamilyProtectionPass = true
		rows[i].FutureLeakProtectionPass = !strings.Contains(rows[i].CohortID, "forward_label")
		rows[i].RouteRateGatePass = btcRegimeETHSOLContextRouteRateGate(rows[i], cfg)

		if rows[i].ComparisonType == BTCRegimeETHSOLContextComparisonBTCPlusLocal {
			btcRegimeETHSOLContextFillBaselineComparison(&rows[i], localBaseline[btcRegimeETHSOLContextBaselineKey(rows[i])], periodSplits)
			rows[i].BTCContextImprovementGatePass = btcRegimeETHSOLContextImprovementGate(rows[i], cfg)
		} else {
			rows[i].BTCContextImprovementGatePass = false
		}

		if !sourcePass {
			reasons = append(reasons, "source_or_coverage_gap")
		}
		if rows[i].ComparisonType != BTCRegimeETHSOLContextComparisonBTCPlusLocal {
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
		if !rows[i].BTCContextImprovementGatePass {
			reasons = append(reasons, "btc_context_improvement_gate_failed")
		}
		if !rows[i].FutureLeakProtectionPass {
			reasons = append(reasons, "future_label_as_feature")
		}
		rows[i].PassesReviewGate = sourcePass &&
			rows[i].ComparisonType == BTCRegimeETHSOLContextComparisonBTCPlusLocal &&
			rows[i].ReviewableCountGatePass &&
			rows[i].SplitStabilityGatePass &&
			rows[i].SplitContributionGatePass &&
			rows[i].RouteRateGatePass &&
			rows[i].BTCContextImprovementGatePass &&
			rows[i].ClosedFamilyProtectionPass &&
			rows[i].FutureLeakProtectionPass &&
			rows[i].Split == fullSplitName
		rows[i].FailureReason = uniqueJoinedReasons(reasons)
	}
}

func btcRegimeETHSOLContextFillBaselineComparison(row *FuturesBTCRegimeETHSOLContextCohortRow, baselineBySplit map[string]*FuturesBTCRegimeETHSOLContextCohortRow, periodSplits []Split) {
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
		if base.UsefulRate < row.BaselineWeakestSplitUsefulRate {
			row.BaselineWeakestSplitUsefulRate = base.UsefulRate
		}
		if base.ToxicRate > row.BaselineWorstSplitToxicRate {
			row.BaselineWorstSplitToxicRate = base.ToxicRate
		}
		if base.UsefulMinusToxicMargin < row.BaselineWeakestSplitMargin {
			row.BaselineWeakestSplitMargin = base.UsefulMinusToxicMargin
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

func btcRegimeETHSOLContextRankingRows(cohorts []FuturesBTCRegimeETHSOLContextCohortRow) []FuturesBTCRegimeETHSOLContextRankingRow {
	rows := []FuturesBTCRegimeETHSOLContextRankingRow{}
	for _, cohort := range cohorts {
		if cohort.Split != fullSplitName || cohort.ComparisonType != BTCRegimeETHSOLContextComparisonBTCPlusLocal {
			continue
		}
		row := FuturesBTCRegimeETHSOLContextRankingRow{
			CohortID:                      cohort.CohortID,
			Symbol:                        cohort.Symbol,
			Timeframe:                     cohort.Timeframe,
			HorizonBars:                   cohort.HorizonBars,
			RouteCandidate:                cohort.RouteCandidate,
			LocalRangeBucketID:            cohort.LocalRangeBucketID,
			BTCRegimeID:                   cohort.BTCRegimeID,
			RelativeStrengthBucket:        cohort.RelativeStrengthBucket,
			PassesGate:                    cohort.PassesReviewGate,
			FullPeriodRows:                cohort.FullPeriodRows,
			WeakestSplitRows:              cohort.WeakestSplitRows,
			MaxSplitContributionRate:      cohort.MaxSplitContributionRate,
			FullUsefulRate:                cohort.FullUsefulRate,
			WeakestSplitUsefulRate:        cohort.WeakestSplitUsefulRate,
			FullToxicRate:                 cohort.FullToxicRate,
			WorstSplitToxicRate:           cohort.WorstSplitToxicRate,
			FullUsefulMinusToxicMargin:    cohort.FullUsefulMinusToxicMargin,
			WeakestSplitMargin:            cohort.WeakestSplitMargin,
			FullUsefulImprovement:         cohort.FullUsefulImprovement,
			WeakestSplitUsefulImprovement: cohort.WeakestSplitUsefulImprovement,
			FullToxicImprovement:          cohort.FullToxicImprovement,
			WeakestSplitToxicImprovement:  cohort.WeakestSplitToxicImprovement,
			FullMarginImprovement:         cohort.FullMarginImprovement,
			WeakestSplitMarginImprovement: cohort.WeakestSplitMarginImprovement,
			DominantForwardLabel:          cohort.DominantForwardLabel,
			DominantForwardLabelRate:      cohort.DominantForwardLabelRate,
			FutureLeakProtectionPass:      cohort.FutureLeakProtectionPass,
			ClosedFamilyProtectionPass:    cohort.ClosedFamilyProtectionPass,
			FailureReason:                 cohort.FailureReason,
		}
		row.RankScore = btcRegimeETHSOLContextRankScore(row)
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
		return rows[i].CohortID < rows[j].CohortID
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func btcRegimeETHSOLContextSummaryRows(result FuturesBTCRegimeETHSOLContextAuditResult, cfg FuturesBTCRegimeETHSOLContextAuditConfig, splits []Split) []FuturesBTCRegimeETHSOLContextSummaryRow {
	sourcePass := btcRegimeETHSOLContextSourceCoveragePass(result.SourceRows, result.CoverageRows)
	rows := []FuturesBTCRegimeETHSOLContextSummaryRow{{
		Split:                        fullSplitName,
		Symbol:                       "all",
		Timeframe:                    "all",
		HorizonBars:                  0,
		SourceRows:                   len(result.SourceRows),
		CoverageRows:                 len(result.CoverageRows),
		BTCStateRows:                 len(result.BTCStateRows),
		LocalStateRows:               len(result.LocalStateRows),
		RelativeStrengthRows:         len(result.RelativeStrengthRows),
		LabelRows:                    len(result.LabelRows),
		CohortRows:                   len(result.CohortRows),
		RankingRows:                  len(result.RankingRows),
		PassingCohorts:               result.PassingCohorts,
		SourceScopePass:              sourcePass,
		CoveragePass:                 sourcePass,
		CommonOutputsZeroTrade:       true,
		BTCContextDiagnosticOnly:     true,
		ETHSOLAuthorityCandidateOnly: true,
		ForwardLabelsAsInputs:        false,
		StopState:                    result.StopState,
	}}
	for _, split := range splits {
		for _, symbol := range []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
			for _, timeframe := range cfg.StateConfig.withDefaults().Timeframes {
				for _, horizon := range rangeStateConstructionLoopHorizons(timeframe, cfg.StateConfig.withDefaults()) {
					row := FuturesBTCRegimeETHSOLContextSummaryRow{
						Split:                        split.Name,
						Symbol:                       symbol,
						Timeframe:                    timeframe,
						HorizonBars:                  horizon,
						SourceRows:                   len(result.SourceRows),
						SourceScopePass:              sourcePass,
						CoveragePass:                 sourcePass,
						CommonOutputsZeroTrade:       true,
						BTCContextDiagnosticOnly:     true,
						ETHSOLAuthorityCandidateOnly: true,
						ForwardLabelsAsInputs:        false,
						StopState:                    result.StopState,
					}
					for _, local := range result.LocalStateRows {
						if local.Symbol == symbol && local.Timeframe == timeframe && rangeStateRowInSplit(local.DecisionCloseTime, split) {
							row.LocalStateRows++
						}
					}
					for _, rs := range result.RelativeStrengthRows {
						if rs.Symbol == symbol && rs.Timeframe == timeframe && rangeStateRowInSplit(rs.Timestamp, split) {
							row.RelativeStrengthRows++
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

func btcRegimeETHSOLContextFillCoverageCounts(result *FuturesBTCRegimeETHSOLContextAuditResult, dataBySymbol map[string]btcRegimeContextSymbolData) {
	for i := range result.CoverageRows {
		row := &result.CoverageRows[i]
		data := dataBySymbol[row.Symbol]
		for _, state := range data.states {
			if state.Timeframe == row.Timeframe {
				row.StateRows++
			}
		}
		if row.Symbol == RangeUniverseSymbolBTCUSDT {
			for _, btc := range result.BTCStateRows {
				if btc.Timeframe == row.Timeframe {
					row.ContextMatchedRows++
				}
			}
			continue
		}
		for _, label := range data.labels {
			if label.Timeframe == row.Timeframe {
				row.LabelRows++
			}
		}
		for _, local := range result.LocalStateRows {
			if local.Symbol == row.Symbol && local.Timeframe == row.Timeframe {
				row.ContextMatchedRows++
			}
		}
		row.MissingBTCContextRows = row.StateRows - row.ContextMatchedRows
		if row.MissingBTCContextRows < 0 {
			row.MissingBTCContextRows = 0
		}
	}
}

func btcRegimeETHSOLContextRequiredSourcesPresent(dataBySymbol map[string]btcRegimeContextSymbolData) bool {
	for _, symbol := range []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		if _, ok := dataBySymbol[symbol]; !ok {
			return false
		}
	}
	return true
}

func btcRegimeETHSOLContextSourceCoveragePass(sources []FuturesBTCRegimeETHSOLContextSourceRow, coverage []FuturesBTCRegimeETHSOLContextCoverageRow) bool {
	if len(sources) != 3 {
		return false
	}
	for _, row := range sources {
		if row.ValidationStatus != "accepted" || !row.SourceFactsPass || row.ForwardLabelsAsSourceInput {
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

func btcRegimeETHSOLContextRole(symbol string) string {
	if symbol == RangeUniverseSymbolBTCUSDT {
		return BTCRegimeETHSOLContextRoleBTCMarketRegime
	}
	return BTCRegimeETHSOLContextRoleAuthorityLocal
}

func btcRegimeETHSOLContextStateIDKey(symbol string, stateID int) string {
	return symbol + "|" + fmt.Sprintf("%d", stateID)
}

func btcRegimeETHSOLContextBTCRegimeID(state FuturesRangeStateConstructionLoopStateRow) string {
	return strings.Join([]string{
		state.Timeframe,
		state.GeometryBucket,
		state.VolBucket,
		state.TrendBucket,
		state.ImpulseBucket,
		state.ParticipationBucket,
		btcRegimeETHSOLContextCloseLocationBucket(state.CloseLocationPct),
		btcRegimeETHSOLContextRangeAgeBucket(state.RangeAgeBars),
	}, "::")
}

func btcRegimeETHSOLContextLocalBucketID(state FuturesRangeStateConstructionLoopStateRow) string {
	return strings.Join([]string{
		state.Timeframe,
		state.GeometryBucket,
		state.VolBucket,
		state.TrendBucket,
		state.ImpulseBucket,
	}, "::")
}

func btcRegimeETHSOLContextCloseLocationBucket(location float64) string {
	switch {
	case !validNumber(location):
		return "close_location_unknown"
	case location < 0.20:
		return "close_near_low"
	case location < 0.45:
		return "close_lower_mid"
	case location <= 0.55:
		return "close_mid"
	case location <= 0.80:
		return "close_upper_mid"
	default:
		return "close_near_high"
	}
}

func btcRegimeETHSOLContextRangeAgeBucket(age int) string {
	switch {
	case age <= 0:
		return "age_unknown"
	case age < 12:
		return "age_young"
	case age < 48:
		return "age_mature"
	default:
		return "age_old"
	}
}

func btcRegimeETHSOLContextRelativeStrengthBucket(short, medium float64) string {
	combined := (short + medium) / 2
	switch {
	case !validNumber(combined):
		return "relative_strength_unknown"
	case combined >= 0.004:
		return "relative_strength_strong"
	case combined <= -0.004:
		return "relative_strength_weak"
	default:
		return "relative_strength_neutral"
	}
}

func btcRegimeETHSOLContextRouteUseful(route string, label FuturesBTCRegimeETHSOLContextLabelRow) bool {
	switch route {
	case RangeStateConstructionLoopRouteRotation:
		return label.RotationUseful
	case RangeStateConstructionLoopRouteContinuation:
		return label.ContinuationUseful
	default:
		return false
	}
}

func btcRegimeETHSOLContextRouteToxic(route string, label FuturesBTCRegimeETHSOLContextLabelRow) bool {
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

func btcRegimeETHSOLContextCohortKeyFor(label FuturesBTCRegimeETHSOLContextLabelRow, split string, route string, comparison string) btcRegimeContextCohortKey {
	local := label.LocalRangeBucketID
	btc := ""
	rs := ""
	if comparison == BTCRegimeETHSOLContextComparisonBTCPlusLocal {
		btc = label.BTCRegimeID
		rs = label.RelativeStrengthBucket
	}
	parts := []string{"btc_regime_eth_sol_context_v1", comparison, label.Symbol, label.Timeframe, fmt.Sprintf("h%d", label.HorizonBars), route, "local=" + local}
	if btc != "" {
		parts = append(parts, "btc="+btc, "rs="+rs)
	}
	return btcRegimeContextCohortKey{
		cohortID:       strings.Join(parts, "|"),
		comparisonType: comparison,
		symbol:         label.Symbol,
		split:          split,
		timeframe:      label.Timeframe,
		horizonBars:    label.HorizonBars,
		routeCandidate: route,
		localBucketID:  local,
		btcRegimeID:    btc,
		rsBucket:       rs,
	}
}

func btcRegimeETHSOLContextBaselineKey(row FuturesBTCRegimeETHSOLContextCohortRow) string {
	return strings.Join([]string{row.Symbol, row.Timeframe, fmt.Sprintf("%d", row.HorizonBars), row.RouteCandidate, row.LocalRangeBucketID}, "|")
}

func btcRegimeETHSOLContextRouteRateGate(row FuturesBTCRegimeETHSOLContextCohortRow, cfg FuturesBTCRegimeETHSOLContextAuditConfig) bool {
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

func btcRegimeETHSOLContextImprovementGate(row FuturesBTCRegimeETHSOLContextCohortRow, cfg FuturesBTCRegimeETHSOLContextAuditConfig) bool {
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

func btcRegimeETHSOLContextRankScore(row FuturesBTCRegimeETHSOLContextRankingRow) float64 {
	switch row.RouteCandidate {
	case RangeStateConstructionLoopRouteNoTradeToxic:
		return row.FullToxicRate + row.WorstSplitToxicRate + row.FullToxicImprovement + row.WeakestSplitToxicImprovement
	default:
		return row.FullUsefulMinusToxicMargin + row.WeakestSplitMargin + row.FullMarginImprovement + row.WeakestSplitMarginImprovement
	}
}

func btcRegimeETHSOLContextPassingRankingCount(rows []FuturesBTCRegimeETHSOLContextRankingRow) int {
	count := 0
	for _, row := range rows {
		if row.PassesGate {
			count++
		}
	}
	return count
}

func btcRegimeETHSOLContextLessCohort(a, b FuturesBTCRegimeETHSOLContextCohortRow) bool {
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
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	return a.CohortID < b.CohortID
}
