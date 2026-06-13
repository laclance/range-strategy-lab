package lab

import "testing"

func TestSRConfirmationTimingDecisionUsesConfirmationCandleOnly(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 105, 97, 100.5),
		testCandle(1, 100.5, 102, 100, 101),
		testCandle(2, 101, 103, 100, 102),
	}
	row := supportTimingRow(0, candles[0], true)

	decision, _, ok := newSRConfirmationTimingDecision(candles, row, SRBoundarySideSupport, 1)
	if !ok {
		t.Fatalf("expected support confirmation decision")
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[2] = testCandle(2, 101, 150, 50, 51)
	changedDecision, _, ok := newSRConfirmationTimingDecision(changedFuture, row, SRBoundarySideSupport, 1)
	if !ok {
		t.Fatalf("expected support confirmation decision after future change")
	}
	if decision != changedDecision {
		t.Fatalf("decision features changed after post-confirmation edit: got %+v want %+v", changedDecision, decision)
	}
	if !decision.ConfirmationFavorableClose || decision.ConfirmationWrongSideClose || !decision.DecisionConfirmationCandidate {
		t.Fatalf("bad support confirmation flags: %+v", decision)
	}
	if decision.ConfirmationCloseLocation != srCloseLocationInZoneAboveLevel {
		t.Fatalf("confirmation close location=%q", decision.ConfirmationCloseLocation)
	}
	if !boundaryAlmostEqual(decision.ConfirmationMovePct, (101.0-100.5)/100.5) {
		t.Fatalf("confirmation move pct=%v", decision.ConfirmationMovePct)
	}
}

func TestSRConfirmationTimingLabelsStartAfterConfirmation(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 105, 97, 100.5),
		testCandle(1, 100.5, 150, 50, 101),
		testCandle(2, 101, 104, 100, 103),
	}
	row := supportTimingRow(0, candles[0], true)

	candidates, summary, err := RunSRConfirmationTimingAudit(candles, []SRAuditRow{row}, SRConfirmationTimingAuditConfig{
		HorizonsBars:          []int{1},
		ConfirmationDelayBars: []int{1},
		DetectorActiveOnly:    true,
	})
	if err != nil {
		t.Fatalf("RunSRConfirmationTimingAudit error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	wantFavorable := (104.0 - 101.0) / 101.0
	wantAdverse := (101.0 - 100.0) / 101.0
	if !boundaryAlmostEqual(candidates[0].LabelAvgFavorablePct, wantFavorable) ||
		!boundaryAlmostEqual(candidates[0].LabelAvgAdversePct, wantAdverse) {
		t.Fatalf("labels should start after confirmation candle: %+v", candidates[0])
	}
	if !boundaryAlmostEqual(summary[0].DecisionCandidateLabelAvgFavorablePct, wantFavorable) ||
		!boundaryAlmostEqual(summary[0].DecisionCandidateLabelAvgAdversePct, wantAdverse) {
		t.Fatalf("decision labels should start after confirmation candle: %+v", summary[0])
	}
}

func TestSRConfirmationTimingSupportResistanceSymmetry(t *testing.T) {
	supportCandles := []Candle{
		testCandle(0, 100, 103, 97, 100.5),
		testCandle(1, 100.5, 102, 100, 101.5),
	}
	supportDecision, _, ok := newSRConfirmationTimingDecision(
		supportCandles,
		supportTimingRow(0, supportCandles[0], true),
		SRBoundarySideSupport,
		1,
	)
	if !ok {
		t.Fatalf("expected support confirmation decision")
	}
	if !supportDecision.ConfirmationFavorableClose || supportDecision.ConfirmationWrongSideClose ||
		!supportDecision.DecisionConfirmationCandidate || supportDecision.ConfirmationCloseLocation != srCloseLocationAboveZone {
		t.Fatalf("bad support confirmation decision: %+v", supportDecision)
	}

	resistanceCandles := []Candle{
		testCandle(0, 100, 103, 97, 99.5),
		testCandle(1, 99.5, 100, 98, 98.5),
	}
	resistanceDecision, _, ok := newSRConfirmationTimingDecision(
		resistanceCandles,
		resistanceTimingRow(0, resistanceCandles[0], true),
		SRBoundarySideResistance,
		1,
	)
	if !ok {
		t.Fatalf("expected resistance confirmation decision")
	}
	if !resistanceDecision.ConfirmationFavorableClose || resistanceDecision.ConfirmationWrongSideClose ||
		!resistanceDecision.DecisionConfirmationCandidate || resistanceDecision.ConfirmationCloseLocation != srCloseLocationBelowZone {
		t.Fatalf("bad resistance confirmation decision: %+v", resistanceDecision)
	}
	if !boundaryAlmostEqual(resistanceDecision.ConfirmationMovePct, (99.5-98.5)/99.5) {
		t.Fatalf("bad resistance move pct: %+v", resistanceDecision)
	}
}

func TestRunSRConfirmationTimingAuditFiltersSkipsAndRequiresSeedRejection(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 105, 97, 100.5),
		testCandle(1, 100.5, 103, 100, 101),
		testCandle(2, 101, 104, 100, 102),
		testCandle(3, 102, 103, 101, 102.5),
	}
	active := supportTimingRow(0, candles[0], true)
	inactive := supportTimingRow(0, candles[0], false)
	tooLate := supportTimingRow(2, candles[2], true)
	notRejection := supportTimingRow(1, candles[1], true)
	notRejection.NearestSupport = 98.5
	notRejection.NearestSupportTop = 99
	notRejection.NearestSupportBottom = 98

	candidates, summary, err := RunSRConfirmationTimingAudit(candles, []SRAuditRow{inactive, active, tooLate, notRejection}, SRConfirmationTimingAuditConfig{
		HorizonsBars:          []int{1},
		ConfirmationDelayBars: []int{1},
		DetectorActiveOnly:    true,
	})
	if err != nil {
		t.Fatalf("RunSRConfirmationTimingAudit error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	if candidates[0].CandidateCount != 1 || summary[0].CandidateCount != 1 {
		t.Fatalf("bad counts after filter/skip: candidates=%+v summary=%+v", candidates[0], summary[0])
	}
}

