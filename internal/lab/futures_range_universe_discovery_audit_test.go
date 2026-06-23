package lab

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRangeUniverseSourceValidationSortsPhysicalNonMonotonicRows(t *testing.T) {
	dir := t.TempDir()
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(2, 102, 103, 101, 102),
		testCandle(1, 101, 102, 100, 101),
	}
	path := writeRangeUniverseTestCSV(t, dir, "solusdt_futures_um_5m_test.csv", candles)
	got, row, err := LoadFuturesRangeUniverseSource(
		FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolSOLUSDT, Path: path, ApprovedPath: path},
		rangeUniverseTestSplits(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if row.PhysicalNonMonotonicCount != 1 || !row.SortedForValidation || !row.AcceptedMonotonic || row.ValidationStatus != "accepted" {
		t.Fatalf("bad source row: %+v", row)
	}
	if len(got) != 3 || !got[0].OpenTime.Before(got[1].OpenTime) || !got[1].OpenTime.Before(got[2].OpenTime) {
		t.Fatalf("accepted candles not sorted: %+v", got)
	}
}

func TestRangeUniverseSourceValidationRejectsUnapprovedGapAndDuplicate(t *testing.T) {
	dir := t.TempDir()
	path := writeRangeUniverseTestCSV(t, dir, "ethusdt_futures_um_5m_test.csv", []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 101, 102, 100, 101),
		testCandle(2, 102, 103, 101, 102),
	})
	_, row, err := LoadFuturesRangeUniverseSource(
		FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolETHUSDT, Path: path, ApprovedPath: filepath.Join(dir, "other_ethusdt_futures_um_5m_test.csv")},
		rangeUniverseTestSplits(),
	)
	if err == nil || !strings.Contains(err.Error(), "approved local path") || row.ValidationStatus != "rejected" {
		t.Fatalf("expected approved path rejection, row=%+v err=%v", row, err)
	}

	gapped := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 101, 102, 100, 101),
		testCandle(3, 103, 104, 102, 103),
	}
	gapPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_gap.csv", gapped)
	_, row, err = LoadFuturesRangeUniverseSource(
		FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolBTCUSDT, Path: gapPath, ApprovedPath: gapPath},
		rangeUniverseTestSplits(),
	)
	if err == nil || !strings.Contains(err.Error(), "missing 5m") || row.ValidationStatus != "rejected" {
		t.Fatalf("expected gap rejection, row=%+v err=%v", row, err)
	}

	duplicate := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 101, 102, 100, 101),
		testCandle(1, 101, 102, 100, 101),
	}
	dupPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_duplicate.csv", duplicate)
	_, row, err = LoadFuturesRangeUniverseSource(
		FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolBTCUSDT, Path: dupPath, ApprovedPath: dupPath},
		rangeUniverseTestSplits(),
	)
	if err == nil || !strings.Contains(err.Error(), "duplicate") || row.ValidationStatus != "rejected" {
		t.Fatalf("expected duplicate rejection, row=%+v err=%v", row, err)
	}
}

func TestRangeUniverseNewFamilyLabels(t *testing.T) {
	cfg := DefaultFuturesRangeCandidateDiscoveryAuditConfig()
	cfg.MinMoveRangeFraction = 0.25
	cfg.QuickInvalidationBars = 1
	episode := rangeDiscoveryTestEpisode()
	frame := rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe1h}

	retestEvent := newRangeDiscoveryDirectionalEvent(1, frame, episode, 1, RangeUniverseFamilyBreakoutRetestAcceptance, RangeDiscoverySideUp, "up", 1, 1, 1, 0, RangeUniverseOutcomeRetestContinuation, RangeUniverseOutcomeRetestReentry, RangeUniverseOutcomeRetestStall)
	retestRow := newRangeDiscoveryCandidateRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 111, 113, 109, 112),
		testCandle(2, 112, 116, 111, 115),
	}, retestEvent, 1, cfg)
	if retestRow.OutcomeClass != RangeDiscoveryClassFavorable || retestRow.OutcomeLabel != RangeUniverseOutcomeRetestContinuation {
		t.Fatalf("bad retest favorable row: %+v", retestRow)
	}
	retestFail := newRangeDiscoveryCandidateRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 111, 113, 109, 112),
		testCandle(2, 112, 112, 105, 109),
	}, retestEvent, 1, cfg)
	if retestFail.OutcomeClass != RangeDiscoveryClassAdverse || retestFail.OutcomeLabel != RangeUniverseOutcomeRetestReentry || !retestFail.QuickInvalidation {
		t.Fatalf("bad retest adverse row: %+v", retestFail)
	}

	structuredEvent := newRangeDiscoveryDirectionalEvent(2, frame, episode, 1, RangeUniverseFamilyStructuredCompressionBreak, RangeDiscoverySideDown, "down", 1, 1, 1, 0, RangeUniverseOutcomeStructuredContinuation, RangeUniverseOutcomeStructuredFailedReentry, RangeUniverseOutcomeStructuredNoContinuation)
	structuredRow := newRangeDiscoveryCandidateRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 99, 100, 96, 97),
		testCandle(2, 97, 98, 93, 94),
	}, structuredEvent, 1, cfg)
	if structuredRow.OutcomeClass != RangeDiscoveryClassFavorable || structuredRow.OutcomeLabel != RangeUniverseOutcomeStructuredContinuation {
		t.Fatalf("bad structured favorable row: %+v", structuredRow)
	}
}

