package lab

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"
)

type rangeWorkbenchFacts struct {
	high  float64
	low   float64
	mid   float64
	width float64
	vwap  float64
}

func (cfg RangeOptimizationWorkbenchConfig) withDefaults() RangeOptimizationWorkbenchConfig {
	defaults := DefaultRangeOptimizationWorkbenchConfig()
	if cfg.ApprovedSourcePath == "" {
		cfg.ApprovedSourcePath = defaults.ApprovedSourcePath
	}
	if cfg.ExpectedSourceRows == 0 {
		cfg.ExpectedSourceRows = defaults.ExpectedSourceRows
	}
	if cfg.ExpectedFirstOpenTime == "" {
		cfg.ExpectedFirstOpenTime = defaults.ExpectedFirstOpenTime
	}
	if cfg.ExpectedLastOpenTime == "" {
		cfg.ExpectedLastOpenTime = defaults.ExpectedLastOpenTime
	}
	if cfg.ExpectedZeroVolumeCount == 0 {
		cfg.ExpectedZeroVolumeCount = defaults.ExpectedZeroVolumeCount
	}
	if cfg.Expected15MRows == 0 {
		cfg.Expected15MRows = defaults.Expected15MRows
	}
	if cfg.Expected15MLastOpenTime == "" {
		cfg.Expected15MLastOpenTime = defaults.Expected15MLastOpenTime
	}
	if cfg.MaxTrials == 0 {
		cfg.MaxTrials = defaults.MaxTrials
	}
	return cfg
}

func rangeWorkbenchSourceRow(manifest SourceManifest, cfg RangeOptimizationWorkbenchConfig, runID string) RangeOptimizationWorkbenchSourceRow {
	row := RangeOptimizationWorkbenchSourceRow{BacktestName: RangeOptimizationWorkbenchName, RunID: runID, Path: manifest.Path, ApprovedPath: cfg.ApprovedSourcePath, Product: manifest.Product, Symbol: manifest.Symbol, Interval: manifest.Interval, RowCount: manifest.RowCount, ExpectedRowCount: cfg.ExpectedSourceRows, FirstOpenTime: manifest.FirstOpenTime, ExpectedFirstOpenTime: cfg.ExpectedFirstOpenTime, LastOpenTime: manifest.LastOpenTime, ExpectedLastOpenTime: cfg.ExpectedLastOpenTime, GapCount: manifest.GapCount, DuplicateCount: manifest.DuplicateCount, ZeroVolumeCount: manifest.ZeroVolumeCount, ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount, ComparisonOnly: manifest.ComparisonOnly, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	failures := []string{}
	if manifest.ValidationStatus != "accepted" || manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly || manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" || filepath.Clean(manifest.Path) != filepath.Clean(cfg.ApprovedSourcePath) || manifest.RowCount != cfg.ExpectedSourceRows || manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime || manifest.LastOpenTime != cfg.ExpectedLastOpenTime || manifest.GapCount != 0 || manifest.DuplicateCount != 0 || manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
		failures = append(failures, "source contract mismatch")
	}
	row.SourceFactsPass = len(failures) == 0
	if !row.SourceFactsPass {
		row.ValidationStatus = "rejected"
		row.ValidationError = strings.Join(failures, "; ")
	}
	return row
}

func rangeWorkbenchResample15M(candles []Candle, cfg RangeOptimizationWorkbenchConfig, sourcePass bool, runID string) ([]Candle, RangeOptimizationWorkbenchCoverageRow) {
	out := []Candle{}
	missingBuckets := 0
	for i := 0; i+2 < len(candles); i += 3 {
		first, second, third := candles[i], candles[i+1], candles[i+2]
		start := first.OpenTime.UTC()
		if start.Minute()%15 != 0 || !second.OpenTime.UTC().Equal(start.Add(5*time.Minute)) || !third.OpenTime.UTC().Equal(start.Add(10*time.Minute)) {
			missingBuckets++
			continue
		}
		high := math.Max(first.High, math.Max(second.High, third.High))
		low := math.Min(first.Low, math.Min(second.Low, third.Low))
		out = append(out, Candle{OpenTime: start, CloseTime: third.CloseTime.UTC(), Open: first.Open, High: high, Low: low, Close: third.Close, Volume: first.Volume + second.Volume + third.Volume})
	}
	row := RangeOptimizationWorkbenchCoverageRow{BacktestName: RangeOptimizationWorkbenchName, RunID: runID, Timeframe: "15m", RowCount: len(out), ExpectedRowCount: cfg.Expected15MRows, ExpectedLastOpenTime: cfg.Expected15MLastOpenTime, ExpectedChildBars: 3, MissingChildBuckets: missingBuckets, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	if len(out) > 0 {
		row.FirstOpenTime = out[0].OpenTime.UTC().Format(timeLayout)
		row.LastOpenTime = out[len(out)-1].OpenTime.UTC().Format(timeLayout)
	}
	row.SourceResamplePass = sourcePass && missingBuckets == 0 && row.RowCount == cfg.Expected15MRows && row.LastOpenTime == cfg.Expected15MLastOpenTime
	if !row.SourceResamplePass {
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("15m resample mismatch rows=%d expected=%d missing_child_buckets=%d last_open=%s expected=%s", row.RowCount, cfg.Expected15MRows, missingBuckets, row.LastOpenTime, cfg.Expected15MLastOpenTime)
	}
	return out, row
}

func rangeWorkbenchRollingFacts(candles []Candle, d int, lookback int) (rangeWorkbenchFacts, bool) {
	if lookback <= 0 || d < lookback {
		return rangeWorkbenchFacts{}, false
	}
	high := candles[d-lookback].High
	low := candles[d-lookback].Low
	weighted := 0.0
	volume := 0.0
	for i := d - lookback; i <= d-1; i++ {
		c := candles[i]
		high = math.Max(high, c.High)
		low = math.Min(low, c.Low)
		typical := (c.High + c.Low + c.Close) / 3
		weighted += typical * c.Volume
		volume += c.Volume
	}
	width := high - low
	if width <= 0 || volume <= 0 {
		return rangeWorkbenchFacts{}, false
	}
	return rangeWorkbenchFacts{high: high, low: low, mid: (high + low) / 2, width: width, vwap: weighted / volume}, true
}

func rangeWorkbenchPreviousDayFacts(candles []Candle, d int) (rangeWorkbenchFacts, bool) {
	priorKey := candles[d].OpenTime.UTC().AddDate(0, 0, -1).Format("2006-01-02")
	count := 0
	high := 0.0
	low := 0.0
	weighted := 0.0
	volume := 0.0
	for i := d - 1; i >= 0; i-- {
		key := candles[i].OpenTime.UTC().Format("2006-01-02")
		if key != priorKey {
			if count > 0 {
				break
			}
			continue
		}
		c := candles[i]
		if count == 0 {
			high = c.High
			low = c.Low
		}
		high = math.Max(high, c.High)
		low = math.Min(low, c.Low)
		typical := (c.High + c.Low + c.Close) / 3
		weighted += typical * c.Volume
		volume += c.Volume
		count++
	}
	width := high - low
	if count < 96 || width <= 0 || volume <= 0 {
		return rangeWorkbenchFacts{}, false
	}
	return rangeWorkbenchFacts{high: high, low: low, mid: (high + low) / 2, width: width, vwap: weighted / volume}, true
}
