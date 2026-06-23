package lab

import (
	"path/filepath"
	"testing"
)

func TestStructuredCompressionBaselineLongSignalTradeAndSummary(t *testing.T) {
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
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: "btcusdt_futures_um_5m.csv", ApprovedPath: "btcusdt_futures_um_5m.csv"}}
	cfg.Candidates = []FuturesRangeUniverseStructuredCompressionCandidateConfig{{
		CandidateID: StructuredCompressionCandidate4HAllH6,
		Timeframe:   RangeDiscoveryTimeframe4h,
		Side:        RangeDiscoverySideAll,
		MaxHoldBars: 2,
	}}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}
	strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles, RangeUniverseSymbolBTCUSDT, cfg.Candidates[0], cfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1", len(signals))
	}
	signal := signals[0]
	if signal.Symbol != RangeUniverseSymbolBTCUSDT || signal.CandidateID != StructuredCompressionCandidate4HAllH6 || signal.Side != Long || signal.BreakoutIndex != 1 || signal.ConfirmationIndex != 2 || signal.EntryIndex != 3 || signal.Stop != 110 || signal.Target != 126 || signal.SkippedReason != "" {
		t.Fatalf("bad long signal: %+v", signal)
	}

	result := RunBacktest(candles, strategy, btCfg)
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	trade := result.Trades[0]
	if trade.Reason != "take_profit" || trade.EntryPrice != 116 || trade.ExitPrice != 126 {
		t.Fatalf("bad long trade: %+v", trade)
	}
	tradeRows := strategy.TradeRows(result.Trades, DefaultSplits())
	if len(tradeRows) != 1 || tradeRows[0].SignalID != signal.SignalID || tradeRows[0].GrossR <= 0 || tradeRows[0].NetR <= 0 {
		t.Fatalf("bad joined trade rows: %+v", tradeRows)
	}
	summary := SummarizeFuturesRangeUniverseStructuredCompressionBaseline(signals, tradeRows, cfg, 1000, DefaultSplits())
	full := structuredCompressionSummaryByKey(summary)[structuredCompressionSummaryKey(StructuredCompressionCandidate4HAllH6, RangeUniverseSymbolBTCUSDT, fullSplitName, "all")]
	aggregate := structuredCompressionSummaryByKey(summary)[structuredCompressionSummaryKey(StructuredCompressionCandidate4HAllH6, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	if full.TotalTrades != 1 || full.SignalCount != 1 || full.NetPnL <= 0 || aggregate.TotalTrades != 1 {
		t.Fatalf("bad summary rows: full=%+v aggregate=%+v", full, aggregate)
	}
}

func TestStructuredCompressionBaselineShortStopFirst(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 103, 104, 97, 98),
		testCandle(2, 97, 98, 93, 94),
		testCandle(3, 94, 101, 83, 90),
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
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolSOLUSDT, Path: "solusdt_futures_um_5m.csv", ApprovedPath: "solusdt_futures_um_5m.csv"}}
	cfg.Candidates = []FuturesRangeUniverseStructuredCompressionCandidateConfig{{
		CandidateID: StructuredCompressionCandidate1HAllH12,
		Timeframe:   RangeDiscoveryTimeframe1h,
		Side:        RangeDiscoverySideAll,
		MaxHoldBars: 2,
	}}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}
	strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles, RangeUniverseSymbolSOLUSDT, cfg.Candidates[0], cfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 || signals[0].Side != Short || signals[0].Stop != 100 || signals[0].Target != 84 {
		t.Fatalf("bad short signal: %+v", signals)
	}
	result := RunBacktest(candles, strategy, btCfg)
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	if result.Trades[0].Reason != "stop_loss" || result.Trades[0].ExitPrice != 100 {
		t.Fatalf("expected stop-first short exit, got %+v", result.Trades[0])
	}
}

func TestStructuredCompressionBaselineSkippedSignals(t *testing.T) {
	episode := rangeRegimeDurabilityEpisode{
		EpisodeID:        1,
		StartIndex:       0,
		EndIndex:         0,
		High:             110,
		Low:              100,
		EndClose:         105,
		RawLengthBars:    1,
		ActiveLengthBars: 1,
		WidthPct:         movePct(10, 105),
	}
	candidate := FuturesRangeUniverseStructuredCompressionCandidateConfig{CandidateID: StructuredCompressionCandidate4HAllH6, Timeframe: RangeDiscoveryTimeframe4h, Side: RangeDiscoverySideAll, MaxHoldBars: 6}
	cfg := DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}

	missingEntry := newFuturesRangeUniverseStructuredCompressionSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 116, 111, 115),
	}, RangeUniverseSymbolBTCUSDT, candidate, episode, 1, 2, 1, cfg, btCfg, DefaultSplits())
	if missingEntry.SkippedReason != "missing_entry_candle" {
		t.Fatalf("missing-entry skip=%q row=%+v", missingEntry.SkippedReason, missingEntry)
	}

	invalidEntry := newFuturesRangeUniverseStructuredCompressionSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 116, 111, 115),
		testCandle(3, 109, 112, 108, 109),
	}, RangeUniverseSymbolBTCUSDT, candidate, episode, 1, 2, 2, cfg, btCfg, DefaultSplits())
	if invalidEntry.SkippedReason != "invalid_entry_geometry" {
		t.Fatalf("invalid-geometry skip=%q row=%+v", invalidEntry.SkippedReason, invalidEntry)
	}

	flatEpisode := episode
	flatEpisode.Low = flatEpisode.High
	flat := newFuturesRangeUniverseStructuredCompressionSignalRow([]Candle{
		testCandle(0, 104, 110, 110, 110),
		testCandle(1, 111, 113, 111, 112),
		testCandle(2, 112, 116, 111, 115),
		testCandle(3, 116, 117, 115, 116),
	}, RangeUniverseSymbolBTCUSDT, candidate, flatEpisode, 1, 2, 3, cfg, btCfg, DefaultSplits())
	if flat.SkippedReason != "non_positive_range_width" {
		t.Fatalf("flat-range skip=%q row=%+v", flat.SkippedReason, flat)
	}
}

