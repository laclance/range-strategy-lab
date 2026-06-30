package lab

import (
	"strings"
	"testing"
)

func TestBTCPostCompressionFixedBacktestBuildSignalsUsesFixedPriorFacts(t *testing.T) {
	candles := []Candle{
		postCompression15MTestCandle(0, 100, 110, 90, 100, 10),
		postCompression15MTestCandle(1, 100, 110, 90, 100, 10),
		postCompression15MTestCandle(2, 100, 100.5, 99.5, 100, 10),
		postCompression15MTestCandle(3, 100, 100.5, 99.5, 100, 10),
		postCompression15MTestCandle(4, 100, 102.2, 100, 102, 12),
		postCompression15MTestCandle(5, 102, 110, 102, 109, 12),
		postCompression15MTestCandle(6, 109, 111, 108, 110, 10),
		postCompression15MTestCandle(7, 110, 112, 109, 111, 10),
	}
	cfg := DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
	cfg.LookbackBars = 2
	cfg.CompressionPercentile = 0.8
	cfg.PercentileReferenceBars = 2
	cfg.BreakoutATRMultiple = 0.2
	cfg.ATRPeriod = 2
	cfg.MaxHoldBars = 2
	cfg.ExpectedRawCandidateRows = 1
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, MaxHoldBars: 2}

	signals, _ := btcPostCompressionFixedBacktestBuildSignals(candles, cfg, btCfg, []Split{{Name: fullSplitName}})
	if len(signals) != 1 {
		t.Fatalf("signals=%d, want 1: %+v", len(signals), signals)
	}
	sig := signals[0]
	if sig.DecisionIndex != 4 || sig.Side != Long || sig.VolumeMode != BTC15MPostCompressionVolumeNone {
		t.Fatalf("unexpected fixed-cell signal: %+v", sig)
	}
	if !postCompressionAlmostEqual(sig.RangeHigh, 100.5) || !postCompressionAlmostEqual(sig.RangeLow, 99.5) {
		t.Fatalf("prior range should use [d-2,d-1], got high=%v low=%v", sig.RangeHigh, sig.RangeLow)
	}
	if !postCompressionAlmostEqual(sig.PriorATR14, 5.75) {
		t.Fatalf("ATR should use d-1, got %v", sig.PriorATR14)
	}
	if !postCompressionAlmostEqual(sig.ExpectedEntryPrice, 102) || !postCompressionAlmostEqual(sig.Stop, 96.25) || !postCompressionAlmostEqual(sig.Target, 113.5) {
		t.Fatalf("unexpected entry geometry: entry=%v stop=%v target=%v", sig.ExpectedEntryPrice, sig.Stop, sig.Target)
	}
	if !sig.PrePositionCandidate || sig.SkippedReason != "" {
		t.Fatalf("signal should be a valid pre-position candidate: %+v", sig)
	}
}

