package lab

import (
	"strings"
	"testing"
	"time"
)

func TestRangeDiscoveryResampleClosedUTCAndRejectsBadParent(t *testing.T) {
	parent := make([]Candle, 10)
	for i := range parent {
		parent[i] = testCandle(i, float64(100+i), float64(105+i), float64(95+i), float64(101+i))
		parent[i].Volume = float64(i + 1)
	}
	frame := rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe15m, interval: 15 * time.Minute, childBars: 3, barsPerDay: 96}
	got, coverage, err := resampleRangeDiscoveryFrame(parent, frame)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || coverage.RowCount != 3 || coverage.PartialFinalChildBars != 1 || !coverage.PartialFinalDropped || !coverage.Complete {
		t.Fatalf("bad resample coverage: bars=%d coverage=%+v", len(got), coverage)
	}
	if !got[0].OpenTime.Equal(parent[0].OpenTime) ||
		!got[0].CloseTime.Equal(parent[0].OpenTime.Add(15*time.Minute-time.Millisecond)) ||
		got[0].Open != parent[0].Open ||
		got[0].High != parent[2].High ||
		got[0].Low != parent[0].Low ||
		got[0].Close != parent[2].Close ||
		got[0].Volume != 6 {
		t.Fatalf("bad first resampled candle: %+v", got[0])
	}

	gapped := append([]Candle(nil), parent...)
	gapped[3].OpenTime = gapped[3].OpenTime.Add(5 * time.Minute)
	gapped[3].CloseTime = gapped[3].CloseTime.Add(5 * time.Minute)
	if _, _, err := resampleRangeDiscoveryFrame(gapped, frame); err == nil || !strings.Contains(err.Error(), "missing 5m") {
		t.Fatalf("expected gap rejection, got %v", err)
	}

	duplicated := append([]Candle(nil), parent...)
	duplicated[2].OpenTime = duplicated[1].OpenTime
	duplicated[2].CloseTime = duplicated[1].CloseTime
	if _, _, err := resampleRangeDiscoveryFrame(duplicated, frame); err == nil {
		t.Fatalf("expected duplicate/monotonic rejection")
	}
}

