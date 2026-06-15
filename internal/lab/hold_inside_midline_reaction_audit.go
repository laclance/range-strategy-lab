package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	HoldInsideMidlineReactionEventTouch       = "mid_touch"
	HoldInsideMidlineReactionEventCloseAcross = "mid_close_across"
)

type HoldInsideMidlineReactionAuditConfig struct {
	HorizonsBars             []int
	QuickInvalidationBars    int
	MaxMidlineEventDelayBars int
	Profiles                 []DetectorSweepProfile
	ContextRules             []DetectorContextRefinementRule
}

type HoldInsideMidlineReactionCandidateRow struct {
	ProfileID                    string  `json:"profile_id"`
	IsBalancedBaseline           bool    `json:"is_balanced_baseline"`
	IsADXComparison              bool    `json:"is_adx_comparison"`
	Percentile                   float64 `json:"percentile"`
	MinConsecutiveBars           int     `json:"min_consecutive_bars"`
	UseBollinger                 bool    `json:"use_bollinger"`
	UseADX                       bool    `json:"use_adx"`
	LookbackDays                 int     `json:"lookback_days"`
	ContextRule                  string  `json:"context_rule"`
	HoldBars                     int     `json:"hold_bars"`
	RequireMid50                 bool    `json:"require_mid_50"`
	EventType                    string  `json:"event_type"`
	Split                        string  `json:"split"`
	SourceEpisodeID              int     `json:"source_episode_id"`
	EpisodeStartIndex            int     `json:"episode_start_index"`
	EpisodeEndIndex              int     `json:"episode_end_index"`
	EpisodeStartTime             string  `json:"episode_start_time"`
	EpisodeEndTime               string  `json:"episode_end_time"`
	RawLengthBars                int     `json:"raw_length_bars"`
	ActiveLengthBars             int     `json:"active_length_bars"`
	RawLengthBucket              string  `json:"raw_length_bucket"`
	ActiveLengthBucket           string  `json:"active_length_bucket"`
	EpisodeHigh                  float64 `json:"episode_high"`
	EpisodeLow                   float64 `json:"episode_low"`
	EpisodeMid                   float64 `json:"episode_mid"`
	EpisodeEndClose              float64 `json:"episode_end_close"`
	EpisodeWidthPct              float64 `json:"episode_width_pct"`
	EpisodeWidthBucket           string  `json:"episode_width_bucket"`
	AvgNormalizedATR             float64 `json:"avg_normalized_atr"`
	EndNormalizedATR             float64 `json:"end_normalized_atr"`
	WidthToATRRatio              float64 `json:"width_to_atr_ratio"`
	WidthToATRBucket             string  `json:"width_to_atr_bucket"`
	HoldDecisionIndex            int     `json:"hold_decision_index"`
	HoldDecisionTime             string  `json:"hold_decision_time"`
	HoldDecisionClose            float64 `json:"hold_decision_close"`
	HoldDecisionClosePosition    float64 `json:"hold_decision_close_position"`
	HoldDecisionCloseBucket      string  `json:"hold_decision_close_position_bucket"`
	HoldDecisionMidSide          string  `json:"hold_decision_mid_side"`
	EventIndex                   int     `json:"event_index"`
	EventTime                    string  `json:"event_time"`
	EventDelayBars               int     `json:"event_delay_bars"`
	EventClose                   float64 `json:"event_close"`
	EventClosePosition           float64 `json:"event_close_position"`
	EventClosePositionBucket     string  `json:"event_close_position_bucket"`
	EventMidSide                 string  `json:"event_mid_side"`
	EventDistanceToHighPct       float64 `json:"event_distance_to_high_pct"`
	EventDistanceToLowPct        float64 `json:"event_distance_to_low_pct"`
	EventDistanceToMidPct        float64 `json:"event_distance_to_mid_pct"`
	HorizonBars                  int     `json:"horizon_bars"`
	LabelWindowStartIndex        int     `json:"label_window_start_index"`
	LabelWindowEndIndex          int     `json:"label_window_end_index"`
	LabelWindowStartTime         string  `json:"label_window_start_time"`
	LabelWindowEndTime           string  `json:"label_window_end_time"`
	LabelReenteredRange          bool    `json:"label_reentered_range"`
	LabelPersistedInsideRange    bool    `json:"label_persisted_inside_range"`
	LabelQuickInvalidated        bool    `json:"label_quick_invalidated"`
	LabelInvalidatedUp           bool    `json:"label_invalidated_up"`
	LabelInvalidatedDown         bool    `json:"label_invalidated_down"`
	LabelTrendedUp               bool    `json:"label_trended_up"`
	LabelTrendedDown             bool    `json:"label_trended_down"`
	LabelTouchedHigh             bool    `json:"label_touched_high"`
	LabelTouchedLow              bool    `json:"label_touched_low"`
	LabelTouchedOppositeHalf     bool    `json:"label_touched_opposite_half"`
	LabelClosedBackAcrossMid     bool    `json:"label_closed_back_across_mid"`
	LabelMidRejectionBeforeBound bool    `json:"label_mid_rejection_before_boundary_touch"`
	LabelBoundBeforeMidRejection bool    `json:"label_boundary_touch_before_mid_rejection"`
}

type HoldInsideMidlineReactionFunnelSummaryRow struct {
	ProfileID             string  `json:"profile_id"`
	IsBalancedBaseline    bool    `json:"is_balanced_baseline"`
	IsADXComparison       bool    `json:"is_adx_comparison"`
	Percentile            float64 `json:"percentile"`
	MinConsecutiveBars    int     `json:"min_consecutive_bars"`
	UseBollinger          bool    `json:"use_bollinger"`
	UseADX                bool    `json:"use_adx"`
	LookbackDays          int     `json:"lookback_days"`
	ContextRule           string  `json:"context_rule"`
	HoldBars              int     `json:"hold_bars"`
	RequireMid50          bool    `json:"require_mid_50"`
	EventType             string  `json:"event_type"`
	Split                 string  `json:"split"`
	SourceHoldCount       int     `json:"source_hold_count"`
	EventCount            int     `json:"event_count"`
	MissingEventCount     int     `json:"missing_event_count"`
	MissingFutureCount    int     `json:"missing_future_count"`
	ReactionEligibleCount int     `json:"reaction_eligible_count"`
	EventRate             float64 `json:"event_rate"`
	ReactionEligibleRate  float64 `json:"reaction_eligible_rate"`
	AvgEventDelayBars     float64 `json:"avg_event_delay_bars"`
}

