package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	srCloseLocationBelowZone        = "below_zone"
	srCloseLocationInZoneBelowLevel = "in_zone_below_boundary"
	srCloseLocationInZoneAboveLevel = "in_zone_above_boundary"
	srCloseLocationAboveZone        = "above_zone"
)

type SRRejectionTimingAuditConfig struct {
	HorizonsBars       []int
	DetectorActiveOnly bool
}

type SRRejectionTimingCandidateRow struct {
	Split                                 string  `json:"split"`
	Side                                  string  `json:"side"`
	HorizonBars                           int     `json:"horizon_bars"`
	CloseLocation                         string  `json:"close_location"`
	TouchedZone                           bool    `json:"touched_zone"`
	PiercedZone                           bool    `json:"pierced_zone"`
	ClosedBack                            bool    `json:"closed_back"`
	DecisionRejectionCandidate            bool    `json:"decision_rejection_candidate"`
	WickBeyondBucket                      string  `json:"wick_beyond_bucket"`
	StrengthBucket                        string  `json:"strength_bucket"`
	DistanceBucket                        string  `json:"distance_bucket"`
	DetectorProfileID                     string  `json:"detector_profile_id"`
	DetectorRawActive                     bool    `json:"detector_raw_active"`
	DetectorActive                        bool    `json:"detector_active"`
	CandidateCount                        int     `json:"candidate_count"`
	AvgScore                              float64 `json:"avg_score"`
	AvgDistancePct                        float64 `json:"avg_distance_pct"`
	AvgWickBeyondPct                      float64 `json:"avg_wick_beyond_pct"`
	LabelCloseBreakCount                  int     `json:"label_close_break_count"`
	LabelWickBreakCount                   int     `json:"label_wick_break_count"`
	LabelReclaimedAfterBreakCount         int     `json:"label_reclaimed_after_break_count"`
	LabelRejectedCount                    int     `json:"label_rejected_count"`
	LabelFavorableGreaterThanAdverseCount int     `json:"label_favorable_greater_than_adverse_count"`
	LabelCloseBreakRate                   float64 `json:"label_close_break_rate"`
	LabelWickBreakRate                    float64 `json:"label_wick_break_rate"`
	LabelReclaimEventRate                 float64 `json:"label_reclaim_event_rate"`
	LabelReclaimGivenCloseBreakRate       float64 `json:"label_reclaim_given_close_break_rate"`
	LabelRejectionRate                    float64 `json:"label_rejection_rate"`
	LabelAvgFavorablePct                  float64 `json:"label_avg_favorable_pct"`
	LabelAvgAdversePct                    float64 `json:"label_avg_adverse_pct"`
	LabelFavorableMinusAdversePct         float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGreaterThanAdverseRate  float64 `json:"label_favorable_greater_than_adverse_rate"`
}

