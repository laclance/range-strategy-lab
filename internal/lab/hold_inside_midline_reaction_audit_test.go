package lab

import "testing"

func TestHoldInsideMidlineReactionDefaults(t *testing.T) {
	cfg := DefaultHoldInsideMidlineReactionAuditConfig()
	if len(cfg.HorizonsBars) != 4 ||
		cfg.HorizonsBars[0] != 1 ||
		cfg.HorizonsBars[1] != 3 ||
		cfg.HorizonsBars[2] != 6 ||
		cfg.HorizonsBars[3] != 12 ||
		cfg.QuickInvalidationBars != 3 ||
		cfg.MaxMidlineEventDelayBars != 12 {
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

func TestHoldInsideMidlineReactionHoldRuleAndMid50Filtering(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.5),
		testCandle(3, 100.5, 101.5, 99.5, 101),
		testCandle(4, 101, 101.7, 100.5, 101.5),
		testCandle(5, 101.5, 101.8, 99.8, 100.2),
		testCandle(6, 100.2, 101, 99.8, 100.5),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false, false, false},
		[]bool{false, true, false, false, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 2,
		ContextRules: []DetectorContextRefinementRule{
			{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
			{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
		},
	}

	candidates, funnel, summary, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineReactionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want one hold_3_inside mid-touch row; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.ContextRule != DetectorContextRuleHold3Inside ||
		row.HoldDecisionIndex != 4 ||
		row.HoldDecisionCloseBucket != decisionClosePositionBucketHigh25 ||
		row.EventIndex != 5 ||
		row.EventType != HoldInsideMidlineReactionEventTouch {
		t.Fatalf("bad delayed candidate: %+v", row)
	}
	for _, row := range funnel {
		if row.ContextRule == DetectorContextRuleHold3InsideMid50 {
			t.Fatalf("mid-50 rule should reject high-25 hold decision before entering the funnel: %+v", row)
		}
	}
	for _, row := range summary {
		if row.ContextRule == DetectorContextRuleHold3InsideMid50 {
			t.Fatalf("mid-50 rule should reject high-25 hold decision before reaction summary: %+v", row)
		}
	}
}

func TestHoldInsideMidlineReactionEventSearchAndDualEmission(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.6, 98.8, 100.2),
		testCandle(4, 100.2, 101, 99.8, 100.5),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false},
		[]bool{false, true, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 3,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineReactionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%d, want same-bar touch and close-across rows; rows=%+v", len(candidates), candidates)
	}
	if candidates[0].EventType != HoldInsideMidlineReactionEventTouch ||
		candidates[1].EventType != HoldInsideMidlineReactionEventCloseAcross ||
		candidates[0].EventIndex != 3 ||
		candidates[1].EventIndex != 3 ||
		candidates[0].LabelWindowStartIndex != 4 ||
		candidates[1].LabelWindowStartIndex != 4 {
		t.Fatalf("bad same-bar dual emission: %+v", candidates)
	}
}

func TestHoldInsideMidlineReactionTouchVersusCloseAcrossSeparation(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.2, 98.8, 99.4),
		testCandle(4, 99.4, 100.5, 99.1, 100.1),
		testCandle(5, 100.1, 101, 99.8, 100.5),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false, false},
		[]bool{false, true, false, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 3,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineReactionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates=%d, want touch and later close-across rows; rows=%+v", len(candidates), candidates)
	}
	if candidates[0].EventType != HoldInsideMidlineReactionEventTouch ||
		candidates[0].EventIndex != 3 ||
		candidates[1].EventType != HoldInsideMidlineReactionEventCloseAcross ||
		candidates[1].EventIndex != 4 {
		t.Fatalf("touch and close-across should be separated by first matching event: %+v", candidates)
	}
}