type HoldInsideMidlineReactionSummaryRow struct {
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
	EventType                         string  `json:"event_type"`
	Split                             string  `json:"split"`
	HorizonBars                       int     `json:"horizon_bars"`
	EventMidSide                      string  `json:"event_mid_side"`
	EventClosePositionBucket          string  `json:"event_close_position_bucket"`
	CandidateCount                    int     `json:"candidate_count"`
	AvgEventDelayBars                 float64 `json:"avg_event_delay_bars"`
	AvgEventClosePosition             float64 `json:"avg_event_close_position"`
	AvgEventDistanceToHighPct         float64 `json:"avg_event_distance_to_high_pct"`
	AvgEventDistanceToLowPct          float64 `json:"avg_event_distance_to_low_pct"`
	AvgEventDistanceToMidPct          float64 `json:"avg_event_distance_to_mid_pct"`
	LabelReenteredRangeCount          int     `json:"label_reentered_range_count"`
	LabelPersistedInsideRangeCount    int     `json:"label_persisted_inside_range_count"`
	LabelQuickInvalidatedCount        int     `json:"label_quick_invalidated_count"`
	LabelInvalidatedUpCount           int     `json:"label_invalidated_up_count"`
	LabelInvalidatedDownCount         int     `json:"label_invalidated_down_count"`
	LabelTrendedUpCount               int     `json:"label_trended_up_count"`
	LabelTrendedDownCount             int     `json:"label_trended_down_count"`
	LabelTouchedHighCount             int     `json:"label_touched_high_count"`
	LabelTouchedLowCount              int     `json:"label_touched_low_count"`
	LabelTouchedOppositeHalfCount     int     `json:"label_touched_opposite_half_count"`
	LabelClosedBackAcrossMidCount     int     `json:"label_closed_back_across_mid_count"`
	LabelMidRejectionBeforeBoundCount int     `json:"label_mid_rejection_before_boundary_touch_count"`
	LabelBoundBeforeMidRejectionCount int     `json:"label_boundary_touch_before_mid_rejection_count"`
	LabelReenteredRangeRate           float64 `json:"label_reentered_range_rate"`
	LabelPersistedInsideRangeRate     float64 `json:"label_persisted_inside_range_rate"`
	LabelQuickInvalidatedRate         float64 `json:"label_quick_invalidated_rate"`
	LabelInvalidatedUpRate            float64 `json:"label_invalidated_up_rate"`
	LabelInvalidatedDownRate          float64 `json:"label_invalidated_down_rate"`
	LabelTrendedUpRate                float64 `json:"label_trended_up_rate"`
	LabelTrendedDownRate              float64 `json:"label_trended_down_rate"`
	LabelTouchedHighRate              float64 `json:"label_touched_high_rate"`
	LabelTouchedLowRate               float64 `json:"label_touched_low_rate"`
	LabelTouchedOppositeHalfRate      float64 `json:"label_touched_opposite_half_rate"`
	LabelClosedBackAcrossMidRate      float64 `json:"label_closed_back_across_mid_rate"`
	LabelMidRejectionBeforeBoundRate  float64 `json:"label_mid_rejection_before_boundary_touch_rate"`
	LabelBoundBeforeMidRejectionRate  float64 `json:"label_boundary_touch_before_mid_rejection_rate"`
}

type HoldInsideMidlineReactionStabilityRow struct {
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
	EventType                             string  `json:"event_type"`
	HorizonBars                           int     `json:"horizon_bars"`
	EventMidSide                          string  `json:"event_mid_side"`
	EventClosePositionBucket              string  `json:"event_close_position_bucket"`
	PeriodSplits                          int     `json:"period_splits"`
	CandidateCount                        int     `json:"candidate_count"`
	CandidateCountMin                     int     `json:"candidate_count_min"`
	CandidateCountMax                     int     `json:"candidate_count_max"`
	CandidateCountDelta                   int     `json:"candidate_count_delta"`
	AvgEventDelayBarsMin                  float64 `json:"avg_event_delay_bars_min"`
	AvgEventDelayBarsMax                  float64 `json:"avg_event_delay_bars_max"`
	AvgEventDelayBarsDelta                float64 `json:"avg_event_delay_bars_delta"`
	LabelReenteredRangeRateMin            float64 `json:"label_reentered_range_rate_min"`
	LabelReenteredRangeRateMax            float64 `json:"label_reentered_range_rate_max"`
	LabelReenteredRangeRateDelta          float64 `json:"label_reentered_range_rate_delta"`
	LabelPersistedInsideRangeRateMin      float64 `json:"label_persisted_inside_range_rate_min"`
	LabelPersistedInsideRangeRateMax      float64 `json:"label_persisted_inside_range_rate_max"`
	LabelPersistedInsideRangeRateDelta    float64 `json:"label_persisted_inside_range_rate_delta"`
	LabelQuickInvalidatedRateMin          float64 `json:"label_quick_invalidated_rate_min"`
	LabelQuickInvalidatedRateMax          float64 `json:"label_quick_invalidated_rate_max"`
	LabelQuickInvalidatedRateDelta        float64 `json:"label_quick_invalidated_rate_delta"`
	LabelTrendedRateMin                   float64 `json:"label_trended_rate_min"`
	LabelTrendedRateMax                   float64 `json:"label_trended_rate_max"`
	LabelTrendedRateDelta                 float64 `json:"label_trended_rate_delta"`
	LabelTouchedHighRateMin               float64 `json:"label_touched_high_rate_min"`
	LabelTouchedHighRateMax               float64 `json:"label_touched_high_rate_max"`
	LabelTouchedHighRateDelta             float64 `json:"label_touched_high_rate_delta"`
	LabelTouchedLowRateMin                float64 `json:"label_touched_low_rate_min"`
	LabelTouchedLowRateMax                float64 `json:"label_touched_low_rate_max"`
	LabelTouchedLowRateDelta              float64 `json:"label_touched_low_rate_delta"`
	LabelTouchedOppositeHalfRateMin       float64 `json:"label_touched_opposite_half_rate_min"`
	LabelTouchedOppositeHalfRateMax       float64 `json:"label_touched_opposite_half_rate_max"`
	LabelTouchedOppositeHalfRateDelta     float64 `json:"label_touched_opposite_half_rate_delta"`
	LabelClosedBackAcrossMidRateMin       float64 `json:"label_closed_back_across_mid_rate_min"`
	LabelClosedBackAcrossMidRateMax       float64 `json:"label_closed_back_across_mid_rate_max"`
	LabelClosedBackAcrossMidRateDelta     float64 `json:"label_closed_back_across_mid_rate_delta"`
	LabelMidRejectionBeforeBoundRateMin   float64 `json:"label_mid_rejection_before_boundary_touch_rate_min"`
	LabelMidRejectionBeforeBoundRateMax   float64 `json:"label_mid_rejection_before_boundary_touch_rate_max"`
	LabelMidRejectionBeforeBoundRateDelta float64 `json:"label_mid_rejection_before_boundary_touch_rate_delta"`
	LabelBoundBeforeMidRejectionRateMin   float64 `json:"label_boundary_touch_before_mid_rejection_rate_min"`
	LabelBoundBeforeMidRejectionRateMax   float64 `json:"label_boundary_touch_before_mid_rejection_rate_max"`
	LabelBoundBeforeMidRejectionRateDelta float64 `json:"label_boundary_touch_before_mid_rejection_rate_delta"`
}

func DefaultHoldInsideMidlineReactionAuditConfig() HoldInsideMidlineReactionAuditConfig {
	lookbackDays := DefaultCompressionRangeDetectorConfig().LookbackDays
	return HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1, 3, 6, 12},
		QuickInvalidationBars:    3,
		MaxMidlineEventDelayBars: 12,
		Profiles:                 defaultHoldInsideDirectionalEdgeProfiles(lookbackDays),
		ContextRules:             defaultHoldInsideDirectionalEdgeRules(),
	}
}

