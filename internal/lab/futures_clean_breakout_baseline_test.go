package lab

import "testing"

func TestFuturesCleanBreakoutBaselineLongSignalTradeAndSummary(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 123, 111, 121),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
	}
	cfg := DefaultFuturesCleanBreakoutBaselineConfig()
	cfg.EventDelayBars = 1
	cfg.Candidates = []FuturesCleanBreakoutBaselineCandidateConfig{{
		CandidateID: CleanBreakoutCandidate4HUp,
		Timeframe:   RangeDiscoveryTimeframe4h,
		Side:        RangeDiscoverySideUp,
		MaxHoldBars: 2,
	}}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}
	strategy, err := newFuturesCleanBreakoutBaselineStrategyFromClassifications(candles, cfg.Candidates[0], cfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1", len(signals))
	}
	signal := signals[0]
	if signal.CandidateID != CleanBreakoutCandidate4HUp || signal.Side != Long || signal.Stop != 110 || signal.Target != 122 || signal.EntryIndex != 2 || signal.SkippedReason != "" {
		t.Fatalf("bad long signal: %+v", signal)
	}

	result := RunBacktest(candles, strategy, btCfg)
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	trade := result.Trades[0]
	if trade.Reason != "take_profit" || trade.EntryPrice != 112 || trade.ExitPrice != 122 {
		t.Fatalf("bad long trade: %+v", trade)
	}
	tradeRows := strategy.TradeRows(result.Trades, DefaultSplits())
	if len(tradeRows) != 1 || tradeRows[0].SignalID != signal.SignalID || tradeRows[0].GrossR <= 0 || tradeRows[0].NetR <= 0 {
		t.Fatalf("bad joined trade rows: %+v", tradeRows)
	}
	summary := SummarizeFuturesCleanBreakoutBaseline(signals, tradeRows, cfg, 1000, DefaultSplits())
	full := cleanBreakoutSummaryByKey(summary)[CleanBreakoutCandidate4HUp+"|"+fullSplitName+"|all"]
	if full.TotalTrades != 1 || full.SignalCount != 1 || full.NetPnL <= 0 {
		t.Fatalf("bad summary row: %+v", full)
	}
}

func TestFuturesCleanBreakoutBaselineShortStopFirst(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 103, 104, 97, 98),
		testCandle(2, 98, 101, 87, 90),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
	}
	cfg := DefaultFuturesCleanBreakoutBaselineConfig()
	cfg.EventDelayBars = 1
	cfg.Candidates = []FuturesCleanBreakoutBaselineCandidateConfig{{
		CandidateID: CleanBreakoutCandidate1HAll,
		Timeframe:   RangeDiscoveryTimeframe1h,
		Side:        RangeDiscoverySideAll,
		MaxHoldBars: 2,
	}}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}
	strategy, err := newFuturesCleanBreakoutBaselineStrategyFromClassifications(candles, cfg.Candidates[0], cfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 || signals[0].Side != Short || signals[0].Stop != 100 || signals[0].Target != 88 {
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

func TestFuturesCleanBreakoutBaselineSkippedSignals(t *testing.T) {
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
	candidate := FuturesCleanBreakoutBaselineCandidateConfig{CandidateID: CleanBreakoutCandidate4HUp, Timeframe: RangeDiscoveryTimeframe4h, Side: RangeDiscoverySideUp, MaxHoldBars: 12}
	cfg := DefaultFuturesCleanBreakoutBaselineConfig()
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}

	missingEntry := newFuturesCleanBreakoutBaselineSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
	}, candidate, episode, 1, 1, cfg, btCfg, DefaultSplits())
	if missingEntry.SkippedReason != "missing_entry_candle" {
		t.Fatalf("missing-entry skip=%q row=%+v", missingEntry.SkippedReason, missingEntry)
	}

	invalidEntry := newFuturesCleanBreakoutBaselineSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 109, 112, 108, 109),
	}, candidate, episode, 1, 2, cfg, btCfg, DefaultSplits())
	if invalidEntry.SkippedReason != "invalid_entry_geometry" {
		t.Fatalf("invalid-geometry skip=%q row=%+v", invalidEntry.SkippedReason, invalidEntry)
	}

	flatEpisode := episode
	flatEpisode.Low = flatEpisode.High
	flat := newFuturesCleanBreakoutBaselineSignalRow([]Candle{
		testCandle(0, 104, 110, 110, 110),
		testCandle(1, 111, 113, 111, 112),
		testCandle(2, 112, 113, 111, 112),
	}, candidate, flatEpisode, 1, 3, cfg, btCfg, DefaultSplits())
	if flat.SkippedReason != "non_positive_range_width" {
		t.Fatalf("flat-range skip=%q row=%+v", flat.SkippedReason, flat)
	}
}

