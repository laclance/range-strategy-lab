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

func TestSizingCapsNotionalAtOneTimesEquity(t *testing.T) {
	size := positionSize(1000, 100, 1, BacktestConfig{RiskPct: 0.02, MaxNotionalPct: 1})
	if size != 10 {
		t.Fatalf("size=%v, want 10 capped by notional", size)
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
