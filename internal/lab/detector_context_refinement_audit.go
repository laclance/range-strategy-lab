package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	DetectorContextRuleEpisodeEnd       = "episode_end"
	DetectorContextRuleHold1Inside      = "hold_1_inside"
	DetectorContextRuleHold3Inside      = "hold_3_inside"
	DetectorContextRuleHold6Inside      = "hold_6_inside"
	DetectorContextRuleHold3InsideMid50 = "hold_3_inside_mid_50"
	decisionClosePositionBucketUnknown  = "unknown"
	decisionClosePositionBucketBelow    = "below_range"
	decisionClosePositionBucketLow25    = "low_25"
	decisionClosePositionBucketMid50    = "mid_50"
	decisionClosePositionBucketHigh25   = "high_25"
	decisionClosePositionBucketAbove    = "above_range"
)

type DetectorContextRefinementAuditConfig struct {
	HorizonsBars          []int
	QuickInvalidationBars int
	Profiles              []DetectorSweepProfile
	ContextRules          []DetectorContextRefinementRule
}

type DetectorContextRefinementRule struct {
	RuleID       string `json:"context_rule"`
	HoldBars     int    `json:"hold_bars"`
	RequireMid50 bool   `json:"require_mid_50"`
}

type DetectorContextRefinementCandidateRow struct {
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
	HorizonBars                 int     `json:"horizon_bars"`
	LabelWindowStartIndex       int     `json:"label_window_start_index"`
	LabelWindowEndIndex         int     `json:"label_window_end_index"`
	LabelWindowStartTime        string  `json:"label_window_start_time"`
	LabelWindowEndTime          string  `json:"label_window_end_time"`
	LabelReenteredRange         bool    `json:"label_reentered_range"`
	LabelPersistedInsideRange   bool    `json:"label_persisted_inside_range"`
	LabelQuickInvalidated       bool    `json:"label_quick_invalidated"`
	LabelInvalidatedUp          bool    `json:"label_invalidated_up"`
	LabelInvalidatedDown        bool    `json:"label_invalidated_down"`
	LabelChopped                bool    `json:"label_chopped"`
	LabelTrendedUp              bool    `json:"label_trended_up"`
	LabelTrendedDown            bool    `json:"label_trended_down"`
	LabelCloseDriftPct          float64 `json:"label_close_drift_pct"`
	LabelMaxUpMovePct           float64 `json:"label_max_up_move_pct"`
	LabelMaxDownMovePct         float64 `json:"label_max_down_move_pct"`
}

