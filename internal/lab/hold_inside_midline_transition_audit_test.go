package lab

import "testing"

func TestHoldInsideMidlineTransitionDefaults(t *testing.T) {
	cfg := DefaultHoldInsideMidlineTransitionAuditConfig()
	if len(cfg.HorizonsBars) != 4 ||
		cfg.HorizonsBars[0] != 1 ||
		cfg.HorizonsBars[1] != 3 ||
		cfg.HorizonsBars[2] != 6 ||
		cfg.HorizonsBars[3] != 12 ||
		cfg.QuickInvalidationBars != 3 {
		t.Fatalf("bad default config: %+v", cfg)
	}
	if len(cfg.Profiles) != 1 || cfg.Profiles[0].ProfileID != BalancedDetectorProfileID {
		t.Fatalf("bad default profiles: %+v", cfg.Profiles)
	}
	wantRules := []DetectorContextRefinementRule{
		{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
		{RuleID: DetectorContextRuleHold6Inside, HoldBars: 6},
		{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
	}
	if len(cfg.ContextRules) != len(wantRules) {
		t.Fatalf("rules=%d, want %d", len(cfg.ContextRules), len(wantRules))
	}
	for i, want := range wantRules {
		if cfg.ContextRules[i] != want {
			t.Fatalf("rule %d=%+v, want %+v", i, cfg.ContextRules[i], want)
		}
	}
}

func TestHoldInsideMidlineTransitionHoldRuleAndMid50Filtering(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.5),
		testCandle(3, 100.5, 101.5, 99.5, 101),
		testCandle(4, 101, 101.7, 100.5, 101.5),
		testCandle(5, 101.5, 102.5, 101, 102),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false, false},
		[]bool{false, true, false, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules: []DetectorContextRefinementRule{
			{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
			{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
		},
	}

	candidates, summary, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineTransitionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want one hold_3_inside row; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.ContextRule != DetectorContextRuleHold3Inside ||
		row.DecisionIndex != 4 ||
		row.DecisionMidSide != holdInsideDecisionMidSideAbove ||
		row.DecisionClosePositionBucket != decisionClosePositionBucketHigh25 {
		t.Fatalf("bad delayed candidate: %+v", row)
	}
	for _, row := range summary {
		if row.ContextRule == DetectorContextRuleHold3InsideMid50 && row.CandidateCount != 0 {
			t.Fatalf("mid-50 rule should reject high-25 decision close: %+v", row)
		}
		if row.SourceEpisodeCount != 1 {
			t.Fatalf("source denominator should include the episode for every mid-side/bucket row: %+v", row)
		}
	}
}

func TestHoldInsideMidlineTransitionLabelsStartAfterDecisionAndSkipMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 104, 100, 103),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false},
		[]bool{false, true, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1, 3},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, summary, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineTransitionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want one row for only available h1 horizon; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.DecisionIndex != 2 ||
		row.LabelWindowStartIndex != 3 ||
		row.LabelWindowEndIndex != 3 ||
		row.HorizonBars != 1 ||
		!row.LabelQuickInvalidated ||
		!row.LabelInvalidatedUp ||
		!row.LabelTrendedUp {
		t.Fatalf("label should start after decision index and skip missing future: %+v", row)
	}
	for _, row := range summary {
		if row.HorizonBars != 1 {
			t.Fatalf("missing-future horizon should be skipped, summary=%+v", summary)
		}
	}
}

func TestHoldInsideMidlineTransitionMidlineLabelsFromDecisionSide(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}

	lowHalfCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.5, 98.8, 100.2),
	}
	lowRow, ok := newHoldInsideMidlineTransitionCandidateRow(lowHalfCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected low-half row")
	}
	if lowRow.DecisionMidSide != holdInsideDecisionMidSideBelow ||
		!lowRow.LabelTouchedMid ||
		!lowRow.LabelClosedAcrossMid ||
		lowRow.LabelFirstMidTouchDelayBars != 1 ||
		lowRow.LabelFirstMidCloseAcrossDelayBars != 1 ||
		!lowRow.LabelMidTouchBeforeBoundaryTouch ||
		!lowRow.LabelMidCrossBeforeBoundaryBreak {
		t.Fatalf("expected low-half decision to touch and close across mid before boundary: %+v", lowRow)
	}

	highHalfCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101.5, 100.5, 101),
		testCandle(3, 101, 101.2, 99.5, 99.8),
	}
	highRow, ok := newHoldInsideMidlineTransitionCandidateRow(highHalfCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected high-half row")
	}
	if highRow.DecisionMidSide != holdInsideDecisionMidSideAbove ||
		!highRow.LabelTouchedMid ||
		!highRow.LabelClosedAcrossMid ||
		highRow.LabelFirstMidTouchDelayBars != 1 ||
		highRow.LabelFirstMidCloseAcrossDelayBars != 1 {
		t.Fatalf("expected high-half decision to touch and close across mid: %+v", highRow)
	}
}

