package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	FuturesImpulseAbsorptionDirectionUp   = "up"
	FuturesImpulseAbsorptionDirectionDown = "down"
	FuturesImpulseAbsorptionDirectionAll  = "all"

	FuturesImpulseAbsorptionOutcomeReclaimFirst       = "reclaim_first"
	FuturesImpulseAbsorptionOutcomeContinuationFirst  = "continuation_first"
	FuturesImpulseAbsorptionOutcomeSameBarAmbiguous   = "same_bar_ambiguous"
	FuturesImpulseAbsorptionOutcomeNone               = "none"
	FuturesImpulseAbsorptionOutcomeMissingFuture      = "missing_future"
	FuturesImpulseAbsorptionBucketAll                 = "all"
	FuturesImpulseAbsorptionStopStateAuditReady       = "impulse_absorption_audit_ready"
	FuturesImpulseAbsorptionStopStateNoViableEdge     = "impulse_absorption_no_viable_edge"
	FuturesImpulseAbsorptionStopStateSourceGap        = "impulse_absorption_source_gap"
	FuturesImpulseAbsorptionStopStateCodegenOrBlocked = "impulse_absorption_codegen_or_test_blocked"
	FuturesImpulseAbsorptionStopStateNeedsReviewOnly  = "impulse_absorption_needs_review_only"
)

type FuturesImpulseAbsorptionAuditConfig struct {
	WarmupBars                   int
	TrueRangePercentileThreshold float64
	VolumePercentileThreshold    float64
	DownClosePositionMax         float64
	UpClosePositionMin           float64
	HorizonsBars                 []int
	QuickContinuationBars        int
}

type FuturesImpulseAbsorptionCandidateRow struct {
	EventID                   int     `json:"event_id"`
	EventIndex                int     `json:"event_index"`
	EventOpenTime             string  `json:"event_open_time"`
	EventCloseTime            string  `json:"event_close_time"`
	Split                     string  `json:"split"`
	Direction                 string  `json:"direction"`
	Open                      float64 `json:"open"`
	High                      float64 `json:"high"`
	Low                       float64 `json:"low"`
	Close                     float64 `json:"close"`
	Volume                    float64 `json:"volume"`
	PreviousClose             float64 `json:"previous_close"`
	TrueRange                 float64 `json:"true_range"`
	TrueRangePercentileRank   float64 `json:"true_range_percentile_rank"`
	VolumePercentileRank      float64 `json:"volume_percentile_rank"`
	TrueRangePercentileBucket string  `json:"true_range_percentile_bucket"`
	VolumePercentileBucket    string  `json:"volume_percentile_bucket"`
	EventRangePct             float64 `json:"event_range_pct"`
	ClosePosition             float64 `json:"close_position"`
	EventMidpoint             float64 `json:"event_midpoint"`
	HorizonBars               int     `json:"horizon_bars"`
	LabelWindowStartIndex     int     `json:"label_window_start_index"`
	LabelWindowEndIndex       int     `json:"label_window_end_index"`
	LabelWindowStartTime      string  `json:"label_window_start_time"`
	LabelWindowEndTime        string  `json:"label_window_end_time"`
	MidpointReclaim           bool    `json:"midpoint_reclaim"`
	ContinuationBeyondExtreme bool    `json:"continuation_beyond_extreme"`
	FirstOutcome              string  `json:"first_outcome"`
	SameBarAmbiguous          bool    `json:"same_bar_ambiguous"`
	QuickContinuation         bool    `json:"quick_continuation"`
	MissingFuture             bool    `json:"missing_future"`
	BarsToReclaim             int     `json:"bars_to_reclaim"`
	BarsToContinuation        int     `json:"bars_to_continuation"`
}

type FuturesImpulseAbsorptionSummaryRow struct {
	Split                          string  `json:"split"`
	Direction                      string  `json:"direction"`
	HorizonBars                    int     `json:"horizon_bars"`
	TrueRangePercentileBucket      string  `json:"true_range_percentile_bucket"`
	VolumePercentileBucket         string  `json:"volume_percentile_bucket"`
	SourceEventCount               int     `json:"source_event_count"`
	LabeledEventCount              int     `json:"labeled_event_count"`
	MissingFutureCount             int     `json:"missing_future_count"`
	MidpointReclaimCount           int     `json:"midpoint_reclaim_count"`
	ContinuationBeyondExtremeCount int     `json:"continuation_beyond_extreme_count"`
	ReclaimFirstCount              int     `json:"reclaim_first_count"`
	ContinuationFirstCount         int     `json:"continuation_first_count"`
	SameBarAmbiguousCount          int     `json:"same_bar_ambiguous_count"`
	QuickContinuationCount         int     `json:"quick_continuation_count"`
	MidpointReclaimRate            float64 `json:"midpoint_reclaim_rate"`
	ContinuationBeyondExtremeRate  float64 `json:"continuation_beyond_extreme_rate"`
	ReclaimFirstRate               float64 `json:"reclaim_first_rate"`
	ContinuationFirstRate          float64 `json:"continuation_first_rate"`
	SameBarAmbiguousRate           float64 `json:"same_bar_ambiguous_rate"`
	QuickContinuationRate          float64 `json:"quick_continuation_rate"`
	AvgBarsToReclaim               float64 `json:"avg_bars_to_reclaim"`
	AvgBarsToContinuation          float64 `json:"avg_bars_to_continuation"`
}

