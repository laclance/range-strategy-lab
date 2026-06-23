package lab

import (
	"path/filepath"
	"testing"
)

func TestBreakoutRetestAcceptanceSelectionPicksDefaultAllSideCandidates(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
	rows := []FuturesRangeUniverseRankingRow{
		{Rank: 20, Timeframe: RangeDiscoveryTimeframe15m, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideDown, HorizonBars: 12, PassesGate: true, RankScore: 99},
		{Rank: 22, Timeframe: RangeDiscoveryTimeframe15m, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideAll, HorizonBars: 12, PassesGate: true, RankScore: 90},
		{Rank: 23, Timeframe: RangeDiscoveryTimeframe15m, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideAll, HorizonBars: 6, PassesGate: true, RankScore: 89},
		{Rank: 24, Timeframe: RangeDiscoveryTimeframe1h, Family: RangeUniverseFamilyStructuredCompressionBreak, Side: RangeDiscoverySideAll, HorizonBars: 12, PassesGate: true, RankScore: 88},
		{Rank: 28, Timeframe: RangeDiscoveryTimeframe1h, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideAll, HorizonBars: 12, PassesGate: true, RankScore: 80},
		{Rank: 29, Timeframe: RangeDiscoveryTimeframe1h, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideAll, HorizonBars: 6, PassesGate: true, RankScore: 79},
	}

	selected := SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidates(rows, cfg)
	if len(selected) != 2 {
		t.Fatalf("selected=%d, want 2: %+v", len(selected), selected)
	}
	if selected[0].CandidateID != BreakoutRetestAcceptanceCandidate15MAllH12 || selected[0].SelectedOrder != 1 || selected[0].SourceRank != 22 {
		t.Fatalf("bad primary selection: %+v", selected[0])
	}
	if selected[1].CandidateID != BreakoutRetestAcceptanceCandidate1HAllH12 || selected[1].SelectedOrder != 2 || selected[1].SourceRank != 28 {
		t.Fatalf("bad secondary selection: %+v", selected[1])
	}
}

func TestBreakoutRetestAcceptanceNoCandidateStopState(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
	selected := SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidates([]FuturesRangeUniverseRankingRow{
		{Rank: 1, Timeframe: RangeDiscoveryTimeframe15m, Family: RangeUniverseFamilyBreakoutRetestAcceptance, Side: RangeDiscoverySideAll, HorizonBars: 12, PassesGate: false},
		{Rank: 2, Timeframe: RangeDiscoveryTimeframe1h, Family: RangeUniverseFamilyStructuredCompressionBreak, Side: RangeDiscoverySideAll, HorizonBars: 12, PassesGate: true},
	}, cfg)
	if len(selected) != 0 {
		t.Fatalf("selected=%d, want 0: %+v", len(selected), selected)
	}
	result := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult{}
	if got := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineStopState(result, cfg, 1000, DefaultSplits()); got != BreakoutRetestAcceptanceStopStateNoRankedCandidate {
		t.Fatalf("stop state=%s, want no candidate", got)
	}
}

func TestBreakoutRetestAcceptanceLongSignalTradeAndSummary(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 114, 109, 112),
		testCandle(3, 116, 127, 115, 126),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}
	cfg := breakoutRetestAcceptanceTestConfig()
	cfg.DiscoveryConfig.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolBTCUSDT, Path: "btcusdt_futures_um_5m.csv", ApprovedPath: "btcusdt_futures_um_5m.csv"}}
	candidate := FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{
		CandidateID:                  BreakoutRetestAcceptanceCandidate15MAllH12,
		SelectedOrder:                1,
		Timeframe:                    RangeDiscoveryTimeframe15m,
		Side:                         RangeDiscoverySideAll,
		MaxHoldBars:                  2,
		TargetRangeWidthMultiple:     1,
		StopBoundaryBufferRangeWidth: 0,
		SourceRank:                   22,
	}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}

	strategy, err := newFuturesRangeUniverseBreakoutRetestAcceptanceStrategyFromClassifications(candles, RangeUniverseSymbolBTCUSDT, candidate, cfg, btCfg, classifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1", len(signals))
	}
	signal := signals[0]
	if signal.Symbol != RangeUniverseSymbolBTCUSDT || signal.CandidateID != BreakoutRetestAcceptanceCandidate15MAllH12 || signal.SelectedOrder != 1 || signal.SourceRank != 22 || signal.Side != Long || signal.BreakoutIndex != 1 || signal.RetestIndex != 2 || signal.EntryIndex != 3 || signal.Stop != 110 || signal.Target != 126 || signal.SkippedReason != "" {
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
	summary := SummarizeFuturesRangeUniverseBreakoutRetestAcceptanceBaseline(signals, tradeRows, cfg, 1000, DefaultSplits())
	byKey := breakoutRetestAcceptanceSummaryByKey(summary)
	full := byKey[breakoutRetestAcceptanceSummaryKey(BreakoutRetestAcceptanceCandidate15MAllH12, RangeUniverseSymbolBTCUSDT, fullSplitName, "all")]
	aggregate := byKey[breakoutRetestAcceptanceSummaryKey(BreakoutRetestAcceptanceCandidate15MAllH12, BreakoutRetestAcceptanceSummaryAggregateSymbol, fullSplitName, "all")]
	if full.TotalTrades != 1 || full.SignalCount != 1 || full.NetPnL <= 0 || aggregate.TotalTrades != 1 {
		t.Fatalf("bad summary rows: full=%+v aggregate=%+v", full, aggregate)
	}
}

