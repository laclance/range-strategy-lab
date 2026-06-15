package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	HoldInsidePaperSideTowardHigh = "toward_high"
	HoldInsidePaperSideTowardLow  = "toward_low"

	holdInsideDecisionPositionBucketAll = "all"
	holdInsideDecisionMidSideBelow      = "below_mid"
	holdInsideDecisionMidSideAt         = "at_mid"
	holdInsideDecisionMidSideAbove      = "above_mid"
)

type HoldInsideDirectionalEdgeAuditConfig struct {
	HorizonsBars          []int
	QuickInvalidationBars int
	Profiles              []DetectorSweepProfile
	ContextRules          []DetectorContextRefinementRule
}

type HoldInsideDirectionalEdgeCandidateRow struct {
	ProfileID                   string  `json:"profile_id"`
	IsBalancedBaseline          bool    `json:"is_balanced_baseline"`
	IsADXComparison             bool    `json:"is_adx_comparison"`
	Percentile                  float64 `json:"percentile"`
	MinConsecutiveBars          int     `json:"min_consecutive_bars"`
	UseBollinger                bool    `json:"use_bollinger"`
	UseADX                      bool    `json:"use_adx"`
	LookbackDays                int     `json:"lookback_days"`
	ContextRule                 string  `json:"context_rule"`
	HoldBars                    int     `json:"hold_bars"`
	RequireMid50                bool    `json:"require_mid_50"`
	Split                       string  `json:"split"`
	PaperSide                   string  `json:"paper_side"`
	SourceEpisodeID             int     `json:"source_episode_id"`
	EpisodeStartIndex           int     `json:"episode_start_index"`
	EpisodeEndIndex             int     `json:"episode_end_index"`
	EpisodeStartTime            string  `json:"episode_start_time"`
	EpisodeEndTime              string  `json:"episode_end_time"`
	RawLengthBars               int     `json:"raw_length_bars"`
	ActiveLengthBars            int     `json:"active_length_bars"`
	RawLengthBucket             string  `json:"raw_length_bucket"`
	ActiveLengthBucket          string  `json:"active_length_bucket"`
	EpisodeHigh                 float64 `json:"episode_high"`
	EpisodeLow                  float64 `json:"episode_low"`
	EpisodeMid                  float64 `json:"episode_mid"`
	EpisodeEndClose             float64 `json:"episode_end_close"`
	EpisodeWidthPct             float64 `json:"episode_width_pct"`
	EpisodeWidthBucket          string  `json:"episode_width_bucket"`
	AvgNormalizedATR            float64 `json:"avg_normalized_atr"`
	EndNormalizedATR            float64 `json:"end_normalized_atr"`
	WidthToATRRatio             float64 `json:"width_to_atr_ratio"`
	WidthToATRBucket            string  `json:"width_to_atr_bucket"`
	DecisionIndex               int     `json:"decision_index"`
	DecisionTime                string  `json:"decision_time"`
	DecisionClose               float64 `json:"decision_close"`
	DecisionClosePosition       float64 `json:"decision_close_position"`
	DecisionClosePositionBucket string  `json:"decision_close_position_bucket"`
	DecisionMidSide             string  `json:"decision_mid_side"`
	DecisionDistanceToHighPct   float64 `json:"decision_distance_to_high_pct"`
	DecisionDistanceToLowPct    float64 `json:"decision_distance_to_low_pct"`
	DecisionDistanceToMidPct    float64 `json:"decision_distance_to_mid_pct"`
	HorizonBars                 int     `json:"horizon_bars"`
	LabelWindowStartIndex       int     `json:"label_window_start_index"`
	LabelWindowEndIndex         int     `json:"label_window_end_index"`
	LabelWindowStartTime        string  `json:"label_window_start_time"`
	LabelWindowEndTime          string  `json:"label_window_end_time"`
	LabelFavorableMovePct       float64 `json:"label_favorable_move_pct"`
	LabelAdverseMovePct         float64 `json:"label_adverse_move_pct"`
	LabelFavorableMinusAdverse  float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGTAdverse     bool    `json:"label_favorable_greater_than_adverse"`
	LabelTouchedMid             bool    `json:"label_touched_mid"`
	LabelClosedAcrossMid        bool    `json:"label_closed_across_mid"`
	LabelSideBoundaryTouch      bool    `json:"label_side_boundary_touch"`
	LabelOppositeCloseBreak     bool    `json:"label_opposite_close_break"`
	LabelReenteredRange         bool    `json:"label_reentered_range"`
	LabelPersistedInsideRange   bool    `json:"label_persisted_inside_range"`
	LabelQuickInvalidated       bool    `json:"label_quick_invalidated"`
	LabelInvalidatedUp          bool    `json:"label_invalidated_up"`
	LabelInvalidatedDown        bool    `json:"label_invalidated_down"`
	LabelTrendedUp              bool    `json:"label_trended_up"`
	LabelTrendedDown            bool    `json:"label_trended_down"`
}

type HoldInsideDirectionalEdgeSummaryRow struct {
	ProfileID                        string  `json:"profile_id"`
	IsBalancedBaseline               bool    `json:"is_balanced_baseline"`
	IsADXComparison                  bool    `json:"is_adx_comparison"`
	Percentile                       float64 `json:"percentile"`
	MinConsecutiveBars               int     `json:"min_consecutive_bars"`
	UseBollinger                     bool    `json:"use_bollinger"`
	UseADX                           bool    `json:"use_adx"`
	LookbackDays                     int     `json:"lookback_days"`
	ContextRule                      string  `json:"context_rule"`
	HoldBars                         int     `json:"hold_bars"`
	RequireMid50                     bool    `json:"require_mid_50"`
	Split                            string  `json:"split"`
	HorizonBars                      int     `json:"horizon_bars"`
	PaperSide                        string  `json:"paper_side"`
	DecisionClosePositionBucket      string  `json:"decision_close_position_bucket"`
	SourceEpisodeCount               int     `json:"source_episode_count"`
	CandidateCount                   int     `json:"candidate_count"`
	CandidateRate                    float64 `json:"candidate_rate"`
	AvgRawLengthBars                 float64 `json:"avg_raw_length_bars"`
	AvgActiveLengthBars              float64 `json:"avg_active_length_bars"`
	AvgEpisodeWidthPct               float64 `json:"avg_episode_width_pct"`
	AvgNormalizedATR                 float64 `json:"avg_normalized_atr"`
	AvgEndNormalizedATR              float64 `json:"avg_end_normalized_atr"`
	AvgWidthToATRRatio               float64 `json:"avg_width_to_atr_ratio"`
	AvgDecisionClosePosition         float64 `json:"avg_decision_close_position"`
	AvgDecisionDistanceToHighPct     float64 `json:"avg_decision_distance_to_high_pct"`
	AvgDecisionDistanceToLowPct      float64 `json:"avg_decision_distance_to_low_pct"`
	AvgDecisionDistanceToMidPct      float64 `json:"avg_decision_distance_to_mid_pct"`
	LabelFavorableGTAdverseCount     int     `json:"label_favorable_greater_than_adverse_count"`
	LabelTouchedMidCount             int     `json:"label_touched_mid_count"`
	LabelClosedAcrossMidCount        int     `json:"label_closed_across_mid_count"`
	LabelSideBoundaryTouchCount      int     `json:"label_side_boundary_touch_count"`
	LabelOppositeCloseBreakCount     int     `json:"label_opposite_close_break_count"`
	LabelQuickInvalidatedCount       int     `json:"label_quick_invalidated_count"`
	LabelInvalidatedUpCount          int     `json:"label_invalidated_up_count"`
	LabelInvalidatedDownCount        int     `json:"label_invalidated_down_count"`
	LabelTrendedUpCount              int     `json:"label_trended_up_count"`
	LabelTrendedDownCount            int     `json:"label_trended_down_count"`
	LabelAvgFavorableMovePct         float64 `json:"label_avg_favorable_move_pct"`
	LabelAvgAdverseMovePct           float64 `json:"label_avg_adverse_move_pct"`
	LabelAvgFavorableMinusAdversePct float64 `json:"label_avg_favorable_minus_adverse_pct"`
	LabelFavorableGTAdverseRate      float64 `json:"label_favorable_greater_than_adverse_rate"`
	LabelTouchedMidRate              float64 `json:"label_touched_mid_rate"`
	LabelClosedAcrossMidRate         float64 `json:"label_closed_across_mid_rate"`
	LabelSideBoundaryTouchRate       float64 `json:"label_side_boundary_touch_rate"`
	LabelOppositeCloseBreakRate      float64 `json:"label_opposite_close_break_rate"`
	LabelQuickInvalidatedRate        float64 `json:"label_quick_invalidated_rate"`
	LabelInvalidatedUpRate           float64 `json:"label_invalidated_up_rate"`
	LabelInvalidatedDownRate         float64 `json:"label_invalidated_down_rate"`
	LabelTrendedUpRate               float64 `json:"label_trended_up_rate"`
	LabelTrendedDownRate             float64 `json:"label_trended_down_rate"`
}