func RunHoldInsideMidlineReactionAudit(candles []Candle, detectorCfg RangeDetectorConfig, cfg HoldInsideMidlineReactionAuditConfig, splits []Split) ([]HoldInsideMidlineReactionCandidateRow, []HoldInsideMidlineReactionFunnelSummaryRow, []HoldInsideMidlineReactionSummaryRow, []HoldInsideMidlineReactionStabilityRow, error) {
	detectorCfg = detectorCfg.withDefaults()
	if err := detectorCfg.validate(); err != nil {
		return nil, nil, nil, nil, err
	}
	cfg = cfg.withDefaults(detectorCfg.LookbackDays)
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, nil, err
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

	candidates := []HoldInsideMidlineReactionCandidateRow{}
	funnelAccumulators := map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator{}
	summaryAccumulators := map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator{}
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
		profileCandidates := runHoldInsideMidlineReactionForProfile(candles, classifications, normalizedATR, profile, cfg, splits, funnelAccumulators, summaryAccumulators)
		candidates = append(candidates, profileCandidates...)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideMidlineReactionCandidate(candidates[i], candidates[j])
	})
	funnelRows := holdInsideMidlineReactionFunnelRows(funnelAccumulators)
	summaryRows := holdInsideMidlineReactionSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideMidlineReactionStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, funnelRows, summaryRows, stabilityRows, nil
}

func runHoldInsideMidlineReactionAuditFromClassifications(candles []Candle, profile DetectorSweepProfile, classifications []RangeClassification, cfg HoldInsideMidlineReactionAuditConfig, splits []Split) ([]HoldInsideMidlineReactionCandidateRow, []HoldInsideMidlineReactionFunnelSummaryRow, []HoldInsideMidlineReactionSummaryRow, []HoldInsideMidlineReactionStabilityRow, error) {
	cfg = cfg.withDefaults(profile.LookbackDays)
	cfg.Profiles = []DetectorSweepProfile{profile}
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}

	normalizedATR := NormalizedATR(candles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
	funnelAccumulators := map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator{}
	summaryAccumulators := map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator{}
	candidates := runHoldInsideMidlineReactionForProfile(candles, classifications, normalizedATR, profile, cfg, splits, funnelAccumulators, summaryAccumulators)
	sort.Slice(candidates, func(i, j int) bool {
		return lessHoldInsideMidlineReactionCandidate(candidates[i], candidates[j])
	})
	funnelRows := holdInsideMidlineReactionFunnelRows(funnelAccumulators)
	summaryRows := holdInsideMidlineReactionSummaryRows(summaryAccumulators)
	stabilityRows := holdInsideMidlineReactionStabilityRows(cfg.Profiles, cfg.ContextRules, cfg.HorizonsBars, summaryRows, splits)
	return candidates, funnelRows, summaryRows, stabilityRows, nil
}

func (cfg HoldInsideMidlineReactionAuditConfig) withDefaults(lookbackDays int) HoldInsideMidlineReactionAuditConfig {
	defaults := DefaultHoldInsideMidlineReactionAuditConfig()
	if lookbackDays > 0 {
		defaults.Profiles = defaultHoldInsideDirectionalEdgeProfiles(lookbackDays)
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickInvalidationBars == 0 {
		cfg.QuickInvalidationBars = defaults.QuickInvalidationBars
	}
	if cfg.MaxMidlineEventDelayBars == 0 {
		cfg.MaxMidlineEventDelayBars = defaults.MaxMidlineEventDelayBars
	}
	if len(cfg.Profiles) == 0 {
		cfg.Profiles = append([]DetectorSweepProfile(nil), defaults.Profiles...)
	}
	if len(cfg.ContextRules) == 0 {
		cfg.ContextRules = append([]DetectorContextRefinementRule(nil), defaults.ContextRules...)
	}
	return cfg
}

func (cfg HoldInsideMidlineReactionAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("hold-inside midline reaction audit horizon bars must be positive")
		}
	}
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("hold-inside midline reaction audit quick invalidation bars must be positive")
	}
	if cfg.MaxMidlineEventDelayBars <= 0 {
		return fmt.Errorf("hold-inside midline reaction audit max midline event delay bars must be positive")
	}
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("hold-inside midline reaction audit profiles must not be empty")
	}
	for _, profile := range cfg.Profiles {
		if profile.ProfileID == "" {
			return fmt.Errorf("hold-inside midline reaction audit profile id must not be empty")
		}
	}
	if len(cfg.ContextRules) == 0 {
		return fmt.Errorf("hold-inside midline reaction audit rules must not be empty")
	}
	for _, rule := range cfg.ContextRules {
		if rule.RuleID == "" {
			return fmt.Errorf("hold-inside midline reaction audit rule id must not be empty")
		}
		if rule.HoldBars <= 0 {
			return fmt.Errorf("hold-inside midline reaction audit rules must use positive hold bars")
		}
	}
	return nil
}

func runHoldInsideMidlineReactionForProfile(candles []Candle, classifications []RangeClassification, normalizedATR []float64, profile DetectorSweepProfile, cfg HoldInsideMidlineReactionAuditConfig, splits []Split, funnelAccumulators map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator, summaryAccumulators map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator) []HoldInsideMidlineReactionCandidateRow {
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, profile.ProfileID)
	candidates := []HoldInsideMidlineReactionCandidateRow{}
	maxHorizon := maxIntInSlice(cfg.HorizonsBars)
	for _, episode := range episodes {
		for _, rule := range cfg.ContextRules {
			holdDecisionIndex := episode.EndIndex + rule.HoldBars
			if holdDecisionIndex < 0 || holdDecisionIndex >= len(candles) {
				continue
			}
			if !detectorContextRulePasses(candles, episode, rule, holdDecisionIndex) {
				continue
			}
			holdSplit := splitNameForCloseTime(candles[holdDecisionIndex].CloseTime, splits)
			episodeMid := (episode.High + episode.Low) / 2
			for _, eventType := range holdInsideMidlineReactionEventTypes() {
				eventIndex, foundEvent := firstHoldInsideMidlineReactionEventIndex(candles, holdDecisionIndex, episodeMid, eventType, cfg.MaxMidlineEventDelayBars)
				missingFuture := foundEvent && eventIndex+maxHorizon >= len(candles)
				addHoldInsideMidlineReactionFunnel(funnelAccumulators, profile, rule, eventType, holdSplit, foundEvent, missingFuture, eventIndex-holdDecisionIndex)
				if holdSplit != fullSplitName {
					addHoldInsideMidlineReactionFunnel(funnelAccumulators, profile, rule, eventType, fullSplitName, foundEvent, missingFuture, eventIndex-holdDecisionIndex)
				}
				if !foundEvent {
					continue
				}
				eventSplit := splitNameForCloseTime(candles[eventIndex].CloseTime, splits)
				for _, horizon := range cfg.HorizonsBars {
					if eventIndex+horizon >= len(candles) {
						continue
					}
					row, ok := newHoldInsideMidlineReactionCandidateRow(candles, profile, rule, eventType, episode, eventSplit, holdDecisionIndex, eventIndex, horizon, cfg.QuickInvalidationBars)
					if !ok {
						continue
					}
					candidates = append(candidates, row)
					for _, combo := range holdInsideMidlineReactionSummaryCombos(row.EventMidSide, row.EventClosePositionBucket) {
						addHoldInsideMidlineReactionCandidate(summaryAccumulators, row, eventSplit, combo.midSide, combo.bucket)
						if eventSplit != fullSplitName {
							addHoldInsideMidlineReactionCandidate(summaryAccumulators, row, fullSplitName, combo.midSide, combo.bucket)
						}
					}
				}
			}
		}
	}
	return candidates
}

