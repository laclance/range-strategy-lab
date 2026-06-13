package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	SRBoundarySideSupport    = "support"
	SRBoundarySideResistance = "resistance"
)

type SRBoundaryAuditConfig struct {
	HorizonsBars       []int
	DetectorActiveOnly bool
}

type SRBoundaryEventRow struct {
	Index                       int     `json:"index"`
	OpenTime                    string  `json:"open_time"`
	CloseTime                   string  `json:"close_time"`
	Split                       string  `json:"split"`
	Side                        string  `json:"side"`
	Close                       float64 `json:"close"`
	BoundaryPrice               float64 `json:"boundary_price"`
	ZoneTop                     float64 `json:"zone_top"`
	ZoneBottom                  float64 `json:"zone_bottom"`
	ZoneWidth                   float64 `json:"zone_width"`
	RejectionThreshold          float64 `json:"rejection_threshold"`
	DistancePct                 float64 `json:"distance_pct"`
	Strength                    int     `json:"strength"`
	StrengthBucket              string  `json:"strength_bucket"`
	Score                       float64 `json:"score"`
	DetectorProfileID           string  `json:"detector_profile_id"`
	DetectorRawActive           bool    `json:"detector_raw_active"`
	DetectorActive              bool    `json:"detector_active"`
	HorizonBars                 int     `json:"horizon_bars"`
	FutureMaxHigh               float64 `json:"future_max_high"`
	FutureMinLow                float64 `json:"future_min_low"`
	FutureClose                 float64 `json:"future_close"`
	FavorableMove               float64 `json:"favorable_move"`
	AdverseMove                 float64 `json:"adverse_move"`
	FavorableMovePct            float64 `json:"favorable_move_pct"`
	AdverseMovePct              float64 `json:"adverse_move_pct"`
	DistanceBucket              string  `json:"distance_bucket"`
	WickBreak                   bool    `json:"wick_break"`
	CloseBreak                  bool    `json:"close_break"`
	ReclaimedAfterBreak         bool    `json:"reclaimed_after_break"`
	Rejected                    bool    `json:"rejected"`
	FavorableGreaterThanAdverse bool    `json:"favorable_greater_than_adverse"`
}

type SRBoundaryQualityRow struct {
	Split                           string  `json:"split"`
	Side                            string  `json:"side"`
	HorizonBars                     int     `json:"horizon_bars"`
	StrengthBucket                  string  `json:"strength_bucket"`
	DistanceBucket                  string  `json:"distance_bucket"`
	EventCount                      int     `json:"event_count"`
	AvgScore                        float64 `json:"avg_score"`
	AvgDistancePct                  float64 `json:"avg_distance_pct"`
	AvgFavorablePct                 float64 `json:"avg_favorable_pct"`
	MedianFavorablePct              float64 `json:"median_favorable_pct"`
	AvgAdversePct                   float64 `json:"avg_adverse_pct"`
	MedianAdversePct                float64 `json:"median_adverse_pct"`
	CloseBreakRate                  float64 `json:"close_break_rate"`
	WickBreakRate                   float64 `json:"wick_break_rate"`
	ReclaimAfterBreakRate           float64 `json:"reclaim_after_break_rate"`
	RejectionRate                   float64 `json:"rejection_rate"`
	FavorableGreaterThanAdverseRate float64 `json:"favorable_greater_than_adverse_rate"`
}

func DefaultSRBoundaryAuditConfig() SRBoundaryAuditConfig {
	return SRBoundaryAuditConfig{
		HorizonsBars:       []int{1, 3, 6, 12},
		DetectorActiveOnly: true,
	}
}

func RunSRBoundaryAudit(candles []Candle, srRows []SRAuditRow, cfg SRBoundaryAuditConfig) ([]SRBoundaryEventRow, []SRBoundaryQualityRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	var events []SRBoundaryEventRow
	for _, row := range srRows {
		if cfg.DetectorActiveOnly && !row.DetectorActive {
			continue
		}
		if row.HasSupport && row.NearSupport {
			events = appendSRBoundaryEvents(events, candles, row, SRBoundarySideSupport, cfg.HorizonsBars)
		}
		if row.HasResistance && row.NearResistance {
			events = appendSRBoundaryEvents(events, candles, row, SRBoundarySideResistance, cfg.HorizonsBars)
		}
	}

	return events, SummarizeSRBoundaryQuality(events), nil
}

func (cfg SRBoundaryAuditConfig) withDefaults() SRBoundaryAuditConfig {
	defaults := DefaultSRBoundaryAuditConfig()
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	return cfg
}

func (cfg SRBoundaryAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("SR boundary audit horizon bars must be positive")
		}
	}
	return nil
}

func appendSRBoundaryEvents(events []SRBoundaryEventRow, candles []Candle, row SRAuditRow, side string, horizons []int) []SRBoundaryEventRow {
	for _, horizon := range horizons {
		event, ok := newSRBoundaryEvent(candles, row, side, horizon)
		if ok {
			events = append(events, event)
		}
	}
	return events
}