func TestBTCPostCompressionFixedBacktestRunnerStopFirstOverlapAndStressNet(t *testing.T) {
	candles := []Candle{
		postCompression15MTestCandle(0, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(1, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(2, 100, 103, 98, 100, 10),
		postCompression15MTestCandle(3, 100, 101, 99, 100, 10),
	}
	cfg := DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
	cfg.MaxHoldBars = 2
	btCfg := BacktestConfig{StartBalance: 1000, RiskPct: 0.01, MaxNotionalPct: 1, FeePct: 0, SlippagePct: 0, MaxHoldBars: 2}
	signals := []BTC15MPostCompressionFixedBacktestSignalRow{
		fixedBacktestTestSignal(1, 2, 100, 99, 102),
	}
	skips := map[btcPostCompressionFixedSkipKey]int{}

	signals, trades := btcPostCompressionFixedBacktestRun(candles, signals, &skips, cfg, btCfg, []Split{{Name: fullSplitName}})
	if len(trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(trades))
	}
	if trades[0].Reason != "stop_loss" {
		t.Fatalf("stop-first ambiguity should exit at stop, got %s", trades[0].Reason)
	}
	if !signals[0].Executed {
		t.Fatalf("executed signal not marked")
	}

	overlapSignals := []BTC15MPostCompressionFixedBacktestSignalRow{
		fixedBacktestTestSignal(0, 1, 100, 90, 130),
		fixedBacktestTestSignal(1, 2, 100, 90, 130),
	}
	overlapCandles := []Candle{
		postCompression15MTestCandle(0, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(1, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(2, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(3, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(4, 100, 101, 99, 100, 10),
	}
	skips = map[btcPostCompressionFixedSkipKey]int{}
	overlapSignals, _ = btcPostCompressionFixedBacktestRun(overlapCandles, overlapSignals, &skips, cfg, btCfg, []Split{{Name: fullSplitName}})
	if overlapSignals[1].SkippedReason != "open_position_already_active" {
		t.Fatalf("second signal should be skipped for open position, got %+v", overlapSignals[1])
	}

	stressTrade := Trade{
		Side:       Long,
		EntryTime:  candles[1].OpenTime.Format(timeLayout),
		ExitTime:   candles[2].CloseTime.Format(timeLayout),
		OpenIndex:  1,
		CloseIndex: 2,
		EntryPrice: 100,
		ExitPrice:  103,
		Stop:       99,
		Target:     103,
		Size:       2,
		GrossPnL:   6,
		NetPnL:     5,
		Fees:       1,
		Slippage:   0.5,
		Reason:     "take_profit",
		Signal:     signals[0].SignalID,
		HoldBars:   1,
	}
	rows := btcPostCompressionFixedBacktestTradeRows([]Trade{stressTrade}, signals, []Split{{Name: fullSplitName}})
	if len(rows) != 1 {
		t.Fatalf("stress trade rows=%d, want 1", len(rows))
	}
	if !postCompressionAlmostEqual(rows[0].ExtraSlippageStressNetPnL, 4.5) {
		t.Fatalf("stress net should subtract slippage again, got %v", rows[0].ExtraSlippageStressNetPnL)
	}
	if !postCompressionAlmostEqual(rows[0].ExtraSlippageStressNetR, 2.25) {
		t.Fatalf("stress R got %v", rows[0].ExtraSlippageStressNetR)
	}
}

func TestBTCPostCompressionFixedBacktestFalsificationStopStates(t *testing.T) {
	passing := BTC15MPostCompressionFixedBacktestFalsification{
		SourceResamplePass:               true,
		CandidateIdentityPass:            true,
		LeakagePass:                      true,
		TradeCountPass:                   true,
		GrossEdgePass:                    true,
		CostedEdgePass:                   true,
		ProfitFactorPass:                 true,
		DrawdownPass:                     true,
		RobustnessPass:                   true,
		OptimizerContaminationPass:       true,
		ClosedFamilyProtectionPass:       true,
		DerivativesVetoContaminationPass: true,
	}
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStatePassedNeedsReview {
		t.Fatalf("passing stop state=%s", got)
	}
	passing.TradeCountPass = false
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStateFailedNoUsableStrategy {
		t.Fatalf("trade-count failure stop state=%s", got)
	}
	passing.TradeCountPass = true
	passing.CandidateIdentityPass = false
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStateFailedSourceOrResample {
		t.Fatalf("candidate mismatch stop state=%s", got)
	}
	passing.CandidateIdentityPass = true
	passing.OptimizerContaminationPass = false
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStateRejectedOptimizerContam {
		t.Fatalf("optimizer stop state=%s", got)
	}
	passing.OptimizerContaminationPass = true
	passing.ClosedFamilyProtectionPass = false
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStateRejectedClosedReslice {
		t.Fatalf("closed-family stop state=%s", got)
	}
	passing.ClosedFamilyProtectionPass = true
	passing.DerivativesVetoContaminationPass = false
	if got := btcPostCompressionFixedBacktestFalsification(passing).StopState; got != BTC15MPostCompressionFixedBacktestStopStateRejectedVetoContam {
		t.Fatalf("veto stop state=%s", got)
	}
}

func TestBTCPostCompressionFixedBacktestFullRunRejectsSourceMismatch(t *testing.T) {
	cfg := DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
	manifest := SourceManifest{
		Path:             cfg.ApprovedSourcePath,
		Product:          "Binance spot",
		Symbol:           RangeUniverseSymbolBTCUSDT,
		Interval:         "5m",
		ValidationStatus: "accepted",
		ComparisonOnly:   true,
	}
	result, err := RunFuturesBTC15MPostCompressionL192Q20M020NoneLongH48Backtest(nil, manifest, cfg, BacktestConfig{}, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != BTC15MPostCompressionFixedBacktestStopStateFailedSourceOrResample {
		t.Fatalf("source mismatch stop state=%s", result.StopState)
	}
	if !strings.Contains(strings.Join(result.Falsification.FailureReasons, ","), "source_or_resample_failed") {
		t.Fatalf("missing source failure reason: %+v", result.Falsification)
	}
}

func TestBTCPostCompressionFixedSummaryUsesStressGates(t *testing.T) {
	cfg := DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
	cfg.MinFullTrades = 1
	cfg.MinSplitTrades = 0
	cfg.FullStressProfitFactorMin = 1.0
	cfg.SplitStressProfitFactorMin = 1.0
	signals := []BTC15MPostCompressionFixedBacktestSignalRow{fixedBacktestTestSignal(0, 1, 100, 99, 103)}
	trades := []BTC15MPostCompressionFixedBacktestTradeRow{{
		SignalID:                  signals[0].SignalID,
		CandidateID:               BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
		Side:                      Long,
		ExitTime:                  postCompression15MTestCandle(1, 100, 101, 99, 100, 1).CloseTime.Format(timeLayout),
		InitialRisk:               2,
		GrossPnL:                  3,
		EngineNetPnL:              2,
		ExtraSlippageStressNetPnL: 1.5,
		Slippage:                  0.5,
		ExitReason:                "take_profit",
	}}
	rows := btcPostCompressionFixedBacktestSummaryRows(signals, nil, trades, cfg, 1000, []Split{{Name: fullSplitName}})
	full := btcPostCompressionFixedSummaryByKey(rows)[fullSplitName+"|"+string(Long)]
	if full.EngineNetPnL != 2 || full.ExtraSlippageStressNetPnL != 1.5 {
		t.Fatalf("summary did not preserve engine/stress net: %+v", full)
	}
	if full.StressProfitFactor != 999.99 {
		t.Fatalf("stress PF=%v", full.StressProfitFactor)
	}
}

func fixedBacktestTestSignal(decisionIndex, entryIndex int, entry, stop, target float64) BTC15MPostCompressionFixedBacktestSignalRow {
	decision := postCompression15MTestCandle(decisionIndex, 100, 101, 99, 100, 10)
	entryCandle := postCompression15MTestCandle(entryIndex, entry, entry+1, entry-1, entry, 10)
	return BTC15MPostCompressionFixedBacktestSignalRow{
		SignalID:             BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID + "_test_" + string(rune('a'+decisionIndex)),
		CandidateID:          BTC15MPostCompressionL192Q20M020NoneLongH48CandidateID,
		Timeframe:            RangeDiscoveryTimeframe15m,
		Split:                fullSplitName,
		DecisionIndex:        decisionIndex,
		DecisionCloseTime:    decision.CloseTime.Format(timeLayout),
		Side:                 Long,
		EntryIndex:           entryIndex,
		EntryOpenTime:        entryCandle.OpenTime.Format(timeLayout),
		EntryOpen:            entry,
		ExpectedEntryPrice:   entry,
		Stop:                 stop,
		Target:               target,
		MaxHoldBars:          2,
		EntryGeometryValid:   true,
		PrePositionCandidate: true,
	}
}
