package lab

import (
	"fmt"
	"sort"
)

type SRFalseBreakReclaimTimingAuditConfig struct {
	HorizonsBars        []int
	MaxBreakDelayBars   int
	MaxReclaimDelayBars int
	DetectorActiveOnly  bool
}

type SRFalseBreakReclaimTimingCandidateRow struct {
	Split                                 string  `json:"split"`
	Side                                  string  `json:"side"`
	BreakDelayBars                        int     `json:"break_delay_bars"`
	ReclaimDelayBars                      int     `json:"reclaim_delay_bars"`
	TotalDelayBars                        int     `json:"total_delay_bars"`
	HorizonBars                           int     `json:"horizon_bars"`
	AnchorCloseLocation                   string  `json:"anchor_close_location"`
	BreakCloseLocation                    string  `json:"break_close_location"`
	ReclaimCloseLocation                  string  `json:"reclaim_close_location"`
	BreakMoveBucket                       string  `json:"break_move_bucket"`
	ReclaimMoveBucket                     string  `json:"reclaim_move_bucket"`
	DecisionFalseBreakReclaimCandidate    bool    `json:"decision_false_break_reclaim_candidate"`
	StrengthBucket                        string  `json:"strength_bucket"`
	DistanceBucket                        string  `json:"distance_bucket"`
	DetectorProfileID                     string  `json:"detector_profile_id"`
	DetectorRawActive                     bool    `json:"detector_raw_active"`
	DetectorActive                        bool    `json:"detector_active"`
	CandidateCount                        int     `json:"candidate_count"`
	AvgScore                              float64 `json:"avg_score"`
	AvgDistancePct                        float64 `json:"avg_distance_pct"`
	AvgBreakMovePct                       float64 `json:"avg_break_move_pct"`
	AvgReclaimMovePct                     float64 `json:"avg_reclaim_move_pct"`
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

type SRFalseBreakReclaimTimingSummaryRow struct {
	Split                                                 string  `json:"split"`
	Side                                                  string  `json:"side"`
	HorizonBars                                           int     `json:"horizon_bars"`
	DetectorProfileID                                     string  `json:"detector_profile_id"`
	DetectorRawActive                                     bool    `json:"detector_raw_active"`
	DetectorActive                                        bool    `json:"detector_active"`
	CandidateCount                                        int     `json:"candidate_count"`
	DecisionFalseBreakReclaimCandidateCount               int     `json:"decision_false_break_reclaim_candidate_count"`
	DecisionFalseBreakReclaimCandidateRate                float64 `json:"decision_false_break_reclaim_candidate_rate"`
	AvgBreakDelayBars                                     float64 `json:"avg_break_delay_bars"`
	AvgReclaimDelayBars                                   float64 `json:"avg_reclaim_delay_bars"`
	AvgTotalDelayBars                                     float64 `json:"avg_total_delay_bars"`
	AvgBreakMovePct                                       float64 `json:"avg_break_move_pct"`
	AvgReclaimMovePct                                     float64 `json:"avg_reclaim_move_pct"`
	LabelCloseBreakRate                                   float64 `json:"label_close_break_rate"`
	LabelWickBreakRate                                    float64 `json:"label_wick_break_rate"`
	LabelReclaimEventRate                                 float64 `json:"label_reclaim_event_rate"`
	LabelReclaimGivenCloseBreakRate                       float64 `json:"label_reclaim_given_close_break_rate"`
	LabelRejectionRate                                    float64 `json:"label_rejection_rate"`
	LabelAvgFavorablePct                                  float64 `json:"label_avg_favorable_pct"`
	LabelAvgAdversePct                                    float64 `json:"label_avg_adverse_pct"`
	LabelFavorableMinusAdversePct                         float64 `json:"label_favorable_minus_adverse_pct"`
	LabelFavorableGreaterThanAdverseRate                  float64 `json:"label_favorable_greater_than_adverse_rate"`
	LabelDecisionCandidateCloseBreakRate                  float64 `json:"label_decision_candidate_close_break_rate"`
	LabelDecisionCandidateWickBreakRate                   float64 `json:"label_decision_candidate_wick_break_rate"`
	LabelDecisionCandidateReclaimEventRate                float64 `json:"label_decision_candidate_reclaim_event_rate"`
	LabelDecisionCandidateReclaimGivenCloseBreakRate      float64 `json:"label_decision_candidate_reclaim_given_close_break_rate"`
	LabelDecisionCandidateRejectionRate                   float64 `json:"label_decision_candidate_rejection_rate"`
	LabelDecisionCandidateAvgFavorablePct                 float64 `json:"label_decision_candidate_avg_favorable_pct"`
	LabelDecisionCandidateAvgAdversePct                   float64 `json:"label_decision_candidate_avg_adverse_pct"`
	LabelDecisionCandidateFavorableMinusAdversePct        float64 `json:"label_decision_candidate_favorable_minus_adverse_pct"`
	LabelDecisionCandidateFavorableGreaterThanAdverseRate float64 `json:"label_decision_candidate_favorable_greater_than_adverse_rate"`
}

func DefaultSRFalseBreakReclaimTimingAuditConfig() SRFalseBreakReclaimTimingAuditConfig {
	boundaryDefaults := DefaultSRBoundaryAuditConfig()
	return SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        append([]int(nil), boundaryDefaults.HorizonsBars...),
		MaxBreakDelayBars:   3,
		MaxReclaimDelayBars: 12,
		DetectorActiveOnly:  boundaryDefaults.DetectorActiveOnly,
	}
}