type DetectorContextRefinementSummaryRow struct {
	ProfileID                      string  `json:"profile_id"`
	IsBalancedBaseline             bool    `json:"is_balanced_baseline"`
	IsADXComparison                bool    `json:"is_adx_comparison"`
	Percentile                     float64 `json:"percentile"`
	MinConsecutiveBars             int     `json:"min_consecutive_bars"`
	UseBollinger                   bool    `json:"use_bollinger"`
	UseADX                         bool    `json:"use_adx"`
	LookbackDays                   int     `json:"lookback_days"`
	ContextRule                    string  `json:"context_rule"`
	HoldBars                       int     `json:"hold_bars"`
	RequireMid50                   bool    `json:"require_mid_50"`
	Split                          string  `json:"split"`
	HorizonBars                    int     `json:"horizon_bars"`
	SourceEpisodeCount             int     `json:"source_episode_count"`
	CandidateCount                 int     `json:"candidate_count"`
	CandidateRate                  float64 `json:"candidate_rate"`
	AvgRawLengthBars               float64 `json:"avg_raw_length_bars"`
	AvgActiveLengthBars            float64 `json:"avg_active_length_bars"`
	AvgEpisodeWidthPct             float64 `json:"avg_episode_width_pct"`
	AvgNormalizedATR               float64 `json:"avg_normalized_atr"`
	AvgEndNormalizedATR            float64 `json:"avg_end_normalized_atr"`
	AvgWidthToATRRatio             float64 `json:"avg_width_to_atr_ratio"`
	AvgDecisionClosePosition       float64 `json:"avg_decision_close_position"`
	LabelReenteredRangeCount       int     `json:"label_reentered_range_count"`
	LabelPersistedInsideRangeCount int     `json:"label_persisted_inside_range_count"`
	LabelQuickInvalidatedCount     int     `json:"label_quick_invalidated_count"`
	LabelInvalidatedUpCount        int     `json:"label_invalidated_up_count"`
	LabelInvalidatedDownCount      int     `json:"label_invalidated_down_count"`
	LabelChoppedCount              int     `json:"label_chopped_count"`
	LabelTrendedUpCount            int     `json:"label_trended_up_count"`
	LabelTrendedDownCount          int     `json:"label_trended_down_count"`
	LabelReenteredRangeRate        float64 `json:"label_reentered_range_rate"`
	LabelPersistedInsideRangeRate  float64 `json:"label_persisted_inside_range_rate"`
	LabelQuickInvalidatedRate      float64 `json:"label_quick_invalidated_rate"`
	LabelInvalidatedUpRate         float64 `json:"label_invalidated_up_rate"`
	LabelInvalidatedDownRate       float64 `json:"label_invalidated_down_rate"`
	LabelChoppedRate               float64 `json:"label_chopped_rate"`
	LabelTrendedUpRate             float64 `json:"label_trended_up_rate"`
	LabelTrendedDownRate           float64 `json:"label_trended_down_rate"`
	LabelAvgCloseDriftPct          float64 `json:"label_avg_close_drift_pct"`
	LabelAvgMaxUpMovePct           float64 `json:"label_avg_max_up_move_pct"`
	LabelAvgMaxDownMovePct         float64 `json:"label_avg_max_down_move_pct"`
}

type DetectorContextRefinementStabilityRow struct {
	ProfileID                          string  `json:"profile_id"`
	IsBalancedBaseline                 bool    `json:"is_balanced_baseline"`
	IsADXComparison                    bool    `json:"is_adx_comparison"`
	Percentile                         float64 `json:"percentile"`
	MinConsecutiveBars                 int     `json:"min_consecutive_bars"`
	UseBollinger                       bool    `json:"use_bollinger"`
	UseADX                             bool    `json:"use_adx"`
	LookbackDays                       int     `json:"lookback_days"`
	ContextRule                        string  `json:"context_rule"`
	HoldBars                           int     `json:"hold_bars"`
	RequireMid50                       bool    `json:"require_mid_50"`
	HorizonBars                        int     `json:"horizon_bars"`
	PeriodSplits                       int     `json:"period_splits"`
	SourceEpisodeCount                 int     `json:"source_episode_count"`
	SourceEpisodeCountMin              int     `json:"source_episode_count_min"`
	SourceEpisodeCountMax              int     `json:"source_episode_count_max"`
	SourceEpisodeCountDelta            int     `json:"source_episode_count_delta"`
	CandidateCount                     int     `json:"candidate_count"`
	CandidateCountMin                  int     `json:"candidate_count_min"`
	CandidateCountMax                  int     `json:"candidate_count_max"`
	CandidateCountDelta                int     `json:"candidate_count_delta"`
	CandidateRateMin                   float64 `json:"candidate_rate_min"`
	CandidateRateMax                   float64 `json:"candidate_rate_max"`
	CandidateRateDelta                 float64 `json:"candidate_rate_delta"`
	LabelPersistedInsideRangeRateMin   float64 `json:"label_persisted_inside_range_rate_min"`
	LabelPersistedInsideRangeRateMax   float64 `json:"label_persisted_inside_range_rate_max"`
	LabelPersistedInsideRangeRateDelta float64 `json:"label_persisted_inside_range_rate_delta"`
	LabelQuickInvalidatedRateMin       float64 `json:"label_quick_invalidated_rate_min"`
	LabelQuickInvalidatedRateMax       float64 `json:"label_quick_invalidated_rate_max"`
	LabelQuickInvalidatedRateDelta     float64 `json:"label_quick_invalidated_rate_delta"`
	LabelChoppedRateMin                float64 `json:"label_chopped_rate_min"`
	LabelChoppedRateMax                float64 `json:"label_chopped_rate_max"`
	LabelChoppedRateDelta              float64 `json:"label_chopped_rate_delta"`
	LabelTrendedRateMin                float64 `json:"label_trended_rate_min"`
	LabelTrendedRateMax                float64 `json:"label_trended_rate_max"`
	LabelTrendedRateDelta              float64 `json:"label_trended_rate_delta"`
	LabelAvgCloseDriftPctMin           float64 `json:"label_avg_close_drift_pct_min"`
	LabelAvgCloseDriftPctMax           float64 `json:"label_avg_close_drift_pct_max"`
	LabelAvgCloseDriftPctDelta         float64 `json:"label_avg_close_drift_pct_delta"`
}