type FuturesImpulseAbsorptionStabilityRow struct {
	Direction                     string  `json:"direction"`
	HorizonBars                   int     `json:"horizon_bars"`
	TrueRangePercentileBucket     string  `json:"true_range_percentile_bucket"`
	VolumePercentileBucket        string  `json:"volume_percentile_bucket"`
	PeriodSplits                  int     `json:"period_splits"`
	SourceEventCount              int     `json:"source_event_count"`
	SourceEventCountMin           int     `json:"source_event_count_min"`
	SourceEventCountMax           int     `json:"source_event_count_max"`
	SourceEventCountDelta         int     `json:"source_event_count_delta"`
	LabeledEventCountMin          int     `json:"labeled_event_count_min"`
	LabeledEventCountMax          int     `json:"labeled_event_count_max"`
	LabeledEventCountDelta        int     `json:"labeled_event_count_delta"`
	ReclaimFirstRateMin           float64 `json:"reclaim_first_rate_min"`
	ReclaimFirstRateMax           float64 `json:"reclaim_first_rate_max"`
	ReclaimFirstRateDelta         float64 `json:"reclaim_first_rate_delta"`
	ContinuationFirstRateMin      float64 `json:"continuation_first_rate_min"`
	ContinuationFirstRateMax      float64 `json:"continuation_first_rate_max"`
	ContinuationFirstRateDelta    float64 `json:"continuation_first_rate_delta"`
	QuickContinuationRateMin      float64 `json:"quick_continuation_rate_min"`
	QuickContinuationRateMax      float64 `json:"quick_continuation_rate_max"`
	QuickContinuationRateDelta    float64 `json:"quick_continuation_rate_delta"`
	SameBarAmbiguousRateMin       float64 `json:"same_bar_ambiguous_rate_min"`
	SameBarAmbiguousRateMax       float64 `json:"same_bar_ambiguous_rate_max"`
	SameBarAmbiguousRateDelta     float64 `json:"same_bar_ambiguous_rate_delta"`
	ReclaimMinusContinuationMin   float64 `json:"reclaim_minus_continuation_min"`
	ReclaimMinusContinuationMax   float64 `json:"reclaim_minus_continuation_max"`
	ReclaimMinusContinuationDelta float64 `json:"reclaim_minus_continuation_delta"`
}

type futuresImpulseAbsorptionEvent struct {
	EventID                   int
	EventIndex                int
	Split                     string
	Direction                 string
	PreviousClose             float64
	TrueRange                 float64
	TrueRangePercentileRank   float64
	VolumePercentileRank      float64
	TrueRangePercentileBucket string
	VolumePercentileBucket    string
	EventRangePct             float64
	ClosePosition             float64
	EventMidpoint             float64
}

type futuresImpulseAbsorptionLabel struct {
	LabelWindowStartIndex     int
	LabelWindowEndIndex       int
	LabelWindowStartTime      string
	LabelWindowEndTime        string
	MidpointReclaim           bool
	ContinuationBeyondExtreme bool
	FirstOutcome              string
	SameBarAmbiguous          bool
	QuickContinuation         bool
	MissingFuture             bool
	BarsToReclaim             int
	BarsToContinuation        int
}

type futuresImpulseAbsorptionSummaryKey struct {
	split                     string
	direction                 string
	horizonBars               int
	trueRangePercentileBucket string
	volumePercentileBucket    string
}

type futuresImpulseAbsorptionSummaryAccumulator struct {
	key                     futuresImpulseAbsorptionSummaryKey
	sourceEvents            int
	labeledEvents           int
	missingFuture           int
	midpointReclaims        int
	continuations           int
	reclaimFirst            int
	continuationFirst       int
	sameBarAmbiguous        int
	quickContinuation       int
	barsToReclaimSum        float64
	barsToReclaimCount      int
	barsToContinuationSum   float64
	barsToContinuationCount int
}

type futuresImpulseAbsorptionStabilityKey struct {
	direction                 string
	horizonBars               int
	trueRangePercentileBucket string
	volumePercentileBucket    string
}

type futuresImpulseAbsorptionStabilityAccumulator struct {
	key  futuresImpulseAbsorptionStabilityKey
	rows []FuturesImpulseAbsorptionSummaryRow
}

