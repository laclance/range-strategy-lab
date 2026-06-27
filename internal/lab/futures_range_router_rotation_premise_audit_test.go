package lab

import (
	"strings"
	"testing"
	"time"
)

func TestRangeRouterRotationPremiseSourceDependencyAndStopStates(t *testing.T) {
	source := rangeRouterRotationPremiseAcceptedSources()
	coverage := rangeRouterRotationPremiseAcceptedCoverage()
	dependency := rangeRouterRotationPremiseAcceptedDependency()
	if !rangeRouterRotationPremiseSourcePass(source, coverage, dependency) {
		t.Fatalf("expected accepted source/dependency")
	}
	badPath := rangeRouterRotationPremiseAcceptedSources()
	badPath[0].Path = "../binance-bot/data/btcusdt_5m_2021_2026.csv"
	if rangeRouterRotationPremiseSourcePass(badPath, coverage, dependency) {
		t.Fatalf("expected wrong path rejection")
	}
	badCount := rangeRouterRotationPremiseAcceptedSources()
	badCount[0].RowCount--
	if rangeRouterRotationPremiseSourcePass(badCount, coverage, dependency) {
		t.Fatalf("expected wrong count rejection")
	}

	result := FuturesRangeRouterRotationPremiseAuditResult{
		SourceRows:      source,
		CoverageRows:    coverage,
		DependencyRows:  dependency,
		RouterStopState: RangeContextRouterStopStatePassedNeedsRotationSpec,
	}
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStateNoEligibleEvents {
		t.Fatalf("stop=%s, want no eligible events", got)
	}

	result.OutcomeRows = []FuturesRangeRouterRotationPremiseOutcomeRow{{OutcomeComplete: true}}
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStateFailedNoPremise {
		t.Fatalf("stop=%s, want failed premise", got)
	}
	result.RankingRows = []FuturesRangeRouterRotationPremiseRankingRow{{PassesGate: true}}
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStatePassedNeedsTriggerAudit {
		t.Fatalf("stop=%s, want trigger audit", got)
	}
	result.RouterStopState = RangeContextRouterStopStateRejectedClosedReslice
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStateRejectedClosedReslice {
		t.Fatalf("stop=%s, want closed reslice", got)
	}

	result.RouterStopState = RangeContextRouterStopStatePassedNeedsRotationSpec
	result.SourceRows[0].Product = "Binance spot"
	result.SourceRows[0].ComparisonOnly = true
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStateSourceRouterGap {
		t.Fatalf("stop=%s, want source/router gap", got)
	}
}

func TestRangeRouterRotationPremiseDependencyRowRequiresRouterStopAndCohort(t *testing.T) {
	cfg := DefaultFuturesRangeRouterRotationPremiseAuditConfig()
	router := FuturesRangeContextRouterAuditResult{
		StopState: RangeContextRouterStopStatePassedNeedsRotationSpec,
		RankingRows: []FuturesRangeContextRouterRankingRow{{
			CohortID:                    RangeRouterRotationPremiseRequiredRouterCohortID,
			RouterLabel:                 RangeContextRouterLabelTradableRotation,
			PassesGate:                  true,
			FullPeriodRows:              892,
			WeakestSplitRows:            210,
			FullExpectedRouteHitRate:    0.599776,
			WeakestSplitExpectedHitRate: 0.554667,
			FullAdverseRouteHitRate:     0.181614,
			WorstSplitAdverseHitRate:    0.247619,
			DominantForwardLabel:        RangeContextTriageLabelContainedRotation,
			DominantForwardLabelRate:    0.385650,
		}},
	}
	row := rangeRouterRotationPremiseDependencyRow(router, cfg)
	if !row.DependencyPass {
		t.Fatalf("expected dependency pass: %+v", row)
	}

	router.StopState = RangeContextRouterStopStateFailedNoActionableRoute
	row = rangeRouterRotationPremiseDependencyRow(router, cfg)
	if row.DependencyPass || !strings.Contains(row.FailureReason, "router_stop_state_mismatch") {
		t.Fatalf("expected router stop mismatch failure: %+v", row)
	}

	router.StopState = RangeContextRouterStopStatePassedNeedsRotationSpec
	router.RankingRows = nil
	row = rangeRouterRotationPremiseDependencyRow(router, cfg)
	if row.DependencyPass || !strings.Contains(row.FailureReason, "required_rotation_router_cohort_missing_or_failed") {
		t.Fatalf("expected missing cohort failure: %+v", row)
	}
}

