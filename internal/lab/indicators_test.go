package lab

import (
	"math"
	"testing"
)

func TestATR14ReturnsNormalizedWilderATR(t *testing.T) {
	candles := make([]Candle, 14)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	got := ATR14(candles)
	for i := 0; i < 13; i++ {
		if !math.IsNaN(got[i]) {
			t.Fatalf("ATR14[%d]=%v, want NaN warmup", i, got[i])
		}
	}
	if !almostEqual(got[13], 0.02) {
		t.Fatalf("ATR14[13]=%v, want 0.02", got[13])
	}
}

func TestDonchian20Width(t *testing.T) {
	candles := make([]Candle, 20)
	for i := range candles {
		candles[i] = testCandle(i, 100, 100+float64(i), 90, 100)
	}
	got := Donchian20Width(candles)
	if !almostEqual(got[19], 0.29) {
		t.Fatalf("Donchian20Width[19]=%v, want 0.29", got[19])
	}
}

func TestBollinger20Width(t *testing.T) {
	candles := make([]Candle, 20)
	for i := range candles {
		close := float64(i + 1)
		candles[i] = testCandle(i, close, close, close, close)
	}
	got := Bollinger20Width(candles)
	want := 4 * math.Sqrt(33.25) / 10.5
	if !almostEqual(got[19], want) {
		t.Fatalf("Bollinger20Width[19]=%v, want %v", got[19], want)
	}
}

func TestADX14ProducesTrendStrengthWithoutDependency(t *testing.T) {
	candles := make([]Candle, 30)
	for i := range candles {
		base := 100 + float64(i)
		candles[i] = testCandle(i, base, base+1, base-1, base+0.5)
	}
	got := ADX14(candles)
	if !validNumber(got[27]) || got[27] < 90 {
		t.Fatalf("ADX14[27]=%v, want strong valid trend", got[27])
	}
}

func almostEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}