func DefaultFuturesImpulseAbsorptionAuditConfig() FuturesImpulseAbsorptionAuditConfig {
	return FuturesImpulseAbsorptionAuditConfig{
		WarmupBars:                   30 * 24 * 12,
		TrueRangePercentileThreshold: 0.99,
		VolumePercentileThreshold:    0.95,
		DownClosePositionMax:         0.25,
		UpClosePositionMin:           0.75,
		HorizonsBars:                 []int{3, 6, 12, 24},
		QuickContinuationBars:        3,
	}
}

func RunFuturesImpulseAbsorptionAudit(candles []Candle, cfg FuturesImpulseAbsorptionAuditConfig, splits []Split) ([]FuturesImpulseAbsorptionCandidateRow, []FuturesImpulseAbsorptionSummaryRow, []FuturesImpulseAbsorptionStabilityRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	trueRanges := futuresImpulseTrueRanges(candles)
	volumes := futuresImpulseVolumes(candles)
	trueRangeRanks := rollingPriorPercentileRanks(trueRanges, cfg.WarmupBars)
	volumeRanks := rollingPriorPercentileRanks(volumes, cfg.WarmupBars)

	rows := []FuturesImpulseAbsorptionCandidateRow{}
	eventID := 0
	for i := maxInt(1, cfg.WarmupBars); i < len(candles); i++ {
		event, ok := newFuturesImpulseAbsorptionEvent(candles, trueRanges, trueRangeRanks, volumeRanks, splits, cfg, i, eventID+1)
		if !ok {
			continue
		}
		eventID++
		for _, horizon := range cfg.HorizonsBars {
			rows = append(rows, newFuturesImpulseAbsorptionCandidateRow(candles, event, horizon, cfg.QuickContinuationBars))
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessFuturesImpulseAbsorptionCandidate(rows[i], rows[j])
	})

	summaryRows := summarizeFuturesImpulseAbsorption(rows)
	stabilityRows := futuresImpulseAbsorptionStabilityRows(summaryRows, splits)
	return rows, summaryRows, stabilityRows, nil
}

func FuturesImpulseAbsorptionReviewStopState(summaryRows []FuturesImpulseAbsorptionSummaryRow, splits []Split) string {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	periodSplits := futuresImpulsePeriodSplitNames(splits)
	rowsByHorizonSplit := map[int]map[string]FuturesImpulseAbsorptionSummaryRow{}
	for _, row := range summaryRows {
		if row.Direction != FuturesImpulseAbsorptionDirectionAll ||
			row.TrueRangePercentileBucket != FuturesImpulseAbsorptionBucketAll ||
			row.VolumePercentileBucket != FuturesImpulseAbsorptionBucketAll ||
			row.Split == fullSplitName {
			continue
		}
		if rowsByHorizonSplit[row.HorizonBars] == nil {
			rowsByHorizonSplit[row.HorizonBars] = map[string]FuturesImpulseAbsorptionSummaryRow{}
		}
		rowsByHorizonSplit[row.HorizonBars][row.Split] = row
	}
	for _, splitRows := range rowsByHorizonSplit {
		passes := true
		for _, split := range periodSplits {
			row, ok := splitRows[split]
			if !ok || row.SourceEventCount < 100 || row.LabeledEventCount == 0 {
				passes = false
				break
			}
			margin := row.ReclaimFirstRate - row.ContinuationFirstRate
			if margin <= 0 || row.QuickContinuationRate >= row.ReclaimFirstRate || row.SameBarAmbiguousRate > margin {
				passes = false
				break
			}
		}
		if passes {
			return FuturesImpulseAbsorptionStopStateAuditReady
		}
	}
	return FuturesImpulseAbsorptionStopStateNoViableEdge
}

