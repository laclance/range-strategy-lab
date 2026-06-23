package lab

import "testing"

func TestHoldInsideMidlineTouchPrototypeExactSurfaceAndCloseBackSide(t *testing.T) {
	strategy := newPrototypeTestStrategy(t, prototypeLongTimeStopCandles())
	rows := strategy.SignalRows()
	if len(rows) != 1 {
		t.Fatalf("signal rows=%d, want 1: %+v", len(rows), rows)
	}
	row := rows[0]
	if row.ContextRule != DetectorContextRuleHold3Inside ||
		row.EventType != HoldInsideMidlineReactionEventTouch ||
		row.EventClosePositionBucket != decisionClosePositionBucketMid50 ||
		row.Side != Long ||
		row.Stop != 98 ||
		row.Target != 102 ||
		row.MaxHoldBars != 6 ||
		row.EventIndex != 5 ||
		row.EntryIndex != 6 {
		t.Fatalf("bad prototype signal row: %+v", row)
	}
	sig, ok := strategy.OnCandle(StrategyContext{Candles: prototypeLongTimeStopCandles(), Index: 5})
	if !ok || sig.Side != Long || sig.Stop != 98 || sig.Target != 102 || sig.MaxHoldBars != 6 || sig.Reason == "" {
		t.Fatalf("bad event-candle signal: %+v ok=%v", sig, ok)
	}
	if _, ok := strategy.OnCandle(StrategyContext{Candles: prototypeLongTimeStopCandles(), Index: 4}); ok {
		t.Fatalf("strategy signaled before event candle")
	}
}

func TestHoldInsideMidlineTouchPrototypeFiltersMid50AndExactMid(t *testing.T) {
	highBucket := prototypeLongTimeStopCandles()
	highBucket[5] = testCandle(5, 99.2, 101.8, 99.1, 101.5)
	if rows := newPrototypeTestStrategy(t, highBucket).SignalRows(); len(rows) != 0 {
		t.Fatalf("high bucket should be outside exact surface, rows=%+v", rows)
	}

	exactMid := prototypeLongTimeStopCandles()
	exactMid[5] = testCandle(5, 99.2, 100.2, 99.1, 100)
	rows := newPrototypeTestStrategy(t, exactMid).SignalRows()
	if len(rows) != 1 || rows[0].SignalID != "" || rows[0].SkippedReason != "event_close_at_mid" {
		t.Fatalf("exact-mid event should be documented but not signaled: %+v", rows)
	}
}

func TestHoldInsideMidlineTouchPrototypeShortSideAndEntryGeometrySkip(t *testing.T) {
	shortCandles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 100.6),
		testCandle(3, 100.6, 101, 99.8, 100.7),
		testCandle(4, 100.7, 101, 99.8, 100.8),
		testCandle(5, 100.8, 101, 99.8, 100.6),
		testCandle(6, 100.5, 101, 99.5, 100.2),
	}
	strategy := newPrototypeTestStrategy(t, shortCandles)
	rows := strategy.SignalRows()
	if len(rows) != 1 || rows[0].Side != Short || rows[0].Stop != 102 || rows[0].Target != 98 {
		t.Fatalf("bad short close-back signal: %+v", rows)
	}

	invalidEntry := prototypeLongTimeStopCandles()
	invalidEntry[6] = testCandle(6, 103, 103.5, 102.5, 103)
	strategy = newPrototypeTestStrategy(t, invalidEntry)
	rows = strategy.SignalRows()
	if len(rows) != 1 || rows[0].EntryGeometryValid {
		t.Fatalf("expected invalid next-bar entry geometry: %+v", rows)
	}
	result := RunBacktest(invalidEntry, strategy, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1})
	if len(result.Trades) != 0 {
		t.Fatalf("invalid geometry should produce no trades: %+v", result.Trades)
	}
}