func TestHoldInsideMidlineReactionLabelsStartAfterEventAndMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.2, 98.8, 99.5),
		testCandle(4, 99.5, 104, 99.3, 103),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false},
		[]bool{false, true, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1, 3},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 3,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, funnel, summary, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineReactionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want only h1 mid-touch row; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.EventIndex != 3 ||
		row.LabelWindowStartIndex != 4 ||
		row.LabelWindowEndIndex != 4 ||
		row.HorizonBars != 1 ||
		!row.LabelQuickInvalidated ||
		!row.LabelInvalidatedUp ||
		!row.LabelTrendedUp {
		t.Fatalf("label should start after event index and skip missing h3 future: %+v", row)
	}
	foundTouchFunnel := false
	for _, row := range funnel {
		if row.EventType != HoldInsideMidlineReactionEventTouch {
			continue
		}
		foundTouchFunnel = true
		if row.SourceHoldCount != 1 || row.EventCount != 1 || row.MissingFutureCount != 1 || row.ReactionEligibleCount != 0 {
			t.Fatalf("bad touch funnel missing-future accounting: %+v", row)
		}
	}
	if !foundTouchFunnel {
		t.Fatalf("missing touch funnel row: %+v", funnel)
	}
	for _, row := range summary {
		if row.HorizonBars != 1 {
			t.Fatalf("missing-future horizon should be skipped, summary=%+v", summary)
		}
	}
}

func TestHoldInsideMidlineReactionMissingEventAndFunnelDenominators(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 99.5, 98.7, 99.2),
		testCandle(4, 99.2, 99.7, 98.8, 99.4),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false},
		[]bool{false, true, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 2,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, funnel, summary, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runHoldInsideMidlineReactionAuditFromClassifications error: %v", err)
	}
	if len(candidates) != 0 || len(summary) != 0 {
		t.Fatalf("missing event should produce no reaction candidates or summary rows: candidates=%+v summary=%+v", candidates, summary)
	}
	if len(funnel) != 2 {
		t.Fatalf("funnel rows=%d, want one per event type; rows=%+v", len(funnel), funnel)
	}
	for _, row := range funnel {
		if row.SourceHoldCount != 1 ||
			row.EventCount != 0 ||
			row.MissingEventCount != 1 ||
			row.MissingFutureCount != 0 ||
			row.ReactionEligibleCount != 0 ||
			row.EventRate != 0 {
			t.Fatalf("bad missing-event funnel denominator: %+v", row)
		}
	}
}

func TestHoldInsideMidlineReactionEventDecisionFieldsDoNotUseFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 99.5, 98.5, 99),
		testCandle(3, 99, 100.6, 98.8, 100.2),
		testCandle(4, 100.2, 101, 99.8, 100.5),
	}
	changedFuture := append([]Candle(nil), candles...)
	changedFuture[4] = testCandle(4, 100.5, 104, 99.8, 103)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}
	episode := holdInsideTestEpisode()

	row, ok := newHoldInsideMidlineReactionCandidateRow(candles, profile, rule, HoldInsideMidlineReactionEventCloseAcross, episode, fullSplitName, 2, 3, 1, 1)
	if !ok {
		t.Fatalf("expected row")
	}
	changedRow, ok := newHoldInsideMidlineReactionCandidateRow(changedFuture, profile, rule, HoldInsideMidlineReactionEventCloseAcross, episode, fullSplitName, 2, 3, 1, 1)
	if !ok {
		t.Fatalf("expected changed row")
	}
	if row.EpisodeHigh != changedRow.EpisodeHigh ||
		row.EpisodeLow != changedRow.EpisodeLow ||
		row.EpisodeMid != changedRow.EpisodeMid ||
		row.HoldDecisionClose != changedRow.HoldDecisionClose ||
		row.HoldDecisionClosePosition != changedRow.HoldDecisionClosePosition ||
		row.HoldDecisionMidSide != changedRow.HoldDecisionMidSide ||
		row.EventClose != changedRow.EventClose ||
		row.EventClosePosition != changedRow.EventClosePosition ||
		row.EventMidSide != changedRow.EventMidSide ||
		row.EventDistanceToHighPct != changedRow.EventDistanceToHighPct ||
		row.EventDistanceToLowPct != changedRow.EventDistanceToLowPct ||
		row.EventDistanceToMidPct != changedRow.EventDistanceToMidPct {
		t.Fatalf("event decision fields changed after future edit: got %+v want %+v", changedRow, row)
	}
	if row.LabelQuickInvalidated == changedRow.LabelQuickInvalidated {
		t.Fatalf("label should still respond to future edit: got %+v changed %+v", row, changedRow)
	}
}

