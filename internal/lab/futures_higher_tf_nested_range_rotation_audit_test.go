package lab

import (
	"testing"
	"time"
)

func TestNestedRangeRotationSourceAndResampleAcceptance(t *testing.T) {
	candles := make([]Candle, 96)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	manifest, err := ValidateResearchSource("btcusdt_futures_um_5m_test.csv", []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}, candles, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := nestedRotationTestConfig()
	cfg.ApprovedSourcePath = manifest.Path
	cfg.ExpectedParentRows = 2
	cfg.ExpectedParentLastOpenTime = candles[48].OpenTime.UTC().Format(timeLayout)
	cfg.ExpectedChildRows = 8
	cfg.ExpectedChildLastOpenTime = candles[84].OpenTime.UTC().Format(timeLayout)
	cfg.SkipCoverageCountCheck = false

	result, err := RunFuturesHigherTFNestedRangeRotationAudit(candles, manifest, cfg, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if got := result.SourceRows[0].ValidationStatus; got != "accepted" {
		t.Fatalf("source status=%s", got)
	}
	if len(result.CoverageRows) != 2 {
		t.Fatalf("coverage rows=%d, want 2", len(result.CoverageRows))
	}
	for _, row := range result.CoverageRows {
		if row.ValidationStatus != "accepted" || !row.CoverageFactsPass {
			t.Fatalf("bad coverage row: %+v", row)
		}
	}
}

func TestNestedRangeRotationParentFreezeAndInvalidation(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 106, 112, 99, 106),
		testCandle(2, 107, 111, 101, 107),
		testCandle(3, 108, 130, 90, 108),
		testCandle(4, 120, 121, 119, 120),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true},
		{Index: 1, RawActive: true},
		{Index: 2, RawActive: true, Active: true},
		{Index: 3, RawActive: true, Active: true},
		{Index: 4},
	}
	parents := nestedRotationParentRanges(candles, classifications, nestedRotationTestConfig())
	if len(parents) != 1 {
		t.Fatalf("parents=%d, want 1", len(parents))
	}
	parent := parents[0]
	if parent.high != 112 || parent.low != 99 || parent.matureIndex != 2 || parent.rawEndIndex != 3 {
		t.Fatalf("bad frozen parent: %+v", parent)
	}
	if parent.invalidationIndex != 4 || !parent.eligible {
		t.Fatalf("bad parent invalidation: %+v", parent)
	}
}

func TestNestedRangeRotationChildEligibilityAndSkips(t *testing.T) {
	parentCandles := []Candle{testCandle(0, 110, 120, 100, 110)}
	parents := []nestedRotationRange{{
		id: 1, matureIndex: 0, high: 120, low: 100, mid: 110,
		upperQuartile: 115, lowerQuartile: 105, width: 20, eligible: true,
		invalidationIndex: -1,
	}}
	childCandles := []Candle{
		testCandle(1, 101, 102, 100, 101),
		testCandle(2, 101, 102, 100, 101),
		testCandle(3, 110, 111, 109, 110),
		testCandle(4, 122, 125, 121, 123),
		testCandle(5, 122, 125, 121, 123),
		testCandle(6, 108, 118, 102, 110),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1, RawActive: true, Active: true},
		{Index: 2},
		{Index: 3, RawActive: true, Active: true},
		{Index: 4, RawActive: true, Active: true},
		{Index: 5, RawActive: true, Active: true},
	}
	cfg := nestedRotationTestConfig()
	cfg.ChildWidthMaxParentFraction = 0.40
	rows := nestedRotationChildRanges(childCandles, classifications, parentCandles, parents, cfg)
	if len(rows) != 2 {
		t.Fatalf("children=%d, want 2: %+v", len(rows), rows)
	}
	if !rows[0].Eligible || rows[0].CandidateSide != RangeDiscoverySideUp {
		t.Fatalf("first child should be eligible up: %+v", rows[0])
	}
	if rows[1].SkippedReason != "child_not_inside_parent" {
		t.Fatalf("second child skip=%q row=%+v", rows[1].SkippedReason, rows[1])
	}
}