func TestRangeRouterRotationPremiseSegmentsEventsAndLabelsUseClosedInputs(t *testing.T) {
	cfg := DefaultFuturesRangeRouterRotationPremiseAuditConfig()
	cfg.HorizonBars = 3
	cfg.EventSearchBars = 6
	candles := rangeRouterRotationPremiseTestCandles(12, 104, 104.5, 103.5, 104)
	candles[7] = testCandle(7, 101.5, 103.2, 100.5, 102.8)
	candles[8] = testCandle(8, 103, 105.5, 102, 104)

	state := rangeRouterRotationPremiseTestState(1, 5, 100, 110)
	stateByID := map[int]FuturesRangeStateConstructionLoopStateRow{1: state, 2: state}
	routerRows := []FuturesRangeContextRouterRow{
		rangeRouterRotationPremiseTestRouterRow(1, 1, 5, 7),
		rangeRouterRotationPremiseTestRouterRow(2, 2, 6, 7),
	}
	skips := map[rangeStateSkipKey]int{}
	segments, events, outcomes := rangeRouterRotationPremiseBuildSegmentsEvents(routerRows, stateByID, candles, cfg, []Split{{Name: fullSplitName}}, skips)
	if len(segments) != 1 || len(events) != 1 || len(outcomes) != 1 {
		t.Fatalf("segments=%d events=%d outcomes=%d", len(segments), len(events), len(outcomes))
	}
	if segments[0].DuplicateRouterRowsSkipped != 1 || skips[rangeStateSkipKey{timeframe: RangeDiscoveryTimeframe15m, split: fullSplitName, reason: rangeRouterRotationPremiseSkipDuplicateContext}] != 1 {
		t.Fatalf("duplicate skip not counted: segment=%+v skips=%+v", segments[0], skips)
	}
	if segments[0].FrozenRangeHigh != 110 || segments[0].FrozenRangeLow != 100 || segments[0].FrozenBoundsKnownThroughIndex != 5 {
		t.Fatalf("bad frozen bounds: %+v", segments[0])
	}
	if events[0].EventIndex != 7 || events[0].Side != RangeRouterRotationPremiseSideLower || events[0].ForwardLabelsAsEventInput {
		t.Fatalf("bad event: %+v", events[0])
	}
	if !outcomes[0].MidlineRotationFirst || outcomes[0].PrimaryOutcome != RangeRouterRotationPremiseOutcomeMidline || outcomes[0].FutureLabelsUsedAsInput {
		t.Fatalf("bad outcome: %+v", outcomes[0])
	}

	candles[7] = testCandle(7, 108.5, 109.5, 106.8, 107.2)
	candles[8] = testCandle(8, 107, 108, 104.5, 106)
	skips = map[rangeStateSkipKey]int{}
	_, events, outcomes = rangeRouterRotationPremiseBuildSegmentsEvents([]FuturesRangeContextRouterRow{routerRows[0]}, stateByID, candles, cfg, []Split{{Name: fullSplitName}}, skips)
	if len(events) != 1 || events[0].Side != RangeRouterRotationPremiseSideUpper {
		t.Fatalf("expected upper event: %+v", events)
	}
	if !outcomes[0].MidlineRotationFirst {
		t.Fatalf("expected upper mirrored midline outcome: %+v", outcomes[0])
	}

	candles = rangeRouterRotationPremiseTestCandles(8, 104, 104.5, 103.5, 104)
	candles[7] = testCandle(7, 101.5, 103.2, 100.5, 102.8)
	_, events, outcomes = rangeRouterRotationPremiseBuildSegmentsEvents([]FuturesRangeContextRouterRow{routerRows[0]}, stateByID, candles, cfg, []Split{{Name: fullSplitName}}, map[rangeStateSkipKey]int{})
	if len(events) != 1 || len(outcomes) != 1 || !outcomes[0].MissingFuture || outcomes[0].PrimaryOutcome != RangeRouterRotationPremiseOutcomeMissingFuture {
		t.Fatalf("expected missing future outcome: events=%+v outcomes=%+v", events, outcomes)
	}
}