func TestFuturesCleanBreakoutBaselineRunnerResamplesAndFiltersCandidates(t *testing.T) {
	parent := make([]Candle, 48)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	cfg := DefaultFuturesCleanBreakoutBaselineConfig()
	cfg.Candidates = []FuturesCleanBreakoutBaselineCandidateConfig{{
		CandidateID: CleanBreakoutCandidate4HUp,
		Timeframe:   RangeDiscoveryTimeframe4h,
		Side:        RangeDiscoverySideUp,
		MaxHoldBars: 12,
	}}
	result, err := RunFuturesCleanBreakoutBaselineBacktest(parent, cfg, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.CoverageRows) != 1 || result.CoverageRows[0].Timeframe != RangeDiscoveryTimeframe4h || result.CoverageRows[0].RowCount != 1 {
		t.Fatalf("bad coverage rows: %+v", result.CoverageRows)
	}
	if len(result.SummaryRows) != len(DefaultSplits())*3 {
		t.Fatalf("summary rows=%d, want %d", len(result.SummaryRows), len(DefaultSplits())*3)
	}
	if _, err := RunFuturesCleanBreakoutBaselineBacktest(parent, FuturesCleanBreakoutBaselineConfig{EventDelayBars: -1}, BacktestConfig{}, nil); err == nil {
		t.Fatalf("expected invalid config error")
	}
}

func TestFuturesCleanBreakoutBaselineStopStates(t *testing.T) {
	cfg := DefaultFuturesCleanBreakoutBaselineConfig()
	passing := cleanBreakoutStopStateRows(CleanBreakoutCandidate4HUp, RangeDiscoveryTimeframe4h, 100, 150, 75, 1.5)
	if got := FuturesCleanBreakoutBaselineStopState(passing, cfg, 1000, DefaultSplits()); got != CleanBreakoutStopStatePassedNeedsOptimization {
		t.Fatalf("stop state=%s, want pass", got)
	}

	mixed := append(
		cleanBreakoutStopStateRows(CleanBreakoutCandidate4HUp, RangeDiscoveryTimeframe4h, 100, 60, -1, 1.05),
		cleanBreakoutStopStateRows(CleanBreakoutCandidate1HAll, RangeDiscoveryTimeframe1h, 100, 70, -2, 1.1)...,
	)
	if got := FuturesCleanBreakoutBaselineStopState(mixed, cfg, 1000, DefaultSplits()); got != CleanBreakoutStopStateMixedNeedsPortfolioReview {
		t.Fatalf("stop state=%s, want mixed", got)
	}

	failed := cleanBreakoutStopStateRows(CleanBreakoutCandidate4HUp, RangeDiscoveryTimeframe4h, 20, -50, -75, 0.5)
	if got := FuturesCleanBreakoutBaselineStopState(failed, cfg, 1000, DefaultSplits()); got != CleanBreakoutStopStateFailedNoPromotion {
		t.Fatalf("stop state=%s, want failed", got)
	}
}

func cleanBreakoutStopStateRows(candidateID, timeframe string, trades int, gross, net, pf float64) []FuturesCleanBreakoutBaselineSummaryRow {
	rows := []FuturesCleanBreakoutBaselineSummaryRow{}
	for _, split := range DefaultSplits() {
		count := trades
		if split.Name != fullSplitName {
			count = 30
		}
		row := FuturesCleanBreakoutBaselineSummaryRow{
			CandidateID:  candidateID,
			Timeframe:    timeframe,
			Split:        split.Name,
			Side:         "all",
			TotalTrades:  count,
			GrossPnL:     gross,
			NetPnL:       net,
			ProfitFactor: pf,
		}
		if split.Name != fullSplitName {
			row.GrossPnL = 20
			if net >= 0 {
				row.NetPnL = 5
			} else {
				row.NetPnL = -1
			}
		}
		rows = append(rows, row)
	}
	return rows
}
