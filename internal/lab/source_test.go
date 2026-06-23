package lab

import (
	"strings"
	"testing"
	"time"
)

func TestValidateResearchSourceAcceptsFuturesAndPopulatesManifest(t *testing.T) {
	candles := sourceTestCandles(3)
	manifest, err := ValidateResearchSource(
		"../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv",
		sourceTestHeader(),
		candles,
		SourceValidationOptions{Product: SourceProductBinanceUSDMFutures},
	)
	if err != nil {
		t.Fatal(err)
	}

	if manifest.Product != "Binance USDT-M futures" || manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" {
		t.Fatalf("unexpected identity fields: %+v", manifest)
	}
	if manifest.RowCount != 3 || manifest.FirstOpenTime != "2021-01-01T00:00:00Z" || manifest.LastOpenTime != "2021-01-01T00:10:00Z" {
		t.Fatalf("unexpected coverage fields: %+v", manifest)
	}
	if manifest.GapCount != 0 || manifest.DuplicateCount != 0 || !manifest.Monotonic || manifest.ComparisonOnly {
		t.Fatalf("unexpected validation fields: %+v", manifest)
	}
	if manifest.ValidationStatus != "accepted" || len(manifest.Schema) != len(sourceTestHeader()) {
		t.Fatalf("unexpected manifest status/schema: %+v", manifest)
	}
}

func TestValidateResearchSourceRejectsSpotUnlessExplicitComparison(t *testing.T) {
	candles := sourceTestCandles(2)
	_, err := ValidateResearchSource(
		"data/btcusdt_spot_5m_2021_2026.csv",
		sourceTestHeader(),
		candles,
		SourceValidationOptions{Product: SourceProductBinanceSpot},
	)
	if err == nil || !strings.Contains(err.Error(), "comparison-only") {
		t.Fatalf("expected comparison-only rejection, got %v", err)
	}

	manifest, err := ValidateResearchSource(
		"data/btcusdt_spot_5m_2021_2026.csv",
		sourceTestHeader(),
		candles,
		SourceValidationOptions{Product: SourceProductBinanceSpot, AllowSpotComparison: true},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !manifest.ComparisonOnly || manifest.Product != "Binance spot" {
		t.Fatalf("expected explicit spot comparison manifest, got %+v", manifest)
	}
}

func TestValidateResearchSourceRejectsSpotPathDuringFuturesRun(t *testing.T) {
	_, err := ValidateResearchSource(
		"data/btcusdt_spot_5m_2021_2026.csv",
		sourceTestHeader(),
		sourceTestCandles(2),
		SourceValidationOptions{Product: SourceProductBinanceUSDMFutures},
	)
	if err == nil || !strings.Contains(err.Error(), "spot-looking") {
		t.Fatalf("expected spot-looking futures rejection, got %v", err)
	}
}

func TestValidateResearchSourceRejectsGapDuplicateAndIrregularCadence(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func([]Candle)
		wantErr string
	}{
		{
			name: "gap",
			mutate: func(candles []Candle) {
				shiftCandle(&candles[2], 5*time.Minute)
			},
			wantErr: "missing 5m candle",
		},
		{
			name: "duplicate",
			mutate: func(candles []Candle) {
				candles[2].OpenTime = candles[1].OpenTime
				candles[2].CloseTime = candles[1].CloseTime
			},
			wantErr: "duplicate open_time",
		},
		{
			name: "irregular",
			mutate: func(candles []Candle) {
				shiftCandle(&candles[2], time.Minute)
			},
			wantErr: "irregular 5m interval",
		},
		{
			name: "bad value",
			mutate: func(candles []Candle) {
				candles[1].Close = candles[1].High + 1
			},
			wantErr: "close must be inside",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			candles := sourceTestCandles(3)
			tc.mutate(candles)
			manifest, err := ValidateResearchSource(
				"data/btcusdt_futures_um_5m_2021_2026.csv",
				sourceTestHeader(),
				candles,
				SourceValidationOptions{Product: SourceProductBinanceUSDMFutures},
			)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected %q error, got %v", tc.wantErr, err)
			}
			if manifest.ValidationStatus != "rejected" {
				t.Fatalf("expected rejected manifest, got %+v", manifest)
			}
		})
	}
}

func sourceTestHeader() []string {
	return []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}
}

func sourceTestCandles(count int) []Candle {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	candles := make([]Candle, count)
	for i := range candles {
		open := start.Add(time.Duration(i) * sourceInterval)
		candles[i] = Candle{
			OpenTime:  open,
			CloseTime: open.Add(sourceInterval - time.Millisecond),
			Open:      100,
			High:      110,
			Low:       90,
			Close:     105,
			Volume:    1,
		}
	}
	return candles
}

func shiftCandle(candle *Candle, delta time.Duration) {
	candle.OpenTime = candle.OpenTime.Add(delta)
	candle.CloseTime = candle.CloseTime.Add(delta)
}