func TestHoldInsideMidlineTransitionDelaySentinelsAndStrictBeforeBoundaryOrdering(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}

	noMidCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 99.5, 98.8, 99.2),
	}
	noMid, ok := newHoldInsideMidlineTransitionCandidateRow(noMidCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected no-mid row")
	}
	if noMid.LabelTouchedMid ||
		noMid.LabelClosedAcrossMid ||
		noMid.LabelFirstMidTouchDelayBars != -1 ||
		noMid.LabelFirstMidCloseAcrossDelayBars != -1 ||
		noMid.LabelMidTouchBeforeBoundaryTouch ||
		noMid.LabelMidCrossBeforeBoundaryBreak {
		t.Fatalf("expected no-mid sentinels and no ordering labels: %+v", noMid)
	}

	sameBarTouchCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 102, 98.8, 100.2),
	}
	sameBarTouch, ok := newHoldInsideMidlineTransitionCandidateRow(sameBarTouchCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected same-bar touch row")
	}
	if !sameBarTouch.LabelTouchedMid || sameBarTouch.LabelMidTouchBeforeBoundaryTouch {
		t.Fatalf("same-bar mid touch and boundary touch must not count as before: %+v", sameBarTouch)
	}
	if !sameBarTouch.LabelMidCrossBeforeBoundaryBreak {
		t.Fatalf("mid close across should precede absent boundary close break: %+v", sameBarTouch)
	}

	sameBarBreakCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 103, 98.8, 103),
	}
	sameBarBreak, ok := newHoldInsideMidlineTransitionCandidateRow(sameBarBreakCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected same-bar break row")
	}
	if !sameBarBreak.LabelClosedAcrossMid ||
		sameBarBreak.LabelMidTouchBeforeBoundaryTouch ||
		sameBarBreak.LabelMidCrossBeforeBoundaryBreak {
		t.Fatalf("same-bar mid and boundary events must not satisfy strict before labels: %+v", sameBarBreak)
	}
}

func TestHoldInsideMidlineTransitionDecisionFeaturesDoNotUseFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 99),
		testCandle(3, 99, 100.5, 98.8, 100.2),
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[3] = testCandle(3, 100, 500, 1, 250)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}
	episode := holdInsideTestEpisode()

	row, ok := newHoldInsideMidlineTransitionCandidateRow(candles, profile, rule, episode, fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected row")
	}
	changedRow, ok := newHoldInsideMidlineTransitionCandidateRow(changedFuture, profile, rule, episode, fullSplitName, 2, 1, 1)
	if !ok {
		t.Fatalf("expected changed row")
	}
	if row.EpisodeHigh != changedRow.EpisodeHigh ||
		row.EpisodeLow != changedRow.EpisodeLow ||
		row.EpisodeMid != changedRow.EpisodeMid ||
		row.DecisionClose != changedRow.DecisionClose ||
		row.DecisionClosePosition != changedRow.DecisionClosePosition ||
		row.DecisionMidSide != changedRow.DecisionMidSide ||
		row.DecisionDistanceToHighPct != changedRow.DecisionDistanceToHighPct ||
		row.DecisionDistanceToLowPct != changedRow.DecisionDistanceToLowPct ||
		row.DecisionDistanceToMidPct != changedRow.DecisionDistanceToMidPct {
		t.Fatalf("decision features changed after future edit: got %+v want %+v", changedRow, row)
	}
	if row.LabelQuickInvalidated == changedRow.LabelQuickInvalidated {
		t.Fatalf("label should still respond to future edit: got %+v changed %+v", row, changedRow)
	}
}

func TestHoldInsideMidlineTransitionSummaryDenominatorsDoNotDuplicateAcrossAggregations(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 99),
		testCandle(3, 99, 100.5, 98.8, 100.2),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false},
		[]bool{false, true, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, summary, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineTransitionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want one row", len(candidates))
	}
	if len(summary) != 4 {
		t.Fatalf("summary=%d, want mid-side x all and bucket x all rows; rows=%+v", len(summary), summary)
	}
	for _, row := range summary {
		if row.SourceEpisodeCount != 1 || row.CandidateCount != 1 || row.CandidateRate != 1 {
			t.Fatalf("summary aggregations should not duplicate source denominator: %+v", row)
		}
	}
}

