package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	CompressionBreakoutSideUp   = "up"
	CompressionBreakoutSideDown = "down"
)

type CompressionBreakoutAuditConfig struct {
	HorizonsBars         []int
	MaxBreakoutDelayBars int
	DetectorProfileID    string
}

type CompressionBreakoutCandidateRow struct {
	Split                                 string  `json:"split"`
	Side                                  string  `json:"side"`
	BreakoutDelayBars                     int     `json:"breakout_delay_bars"`
	HorizonBars                           int     `json:"horizon_bars"`
	EpisodeRawLengthBucket                string  `json:"episode_raw_length_bucket"`
	EpisodeActiveLengthBucket             string  `json:"episode_active_length_bucket"`
	EpisodeRangeWidthBucket               string  `json:"episode_range_width_bucket"`
	BreakoutMoveBucket                    string  `json:"breakout_move_bucket"`
	DecisionTrueRangeExpansionBucket      string  `json:"decision_true_range_expansion_bucket"`
	DetectorProfileID                     string  `json:"detector_profile_id"`
	CandidateCount                        int     `json:"candidate_count"`
	AvgEpisodeRawLengthBars               float64 `json:"avg_episode_raw_length_bars"`
	AvgEpisodeActiveLengthBars            float64 `json:"avg_episode_active_length_bars"`
	AvgEpisodeRangeWidthPct               float64 `json:"avg_episode_range_width_pct"`
	AvgBreakoutMovePct                    float64 `json:"avg_breakout_move_pct"`
	AvgDecisionTrueRangeATR               float64 `json:"avg_decision_true_range_atr"`
	LabelReenteredRangeCount              int     `json:"label_reentered_range_count"`
	LabelOppositeCloseBreakCount          int     `json:"label_opposite_close_break_count"`
	LabelFavorableGreaterThanAdverseCount int     `json:"label_favorable_greater_than_adverse_count"`
	LabelReenteredRangeRate               float64 `json:"label_reentered_range_rate"`
	LabelOppositeCloseBreakRate           float64 `json:"label_opposite_close_break_rate"`
	LabelAvgFavorablePct                  float64 `json:"label_avg_favorable_pct"`
	LabelAvgAdversePct                    float64 `json:"label_avg_adverse_pct"`
	LabelFavorableMinusAdversePct         float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGreaterThanAdverseRate  float64 `json:"label_favorable_greater_than_adverse_rate"`
}

type CompressionBreakoutSummaryRow struct {
	Split                                string  `json:"split"`
	Side                                 string  `json:"side"`
	HorizonBars                          int     `json:"horizon_bars"`
	DetectorProfileID                    string  `json:"detector_profile_id"`
	CandidateCount                       int     `json:"candidate_count"`
	AvgBreakoutDelayBars                 float64 `json:"avg_breakout_delay_bars"`
	AvgEpisodeRawLengthBars              float64 `json:"avg_episode_raw_length_bars"`
	AvgEpisodeActiveLengthBars           float64 `json:"avg_episode_active_length_bars"`
	AvgEpisodeRangeWidthPct              float64 `json:"avg_episode_range_width_pct"`
	AvgBreakoutMovePct                   float64 `json:"avg_breakout_move_pct"`
	AvgDecisionTrueRangeATR              float64 `json:"avg_decision_true_range_atr"`
	LabelReenteredRangeRate              float64 `json:"label_reentered_range_rate"`
	LabelOppositeCloseBreakRate          float64 `json:"label_opposite_close_break_rate"`
	LabelAvgFavorablePct                 float64 `json:"label_avg_favorable_pct"`
	LabelAvgAdversePct                   float64 `json:"label_avg_adverse_pct"`
	LabelFavorableMinusAdversePct        float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGreaterThanAdverseRate float64 `json:"label_favorable_greater_than_adverse_rate"`
}

func DefaultCompressionBreakoutAuditConfig() CompressionBreakoutAuditConfig {
	return CompressionBreakoutAuditConfig{
		HorizonsBars:         []int{1, 3, 6, 12},
		MaxBreakoutDelayBars: 12,
		DetectorProfileID:    BalancedDetectorProfileID,
	}
}