func (cfg FuturesImpulseAbsorptionAuditConfig) withDefaults() FuturesImpulseAbsorptionAuditConfig {
	defaults := DefaultFuturesImpulseAbsorptionAuditConfig()
	if cfg.WarmupBars == 0 {
		cfg.WarmupBars = defaults.WarmupBars
	}
	if cfg.TrueRangePercentileThreshold == 0 {
		cfg.TrueRangePercentileThreshold = defaults.TrueRangePercentileThreshold
	}
	if cfg.VolumePercentileThreshold == 0 {
		cfg.VolumePercentileThreshold = defaults.VolumePercentileThreshold
	}
	if cfg.DownClosePositionMax == 0 {
		cfg.DownClosePositionMax = defaults.DownClosePositionMax
	}
	if cfg.UpClosePositionMin == 0 {
		cfg.UpClosePositionMin = defaults.UpClosePositionMin
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickContinuationBars == 0 {
		cfg.QuickContinuationBars = defaults.QuickContinuationBars
	}
	return cfg
}

func (cfg FuturesImpulseAbsorptionAuditConfig) validate() error {
	if cfg.WarmupBars <= 0 {
		return fmt.Errorf("futures impulse absorption audit warmup bars must be positive")
	}
	if cfg.TrueRangePercentileThreshold <= 0 || cfg.TrueRangePercentileThreshold >= 1 {
		return fmt.Errorf("futures impulse absorption audit true-range percentile threshold must be between 0 and 1")
	}
	if cfg.VolumePercentileThreshold <= 0 || cfg.VolumePercentileThreshold >= 1 {
		return fmt.Errorf("futures impulse absorption audit volume percentile threshold must be between 0 and 1")
	}
	if cfg.DownClosePositionMax <= 0 || cfg.DownClosePositionMax >= 1 {
		return fmt.Errorf("futures impulse absorption audit down close-position max must be between 0 and 1")
	}
	if cfg.UpClosePositionMin <= 0 || cfg.UpClosePositionMin >= 1 {
		return fmt.Errorf("futures impulse absorption audit up close-position min must be between 0 and 1")
	}
	if cfg.DownClosePositionMax >= cfg.UpClosePositionMin {
		return fmt.Errorf("futures impulse absorption audit close-position cutoffs must leave a middle skip zone")
	}
	if cfg.QuickContinuationBars <= 0 {
		return fmt.Errorf("futures impulse absorption audit quick continuation bars must be positive")
	}
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("futures impulse absorption audit horizon bars must be positive")
		}
	}
	return nil
}

func futuresImpulseTrueRanges(candles []Candle) []float64 {
	out := nanSlice(len(candles))
	for i, candle := range candles {
		value := candle.High - candle.Low
		if i > 0 {
			value = math.Max(value, math.Max(math.Abs(candle.High-candles[i-1].Close), math.Abs(candle.Low-candles[i-1].Close)))
		}
		out[i] = value
	}
	return out
}

func futuresImpulseVolumes(candles []Candle) []float64 {
	out := nanSlice(len(candles))
	for i, candle := range candles {
		out[i] = candle.Volume
	}
	return out
}

func rollingPriorPercentileRanks(values []float64, lookback int) []float64 {
	out := nanSlice(len(values))
	if lookback <= 0 {
		return out
	}
	window := make([]float64, 0, lookback)
	for i, value := range values {
		if i >= lookback && len(window) > 0 && validNumber(value) {
			out[i] = percentileRankFromSorted(window, value)
		}
		if validNumber(value) {
			window = insertSorted(window, value)
		}
		if i >= lookback && validNumber(values[i-lookback]) {
			window = removeSorted(window, values[i-lookback])
		}
	}
	return out
}

func percentileRankFromSorted(sortedValues []float64, value float64) float64 {
	if len(sortedValues) == 0 || !validNumber(value) {
		return math.NaN()
	}
	countLTE := sort.Search(len(sortedValues), func(i int) bool {
		return sortedValues[i] > value
	})
	return float64(countLTE) / float64(len(sortedValues))
}

func newFuturesImpulseAbsorptionEvent(candles []Candle, trueRanges, trueRangeRanks, volumeRanks []float64, splits []Split, cfg FuturesImpulseAbsorptionAuditConfig, index int, eventID int) (futuresImpulseAbsorptionEvent, bool) {
	if index <= 0 || index >= len(candles) ||
		index >= len(trueRanges) || index >= len(trueRangeRanks) || index >= len(volumeRanks) ||
		!validNumber(trueRanges[index]) || !validNumber(trueRangeRanks[index]) || !validNumber(volumeRanks[index]) ||
		trueRangeRanks[index] < cfg.TrueRangePercentileThreshold ||
		volumeRanks[index] < cfg.VolumePercentileThreshold {
		return futuresImpulseAbsorptionEvent{}, false
	}
	candle := candles[index]
	rangeValue := candle.High - candle.Low
	if rangeValue <= 0 {
		return futuresImpulseAbsorptionEvent{}, false
	}
	closePosition := (candle.Close - candle.Low) / rangeValue
	direction := ""
	switch {
	case closePosition <= cfg.DownClosePositionMax:
		direction = FuturesImpulseAbsorptionDirectionDown
	case closePosition >= cfg.UpClosePositionMin:
		direction = FuturesImpulseAbsorptionDirectionUp
	default:
		return futuresImpulseAbsorptionEvent{}, false
	}
	return futuresImpulseAbsorptionEvent{
		EventID:                   eventID,
		EventIndex:                index,
		Split:                     splitNameForCloseTime(candle.CloseTime, splits),
		Direction:                 direction,
		PreviousClose:             candles[index-1].Close,
		TrueRange:                 trueRanges[index],
		TrueRangePercentileRank:   trueRangeRanks[index],
		VolumePercentileRank:      volumeRanks[index],
		TrueRangePercentileBucket: futuresImpulseTrueRangeRankBucket(trueRangeRanks[index]),
		VolumePercentileBucket:    futuresImpulseVolumeRankBucket(volumeRanks[index]),
		EventRangePct:             movePct(rangeValue, candle.Close),
		ClosePosition:             closePosition,
		EventMidpoint:             (candle.High + candle.Low) / 2,
	}, true
}

