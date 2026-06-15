package lab

import (
	"fmt"
	"math"
	"sort"
)

const holdInsideDecisionMidSideAll = "all"

type HoldInsideMidlineTransitionAuditConfig struct {
	HorizonsBars          []int
	QuickInvalidationBars int
	Profiles              []DetectorSweepProfile
	ContextRules          []DetectorContextRefinementRule
}

type HoldInsideMidlineTransitionCandidateRow struct {
	ProfileID                         string  `json:"profile_id"`
	IsBalancedBaseline                bool    `json:"is_balanced_baseline"`
	IsADXComparison                   bool    `json:"is_adx_comparison"`
	Percentile                        float64 `json:"percentile"`
	MinConsecutiveBars                int     `json:"min_consecutive_bars"`
	UseBollinger                      bool    `json:"use_bollinger"`
	UseADX                            bool    `json:"use_adx"`
	LookbackDays                      int     `json:"lookback_days"`
	ContextRule                       string  `json:"context_rule"`
	HoldBars                          int     `json:"hold_bars"`
	RequireMid50                      bool    `json:"require_mid_50"`
	Split                             string  `json:"split"`
	SourceEpisodeID                   int     `json:"source_episode_id"`
	EpisodeStartIndex                 int     `json:"episode_start_index"`
	EpisodeEndIndex                   int     `json:"episode_end_index"`
	EpisodeStartTime                  string  `json:"episode_start_time"`
	EpisodeEndTime                    string  `json:"episode_end_time"`
	RawLengthBars                     int     `json:"raw_length_bars"`
	ActiveLengthBars                  int     `json:"active_length_bars"`
	RawLengthBucket                   string  `json:"raw_length_bucket"`
	ActiveLengthBucket                string  `json:"active_length_bucket"`
	EpisodeHigh                       float64 `json:"episode_high"`
	EpisodeLow                        float64 `json:"episode_low"`
	EpisodeMid                        float64 `json:"episode_mid"`
	EpisodeEndClose                   float64 `json:"episode_end_close"`
	EpisodeWidthPct                   float64 `json:"episode_width_pct"`
	EpisodeWidthBucket                string  `json:"episode_width_bucket"`
	AvgNormalizedATR                  float64 `json:"avg_normalized_atr"`
	EndNormalizedATR                  float64 `json:"end_normalized_atr"`
	WidthToATRRatio                   float64 `json:"width_to_atr_ratio"`
	WidthToATRBucket                  string  `json:"width_to_atr_bucket"`
	DecisionIndex                     int     `json:"decision_index"`
	DecisionTime                      string  `json:"decision_time"`
	DecisionClose                     float64 `json:"decision_close"`
	DecisionClosePosition             float64 `json:"decision_close_position"`
	DecisionClosePositionBucket       string  `json:"decision_close_position_bucket"`
	DecisionMidSide                   string  `json:"decision_mid_side"`
	DecisionDistanceToHighPct         float64 `json:"decision_distance_to_high_pct"`
	DecisionDistanceToLowPct          float64 `json:"decision_distance_to_low_pct"`
	DecisionDistanceToMidPct          float64 `json:"decision_distance_to_mid_pct"`
	HorizonBars                       int     `json:"horizon_bars"`
	LabelWindowStartIndex             int     `json:"label_window_start_index"`
	LabelWindowEndIndex               int     `json:"label_window_end_index"`
	LabelWindowStartTime              string  `json:"label_window_start_time"`
	LabelWindowEndTime                string  `json:"label_window_end_time"`
	LabelTouchedMid                   bool    `json:"label_touched_mid"`
	LabelClosedAcrossMid              bool    `json:"label_closed_across_mid"`
	LabelFirstMidTouchDelayBars       int     `json:"label_first_mid_touch_delay_bars"`
	LabelFirstMidCloseAcrossDelayBars int     `json:"label_first_mid_close_across_delay_bars"`
	LabelMidTouchBeforeBoundaryTouch  bool    `json:"label_mid_touch_before_boundary_touch"`
	LabelMidCrossBeforeBoundaryBreak  bool    `json:"label_mid_cross_before_boundary_close_break"`
	LabelReenteredRange               bool    `json:"label_reentered_range"`
	LabelPersistedInsideRange         bool    `json:"label_persisted_inside_range"`
	LabelQuickInvalidated             bool    `json:"label_quick_invalidated"`
	LabelInvalidatedUp                bool    `json:"label_invalidated_up"`
	LabelInvalidatedDown              bool    `json:"label_invalidated_down"`
	LabelTrendedUp                    bool    `json:"label_trended_up"`
	LabelTrendedDown                  bool    `json:"label_trended_down"`
}

type HoldInsideMidlineTransitionSummaryRow struct {
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
	Split                                 string  `json:"split"`
	HorizonBars                           int     `json:"horizon_bars"`
	DecisionMidSide                       string  `json:"decision_mid_side"`
	DecisionClosePositionBucket           string  `json:"decision_close_position_bucket"`
	SourceEpisodeCount                    int     `json:"source_episode_count"`
	CandidateCount                        int     `json:"candidate_count"`
	CandidateRate                         float64 `json:"candidate_rate"`
	AvgRawLengthBars                      float64 `json:"avg_raw_length_bars"`
	AvgActiveLengthBars                   float64 `json:"avg_active_length_bars"`
	AvgEpisodeWidthPct                    float64 `json:"avg_episode_width_pct"`
	AvgNormalizedATR                      float64 `json:"avg_normalized_atr"`
	AvgEndNormalizedATR                   float64 `json:"avg_end_normalized_atr"`
	AvgWidthToATRRatio                    float64 `json:"avg_width_to_atr_ratio"`
	AvgDecisionClosePosition              float64 `json:"avg_decision_close_position"`
	AvgDecisionDistanceToHighPct          float64 `json:"avg_decision_distance_to_high_pct"`
	AvgDecisionDistanceToLowPct           float64 `json:"avg_decision_distance_to_low_pct"`
	AvgDecisionDistanceToMidPct           float64 `json:"avg_decision_distance_to_mid_pct"`
	LabelTouchedMidCount                  int     `json:"label_touched_mid_count"`
	LabelClosedAcrossMidCount             int     `json:"label_closed_across_mid_count"`
	LabelMidTouchBeforeBoundaryTouchCount int     `json:"label_mid_touch_before_boundary_touch_count"`
	LabelMidCrossBeforeBoundaryBreakCount int     `json:"label_mid_cross_before_boundary_close_break_count"`
	LabelReenteredRangeCount              int     `json:"label_reentered_range_count"`
	LabelPersistedInsideRangeCount        int     `json:"label_persisted_inside_range_count"`
	LabelQuickInvalidatedCount            int     `json:"label_quick_invalidated_count"`
	LabelInvalidatedUpCount               int     `json:"label_invalidated_up_count"`
	LabelInvalidatedDownCount             int     `json:"label_invalidated_down_count"`
	LabelTrendedUpCount                   int     `json:"label_trended_up_count"`
	LabelTrendedDownCount                 int     `json:"label_trended_down_count"`
	LabelTouchedMidRate                   float64 `json:"label_touched_mid_rate"`
	LabelClosedAcrossMidRate              float64 `json:"label_closed_across_mid_rate"`
	LabelMidTouchBeforeBoundaryTouchRate  float64 `json:"label_mid_touch_before_boundary_touch_rate"`
	LabelMidCrossBeforeBoundaryBreakRate  float64 `json:"label_mid_cross_before_boundary_close_break_rate"`
	LabelReenteredRangeRate               float64 `json:"label_reentered_range_rate"`
	LabelPersistedInsideRangeRate         float64 `json:"label_persisted_inside_range_rate"`
	LabelQuickInvalidatedRate             float64 `json:"label_quick_invalidated_rate"`
	LabelInvalidatedUpRate                float64 `json:"label_invalidated_up_rate"`
	LabelInvalidatedDownRate              float64 `json:"label_invalidated_down_rate"`
	LabelTrendedUpRate                    float64 `json:"label_trended_up_rate"`
	LabelTrendedDownRate                  float64 `json:"label_trended_down_rate"`
	LabelAvgFirstMidTouchDelayBars        float64 `json:"label_avg_first_mid_touch_delay_bars"`
	LabelAvgFirstMidCloseAcrossDelayBars  float64 `json:"label_avg_first_mid_close_across_delay_bars"`
}