func RunCompressionBreakoutAudit(candles []Candle, cfg CompressionBreakoutAuditConfig, splits []Split) ([]CompressionBreakoutCandidateRow, []CompressionBreakoutSummaryRow, error) {
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
	return runCompressionBreakoutAuditFromClassifications(candles, classifications, cfg, splits)
}

func runCompressionBreakoutAuditFromClassifications(candles []Candle, classifications []RangeClassification, cfg CompressionBreakoutAuditConfig, splits []Split) ([]CompressionBreakoutCandidateRow, []CompressionBreakoutSummaryRow, error) {
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

	atr := ATR(candles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
	var events []compressionBreakoutEvent
	for _, episode := range compressionBreakoutEpisodes(candles, classifications, splits, cfg.DetectorProfileID) {
		decision, ok := newCompressionBreakoutDecision(candles, atr, episode, cfg.MaxBreakoutDelayBars)
		if !ok {
			continue
		}
		for _, horizon := range cfg.HorizonsBars {
			label, ok := newCompressionBreakoutLabel(candles, decision, horizon)
			if !ok {
				continue
			}
			events = append(events, compressionBreakoutEvent{
				compressionBreakoutDecision:      decision,
				HorizonBars:                      horizon,
				LabelReenteredRange:              label.ReenteredRange,
				LabelOppositeCloseBreak:          label.OppositeCloseBreak,
				LabelFavorableGreaterThanAdverse: label.FavorableGreaterThanAdverse,
				LabelFavorableMovePct:            label.FavorableMovePct,
				LabelAdverseMovePct:              label.AdverseMovePct,
			})
		}
	}
	return summarizeCompressionBreakoutCandidates(events), summarizeCompressionBreakoutSummary(events), nil
}

func (cfg CompressionBreakoutAuditConfig) withDefaults() CompressionBreakoutAuditConfig {
	defaults := DefaultCompressionBreakoutAuditConfig()
	if len(cfg.HorizonsBars) == 0 && cfg.MaxBreakoutDelayBars == 0 && cfg.DetectorProfileID == "" {
		return defaults
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.MaxBreakoutDelayBars == 0 {
		cfg.MaxBreakoutDelayBars = defaults.MaxBreakoutDelayBars
	}
	if cfg.DetectorProfileID == "" {
		cfg.DetectorProfileID = defaults.DetectorProfileID
	}
	return cfg
}

func (cfg CompressionBreakoutAuditConfig) validate() error {
	if cfg.MaxBreakoutDelayBars <= 0 {
		return fmt.Errorf("compression breakout audit max breakout delay bars must be positive")
	}
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("compression breakout audit horizon bars must be positive")
		}
	}
	return nil
}

type compressionBreakoutEpisode struct {
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
	RangeWidthPct      float64
	RangeWidthBucket   string
	DetectorProfileID  string
}

type compressionBreakoutDecision struct {
	compressionBreakoutEpisode
	BreakoutIndex                    int
	BreakoutDelayBars                int
	Side                             string
	BreakoutClose                    float64
	BreakoutMovePct                  float64
	BreakoutMoveBucket               string
	DecisionTrueRangeATR             float64
	DecisionTrueRangeExpansionBucket string
}

type compressionBreakoutLabel struct {
	ReenteredRange              bool
	OppositeCloseBreak          bool
	FavorableGreaterThanAdverse bool
	FavorableMovePct            float64
	AdverseMovePct              float64
}

type compressionBreakoutEvent struct {
	compressionBreakoutDecision
	HorizonBars                      int
	LabelReenteredRange              bool
	LabelOppositeCloseBreak          bool
	LabelFavorableGreaterThanAdverse bool
	LabelFavorableMovePct            float64
	LabelAdverseMovePct              float64
}

func compressionBreakoutEpisodes(candles []Candle, classifications []RangeClassification, splits []Split, detectorProfileID string) []compressionBreakoutEpisode {
	var episodes []compressionBreakoutEpisode
	for i := 0; i < len(classifications); {
		if !classifications[i].RawActive {
			i++
			continue
		}
		start := i
		activeLength := 0
		high := candles[i].High
		low := candles[i].Low
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
			i++
		}
		end := i - 1
		if activeLength == 0 {
			continue
		}
		rawLength := end - start + 1
		endClose := candles[end].Close
		rangeWidthPct := 0.0
		if endClose > 0 {
			rangeWidthPct = (high - low) / endClose
		}
		episodes = append(episodes, compressionBreakoutEpisode{
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
			RangeWidthPct:      rangeWidthPct,
			RangeWidthBucket:   rangeWidthBucket(rangeWidthPct),
			DetectorProfileID:  detectorProfileID,
		})
	}
	return episodes
}