type SRRejectionTimingSummaryRow struct {
	Split                                                 string  `json:"split"`
	Side                                                  string  `json:"side"`
	HorizonBars                                           int     `json:"horizon_bars"`
	DetectorProfileID                                     string  `json:"detector_profile_id"`
	DetectorRawActive                                     bool    `json:"detector_raw_active"`
	DetectorActive                                        bool    `json:"detector_active"`
	CandidateCount                                        int     `json:"candidate_count"`
	TouchedCount                                          int     `json:"touched_count"`
	PiercedCount                                          int     `json:"pierced_count"`
	ClosedBackCount                                       int     `json:"closed_back_count"`
	DecisionRejectionCandidateCount                       int     `json:"decision_rejection_candidate_count"`
	TouchedRate                                           float64 `json:"touched_rate"`
	PiercedRate                                           float64 `json:"pierced_rate"`
	ClosedBackRate                                        float64 `json:"closed_back_rate"`
	DecisionRejectionCandidateRate                        float64 `json:"decision_rejection_candidate_rate"`
	LabelCloseBreakRate                                   float64 `json:"label_close_break_rate"`
	LabelWickBreakRate                                    float64 `json:"label_wick_break_rate"`
	LabelReclaimEventRate                                 float64 `json:"label_reclaim_event_rate"`
	LabelReclaimGivenCloseBreakRate                       float64 `json:"label_reclaim_given_close_break_rate"`
	LabelRejectionRate                                    float64 `json:"label_rejection_rate"`
	LabelAvgFavorablePct                                  float64 `json:"label_avg_favorable_pct"`
	LabelAvgAdversePct                                    float64 `json:"label_avg_adverse_pct"`
	LabelFavorableMinusAdversePct                         float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGreaterThanAdverseRate                  float64 `json:"label_favorable_greater_than_adverse_rate"`
	DecisionCandidateLabelCloseBreakRate                  float64 `json:"decision_candidate_label_close_break_rate"`
	DecisionCandidateLabelWickBreakRate                   float64 `json:"decision_candidate_label_wick_break_rate"`
	DecisionCandidateLabelReclaimEventRate                float64 `json:"decision_candidate_label_reclaim_event_rate"`
	DecisionCandidateLabelReclaimGivenCloseBreakRate      float64 `json:"decision_candidate_label_reclaim_given_close_break_rate"`
	DecisionCandidateLabelRejectionRate                   float64 `json:"decision_candidate_label_rejection_rate"`
	DecisionCandidateLabelAvgFavorablePct                 float64 `json:"decision_candidate_label_avg_favorable_pct"`
	DecisionCandidateLabelAvgAdversePct                   float64 `json:"decision_candidate_label_avg_adverse_pct"`
	DecisionCandidateLabelFavorableMinusAdversePct        float64 `json:"decision_candidate_label_favorable_minus_adverse_pct"`
	DecisionCandidateLabelFavorableGreaterThanAdverseRate float64 `json:"decision_candidate_label_favorable_greater_than_adverse_rate"`
}

func DefaultSRRejectionTimingAuditConfig() SRRejectionTimingAuditConfig {
	boundaryDefaults := DefaultSRBoundaryAuditConfig()
	return SRRejectionTimingAuditConfig{
		HorizonsBars:       append([]int(nil), boundaryDefaults.HorizonsBars...),
		DetectorActiveOnly: boundaryDefaults.DetectorActiveOnly,
	}
}

func RunSRRejectionTimingAudit(candles []Candle, srRows []SRAuditRow, cfg SRRejectionTimingAuditConfig) ([]SRRejectionTimingCandidateRow, []SRRejectionTimingSummaryRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	var events []srRejectionTimingEvent
	for _, row := range srRows {
		if cfg.DetectorActiveOnly && !row.DetectorActive {
			continue
		}
		if row.HasSupport && row.NearSupport {
			events = appendSRRejectionTimingEvents(events, candles, row, SRBoundarySideSupport, cfg.HorizonsBars)
		}
		if row.HasResistance && row.NearResistance {
			events = appendSRRejectionTimingEvents(events, candles, row, SRBoundarySideResistance, cfg.HorizonsBars)
		}
	}

	return summarizeSRRejectionTimingCandidates(events), summarizeSRRejectionTimingSummary(events), nil
}

func (cfg SRRejectionTimingAuditConfig) withDefaults() SRRejectionTimingAuditConfig {
	defaults := DefaultSRRejectionTimingAuditConfig()
	if len(cfg.HorizonsBars) == 0 && !cfg.DetectorActiveOnly {
		return defaults
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	return cfg
}

func (cfg SRRejectionTimingAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("SR rejection timing audit horizon bars must be positive")
		}
	}
	return nil
}

type srRejectionTimingDecision struct {
	Split                      string
	Side                       string
	CloseLocation              string
	TouchedZone                bool
	PiercedZone                bool
	ClosedBack                 bool
	DecisionRejectionCandidate bool
	WickBeyondPct              float64
	WickBeyondBucket           string
	StrengthBucket             string
	DistanceBucket             string
	Score                      float64
	DistancePct                float64
	DetectorProfileID          string
	DetectorRawActive          bool
	DetectorActive             bool
}