type HoldInsideDirectionalEdgeStabilityRow struct {
	ProfileID                             string  `json:"profile_id"`
	IsBalancedBaseline                    bool    `json:"is_balanced_baseline"`
	IsADXComparison                       bool    `json:"is_adx_comparison"`
	Percentile                            float64 `json:"percentile"`
	MinConsecutiveBars                    int     `json:"min_consecutive_bars"`
	UseBollinger                          bool    `json:"use_bollinger"`
	UseADX                                bool    `json:"use_adx"`
	LookbackDays                          int     `json:"lookback_days"`
	ContextRule                           string  `json:"context_rule"`
	HoldBars                              int     `json:"hold_bars"`
	RequireMid50                          bool    `json:"require_mid_50"`
	HorizonBars                           int     `json:"horizon_bars"`
	PaperSide                             string  `json:"paper_side"`
	DecisionClosePositionBucket           string  `json:"decision_close_position_bucket"`
	PeriodSplits                          int     `json:"period_splits"`
	SourceEpisodeCount                    int     `json:"source_episode_count"`
	SourceEpisodeCountMin                 int     `json:"source_episode_count_min"`
	SourceEpisodeCountMax                 int     `json:"source_episode_count_max"`
	SourceEpisodeCountDelta               int     `json:"source_episode_count_delta"`
	CandidateCount                        int     `json:"candidate_count"`
	CandidateCountMin                     int     `json:"candidate_count_min"`
	CandidateCountMax                     int     `json:"candidate_count_max"`
	CandidateCountDelta                   int     `json:"candidate_count_delta"`
	CandidateRateMin                      float64 `json:"candidate_rate_min"`
	CandidateRateMax                      float64 `json:"candidate_rate_max"`
	CandidateRateDelta                    float64 `json:"candidate_rate_delta"`
	LabelFavorableGTAdverseRateMin        float64 `json:"label_favorable_greater_than_adverse_rate_min"`
	LabelFavorableGTAdverseRateMax        float64 `json:"label_favorable_greater_than_adverse_rate_max"`
	LabelFavorableGTAdverseRateDelta      float64 `json:"label_favorable_greater_than_adverse_rate_delta"`
	LabelAvgFavorableMinusAdversePctMin   float64 `json:"label_avg_favorable_minus_adverse_pct_min"`
	LabelAvgFavorableMinusAdversePctMax   float64 `json:"label_avg_favorable_minus_adverse_pct_max"`
	LabelAvgFavorableMinusAdversePctDelta float64 `json:"label_avg_favorable_minus_adverse_pct_delta"`
	LabelAvgFavorableMovePctMin           float64 `json:"label_avg_favorable_move_pct_min"`
	LabelAvgFavorableMovePctMax           float64 `json:"label_avg_favorable_move_pct_max"`
	LabelAvgFavorableMovePctDelta         float64 `json:"label_avg_favorable_move_pct_delta"`
	LabelAvgAdverseMovePctMin             float64 `json:"label_avg_adverse_move_pct_min"`
	LabelAvgAdverseMovePctMax             float64 `json:"label_avg_adverse_move_pct_max"`
	LabelAvgAdverseMovePctDelta           float64 `json:"label_avg_adverse_move_pct_delta"`
	LabelTouchedMidRateMin                float64 `json:"label_touched_mid_rate_min"`
	LabelTouchedMidRateMax                float64 `json:"label_touched_mid_rate_max"`
	LabelTouchedMidRateDelta              float64 `json:"label_touched_mid_rate_delta"`
	LabelClosedAcrossMidRateMin           float64 `json:"label_closed_across_mid_rate_min"`
	LabelClosedAcrossMidRateMax           float64 `json:"label_closed_across_mid_rate_max"`
	LabelClosedAcrossMidRateDelta         float64 `json:"label_closed_across_mid_rate_delta"`
	LabelSideBoundaryTouchRateMin         float64 `json:"label_side_boundary_touch_rate_min"`
	LabelSideBoundaryTouchRateMax         float64 `json:"label_side_boundary_touch_rate_max"`
	LabelSideBoundaryTouchRateDelta       float64 `json:"label_side_boundary_touch_rate_delta"`
	LabelOppositeCloseBreakRateMin        float64 `json:"label_opposite_close_break_rate_min"`
	LabelOppositeCloseBreakRateMax        float64 `json:"label_opposite_close_break_rate_max"`
	LabelOppositeCloseBreakRateDelta      float64 `json:"label_opposite_close_break_rate_delta"`
	LabelQuickInvalidatedRateMin          float64 `json:"label_quick_invalidated_rate_min"`
	LabelQuickInvalidatedRateMax          float64 `json:"label_quick_invalidated_rate_max"`
	LabelQuickInvalidatedRateDelta        float64 `json:"label_quick_invalidated_rate_delta"`
}

