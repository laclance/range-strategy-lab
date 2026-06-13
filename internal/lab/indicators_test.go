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

func TestIndicatorWarmupAndInvalidPeriodPaths(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	checkNaN := func(name string, values []float64) {
		if len(values) != len(candles) {
			t.Fatalf("%s len=%d, want %d", name, len(values), len(candles))
		}
		if !math.IsNaN(values[0]) {
			t.Fatalf("%s[0]=%v, want NaN", name, values[0])
		}
	}
	checkNaN("ATR invalid", ATR(candles, 0))
	checkNaN("ATR warmup", ATR(candles, 2))
	checkNaN("Donchian invalid", DonchianWidth(candles, 0))
	checkNaN("Donchian warmup", DonchianWidth(candles, 2))
	checkNaN("Bollinger invalid", BollingerWidth(candles, 0))
	checkNaN("Bollinger warmup", BollingerWidth(candles, 2))
	checkNaN("ADX invalid", ADX(candles, 0))
	checkNaN("ADX warmup", ADX(candles, 2))
}

func TestNormalizedATRSkipsNonPositiveClose(t *testing.T) {
	candles := make([]Candle, 2)
	candles[0] = testCandle(0, 0, 1, -1, 0)
	candles[1] = testCandle(1, 0, 1, -1, 0)
	got := NormalizedATR(candles, 1)
	if !math.IsNaN(got[0]) || !math.IsNaN(got[1]) {
		t.Fatalf("normalized ATR with non-positive close=%v, want NaNs", got)
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

func TestBollingerWidthSkipsNonPositiveMean(t *testing.T) {
	candles := []Candle{
		testCandle(0, -1, -1, -1, -1),
		testCandle(1, 1, 1, 1, 1),
	}
	got := BollingerWidth(candles, 2)
	if !math.IsNaN(got[1]) {
		t.Fatalf("BollingerWidth[1]=%v, want NaN for zero mean", got[1])
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

func TestADXHandlesFlatAndGaplessSeries(t *testing.T) {
	flat := make([]Candle, 6)
	for i := range flat {
		flat[i] = testCandle(i, 100, 100, 100, 100)
	}
	got := ADX(flat, 2)
	if got[3] != 0 {
		t.Fatalf("flat ADX=%v, want 0", got[3])
	}

	rangeOnly := make([]Candle, 6)
	for i := range rangeOnly {
		rangeOnly[i] = testCandle(i, 100, 101, 99, 100)
	}
	got = ADX(rangeOnly, 2)
	if got[3] != 0 {
		t.Fatalf("range-only ADX=%v, want 0 when DI denominator is zero", got[3])
	}
}

func almostEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}