func newFuturesImpulseAbsorptionCandidateRow(candles []Candle, event futuresImpulseAbsorptionEvent, horizon int, quickContinuationBars int) FuturesImpulseAbsorptionCandidateRow {
	eventCandle := candles[event.EventIndex]
	label := newFuturesImpulseAbsorptionLabel(candles, event, horizon, quickContinuationBars)
	return FuturesImpulseAbsorptionCandidateRow{
		EventID:                   event.EventID,
		EventIndex:                event.EventIndex,
		EventOpenTime:             eventCandle.OpenTime.Format(timeLayout),
		EventCloseTime:            eventCandle.CloseTime.Format(timeLayout),
		Split:                     event.Split,
		Direction:                 event.Direction,
		Open:                      eventCandle.Open,
		High:                      eventCandle.High,
		Low:                       eventCandle.Low,
		Close:                     eventCandle.Close,
		Volume:                    eventCandle.Volume,
		PreviousClose:             event.PreviousClose,
		TrueRange:                 event.TrueRange,
		TrueRangePercentileRank:   event.TrueRangePercentileRank,
		VolumePercentileRank:      event.VolumePercentileRank,
		TrueRangePercentileBucket: event.TrueRangePercentileBucket,
		VolumePercentileBucket:    event.VolumePercentileBucket,
		EventRangePct:             event.EventRangePct,
		ClosePosition:             event.ClosePosition,
		EventMidpoint:             event.EventMidpoint,
		HorizonBars:               horizon,
		LabelWindowStartIndex:     label.LabelWindowStartIndex,
		LabelWindowEndIndex:       label.LabelWindowEndIndex,
		LabelWindowStartTime:      label.LabelWindowStartTime,
		LabelWindowEndTime:        label.LabelWindowEndTime,
		MidpointReclaim:           label.MidpointReclaim,
		ContinuationBeyondExtreme: label.ContinuationBeyondExtreme,
		FirstOutcome:              label.FirstOutcome,
		SameBarAmbiguous:          label.SameBarAmbiguous,
		QuickContinuation:         label.QuickContinuation,
		MissingFuture:             label.MissingFuture,
		BarsToReclaim:             label.BarsToReclaim,
		BarsToContinuation:        label.BarsToContinuation,
	}
}

func newFuturesImpulseAbsorptionLabel(candles []Candle, event futuresImpulseAbsorptionEvent, horizon int, quickContinuationBars int) futuresImpulseAbsorptionLabel {
	start := event.EventIndex + 1
	end := event.EventIndex + horizon
	label := futuresImpulseAbsorptionLabel{
		LabelWindowStartIndex: -1,
		LabelWindowEndIndex:   -1,
		FirstOutcome:          FuturesImpulseAbsorptionOutcomeMissingFuture,
		MissingFuture:         true,
	}
	if start < len(candles) {
		label.LabelWindowStartIndex = start
		label.LabelWindowStartTime = candles[start].CloseTime.Format(timeLayout)
		label.LabelWindowEndIndex = minInt(end, len(candles)-1)
		label.LabelWindowEndTime = candles[label.LabelWindowEndIndex].CloseTime.Format(timeLayout)
	}
	if horizon <= 0 || end >= len(candles) {
		return label
	}

	label.MissingFuture = false
	label.FirstOutcome = FuturesImpulseAbsorptionOutcomeNone
	label.LabelWindowEndIndex = end
	label.LabelWindowEndTime = candles[end].CloseTime.Format(timeLayout)
	quickLimit := minInt(horizon, quickContinuationBars)
	for i := start; i <= end; i++ {
		candle := candles[i]
		delay := i - event.EventIndex
		reclaim, continuation := futuresImpulseAbsorptionTouches(candle, candles[event.EventIndex], event)
		if reclaim && !label.MidpointReclaim {
			label.MidpointReclaim = true
			label.BarsToReclaim = delay
		}
		if continuation && !label.ContinuationBeyondExtreme {
			label.ContinuationBeyondExtreme = true
			label.BarsToContinuation = delay
		}
		if continuation && delay <= quickLimit {
			label.QuickContinuation = true
		}
		if label.FirstOutcome != FuturesImpulseAbsorptionOutcomeNone {
			continue
		}
		switch {
		case reclaim && continuation:
			label.FirstOutcome = FuturesImpulseAbsorptionOutcomeSameBarAmbiguous
			label.SameBarAmbiguous = true
		case reclaim:
			label.FirstOutcome = FuturesImpulseAbsorptionOutcomeReclaimFirst
		case continuation:
			label.FirstOutcome = FuturesImpulseAbsorptionOutcomeContinuationFirst
		}
	}
	return label
}