func DefaultHoldInsideDirectionalEdgeAuditConfig() HoldInsideDirectionalEdgeAuditConfig {
	lookbackDays := DefaultCompressionRangeDetectorConfig().LookbackDays
	return HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1, 3, 6, 12},
		QuickInvalidationBars: 3,
		Profiles:              defaultHoldInsideDirectionalEdgeProfiles(lookbackDays),
		ContextRules:          defaultHoldInsideDirectionalEdgeRules(),
	}
}

func RunHoldInsideDirectionalEdgeAudit(candles []Candle, detectorCfg RangeDetectorConfig, cfg HoldInsideDirectionalEdgeAuditConfig, splits []Split) ([]HoldInsideDirectionalEdgeCandidateRow, []HoldInsideDirectionalEdgeSummaryRow, []HoldInsideDirectionalEdgeStabilityRow, error) {
	detectorCfg = detectorCfg.withDefaults()
	if err := detectorCfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	cfg = cfg.withDefaults(detectorCfg.LookbackDays)
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	lookbackBars := detectorCfg.LookbackDays * detectorCfg.BarsPerDay
	atr := NormalizedATR(candles, detectorCfg.ATRPeriod)
	donchian := DonchianWidth(candles, detectorCfg.DonchianPeriod)
	bollinger := BollingerWidth(candles, detectorCfg.BollingerPeriod)
	adx := ADX(candles, detectorCfg.ADXPeriod)
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)

	percentiles := detectorContextPercentiles(cfg.Profiles)
	atrThresholds := make(map[float64][]float64, len(percentiles))
	donchianThresholds := make(map[float64][]float64, len(percentiles))
	bollingerThresholds := make(map[float64][]float64, len(percentiles))
	adxThresholds := make(map[float64][]float64, len(percentiles))
	for _, percentile := range percentiles {
		atrThresholds[percentile] = rollingPriorPercentile(atr, lookbackBars, percentile)
		donchianThresholds[percentile] = rollingPriorPercentile(donchian, lookbackBars, percentile)
		bollingerThresholds[percentile] = rollingPriorPercentile(bollinger, lookbackBars, percentile)
		adxThresholds[percentile] = rollingPriorPercentile(adx, lookbackBars, percentile)
	}

	candidates := []HoldInsideDirectionalEdgeCandidateRow{}
	summaryAccumulators := map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator{}
	for _, profile := range cfg.Profiles {
		classifications := classifyDetectorSweepProfile(
			candles,
			profile,
			atr,
			donchian,
			bollinger,
			adx,
			atrThresholds[profile.Percentile],
			donchianThresholds[profile.Percentile],
			bollingerThresholds[profile.Percentile],
			adxThresholds[profile.Percentile],
		)
		profileCandidates := runHoldInsideDirectionalEdgeForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
		candidates = append(candidates, profileCandidates...)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideDirectionalEdgeCandidate(candidates[i], candidates[j])
	})
	summaryRows := holdInsideDirectionalEdgeSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideDirectionalEdgeStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func runHoldInsideDirectionalEdgeAuditFromClassifications(candles []Candle, profile DetectorSweepProfile, classifications []RangeClassification, cfg HoldInsideDirectionalEdgeAuditConfig, splits []Split) ([]HoldInsideDirectionalEdgeCandidateRow, []HoldInsideDirectionalEdgeSummaryRow, []HoldInsideDirectionalEdgeStabilityRow, error) {
	cfg = cfg.withDefaults(profile.LookbackDays)
	cfg.Profiles = []DetectorSweepProfile{profile}
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}

	normalizedATR := NormalizedATR(candles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
	summaryAccumulators := map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator{}
	candidates := runHoldInsideDirectionalEdgeForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideDirectionalEdgeCandidate(candidates[i], candidates[j])
	})
	summaryRows := holdInsideDirectionalEdgeSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideDirectionalEdgeStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func (cfg HoldInsideDirectionalEdgeAuditConfig) withDefaults(lookbackDays int) HoldInsideDirectionalEdgeAuditConfig {
	defaults := DefaultHoldInsideDirectionalEdgeAuditConfig()
	if lookbackDays > 0 {
		defaults.Profiles = defaultHoldInsideDirectionalEdgeProfiles(lookbackDays)
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickInvalidationBars == 0 {
		cfg.QuickInvalidationBars = defaults.QuickInvalidationBars
	}
	if len(cfg.Profiles) == 0 {
		cfg.Profiles = append([]DetectorSweepProfile(nil), defaults.Profiles...)
	}
	if len(cfg.ContextRules) == 0 {
		cfg.ContextRules = append([]DetectorContextRefinementRule(nil), defaults.ContextRules...)
	}
	return cfg
}

func (cfg HoldInsideDirectionalEdgeAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("hold-inside directional edge audit horizon bars must be positive")
		}
	}
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("hold-inside directional edge audit quick invalidation bars must be positive")
	}
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("hold-inside directional edge audit profiles must not be empty")
	}
	for _, profile := range cfg.Profiles {
		if profile.ProfileID == "" {
			return fmt.Errorf("hold-inside directional edge audit profile id must not be empty")
		}
	}
	if len(cfg.ContextRules) == 0 {
		return fmt.Errorf("hold-inside directional edge audit rules must not be empty")
	}
	for _, rule := range cfg.ContextRules {
		if rule.RuleID == "" {
			return fmt.Errorf("hold-inside directional edge audit rule id must not be empty")
		}
		if rule.HoldBars <= 0 {
			return fmt.Errorf("hold-inside directional edge audit rules must use positive hold bars")
		}
		if rule.RequireMid50 && rule.HoldBars == 0 {
			return fmt.Errorf("hold-inside directional edge audit mid-50 rule requires positive hold bars")
		}
	}
	return nil
}

func defaultHoldInsideDirectionalEdgeProfiles(lookbackDays int) []DetectorSweepProfile {
	if lookbackDays <= 0 {
		lookbackDays = DefaultCompressionRangeDetectorConfig().LookbackDays
	}
	return []DetectorSweepProfile{
		newDetectorSweepProfile(0.30, 12, true, false, lookbackDays, false),
	}
}

