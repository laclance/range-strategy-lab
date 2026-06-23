package lab

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"
)

const (
	SourceProductBinanceUSDMFutures = "binance-usdm-futures"
	SourceProductBinanceSpot        = "binance-spot"
)

const sourceInterval = 5 * time.Minute

type SourceValidationOptions struct {
	Product             string
	AllowSpotComparison bool
}

type SourceManifest struct {
	Path               string   `json:"path"`
	Venue              string   `json:"venue"`
	Product            string   `json:"product"`
	Symbol             string   `json:"symbol"`
	Interval           string   `json:"interval"`
	RowCount           int      `json:"row_count"`
	FirstOpenTime      string   `json:"first_open_time"`
	LastOpenTime       string   `json:"last_open_time"`
	Schema             []string `json:"schema"`
	TimestampSemantics string   `json:"timestamp_semantics"`
	FinalityRule       string   `json:"finality_rule"`
	GapCount           int      `json:"gap_count"`
	DuplicateCount     int      `json:"duplicate_count"`
	ZeroVolumeCount    int      `json:"zero_volume_count"`
	Monotonic          bool     `json:"monotonic"`
	ComparisonOnly     bool     `json:"comparison_only"`
	ValidationStatus   string   `json:"validation_status"`
	ValidationError    string   `json:"validation_error,omitempty"`
}

func LoadResearchSourceCSV(path string, opts SourceValidationOptions) ([]Candle, SourceManifest, error) {
	candles, header, err := LoadCSVWithHeader(path)
	if err != nil {
		return nil, SourceManifest{}, err
	}
	manifest, err := ValidateResearchSource(path, header, candles, opts)
	if err != nil {
		return nil, manifest, err
	}
	return candles, manifest, nil
}

func ValidateResearchSource(path string, header []string, candles []Candle, opts SourceValidationOptions) (SourceManifest, error) {
	product, productLabel, err := normalizeSourceProduct(opts.Product)
	if err != nil {
		return SourceManifest{}, err
	}

	manifest := SourceManifest{
		Path:               path,
		Venue:              "Binance",
		Product:            productLabel,
		Symbol:             "BTCUSDT",
		Interval:           "5m",
		Schema:             normalizeSchema(header),
		TimestampSemantics: "open_time",
		FinalityRule:       "closed 5m candles only; close_time must equal open_time plus 5m minus 1ms",
		Monotonic:          true,
		ComparisonOnly:     product == SourceProductBinanceSpot,
		ValidationStatus:   "accepted",
	}
	reject := func(format string, args ...any) (SourceManifest, error) {
		manifest.ValidationStatus = "rejected"
		manifest.ValidationError = fmt.Sprintf(format, args...)
		return manifest, fmt.Errorf("%s", manifest.ValidationError)
	}

	if err := validateSourcePath(path, product, opts.AllowSpotComparison); err != nil {
		return reject(err.Error())
	}
	if len(candles) == 0 {
		return reject("source CSV contains no candles")
	}

	manifest.RowCount = len(candles)
	manifest.FirstOpenTime = candles[0].OpenTime.UTC().Format(time.RFC3339)
	manifest.LastOpenTime = candles[len(candles)-1].OpenTime.UTC().Format(time.RFC3339)

	seen := make(map[time.Time]struct{}, len(candles))
	for i, candle := range candles {
		if err := validateCandleValues(i, candle); err != nil {
			return reject(err.Error())
		}
		if candle.Volume == 0 {
			manifest.ZeroVolumeCount++
		}

		openTime := candle.OpenTime.UTC()
		if _, ok := seen[openTime]; ok {
			manifest.DuplicateCount++
		}
		seen[openTime] = struct{}{}

		expectedClose := openTime.Add(sourceInterval - time.Millisecond)
		if !candle.CloseTime.UTC().Equal(expectedClose) {
			return reject("row %d close_time=%s does not match expected closed 5m candle close_time=%s", i, candle.CloseTime.UTC().Format(time.RFC3339Nano), expectedClose.Format(time.RFC3339Nano))
		}

		if i == 0 {
			continue
		}
		prevOpen := candles[i-1].OpenTime.UTC()
		if !openTime.After(prevOpen) {
			manifest.Monotonic = false
			continue
		}
		diff := openTime.Sub(prevOpen)
		if diff != sourceInterval {
			if diff > sourceInterval && diff%sourceInterval == 0 {
				manifest.GapCount += int(diff/sourceInterval) - 1
				continue
			}
			return reject("irregular 5m interval between %s and %s: got %s", prevOpen.Format(time.RFC3339), openTime.Format(time.RFC3339), diff)
		}
	}

	if manifest.DuplicateCount > 0 {
		return reject("source CSV contains %d duplicate open_time value(s)", manifest.DuplicateCount)
	}
	if !manifest.Monotonic {
		return reject("source CSV open_time values must be strictly increasing")
	}
	if manifest.GapCount > 0 {
		return reject("source CSV has %d missing 5m candle(s)", manifest.GapCount)
	}

	return manifest, nil
}

