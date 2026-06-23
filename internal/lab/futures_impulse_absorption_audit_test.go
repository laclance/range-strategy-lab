package lab

import "testing"

func TestFuturesImpulseAbsorptionDefaults(t *testing.T) {
	cfg := DefaultFuturesImpulseAbsorptionAuditConfig()
	if cfg.WarmupBars != 8640 ||
		cfg.TrueRangePercentileThreshold != 0.99 ||
		cfg.VolumePercentileThreshold != 0.95 ||
		cfg.DownClosePositionMax != 0.25 ||
		cfg.UpClosePositionMin != 0.75 ||
		cfg.QuickContinuationBars != 3 {
		t.Fatalf("bad defaults: %+v", cfg)
	}
	if len(cfg.HorizonsBars) != 4 ||
		cfg.HorizonsBars[0] != 3 ||
		cfg.HorizonsBars[1] != 6 ||
		cfg.HorizonsBars[2] != 12 ||
		cfg.HorizonsBars[3] != 24 {
		t.Fatalf("bad default horizons: %+v", cfg.HorizonsBars)
	}
}

func TestRollingPriorPercentileRanksExcludeCurrentBar(t *testing.T) {
	got := rollingPriorPercentileRanks([]float64{1, 100, 2}, 2)
	if !almostEqual(got[2], 0.5) {
		t.Fatalf("rank=%v, want 0.5 from prior [1,100] only", got[2])
	}
}

func TestFuturesImpulseAbsorptionDirectionCutoffsAndZeroRangeSkip(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100),
		testCandle(3, 100, 110, 100, 107.5),
		testCandle(4, 107.5, 108, 104, 105),
		testCandle(5, 105, 106, 104, 105),
		testCandle(6, 105, 115, 105, 110),
		testCandle(7, 107.5, 108, 104, 105),
		testCandle(8, 105, 110, 100, 102.5),
		testCandle(9, 102.5, 106, 100, 105),
		testCandle(10, 105, 105, 105, 105),
		testCandle(11, 105, 106, 104, 105),
	}
	for i := range candles {
		candles[i].Volume = float64(i + 1)
	}
	cfg := FuturesImpulseAbsorptionAuditConfig{
		WarmupBars:                   3,
		TrueRangePercentileThreshold: 0.5,
		VolumePercentileThreshold:    0.5,
		DownClosePositionMax:         0.25,
		UpClosePositionMin:           0.75,
		HorizonsBars:                 []int{1},
		QuickContinuationBars:        1,
	}

	rows, _, _, err := RunFuturesImpulseAbsorptionAudit(candles, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	directionsByIndex := map[int]string{}
	for _, row := range rows {
		directionsByIndex[row.EventIndex] = row.Direction
	}
	if directionsByIndex[3] != FuturesImpulseAbsorptionDirectionUp {
		t.Fatalf("index 3 direction=%q, want up at exact 0.75 cutoff; rows=%+v", directionsByIndex[3], rows)
	}
	if directionsByIndex[6] != "" {
		t.Fatalf("index 6 exact-mid event should be skipped; rows=%+v", rows)
	}
	if directionsByIndex[8] != FuturesImpulseAbsorptionDirectionDown {
		t.Fatalf("index 8 direction=%q, want down at exact 0.25 cutoff; rows=%+v", directionsByIndex[8], rows)
	}
	if directionsByIndex[10] != "" {
		t.Fatalf("zero-range event should be skipped; rows=%+v", rows)
	}
}

func TestFuturesImpulseAbsorptionLabelsBothDirections(t *testing.T) {
	upCandles := []Candle{
		testCandle(0, 100, 110, 100, 108),
		testCandle(1, 108, 109, 104, 106),
		testCandle(2, 106, 111, 105, 107),
	}
	upEvent := futuresImpulseAbsorptionEvent{
		EventIndex:    0,
		Direction:     FuturesImpulseAbsorptionDirectionUp,
		EventMidpoint: 105,
	}
	upLabel := newFuturesImpulseAbsorptionLabel(upCandles, upEvent, 2, 1)
	if upLabel.FirstOutcome != FuturesImpulseAbsorptionOutcomeReclaimFirst ||
		!upLabel.MidpointReclaim ||
		!upLabel.ContinuationBeyondExtreme ||
		upLabel.BarsToReclaim != 1 ||
		upLabel.BarsToContinuation != 2 ||
		upLabel.QuickContinuation {
		t.Fatalf("bad up label: %+v", upLabel)
	}

	downCandles := []Candle{
		testCandle(0, 100, 110, 100, 102),
		testCandle(1, 102, 104, 99, 101),
	}
	downEvent := futuresImpulseAbsorptionEvent{
		EventIndex:    0,
		Direction:     FuturesImpulseAbsorptionDirectionDown,
		EventMidpoint: 105,
	}
	downLabel := newFuturesImpulseAbsorptionLabel(downCandles, downEvent, 1, 3)
	if downLabel.FirstOutcome != FuturesImpulseAbsorptionOutcomeContinuationFirst ||
		downLabel.MidpointReclaim ||
		!downLabel.ContinuationBeyondExtreme ||
		!downLabel.QuickContinuation ||
		downLabel.BarsToContinuation != 1 {
		t.Fatalf("bad down label: %+v", downLabel)
	}
}