type HoldInsideMidlineTransitionStabilityRow struct {
	ProfileID                                 string  `json:"profile_id"`
	IsBalancedBaseline                        bool    `json:"is_balanced_baseline"`
	IsADXComparison                           bool    `json:"is_adx_comparison"`
	Percentile                                float64 `json:"percentile"`
	MinConsecutiveBars                        int     `json:"min_consecutive_bars"`
	UseBollinger                              bool    `json:"use_bollinger"`
	UseADX                                    bool    `json:"use_adx"`
	LookbackDays                              int     `json:"lookback_days"`
	ContextRule                               string  `json:"context_rule"`
	HoldBars                                  int     `json:"hold_bars"`
	RequireMid50                              bool    `json:"require_mid_50"`
	HorizonBars                               int     `json:"horizon_bars"`
	DecisionMidSide                           string  `json:"decision_mid_side"`
	DecisionClosePositionBucket               string  `json:"decision_close_position_bucket"`
	PeriodSplits                              int     `json:"period_splits"`
	SourceEpisodeCount                        int     `json:"source_episode_count"`
	SourceEpisodeCountMin                     int     `json:"source_episode_count_min"`
	SourceEpisodeCountMax                     int     `json:"source_episode_count_max"`
	SourceEpisodeCountDelta                   int     `json:"source_episode_count_delta"`
	CandidateCount                            int     `json:"candidate_count"`
	CandidateCountMin                         int     `json:"candidate_count_min"`
	CandidateCountMax                         int     `json:"candidate_count_max"`
	CandidateCountDelta                       int     `json:"candidate_count_delta"`
	CandidateRateMin                          float64 `json:"candidate_rate_min"`
	CandidateRateMax                          float64 `json:"candidate_rate_max"`
	CandidateRateDelta                        float64 `json:"candidate_rate_delta"`
	LabelTouchedMidRateMin                    float64 `json:"label_touched_mid_rate_min"`
	LabelTouchedMidRateMax                    float64 `json:"label_touched_mid_rate_max"`
	LabelTouchedMidRateDelta                  float64 `json:"label_touched_mid_rate_delta"`
	LabelClosedAcrossMidRateMin               float64 `json:"label_closed_across_mid_rate_min"`
	LabelClosedAcrossMidRateMax               float64 `json:"label_closed_across_mid_rate_max"`
	LabelClosedAcrossMidRateDelta             float64 `json:"label_closed_across_mid_rate_delta"`
	LabelMidTouchBeforeBoundaryTouchRateMin   float64 `json:"label_mid_touch_before_boundary_touch_rate_min"`
	LabelMidTouchBeforeBoundaryTouchRateMax   float64 `json:"label_mid_touch_before_boundary_touch_rate_max"`
	LabelMidTouchBeforeBoundaryTouchRateDelta float64 `json:"label_mid_touch_before_boundary_touch_rate_delta"`
	LabelMidCrossBeforeBoundaryBreakRateMin   float64 `json:"label_mid_cross_before_boundary_close_break_rate_min"`
	LabelMidCrossBeforeBoundaryBreakRateMax   float64 `json:"label_mid_cross_before_boundary_close_break_rate_max"`
	LabelMidCrossBeforeBoundaryBreakRateDelta float64 `json:"label_mid_cross_before_boundary_close_break_rate_delta"`
	LabelReenteredRangeRateMin                float64 `json:"label_reentered_range_rate_min"`
	LabelReenteredRangeRateMax                float64 `json:"label_reentered_range_rate_max"`
	LabelReenteredRangeRateDelta              float64 `json:"label_reentered_range_rate_delta"`
	LabelPersistedInsideRangeRateMin          float64 `json:"label_persisted_inside_range_rate_min"`
	LabelPersistedInsideRangeRateMax          float64 `json:"label_persisted_inside_range_rate_max"`
	LabelPersistedInsideRangeRateDelta        float64 `json:"label_persisted_inside_range_rate_delta"`
	LabelQuickInvalidatedRateMin              float64 `json:"label_quick_invalidated_rate_min"`
	LabelQuickInvalidatedRateMax              float64 `json:"label_quick_invalidated_rate_max"`
	LabelQuickInvalidatedRateDelta            float64 `json:"label_quick_invalidated_rate_delta"`
	LabelInvalidatedUpRateMin                 float64 `json:"label_invalidated_up_rate_min"`
	LabelInvalidatedUpRateMax                 float64 `json:"label_invalidated_up_rate_max"`
	LabelInvalidatedUpRateDelta               float64 `json:"label_invalidated_up_rate_delta"`
	LabelInvalidatedDownRateMin               float64 `json:"label_invalidated_down_rate_min"`
	LabelInvalidatedDownRateMax               float64 `json:"label_invalidated_down_rate_max"`
	LabelInvalidatedDownRateDelta             float64 `json:"label_invalidated_down_rate_delta"`
	LabelTrendedUpRateMin                     float64 `json:"label_trended_up_rate_min"`
	LabelTrendedUpRateMax                     float64 `json:"label_trended_up_rate_max"`
	LabelTrendedUpRateDelta                   float64 `json:"label_trended_up_rate_delta"`
	LabelTrendedDownRateMin                   float64 `json:"label_trended_down_rate_min"`
	LabelTrendedDownRateMax                   float64 `json:"label_trended_down_rate_max"`
	LabelTrendedDownRateDelta                 float64 `json:"label_trended_down_rate_delta"`
	LabelTrendedRateMin                       float64 `json:"label_trended_rate_min"`
	LabelTrendedRateMax                       float64 `json:"label_trended_rate_max"`
	LabelTrendedRateDelta                     float64 `json:"label_trended_rate_delta"`
	LabelAvgFirstMidTouchDelayBarsMin         float64 `json:"label_avg_first_mid_touch_delay_bars_min"`
	LabelAvgFirstMidTouchDelayBarsMax         float64 `json:"label_avg_first_mid_touch_delay_bars_max"`
	LabelAvgFirstMidTouchDelayBarsDelta       float64 `json:"label_avg_first_mid_touch_delay_bars_delta"`
	LabelAvgFirstMidCloseAcrossDelayBarsMin   float64 `json:"label_avg_first_mid_close_across_delay_bars_min"`
	LabelAvgFirstMidCloseAcrossDelayBarsMax   float64 `json:"label_avg_first_mid_close_across_delay_bars_max"`
	LabelAvgFirstMidCloseAcrossDelayBarsDelta float64 `json:"label_avg_first_mid_close_across_delay_bars_delta"`
}