func TestSRConfirmationTimingSummaryDenominatorsAndZeroDecisionRates(t *testing.T) {
	events := []srConfirmationTimingEvent{
		{
			srConfirmationTimingDecision: srConfirmationTimingDecision{
				Split:                         "2021_2022_stress",
				Side:                          SRBoundarySideSupport,
				ConfirmationDelayBars:         1,
				ConfirmationFavorableClose:    true,
				DecisionConfirmationCandidate: true,
				DetectorProfileID:             BalancedDetectorProfileID,
				DetectorActive:                true,
			},
			HorizonBars:                      3,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.04,
			LabelAdverseMovePct:              0.01,
		},
		{
			srConfirmationTimingDecision: srConfirmationTimingDecision{
				Split:                      "2021_2022_stress",
				Side:                       SRBoundarySideSupport,
				ConfirmationDelayBars:      1,
				ConfirmationFavorableClose: true,
				ConfirmationWrongSideClose: true,
				DetectorProfileID:          BalancedDetectorProfileID,
				DetectorActive:             true,
			},
			HorizonBars:              3,
			LabelCloseBreak:          true,
			LabelReclaimedAfterBreak: true,
			LabelFavorableMovePct:    0.02,
			LabelAdverseMovePct:      0.03,
		},
		{
			srConfirmationTimingDecision: srConfirmationTimingDecision{
				Split:                 "2021_2022_stress",
				Side:                  SRBoundarySideSupport,
				ConfirmationDelayBars: 1,
				DetectorProfileID:     BalancedDetectorProfileID,
				DetectorActive:        true,
			},
			HorizonBars:           6,
			LabelRejected:         true,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRConfirmationTimingSummary(events)
	if len(rows) != 2 {
		t.Fatalf("summary rows=%d, want 2", len(rows))
	}
	first := rows[0]
	if first.HorizonBars != 3 || first.CandidateCount != 2 ||
		first.ConfirmationFavorableCloseCount != 2 || first.ConfirmationWrongSideCloseCount != 1 ||
		first.DecisionConfirmationCandidateCount != 1 {
		t.Fatalf("bad first summary counts: %+v", first)
	}
	if !boundaryAlmostEqual(first.ConfirmationFavorableCloseRate, 1) ||
		!boundaryAlmostEqual(first.ConfirmationWrongSideCloseRate, 0.5) ||
		!boundaryAlmostEqual(first.DecisionConfirmationCandidateRate, 0.5) ||
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
	if second.HorizonBars != 6 || second.DecisionConfirmationCandidateCount != 0 {
		t.Fatalf("bad second summary counts: %+v", second)
	}
	if second.DecisionCandidateLabelRejectionRate != 0 ||
		second.DecisionCandidateLabelAvgFavorablePct != 0 ||
		second.DecisionCandidateLabelFavorableGreaterThanAdverseRate != 0 {
		t.Fatalf("zero decision-candidate denominators should produce zero rates: %+v", second)
	}
}

func TestSRConfirmationTimingCandidateAggregationAndInvalidConfig(t *testing.T) {
	events := []srConfirmationTimingEvent{
		{
			srConfirmationTimingDecision: srConfirmationTimingDecision{
				Split:                         "2021_2022_stress",
				Side:                          SRBoundarySideSupport,
				ConfirmationDelayBars:         1,
				SeedCloseLocation:             srCloseLocationInZoneAboveLevel,
				SeedPiercedZone:               true,
				SeedWickBeyondBucket:          "0_5bp",
				ConfirmationCloseLocation:     srCloseLocationAboveZone,
				ConfirmationFavorableClose:    true,
				DecisionConfirmationCandidate: true,
				StrengthBucket:                "2",
				DistanceBucket:                "0_5bp",
				DetectorProfileID:             BalancedDetectorProfileID,
				DetectorActive:                true,
				Score:                         2,
				DistancePct:                   0.0004,
				SeedWickBeyondPct:             0.0002,
				ConfirmationMovePct:           0.001,
			},
			HorizonBars:                      1,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.03,
			LabelAdverseMovePct:              0.01,
		},
		{
			srConfirmationTimingDecision: srConfirmationTimingDecision{
				Split:                         "2021_2022_stress",
				Side:                          SRBoundarySideSupport,
				ConfirmationDelayBars:         1,
				SeedCloseLocation:             srCloseLocationInZoneAboveLevel,
				SeedPiercedZone:               true,
				SeedWickBeyondBucket:          "0_5bp",
				ConfirmationCloseLocation:     srCloseLocationAboveZone,
				ConfirmationFavorableClose:    true,
				DecisionConfirmationCandidate: true,
				StrengthBucket:                "2",
				DistanceBucket:                "0_5bp",
				DetectorProfileID:             BalancedDetectorProfileID,
				DetectorActive:                true,
				Score:                         4,
				DistancePct:                   0.0006,
				SeedWickBeyondPct:             0.0004,
				ConfirmationMovePct:           0.003,
			},
			HorizonBars:           1,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRConfirmationTimingCandidates(events)
	if len(rows) != 1 {
		t.Fatalf("candidate rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.CandidateCount != 2 || row.LabelRejectedCount != 1 ||
		!boundaryAlmostEqual(row.AvgScore, 3) ||
		!boundaryAlmostEqual(row.AvgDistancePct, 0.0005) ||
		!boundaryAlmostEqual(row.AvgSeedWickBeyondPct, 0.0003) ||
		!boundaryAlmostEqual(row.AvgConfirmationMovePct, 0.002) ||
		!boundaryAlmostEqual(row.LabelRejectionRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelFavorableGreaterThanAdverseRate, 0.5) {
		t.Fatalf("bad candidate aggregation: %+v", row)
	}

	if _, _, err := RunSRConfirmationTimingAudit(nil, nil, SRConfirmationTimingAuditConfig{HorizonsBars: []int{1, 0}, ConfirmationDelayBars: []int{1}}); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, err := RunSRConfirmationTimingAudit(nil, nil, SRConfirmationTimingAuditConfig{HorizonsBars: []int{1}, ConfirmationDelayBars: []int{0}}); err == nil {
		t.Fatalf("expected invalid confirmation delay error")
	}
}