func TestFuturesImpulseAbsorptionSameBarAmbiguityAndMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 110, 100, 108),
		testCandle(1, 108, 111, 104, 106),
	}
	event := futuresImpulseAbsorptionEvent{
		EventIndex:    0,
		Direction:     FuturesImpulseAbsorptionDirectionUp,
		EventMidpoint: 105,
	}
	sameBar := newFuturesImpulseAbsorptionLabel(candles, event, 1, 1)
	if sameBar.FirstOutcome != FuturesImpulseAbsorptionOutcomeSameBarAmbiguous ||
		!sameBar.SameBarAmbiguous ||
		!sameBar.MidpointReclaim ||
		!sameBar.ContinuationBeyondExtreme {
		t.Fatalf("bad same-bar label: %+v", sameBar)
	}

	missing := newFuturesImpulseAbsorptionLabel(candles, event, 3, 1)
	if missing.FirstOutcome != FuturesImpulseAbsorptionOutcomeMissingFuture ||
		!missing.MissingFuture ||
		missing.LabelWindowStartIndex != 1 ||
		missing.LabelWindowEndIndex != 1 {
		t.Fatalf("bad missing-future label: %+v", missing)
	}
}

func TestFuturesImpulseAbsorptionSummaryRollupsAndStability(t *testing.T) {
	rows := []FuturesImpulseAbsorptionCandidateRow{
		impulseCandidate("2021_2022_stress", FuturesImpulseAbsorptionDirectionUp, 3, "p99_99_5", "p95_99", FuturesImpulseAbsorptionOutcomeReclaimFirst),
		impulseCandidate("2021_2022_stress", FuturesImpulseAbsorptionDirectionDown, 3, "p99_5_99_9", "p99_99_5", FuturesImpulseAbsorptionOutcomeContinuationFirst),
		impulseCandidate("2023_2024_oos", FuturesImpulseAbsorptionDirectionUp, 3, "p99_99_5", "p95_99", FuturesImpulseAbsorptionOutcomeSameBarAmbiguous),
		impulseCandidate("2025_2026_recent", FuturesImpulseAbsorptionDirectionDown, 3, "p99_5_99_9", "p99_99_5", FuturesImpulseAbsorptionOutcomeReclaimFirst),
	}
	rows[1].QuickContinuation = true
	rows[2].SameBarAmbiguous = true
	rows[3].MissingFuture = true
	rows[3].FirstOutcome = FuturesImpulseAbsorptionOutcomeMissingFuture
	rows[3].MidpointReclaim = false
	rows[3].ContinuationBeyondExtreme = false

	summary := summarizeFuturesImpulseAbsorption(rows)
	fullAll, ok := findImpulseSummary(summary, fullSplitName, FuturesImpulseAbsorptionDirectionAll, 3, FuturesImpulseAbsorptionBucketAll, FuturesImpulseAbsorptionBucketAll)
	if !ok {
		t.Fatalf("missing full all/all summary: %+v", summary)
	}
	if fullAll.SourceEventCount != 4 ||
		fullAll.LabeledEventCount != 3 ||
		fullAll.MissingFutureCount != 1 ||
		fullAll.ReclaimFirstCount != 1 ||
		fullAll.ContinuationFirstCount != 1 ||
		fullAll.SameBarAmbiguousCount != 1 ||
		fullAll.QuickContinuationCount != 1 {
		t.Fatalf("bad full rollup: %+v", fullAll)
	}

	stability := futuresImpulseAbsorptionStabilityRows(summary, DefaultSplits())
	allStability, ok := findImpulseStability(stability, FuturesImpulseAbsorptionDirectionAll, 3, FuturesImpulseAbsorptionBucketAll, FuturesImpulseAbsorptionBucketAll)
	if !ok {
		t.Fatalf("missing all/all stability row: %+v", stability)
	}
	if allStability.PeriodSplits != 3 ||
		allStability.SourceEventCount != 4 ||
		allStability.SourceEventCountMin != 1 ||
		allStability.SourceEventCountMax != 2 {
		t.Fatalf("bad stability row, full split may have leaked in: %+v", allStability)
	}
}