type srRejectionTimingEvent struct {
	srRejectionTimingDecision
	HorizonBars                      int
	LabelCloseBreak                  bool
	LabelWickBreak                   bool
	LabelReclaimedAfterBreak         bool
	LabelRejected                    bool
	LabelFavorableGreaterThanAdverse bool
	LabelFavorableMovePct            float64
	LabelAdverseMovePct              float64
}

func appendSRRejectionTimingEvents(events []srRejectionTimingEvent, candles []Candle, row SRAuditRow, side string, horizons []int) []srRejectionTimingEvent {
	decision, ok := newSRRejectionTimingDecision(candles, row, side)
	if !ok {
		return events
	}
	for _, horizon := range horizons {
		label, ok := newSRBoundaryEvent(candles, row, side, horizon)
		if !ok {
			continue
		}
		events = append(events, srRejectionTimingEvent{
			srRejectionTimingDecision:        decision,
			HorizonBars:                      horizon,
			LabelCloseBreak:                  label.CloseBreak,
			LabelWickBreak:                   label.WickBreak,
			LabelReclaimedAfterBreak:         label.ReclaimedAfterBreak,
			LabelRejected:                    label.Rejected,
			LabelFavorableGreaterThanAdverse: label.FavorableGreaterThanAdverse,
			LabelFavorableMovePct:            label.FavorableMovePct,
			LabelAdverseMovePct:              label.AdverseMovePct,
		})
	}
	return events
}

func newSRRejectionTimingDecision(candles []Candle, row SRAuditRow, side string) (srRejectionTimingDecision, bool) {
	if row.Index < 0 || row.Index >= len(candles) {
		return srRejectionTimingDecision{}, false
	}
	candle := candles[row.Index]
	decision := srRejectionTimingDecision{
		Split:             row.Split,
		Side:              side,
		DetectorProfileID: row.DetectorProfileID,
		DetectorRawActive: row.DetectorRawActive,
		DetectorActive:    row.DetectorActive,
	}

	switch side {
	case SRBoundarySideSupport:
		decision.Score = row.NearestSupportScore
		decision.DistancePct = row.NearestSupportDistancePct
		decision.StrengthBucket = strengthBucket(row.NearestSupportStrength)
		decision.DistanceBucket = distanceBucket(row.NearestSupportDistancePct)
		decision.TouchedZone = candle.Low <= row.NearestSupportTop
		decision.PiercedZone = candle.Low < row.NearestSupportBottom
		decision.ClosedBack = candle.Close >= row.NearestSupport
		decision.CloseLocation = closeLocationForBoundary(side, candle.Close, row.NearestSupport, row.NearestSupportTop, row.NearestSupportBottom)
		decision.WickBeyondPct = movePct(math.Max(0, row.NearestSupportBottom-candle.Low), candle.Close)
	case SRBoundarySideResistance:
		decision.Score = row.NearestResistanceScore
		decision.DistancePct = row.NearestResistanceDistancePct
		decision.StrengthBucket = strengthBucket(row.NearestResistanceStrength)
		decision.DistanceBucket = distanceBucket(row.NearestResistanceDistancePct)
		decision.TouchedZone = candle.High >= row.NearestResistanceBottom
		decision.PiercedZone = candle.High > row.NearestResistanceTop
		decision.ClosedBack = candle.Close <= row.NearestResistance
		decision.CloseLocation = closeLocationForBoundary(side, candle.Close, row.NearestResistance, row.NearestResistanceTop, row.NearestResistanceBottom)
		decision.WickBeyondPct = movePct(math.Max(0, candle.High-row.NearestResistanceTop), candle.Close)
	default:
		return srRejectionTimingDecision{}, false
	}

	decision.DecisionRejectionCandidate = decision.TouchedZone && decision.ClosedBack
	decision.WickBeyondBucket = distanceBucket(decision.WickBeyondPct)
	return decision, true
}

