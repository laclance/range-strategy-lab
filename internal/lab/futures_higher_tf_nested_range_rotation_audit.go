package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const (
	FuturesHigherTFNestedRangeRotationAuditName = "futures_higher_tf_nested_range_rotation_audit"

	HigherTFNestedRangeRotationStopStateSourceGap           = "higher_tf_nested_range_rotation_audit_source_gap"
	HigherTFNestedRangeRotationStopStateNoCandidateEvents   = "higher_tf_nested_range_rotation_audit_no_candidate_events"
	HigherTFNestedRangeRotationStopStateClosedFamilyReslice = "higher_tf_nested_range_rotation_audit_rejected_as_closed_family_reslice"
	HigherTFNestedRangeRotationStopStateFailedNoBaseline    = "higher_tf_nested_range_rotation_audit_failed_no_baseline"
	HigherTFNestedRangeRotationStopStateReadyForBaseline    = "higher_tf_nested_range_rotation_audit_ready_for_baseline_brief"

	HigherTFNestedRangeRotationEventUp   = "nested_rotation_up"
	HigherTFNestedRangeRotationEventDown = "nested_rotation_down"

	HigherTFNestedRangeRotationOutcomeFavorableMidpoint         = "favorable_midpoint"
	HigherTFNestedRangeRotationOutcomeFavorableFarQuartile      = "favorable_far_quartile"
	HigherTFNestedRangeRotationOutcomeAdverseChildInvalidation  = "adverse_child_invalidation"
	HigherTFNestedRangeRotationOutcomeAdverseParentInvalidation = "adverse_parent_invalidation"
	HigherTFNestedRangeRotationOutcomeNoResolution              = "no_resolution"
	HigherTFNestedRangeRotationOutcomeMissingFuture             = "missing_future"

	nestedRotationSourcePath      = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"
	nestedRotationExpectedRows    = 573984
	nestedRotationExpectedFirst   = "2021-01-01T00:00:00Z"
	nestedRotationExpectedLast    = "2026-06-16T23:55:00Z"
	nestedRotationExpectedZeroVol = 66
	nestedRotationExpected1HRows  = 47832
	nestedRotationExpected4HRows  = 11958
	nestedRotationExpected1HLast  = "2026-06-16T23:00:00Z"
	nestedRotationExpected4HLast  = "2026-06-16T20:00:00Z"
)

type FuturesHigherTFNestedRangeRotationAuditConfig struct {
	SourcePath                   string
	ApprovedSourcePath           string
	ExpectedSourceRows           int
	ExpectedFirstOpenTime        string
	ExpectedLastOpenTime         string
	ExpectedGapCount             int
	ExpectedDuplicateCount       int
	ExpectedZeroVolumeCount      int
	SkipSourceFactCheck          bool
	ParentTimeframe              string
	ChildTimeframe               string
	ExpectedParentRows           int
	ExpectedParentLastOpenTime   string
	ExpectedChildRows            int
	ExpectedChildLastOpenTime    string
	SkipCoverageCountCheck       bool
	DetectorLookbackDays         int
	DetectorLookbackBarsOverride int
	DetectorPercentile           float64
	DetectorMinConsecutiveBars   int
	ChildWidthMaxParentFraction  float64
	OutcomeHorizonBars           int
	QuickInvalidationBars        int
	MinFullEvents                int
	MinSplitEvents               int
	MinSideEvents                int
	MinFarQuartileRate           float64
}

type FuturesHigherTFNestedRangeRotationAuditResult struct {
	SourceRows      []FuturesHigherTFNestedRangeRotationSourceRow      `json:"source_rows"`
	CoverageRows    []FuturesHigherTFNestedRangeRotationCoverageRow    `json:"coverage_rows"`
	ParentRangeRows []FuturesHigherTFNestedRangeRotationParentRangeRow `json:"parent_range_rows"`
	ChildRangeRows  []FuturesHigherTFNestedRangeRotationChildRangeRow  `json:"child_range_rows"`
	EventRows       []FuturesHigherTFNestedRangeRotationEventRow       `json:"event_rows"`
	SummaryRows     []FuturesHigherTFNestedRangeRotationSummaryRow     `json:"summary_rows"`
	StopState       string                                             `json:"stop_state"`
}