func TestHoldInsideMidlineTransitionStabilityExcludesFullSplit(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}
	summary := []HoldInsideMidlineTransitionSummaryRow{
		holdInsideMidlineSummaryForStability("2021_2022_stress", 10, 5, 0.5, 0.1),
		holdInsideMidlineSummaryForStability("2023_2024_oos", 20, 10, 0.5, 0.2),
		holdInsideMidlineSummaryForStability("2025_2026_recent", 30, 15, 0.5, 0.3),
		holdInsideMidlineSummaryForStability(fullSplitName, 1000, 900, 0.9, 9.0),
	}

	rows := holdInsideMidlineTransitionStabilityRows([]DetectorSweepProfile{profile}, []DetectorContextRefinementRule{rule}, []int{12}, summary, DefaultSplits())
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1: %+v", len(rows), rows)
	}
	row := rows[0]
	if row.SourceEpisodeCount != 60 ||
		row.SourceEpisodeCountMin != 10 ||
		row.SourceEpisodeCountMax != 30 ||
		row.SourceEpisodeCountDelta != 20 ||
		row.CandidateCount != 30 ||
		row.CandidateCountMin != 5 ||
		row.CandidateCountMax != 15 ||
		!boundaryAlmostEqual(row.LabelTouchedMidRateMax, 0.3) ||
		!boundaryAlmostEqual(row.LabelTouchedMidRateDelta, 0.2) {
		t.Fatalf("bad stability row, full split may have leaked in: %+v", row)
	}
}

func TestHoldInsideMidlineTransitionCandidateSorting(t *testing.T) {
	rows := []HoldInsideMidlineTransitionCandidateRow{
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold6Inside, Split: fullSplitName, SourceEpisodeID: 2, HorizonBars: 12},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, HorizonBars: 6},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, HorizonBars: 3},
	}
	if !lessHoldInsideMidlineTransitionCandidate(rows[1], rows[0]) ||
		!lessHoldInsideMidlineTransitionCandidate(rows[2], rows[1]) {
		t.Fatalf("unexpected candidate sort order")
	}
}

func TestHoldInsideMidlineTransitionRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	classifications := testCompressionClassifications([]bool{false}, []bool{false})

	if _, _, _, err := RunHoldInsideMidlineTransitionAudit(candles, RangeDetectorConfig{ATRPeriod: -1}, HoldInsideMidlineTransitionAuditConfig{}, nil); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
	if _, _, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{0},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: -1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid quick invalidation error")
	}
	if _, _, _, err := runHoldInsideMidlineTransitionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineTransitionAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleEpisodeEnd}},
	}, nil); err == nil {
		t.Fatalf("expected positive hold-bars rule error")
	}
}

func holdInsideMidlineSummaryForStability(split string, sourceEpisodes, candidates int, candidateRate, touchedMid float64) HoldInsideMidlineTransitionSummaryRow {
	return HoldInsideMidlineTransitionSummaryRow{
		ProfileID:                            BalancedDetectorProfileID,
		IsBalancedBaseline:                   true,
		ContextRule:                          DetectorContextRuleHold3Inside,
		HoldBars:                             3,
		Split:                                split,
		HorizonBars:                          12,
		DecisionMidSide:                      holdInsideDecisionMidSideAll,
		DecisionClosePositionBucket:          holdInsideDecisionPositionBucketAll,
		SourceEpisodeCount:                   sourceEpisodes,
		CandidateCount:                       candidates,
		CandidateRate:                        candidateRate,
		LabelTouchedMidRate:                  touchedMid,
		LabelClosedAcrossMidRate:             touchedMid,
		LabelMidTouchBeforeBoundaryTouchRate: touchedMid,
		LabelMidCrossBeforeBoundaryBreakRate: touchedMid,
		LabelPersistedInsideRangeRate:        touchedMid,
		LabelQuickInvalidatedRate:            touchedMid,
		LabelTrendedUpRate:                   touchedMid / 2,
		LabelTrendedDownRate:                 touchedMid / 2,
		LabelAvgFirstMidTouchDelayBars:       touchedMid,
		LabelAvgFirstMidCloseAcrossDelayBars: touchedMid,
	}
}