func DefaultHoldInsideMidlineTransitionAuditConfig() HoldInsideMidlineTransitionAuditConfig {
	lookbackDays := DefaultCompressionRangeDetectorConfig().LookbackDays
	return HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1, 3, 6, 12},
		QuickInvalidationBars: 3,
		Profiles:              defaultHoldInsideDirectionalEdgeProfiles(lookbackDays),
		ContextRules:          defaultHoldInsideDirectionalEdgeRules(),
	}
}

func RunHoldInsideMidlineTransitionAudit(candles []Candle, detectorCfg RangeDetectorConfig, cfg HoldInsideMidlineTransitionAuditConfig, splits []Split) ([]HoldInsideMidlineTransitionCandidateRow, []HoldInsideMidlineTransitionSummaryRow, []HoldInsideMidlineTransitionStabilityRow, error) {
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

	candidates := []HoldInsideMidlineTransitionCandidateRow{}
	summaryAccumulators := map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator{}
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
		profileCandidates := runHoldInsideMidlineTransitionForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
		candidates = append(candidates, profileCandidates...)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideMidlineTransitionCandidate(candidates[i], candidates[j])
	})
	summaryRows := holdInsideMidlineTransitionSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideMidlineTransitionStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func runHoldInsideMidlineTransitionAuditFromClassifications(candles []Candle, profile DetectorSweepProfile, classifications []RangeClassification, cfg HoldInsideMidlineTransitionAuditConfig, splits []Split) ([]HoldInsideMidlineTransitionCandidateRow, []HoldInsideMidlineTransitionSummaryRow, []HoldInsideMidlineTransitionStabilityRow, error) {
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
	summaryAccumulators := map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator{}
	candidates := runHoldInsideMidlineTransitionForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideMidlineTransitionCandidate(candidates[i], candidates[j])
	})
	summaryRows := holdInsideMidlineTransitionSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideMidlineTransitionStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func (cfg HoldInsideMidlineTransitionAuditConfig) withDefaults(lookbackDays int) HoldInsideMidlineTransitionAuditConfig {
	defaults := DefaultHoldInsideMidlineTransitionAuditConfig()
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

func (cfg HoldInsideMidlineTransitionAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("hold-inside midline transition audit horizon bars must be positive")
		}
	}
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("hold-inside midline transition audit quick invalidation bars must be positive")
	}
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("hold-inside midline transition audit profiles must not be empty")
	}
	for _, profile := range cfg.Profiles {
		if profile.ProfileID == "" {
			return fmt.Errorf("hold-inside midline transition audit profile id must not be empty")
		}
	}
	if len(cfg.ContextRules) == 0 {
		return fmt.Errorf("hold-inside midline transition audit rules must not be empty")
	}
	for _, rule := range cfg.ContextRules {
		if rule.RuleID == "" {
			return fmt.Errorf("hold-inside midline transition audit rule id must not be empty")
		}
		if rule.HoldBars <= 0 {
			return fmt.Errorf("hold-inside midline transition audit rules must use positive hold bars")
		}
		if rule.RequireMid50 && rule.HoldBars == 0 {
			return fmt.Errorf("hold-inside midline transition audit mid-50 rule requires positive hold bars")
		}
	}
	return nil
}

func runHoldInsideMidlineTransitionForProfile(candles []Candle, classifications []RangeClassification, normalizedATR []float64, profile DetectorSweepProfile, cfg HoldInsideMidlineTransitionAuditConfig, splits []Split, summaryAccumulators map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator) []HoldInsideMidlineTransitionCandidateRow {
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, profile.ProfileID)
	candidates := []HoldInsideMidlineTransitionCandidateRow{}
	for _, episode := range episodes {
		for _, rule := range cfg.ContextRules {
			decisionIndex := episode.EndIndex + rule.HoldBars
			if decisionIndex < 0 || decisionIndex >= len(candles) {
				continue
			}
			split := splitNameForCloseTime(candles[decisionIndex].CloseTime, splits)
			episodeMid := (episode.High + episode.Low) / 2
			decisionMidSide := holdInsideDecisionMidSide(candles[decisionIndex].Close, episodeMid)
			decisionBucket := decisionClosePositionBucket(decisionClosePosition(candles[decisionIndex].Close, episode.Low, episode.High))
			passesRule := detectorContextRulePasses(candles, episode, rule, decisionIndex)
			for _, horizon := range cfg.HorizonsBars {
				if decisionIndex+horizon >= len(candles) {
					continue
				}
				for _, combo := range holdInsideMidlineTransitionSummaryCombos(decisionMidSide, decisionBucket) {
					addHoldInsideMidlineTransitionSource(summaryAccumulators, profile, rule, split, horizon, combo.midSide, combo.bucket)
					if split != fullSplitName {
						addHoldInsideMidlineTransitionSource(summaryAccumulators, profile, rule, fullSplitName, horizon, combo.midSide, combo.bucket)
					}
				}
				if !passesRule {
					continue
				}
				row, ok := newHoldInsideMidlineTransitionCandidateRow(candles, profile, rule, episode, split, decisionIndex, horizon, cfg.QuickInvalidationBars)
				if !ok {
					continue
				}
				candidates = append(candidates, row)
				for _, combo := range holdInsideMidlineTransitionSummaryCombos(row.DecisionMidSide, row.DecisionClosePositionBucket) {
					addHoldInsideMidlineTransitionCandidate(summaryAccumulators, row, split, combo.midSide, combo.bucket)
					if split != fullSplitName {
						addHoldInsideMidlineTransitionCandidate(summaryAccumulators, row, fullSplitName, combo.midSide, combo.bucket)
					}
				}
			}
		}
	}
	return candidates
}