func TestRangeRouterRotationPremiseOutcomePrecedenceIsConservative(t *testing.T) {
	cfg := DefaultFuturesRangeRouterRotationPremiseAuditConfig()
	cfg.HorizonBars = 3
	candles := rangeRouterRotationPremiseTestCandles(5, 104, 104.5, 103.5, 104)
	event := FuturesRangeRouterRotationPremiseEventRow{
		EventID:             1,
		SegmentID:           1,
		Timeframe:           RangeDiscoveryTimeframe15m,
		HorizonBars:         3,
		Split:               fullSplitName,
		Side:                RangeRouterRotationPremiseSideLower,
		EventIndex:          1,
		FrozenRangeLow:      100,
		FrozenRangeHigh:     110,
		FrozenRangeMid:      105,
		FrozenUpperQuartile: 107.5,
		FrozenLowerQuartile: 102.5,
		FrozenRangeWidth:    10,
		BoundaryZoneWidth:   2,
	}
	candles[2] = testCandle(2, 104, 106, 98, 99)
	outcome := rangeRouterRotationPremiseOutcomeRow(candles, event, 1, cfg)
	if !outcome.HardAdverse || outcome.PrimaryOutcome != RangeRouterRotationPremiseOutcomeBoundaryFailure || outcome.MidlineRotationFirst {
		t.Fatalf("same-candle adverse should win conservatively: %+v", outcome)
	}

	candles[2] = testCandle(2, 101, 102, 100.5, 101.5)
	candles[3] = testCandle(3, 101.2, 102, 100.2, 101.6)
	candles[4] = testCandle(4, 101.1, 102, 100.1, 101.7)
	outcome = rangeRouterRotationPremiseOutcomeRow(candles, event, 1, cfg)
	if outcome.PrimaryOutcome != RangeRouterRotationPremiseOutcomeBoundaryChopNoRotation {
		t.Fatalf("expected boundary chop, got %+v", outcome)
	}
}

func TestRangeRouterRotationPremiseCohortGatesAndStateCarry(t *testing.T) {
	cfg := DefaultFuturesRangeRouterRotationPremiseAuditConfig()
	segments := []FuturesRangeRouterRotationPremiseContextSegmentRow{}
	for i := 0; i < 250; i++ {
		segments = append(segments, FuturesRangeRouterRotationPremiseContextSegmentRow{Split: rangeRouterRotationPremiseSplitForOrdinal(i), Timeframe: RangeDiscoveryTimeframe15m})
	}
	outcomes := []FuturesRangeRouterRotationPremiseOutcomeRow{}
	for i := 0; i < 150; i++ {
		split := rangeRouterRotationPremiseSplitForOrdinal(i)
		side := RangeRouterRotationPremiseSideLower
		if i%2 == 0 {
			side = RangeRouterRotationPremiseSideUpper
		}
		outcomes = append(outcomes, FuturesRangeRouterRotationPremiseOutcomeRow{
			Split:                split,
			Side:                 side,
			Timeframe:            RangeDiscoveryTimeframe15m,
			HorizonBars:          24,
			PrimaryOutcome:       RangeRouterRotationPremiseOutcomeMidline,
			MidlineRotationFirst: true,
			OutcomeComplete:      true,
			AllFamiliesID:        "state_" + string(rune('a'+i%3)),
		})
	}
	cohorts := rangeRouterRotationPremiseCohortRows(segments, outcomes, cfg, DefaultSplits(), true)
	rankings := rangeRouterRotationPremiseRankingRows(cohorts)
	if len(rankings) == 0 || !rankings[0].PassesGate {
		t.Fatalf("expected passing ranking, got %+v", rankings)
	}
	result := FuturesRangeRouterRotationPremiseAuditResult{
		SourceRows:      rangeRouterRotationPremiseAcceptedSources(),
		CoverageRows:    rangeRouterRotationPremiseAcceptedCoverage(),
		DependencyRows:  rangeRouterRotationPremiseAcceptedDependency(),
		OutcomeRows:     outcomes,
		RankingRows:     rankings,
		RouterStopState: RangeContextRouterStopStatePassedNeedsRotationSpec,
	}
	if got := FuturesRangeRouterRotationPremiseAuditStopState(result); got != RangeRouterRotationPremiseStopStatePassedNeedsTriggerAudit {
		t.Fatalf("stop=%s, want trigger audit", got)
	}

	for i := range outcomes {
		outcomes[i].AllFamiliesID = "single_state"
	}
	cohorts = rangeRouterRotationPremiseCohortRows(segments, outcomes, cfg, DefaultSplits(), true)
	rankings = rangeRouterRotationPremiseRankingRows(cohorts)
	if len(rankings) == 0 || rankings[0].PassesGate || !strings.Contains(rankings[0].FailureReason, "single_state_rollup_carry") {
		t.Fatalf("expected state carry failure, got %+v", rankings)
	}
}