func RunSRFalseBreakReclaimTimingAudit(candles []Candle, srRows []SRAuditRow, cfg SRFalseBreakReclaimTimingAuditConfig) ([]SRFalseBreakReclaimTimingCandidateRow, []SRFalseBreakReclaimTimingSummaryRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	var events []srFalseBreakReclaimTimingEvent
	for _, row := range srRows {
		if cfg.DetectorActiveOnly && !row.DetectorActive {
			continue
		}
		if row.HasSupport && row.NearSupport {
			events = appendSRFalseBreakReclaimTimingEvents(events, candles, row, SRBoundarySideSupport, cfg)
		}
		if row.HasResistance && row.NearResistance {
			events = appendSRFalseBreakReclaimTimingEvents(events, candles, row, SRBoundarySideResistance, cfg)
		}
	}

	return summarizeSRFalseBreakReclaimTimingCandidates(events), summarizeSRFalseBreakReclaimTimingSummary(events), nil
}

func (cfg SRFalseBreakReclaimTimingAuditConfig) withDefaults() SRFalseBreakReclaimTimingAuditConfig {
	defaults := DefaultSRFalseBreakReclaimTimingAuditConfig()
	if len(cfg.HorizonsBars) == 0 && cfg.MaxBreakDelayBars == 0 && cfg.MaxReclaimDelayBars == 0 && !cfg.DetectorActiveOnly {
		return defaults
	}
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.MaxBreakDelayBars == 0 {
		cfg.MaxBreakDelayBars = defaults.MaxBreakDelayBars
	}
	if cfg.MaxReclaimDelayBars == 0 {
		cfg.MaxReclaimDelayBars = defaults.MaxReclaimDelayBars
	}
	return cfg
}

func (cfg SRFalseBreakReclaimTimingAuditConfig) validate() error {
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("SR false-break reclaim timing audit horizon bars must be positive")
		}
	}
	if cfg.MaxBreakDelayBars <= 0 {
		return fmt.Errorf("SR false-break reclaim timing audit max break delay bars must be positive")
	}
	if cfg.MaxReclaimDelayBars <= 0 {
		return fmt.Errorf("SR false-break reclaim timing audit max reclaim delay bars must be positive")
	}
	return nil
}

type srFalseBreakReclaimTimingDecision struct {
	Split                              string
	Side                               string
	BreakDelayBars                     int
	ReclaimDelayBars                   int
	TotalDelayBars                     int
	AnchorCloseLocation                string
	BreakCloseLocation                 string
	ReclaimCloseLocation               string
	BreakMovePct                       float64
	BreakMoveBucket                    string
	ReclaimMovePct                     float64
	ReclaimMoveBucket                  string
	DecisionFalseBreakReclaimCandidate bool
	StrengthBucket                     string
	DistanceBucket                     string
	Score                              float64
	DistancePct                        float64
	DetectorProfileID                  string
	DetectorRawActive                  bool
	DetectorActive                     bool
}