func newHoldInsideMidlineReactionCandidateRow(candles []Candle, profile DetectorSweepProfile, rule DetectorContextRefinementRule, eventType string, episode rangeRegimeDurabilityEpisode, split string, holdDecisionIndex, eventIndex, horizon, quickInvalidationBars int) (HoldInsideMidlineReactionCandidateRow, bool) {
	if horizon <= 0 || quickInvalidationBars <= 0 || holdDecisionIndex < 0 || eventIndex <= holdDecisionIndex || eventIndex+horizon >= len(candles) || !validHoldInsideMidlineReactionEventType(eventType) {
		return HoldInsideMidlineReactionCandidateRow{}, false
	}
	holdDecision := candles[holdDecisionIndex]
	event := candles[eventIndex]
	episodeMid := (episode.High + episode.Low) / 2
	holdPosition := decisionClosePosition(holdDecision.Close, episode.Low, episode.High)
	holdPositionValue := holdPosition
	if !validNumber(holdPositionValue) {
		holdPositionValue = 0
	}
	eventPosition := decisionClosePosition(event.Close, episode.Low, episode.High)
	eventPositionValue := eventPosition
	if !validNumber(eventPositionValue) {
		eventPositionValue = 0
	}
	startIndex := eventIndex + 1
	endIndex := eventIndex + horizon
	future := candles[startIndex : endIndex+1]
	label := newHoldInsideMidlineReactionLabel(future, event.Close, episode.Low, episode.High, episodeMid, quickInvalidationBars)

	return HoldInsideMidlineReactionCandidateRow{
		ProfileID:                    profile.ProfileID,
		IsBalancedBaseline:           profile.IsBalancedBaseline,
		IsADXComparison:              profile.IsADXComparison,
		Percentile:                   profile.Percentile,
		MinConsecutiveBars:           profile.MinConsecutiveBars,
		UseBollinger:                 profile.UseBollinger,
		UseADX:                       profile.UseADX,
		LookbackDays:                 profile.LookbackDays,
		ContextRule:                  rule.RuleID,
		HoldBars:                     rule.HoldBars,
		RequireMid50:                 rule.RequireMid50,
		EventType:                    eventType,
		Split:                        split,
		SourceEpisodeID:              episode.EpisodeID,
		EpisodeStartIndex:            episode.StartIndex,
		EpisodeEndIndex:              episode.EndIndex,
		EpisodeStartTime:             candles[episode.StartIndex].CloseTime.Format(timeLayout),
		EpisodeEndTime:               candles[episode.EndIndex].CloseTime.Format(timeLayout),
		RawLengthBars:                episode.RawLengthBars,
		ActiveLengthBars:             episode.ActiveLengthBars,
		RawLengthBucket:              episode.RawLengthBucket,
		ActiveLengthBucket:           episode.ActiveLengthBucket,
		EpisodeHigh:                  episode.High,
		EpisodeLow:                   episode.Low,
		EpisodeMid:                   episodeMid,
		EpisodeEndClose:              episode.EndClose,
		EpisodeWidthPct:              episode.WidthPct,
		EpisodeWidthBucket:           episode.WidthBucket,
		AvgNormalizedATR:             episode.AvgNormalizedATR,
		EndNormalizedATR:             episode.EndNormalizedATR,
		WidthToATRRatio:              episode.WidthToATRRatio,
		WidthToATRBucket:             episode.WidthToATRBucket,
		HoldDecisionIndex:            holdDecisionIndex,
		HoldDecisionTime:             holdDecision.CloseTime.Format(timeLayout),
		HoldDecisionClose:            holdDecision.Close,
		HoldDecisionClosePosition:    holdPositionValue,
		HoldDecisionCloseBucket:      decisionClosePositionBucket(holdPosition),
		HoldDecisionMidSide:          holdInsideDecisionMidSide(holdDecision.Close, episodeMid),
		EventIndex:                   eventIndex,
		EventTime:                    event.CloseTime.Format(timeLayout),
		EventDelayBars:               eventIndex - holdDecisionIndex,
		EventClose:                   event.Close,
		EventClosePosition:           eventPositionValue,
		EventClosePositionBucket:     decisionClosePositionBucket(eventPosition),
		EventMidSide:                 holdInsideDecisionMidSide(event.Close, episodeMid),
		EventDistanceToHighPct:       movePct(math.Max(0, episode.High-event.Close), event.Close),
		EventDistanceToLowPct:        movePct(math.Max(0, event.Close-episode.Low), event.Close),
		EventDistanceToMidPct:        movePct(math.Abs(event.Close-episodeMid), event.Close),
		HorizonBars:                  horizon,
		LabelWindowStartIndex:        startIndex,
		LabelWindowEndIndex:          endIndex,
		LabelWindowStartTime:         candles[startIndex].CloseTime.Format(timeLayout),
		LabelWindowEndTime:           candles[endIndex].CloseTime.Format(timeLayout),
		LabelReenteredRange:          label.ReenteredRange,
		LabelPersistedInsideRange:    label.PersistedInsideRange,
		LabelQuickInvalidated:        label.QuickInvalidated,
		LabelInvalidatedUp:           label.InvalidatedUp,
		LabelInvalidatedDown:         label.InvalidatedDown,
		LabelTrendedUp:               label.TrendedUp,
		LabelTrendedDown:             label.TrendedDown,
		LabelTouchedHigh:             label.TouchedHigh,
		LabelTouchedLow:              label.TouchedLow,
		LabelTouchedOppositeHalf:     label.TouchedOppositeHalf,
		LabelClosedBackAcrossMid:     label.ClosedBackAcrossMid,
		LabelMidRejectionBeforeBound: label.MidRejectionBeforeBoundaryTouch,
		LabelBoundBeforeMidRejection: label.BoundaryTouchBeforeMidRejection,
	}, true
}

type holdInsideMidlineReactionLabel struct {
	ReenteredRange                  bool
	PersistedInsideRange            bool
	QuickInvalidated                bool
	InvalidatedUp                   bool
	InvalidatedDown                 bool
	TrendedUp                       bool
	TrendedDown                     bool
	TouchedHigh                     bool
	TouchedLow                      bool
	TouchedOppositeHalf             bool
	ClosedBackAcrossMid             bool
	MidRejectionBeforeBoundaryTouch bool
	BoundaryTouchBeforeMidRejection bool
}

func newHoldInsideMidlineReactionLabel(future []Candle, eventClose, low, high, mid float64, quickInvalidationBars int) holdInsideMidlineReactionLabel {
	futureMaxHigh, futureMinLow := futureHighLow(future)
	lastClose := future[len(future)-1].Close
	maxUpMovePct := movePct(math.Max(0, futureMaxHigh-eventClose), eventClose)
	maxDownMovePct := movePct(math.Max(0, eventClose-futureMinLow), eventClose)
	label := holdInsideMidlineReactionLabel{PersistedInsideRange: true}
	firstBoundaryTouchDelay := -1
	firstMidRejectionDelay := -1
	quickLimit := minInt(len(future), quickInvalidationBars)
	for i, candle := range future {
		delay := i + 1
		if candle.High >= high {
			label.TouchedHigh = true
			if firstBoundaryTouchDelay == -1 {
				firstBoundaryTouchDelay = delay
			}
		}
		if candle.Low <= low {
			label.TouchedLow = true
			if firstBoundaryTouchDelay == -1 {
				firstBoundaryTouchDelay = delay
			}
		}
		touchedOppositeHalf, closedBackAcrossMid := holdInsideMidLabels(candle, eventClose, mid)
		if touchedOppositeHalf {
			label.TouchedOppositeHalf = true
		}
		if closedBackAcrossMid {
			label.ClosedBackAcrossMid = true
			if firstMidRejectionDelay == -1 {
				firstMidRejectionDelay = delay
			}
		}
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
	}
	label.MidRejectionBeforeBoundaryTouch = firstMidRejectionDelay != -1 &&
		(firstBoundaryTouchDelay == -1 || firstMidRejectionDelay < firstBoundaryTouchDelay)
	label.BoundaryTouchBeforeMidRejection = firstBoundaryTouchDelay != -1 &&
		(firstMidRejectionDelay == -1 || firstBoundaryTouchDelay < firstMidRejectionDelay)
	label.TrendedUp = label.InvalidatedUp && lastClose > high && maxUpMovePct > maxDownMovePct
	label.TrendedDown = label.InvalidatedDown && lastClose < low && maxDownMovePct > maxUpMovePct
	return label
}