func rangeRouterRotationPremiseAcceptedSources() []FuturesRangeRouterRotationPremiseSourceRow {
	return []FuturesRangeRouterRotationPremiseSourceRow{{
		SourceResamplePass:      true,
		RouterDependencyPass:    true,
		Closed15MResamplePass:   true,
		Path:                    rangeStateConstructionLoopSourcePath,
		ApprovedPath:            rangeStateConstructionLoopSourcePath,
		Product:                 "Binance USDT-M futures",
		Symbol:                  "BTCUSDT",
		Interval:                "5m",
		RowCount:                rangeStateConstructionLoopExpectedRows,
		ExpectedRowCount:        rangeStateConstructionLoopExpectedRows,
		FirstOpenTime:           rangeStateConstructionLoopExpectedFirst,
		ExpectedFirstOpenTime:   rangeStateConstructionLoopExpectedFirst,
		LastOpenTime:            rangeStateConstructionLoopExpectedLast,
		ExpectedLastOpenTime:    rangeStateConstructionLoopExpectedLast,
		GapCount:                0,
		ExpectedGapCount:        0,
		DuplicateCount:          0,
		ExpectedDuplicateCount:  0,
		ZeroVolumeCount:         rangeStateConstructionLoopExpectedZeroVol,
		ExpectedZeroVolumeCount: rangeStateConstructionLoopExpectedZeroVol,
		SourceFactsPass:         true,
		ValidationStatus:        "accepted",
	}}
}

func rangeRouterRotationPremiseAcceptedCoverage() []FuturesRangeRouterRotationPremiseCoverageRow {
	return []FuturesRangeRouterRotationPremiseCoverageRow{{
		Timeframe:          RangeDiscoveryTimeframe15m,
		SourceResamplePass: true,
		CoverageFactsPass:  true,
		Complete:           true,
		ValidationStatus:   "accepted",
	}}
}

func rangeRouterRotationPremiseAcceptedDependency() []FuturesRangeRouterRotationPremiseRouterDependencyRow {
	return []FuturesRangeRouterRotationPremiseRouterDependencyRow{{
		RequiredRouterStopState:  RangeContextRouterStopStatePassedNeedsRotationSpec,
		ActualRouterStopState:    RangeContextRouterStopStatePassedNeedsRotationSpec,
		DependencyPass:           true,
		RequiredCohortID:         RangeRouterRotationPremiseRequiredRouterCohortID,
		RequiredCohortPassesGate: true,
	}}
}

func rangeRouterRotationPremiseTestCandles(n int, open, high, low, close float64) []Candle {
	candles := make([]Candle, n)
	for i := range candles {
		openTime := time.Date(2021, 1, 1, 0, i*15, 0, 0, time.UTC)
		candles[i] = Candle{
			OpenTime:  openTime,
			CloseTime: openTime.Add(15*time.Minute - time.Millisecond),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    1,
		}
	}
	return candles
}

func rangeRouterRotationPremiseTestState(id, decision int, low, high float64) FuturesRangeStateConstructionLoopStateRow {
	width := high - low
	return FuturesRangeStateConstructionLoopStateRow{
		StateRowID:        id,
		Timeframe:         RangeDiscoveryTimeframe15m,
		Split:             fullSplitName,
		RangeEpisodeID:    7,
		RangeStartIndex:   1,
		DecisionIndex:     decision,
		RangeStartTime:    "2021-01-01T00:14:59Z",
		DecisionCloseTime: rangeRouterRotationPremiseTestCandles(decision+1, 104, 104.5, 103.5, 104)[decision].CloseTime.UTC().Format(timeLayout),
		RangeHigh:         high,
		RangeLow:          low,
		RangeMid:          low + width/2,
		RangeWidth:        width,
		StateID:           "state",
		AllFamiliesID:     "state",
		Eligible:          true,
	}
}

func rangeRouterRotationPremiseTestRouterRow(routerID, stateID, decision, episode int) FuturesRangeContextRouterRow {
	return FuturesRangeContextRouterRow{
		RouterRowID:                routerID,
		StateRowID:                 stateID,
		Timeframe:                  RangeDiscoveryTimeframe15m,
		Split:                      fullSplitName,
		RangeEpisodeID:             episode,
		RangeStartIndex:            1,
		DecisionIndex:              decision,
		RangeStartTime:             "2021-01-01T00:14:59Z",
		DecisionCloseTime:          rangeRouterRotationPremiseTestCandles(decision+1, 104, 104.5, 103.5, 104)[decision].CloseTime.UTC().Format(timeLayout),
		StateID:                    "state",
		AllFamiliesID:              "state",
		RouterLabel:                RangeContextRouterLabelTradableRotation,
		RouterReason:               "matched_passing_rule",
		Actionable:                 true,
		MatchedRuleCount:           1,
		MatchedRuleIDs:             rangeRouterRotationPremiseRulePrefix + "all_families|state",
		ClosedCandleOnly:           true,
		ForwardLabelsAsRouterInput: false,
		ForwardLabelColumnsPresent: false,
	}
}

func rangeRouterRotationPremiseSplitForOrdinal(i int) string {
	switch i % 3 {
	case 0:
		return "2021_2022_stress"
	case 1:
		return "2023_2024_oos"
	default:
		return "2025_2026_recent"
	}
}