func DefaultDetectorContextRefinementAuditConfig() DetectorContextRefinementAuditConfig {
	return DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1, 3, 6, 12},
		QuickInvalidationBars: 3,
		Profiles:              DefaultDetectorContextRefinementProfiles(DefaultCompressionRangeDetectorConfig().LookbackDays),
		ContextRules:          defaultDetectorContextRefinementRules(),
	}
}

func DefaultDetectorContextRefinementProfiles(lookbackDays int) []DetectorSweepProfile {
	if lookbackDays <= 0 {
		lookbackDays = DefaultCompressionRangeDetectorConfig().LookbackDays
	}
	return []DetectorSweepProfile{
		newDetectorSweepProfile(0.30, 12, true, false, lookbackDays, false),
		newDetectorSweepProfile(0.30, 12, true, true, lookbackDays, true),
		newDetectorSweepProfile(0.20, 24, true, false, lookbackDays, false),
		newDetectorSweepProfile(0.30, 24, true, false, lookbackDays, false),
		newDetectorSweepProfile(0.40, 24, true, false, lookbackDays, false),
		newDetectorSweepProfile(0.20, 24, true, true, lookbackDays, false),
		newDetectorSweepProfile(0.30, 24, true, true, lookbackDays, false),
		newDetectorSweepProfile(0.40, 24, true, true, lookbackDays, false),
	}
}

func RunDetectorContextRefinementAudit(candles []Candle, detectorCfg RangeDetectorConfig, cfg DetectorContextRefinementAuditConfig, splits []Split) ([]DetectorContextRefinementCandidateRow, []DetectorContextRefinementSummaryRow, []DetectorContextRefinementStabilityRow, error) {
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

	candidates := []DetectorContextRefinementCandidateRow{}
	summaryAccumulators := map[detectorContextSummaryKey]*detectorContextSummaryAccumulator{}
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
		profileCandidates := runDetectorContextRefinementForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
		candidates = append(candidates, profileCandidates...)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return lessDetectorContextCandidate(candidates[i], candidates[j])
	})
	summaryRows := detectorContextSummaryRows(summaryAccumulators)
	stabilityRows := detectorContextStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func runDetectorContextRefinementFromClassifications(candles []Candle, profile DetectorSweepProfile, classifications []RangeClassification, cfg DetectorContextRefinementAuditConfig, splits []Split) ([]DetectorContextRefinementCandidateRow, []DetectorContextRefinementSummaryRow, []DetectorContextRefinementStabilityRow, error) {
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
	summaryAccumulators := map[detectorContextSummaryKey]*detectorContextSummaryAccumulator{}
	candidates := runDetectorContextRefinementForProfile(candles, classifications, normalizedATR, profile, cfg, splits, summaryAccumulators)
	sort.Slice(candidates, func(i, j int) bool {
		return lessDetectorContextCandidate(candidates[i], candidates[j])
	})
	summaryRows := detectorContextSummaryRows(summaryAccumulators)
	stabilityRows := detectorContextStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, summaryRows, stabilityRows, nil
}

