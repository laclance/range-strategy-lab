package lab

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStructuredCompressionStrategyReplayFixedConfigAndNoGrid(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	if cfg.ConfigID != StructuredCompressionStrategyReplayConfigID {
		t.Fatalf("config id=%s", cfg.ConfigID)
	}
	if cfg.SymbolSet != StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL {
		t.Fatalf("symbol set=%s", cfg.SymbolSet)
	}
	if cfg.ConfirmationWindowBars != 2 || cfg.MaxHoldBars != 12 || cfg.TargetRangeWidthMultiple != 1 || cfg.StopBoundaryBufferRangeWidth != 0 {
		t.Fatalf("unexpected fixed replay knobs: %+v", cfg)
	}
	if cfg.DetectorLookbackDays != 20 || cfg.DetectorPercentile != 0.30 || cfg.DetectorMinConsecutiveBars != 12 {
		t.Fatalf("unexpected detector config: %+v", cfg)
	}

	bad := cfg
	bad.SymbolSet = StructuredCompressionOptimizationSymbolSetAll
	if err := bad.validate(); err == nil {
		t.Fatalf("expected non-frozen symbol set rejection")
	}
	bad = cfg
	bad.ConfigID = "other"
	if err := bad.validate(); err == nil {
		t.Fatalf("expected non-frozen config id rejection")
	}
}

func TestStructuredCompressionStrategyReplayRunnerValidatesSourcesAndResamples(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 48)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	btc := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_strategy.csv", parent)
	eth := writeRangeUniverseTestCSV(t, dir, "ethusdt_futures_um_5m_strategy.csv", parent)
	sol := writeRangeUniverseTestCSV(t, dir, "solusdt_futures_um_5m_strategy.csv", parent)
	cfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	cfg.Sources = []FuturesRangeUniverseSourceConfig{
		{Symbol: RangeUniverseSymbolBTCUSDT, Path: btc, ApprovedPath: btc, SkipSplitEligibilityCheck: true},
		{Symbol: RangeUniverseSymbolETHUSDT, Path: eth, ApprovedPath: eth, SkipSplitEligibilityCheck: true},
		{Symbol: RangeUniverseSymbolSOLUSDT, Path: sol, ApprovedPath: sol, SkipSplitEligibilityCheck: true},
	}
	cfg.DetectorMinConsecutiveBars = 1
	result, err := RunFuturesRangeUniverseStructuredCompressionStrategyReplay(cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 3 || len(result.CoverageRows) != 3 {
		t.Fatalf("source rows=%d coverage rows=%d", len(result.SourceRows), len(result.CoverageRows))
	}
	for _, row := range result.CoverageRows {
		if row.Timeframe != RangeDiscoveryTimeframe4h || row.RowCount != 1 || row.ValidationStatus != "accepted" {
			t.Fatalf("bad coverage row: %+v", row)
		}
	}
	if len(result.SummaryRows) != 48 {
		t.Fatalf("summary rows=%d, want 48", len(result.SummaryRows))
	}

	badPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_spot_5m_strategy.csv", parent)
	cfg.Sources[0] = FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolBTCUSDT, Path: badPath, ApprovedPath: badPath, SkipSplitEligibilityCheck: true}
	if _, err := RunFuturesRangeUniverseStructuredCompressionStrategyReplay(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected spot-looking source rejection")
	}
	cfg.Sources[0] = FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolBTCUSDT, Path: btc, ApprovedPath: filepath.Join(dir, "other_btcusdt_futures_um_5m.csv"), SkipSplitEligibilityCheck: true}
	if _, err := RunFuturesRangeUniverseStructuredCompressionStrategyReplay(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected approved-path rejection")
	}
}

func TestStructuredCompressionStrategyReplaySignalsUseFrozenCandidate(t *testing.T) {
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
	cfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	baseCfg := cfg.baselineConfig()
	candidate := baseCfg.Candidates[0]
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 12}
	strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles, RangeUniverseSymbolETHUSDT, candidate, baseCfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1", len(signals))
	}
	signal := signals[0]
	if signal.CandidateID != StructuredCompressionCandidate4HAllH12 || signal.MaxHoldBars != 12 || signal.Side != Long || signal.EntryIndex != 3 || signal.Stop != 110 || signal.Target != 126 {
		t.Fatalf("bad frozen signal: %+v", signal)
	}
	result := RunBacktest(candles, strategy, btCfg)
	if len(result.Trades) != 1 || result.Trades[0].Reason != "take_profit" {
		t.Fatalf("bad frozen replay trade: %+v", result.Trades)
	}
}