func newCompressionBreakoutDecision(candles []Candle, atr []float64, episode compressionBreakoutEpisode, maxDelay int) (compressionBreakoutDecision, bool) {
	if episode.EndIndex < 0 || episode.EndIndex >= len(candles)-1 || maxDelay <= 0 {
		return compressionBreakoutDecision{}, false
	}
	last := episode.EndIndex + maxDelay
	if last >= len(candles) {
		last = len(candles) - 1
	}
	for i := episode.EndIndex + 1; i <= last; i++ {
		side := ""
		move := 0.0
		switch {
		case candles[i].Close > episode.High:
			side = CompressionBreakoutSideUp
			move = candles[i].Close - episode.High
		case candles[i].Close < episode.Low:
			side = CompressionBreakoutSideDown
			move = episode.Low - candles[i].Close
		default:
			continue
		}
		movePct := movePct(move, candles[i].Close)
		decisionTrueRangeATR := trueRangeATR(candles, atr, i)
		return compressionBreakoutDecision{
			compressionBreakoutEpisode:       episode,
			BreakoutIndex:                    i,
			BreakoutDelayBars:                i - episode.EndIndex,
			Side:                             side,
			BreakoutClose:                    candles[i].Close,
			BreakoutMovePct:                  movePct,
			BreakoutMoveBucket:               distanceBucket(movePct),
			DecisionTrueRangeATR:             decisionTrueRangeATR,
			DecisionTrueRangeExpansionBucket: trueRangeExpansionBucket(decisionTrueRangeATR),
		}, true
	}
	return compressionBreakoutDecision{}, false
}

func newCompressionBreakoutLabel(candles []Candle, decision compressionBreakoutDecision, horizon int) (compressionBreakoutLabel, bool) {
	if horizon <= 0 || decision.BreakoutIndex < 0 || decision.BreakoutIndex+horizon >= len(candles) {
		return compressionBreakoutLabel{}, false
	}
	future := candles[decision.BreakoutIndex+1 : decision.BreakoutIndex+horizon+1]
	futureMaxHigh, futureMinLow := futureHighLow(future)
	label := compressionBreakoutLabel{}
	for _, candle := range future {
		if candle.Close >= decision.Low && candle.Close <= decision.High {
			label.ReenteredRange = true
		}
		switch decision.Side {
		case CompressionBreakoutSideUp:
			if candle.Close < decision.Low {
				label.OppositeCloseBreak = true
			}
		case CompressionBreakoutSideDown:
			if candle.Close > decision.High {
				label.OppositeCloseBreak = true
			}
		default:
			return compressionBreakoutLabel{}, false
		}
	}
	switch decision.Side {
	case CompressionBreakoutSideUp:
		label.FavorableMovePct = movePct(math.Max(0, futureMaxHigh-decision.BreakoutClose), decision.BreakoutClose)
		label.AdverseMovePct = movePct(math.Max(0, decision.BreakoutClose-futureMinLow), decision.BreakoutClose)
	case CompressionBreakoutSideDown:
		label.FavorableMovePct = movePct(math.Max(0, decision.BreakoutClose-futureMinLow), decision.BreakoutClose)
		label.AdverseMovePct = movePct(math.Max(0, futureMaxHigh-decision.BreakoutClose), decision.BreakoutClose)
	}
	label.FavorableGreaterThanAdverse = label.FavorableMovePct > label.AdverseMovePct
	return label, true
}

type compressionBreakoutCandidateKey struct {
	split                            string
	side                             string
	breakoutDelayBars                int
	horizonBars                      int
	episodeRawLengthBucket           string
	episodeActiveLengthBucket        string
	episodeRangeWidthBucket          string
	breakoutMoveBucket               string
	decisionTrueRangeExpansionBucket string
	detectorProfileID                string
}

type compressionBreakoutSummaryKey struct {
	split             string
	side              string
	horizonBars       int
	detectorProfileID string
}