func closeLocationForBoundary(side string, close, boundary, zoneTop, zoneBottom float64) string {
	switch side {
	case SRBoundarySideSupport:
		switch {
		case close < zoneBottom:
			return srCloseLocationBelowZone
		case close < boundary:
			return srCloseLocationInZoneBelowLevel
		case close <= zoneTop:
			return srCloseLocationInZoneAboveLevel
		default:
			return srCloseLocationAboveZone
		}
	case SRBoundarySideResistance:
		switch {
		case close > zoneTop:
			return srCloseLocationAboveZone
		case close > boundary:
			return srCloseLocationInZoneAboveLevel
		case close >= zoneBottom:
			return srCloseLocationInZoneBelowLevel
		default:
			return srCloseLocationBelowZone
		}
	default:
		return ""
	}
}

type srRejectionTimingCandidateKey struct {
	split                      string
	side                       string
	horizonBars                int
	closeLocation              string
	touchedZone                bool
	piercedZone                bool
	closedBack                 bool
	decisionRejectionCandidate bool
	wickBeyondBucket           string
	strengthBucket             string
	distanceBucket             string
	detectorProfileID          string
	detectorRawActive          bool
	detectorActive             bool
}

type srRejectionTimingSummaryKey struct {
	split             string
	side              string
	horizonBars       int
	detectorProfileID string
	detectorRawActive bool
	detectorActive    bool
}

type srRejectionTimingCandidateAccumulator struct {
	key              srRejectionTimingCandidateKey
	candidates       int
	scoreSum         float64
	distancePctSum   float64
	wickBeyondPctSum float64
	labels           srRejectionTimingLabelAccumulator
}

type srRejectionTimingSummaryAccumulator struct {
	key                         srRejectionTimingSummaryKey
	candidates                  int
	touched                     int
	pierced                     int
	closedBack                  int
	decisionRejectionCandidates int
	labels                      srRejectionTimingLabelAccumulator
	decisionCandidateLabels     srRejectionTimingLabelAccumulator
}

type srRejectionTimingLabelAccumulator struct {
	events                       int
	closeBreaks                  int
	wickBreaks                   int
	reclaimsAfterBreak           int
	rejections                   int
	favorableGreaterThanAdverses int
	favorablePctSum              float64
	adversePctSum                float64
}

func summarizeSRRejectionTimingCandidates(events []srRejectionTimingEvent) []SRRejectionTimingCandidateRow {
	accumulators := map[srRejectionTimingCandidateKey]*srRejectionTimingCandidateAccumulator{}
	for _, event := range events {
		key := srRejectionTimingCandidateKey{
			split:                      event.Split,
			side:                       event.Side,
			horizonBars:                event.HorizonBars,
			closeLocation:              event.CloseLocation,
			touchedZone:                event.TouchedZone,
			piercedZone:                event.PiercedZone,
			closedBack:                 event.ClosedBack,
			decisionRejectionCandidate: event.DecisionRejectionCandidate,
			wickBeyondBucket:           event.WickBeyondBucket,
			strengthBucket:             event.StrengthBucket,
			distanceBucket:             event.DistanceBucket,
			detectorProfileID:          event.DetectorProfileID,
			detectorRawActive:          event.DetectorRawActive,
			detectorActive:             event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srRejectionTimingCandidateAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRRejectionTimingCandidateRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRRejectionTimingCandidateRow(rows[i], rows[j])
	})
	return rows
}