func newHoldInsideMidlineTransitionCandidateRow(candles []Candle, profile DetectorSweepProfile, rule DetectorContextRefinementRule, episode rangeRegimeDurabilityEpisode, split string, decisionIndex, horizon, quickInvalidationBars int) (HoldInsideMidlineTransitionCandidateRow, bool) {
	if horizon <= 0 || quickInvalidationBars <= 0 || decisionIndex < 0 || decisionIndex+horizon >= len(candles) {
		return HoldInsideMidlineTransitionCandidateRow{}, false
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
	label := newHoldInsideMidlineTransitionLabel(future, decision.Close, episode.Low, episode.High, episodeMid, quickInvalidationBars)

	return HoldInsideMidlineTransitionCandidateRow{
		ProfileID:                         profile.ProfileID,
		IsBalancedBaseline:                profile.IsBalancedBaseline,
		IsADXComparison:                   profile.IsADXComparison,
		Percentile:                        profile.Percentile,
		MinConsecutiveBars:                profile.MinConsecutiveBars,
		UseBollinger:                      profile.UseBollinger,
		UseADX:                            profile.UseADX,
		LookbackDays:                      profile.LookbackDays,
		ContextRule:                       rule.RuleID,
		HoldBars:                          rule.HoldBars,
		RequireMid50:                      rule.RequireMid50,
		Split:                             split,
		SourceEpisodeID:                   episode.EpisodeID,
		EpisodeStartIndex:                 episode.StartIndex,
		EpisodeEndIndex:                   episode.EndIndex,
		EpisodeStartTime:                  candles[episode.StartIndex].CloseTime.Format(timeLayout),
		EpisodeEndTime:                    candles[episode.EndIndex].CloseTime.Format(timeLayout),
		RawLengthBars:                     episode.RawLengthBars,
		ActiveLengthBars:                  episode.ActiveLengthBars,
		RawLengthBucket:                   episode.RawLengthBucket,
		ActiveLengthBucket:                episode.ActiveLengthBucket,
		EpisodeHigh:                       episode.High,
		EpisodeLow:                        episode.Low,
		EpisodeMid:                        episodeMid,
		EpisodeEndClose:                   episode.EndClose,
		EpisodeWidthPct:                   episode.WidthPct,
		EpisodeWidthBucket:                episode.WidthBucket,
		AvgNormalizedATR:                  episode.AvgNormalizedATR,
		EndNormalizedATR:                  episode.EndNormalizedATR,
		WidthToATRRatio:                   episode.WidthToATRRatio,
		WidthToATRBucket:                  episode.WidthToATRBucket,
		DecisionIndex:                     decisionIndex,
		DecisionTime:                      decision.CloseTime.Format(timeLayout),
		DecisionClose:                     decision.Close,
		DecisionClosePosition:             decisionPositionValue,
		DecisionClosePositionBucket:       decisionClosePositionBucket(decisionPosition),
		DecisionMidSide:                   holdInsideDecisionMidSide(decision.Close, episodeMid),
		DecisionDistanceToHighPct:         movePct(math.Max(0, episode.High-decision.Close), decision.Close),
		DecisionDistanceToLowPct:          movePct(math.Max(0, decision.Close-episode.Low), decision.Close),
		DecisionDistanceToMidPct:          movePct(math.Abs(decision.Close-episodeMid), decision.Close),
		HorizonBars:                       horizon,
		LabelWindowStartIndex:             startIndex,
		LabelWindowEndIndex:               endIndex,
		LabelWindowStartTime:              candles[startIndex].CloseTime.Format(timeLayout),
		LabelWindowEndTime:                candles[endIndex].CloseTime.Format(timeLayout),
		LabelTouchedMid:                   label.TouchedMid,
		LabelClosedAcrossMid:              label.ClosedAcrossMid,
		LabelFirstMidTouchDelayBars:       label.FirstMidTouchDelayBars,
		LabelFirstMidCloseAcrossDelayBars: label.FirstMidCloseAcrossDelayBars,
		LabelMidTouchBeforeBoundaryTouch:  label.MidTouchBeforeBoundaryTouch,
		LabelMidCrossBeforeBoundaryBreak:  label.MidCrossBeforeBoundaryBreak,
		LabelReenteredRange:               label.ReenteredRange,
		LabelPersistedInsideRange:         label.PersistedInsideRange,
		LabelQuickInvalidated:             label.QuickInvalidated,
		LabelInvalidatedUp:                label.InvalidatedUp,
		LabelInvalidatedDown:              label.InvalidatedDown,
		LabelTrendedUp:                    label.TrendedUp,
		LabelTrendedDown:                  label.TrendedDown,
	}, true
}

type holdInsideMidlineTransitionLabel struct {
	TouchedMid                   bool
	ClosedAcrossMid              bool
	FirstMidTouchDelayBars       int
	FirstMidCloseAcrossDelayBars int
	MidTouchBeforeBoundaryTouch  bool
	MidCrossBeforeBoundaryBreak  bool
	ReenteredRange               bool
	PersistedInsideRange         bool
	QuickInvalidated             bool
	InvalidatedUp                bool
	InvalidatedDown              bool
	TrendedUp                    bool
	TrendedDown                  bool
}

func newHoldInsideMidlineTransitionLabel(future []Candle, decisionClose, low, high, mid float64, quickInvalidationBars int) holdInsideMidlineTransitionLabel {
	futureMaxHigh, futureMinLow := futureHighLow(future)
	lastClose := future[len(future)-1].Close
	maxUpMovePct := movePct(math.Max(0, futureMaxHigh-decisionClose), decisionClose)
	maxDownMovePct := movePct(math.Max(0, decisionClose-futureMinLow), decisionClose)
	label := holdInsideMidlineTransitionLabel{
		FirstMidTouchDelayBars:       -1,
		FirstMidCloseAcrossDelayBars: -1,
		PersistedInsideRange:         true,
	}
	firstBoundaryTouchDelay := -1
	firstBoundaryBreakDelay := -1
	quickLimit := minInt(len(future), quickInvalidationBars)
	for i, candle := range future {
		delay := i + 1
		touchedMid, closedAcrossMid := holdInsideMidLabels(candle, decisionClose, mid)
		if touchedMid && label.FirstMidTouchDelayBars == -1 {
			label.FirstMidTouchDelayBars = delay
		}
		if closedAcrossMid && label.FirstMidCloseAcrossDelayBars == -1 {
			label.FirstMidCloseAcrossDelayBars = delay
		}
		if (candle.High >= high || candle.Low <= low) && firstBoundaryTouchDelay == -1 {
			firstBoundaryTouchDelay = delay
		}
		if closeInsideRange(candle.Close, low, high) {
			label.ReenteredRange = true
		} else {
			label.PersistedInsideRange = false
		}
		if candle.Close > high {
			label.InvalidatedUp = true
			if firstBoundaryBreakDelay == -1 {
				firstBoundaryBreakDelay = delay
			}
			if i < quickLimit {
				label.QuickInvalidated = true
			}
		}
		if candle.Close < low {
			label.InvalidatedDown = true
			if firstBoundaryBreakDelay == -1 {
				firstBoundaryBreakDelay = delay
			}
			if i < quickLimit {
				label.QuickInvalidated = true
			}
		}
	}
	label.TouchedMid = label.FirstMidTouchDelayBars != -1
	label.ClosedAcrossMid = label.FirstMidCloseAcrossDelayBars != -1
	label.MidTouchBeforeBoundaryTouch = label.FirstMidTouchDelayBars != -1 &&
		(firstBoundaryTouchDelay == -1 || label.FirstMidTouchDelayBars < firstBoundaryTouchDelay)
	label.MidCrossBeforeBoundaryBreak = label.FirstMidCloseAcrossDelayBars != -1 &&
		(firstBoundaryBreakDelay == -1 || label.FirstMidCloseAcrossDelayBars < firstBoundaryBreakDelay)
	label.TrendedUp = label.InvalidatedUp && lastClose > high && maxUpMovePct > maxDownMovePct
	label.TrendedDown = label.InvalidatedDown && lastClose < low && maxDownMovePct > maxUpMovePct
	return label
}

type holdInsideMidlineTransitionSummaryCombo struct {
	midSide string
	bucket  string
}

func holdInsideMidlineTransitionSummaryCombos(midSide, bucket string) []holdInsideMidlineTransitionSummaryCombo {
	return []holdInsideMidlineTransitionSummaryCombo{
		{midSide: midSide, bucket: bucket},
		{midSide: midSide, bucket: holdInsideDecisionPositionBucketAll},
		{midSide: holdInsideDecisionMidSideAll, bucket: bucket},
		{midSide: holdInsideDecisionMidSideAll, bucket: holdInsideDecisionPositionBucketAll},
	}
}

type holdInsideMidlineTransitionSummaryKey struct {
	profileID                   string
	contextRule                 string
	split                       string
	horizonBars                 int
	decisionMidSide             string
	decisionClosePositionBucket string
}

type holdInsideMidlineTransitionSummaryAccumulator struct {
	profile                       DetectorSweepProfile
	rule                          DetectorContextRefinementRule
	split                         string
	horizonBars                   int
	decisionMidSide               string
	decisionClosePositionBucket   string
	sourceEpisodes                int
	candidates                    int
	rawLengthSum                  float64
	activeLengthSum               float64
	widthPctSum                   float64
	avgATRSum                     float64
	endATRSum                     float64
	widthToATRRatioSum            float64
	decisionPositionSum           float64
	decisionDistanceToHighPctSum  float64
	decisionDistanceToLowPctSum   float64
	decisionDistanceToMidPctSum   float64
	touchedMid                    int
	closedAcrossMid               int
	midTouchBeforeBoundaryTouch   int
	midCrossBeforeBoundaryBreak   int
	reenteredRange                int
	persistedInsideRange          int
	quickInvalidated              int
	invalidatedUp                 int
	invalidatedDown               int
	trendedUp                     int
	trendedDown                   int
	firstMidTouchDelaySum         float64
	firstMidTouchDelayCount       int
	firstMidCloseAcrossDelaySum   float64
	firstMidCloseAcrossDelayCount int
}

func addHoldInsideMidlineTransitionSource(accumulators map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int, midSide, bucket string) {
	acc := holdInsideMidlineTransitionSummaryAccumulatorFor(accumulators, profile, rule, split, horizon, midSide, bucket)
	acc.sourceEpisodes++
}

func addHoldInsideMidlineTransitionCandidate(accumulators map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator, row HoldInsideMidlineTransitionCandidateRow, split, midSide, bucket string) {
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
	acc := holdInsideMidlineTransitionSummaryAccumulatorFor(accumulators, profile, rule, split, row.HorizonBars, midSide, bucket)
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
	if row.LabelTouchedMid {
		acc.touchedMid++
	}
	if row.LabelClosedAcrossMid {
		acc.closedAcrossMid++
	}
	if row.LabelMidTouchBeforeBoundaryTouch {
		acc.midTouchBeforeBoundaryTouch++
	}
	if row.LabelMidCrossBeforeBoundaryBreak {
		acc.midCrossBeforeBoundaryBreak++
	}
	if row.LabelReenteredRange {
		acc.reenteredRange++
	}
	if row.LabelPersistedInsideRange {
		acc.persistedInsideRange++
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
	if row.LabelFirstMidTouchDelayBars >= 0 {
		acc.firstMidTouchDelaySum += float64(row.LabelFirstMidTouchDelayBars)
		acc.firstMidTouchDelayCount++
	}
	if row.LabelFirstMidCloseAcrossDelayBars >= 0 {
		acc.firstMidCloseAcrossDelaySum += float64(row.LabelFirstMidCloseAcrossDelayBars)
		acc.firstMidCloseAcrossDelayCount++
	}
}

func holdInsideMidlineTransitionSummaryAccumulatorFor(accumulators map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int, midSide, bucket string) *holdInsideMidlineTransitionSummaryAccumulator {
	key := holdInsideMidlineTransitionSummaryKey{
		profileID:                   profile.ProfileID,
		contextRule:                 rule.RuleID,
		split:                       split,
		horizonBars:                 horizon,
		decisionMidSide:             midSide,
		decisionClosePositionBucket: bucket,
	}
	acc := accumulators[key]
	if acc == nil {
		acc = &holdInsideMidlineTransitionSummaryAccumulator{
			profile:                     profile,
			rule:                        rule,
			split:                       split,
			horizonBars:                 horizon,
			decisionMidSide:             midSide,
			decisionClosePositionBucket: bucket,
		}
		accumulators[key] = acc
	}
	return acc
}

func holdInsideMidlineTransitionSummaryRows(accumulators map[holdInsideMidlineTransitionSummaryKey]*holdInsideMidlineTransitionSummaryAccumulator) []HoldInsideMidlineTransitionSummaryRow {
	rows := make([]HoldInsideMidlineTransitionSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideMidlineTransitionSummary(rows[i], rows[j])
	})
	return rows
}

func (acc holdInsideMidlineTransitionSummaryAccumulator) row() HoldInsideMidlineTransitionSummaryRow {
	row := HoldInsideMidlineTransitionSummaryRow{
		ProfileID:                             acc.profile.ProfileID,
		IsBalancedBaseline:                    acc.profile.IsBalancedBaseline,
		IsADXComparison:                       acc.profile.IsADXComparison,
		Percentile:                            acc.profile.Percentile,
		MinConsecutiveBars:                    acc.profile.MinConsecutiveBars,
		UseBollinger:                          acc.profile.UseBollinger,
		UseADX:                                acc.profile.UseADX,
		LookbackDays:                          acc.profile.LookbackDays,
		ContextRule:                           acc.rule.RuleID,
		HoldBars:                              acc.rule.HoldBars,
		RequireMid50:                          acc.rule.RequireMid50,
		Split:                                 acc.split,
		HorizonBars:                           acc.horizonBars,
		DecisionMidSide:                       acc.decisionMidSide,
		DecisionClosePositionBucket:           acc.decisionClosePositionBucket,
		SourceEpisodeCount:                    acc.sourceEpisodes,
		CandidateCount:                        acc.candidates,
		LabelTouchedMidCount:                  acc.touchedMid,
		LabelClosedAcrossMidCount:             acc.closedAcrossMid,
		LabelMidTouchBeforeBoundaryTouchCount: acc.midTouchBeforeBoundaryTouch,
		LabelMidCrossBeforeBoundaryBreakCount: acc.midCrossBeforeBoundaryBreak,
		LabelReenteredRangeCount:              acc.reenteredRange,
		LabelPersistedInsideRangeCount:        acc.persistedInsideRange,
		LabelQuickInvalidatedCount:            acc.quickInvalidated,
		LabelInvalidatedUpCount:               acc.invalidatedUp,
		LabelInvalidatedDownCount:             acc.invalidatedDown,
		LabelTrendedUpCount:                   acc.trendedUp,
		LabelTrendedDownCount:                 acc.trendedDown,
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
	row.LabelTouchedMidRate = float64(acc.touchedMid) / count
	row.LabelClosedAcrossMidRate = float64(acc.closedAcrossMid) / count
	row.LabelMidTouchBeforeBoundaryTouchRate = float64(acc.midTouchBeforeBoundaryTouch) / count
	row.LabelMidCrossBeforeBoundaryBreakRate = float64(acc.midCrossBeforeBoundaryBreak) / count
	row.LabelReenteredRangeRate = float64(acc.reenteredRange) / count
	row.LabelPersistedInsideRangeRate = float64(acc.persistedInsideRange) / count
	row.LabelQuickInvalidatedRate = float64(acc.quickInvalidated) / count
	row.LabelInvalidatedUpRate = float64(acc.invalidatedUp) / count
	row.LabelInvalidatedDownRate = float64(acc.invalidatedDown) / count
	row.LabelTrendedUpRate = float64(acc.trendedUp) / count
	row.LabelTrendedDownRate = float64(acc.trendedDown) / count
	if acc.firstMidTouchDelayCount > 0 {
		row.LabelAvgFirstMidTouchDelayBars = acc.firstMidTouchDelaySum / float64(acc.firstMidTouchDelayCount)
	}
	if acc.firstMidCloseAcrossDelayCount > 0 {
		row.LabelAvgFirstMidCloseAcrossDelayBars = acc.firstMidCloseAcrossDelaySum / float64(acc.firstMidCloseAcrossDelayCount)
	}
	return row
}

func holdInsideMidlineTransitionStabilityRows(profiles []DetectorSweepProfile, rules []DetectorContextRefinementRule, horizons []int, summaryRows []HoldInsideMidlineTransitionSummaryRow, splits []Split) []HoldInsideMidlineTransitionStabilityRow {
	periodSplits := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplits = append(periodSplits, split)
		}
	}
	byKey := map[holdInsideMidlineTransitionSummaryKey]HoldInsideMidlineTransitionSummaryRow{}
	combosByBase := map[holdInsideMidlineTransitionStabilityBaseKey]map[holdInsideMidlineTransitionSummaryCombo]bool{}
	for _, row := range summaryRows {
		key := holdInsideMidlineTransitionSummaryKey{
			profileID:                   row.ProfileID,
			contextRule:                 row.ContextRule,
			split:                       row.Split,
			horizonBars:                 row.HorizonBars,
			decisionMidSide:             row.DecisionMidSide,
			decisionClosePositionBucket: row.DecisionClosePositionBucket,
		}
		byKey[key] = row
		base := holdInsideMidlineTransitionStabilityBaseKey{
			profileID:   row.ProfileID,
			contextRule: row.ContextRule,
			horizonBars: row.HorizonBars,
		}
		if combosByBase[base] == nil {
			combosByBase[base] = map[holdInsideMidlineTransitionSummaryCombo]bool{}
		}
		combosByBase[base][holdInsideMidlineTransitionSummaryCombo{midSide: row.DecisionMidSide, bucket: row.DecisionClosePositionBucket}] = true
	}

	rows := []HoldInsideMidlineTransitionStabilityRow{}
	for _, profile := range profiles {
		for _, rule := range rules {
			for _, horizon := range horizons {
				base := holdInsideMidlineTransitionStabilityBaseKey{
					profileID:   profile.ProfileID,
					contextRule: rule.RuleID,
					horizonBars: horizon,
				}
				for _, combo := range holdInsideMidlineTransitionSortedCombos(combosByBase[base]) {
					rows = append(rows, newHoldInsideMidlineTransitionStabilityRow(profile, rule, horizon, combo.midSide, combo.bucket, periodSplits, byKey))
				}
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideMidlineTransitionStability(rows[i], rows[j])
	})
	return rows
}

type holdInsideMidlineTransitionStabilityBaseKey struct {
	profileID   string
	contextRule string
	horizonBars int
}

func newHoldInsideMidlineTransitionStabilityRow(profile DetectorSweepProfile, rule DetectorContextRefinementRule, horizon int, midSide, bucket string, periodSplits []Split, byKey map[holdInsideMidlineTransitionSummaryKey]HoldInsideMidlineTransitionSummaryRow) HoldInsideMidlineTransitionStabilityRow {
	row := HoldInsideMidlineTransitionStabilityRow{
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
		DecisionMidSide:             midSide,
		DecisionClosePositionBucket: bucket,
		PeriodSplits:                len(periodSplits),
	}
	first := true
	for _, split := range periodSplits {
		summary := byKey[holdInsideMidlineTransitionSummaryKey{
			profileID:                   profile.ProfileID,
			contextRule:                 rule.RuleID,
			split:                       split.Name,
			horizonBars:                 horizon,
			decisionMidSide:             midSide,
			decisionClosePositionBucket: bucket,
		}]
		if first {
			row.SourceEpisodeCountMin = summary.SourceEpisodeCount
			row.SourceEpisodeCountMax = summary.SourceEpisodeCount
			row.CandidateCountMin = summary.CandidateCount
			row.CandidateCountMax = summary.CandidateCount
			row.CandidateRateMin = summary.CandidateRate
			row.CandidateRateMax = summary.CandidateRate
			row.LabelTouchedMidRateMin = summary.LabelTouchedMidRate
			row.LabelTouchedMidRateMax = summary.LabelTouchedMidRate
			row.LabelClosedAcrossMidRateMin = summary.LabelClosedAcrossMidRate
			row.LabelClosedAcrossMidRateMax = summary.LabelClosedAcrossMidRate
			row.LabelMidTouchBeforeBoundaryTouchRateMin = summary.LabelMidTouchBeforeBoundaryTouchRate
			row.LabelMidTouchBeforeBoundaryTouchRateMax = summary.LabelMidTouchBeforeBoundaryTouchRate
			row.LabelMidCrossBeforeBoundaryBreakRateMin = summary.LabelMidCrossBeforeBoundaryBreakRate
			row.LabelMidCrossBeforeBoundaryBreakRateMax = summary.LabelMidCrossBeforeBoundaryBreakRate
			row.LabelReenteredRangeRateMin = summary.LabelReenteredRangeRate
			row.LabelReenteredRangeRateMax = summary.LabelReenteredRangeRate
			row.LabelPersistedInsideRangeRateMin = summary.LabelPersistedInsideRangeRate
			row.LabelPersistedInsideRangeRateMax = summary.LabelPersistedInsideRangeRate
			row.LabelQuickInvalidatedRateMin = summary.LabelQuickInvalidatedRate
			row.LabelQuickInvalidatedRateMax = summary.LabelQuickInvalidatedRate
			row.LabelInvalidatedUpRateMin = summary.LabelInvalidatedUpRate
			row.LabelInvalidatedUpRateMax = summary.LabelInvalidatedUpRate
			row.LabelInvalidatedDownRateMin = summary.LabelInvalidatedDownRate
			row.LabelInvalidatedDownRateMax = summary.LabelInvalidatedDownRate
			row.LabelTrendedUpRateMin = summary.LabelTrendedUpRate
			row.LabelTrendedUpRateMax = summary.LabelTrendedUpRate
			row.LabelTrendedDownRateMin = summary.LabelTrendedDownRate
			row.LabelTrendedDownRateMax = summary.LabelTrendedDownRate
			row.LabelTrendedRateMin = summary.LabelTrendedUpRate + summary.LabelTrendedDownRate
			row.LabelTrendedRateMax = summary.LabelTrendedUpRate + summary.LabelTrendedDownRate
			row.LabelAvgFirstMidTouchDelayBarsMin = summary.LabelAvgFirstMidTouchDelayBars
			row.LabelAvgFirstMidTouchDelayBarsMax = summary.LabelAvgFirstMidTouchDelayBars
			row.LabelAvgFirstMidCloseAcrossDelayBarsMin = summary.LabelAvgFirstMidCloseAcrossDelayBars
			row.LabelAvgFirstMidCloseAcrossDelayBarsMax = summary.LabelAvgFirstMidCloseAcrossDelayBars
			first = false
		} else {
			trendedRate := summary.LabelTrendedUpRate + summary.LabelTrendedDownRate
			row.SourceEpisodeCountMin = minInt(row.SourceEpisodeCountMin, summary.SourceEpisodeCount)
			row.SourceEpisodeCountMax = maxInt(row.SourceEpisodeCountMax, summary.SourceEpisodeCount)
			row.CandidateCountMin = minInt(row.CandidateCountMin, summary.CandidateCount)
			row.CandidateCountMax = maxInt(row.CandidateCountMax, summary.CandidateCount)
			row.CandidateRateMin = minFloat(row.CandidateRateMin, summary.CandidateRate)
			row.CandidateRateMax = maxFloat(row.CandidateRateMax, summary.CandidateRate)
			row.LabelTouchedMidRateMin = minFloat(row.LabelTouchedMidRateMin, summary.LabelTouchedMidRate)
			row.LabelTouchedMidRateMax = maxFloat(row.LabelTouchedMidRateMax, summary.LabelTouchedMidRate)
			row.LabelClosedAcrossMidRateMin = minFloat(row.LabelClosedAcrossMidRateMin, summary.LabelClosedAcrossMidRate)
			row.LabelClosedAcrossMidRateMax = maxFloat(row.LabelClosedAcrossMidRateMax, summary.LabelClosedAcrossMidRate)
			row.LabelMidTouchBeforeBoundaryTouchRateMin = minFloat(row.LabelMidTouchBeforeBoundaryTouchRateMin, summary.LabelMidTouchBeforeBoundaryTouchRate)
			row.LabelMidTouchBeforeBoundaryTouchRateMax = maxFloat(row.LabelMidTouchBeforeBoundaryTouchRateMax, summary.LabelMidTouchBeforeBoundaryTouchRate)
			row.LabelMidCrossBeforeBoundaryBreakRateMin = minFloat(row.LabelMidCrossBeforeBoundaryBreakRateMin, summary.LabelMidCrossBeforeBoundaryBreakRate)
			row.LabelMidCrossBeforeBoundaryBreakRateMax = maxFloat(row.LabelMidCrossBeforeBoundaryBreakRateMax, summary.LabelMidCrossBeforeBoundaryBreakRate)
			row.LabelReenteredRangeRateMin = minFloat(row.LabelReenteredRangeRateMin, summary.LabelReenteredRangeRate)
			row.LabelReenteredRangeRateMax = maxFloat(row.LabelReenteredRangeRateMax, summary.LabelReenteredRangeRate)
			row.LabelPersistedInsideRangeRateMin = minFloat(row.LabelPersistedInsideRangeRateMin, summary.LabelPersistedInsideRangeRate)
			row.LabelPersistedInsideRangeRateMax = maxFloat(row.LabelPersistedInsideRangeRateMax, summary.LabelPersistedInsideRangeRate)
			row.LabelQuickInvalidatedRateMin = minFloat(row.LabelQuickInvalidatedRateMin, summary.LabelQuickInvalidatedRate)
			row.LabelQuickInvalidatedRateMax = maxFloat(row.LabelQuickInvalidatedRateMax, summary.LabelQuickInvalidatedRate)
			row.LabelInvalidatedUpRateMin = minFloat(row.LabelInvalidatedUpRateMin, summary.LabelInvalidatedUpRate)
			row.LabelInvalidatedUpRateMax = maxFloat(row.LabelInvalidatedUpRateMax, summary.LabelInvalidatedUpRate)
			row.LabelInvalidatedDownRateMin = minFloat(row.LabelInvalidatedDownRateMin, summary.LabelInvalidatedDownRate)
			row.LabelInvalidatedDownRateMax = maxFloat(row.LabelInvalidatedDownRateMax, summary.LabelInvalidatedDownRate)
			row.LabelTrendedUpRateMin = minFloat(row.LabelTrendedUpRateMin, summary.LabelTrendedUpRate)
			row.LabelTrendedUpRateMax = maxFloat(row.LabelTrendedUpRateMax, summary.LabelTrendedUpRate)
			row.LabelTrendedDownRateMin = minFloat(row.LabelTrendedDownRateMin, summary.LabelTrendedDownRate)
			row.LabelTrendedDownRateMax = maxFloat(row.LabelTrendedDownRateMax, summary.LabelTrendedDownRate)
			row.LabelTrendedRateMin = minFloat(row.LabelTrendedRateMin, trendedRate)
			row.LabelTrendedRateMax = maxFloat(row.LabelTrendedRateMax, trendedRate)
			row.LabelAvgFirstMidTouchDelayBarsMin = minFloat(row.LabelAvgFirstMidTouchDelayBarsMin, summary.LabelAvgFirstMidTouchDelayBars)
			row.LabelAvgFirstMidTouchDelayBarsMax = maxFloat(row.LabelAvgFirstMidTouchDelayBarsMax, summary.LabelAvgFirstMidTouchDelayBars)
			row.LabelAvgFirstMidCloseAcrossDelayBarsMin = minFloat(row.LabelAvgFirstMidCloseAcrossDelayBarsMin, summary.LabelAvgFirstMidCloseAcrossDelayBars)
			row.LabelAvgFirstMidCloseAcrossDelayBarsMax = maxFloat(row.LabelAvgFirstMidCloseAcrossDelayBarsMax, summary.LabelAvgFirstMidCloseAcrossDelayBars)
		}
		row.SourceEpisodeCount += summary.SourceEpisodeCount
		row.CandidateCount += summary.CandidateCount
	}
	row.SourceEpisodeCountDelta = row.SourceEpisodeCountMax - row.SourceEpisodeCountMin
	row.CandidateCountDelta = row.CandidateCountMax - row.CandidateCountMin
	row.CandidateRateDelta = row.CandidateRateMax - row.CandidateRateMin
	row.LabelTouchedMidRateDelta = row.LabelTouchedMidRateMax - row.LabelTouchedMidRateMin
	row.LabelClosedAcrossMidRateDelta = row.LabelClosedAcrossMidRateMax - row.LabelClosedAcrossMidRateMin
	row.LabelMidTouchBeforeBoundaryTouchRateDelta = row.LabelMidTouchBeforeBoundaryTouchRateMax - row.LabelMidTouchBeforeBoundaryTouchRateMin
	row.LabelMidCrossBeforeBoundaryBreakRateDelta = row.LabelMidCrossBeforeBoundaryBreakRateMax - row.LabelMidCrossBeforeBoundaryBreakRateMin
	row.LabelReenteredRangeRateDelta = row.LabelReenteredRangeRateMax - row.LabelReenteredRangeRateMin
	row.LabelPersistedInsideRangeRateDelta = row.LabelPersistedInsideRangeRateMax - row.LabelPersistedInsideRangeRateMin
	row.LabelQuickInvalidatedRateDelta = row.LabelQuickInvalidatedRateMax - row.LabelQuickInvalidatedRateMin
	row.LabelInvalidatedUpRateDelta = row.LabelInvalidatedUpRateMax - row.LabelInvalidatedUpRateMin
	row.LabelInvalidatedDownRateDelta = row.LabelInvalidatedDownRateMax - row.LabelInvalidatedDownRateMin
	row.LabelTrendedUpRateDelta = row.LabelTrendedUpRateMax - row.LabelTrendedUpRateMin
	row.LabelTrendedDownRateDelta = row.LabelTrendedDownRateMax - row.LabelTrendedDownRateMin
	row.LabelTrendedRateDelta = row.LabelTrendedRateMax - row.LabelTrendedRateMin
	row.LabelAvgFirstMidTouchDelayBarsDelta = row.LabelAvgFirstMidTouchDelayBarsMax - row.LabelAvgFirstMidTouchDelayBarsMin
	row.LabelAvgFirstMidCloseAcrossDelayBarsDelta = row.LabelAvgFirstMidCloseAcrossDelayBarsMax - row.LabelAvgFirstMidCloseAcrossDelayBarsMin
	return row
}

func holdInsideMidlineTransitionSortedCombos(comboSet map[holdInsideMidlineTransitionSummaryCombo]bool) []holdInsideMidlineTransitionSummaryCombo {
	if len(comboSet) == 0 {
		return nil
	}
	combos := make([]holdInsideMidlineTransitionSummaryCombo, 0, len(comboSet))
	for combo := range comboSet {
		combos = append(combos, combo)
	}
	sort.Slice(combos, func(i, j int) bool {
		if holdInsideDecisionMidSideSortKey(combos[i].midSide) != holdInsideDecisionMidSideSortKey(combos[j].midSide) {
			return holdInsideDecisionMidSideSortKey(combos[i].midSide) < holdInsideDecisionMidSideSortKey(combos[j].midSide)
		}
		return holdInsideDecisionPositionBucketSortKey(combos[i].bucket) < holdInsideDecisionPositionBucketSortKey(combos[j].bucket)
	})
	return combos
}

func lessHoldInsideMidlineTransitionCandidate(a, b HoldInsideMidlineTransitionCandidateRow) bool {
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
	return a.HorizonBars < b.HorizonBars
}

func lessHoldInsideMidlineTransitionSummary(a, b HoldInsideMidlineTransitionSummaryRow) bool {
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
	if holdInsideDecisionMidSideSortKey(a.DecisionMidSide) != holdInsideDecisionMidSideSortKey(b.DecisionMidSide) {
		return holdInsideDecisionMidSideSortKey(a.DecisionMidSide) < holdInsideDecisionMidSideSortKey(b.DecisionMidSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.DecisionClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.DecisionClosePositionBucket)
}

func lessHoldInsideMidlineTransitionStability(a, b HoldInsideMidlineTransitionStabilityRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if holdInsideDecisionMidSideSortKey(a.DecisionMidSide) != holdInsideDecisionMidSideSortKey(b.DecisionMidSide) {
		return holdInsideDecisionMidSideSortKey(a.DecisionMidSide) < holdInsideDecisionMidSideSortKey(b.DecisionMidSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.DecisionClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.DecisionClosePositionBucket)
}

func holdInsideDecisionMidSideSortKey(side string) int {
	switch side {
	case holdInsideDecisionMidSideAll:
		return 0
	case holdInsideDecisionMidSideBelow:
		return 1
	case holdInsideDecisionMidSideAt:
		return 2
	case holdInsideDecisionMidSideAbove:
		return 3
	default:
		return 99
	}
}