func (cfg DetectorContextRefinementAuditConfig) withDefaults(lookbackDays int) DetectorContextRefinementAuditConfig {
	defaults := DefaultDetectorContextRefinementAuditConfig()
	if lookbackDays > 0 {
		defaults.Profiles = DefaultDetectorContextRefinementProfiles(lookbackDays)
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

func (cfg DetectorContextRefinementAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("detector context refinement horizon bars must be positive")
		}
	}
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("detector context refinement quick invalidation bars must be positive")
	}
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("detector context refinement profiles must not be empty")
	}
	for _, profile := range cfg.Profiles {
		if profile.ProfileID == "" {
			return fmt.Errorf("detector context refinement profile id must not be empty")
		}
	}
	if len(cfg.ContextRules) == 0 {
		return fmt.Errorf("detector context refinement rules must not be empty")
	}
	for _, rule := range cfg.ContextRules {
		if rule.RuleID == "" {
			return fmt.Errorf("detector context refinement rule id must not be empty")
		}
		if rule.HoldBars < 0 {
			return fmt.Errorf("detector context refinement hold bars must not be negative")
		}
		if rule.RequireMid50 && rule.HoldBars == 0 {
			return fmt.Errorf("detector context refinement mid-50 rule requires positive hold bars")
		}
	}
	return nil
}