func TestRangeDiscoveryFamilyEventsAndLabels(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 106, 100, 104),
		testCandle(2, 104, 108, 103, 107),
		testCandle(3, 107, 108, 99, 100),
	}
	episode := rangeDiscoveryTestEpisode()
	frame := rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe5m}
	cfg := FuturesRangeCandidateDiscoveryAuditConfig{
		HorizonsBars:               []int{1},
		QuickInvalidationBars:      1,
		MaxEventDelayBars:          3,
		ReentryWindowBars:          2,
		MinMoveRangeFraction:       0.25,
		LargeWickFraction:          0.40,
		RoundTripCostPct:           0.001,
		MinCandidatesPerSplit:      1,
		DetectorLookbackDays:       1,
		DetectorPercentile:         0.3,
		DetectorMinConsecutiveBars: 1,
	}
	event := newRangeDiscoveryDirectionalEvent(1, frame, episode, 1, RangeDiscoveryFamilyBoundaryTouch, RangeDiscoverySideSupport, "up", 1, 0, 0, 0, RangeDiscoveryOutcomeRejectInward, RangeDiscoveryOutcomeAcceptOutside, RangeDiscoveryOutcomeStallNone)
	row := newRangeDiscoveryCandidateRow(candles, event, 1, cfg)
	if row.OutcomeLabel != RangeDiscoveryOutcomeRejectInward || row.OutcomeClass != RangeDiscoveryClassFavorable || !row.FavorableGreaterThanAdverse {
		t.Fatalf("bad boundary favorable row: %+v", row)
	}
	adverseCandles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 106, 100, 104),
		testCandle(2, 104, 105, 98, 99),
	}
	adverse := newRangeDiscoveryCandidateRow(adverseCandles, event, 1, cfg)
	if adverse.OutcomeLabel != RangeDiscoveryOutcomeAcceptOutside || adverse.OutcomeClass != RangeDiscoveryClassAdverse {
		t.Fatalf("bad boundary adverse row: %+v", adverse)
	}

	nextEventID := 0
	wickEvents := rangeDiscoveryWickRejectionEvents([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 104, 105, 98, 103),
		testCandle(2, 103, 108, 102, 107),
	}, episode, frame, cfg, &nextEventID)
	if len(wickEvents) != 1 || wickEvents[0].family != RangeDiscoveryFamilyWickRejection || wickEvents[0].side != RangeDiscoverySideSupport {
		t.Fatalf("bad wick events: %+v", wickEvents)
	}

	nextEventID = 0
	reentryEvents := rangeDiscoveryFailedBreakReentryEvents([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 113, 106, 112),
		testCandle(2, 112, 112, 104, 108),
		testCandle(3, 108, 109, 101, 102),
	}, episode, frame, cfg, &nextEventID)
	if len(reentryEvents) != 1 || reentryEvents[0].family != RangeDiscoveryFamilyFailedBreakReentry || reentryEvents[0].reentryDelayBars != 1 {
		t.Fatalf("bad reentry events: %+v", reentryEvents)
	}
	reentryRow := newRangeDiscoveryCandidateRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 113, 106, 112),
		testCandle(2, 112, 112, 104, 108),
		testCandle(3, 108, 108, 100, 102),
	}, reentryEvents[0], 1, cfg)
	if reentryRow.OutcomeLabel != RangeDiscoveryOutcomeReentryContinuation {
		t.Fatalf("bad reentry label: %+v", reentryRow)
	}

	nextEventID = 0
	breakoutEvents := rangeDiscoveryCleanBreakoutEvents([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 113, 106, 112),
		testCandle(2, 112, 116, 111, 115),
	}, episode, rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe15m}, cfg, &nextEventID)
	if len(breakoutEvents) < 1 || breakoutEvents[0].family != RangeDiscoveryFamilyCleanBreakout || breakoutEvents[0].side != RangeDiscoverySideUp {
		t.Fatalf("bad breakout events: %+v", breakoutEvents)
	}
	breakoutRow := newRangeDiscoveryCandidateRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 113, 106, 112),
		testCandle(2, 112, 116, 111, 115),
	}, breakoutEvents[0], 1, cfg)
	if breakoutRow.OutcomeLabel != RangeDiscoveryOutcomeContinuationBeyondBreak {
		t.Fatalf("bad breakout label: %+v", breakoutRow)
	}
}

func TestRangeDiscoveryBalanceLabelsMissingFutureAndInvalidConfig(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 105, 109, 101, 106),
		testCandle(2, 106, 108, 99, 101),
	}
	episode := rangeDiscoveryTestEpisode()
	cfg := DefaultFuturesRangeCandidateDiscoveryAuditConfig()
	cfg.MinMoveRangeFraction = 0.25
	event := rangeDiscoveryEvent{
		eventID:        1,
		timeframe:      RangeDiscoveryTimeframe1h,
		family:         RangeDiscoveryFamilyBalancePersistence,
		side:           RangeDiscoverySideBalance,
		eventIndex:     0,
		episode:        episode,
		favorableLabel: RangeDiscoveryOutcomeInternalRotation,
		adverseLabel:   RangeDiscoveryOutcomeExpansionFailure,
		neutralLabel:   RangeDiscoveryOutcomeInsidePersistence,
	}
	row := newRangeDiscoveryCandidateRow(candles, event, 2, cfg)
	if row.OutcomeLabel != RangeDiscoveryOutcomeInternalRotation || row.OutcomeClass != RangeDiscoveryClassFavorable {
		t.Fatalf("bad balance label: %+v", row)
	}
	missing := newRangeDiscoveryCandidateRow(candles, event, 10, cfg)
	if !missing.MissingFuture || missing.OutcomeLabel != RangeDiscoveryOutcomeMissingFuture {
		t.Fatalf("bad missing future row: %+v", missing)
	}
	if _, _, _, _, _, err := RunFuturesRangeCandidateDiscoveryAudit(candles, FuturesRangeCandidateDiscoveryAuditConfig{QuickInvalidationBars: -1}, nil); err == nil {
		t.Fatalf("expected invalid config error")
	}
}

