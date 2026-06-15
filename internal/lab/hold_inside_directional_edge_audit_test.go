package lab

import "testing"

func TestHoldInsideDirectionalEdgeDefaults(t *testing.T) {
	cfg := DefaultHoldInsideDirectionalEdgeAuditConfig()
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

func TestHoldInsideDirectionalEdgeHoldRuleAndMid50Filtering(t *testing.T) {
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
	cfg := HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules: []DetectorContextRefinementRule{
			{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
			{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
		},
	}

	candidates, summary, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideDirectionalEdgeAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%d, want two paper-side rows for hold_3_inside; rows=%+v", len(candidates), candidates)
	}
	for _, row := range candidates {
		if row.ContextRule != DetectorContextRuleHold3Inside ||
			row.DecisionIndex != 4 ||
			row.DecisionClosePositionBucket != decisionClosePositionBucketHigh25 {
			t.Fatalf("bad delayed candidate: %+v", row)
		}
	}
	for _, row := range summary {
		if row.ContextRule == DetectorContextRuleHold3InsideMid50 && row.CandidateCount != 0 {
			t.Fatalf("mid-50 rule should reject high-25 decision close: %+v", row)
		}
		if row.SourceEpisodeCount != 1 {
			t.Fatalf("source denominator should include the episode for every side/bucket: %+v", row)
		}
	}
}

func TestHoldInsideDirectionalEdgeLabelsStartAfterDecisionAndSkipMissingFuture(t *testing.T) {
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
	cfg := HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1, 3},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, summary, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideDirectionalEdgeAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%d, want two paper-side rows for only available h1 horizon; rows=%+v", len(candidates), candidates)
	}
	for _, row := range candidates {
		if row.DecisionIndex != 2 ||
			row.LabelWindowStartIndex != 3 ||
			row.LabelWindowEndIndex != 3 ||
			row.HorizonBars != 1 ||
			!row.LabelQuickInvalidated ||
			!row.LabelInvalidatedUp ||
			!row.LabelTrendedUp {
			t.Fatalf("label should start after decision index and skip missing future: %+v", row)
		}
	}
	for _, row := range summary {
		if row.HorizonBars != 1 {
			t.Fatalf("missing-future horizon should be skipped, summary=%+v", summary)
		}
	}
}

func TestHoldInsideDirectionalEdgePaperSideSymmetry(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 104, 99, 101),
	}
	episode := holdInsideTestEpisode()
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}

	towardHigh, ok := newHoldInsideDirectionalEdgeCandidateRow(candles, profile, rule, episode, fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardHigh)
	if !ok {
		t.Fatalf("expected toward-high row")
	}
	towardLow, ok := newHoldInsideDirectionalEdgeCandidateRow(candles, profile, rule, episode, fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardLow)
	if !ok {
		t.Fatalf("expected toward-low row")
	}
	if !boundaryAlmostEqual(towardHigh.LabelFavorableMovePct, 0.04) ||
		!boundaryAlmostEqual(towardHigh.LabelAdverseMovePct, 0.01) ||
		!boundaryAlmostEqual(towardLow.LabelFavorableMovePct, 0.01) ||
		!boundaryAlmostEqual(towardLow.LabelAdverseMovePct, 0.04) {
		t.Fatalf("bad paper-side symmetry: high=%+v low=%+v", towardHigh, towardLow)
	}
	if !towardHigh.LabelFavorableGTAdverse || towardLow.LabelFavorableGTAdverse {
		t.Fatalf("bad favorable/adverse direction labels: high=%+v low=%+v", towardHigh, towardLow)
	}
}

