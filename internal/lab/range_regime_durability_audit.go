package lab

import (
	"fmt"
	"math"
	"sort"
)

const fullSplitName = "full_2021_2026"

type RangeRegimeDurabilityAuditConfig struct {
	HorizonsBars          []int
	QuickInvalidationBars int
	DetectorProfileID     string
}

type RangeRegimeDurabilityEpisodeRow struct {
	Split                     string  `json:"split"`
	EpisodeID                 int     `json:"episode_id"`
	StartIndex                int     `json:"start_index"`
	EndIndex                  int     `json:"end_index"`
	StartTime                 string  `json:"start_time"`
	EndTime                   string  `json:"end_time"`
	HorizonBars               int     `json:"horizon_bars"`
	DetectorProfileID         string  `json:"detector_profile_id"`
	RawLengthBars             int     `json:"raw_length_bars"`
	ActiveLengthBars          int     `json:"active_length_bars"`
	RawLengthBucket           string  `json:"raw_length_bucket"`
	ActiveLengthBucket        string  `json:"active_length_bucket"`
	EpisodeHigh               float64 `json:"episode_high"`
	EpisodeLow                float64 `json:"episode_low"`
	EpisodeEndClose           float64 `json:"episode_end_close"`
	EpisodeWidthPct           float64 `json:"episode_width_pct"`
	EpisodeWidthBucket        string  `json:"episode_width_bucket"`
	AvgNormalizedATR          float64 `json:"avg_normalized_atr"`
	EndNormalizedATR          float64 `json:"end_normalized_atr"`
	WidthToATRRatio           float64 `json:"width_to_atr_ratio"`
	WidthToATRBucket          string  `json:"width_to_atr_bucket"`
	LabelWindowStartIndex     int     `json:"label_window_start_index"`
	LabelWindowEndIndex       int     `json:"label_window_end_index"`
	LabelWindowStartTime      string  `json:"label_window_start_time"`
	LabelWindowEndTime        string  `json:"label_window_end_time"`
	LabelReenteredRange       bool    `json:"label_reentered_range"`
	LabelPersistedInsideRange bool    `json:"label_persisted_inside_range"`
	LabelQuickInvalidated     bool    `json:"label_quick_invalidated"`
	LabelInvalidatedUp        bool    `json:"label_invalidated_up"`
	LabelInvalidatedDown      bool    `json:"label_invalidated_down"`
	LabelChopped              bool    `json:"label_chopped"`
	LabelTrendedUp            bool    `json:"label_trended_up"`
	LabelTrendedDown          bool    `json:"label_trended_down"`
	LabelCloseDriftPct        float64 `json:"label_close_drift_pct"`
	LabelMaxUpMovePct         float64 `json:"label_max_up_move_pct"`
	LabelMaxDownMovePct       float64 `json:"label_max_down_move_pct"`
}

type RangeRegimeDurabilitySummaryRow struct {
	Split                          string  `json:"split"`
	HorizonBars                    int     `json:"horizon_bars"`
	RawLengthBucket                string  `json:"raw_length_bucket"`
	ActiveLengthBucket             string  `json:"active_length_bucket"`
	EpisodeWidthBucket             string  `json:"episode_width_bucket"`
	WidthToATRBucket               string  `json:"width_to_atr_bucket"`
	DetectorProfileID              string  `json:"detector_profile_id"`
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

func DefaultRangeRegimeDurabilityAuditConfig() RangeRegimeDurabilityAuditConfig {
	return RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{1, 3, 6, 12},
		QuickInvalidationBars: 3,
		DetectorProfileID:     BalancedDetectorProfileID,
	}
}

func RunRangeRegimeDurabilityAudit(candles []Candle, cfg RangeRegimeDurabilityAuditConfig, splits []Split) ([]RangeRegimeDurabilityEpisodeRow, []RangeRegimeDurabilitySummaryRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	detector := CompressionRangeDetector{Config: DefaultCompressionRangeDetectorConfig()}
	classifications, err := detector.Classify(candles)
	if err != nil {
		return nil, nil, err
	}
	return runRangeRegimeDurabilityAuditFromClassifications(candles, classifications, cfg, splits)
}

