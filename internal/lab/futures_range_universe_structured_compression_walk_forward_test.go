package lab

import (
	"strings"
	"testing"
)

func TestStructuredCompressionWalkForwardFoldDefinitionsAndGridReuse(t *testing.T) {
	folds := FuturesRangeUniverseStructuredCompressionWalkForwardFolds()
	if len(folds) != 3 {
		t.Fatalf("folds=%d, want 3", len(folds))
	}
	if folds[0].FoldID != StructuredCompressionWalkForwardFoldStressToOOS || strings.Join(folds[0].TrainSplits, "+") != "2021_2022_stress" || folds[0].TestSplit != "2023_2024_oos" {
		t.Fatalf("bad first fold: %+v", folds[0])
	}
	if folds[1].FoldID != StructuredCompressionWalkForwardFoldStressOOSToRecent || strings.Join(folds[1].TrainSplits, "+") != "2021_2022_stress+2023_2024_oos" || folds[1].TestSplit != "2025_2026_recent" {
		t.Fatalf("bad second fold: %+v", folds[1])
	}
	if folds[2].FoldID != StructuredCompressionWalkForwardFoldOOSToRecent || strings.Join(folds[2].TrainSplits, "+") != "2023_2024_oos" || folds[2].TestSplit != "2025_2026_recent" {
		t.Fatalf("bad third fold: %+v", folds[2])
	}

	cfg := DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig()
	grid := structuredCompressionOptimizationGridConfigs(cfg.OptimizationConfig)
	if len(grid) != 216 {
		t.Fatalf("walk-forward grid rows=%d, want existing 216-row optimization grid", len(grid))
	}
	if cfg.FrozenConfigID != StructuredCompressionStrategyReplayConfigID {
		t.Fatalf("frozen config=%s", cfg.FrozenConfigID)
	}
}

func TestStructuredCompressionWalkForwardSelectionAndFrozenComparison(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig()
	cfg.MinTrainSplitTrades = 1
	cfg.MinMultiSplitTrainTrades = 2
	cfg.MinTestTrades = 1

	fold := FuturesRangeUniverseStructuredCompressionWalkForwardFolds()[0]
	frozenID := cfg.FrozenConfigID
	ethSolID := structuredCompressionOptimizationConfigID(StructuredCompressionOptimizationSymbolSetETHSOL, 2, 12, 1.25, 0)
	allID := structuredCompressionOptimizationConfigID(StructuredCompressionOptimizationSymbolSetAll, 2, 12, 1.25, 0)
	trades := []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow{}
	trades = append(trades, walkForwardTestTrades(frozenID, StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL, RangeUniverseSymbolETHUSDT, true, false, "2022-01-01T00:00:00Z", 5, 1)...)
	trades = append(trades, walkForwardTestTrades(frozenID, StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL, RangeUniverseSymbolSOLUSDT, true, false, "2022-01-02T00:00:00Z", 5, 2)...)
	trades = append(trades, walkForwardTestTrades(frozenID, StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL, RangeUniverseSymbolETHUSDT, true, false, "2023-01-01T00:00:00Z", 10, 3)...)
	trades = append(trades, walkForwardTestTrades(frozenID, StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL, RangeUniverseSymbolSOLUSDT, true, false, "2023-01-02T00:00:00Z", 10, 4)...)
	trades = append(trades, walkForwardTestTrades(frozenID, StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL, RangeUniverseSymbolBTCUSDT, false, true, "2023-01-03T00:00:00Z", -5, 5)...)

	trades = append(trades, walkForwardTestTrades(ethSolID, StructuredCompressionOptimizationSymbolSetETHSOL, RangeUniverseSymbolETHUSDT, true, false, "2022-02-01T00:00:00Z", 20, 6)...)
	trades = append(trades, walkForwardTestTrades(ethSolID, StructuredCompressionOptimizationSymbolSetETHSOL, RangeUniverseSymbolSOLUSDT, true, false, "2022-02-02T00:00:00Z", 20, 7)...)
	trades = append(trades, walkForwardTestTrades(ethSolID, StructuredCompressionOptimizationSymbolSetETHSOL, RangeUniverseSymbolETHUSDT, true, false, "2023-02-01T00:00:00Z", 15, 8)...)
	trades = append(trades, walkForwardTestTrades(ethSolID, StructuredCompressionOptimizationSymbolSetETHSOL, RangeUniverseSymbolSOLUSDT, true, false, "2023-02-02T00:00:00Z", 15, 9)...)

	trades = append(trades, walkForwardTestTrades(allID, StructuredCompressionOptimizationSymbolSetAll, RangeUniverseSymbolBTCUSDT, true, false, "2022-03-01T00:00:00Z", 100, 10)...)
	trades = append(trades, walkForwardTestTrades(allID, StructuredCompressionOptimizationSymbolSetAll, RangeUniverseSymbolETHUSDT, true, false, "2022-03-02T00:00:00Z", 100, 11)...)
	trades = append(trades, walkForwardTestTrades(allID, StructuredCompressionOptimizationSymbolSetAll, RangeUniverseSymbolSOLUSDT, true, false, "2022-03-03T00:00:00Z", 100, 12)...)

	rankings := structuredCompressionWalkForwardRankingRows(fold, structuredCompressionWalkForwardGridByID(structuredCompressionOptimizationGridConfigs(cfg.OptimizationConfig)), trades, cfg, 1000, DefaultSplits())
	foldRow := structuredCompressionWalkForwardFoldRow(fold, rankings, cfg)
	if foldRow.SelectedConfigID != ethSolID {
		t.Fatalf("selected=%s, want %s; fold=%+v", foldRow.SelectedConfigID, ethSolID, foldRow)
	}
	if foldRow.RequiresBTCUSDTAuthority {
		t.Fatalf("selected config should not require BTC authority: %+v", foldRow)
	}
	if !foldRow.SelectedFrozenOrEquivalent || !foldRow.PassesFoldGate {
		t.Fatalf("expected selected ETH/SOL h12 row to pass as frozen-equivalent: %+v", foldRow)
	}
	if foldRow.FrozenRank == 0 || foldRow.FrozenTestTrades != 2 || foldRow.FrozenTestNetPnL != 20 {
		t.Fatalf("bad frozen comparison row: %+v", foldRow)
	}

	metrics := structuredCompressionWalkForwardMetricsForSplits(trades, frozenID, []string{"2023_2024_oos"}, DefaultSplits(), 1000)
	if metrics.All.TotalTrades != 2 || metrics.All.NetPnL != 20 {
		t.Fatalf("authority aggregate should exclude BTC diagnostic trade: %+v", metrics.All)
	}
	if metrics.BTC.TotalTrades != 1 || metrics.BTC.NetPnL != -5 {
		t.Fatalf("BTC diagnostic metric missing: %+v", metrics.BTC)
	}

	var allRow FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow
	for _, row := range rankings {
		if row.ConfigID == allID {
			allRow = row
			break
		}
	}
	if !allRow.HistoricalComparisonOnly || allRow.PassesTrainingGate || !strings.Contains(allRow.FailureReason, "btcusdt_authority_required") {
		t.Fatalf("BTC authority row should be comparison-only: %+v", allRow)
	}
}

