package lab

import (
	"math"
	"testing"
	"time"
)

func TestRunBacktestUsesNextBarOpenAndStopFirst(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 101, 106, 94, 102),
	}
	strategy := fixedStrategy{signal: Signal{Side: Long, Stop: 95, Target: 105, MaxHoldBars: 10, Reason: "test"}}
	result := RunBacktest(candles, strategy, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
	})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	tr := result.Trades[0]
	if tr.EntryPrice != 101 {
		t.Fatalf("entry=%v, want next bar open 101", tr.EntryPrice)
	}
	if tr.Reason != "stop_loss" || tr.ExitPrice != 95 {
		t.Fatalf("expected stop-first at 95, got %+v", tr)
	}
}

func TestRunBacktestShortTargetAndSlippage(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 88, 90),
	}
	strategy := fixedStrategy{signal: Signal{Side: Short, Stop: 105, Target: 90, MaxHoldBars: 10, Reason: "short"}}
	result := RunBacktest(candles, strategy, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
		SlippagePct:    0.001,
	})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	tr := result.Trades[0]
	if tr.Side != Short || tr.Reason != "take_profit" {
		t.Fatalf("bad short trade: %+v", tr)
	}
	if tr.EntryPrice >= candles[1].Open {
		t.Fatalf("short entry should slip below raw open, got %v raw %v", tr.EntryPrice, candles[1].Open)
	}
	if tr.ExitPrice <= tr.Target {
		t.Fatalf("short exit should slip above raw target, got %v raw %v", tr.ExitPrice, tr.Target)
	}
}

func TestRunBacktestUsesTimeStopAndConfigMaxHoldDefault(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 101, 99, 100),
		testCandle(2, 100, 101, 99, 100),
	}
	strategy := fixedStrategy{signal: Signal{Side: Long, Stop: 90, Target: 110, Reason: "time"}}
	result := RunBacktest(candles, strategy, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
		MaxHoldBars:    1,
	})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	tr := result.Trades[0]
	if tr.Reason != "time_stop" || tr.HoldBars != 1 || tr.CloseIndex != 2 {
		t.Fatalf("bad time stop trade: %+v", tr)
	}
}

func TestRunBacktestForceClosesOpenPosition(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 101, 99, 100),
	}
	strategy := fixedStrategy{signal: Signal{Side: Long, Stop: 90, Target: 110, MaxHoldBars: 0, Reason: "force"}}
	result := RunBacktest(candles, strategy, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
		MaxHoldBars:    0,
	})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	if result.Trades[0].Reason != "force_close" {
		t.Fatalf("expected force close, got %+v", result.Trades[0])
	}
	if len(result.EquityCurve) != len(candles)+2 {
		t.Fatalf("equity curve len=%d, want %d including force close point", len(result.EquityCurve), len(candles)+2)
	}
}

func TestRunBacktestSkipsInvalidSignalsAndZeroSizing(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
	}
	invalidSignals := []Signal{
		{Side: Direction("flat"), Stop: 90, Target: 110},
		{Side: Long, Stop: 0, Target: 110},
		{Side: Long, Stop: 100, Target: 90},
		{Side: Short, Stop: 90, Target: 100},
	}
	for _, sig := range invalidSignals {
		result := RunBacktest(candles, fixedStrategy{signal: sig}, BacktestConfig{
			StartBalance:   1000,
			RiskPct:        0.01,
			MaxNotionalPct: 1,
		})
		if len(result.Trades) != 0 {
			t.Fatalf("invalid signal %+v produced trades %+v", sig, result.Trades)
		}
	}

	result := RunBacktest(candles, fixedStrategy{signal: Signal{Side: Long, Stop: 100, Target: 110}}, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
	})
	if len(result.Trades) != 0 {
		t.Fatalf("zero stop distance produced trades %+v", result.Trades)
	}
}

func TestSizingCapsNotionalAtOneTimesEquity(t *testing.T) {
	size := positionSize(1000, 100, 1, BacktestConfig{RiskPct: 0.02, MaxNotionalPct: 1})
	if size != 10 {
		t.Fatalf("size=%v, want 10 capped by notional", size)
	}
}

func TestPositionSizeRejectsInvalidInputsAndAllowsUncappedSize(t *testing.T) {
	for _, size := range []float64{
		positionSize(0, 100, 1, BacktestConfig{RiskPct: 0.01}),
		positionSize(1000, 0, 1, BacktestConfig{RiskPct: 0.01}),
		positionSize(1000, 100, 0, BacktestConfig{RiskPct: 0.01}),
		positionSize(1000, 100, 1, BacktestConfig{RiskPct: 0}),
	} {
		if size != 0 {
			t.Fatalf("invalid sizing returned %v, want 0", size)
		}
	}
	size := positionSize(1000, 100, 5, BacktestConfig{RiskPct: 0.01, MaxNotionalPct: 0})
	if size != 2 {
		t.Fatalf("uncapped size=%v, want 2", size)
	}
}

