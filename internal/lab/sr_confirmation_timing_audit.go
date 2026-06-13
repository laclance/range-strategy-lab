package lab

import (
	"fmt"
	"sort"
)

type SRConfirmationTimingAuditConfig struct {
	HorizonsBars          []int
	ConfirmationDelayBars []int
	DetectorActiveOnly    bool
}

type SRConfirmationTimingCandidateRow struct {
	Split                                 string  `json:"split"`
	Side                                  string  `json:"side"`
	ConfirmationDelayBars                 int     `json:"confirmation_delay_bars"`
	HorizonBars                           int     `json:"horizon_bars"`
	SeedCloseLocation                     string  `json:"seed_close_location"`
	SeedPiercedZone                       bool    `json:"seed_pierced_zone"`
	SeedWickBeyondBucket                  string  `json:"seed_wick_beyond_bucket"`
	ConfirmationCloseLocation             string  `json:"confirmation_close_location"`
	ConfirmationFavorableClose            bool    `json:"confirmation_favorable_close"`
	ConfirmationWrongSideClose            bool    `json:"confirmation_wrong_side_close"`
	DecisionConfirmationCandidate         bool    `json:"decision_confirmation_candidate"`
	StrengthBucket                        string  `json:"strength_bucket"`
	DistanceBucket                        string  `json:"distance_bucket"`
	DetectorProfileID                     string  `json:"detector_profile_id"`
	DetectorRawActive                     bool    `json:"detector_raw_active"`
	DetectorActive                        bool    `json:"detector_active"`
	CandidateCount                        int     `json:"candidate_count"`
	AvgScore                              float64 `json:"avg_score"`
	AvgDistancePct                        float64 `json:"avg_distance_pct"`
	AvgSeedWickBeyondPct                  float64 `json:"avg_seed_wick_beyond_pct"`
	AvgConfirmationMovePct                float64 `json:"avg_confirmation_move_pct"`
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

type SRConfirmationTimingSummaryRow struct {
	Split                                                 string  `json:"split"`
	Side                                                  string  `json:"side"`
	ConfirmationDelayBars                                 int     `json:"confirmation_delay_bars"`
	HorizonBars                                           int     `json:"horizon_bars"`
	DetectorProfileID                                     string  `json:"detector_profile_id"`
	DetectorRawActive                                     bool    `json:"detector_raw_active"`
	DetectorActive                                        bool    `json:"detector_active"`
	CandidateCount                                        int     `json:"candidate_count"`
	ConfirmationFavorableCloseCount                       int     `json:"confirmation_favorable_close_count"`
	ConfirmationWrongSideCloseCount                       int     `json:"confirmation_wrong_side_close_count"`
	DecisionConfirmationCandidateCount                    int     `json:"decision_confirmation_candidate_count"`
	ConfirmationFavorableCloseRate                        float64 `json:"confirmation_favorable_close_rate"`
	ConfirmationWrongSideCloseRate                        float64 `json:"confirmation_wrong_side_close_rate"`
	DecisionConfirmationCandidateRate                     float64 `json:"decision_confirmation_candidate_rate"`
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

func DefaultSRConfirmationTimingAuditConfig() SRConfirmationTimingAuditConfig {
	boundaryDefaults := DefaultSRBoundaryAuditConfig()
	return SRConfirmationTimingAuditConfig{
		HorizonsBars:          append([]int(nil), boundaryDefaults.HorizonsBars...),
		ConfirmationDelayBars: []int{1, 2, 3},
		DetectorActiveOnly:    boundaryDefaults.DetectorActiveOnly,
	}
}

func RunSRConfirmationTimingAudit(candles []Candle, srRows []SRAuditRow, cfg SRConfirmationTimingAuditConfig) ([]SRConfirmationTimingCandidateRow, []SRConfirmationTimingSummaryRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	var events []srConfirmationTimingEvent
	for _, row := range srRows {
		if cfg.DetectorActiveOnly && !row.DetectorActive {
			continue
		}
		if row.HasSupport && row.NearSupport {
			events = appendSRConfirmationTimingEvents(events, candles, row, SRBoundarySideSupport, cfg.ConfirmationDelayBars, cfg.HorizonsBars)
		}
		if row.HasResistance && row.NearResistance {
			events = appendSRConfirmationTimingEvents(events, candles, row, SRBoundarySideResistance, cfg.ConfirmationDelayBars, cfg.HorizonsBars)
		}
	}

	return summarizeSRConfirmationTimingCandidates(events), summarizeSRConfirmationTimingSummary(events), nil
}

func (cfg SRConfirmationTimingAuditConfig) withDefaults() SRConfirmationTimingAuditConfig {
	defaults := DefaultSRConfirmationTimingAuditConfig()
	if len(cfg.HorizonsBars) == 0 && len(cfg.ConfirmationDelayBars) == 0 && !cfg.DetectorActiveOnly {
		return defaults
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if len(cfg.ConfirmationDelayBars) == 0 {
		cfg.ConfirmationDelayBars = append([]int(nil), defaults.ConfirmationDelayBars...)
	}
	return cfg
}

func (cfg SRConfirmationTimingAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("SR confirmation timing audit horizon bars must be positive")
		}
	}
	for _, delay := range cfg.ConfirmationDelayBars {
		if delay <= 0 {
			return fmt.Errorf("SR confirmation timing audit confirmation delay bars must be positive")
		}
	}
	return nil
}

type srConfirmationTimingDecision struct {
	Split                         string
	Side                          string
	ConfirmationDelayBars         int
	SeedCloseLocation             string
	SeedPiercedZone               bool
	SeedWickBeyondPct             float64
	SeedWickBeyondBucket          string
	ConfirmationCloseLocation     string
	ConfirmationFavorableClose    bool
	ConfirmationWrongSideClose    bool
	DecisionConfirmationCandidate bool
	ConfirmationMovePct           float64
	StrengthBucket                string
	DistanceBucket                string
	Score                         float64
	DistancePct                   float64
	DetectorProfileID             string
	DetectorRawActive             bool
	DetectorActive                bool
}

type srConfirmationTimingEvent struct {
	srConfirmationTimingDecision
	HorizonBars                      int
	LabelCloseBreak                  bool
	LabelWickBreak                   bool
	LabelReclaimedAfterBreak         bool
	LabelRejected                    bool
	LabelFavorableGreaterThanAdverse bool
	LabelFavorableMovePct            float64
	LabelAdverseMovePct              float64
}

func appendSRConfirmationTimingEvents(events []srConfirmationTimingEvent, candles []Candle, row SRAuditRow, side string, delays, horizons []int) []srConfirmationTimingEvent {
	for _, delay := range delays {
		decision, labelRow, ok := newSRConfirmationTimingDecision(candles, row, side, delay)
		if !ok {
			continue
		}
		for _, horizon := range horizons {
			label, ok := newSRBoundaryEvent(candles, labelRow, side, horizon)
			if !ok {
				continue
			}
			events = append(events, srConfirmationTimingEvent{
				srConfirmationTimingDecision:     decision,
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
	}
	return events
}

func newSRConfirmationTimingDecision(candles []Candle, row SRAuditRow, side string, delay int) (srConfirmationTimingDecision, SRAuditRow, bool) {
	if delay <= 0 || row.Index < 0 || row.Index >= len(candles) || row.Index+delay >= len(candles) {
		return srConfirmationTimingDecision{}, SRAuditRow{}, false
	}
	seed, ok := newSRRejectionTimingDecision(candles, row, side)
	if !ok || !seed.DecisionRejectionCandidate {
		return srConfirmationTimingDecision{}, SRAuditRow{}, false
	}

	confirmationIndex := row.Index + delay
	confirmationCandle := candles[confirmationIndex]
	labelRow := row
	labelRow.Index = confirmationIndex
	labelRow.OpenTime = confirmationCandle.OpenTime.Format(timeLayout)
	labelRow.CloseTime = confirmationCandle.CloseTime.Format(timeLayout)
	labelRow.Split = splitNameForCloseTime(confirmationCandle.CloseTime, DefaultSplits())
	labelRow.Close = confirmationCandle.Close

	decision := srConfirmationTimingDecision{
		Split:                 labelRow.Split,
		Side:                  side,
		ConfirmationDelayBars: delay,
		SeedCloseLocation:     seed.CloseLocation,
		SeedPiercedZone:       seed.PiercedZone,
		SeedWickBeyondPct:     seed.WickBeyondPct,
		SeedWickBeyondBucket:  seed.WickBeyondBucket,
		StrengthBucket:        seed.StrengthBucket,
		DistanceBucket:        seed.DistanceBucket,
		Score:                 seed.Score,
		DistancePct:           seed.DistancePct,
		DetectorProfileID:     seed.DetectorProfileID,
		DetectorRawActive:     seed.DetectorRawActive,
		DetectorActive:        seed.DetectorActive,
	}

	switch side {
	case SRBoundarySideSupport:
		decision.ConfirmationCloseLocation = closeLocationForBoundary(side, confirmationCandle.Close, row.NearestSupport, row.NearestSupportTop, row.NearestSupportBottom)
		decision.ConfirmationFavorableClose = confirmationCandle.Close > row.Close
		decision.ConfirmationWrongSideClose = confirmationCandle.Close < row.NearestSupportBottom
		decision.ConfirmationMovePct = movePct(confirmationCandle.Close-row.Close, row.Close)
	case SRBoundarySideResistance:
		decision.ConfirmationCloseLocation = closeLocationForBoundary(side, confirmationCandle.Close, row.NearestResistance, row.NearestResistanceTop, row.NearestResistanceBottom)
		decision.ConfirmationFavorableClose = confirmationCandle.Close < row.Close
		decision.ConfirmationWrongSideClose = confirmationCandle.Close > row.NearestResistanceTop
		decision.ConfirmationMovePct = movePct(row.Close-confirmationCandle.Close, row.Close)
	default:
		return srConfirmationTimingDecision{}, SRAuditRow{}, false
	}

	decision.DecisionConfirmationCandidate = decision.ConfirmationFavorableClose && !decision.ConfirmationWrongSideClose
	return decision, labelRow, true
}

type srConfirmationTimingCandidateKey struct {
	split                         string
	side                          string
	confirmationDelayBars         int
	horizonBars                   int
	seedCloseLocation             string
	seedPiercedZone               bool
	seedWickBeyondBucket          string
	confirmationCloseLocation     string
	confirmationFavorableClose    bool
	confirmationWrongSideClose    bool
	decisionConfirmationCandidate bool
	strengthBucket                string
	distanceBucket                string
	detectorProfileID             string
	detectorRawActive             bool
	detectorActive                bool
}

type srConfirmationTimingSummaryKey struct {
	split                 string
	side                  string
	confirmationDelayBars int
	horizonBars           int
	detectorProfileID     string
	detectorRawActive     bool
	detectorActive        bool
}

type srConfirmationTimingCandidateAccumulator struct {
	key                    srConfirmationTimingCandidateKey
	candidates             int
	scoreSum               float64
	distancePctSum         float64
	seedWickBeyondPctSum   float64
	confirmationMovePctSum float64
	labels                 srConfirmationTimingLabelAccumulator
}

type srConfirmationTimingSummaryAccumulator struct {
	key                            srConfirmationTimingSummaryKey
	candidates                     int
	confirmationFavorableCloses    int
	confirmationWrongSideCloses    int
	decisionConfirmationCandidates int
	labels                         srConfirmationTimingLabelAccumulator
	decisionCandidateLabels        srConfirmationTimingLabelAccumulator
}

type srConfirmationTimingLabelAccumulator struct {
	events                       int
	closeBreaks                  int
	wickBreaks                   int
	reclaimsAfterBreak           int
	rejections                   int
	favorableGreaterThanAdverses int
	favorablePctSum              float64
	adversePctSum                float64
}

func summarizeSRConfirmationTimingCandidates(events []srConfirmationTimingEvent) []SRConfirmationTimingCandidateRow {
	accumulators := map[srConfirmationTimingCandidateKey]*srConfirmationTimingCandidateAccumulator{}
	for _, event := range events {
		key := srConfirmationTimingCandidateKey{
			split:                         event.Split,
			side:                          event.Side,
			confirmationDelayBars:         event.ConfirmationDelayBars,
			horizonBars:                   event.HorizonBars,
			seedCloseLocation:             event.SeedCloseLocation,
			seedPiercedZone:               event.SeedPiercedZone,
			seedWickBeyondBucket:          event.SeedWickBeyondBucket,
			confirmationCloseLocation:     event.ConfirmationCloseLocation,
			confirmationFavorableClose:    event.ConfirmationFavorableClose,
			confirmationWrongSideClose:    event.ConfirmationWrongSideClose,
			decisionConfirmationCandidate: event.DecisionConfirmationCandidate,
			strengthBucket:                event.StrengthBucket,
			distanceBucket:                event.DistanceBucket,
			detectorProfileID:             event.DetectorProfileID,
			detectorRawActive:             event.DetectorRawActive,
			detectorActive:                event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srConfirmationTimingCandidateAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRConfirmationTimingCandidateRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRConfirmationTimingCandidateRow(rows[i], rows[j])
	})
	return rows
}

func summarizeSRConfirmationTimingSummary(events []srConfirmationTimingEvent) []SRConfirmationTimingSummaryRow {
	accumulators := map[srConfirmationTimingSummaryKey]*srConfirmationTimingSummaryAccumulator{}
	for _, event := range events {
		key := srConfirmationTimingSummaryKey{
			split:                 event.Split,
			side:                  event.Side,
			confirmationDelayBars: event.ConfirmationDelayBars,
			horizonBars:           event.HorizonBars,
			detectorProfileID:     event.DetectorProfileID,
			detectorRawActive:     event.DetectorRawActive,
			detectorActive:        event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srConfirmationTimingSummaryAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRConfirmationTimingSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRConfirmationTimingSummaryRow(rows[i], rows[j])
	})
	return rows
}

func (acc *srConfirmationTimingCandidateAccumulator) add(event srConfirmationTimingEvent) {
	acc.candidates++
	acc.scoreSum += event.Score
	acc.distancePctSum += event.DistancePct
	acc.seedWickBeyondPctSum += event.SeedWickBeyondPct
	acc.confirmationMovePctSum += event.ConfirmationMovePct
	acc.labels.add(event)
}

func (acc srConfirmationTimingCandidateAccumulator) row() SRConfirmationTimingCandidateRow {
	row := SRConfirmationTimingCandidateRow{
		Split:                         acc.key.split,
		Side:                          acc.key.side,
		ConfirmationDelayBars:         acc.key.confirmationDelayBars,
		HorizonBars:                   acc.key.horizonBars,
		SeedCloseLocation:             acc.key.seedCloseLocation,
		SeedPiercedZone:               acc.key.seedPiercedZone,
		SeedWickBeyondBucket:          acc.key.seedWickBeyondBucket,
		ConfirmationCloseLocation:     acc.key.confirmationCloseLocation,
		ConfirmationFavorableClose:    acc.key.confirmationFavorableClose,
		ConfirmationWrongSideClose:    acc.key.confirmationWrongSideClose,
		DecisionConfirmationCandidate: acc.key.decisionConfirmationCandidate,
		StrengthBucket:                acc.key.strengthBucket,
		DistanceBucket:                acc.key.distanceBucket,
		DetectorProfileID:             acc.key.detectorProfileID,
		DetectorRawActive:             acc.key.detectorRawActive,
		DetectorActive:                acc.key.detectorActive,
		CandidateCount:                acc.candidates,
	}
	if acc.candidates > 0 {
		row.AvgScore = acc.scoreSum / float64(acc.candidates)
		row.AvgDistancePct = acc.distancePctSum / float64(acc.candidates)
		row.AvgSeedWickBeyondPct = acc.seedWickBeyondPctSum / float64(acc.candidates)
		row.AvgConfirmationMovePct = acc.confirmationMovePctSum / float64(acc.candidates)
	}
	acc.labels.addToCandidateRow(&row)
	return row
}

func (acc *srConfirmationTimingSummaryAccumulator) add(event srConfirmationTimingEvent) {
	acc.candidates++
	if event.ConfirmationFavorableClose {
		acc.confirmationFavorableCloses++
	}
	if event.ConfirmationWrongSideClose {
		acc.confirmationWrongSideCloses++
	}
	if event.DecisionConfirmationCandidate {
		acc.decisionConfirmationCandidates++
		acc.decisionCandidateLabels.add(event)
	}
	acc.labels.add(event)
}

func (acc srConfirmationTimingSummaryAccumulator) row() SRConfirmationTimingSummaryRow {
	row := SRConfirmationTimingSummaryRow{
		Split:                              acc.key.split,
		Side:                               acc.key.side,
		ConfirmationDelayBars:              acc.key.confirmationDelayBars,
		HorizonBars:                        acc.key.horizonBars,
		DetectorProfileID:                  acc.key.detectorProfileID,
		DetectorRawActive:                  acc.key.detectorRawActive,
		DetectorActive:                     acc.key.detectorActive,
		CandidateCount:                     acc.candidates,
		ConfirmationFavorableCloseCount:    acc.confirmationFavorableCloses,
		ConfirmationWrongSideCloseCount:    acc.confirmationWrongSideCloses,
		DecisionConfirmationCandidateCount: acc.decisionConfirmationCandidates,
	}
	if acc.candidates > 0 {
		row.ConfirmationFavorableCloseRate = float64(acc.confirmationFavorableCloses) / float64(acc.candidates)
		row.ConfirmationWrongSideCloseRate = float64(acc.confirmationWrongSideCloses) / float64(acc.candidates)
		row.DecisionConfirmationCandidateRate = float64(acc.decisionConfirmationCandidates) / float64(acc.candidates)
	}
	acc.labels.addToSummaryRow(&row)
	acc.decisionCandidateLabels.addToDecisionCandidateSummaryRow(&row)
	return row
}

func (acc *srConfirmationTimingLabelAccumulator) add(event srConfirmationTimingEvent) {
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

func (acc srConfirmationTimingLabelAccumulator) avgFavorablePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.favorablePctSum / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) avgAdversePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.adversePctSum / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) favorableGreaterThanAdverseRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.favorableGreaterThanAdverses) / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) closeBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.closeBreaks) / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) wickBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.wickBreaks) / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) reclaimEventRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) reclaimGivenCloseBreakRate() float64 {
	if acc.closeBreaks == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.closeBreaks)
}