func newSRBoundaryEvent(candles []Candle, row SRAuditRow, side string, horizon int) (SRBoundaryEventRow, bool) {
	if row.Index < 0 || row.Index+horizon >= len(candles) {
		return SRBoundaryEventRow{}, false
	}

	future := candles[row.Index+1 : row.Index+horizon+1]
	futureMaxHigh, futureMinLow := futureHighLow(future)
	futureClose := future[len(future)-1].Close

	event := SRBoundaryEventRow{
		Index:             row.Index,
		OpenTime:          row.OpenTime,
		CloseTime:         row.CloseTime,
		Split:             row.Split,
		Side:              side,
		Close:             row.Close,
		DetectorProfileID: row.DetectorProfileID,
		DetectorRawActive: row.DetectorRawActive,
		DetectorActive:    row.DetectorActive,
		HorizonBars:       horizon,
		FutureMaxHigh:     futureMaxHigh,
		FutureMinLow:      futureMinLow,
		FutureClose:       futureClose,
	}

	switch side {
	case SRBoundarySideSupport:
		event.BoundaryPrice = row.NearestSupport
		event.ZoneTop = row.NearestSupportTop
		event.ZoneBottom = row.NearestSupportBottom
		event.DistancePct = row.NearestSupportDistancePct
		event.Strength = row.NearestSupportStrength
		event.Score = row.NearestSupportScore
		event.FavorableMove = futureMaxHigh - row.Close
		event.AdverseMove = row.Close - futureMinLow
		event.WickBreak = futureMinLow < row.NearestSupportBottom
		event.CloseBreak, event.ReclaimedAfterBreak = supportCloseBreakAndReclaim(future, row.NearestSupportBottom, row.NearestSupport)
	case SRBoundarySideResistance:
		event.BoundaryPrice = row.NearestResistance
		event.ZoneTop = row.NearestResistanceTop
		event.ZoneBottom = row.NearestResistanceBottom
		event.DistancePct = row.NearestResistanceDistancePct
		event.Strength = row.NearestResistanceStrength
		event.Score = row.NearestResistanceScore
		event.FavorableMove = row.Close - futureMinLow
		event.AdverseMove = futureMaxHigh - row.Close
		event.WickBreak = futureMaxHigh > row.NearestResistanceTop
		event.CloseBreak, event.ReclaimedAfterBreak = resistanceCloseBreakAndReclaim(future, row.NearestResistanceTop, row.NearestResistance)
	default:
		return SRBoundaryEventRow{}, false
	}

	event.ZoneWidth = math.Abs(event.ZoneTop - event.ZoneBottom)
	event.RejectionThreshold = rejectionThreshold(event.ZoneWidth, row.Close)
	event.FavorableMovePct = movePct(event.FavorableMove, row.Close)
	event.AdverseMovePct = movePct(event.AdverseMove, row.Close)
	event.StrengthBucket = strengthBucket(event.Strength)
	event.DistanceBucket = distanceBucket(event.DistancePct)
	event.Rejected = !event.CloseBreak && event.FavorableMove >= event.RejectionThreshold
	event.FavorableGreaterThanAdverse = event.FavorableMove > event.AdverseMove
	return event, true
}

func futureHighLow(candles []Candle) (float64, float64) {
	maxHigh := candles[0].High
	minLow := candles[0].Low
	for _, candle := range candles[1:] {
		if candle.High > maxHigh {
			maxHigh = candle.High
		}
		if candle.Low < minLow {
			minLow = candle.Low
		}
	}
	return maxHigh, minLow
}

func supportCloseBreakAndReclaim(future []Candle, zoneBottom, supportPrice float64) (bool, bool) {
	for i, candle := range future {
		if candle.Close >= zoneBottom {
			continue
		}
		for _, later := range future[i+1:] {
			if later.Close >= supportPrice {
				return true, true
			}
		}
		return true, false
	}
	return false, false
}

func resistanceCloseBreakAndReclaim(future []Candle, zoneTop, resistancePrice float64) (bool, bool) {
	for i, candle := range future {
		if candle.Close <= zoneTop {
			continue
		}
		for _, later := range future[i+1:] {
			if later.Close <= resistancePrice {
				return true, true
			}
		}
		return true, false
	}
	return false, false
}

func rejectionThreshold(zoneWidth, close float64) float64 {
	if zoneWidth > 0 {
		return zoneWidth
	}
	if close <= 0 {
		return 0
	}
	return close * 0.001
}

func movePct(move, close float64) float64 {
	if close <= 0 {
		return 0
	}
	return move / close
}

func strengthBucket(strength int) string {
	if strength >= 4 {
		return "4plus"
	}
	if strength == 3 {
		return "3"
	}
	return "2"
}

func distanceBucket(distancePct float64) string {
	switch {
	case distancePct <= 0.0005:
		return "0_5bp"
	case distancePct <= 0.0010:
		return "5_10bp"
	case distancePct <= 0.0020:
		return "10_20bp"
	default:
		return "gt_20bp"
	}
}

type srBoundaryQualityKey struct {
	split          string
	side           string
	horizonBars    int
	strengthBucket string
	distanceBucket string
}