func TestCostsAndSlippage(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 112, 99, 110),
	}
	strategy := fixedStrategy{signal: Signal{Side: Long, Stop: 95, Target: 110, MaxHoldBars: 10, Reason: "test"}}
	result := RunBacktest(candles, strategy, BacktestConfig{
		StartBalance:   1000,
		RiskPct:        0.01,
		MaxNotionalPct: 1,
		FeePct:         0.001,
		SlippagePct:    0.001,
	})
	if len(result.Trades) != 1 {
		t.Fatalf("trades=%d, want 1", len(result.Trades))
	}
	tr := result.Trades[0]
	if tr.Fees <= 0 || tr.Slippage <= 0 {
		t.Fatalf("expected fees and slippage, got %+v", tr)
	}
	if tr.NetPnL >= tr.GrossPnL {
		t.Fatalf("net should be below gross after fees, got gross=%v net=%v", tr.GrossPnL, tr.NetPnL)
	}
}

func TestApplySlippageAllBranches(t *testing.T) {
	cases := []struct {
		side    Direction
		entry   bool
		want    float64
		message string
	}{
		{side: Long, entry: true, want: 101, message: "long entry"},
		{side: Long, entry: false, want: 99, message: "long exit"},
		{side: Short, entry: true, want: 99, message: "short entry"},
		{side: Short, entry: false, want: 101, message: "short exit"},
	}
	for _, tc := range cases {
		got := applySlippage(100, 0.01, tc.side, tc.entry)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Fatalf("%s slippage=%v, want %v", tc.message, got, tc.want)
		}
	}
	if got := applySlippage(100, 0, Long, true); got != 100 {
		t.Fatalf("zero slippage=%v, want 100", got)
	}
}

func TestMarkToMarketNilLongAndShort(t *testing.T) {
	if got := markToMarket(1000, nil, 100); got != 1000 {
		t.Fatalf("nil position MTM=%v, want 1000", got)
	}
	long := &Position{Side: Long, EntryPrice: 100, Size: 2, EntryFee: 1}
	if got := markToMarket(1000, long, 105); got != 1009 {
		t.Fatalf("long MTM=%v, want 1009", got)
	}
	short := &Position{Side: Short, EntryPrice: 100, Size: 2, EntryFee: 1}
	if got := markToMarket(1000, short, 95); got != 1009 {
		t.Fatalf("short MTM=%v, want 1009", got)
	}
}

func TestSummarizeSplitsByCloseTimeAndSide(t *testing.T) {
	trades := []Trade{
		{Side: Long, ExitTime: "2022-01-01T00:00:00Z", NetPnL: 10, GrossPnL: 10, HoldBars: 2},
		{Side: Short, ExitTime: "2024-01-01T00:00:00Z", NetPnL: -5, GrossPnL: -5, HoldBars: 4},
	}
	rows := SummarizeSplits(trades, 1000, DefaultSplits())
	got := map[string]SummaryRow{}
	for _, row := range rows {
		got[row.Split+"|"+row.Side] = row
	}
	if got["2021_2022_stress|all"].NetPnL != 10 {
		t.Fatalf("bad stress row: %+v", got["2021_2022_stress|all"])
	}
	if got["full_2021_2026|short"].NetPnL != -5 {
		t.Fatalf("bad full short row: %+v", got["full_2021_2026|short"])
	}
}

func TestMaxDrawdown(t *testing.T) {
	got := MaxDrawdown([]float64{1000, 1100, 990, 1200})
	if math.Abs(got-0.1) > 1e-9 {
		t.Fatalf("dd=%v, want 0.1", got)
	}
}

func TestMaxDrawdownHandlesEmptyAndNonPositivePeak(t *testing.T) {
	if got := MaxDrawdown(nil); got != 0 {
		t.Fatalf("empty drawdown=%v, want 0", got)
	}
	if got := MaxDrawdown([]float64{0, -10}); got != 0 {
		t.Fatalf("non-positive peak drawdown=%v, want 0", got)
	}
}

func TestStrategyContextAndEmptyStrategy(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	ctx := StrategyContext{Candles: candles, Index: 0}
	if ctx.Candle() != candles[0] {
		t.Fatalf("context candle mismatch")
	}
	strategy := EmptyStrategy{}
	if strategy.Name() != "empty" {
		t.Fatalf("empty strategy name=%q", strategy.Name())
	}
	if sig, ok := strategy.OnCandle(ctx); ok || sig != (Signal{}) {
		t.Fatalf("empty strategy signal=%+v ok=%v, want zero false", sig, ok)
	}
}

type fixedStrategy struct {
	signal Signal
}

func (f fixedStrategy) Name() string {
	return "fixed"
}

func (f fixedStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	return f.signal, ctx.Index == 0
}

func testCandle(i int, open, high, low, close float64) Candle {
	openTime := time.Date(2021, 1, 1, 0, i*5, 0, 0, time.UTC)
	return Candle{
		OpenTime:  openTime,
		CloseTime: openTime.Add(5*time.Minute - time.Millisecond),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    1,
	}
}