func firstHoldInsideMidlineReactionEventIndex(candles []Candle, holdDecisionIndex int, mid float64, eventType string, maxDelay int) (int, bool) {
	if holdDecisionIndex < 0 || holdDecisionIndex >= len(candles) || maxDelay <= 0 || !validHoldInsideMidlineReactionEventType(eventType) {
		return -1, false
	}
	decisionClose := candles[holdDecisionIndex].Close
	lastIndex := minInt(len(candles)-1, holdDecisionIndex+maxDelay)
	for i := holdDecisionIndex + 1; i <= lastIndex; i++ {
		touchedMid, closedAcrossMid := holdInsideMidLabels(candles[i], decisionClose, mid)
		if eventType == HoldInsideMidlineReactionEventTouch && touchedMid {
			return i, true
		}
		if eventType == HoldInsideMidlineReactionEventCloseAcross && closedAcrossMid {
			return i, true
		}
	}
	return -1, false
}

func holdInsideMidlineReactionEventTypes() []string {
	return []string{HoldInsideMidlineReactionEventTouch, HoldInsideMidlineReactionEventCloseAcross}
}

func validHoldInsideMidlineReactionEventType(eventType string) bool {
	return eventType == HoldInsideMidlineReactionEventTouch || eventType == HoldInsideMidlineReactionEventCloseAcross
}

type holdInsideMidlineReactionFunnelKey struct {
	profileID   string
	contextRule string
	eventType   string
	split       string
}

type holdInsideMidlineReactionFunnelAccumulator struct {
	profile               DetectorSweepProfile
	rule                  DetectorContextRefinementRule
	eventType             string
	split                 string
	sourceHoldCount       int
	eventCount            int
	missingEventCount     int
	missingFutureCount    int
	reactionEligibleCount int
	eventDelaySum         float64
}

func addHoldInsideMidlineReactionFunnel(accumulators map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, eventType, split string, foundEvent, missingFuture bool, eventDelayBars int) {
	acc := holdInsideMidlineReactionFunnelAccumulatorFor(accumulators, profile, rule, eventType, split)
	acc.sourceHoldCount++
	if !foundEvent {
		acc.missingEventCount++
		return
	}
	acc.eventCount++
	acc.eventDelaySum += float64(eventDelayBars)
	if missingFuture {
		acc.missingFutureCount++
	} else {
		acc.reactionEligibleCount++
	}
}

func holdInsideMidlineReactionFunnelAccumulatorFor(accumulators map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, eventType, split string) *holdInsideMidlineReactionFunnelAccumulator {
	key := holdInsideMidlineReactionFunnelKey{
		profileID:   profile.ProfileID,
		contextRule: rule.RuleID,
		eventType:   eventType,
		split:       split,
	}
	acc := accumulators[key]
	if acc == nil {
		acc = &holdInsideMidlineReactionFunnelAccumulator{
			profile:   profile,
			rule:      rule,
			eventType: eventType,
			split:     split,
		}
		accumulators[key] = acc
	}
	return acc
}

func holdInsideMidlineReactionFunnelRows(accumulators map[holdInsideMidlineReactionFunnelKey]*holdInsideMidlineReactionFunnelAccumulator) []HoldInsideMidlineReactionFunnelSummaryRow {
	rows := make([]HoldInsideMidlineReactionFunnelSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideMidlineReactionFunnel(rows[i], rows[j])
	})
	return rows
}

func (acc holdInsideMidlineReactionFunnelAccumulator) row() HoldInsideMidlineReactionFunnelSummaryRow {
	row := HoldInsideMidlineReactionFunnelSummaryRow{
		ProfileID:             acc.profile.ProfileID,
		IsBalancedBaseline:    acc.profile.IsBalancedBaseline,
		IsADXComparison:       acc.profile.IsADXComparison,
		Percentile:            acc.profile.Percentile,
		MinConsecutiveBars:    acc.profile.MinConsecutiveBars,
		UseBollinger:          acc.profile.UseBollinger,
		UseADX:                acc.profile.UseADX,
		LookbackDays:          acc.profile.LookbackDays,
		ContextRule:           acc.rule.RuleID,
		HoldBars:              acc.rule.HoldBars,
		RequireMid50:          acc.rule.RequireMid50,
		EventType:             acc.eventType,
		Split:                 acc.split,
		SourceHoldCount:       acc.sourceHoldCount,
		EventCount:            acc.eventCount,
		MissingEventCount:     acc.missingEventCount,
		MissingFutureCount:    acc.missingFutureCount,
		ReactionEligibleCount: acc.reactionEligibleCount,
	}
	if acc.sourceHoldCount > 0 {
		row.EventRate = float64(acc.eventCount) / float64(acc.sourceHoldCount)
		row.ReactionEligibleRate = float64(acc.reactionEligibleCount) / float64(acc.sourceHoldCount)
	}
	if acc.eventCount > 0 {
		row.AvgEventDelayBars = acc.eventDelaySum / float64(acc.eventCount)
	}
	return row
}

type holdInsideMidlineReactionSummaryCombo struct {
	midSide string
	bucket  string
}

func holdInsideMidlineReactionSummaryCombos(midSide, bucket string) []holdInsideMidlineReactionSummaryCombo {
	return []holdInsideMidlineReactionSummaryCombo{
		{midSide: midSide, bucket: bucket},
		{midSide: midSide, bucket: holdInsideDecisionPositionBucketAll},
		{midSide: holdInsideDecisionMidSideAll, bucket: bucket},
		{midSide: holdInsideDecisionMidSideAll, bucket: holdInsideDecisionPositionBucketAll},
	}
}

type holdInsideMidlineReactionSummaryKey struct {
	profileID                string
	contextRule              string
	eventType                string
	split                    string
	horizonBars              int
	eventMidSide             string
	eventClosePositionBucket string
}

type holdInsideMidlineReactionSummaryAccumulator struct {
	profile                    DetectorSweepProfile
	rule                       DetectorContextRefinementRule
	eventType                  string
	split                      string
	horizonBars                int
	eventMidSide               string
	eventClosePositionBucket   string
	candidates                 int
	eventDelaySum              float64
	eventPositionSum           float64
	eventDistanceToHighPctSum  float64
	eventDistanceToLowPctSum   float64
	eventDistanceToMidPctSum   float64
	reenteredRange             int
	persistedInsideRange       int
	quickInvalidated           int
	invalidatedUp              int
	invalidatedDown            int
	trendedUp                  int
	trendedDown                int
	touchedHigh                int
	touchedLow                 int
	touchedOppositeHalf        int
	closedBackAcrossMid        int
	midRejectionBeforeBoundary int
	boundaryBeforeMidRejection int
}