func TestStructuredCompressionWalkForwardStopStates(t *testing.T) {
	passing := []FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow{
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true, SelectedExactFrozen: true},
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true, SelectedExactFrozen: true},
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true},
	}
	if got := FuturesRangeUniverseStructuredCompressionWalkForwardStopState(passing); got != StructuredCompressionWalkForwardStopStatePassedCandidatePackage {
		t.Fatalf("stop=%s, want passed", got)
	}

	reviewOnly := []FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow{
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true},
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true},
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true},
	}
	if got := FuturesRangeUniverseStructuredCompressionWalkForwardStopState(reviewOnly); got != StructuredCompressionWalkForwardStopStateReviewOnlyNoStrategy {
		t.Fatalf("stop=%s, want review-only", got)
	}

	fragile := []FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow{
		{PassesFoldGate: true, SelectedFrozenOrEquivalent: true, SelectedExactFrozen: true},
		{PassesFoldGate: false},
		{PassesFoldGate: false},
	}
	if got := FuturesRangeUniverseStructuredCompressionWalkForwardStopState(fragile); got != StructuredCompressionWalkForwardStopStateFragileNeedsReview {
		t.Fatalf("stop=%s, want fragile", got)
	}
}

func walkForwardTestTrades(configID, symbolSet, symbol string, authority bool, diagnostic bool, exitTime string, net float64, seq int) []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow {
	side := Long
	if seq%2 == 0 {
		side = Short
	}
	return []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow{{
		ConfigID:          configID,
		SymbolSet:         symbolSet,
		AuthoritySymbols:  "ETHUSDT,SOLUSDT",
		DiagnosticSymbols: "BTCUSDT",
		IsAuthority:       authority,
		IsDiagnostic:      diagnostic,
		FuturesRangeUniverseStructuredCompressionTradeRow: FuturesRangeUniverseStructuredCompressionTradeRow{
			SignalID:    configID + "_" + symbol + "_" + exitTime,
			Symbol:      symbol,
			CandidateID: StructuredCompressionCandidate4HAllH6,
			Timeframe:   RangeDiscoveryTimeframe4h,
			Family:      RangeUniverseFamilyStructuredCompressionBreak,
			Side:        side,
			EntryTime:   exitTime,
			ExitTime:    exitTime,
			GrossPnL:    net,
			NetPnL:      net,
			GrossR:      net,
			NetR:        net,
			InitialRisk: 1,
			ExitReason:  "test",
			HoldBars:    1,
		},
	}}
}