func summarizeSRRejectionTimingSummary(events []srRejectionTimingEvent) []SRRejectionTimingSummaryRow {
	accumulators := map[srRejectionTimingSummaryKey]*srRejectionTimingSummaryAccumulator{}
	for _, event := range events {
		key := srRejectionTimingSummaryKey{
			split:             event.Split,
			side:              event.Side,
			horizonBars:       event.HorizonBars,
			detectorProfileID: event.DetectorProfileID,
			detectorRawActive: event.DetectorRawActive,
			detectorActive:    event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srRejectionTimingSummaryAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRRejectionTimingSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRRejectionTimingSummaryRow(rows[i], rows[j])
	})
	return rows
}

func (acc *srRejectionTimingCandidateAccumulator) add(event srRejectionTimingEvent) {
	acc.candidates++
	acc.scoreSum += event.Score
	acc.distancePctSum += event.DistancePct
	acc.wickBeyondPctSum += event.WickBeyondPct
	acc.labels.add(event)
}

func (acc srRejectionTimingCandidateAccumulator) row() SRRejectionTimingCandidateRow {
	row := SRRejectionTimingCandidateRow{
		Split:                      acc.key.split,
		Side:                       acc.key.side,
		HorizonBars:                acc.key.horizonBars,
		CloseLocation:              acc.key.closeLocation,
		TouchedZone:                acc.key.touchedZone,
		PiercedZone:                acc.key.piercedZone,
		ClosedBack:                 acc.key.closedBack,
		DecisionRejectionCandidate: acc.key.decisionRejectionCandidate,
		WickBeyondBucket:           acc.key.wickBeyondBucket,
		StrengthBucket:             acc.key.strengthBucket,
		DistanceBucket:             acc.key.distanceBucket,
		DetectorProfileID:          acc.key.detectorProfileID,
		DetectorRawActive:          acc.key.detectorRawActive,
		DetectorActive:             acc.key.detectorActive,
		CandidateCount:             acc.candidates,
	}
	if acc.candidates > 0 {
		row.AvgScore = acc.scoreSum / float64(acc.candidates)
		row.AvgDistancePct = acc.distancePctSum / float64(acc.candidates)
		row.AvgWickBeyondPct = acc.wickBeyondPctSum / float64(acc.candidates)
	}
	acc.labels.addToCandidateRow(&row)
	return row
}

func (acc *srRejectionTimingSummaryAccumulator) add(event srRejectionTimingEvent) {
	acc.candidates++
	if event.TouchedZone {
		acc.touched++
	}
	if event.PiercedZone {
		acc.pierced++
	}
	if event.ClosedBack {
		acc.closedBack++
	}
	if event.DecisionRejectionCandidate {
		acc.decisionRejectionCandidates++
		acc.decisionCandidateLabels.add(event)
	}
	acc.labels.add(event)
}

func (acc srRejectionTimingSummaryAccumulator) row() SRRejectionTimingSummaryRow {
	row := SRRejectionTimingSummaryRow{
		Split:                           acc.key.split,
		Side:                            acc.key.side,
		HorizonBars:                     acc.key.horizonBars,
		DetectorProfileID:               acc.key.detectorProfileID,
		DetectorRawActive:               acc.key.detectorRawActive,
		DetectorActive:                  acc.key.detectorActive,
		CandidateCount:                  acc.candidates,
		TouchedCount:                    acc.touched,
		PiercedCount:                    acc.pierced,
		ClosedBackCount:                 acc.closedBack,
		DecisionRejectionCandidateCount: acc.decisionRejectionCandidates,
	}
	if acc.candidates > 0 {
		row.TouchedRate = float64(acc.touched) / float64(acc.candidates)
		row.PiercedRate = float64(acc.pierced) / float64(acc.candidates)
		row.ClosedBackRate = float64(acc.closedBack) / float64(acc.candidates)
		row.DecisionRejectionCandidateRate = float64(acc.decisionRejectionCandidates) / float64(acc.candidates)
	}
	acc.labels.addToSummaryRow(&row)
	acc.decisionCandidateLabels.addToDecisionCandidateSummaryRow(&row)
	return row
}

