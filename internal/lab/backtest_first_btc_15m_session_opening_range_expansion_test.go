package lab

import (
	"math"
	"testing"
	"time"
)

func TestBTC15MSessionOpeningRangeBuildsOneLongSignalPerUTCDate(t *testing.T) {
	candles := sessionOpeningRangeTest15MCandles()
	cfg := DefaultBacktestFirstBTC15MSessionOpeningRangeExpansionConfig()
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, FeePct: 0.0004, SlippagePct: 0.000116, MaxHoldBars: cfg.MaxHoldBars}

	sessionRows, signals, skips := btc15MSessionOpeningRangeBuildSignals(candles, cfg, btCfg, []Split{{Name: fullSplitName}})
	if len(sessionRows) != 1 || !sessionRows[0].RangeReady || !sessionRows[0].RangeWidthPositive {
		t.Fatalf("unexpected session rows: %+v", sessionRows)
	}
	if len(signals) != 1 {
		t.Fatalf("expected exactly one signal for the UTC date, got %d skips=%+v", len(signals), btc15MSessionOpeningRangeSkipRows(skips))
	}
	signal := signals[0]
	if signal.Side != Long {
		t.Fatalf("side=%s, want long", signal.Side)
	}
	if got := signal.DecisionOpenTime; got != "2021-01-01T14:30:00Z" {
		t.Fatalf("decision open=%s", got)
	}
	if signal.EntryIndex != signal.DecisionIndex+1 {
		t.Fatalf("entry index=%d decision=%d", signal.EntryIndex, signal.DecisionIndex)
	}
	if !signal.PrePositionCandidate || signal.ForwardLabelsAsInput || signal.UsesFutureRowsForFeatures || signal.DerivativesVetoUsed || signal.OptimizerSelectionUsed {
		t.Fatalf("unexpected signal flags: %+v", signal)
	}
	wantEntry := applySlippage(candles[signal.EntryIndex].Open, btCfg.SlippagePct, Long, true)
	if !sessionOpeningRangeAlmostEqual(signal.ExpectedEntryPrice, wantEntry, 1e-9) {
		t.Fatalf("entry=%f want=%f", signal.ExpectedEntryPrice, wantEntry)
	}
	wantStop := signal.OpeningRangeLow - cfg.AcceptanceATRMultiple*signal.ATR14
	if !sessionOpeningRangeAlmostEqual(signal.Stop, wantStop, 1e-9) {
		t.Fatalf("stop=%f want=%f", signal.Stop, wantStop)
	}
	wantTarget := signal.ExpectedEntryPrice + cfg.TargetR*(signal.ExpectedEntryPrice-signal.Stop)
	if !sessionOpeningRangeAlmostEqual(signal.Target, wantTarget, 1e-9) {
		t.Fatalf("target=%f want=%f", signal.Target, wantTarget)
	}

	signals, trades := btc15MSessionOpeningRangeRun(candles, signals, &skips, cfg, btCfg)
	if len(trades) != 1 {
		t.Fatalf("trades=%d want 1", len(trades))
	}
	if !signals[0].Executed {
		t.Fatalf("signal was not marked executed: %+v", signals[0])
	}
	if trades[0].Reason != "take_profit" {
		t.Fatalf("exit reason=%s want take_profit", trades[0].Reason)
	}
}

func TestRunBacktestFirstBTC15MSessionOpeningRangeExpansionResamplesAndReports(t *testing.T) {
	base15m := sessionOpeningRangeTest15MCandles()
	parent5m := sessionOpeningRangeExpand15MTo5M(base15m)
	manifest := SourceManifest{
		Path:             "btcusdt_futures_um_5m_test.csv",
		Venue:            "Binance",
		Product:          "Binance USDT-M futures",
		Symbol:           "BTCUSDT",
		Interval:         "5m",
		RowCount:         len(parent5m),
		FirstOpenTime:    parent5m[0].OpenTime.UTC().Format(timeLayout),
		LastOpenTime:     parent5m[len(parent5m)-1].OpenTime.UTC().Format(timeLayout),
		Monotonic:        true,
		ValidationStatus: "accepted",
	}
	cfg := DefaultBacktestFirstBTC15MSessionOpeningRangeExpansionConfig()
	cfg.ApprovedSourcePath = manifest.Path
	cfg.SkipSourceFactCheck = true
	cfg.SkipCoverageCountCheck = true
	cfg.MinFullTrades = 1
	cfg.MinSplitTrades = 1
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, FeePct: 0.0004, SlippagePct: 0.000116, MaxHoldBars: cfg.MaxHoldBars}

	result, err := RunBacktestFirstBTC15MSessionOpeningRangeExpansion(parent5m, manifest, cfg, btCfg, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.CoverageRows) != 1 || !result.CoverageRows[0].SourceResamplePass {
		t.Fatalf("coverage did not pass: %+v", result.CoverageRows)
	}
	if len(result.SessionRangeRows) != 1 || len(result.SignalRows) != 1 || len(result.Trades) != 1 || len(result.TradeRows) != 1 {
		t.Fatalf("unexpected result counts sessions=%d signals=%d trades=%d trade_rows=%d", len(result.SessionRangeRows), len(result.SignalRows), len(result.Trades), len(result.TradeRows))
	}
	if result.Falsification.SourceResamplePass != true || result.Falsification.LeakagePass != true || result.Falsification.OptimizerContaminationPass != true || result.Falsification.DerivativesVetoContaminationPass != true {
		t.Fatalf("unexpected falsification contamination flags: %+v", result.Falsification)
	}
	if result.StopState == "" {
		t.Fatalf("missing stop state")
	}
}

func sessionOpeningRangeTest15MCandles() []Candle {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]Candle, 80)
	for i := range out {
		openTime := start.Add(time.Duration(i) * 15 * time.Minute)
		out[i] = Candle{
			OpenTime:  openTime,
			CloseTime: openTime.Add(15*time.Minute - time.Millisecond),
			Open:      100,
			High:      100.8,
			Low:       99.2,
			Close:     100,
			Volume:    100,
		}
	}
	for i := 54; i <= 57; i++ {
		out[i].High = 101
		out[i].Low = 99
	}
	out[58].Open = 100
	out[58].High = 102.4
	out[58].Low = 99.8
	out[58].Close = 102
	out[59].Open = 102
	out[59].High = 108
	out[59].Low = 101.5
	out[59].Close = 107
	return out
}

func sessionOpeningRangeExpand15MTo5M(candles []Candle) []Candle {
	out := make([]Candle, 0, len(candles)*3)
	for _, candle := range candles {
		for offset := 0; offset < 3; offset++ {
			openTime := candle.OpenTime.Add(time.Duration(offset) * 5 * time.Minute)
			out = append(out, Candle{
				OpenTime:  openTime,
				CloseTime: openTime.Add(5*time.Minute - time.Millisecond),
				Open:      candle.Open,
				High:      candle.High,
				Low:       candle.Low,
				Close:     candle.Close,
				Volume:    candle.Volume / 3,
			})
		}
	}
	return out
}

func sessionOpeningRangeAlmostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