func defaultHoldInsideDirectionalEdgeRules() []DetectorContextRefinementRule {
	return []DetectorContextRefinementRule{
		{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
		{RuleID: DetectorContextRuleHold6Inside, HoldBars: 6},
		{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
	}
}

func runHoldInsideDirectionalEdgeForProfile(candles []Candle, classifications []RangeClassification, normalizedATR []float64, profile DetectorSweepProfile, cfg HoldInsideDirectionalEdgeAuditConfig, splits []Split, summaryAccumulators map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator) []HoldInsideDirectionalEdgeCandidateRow {
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, profile.ProfileID)
	candidates := []HoldInsideDirectionalEdgeCandidateRow{}
	for _, episode := range episodes {
		for _, rule := range cfg.ContextRules {
			decisionIndex := episode.EndIndex + rule.HoldBars
			if decisionIndex < 0 || decisionIndex >= len(candles) {
				continue
			}
			split := splitNameForCloseTime(candles[decisionIndex].CloseTime, splits)
			decisionBucket := decisionClosePositionBucket(decisionClosePosition(candles[decisionIndex].Close, episode.Low, episode.High))
			passesRule := detectorContextRulePasses(candles, episode, rule, decisionIndex)
			for _, horizon := range cfg.HorizonsBars {
				if decisionIndex+horizon >= len(candles) {
					continue
				}
				for _, paperSide := range holdInsidePaperSides() {
					addHoldInsideDirectionalEdgeSource(summaryAccumulators, profile, rule, split, horizon, paperSide, decisionBucket)
					addHoldInsideDirectionalEdgeSource(summaryAccumulators, profile, rule, split, horizon, paperSide, holdInsideDecisionPositionBucketAll)
					if split != fullSplitName {
						addHoldInsideDirectionalEdgeSource(summaryAccumulators, profile, rule, fullSplitName, horizon, paperSide, decisionBucket)
						addHoldInsideDirectionalEdgeSource(summaryAccumulators, profile, rule, fullSplitName, horizon, paperSide, holdInsideDecisionPositionBucketAll)
					}
					if !passesRule {
						continue
					}
					row, ok := newHoldInsideDirectionalEdgeCandidateRow(candles, profile, rule, episode, split, decisionIndex, horizon, cfg.QuickInvalidationBars, paperSide)
					if !ok {
						continue
					}
					candidates = append(candidates, row)
					addHoldInsideDirectionalEdgeCandidate(summaryAccumulators, row, split, row.DecisionClosePositionBucket)
					addHoldInsideDirectionalEdgeCandidate(summaryAccumulators, row, split, holdInsideDecisionPositionBucketAll)
					if split != fullSplitName {
						addHoldInsideDirectionalEdgeCandidate(summaryAccumulators, row, fullSplitName, row.DecisionClosePositionBucket)
						addHoldInsideDirectionalEdgeCandidate(summaryAccumulators, row, fullSplitName, holdInsideDecisionPositionBucketAll)
					}
				}
			}
		}
	}
	return candidates
}

func newHoldInsideDirectionalEdgeCandidateRow(candles []Candle, profile DetectorSweepProfile, rule DetectorContextRefinementRule, episode rangeRegimeDurabilityEpisode, split string, decisionIndex, horizon, quickInvalidationBars int, paperSide string) (HoldInsideDirectionalEdgeCandidateRow, bool) {
	if horizon <= 0 || quickInvalidationBars <= 0 || decisionIndex < 0 || decisionIndex+horizon >= len(candles) || !validHoldInsidePaperSide(paperSide) {
		return HoldInsideDirectionalEdgeCandidateRow{}, false
	}
	decision := candles[decisionIndex]
	episodeMid := (episode.High + episode.Low) / 2
	decisionPosition := decisionClosePosition(decision.Close, episode.Low, episode.High)
	decisionPositionValue := decisionPosition
	if !validNumber(decisionPositionValue) {
		decisionPositionValue = 0
	}
	startIndex := decisionIndex + 1
	endIndex := decisionIndex + horizon
	future := candles[startIndex : endIndex+1]
	futureMaxHigh, futureMinLow := futureHighLow(future)
	label := newHoldInsideDirectionalEdgeLabel(future, decision.Close, episode.Low, episode.High, episodeMid, paperSide, quickInvalidationBars)

	return HoldInsideDirectionalEdgeCandidateRow{
		ProfileID:                   profile.ProfileID,
		IsBalancedBaseline:          profile.IsBalancedBaseline,
		IsADXComparison:             profile.IsADXComparison,
		Percentile:                  profile.Percentile,
		MinConsecutiveBars:          profile.MinConsecutiveBars,
		UseBollinger:                profile.UseBollinger,
		UseADX:                      profile.UseADX,
		LookbackDays:                profile.LookbackDays,
		ContextRule:                 rule.RuleID,
		HoldBars:                    rule.HoldBars,
		RequireMid50:                rule.RequireMid50,
		Split:                       split,
		PaperSide:                   paperSide,
		SourceEpisodeID:             episode.EpisodeID,
		EpisodeStartIndex:           episode.StartIndex,
		EpisodeEndIndex:             episode.EndIndex,
		EpisodeStartTime:            candles[episode.StartIndex].CloseTime.Format(timeLayout),
		EpisodeEndTime:              candles[episode.EndIndex].CloseTime.Format(timeLayout),
		RawLengthBars:               episode.RawLengthBars,
		ActiveLengthBars:            episode.ActiveLengthBars,
		RawLengthBucket:             episode.RawLengthBucket,
		ActiveLengthBucket:          episode.ActiveLengthBucket,
		EpisodeHigh:                 episode.High,
		EpisodeLow:                  episode.Low,
		EpisodeMid:                  episodeMid,
		EpisodeEndClose:             episode.EndClose,
		EpisodeWidthPct:             episode.WidthPct,
		EpisodeWidthBucket:          episode.WidthBucket,
		AvgNormalizedATR:            episode.AvgNormalizedATR,
		EndNormalizedATR:            episode.EndNormalizedATR,
		WidthToATRRatio:             episode.WidthToATRRatio,
		WidthToATRBucket:            episode.WidthToATRBucket,
		DecisionIndex:               decisionIndex,
		DecisionTime:                decision.CloseTime.Format(timeLayout),
		DecisionClose:               decision.Close,
		DecisionClosePosition:       decisionPositionValue,
		DecisionClosePositionBucket: decisionClosePositionBucket(decisionPosition),
		DecisionMidSide:             holdInsideDecisionMidSide(decision.Close, episodeMid),
		DecisionDistanceToHighPct:   movePct(math.Max(0, episode.High-decision.Close), decision.Close),
		DecisionDistanceToLowPct:    movePct(math.Max(0, decision.Close-episode.Low), decision.Close),
		DecisionDistanceToMidPct:    movePct(math.Abs(decision.Close-episodeMid), decision.Close),
		HorizonBars:                 horizon,
		LabelWindowStartIndex:       startIndex,
		LabelWindowEndIndex:         endIndex,
		LabelWindowStartTime:        candles[startIndex].CloseTime.Format(timeLayout),
		LabelWindowEndTime:          candles[endIndex].CloseTime.Format(timeLayout),
		LabelFavorableMovePct:       label.FavorableMovePct,
		LabelAdverseMovePct:         label.AdverseMovePct,
		LabelFavorableMinusAdverse:  label.FavorableMovePct - label.AdverseMovePct,
		LabelFavorableGTAdverse:     label.FavorableMovePct > label.AdverseMovePct,
		LabelTouchedMid:             label.TouchedMid,
		LabelClosedAcrossMid:        label.ClosedAcrossMid,
		LabelSideBoundaryTouch:      holdInsideSideBoundaryTouch(paperSide, futureMaxHigh, futureMinLow, episode.High, episode.Low),
		LabelOppositeCloseBreak:     label.OppositeCloseBreak,
		LabelReenteredRange:         label.ReenteredRange,
		LabelPersistedInsideRange:   label.PersistedInsideRange,
		LabelQuickInvalidated:       label.QuickInvalidated,
		LabelInvalidatedUp:          label.InvalidatedUp,
		LabelInvalidatedDown:        label.InvalidatedDown,
		LabelTrendedUp:              label.TrendedUp,
		LabelTrendedDown:            label.TrendedDown,
	}, true
}

type holdInsideDirectionalEdgeLabel struct {
	FavorableMovePct     float64
	AdverseMovePct       float64
	TouchedMid           bool
	ClosedAcrossMid      bool
	OppositeCloseBreak   bool
	ReenteredRange       bool
	PersistedInsideRange bool
	QuickInvalidated     bool
	InvalidatedUp        bool
	InvalidatedDown      bool
	TrendedUp            bool
	TrendedDown          bool
}

func newHoldInsideDirectionalEdgeLabel(future []Candle, decisionClose, low, high, mid float64, paperSide string, quickInvalidationBars int) holdInsideDirectionalEdgeLabel {
	futureMaxHigh, futureMinLow := futureHighLow(future)
	lastClose := future[len(future)-1].Close
	maxUpMovePct := movePct(math.Max(0, futureMaxHigh-decisionClose), decisionClose)
	maxDownMovePct := movePct(math.Max(0, decisionClose-futureMinLow), decisionClose)
	label := holdInsideDirectionalEdgeLabel{PersistedInsideRange: true}
	switch paperSide {
	case HoldInsidePaperSideTowardHigh:
		label.FavorableMovePct = maxUpMovePct
		label.AdverseMovePct = maxDownMovePct
	case HoldInsidePaperSideTowardLow:
		label.FavorableMovePct = maxDownMovePct
		label.AdverseMovePct = maxUpMovePct
	}

	quickLimit := minInt(len(future), quickInvalidationBars)
	for i, candle := range future {
		if closeInsideRange(candle.Close, low, high) {
			label.ReenteredRange = true
		} else {
			label.PersistedInsideRange = false
		}
		if candle.Close > high {
			label.InvalidatedUp = true
			if i < quickLimit {
				label.QuickInvalidated = true
			}
		}
		if candle.Close < low {
			label.InvalidatedDown = true
			if i < quickLimit {
				label.QuickInvalidated = true
			}
		}
		switch paperSide {
		case HoldInsidePaperSideTowardHigh:
			if candle.Close < low {
				label.OppositeCloseBreak = true
			}
		case HoldInsidePaperSideTowardLow:
			if candle.Close > high {
				label.OppositeCloseBreak = true
			}
		}
		touchedMid, closedAcrossMid := holdInsideMidLabels(candle, decisionClose, mid)
		label.TouchedMid = label.TouchedMid || touchedMid
		label.ClosedAcrossMid = label.ClosedAcrossMid || closedAcrossMid
	}
	label.TrendedUp = label.InvalidatedUp && lastClose > high && maxUpMovePct > maxDownMovePct
	label.TrendedDown = label.InvalidatedDown && lastClose < low && maxDownMovePct > maxUpMovePct
	return label
}

func holdInsidePaperSides() []string {
	return []string{HoldInsidePaperSideTowardHigh, HoldInsidePaperSideTowardLow}
}

func validHoldInsidePaperSide(side string) bool {
	return side == HoldInsidePaperSideTowardHigh || side == HoldInsidePaperSideTowardLow
}

func holdInsideDecisionMidSide(close, mid float64) string {
	switch {
	case close < mid:
		return holdInsideDecisionMidSideBelow
	case close > mid:
		return holdInsideDecisionMidSideAbove
	default:
		return holdInsideDecisionMidSideAt
	}
}

func holdInsideMidLabels(candle Candle, decisionClose, mid float64) (bool, bool) {
	switch holdInsideDecisionMidSide(decisionClose, mid) {
	case holdInsideDecisionMidSideBelow:
		return candle.High >= mid, candle.Close >= mid
	case holdInsideDecisionMidSideAbove:
		return candle.Low <= mid, candle.Close <= mid
	default:
		return candle.Low <= mid && candle.High >= mid, candle.Close == mid
	}
}

func holdInsideSideBoundaryTouch(paperSide string, futureMaxHigh, futureMinLow, high, low float64) bool {
	switch paperSide {
	case HoldInsidePaperSideTowardHigh:
		return futureMaxHigh >= high
	case HoldInsidePaperSideTowardLow:
		return futureMinLow <= low
	default:
		return false
	}
}

type holdInsideDirectionalEdgeSummaryKey struct {
	profileID                   string
	contextRule                 string
	split                       string
	horizonBars                 int
	paperSide                   string
	decisionClosePositionBucket string
}

type holdInsideDirectionalEdgeSummaryAccumulator struct {
	profile                      DetectorSweepProfile
	rule                         DetectorContextRefinementRule
	split                        string
	horizonBars                  int
	paperSide                    string
	decisionClosePositionBucket  string
	sourceEpisodes               int
	candidates                   int
	rawLengthSum                 float64
	activeLengthSum              float64
	widthPctSum                  float64
	avgATRSum                    float64
	endATRSum                    float64
	widthToATRRatioSum           float64
	decisionPositionSum          float64
	decisionDistanceToHighPctSum float64
	decisionDistanceToLowPctSum  float64
	decisionDistanceToMidPctSum  float64
	favorableGTAdverse           int
	touchedMid                   int
	closedAcrossMid              int
	sideBoundaryTouch            int
	oppositeCloseBreak           int
	quickInvalidated             int
	invalidatedUp                int
	invalidatedDown              int
	trendedUp                    int
	trendedDown                  int
	favorableMovePctSum          float64
	adverseMovePctSum            float64
	favorableMinusAdverseSum     float64
}

func addHoldInsideDirectionalEdgeSource(accumulators map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int, paperSide, decisionBucket string) {
	acc := holdInsideDirectionalEdgeSummaryAccumulatorFor(accumulators, profile, rule, split, horizon, paperSide, decisionBucket)
	acc.sourceEpisodes++
}

func addHoldInsideDirectionalEdgeCandidate(accumulators map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator, row HoldInsideDirectionalEdgeCandidateRow, split, decisionBucket string) {
	profile := DetectorSweepProfile{
		ProfileID:          row.ProfileID,
		IsBalancedBaseline: row.IsBalancedBaseline,
		IsADXComparison:    row.IsADXComparison,
		Percentile:         row.Percentile,
		MinConsecutiveBars: row.MinConsecutiveBars,
		UseBollinger:       row.UseBollinger,
		UseADX:             row.UseADX,
		LookbackDays:       row.LookbackDays,
	}
	rule := DetectorContextRefinementRule{
		RuleID:       row.ContextRule,
		HoldBars:     row.HoldBars,
		RequireMid50: row.RequireMid50,
	}
	acc := holdInsideDirectionalEdgeSummaryAccumulatorFor(accumulators, profile, rule, split, row.HorizonBars, row.PaperSide, decisionBucket)
	acc.candidates++
	acc.rawLengthSum += float64(row.RawLengthBars)
	acc.activeLengthSum += float64(row.ActiveLengthBars)
	acc.widthPctSum += row.EpisodeWidthPct
	acc.avgATRSum += row.AvgNormalizedATR
	acc.endATRSum += row.EndNormalizedATR
	acc.widthToATRRatioSum += row.WidthToATRRatio
	if row.DecisionClosePositionBucket != decisionClosePositionBucketUnknown {
		acc.decisionPositionSum += row.DecisionClosePosition
	}
	acc.decisionDistanceToHighPctSum += row.DecisionDistanceToHighPct
	acc.decisionDistanceToLowPctSum += row.DecisionDistanceToLowPct
	acc.decisionDistanceToMidPctSum += row.DecisionDistanceToMidPct
	if row.LabelFavorableGTAdverse {
		acc.favorableGTAdverse++
	}
	if row.LabelTouchedMid {
		acc.touchedMid++
	}
	if row.LabelClosedAcrossMid {
		acc.closedAcrossMid++
	}
	if row.LabelSideBoundaryTouch {
		acc.sideBoundaryTouch++
	}
	if row.LabelOppositeCloseBreak {
		acc.oppositeCloseBreak++
	}
	if row.LabelQuickInvalidated {
		acc.quickInvalidated++
	}
	if row.LabelInvalidatedUp {
		acc.invalidatedUp++
	}
	if row.LabelInvalidatedDown {
		acc.invalidatedDown++
	}
	if row.LabelTrendedUp {
		acc.trendedUp++
	}
	if row.LabelTrendedDown {
		acc.trendedDown++
	}
	acc.favorableMovePctSum += row.LabelFavorableMovePct
	acc.adverseMovePctSum += row.LabelAdverseMovePct
	acc.favorableMinusAdverseSum += row.LabelFavorableMinusAdverse
}

func holdInsideDirectionalEdgeSummaryAccumulatorFor(accumulators map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int, paperSide, decisionBucket string) *holdInsideDirectionalEdgeSummaryAccumulator {
	key := holdInsideDirectionalEdgeSummaryKey{
		profileID:                   profile.ProfileID,
		contextRule:                 rule.RuleID,
		split:                       split,
		horizonBars:                 horizon,
		paperSide:                   paperSide,
		decisionClosePositionBucket: decisionBucket,
	}
	acc := accumulators[key]
	if acc == nil {
		acc = &holdInsideDirectionalEdgeSummaryAccumulator{
			profile:                     profile,
			rule:                        rule,
			split:                       split,
			horizonBars:                 horizon,
			paperSide:                   paperSide,
			decisionClosePositionBucket: decisionBucket,
		}
		accumulators[key] = acc
	}
	return acc
}

func holdInsideDirectionalEdgeSummaryRows(accumulators map[holdInsideDirectionalEdgeSummaryKey]*holdInsideDirectionalEdgeSummaryAccumulator) []HoldInsideDirectionalEdgeSummaryRow {
	rows := make([]HoldInsideDirectionalEdgeSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideDirectionalEdgeSummary(rows[i], rows[j])
	})
	return rows
}