type srFalseBreakReclaimTimingEvent struct {
	srFalseBreakReclaimTimingDecision
	HorizonBars                      int
	LabelCloseBreak                  bool
	LabelWickBreak                   bool
	LabelReclaimedAfterBreak         bool
	LabelRejected                    bool
	LabelFavorableGreaterThanAdverse bool
	LabelFavorableMovePct            float64
	LabelAdverseMovePct              float64
}

func appendSRFalseBreakReclaimTimingEvents(events []srFalseBreakReclaimTimingEvent, candles []Candle, row SRAuditRow, side string, cfg SRFalseBreakReclaimTimingAuditConfig) []srFalseBreakReclaimTimingEvent {
	decision, labelRow, ok := newSRFalseBreakReclaimTimingDecision(candles, row, side, cfg.MaxBreakDelayBars, cfg.MaxReclaimDelayBars)
	if !ok {
		return events
	}
	for _, horizon := range cfg.HorizonsBars {
		label, ok := newSRBoundaryEvent(candles, labelRow, side, horizon)
		if !ok {
			continue
		}
		events = append(events, srFalseBreakReclaimTimingEvent{
			srFalseBreakReclaimTimingDecision: decision,
			HorizonBars:                       horizon,
			LabelCloseBreak:                   label.CloseBreak,
			LabelWickBreak:                    label.WickBreak,
			LabelReclaimedAfterBreak:          label.ReclaimedAfterBreak,
			LabelRejected:                     label.Rejected,
			LabelFavorableGreaterThanAdverse:  label.FavorableGreaterThanAdverse,
			LabelFavorableMovePct:             label.FavorableMovePct,
			LabelAdverseMovePct:               label.AdverseMovePct,
		})
	}
	return events
}

func newSRFalseBreakReclaimTimingDecision(candles []Candle, row SRAuditRow, side string, maxBreakDelay, maxReclaimDelay int) (srFalseBreakReclaimTimingDecision, SRAuditRow, bool) {
	if maxBreakDelay <= 0 || maxReclaimDelay <= 0 || row.Index < 0 || row.Index >= len(candles) {
		return srFalseBreakReclaimTimingDecision{}, SRAuditRow{}, false
	}
	breakIndex, reclaimIndex, ok := findSRFalseBreakReclaim(candles, row, side, maxBreakDelay, maxReclaimDelay)
	if !ok {
		return srFalseBreakReclaimTimingDecision{}, SRAuditRow{}, false
	}

	breakCandle := candles[breakIndex]
	reclaimCandle := candles[reclaimIndex]
	labelRow := row
	labelRow.Index = reclaimIndex
	labelRow.OpenTime = reclaimCandle.OpenTime.Format(timeLayout)
	labelRow.CloseTime = reclaimCandle.CloseTime.Format(timeLayout)
	labelRow.Split = splitNameForCloseTime(reclaimCandle.CloseTime, DefaultSplits())
	labelRow.Close = reclaimCandle.Close

	decision := srFalseBreakReclaimTimingDecision{
		Split:                              labelRow.Split,
		Side:                               side,
		BreakDelayBars:                     breakIndex - row.Index,
		ReclaimDelayBars:                   reclaimIndex - breakIndex,
		TotalDelayBars:                     reclaimIndex - row.Index,
		DecisionFalseBreakReclaimCandidate: true,
		DetectorProfileID:                  row.DetectorProfileID,
		DetectorRawActive:                  row.DetectorRawActive,
		DetectorActive:                     row.DetectorActive,
	}

	switch side {
	case SRBoundarySideSupport:
		decision.AnchorCloseLocation = closeLocationForBoundary(side, row.Close, row.NearestSupport, row.NearestSupportTop, row.NearestSupportBottom)
		decision.BreakCloseLocation = closeLocationForBoundary(side, breakCandle.Close, row.NearestSupport, row.NearestSupportTop, row.NearestSupportBottom)
		decision.ReclaimCloseLocation = closeLocationForBoundary(side, reclaimCandle.Close, row.NearestSupport, row.NearestSupportTop, row.NearestSupportBottom)
		decision.BreakMovePct = movePct(row.NearestSupportBottom-breakCandle.Close, row.Close)
		decision.ReclaimMovePct = movePct(reclaimCandle.Close-breakCandle.Close, row.Close)
		decision.StrengthBucket = strengthBucket(row.NearestSupportStrength)
		decision.DistanceBucket = distanceBucket(row.NearestSupportDistancePct)
		decision.Score = row.NearestSupportScore
		decision.DistancePct = row.NearestSupportDistancePct
	case SRBoundarySideResistance:
		decision.AnchorCloseLocation = closeLocationForBoundary(side, row.Close, row.NearestResistance, row.NearestResistanceTop, row.NearestResistanceBottom)
		decision.BreakCloseLocation = closeLocationForBoundary(side, breakCandle.Close, row.NearestResistance, row.NearestResistanceTop, row.NearestResistanceBottom)
		decision.ReclaimCloseLocation = closeLocationForBoundary(side, reclaimCandle.Close, row.NearestResistance, row.NearestResistanceTop, row.NearestResistanceBottom)
		decision.BreakMovePct = movePct(breakCandle.Close-row.NearestResistanceTop, row.Close)
		decision.ReclaimMovePct = movePct(breakCandle.Close-reclaimCandle.Close, row.Close)
		decision.StrengthBucket = strengthBucket(row.NearestResistanceStrength)
		decision.DistanceBucket = distanceBucket(row.NearestResistanceDistancePct)
		decision.Score = row.NearestResistanceScore
		decision.DistancePct = row.NearestResistanceDistancePct
	default:
		return srFalseBreakReclaimTimingDecision{}, SRAuditRow{}, false
	}

	decision.BreakMoveBucket = distanceBucket(decision.BreakMovePct)
	decision.ReclaimMoveBucket = distanceBucket(decision.ReclaimMovePct)
	return decision, labelRow, true
}