func runRangeRegimeDurabilityAuditFromClassifications(candles []Candle, classifications []RangeClassification, cfg RangeRegimeDurabilityAuditConfig, splits []Split) ([]RangeRegimeDurabilityEpisodeRow, []RangeRegimeDurabilitySummaryRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}

	normalizedATR := NormalizedATR(candles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, cfg.DetectorProfileID)
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
	return rows, summarizeRangeRegimeDurability(rows), nil
}

func (cfg RangeRegimeDurabilityAuditConfig) withDefaults() RangeRegimeDurabilityAuditConfig {
	defaults := DefaultRangeRegimeDurabilityAuditConfig()
	if len(cfg.HorizonsBars) == 0 && cfg.QuickInvalidationBars == 0 && cfg.DetectorProfileID == "" {
		return defaults
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickInvalidationBars == 0 {
		cfg.QuickInvalidationBars = defaults.QuickInvalidationBars
	}
	if cfg.DetectorProfileID == "" {
		cfg.DetectorProfileID = defaults.DetectorProfileID
	}
	return cfg
}

func (cfg RangeRegimeDurabilityAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("range regime durability audit horizon bars must be positive")
		}
	}
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("range regime durability audit quick invalidation bars must be positive")
	}
	return nil
}

type rangeRegimeDurabilityEpisode struct {
	EpisodeID          int
	Split              string
	StartIndex         int
	EndIndex           int
	High               float64
	Low                float64
	EndClose           float64
	RawLengthBars      int
	ActiveLengthBars   int
	RawLengthBucket    string
	ActiveLengthBucket string
	WidthPct           float64
	WidthBucket        string
	AvgNormalizedATR   float64
	EndNormalizedATR   float64
	WidthToATRRatio    float64
	WidthToATRBucket   string
	DetectorProfileID  string
}

func rangeRegimeDurabilityEpisodes(candles []Candle, classifications []RangeClassification, normalizedATR []float64, splits []Split, detectorProfileID string) []rangeRegimeDurabilityEpisode {
	var episodes []rangeRegimeDurabilityEpisode
	for i := 0; i < len(classifications); {
		if !classifications[i].RawActive {
			i++
			continue
		}
		start := i
		activeLength := 0
		high := candles[i].High
		low := candles[i].Low
		atrSum := 0.0
		atrCount := 0
		for i < len(classifications) && classifications[i].RawActive {
			if classifications[i].Active {
				activeLength++
			}
			if candles[i].High > high {
				high = candles[i].High
			}
			if candles[i].Low < low {
				low = candles[i].Low
			}
			if i < len(normalizedATR) && validNumber(normalizedATR[i]) {
				atrSum += normalizedATR[i]
				atrCount++
			}
			i++
		}
		end := i - 1
		if activeLength == 0 {
			continue
		}
		rawLength := end - start + 1
		endClose := candles[end].Close
		widthPct := 0.0
		if endClose > 0 {
			widthPct = (high - low) / endClose
		}
		avgATR := 0.0
		if atrCount > 0 {
			avgATR = atrSum / float64(atrCount)
		}
		endATR := 0.0
		if end < len(normalizedATR) && validNumber(normalizedATR[end]) {
			endATR = normalizedATR[end]
		}
		widthToATRRatio := 0.0
		if avgATR > 0 {
			widthToATRRatio = widthPct / avgATR
		}
		episodes = append(episodes, rangeRegimeDurabilityEpisode{
			EpisodeID:          len(episodes) + 1,
			Split:              splitNameForCloseTime(candles[end].CloseTime, splits),
			StartIndex:         start,
			EndIndex:           end,
			High:               high,
			Low:                low,
			EndClose:           endClose,
			RawLengthBars:      rawLength,
			ActiveLengthBars:   activeLength,
			RawLengthBucket:    barLengthBucket(rawLength),
			ActiveLengthBucket: barLengthBucket(activeLength),
			WidthPct:           widthPct,
			WidthBucket:        rangeWidthBucket(widthPct),
			AvgNormalizedATR:   avgATR,
			EndNormalizedATR:   endATR,
			WidthToATRRatio:    widthToATRRatio,
			WidthToATRBucket:   widthToATRBucket(widthToATRRatio),
			DetectorProfileID:  detectorProfileID,
		})
	}
	return episodes
}

