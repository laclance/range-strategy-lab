package lab

import "testing"

func TestDetectorContextRefinementDefaults(t *testing.T) {
	cfg := DefaultDetectorContextRefinementAuditConfig()
	if len(cfg.HorizonsBars) != 4 ||
		cfg.HorizonsBars[0] != 1 ||
		cfg.HorizonsBars[1] != 3 ||
		cfg.HorizonsBars[2] != 6 ||
		cfg.HorizonsBars[3] != 12 ||
		cfg.QuickInvalidationBars != 3 {
		t.Fatalf("bad default config: %+v", cfg)
	}

	profiles := DefaultDetectorContextRefinementProfiles(20)
	wantProfiles := []string{
		"p30_c12_bollinger_on_adx_off",
		"p30_c12_bollinger_on_adx_on",
		"p20_c24_bollinger_on_adx_off",
		"p30_c24_bollinger_on_adx_off",
		"p40_c24_bollinger_on_adx_off",
		"p20_c24_bollinger_on_adx_on",
		"p30_c24_bollinger_on_adx_on",
		"p40_c24_bollinger_on_adx_on",
	}
	if len(profiles) != len(wantProfiles) {
		t.Fatalf("profiles=%d, want %d", len(profiles), len(wantProfiles))
	}
	for i, want := range wantProfiles {
		if profiles[i].ProfileID != want {
			t.Fatalf("profile %d=%q, want %q", i, profiles[i].ProfileID, want)
		}
	}

	wantRules := []DetectorContextRefinementRule{
		{RuleID: DetectorContextRuleEpisodeEnd, HoldBars: 0},
		{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1},
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

func TestDetectorContextRefinementDelayedRulesAndMid50Filter(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.5),
		testCandle(3, 100.5, 101.5, 99.5, 101),
		testCandle(4, 101, 101.7, 100.5, 101.5),
		testCandle(5, 101.5, 102.5, 101, 102.2),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false, false},
		[]bool{false, true, false, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules: []DetectorContextRefinementRule{
			{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3},
			{RuleID: DetectorContextRuleHold3InsideMid50, HoldBars: 3, RequireMid50: true},
		},
	}

	candidates, summary, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runDetectorContextRefinementFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want only hold_3_inside candidate; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.ContextRule != DetectorContextRuleHold3Inside ||
		row.DecisionIndex != 4 ||
		row.DecisionClosePositionBucket != decisionClosePositionBucketHigh25 {
		t.Fatalf("bad delayed candidate: %+v", row)
	}
	if len(summary) != 2 {
		t.Fatalf("summary rows=%d, want 2", len(summary))
	}
	for _, row := range summary {
		if row.SourceEpisodeCount != 1 {
			t.Fatalf("source denominator should include the episode for both rules: %+v", row)
		}
		if row.ContextRule == DetectorContextRuleHold3InsideMid50 && row.CandidateCount != 0 {
			t.Fatalf("mid-50 rule should reject high-25 decision close: %+v", row)
		}
	}
}

func TestDetectorContextRefinementDelayedRuleSkipsIfPreDecisionCloseLeavesRange(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.5),
		testCandle(3, 100.5, 101, 97, 97.5),
		testCandle(4, 97.5, 101, 97, 100),
		testCandle(5, 100, 101, 99, 100),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, false, false},
		[]bool{false, true, false, false, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}},
	}

	candidates, summary, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runDetectorContextRefinementFromClassifications error: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("expected no candidate after pre-decision close left range: %+v", candidates)
	}
	if len(summary) != 1 || summary[0].SourceEpisodeCount != 1 || summary[0].CandidateCount != 0 || summary[0].CandidateRate != 0 {
		t.Fatalf("bad rejected summary: %+v", summary)
	}
}

func TestDetectorContextRefinementLabelsStartAfterDecisionAndSkipMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.5),
		testCandle(3, 100.5, 104, 101, 103),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false},
		[]bool{false, true, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1, 3},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	candidates, summary, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runDetectorContextRefinementFromClassifications error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidates=%d, want only available h1 row; rows=%+v", len(candidates), candidates)
	}
	row := candidates[0]
	if row.DecisionIndex != 2 ||
		row.LabelWindowStartIndex != 3 ||
		row.LabelWindowEndIndex != 3 ||
		!row.LabelQuickInvalidated ||
		!row.LabelInvalidatedUp ||
		!row.LabelTrendedUp {
		t.Fatalf("label should start after decision index: %+v", row)
	}
	if len(summary) != 1 || summary[0].HorizonBars != 1 {
		t.Fatalf("missing-future horizon should be skipped, summary=%+v", summary)
	}
}