func TestNestedRangeRotationUpDownEventsAndDuplicate(t *testing.T) {
	cfg := nestedRotationTestConfig()
	cfg.OutcomeHorizonBars = 3
	candles := []Candle{
		testCandle(0, 102, 104, 100, 102),
		testCandle(1, 104, 105, 103, 105),
		testCandle(2, 105, 111, 104, 110),
		testCandle(3, 110, 116, 109, 114),
		testCandle(4, 114, 116, 113, 114),
	}
	child := nestedRotationEligibleChild(RangeDiscoverySideUp)
	rows := nestedRotationEventRows(candles, []FuturesHigherTFNestedRangeRotationChildRangeRow{child}, cfg, DefaultSplits())
	if len(rows) != 2 {
		t.Fatalf("events=%d, want valid+duplicate: %+v", len(rows), rows)
	}
	if rows[0].SkippedReason != "" || rows[0].EventType != HigherTFNestedRangeRotationEventUp || !rows[0].FavorableMidpoint || !rows[0].FavorableFarQuartile {
		t.Fatalf("bad up event: %+v", rows[0])
	}
	if rows[1].SkippedReason != "duplicate_child_event" {
		t.Fatalf("bad duplicate row: %+v", rows[1])
	}

	downCandles := []Candle{
		testCandle(0, 118, 120, 116, 118),
		testCandle(1, 116, 117, 115, 115),
		testCandle(2, 115, 116, 109, 110),
		testCandle(3, 110, 111, 104, 106),
		testCandle(4, 106, 107, 104, 106),
	}
	downChild := nestedRotationEligibleChild(RangeDiscoverySideDown)
	downRows := nestedRotationEventRows(downCandles, []FuturesHigherTFNestedRangeRotationChildRangeRow{downChild}, cfg, DefaultSplits())
	if len(downRows) == 0 || downRows[0].SkippedReason != "" || downRows[0].EventType != HigherTFNestedRangeRotationEventDown || !downRows[0].FavorableMidpoint {
		t.Fatalf("bad down event: %+v", downRows)
	}
}

func TestNestedRangeRotationOutcomeAdverseFirstQuickAndMissing(t *testing.T) {
	cfg := nestedRotationTestConfig()
	cfg.OutcomeHorizonBars = 2
	cfg.QuickInvalidationBars = 1
	child := nestedRotationEligibleChild(RangeDiscoverySideUp)
	candles := []Candle{
		testCandle(0, 102, 104, 100, 102),
		testCandle(1, 104, 105, 103, 105),
		testCandle(2, 105, 111, 99, 106),
		testCandle(3, 106, 107, 105, 106),
	}
	row := nestedRotationNewEventRow(1, candles, child, 1, DefaultSplits())
	row.EventType = HigherTFNestedRangeRotationEventUp
	row.Side = RangeDiscoverySideUp
	row = nestedRotationLabelEvent(row, candles, cfg)
	if row.OutcomeLabel != HigherTFNestedRangeRotationOutcomeAdverseChildInvalidation || !row.AdverseChildInvalidation || !row.QuickInvalidation || row.FavorableMidpoint {
		t.Fatalf("expected adverse-first quick invalidation, got %+v", row)
	}

	missing := nestedRotationNewEventRow(2, candles[:2], child, 1, DefaultSplits())
	missing.EventType = HigherTFNestedRangeRotationEventUp
	missing.Side = RangeDiscoverySideUp
	missing = nestedRotationLabelEvent(missing, candles[:2], cfg)
	if !missing.MissingFuture || missing.OutcomeLabel != HigherTFNestedRangeRotationOutcomeMissingFuture {
		t.Fatalf("expected missing future, got %+v", missing)
	}
}