func (acc holdInsideDirectionalEdgeSummaryAccumulator) row() HoldInsideDirectionalEdgeSummaryRow {
	row := HoldInsideDirectionalEdgeSummaryRow{
		ProfileID:                    acc.profile.ProfileID,
		IsBalancedBaseline:           acc.profile.IsBalancedBaseline,
		IsADXComparison:              acc.profile.IsADXComparison,
		Percentile:                   acc.profile.Percentile,
		MinConsecutiveBars:           acc.profile.MinConsecutiveBars,
		UseBollinger:                 acc.profile.UseBollinger,
		UseADX:                       acc.profile.UseADX,
		LookbackDays:                 acc.profile.LookbackDays,
		ContextRule:                  acc.rule.RuleID,
		HoldBars:                     acc.rule.HoldBars,
		RequireMid50:                 acc.rule.RequireMid50,
		Split:                        acc.split,
		HorizonBars:                  acc.horizonBars,
		PaperSide:                    acc.paperSide,
		DecisionClosePositionBucket:  acc.decisionClosePositionBucket,
		SourceEpisodeCount:           acc.sourceEpisodes,
		CandidateCount:               acc.candidates,
		LabelFavorableGTAdverseCount: acc.favorableGTAdverse,
		LabelTouchedMidCount:         acc.touchedMid,
		LabelClosedAcrossMidCount:    acc.closedAcrossMid,
		LabelSideBoundaryTouchCount:  acc.sideBoundaryTouch,
		LabelOppositeCloseBreakCount: acc.oppositeCloseBreak,
		LabelQuickInvalidatedCount:   acc.quickInvalidated,
		LabelInvalidatedUpCount:      acc.invalidatedUp,
		LabelInvalidatedDownCount:    acc.invalidatedDown,
		LabelTrendedUpCount:          acc.trendedUp,
		LabelTrendedDownCount:        acc.trendedDown,
	}
	if acc.sourceEpisodes > 0 {
		row.CandidateRate = float64(acc.candidates) / float64(acc.sourceEpisodes)
	}
	if acc.candidates == 0 {
		return row
	}
	count := float64(acc.candidates)
	row.AvgRawLengthBars = acc.rawLengthSum / count
	row.AvgActiveLengthBars = acc.activeLengthSum / count
	row.AvgEpisodeWidthPct = acc.widthPctSum / count
	row.AvgNormalizedATR = acc.avgATRSum / count
	row.AvgEndNormalizedATR = acc.endATRSum / count
	row.AvgWidthToATRRatio = acc.widthToATRRatioSum / count
	row.AvgDecisionClosePosition = acc.decisionPositionSum / count
	row.AvgDecisionDistanceToHighPct = acc.decisionDistanceToHighPctSum / count
	row.AvgDecisionDistanceToLowPct = acc.decisionDistanceToLowPctSum / count
	row.AvgDecisionDistanceToMidPct = acc.decisionDistanceToMidPctSum / count
	row.LabelAvgFavorableMovePct = acc.favorableMovePctSum / count
	row.LabelAvgAdverseMovePct = acc.adverseMovePctSum / count
	row.LabelAvgFavorableMinusAdversePct = acc.favorableMinusAdverseSum / count
	row.LabelFavorableGTAdverseRate = float64(acc.favorableGTAdverse) / count
	row.LabelTouchedMidRate = float64(acc.touchedMid) / count
	row.LabelClosedAcrossMidRate = float64(acc.closedAcrossMid) / count
	row.LabelSideBoundaryTouchRate = float64(acc.sideBoundaryTouch) / count
	row.LabelOppositeCloseBreakRate = float64(acc.oppositeCloseBreak) / count
	row.LabelQuickInvalidatedRate = float64(acc.quickInvalidated) / count
	row.LabelInvalidatedUpRate = float64(acc.invalidatedUp) / count
	row.LabelInvalidatedDownRate = float64(acc.invalidatedDown) / count
	row.LabelTrendedUpRate = float64(acc.trendedUp) / count
	row.LabelTrendedDownRate = float64(acc.trendedDown) / count
	return row
}

