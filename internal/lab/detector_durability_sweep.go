package lab

import "sort"

type DetectorDurabilitySweepRow struct {
	ProfileID                      string  `json:"profile_id"`
	IsBalancedBaseline             bool    `json:"is_balanced_baseline"`
	IsADXComparison                bool    `json:"is_adx_comparison"`
	Percentile                     float64 `json:"percentile"`
	MinConsecutiveBars             int     `json:"min_consecutive_bars"`
	UseBollinger                   bool    `json:"use_bollinger"`
	UseADX                         bool    `json:"use_adx"`
	LookbackDays                   int     `json:"lookback_days"`
	Split                          string  `json:"split"`
	HorizonBars                    int     `json:"horizon_bars"`
	ActiveBars                     int     `json:"active_bars"`
	TotalBars                      int     `json:"total_bars"`
	DutyCycle                      float64 `json:"duty_cycle"`
	DetectorEpisodes               int     `json:"detector_episodes"`
	AvgDetectorEpisodeLength       float64 `json:"avg_detector_episode_length"`
	MedianDetectorEpisodeLength    float64 `json:"median_detector_episode_length"`
	LongestDetectorEpisodeLength   int     `json:"longest_detector_episode_length"`
	DurabilityEpisodeCount         int     `json:"durability_episode_count"`
	AvgRawLengthBars               float64 `json:"avg_raw_length_bars"`
	AvgActiveLengthBars            float64 `json:"avg_active_length_bars"`
	AvgEpisodeWidthPct             float64 `json:"avg_episode_width_pct"`
	AvgNormalizedATR               float64 `json:"avg_normalized_atr"`
	AvgEndNormalizedATR            float64 `json:"avg_end_normalized_atr"`
	AvgWidthToATRRatio             float64 `json:"avg_width_to_atr_ratio"`
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

type DetectorDurabilitySliceRow struct {
	ProfileID                      string  `json:"profile_id"`
	IsBalancedBaseline             bool    `json:"is_balanced_baseline"`
	IsADXComparison                bool    `json:"is_adx_comparison"`
	Percentile                     float64 `json:"percentile"`
	MinConsecutiveBars             int     `json:"min_consecutive_bars"`
	UseBollinger                   bool    `json:"use_bollinger"`
	UseADX                         bool    `json:"use_adx"`
	LookbackDays                   int     `json:"lookback_days"`
	Split                          string  `json:"split"`
	HorizonBars                    int     `json:"horizon_bars"`
	RawLengthBucket                string  `json:"raw_length_bucket"`
	ActiveLengthBucket             string  `json:"active_length_bucket"`
	EpisodeWidthBucket             string  `json:"episode_width_bucket"`
	WidthToATRBucket               string  `json:"width_to_atr_bucket"`
	EpisodeCount                   int     `json:"episode_count"`
	AvgRawLengthBars               float64 `json:"avg_raw_length_bars"`
	AvgActiveLengthBars            float64 `json:"avg_active_length_bars"`
	AvgEpisodeWidthPct             float64 `json:"avg_episode_width_pct"`
	AvgNormalizedATR               float64 `json:"avg_normalized_atr"`
	AvgEndNormalizedATR            float64 `json:"avg_end_normalized_atr"`
	AvgWidthToATRRatio             float64 `json:"avg_width_to_atr_ratio"`
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

type DetectorDurabilityStabilityRow struct {
	ProfileID                          string  `json:"profile_id"`
	IsBalancedBaseline                 bool    `json:"is_balanced_baseline"`
	IsADXComparison                    bool    `json:"is_adx_comparison"`
	Percentile                         float64 `json:"percentile"`
	MinConsecutiveBars                 int     `json:"min_consecutive_bars"`
	UseBollinger                       bool    `json:"use_bollinger"`
	UseADX                             bool    `json:"use_adx"`
	LookbackDays                       int     `json:"lookback_days"`
	HorizonBars                        int     `json:"horizon_bars"`
	PeriodSplits                       int     `json:"period_splits"`
	PeriodEpisodeCount                 int     `json:"period_episode_count"`
	EpisodeCountMin                    int     `json:"episode_count_min"`
	EpisodeCountMax                    int     `json:"episode_count_max"`
	EpisodeCountDelta                  int     `json:"episode_count_delta"`
	DutyCycleMin                       float64 `json:"duty_cycle_min"`
	DutyCycleMax                       float64 `json:"duty_cycle_max"`
	DutyCycleDelta                     float64 `json:"duty_cycle_delta"`
	LabelReenteredRangeRateMin         float64 `json:"label_reentered_range_rate_min"`
	LabelReenteredRangeRateMax         float64 `json:"label_reentered_range_rate_max"`
	LabelReenteredRangeRateDelta       float64 `json:"label_reentered_range_rate_delta"`
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

type detectorDurabilityAccumulator struct {
	episodes              int
	rawLengthSum          float64
	activeLengthSum       float64
	widthPctSum           float64
	avgATRSum             float64
	endATRSum             float64
	widthToATRRatioSum    float64
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

func RunDetectorDurabilitySweep(candles []Candle, detectorCfg RangeDetectorConfig, durabilityCfg RangeRegimeDurabilityAuditConfig, splits []Split) ([]DetectorDurabilitySweepRow, []DetectorDurabilitySliceRow, []DetectorDurabilityStabilityRow, error) {
	detectorCfg = detectorCfg.withDefaults()
	if err := detectorCfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	durabilityCfg = durabilityCfg.withDefaults()
	if err := durabilityCfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	profiles := DefaultDetectorSweepProfiles(detectorCfg.LookbackDays)
	lookbackBars := detectorCfg.LookbackDays * detectorCfg.BarsPerDay
	atr := NormalizedATR(candles, detectorCfg.ATRPeriod)
	donchian := DonchianWidth(candles, detectorCfg.DonchianPeriod)
	bollinger := BollingerWidth(candles, detectorCfg.BollingerPeriod)
	adx := ADX(candles, detectorCfg.ADXPeriod)
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)

	percentiles := []float64{0.20, 0.30, 0.40}
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

	broadRows := make([]DetectorDurabilitySweepRow, 0, len(profiles)*len(splits)*len(durabilityCfg.HorizonsBars))
	sliceRows := []DetectorDurabilitySliceRow{}
	for _, profile := range profiles {
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
		dutyRows, _ := SummarizeDetectorSplits(candles, classifications, splits)
		episodeRows := detectorDurabilityEpisodeRows(candles, classifications, normalizedATR, durabilityCfg, splits, profile.ProfileID)
		broadRows = append(broadRows, detectorDurabilitySweepRows(profile, dutyRows, episodeRows, durabilityCfg.HorizonsBars, splits)...)
		sliceRows = append(sliceRows, detectorDurabilitySliceRows(profile, summarizeRangeRegimeDurability(episodeRows))...)
	}
	stabilityRows := detectorDurabilityStabilityRows(profiles, broadRows, durabilityCfg.HorizonsBars, splits)
	return broadRows, sliceRows, stabilityRows, nil
}

func detectorDurabilityEpisodeRows(candles []Candle, classifications []RangeClassification, normalizedATR []float64, cfg RangeRegimeDurabilityAuditConfig, splits []Split, detectorProfileID string) []RangeRegimeDurabilityEpisodeRow {
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, detectorProfileID)
	rows := make([]RangeRegimeDurabilityEpisodeRow, 0, len(episodes)*len(cfg.HorizonsBars))
	for _, episode := range episodes {
		for _, horizon := range cfg.HorizonsBars {
			row, ok := newRangeRegimeDurabilityEpisodeRow(candles, episode, horizon, cfg.QuickInvalidationBars)
			if !ok {
				continue
			}
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessRangeRegimeDurabilityEpisodeRow(rows[i], rows[j])
	})
	return rows
}

func detectorDurabilitySweepRows(profile DetectorSweepProfile, dutyRows []DetectorDutyCycleRow, episodeRows []RangeRegimeDurabilityEpisodeRow, horizons []int, splits []Split) []DetectorDurabilitySweepRow {
	dutyBySplit := map[string]DetectorDutyCycleRow{}
	for _, row := range dutyRows {
		dutyBySplit[row.Split] = row
	}

	accumulators := map[detectorDurabilityKey]*detectorDurabilityAccumulator{}
	for _, row := range episodeRows {
		addDetectorDurabilityRow(accumulators, row.Split, row.HorizonBars, row)
		if row.Split != fullSplitName {
			addDetectorDurabilityRow(accumulators, fullSplitName, row.HorizonBars, row)
		}
	}

	rows := make([]DetectorDurabilitySweepRow, 0, len(splits)*len(horizons))
	for _, split := range splits {
		dutyRow := dutyBySplit[split.Name]
		for _, horizon := range horizons {
			key := detectorDurabilityKey{split: split.Name, horizonBars: horizon}
			row := DetectorDurabilitySweepRow{
				ProfileID:                    profile.ProfileID,
				IsBalancedBaseline:           profile.IsBalancedBaseline,
				IsADXComparison:              profile.IsADXComparison,
				Percentile:                   profile.Percentile,
				MinConsecutiveBars:           profile.MinConsecutiveBars,
				UseBollinger:                 profile.UseBollinger,
				UseADX:                       profile.UseADX,
				LookbackDays:                 profile.LookbackDays,
				Split:                        split.Name,
				HorizonBars:                  horizon,
				ActiveBars:                   dutyRow.ActiveBars,
				TotalBars:                    dutyRow.TotalBars,
				DutyCycle:                    dutyRow.DutyCycle,
				DetectorEpisodes:             dutyRow.Episodes,
				AvgDetectorEpisodeLength:     dutyRow.AvgEpisodeLength,
				MedianDetectorEpisodeLength:  dutyRow.MedianEpisodeLength,
				LongestDetectorEpisodeLength: dutyRow.LongestEpisodeLength,
			}
			if acc := accumulators[key]; acc != nil {
				acc.addToSweepRow(&row)
			}
			rows = append(rows, row)
		}
	}
	return rows
}

type detectorDurabilityKey struct {
	split       string
	horizonBars int
}

func addDetectorDurabilityRow(accumulators map[detectorDurabilityKey]*detectorDurabilityAccumulator, split string, horizon int, row RangeRegimeDurabilityEpisodeRow) {
	key := detectorDurabilityKey{split: split, horizonBars: horizon}
	acc := accumulators[key]
	if acc == nil {
		acc = &detectorDurabilityAccumulator{}
		accumulators[key] = acc
	}
	acc.add(row)
}

func detectorDurabilitySliceRows(profile DetectorSweepProfile, summaryRows []RangeRegimeDurabilitySummaryRow) []DetectorDurabilitySliceRow {
	rows := make([]DetectorDurabilitySliceRow, 0, len(summaryRows))
	for _, summary := range summaryRows {
		rows = append(rows, DetectorDurabilitySliceRow{
			ProfileID:                      profile.ProfileID,
			IsBalancedBaseline:             profile.IsBalancedBaseline,
			IsADXComparison:                profile.IsADXComparison,
			Percentile:                     profile.Percentile,
			MinConsecutiveBars:             profile.MinConsecutiveBars,
			UseBollinger:                   profile.UseBollinger,
			UseADX:                         profile.UseADX,
			LookbackDays:                   profile.LookbackDays,
			Split:                          summary.Split,
			HorizonBars:                    summary.HorizonBars,
			RawLengthBucket:                summary.RawLengthBucket,
			ActiveLengthBucket:             summary.ActiveLengthBucket,
			EpisodeWidthBucket:             summary.EpisodeWidthBucket,
			WidthToATRBucket:               summary.WidthToATRBucket,
			EpisodeCount:                   summary.EpisodeCount,
			AvgRawLengthBars:               summary.AvgRawLengthBars,
			AvgActiveLengthBars:            summary.AvgActiveLengthBars,
			AvgEpisodeWidthPct:             summary.AvgEpisodeWidthPct,
			AvgNormalizedATR:               summary.AvgNormalizedATR,
			AvgEndNormalizedATR:            summary.AvgEndNormalizedATR,
			AvgWidthToATRRatio:             summary.AvgWidthToATRRatio,
			LabelReenteredRangeCount:       summary.LabelReenteredRangeCount,
			LabelPersistedInsideRangeCount: summary.LabelPersistedInsideRangeCount,
			LabelQuickInvalidatedCount:     summary.LabelQuickInvalidatedCount,
			LabelInvalidatedUpCount:        summary.LabelInvalidatedUpCount,
			LabelInvalidatedDownCount:      summary.LabelInvalidatedDownCount,
			LabelChoppedCount:              summary.LabelChoppedCount,
			LabelTrendedUpCount:            summary.LabelTrendedUpCount,
			LabelTrendedDownCount:          summary.LabelTrendedDownCount,
			LabelReenteredRangeRate:        summary.LabelReenteredRangeRate,
			LabelPersistedInsideRangeRate:  summary.LabelPersistedInsideRangeRate,
			LabelQuickInvalidatedRate:      summary.LabelQuickInvalidatedRate,
			LabelInvalidatedUpRate:         summary.LabelInvalidatedUpRate,
			LabelInvalidatedDownRate:       summary.LabelInvalidatedDownRate,
			LabelChoppedRate:               summary.LabelChoppedRate,
			LabelTrendedUpRate:             summary.LabelTrendedUpRate,
			LabelTrendedDownRate:           summary.LabelTrendedDownRate,
			LabelAvgCloseDriftPct:          summary.LabelAvgCloseDriftPct,
			LabelAvgMaxUpMovePct:           summary.LabelAvgMaxUpMovePct,
			LabelAvgMaxDownMovePct:         summary.LabelAvgMaxDownMovePct,
		})
	}
	return rows
}

func detectorDurabilityStabilityRows(profiles []DetectorSweepProfile, broadRows []DetectorDurabilitySweepRow, horizons []int, splits []Split) []DetectorDurabilityStabilityRow {
	rowsByKey := map[detectorDurabilityProfileHorizonKey]map[string]DetectorDurabilitySweepRow{}
	for _, row := range broadRows {
		key := detectorDurabilityProfileHorizonKey{profileID: row.ProfileID, horizonBars: row.HorizonBars}
		if rowsByKey[key] == nil {
			rowsByKey[key] = map[string]DetectorDurabilitySweepRow{}
		}
		rowsByKey[key][row.Split] = row
	}

	periodSplits := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplits = append(periodSplits, split)
		}
	}

	out := make([]DetectorDurabilityStabilityRow, 0, len(profiles)*len(horizons))
	for _, profile := range profiles {
		for _, horizon := range horizons {
			key := detectorDurabilityProfileHorizonKey{profileID: profile.ProfileID, horizonBars: horizon}
			row := newDetectorDurabilityStabilityRow(profile, horizon, rowsByKey[key], periodSplits)
			out = append(out, row)
		}
	}
	return out
}