func TestDetectorContextRefinementSummaryCandidateRateUsesSourceEpisodes(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 101, 99, 100),
		testCandle(4, 100, 101, 99, 100),
		testCandle(5, 100, 102, 98, 100),
		testCandle(6, 100, 103, 102.5, 103),
		testCandle(7, 103, 104, 102, 103),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, false, true, true, false, false},
		[]bool{false, true, false, false, false, true, false, false},
	)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}},
	}

	_, summary, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runDetectorContextRefinementFromClassifications error: %v", err)
	}
	if len(summary) != 1 {
		t.Fatalf("summary rows=%d, want 1", len(summary))
	}
	row := summary[0]
	if row.SourceEpisodeCount != 2 ||
		row.CandidateCount != 1 ||
		!boundaryAlmostEqual(row.CandidateRate, 0.5) {
		t.Fatalf("candidate rate should use source episode denominator: %+v", row)
	}
}

func TestDetectorContextRefinementStabilityRowsCalculateMinMaxDelta(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, true, 20, true)
	rule := DetectorContextRefinementRule{RuleID: DetectorContextRuleHold1Inside, HoldBars: 1}
	summary := []DetectorContextRefinementSummaryRow{
		{
			ProfileID:                     profile.ProfileID,
			ContextRule:                   rule.RuleID,
			Split:                         "2021_2022_stress",
			HorizonBars:                   12,
			SourceEpisodeCount:            20,
			CandidateCount:                10,
			CandidateRate:                 0.50,
			LabelPersistedInsideRangeRate: 0.20,
			LabelQuickInvalidatedRate:     0.60,
			LabelChoppedRate:              0.10,
			LabelTrendedUpRate:            0.20,
			LabelTrendedDownRate:          0.10,
			LabelAvgCloseDriftPct:         -0.01,
		},
		{
			ProfileID:                     profile.ProfileID,
			ContextRule:                   rule.RuleID,
			Split:                         "2023_2024_oos",
			HorizonBars:                   12,
			SourceEpisodeCount:            30,
			CandidateCount:                15,
			CandidateRate:                 0.50,
			LabelPersistedInsideRangeRate: 0.30,
			LabelQuickInvalidatedRate:     0.50,
			LabelChoppedRate:              0.20,
			LabelTrendedUpRate:            0.10,
			LabelTrendedDownRate:          0.20,
			LabelAvgCloseDriftPct:         0.02,
		},
		{
			ProfileID:                     profile.ProfileID,
			ContextRule:                   rule.RuleID,
			Split:                         "2025_2026_recent",
			HorizonBars:                   12,
			SourceEpisodeCount:            25,
			CandidateCount:                5,
			CandidateRate:                 0.20,
			LabelPersistedInsideRangeRate: 0.25,
			LabelQuickInvalidatedRate:     0.70,
			LabelChoppedRate:              0.05,
			LabelTrendedUpRate:            0.30,
			LabelTrendedDownRate:          0.15,
			LabelAvgCloseDriftPct:         0.01,
		},
	}

	rows := detectorContextStabilityRows([]DetectorSweepProfile{profile}, []DetectorContextRefinementRule{rule}, []int{12}, summary, DefaultSplits())
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if !row.IsADXComparison ||
		row.PeriodSplits != 3 ||
		row.SourceEpisodeCount != 75 ||
		row.SourceEpisodeCountMin != 20 ||
		row.SourceEpisodeCountMax != 30 ||
		row.SourceEpisodeCountDelta != 10 ||
		row.CandidateCount != 30 ||
		row.CandidateCountMin != 5 ||
		row.CandidateCountMax != 15 ||
		!boundaryAlmostEqual(row.CandidateRateDelta, 0.30) ||
		!boundaryAlmostEqual(row.LabelPersistedInsideRangeRateMin, 0.20) ||
		!boundaryAlmostEqual(row.LabelQuickInvalidatedRateDelta, 0.20) ||
		!boundaryAlmostEqual(row.LabelTrendedRateMax, 0.45) ||
		!boundaryAlmostEqual(row.LabelAvgCloseDriftPctDelta, 0.03) {
		t.Fatalf("bad stability row: %+v", row)
	}
}

func TestDetectorContextRefinementRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	classifications := testCompressionClassifications([]bool{false}, []bool{false})

	if _, _, _, err := RunDetectorContextRefinementAudit(candles, RangeDetectorConfig{ATRPeriod: -1}, DetectorContextRefinementAuditConfig{}, nil); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
	if _, _, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{0},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleEpisodeEnd}},
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: -1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleEpisodeEnd}},
	}, nil); err == nil {
		t.Fatalf("expected invalid quick invalidation error")
	}
	if _, _, _, err := runDetectorContextRefinementFromClassifications(candles, profile, classifications, DetectorContextRefinementAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		ContextRules:          []DetectorContextRefinementRule{{RuleID: DetectorContextRuleHold3InsideMid50, RequireMid50: true}},
	}, nil); err == nil {
		t.Fatalf("expected invalid mid-50 rule error")
	}
}