func holdInsideDirectionalEdgeStabilityRows(profiles []DetectorSweepProfile, rules []DetectorContextRefinementRule, horizons []int, summaryRows []HoldInsideDirectionalEdgeSummaryRow, splits []Split) []HoldInsideDirectionalEdgeStabilityRow {
	periodSplits := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplits = append(periodSplits, split)
		}
	}
	byKey := map[holdInsideDirectionalEdgeSummaryKey]HoldInsideDirectionalEdgeSummaryRow{}
	bucketsByBase := map[holdInsideDirectionalEdgeStabilityBaseKey]map[string]bool{}
	for _, row := range summaryRows {
		key := holdInsideDirectionalEdgeSummaryKey{
			profileID:                   row.ProfileID,
			contextRule:                 row.ContextRule,
			split:                       row.Split,
			horizonBars:                 row.HorizonBars,
			paperSide:                   row.PaperSide,
			decisionClosePositionBucket: row.DecisionClosePositionBucket,
		}
		byKey[key] = row
		base := holdInsideDirectionalEdgeStabilityBaseKey{
			profileID:   row.ProfileID,
			contextRule: row.ContextRule,
			horizonBars: row.HorizonBars,
			paperSide:   row.PaperSide,
		}
		if bucketsByBase[base] == nil {
			bucketsByBase[base] = map[string]bool{}
		}
		bucketsByBase[base][row.DecisionClosePositionBucket] = true
	}

	rows := []HoldInsideDirectionalEdgeStabilityRow{}
	for _, profile := range profiles {
		for _, rule := range rules {
			for _, horizon := range horizons {
				for _, paperSide := range holdInsidePaperSides() {
					base := holdInsideDirectionalEdgeStabilityBaseKey{
						profileID:   profile.ProfileID,
						contextRule: rule.RuleID,
						horizonBars: horizon,
						paperSide:   paperSide,
					}
					buckets := holdInsideDirectionalEdgeSortedBuckets(bucketsByBase[base])
					for _, bucket := range buckets {
						rows = append(rows, newHoldInsideDirectionalEdgeStabilityRow(profile, rule, horizon, paperSide, bucket, periodSplits, byKey))
					}
				}
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideDirectionalEdgeStability(rows[i], rows[j])
	})
	return rows
}