func defaultDetectorContextRefinementRules() []DetectorContextRefinementRule {
	return []DetectorContextRefinementRule{
		{RuleID: DetectorContextRuleEpisodeEnd, HoldBars: 0},
		{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1},
		{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
		{RuleID: DetectorContextRuleHold6Inside, HoldBars: 6},
		{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
	}
}

func detectorContextPercentiles(profiles []DetectorSweepProfile) []float64 {
	seen := map[float64]bool{}
	var out []float64
	for _, profile := range profiles {
		if seen[profile.Percentile] {
			continue
		}
		seen[profile.Percentile] = true
		out = append(out, profile.Percentile)
	}
	sort.Float64s(out)
	return out
}

func runDetectorContextRefinementForProfile(candles []Candle, classifications []RangeClassification, normalizedATR []float64, profile DetectorSweepProfile, cfg DetectorContextRefinementAuditConfig, splits []Split, summaryAccumulators map[detectorContextSummaryKey]*detectorContextSummaryAccumulator) []DetectorContextRefinementCandidateRow {
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, profile.ProfileID)
	candidates := []DetectorContextRefinementCandidateRow{}
	for _, episode := range episodes {
		for _, rule := range cfg.ContextRules {
			decisionIndex := episode.EndIndex + rule.HoldBars
			if decisionIndex < 0 || decisionIndex >= len(candles) {
				continue
			}
			split := splitNameForCloseTime(candles[decisionIndex].CloseTime, splits)
			passesRule := detectorContextRulePasses(candles, episode, rule, decisionIndex)
			for _, horizon := range cfg.HorizonsBars {
				if decisionIndex+horizon >= len(candles) {
					continue
				}
				addDetectorContextSource(summaryAccumulators, profile, rule, split, horizon)
				if split != fullSplitName {
					addDetectorContextSource(summaryAccumulators, profile, rule, fullSplitName, horizon)
				}
				if !passesRule {
					continue
				}
				row := newDetectorContextCandidateRow(candles, profile, rule, episode, split, decisionIndex, horizon, cfg.QuickInvalidationBars)
				candidates = append(candidates, row)
				addDetectorContextCandidate(summaryAccumulators, row, split)
				if split != fullSplitName {
					addDetectorContextCandidate(summaryAccumulators, row, fullSplitName)
				}
			}
		}
	}
	return candidates
}

func detectorContextRulePasses(candles []Candle, episode rangeRegimeDurabilityEpisode, rule DetectorContextRefinementRule, decisionIndex int) bool {
	if rule.HoldBars > 0 {
		for i := episode.EndIndex + 1; i <= decisionIndex; i++ {
			if i < 0 || i >= len(candles) || !closeInsideRange(candles[i].Close, episode.Low, episode.High) {
				return false
			}
		}
	}
	if rule.RequireMid50 {
		position := decisionClosePosition(candles[decisionIndex].Close, episode.Low, episode.High)
		return validNumber(position) && position >= 0.25 && position <= 0.75
	}
	return true
}

func newDetectorContextCandidateRow(candles []Candle, profile DetectorSweepProfile, rule DetectorContextRefinementRule, episode rangeRegimeDurabilityEpisode, split string, decisionIndex, horizon, quickInvalidationBars int) DetectorContextRefinementCandidateRow {
	decision := candles[decisionIndex]
	decisionPosition := decisionClosePosition(decision.Close, episode.Low, episode.High)
	decisionPositionValue := decisionPosition
	if !validNumber(decisionPositionValue) {
		decisionPositionValue = 0
	}
	startIndex := decisionIndex + 1
	endIndex := decisionIndex + horizon
	future := candles[startIndex : endIndex+1]
	futureMaxHigh, futureMinLow := futureHighLow(future)
	lastClose := candles[endIndex].Close

	row := DetectorContextRefinementCandidateRow{
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
		HorizonBars:                 horizon,
		LabelWindowStartIndex:       startIndex,
		LabelWindowEndIndex:         endIndex,
		LabelWindowStartTime:        candles[startIndex].CloseTime.Format(timeLayout),
		LabelWindowEndTime:          candles[endIndex].CloseTime.Format(timeLayout),
		LabelCloseDriftPct:          closeDriftPct(lastClose, decision.Close),
		LabelMaxUpMovePct:           movePct(math.Max(0, futureMaxHigh-decision.Close), decision.Close),
		LabelMaxDownMovePct:         movePct(math.Max(0, decision.Close-futureMinLow), decision.Close),
	}

	persistedInside := true
	quickLimit := minInt(len(future), quickInvalidationBars)
	for i, candle := range future {
		inside := closeInsideRange(candle.Close, episode.Low, episode.High)
		if inside {
			row.LabelReenteredRange = true
		} else {
			persistedInside = false
		}
		if candle.Close > episode.High {
			row.LabelInvalidatedUp = true
			if i < quickLimit {
				row.LabelQuickInvalidated = true
			}
		}
		if candle.Close < episode.Low {
			row.LabelInvalidatedDown = true
			if i < quickLimit {
				row.LabelQuickInvalidated = true
			}
		}
	}
	row.LabelPersistedInsideRange = persistedInside
	row.LabelTrendedUp = row.LabelInvalidatedUp && lastClose > episode.High && row.LabelMaxUpMovePct > row.LabelMaxDownMovePct
	row.LabelTrendedDown = row.LabelInvalidatedDown && lastClose < episode.Low && row.LabelMaxDownMovePct > row.LabelMaxUpMovePct
	row.LabelChopped = row.LabelReenteredRange &&
		(row.LabelInvalidatedUp || row.LabelInvalidatedDown) &&
		!row.LabelTrendedUp &&
		!row.LabelTrendedDown
	return row
}

func closeInsideRange(close, low, high float64) bool {
	return close >= low && close <= high
}

func decisionClosePosition(close, low, high float64) float64 {
	if high <= low {
		return math.NaN()
	}
	return (close - low) / (high - low)
}

func decisionClosePositionBucket(position float64) string {
	switch {
	case !validNumber(position):
		return decisionClosePositionBucketUnknown
	case position < 0:
		return decisionClosePositionBucketBelow
	case position < 0.25:
		return decisionClosePositionBucketLow25
	case position <= 0.75:
		return decisionClosePositionBucketMid50
	case position <= 1:
		return decisionClosePositionBucketHigh25
	default:
		return decisionClosePositionBucketAbove
	}
}

type detectorContextSummaryKey struct {
	profileID   string
	contextRule string
	split       string
	horizonBars int
}

type detectorContextSummaryAccumulator struct {
	profile               DetectorSweepProfile
	rule                  DetectorContextRefinementRule
	split                 string
	horizonBars           int
	sourceEpisodes        int
	candidates            int
	rawLengthSum          float64
	activeLengthSum       float64
	widthPctSum           float64
	avgATRSum             float64
	endATRSum             float64
	widthToATRRatioSum    float64
	decisionPositionSum   float64
	reenteredRanges       int
	persistedInsideRanges int
	quickInvalidated      int
	invalidatedUp         int
	invalidatedDown       int
	chopped               int
	trendedUp             int
	trendedDown           int
	closeDriftPctSum      float64
	maxUpMovePctSum       float64
	maxDownMovePctSum     float64
}

func addDetectorContextSource(accumulators map[detectorContextSummaryKey]*detectorContextSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int) {
	acc := detectorContextSummaryAccumulatorFor(accumulators, profile, rule, split, horizon)
	acc.sourceEpisodes++
}

func addDetectorContextCandidate(accumulators map[detectorContextSummaryKey]*detectorContextSummaryAccumulator, row DetectorContextRefinementCandidateRow, split string) {
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
	acc := detectorContextSummaryAccumulatorFor(accumulators, profile, rule, split, row.HorizonBars)
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
	if row.LabelReenteredRange {
		acc.reenteredRanges++
	}
	if row.LabelPersistedInsideRange {
		acc.persistedInsideRanges++
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
	if row.LabelChopped {
		acc.chopped++
	}
	if row.LabelTrendedUp {
		acc.trendedUp++
	}
	if row.LabelTrendedDown {
		acc.trendedDown++
	}
	acc.closeDriftPctSum += row.LabelCloseDriftPct
	acc.maxUpMovePctSum += row.LabelMaxUpMovePct
	acc.maxDownMovePctSum += row.LabelMaxDownMovePct
}

func detectorContextSummaryAccumulatorFor(accumulators map[detectorContextSummaryKey]*detectorContextSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, split string, horizon int) *detectorContextSummaryAccumulator {
	key := detectorContextSummaryKey{profileID: profile.ProfileID, contextRule: rule.RuleID, split: split, horizonBars: horizon}
	acc := accumulators[key]
	if acc == nil {
		acc = &detectorContextSummaryAccumulator{
			profile:     profile,
			rule:        rule,
			split:       split,
			horizonBars: horizon,
		}
		accumulators[key] = acc
	}
	return acc
}

func detectorContextSummaryRows(accumulators map[detectorContextSummaryKey]*detectorContextSummaryAccumulator) []DetectorContextRefinementSummaryRow {
	rows := make([]DetectorContextRefinementSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessDetectorContextSummary(rows[i], rows[j])
	})
	return rows
}

func (acc detectorContextSummaryAccumulator) row() DetectorContextRefinementSummaryRow {
	row := DetectorContextRefinementSummaryRow{
		ProfileID:                      acc.profile.ProfileID,
		IsBalancedBaseline:             acc.profile.IsBalancedBaseline,
		IsADXComparison:                acc.profile.IsADXComparison,
		Percentile:                     acc.profile.Percentile,
		MinConsecutiveBars:             acc.profile.MinConsecutiveBars,
		UseBollinger:                   acc.profile.UseBollinger,
		UseADX:                         acc.profile.UseADX,
		LookbackDays:                   acc.profile.LookbackDays,
		ContextRule:                    acc.rule.RuleID,
		HoldBars:                       acc.rule.HoldBars,
		RequireMid50:                   acc.rule.RequireMid50,
		Split:                          acc.split,
		HorizonBars:                    acc.horizonBars,
		SourceEpisodeCount:             acc.sourceEpisodes,
		CandidateCount:                 acc.candidates,
		LabelReenteredRangeCount:       acc.reenteredRanges,
		LabelPersistedInsideRangeCount: acc.persistedInsideRanges,
		LabelQuickInvalidatedCount:     acc.quickInvalidated,
		LabelInvalidatedUpCount:        acc.invalidatedUp,
		LabelInvalidatedDownCount:      acc.invalidatedDown,
		LabelChoppedCount:              acc.chopped,
		LabelTrendedUpCount:            acc.trendedUp,
		LabelTrendedDownCount:          acc.trendedDown,
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
	row.LabelReenteredRangeRate = float64(acc.reenteredRanges) / count
	row.LabelPersistedInsideRangeRate = float64(acc.persistedInsideRanges) / count
	row.LabelQuickInvalidatedRate = float64(acc.quickInvalidated) / count
	row.LabelInvalidatedUpRate = float64(acc.invalidatedUp) / count
	row.LabelInvalidatedDownRate = float64(acc.invalidatedDown) / count
	row.LabelChoppedRate = float64(acc.chopped) / count
	row.LabelTrendedUpRate = float64(acc.trendedUp) / count
	row.LabelTrendedDownRate = float64(acc.trendedDown) / count
	row.LabelAvgCloseDriftPct = acc.closeDriftPctSum / count
	row.LabelAvgMaxUpMovePct = acc.maxUpMovePctSum / count
	row.LabelAvgMaxDownMovePct = acc.maxDownMovePctSum / count
	return row
}

func detectorContextStabilityRows(profiles []DetectorSweepProfile, rules []DetectorContextRefinementRule, horizons []int, summaryRows []DetectorContextRefinementSummaryRow, splits []Split) []DetectorContextRefinementStabilityRow {
	periodSplits := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplits = append(periodSplits, split)
		}
	}
	byKey := map[detectorContextSummaryKey]DetectorContextRefinementSummaryRow{}
	for _, row := range summaryRows {
		byKey[detectorContextSummaryKey{profileID: row.ProfileID, contextRule: row.ContextRule, split: row.Split, horizonBars: row.HorizonBars}] = row
	}

	rows := make([]DetectorContextRefinementStabilityRow, 0, len(profiles)*len(rules)*len(horizons))
	for _, profile := range profiles {
		for _, rule := range rules {
			for _, horizon := range horizons {
				rows = append(rows, newDetectorContextStabilityRow(profile, rule, horizon, periodSplits, byKey))
			}
		}
	}
	return rows
}