func TestHoldInsideMidlineTouchPrototypeBacktestTimeStopAndStopFirst(t *testing.T) {
	candles := prototypeLongTimeStopCandles()
	strategy := newPrototypeTestStrategy(t, candles)
	result := RunBacktest(candles, strategy, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1: %+v", len(result.Trades), result.Trades)
	}
	trade := result.Trades[0]
	if trade.OpenIndex != 6 || trade.CloseIndex != 12 || trade.Reason != "time_stop" || trade.HoldBars != 6 {
		t.Fatalf("bad time-stop trade: %+v", trade)
	}
	if trade.EntryPrice != candles[6].Open {
		t.Fatalf("entry should use next bar open, got %v want %v", trade.EntryPrice, candles[6].Open)
	}

	stopFirst := prototypeLongTimeStopCandles()
	stopFirst[6] = testCandle(6, 99.5, 103, 97, 100)
	result = RunBacktest(stopFirst, newPrototypeTestStrategy(t, stopFirst), BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1})
	if len(result.Trades) != 1 || result.Trades[0].Reason != "stop_loss" {
		t.Fatalf("same-bar stop/target ambiguity should be stop-first: %+v", result.Trades)
	}
}

func TestHoldInsideMidlineTouchPrototypeTradeRowsAndSummary(t *testing.T) {
	candles := prototypeLongTimeStopCandles()
	strategy := newPrototypeTestStrategy(t, candles)
	result := RunBacktest(candles, strategy, BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1})
	tradeRows := strategy.TradeRows(result.Trades, DefaultSplits())
	if len(tradeRows) != 1 ||
		tradeRows[0].SignalID == "" ||
		tradeRows[0].InitialRisk <= 0 ||
		tradeRows[0].NetR == 0 ||
		tradeRows[0].EntrySplit != "2021_2022_stress" ||
		tradeRows[0].CloseSplit != "2021_2022_stress" {
		t.Fatalf("bad joined trade row: %+v", tradeRows)
	}

	rows := SummarizeHoldInsideMidlineTouchPrototype([]HoldInsideMidlineTouchPrototypeTradeRow{
		{Side: Long, ExitTime: "2022-01-01T00:00:00Z", GrossPnL: 12, NetPnL: 10, Fees: 1, Slippage: 1, InitialRisk: 10, GrossR: 1.2, NetR: 1, HoldBars: 2},
		{Side: Short, ExitTime: "2024-01-01T00:00:00Z", GrossPnL: -3, NetPnL: -5, Fees: 1, Slippage: 1, InitialRisk: 10, GrossR: -0.3, NetR: -0.5, HoldBars: 4},
	}, 1000, DefaultSplits())
	byKey := map[string]HoldInsideMidlineTouchPrototypeSummaryRow{}
	for _, row := range rows {
		byKey[row.Split+"|"+row.Side] = row
	}
	full := byKey[fullSplitName+"|all"]
	if full.TotalTrades != 2 ||
		!boundaryAlmostEqual(full.AvgNetR, 0.25) ||
		!boundaryAlmostEqual(full.AvgGrossR, 0.45) ||
		full.TotalCosts != 4 {
		t.Fatalf("bad full summary row: %+v", full)
	}
	if !byKey["2023_2024_oos|all"].IsWorstPeriodSplit {
		t.Fatalf("expected 2023_2024_oos to be marked as worst split: %+v", byKey["2023_2024_oos|all"])
	}
}

func newPrototypeTestStrategy(t *testing.T, candles []Candle) HoldInsideMidlineTouchPrototypeStrategy {
	t.Helper()
	raw := make([]bool, len(candles))
	active := make([]bool, len(candles))
	raw[0], raw[1] = true, true
	active[1] = true
	classifications := testCompressionClassifications(raw, active)
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	strategy, err := newHoldInsideMidlineTouchPrototypeStrategyFromClassifications(candles, classifications, nil, HoldInsideMidlineTouchPrototypeConfig{
		Profile: profile,
	}, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	return strategy
}

func prototypeLongTimeStopCandles() []Candle {
	return []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 100),
		testCandle(2, 100, 101, 99, 99),
		testCandle(3, 99, 100.5, 98.8, 99.5),
		testCandle(4, 99.5, 100.2, 98.9, 99.2),
		testCandle(5, 99.2, 100.1, 99.1, 99.4),
		testCandle(6, 99.5, 101, 99.2, 99.8),
		testCandle(7, 99.8, 101, 99.1, 100),
		testCandle(8, 100, 101, 99.1, 100.1),
		testCandle(9, 100.1, 101, 99.1, 100.2),
		testCandle(10, 100.2, 101, 99.1, 100.1),
		testCandle(11, 100.1, 101, 99.1, 100),
		testCandle(12, 100, 101, 99.1, 100.2),
	}
}