type holdInsideDirectionalEdgeStabilityBaseKey struct {
	profileID   string
	contextRule string
	horizonBars int
	paperSide   string
}

func newHoldInsideDirectionalEdgeStabilityRow(profile DetectorSweepProfile, rule DetectorContextRefinementRule, horizon int, paperSide, decisionBucket string, periodSplits []Split, byKey map[holdInsideDirectionalEdgeSummaryKey]HoldInsideDirectionalEdgeSummaryRow) HoldInsideDirectionalEdgeStabilityRow {
	row := HoldInsideDirectionalEdgeStabilityRow{
		ProfileID:                   profile.ProfileID,
		IsBalancedBaseline:          profile.IsBalancedBaseline,
		IsADXComparison:             profile.IsADXComparison,
		Percentile:                  profile.Percentile,
		MinConsecutiveBars:          profile.MinConsecutiveBars,
		UseBollinger:                profile.UseBollinger,
		UseADX:                      profile.UseADX,
		LookbackDays:                profile.LookbackDays,
		ContextRule:                 rule.RuleID,
		HoldBars:                    rule.HoldBars,
		RequireMid50:                rule.RequireMid50,
		HorizonBars:                 horizon,
		PaperSide:                   paperSide,
		DecisionClosePositionBucket: decisionBucket,
		PeriodSplits:                len(periodSplits),
	}
	first := true
	for _, split := range periodSplits {
		summary := byKey[holdInsideDirectionalEdgeSummaryKey{
			profileID:                   profile.ProfileID,
			contextRule:                 rule.RuleID,
			split:                       split.Name,
			horizonBars:                 horizon,
			paperSide:                   paperSide,
			decisionClosePositionBucket: decisionBucket,
		}]
		if first {
			row.SourceEpisodeCountMin = summary.SourceEpisodeCount
			row.SourceEpisodeCountMax = summary.SourceEpisodeCount
			row.CandidateCountMin = summary.CandidateCount
			row.CandidateCountMax = summary.CandidateCount
			row.CandidateRateMin = summary.CandidateRate
			row.CandidateRateMax = summary.CandidateRate
			row.LabelFavorableGTAdverseRateMin = summary.LabelFavorableGTAdverseRate
			row.LabelFavorableGTAdverseRateMax = summary.LabelFavorableGTAdverseRate
			row.LabelAvgFavorableMinusAdversePctMin = summary.LabelAvgFavorableMinusAdversePct
			row.LabelAvgFavorableMinusAdversePctMax = summary.LabelAvgFavorableMinusAdversePct
			row.LabelAvgFavorableMovePctMin = summary.LabelAvgFavorableMovePct
			row.LabelAvgFavorableMovePctMax = summary.LabelAvgFavorableMovePct
			row.LabelAvgAdverseMovePctMin = summary.LabelAvgAdverseMovePct
			row.LabelAvgAdverseMovePctMax = summary.LabelAvgAdverseMovePct
			row.LabelTouchedMidRateMin = summary.LabelTouchedMidRate
			row.LabelTouchedMidRateMax = summary.LabelTouchedMidRate
			row.LabelClosedAcrossMidRateMin = summary.LabelClosedAcrossMidRate
			row.LabelClosedAcrossMidRateMax = summary.LabelClosedAcrossMidRate
			row.LabelSideBoundaryTouchRateMin = summary.LabelSideBoundaryTouchRate
			row.LabelSideBoundaryTouchRateMax = summary.LabelSideBoundaryTouchRate
			row.LabelOppositeCloseBreakRateMin = summary.LabelOppositeCloseBreakRate
			row.LabelOppositeCloseBreakRateMax = summary.LabelOppositeCloseBreakRate
			row.LabelQuickInvalidatedRateMin = summary.LabelQuickInvalidatedRate
			row.LabelQuickInvalidatedRateMax = summary.LabelQuickInvalidatedRate
			first = false
		} else {
			row.SourceEpisodeCountMin = minInt(row.SourceEpisodeCountMin, summary.SourceEpisodeCount)
			row.SourceEpisodeCountMax = maxInt(row.SourceEpisodeCountMax, summary.SourceEpisodeCount)
			row.CandidateCountMin = minInt(row.CandidateCountMin, summary.CandidateCount)
			row.CandidateCountMax = maxInt(row.CandidateCountMax, summary.CandidateCount)
			row.CandidateRateMin = minFloat(row.CandidateRateMin, summary.CandidateRate)
			row.CandidateRateMax = maxFloat(row.CandidateRateMax, summary.CandidateRate)
			row.LabelFavorableGTAdverseRateMin = minFloat(row.LabelFavorableGTAdverseRateMin, summary.LabelFavorableGTAdverseRate)
			row.LabelFavorableGTAdverseRateMax = maxFloat(row.LabelFavorableGTAdverseRateMax, summary.LabelFavorableGTAdverseRate)
			row.LabelAvgFavorableMinusAdversePctMin = minFloat(row.LabelAvgFavorableMinusAdversePctMin, summary.LabelAvgFavorableMinusAdversePct)
			row.LabelAvgFavorableMinusAdversePctMax = maxFloat(row.LabelAvgFavorableMinusAdversePctMax, summary.LabelAvgFavorableMinusAdversePct)
			row.LabelAvgFavorableMovePctMin = minFloat(row.LabelAvgFavorableMovePctMin, summary.LabelAvgFavorableMovePct)
			row.LabelAvgFavorableMovePctMax = maxFloat(row.LabelAvgFavorableMovePctMax, summary.LabelAvgFavorableMovePct)
			row.LabelAvgAdverseMovePctMin = minFloat(row.LabelAvgAdverseMovePctMin, summary.LabelAvgAdverseMovePct)
			row.LabelAvgAdverseMovePctMax = maxFloat(row.LabelAvgAdverseMovePctMax, summary.LabelAvgAdverseMovePct)
			row.LabelTouchedMidRateMin = minFloat(row.LabelTouchedMidRateMin, summary.LabelTouchedMidRate)
			row.LabelTouchedMidRateMax = maxFloat(row.LabelTouchedMidRateMax, summary.LabelTouchedMidRate)
			row.LabelClosedAcrossMidRateMin = minFloat(row.LabelClosedAcrossMidRateMin, summary.LabelClosedAcrossMidRate)
			row.LabelClosedAcrossMidRateMax = maxFloat(row.LabelClosedAcrossMidRateMax, summary.LabelClosedAcrossMidRate)
			row.LabelSideBoundaryTouchRateMin = minFloat(row.LabelSideBoundaryTouchRateMin, summary.LabelSideBoundaryTouchRate)
			row.LabelSideBoundaryTouchRateMax = maxFloat(row.LabelSideBoundaryTouchRateMax, summary.LabelSideBoundaryTouchRate)
			row.LabelOppositeCloseBreakRateMin = minFloat(row.LabelOppositeCloseBreakRateMin, summary.LabelOppositeCloseBreakRate)
			row.LabelOppositeCloseBreakRateMax = maxFloat(row.LabelOppositeCloseBreakRateMax, summary.LabelOppositeCloseBreakRate)
			row.LabelQuickInvalidatedRateMin = minFloat(row.LabelQuickInvalidatedRateMin, summary.LabelQuickInvalidatedRate)
			row.LabelQuickInvalidatedRateMax = maxFloat(row.LabelQuickInvalidatedRateMax, summary.LabelQuickInvalidatedRate)
		}
		row.SourceEpisodeCount += summary.SourceEpisodeCount
		row.CandidateCount += summary.CandidateCount
	}
	row.SourceEpisodeCountDelta = row.SourceEpisodeCountMax - row.SourceEpisodeCountMin
	row.CandidateCountDelta = row.CandidateCountMax - row.CandidateCountMin
	row.CandidateRateDelta = row.CandidateRateMax - row.CandidateRateMin
	row.LabelFavorableGTAdverseRateDelta = row.LabelFavorableGTAdverseRateMax - row.LabelFavorableGTAdverseRateMin
	row.LabelAvgFavorableMinusAdversePctDelta = row.LabelAvgFavorableMinusAdversePctMax - row.LabelAvgFavorableMinusAdversePctMin
	row.LabelAvgFavorableMovePctDelta = row.LabelAvgFavorableMovePctMax - row.LabelAvgFavorableMovePctMin
	row.LabelAvgAdverseMovePctDelta = row.LabelAvgAdverseMovePctMax - row.LabelAvgAdverseMovePctMin
	row.LabelTouchedMidRateDelta = row.LabelTouchedMidRateMax - row.LabelTouchedMidRateMin
	row.LabelClosedAcrossMidRateDelta = row.LabelClosedAcrossMidRateMax - row.LabelClosedAcrossMidRateMin
	row.LabelSideBoundaryTouchRateDelta = row.LabelSideBoundaryTouchRateMax - row.LabelSideBoundaryTouchRateMin
	row.LabelOppositeCloseBreakRateDelta = row.LabelOppositeCloseBreakRateMax - row.LabelOppositeCloseBreakRateMin
	row.LabelQuickInvalidatedRateDelta = row.LabelQuickInvalidatedRateMax - row.LabelQuickInvalidatedRateMin
	return row
}