func addHoldInsideMidlineReactionCandidate(accumulators map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator, row HoldInsideMidlineReactionCandidateRow, split, midSide, bucket string) {
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
	acc := holdInsideMidlineReactionSummaryAccumulatorFor(accumulators, profile, rule, row.EventType, split, row.HorizonBars, midSide, bucket)
	acc.candidates++
	acc.eventDelaySum += float64(row.EventDelayBars)
	if row.EventClosePositionBucket != decisionClosePositionBucketUnknown {
		acc.eventPositionSum += row.EventClosePosition
	}
	acc.eventDistanceToHighPctSum += row.EventDistanceToHighPct
	acc.eventDistanceToLowPctSum += row.EventDistanceToLowPct
	acc.eventDistanceToMidPctSum += row.EventDistanceToMidPct
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
	if row.LabelTouchedHigh {
		acc.touchedHigh++
	}
	if row.LabelTouchedLow {
		acc.touchedLow++
	}
	if row.LabelTouchedOppositeHalf {
		acc.touchedOppositeHalf++
	}
	if row.LabelClosedBackAcrossMid {
		acc.closedBackAcrossMid++
	}
	if row.LabelMidRejectionBeforeBound {
		acc.midRejectionBeforeBoundary++
	}
	if row.LabelBoundBeforeMidRejection {
		acc.boundaryBeforeMidRejection++
	}
}

func holdInsideMidlineReactionSummaryAccumulatorFor(accumulators map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator, profile DetectorSweepProfile, rule DetectorContextRefinementRule, eventType, split string, horizon int, midSide, bucket string) *holdInsideMidlineReactionSummaryAccumulator {
	key := holdInsideMidlineReactionSummaryKey{
		profileID:                profile.ProfileID,
		contextRule:              rule.RuleID,
		eventType:                eventType,
		split:                    split,
		horizonBars:              horizon,
		eventMidSide:             midSide,
		eventClosePositionBucket: bucket,
	}
	acc := accumulators[key]
	if acc == nil {
		acc = &holdInsideMidlineReactionSummaryAccumulator{
			profile:                  profile,
			rule:                     rule,
			eventType:                eventType,
			split:                    split,
			horizonBars:              horizon,
			eventMidSide:             midSide,
			eventClosePositionBucket: bucket,
		}
		accumulators[key] = acc
	}
	return acc
}

func holdInsideMidlineReactionSummaryRows(accumulators map[holdInsideMidlineReactionSummaryKey]*holdInsideMidlineReactionSummaryAccumulator) []HoldInsideMidlineReactionSummaryRow {
	rows := make([]HoldInsideMidlineReactionSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideMidlineReactionSummary(rows[i], rows[j])
	})
	return rows
}

func (acc holdInsideMidlineReactionSummaryAccumulator) row() HoldInsideMidlineReactionSummaryRow {
	row := HoldInsideMidlineReactionSummaryRow{
		ProfileID:                         acc.profile.ProfileID,
		IsBalancedBaseline:                acc.profile.IsBalancedBaseline,
		IsADXComparison:                   acc.profile.IsADXComparison,
		Percentile:                        acc.profile.Percentile,
		MinConsecutiveBars:                acc.profile.MinConsecutiveBars,
		UseBollinger:                      acc.profile.UseBollinger,
		UseADX:                            acc.profile.UseADX,
		LookbackDays:                      acc.profile.LookbackDays,
		ContextRule:                       acc.rule.RuleID,
		HoldBars:                          acc.rule.HoldBars,
		RequireMid50:                      acc.rule.RequireMid50,
		EventType:                         acc.eventType,
		Split:                             acc.split,
		HorizonBars:                       acc.horizonBars,
		EventMidSide:                      acc.eventMidSide,
		EventClosePositionBucket:          acc.eventClosePositionBucket,
		CandidateCount:                    acc.candidates,
		LabelReenteredRangeCount:          acc.reenteredRange,
		LabelPersistedInsideRangeCount:    acc.persistedInsideRange,
		LabelQuickInvalidatedCount:        acc.quickInvalidated,
		LabelInvalidatedUpCount:           acc.invalidatedUp,
		LabelInvalidatedDownCount:         acc.invalidatedDown,
		LabelTrendedUpCount:               acc.trendedUp,
		LabelTrendedDownCount:             acc.trendedDown,
		LabelTouchedHighCount:             acc.touchedHigh,
		LabelTouchedLowCount:              acc.touchedLow,
		LabelTouchedOppositeHalfCount:     acc.touchedOppositeHalf,
		LabelClosedBackAcrossMidCount:     acc.closedBackAcrossMid,
		LabelMidRejectionBeforeBoundCount: acc.midRejectionBeforeBoundary,
		LabelBoundBeforeMidRejectionCount: acc.boundaryBeforeMidRejection,
	}
	if acc.candidates == 0 {
		return row
	}
	count := float64(acc.candidates)
	row.AvgEventDelayBars = acc.eventDelaySum / count
	row.AvgEventClosePosition = acc.eventPositionSum / count
	row.AvgEventDistanceToHighPct = acc.eventDistanceToHighPctSum / count
	row.AvgEventDistanceToLowPct = acc.eventDistanceToLowPctSum / count
	row.AvgEventDistanceToMidPct = acc.eventDistanceToMidPctSum / count
	row.LabelReenteredRangeRate = float64(acc.reenteredRange) / count
	row.LabelPersistedInsideRangeRate = float64(acc.persistedInsideRange) / count
	row.LabelQuickInvalidatedRate = float64(acc.quickInvalidated) / count
	row.LabelInvalidatedUpRate = float64(acc.invalidatedUp) / count
	row.LabelInvalidatedDownRate = float64(acc.invalidatedDown) / count
	row.LabelTrendedUpRate = float64(acc.trendedUp) / count
	row.LabelTrendedDownRate = float64(acc.trendedDown) / count
	row.LabelTouchedHighRate = float64(acc.touchedHigh) / count
	row.LabelTouchedLowRate = float64(acc.touchedLow) / count
	row.LabelTouchedOppositeHalfRate = float64(acc.touchedOppositeHalf) / count
	row.LabelClosedBackAcrossMidRate = float64(acc.closedBackAcrossMid) / count
	row.LabelMidRejectionBeforeBoundRate = float64(acc.midRejectionBeforeBoundary) / count
	row.LabelBoundBeforeMidRejectionRate = float64(acc.boundaryBeforeMidRejection) / count
	return row
}