type compressionBreakoutAccumulator struct {
	candidates                   int
	breakoutDelaySum             float64
	rawLengthSum                 float64
	activeLengthSum              float64
	rangeWidthPctSum             float64
	breakoutMovePctSum           float64
	decisionTrueRangeATRSum      float64
	reenteredRanges              int
	oppositeCloseBreaks          int
	favorableGreaterThanAdverses int
	favorableMovePctSum          float64
	adverseMovePctSum            float64
}

type compressionBreakoutCandidateAccumulator struct {
	key compressionBreakoutCandidateKey
	compressionBreakoutAccumulator
}

type compressionBreakoutSummaryAccumulator struct {
	key compressionBreakoutSummaryKey
	compressionBreakoutAccumulator
}

func summarizeCompressionBreakoutCandidates(events []compressionBreakoutEvent) []CompressionBreakoutCandidateRow {
	accumulators := map[compressionBreakoutCandidateKey]*compressionBreakoutCandidateAccumulator{}
	for _, event := range events {
		key := compressionBreakoutCandidateKey{
			split:                            event.Split,
			side:                             event.Side,
			breakoutDelayBars:                event.BreakoutDelayBars,
			horizonBars:                      event.HorizonBars,
			episodeRawLengthBucket:           event.RawLengthBucket,
			episodeActiveLengthBucket:        event.ActiveLengthBucket,
			episodeRangeWidthBucket:          event.RangeWidthBucket,
			breakoutMoveBucket:               event.BreakoutMoveBucket,
			decisionTrueRangeExpansionBucket: event.DecisionTrueRangeExpansionBucket,
			detectorProfileID:                event.DetectorProfileID,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &compressionBreakoutCandidateAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}
	rows := make([]CompressionBreakoutCandidateRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessCompressionBreakoutCandidateRow(rows[i], rows[j])
	})
	return rows
}

func summarizeCompressionBreakoutSummary(events []compressionBreakoutEvent) []CompressionBreakoutSummaryRow {
	accumulators := map[compressionBreakoutSummaryKey]*compressionBreakoutSummaryAccumulator{}
	for _, event := range events {
		key := compressionBreakoutSummaryKey{
			split:             event.Split,
			side:              event.Side,
			horizonBars:       event.HorizonBars,
			detectorProfileID: event.DetectorProfileID,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &compressionBreakoutSummaryAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}
	rows := make([]CompressionBreakoutSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessCompressionBreakoutSummaryRow(rows[i], rows[j])
	})
	return rows
}

func (acc *compressionBreakoutAccumulator) add(event compressionBreakoutEvent) {
	acc.candidates++
	acc.breakoutDelaySum += float64(event.BreakoutDelayBars)
	acc.rawLengthSum += float64(event.RawLengthBars)
	acc.activeLengthSum += float64(event.ActiveLengthBars)
	acc.rangeWidthPctSum += event.RangeWidthPct
	acc.breakoutMovePctSum += event.BreakoutMovePct
	acc.decisionTrueRangeATRSum += event.DecisionTrueRangeATR
	if event.LabelReenteredRange {
		acc.reenteredRanges++
	}
	if event.LabelOppositeCloseBreak {
		acc.oppositeCloseBreaks++
	}
	if event.LabelFavorableGreaterThanAdverse {
		acc.favorableGreaterThanAdverses++
	}
	acc.favorableMovePctSum += event.LabelFavorableMovePct
	acc.adverseMovePctSum += event.LabelAdverseMovePct
}

func (acc compressionBreakoutCandidateAccumulator) row() CompressionBreakoutCandidateRow {
	row := CompressionBreakoutCandidateRow{
		Split:                                 acc.key.split,
		Side:                                  acc.key.side,
		BreakoutDelayBars:                     acc.key.breakoutDelayBars,
		HorizonBars:                           acc.key.horizonBars,
		EpisodeRawLengthBucket:                acc.key.episodeRawLengthBucket,
		EpisodeActiveLengthBucket:             acc.key.episodeActiveLengthBucket,
		EpisodeRangeWidthBucket:               acc.key.episodeRangeWidthBucket,
		BreakoutMoveBucket:                    acc.key.breakoutMoveBucket,
		DecisionTrueRangeExpansionBucket:      acc.key.decisionTrueRangeExpansionBucket,
		DetectorProfileID:                     acc.key.detectorProfileID,
		CandidateCount:                        acc.candidates,
		LabelReenteredRangeCount:              acc.reenteredRanges,
		LabelOppositeCloseBreakCount:          acc.oppositeCloseBreaks,
		LabelFavorableGreaterThanAdverseCount: acc.favorableGreaterThanAdverses,
	}
	acc.addAveragesToCandidateRow(&row)
	return row
}