func newDetectorContextStabilityRow(profile DetectorSweepProfile, rule DetectorContextRefinementRule, horizon int, periodSplits []Split, byKey map[detectorContextSummaryKey]DetectorContextRefinementSummaryRow) DetectorContextRefinementStabilityRow {
	row := DetectorContextRefinementStabilityRow{
		ProfileID:          profile.ProfileID,
		IsBalancedBaseline: profile.IsBalancedBaseline,
		IsADXComparison:    profile.IsADXComparison,
		Percentile:         profile.Percentile,
		MinConsecutiveBars: profile.MinConsecutiveBars,
		UseBollinger:       profile.UseBollinger,
		UseADX:             profile.UseADX,
		LookbackDays:       profile.LookbackDays,
		ContextRule:        rule.RuleID,
		HoldBars:           rule.HoldBars,
		RequireMid50:       rule.RequireMid50,
		HorizonBars:        horizon,
		PeriodSplits:       len(periodSplits),
	}
	first := true
	for _, split := range periodSplits {
		summary := byKey[detectorContextSummaryKey{profileID: profile.ProfileID, contextRule: rule.RuleID, split: split.Name, horizonBars: horizon}]
		trendedRate := summary.LabelTrendedUpRate + summary.LabelTrendedDownRate
		if first {
			row.SourceEpisodeCountMin = summary.SourceEpisodeCount
			row.SourceEpisodeCountMax = summary.SourceEpisodeCount
			row.CandidateCountMin = summary.CandidateCount
			row.CandidateCountMax = summary.CandidateCount
			row.CandidateRateMin = summary.CandidateRate
			row.CandidateRateMax = summary.CandidateRate
			row.LabelPersistedInsideRangeRateMin = summary.LabelPersistedInsideRangeRate
			row.LabelPersistedInsideRangeRateMax = summary.LabelPersistedInsideRangeRate
			row.LabelQuickInvalidatedRateMin = summary.LabelQuickInvalidatedRate
			row.LabelQuickInvalidatedRateMax = summary.LabelQuickInvalidatedRate
			row.LabelChoppedRateMin = summary.LabelChoppedRate
			row.LabelChoppedRateMax = summary.LabelChoppedRate
			row.LabelTrendedRateMin = trendedRate
			row.LabelTrendedRateMax = trendedRate
			row.LabelAvgCloseDriftPctMin = summary.LabelAvgCloseDriftPct
			row.LabelAvgCloseDriftPctMax = summary.LabelAvgCloseDriftPct
			first = false
		} else {
			row.SourceEpisodeCountMin = minInt(row.SourceEpisodeCountMin, summary.SourceEpisodeCount)
			row.SourceEpisodeCountMax = maxInt(row.SourceEpisodeCountMax, summary.SourceEpisodeCount)
			row.CandidateCountMin = minInt(row.CandidateCountMin, summary.CandidateCount)
			row.CandidateCountMax = maxInt(row.CandidateCountMax, summary.CandidateCount)
			row.CandidateRateMin = minFloat(row.CandidateRateMin, summary.CandidateRate)
			row.CandidateRateMax = maxFloat(row.CandidateRateMax, summary.CandidateRate)
			row.LabelPersistedInsideRangeRateMin = minFloat(row.LabelPersistedInsideRangeRateMin, summary.LabelPersistedInsideRangeRate)
			row.LabelPersistedInsideRangeRateMax = maxFloat(row.LabelPersistedInsideRangeRateMax, summary.LabelPersistedInsideRangeRate)
			row.LabelQuickInvalidatedRateMin = minFloat(row.LabelQuickInvalidatedRateMin, summary.LabelQuickInvalidatedRate)
			row.LabelQuickInvalidatedRateMax = maxFloat(row.LabelQuickInvalidatedRateMax, summary.LabelQuickInvalidatedRate)
			row.LabelChoppedRateMin = minFloat(row.LabelChoppedRateMin, summary.LabelChoppedRate)
			row.LabelChoppedRateMax = maxFloat(row.LabelChoppedRateMax, summary.LabelChoppedRate)
			row.LabelTrendedRateMin = minFloat(row.LabelTrendedRateMin, trendedRate)
			row.LabelTrendedRateMax = maxFloat(row.LabelTrendedRateMax, trendedRate)
			row.LabelAvgCloseDriftPctMin = minFloat(row.LabelAvgCloseDriftPctMin, summary.LabelAvgCloseDriftPct)
			row.LabelAvgCloseDriftPctMax = maxFloat(row.LabelAvgCloseDriftPctMax, summary.LabelAvgCloseDriftPct)
		}
		row.SourceEpisodeCount += summary.SourceEpisodeCount
		row.CandidateCount += summary.CandidateCount
	}
	row.SourceEpisodeCountDelta = row.SourceEpisodeCountMax - row.SourceEpisodeCountMin
	row.CandidateCountDelta = row.CandidateCountMax - row.CandidateCountMin
	row.CandidateRateDelta = row.CandidateRateMax - row.CandidateRateMin
	row.LabelPersistedInsideRangeRateDelta = row.LabelPersistedInsideRangeRateMax - row.LabelPersistedInsideRangeRateMin
	row.LabelQuickInvalidatedRateDelta = row.LabelQuickInvalidatedRateMax - row.LabelQuickInvalidatedRateMin
	row.LabelChoppedRateDelta = row.LabelChoppedRateMax - row.LabelChoppedRateMin
	row.LabelTrendedRateDelta = row.LabelTrendedRateMax - row.LabelTrendedRateMin
	row.LabelAvgCloseDriftPctDelta = row.LabelAvgCloseDriftPctMax - row.LabelAvgCloseDriftPctMin
	return row
}

func lessDetectorContextCandidate(a, b DetectorContextRefinementCandidateRow) bool {
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

func lessDetectorContextSummary(a, b DetectorContextRefinementSummaryRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	return a.HorizonBars < b.HorizonBars
}

func detectorContextRuleSortKey(rule string) int {
	switch rule {
	case DetectorContextRuleEpisodeEnd:
		return 0
	case DetectorContextRuleHold1Inside:
		return 1
	case DetectorContextRuleHold3Inside:
		return 2
	case DetectorContextRuleHold6Inside:
		return 3
	case DetectorContextRuleHold3InsideMid50:
		return 4
	default:
		return 99
	}
}
