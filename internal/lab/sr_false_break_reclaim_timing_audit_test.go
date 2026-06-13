package lab

import "testing"

func TestSRFalseBreakReclaimTimingDecisionUsesReclaimCandleOnly(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 100.5, 97, 97.5),
		testCandle(2, 97.5, 101, 97, 100.5),
		testCandle(3, 100.5, 104, 99, 103),
	}
	row := supportTimingRow(0, candles[0], true)

	decision, _, ok := newSRFalseBreakReclaimTimingDecision(candles, row, SRBoundarySideSupport, 3, 12)
	if !ok {
		t.Fatalf("expected support false-break reclaim decision")
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[3] = testCandle(3, 100.5, 150, 50, 51)
	changedDecision, _, ok := newSRFalseBreakReclaimTimingDecision(changedFuture, row, SRBoundarySideSupport, 3, 12)
	if !ok {
		t.Fatalf("expected support false-break reclaim decision after future change")
	}
	if decision != changedDecision {
		t.Fatalf("decision features changed after post-reclaim edit: got %+v want %+v", changedDecision, decision)
	}

	if decision.BreakDelayBars != 1 || decision.ReclaimDelayBars != 1 || decision.TotalDelayBars != 2 {
		t.Fatalf("bad delays: %+v", decision)
	}
	if decision.AnchorCloseLocation != srCloseLocationInZoneAboveLevel ||
		decision.BreakCloseLocation != srCloseLocationBelowZone ||
		decision.ReclaimCloseLocation != srCloseLocationInZoneAboveLevel {
		t.Fatalf("bad close locations: %+v", decision)
	}
	if !decision.DecisionFalseBreakReclaimCandidate {
		t.Fatalf("expected decision candidate: %+v", decision)
	}
	if !boundaryAlmostEqual(decision.BreakMovePct, (98.0-97.5)/100.0) ||
		!boundaryAlmostEqual(decision.ReclaimMovePct, (100.5-97.5)/100.0) {
		t.Fatalf("bad move pcts: %+v", decision)
	}
}

func TestSRFalseBreakReclaimTimingLabelsStartAfterReclaim(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 100.5, 97, 97.5),
		testCandle(2, 97.5, 150, 50, 100.5),
		testCandle(3, 100.5, 104, 99, 103),
	}
	row := supportTimingRow(0, candles[0], true)

	candidates, summary, err := RunSRFalseBreakReclaimTimingAudit(candles, []SRAuditRow{row}, SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        []int{1},
		MaxBreakDelayBars:   3,
		MaxReclaimDelayBars: 12,
		DetectorActiveOnly:  true,
	})
	if err != nil {
		t.Fatalf("RunSRFalseBreakReclaimTimingAudit error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	wantFavorable := (104.0 - 100.5) / 100.5
	wantAdverse := (100.5 - 99.0) / 100.5
	if !boundaryAlmostEqual(candidates[0].LabelAvgFavorablePct, wantFavorable) ||
		!boundaryAlmostEqual(candidates[0].LabelAvgAdversePct, wantAdverse) {
		t.Fatalf("candidate labels should start after reclaim candle: %+v", candidates[0])
	}
	if !boundaryAlmostEqual(summary[0].LabelDecisionCandidateAvgFavorablePct, wantFavorable) ||
		!boundaryAlmostEqual(summary[0].LabelDecisionCandidateAvgAdversePct, wantAdverse) {
		t.Fatalf("summary labels should start after reclaim candle: %+v", summary[0])
	}
}

func TestSRFalseBreakReclaimTimingSupportResistanceSymmetry(t *testing.T) {
	supportCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 100.5, 97, 97.5),
		testCandle(2, 97.5, 101, 97, 100.5),
	}
	supportDecision, _, ok := newSRFalseBreakReclaimTimingDecision(
		supportCandles,
		supportTimingRow(0, supportCandles[0], true),
		SRBoundarySideSupport,
		3,
		12,
	)
	if !ok {
		t.Fatalf("expected support false-break reclaim decision")
	}
	if supportDecision.BreakCloseLocation != srCloseLocationBelowZone ||
		supportDecision.ReclaimCloseLocation != srCloseLocationInZoneAboveLevel ||
		!boundaryAlmostEqual(supportDecision.BreakMovePct, 0.005) ||
		!boundaryAlmostEqual(supportDecision.ReclaimMovePct, 0.03) {
		t.Fatalf("bad support symmetry decision: %+v", supportDecision)
	}

	resistanceCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 103, 99, 102.5),
		testCandle(2, 102.5, 103, 99, 99.5),
	}
	resistanceDecision, _, ok := newSRFalseBreakReclaimTimingDecision(
		resistanceCandles,
		resistanceTimingRow(0, resistanceCandles[0], true),
		SRBoundarySideResistance,
		3,
		12,
	)
	if !ok {
		t.Fatalf("expected resistance false-break reclaim decision")
	}
	if resistanceDecision.BreakCloseLocation != srCloseLocationAboveZone ||
		resistanceDecision.ReclaimCloseLocation != srCloseLocationInZoneBelowLevel ||
		!boundaryAlmostEqual(resistanceDecision.BreakMovePct, 0.005) ||
		!boundaryAlmostEqual(resistanceDecision.ReclaimMovePct, 0.03) {
		t.Fatalf("bad resistance symmetry decision: %+v", resistanceDecision)
	}
}