func futuresImpulseAbsorptionTouches(candle Candle, eventCandle Candle, event futuresImpulseAbsorptionEvent) (bool, bool) {
	switch event.Direction {
	case FuturesImpulseAbsorptionDirectionUp:
		return candle.Low <= event.EventMidpoint, candle.High > eventCandle.High
	case FuturesImpulseAbsorptionDirectionDown:
		return candle.High >= event.EventMidpoint, candle.Low < eventCandle.Low
	default:
		return false, false
	}
}

func summarizeFuturesImpulseAbsorption(rows []FuturesImpulseAbsorptionCandidateRow) []FuturesImpulseAbsorptionSummaryRow {
	accumulators := map[futuresImpulseAbsorptionSummaryKey]*futuresImpulseAbsorptionSummaryAccumulator{}
	for _, row := range rows {
		for _, split := range futuresImpulseAbsorptionSplitCombos(row.Split) {
			for _, direction := range []string{row.Direction, FuturesImpulseAbsorptionDirectionAll} {
				for _, trBucket := range []string{row.TrueRangePercentileBucket, FuturesImpulseAbsorptionBucketAll} {
					for _, volumeBucket := range []string{row.VolumePercentileBucket, FuturesImpulseAbsorptionBucketAll} {
						key := futuresImpulseAbsorptionSummaryKey{
							split:                     split,
							direction:                 direction,
							horizonBars:               row.HorizonBars,
							trueRangePercentileBucket: trBucket,
							volumePercentileBucket:    volumeBucket,
						}
						acc := accumulators[key]
						if acc == nil {
							acc = &futuresImpulseAbsorptionSummaryAccumulator{key: key}
							accumulators[key] = acc
						}
						acc.add(row)
					}
				}
			}
		}
	}
	out := make([]FuturesImpulseAbsorptionSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		out = append(out, acc.row())
	}
	sort.Slice(out, func(i, j int) bool {
		return lessFuturesImpulseAbsorptionSummary(out[i], out[j])
	})
	return out
}

func futuresImpulseAbsorptionStabilityRows(summaryRows []FuturesImpulseAbsorptionSummaryRow, splits []Split) []FuturesImpulseAbsorptionStabilityRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	periodSplitSet := map[string]bool{}
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplitSet[split.Name] = true
		}
	}
	accumulators := map[futuresImpulseAbsorptionStabilityKey]*futuresImpulseAbsorptionStabilityAccumulator{}
	for _, row := range summaryRows {
		if !periodSplitSet[row.Split] {
			continue
		}
		key := futuresImpulseAbsorptionStabilityKey{
			direction:                 row.Direction,
			horizonBars:               row.HorizonBars,
			trueRangePercentileBucket: row.TrueRangePercentileBucket,
			volumePercentileBucket:    row.VolumePercentileBucket,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &futuresImpulseAbsorptionStabilityAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.rows = append(acc.rows, row)
	}
	out := make([]FuturesImpulseAbsorptionStabilityRow, 0, len(accumulators))
	for _, acc := range accumulators {
		out = append(out, acc.row())
	}
	sort.Slice(out, func(i, j int) bool {
		return lessFuturesImpulseAbsorptionStability(out[i], out[j])
	})
	return out
}

func (acc *futuresImpulseAbsorptionSummaryAccumulator) add(row FuturesImpulseAbsorptionCandidateRow) {
	acc.sourceEvents++
	if row.MissingFuture {
		acc.missingFuture++
		return
	}
	acc.labeledEvents++
	if row.MidpointReclaim {
		acc.midpointReclaims++
		acc.barsToReclaimSum += float64(row.BarsToReclaim)
		acc.barsToReclaimCount++
	}
	if row.ContinuationBeyondExtreme {
		acc.continuations++
		acc.barsToContinuationSum += float64(row.BarsToContinuation)
		acc.barsToContinuationCount++
	}
	if row.QuickContinuation {
		acc.quickContinuation++
	}
	switch row.FirstOutcome {
	case FuturesImpulseAbsorptionOutcomeReclaimFirst:
		acc.reclaimFirst++
	case FuturesImpulseAbsorptionOutcomeContinuationFirst:
		acc.continuationFirst++
	case FuturesImpulseAbsorptionOutcomeSameBarAmbiguous:
		acc.sameBarAmbiguous++
	}
}