type srBoundaryQualityAccumulator struct {
	key                          srBoundaryQualityKey
	events                       int
	scoreSum                     float64
	distancePctSum               float64
	favorablePctSum              float64
	adversePctSum                float64
	favorablePcts                []float64
	adversePcts                  []float64
	closeBreaks                  int
	wickBreaks                   int
	reclaimsAfterBreak           int
	rejections                   int
	favorableGreaterThanAdverses int
}

func SummarizeSRBoundaryQuality(events []SRBoundaryEventRow) []SRBoundaryQualityRow {
	accumulators := map[srBoundaryQualityKey]*srBoundaryQualityAccumulator{}
	for _, event := range events {
		key := srBoundaryQualityKey{
			split:          event.Split,
			side:           event.Side,
			horizonBars:    event.HorizonBars,
			strengthBucket: event.StrengthBucket,
			distanceBucket: event.DistanceBucket,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srBoundaryQualityAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRBoundaryQualityRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRBoundaryQualityRow(rows[i], rows[j])
	})
	return rows
}

func (acc *srBoundaryQualityAccumulator) add(event SRBoundaryEventRow) {
	acc.events++
	acc.scoreSum += event.Score
	acc.distancePctSum += event.DistancePct
	acc.favorablePctSum += event.FavorableMovePct
	acc.adversePctSum += event.AdverseMovePct
	acc.favorablePcts = append(acc.favorablePcts, event.FavorableMovePct)
	acc.adversePcts = append(acc.adversePcts, event.AdverseMovePct)
	if event.CloseBreak {
		acc.closeBreaks++
	}
	if event.WickBreak {
		acc.wickBreaks++
	}
	if event.ReclaimedAfterBreak {
		acc.reclaimsAfterBreak++
	}
	if event.Rejected {
		acc.rejections++
	}
	if event.FavorableGreaterThanAdverse {
		acc.favorableGreaterThanAdverses++
	}
}

func (acc srBoundaryQualityAccumulator) row() SRBoundaryQualityRow {
	row := SRBoundaryQualityRow{
		Split:              acc.key.split,
		Side:               acc.key.side,
		HorizonBars:        acc.key.horizonBars,
		StrengthBucket:     acc.key.strengthBucket,
		DistanceBucket:     acc.key.distanceBucket,
		EventCount:         acc.events,
		MedianFavorablePct: medianFloat(acc.favorablePcts),
		MedianAdversePct:   medianFloat(acc.adversePcts),
	}
	if acc.events > 0 {
		row.AvgScore = acc.scoreSum / float64(acc.events)
		row.AvgDistancePct = acc.distancePctSum / float64(acc.events)
		row.AvgFavorablePct = acc.favorablePctSum / float64(acc.events)
		row.AvgAdversePct = acc.adversePctSum / float64(acc.events)
		row.CloseBreakRate = float64(acc.closeBreaks) / float64(acc.events)
		row.WickBreakRate = float64(acc.wickBreaks) / float64(acc.events)
		row.RejectionRate = float64(acc.rejections) / float64(acc.events)
		row.FavorableGreaterThanAdverseRate = float64(acc.favorableGreaterThanAdverses) / float64(acc.events)
	}
	if acc.closeBreaks > 0 {
		row.ReclaimAfterBreakRate = float64(acc.reclaimsAfterBreak) / float64(acc.closeBreaks)
	}
	return row
}

func medianFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sortedValues := append([]float64(nil), values...)
	sort.Float64s(sortedValues)
	mid := len(sortedValues) / 2
	if len(sortedValues)%2 == 1 {
		return sortedValues[mid]
	}
	return (sortedValues[mid-1] + sortedValues[mid]) / 2
}

func lessSRBoundaryQualityRow(a, b SRBoundaryQualityRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if strengthBucketSortKey(a.StrengthBucket) != strengthBucketSortKey(b.StrengthBucket) {
		return strengthBucketSortKey(a.StrengthBucket) < strengthBucketSortKey(b.StrengthBucket)
	}
	return distanceBucketSortKey(a.DistanceBucket) < distanceBucketSortKey(b.DistanceBucket)
}

func splitSortKey(split string) int {
	switch split {
	case "2021_2022_stress":
		return 0
	case "2023_2024_oos":
		return 1
	case "2025_2026_recent":
		return 2
	case "full_2021_2026":
		return 3
	default:
		return 99
	}
}

func sideSortKey(side string) int {
	switch side {
	case SRBoundarySideSupport:
		return 0
	case SRBoundarySideResistance:
		return 1
	default:
		return 99
	}
}

func strengthBucketSortKey(bucket string) int {
	switch bucket {
	case "2":
		return 0
	case "3":
		return 1
	case "4plus":
		return 2
	default:
		return 99
	}
}

func distanceBucketSortKey(bucket string) int {
	switch bucket {
	case "0_5bp":
		return 0
	case "5_10bp":
		return 1
	case "10_20bp":
		return 2
	case "gt_20bp":
		return 3
	default:
		return 99
	}
}