func TestRunSRFalseBreakReclaimTimingAuditFiltersSkipsAndValidatesConfig(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 100.5, 97, 97.5),
		testCandle(2, 97.5, 101, 97, 100.5),
		testCandle(3, 100.5, 104, 99, 103),
	}
	active := supportTimingRow(0, candles[0], true)
	inactive := supportTimingRow(0, candles[0], false)
	noBreak := supportTimingRow(0, candles[0], true)
	noBreak.NearestSupportBottom = 97
	noReclaim := supportTimingRow(0, candles[0], true)
	noReclaim.NearestSupport = 101
	tooLate := supportTimingRow(1, candles[1], true)

	candidates, summary, err := RunSRFalseBreakReclaimTimingAudit(candles, []SRAuditRow{inactive, active, noBreak, noReclaim, tooLate}, SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        []int{1, 2},
		MaxBreakDelayBars:   3,
		MaxReclaimDelayBars: 12,
		DetectorActiveOnly:  true,
	})
	if err != nil {
		t.Fatalf("RunSRFalseBreakReclaimTimingAudit error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	if candidates[0].HorizonBars != 1 || candidates[0].CandidateCount != 1 ||
		summary[0].CandidateCount != 1 || summary[0].DecisionFalseBreakReclaimCandidateCount != 1 {
		t.Fatalf("bad rows after filter/skip: candidates=%+v summary=%+v", candidates[0], summary[0])
	}

	defaults := SRFalseBreakReclaimTimingAuditConfig{}.withDefaults()
	if len(defaults.HorizonsBars) != 4 || defaults.MaxBreakDelayBars != 3 ||
		defaults.MaxReclaimDelayBars != 12 || !defaults.DetectorActiveOnly {
		t.Fatalf("bad defaults: %+v", defaults)
	}
	if _, _, err := RunSRFalseBreakReclaimTimingAudit(nil, nil, SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        []int{0},
		MaxBreakDelayBars:   3,
		MaxReclaimDelayBars: 12,
	}); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, err := RunSRFalseBreakReclaimTimingAudit(nil, nil, SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        []int{1},
		MaxBreakDelayBars:   -1,
		MaxReclaimDelayBars: 12,
	}); err == nil {
		t.Fatalf("expected invalid break delay error")
	}
	if _, _, err := RunSRFalseBreakReclaimTimingAudit(nil, nil, SRFalseBreakReclaimTimingAuditConfig{
		HorizonsBars:        []int{1},
		MaxBreakDelayBars:   3,
		MaxReclaimDelayBars: -1,
	}); err == nil {
		t.Fatalf("expected invalid reclaim delay error")
	}
	if _, _, ok := newSRFalseBreakReclaimTimingDecision(candles, active, "bad", 3, 12); ok {
		t.Fatalf("expected invalid side to reject decision")
	}
}