type detectorDurabilityProfileHorizonKey struct {
	profileID   string
	horizonBars int
}

func newDetectorDurabilityStabilityRow(profile DetectorSweepProfile, horizon int, rowsBySplit map[string]DetectorDurabilitySweepRow, periodSplits []Split) DetectorDurabilityStabilityRow {
	row := DetectorDurabilityStabilityRow{
		ProfileID:          profile.ProfileID,
		IsBalancedBaseline: profile.IsBalancedBaseline,
		IsADXComparison:    profile.IsADXComparison,
		Percentile:         profile.Percentile,
		MinConsecutiveBars: profile.MinConsecutiveBars,
		UseBollinger:       profile.UseBollinger,
		UseADX:             profile.UseADX,
		LookbackDays:       profile.LookbackDays,
		HorizonBars:        horizon,
		PeriodSplits:       len(periodSplits),
	}
	first := true
	for _, split := range periodSplits {
		broadRow := rowsBySplit[split.Name]
		trendedRate := broadRow.LabelTrendedUpRate + broadRow.LabelTrendedDownRate
		if first {
			row.EpisodeCountMin = broadRow.DurabilityEpisodeCount
			row.EpisodeCountMax = broadRow.DurabilityEpisodeCount
			row.DutyCycleMin = broadRow.DutyCycle
			row.DutyCycleMax = broadRow.DutyCycle
			row.LabelReenteredRangeRateMin = broadRow.LabelReenteredRangeRate
			row.LabelReenteredRangeRateMax = broadRow.LabelReenteredRangeRate
			row.LabelPersistedInsideRangeRateMin = broadRow.LabelPersistedInsideRangeRate
			row.LabelPersistedInsideRangeRateMax = broadRow.LabelPersistedInsideRangeRate
			row.LabelQuickInvalidatedRateMin = broadRow.LabelQuickInvalidatedRate
			row.LabelQuickInvalidatedRateMax = broadRow.LabelQuickInvalidatedRate
			row.LabelChoppedRateMin = broadRow.LabelChoppedRate
			row.LabelChoppedRateMax = broadRow.LabelChoppedRate
			row.LabelTrendedRateMin = trendedRate
			row.LabelTrendedRateMax = trendedRate
			row.LabelAvgCloseDriftPctMin = broadRow.LabelAvgCloseDriftPct
			row.LabelAvgCloseDriftPctMax = broadRow.LabelAvgCloseDriftPct
			first = false
		} else {
			row.EpisodeCountMin = minInt(row.EpisodeCountMin, broadRow.DurabilityEpisodeCount)
			if broadRow.DurabilityEpisodeCount > row.EpisodeCountMax {
				row.EpisodeCountMax = broadRow.DurabilityEpisodeCount
			}
			row.DutyCycleMin = minFloat(row.DutyCycleMin, broadRow.DutyCycle)
			row.DutyCycleMax = maxFloat(row.DutyCycleMax, broadRow.DutyCycle)
			row.LabelReenteredRangeRateMin = minFloat(row.LabelReenteredRangeRateMin, broadRow.LabelReenteredRangeRate)
			row.LabelReenteredRangeRateMax = maxFloat(row.LabelReenteredRangeRateMax, broadRow.LabelReenteredRangeRate)
			row.LabelPersistedInsideRangeRateMin = minFloat(row.LabelPersistedInsideRangeRateMin, broadRow.LabelPersistedInsideRangeRate)
			row.LabelPersistedInsideRangeRateMax = maxFloat(row.LabelPersistedInsideRangeRateMax, broadRow.LabelPersistedInsideRangeRate)
			row.LabelQuickInvalidatedRateMin = minFloat(row.LabelQuickInvalidatedRateMin, broadRow.LabelQuickInvalidatedRate)
			row.LabelQuickInvalidatedRateMax = maxFloat(row.LabelQuickInvalidatedRateMax, broadRow.LabelQuickInvalidatedRate)
			row.LabelChoppedRateMin = minFloat(row.LabelChoppedRateMin, broadRow.LabelChoppedRate)
			row.LabelChoppedRateMax = maxFloat(row.LabelChoppedRateMax, broadRow.LabelChoppedRate)
			row.LabelTrendedRateMin = minFloat(row.LabelTrendedRateMin, trendedRate)
			row.LabelTrendedRateMax = maxFloat(row.LabelTrendedRateMax, trendedRate)
			row.LabelAvgCloseDriftPctMin = minFloat(row.LabelAvgCloseDriftPctMin, broadRow.LabelAvgCloseDriftPct)
			row.LabelAvgCloseDriftPctMax = maxFloat(row.LabelAvgCloseDriftPctMax, broadRow.LabelAvgCloseDriftPct)
		}
		row.PeriodEpisodeCount += broadRow.DurabilityEpisodeCount
	}
	row.EpisodeCountDelta = row.EpisodeCountMax - row.EpisodeCountMin
	row.DutyCycleDelta = row.DutyCycleMax - row.DutyCycleMin
	row.LabelReenteredRangeRateDelta = row.LabelReenteredRangeRateMax - row.LabelReenteredRangeRateMin
	row.LabelPersistedInsideRangeRateDelta = row.LabelPersistedInsideRangeRateMax - row.LabelPersistedInsideRangeRateMin
	row.LabelQuickInvalidatedRateDelta = row.LabelQuickInvalidatedRateMax - row.LabelQuickInvalidatedRateMin
	row.LabelChoppedRateDelta = row.LabelChoppedRateMax - row.LabelChoppedRateMin
	row.LabelTrendedRateDelta = row.LabelTrendedRateMax - row.LabelTrendedRateMin
	row.LabelAvgCloseDriftPctDelta = row.LabelAvgCloseDriftPctMax - row.LabelAvgCloseDriftPctMin
	return row
}