func TestHoldInsideMidlineReactionStabilityExcludesFullSplit(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}
	summary := []HoldInsideMidlineReactionSummaryRow{
		holdInsideMidlineReactionSummaryForStability("2021_2022_stress", 10, 0.1),
		holdInsideMidlineReactionSummaryForStability("2023_2024_oos", 20, 0.2),
		holdInsideMidlineReactionSummaryForStability("2025_2026_recent", 30, 0.3),
		holdInsideMidlineReactionSummaryForStability(fullSplitName, 1000, 9.0),
	}

	rows := holdInsideMidlineReactionStabilityRows([]DetectorSweepProfile{profile}, []DetectorContextRefinementRule{rule}, []int{12}, summary, DefaultSplits())
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1: %+v", len(rows), rows)
	}
	row := rows[0]
	if row.CandidateCount != 60 ||
		row.CandidateCountMin != 10 ||
		row.CandidateCountMax != 30 ||
		row.CandidateCountDelta != 20 ||
		!boundaryAlmostEqual(row.LabelTouchedHighRateMax, 0.3) ||
		!boundaryAlmostEqual(row.LabelTouchedHighRateDelta, 0.2) {
		t.Fatalf("bad stability row, full split may have leaked in: %+v", row)
	}
}

func TestHoldInsideMidlineReactionCandidateSorting(t *testing.T) {
	rows := []HoldInsideMidlineReactionCandidateRow{
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold6Inside, Split: fullSplitName, SourceEpisodeID: 2, EventType: HoldInsideMidlineReactionEventCloseAcross, EventIndex: 5, HorizonBars: 12},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, EventType: HoldInsideMidlineReactionEventCloseAcross, EventIndex: 3, HorizonBars: 3},
		{ProfileID: BalancedDetectorProfileID, ContextRule: DetectorContextRuleHold3Inside, Split: "2021_2022_stress", SourceEpisodeID: 1, EventType: HoldInsideMidlineReactionEventTouch, EventIndex: 3, HorizonBars: 6},
	}
	if !lessHoldInsideMidlineReactionCandidate(rows[1], rows[0]) ||
		!lessHoldInsideMidlineReactionCandidate(rows[2], rows[1]) {
		t.Fatalf("unexpected candidate sort order")
	}
}

func TestHoldInsideMidlineReactionRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	classifications := testCompressionClassifications([]bool{false}, []bool{false})

	if _, _, _, _, err := RunHoldInsideMidlineReactionAudit(candles, RangeDetectorConfig{ATRPeriod: -1}, HoldInsideMidlineReactionAuditConfig{}, nil); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
	if _, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{0},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 1,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    -1,
		MaxMidlineEventDelayBars: 1,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid quick invalidation error")
	}
	if _, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: -1,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}, nil); err == nil {
		t.Fatalf("expected invalid max event delay error")
	}
	if _, _, _, _, err := runHoldInsideMidlineReactionAuditFromClassifications(candles, profile, classifications, HoldInsideMidlineReactionAuditConfig{
		HorizonsBars:             []int{1},
		QuickInvalidationBars:    1,
		MaxMidlineEventDelayBars: 1,
		ContextRules:             []DetectorContextRefinementRule{{RuleID: DetectorContextRuleEpisodeEnd}},
	}, nil); err == nil {
		t.Fatalf("expected positive hold-bars rule error")
	}
}

func holdInsideMidlineReactionSummaryForStability(split string, candidates int, rate float64) HoldInsideMidlineReactionSummaryRow {
	return HoldInsideMidlineReactionSummaryRow{
		ProfileID:                        BalancedDetectorProfileID,
		IsBalancedBaseline:               true,
		ContextRule:                      DetectorContextRuleHold3Inside,
		HoldBars:                         3,
		EventType:                        HoldInsideMidlineReactionEventTouch,
		Split:                            split,
		HorizonBars:                      12,
		EventMidSide:                     holdInsideDecisionMidSideAll,
		EventClosePositionBucket:         holdInsideDecisionPositionBucketAll,
		CandidateCount:                   candidates,
		AvgEventDelayBars:                rate,
		LabelReenteredRangeRate:          rate,
		LabelPersistedInsideRangeRate:    rate,
		LabelQuickInvalidatedRate:        rate,
		LabelTrendedUpRate:               rate / 2,
		LabelTrendedDownRate:             rate / 2,
		LabelTouchedHighRate:             rate,
		LabelTouchedLowRate:              rate,
		LabelTouchedOppositeHalfRate:     rate,
		LabelClosedBackAcrossMidRate:     rate,
		LabelMidRejectionBeforeBoundRate: rate,
		LabelBoundBeforeMidRejectionRate: rate,
	}
}
