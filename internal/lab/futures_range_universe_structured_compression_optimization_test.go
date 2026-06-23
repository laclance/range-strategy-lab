package lab

import (
	"path/filepath"
	"testing"
)

func TestStructuredCompressionOptimizationGridBoundsAndSymbolSets(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig()
	grid := structuredCompressionOptimizationGridConfigs(cfg)
	if len(grid) != 216 {
		t.Fatalf("grid rows=%d, want 216", len(grid))
	}

	foundDiagnostic := false
	for _, row := range grid {
		if row.SymbolSet != StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL {
			continue
		}
		foundDiagnostic = true
		if !stringInSlice(RangeUniverseSymbolETHUSDT, row.AuthoritySymbols) || !stringInSlice(RangeUniverseSymbolSOLUSDT, row.AuthoritySymbols) || stringInSlice(RangeUniverseSymbolBTCUSDT, row.AuthoritySymbols) {
			t.Fatalf("bad diagnostic authority symbols: %+v", row)
		}
		if !stringInSlice(RangeUniverseSymbolBTCUSDT, row.DiagnosticSymbols) {
			t.Fatalf("BTC should be diagnostic: %+v", row)
		}
	}
	if !foundDiagnostic {
		t.Fatalf("expected diagnostic symbol-set rows")
	}
}

func TestStructuredCompressionOptimizationTargetAndStopBufferGeometry(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 116, 111, 115),
		testCandle(3, 116, 127, 115, 126),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}
	cfg := DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
	cfg.EventDelayBars = 1
	cfg.ConfirmationWindowBars = 1
	candidate := FuturesRangeUniverseStructuredCompressionCandidateConfig{
		CandidateID:                  StructuredCompressionCandidate4HAllH6,
		Timeframe:                    RangeDiscoveryTimeframe4h,
		Side:                         RangeDiscoverySideAll,
		MaxHoldBars:                  2,
		TargetRangeWidthMultiple:     0.75,
		StopBoundaryBufferRangeWidth: 0.10,
	}
	strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles, RangeUniverseSymbolETHUSDT, candidate, cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1", len(signals))
	}
	signal := signals[0]
	if signal.Stop != 109 || signal.Target != 123.5 {
		t.Fatalf("unexpected buffered long geometry: stop=%v target=%v row=%+v", signal.Stop, signal.Target, signal)
	}
}

func TestStructuredCompressionOptimizationSourceReuseAndResample(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 48)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_optimization.csv", parent)
	cfg := DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig()
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: path, ApprovedPath: path, SkipSplitEligibilityCheck: true}}
	cfg.ConfirmationWindowBars = []int{2}
	cfg.MaxHoldBars = []int{4}
	cfg.TargetRangeWidthMultiples = []float64{1}
	cfg.StopBoundaryBufferRangeWidths = []float64{0}
	cfg.SymbolSets = []string{StructuredCompressionOptimizationSymbolSetAll}
	cfg.DetectorMinConsecutiveBars = 1
	result, err := RunFuturesRangeUniverseStructuredCompressionOptimization(cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 1 || result.SourceRows[0].ValidationStatus != "accepted" {
		t.Fatalf("bad source rows: %+v", result.SourceRows)
	}
	if len(result.CoverageRows) != 1 || result.CoverageRows[0].Timeframe != RangeDiscoveryTimeframe4h || result.CoverageRows[0].RowCount != 1 {
		t.Fatalf("bad coverage rows: %+v", result.CoverageRows)
	}
	if len(result.GridRows) != 1 || len(result.RankingRows) != 1 {
		t.Fatalf("bad grid/ranking rows: grid=%d rankings=%d", len(result.GridRows), len(result.RankingRows))
	}

	badPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_spot_5m_optimization.csv", parent)
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: badPath, ApprovedPath: badPath, SkipSplitEligibilityCheck: true}}
	if _, err := RunFuturesRangeUniverseStructuredCompressionOptimization(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected spot-looking source rejection")
	}
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: path, ApprovedPath: filepath.Join(dir, "other_btcusdt_futures_um_5m.csv"), SkipSplitEligibilityCheck: true}}
	if _, err := RunFuturesRangeUniverseStructuredCompressionOptimization(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected approved-path rejection")
	}
}

func TestStructuredCompressionOptimizationRankingAndStopStates(t *testing.T) {
	rows := []FuturesRangeUniverseStructuredCompressionOptimizationGridRow{
		{ConfigID: "near_a", NearViableForPortfolio: true, RankScore: 10},
		{ConfigID: "pass_b", PassesGate: true, RankScore: 5},
		{ConfigID: "near_c", NearViableForPortfolio: true, RankScore: 7},
	}
	rankings := structuredCompressionOptimizationRankingRows(rows)
	if rankings[0].ConfigID != "pass_b" || rankings[0].Rank != 1 {
		t.Fatalf("passing row should rank first: %+v", rankings)
	}
	if got := FuturesRangeUniverseStructuredCompressionOptimizationStopState(rankings); got != StructuredCompressionOptimizationStopStatePassedStrategySpec {
		t.Fatalf("stop=%s, want pass", got)
	}
	mixed := structuredCompressionOptimizationRankingRows([]FuturesRangeUniverseStructuredCompressionOptimizationGridRow{
		{ConfigID: "near_a", NearViableForPortfolio: true, RankScore: 10},
		{ConfigID: "near_b", NearViableForPortfolio: true, RankScore: 9},
	})
	if got := FuturesRangeUniverseStructuredCompressionOptimizationStopState(mixed); got != StructuredCompressionOptimizationStopStateMixedPortfolioReview {
		t.Fatalf("stop=%s, want mixed", got)
	}
	failed := structuredCompressionOptimizationRankingRows([]FuturesRangeUniverseStructuredCompressionOptimizationGridRow{{ConfigID: "failed"}})
	if got := FuturesRangeUniverseStructuredCompressionOptimizationStopState(failed); got != StructuredCompressionOptimizationStopStateFailedNoPromotion {
		t.Fatalf("stop=%s, want failed", got)
	}
}