func (acc futuresImpulseAbsorptionSummaryAccumulator) row() FuturesImpulseAbsorptionSummaryRow {
	row := FuturesImpulseAbsorptionSummaryRow{
		Split:                          acc.key.split,
		Direction:                      acc.key.direction,
		HorizonBars:                    acc.key.horizonBars,
		TrueRangePercentileBucket:      acc.key.trueRangePercentileBucket,
		VolumePercentileBucket:         acc.key.volumePercentileBucket,
		SourceEventCount:               acc.sourceEvents,
		LabeledEventCount:              acc.labeledEvents,
		MissingFutureCount:             acc.missingFuture,
		MidpointReclaimCount:           acc.midpointReclaims,
		ContinuationBeyondExtremeCount: acc.continuations,
		ReclaimFirstCount:              acc.reclaimFirst,
		ContinuationFirstCount:         acc.continuationFirst,
		SameBarAmbiguousCount:          acc.sameBarAmbiguous,
		QuickContinuationCount:         acc.quickContinuation,
	}
	if acc.labeledEvents > 0 {
		denom := float64(acc.labeledEvents)
		row.MidpointReclaimRate = float64(acc.midpointReclaims) / denom
		row.ContinuationBeyondExtremeRate = float64(acc.continuations) / denom
		row.ReclaimFirstRate = float64(acc.reclaimFirst) / denom
		row.ContinuationFirstRate = float64(acc.continuationFirst) / denom
		row.SameBarAmbiguousRate = float64(acc.sameBarAmbiguous) / denom
		row.QuickContinuationRate = float64(acc.quickContinuation) / denom
	}
	if acc.barsToReclaimCount > 0 {
		row.AvgBarsToReclaim = acc.barsToReclaimSum / float64(acc.barsToReclaimCount)
	}
	if acc.barsToContinuationCount > 0 {
		row.AvgBarsToContinuation = acc.barsToContinuationSum / float64(acc.barsToContinuationCount)
	}
	return row
}

func (acc futuresImpulseAbsorptionStabilityAccumulator) row() FuturesImpulseAbsorptionStabilityRow {
	row := FuturesImpulseAbsorptionStabilityRow{
		Direction:                   acc.key.direction,
		HorizonBars:                 acc.key.horizonBars,
		TrueRangePercentileBucket:   acc.key.trueRangePercentileBucket,
		VolumePercentileBucket:      acc.key.volumePercentileBucket,
		PeriodSplits:                len(acc.rows),
		SourceEventCountMin:         math.MaxInt,
		LabeledEventCountMin:        math.MaxInt,
		ReclaimFirstRateMin:         math.Inf(1),
		ContinuationFirstRateMin:    math.Inf(1),
		QuickContinuationRateMin:    math.Inf(1),
		SameBarAmbiguousRateMin:     math.Inf(1),
		ReclaimMinusContinuationMin: math.Inf(1),
	}
	for _, summary := range acc.rows {
		row.SourceEventCount += summary.SourceEventCount
		row.SourceEventCountMin = minInt(row.SourceEventCountMin, summary.SourceEventCount)
		row.SourceEventCountMax = maxInt(row.SourceEventCountMax, summary.SourceEventCount)
		row.LabeledEventCountMin = minInt(row.LabeledEventCountMin, summary.LabeledEventCount)
		row.LabeledEventCountMax = maxInt(row.LabeledEventCountMax, summary.LabeledEventCount)
		row.ReclaimFirstRateMin = math.Min(row.ReclaimFirstRateMin, summary.ReclaimFirstRate)
		row.ReclaimFirstRateMax = math.Max(row.ReclaimFirstRateMax, summary.ReclaimFirstRate)
		row.ContinuationFirstRateMin = math.Min(row.ContinuationFirstRateMin, summary.ContinuationFirstRate)
		row.ContinuationFirstRateMax = math.Max(row.ContinuationFirstRateMax, summary.ContinuationFirstRate)
		row.QuickContinuationRateMin = math.Min(row.QuickContinuationRateMin, summary.QuickContinuationRate)
		row.QuickContinuationRateMax = math.Max(row.QuickContinuationRateMax, summary.QuickContinuationRate)
		row.SameBarAmbiguousRateMin = math.Min(row.SameBarAmbiguousRateMin, summary.SameBarAmbiguousRate)
		row.SameBarAmbiguousRateMax = math.Max(row.SameBarAmbiguousRateMax, summary.SameBarAmbiguousRate)
		margin := summary.ReclaimFirstRate - summary.ContinuationFirstRate
		row.ReclaimMinusContinuationMin = math.Min(row.ReclaimMinusContinuationMin, margin)
		row.ReclaimMinusContinuationMax = math.Max(row.ReclaimMinusContinuationMax, margin)
	}
	if len(acc.rows) == 0 {
		row.SourceEventCountMin = 0
		row.LabeledEventCountMin = 0
		row.ReclaimFirstRateMin = 0
		row.ContinuationFirstRateMin = 0
		row.QuickContinuationRateMin = 0
		row.SameBarAmbiguousRateMin = 0
		row.ReclaimMinusContinuationMin = 0
	}
	row.SourceEventCountDelta = row.SourceEventCountMax - row.SourceEventCountMin
	row.LabeledEventCountDelta = row.LabeledEventCountMax - row.LabeledEventCountMin
	row.ReclaimFirstRateDelta = row.ReclaimFirstRateMax - row.ReclaimFirstRateMin
	row.ContinuationFirstRateDelta = row.ContinuationFirstRateMax - row.ContinuationFirstRateMin
	row.QuickContinuationRateDelta = row.QuickContinuationRateMax - row.QuickContinuationRateMin
	row.SameBarAmbiguousRateDelta = row.SameBarAmbiguousRateMax - row.SameBarAmbiguousRateMin
	row.ReclaimMinusContinuationDelta = row.ReclaimMinusContinuationMax - row.ReclaimMinusContinuationMin
	return row
}