func (acc compressionBreakoutSummaryAccumulator) row() CompressionBreakoutSummaryRow {
	row := CompressionBreakoutSummaryRow{
		Split:             acc.key.split,
		Side:              acc.key.side,
		HorizonBars:       acc.key.horizonBars,
		DetectorProfileID: acc.key.detectorProfileID,
		CandidateCount:    acc.candidates,
	}
	acc.addAveragesToSummaryRow(&row)
	return row
}

func (acc compressionBreakoutAccumulator) addAveragesToCandidateRow(row *CompressionBreakoutCandidateRow) {
	if acc.candidates == 0 {
		return
	}
	count := float64(acc.candidates)
	row.AvgEpisodeRawLengthBars = acc.rawLengthSum / count
	row.AvgEpisodeActiveLengthBars = acc.activeLengthSum / count
	row.AvgEpisodeRangeWidthPct = acc.rangeWidthPctSum / count
	row.AvgBreakoutMovePct = acc.breakoutMovePctSum / count
	row.AvgDecisionTrueRangeATR = acc.decisionTrueRangeATRSum / count
	row.LabelReenteredRangeRate = float64(acc.reenteredRanges) / count
	row.LabelOppositeCloseBreakRate = float64(acc.oppositeCloseBreaks) / count
	row.LabelAvgFavorablePct = acc.favorableMovePctSum / count
	row.LabelAvgAdversePct = acc.adverseMovePctSum / count
	row.LabelFavorableMinusAdversePct = row.LabelAvgFavorablePct - row.LabelAvgAdversePct
	row.LabelFavorableGreaterThanAdverseRate = float64(acc.favorableGreaterThanAdverses) / count
}

func (acc compressionBreakoutAccumulator) addAveragesToSummaryRow(row *CompressionBreakoutSummaryRow) {
	if acc.candidates == 0 {
		return
	}
	count := float64(acc.candidates)
	row.AvgBreakoutDelayBars = acc.breakoutDelaySum / count
	row.AvgEpisodeRawLengthBars = acc.rawLengthSum / count
	row.AvgEpisodeActiveLengthBars = acc.activeLengthSum / count
	row.AvgEpisodeRangeWidthPct = acc.rangeWidthPctSum / count
	row.AvgBreakoutMovePct = acc.breakoutMovePctSum / count
	row.AvgDecisionTrueRangeATR = acc.decisionTrueRangeATRSum / count
	row.LabelReenteredRangeRate = float64(acc.reenteredRanges) / count
	row.LabelOppositeCloseBreakRate = float64(acc.oppositeCloseBreaks) / count
	row.LabelAvgFavorablePct = acc.favorableMovePctSum / count
	row.LabelAvgAdversePct = acc.adverseMovePctSum / count
	row.LabelFavorableMinusAdversePct = row.LabelAvgFavorablePct - row.LabelAvgAdversePct
	row.LabelFavorableGreaterThanAdverseRate = float64(acc.favorableGreaterThanAdverses) / count
}

func barLengthBucket(length int) string {
	switch {
	case length < 12:
		return "lt_12"
	case length < 24:
		return "12_23"
	case length < 48:
		return "24_47"
	default:
		return "48plus"
	}
}

func rangeWidthBucket(widthPct float64) string {
	switch {
	case widthPct <= 0.0010:
		return "0_10bp"
	case widthPct <= 0.0025:
		return "10_25bp"
	case widthPct <= 0.0050:
		return "25_50bp"
	default:
		return "gt_50bp"
	}
}

func trueRangeExpansionBucket(ratio float64) string {
	switch {
	case !validNumber(ratio) || ratio <= 0:
		return "unknown"
	case ratio < 1:
		return "lt_1x"
	case ratio < 1.5:
		return "1_1_5x"
	case ratio < 2:
		return "1_5_2x"
	default:
		return "gt_2x"
	}
}