func TestBreakoutRetestAcceptanceShortStopFirst(t *testing.T) {
	candles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 103, 104, 97, 98),
		testCandle(2, 97, 101, 96, 98),
		testCandle(3, 94, 101, 83, 90),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}
	cfg := breakoutRetestAcceptanceTestConfig()
	cfg.DiscoveryConfig.Discovery.MaxEventDelayBars = 1
	cfg.DiscoveryConfig.Sources = []FuturesRangeUniverseSourceConfig{{Symbol: RangeUniverseSymbolSOLUSDT, Path: "solusdt_futures_um_5m.csv", ApprovedPath: "solusdt_futures_um_5m.csv"}}
	candidate := FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{CandidateID: BreakoutRetestAcceptanceCandidate1HAllH12, SelectedOrder: 2, Timeframe: RangeDiscoveryTimeframe1h, Side: RangeDiscoverySideAll, MaxHoldBars: 2}
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}

	strategy, err := newFuturesRangeUniverseBreakoutRetestAcceptanceStrategyFromClassifications(candles, RangeUniverseSymbolSOLUSDT, candidate, cfg, btCfg, classifications, DefaultSplits())
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

func TestBreakoutRetestAcceptanceSkippedAndDuplicateSignals(t *testing.T) {
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
	candidate := FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{CandidateID: BreakoutRetestAcceptanceCandidate15MAllH12, SelectedOrder: 1, Timeframe: RangeDiscoveryTimeframe15m, Side: RangeDiscoverySideAll, MaxHoldBars: 12}
	cfg := breakoutRetestAcceptanceTestConfig()
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1}
	event := rangeDiscoveryEvent{eventIndex: 2, episode: episode, side: RangeDiscoverySideUp, eventDelayBars: 2, breakoutDelayBars: 1, reentryDelayBars: 1}

	missingEntry := newFuturesRangeUniverseBreakoutRetestAcceptanceSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 114, 109, 112),
	}, RangeUniverseSymbolBTCUSDT, candidate, event, 1, cfg, btCfg, DefaultSplits())
	if missingEntry.SkippedReason != "missing_entry_candle" {
		t.Fatalf("missing-entry skip=%q row=%+v", missingEntry.SkippedReason, missingEntry)
	}

	invalidEntry := newFuturesRangeUniverseBreakoutRetestAcceptanceSignalRow([]Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 114, 109, 112),
		testCandle(3, 109, 112, 108, 109),
	}, RangeUniverseSymbolBTCUSDT, candidate, event, 2, cfg, btCfg, DefaultSplits())
	if invalidEntry.SkippedReason != "invalid_entry_geometry" {
		t.Fatalf("invalid-geometry skip=%q row=%+v", invalidEntry.SkippedReason, invalidEntry)
	}

	flatEpisode := episode
	flatEpisode.Low = flatEpisode.High
	flatEvent := event
	flatEvent.episode = flatEpisode
	flat := newFuturesRangeUniverseBreakoutRetestAcceptanceSignalRow([]Candle{
		testCandle(0, 104, 110, 110, 110),
		testCandle(1, 111, 113, 111, 112),
		testCandle(2, 112, 114, 109, 112),
		testCandle(3, 116, 117, 115, 116),
	}, RangeUniverseSymbolBTCUSDT, candidate, flatEvent, 3, cfg, btCfg, DefaultSplits())
	if flat.SkippedReason != "non_positive_range_width" {
		t.Fatalf("flat-range skip=%q row=%+v", flat.SkippedReason, flat)
	}

	duplicateCandles := []Candle{
		testCandle(0, 104, 110, 100, 105),
		testCandle(1, 107, 113, 106, 112),
		testCandle(2, 112, 114, 111, 113),
		testCandle(3, 112, 114, 109, 112),
		testCandle(4, 116, 127, 115, 126),
	}
	duplicateClassifications := []RangeClassification{
		{Index: 0, RawActive: true, Active: true},
		{Index: 1},
		{Index: 2},
		{Index: 3},
		{Index: 4},
	}
	strategy, err := newFuturesRangeUniverseBreakoutRetestAcceptanceStrategyFromClassifications(duplicateCandles, RangeUniverseSymbolBTCUSDT, candidate, cfg, btCfg, duplicateClassifications, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	signals := strategy.SignalRows()
	if len(signals) != 2 {
		t.Fatalf("signals=%d, want duplicate pair: %+v", len(signals), signals)
	}
	if signals[0].SkippedReason != "" || signals[1].SkippedReason != "duplicate_retest_index" {
		t.Fatalf("bad duplicate handling: %+v", signals)
	}
}