func TestStructuredCompressionBaselineRunnerValidatesSourcesAndResamples(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 48)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_structured.csv", parent)
	cfg := DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: path, ApprovedPath: path, SkipSplitEligibilityCheck: true}}
	cfg.Candidates = []FuturesRangeUniverseStructuredCompressionCandidateConfig{{
		CandidateID: StructuredCompressionCandidate4HAllH6,
		Timeframe:   RangeDiscoveryTimeframe4h,
		Side:        RangeDiscoverySideAll,
		MaxHoldBars: 6,
	}}
	result, err := RunFuturesRangeUniverseStructuredCompressionBaselineBacktest(cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 1 || result.SourceRows[0].ValidationStatus != "accepted" {
		t.Fatalf("bad source rows: %+v", result.SourceRows)
	}
	if len(result.CoverageRows) != 1 || result.CoverageRows[0].Symbol != RangeUniverseSymbolBTCUSDT || result.CoverageRows[0].Timeframe != RangeDiscoveryTimeframe4h || result.CoverageRows[0].RowCount != 1 {
		t.Fatalf("bad coverage rows: %+v", result.CoverageRows)
	}
	if len(result.SummaryRows) != 24 {
		t.Fatalf("summary rows=%d, want 24", len(result.SummaryRows))
	}

	badPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_spot_5m_structured.csv", parent)
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: badPath, ApprovedPath: badPath, SkipSplitEligibilityCheck: true}}
	if _, err := RunFuturesRangeUniverseStructuredCompressionBaselineBacktest(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected spot-looking source rejection")
	}
	cfg.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: path, ApprovedPath: filepath.Join(dir, "other_btcusdt_futures_um_5m.csv"), SkipSplitEligibilityCheck: true}}
	if _, err := RunFuturesRangeUniverseStructuredCompressionBaselineBacktest(cfg, BacktestConfig{}, DefaultSplits()); err == nil {
		t.Fatalf("expected approved-path rejection")
	}
}

func TestStructuredCompressionBaselineStopStates(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
	passing := structuredCompressionStopStateRows(StructuredCompressionCandidate4HAllH6, RangeDiscoveryTimeframe4h, 100, 160, 80, 1.5, true)
	if got := FuturesRangeUniverseStructuredCompressionBaselineStopState(passing, cfg, 1000, DefaultSplits()); got != StructuredCompressionStopStatePassedNeedsOptimize {
		t.Fatalf("stop state=%s, want pass", got)
	}

	mixed := append(
		structuredCompressionStopStateRows(StructuredCompressionCandidate4HAllH6, RangeDiscoveryTimeframe4h, 100, 60, -1, 1.05, false),
		structuredCompressionStopStateRows(StructuredCompressionCandidate1HAllH12, RangeDiscoveryTimeframe1h, 100, 70, -2, 1.1, false)...,
	)
	if got := FuturesRangeUniverseStructuredCompressionBaselineStopState(mixed, cfg, 1000, DefaultSplits()); got != StructuredCompressionStopStateMixedPortfolioReview {
		t.Fatalf("stop state=%s, want mixed", got)
	}

	failed := structuredCompressionStopStateRows(StructuredCompressionCandidate4HAllH6, RangeDiscoveryTimeframe4h, 20, -50, -75, 0.5, false)
	if got := FuturesRangeUniverseStructuredCompressionBaselineStopState(failed, cfg, 1000, DefaultSplits()); got != StructuredCompressionStopStateFailedNoPromotion {
		t.Fatalf("stop state=%s, want failed", got)
	}
}

func structuredCompressionStopStateRows(candidateID, timeframe string, trades int, gross, net, pf float64, includeSymbolPasses bool) []FuturesRangeUniverseStructuredCompressionSummaryRow {
	rows := []FuturesRangeUniverseStructuredCompressionSummaryRow{}
	symbols := []string{StructuredCompressionSummaryAggregateSymbol}
	if includeSymbolPasses {
		symbols = append(symbols, RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT)
	}
	for _, symbol := range symbols {
		for _, split := range DefaultSplits() {
			count := trades
			rowGross := gross
			rowNet := net
			rowPF := pf
			if split.Name != fullSplitName {
				count = 30
				rowGross = 20
				if net >= 0 {
					rowNet = 5
				} else {
					rowNet = -1
				}
				rowPF = 1.1
			}
			rows = append(rows, FuturesRangeUniverseStructuredCompressionSummaryRow{
				CandidateID:  candidateID,
				Symbol:       symbol,
				Timeframe:    timeframe,
				Split:        split.Name,
				Side:         "all",
				TotalTrades:  count,
				GrossPnL:     rowGross,
				NetPnL:       rowNet,
				ProfitFactor: rowPF,
			})
		}
	}
	return rows
}