func TestFuturesImpulseAbsorptionReviewStopState(t *testing.T) {
	rows := []FuturesImpulseAbsorptionSummaryRow{
		impulseSummaryForGate("2021_2022_stress", 3, 100, 0.60, 0.30, 0.20, 0.10),
		impulseSummaryForGate("2023_2024_oos", 3, 110, 0.55, 0.35, 0.20, 0.05),
		impulseSummaryForGate("2025_2026_recent", 3, 120, 0.58, 0.32, 0.18, 0.10),
	}
	if got := FuturesImpulseAbsorptionReviewStopState(rows, DefaultSplits()); got != FuturesImpulseAbsorptionStopStateAuditReady {
		t.Fatalf("stop state=%s, want audit_ready", got)
	}
	rows[2].SourceEventCount = 99
	if got := FuturesImpulseAbsorptionReviewStopState(rows, DefaultSplits()); got != FuturesImpulseAbsorptionStopStateNoViableEdge {
		t.Fatalf("stop state=%s, want no_viable_edge", got)
	}
}

func TestFuturesImpulseAbsorptionRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	if _, _, _, err := RunFuturesImpulseAbsorptionAudit(candles, FuturesImpulseAbsorptionAuditConfig{WarmupBars: -1}, nil); err == nil {
		t.Fatalf("expected invalid warmup error")
	}
	if _, _, _, err := RunFuturesImpulseAbsorptionAudit(candles, FuturesImpulseAbsorptionAuditConfig{
		WarmupBars:                   1,
		TrueRangePercentileThreshold: 1,
		VolumePercentileThreshold:    0.5,
		DownClosePositionMax:         0.25,
		UpClosePositionMin:           0.75,
		HorizonsBars:                 []int{1},
		QuickContinuationBars:        1,
	}, nil); err == nil {
		t.Fatalf("expected invalid threshold error")
	}
	if _, _, _, err := RunFuturesImpulseAbsorptionAudit(candles, FuturesImpulseAbsorptionAuditConfig{
		WarmupBars:                   1,
		TrueRangePercentileThreshold: 0.5,
		VolumePercentileThreshold:    0.5,
		DownClosePositionMax:         0.75,
		UpClosePositionMin:           0.75,
		HorizonsBars:                 []int{1},
		QuickContinuationBars:        1,
	}, nil); err == nil {
		t.Fatalf("expected invalid cutoff error")
	}
}

func impulseCandidate(split, direction string, horizon int, trBucket, volumeBucket, outcome string) FuturesImpulseAbsorptionCandidateRow {
	row := FuturesImpulseAbsorptionCandidateRow{
		Split:                     split,
		Direction:                 direction,
		HorizonBars:               horizon,
		TrueRangePercentileBucket: trBucket,
		VolumePercentileBucket:    volumeBucket,
		FirstOutcome:              outcome,
		MidpointReclaim:           outcome == FuturesImpulseAbsorptionOutcomeReclaimFirst || outcome == FuturesImpulseAbsorptionOutcomeSameBarAmbiguous,
		ContinuationBeyondExtreme: outcome == FuturesImpulseAbsorptionOutcomeContinuationFirst || outcome == FuturesImpulseAbsorptionOutcomeSameBarAmbiguous,
		BarsToReclaim:             1,
		BarsToContinuation:        1,
	}
	return row
}

func findImpulseSummary(rows []FuturesImpulseAbsorptionSummaryRow, split, direction string, horizon int, trBucket, volumeBucket string) (FuturesImpulseAbsorptionSummaryRow, bool) {
	for _, row := range rows {
		if row.Split == split &&
			row.Direction == direction &&
			row.HorizonBars == horizon &&
			row.TrueRangePercentileBucket == trBucket &&
			row.VolumePercentileBucket == volumeBucket {
			return row, true
		}
	}
	return FuturesImpulseAbsorptionSummaryRow{}, false
}

func findImpulseStability(rows []FuturesImpulseAbsorptionStabilityRow, direction string, horizon int, trBucket, volumeBucket string) (FuturesImpulseAbsorptionStabilityRow, bool) {
	for _, row := range rows {
		if row.Direction == direction &&
			row.HorizonBars == horizon &&
			row.TrueRangePercentileBucket == trBucket &&
			row.VolumePercentileBucket == volumeBucket {
			return row, true
		}
	}
	return FuturesImpulseAbsorptionStabilityRow{}, false
}

func impulseSummaryForGate(split string, horizon, count int, reclaimRate, continuationRate, quickRate, ambiguousRate float64) FuturesImpulseAbsorptionSummaryRow {
	return FuturesImpulseAbsorptionSummaryRow{
		Split:                     split,
		Direction:                 FuturesImpulseAbsorptionDirectionAll,
		HorizonBars:               horizon,
		TrueRangePercentileBucket: FuturesImpulseAbsorptionBucketAll,
		VolumePercentileBucket:    FuturesImpulseAbsorptionBucketAll,
		SourceEventCount:          count,
		LabeledEventCount:         count,
		ReclaimFirstRate:          reclaimRate,
		ContinuationFirstRate:     continuationRate,
		QuickContinuationRate:     quickRate,
		SameBarAmbiguousRate:      ambiguousRate,
	}
}
