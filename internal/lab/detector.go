package lab

import (
	"fmt"
	"math"
	"sort"
)

type RangeDetector interface {
	Name() string
	Classify(candles []Candle) ([]RangeClassification, error)
}

type RangeDetectorConfig struct {
	ATRPeriod          int
	DonchianPeriod     int
	BollingerPeriod    int
	ADXPeriod          int
	LookbackDays       int
	BarsPerDay         int
	Percentile         float64
	MinConsecutiveBars int
	UseBollinger       bool
	UseADX             bool
}

type CompressionRangeDetector struct {
	Config RangeDetectorConfig
}

type RangeClassification struct {
	Index     int    `json:"index"`
	CloseTime string `json:"close_time"`
	RawActive bool   `json:"raw_active"`
	Active    bool   `json:"active"`
}

func DefaultCompressionRangeDetectorConfig() RangeDetectorConfig {
	return RangeDetectorConfig{
		ATRPeriod:          14,
		DonchianPeriod:     20,
		BollingerPeriod:    20,
		ADXPeriod:          14,
		LookbackDays:       20,
		BarsPerDay:         24 * 12,
		Percentile:         0.30,
		MinConsecutiveBars: 12,
		UseBollinger:       true,
		UseADX:             false,
	}
}

func (d CompressionRangeDetector) Name() string {
	return "compression_range"
}

func (d CompressionRangeDetector) Classify(candles []Candle) ([]RangeClassification, error) {
	cfg := d.Config.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	lookbackBars := cfg.LookbackDays * cfg.BarsPerDay
	atr := NormalizedATR(candles, cfg.ATRPeriod)
	donchian := DonchianWidth(candles, cfg.DonchianPeriod)

	atrThresholds := rollingPriorPercentile(atr, lookbackBars, cfg.Percentile)
	donchianThresholds := rollingPriorPercentile(donchian, lookbackBars, cfg.Percentile)
	var bollinger, bollingerThresholds []float64
	if cfg.UseBollinger {
		bollinger = BollingerWidth(candles, cfg.BollingerPeriod)
		bollingerThresholds = rollingPriorPercentile(bollinger, lookbackBars, cfg.Percentile)
	}
	var adx, adxThresholds []float64
	if cfg.UseADX {
		adx = ADX(candles, cfg.ADXPeriod)
		adxThresholds = rollingPriorPercentile(adx, lookbackBars, cfg.Percentile)
	}

	raw := make([]bool, len(candles))
	for i := range candles {
		raw[i] = metricLow(atr[i], atrThresholds[i]) &&
			metricLow(donchian[i], donchianThresholds[i])
		if raw[i] && cfg.UseBollinger {
			raw[i] = metricLow(bollinger[i], bollingerThresholds[i])
		}
		if raw[i] && cfg.UseADX {
			raw[i] = metricLow(adx[i], adxThresholds[i])
		}
	}
	active := markConsecutiveActive(raw, cfg.MinConsecutiveBars)

	out := make([]RangeClassification, len(candles))
	for i, c := range candles {
		out[i] = RangeClassification{
			Index:     i,
			CloseTime: c.CloseTime.Format(timeLayout),
			RawActive: raw[i],
			Active:    active[i],
		}
	}
	return out, nil
}

func (cfg RangeDetectorConfig) withDefaults() RangeDetectorConfig {
	defaults := DefaultCompressionRangeDetectorConfig()
	if cfg == (RangeDetectorConfig{}) {
		return defaults
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.DonchianPeriod == 0 {
		cfg.DonchianPeriod = defaults.DonchianPeriod
	}
	if cfg.BollingerPeriod == 0 {
		cfg.BollingerPeriod = defaults.BollingerPeriod
	}
	if cfg.ADXPeriod == 0 {
		cfg.ADXPeriod = defaults.ADXPeriod
	}
	if cfg.LookbackDays == 0 {
		cfg.LookbackDays = defaults.LookbackDays
	}
	if cfg.BarsPerDay == 0 {
		cfg.BarsPerDay = defaults.BarsPerDay
	}
	if cfg.Percentile == 0 {
		cfg.Percentile = defaults.Percentile
	}
	if cfg.MinConsecutiveBars == 0 {
		cfg.MinConsecutiveBars = defaults.MinConsecutiveBars
	}
	return cfg
}

func (cfg RangeDetectorConfig) validate() error {
	if cfg.ATRPeriod <= 0 {
		return fmt.Errorf("detector ATR period must be positive")
	}
	if cfg.DonchianPeriod <= 0 {
		return fmt.Errorf("detector Donchian period must be positive")
	}
	if cfg.BollingerPeriod <= 0 {
		return fmt.Errorf("detector Bollinger period must be positive")
	}
	if cfg.ADXPeriod <= 0 {
		return fmt.Errorf("detector ADX period must be positive")
	}
	if cfg.LookbackDays <= 0 {
		return fmt.Errorf("detector lookback days must be positive")
	}
	if cfg.BarsPerDay <= 0 {
		return fmt.Errorf("detector bars per day must be positive")
	}
	if cfg.Percentile <= 0 || cfg.Percentile >= 1 {
		return fmt.Errorf("detector percentile must be between 0 and 1")
	}
	if cfg.MinConsecutiveBars <= 0 {
		return fmt.Errorf("detector min consecutive bars must be positive")
	}
	return nil
}

func metricLow(value, threshold float64) bool {
	return validNumber(value) && validNumber(threshold) && value <= threshold
}

func markConsecutiveActive(raw []bool, minConsecutive int) []bool {
	active := make([]bool, len(raw))
	consecutive := 0
	for i, ok := range raw {
		if !ok {
			consecutive = 0
			continue
		}
		consecutive++
		active[i] = consecutive >= minConsecutive
	}
	return active
}

func rollingPriorPercentile(values []float64, lookback int, percentile float64) []float64 {
	out := nanSlice(len(values))
	if lookback <= 0 || percentile <= 0 || percentile >= 1 {
		return out
	}
	window := make([]float64, 0, lookback)
	for i, value := range values {
		if i >= lookback && len(window) > 0 {
			out[i] = percentileFromSorted(window, percentile)
		}
		if validNumber(value) {
			window = insertSorted(window, value)
		}
		if i >= lookback && validNumber(values[i-lookback]) {
			window = removeSorted(window, values[i-lookback])
		}
	}
	return out
}

func percentileFromSorted(sortedValues []float64, percentile float64) float64 {
	if len(sortedValues) == 0 {
		return math.NaN()
	}
	rank := int(math.Ceil(percentile*float64(len(sortedValues)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(sortedValues) {
		rank = len(sortedValues) - 1
	}
	return sortedValues[rank]
}

func insertSorted(values []float64, value float64) []float64 {
	idx := sort.SearchFloat64s(values, value)
	values = append(values, 0)
	copy(values[idx+1:], values[idx:])
	values[idx] = value
	return values
}

func removeSorted(values []float64, value float64) []float64 {
	idx := sort.SearchFloat64s(values, value)
	if idx >= len(values) || values[idx] != value {
		return values
	}
	copy(values[idx:], values[idx+1:])
	return values[:len(values)-1]
}