func holdInsideMidlineReactionStabilityRows(profiles []DetectorSweepProfile, rules []DetectorContextRefinementRule, horizons []int, summaryRows []HoldInsideMidlineReactionSummaryRow, splits []Split) []HoldInsideMidlineReactionStabilityRow {
	periodSplits := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplits = append(periodSplits, split)
		}
	}
	byKey := map[holdInsideMidlineReactionSummaryKey]HoldInsideMidlineReactionSummaryRow{}
	combosByBase := map[holdInsideMidlineReactionStabilityBaseKey]map[holdInsideMidlineReactionSummaryCombo]bool{}
	for _, row := range summaryRows {
		key := holdInsideMidlineReactionSummaryKey{
			profileID:                row.ProfileID,
			contextRule:              row.ContextRule,
			eventType:                row.EventType,
			split:                    row.Split,
			horizonBars:              row.HorizonBars,
			eventMidSide:             row.EventMidSide,
			eventClosePositionBucket: row.EventClosePositionBucket,
		}
		byKey[key] = row
		base := holdInsideMidlineReactionStabilityBaseKey{
			profileID:   row.ProfileID,
			contextRule: row.ContextRule,
			eventType:   row.EventType,
			horizonBars: row.HorizonBars,
		}
		if combosByBase[base] == nil {
			combosByBase[base] = map[holdInsideMidlineReactionSummaryCombo]bool{}
		}
		combosByBase[base][holdInsideMidlineReactionSummaryCombo{midSide: row.EventMidSide, bucket: row.EventClosePositionBucket}] = true
	}

	rows := []HoldInsideMidlineReactionStabilityRow{}
	for _, profile := range profiles {
		for _, rule := range rules {
			for _, eventType := range holdInsideMidlineReactionEventTypes() {
				for _, horizon := range horizons {
					base := holdInsideMidlineReactionStabilityBaseKey{
						profileID:   profile.ProfileID,
						contextRule: rule.RuleID,
						eventType:   eventType,
						horizonBars: horizon,
					}
					for _, combo := range holdInsideMidlineReactionSortedCombos(combosByBase[base]) {
						rows = append(rows, newHoldInsideMidlineReactionStabilityRow(profile, rule, eventType, horizon, combo.midSide, combo.bucket, periodSplits, byKey))
					}
				}
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessHoldInsideMidlineReactionStability(rows[i], rows[j])
	})
	return rows
}

type holdInsideMidlineReactionStabilityBaseKey struct {
	profileID   string
	contextRule string
	eventType   string
	horizonBars int
}

func newHoldInsideMidlineReactionStabilityRow(profile DetectorSweepProfile, rule DetectorContextRefinementRule, eventType string, horizon int, midSide, bucket string, periodSplits []Split, byKey map[holdInsideMidlineReactionSummaryKey]HoldInsideMidlineReactionSummaryRow) HoldInsideMidlineReactionStabilityRow {
	row := HoldInsideMidlineReactionStabilityRow{
		ProfileID:                profile.ProfileID,
		IsBalancedBaseline:       profile.IsBalancedBaseline,
		IsADXComparison:          profile.IsADXComparison,
		Percentile:               profile.Percentile,
		MinConsecutiveBars:       profile.MinConsecutiveBars,
		UseBollinger:             profile.UseBollinger,
		UseADX:                   profile.UseADX,
		LookbackDays:             profile.LookbackDays,
		ContextRule:              rule.RuleID,
		HoldBars:                 rule.HoldBars,
		RequireMid50:             rule.RequireMid50,
		EventType:                eventType,
		HorizonBars:              horizon,
		EventMidSide:             midSide,
		EventClosePositionBucket: bucket,
		PeriodSplits:             len(periodSplits),
	}
	first := true
	for _, split := range periodSplits {
		summary := byKey[holdInsideMidlineReactionSummaryKey{
			profileID:                profile.ProfileID,
			contextRule:              rule.RuleID,
			eventType:                eventType,
			split:                    split.Name,
			horizonBars:              horizon,
			eventMidSide:             midSide,
			eventClosePositionBucket: bucket,
		}]
		trendedRate := summary.LabelTrendedUpRate + summary.LabelTrendedDownRate
		if first {
			row.CandidateCountMin = summary.CandidateCount
			row.CandidateCountMax = summary.CandidateCount
			row.AvgEventDelayBarsMin = summary.AvgEventDelayBars
			row.AvgEventDelayBarsMax = summary.AvgEventDelayBars
			row.LabelReenteredRangeRateMin = summary.LabelReenteredRangeRate
			row.LabelReenteredRangeRateMax = summary.LabelReenteredRangeRate
			row.LabelPersistedInsideRangeRateMin = summary.LabelPersistedInsideRangeRate
			row.LabelPersistedInsideRangeRateMax = summary.LabelPersistedInsideRangeRate
			row.LabelQuickInvalidatedRateMin = summary.LabelQuickInvalidatedRate
			row.LabelQuickInvalidatedRateMax = summary.LabelQuickInvalidatedRate
			row.LabelTrendedRateMin = trendedRate
			row.LabelTrendedRateMax = trendedRate
			row.LabelTouchedHighRateMin = summary.LabelTouchedHighRate
			row.LabelTouchedHighRateMax = summary.LabelTouchedHighRate
			row.LabelTouchedLowRateMin = summary.LabelTouchedLowRate
			row.LabelTouchedLowRateMax = summary.LabelTouchedLowRate
			row.LabelTouchedOppositeHalfRateMin = summary.LabelTouchedOppositeHalfRate
			row.LabelTouchedOppositeHalfRateMax = summary.LabelTouchedOppositeHalfRate
			row.LabelClosedBackAcrossMidRateMin = summary.LabelClosedBackAcrossMidRate
			row.LabelClosedBackAcrossMidRateMax = summary.LabelClosedBackAcrossMidRate
			row.LabelMidRejectionBeforeBoundRateMin = summary.LabelMidRejectionBeforeBoundRate
			row.LabelMidRejectionBeforeBoundRateMax = summary.LabelMidRejectionBeforeBoundRate
			row.LabelBoundBeforeMidRejectionRateMin = summary.LabelBoundBeforeMidRejectionRate
			row.LabelBoundBeforeMidRejectionRateMax = summary.LabelBoundBeforeMidRejectionRate
			first = false
		} else {
			row.CandidateCountMin = minInt(row.CandidateCountMin, summary.CandidateCount)
			row.CandidateCountMax = maxInt(row.CandidateCountMax, summary.CandidateCount)
			row.AvgEventDelayBarsMin = minFloat(row.AvgEventDelayBarsMin, summary.AvgEventDelayBars)
			row.AvgEventDelayBarsMax = maxFloat(row.AvgEventDelayBarsMax, summary.AvgEventDelayBars)
			row.LabelReenteredRangeRateMin = minFloat(row.LabelReenteredRangeRateMin, summary.LabelReenteredRangeRate)
			row.LabelReenteredRangeRateMax = maxFloat(row.LabelReenteredRangeRateMax, summary.LabelReenteredRangeRate)
			row.LabelPersistedInsideRangeRateMin = minFloat(row.LabelPersistedInsideRangeRateMin, summary.LabelPersistedInsideRangeRate)
			row.LabelPersistedInsideRangeRateMax = maxFloat(row.LabelPersistedInsideRangeRateMax, summary.LabelPersistedInsideRangeRate)
			row.LabelQuickInvalidatedRateMin = minFloat(row.LabelQuickInvalidatedRateMin, summary.LabelQuickInvalidatedRate)
			row.LabelQuickInvalidatedRateMax = maxFloat(row.LabelQuickInvalidatedRateMax, summary.LabelQuickInvalidatedRate)
			row.LabelTrendedRateMin = minFloat(row.LabelTrendedRateMin, trendedRate)
			row.LabelTrendedRateMax = maxFloat(row.LabelTrendedRateMax, trendedRate)
			row.LabelTouchedHighRateMin = minFloat(row.LabelTouchedHighRateMin, summary.LabelTouchedHighRate)
			row.LabelTouchedHighRateMax = maxFloat(row.LabelTouchedHighRateMax, summary.LabelTouchedHighRate)
			row.LabelTouchedLowRateMin = minFloat(row.LabelTouchedLowRateMin, summary.LabelTouchedLowRate)
			row.LabelTouchedLowRateMax = maxFloat(row.LabelTouchedLowRateMax, summary.LabelTouchedLowRate)
			row.LabelTouchedOppositeHalfRateMin = minFloat(row.LabelTouchedOppositeHalfRateMin, summary.LabelTouchedOppositeHalfRate)
			row.LabelTouchedOppositeHalfRateMax = maxFloat(row.LabelTouchedOppositeHalfRateMax, summary.LabelTouchedOppositeHalfRate)
			row.LabelClosedBackAcrossMidRateMin = minFloat(row.LabelClosedBackAcrossMidRateMin, summary.LabelClosedBackAcrossMidRate)
			row.LabelClosedBackAcrossMidRateMax = maxFloat(row.LabelClosedBackAcrossMidRateMax, summary.LabelClosedBackAcrossMidRate)
			row.LabelMidRejectionBeforeBoundRateMin = minFloat(row.LabelMidRejectionBeforeBoundRateMin, summary.LabelMidRejectionBeforeBoundRate)
			row.LabelMidRejectionBeforeBoundRateMax = maxFloat(row.LabelMidRejectionBeforeBoundRateMax, summary.LabelMidRejectionBeforeBoundRate)
			row.LabelBoundBeforeMidRejectionRateMin = minFloat(row.LabelBoundBeforeMidRejectionRateMin, summary.LabelBoundBeforeMidRejectionRate)
			row.LabelBoundBeforeMidRejectionRateMax = maxFloat(row.LabelBoundBeforeMidRejectionRateMax, summary.LabelBoundBeforeMidRejectionRate)
		}
		row.CandidateCount += summary.CandidateCount
	}
	row.CandidateCountDelta = row.CandidateCountMax - row.CandidateCountMin
	row.AvgEventDelayBarsDelta = row.AvgEventDelayBarsMax - row.AvgEventDelayBarsMin
	row.LabelReenteredRangeRateDelta = row.LabelReenteredRangeRateMax - row.LabelReenteredRangeRateMin
	row.LabelPersistedInsideRangeRateDelta = row.LabelPersistedInsideRangeRateMax - row.LabelPersistedInsideRangeRateMin
	row.LabelQuickInvalidatedRateDelta = row.LabelQuickInvalidatedRateMax - row.LabelQuickInvalidatedRateMin
	row.LabelTrendedRateDelta = row.LabelTrendedRateMax - row.LabelTrendedRateMin
	row.LabelTouchedHighRateDelta = row.LabelTouchedHighRateMax - row.LabelTouchedHighRateMin
	row.LabelTouchedLowRateDelta = row.LabelTouchedLowRateMax - row.LabelTouchedLowRateMin
	row.LabelTouchedOppositeHalfRateDelta = row.LabelTouchedOppositeHalfRateMax - row.LabelTouchedOppositeHalfRateMin
	row.LabelClosedBackAcrossMidRateDelta = row.LabelClosedBackAcrossMidRateMax - row.LabelClosedBackAcrossMidRateMin
	row.LabelMidRejectionBeforeBoundRateDelta = row.LabelMidRejectionBeforeBoundRateMax - row.LabelMidRejectionBeforeBoundRateMin
	row.LabelBoundBeforeMidRejectionRateDelta = row.LabelBoundBeforeMidRejectionRateMax - row.LabelBoundBeforeMidRejectionRateMin
	return row
}