func (acc *srRejectionTimingLabelAccumulator) add(event srRejectionTimingEvent) {
	acc.events++
	if event.LabelCloseBreak {
		acc.closeBreaks++
	}
	if event.LabelWickBreak {
		acc.wickBreaks++
	}
	if event.LabelReclaimedAfterBreak {
		acc.reclaimsAfterBreak++
	}
	if event.LabelRejected {
		acc.rejections++
	}
	if event.LabelFavorableGreaterThanAdverse {
		acc.favorableGreaterThanAdverses++
	}
	acc.favorablePctSum += event.LabelFavorableMovePct
	acc.adversePctSum += event.LabelAdverseMovePct
}

func (acc srRejectionTimingLabelAccumulator) avgFavorablePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.favorablePctSum / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) avgAdversePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.adversePctSum / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) favorableGreaterThanAdverseRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.favorableGreaterThanAdverses) / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) closeBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.closeBreaks) / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) wickBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.wickBreaks) / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) reclaimEventRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) reclaimGivenCloseBreakRate() float64 {
	if acc.closeBreaks == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.closeBreaks)
}

func (acc srRejectionTimingLabelAccumulator) rejectionRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.rejections) / float64(acc.events)
}

func (acc srRejectionTimingLabelAccumulator) addToCandidateRow(row *SRRejectionTimingCandidateRow) {
	row.LabelCloseBreakCount = acc.closeBreaks
	row.LabelWickBreakCount = acc.wickBreaks
	row.LabelReclaimedAfterBreakCount = acc.reclaimsAfterBreak
	row.LabelRejectedCount = acc.rejections
	row.LabelFavorableGreaterThanAdverseCount = acc.favorableGreaterThanAdverses
	row.LabelCloseBreakRate = acc.closeBreakRate()
	row.LabelWickBreakRate = acc.wickBreakRate()
	row.LabelReclaimEventRate = acc.reclaimEventRate()
	row.LabelReclaimGivenCloseBreakRate = acc.reclaimGivenCloseBreakRate()
	row.LabelRejectionRate = acc.rejectionRate()
	row.LabelAvgFavorablePct = acc.avgFavorablePct()
	row.LabelAvgAdversePct = acc.avgAdversePct()
	row.LabelFavorableMinusAdversePct = row.LabelAvgFavorablePct - row.LabelAvgAdversePct
	row.LabelFavorableGreaterThanAdverseRate = acc.favorableGreaterThanAdverseRate()
}

func (acc srRejectionTimingLabelAccumulator) addToSummaryRow(row *SRRejectionTimingSummaryRow) {
	row.LabelCloseBreakRate = acc.closeBreakRate()
	row.LabelWickBreakRate = acc.wickBreakRate()
	row.LabelReclaimEventRate = acc.reclaimEventRate()
	row.LabelReclaimGivenCloseBreakRate = acc.reclaimGivenCloseBreakRate()
	row.LabelRejectionRate = acc.rejectionRate()
	row.LabelAvgFavorablePct = acc.avgFavorablePct()
	row.LabelAvgAdversePct = acc.avgAdversePct()
	row.LabelFavorableMinusAdversePct = row.LabelAvgFavorablePct - row.LabelAvgAdversePct
	row.LabelFavorableGreaterThanAdverseRate = acc.favorableGreaterThanAdverseRate()
}

func (acc srRejectionTimingLabelAccumulator) addToDecisionCandidateSummaryRow(row *SRRejectionTimingSummaryRow) {
	row.DecisionCandidateLabelCloseBreakRate = acc.closeBreakRate()
	row.DecisionCandidateLabelWickBreakRate = acc.wickBreakRate()
	row.DecisionCandidateLabelReclaimEventRate = acc.reclaimEventRate()
	row.DecisionCandidateLabelReclaimGivenCloseBreakRate = acc.reclaimGivenCloseBreakRate()
	row.DecisionCandidateLabelRejectionRate = acc.rejectionRate()
	row.DecisionCandidateLabelAvgFavorablePct = acc.avgFavorablePct()
	row.DecisionCandidateLabelAvgAdversePct = acc.avgAdversePct()
	row.DecisionCandidateLabelFavorableMinusAdversePct = row.DecisionCandidateLabelAvgFavorablePct - row.DecisionCandidateLabelAvgAdversePct
	row.DecisionCandidateLabelFavorableGreaterThanAdverseRate = acc.favorableGreaterThanAdverseRate()
}

