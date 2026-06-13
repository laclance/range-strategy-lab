package lab

import (
	"fmt"
	"math"
	"time"

	sr "github.com/laclance/go-sr"
)

const BalancedDetectorProfileID = "p30_c12_bollinger_on_adx_off"

type SRAuditConfig struct {
	Timeframe         string
	Mode              string
	LookbackBars      int
	MinStrength       int
	DetectorProfileID string
}

type SRAuditRow struct {
	Index                           int     `json:"index"`
	OpenTime                        string  `json:"open_time"`
	CloseTime                       string  `json:"close_time"`
	Split                           string  `json:"split"`
	Close                           float64 `json:"close"`
	Timeframe                       string  `json:"timeframe"`
	Mode                            string  `json:"mode"`
	LookbackBars                    int     `json:"lookback_bars"`
	WarmupBars                      int     `json:"warmup_bars"`
	MinStrength                     int     `json:"min_strength"`
	DetectorProfileID               string  `json:"detector_profile_id"`
	DetectorRawActive               bool    `json:"detector_raw_active"`
	DetectorActive                  bool    `json:"detector_active"`
	QualifiedZoneCount              int     `json:"qualified_zone_count"`
	RawZoneCount                    int     `json:"raw_zone_count"`
	HasSupport                      bool    `json:"has_support"`
	NearSupport                     bool    `json:"near_support"`
	NearestSupport                  float64 `json:"nearest_support"`
	NearestSupportDistance          float64 `json:"nearest_support_distance"`
	NearestSupportDistancePct       float64 `json:"nearest_support_distance_pct"`
	NearestSupportStrength          int     `json:"nearest_support_strength"`
	NearestSupportScore             float64 `json:"nearest_support_score"`
	NearestSupportTop               float64 `json:"nearest_support_top"`
	NearestSupportBottom            float64 `json:"nearest_support_bottom"`
	NearestSupportLastTouchIndex    int     `json:"nearest_support_last_touch_index"`
	NearestSupportSourcePivots      []int   `json:"nearest_support_source_pivots"`
	HasResistance                   bool    `json:"has_resistance"`
	NearResistance                  bool    `json:"near_resistance"`
	NearestResistance               float64 `json:"nearest_resistance"`
	NearestResistanceDistance       float64 `json:"nearest_resistance_distance"`
	NearestResistanceDistancePct    float64 `json:"nearest_resistance_distance_pct"`
	NearestResistanceStrength       int     `json:"nearest_resistance_strength"`
	NearestResistanceScore          float64 `json:"nearest_resistance_score"`
	NearestResistanceTop            float64 `json:"nearest_resistance_top"`
	NearestResistanceBottom         float64 `json:"nearest_resistance_bottom"`
	NearestResistanceLastTouchIndex int     `json:"nearest_resistance_last_touch_index"`
	NearestResistanceSourcePivots   []int   `json:"nearest_resistance_source_pivots"`
}

func DefaultSRAuditConfig() SRAuditConfig {
	return SRAuditConfig{
		Timeframe:         "5m",
		Mode:              string(sr.ModeZones),
		LookbackBars:      120,
		MinStrength:       2,
		DetectorProfileID: BalancedDetectorProfileID,
	}
}

func RunSRAudit(candles []Candle, cfg SRAuditConfig, splits []Split) ([]SRAuditRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	detector := CompressionRangeDetector{Config: DefaultCompressionRangeDetectorConfig()}
	classifications, err := detector.Classify(candles)
	if err != nil {
		return nil, err
	}

	mode := sr.Mode(cfg.Mode)
	warmupBars := sr.WarmupCandles(cfg.LookbackBars, mode)
	srCandles := toSRCandles(candles)
	rows := make([]SRAuditRow, 0, maxInt(0, len(candles)-warmupBars))

	for i := warmupBars; i < len(candles); i++ {
		levels, err := sr.Compute(srCandles[:i+1], sr.Options{
			Timeframe:   cfg.Timeframe,
			Lookback:    cfg.LookbackBars,
			Mode:        mode,
			MinStrength: cfg.MinStrength,
		})
		if err != nil {
			return nil, fmt.Errorf("compute SR audit at index %d: %w", i, err)
		}
		rows = append(rows, newSRAuditRow(candles[i], i, cfg, warmupBars, levels, classificationAt(classifications, i), splits))
	}
	return rows, nil
}

func (cfg SRAuditConfig) withDefaults() SRAuditConfig {
	defaults := DefaultSRAuditConfig()
	if cfg.Timeframe == "" {
		cfg.Timeframe = defaults.Timeframe
	}
	if cfg.Mode == "" {
		cfg.Mode = defaults.Mode
	}
	if cfg.LookbackBars == 0 {
		cfg.LookbackBars = defaults.LookbackBars
	}
	if cfg.MinStrength == 0 {
		cfg.MinStrength = defaults.MinStrength
	}
	if cfg.DetectorProfileID == "" {
		cfg.DetectorProfileID = defaults.DetectorProfileID
	}
	return cfg
}

