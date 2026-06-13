package lab

import "testing"

func TestSRRejectionTimingDecisionUsesCurrentCandleOnly(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 105, 97, 100.5),
		testCandle(1, 100.5, 104, 99, 103),
	}
	row := supportTimingRow(0, candles[0], true)

	decision, ok := newSRRejectionTimingDecision(candles, row, SRBoundarySideSupport)
	if !ok {
		t.Fatalf("expected support timing decision")
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[1] = testCandle(1, 100.5, 120, 80, 81)
	changedDecision, ok := newSRRejectionTimingDecision(changedFuture, row, SRBoundarySideSupport)
	if !ok {
		t.Fatalf("expected support timing decision after future change")
	}
	if decision != changedDecision {
		t.Fatalf("decision features changed after future-only edit: got %+v want %+v", changedDecision, decision)
	}

	if !decision.TouchedZone || !decision.PiercedZone || !decision.ClosedBack || !decision.DecisionRejectionCandidate {
		t.Fatalf("bad support decision flags: %+v", decision)
	}
	if decision.CloseLocation != srCloseLocationInZoneAboveLevel {
		t.Fatalf("support close location=%q", decision.CloseLocation)
	}
	if !boundaryAlmostEqual(decision.WickBeyondPct, (98.0-97.0)/100.5) || decision.WickBeyondBucket != "gt_20bp" {
		t.Fatalf("bad support wick bucket/pct: %+v", decision)
	}
}

func TestSRRejectionTimingDecisionSupportResistanceSymmetry(t *testing.T) {
	supportCandle := testCandle(0, 100, 103, 97, 100.5)
	supportDecision, ok := newSRRejectionTimingDecision(
		[]Candle{supportCandle},
		supportTimingRow(0, supportCandle, true),
		SRBoundarySideSupport,
	)
	if !ok {
		t.Fatalf("expected support decision")
	}
	if !supportDecision.TouchedZone || !supportDecision.PiercedZone || !supportDecision.ClosedBack ||
		supportDecision.CloseLocation != srCloseLocationInZoneAboveLevel {
		t.Fatalf("bad support symmetry decision: %+v", supportDecision)
	}

	resistanceCandle := testCandle(0, 100, 103, 97, 99.5)
	resistanceDecision, ok := newSRRejectionTimingDecision(
		[]Candle{resistanceCandle},
		resistanceTimingRow(0, resistanceCandle, true),
		SRBoundarySideResistance,
	)
	if !ok {
		t.Fatalf("expected resistance decision")
	}
	if !resistanceDecision.TouchedZone || !resistanceDecision.PiercedZone || !resistanceDecision.ClosedBack ||
		resistanceDecision.CloseLocation != srCloseLocationInZoneBelowLevel {
		t.Fatalf("bad resistance symmetry decision: %+v", resistanceDecision)
	}
	if !boundaryAlmostEqual(resistanceDecision.WickBeyondPct, (103.0-102.0)/99.5) ||
		resistanceDecision.WickBeyondBucket != "gt_20bp" {
		t.Fatalf("bad resistance wick bucket/pct: %+v", resistanceDecision)
	}
}

func TestSRRejectionTimingDecisionTouchPierceClosedBackBuckets(t *testing.T) {
	aboveSupport := testCandle(0, 102, 103, 101.5, 102)
	decision, ok := newSRRejectionTimingDecision(
		[]Candle{aboveSupport},
		supportTimingRow(0, aboveSupport, true),
		SRBoundarySideSupport,
	)
	if !ok {
		t.Fatalf("expected support decision")
	}
	if decision.TouchedZone || decision.PiercedZone || !decision.ClosedBack ||
		decision.DecisionRejectionCandidate || decision.CloseLocation != srCloseLocationAboveZone ||
		decision.WickBeyondBucket != "0_5bp" {
		t.Fatalf("bad untouched support decision: %+v", decision)
	}

	inZoneBelowSupport := testCandle(0, 100, 101, 100.5, 99.5)
	decision, ok = newSRRejectionTimingDecision(
		[]Candle{inZoneBelowSupport},
		supportTimingRow(0, inZoneBelowSupport, true),
		SRBoundarySideSupport,
	)
	if !ok {
		t.Fatalf("expected support decision")
	}
	if !decision.TouchedZone || decision.PiercedZone || decision.ClosedBack ||
		decision.DecisionRejectionCandidate || decision.CloseLocation != srCloseLocationInZoneBelowLevel {
		t.Fatalf("bad touched support decision: %+v", decision)
	}
	if decision.StrengthBucket != "3" || decision.DistanceBucket != "5_10bp" {
		t.Fatalf("bad support buckets: %+v", decision)
	}

	if _, ok := newSRRejectionTimingDecision(nil, SRAuditRow{Index: 0}, SRBoundarySideSupport); ok {
		t.Fatalf("expected missing candle to reject decision")
	}
	if _, ok := newSRRejectionTimingDecision([]Candle{testCandle(0, 1, 1, 1, 1)}, SRAuditRow{Index: 0}, "bad"); ok {
		t.Fatalf("expected invalid side to reject decision")
	}
}

