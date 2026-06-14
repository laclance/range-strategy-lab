package lab

import "testing"

func TestCompressionBreakoutEpisodesUseRawRunThatEventuallyBecomesActive(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 103, 97, 102),
		testCandle(3, 102, 104, 96, 103),
		testCandle(4, 103, 105, 95, 104),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, true, true},
		[]bool{false, true, false, false, false},
	)

	episodes := compressionBreakoutEpisodes(candles, classifications, []Split{{Name: "full_2021_2026"}}, BalancedDetectorProfileID)
	if len(episodes) != 1 {
		t.Fatalf("episodes=%d, want 1", len(episodes))
	}
	episode := episodes[0]
	if episode.StartIndex != 0 || episode.EndIndex != 1 ||
		episode.High != 102 || episode.Low != 98 ||
		episode.RawLengthBars != 2 || episode.ActiveLengthBars != 1 {
		t.Fatalf("bad episode: %+v", episode)
	}
}

func TestCompressionBreakoutDecisionUsesFrozenEpisodeBounds(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 200, 50, 103),
		testCandle(4, 103, 250, 20, 220),
	}
	episode := compressionBreakoutEpisode{
		Split:              "full_2021_2026",
		StartIndex:         0,
		EndIndex:           2,
		High:               102.5,
		Low:                97,
		EndClose:           102,
		RawLengthBars:      3,
		ActiveLengthBars:   1,
		RawLengthBucket:    "lt_12",
		ActiveLengthBucket: "lt_12",
		RangeWidthBucket:   "gt_50bp",
		DetectorProfileID:  BalancedDetectorProfileID,
	}

	decision, ok := newCompressionBreakoutDecision(candles, nil, episode, 12)
	if !ok {
		t.Fatalf("expected breakout decision")
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[4] = testCandle(4, 103, 500, 1, 2)
	changedDecision, ok := newCompressionBreakoutDecision(changedFuture, nil, episode, 12)
	if !ok {
		t.Fatalf("expected breakout decision after future change")
	}
	if decision != changedDecision {
		t.Fatalf("decision changed after post-breakout edit: got %+v want %+v", changedDecision, decision)
	}
	if decision.Side != CompressionBreakoutSideUp || decision.BreakoutIndex != 3 || decision.BreakoutDelayBars != 1 {
		t.Fatalf("bad breakout decision: %+v", decision)
	}
	if !boundaryAlmostEqual(decision.BreakoutMovePct, (103.0-102.5)/103.0) {
		t.Fatalf("bad breakout move: %+v", decision)
	}
}

func TestCompressionBreakoutLabelsStartAfterBreakout(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 200, 50, 103),
		testCandle(4, 103, 104, 101, 101.5),
	}
	decision := compressionBreakoutDecision{
		compressionBreakoutEpisode: compressionBreakoutEpisode{
			High: 102.5,
			Low:  97,
		},
		BreakoutIndex: 3,
		Side:          CompressionBreakoutSideUp,
		BreakoutClose: 103,
	}

	label, ok := newCompressionBreakoutLabel(candles, decision, 1)
	if !ok {
		t.Fatalf("expected label")
	}
	wantFavorable := (104.0 - 103.0) / 103.0
	wantAdverse := (103.0 - 101.0) / 103.0
	if !label.ReenteredRange || label.OppositeCloseBreak {
		t.Fatalf("bad range labels: %+v", label)
	}
	if !boundaryAlmostEqual(label.FavorableMovePct, wantFavorable) ||
		!boundaryAlmostEqual(label.AdverseMovePct, wantAdverse) {
		t.Fatalf("label used breakout candle or wrong future window: %+v", label)
	}
}

func TestCompressionBreakoutFirstCloseBreakAndSideSymmetry(t *testing.T) {
	upCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 120, 96, 102),
		testCandle(4, 102, 104, 101, 103),
	}
	upEpisode := compressionBreakoutEpisode{EndIndex: 2, High: 102.5, Low: 97}
	upDecision, ok := newCompressionBreakoutDecision(upCandles, nil, upEpisode, 12)
	if !ok {
		t.Fatalf("expected up breakout")
	}
	if upDecision.Side != CompressionBreakoutSideUp || upDecision.BreakoutIndex != 4 || upDecision.BreakoutDelayBars != 2 {
		t.Fatalf("bad up breakout: %+v", upDecision)
	}

	downCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 103, 80, 98),
		testCandle(4, 98, 99, 94, 96.5),
	}
	downEpisode := compressionBreakoutEpisode{EndIndex: 2, High: 102.5, Low: 97}
	downDecision, ok := newCompressionBreakoutDecision(downCandles, nil, downEpisode, 12)
	if !ok {
		t.Fatalf("expected down breakout")
	}
	if downDecision.Side != CompressionBreakoutSideDown || downDecision.BreakoutIndex != 4 || downDecision.BreakoutDelayBars != 2 {
		t.Fatalf("bad down breakout: %+v", downDecision)
	}
}