func TestNestedRangeRotationSummaryGatesAndStopStates(t *testing.T) {
	cfg := nestedRotationTestConfig()
	cfg.MinFullEvents = 3
	cfg.MinSplitEvents = 1
	cfg.MinSideEvents = 1
	cfg.MinFarQuartileRate = 0.20
	result := FuturesHigherTFNestedRangeRotationAuditResult{
		SourceRows:   []FuturesHigherTFNestedRangeRotationSourceRow{{ValidationStatus: "accepted"}},
		CoverageRows: []FuturesHigherTFNestedRangeRotationCoverageRow{{CoverageFactsPass: true, FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{ValidationStatus: "accepted"}}},
		EventRows: []FuturesHigherTFNestedRangeRotationEventRow{
			nestedRotationGateEvent("2021-06-01T00:00:00Z", RangeDiscoverySideUp),
			nestedRotationGateEvent("2023-06-01T00:00:00Z", RangeDiscoverySideDown),
			nestedRotationGateEvent("2025-06-01T00:00:00Z", RangeDiscoverySideUp),
		},
	}
	result.SummaryRows = SummarizeFuturesHigherTFNestedRangeRotationAudit(result, cfg, DefaultSplits())
	if got := FuturesHigherTFNestedRangeRotationAuditStopState(result); got != HigherTFNestedRangeRotationStopStateReadyForBaseline {
		t.Fatalf("stop=%s, want ready", got)
	}

	result.EventRows[0].FavorableMidpoint = false
	result.EventRows[0].FavorableFarQuartile = false
	result.EventRows[0].AdverseChildInvalidation = true
	result.EventRows[0].QuickInvalidation = true
	result.EventRows[0].OutcomeLabel = HigherTFNestedRangeRotationOutcomeAdverseChildInvalidation
	result.SummaryRows = SummarizeFuturesHigherTFNestedRangeRotationAudit(result, cfg, DefaultSplits())
	if got := FuturesHigherTFNestedRangeRotationAuditStopState(result); got != HigherTFNestedRangeRotationStopStateFailedNoBaseline {
		t.Fatalf("stop=%s, want failed", got)
	}

	empty := FuturesHigherTFNestedRangeRotationAuditResult{
		SourceRows:   []FuturesHigherTFNestedRangeRotationSourceRow{{ValidationStatus: "accepted"}},
		CoverageRows: []FuturesHigherTFNestedRangeRotationCoverageRow{{CoverageFactsPass: true, FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{ValidationStatus: "accepted"}}},
	}
	if got := FuturesHigherTFNestedRangeRotationAuditStopState(empty); got != HigherTFNestedRangeRotationStopStateNoCandidateEvents {
		t.Fatalf("empty stop=%s, want no candidate", got)
	}
}

func nestedRotationTestConfig() FuturesHigherTFNestedRangeRotationAuditConfig {
	cfg := DefaultFuturesHigherTFNestedRangeRotationAuditConfig()
	cfg.SkipSourceFactCheck = true
	cfg.SkipCoverageCountCheck = true
	cfg.DetectorLookbackBarsOverride = 1
	cfg.DetectorMinConsecutiveBars = 1
	return cfg
}

func nestedRotationEligibleChild(side string) FuturesHigherTFNestedRangeRotationChildRangeRow {
	row := FuturesHigherTFNestedRangeRotationChildRangeRow{
		ChildRangeID:             1,
		ParentRangeID:            1,
		CandidateSide:            side,
		MatureIndex:              0,
		ParentHigh:               120,
		ParentLow:                100,
		ParentMid:                110,
		ParentUpperQuartile:      115,
		ParentLowerQuartile:      105,
		ParentWidth:              20,
		Eligible:                 true,
		ChildWidthParentFraction: 0.20,
	}
	if side == RangeDiscoverySideDown {
		row.High = 120
		row.Low = 116
		row.Mid = 118
		row.Width = 4
		return row
	}
	row.High = 104
	row.Low = 100
	row.Mid = 102
	row.Width = 4
	return row
}

func nestedRotationGateEvent(closeTime string, side string) FuturesHigherTFNestedRangeRotationEventRow {
	t, _ := time.Parse(time.RFC3339, closeTime)
	return FuturesHigherTFNestedRangeRotationEventRow{
		EventCloseTime:                t.Format(timeLayout),
		Side:                          side,
		FavorableMidpoint:             true,
		FavorableFarQuartile:          true,
		OutcomeLabel:                  HigherTFNestedRangeRotationOutcomeFavorableMidpoint,
		FavorableExcursionParentWidth: 0.50,
		AdverseExcursionParentWidth:   0.10,
	}
}