func TestRangeDiscoverySummaryStabilityRankingAndStopState(t *testing.T) {
	rows := []FuturesRangeDiscoveryCandidateRow{}
	for _, split := range []string{"2021_2022_stress", "2023_2024_oos", "2025_2026_recent"} {
		for i := 0; i < 100; i++ {
			rows = append(rows, rangeDiscoveryCandidateForGate(split, RangeDiscoveryClassFavorable, false))
		}
		for i := 0; i < 20; i++ {
			rows = append(rows, rangeDiscoveryCandidateForGate(split, RangeDiscoveryClassAdverse, false))
		}
		for i := 0; i < 10; i++ {
			rows = append(rows, rangeDiscoveryCandidateForGate(split, RangeDiscoveryClassNeutral, true))
		}
	}
	summary := summarizeRangeDiscovery(rows, 0.001)
	stability := rangeDiscoveryStabilityRows(summary, DefaultSplits())
	stabilityRow, ok := findRangeDiscoveryStability(stability, RangeDiscoveryTimeframe15m, RangeDiscoveryFamilyCleanBreakout, RangeDiscoverySideUp, 3)
	if !ok {
		t.Fatalf("missing stability row: %+v", stability)
	}
	if stabilityRow.PeriodSplits != 3 || stabilityRow.CandidateCount != 390 || stabilityRow.CandidateCountMin != 130 {
		t.Fatalf("bad stability row, full split may have leaked in: %+v", stabilityRow)
	}
	cfg := DefaultFuturesRangeCandidateDiscoveryAuditConfig()
	cfg.MinCandidatesPerSplit = 100
	rankings := rangeDiscoveryRankingRows(summary, stability, cfg, DefaultSplits())
	if got := FuturesRangeDiscoveryReviewStopState(rankings); got != RangeDiscoveryStopStateAuditReady {
		t.Fatalf("stop state=%s, want ready; rankings=%+v", got, rankings[:minInt(3, len(rankings))])
	}

	small := summarizeRangeDiscovery(rows[:99], 0.001)
	smallStability := rangeDiscoveryStabilityRows(small, DefaultSplits())
	smallRankings := rangeDiscoveryRankingRows(small, smallStability, cfg, DefaultSplits())
	if got := FuturesRangeDiscoveryReviewStopState(smallRankings); got != RangeDiscoveryStopStateNoBacktestCandidate {
		t.Fatalf("stop state=%s, want no candidate", got)
	}
}

func rangeDiscoveryTestEpisode() rangeRegimeDurabilityEpisode {
	return rangeRegimeDurabilityEpisode{
		EpisodeID:          1,
		Split:              fullSplitName,
		StartIndex:         0,
		EndIndex:           0,
		High:               110,
		Low:                100,
		EndClose:           105,
		RawLengthBars:      12,
		ActiveLengthBars:   12,
		RawLengthBucket:    barLengthBucket(12),
		ActiveLengthBucket: barLengthBucket(12),
		WidthPct:           movePct(10, 105),
		WidthBucket:        rangeWidthBucket(movePct(10, 105)),
		DetectorProfileID:  BalancedDetectorProfileID,
	}
}

func rangeDiscoveryCandidateForGate(split, class string, quick bool) FuturesRangeDiscoveryCandidateRow {
	label := RangeDiscoveryOutcomeNoExtension
	if class == RangeDiscoveryClassFavorable {
		label = RangeDiscoveryOutcomeContinuationBeyondBreak
	}
	if class == RangeDiscoveryClassAdverse {
		label = RangeDiscoveryOutcomeFailedBreakReentry
	}
	return FuturesRangeDiscoveryCandidateRow{
		Split:                    split,
		Timeframe:                RangeDiscoveryTimeframe15m,
		Family:                   RangeDiscoveryFamilyCleanBreakout,
		Side:                     RangeDiscoverySideUp,
		HorizonBars:              3,
		OutcomeClass:             class,
		OutcomeLabel:             label,
		FavorableMovePct:         0.010,
		AdverseMovePct:           0.002,
		FavorableMinusAdversePct: 0.008,
		QuickInvalidation:        quick,
		RangeWidthPct:            0.020,
	}
}

func findRangeDiscoveryStability(rows []FuturesRangeDiscoveryStabilityRow, timeframe, family, side string, horizon int) (FuturesRangeDiscoveryStabilityRow, bool) {
	for _, row := range rows {
		if row.Timeframe == timeframe && row.Family == family && row.Side == side && row.HorizonBars == horizon {
			return row, true
		}
	}
	return FuturesRangeDiscoveryStabilityRow{}, false
}