func TestStructuredCompressionStrategyReplaySummaryAuthorityOnlyAggregate(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	splits := []Split{{Name: fullSplitName}}
	signals := []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow{
		replayTestSignal(RangeUniverseSymbolBTCUSDT, false, true),
		replayTestSignal(RangeUniverseSymbolETHUSDT, true, false),
		replayTestSignal(RangeUniverseSymbolSOLUSDT, true, false),
	}
	trades := []FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow{
		replayTestTrade(RangeUniverseSymbolBTCUSDT, false, true, 1),
		replayTestTrade(RangeUniverseSymbolETHUSDT, true, false, 2),
		replayTestTrade(RangeUniverseSymbolSOLUSDT, true, false, 3),
	}
	cfg.ExpectedFullTrades = 2
	cfg.ExpectedOOSTrades = 0
	cfg.ExpectedRecentTrades = 0
	cfg.ExpectedFullNetPnL = 20
	cfg.ExpectedFullProfitFactor = 999.99
	cfg.ExpectedFullMaxDrawdown = 0
	summary := SummarizeFuturesRangeUniverseStructuredCompressionStrategyReplay(signals, trades, cfg, 1000, splits)
	byKey := structuredCompressionStrategyReplaySummaryByKey(summary)
	aggregate := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	btc := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, RangeUniverseSymbolBTCUSDT, fullSplitName, "all")]
	if aggregate.TotalTrades != 2 || aggregate.SignalCount != 2 || aggregate.NetPnL != 20 {
		t.Fatalf("aggregate should be authority-only: %+v", aggregate)
	}
	if btc.TotalTrades != 1 || btc.IsAuthority || !btc.IsDiagnostic {
		t.Fatalf("BTC should remain diagnostic-only: %+v", btc)
	}
}

func TestStructuredCompressionStrategyReplayStopStates(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	passing := replayStopStateRows(cfg, 129, 54, 43, 32, 573.87, 1.8089, 0.0982)
	if got := FuturesRangeUniverseStructuredCompressionStrategyReplayStopState(passing, cfg, 1000, DefaultSplits()); got != StructuredCompressionStrategyReplayStopStatePassedWalkForward {
		t.Fatalf("stop=%s, want pass", got)
	}

	mismatch := replayStopStateRows(cfg, 128, 54, 43, 32, 573.87, 1.8089, 0.0982)
	if got := FuturesRangeUniverseStructuredCompressionStrategyReplayStopState(mismatch, cfg, 1000, DefaultSplits()); got != StructuredCompressionStrategyReplayStopStateRegressionOrMismatch {
		t.Fatalf("stop=%s, want mismatch", got)
	}

	failedCfg := cfg
	failedCfg.ExpectedFullTrades = 20
	failedCfg.ExpectedOOSTrades = 10
	failedCfg.ExpectedRecentTrades = 10
	failedCfg.ExpectedFullNetPnL = -20
	failedCfg.ExpectedFullProfitFactor = 0.8
	failedCfg.ExpectedFullMaxDrawdown = 0.05
	failed := replayStopStateRows(failedCfg, 20, 10, 10, 10, -20, 0.8, 0.05)
	if got := FuturesRangeUniverseStructuredCompressionStrategyReplayStopState(failed, failedCfg, 1000, DefaultSplits()); got != StructuredCompressionStrategyReplayStopStateFailedNoPromotion {
		t.Fatalf("stop=%s, want failed", got)
	}
}

func replayTestSignal(symbol string, authority bool, diagnostic bool) FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow {
	return FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow{
		ConfigID:          StructuredCompressionStrategyReplayConfigID,
		SymbolSet:         StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL,
		AuthoritySymbols:  "ETHUSDT,SOLUSDT",
		DiagnosticSymbols: "BTCUSDT",
		IsAuthority:       authority,
		IsDiagnostic:      diagnostic,
		FuturesRangeUniverseStructuredCompressionSignalRow: FuturesRangeUniverseStructuredCompressionSignalRow{
			SignalID:              strings.ToLower(symbol) + "_signal",
			Symbol:                symbol,
			CandidateID:           StructuredCompressionCandidate4HAllH12,
			Timeframe:             RangeDiscoveryTimeframe4h,
			Family:                RangeUniverseFamilyStructuredCompressionBreak,
			ConfirmationCloseTime: "2024-01-01T04:00:00Z",
			Side:                  Long,
		},
	}
}