func TestSRFalseBreakReclaimTimingCandidateAggregation(t *testing.T) {
	events := []srFalseBreakReclaimTimingEvent{
		{
			srFalseBreakReclaimTimingDecision: srFalseBreakReclaimTimingDecision{
				Split:                              "2021_2022_stress",
				Side:                               SRBoundarySideSupport,
				BreakDelayBars:                     1,
				ReclaimDelayBars:                   2,
				TotalDelayBars:                     3,
				AnchorCloseLocation:                srCloseLocationInZoneAboveLevel,
				BreakCloseLocation:                 srCloseLocationBelowZone,
				ReclaimCloseLocation:               srCloseLocationInZoneAboveLevel,
				BreakMoveBucket:                    "0_5bp",
				ReclaimMoveBucket:                  "gt_20bp",
				DecisionFalseBreakReclaimCandidate: true,
				StrengthBucket:                     "2",
				DistanceBucket:                     "0_5bp",
				DetectorProfileID:                  BalancedDetectorProfileID,
				DetectorActive:                     true,
				Score:                              2,
				DistancePct:                        0.0004,
				BreakMovePct:                       0.002,
				ReclaimMovePct:                     0.003,
			},
			HorizonBars:                      1,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.03,
			LabelAdverseMovePct:              0.01,
		},
		{
			srFalseBreakReclaimTimingDecision: srFalseBreakReclaimTimingDecision{
				Split:                              "2021_2022_stress",
				Side:                               SRBoundarySideSupport,
				BreakDelayBars:                     1,
				ReclaimDelayBars:                   2,
				TotalDelayBars:                     3,
				AnchorCloseLocation:                srCloseLocationInZoneAboveLevel,
				BreakCloseLocation:                 srCloseLocationBelowZone,
				ReclaimCloseLocation:               srCloseLocationInZoneAboveLevel,
				BreakMoveBucket:                    "0_5bp",
				ReclaimMoveBucket:                  "gt_20bp",
				DecisionFalseBreakReclaimCandidate: true,
				StrengthBucket:                     "2",
				DistanceBucket:                     "0_5bp",
				DetectorProfileID:                  BalancedDetectorProfileID,
				DetectorActive:                     true,
				Score:                              4,
				DistancePct:                        0.0006,
				BreakMovePct:                       0.004,
				ReclaimMovePct:                     0.005,
			},
			HorizonBars:           1,
			LabelCloseBreak:       true,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRFalseBreakReclaimTimingCandidates(events)
	if len(rows) != 1 {
		t.Fatalf("candidate rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.CandidateCount != 2 || row.LabelRejectedCount != 1 ||
		!boundaryAlmostEqual(row.AvgScore, 3) ||
		!boundaryAlmostEqual(row.AvgDistancePct, 0.0005) ||
		!boundaryAlmostEqual(row.AvgBreakMovePct, 0.003) ||
		!boundaryAlmostEqual(row.AvgReclaimMovePct, 0.004) ||
		!boundaryAlmostEqual(row.LabelCloseBreakRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelRejectionRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelFavorableGreaterThanAdverseRate, 0.5) {
		t.Fatalf("bad candidate aggregation: %+v", row)
	}
}

func TestSRFalseBreakReclaimTimingSummaryDenominatorsAndZeroDecisionRates(t *testing.T) {
	events := []srFalseBreakReclaimTimingEvent{
		{
			srFalseBreakReclaimTimingDecision: srFalseBreakReclaimTimingDecision{
				Split:                              "2021_2022_stress",
				Side:                               SRBoundarySideSupport,
				BreakDelayBars:                     1,
				ReclaimDelayBars:                   2,
				TotalDelayBars:                     3,
				DecisionFalseBreakReclaimCandidate: true,
				DetectorProfileID:                  BalancedDetectorProfileID,
				DetectorActive:                     true,
				BreakMovePct:                       0.002,
				ReclaimMovePct:                     0.004,
			},
			HorizonBars:                      3,
			LabelRejected:                    true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.04,
			LabelAdverseMovePct:              0.01,
		},
		{
			srFalseBreakReclaimTimingDecision: srFalseBreakReclaimTimingDecision{
				Split:             "2021_2022_stress",
				Side:              SRBoundarySideSupport,
				BreakDelayBars:    2,
				ReclaimDelayBars:  3,
				TotalDelayBars:    5,
				DetectorProfileID: BalancedDetectorProfileID,
				DetectorActive:    true,
				BreakMovePct:      0.004,
				ReclaimMovePct:    0.006,
			},
			HorizonBars:              3,
			LabelCloseBreak:          true,
			LabelReclaimedAfterBreak: true,
			LabelFavorableMovePct:    0.02,
			LabelAdverseMovePct:      0.03,
		},
		{
			srFalseBreakReclaimTimingDecision: srFalseBreakReclaimTimingDecision{
				Split:             "2021_2022_stress",
				Side:              SRBoundarySideSupport,
				BreakDelayBars:    1,
				ReclaimDelayBars:  1,
				TotalDelayBars:    2,
				DetectorProfileID: BalancedDetectorProfileID,
				DetectorActive:    true,
			},
			HorizonBars:           6,
			LabelRejected:         true,
			LabelFavorableMovePct: 0.01,
			LabelAdverseMovePct:   0.02,
		},
	}

	rows := summarizeSRFalseBreakReclaimTimingSummary(events)
	if len(rows) != 2 {
		t.Fatalf("summary rows=%d, want 2", len(rows))
	}
	first := rows[0]
	if first.HorizonBars != 3 || first.CandidateCount != 2 ||
		first.DecisionFalseBreakReclaimCandidateCount != 1 {
		t.Fatalf("bad first summary counts: %+v", first)
	}
	if !boundaryAlmostEqual(first.DecisionFalseBreakReclaimCandidateRate, 0.5) ||
		!boundaryAlmostEqual(first.AvgBreakDelayBars, 1.5) ||
		!boundaryAlmostEqual(first.AvgReclaimDelayBars, 2.5) ||
		!boundaryAlmostEqual(first.AvgTotalDelayBars, 4) ||
		!boundaryAlmostEqual(first.AvgBreakMovePct, 0.003) ||
		!boundaryAlmostEqual(first.AvgReclaimMovePct, 0.005) {
		t.Fatalf("bad first summary averages: %+v", first)
	}
	if !boundaryAlmostEqual(first.LabelCloseBreakRate, 0.5) ||
		!boundaryAlmostEqual(first.LabelReclaimGivenCloseBreakRate, 1) ||
		!boundaryAlmostEqual(first.LabelAvgFavorablePct, 0.03) ||
		!boundaryAlmostEqual(first.LabelAvgAdversePct, 0.02) ||
		!boundaryAlmostEqual(first.LabelDecisionCandidateAvgFavorablePct, 0.04) ||
		!boundaryAlmostEqual(first.LabelDecisionCandidateAvgAdversePct, 0.01) ||
		!boundaryAlmostEqual(first.LabelDecisionCandidateRejectionRate, 1) {
		t.Fatalf("bad first summary label rates: %+v", first)
	}

	second := rows[1]
	if second.HorizonBars != 6 || second.DecisionFalseBreakReclaimCandidateCount != 0 {
		t.Fatalf("bad second summary counts: %+v", second)
	}
	if second.LabelDecisionCandidateRejectionRate != 0 ||
		second.LabelDecisionCandidateAvgFavorablePct != 0 ||
		second.LabelDecisionCandidateFavorableGreaterThanAdverseRate != 0 {
		t.Fatalf("zero decision-candidate denominators should produce zero rates: %+v", second)
	}
}