func TestHoldInsideDirectionalEdgeMidlineLabelsFromDecisionSide(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}

	lowHalfCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.5, 98.8, 100.2),
	}
	lowRow, ok := newHoldInsideDirectionalEdgeCandidateRow(lowHalfCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardHigh)
	if !ok {
		t.Fatalf("expected low-half row")
	}
	if lowRow.DecisionMidSide != holdInsideDecisionMidSideBelow || !lowRow.LabelTouchedMid || !lowRow.LabelClosedAcrossMid {
		t.Fatalf("expected low-half decision to touch and close across mid: %+v", lowRow)
	}

	highHalfCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101.5, 100.5, 101),
		testCandle(3, 101, 101.2, 99.5, 99.8),
	}
	highRow, ok := newHoldInsideDirectionalEdgeCandidateRow(highHalfCandles, profile, rule, holdInsideTestEpisode(), fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardLow)
	if !ok {
		t.Fatalf("expected high-half row")
	}
	if highRow.DecisionMidSide != holdInsideDecisionMidSideAbove || !highRow.LabelTouchedMid || !highRow.LabelClosedAcrossMid {
		t.Fatalf("expected high-half decision to touch and close across mid: %+v", highRow)
	}
}

func TestHoldInsideDirectionalEdgeDecisionFeaturesDoNotUseFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 104, 99, 101),
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[3] = testCandle(3, 100, 500, 1, 250)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}
	episode := holdInsideTestEpisode()

	row, ok := newHoldInsideDirectionalEdgeCandidateRow(candles, profile, rule, episode, fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardHigh)
	if !ok {
		t.Fatalf("expected row")
	}
	changedRow, ok := newHoldInsideDirectionalEdgeCandidateRow(changedFuture, profile, rule, episode, fullSplitName, 2, 1, 1, HoldInsidePaperSideTowardHigh)
	if !ok {
		t.Fatalf("expected changed row")
	}
	if row.EpisodeHigh != changedRow.EpisodeHigh ||
		row.EpisodeLow != changedRow.EpisodeLow ||
		row.EpisodeMid != changedRow.EpisodeMid ||
		row.DecisionClose != changedRow.DecisionClose ||
		row.DecisionClosePosition != changedRow.DecisionClosePosition ||
		row.DecisionDistanceToHighPct != changedRow.DecisionDistanceToHighPct ||
		row.DecisionDistanceToLowPct != changedRow.DecisionDistanceToLowPct ||
		row.DecisionDistanceToMidPct != changedRow.DecisionDistanceToMidPct {
		t.Fatalf("decision features changed after future edit: got %+v want %+v", changedRow, row)
	}
	if row.LabelFavorableMovePct == changedRow.LabelFavorableMovePct {
		t.Fatalf("label should still respond to future edit: got %+v changed %+v", row, changedRow)
	}
}

func TestHoldInsideDirectionalEdgeSummaryDenominatorsDoNotDoublePaperSides(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 104, 99, 101),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false},
		[]bool{false, true, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, summary, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideDirectionalEdgeAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%d, want two paper-side rows", len(candidates))
	}
	if len(summary) != 4 {
		t.Fatalf("summary=%d, want side x all/bucket rows; rows=%+v", len(summary), summary)
	}
	for _, row := range summary {
		if row.SourceEpisodeCount != 1 || row.CandidateCount != 1 || row.CandidateRate != 1 {
			t.Fatalf("paper-side expansion should not double source denominator: %+v", row)
		}
	}
}

func TestHoldInsideDirectionalEdgeStabilityExcludesFullSplit(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}
	summary := []HoldInsideDirectionalEdgeSummaryRow{
		holdInsideSummaryForStability("2021_2022_stress", 10, 5, 0.5, 0.1),
		holdInsideSummaryForStability("2023_2024_oos", 20, 10, 0.5, 0.2),
		holdInsideSummaryForStability("2025_2026_recent", 30, 15, 0.5, 0.3),
		holdInsideSummaryForStability(fullSplitName, 1000, 900, 0.9, 9.0),
	}

	rows := holdInsideDirectionalEdgeStabilityRows([]DetectorSweepProfile{profile}, []DetectorContextRefinementRule{rule}, []int{12}, summary, DefaultSplits())
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
		!boundaryAlmostEqual(row.LabelAvgFavorableMinusAdversePctMax, 0.3) ||
		!boundaryAlmostEqual(row.LabelAvgFavorableMinusAdversePctDelta, 0.2) {
		t.Fatalf("bad stability row, full split may have leaked in: %+v", row)
	}
}