func trueRangeATR(candles []Candle, atr []float64, index int) float64 {
	if index < 0 || index >= len(candles) || index >= len(atr) || !validNumber(atr[index]) || atr[index] <= 0 {
		return 0
	}
	return trueRangeAt(candles, index) / atr[index]
}

func trueRangeAt(candles []Candle, index int) float64 {
	if index < 0 || index >= len(candles) {
		return 0
	}
	candle := candles[index]
	tr := candle.High - candle.Low
	if index == 0 {
		return tr
	}
	highGap := math.Abs(candle.High - candles[index-1].Close)
	lowGap := math.Abs(candle.Low - candles[index-1].Close)
	return math.Max(tr, math.Max(highGap, lowGap))
}

func lessCompressionBreakoutCandidateRow(a, b CompressionBreakoutCandidateRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if compressionBreakoutSideSortKey(a.Side) != compressionBreakoutSideSortKey(b.Side) {
		return compressionBreakoutSideSortKey(a.Side) < compressionBreakoutSideSortKey(b.Side)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.BreakoutDelayBars != b.BreakoutDelayBars {
		return a.BreakoutDelayBars < b.BreakoutDelayBars
	}
	if a.DetectorProfileID != b.DetectorProfileID {
		return a.DetectorProfileID < b.DetectorProfileID
	}
	if barLengthBucketSortKey(a.EpisodeRawLengthBucket) != barLengthBucketSortKey(b.EpisodeRawLengthBucket) {
		return barLengthBucketSortKey(a.EpisodeRawLengthBucket) < barLengthBucketSortKey(b.EpisodeRawLengthBucket)
	}
	if barLengthBucketSortKey(a.EpisodeActiveLengthBucket) != barLengthBucketSortKey(b.EpisodeActiveLengthBucket) {
		return barLengthBucketSortKey(a.EpisodeActiveLengthBucket) < barLengthBucketSortKey(b.EpisodeActiveLengthBucket)
	}
	if rangeWidthBucketSortKey(a.EpisodeRangeWidthBucket) != rangeWidthBucketSortKey(b.EpisodeRangeWidthBucket) {
		return rangeWidthBucketSortKey(a.EpisodeRangeWidthBucket) < rangeWidthBucketSortKey(b.EpisodeRangeWidthBucket)
	}
	if distanceBucketSortKey(a.BreakoutMoveBucket) != distanceBucketSortKey(b.BreakoutMoveBucket) {
		return distanceBucketSortKey(a.BreakoutMoveBucket) < distanceBucketSortKey(b.BreakoutMoveBucket)
	}
	return trueRangeExpansionBucketSortKey(a.DecisionTrueRangeExpansionBucket) < trueRangeExpansionBucketSortKey(b.DecisionTrueRangeExpansionBucket)
}

func lessCompressionBreakoutSummaryRow(a, b CompressionBreakoutSummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if compressionBreakoutSideSortKey(a.Side) != compressionBreakoutSideSortKey(b.Side) {
		return compressionBreakoutSideSortKey(a.Side) < compressionBreakoutSideSortKey(b.Side)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	return a.DetectorProfileID < b.DetectorProfileID
}

func compressionBreakoutSideSortKey(side string) int {
	switch side {
	case CompressionBreakoutSideUp:
		return 0
	case CompressionBreakoutSideDown:
		return 1
	default:
		return 99
	}
}

func barLengthBucketSortKey(bucket string) int {
	switch bucket {
	case "lt_12":
		return 0
	case "12_23":
		return 1
	case "24_47":
		return 2
	case "48plus":
		return 3
	default:
		return 99
	}
}

func rangeWidthBucketSortKey(bucket string) int {
	switch bucket {
	case "0_10bp":
		return 0
	case "10_25bp":
		return 1
	case "25_50bp":
		return 2
	case "gt_50bp":
		return 3
	default:
		return 99
	}
}

func trueRangeExpansionBucketSortKey(bucket string) int {
	switch bucket {
	case "lt_1x":
		return 0
	case "1_1_5x":
		return 1
	case "1_5_2x":
		return 2
	case "gt_2x":
		return 3
	case "unknown":
		return 4
	default:
		return 99
	}
}