func holdInsideDirectionalEdgeSortedBuckets(bucketSet map[string]bool) []string {
	if len(bucketSet) == 0 {
		return nil
	}
	buckets := make([]string, 0, len(bucketSet))
	for bucket := range bucketSet {
		buckets = append(buckets, bucket)
	}
	sort.Slice(buckets, func(i, j int) bool {
		return holdInsideDecisionPositionBucketSortKey(buckets[i]) < holdInsideDecisionPositionBucketSortKey(buckets[j])
	})
	return buckets
}

func lessHoldInsideDirectionalEdgeCandidate(a, b HoldInsideDirectionalEdgeCandidateRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.SourceEpisodeID != b.SourceEpisodeID {
		return a.SourceEpisodeID < b.SourceEpisodeID
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	return holdInsidePaperSideSortKey(a.PaperSide) < holdInsidePaperSideSortKey(b.PaperSide)
}

func lessHoldInsideDirectionalEdgeSummary(a, b HoldInsideDirectionalEdgeSummaryRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if holdInsidePaperSideSortKey(a.PaperSide) != holdInsidePaperSideSortKey(b.PaperSide) {
		return holdInsidePaperSideSortKey(a.PaperSide) < holdInsidePaperSideSortKey(b.PaperSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.DecisionClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.DecisionClosePositionBucket)
}

func lessHoldInsideDirectionalEdgeStability(a, b HoldInsideDirectionalEdgeStabilityRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if holdInsidePaperSideSortKey(a.PaperSide) != holdInsidePaperSideSortKey(b.PaperSide) {
		return holdInsidePaperSideSortKey(a.PaperSide) < holdInsidePaperSideSortKey(b.PaperSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.DecisionClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.DecisionClosePositionBucket)
}

func holdInsidePaperSideSortKey(side string) int {
	switch side {
	case HoldInsidePaperSideTowardHigh:
		return 0
	case HoldInsidePaperSideTowardLow:
		return 1
	default:
		return 99
	}
}

func holdInsideDecisionPositionBucketSortKey(bucket string) int {
	switch bucket {
	case holdInsideDecisionPositionBucketAll:
		return 0
	case decisionClosePositionBucketBelow:
		return 1
	case decisionClosePositionBucketLow25:
		return 2
	case decisionClosePositionBucketMid50:
		return 3
	case decisionClosePositionBucketHigh25:
		return 4
	case decisionClosePositionBucketAbove:
		return 5
	case decisionClosePositionBucketUnknown:
		return 6
	default:
		return 99
	}
}