func holdInsideMidlineReactionSortedCombos(comboSet map[holdInsideMidlineReactionSummaryCombo]bool) []holdInsideMidlineReactionSummaryCombo {
	if len(comboSet) == 0 {
		return nil
	}
	combos := make([]holdInsideMidlineReactionSummaryCombo, 0, len(comboSet))
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

func lessHoldInsideMidlineReactionCandidate(a, b HoldInsideMidlineReactionCandidateRow) bool {
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
	if holdInsideMidlineReactionEventTypeSortKey(a.EventType) != holdInsideMidlineReactionEventTypeSortKey(b.EventType) {
		return holdInsideMidlineReactionEventTypeSortKey(a.EventType) < holdInsideMidlineReactionEventTypeSortKey(b.EventType)
	}
	if a.EventIndex != b.EventIndex {
		return a.EventIndex < b.EventIndex
	}
	return a.HorizonBars < b.HorizonBars
}

func lessHoldInsideMidlineReactionFunnel(a, b HoldInsideMidlineReactionFunnelSummaryRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if holdInsideMidlineReactionEventTypeSortKey(a.EventType) != holdInsideMidlineReactionEventTypeSortKey(b.EventType) {
		return holdInsideMidlineReactionEventTypeSortKey(a.EventType) < holdInsideMidlineReactionEventTypeSortKey(b.EventType)
	}
	return splitSortKey(a.Split) < splitSortKey(b.Split)
}

func lessHoldInsideMidlineReactionSummary(a, b HoldInsideMidlineReactionSummaryRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if holdInsideMidlineReactionEventTypeSortKey(a.EventType) != holdInsideMidlineReactionEventTypeSortKey(b.EventType) {
		return holdInsideMidlineReactionEventTypeSortKey(a.EventType) < holdInsideMidlineReactionEventTypeSortKey(b.EventType)
	}
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if holdInsideDecisionMidSideSortKey(a.EventMidSide) != holdInsideDecisionMidSideSortKey(b.EventMidSide) {
		return holdInsideDecisionMidSideSortKey(a.EventMidSide) < holdInsideDecisionMidSideSortKey(b.EventMidSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.EventClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.EventClosePositionBucket)
}

func lessHoldInsideMidlineReactionStability(a, b HoldInsideMidlineReactionStabilityRow) bool {
	if a.ProfileID != b.ProfileID {
		return a.ProfileID < b.ProfileID
	}
	if detectorContextRuleSortKey(a.ContextRule) != detectorContextRuleSortKey(b.ContextRule) {
		return detectorContextRuleSortKey(a.ContextRule) < detectorContextRuleSortKey(b.ContextRule)
	}
	if holdInsideMidlineReactionEventTypeSortKey(a.EventType) != holdInsideMidlineReactionEventTypeSortKey(b.EventType) {
		return holdInsideMidlineReactionEventTypeSortKey(a.EventType) < holdInsideMidlineReactionEventTypeSortKey(b.EventType)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if holdInsideDecisionMidSideSortKey(a.EventMidSide) != holdInsideDecisionMidSideSortKey(b.EventMidSide) {
		return holdInsideDecisionMidSideSortKey(a.EventMidSide) < holdInsideDecisionMidSideSortKey(b.EventMidSide)
	}
	return holdInsideDecisionPositionBucketSortKey(a.EventClosePositionBucket) < holdInsideDecisionPositionBucketSortKey(b.EventClosePositionBucket)
}

func holdInsideMidlineReactionEventTypeSortKey(eventType string) int {
	switch eventType {
	case HoldInsideMidlineReactionEventTouch:
		return 0
	case HoldInsideMidlineReactionEventCloseAcross:
		return 1
	default:
		return 99
	}
}

func maxIntInSlice(values []int) int {
	maxValue := 0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}