func newRangeRegimeDurabilityEpisodeRow(candles []Candle, episode rangeRegimeDurabilityEpisode, horizon, quickInvalidationBars int) (RangeRegimeDurabilityEpisodeRow, bool) {
	if horizon <= 0 || quickInvalidationBars <= 0 || episode.EndIndex < 0 || episode.EndIndex+horizon >= len(candles) {
		return RangeRegimeDurabilityEpisodeRow{}, false
	}
	startIndex := episode.EndIndex + 1
	endIndex := episode.EndIndex + horizon
	future := candles[startIndex : endIndex+1]
	futureMaxHigh, futureMinLow := futureHighLow(future)
	lastClose := candles[endIndex].Close

	row := RangeRegimeDurabilityEpisodeRow{
		Split:                 episode.Split,
		EpisodeID:             episode.EpisodeID,
		StartIndex:            episode.StartIndex,
		EndIndex:              episode.EndIndex,
		StartTime:             candles[episode.StartIndex].CloseTime.Format(timeLayout),
		EndTime:               candles[episode.EndIndex].CloseTime.Format(timeLayout),
		HorizonBars:           horizon,
		DetectorProfileID:     episode.DetectorProfileID,
		RawLengthBars:         episode.RawLengthBars,
		ActiveLengthBars:      episode.ActiveLengthBars,
		RawLengthBucket:       episode.RawLengthBucket,
		ActiveLengthBucket:    episode.ActiveLengthBucket,
		EpisodeHigh:           episode.High,
		EpisodeLow:            episode.Low,
		EpisodeEndClose:       episode.EndClose,
		EpisodeWidthPct:       episode.WidthPct,
		EpisodeWidthBucket:    episode.WidthBucket,
		AvgNormalizedATR:      episode.AvgNormalizedATR,
		EndNormalizedATR:      episode.EndNormalizedATR,
		WidthToATRRatio:       episode.WidthToATRRatio,
		WidthToATRBucket:      episode.WidthToATRBucket,
		LabelWindowStartIndex: startIndex,
		LabelWindowEndIndex:   endIndex,
		LabelWindowStartTime:  candles[startIndex].CloseTime.Format(timeLayout),
		LabelWindowEndTime:    candles[endIndex].CloseTime.Format(timeLayout),
		LabelCloseDriftPct:    closeDriftPct(lastClose, episode.EndClose),
		LabelMaxUpMovePct:     movePct(math.Max(0, futureMaxHigh-episode.EndClose), episode.EndClose),
		LabelMaxDownMovePct:   movePct(math.Max(0, episode.EndClose-futureMinLow), episode.EndClose),
	}

	persistedInside := true
	quickLimit := minInt(len(future), quickInvalidationBars)
	for i, candle := range future {
		inside := candle.Close >= episode.Low && candle.Close <= episode.High
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
	return row, true
}

type rangeRegimeDurabilitySummaryKey struct {
	split              string
	horizonBars        int
	rawLengthBucket    string
	activeLengthBucket string
	widthBucket        string
	widthToATRBucket   string
	detectorProfileID  string
}

type rangeRegimeDurabilitySummaryAccumulator struct {
	key                   rangeRegimeDurabilitySummaryKey
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

func summarizeRangeRegimeDurability(rows []RangeRegimeDurabilityEpisodeRow) []RangeRegimeDurabilitySummaryRow {
	accumulators := map[rangeRegimeDurabilitySummaryKey]*rangeRegimeDurabilitySummaryAccumulator{}
	for _, row := range rows {
		addRangeRegimeDurabilitySummaryRow(accumulators, row, row.Split)
		if row.Split != fullSplitName {
			addRangeRegimeDurabilitySummaryRow(accumulators, row, fullSplitName)
		}
	}
	out := make([]RangeRegimeDurabilitySummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		out = append(out, acc.row())
	}
	sort.Slice(out, func(i, j int) bool {
		return lessRangeRegimeDurabilitySummaryRow(out[i], out[j])
	})
	return out
}

func addRangeRegimeDurabilitySummaryRow(accumulators map[rangeRegimeDurabilitySummaryKey]*rangeRegimeDurabilitySummaryAccumulator, row RangeRegimeDurabilityEpisodeRow, split string) {
	key := rangeRegimeDurabilitySummaryKey{
		split:              split,
		horizonBars:        row.HorizonBars,
		rawLengthBucket:    row.RawLengthBucket,
		activeLengthBucket: row.ActiveLengthBucket,
		widthBucket:        row.EpisodeWidthBucket,
		widthToATRBucket:   row.WidthToATRBucket,
		detectorProfileID:  row.DetectorProfileID,
	}
	acc := accumulators[key]
	if acc == nil {
		acc = &rangeRegimeDurabilitySummaryAccumulator{key: key}
		accumulators[key] = acc
	}
	acc.add(row)
}

func (acc *rangeRegimeDurabilitySummaryAccumulator) add(row RangeRegimeDurabilityEpisodeRow) {
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

func (acc rangeRegimeDurabilitySummaryAccumulator) row() RangeRegimeDurabilitySummaryRow {
	row := RangeRegimeDurabilitySummaryRow{
		Split:                          acc.key.split,
		HorizonBars:                    acc.key.horizonBars,
		RawLengthBucket:                acc.key.rawLengthBucket,
		ActiveLengthBucket:             acc.key.activeLengthBucket,
		EpisodeWidthBucket:             acc.key.widthBucket,
		WidthToATRBucket:               acc.key.widthToATRBucket,
		DetectorProfileID:              acc.key.detectorProfileID,
		EpisodeCount:                   acc.episodes,
		LabelReenteredRangeCount:       acc.reenteredRanges,
		LabelPersistedInsideRangeCount: acc.persistedInsideRanges,
		LabelQuickInvalidatedCount:     acc.quickInvalidated,
		LabelInvalidatedUpCount:        acc.invalidatedUp,
		LabelInvalidatedDownCount:      acc.invalidatedDown,
		LabelChoppedCount:              acc.chopped,
		LabelTrendedUpCount:            acc.trendedUp,
		LabelTrendedDownCount:          acc.trendedDown,
	}
	if acc.episodes == 0 {
		return row
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
	return row
}

func widthToATRBucket(ratio float64) string {
	switch {
	case !validNumber(ratio) || ratio <= 0:
		return "unknown"
	case ratio < 0.5:
		return "lt_0_5x"
	case ratio < 1:
		return "0_5_1x"
	case ratio < 2:
		return "1_2x"
	case ratio < 4:
		return "2_4x"
	default:
		return "gt_4x"
	}
}

func widthToATRBucketSortKey(bucket string) int {
	switch bucket {
	case "lt_0_5x":
		return 0
	case "0_5_1x":
		return 1
	case "1_2x":
		return 2
	case "2_4x":
		return 3
	case "gt_4x":
		return 4
	case "unknown":
		return 5
	default:
		return 99
	}
}

func closeDriftPct(close, startClose float64) float64 {
	if startClose <= 0 {
		return 0
	}
	return (close - startClose) / startClose
}

func lessRangeRegimeDurabilityEpisodeRow(a, b RangeRegimeDurabilityEpisodeRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.EpisodeID != b.EpisodeID {
		return a.EpisodeID < b.EpisodeID
	}
	return a.HorizonBars < b.HorizonBars
}

func lessRangeRegimeDurabilitySummaryRow(a, b RangeRegimeDurabilitySummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.DetectorProfileID != b.DetectorProfileID {
		return a.DetectorProfileID < b.DetectorProfileID
	}
	if barLengthBucketSortKey(a.RawLengthBucket) != barLengthBucketSortKey(b.RawLengthBucket) {
		return barLengthBucketSortKey(a.RawLengthBucket) < barLengthBucketSortKey(b.RawLengthBucket)
	}
	if barLengthBucketSortKey(a.ActiveLengthBucket) != barLengthBucketSortKey(b.ActiveLengthBucket) {
		return barLengthBucketSortKey(a.ActiveLengthBucket) < barLengthBucketSortKey(b.ActiveLengthBucket)
	}
	if rangeWidthBucketSortKey(a.EpisodeWidthBucket) != rangeWidthBucketSortKey(b.EpisodeWidthBucket) {
		return rangeWidthBucketSortKey(a.EpisodeWidthBucket) < rangeWidthBucketSortKey(b.EpisodeWidthBucket)
	}
	return widthToATRBucketSortKey(a.WidthToATRBucket) < widthToATRBucketSortKey(b.WidthToATRBucket)
}