func lessSRRejectionTimingCandidateRow(a, b SRRejectionTimingCandidateRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.DetectorProfileID != b.DetectorProfileID {
		return a.DetectorProfileID < b.DetectorProfileID
	}
	if boolSortKey(a.DetectorRawActive) != boolSortKey(b.DetectorRawActive) {
		return boolSortKey(a.DetectorRawActive) < boolSortKey(b.DetectorRawActive)
	}
	if boolSortKey(a.DetectorActive) != boolSortKey(b.DetectorActive) {
		return boolSortKey(a.DetectorActive) < boolSortKey(b.DetectorActive)
	}
	if closeLocationSortKey(a.CloseLocation) != closeLocationSortKey(b.CloseLocation) {
		return closeLocationSortKey(a.CloseLocation) < closeLocationSortKey(b.CloseLocation)
	}
	if boolSortKey(a.TouchedZone) != boolSortKey(b.TouchedZone) {
		return boolSortKey(a.TouchedZone) < boolSortKey(b.TouchedZone)
	}
	if boolSortKey(a.PiercedZone) != boolSortKey(b.PiercedZone) {
		return boolSortKey(a.PiercedZone) < boolSortKey(b.PiercedZone)
	}
	if boolSortKey(a.ClosedBack) != boolSortKey(b.ClosedBack) {
		return boolSortKey(a.ClosedBack) < boolSortKey(b.ClosedBack)
	}
	if boolSortKey(a.DecisionRejectionCandidate) != boolSortKey(b.DecisionRejectionCandidate) {
		return boolSortKey(a.DecisionRejectionCandidate) < boolSortKey(b.DecisionRejectionCandidate)
	}
	if distanceBucketSortKey(a.WickBeyondBucket) != distanceBucketSortKey(b.WickBeyondBucket) {
		return distanceBucketSortKey(a.WickBeyondBucket) < distanceBucketSortKey(b.WickBeyondBucket)
	}
	if strengthBucketSortKey(a.StrengthBucket) != strengthBucketSortKey(b.StrengthBucket) {
		return strengthBucketSortKey(a.StrengthBucket) < strengthBucketSortKey(b.StrengthBucket)
	}
	return distanceBucketSortKey(a.DistanceBucket) < distanceBucketSortKey(b.DistanceBucket)
}

func lessSRRejectionTimingSummaryRow(a, b SRRejectionTimingSummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.HorizonBars != b.HorizonBars {
		return a.HorizonBars < b.HorizonBars
	}
	if a.DetectorProfileID != b.DetectorProfileID {
		return a.DetectorProfileID < b.DetectorProfileID
	}
	if boolSortKey(a.DetectorRawActive) != boolSortKey(b.DetectorRawActive) {
		return boolSortKey(a.DetectorRawActive) < boolSortKey(b.DetectorRawActive)
	}
	return boolSortKey(a.DetectorActive) < boolSortKey(b.DetectorActive)
}

func closeLocationSortKey(location string) int {
	switch location {
	case srCloseLocationBelowZone:
		return 0
	case srCloseLocationInZoneBelowLevel:
		return 1
	case srCloseLocationInZoneAboveLevel:
		return 2
	case srCloseLocationAboveZone:
		return 3
	default:
		return 99
	}
}

func boolSortKey(value bool) int {
	if value {
		return 1
	}
	return 0
}