func (acc *detectorDurabilityAccumulator) add(row RangeRegimeDurabilityEpisodeRow) {
	acc.episodes++
	acc.rawLengthSum += float64(row.RawLengthBars)
	acc.activeLengthSum += float64(row.ActiveLengthBars)
	acc.widthPctSum += row.EpisodeWidthPct
	acc.avgATRSum += row.AvgNormalizedATR
	acc.endATRSum += row.EndNormalizedATR
	acc.widthToATRRatioSum += row.WidthToATRRatio
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

func (acc detectorDurabilityAccumulator) addToSweepRow(row *DetectorDurabilitySweepRow) {
	row.DurabilityEpisodeCount = acc.episodes
	row.LabelReenteredRangeCount = acc.reenteredRanges
	row.LabelPersistedInsideRangeCount = acc.persistedInsideRanges
	row.LabelQuickInvalidatedCount = acc.quickInvalidated
	row.LabelInvalidatedUpCount = acc.invalidatedUp
	row.LabelInvalidatedDownCount = acc.invalidatedDown
	row.LabelChoppedCount = acc.chopped
	row.LabelTrendedUpCount = acc.trendedUp
	row.LabelTrendedDownCount = acc.trendedDown
	if acc.episodes == 0 {
		return
	}
	count := float64(acc.episodes)
	row.AvgRawLengthBars = acc.rawLengthSum / count
	row.AvgActiveLengthBars = acc.activeLengthSum / count
	row.AvgEpisodeWidthPct = acc.widthPctSum / count
	row.AvgNormalizedATR = acc.avgATRSum / count
	row.AvgEndNormalizedATR = acc.endATRSum / count
	row.AvgWidthToATRRatio = acc.widthToATRRatioSum / count
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
}

func minFloat(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

func maxFloat(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}