func TestBreakoutRetestAcceptanceSourceValidationAndSelectedResamples(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 48)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_breakout_retest.csv", parent)
	source := FuturesRangeUniverseSourceConfig{Symbol: RangeUniverseSymbolBTCUSDT, Path: path, ApprovedPath: path, SkipSplitEligibilityCheck: true}
	candles, sourceRow, err := LoadFuturesRangeUniverseSource(source, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if sourceRow.ValidationStatus != "accepted" || sourceRow.RowCount != 48 || sourceRow.ComparisonOnly {
		t.Fatalf("bad source row: %+v", sourceRow)
	}

	frame15m, ok := breakoutRetestAcceptanceFrameDef(RangeDiscoveryTimeframe15m)
	if !ok {
		t.Fatalf("missing 15m frame")
	}
	resampled15m, coverage15m, err := resampleRangeDiscoveryFrame(candles, frame15m)
	if err != nil {
		t.Fatal(err)
	}
	if len(resampled15m) != 16 || coverage15m.RowCount != 16 || !coverage15m.Complete || coverage15m.ValidationStatus != "accepted" {
		t.Fatalf("bad 15m coverage: rows=%d coverage=%+v", len(resampled15m), coverage15m)
	}

	frame1h, ok := breakoutRetestAcceptanceFrameDef(RangeDiscoveryTimeframe1h)
	if !ok {
		t.Fatalf("missing 1h frame")
	}
	resampled1h, coverage1h, err := resampleRangeDiscoveryFrame(candles, frame1h)
	if err != nil {
		t.Fatal(err)
	}
	if len(resampled1h) != 4 || coverage1h.RowCount != 4 || !coverage1h.Complete || coverage1h.ValidationStatus != "accepted" {
		t.Fatalf("bad 1h coverage: rows=%d coverage=%+v", len(resampled1h), coverage1h)
	}

	badPath := writeRangeUniverseTestCSV(t, dir, "btcusdt_spot_5m_breakout_retest.csv", parent)
	source.Path = badPath
	source.ApprovedPath = badPath
	if _, _, err := LoadFuturesRangeUniverseSource(source, DefaultSplits()); err == nil {
		t.Fatalf("expected spot-looking source rejection")
	}
	source.Path = path
	source.ApprovedPath = filepath.Join(dir, "other_btcusdt_futures_um_5m.csv")
	if _, _, err := LoadFuturesRangeUniverseSource(source, DefaultSplits()); err == nil {
		t.Fatalf("expected approved-path rejection")
	}
}

func TestBreakoutRetestAcceptanceStopStates(t *testing.T) {
	cfg := DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
	result := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult{
		SelectionRows: []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow{{
			SelectedOrder: 1,
			CandidateID:   BreakoutRetestAcceptanceCandidate15MAllH12,
			Timeframe:     RangeDiscoveryTimeframe15m,
			Family:        RangeUniverseFamilyBreakoutRetestAcceptance,
			Side:          RangeDiscoverySideAll,
			HorizonBars:   12,
			PassesGate:    true,
		}},
		SummaryRows: breakoutRetestAcceptanceStopStateRows(BreakoutRetestAcceptanceCandidate15MAllH12, RangeDiscoveryTimeframe15m, 120, 180, 120, 1.6, true),
	}
	if got := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineStopState(result, cfg, 1000, DefaultSplits()); got != BreakoutRetestAcceptanceStopStatePassedNeedsRobust {
		t.Fatalf("stop state=%s, want pass", got)
	}

	result.SummaryRows = breakoutRetestAcceptanceStopStateRows(BreakoutRetestAcceptanceCandidate15MAllH12, RangeDiscoveryTimeframe15m, 80, 20, -10, 0.8, false)
	if got := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineStopState(result, cfg, 1000, DefaultSplits()); got != BreakoutRetestAcceptanceStopStateFailedNoPromotion {
		t.Fatalf("stop state=%s, want failed", got)
	}
}

func breakoutRetestAcceptanceTestConfig() FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig {
	cfg := DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
	cfg.DiscoveryConfig.Discovery.MaxEventDelayBars = 2
	cfg.DiscoveryConfig.Discovery.ReentryWindowBars = 2
	cfg.DiscoveryConfig.Discovery.DetectorLookbackBarsOverride = 2
	cfg.DiscoveryConfig.Discovery.DetectorMinConsecutiveBars = 1
	return cfg
}

func breakoutRetestAcceptanceStopStateRows(candidateID, timeframe string, trades int, gross, net, pf float64, includeSymbolPasses bool) []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow {
	rows := []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow{}
	symbols := []string{BreakoutRetestAcceptanceSummaryAggregateSymbol}
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
			rows = append(rows, FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow{
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