func findSRFalseBreakReclaim(candles []Candle, row SRAuditRow, side string, maxBreakDelay, maxReclaimDelay int) (int, int, bool) {
	breakLimit := minInt(len(candles)-1, row.Index+maxBreakDelay)
	for breakIndex := row.Index + 1; breakIndex <= breakLimit; breakIndex++ {
		if !isSRFalseBreakClose(candles[breakIndex], row, side) {
			continue
		}
		reclaimLimit := minInt(len(candles)-1, breakIndex+maxReclaimDelay)
		for reclaimIndex := breakIndex + 1; reclaimIndex <= reclaimLimit; reclaimIndex++ {
			if isSRReclaimClose(candles[reclaimIndex], row, side) {
				return breakIndex, reclaimIndex, true
			}
		}
		return 0, 0, false
	}
	return 0, 0, false
}

func isSRFalseBreakClose(candle Candle, row SRAuditRow, side string) bool {
	switch side {
	case SRBoundarySideSupport:
		return candle.Close < row.NearestSupportBottom
	case SRBoundarySideResistance:
		return candle.Close > row.NearestResistanceTop
	default:
		return false
	}
}

func isSRReclaimClose(candle Candle, row SRAuditRow, side string) bool {
	switch side {
	case SRBoundarySideSupport:
		return candle.Close >= row.NearestSupport
	case SRBoundarySideResistance:
		return candle.Close <= row.NearestResistance
	default:
		return false
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type srFalseBreakReclaimTimingCandidateKey struct {
	split                              string
	side                               string
	breakDelayBars                     int
	reclaimDelayBars                   int
	totalDelayBars                     int
	horizonBars                        int
	anchorCloseLocation                string
	breakCloseLocation                 string
	reclaimCloseLocation               string
	breakMoveBucket                    string
	reclaimMoveBucket                  string
	decisionFalseBreakReclaimCandidate bool
	strengthBucket                     string
	distanceBucket                     string
	detectorProfileID                  string
	detectorRawActive                  bool
	detectorActive                     bool
}

type srFalseBreakReclaimTimingSummaryKey struct {
	split             string
	side              string
	horizonBars       int
	detectorProfileID string
	detectorRawActive bool
	detectorActive    bool
}

type srFalseBreakReclaimTimingCandidateAccumulator struct {
	key               srFalseBreakReclaimTimingCandidateKey
	candidates        int
	scoreSum          float64
	distancePctSum    float64
	breakMovePctSum   float64
	reclaimMovePctSum float64
	labels            srFalseBreakReclaimTimingLabelAccumulator
}

type srFalseBreakReclaimTimingSummaryAccumulator struct {
	key                                 srFalseBreakReclaimTimingSummaryKey
	candidates                          int
	decisionFalseBreakReclaimCandidates int
	breakDelayBarsSum                   int
	reclaimDelayBarsSum                 int
	totalDelayBarsSum                   int
	breakMovePctSum                     float64
	reclaimMovePctSum                   float64
	labels                              srFalseBreakReclaimTimingLabelAccumulator
	decisionCandidateLabels             srFalseBreakReclaimTimingLabelAccumulator
}

type srFalseBreakReclaimTimingLabelAccumulator struct {
	events                       int
	closeBreaks                  int
	wickBreaks                   int
	reclaimsAfterBreak           int
	rejections                   int
	favorableGreaterThanAdverses int
	favorablePctSum              float64
	adversePctSum                float64
}

func summarizeSRFalseBreakReclaimTimingCandidates(events []srFalseBreakReclaimTimingEvent) []SRFalseBreakReclaimTimingCandidateRow {
	accumulators := map[srFalseBreakReclaimTimingCandidateKey]*srFalseBreakReclaimTimingCandidateAccumulator{}
	for _, event := range events {
		key := srFalseBreakReclaimTimingCandidateKey{
			split:                              event.Split,
			side:                               event.Side,
			breakDelayBars:                     event.BreakDelayBars,
			reclaimDelayBars:                   event.ReclaimDelayBars,
			totalDelayBars:                     event.TotalDelayBars,
			horizonBars:                        event.HorizonBars,
			anchorCloseLocation:                event.AnchorCloseLocation,
			breakCloseLocation:                 event.BreakCloseLocation,
			reclaimCloseLocation:               event.ReclaimCloseLocation,
			breakMoveBucket:                    event.BreakMoveBucket,
			reclaimMoveBucket:                  event.ReclaimMoveBucket,
			decisionFalseBreakReclaimCandidate: event.DecisionFalseBreakReclaimCandidate,
			strengthBucket:                     event.StrengthBucket,
			distanceBucket:                     event.DistanceBucket,
			detectorProfileID:                  event.DetectorProfileID,
			detectorRawActive:                  event.DetectorRawActive,
			detectorActive:                     event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srFalseBreakReclaimTimingCandidateAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRFalseBreakReclaimTimingCandidateRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRFalseBreakReclaimTimingCandidateRow(rows[i], rows[j])
	})
	return rows
}

func summarizeSRFalseBreakReclaimTimingSummary(events []srFalseBreakReclaimTimingEvent) []SRFalseBreakReclaimTimingSummaryRow {
	accumulators := map[srFalseBreakReclaimTimingSummaryKey]*srFalseBreakReclaimTimingSummaryAccumulator{}
	for _, event := range events {
		key := srFalseBreakReclaimTimingSummaryKey{
			split:             event.Split,
			side:              event.Side,
			horizonBars:       event.HorizonBars,
			detectorProfileID: event.DetectorProfileID,
			detectorRawActive: event.DetectorRawActive,
			detectorActive:    event.DetectorActive,
		}
		acc := accumulators[key]
		if acc == nil {
			acc = &srFalseBreakReclaimTimingSummaryAccumulator{key: key}
			accumulators[key] = acc
		}
		acc.add(event)
	}

	rows := make([]SRFalseBreakReclaimTimingSummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		rows = append(rows, acc.row())
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessSRFalseBreakReclaimTimingSummaryRow(rows[i], rows[j])
	})
	return rows
}

func (acc *srFalseBreakReclaimTimingCandidateAccumulator) add(event srFalseBreakReclaimTimingEvent) {
	acc.candidates++
	acc.scoreSum += event.Score
	acc.distancePctSum += event.DistancePct
	acc.breakMovePctSum += event.BreakMovePct
	acc.reclaimMovePctSum += event.ReclaimMovePct
	acc.labels.add(event)
}

func (acc srFalseBreakReclaimTimingCandidateAccumulator) row() SRFalseBreakReclaimTimingCandidateRow {
	row := SRFalseBreakReclaimTimingCandidateRow{
		Split:                              acc.key.split,
		Side:                               acc.key.side,
		BreakDelayBars:                     acc.key.breakDelayBars,
		ReclaimDelayBars:                   acc.key.reclaimDelayBars,
		TotalDelayBars:                     acc.key.totalDelayBars,
		HorizonBars:                        acc.key.horizonBars,
		AnchorCloseLocation:                acc.key.anchorCloseLocation,
		BreakCloseLocation:                 acc.key.breakCloseLocation,
		ReclaimCloseLocation:               acc.key.reclaimCloseLocation,
		BreakMoveBucket:                    acc.key.breakMoveBucket,
		ReclaimMoveBucket:                  acc.key.reclaimMoveBucket,
		DecisionFalseBreakReclaimCandidate: acc.key.decisionFalseBreakReclaimCandidate,
		StrengthBucket:                     acc.key.strengthBucket,
		DistanceBucket:                     acc.key.distanceBucket,
		DetectorProfileID:                  acc.key.detectorProfileID,
		DetectorRawActive:                  acc.key.detectorRawActive,
		DetectorActive:                     acc.key.detectorActive,
		CandidateCount:                     acc.candidates,
	}
	if acc.candidates > 0 {
		row.AvgScore = acc.scoreSum / float64(acc.candidates)
		row.AvgDistancePct = acc.distancePctSum / float64(acc.candidates)
		row.AvgBreakMovePct = acc.breakMovePctSum / float64(acc.candidates)
		row.AvgReclaimMovePct = acc.reclaimMovePctSum / float64(acc.candidates)
	}
	acc.labels.addToCandidateRow(&row)
	return row
}

func (acc *srFalseBreakReclaimTimingSummaryAccumulator) add(event srFalseBreakReclaimTimingEvent) {
	acc.candidates++
	acc.breakDelayBarsSum += event.BreakDelayBars
	acc.reclaimDelayBarsSum += event.ReclaimDelayBars
	acc.totalDelayBarsSum += event.TotalDelayBars
	acc.breakMovePctSum += event.BreakMovePct
	acc.reclaimMovePctSum += event.ReclaimMovePct
	if event.DecisionFalseBreakReclaimCandidate {
		acc.decisionFalseBreakReclaimCandidates++
		acc.decisionCandidateLabels.add(event)
	}
	acc.labels.add(event)
}

func (acc srFalseBreakReclaimTimingSummaryAccumulator) row() SRFalseBreakReclaimTimingSummaryRow {
	row := SRFalseBreakReclaimTimingSummaryRow{
		Split:                                   acc.key.split,
		Side:                                    acc.key.side,
		HorizonBars:                             acc.key.horizonBars,
		DetectorProfileID:                       acc.key.detectorProfileID,
		DetectorRawActive:                       acc.key.detectorRawActive,
		DetectorActive:                          acc.key.detectorActive,
		CandidateCount:                          acc.candidates,
		DecisionFalseBreakReclaimCandidateCount: acc.decisionFalseBreakReclaimCandidates,
	}
	if acc.candidates > 0 {
		row.DecisionFalseBreakReclaimCandidateRate = float64(acc.decisionFalseBreakReclaimCandidates) / float64(acc.candidates)
		row.AvgBreakDelayBars = float64(acc.breakDelayBarsSum) / float64(acc.candidates)
		row.AvgReclaimDelayBars = float64(acc.reclaimDelayBarsSum) / float64(acc.candidates)
		row.AvgTotalDelayBars = float64(acc.totalDelayBarsSum) / float64(acc.candidates)
		row.AvgBreakMovePct = acc.breakMovePctSum / float64(acc.candidates)
		row.AvgReclaimMovePct = acc.reclaimMovePctSum / float64(acc.candidates)
	}
	acc.labels.addToSummaryRow(&row)
	acc.decisionCandidateLabels.addToDecisionCandidateSummaryRow(&row)
	return row
}

func (acc *srFalseBreakReclaimTimingLabelAccumulator) add(event srFalseBreakReclaimTimingEvent) {
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

func (acc srFalseBreakReclaimTimingLabelAccumulator) avgFavorablePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.favorablePctSum / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) avgAdversePct() float64 {
	if acc.events == 0 {
		return 0
	}
	return acc.adversePctSum / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) favorableGreaterThanAdverseRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.favorableGreaterThanAdverses) / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) closeBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.closeBreaks) / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) wickBreakRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.wickBreaks) / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) reclaimEventRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) reclaimGivenCloseBreakRate() float64 {
	if acc.closeBreaks == 0 {
		return 0
	}
	return float64(acc.reclaimsAfterBreak) / float64(acc.closeBreaks)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) rejectionRate() float64 {
	if acc.events == 0 {
		return 0
	}
	return float64(acc.rejections) / float64(acc.events)
}