func TestRunSRRejectionTimingAuditFiltersDetectorActiveAndSkipsMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 97, 100.5),
		testCandle(2, 100.5, 104, 100, 103),
	}
	inactive := supportTimingRow(0, candles[0], false)
	active := supportTimingRow(1, candles[1], true)

	candidates, summary, err := RunSRRejectionTimingAudit(candles, []SRAuditRow{inactive, active}, SRRejectionTimingAuditConfig{
		HorizonsBars:       []int{1, 2},
		DetectorActiveOnly: true,
	})
	if err != nil {
		t.Fatalf("RunSRRejectionTimingAudit error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	if candidates[0].HorizonBars != 1 || candidates[0].CandidateCount != 1 {
		t.Fatalf("bad candidate after filter/skip: %+v", candidates[0])
	}
	if summary[0].CandidateCount != 1 || summary[0].TouchedCount != 1 || summary[0].PiercedCount != 1 {
		t.Fatalf("bad summary after filter/skip: %+v", summary[0])
	}
}

func TestSRRejectionTimingSummaryDenominatorsAndZeroDecisionRates(t *testing.T) {
	events := []srRejectionTimingEvent{
		{
			srRejectionTimingDecision: srRejectionTimingDecision{
				Split:                      "2021_2022_stress",
				Side:                       SRBoundarySideSupport,
				TouchedZone:                true,
				ClosedBack:                 true,
				DecisionRejectionCandidate: true,
				DetectorProfileID:          BalancedDetectorProfileID,
				DetectorActive:             true,
			},
			HorizonBars:                      3,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.04,
			LabelAdverseMovePct:              0.01,
		},
		{
			srRejectionTimingDecision: srRejectionTimingDecision{
				Split:             "2021_2022_stress",
				Side:              SRBoundarySideSupport,
				DetectorProfileID: BalancedDetectorProfileID,
				DetectorActive:    true,
			},
			HorizonBars:              3,
			LabelCloseBreak:          true,
			LabelReclaimedAfterBreak: true,
			LabelRejected:            true,
			LabelFavorableMovePct:    0.02,
			LabelAdverseMovePct:      0.03,
		},
		{
			srRejectionTimingDecision: srRejectionTimingDecision{
				Split:             "2021_2022_stress",
				Side:              SRBoundarySideSupport,
				DetectorProfileID: BalancedDetectorProfileID,
				DetectorActive:    true,
			},
			HorizonBars:           6,
			LabelRejected:         true,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRRejectionTimingSummary(events)
	if len(rows) != 2 {
		t.Fatalf("summary rows=%d, want 2", len(rows))
	}
	first := rows[0]
	if first.HorizonBars != 3 || first.CandidateCount != 2 ||
		first.TouchedCount != 1 || first.ClosedBackCount != 1 || first.DecisionRejectionCandidateCount != 1 {
		t.Fatalf("bad first summary counts: %+v", first)
	}
	if !boundaryAlmostEqual(first.TouchedRate, 0.5) ||
		!boundaryAlmostEqual(first.DecisionRejectionCandidateRate, 0.5) ||
		!boundaryAlmostEqual(first.LabelRejectionRate, 1) ||
		!boundaryAlmostEqual(first.LabelCloseBreakRate, 0.5) ||
		!boundaryAlmostEqual(first.LabelReclaimGivenCloseBreakRate, 1) {
		t.Fatalf("bad first summary rates: %+v", first)
	}
	if !boundaryAlmostEqual(first.LabelAvgFavorablePct, 0.03) ||
		!boundaryAlmostEqual(first.LabelAvgAdversePct, 0.02) ||
		!boundaryAlmostEqual(first.DecisionCandidateLabelAvgFavorablePct, 0.04) ||
		!boundaryAlmostEqual(first.DecisionCandidateLabelAvgAdversePct, 0.01) ||
		!boundaryAlmostEqual(first.DecisionCandidateLabelRejectionRate, 1) {
		t.Fatalf("bad first summary label averages: %+v", first)
	}

	second := rows[1]
	if second.HorizonBars != 6 || second.DecisionRejectionCandidateCount != 0 {
		t.Fatalf("bad second summary counts: %+v", second)
	}
	if second.DecisionCandidateLabelRejectionRate != 0 ||
		second.DecisionCandidateLabelAvgFavorablePct != 0 ||
		second.DecisionCandidateLabelFavorableGreaterThanAdverseRate != 0 {
		t.Fatalf("zero decision-candidate denominators should produce zero rates: %+v", second)
	}
}

func TestSRRejectionTimingCandidateAggregationAndInvalidConfig(t *testing.T) {
	events := []srRejectionTimingEvent{
		{
			srRejectionTimingDecision: srRejectionTimingDecision{
				Split:                      "2021_2022_stress",
				Side:                       SRBoundarySideSupport,
				CloseLocation:              srCloseLocationInZoneAboveLevel,
				TouchedZone:                true,
				ClosedBack:                 true,
				DecisionRejectionCandidate: true,
				WickBeyondBucket:           "0_5bp",
				StrengthBucket:             "2",
				DistanceBucket:             "0_5bp",
				DetectorProfileID:          BalancedDetectorProfileID,
				DetectorActive:             true,
				Score:                      2,
				DistancePct:                0.0004,
				WickBeyondPct:              0.0002,
			},
			HorizonBars:                      1,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.03,
			LabelAdverseMovePct:              0.01,
		},
		{
			srRejectionTimingDecision: srRejectionTimingDecision{
				Split:                      "2021_2022_stress",
				Side:                       SRBoundarySideSupport,
				CloseLocation:              srCloseLocationInZoneAboveLevel,
				TouchedZone:                true,
				ClosedBack:                 true,
				DecisionRejectionCandidate: true,
				WickBeyondBucket:           "0_5bp",
				StrengthBucket:             "2",
				DistanceBucket:             "0_5bp",
				DetectorProfileID:          BalancedDetectorProfileID,
				DetectorActive:             true,
				Score:                      4,
				DistancePct:                0.0006,
				WickBeyondPct:              0.0004,
			},
			HorizonBars:           1,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRRejectionTimingCandidates(events)
	if len(rows) != 1 {
		t.Fatalf("candidate rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.CandidateCount != 2 || row.LabelRejectedCount != 1 ||
		!boundaryAlmostEqual(row.AvgScore, 3) ||
		!boundaryAlmostEqual(row.AvgDistancePct, 0.0005) ||
		!boundaryAlmostEqual(row.AvgWickBeyondPct, 0.0003) ||
		!boundaryAlmostEqual(row.LabelRejectionRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelFavorableGreaterThanAdverseRate, 0.5) {
		t.Fatalf("bad candidate aggregation: %+v", row)
	}

	if _, _, err := RunSRRejectionTimingAudit(nil, nil, SRRejectionTimingAuditConfig{HorizonsBars: []int{1, 0}}); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
}

func supportTimingRow(index int, candle Candle, detectorActive bool) SRAuditRow {
	return SRAuditRow{
		Index:                     index,
		OpenTime:                  candle.OpenTime.Format(timeLayout),
		CloseTime:                 candle.CloseTime.Format(timeLayout),
		Split:                     "2021_2022_stress",
		Close:                     candle.Close,
		DetectorProfileID:         BalancedDetectorProfileID,
		DetectorRawActive:         detectorActive,
		DetectorActive:            detectorActive,
		HasSupport:                true,
		NearSupport:               true,
		NearestSupport:            100,
		NearestSupportTop:         101,
		NearestSupportBottom:      98,
		NearestSupportDistancePct: 0.0008,
		NearestSupportStrength:    3,
		NearestSupportScore:       4,
	}
}

func resistanceTimingRow(index int, candle Candle, detectorActive bool) SRAuditRow {
	return SRAuditRow{
		Index:                        index,
		OpenTime:                     candle.OpenTime.Format(timeLayout),
		CloseTime:                    candle.CloseTime.Format(timeLayout),
		Split:                        "2021_2022_stress",
		Close:                        candle.Close,
		DetectorProfileID:            BalancedDetectorProfileID,
		DetectorRawActive:            detectorActive,
		DetectorActive:               detectorActive,
		HasResistance:                true,
		NearResistance:               true,
		NearestResistance:            100,
		NearestResistanceTop:         102,
		NearestResistanceBottom:      99,
		NearestResistanceDistancePct: 0.0008,
		NearestResistanceStrength:    3,
		NearestResistanceScore:       4,
	}
}