func (acc srConfirmationTimingLabelAccumulator) rejectionRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.rejections) / float64(acc.events)
}

func (acc srConfirmationTimingLabelAccumulator) addToCandidateRow(row *SRConfirmationTimingCandidateRow) {
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

func (acc srConfirmationTimingLabelAccumulator) addToSummaryRow(row *SRConfirmationTimingSummaryRow) {
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

func (acc srConfirmationTimingLabelAccumulator) addToDecisionCandidateSummaryRow(row *SRConfirmationTimingSummaryRow) {
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

func lessSRConfirmationTimingCandidateRow(a, b SRConfirmationTimingCandidateRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.ConfirmationDelayBars != b.ConfirmationDelayBars {
		return a.ConfirmationDelayBars < b.ConfirmationDelayBars
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
	if closeLocationSortKey(a.SeedCloseLocation) != closeLocationSortKey(b.SeedCloseLocation) {
		return closeLocationSortKey(a.SeedCloseLocation) < closeLocationSortKey(b.SeedCloseLocation)
	}
	if boolSortKey(a.SeedPiercedZone) != boolSortKey(b.SeedPiercedZone) {
		return boolSortKey(a.SeedPiercedZone) < boolSortKey(b.SeedPiercedZone)
	}
	if distanceBucketSortKey(a.SeedWickBeyondBucket) != distanceBucketSortKey(b.SeedWickBeyondBucket) {
		return distanceBucketSortKey(a.SeedWickBeyondBucket) < distanceBucketSortKey(b.SeedWickBeyondBucket)
	}
	if closeLocationSortKey(a.ConfirmationCloseLocation) != closeLocationSortKey(b.ConfirmationCloseLocation) {
		return closeLocationSortKey(a.ConfirmationCloseLocation) < closeLocationSortKey(b.ConfirmationCloseLocation)
	}
	if boolSortKey(a.ConfirmationFavorableClose) != boolSortKey(b.ConfirmationFavorableClose) {
		return boolSortKey(a.ConfirmationFavorableClose) < boolSortKey(b.ConfirmationFavorableClose)
	}
	if boolSortKey(a.ConfirmationWrongSideClose) != boolSortKey(b.ConfirmationWrongSideClose) {
		return boolSortKey(a.ConfirmationWrongSideClose) < boolSortKey(b.ConfirmationWrongSideClose)
	}
	if boolSortKey(a.DecisionConfirmationCandidate) != boolSortKey(b.DecisionConfirmationCandidate) {
		return boolSortKey(a.DecisionConfirmationCandidate) < boolSortKey(b.DecisionConfirmationCandidate)
	}
	if strengthBucketSortKey(a.StrengthBucket) != strengthBucketSortKey(b.StrengthBucket) {
		return strengthBucketSortKey(a.StrengthBucket) < strengthBucketSortKey(b.StrengthBucket)
	}
	return distanceBucketSortKey(a.DistanceBucket) < distanceBucketSortKey(b.DistanceBucket)
}

func lessSRConfirmationTimingSummaryRow(a, b SRConfirmationTimingSummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.ConfirmationDelayBars != b.ConfirmationDelayBars {
		return a.ConfirmationDelayBars < b.ConfirmationDelayBars
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