type FuturesHigherTFNestedRangeRotationSourceRow struct {
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

type FuturesHigherTFNestedRangeRotationCoverageRow struct {
	Role                 string `json:"role"`
	ExpectedRowCount     int    `json:"expected_row_count"`
	ExpectedLastOpenTime string `json:"expected_last_open_time"`
	CoverageFactsPass    bool   `json:"coverage_facts_pass"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesHigherTFNestedRangeRotationParentRangeRow struct {
	ParentRangeID         int     `json:"parent_range_id"`
	Timeframe             string  `json:"timeframe"`
	StartIndex            int     `json:"start_index"`
	MatureIndex           int     `json:"mature_index"`
	RawEndIndex           int     `json:"raw_end_index"`
	StartCloseTime        string  `json:"start_close_time"`
	MatureCloseTime       string  `json:"mature_close_time"`
	RawEndCloseTime       string  `json:"raw_end_close_time"`
	InvalidationIndex     int     `json:"invalidation_index"`
	InvalidationCloseTime string  `json:"invalidation_close_time,omitempty"`
	Invalidated           bool    `json:"invalidated"`
	High                  float64 `json:"high"`
	Low                   float64 `json:"low"`
	Mid                   float64 `json:"mid"`
	UpperQuartile         float64 `json:"upper_quartile"`
	LowerQuartile         float64 `json:"lower_quartile"`
	Width                 float64 `json:"width"`
	Eligible              bool    `json:"eligible"`
	SkippedReason         string  `json:"skipped_reason,omitempty"`
}

type FuturesHigherTFNestedRangeRotationChildRangeRow struct {
	ChildRangeID             int     `json:"child_range_id"`
	ParentRangeID            int     `json:"parent_range_id"`
	Timeframe                string  `json:"timeframe"`
	CandidateSide            string  `json:"candidate_side"`
	StartIndex               int     `json:"start_index"`
	MatureIndex              int     `json:"mature_index"`
	RawEndIndex              int     `json:"raw_end_index"`
	StartCloseTime           string  `json:"start_close_time"`
	MatureCloseTime          string  `json:"mature_close_time"`
	RawEndCloseTime          string  `json:"raw_end_close_time"`
	High                     float64 `json:"high"`
	Low                      float64 `json:"low"`
	Mid                      float64 `json:"mid"`
	Width                    float64 `json:"width"`
	ParentHigh               float64 `json:"parent_high"`
	ParentLow                float64 `json:"parent_low"`
	ParentMid                float64 `json:"parent_mid"`
	ParentUpperQuartile      float64 `json:"parent_upper_quartile"`
	ParentLowerQuartile      float64 `json:"parent_lower_quartile"`
	ParentWidth              float64 `json:"parent_width"`
	ChildWidthParentFraction float64 `json:"child_width_parent_fraction"`
	Eligible                 bool    `json:"eligible"`
	SkippedReason            string  `json:"skipped_reason,omitempty"`
}

type FuturesHigherTFNestedRangeRotationEventRow struct {
	EventID                       int     `json:"event_id"`
	ChildRangeID                  int     `json:"child_range_id"`
	ParentRangeID                 int     `json:"parent_range_id"`
	EventType                     string  `json:"event_type"`
	Side                          string  `json:"side"`
	Split                         string  `json:"split"`
	EventIndex                    int     `json:"event_index"`
	EventOpenTime                 string  `json:"event_open_time"`
	EventCloseTime                string  `json:"event_close_time"`
	EventOpen                     float64 `json:"event_open"`
	EventHigh                     float64 `json:"event_high"`
	EventLow                      float64 `json:"event_low"`
	EventClose                    float64 `json:"event_close"`
	ParentHigh                    float64 `json:"parent_high"`
	ParentLow                     float64 `json:"parent_low"`
	ParentMid                     float64 `json:"parent_mid"`
	ParentUpperQuartile           float64 `json:"parent_upper_quartile"`
	ParentLowerQuartile           float64 `json:"parent_lower_quartile"`
	ParentWidth                   float64 `json:"parent_width"`
	ChildHigh                     float64 `json:"child_high"`
	ChildLow                      float64 `json:"child_low"`
	ChildMid                      float64 `json:"child_mid"`
	ChildWidth                    float64 `json:"child_width"`
	ChildWidthParentFraction      float64 `json:"child_width_parent_fraction"`
	LabelWindowStartIndex         int     `json:"label_window_start_index"`
	LabelWindowEndIndex           int     `json:"label_window_end_index"`
	LabelWindowStartTime          string  `json:"label_window_start_time"`
	LabelWindowEndTime            string  `json:"label_window_end_time"`
	OutcomeLabel                  string  `json:"outcome_label"`
	FavorableMidpoint             bool    `json:"favorable_midpoint"`
	FavorableFarQuartile          bool    `json:"favorable_far_quartile"`
	AdverseChildInvalidation      bool    `json:"adverse_child_invalidation"`
	AdverseParentInvalidation     bool    `json:"adverse_parent_invalidation"`
	NoResolution                  bool    `json:"no_resolution"`
	QuickInvalidation             bool    `json:"quick_invalidation"`
	MissingFuture                 bool    `json:"missing_future"`
	BarsToMidpoint                int     `json:"bars_to_midpoint"`
	BarsToFarQuartile             int     `json:"bars_to_far_quartile"`
	BarsToAdverseChild            int     `json:"bars_to_adverse_child"`
	BarsToAdverseParent           int     `json:"bars_to_adverse_parent"`
	FavorableExcursionParentWidth float64 `json:"favorable_excursion_parent_width"`
	AdverseExcursionParentWidth   float64 `json:"adverse_excursion_parent_width"`
	SkippedReason                 string  `json:"skipped_reason,omitempty"`
}

type FuturesHigherTFNestedRangeRotationSummaryRow struct {
	Split                            string  `json:"split"`
	Side                             string  `json:"side"`
	ChildCount                       int     `json:"child_count"`
	SkippedChildCount                int     `json:"skipped_child_count"`
	EventCount                       int     `json:"event_count"`
	SkippedEventCount                int     `json:"skipped_event_count"`
	FavorableMidpointCount           int     `json:"favorable_midpoint_count"`
	FavorableFarQuartileCount        int     `json:"favorable_far_quartile_count"`
	AdverseChildInvalidationCount    int     `json:"adverse_child_invalidation_count"`
	AdverseParentInvalidationCount   int     `json:"adverse_parent_invalidation_count"`
	NoResolutionCount                int     `json:"no_resolution_count"`
	MissingFutureCount               int     `json:"missing_future_count"`
	QuickInvalidationCount           int     `json:"quick_invalidation_count"`
	FavorableMidpointRate            float64 `json:"favorable_midpoint_rate"`
	FavorableFarQuartileRate         float64 `json:"favorable_far_quartile_rate"`
	AdverseChildInvalidationRate     float64 `json:"adverse_child_invalidation_rate"`
	AdverseParentInvalidationRate    float64 `json:"adverse_parent_invalidation_rate"`
	NoResolutionRate                 float64 `json:"no_resolution_rate"`
	MissingFutureRate                float64 `json:"missing_future_rate"`
	QuickInvalidationRate            float64 `json:"quick_invalidation_rate"`
	AvgFavorableExcursionParentWidth float64 `json:"avg_favorable_excursion_parent_width"`
	AvgAdverseExcursionParentWidth   float64 `json:"avg_adverse_excursion_parent_width"`
	MinFullEvents                    int     `json:"min_full_events"`
	MinSplitEvents                   int     `json:"min_split_events"`
	MinSideEvents                    int     `json:"min_side_events"`
	MinFarQuartileRate               float64 `json:"min_far_quartile_rate"`
	SourceResamplePass               bool    `json:"source_resample_pass"`
	CountGatePass                    bool    `json:"count_gate_pass"`
	SideGatePass                     bool    `json:"side_gate_pass"`
	FavorableBeatsAdversePass        bool    `json:"favorable_beats_adverse_pass"`
	FavorableBeatsQuickPass          bool    `json:"favorable_beats_quick_pass"`
	FarQuartileGatePass              bool    `json:"far_quartile_gate_pass"`
	ExcursionGatePass                bool    `json:"excursion_gate_pass"`
	RegimeSpreadPass                 bool    `json:"regime_spread_pass"`
	PassesReviewGate                 bool    `json:"passes_review_gate"`
	Caveat                           string  `json:"caveat,omitempty"`
	FailureReason                    string  `json:"failure_reason,omitempty"`
	StopState                        string  `json:"stop_state"`
}

type nestedRotationRange struct {
	id                int
	startIndex        int
	matureIndex       int
	rawEndIndex       int
	high              float64
	low               float64
	mid               float64
	upperQuartile     float64
	lowerQuartile     float64
	width             float64
	eligible          bool
	skippedReason     string
	invalidationIndex int
}

func DefaultFuturesHigherTFNestedRangeRotationAuditConfig() FuturesHigherTFNestedRangeRotationAuditConfig {
	return FuturesHigherTFNestedRangeRotationAuditConfig{
		SourcePath:                  nestedRotationSourcePath,
		ApprovedSourcePath:          nestedRotationSourcePath,
		ExpectedSourceRows:          nestedRotationExpectedRows,
		ExpectedFirstOpenTime:       nestedRotationExpectedFirst,
		ExpectedLastOpenTime:        nestedRotationExpectedLast,
		ExpectedGapCount:            0,
		ExpectedDuplicateCount:      0,
		ExpectedZeroVolumeCount:     nestedRotationExpectedZeroVol,
		ParentTimeframe:             RangeDiscoveryTimeframe4h,
		ChildTimeframe:              RangeDiscoveryTimeframe1h,
		ExpectedParentRows:          nestedRotationExpected4HRows,
		ExpectedParentLastOpenTime:  nestedRotationExpected4HLast,
		ExpectedChildRows:           nestedRotationExpected1HRows,
		ExpectedChildLastOpenTime:   nestedRotationExpected1HLast,
		DetectorLookbackDays:        20,
		DetectorPercentile:          0.30,
		DetectorMinConsecutiveBars:  12,
		ChildWidthMaxParentFraction: 0.40,
		OutcomeHorizonBars:          24,
		QuickInvalidationBars:       6,
		MinFullEvents:               100,
		MinSplitEvents:              25,
		MinSideEvents:               25,
		MinFarQuartileRate:          0.20,
	}
}

func RunFuturesHigherTFNestedRangeRotationAudit(candles []Candle, manifest SourceManifest, cfg FuturesHigherTFNestedRangeRotationAuditConfig, splits []Split) (FuturesHigherTFNestedRangeRotationAuditResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesHigherTFNestedRangeRotationAuditResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesHigherTFNestedRangeRotationAuditResult{}
	sourceRow := nestedRotationSourceRow(manifest, cfg)
	result.SourceRows = append(result.SourceRows, sourceRow)
	if sourceRow.ValidationStatus != "accepted" {
		result.StopState = HigherTFNestedRangeRotationStopStateSourceGap
		result.SummaryRows = SummarizeFuturesHigherTFNestedRangeRotationAudit(result, cfg, splits)
		return result, nil
	}

	parentFrame, err := nestedRotationFrameDef(cfg.ParentTimeframe)
	if err != nil {
		return result, err
	}
	childFrame, err := nestedRotationFrameDef(cfg.ChildTimeframe)
	if err != nil {
		return result, err
	}
	parentCandles, parentCoverage, parentErr := resampleRangeDiscoveryFrame(candles, parentFrame)
	parentCoverageRow := nestedRotationCoverageRow("parent", parentCoverage, cfg.ExpectedParentRows, cfg.ExpectedParentLastOpenTime, cfg)
	result.CoverageRows = append(result.CoverageRows, parentCoverageRow)
	childCandles, childCoverage, childErr := resampleRangeDiscoveryFrame(candles, childFrame)
	childCoverageRow := nestedRotationCoverageRow("child", childCoverage, cfg.ExpectedChildRows, cfg.ExpectedChildLastOpenTime, cfg)
	result.CoverageRows = append(result.CoverageRows, childCoverageRow)
	if parentErr != nil || childErr != nil || !parentCoverageRow.CoverageFactsPass || !childCoverageRow.CoverageFactsPass {
		result.StopState = HigherTFNestedRangeRotationStopStateSourceGap
		result.SummaryRows = SummarizeFuturesHigherTFNestedRangeRotationAudit(result, cfg, splits)
		return result, nil
	}

	parentClassifications, err := (CompressionRangeDetector{Config: nestedRotationDetectorConfig(cfg, parentFrame.barsPerDay)}).Classify(parentCandles)
	if err != nil {
		return result, err
	}
	childClassifications, err := (CompressionRangeDetector{Config: nestedRotationDetectorConfig(cfg, childFrame.barsPerDay)}).Classify(childCandles)
	if err != nil {
		return result, err
	}

	parentRanges := nestedRotationParentRanges(parentCandles, parentClassifications, cfg)
	result.ParentRangeRows = nestedRotationParentRows(parentCandles, parentRanges, cfg.ParentTimeframe)
	childRanges := nestedRotationChildRanges(childCandles, childClassifications, parentCandles, parentRanges, cfg)
	result.ChildRangeRows = childRanges
	result.EventRows = nestedRotationEventRows(childCandles, childRanges, cfg, splits)
	result.SummaryRows = SummarizeFuturesHigherTFNestedRangeRotationAudit(result, cfg, splits)
	result.StopState = FuturesHigherTFNestedRangeRotationAuditStopState(result)
	for i := range result.SummaryRows {
		result.SummaryRows[i].StopState = result.StopState
	}
	return result, nil
}

func FuturesHigherTFNestedRangeRotationAuditStopState(result FuturesHigherTFNestedRangeRotationAuditResult) string {
	if result.StopState == HigherTFNestedRangeRotationStopStateSourceGap || result.StopState == HigherTFNestedRangeRotationStopStateClosedFamilyReslice {
		return result.StopState
	}
	for _, row := range result.SourceRows {
		if row.ValidationStatus != "accepted" {
			return HigherTFNestedRangeRotationStopStateSourceGap
		}
	}
	for _, row := range result.CoverageRows {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass {
			return HigherTFNestedRangeRotationStopStateSourceGap
		}
	}
	validEvents := 0
	for _, row := range result.EventRows {
		if row.SkippedReason == "" && !row.MissingFuture {
			validEvents++
		}
	}
	if validEvents == 0 {
		return HigherTFNestedRangeRotationStopStateNoCandidateEvents
	}
	for _, row := range result.SummaryRows {
		if row.Split == fullSplitName && row.Side == "all" && row.PassesReviewGate {
			return HigherTFNestedRangeRotationStopStateReadyForBaseline
		}
	}
	return HigherTFNestedRangeRotationStopStateFailedNoBaseline
}

func SummarizeFuturesHigherTFNestedRangeRotationAudit(result FuturesHigherTFNestedRangeRotationAuditResult, cfg FuturesHigherTFNestedRangeRotationAuditConfig, splits []Split) []FuturesHigherTFNestedRangeRotationSummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	sourcePass := true
	for _, row := range result.SourceRows {
		if row.ValidationStatus != "accepted" {
			sourcePass = false
		}
	}
	for _, row := range result.CoverageRows {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass {
			sourcePass = false
		}
	}

	rows := make([]FuturesHigherTFNestedRangeRotationSummaryRow, 0, len(splits)*3)
	for _, split := range splits {
		for _, side := range []string{"all", RangeDiscoverySideUp, RangeDiscoverySideDown} {
			row := summarizeNestedRotationSplitSide(result, cfg, split, side)
			row.SourceResamplePass = sourcePass
			row.MinFullEvents = cfg.MinFullEvents
			row.MinSplitEvents = cfg.MinSplitEvents
			row.MinSideEvents = cfg.MinSideEvents
			row.MinFarQuartileRate = cfg.MinFarQuartileRate
			rows = append(rows, row)
		}
	}
	markNestedRotationReviewGates(rows, cfg, sourcePass, splits)
	return rows
}

func (cfg FuturesHigherTFNestedRangeRotationAuditConfig) withDefaults() FuturesHigherTFNestedRangeRotationAuditConfig {
	defaults := DefaultFuturesHigherTFNestedRangeRotationAuditConfig()
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
	if cfg.ParentTimeframe == "" {
		cfg.ParentTimeframe = defaults.ParentTimeframe
	}
	if cfg.ChildTimeframe == "" {
		cfg.ChildTimeframe = defaults.ChildTimeframe
	}
	if cfg.ExpectedParentRows == 0 {
		cfg.ExpectedParentRows = defaults.ExpectedParentRows
	}
	if cfg.ExpectedParentLastOpenTime == "" {
		cfg.ExpectedParentLastOpenTime = defaults.ExpectedParentLastOpenTime
	}
	if cfg.ExpectedChildRows == 0 {
		cfg.ExpectedChildRows = defaults.ExpectedChildRows
	}
	if cfg.ExpectedChildLastOpenTime == "" {
		cfg.ExpectedChildLastOpenTime = defaults.ExpectedChildLastOpenTime
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
	if cfg.ChildWidthMaxParentFraction == 0 {
		cfg.ChildWidthMaxParentFraction = defaults.ChildWidthMaxParentFraction
	}
	if cfg.OutcomeHorizonBars == 0 {
		cfg.OutcomeHorizonBars = defaults.OutcomeHorizonBars
	}
	if cfg.QuickInvalidationBars == 0 {
		cfg.QuickInvalidationBars = defaults.QuickInvalidationBars
	}
	if cfg.MinFullEvents == 0 {
		cfg.MinFullEvents = defaults.MinFullEvents
	}
	if cfg.MinSplitEvents == 0 {
		cfg.MinSplitEvents = defaults.MinSplitEvents
	}
	if cfg.MinSideEvents == 0 {
		cfg.MinSideEvents = defaults.MinSideEvents
	}
	if cfg.MinFarQuartileRate == 0 {
		cfg.MinFarQuartileRate = defaults.MinFarQuartileRate
	}
	return cfg
}

func (cfg FuturesHigherTFNestedRangeRotationAuditConfig) validate() error {
	if cfg.SourcePath == "" || cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("nested range rotation source paths must not be empty")
	}
	if cfg.ParentTimeframe != RangeDiscoveryTimeframe4h || cfg.ChildTimeframe != RangeDiscoveryTimeframe1h {
		return fmt.Errorf("nested range rotation requires parent=%s child=%s", RangeDiscoveryTimeframe4h, RangeDiscoveryTimeframe1h)
	}
	if cfg.DetectorLookbackDays <= 0 && cfg.DetectorLookbackBarsOverride <= 0 {
		return fmt.Errorf("nested range rotation detector lookback must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("nested range rotation detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("nested range rotation detector min consecutive bars must be positive")
	}
	if cfg.ChildWidthMaxParentFraction <= 0 || cfg.ChildWidthMaxParentFraction > 1 {
		return fmt.Errorf("nested range rotation child width fraction must be in (0,1]")
	}
	if cfg.OutcomeHorizonBars <= 0 || cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("nested range rotation horizons must be positive")
	}
	if cfg.MinFullEvents <= 0 || cfg.MinSplitEvents <= 0 || cfg.MinSideEvents <= 0 {
		return fmt.Errorf("nested range rotation review count gates must be positive")
	}
	if cfg.MinFarQuartileRate <= 0 || cfg.MinFarQuartileRate > 1 {
		return fmt.Errorf("nested range rotation far-quartile gate must be in (0,1]")
	}
	return nil
}

func nestedRotationSourceRow(manifest SourceManifest, cfg FuturesHigherTFNestedRangeRotationAuditConfig) FuturesHigherTFNestedRangeRotationSourceRow {
	row := FuturesHigherTFNestedRangeRotationSourceRow{
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
		reject(fmt.Sprintf("source must be BTCUSDT Binance USDT-M futures 5m comparison_only=false; got product=%q symbol=%q interval=%q comparison_only=%t", manifest.Product, manifest.Symbol, manifest.Interval, manifest.ComparisonOnly))
		return row
	}
	if !sameCleanPath(manifest.Path, cfg.ApprovedSourcePath) {
		reject(fmt.Sprintf("source path %q is not approved path %q", manifest.Path, cfg.ApprovedSourcePath))
		return row
	}
	if cfg.SkipSourceFactCheck {
		return row
	}
	switch {
	case manifest.RowCount != cfg.ExpectedSourceRows:
		reject(fmt.Sprintf("source row_count=%d does not match expected %d", manifest.RowCount, cfg.ExpectedSourceRows))
	case manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime:
		reject(fmt.Sprintf("source first_open_time=%s does not match expected %s", manifest.FirstOpenTime, cfg.ExpectedFirstOpenTime))
	case manifest.LastOpenTime != cfg.ExpectedLastOpenTime:
		reject(fmt.Sprintf("source last_open_time=%s does not match expected %s", manifest.LastOpenTime, cfg.ExpectedLastOpenTime))
	case manifest.GapCount != cfg.ExpectedGapCount:
		reject(fmt.Sprintf("source gap_count=%d does not match expected %d", manifest.GapCount, cfg.ExpectedGapCount))
	case manifest.DuplicateCount != cfg.ExpectedDuplicateCount:
		reject(fmt.Sprintf("source duplicate_count=%d does not match expected %d", manifest.DuplicateCount, cfg.ExpectedDuplicateCount))
	case manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount:
		reject(fmt.Sprintf("source zero_volume_count=%d does not match expected %d", manifest.ZeroVolumeCount, cfg.ExpectedZeroVolumeCount))
	}
	return row
}

func nestedRotationCoverageRow(role string, coverage FuturesRangeDiscoveryCoverageRow, expectedRows int, expectedLast string, cfg FuturesHigherTFNestedRangeRotationAuditConfig) FuturesHigherTFNestedRangeRotationCoverageRow {
	row := FuturesHigherTFNestedRangeRotationCoverageRow{
		Role:                             role,
		ExpectedRowCount:                 expectedRows,
		ExpectedLastOpenTime:             expectedLast,
		CoverageFactsPass:                true,
		FuturesRangeDiscoveryCoverageRow: coverage,
	}
	if coverage.ValidationStatus != "accepted" || !coverage.Complete {
		row.CoverageFactsPass = false
		return row
	}
	if cfg.SkipCoverageCountCheck {
		return row
	}
	if coverage.RowCount != expectedRows {
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("%s row_count=%d does not match expected %d", role, coverage.RowCount, expectedRows)
	}
	if coverage.LastOpenTime != expectedLast {
		row.CoverageFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("%s last_open_time=%s does not match expected %s", role, coverage.LastOpenTime, expectedLast)
	}
	return row
}

func nestedRotationFrameDef(timeframe string) (rangeDiscoveryFrameDef, error) {
	for _, frame := range rangeDiscoveryFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, nil
		}
	}
	return rangeDiscoveryFrameDef{}, fmt.Errorf("unsupported nested range rotation timeframe %q", timeframe)
}

func nestedRotationDetectorConfig(cfg FuturesHigherTFNestedRangeRotationAuditConfig, barsPerDay int) RangeDetectorConfig {
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

func nestedRotationParentRanges(candles []Candle, classifications []RangeClassification, cfg FuturesHigherTFNestedRangeRotationAuditConfig) []nestedRotationRange {
	ranges := nestedRotationFrozenRanges(candles, classifications)
	for i := range ranges {
		if ranges[i].width <= 0 || !finitePositive(ranges[i].width) {
			ranges[i].eligible = false
			ranges[i].skippedReason = "non_positive_parent_width"
			ranges[i].invalidationIndex = -1
			continue
		}
		ranges[i].eligible = true
		ranges[i].invalidationIndex = -1
		for j := ranges[i].matureIndex + 1; j < len(candles); j++ {
			if candles[j].Close > ranges[i].high || candles[j].Close < ranges[i].low {
				ranges[i].invalidationIndex = j
				break
			}
		}
	}
	return ranges
}

func nestedRotationFrozenRanges(candles []Candle, classifications []RangeClassification) []nestedRotationRange {
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	ranges := []nestedRotationRange{}
	for i := 0; i < len(classifications); {
		if !classifications[i].RawActive {
			i++
			continue
		}
		start := i
		firstActive := -1
		high := candles[i].High
		low := candles[i].Low
		freezeHigh := high
		freezeLow := low
		for i < len(classifications) && classifications[i].RawActive {
			if candles[i].High > high {
				high = candles[i].High
			}
			if candles[i].Low < low {
				low = candles[i].Low
			}
			if classifications[i].Active && firstActive == -1 {
				firstActive = i
				freezeHigh = high
				freezeLow = low
			}
			i++
		}
		rawEnd := i - 1
		if firstActive < 0 {
			continue
		}
		width := freezeHigh - freezeLow
		ranges = append(ranges, nestedRotationRange{
			id:                len(ranges) + 1,
			startIndex:        start,
			matureIndex:       firstActive,
			rawEndIndex:       rawEnd,
			high:              freezeHigh,
			low:               freezeLow,
			mid:               (freezeHigh + freezeLow) / 2,
			upperQuartile:     freezeLow + width*0.75,
			lowerQuartile:     freezeLow + width*0.25,
			width:             width,
			invalidationIndex: -1,
		})
	}
	return ranges
}

func nestedRotationParentRows(candles []Candle, ranges []nestedRotationRange, timeframe string) []FuturesHigherTFNestedRangeRotationParentRangeRow {
	rows := make([]FuturesHigherTFNestedRangeRotationParentRangeRow, 0, len(ranges))
	for _, r := range ranges {
		row := FuturesHigherTFNestedRangeRotationParentRangeRow{
			ParentRangeID:     r.id,
			Timeframe:         timeframe,
			StartIndex:        r.startIndex,
			MatureIndex:       r.matureIndex,
			RawEndIndex:       r.rawEndIndex,
			High:              r.high,
			Low:               r.low,
			Mid:               r.mid,
			UpperQuartile:     r.upperQuartile,
			LowerQuartile:     r.lowerQuartile,
			Width:             r.width,
			Eligible:          r.eligible,
			SkippedReason:     r.skippedReason,
			InvalidationIndex: r.invalidationIndex,
		}
		if r.startIndex >= 0 && r.startIndex < len(candles) {
			row.StartCloseTime = candles[r.startIndex].CloseTime.UTC().Format(timeLayout)
		}
		if r.matureIndex >= 0 && r.matureIndex < len(candles) {
			row.MatureCloseTime = candles[r.matureIndex].CloseTime.UTC().Format(timeLayout)
		}
		if r.rawEndIndex >= 0 && r.rawEndIndex < len(candles) {
			row.RawEndCloseTime = candles[r.rawEndIndex].CloseTime.UTC().Format(timeLayout)
		}
		if r.invalidationIndex >= 0 && r.invalidationIndex < len(candles) {
			row.Invalidated = true
			row.InvalidationCloseTime = candles[r.invalidationIndex].CloseTime.UTC().Format(timeLayout)
		}
		rows = append(rows, row)
	}
	return rows
}

func nestedRotationChildRanges(childCandles []Candle, childClassifications []RangeClassification, parentCandles []Candle, parents []nestedRotationRange, cfg FuturesHigherTFNestedRangeRotationAuditConfig) []FuturesHigherTFNestedRangeRotationChildRangeRow {
	childFrozen := nestedRotationFrozenRanges(childCandles, childClassifications)
	rows := make([]FuturesHigherTFNestedRangeRotationChildRangeRow, 0, len(childFrozen))
	for _, child := range childFrozen {
		row := nestedRotationChildRangeRow(childCandles, child, cfg.ChildTimeframe)
		parent, ok := nestedRotationParentForChild(childCandles[child.matureIndex].CloseTime, parentCandles, parents)
		if !ok {
			row.SkippedReason = "no_valid_parent"
			rows = append(rows, row)
			continue
		}
		row.ParentRangeID = parent.id
		row.ParentHigh = parent.high
		row.ParentLow = parent.low
		row.ParentMid = parent.mid
		row.ParentUpperQuartile = parent.upperQuartile
		row.ParentLowerQuartile = parent.lowerQuartile
		row.ParentWidth = parent.width
		if parent.width > 0 {
			row.ChildWidthParentFraction = row.Width / parent.width
		}
		switch {
		case row.Width <= 0 || !finitePositive(row.Width):
			row.SkippedReason = "non_positive_child_width"
		case parent.width <= 0 || !finitePositive(parent.width):
			row.SkippedReason = "non_positive_parent_width"
		case row.Low < parent.low || row.High > parent.high:
			row.SkippedReason = "child_not_inside_parent"
		case row.ChildWidthParentFraction > cfg.ChildWidthMaxParentFraction:
			row.SkippedReason = "child_width_above_40pct_parent"
		case row.Mid < parent.mid:
			row.CandidateSide = RangeDiscoverySideUp
			row.Eligible = true
		case row.Mid > parent.mid:
			row.CandidateSide = RangeDiscoverySideDown
			row.Eligible = true
		default:
			row.SkippedReason = "child_midpoint_wrong_parent_half"
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].MatureIndex < rows[j].MatureIndex
	})
	return rows
}

func nestedRotationChildRangeRow(candles []Candle, child nestedRotationRange, timeframe string) FuturesHigherTFNestedRangeRotationChildRangeRow {
	row := FuturesHigherTFNestedRangeRotationChildRangeRow{
		ChildRangeID: child.id,
		Timeframe:    timeframe,
		StartIndex:   child.startIndex,
		MatureIndex:  child.matureIndex,
		RawEndIndex:  child.rawEndIndex,
		High:         child.high,
		Low:          child.low,
		Mid:          child.mid,
		Width:        child.width,
	}
	if child.startIndex >= 0 && child.startIndex < len(candles) {
		row.StartCloseTime = candles[child.startIndex].CloseTime.UTC().Format(timeLayout)
	}
	if child.matureIndex >= 0 && child.matureIndex < len(candles) {
		row.MatureCloseTime = candles[child.matureIndex].CloseTime.UTC().Format(timeLayout)
	}
	if child.rawEndIndex >= 0 && child.rawEndIndex < len(candles) {
		row.RawEndCloseTime = candles[child.rawEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	return row
}

func nestedRotationParentForChild(childClose time.Time, parentCandles []Candle, parents []nestedRotationRange) (nestedRotationRange, bool) {
	var selected nestedRotationRange
	ok := false
	for _, parent := range parents {
		if !parent.eligible || parent.matureIndex < 0 || parent.matureIndex >= len(parentCandles) {
			continue
		}
		parentMature := parentCandles[parent.matureIndex].CloseTime
		if childClose.Before(parentMature) {
			continue
		}
		if parent.invalidationIndex >= 0 && parent.invalidationIndex < len(parentCandles) && !childClose.Before(parentCandles[parent.invalidationIndex].CloseTime) {
			continue
		}
		if !ok || parentMature.After(parentCandles[selected.matureIndex].CloseTime) {
			selected = parent
			ok = true
		}
	}
	return selected, ok
}

func nestedRotationEventRows(candles []Candle, children []FuturesHigherTFNestedRangeRotationChildRangeRow, cfg FuturesHigherTFNestedRangeRotationAuditConfig, splits []Split) []FuturesHigherTFNestedRangeRotationEventRow {
	rows := []FuturesHigherTFNestedRangeRotationEventRow{}
	nextID := 0
	for _, child := range children {
		if !child.Eligible {
			continue
		}
		attempted := false
		for i := child.MatureIndex + 1; i < len(candles); i++ {
			candle := candles[i]
			if child.ParentRangeID > 0 && (candle.Close < child.ParentLow || candle.Close > child.ParentHigh) {
				break
			}
			breakUp := candle.Close > child.High
			breakDown := candle.Close < child.Low
			if !breakUp && !breakDown {
				continue
			}
			nextID++
			row := nestedRotationNewEventRow(nextID, candles, child, i, splits)
			if attempted {
				row.SkippedReason = "duplicate_child_event"
				rows = append(rows, row)
				break
			}
			attempted = true
			if breakUp {
				row.EventType = HigherTFNestedRangeRotationEventUp
				row.Side = RangeDiscoverySideUp
			} else {
				row.EventType = HigherTFNestedRangeRotationEventDown
				row.Side = RangeDiscoverySideDown
			}
			switch {
			case row.Side != child.CandidateSide:
				row.SkippedReason = "event_wrong_rotation_side"
			case candle.High > child.ParentHigh || candle.Low < child.ParentLow:
				row.SkippedReason = "event_outside_parent"
			case row.Side == RangeDiscoverySideUp && candle.Close >= child.ParentMid:
				row.SkippedReason = "event_beyond_parent_midpoint"
			case row.Side == RangeDiscoverySideDown && candle.Close <= child.ParentMid:
				row.SkippedReason = "event_beyond_parent_midpoint"
			default:
				row = nestedRotationLabelEvent(row, candles, cfg)
			}
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].EventIndex != rows[j].EventIndex {
			return rows[i].EventIndex < rows[j].EventIndex
		}
		return rows[i].EventID < rows[j].EventID
	})
	return rows
}

func nestedRotationNewEventRow(id int, candles []Candle, child FuturesHigherTFNestedRangeRotationChildRangeRow, eventIndex int, splits []Split) FuturesHigherTFNestedRangeRotationEventRow {
	candle := candles[eventIndex]
	return FuturesHigherTFNestedRangeRotationEventRow{
		EventID:                  id,
		ChildRangeID:             child.ChildRangeID,
		ParentRangeID:            child.ParentRangeID,
		EventIndex:               eventIndex,
		EventOpenTime:            candle.OpenTime.UTC().Format(timeLayout),
		EventCloseTime:           candle.CloseTime.UTC().Format(timeLayout),
		EventOpen:                candle.Open,
		EventHigh:                candle.High,
		EventLow:                 candle.Low,
		EventClose:               candle.Close,
		Split:                    splitNameForCloseTime(candle.CloseTime, splits),
		ParentHigh:               child.ParentHigh,
		ParentLow:                child.ParentLow,
		ParentMid:                child.ParentMid,
		ParentUpperQuartile:      child.ParentUpperQuartile,
		ParentLowerQuartile:      child.ParentLowerQuartile,
		ParentWidth:              child.ParentWidth,
		ChildHigh:                child.High,
		ChildLow:                 child.Low,
		ChildMid:                 child.Mid,
		ChildWidth:               child.Width,
		ChildWidthParentFraction: child.ChildWidthParentFraction,
		LabelWindowStartIndex:    -1,
		LabelWindowEndIndex:      -1,
		BarsToMidpoint:           -1,
		BarsToFarQuartile:        -1,
		BarsToAdverseChild:       -1,
		BarsToAdverseParent:      -1,
	}
}

func nestedRotationLabelEvent(row FuturesHigherTFNestedRangeRotationEventRow, candles []Candle, cfg FuturesHigherTFNestedRangeRotationAuditConfig) FuturesHigherTFNestedRangeRotationEventRow {
	start := row.EventIndex + 1
	end := row.EventIndex + cfg.OutcomeHorizonBars
	if start < len(candles) {
		row.LabelWindowStartIndex = start
		row.LabelWindowStartTime = candles[start].CloseTime.UTC().Format(timeLayout)
		row.LabelWindowEndIndex = minInt(end, len(candles)-1)
		row.LabelWindowEndTime = candles[row.LabelWindowEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	if end >= len(candles) || row.ParentWidth <= 0 {
		row.MissingFuture = true
		row.OutcomeLabel = HigherTFNestedRangeRotationOutcomeMissingFuture
		return row
	}
	row.LabelWindowEndIndex = end
	row.LabelWindowEndTime = candles[end].CloseTime.UTC().Format(timeLayout)
	midpointReached := false
	farTarget := row.ParentUpperQuartile
	if row.Side == RangeDiscoverySideDown {
		farTarget = row.ParentLowerQuartile
	}
	for i := start; i <= end; i++ {
		candle := candles[i]
		delay := i - row.EventIndex
		if row.Side == RangeDiscoverySideUp {
			row.FavorableExcursionParentWidth = math.Max(row.FavorableExcursionParentWidth, math.Max(0, candle.High-row.EventClose)/row.ParentWidth)
			row.AdverseExcursionParentWidth = math.Max(row.AdverseExcursionParentWidth, math.Max(0, row.EventClose-candle.Low)/row.ParentWidth)
			adverseParent := candle.Close < row.ParentLow
			adverseChild := candle.Low <= row.ChildLow
			farTouch := candle.High >= farTarget
			midTouch := candle.High >= row.ParentMid
			row = nestedRotationApplyOutcomeStep(row, delay, midpointReached, adverseParent, adverseChild, midTouch, farTouch, cfg)
		} else {
			row.FavorableExcursionParentWidth = math.Max(row.FavorableExcursionParentWidth, math.Max(0, row.EventClose-candle.Low)/row.ParentWidth)
			row.AdverseExcursionParentWidth = math.Max(row.AdverseExcursionParentWidth, math.Max(0, candle.High-row.EventClose)/row.ParentWidth)
			adverseParent := candle.Close > row.ParentHigh
			adverseChild := candle.High >= row.ChildHigh
			farTouch := candle.Low <= farTarget
			midTouch := candle.Low <= row.ParentMid
			row = nestedRotationApplyOutcomeStep(row, delay, midpointReached, adverseParent, adverseChild, midTouch, farTouch, cfg)
		}
		if row.FavorableMidpoint {
			midpointReached = true
		}
		if row.OutcomeLabel == HigherTFNestedRangeRotationOutcomeAdverseChildInvalidation || row.OutcomeLabel == HigherTFNestedRangeRotationOutcomeAdverseParentInvalidation {
			break
		}
	}
	if row.OutcomeLabel == "" {
		row.NoResolution = true
		row.OutcomeLabel = HigherTFNestedRangeRotationOutcomeNoResolution
	}
	return row
}

func nestedRotationApplyOutcomeStep(row FuturesHigherTFNestedRangeRotationEventRow, delay int, midpointReached bool, adverseParent bool, adverseChild bool, midTouch bool, farTouch bool, cfg FuturesHigherTFNestedRangeRotationAuditConfig) FuturesHigherTFNestedRangeRotationEventRow {
	if !midpointReached {
		if adverseParent || adverseChild {
			if adverseParent {
				row.AdverseParentInvalidation = true
				row.BarsToAdverseParent = delay
				row.OutcomeLabel = HigherTFNestedRangeRotationOutcomeAdverseParentInvalidation
			}
			if adverseChild {
				row.AdverseChildInvalidation = true
				row.BarsToAdverseChild = delay
				if row.OutcomeLabel == "" {
					row.OutcomeLabel = HigherTFNestedRangeRotationOutcomeAdverseChildInvalidation
				}
				if delay <= cfg.QuickInvalidationBars {
					row.QuickInvalidation = true
				}
			}
			return row
		}
		if midTouch {
			row.FavorableMidpoint = true
			row.BarsToMidpoint = delay
			row.OutcomeLabel = HigherTFNestedRangeRotationOutcomeFavorableMidpoint
			if farTouch {
				row.FavorableFarQuartile = true
				row.BarsToFarQuartile = delay
			}
		}
		return row
	}
	if !row.FavorableFarQuartile {
		if adverseParent || adverseChild {
			return row
		}
		if farTouch {
			row.FavorableFarQuartile = true
			row.BarsToFarQuartile = delay
		}
	}
	return row
}

func summarizeNestedRotationSplitSide(result FuturesHigherTFNestedRangeRotationAuditResult, cfg FuturesHigherTFNestedRangeRotationAuditConfig, split Split, side string) FuturesHigherTFNestedRangeRotationSummaryRow {
	row := FuturesHigherTFNestedRangeRotationSummaryRow{Split: split.Name, Side: side}
	for _, child := range result.ChildRangeRows {
		if side != "all" && child.CandidateSide != side {
			continue
		}
		childTime, err := parseTime(child.MatureCloseTime)
		if err != nil || !split.Contains(childTime) {
			continue
		}
		row.ChildCount++
		if !child.Eligible {
			row.SkippedChildCount++
		}
	}
	for _, event := range result.EventRows {
		if side != "all" && event.Side != side {
			continue
		}
		eventTime, err := parseTime(event.EventCloseTime)
		if err != nil || !split.Contains(eventTime) {
			continue
		}
		if event.SkippedReason != "" {
			row.SkippedEventCount++
			continue
		}
		if event.MissingFuture {
			row.MissingFutureCount++
			continue
		}
		row.EventCount++
		if event.FavorableMidpoint {
			row.FavorableMidpointCount++
		}
		if event.FavorableFarQuartile {
			row.FavorableFarQuartileCount++
		}
		if event.AdverseChildInvalidation {
			row.AdverseChildInvalidationCount++
		}
		if event.AdverseParentInvalidation {
			row.AdverseParentInvalidationCount++
		}
		if event.NoResolution {
			row.NoResolutionCount++
		}
		if event.QuickInvalidation {
			row.QuickInvalidationCount++
		}
		row.AvgFavorableExcursionParentWidth += event.FavorableExcursionParentWidth
		row.AvgAdverseExcursionParentWidth += event.AdverseExcursionParentWidth
	}
	if row.EventCount > 0 {
		denom := float64(row.EventCount)
		row.FavorableMidpointRate = float64(row.FavorableMidpointCount) / denom
		row.FavorableFarQuartileRate = float64(row.FavorableFarQuartileCount) / denom
		row.AdverseChildInvalidationRate = float64(row.AdverseChildInvalidationCount) / denom
		row.AdverseParentInvalidationRate = float64(row.AdverseParentInvalidationCount) / denom
		row.NoResolutionRate = float64(row.NoResolutionCount) / denom
		row.QuickInvalidationRate = float64(row.QuickInvalidationCount) / denom
		row.MissingFutureRate = float64(row.MissingFutureCount) / denom
		row.AvgFavorableExcursionParentWidth /= denom
		row.AvgAdverseExcursionParentWidth /= denom
	}
	return row
}

func markNestedRotationReviewGates(rows []FuturesHigherTFNestedRangeRotationSummaryRow, cfg FuturesHigherTFNestedRangeRotationAuditConfig, sourcePass bool, splits []Split) {
	byKey := map[string]*FuturesHigherTFNestedRangeRotationSummaryRow{}
	for i := range rows {
		key := rows[i].Split + "|" + rows[i].Side
		byKey[key] = &rows[i]
	}
	full := byKey[fullSplitName+"|all"]
	up := byKey[fullSplitName+"|"+RangeDiscoverySideUp]
	down := byKey[fullSplitName+"|"+RangeDiscoverySideDown]
	fullCountPass := full != nil && full.EventCount >= cfg.MinFullEvents
	sidePass := false
	caveat := ""
	if up != nil && down != nil {
		switch {
		case up.EventCount >= cfg.MinSideEvents && down.EventCount >= cfg.MinSideEvents:
			sidePass = true
		case up.EventCount >= cfg.MinSideEvents:
			sidePass = true
			caveat = "exclude_down_future_baseline"
		case down.EventCount >= cfg.MinSideEvents:
			sidePass = true
			caveat = "exclude_up_future_baseline"
		}
	}
	periodRows := []*FuturesHigherTFNestedRangeRotationSummaryRow{}
	for _, split := range rangeDiscoveryPeriodSplits(splits) {
		if row := byKey[split.Name+"|all"]; row != nil {
			periodRows = append(periodRows, row)
		}
	}
	splitCountPass := len(periodRows) == len(rangeDiscoveryPeriodSplits(splits))
	favAdvPass := splitCountPass
	favQuickPass := splitCountPass
	farPass := splitCountPass
	excursionPass := splitCountPass
	for _, row := range periodRows {
		if row.EventCount < cfg.MinSplitEvents {
			splitCountPass = false
		}
		if row.FavorableMidpointRate <= row.AdverseChildInvalidationRate {
			favAdvPass = false
		}
		if row.FavorableMidpointRate <= row.QuickInvalidationRate {
			favQuickPass = false
		}
		if row.FavorableFarQuartileRate < cfg.MinFarQuartileRate {
			farPass = false
		}
		if row.AvgFavorableExcursionParentWidth <= row.AvgAdverseExcursionParentWidth {
			excursionPass = false
		}
	}
	regimePass := splitCountPass
	globalPass := sourcePass && fullCountPass && splitCountPass && sidePass && favAdvPass && favQuickPass && farPass && excursionPass && regimePass
	for i := range rows {
		rows[i].SourceResamplePass = sourcePass
		rows[i].CountGatePass = fullCountPass
		if rows[i].Split != fullSplitName {
			rows[i].CountGatePass = rows[i].EventCount >= cfg.MinSplitEvents
		}
		rows[i].SideGatePass = sidePass
		rows[i].FavorableBeatsAdversePass = favAdvPass
		rows[i].FavorableBeatsQuickPass = favQuickPass
		rows[i].FarQuartileGatePass = farPass
		rows[i].ExcursionGatePass = excursionPass
		rows[i].RegimeSpreadPass = regimePass
		rows[i].PassesReviewGate = globalPass && rows[i].Split == fullSplitName && rows[i].Side == "all"
		rows[i].Caveat = caveat
		rows[i].FailureReason = nestedRotationFailureReason(sourcePass, fullCountPass, splitCountPass, sidePass, favAdvPass, favQuickPass, farPass, excursionPass, regimePass)
	}
}

func nestedRotationFailureReason(sourcePass, fullCountPass, splitCountPass, sidePass, favAdvPass, favQuickPass, farPass, excursionPass, regimePass bool) string {
	reasons := []string{}
	if !sourcePass {
		reasons = append(reasons, "source_or_resample_gap")
	}
	if !fullCountPass {
		reasons = append(reasons, "inadequate_full_events")
	}
	if !splitCountPass {
		reasons = append(reasons, "inadequate_split_events")
	}
	if !sidePass {
		reasons = append(reasons, "inadequate_side_events")
	}
	if !favAdvPass {
		reasons = append(reasons, "favorable_midpoint_not_above_adverse_child")
	}
	if !favQuickPass {
		reasons = append(reasons, "favorable_midpoint_not_above_quick_invalidation")
	}
	if !farPass {
		reasons = append(reasons, "far_quartile_below_gate")
	}
	if !excursionPass {
		reasons = append(reasons, "favorable_excursion_not_above_adverse")
	}
	if !regimePass {
		reasons = append(reasons, "regime_spread_failed")
	}
	return strings.Join(reasons, ";")
}