func normalizeSourceProduct(product string) (string, string, error) {
	switch strings.TrimSpace(strings.ToLower(product)) {
	case SourceProductBinanceUSDMFutures:
		return SourceProductBinanceUSDMFutures, "Binance USDT-M futures", nil
	case SourceProductBinanceSpot:
		return SourceProductBinanceSpot, "Binance spot", nil
	default:
		return "", "", fmt.Errorf("invalid -source-product %q; use %q or %q", product, SourceProductBinanceUSDMFutures, SourceProductBinanceSpot)
	}
}

func validateSourcePath(path string, product string, allowSpotComparison bool) error {
	base := strings.ToLower(filepath.Base(path))
	if !strings.Contains(base, "btcusdt") {
		return fmt.Errorf("source CSV path must identify BTCUSDT")
	}
	if !strings.Contains(base, "5m") {
		return fmt.Errorf("source CSV path must identify 5m candles")
	}

	if product == SourceProductBinanceUSDMFutures && strings.Contains(base, "spot") {
		return fmt.Errorf("spot-looking CSV path %q cannot be used for Binance USDT-M futures research", path)
	}
	if product == SourceProductBinanceSpot {
		if !allowSpotComparison {
			return fmt.Errorf("spot CSVs are comparison-only; rerun with -allow-spot-comparison to make that explicit")
		}
		if strings.Contains(base, "futures") || strings.Contains(base, "usdm") || strings.Contains(base, "usd-m") || strings.Contains(base, "_um_") {
			return fmt.Errorf("futures-looking CSV path %q cannot be labeled as Binance spot", path)
		}
	}
	return nil
}

func normalizeSchema(header []string) []string {
	out := make([]string, len(header))
	for i, value := range header {
		out[i] = strings.ToLower(strings.TrimSpace(value))
	}
	return out
}

func validateCandleValues(index int, candle Candle) error {
	prices := map[string]float64{
		"open":  candle.Open,
		"high":  candle.High,
		"low":   candle.Low,
		"close": candle.Close,
	}
	for name, value := range prices {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("row %d %s must be finite", index, name)
		}
		if value <= 0 {
			return fmt.Errorf("row %d %s must be positive", index, name)
		}
	}
	if math.IsNaN(candle.Volume) || math.IsInf(candle.Volume, 0) {
		return fmt.Errorf("row %d volume must be finite", index)
	}
	if candle.Volume < 0 {
		return fmt.Errorf("row %d volume cannot be negative", index)
	}
	if candle.Low > candle.High {
		return fmt.Errorf("row %d low cannot exceed high", index)
	}
	if candle.Open < candle.Low || candle.Open > candle.High {
		return fmt.Errorf("row %d open must be inside high/low range", index)
	}
	if candle.Close < candle.Low || candle.Close > candle.High {
		return fmt.Errorf("row %d close must be inside high/low range", index)
	}
	return nil
}