func TestHoldInsideDirectionalEdgeCandidateSorting(t *testing.T) {
	rows := []HoldInsideDirectionalEdgeCandidateRow{
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold6Inside, Split: fullSplitName, SourceEpisodeID: 2, HorizonBars: 12, PaperSide: HoldInsidePaperSideTowardLow},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, HorizonBars: 3, PaperSide: HoldInsidePaperSideTowardLow},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, HorizonBars: 3, PaperSide: HoldInsidePaperSideTowardHigh},
	}
	if !lessHoldInsideDirectionalEdgeCandidate(rows[1], rows[0]) ||
		!lessHoldInsideDirectionalEdgeCandidate(rows[2], rows[1]) {
		t.Fatalf("unexpected candidate sort order")
	}
}

func TestHoldInsideDirectionalEdgeRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	classifications := testCompressionClassifications([]bool{false}, []bool{false})

	if _, _, _, err := RunHoldInsideDirectionalEdgeAudit(candles, RangeDetectorConfig{ATRPeriod: -1}, HoldInsideDirectionalEdgeAuditConfig{}, nil); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
	if _, _, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{0},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: -1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid quick invalidation error")
	}
	if _, _, _, err := runHoldInsideDirectionalEdgeAuditFromClassifications(candles, profile, classifications, HoldInsideDirectionalEdgeAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleEpisodeEnd}},
	}, nil); err == nil {
		t.Fatalf("expected positive hold-bars rule error")
	}
	if _, ok := newHoldInsideDirectionalEdgeCandidateRow(candles, profile, DetectorContextRefinementRule{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}, holdInsideTestEpisode(), fullSplitName, 0, 1, 1, "bad"); ok {
		t.Fatalf("expected invalid paper side rejection")
	}
}

func holdInsideTestEpisode() rangeRegimeDurabilityEpisode {
	return rangeRegimeDurabilityEpisode{
		EpisodeID:          1,
		Split:              fullSplitName,
		StartIndex:         0,
		EndIndex:           1,
		High:               102,
		Low:                98,
		EndClose:           100,
		RawLengthBars:      2,
		ActiveLengthBars:   1,
		RawLengthBucket:    "lt_12",
		ActiveLengthBucket: "lt_12",
		WidthPct:           0.04,
		WidthBucket:        "gt_50bp",
		DetectorProfileID:  BalancedDetectorProfileID,
	}
}

func holdInsideSummaryForStability(split string, sourceEpisodes, candidates int, candidateRate, favorableMinusAdverse float64) HoldInsideDirectionalEdgeSummaryRow {
	return HoldInsideDirectionalEdgeSummaryRow{
		ProfileID:                        BalancedDetectorProfileID,
		IsBalancedBaseline:               true,
		ContextRule:                      DetectorContextRuleHold3Inside,
		HoldBars:                         3,
		Split:                            split,
		HorizonBars:                      12,
		PaperSide:                        HoldInsidePaperSideTowardHigh,
		DecisionClosePositionBucket:      holdInsideDecisionPositionBucketAll,
		SourceEpisodeCount:               sourceEpisodes,
		CandidateCount:                   candidates,
		CandidateRate:                    candidateRate,
		LabelAvgFavorableMinusAdversePct: favorableMinusAdverse,
		LabelFavorableGTAdverseRate:      favorableMinusAdverse,
		LabelAvgFavorableMovePct:         favorableMinusAdverse + 0.1,
		LabelAvgAdverseMovePct:           0.1,
		LabelTouchedMidRate:              favorableMinusAdverse,
		LabelClosedAcrossMidRate:         favorableMinusAdverse,
		LabelSideBoundaryTouchRate:       favorableMinusAdverse,
		LabelOppositeCloseBreakRate:      favorableMinusAdverse,
		LabelQuickInvalidatedRate:        favorableMinusAdverse,
	}
}