func futuresImpulseAbsorptionSplitCombos(split string) []string {
	if split == "" || split == fullSplitName {
		return []string{fullSplitName}
	}
	return []string{split, fullSplitName}
}

func futuresImpulsePeriodSplitNames(splits []Split) []string {
	out := []string{}
	for _, split := range splits {
		if split.Name != fullSplitName {
			out = append(out, split.Name)
		}
	}
	return out
}

func futuresImpulseTrueRangeRankBucket(rank float64) string {
	switch {
	case rank >= 0.999:
		return "p99_9_plus"
	case rank >= 0.995:
		return "p99_5_99_9"
	case rank >= 0.99:
		return "p99_99_5"
	default:
		return "below_p99"
	}
}

func futuresImpulseVolumeRankBucket(rank float64) string {
	switch {
	case rank >= 0.999:
		return "p99_9_plus"
	case rank >= 0.995:
		return "p99_5_99_9"
	case rank >= 0.99:
		return "p99_99_5"
	case rank >= 0.95:
		return "p95_99"
	default:
		return "below_p95"
	}
}

func lessFuturesImpulseAbsorptionCandidate(a, b FuturesImpulseAbsorptionCandidateRow) bool {
	if a.EventIndex != b.EventIndex {
		return a.EventIndex < b.EventIndex
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	return futuresImpulseDirectionSortKey(a.Direction) < futuresImpulseDirectionSortKey(b.Direction)
}

func lessFuturesImpulseAbsorptionSummary(a, b FuturesImpulseAbsorptionSummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if futuresImpulseDirectionSortKey(a.Direction) != futuresImpulseDirectionSortKey(b.Direction) {
		return futuresImpulseDirectionSortKey(a.Direction) < futuresImpulseDirectionSortKey(b.Direction)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if futuresImpulseRankBucketSortKey(a.TrueRangePercentileBucket) != futuresImpulseRankBucketSortKey(b.TrueRangePercentileBucket) {
		return futuresImpulseRankBucketSortKey(a.TrueRangePercentileBucket) < futuresImpulseRankBucketSortKey(b.TrueRangePercentileBucket)
	}
	return futuresImpulseRankBucketSortKey(a.VolumePercentileBucket) < futuresImpulseRankBucketSortKey(b.VolumePercentileBucket)
}

func lessFuturesImpulseAbsorptionStability(a, b FuturesImpulseAbsorptionStabilityRow) bool {
	if futuresImpulseDirectionSortKey(a.Direction) != futuresImpulseDirectionSortKey(b.Direction) {
		return futuresImpulseDirectionSortKey(a.Direction) < futuresImpulseDirectionSortKey(b.Direction)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if futuresImpulseRankBucketSortKey(a.TrueRangePercentileBucket) != futuresImpulseRankBucketSortKey(b.TrueRangePercentileBucket) {
		return futuresImpulseRankBucketSortKey(a.TrueRangePercentileBucket) < futuresImpulseRankBucketSortKey(b.TrueRangePercentileBucket)
	}
	return futuresImpulseRankBucketSortKey(a.VolumePercentileBucket) < futuresImpulseRankBucketSortKey(b.VolumePercentileBucket)
}

func futuresImpulseDirectionSortKey(direction string) int {
	switch direction {
	case FuturesImpulseAbsorptionDirectionAll:
		return 0
	case FuturesImpulseAbsorptionDirectionDown:
		return 1
	case FuturesImpulseAbsorptionDirectionUp:
		return 2
	default:
		return 99
	}
}

func futuresImpulseRankBucketSortKey(bucket string) int {
	switch bucket {
	case FuturesImpulseAbsorptionBucketAll:
		return 0
	case "p95_99":
		return 1
	case "p99_99_5":
		return 2
	case "p99_5_99_9":
		return 3
	case "p99_9_plus":
		return 4
	default:
		return 99
	}
}