func replayTestTrade(symbol string, authority bool, diagnostic bool, sequence int) FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow {
	return FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow{
		ConfigID:          StructuredCompressionStrategyReplayConfigID,
		SymbolSet:         StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL,
		AuthoritySymbols:  "ETHUSDT,SOLUSDT",
		DiagnosticSymbols: "BTCUSDT",
		IsAuthority:       authority,
		IsDiagnostic:      diagnostic,
		FuturesRangeUniverseStructuredCompressionTradeRow: FuturesRangeUniverseStructuredCompressionTradeRow{
			SignalID:    strings.ToLower(symbol) + "_signal",
			Symbol:      symbol,
			CandidateID: StructuredCompressionCandidate4HAllH12,
			Timeframe:   RangeDiscoveryTimeframe4h,
			Family:      RangeUniverseFamilyStructuredCompressionBreak,
			CloseSplit:  fullSplitName,
			Side:        Long,
			ExitTime:    "2024-01-01T08:00:00Z",
			GrossPnL:    10,
			NetPnL:      10,
			GrossR:      1,
			NetR:        1,
			ExitReason:  "take_profit",
			HoldBars:    sequence,
		},
	}
}

func replayStopStateRows(cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig, fullTrades, stressTrades, oosTrades, recentTrades int, fullNet, pf, maxDD float64) []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow {
	rows := []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow{}
	for _, symbol := range []string{StructuredCompressionSummaryAggregateSymbol, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT, RangeUniverseSymbolBTCUSDT} {
		for _, split := range DefaultSplits() {
			count := fullTrades
			net := fullNet
			rowPF := pf
			if split.Name == "2021_2022_stress" {
				count = stressTrades
				net = 100
			} else if split.Name == "2023_2024_oos" {
				count = oosTrades
				net = 100
			} else if split.Name == "2025_2026_recent" {
				count = recentTrades
				net = 100
			}
			if symbol == RangeUniverseSymbolBTCUSDT {
				net = -10
				rowPF = 0.8
			}
			rows = append(rows, FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow{
				ConfigID:         cfg.ConfigID,
				SymbolSet:        cfg.SymbolSet,
				IsAuthority:      symbol != RangeUniverseSymbolBTCUSDT,
				IsDiagnostic:     symbol == RangeUniverseSymbolBTCUSDT,
				PassesReplayGate: true,
				FuturesRangeUniverseStructuredCompressionSummaryRow: FuturesRangeUniverseStructuredCompressionSummaryRow{
					CandidateID:    StructuredCompressionCandidate4HAllH12,
					Symbol:         symbol,
					Timeframe:      RangeDiscoveryTimeframe4h,
					Split:          split.Name,
					Side:           "all",
					TotalTrades:    count,
					GrossPnL:       net + 10,
					NetPnL:         net,
					ProfitFactor:   rowPF,
					MaxDrawdown:    maxDD,
					AvgNetR:        0.1,
					AvgInitialRisk: 1,
				},
			})
		}
		for _, side := range []string{string(Long), string(Short)} {
			rows = append(rows, FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow{
				ConfigID:         cfg.ConfigID,
				SymbolSet:        cfg.SymbolSet,
				IsAuthority:      symbol != RangeUniverseSymbolBTCUSDT,
				IsDiagnostic:     symbol == RangeUniverseSymbolBTCUSDT,
				PassesReplayGate: true,
				FuturesRangeUniverseStructuredCompressionSummaryRow: FuturesRangeUniverseStructuredCompressionSummaryRow{
					CandidateID:  StructuredCompressionCandidate4HAllH12,
					Symbol:       symbol,
					Timeframe:    RangeDiscoveryTimeframe4h,
					Split:        fullSplitName,
					Side:         side,
					TotalTrades:  30,
					GrossPnL:     100,
					NetPnL:       50,
					ProfitFactor: 1.5,
					MaxDrawdown:  maxDD,
				},
			})
		}
	}
	return rows
}