func (acc srFalseBreakReclaimTimingLabelAccumulator) addToCandidateRow(row *SRFalseBreakReclaimTimingCandidateRow) {
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

func (acc srFalseBreakReclaimTimingLabelAccumulator) addToSummaryRow(row *SRFalseBreakReclaimTimingSummaryRow) {
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

func (acc srFalseBreakReclaimTimingLabelAccumulator) addToDecisionCandidateSummaryRow(row *SRFalseBreakReclaimTimingSummaryRow) {
	row.LabelDecisionCandidateCloseBreakRate = acc.closeBreakRate()
	row.LabelDecisionCandidateWickBreakRate = acc.wickBreakRate()
	row.LabelDecisionCandidateReclaimEventRate = acc.reclaimEventRate()
	row.LabelDecisionCandidateReclaimGivenCloseBreakRate = acc.reclaimGivenCloseBreakRate()
	row.LabelDecisionCandidateRejectionRate = acc.rejectionRate()
	row.LabelDecisionCandidateAvgFavorablePct = acc.avgFavorablePct()
	row.LabelDecisionCandidateAvgAdversePct = acc.avgAdversePct()
	row.LabelDecisionCandidateFavorableMinusAdversePct = row.LabelDecisionCandidateAvgFavorablePct - row.LabelDecisionCandidateAvgAdversePct
	row.LabelDecisionCandidateFavorableGreaterThanAdverseRate = acc.favorableGreaterThanAdverseRate()
}

func lessSRFalseBreakReclaimTimingCandidateRow(a, b SRFalseBreakReclaimTimingCandidateRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if sideSortKey(a.Side) != sideSortKey(b.Side) {
		return sideSortKey(a.Side) < sideSortKey(b.Side)
	}
	if a.BreakDelayBars != b.BreakDelayBars {
		return a.BreakDelayBars < b.BreakDelayBars
	}
	if a.ReclaimDelayBars != b.ReclaimDelayBars {
		return a.ReclaimDelayBars < b.ReclaimDelayBars
	}
	if a.TotalDelayBars != b.TotalDelayBars {
		return a.TotalDelayBars < b.TotalDelayBars
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
	if closeLocationSortKey(a.AnchorCloseLocation) != closeLocationSortKey(b.AnchorCloseLocation) {
		return closeLocationSortKey(a.AnchorCloseLocation) < closeLocationSortKey(b.AnchorCloseLocation)
	}
	if closeLocationSortKey(a.BreakCloseLocation) != closeLocationSortKey(b.BreakCloseLocation) {
		return closeLocationSortKey(a.BreakCloseLocation) < closeLocationSortKey(b.BreakCloseLocation)
	}
	if closeLocationSortKey(a.ReclaimCloseLocation) != closeLocationSortKey(b.ReclaimCloseLocation) {
		return closeLocationSortKey(a.ReclaimCloseLocation) < closeLocationSortKey(b.ReclaimCloseLocation)
	}
	if distanceBucketSortKey(a.BreakMoveBucket) != distanceBucketSortKey(b.BreakMoveBucket) {
		return distanceBucketSortKey(a.BreakMoveBucket) < distanceBucketSortKey(b.BreakMoveBucket)
	}
	if distanceBucketSortKey(a.ReclaimMoveBucket) != distanceBucketSortKey(b.ReclaimMoveBucket) {
		return distanceBucketSortKey(a.ReclaimMoveBucket) < distanceBucketSortKey(b.ReclaimMoveBucket)
	}
	if boolSortKey(a.DecisionFalseBreakReclaimCandidate) != boolSortKey(b.DecisionFalseBreakReclaimCandidate) {
		return boolSortKey(a.DecisionFalseBreakReclaimCandidate) < boolSortKey(b.DecisionFalseBreakReclaimCandidate)
	}
	if strengthBucketSortKey(a.StrengthBucket) != strengthBucketSortKey(b.StrengthBucket) {
		return strengthBucketSortKey(a.StrengthBucket) < strengthBucketSortKey(b.StrengthBucket)
	}
	return distanceBucketSortKey(a.DistanceBucket) < distanceBucketSortKey(b.DistanceBucket)
}

func lessSRFalseBreakReclaimTimingSummaryRow(a, b SRFalseBreakReclaimTimingSummaryRow) bool {
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