func TestRunCompressionBreakoutAuditSkipsMissingFutureAndInvalidConfig(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 104, 101, 103),
		testCandle(4, 103, 104, 100, 101),
		testCandle(5, 101, 102, 100, 101),
		testCandle(6, 101, 102, 100, 101),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, true, false, false, true, true},
		[]bool{false, false, true, false, false, false, true},
	)

	candidates, summary, err := runCompressionBreakoutAuditFromClassifications(candles, classifications, CompressionBreakoutAuditConfig{
		HorizonsBars:         []int{1, 4},
		MaxBreakoutDelayBars: 12,
		DetectorProfileID:    BalancedDetectorProfileID,
	}, []Split{{Name: "full_2021_2026"}})
	if err != nil {
		t.Fatalf("runCompressionBreakoutAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 || len(summary) != 1 {
		t.Fatalf("rows candidates=%d summary=%d, want 1/1", len(candidates), len(summary))
	}
	if candidates[0].HorizonBars != 1 || candidates[0].CandidateCount != 1 ||
		summary[0].HorizonBars != 1 || summary[0].CandidateCount != 1 {
		t.Fatalf("bad rows after skip: candidates=%+v summary=%+v", candidates[0], summary[0])
	}

	defaults := CompressionBreakoutAuditConfig{}.withDefaults()
	if len(defaults.HorizonsBars) != 4 || defaults.MaxBreakoutDelayBars != 12 ||
		defaults.DetectorProfileID != BalancedDetectorProfileID {
		t.Fatalf("bad defaults: %+v", defaults)
	}
	if _, _, err := runCompressionBreakoutAuditFromClassifications(nil, nil, CompressionBreakoutAuditConfig{
		HorizonsBars:         []int{0},
		MaxBreakoutDelayBars: 12,
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, err := runCompressionBreakoutAuditFromClassifications(nil, nil, CompressionBreakoutAuditConfig{
		HorizonsBars:         []int{1},
		MaxBreakoutDelayBars: -1,
	}, nil); err == nil {
		t.Fatalf("expected invalid max delay error")
	}
	if _, ok := newCompressionBreakoutLabel(candles, compressionBreakoutDecision{BreakoutIndex: 0, Side: "bad"}, 1); ok {
		t.Fatalf("expected invalid side label rejection")
	}
}

func TestRunCompressionBreakoutAuditSkipsNoBreakEpisodes(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 102.5, 97, 102),
		testCandle(3, 102, 102.4, 98, 101),
		testCandle(4, 101, 102.2, 98.5, 100),
		testCandle(5, 100, 102, 99, 101),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, true, false, false, false},
		[]bool{false, false, true, false, false, false},
	)

	candidates, summary, err := runCompressionBreakoutAuditFromClassifications(candles, classifications, CompressionBreakoutAuditConfig{
		HorizonsBars:         []int{1},
		MaxBreakoutDelayBars: 3,
		DetectorProfileID:    BalancedDetectorProfileID,
	}, []Split{{Name: "full_2021_2026"}})
	if err != nil {
		t.Fatalf("runCompressionBreakoutAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 0 || len(summary) != 0 {
		t.Fatalf("expected no rows for no-break episode: candidates=%+v summary=%+v", candidates, summary)
	}
}

