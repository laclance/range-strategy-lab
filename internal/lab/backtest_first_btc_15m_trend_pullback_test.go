package lab

import (
	"math"
	"testing"
	"time"
)

func TestBTC15MTrendPullbackBuildsLongSignalFromFixedPacketRules(t *testing.T) {
	candles := trendPullbackTestCandles(230, 100, 0.5)
	// Force one shallow pullback into the EMA band while preserving the close
	// above EMA50; the next candle is the continuation trigger.
	candles[218].Low = candles[218].Close - 10

	cfg := DefaultBacktestFirstBTC15MTrendPullbackContinuationConfig()
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, FeePct: 0.0004, SlippagePct: 0.000116, MaxHoldBars: cfg.MaxHoldBars}
	signals, skips := btc15MTrendPullbackBuildSignals(candles, cfg, btCfg, []Split{{Name: fullSplitName}})

	var signal BTC15MTrendPullbackContinuationSignalRow
	found := false
	for _, row := range signals {
		if row.PrePositionCandidate && row.Side == Long {
			signal = row
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected long pre-position candidate, signals=%d skips=%v", len(signals), btc15MTrendPullbackSkipRows(skips))
	}
	if signal.DecisionIndex != 219 {
		t.Fatalf("decision index=%d, want 219", signal.DecisionIndex)
	}
	if signal.EntryIndex != signal.DecisionIndex+1 {
		t.Fatalf("entry index=%d, decision=%d", signal.EntryIndex, signal.DecisionIndex)
	}
	wantEntry := applySlippage(candles[signal.EntryIndex].Open, btCfg.SlippagePct, Long, true)
	if !trendPullbackAlmostEqual(signal.ExpectedEntryPrice, wantEntry, 1e-9) {
		t.Fatalf("expected entry=%f, want %f", signal.ExpectedEntryPrice, wantEntry)
	}
	wantTarget := signal.ExpectedEntryPrice + cfg.TargetR*(signal.ExpectedEntryPrice-signal.Stop)
	if !trendPullbackAlmostEqual(signal.Target, wantTarget, 1e-9) {
		t.Fatalf("target=%f, want %f", signal.Target, wantTarget)
	}
	if signal.ForwardLabelsAsInput || signal.UsesFutureRowsForFeatures || signal.DerivativesVetoUsed || signal.OptimizerSelectionUsed {
		t.Fatalf("signal used forbidden inputs: %+v", signal)
	}
}

func trendPullbackTestCandles(n int, start, step float64) []Candle {
	out := make([]Candle, n)
	ts := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range out {
		close := start + float64(i)*step
		out[i] = Candle{
			OpenTime:  ts.Add(time.Duration(i) * 15 * time.Minute),
			CloseTime: ts.Add(time.Duration(i+1)*15*time.Minute - time.Millisecond),
			Open:      close - 0.1,
			High:      close + 0.2,
			Low:       close - 0.2,
			Close:     close,
			Volume:    100,
		}
	}
	return out
}

func trendPullbackAlmostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