func (cfg SRAuditConfig) validate() error {
	if cfg.Timeframe != "5m" {
		return fmt.Errorf("SR audit timeframe must be 5m")
	}
	if cfg.Mode != string(sr.ModeZones) {
		return fmt.Errorf("SR audit mode must be zone")
	}
	if cfg.LookbackBars <= 0 {
		return fmt.Errorf("SR audit lookback bars must be positive")
	}
	if cfg.MinStrength <= 0 {
		return fmt.Errorf("SR audit min strength must be positive")
	}
	return nil
}

func toSRCandles(candles []Candle) []sr.Candle {
	out := make([]sr.Candle, len(candles))
	for i, c := range candles {
		out[i] = sr.Candle{
			OpenTime:  c.OpenTime,
			CloseTime: c.CloseTime,
			Open:      c.Open,
			High:      c.High,
			Low:       c.Low,
			Close:     c.Close,
			Volume:    c.Volume,
		}
	}
	return out
}

func newSRAuditRow(c Candle, index int, cfg SRAuditConfig, warmupBars int, levels sr.Levels, classification RangeClassification, splits []Split) SRAuditRow {
	row := SRAuditRow{
		Index:                         index,
		OpenTime:                      c.OpenTime.Format(timeLayout),
		CloseTime:                     c.CloseTime.Format(timeLayout),
		Split:                         splitNameForCloseTime(c.CloseTime, splits),
		Close:                         c.Close,
		Timeframe:                     cfg.Timeframe,
		Mode:                          cfg.Mode,
		LookbackBars:                  cfg.LookbackBars,
		WarmupBars:                    warmupBars,
		MinStrength:                   cfg.MinStrength,
		DetectorProfileID:             cfg.DetectorProfileID,
		DetectorRawActive:             classification.RawActive,
		DetectorActive:                classification.Active,
		QualifiedZoneCount:            len(levels.Levels),
		RawZoneCount:                  len(levels.RawZones),
		NearestSupportSourcePivots:    []int{},
		NearestResistanceSourcePivots: []int{},
	}

	support, hasSupport, supportDistance, resistance, hasResistance, resistanceDistance := nearestSRZones(levels.Levels, c.Close)
	if hasSupport {
		row.HasSupport = true
		row.NearSupport = levels.NearSupport
		row.NearestSupport = support.Price
		row.NearestSupportDistance = supportDistance
		row.NearestSupportDistancePct = distancePct(supportDistance, c.Close)
		row.NearestSupportStrength = support.Strength
		row.NearestSupportScore = support.Score
		row.NearestSupportTop = support.Top
		row.NearestSupportBottom = support.Bottom
		row.NearestSupportLastTouchIndex = support.LastTouchIndex
		row.NearestSupportSourcePivots = append([]int(nil), support.SourcePivotIndexes...)
	}
	if hasResistance {
		row.HasResistance = true
		row.NearResistance = levels.NearResistance
		row.NearestResistance = resistance.Price
		row.NearestResistanceDistance = resistanceDistance
		row.NearestResistanceDistancePct = distancePct(resistanceDistance, c.Close)
		row.NearestResistanceStrength = resistance.Strength
		row.NearestResistanceScore = resistance.Score
		row.NearestResistanceTop = resistance.Top
		row.NearestResistanceBottom = resistance.Bottom
		row.NearestResistanceLastTouchIndex = resistance.LastTouchIndex
		row.NearestResistanceSourcePivots = append([]int(nil), resistance.SourcePivotIndexes...)
	}
	return row
}

func classificationAt(classifications []RangeClassification, index int) RangeClassification {
	if index < 0 || index >= len(classifications) {
		return RangeClassification{}
	}
	return classifications[index]
}

func nearestSRZones(levels []sr.Level, price float64) (sr.Level, bool, float64, sr.Level, bool, float64) {
	nearestSupportDistance := math.MaxFloat64
	nearestResistanceDistance := math.MaxFloat64
	var nearestSupport sr.Level
	var nearestResistance sr.Level
	hasSupport := false
	hasResistance := false

	for _, level := range levels {
		supportSide := price >= level.Price
		resistanceSide := price <= level.Price
		if supportSide && resistanceSide {
			if level.IsHigh {
				supportSide = false
			} else {
				resistanceSide = false
			}
		}

		if supportSide {
			distance := math.Max(0, price-level.Price)
			if distance < nearestSupportDistance {
				nearestSupportDistance = distance
				nearestSupport = level
				hasSupport = true
			}
		}
		if resistanceSide {
			distance := math.Max(0, level.Price-price)
			if distance < nearestResistanceDistance {
				nearestResistanceDistance = distance
				nearestResistance = level
				hasResistance = true
			}
		}
	}

	if !hasSupport {
		nearestSupportDistance = 0
	}
	if !hasResistance {
		nearestResistanceDistance = 0
	}
	return nearestSupport, hasSupport, nearestSupportDistance, nearestResistance, hasResistance, nearestResistanceDistance
}

func splitNameForCloseTime(closeTime time.Time, splits []Split) string {
	for _, split := range splits {
		if split.Name != "full_2021_2026" && split.Contains(closeTime) {
			return split.Name
		}
	}
	for _, split := range splits {
		if split.Contains(closeTime) {
			return split.Name
		}
	}
	return ""
}

func distancePct(distance, price float64) float64 {
	if price <= 0 {
		return 0
	}
	return distance / price
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