func TestCompressionBreakoutCandidateAggregation(t *testing.T) {
	events := []compressionBreakoutEvent{
		{
			compressionBreakoutDecision: compressionBreakoutDecision{
				compressionBreakoutEpisode: compressionBreakoutEpisode{
					Split:              "2021_2022_stress",
					RawLengthBars:      12,
					ActiveLengthBars:   1,
					RawLengthBucket:    "12_23",
					ActiveLengthBucket: "lt_12",
					RangeWidthPct:      0.002,
					RangeWidthBucket:   "10_25bp",
					DetectorProfileID:  BalancedDetectorProfileID,
				},
				Side:                             CompressionBreakoutSideUp,
				BreakoutDelayBars:                1,
				BreakoutMovePct:                  0.001,
				BreakoutMoveBucket:               "5_10bp",
				DecisionTrueRangeATR:             1.2,
				DecisionTrueRangeExpansionBucket: "1_1_5x",
			},
			HorizonBars:                      3,
			LabelReenteredRange:              true,
			LabelFavorableGreaterThanAdverse: true,
			LabelFavorableMovePct:            0.03,
			LabelAdverseMovePct:              0.01,
		},
		{
			compressionBreakoutDecision: compressionBreakoutDecision{
				compressionBreakoutEpisode: compressionBreakoutEpisode{
					Split:              "2021_2022_stress",
					RawLengthBars:      24,
					ActiveLengthBars:   2,
					RawLengthBucket:    "12_23",
					ActiveLengthBucket: "lt_12",
					RangeWidthPct:      0.004,
					RangeWidthBucket:   "10_25bp",
					DetectorProfileID:  BalancedDetectorProfileID,
				},
				Side:                             CompressionBreakoutSideUp,
				BreakoutDelayBars:                1,
				BreakoutMovePct:                  0.003,
				BreakoutMoveBucket:               "5_10bp",
				DecisionTrueRangeATR:             1.6,
				DecisionTrueRangeExpansionBucket: "1_1_5x",
			},
			HorizonBars:             3,
			LabelOppositeCloseBreak: true,
			LabelFavorableMovePct:   0.01,
			LabelAdverseMovePct:     0.02,
		},
	}

	candidates := summarizeCompressionBreakoutCandidates(events)
	if len(candidates) != 1 {
		t.Fatalf("candidate rows=%d, want 1", len(candidates))
	}
	row := candidates[0]
	if row.CandidateCount != 2 ||
		row.LabelReenteredRangeCount != 1 ||
		row.LabelOppositeCloseBreakCount != 1 ||
		row.LabelFavorableGreaterThanAdverseCount != 1 ||
		!boundaryAlmostEqual(row.AvgEpisodeRawLengthBars, 18) ||
		!boundaryAlmostEqual(row.AvgEpisodeActiveLengthBars, 1.5) ||
		!boundaryAlmostEqual(row.AvgEpisodeRangeWidthPct, 0.003) ||
		!boundaryAlmostEqual(row.AvgBreakoutMovePct, 0.002) ||
		!boundaryAlmostEqual(row.AvgDecisionTrueRangeATR, 1.4) ||
		!boundaryAlmostEqual(row.LabelReenteredRangeRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelOppositeCloseBreakRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelFavorableMinusAdversePct, 0.005) ||
		!boundaryAlmostEqual(row.LabelFavorableGreaterThanAdverseRate, 0.5) {
		t.Fatalf("bad candidate aggregation: %+v", row)
	}

	summary := summarizeCompressionBreakoutSummary(events)
	if len(summary) != 1 {
		t.Fatalf("summary rows=%d, want 1", len(summary))
	}
	if !boundaryAlmostEqual(summary[0].AvgBreakoutDelayBars, 1) ||
		!boundaryAlmostEqual(summary[0].AvgDecisionTrueRangeATR, 1.4) {
		t.Fatalf("bad summary aggregation: %+v", summary[0])
	}
}

func TestCompressionBreakoutRowsSortDeterministically(t *testing.T) {
	events := []compressionBreakoutEvent{
		compressionBreakoutSortEvent("full_2021_2026", CompressionBreakoutSideDown, 12, 2),
		compressionBreakoutSortEvent("2021_2022_stress", CompressionBreakoutSideUp, 3, 2),
		compressionBreakoutSortEvent("2021_2022_stress", CompressionBreakoutSideUp, 1, 1),
	}

	rows := summarizeCompressionBreakoutCandidates(events)
	if len(rows) != 3 {
		t.Fatalf("rows=%d, want 3", len(rows))
	}
	want := []CompressionBreakoutCandidateRow{
		{Split: "2021_2022_stress", Side: CompressionBreakoutSideUp, HorizonBars: 1, BreakoutDelayBars: 1},
		{Split: "2021_2022_stress", Side: CompressionBreakoutSideUp, HorizonBars: 3, BreakoutDelayBars: 2},
		{Split: "full_2021_2026", Side: CompressionBreakoutSideDown, HorizonBars: 12, BreakoutDelayBars: 2},
	}
	for i := range want {
		if rows[i].Split != want[i].Split ||
			rows[i].Side != want[i].Side ||
			rows[i].HorizonBars != want[i].HorizonBars ||
			rows[i].BreakoutDelayBars != want[i].BreakoutDelayBars {
			t.Fatalf("row %d=%+v, want key %+v", i, rows[i], want[i])
		}
	}
}

func compressionBreakoutSortEvent(split, side string, horizon, delay int) compressionBreakoutEvent {
	return compressionBreakoutEvent{
		compressionBreakoutDecision: compressionBreakoutDecision{
			compressionBreakoutEpisode: compressionBreakoutEpisode{
				Split:              split,
				RawLengthBucket:    "12_23",
				ActiveLengthBucket: "lt_12",
				RangeWidthBucket:   "10_25bp",
				DetectorProfileID:  BalancedDetectorProfileID,
			},
			Side:                             side,
			BreakoutDelayBars:                delay,
			BreakoutMoveBucket:               "5_10bp",
			DecisionTrueRangeExpansionBucket: "1_1_5x",
		},
		HorizonBars: horizon,
	}
}

func testCompressionClassifications(raw, active []bool) []RangeClassification {
	out := make([]RangeClassification, len(raw))
	for i := range raw {
		out[i] = RangeClassification{Index: i, RawActive: raw[i], Active: active[i]}
	}
	return out
}
