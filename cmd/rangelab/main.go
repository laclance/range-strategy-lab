package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"range-strategy-lab/internal/lab"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	csvPath := flag.String("csv", "data/btcusdt_spot_5m_2021_2026.csv", "5m BTCUSDT candle CSV")
	outDir := flag.String("out-dir", "results/smoke", "output directory")
	startBalance := flag.Float64("start-balance", 1000, "starting balance")
	riskPct := flag.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := flag.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := flag.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := flag.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	maxHoldBars := flag.Int("max-hold-bars", 24, "default max hold bars stamped on placeholder signals")
	detector := flag.Bool("detector", false, "write detector-only range diagnostics")
	detectorSweep := flag.Bool("detector-sweep", false, "write detector sweep/audit diagnostics")
	detectorLookbackDays := flag.Int("detector-lookback-days", 20, "range detector trailing lookback in days")
	detectorPercentile := flag.Float64("detector-percentile", 0.30, "range detector low-compression percentile threshold")
	detectorMinConsecutiveBars := flag.Int("detector-min-consecutive-bars", 12, "range detector confirmed raw-active bars before active")
	detectorUseBollinger := flag.Bool("detector-use-bollinger", true, "include Bollinger20 width in range detector")
	detectorUseADX := flag.Bool("detector-use-adx", false, "include ADX14 in range detector")
	flag.Parse()

	candles, err := lab.LoadCSV(*csvPath)
	if err != nil {
		return err
	}
	if len(candles) == 0 {
		return fmt.Errorf("no candles loaded")
	}

	cfg := lab.BacktestConfig{
		StartBalance:   *startBalance,
		RiskPct:        *riskPct,
		MaxNotionalPct: *maxNotionalPct,
		FeePct:         *feePct,
		SlippagePct:    *slippagePct,
		MaxHoldBars:    *maxHoldBars,
	}

	strategy := lab.EmptyStrategy{}
	result := lab.RunBacktest(candles, strategy, cfg)
	summaries := lab.SummarizeSplits(result.Trades, *startBalance, lab.DefaultSplits())

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "trades.json"), result.Trades); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "summary.json"), summaries); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(*outDir, "summary.csv"), summaries); err != nil {
		return err
	}
	if *detector || *detectorSweep {
		if *detectorLookbackDays <= 0 {
			return fmt.Errorf("detector lookback days must be positive")
		}
	}
	if *detector {
		if *detectorPercentile <= 0 || *detectorPercentile >= 1 {
			return fmt.Errorf("detector percentile must be between 0 and 1")
		}
		if *detectorMinConsecutiveBars <= 0 {
			return fmt.Errorf("detector min consecutive bars must be positive")
		}

		cfg := lab.DefaultCompressionRangeDetectorConfig()
		cfg.LookbackDays = *detectorLookbackDays
		cfg.Percentile = *detectorPercentile
		cfg.MinConsecutiveBars = *detectorMinConsecutiveBars
		cfg.UseBollinger = *detectorUseBollinger
		cfg.UseADX = *detectorUseADX

		rangeDetector := lab.CompressionRangeDetector{Config: cfg}
		classifications, err := rangeDetector.Classify(candles)
		if err != nil {
			return err
		}
		dutyRows, episodes := lab.SummarizeDetectorSplits(candles, classifications, lab.DefaultSplits())
		if err := writeJSON(filepath.Join(*outDir, "detector_duty_cycle.json"), dutyRows); err != nil {
			return err
		}
		if err := writeDetectorDutyCycleCSV(filepath.Join(*outDir, "detector_duty_cycle.csv"), dutyRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_episodes.json"), episodes); err != nil {
			return err
		}
		if err := writeRangeEpisodesCSV(filepath.Join(*outDir, "range_episodes.csv"), episodes); err != nil {
			return err
		}
		for _, row := range dutyRows {
			if row.Split == "full_2021_2026" {
				fmt.Printf("detector=%s active_bars=%d total_bars=%d duty_cycle=%.4f episodes=%d\n",
					rangeDetector.Name(), row.ActiveBars, row.TotalBars, row.DutyCycle, row.Episodes)
				break
			}
		}
	}
	if *detectorSweep {
		cfg := lab.DefaultCompressionRangeDetectorConfig()
		cfg.LookbackDays = *detectorLookbackDays

		sweepRows, err := lab.RunDetectorSweep(candles, cfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_sweep.json"), sweepRows); err != nil {
			return err
		}
		if err := writeDetectorSweepCSV(filepath.Join(*outDir, "detector_sweep.csv"), sweepRows); err != nil {
			return err
		}
		for _, row := range sweepRows {
			if row.IsBalancedBaseline && row.Split == "full_2021_2026" {
				fmt.Printf("detector_sweep profiles=%d rows=%d baseline_active_bars=%d baseline_total_bars=%d baseline_duty_cycle=%.4f baseline_episodes=%d\n",
					len(lab.DefaultDetectorSweepProfiles(*detectorLookbackDays)),
					len(sweepRows),
					row.ActiveBars,
					row.TotalBars,
					row.DutyCycle,
					row.Episodes,
				)
				break
			}
		}
	}

	first := candles[0].OpenTime.Format(time.RFC3339)
	last := candles[len(candles)-1].CloseTime.Format(time.RFC3339)
	fmt.Printf("loaded %d candles from %s to %s\n", len(candles), first, last)
	fmt.Printf("strategy=%s trades=%d output=%s\n", strategy.Name(), len(result.Trades), *outDir)
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func writeSummaryCSV(path string, rows []lab.SummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"side",
		"total_trades",
		"wins",
		"losses",
		"win_rate",
		"gross_pnl",
		"net_pnl",
		"total_costs",
		"profit_factor",
		"gross_profit_factor",
		"max_drawdown",
		"expectancy",
		"avg_hold_bars",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.TotalTrades),
			strconv.Itoa(row.Wins),
			strconv.Itoa(row.Losses),
			formatFloat(row.WinRate),
			formatFloat(row.GrossPnL),
			formatFloat(row.NetPnL),
			formatFloat(row.TotalCosts),
			formatFloat(row.ProfitFactor),
			formatFloat(row.GrossProfitFactor),
			formatFloat(row.MaxDrawdown),
			formatFloat(row.Expectancy),
			formatFloat(row.AvgHoldBars),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDutyCycleCSV(path string, rows []lab.DetectorDutyCycleRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"episodes",
		"avg_episode_length",
		"median_episode_length",
		"longest_episode_length",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.Episodes),
			formatFloat(row.AvgEpisodeLength),
			formatFloat(row.MedianEpisodeLength),
			strconv.Itoa(row.LongestEpisodeLength),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeEpisodesCSV(path string, episodes []lab.RangeEpisode) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"start_index",
		"end_index",
		"start_time",
		"end_time",
		"length_bars",
	}); err != nil {
		return err
	}
	for _, episode := range episodes {
		if err := w.Write([]string{
			episode.Split,
			strconv.Itoa(episode.StartIndex),
			strconv.Itoa(episode.EndIndex),
			episode.StartTime,
			episode.EndTime,
			strconv.Itoa(episode.LengthBars),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorSweepCSV(path string, rows []lab.DetectorSweepRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"profile_id",
		"is_balanced_baseline",
		"is_adx_comparison",
		"percentile",
		"min_consecutive_bars",
		"use_bollinger",
		"use_adx",
		"lookback_days",
		"split",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"episodes",
		"avg_episode_length",
		"median_episode_length",
		"longest_episode_length",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.ProfileID,
			strconv.FormatBool(row.IsBalancedBaseline),
			strconv.FormatBool(row.IsADXComparison),
			formatFloat(row.Percentile),
			strconv.Itoa(row.MinConsecutiveBars),
			strconv.FormatBool(row.UseBollinger),
			strconv.FormatBool(row.UseADX),
			strconv.Itoa(row.LookbackDays),
			row.Split,
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.Episodes),
			formatFloat(row.AvgEpisodeLength),
			formatFloat(row.MedianEpisodeLength),
			strconv.Itoa(row.LongestEpisodeLength),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