func TestRangeUniverseSummaryRankingAndStopStates(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseDiscoveryAuditConfig()
	cfg.Discovery.MinCandidatesPerSplit = 100
	cfg.Sources = []FuturesRangeUniverseSourceConfig{
		{Symbol: RangeUniverseSymbolBTCUSDT, Path: "btcusdt_futures_um_5m.csv", ApprovedPath: "btcusdt_futures_um_5m.csv"},
		{Symbol: RangeUniverseSymbolETHUSDT, Path: "ethusdt_futures_um_5m.csv", ApprovedPath: "ethusdt_futures_um_5m.csv"},
		{Symbol: RangeUniverseSymbolSOLUSDT, Path: "solusdt_futures_um_5m.csv", ApprovedPath: "solusdt_futures_um_5m.csv"},
	}
	rows := []FuturesRangeUniverseCandidateRow{}
	for _, symbol := range []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT} {
		for _, split := range []string{"2021_2022_stress", "2023_2024_oos", "2025_2026_recent"} {
			for i := 0; i < 110; i++ {
				rows = append(rows, rangeUniverseCandidateForGate(symbol, split, RangeDiscoveryClassFavorable, false))
			}
			for i := 0; i < 20; i++ {
				rows = append(rows, rangeUniverseCandidateForGate(symbol, split, RangeDiscoveryClassAdverse, false))
			}
			for i := 0; i < 5; i++ {
				rows = append(rows, rangeUniverseCandidateForGate(symbol, split, RangeDiscoveryClassNeutral, true))
			}
		}
	}
	summary := summarizeRangeUniverseDiscovery(rows, cfg.Discovery.RoundTripCostPct)
	stability := rangeUniverseStabilityRows(summary, DefaultSplits())
	rankings := rangeUniverseRankingRows(summary, stability, cfg, DefaultSplits())
	if got := FuturesRangeUniverseReviewStopState(rankings); got != RangeUniverseStopStateAuditReady {
		t.Fatalf("stop state=%s, want ready; rankings=%+v", got, rankings[:minInt(3, len(rankings))])
	}
	if len(rankings) == 0 || !rankings[0].BTCUSDTGatePass || rankings[0].TransferSymbolGatePassCount != 1 || !rankings[0].PassesGate {
		t.Fatalf("bad passing ranking row: %+v", rankings)
	}

	onlyTransfer := []FuturesRangeUniverseCandidateRow{}
	for _, split := range []string{"2021_2022_stress", "2023_2024_oos", "2025_2026_recent"} {
		for i := 0; i < 110; i++ {
			onlyTransfer = append(onlyTransfer, rangeUniverseCandidateForGate(RangeUniverseSymbolETHUSDT, split, RangeDiscoveryClassFavorable, false))
		}
	}
	summary = summarizeRangeUniverseDiscovery(onlyTransfer, cfg.Discovery.RoundTripCostPct)
	stability = rangeUniverseStabilityRows(summary, DefaultSplits())
	rankings = rangeUniverseRankingRows(summary, stability, cfg, DefaultSplits())
	if got := FuturesRangeUniverseReviewStopState(rankings); got != RangeUniverseStopStateNoBacktestCandidate {
		t.Fatalf("stop state=%s, want no candidate", got)
	}
}

func rangeUniverseCandidateForGate(symbol, split, class string, quick bool) FuturesRangeUniverseCandidateRow {
	label := RangeUniverseOutcomeRetestStall
	if class == RangeDiscoveryClassFavorable {
		label = RangeUniverseOutcomeRetestContinuation
	}
	if class == RangeDiscoveryClassAdverse {
		label = RangeUniverseOutcomeRetestReentry
	}
	return FuturesRangeUniverseCandidateRow{
		Symbol: symbol,
		FuturesRangeDiscoveryCandidateRow: FuturesRangeDiscoveryCandidateRow{
			Split:                    split,
			Timeframe:                RangeDiscoveryTimeframe1h,
			Family:                   RangeUniverseFamilyBreakoutRetestAcceptance,
			Side:                     RangeDiscoverySideUp,
			HorizonBars:              12,
			OutcomeClass:             class,
			OutcomeLabel:             label,
			FavorableMovePct:         0.020,
			AdverseMovePct:           0.002,
			FavorableMinusAdversePct: 0.018,
			QuickInvalidation:        quick,
			RangeWidthPct:            0.030,
		},
	}
}

func rangeUniverseTestSplits() []Split {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Split{
		{Name: "2021_2022_stress", Start: start, End: start.Add(5 * time.Minute)},
		{Name: "2023_2024_oos", Start: start.Add(5 * time.Minute), End: start.Add(10 * time.Minute)},
		{Name: "2025_2026_recent", Start: start.Add(10 * time.Minute), End: start.Add(15 * time.Minute)},
		{Name: fullSplitName},
	}
}

func writeRangeUniverseTestCSV(t *testing.T, dir string, name string, candles []Candle) string {
	t.Helper()
	path := filepath.Join(dir, name)
	data := "open_time,open,high,low,close,volume,close_time\n"
	for _, candle := range candles {
		data += fmt.Sprintf("%d,%v,%v,%v,%v,%v,%d\n",
			candle.OpenTime.UnixMilli(),
			candle.Open,
			candle.High,
			candle.Low,
			candle.Close,
			candle.Volume,
			candle.CloseTime.UnixMilli(),
		)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
