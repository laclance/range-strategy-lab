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
	detectorDurabilitySweep := flag.Bool("detector-durability-sweep", false, "write non-trading detector durability sweep diagnostics")
	srAudit := flag.Bool("sr-audit", false, "write go-sr support/resistance audit diagnostics")
	srBoundaryAudit := flag.Bool("sr-boundary-audit", false, "write non-trading SR boundary quality diagnostics")
	srBoundaryInspect := flag.Bool("sr-boundary-inspect", false, "write compact non-trading SR boundary candidate comparison diagnostics")
	srRejectionTimingAudit := flag.Bool("sr-rejection-timing-audit", false, "write compact non-trading SR rejection timing diagnostics")
	srConfirmationTimingAudit := flag.Bool("sr-confirmation-timing-audit", false, "write compact non-trading SR confirmation timing diagnostics")
	srFalseBreakReclaimTimingAudit := flag.Bool("sr-false-break-reclaim-timing-audit", false, "write compact non-trading SR false-break reclaim timing diagnostics")
	compressionBreakoutAudit := flag.Bool("compression-breakout-audit", false, "write compact non-trading compression breakout diagnostics")
	rangeRegimeDurabilityAudit := flag.Bool("range-regime-durability-audit", false, "write compact non-trading range regime durability diagnostics")
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
	var srRows []lab.SRAuditRow
	srCfg := lab.DefaultSRAuditConfig()
	if *srAudit || *srBoundaryAudit || *srBoundaryInspect || *srRejectionTimingAudit || *srConfirmationTimingAudit || *srFalseBreakReclaimTimingAudit {
		var err error
		srRows, err = lab.RunSRAudit(candles, srCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
	}
	if *srAudit {
		if err := writeJSON(filepath.Join(*outDir, "sr_touch_audit.json"), srRows); err != nil {
			return err
		}
		if err := writeSRTouchAuditCSV(filepath.Join(*outDir, "sr_touch_audit.csv"), srRows); err != nil {
			return err
		}
		nearSupportRows := 0
		nearResistanceRows := 0
		for _, row := range srRows {
			if row.NearSupport {
				nearSupportRows++
			}
			if row.NearResistance {
				nearResistanceRows++
			}
		}
		warmupBars := 0
		if len(srRows) > 0 {
			warmupBars = srRows[0].WarmupBars
		}
		fmt.Printf("sr_audit rows=%d lookback_bars=%d warmup_bars=%d near_support_rows=%d near_resistance_rows=%d\n",
			len(srRows),
			srCfg.LookbackBars,
			warmupBars,
			nearSupportRows,
			nearResistanceRows,
		)
	}
	if *srBoundaryAudit || *srBoundaryInspect {
		boundaryCfg := lab.DefaultSRBoundaryAuditConfig()
		events, qualityRows, err := lab.RunSRBoundaryAudit(candles, srRows, boundaryCfg)
		if err != nil {
			return err
		}
		if *srBoundaryAudit {
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_events.json"), events); err != nil {
				return err
			}
			if err := writeSRBoundaryEventsCSV(filepath.Join(*outDir, "sr_boundary_events.csv"), events); err != nil {
				return err
			}
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_quality.json"), qualityRows); err != nil {
				return err
			}
			if err := writeSRBoundaryQualityCSV(filepath.Join(*outDir, "sr_boundary_quality.csv"), qualityRows); err != nil {
				return err
			}
			fmt.Printf("sr_boundary_audit events=%d summary_rows=%d horizons=%s detector_active_only=%t\n",
				len(events),
				len(qualityRows),
				formatIntSlice(boundaryCfg.HorizonsBars),
				boundaryCfg.DetectorActiveOnly,
			)
		}
		if *srBoundaryInspect {
			comparisonRows := lab.SummarizeSRBoundaryCandidateComparison(events)
			if err := writeJSON(filepath.Join(*outDir, "sr_boundary_candidate_comparison.json"), comparisonRows); err != nil {
				return err
			}
			if err := writeSRBoundaryCandidateComparisonCSV(filepath.Join(*outDir, "sr_boundary_candidate_comparison.csv"), comparisonRows); err != nil {
				return err
			}
			fmt.Printf("sr_boundary_inspect events=%d comparison_rows=%d horizons=%s detector_active_only=%t\n",
				len(events),
				len(comparisonRows),
				formatIntSlice(boundaryCfg.HorizonsBars),
				boundaryCfg.DetectorActiveOnly,
			)
		}
	}
	if *srRejectionTimingAudit {
		timingCfg := lab.DefaultSRRejectionTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRRejectionTimingAudit(candles, srRows, timingCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_rejection_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRRejectionTimingCandidatesCSV(filepath.Join(*outDir, "sr_rejection_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_rejection_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRRejectionTimingSummaryCSV(filepath.Join(*outDir, "sr_rejection_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_rejection_timing_audit candidate_rows=%d summary_rows=%d horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			formatIntSlice(timingCfg.HorizonsBars),
			timingCfg.DetectorActiveOnly,
		)
	}
	if *srConfirmationTimingAudit {
		confirmationCfg := lab.DefaultSRConfirmationTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRConfirmationTimingAudit(candles, srRows, confirmationCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_confirmation_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRConfirmationTimingCandidatesCSV(filepath.Join(*outDir, "sr_confirmation_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_confirmation_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRConfirmationTimingSummaryCSV(filepath.Join(*outDir, "sr_confirmation_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_confirmation_timing_audit candidate_rows=%d summary_rows=%d delays=%s horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			formatIntSlice(confirmationCfg.ConfirmationDelayBars),
			formatIntSlice(confirmationCfg.HorizonsBars),
			confirmationCfg.DetectorActiveOnly,
		)
	}
	if *srFalseBreakReclaimTimingAudit {
		falseBreakCfg := lab.DefaultSRFalseBreakReclaimTimingAuditConfig()
		candidateRows, summaryRows, err := lab.RunSRFalseBreakReclaimTimingAudit(candles, srRows, falseBreakCfg)
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_false_break_reclaim_timing_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeSRFalseBreakReclaimTimingCandidatesCSV(filepath.Join(*outDir, "sr_false_break_reclaim_timing_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "sr_false_break_reclaim_timing_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeSRFalseBreakReclaimTimingSummaryCSV(filepath.Join(*outDir, "sr_false_break_reclaim_timing_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("sr_false_break_reclaim_timing_audit candidate_rows=%d summary_rows=%d max_break_delay=%d max_reclaim_delay=%d horizons=%s detector_active_only=%t\n",
			len(candidateRows),
			len(summaryRows),
			falseBreakCfg.MaxBreakDelayBars,
			falseBreakCfg.MaxReclaimDelayBars,
			formatIntSlice(falseBreakCfg.HorizonsBars),
			falseBreakCfg.DetectorActiveOnly,
		)
	}
	if *compressionBreakoutAudit {
		breakoutCfg := lab.DefaultCompressionBreakoutAuditConfig()
		candidateRows, summaryRows, err := lab.RunCompressionBreakoutAudit(candles, breakoutCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "compression_breakout_candidates.json"), candidateRows); err != nil {
			return err
		}
		if err := writeCompressionBreakoutCandidatesCSV(filepath.Join(*outDir, "compression_breakout_candidates.csv"), candidateRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "compression_breakout_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeCompressionBreakoutSummaryCSV(filepath.Join(*outDir, "compression_breakout_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("compression_breakout_audit candidate_rows=%d summary_rows=%d max_breakout_delay=%d horizons=%s detector_profile_id=%s\n",
			len(candidateRows),
			len(summaryRows),
			breakoutCfg.MaxBreakoutDelayBars,
			formatIntSlice(breakoutCfg.HorizonsBars),
			breakoutCfg.DetectorProfileID,
		)
	}
	if *rangeRegimeDurabilityAudit {
		durabilityCfg := lab.DefaultRangeRegimeDurabilityAuditConfig()
		episodeRows, summaryRows, err := lab.RunRangeRegimeDurabilityAudit(candles, durabilityCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_regime_durability_episodes.json"), episodeRows); err != nil {
			return err
		}
		if err := writeRangeRegimeDurabilityEpisodesCSV(filepath.Join(*outDir, "range_regime_durability_episodes.csv"), episodeRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "range_regime_durability_summary.json"), summaryRows); err != nil {
			return err
		}
		if err := writeRangeRegimeDurabilitySummaryCSV(filepath.Join(*outDir, "range_regime_durability_summary.csv"), summaryRows); err != nil {
			return err
		}
		fmt.Printf("range_regime_durability_audit episode_rows=%d summary_rows=%d quick_invalidation_bars=%d horizons=%s detector_profile_id=%s\n",
			len(episodeRows),
			len(summaryRows),
			durabilityCfg.QuickInvalidationBars,
			formatIntSlice(durabilityCfg.HorizonsBars),
			durabilityCfg.DetectorProfileID,
		)
	}
	if *detector || *detectorSweep || *detectorDurabilitySweep {
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
	if *detectorDurabilitySweep {
		detectorCfg := lab.DefaultCompressionRangeDetectorConfig()
		detectorCfg.LookbackDays = *detectorLookbackDays
		durabilityCfg := lab.DefaultRangeRegimeDurabilityAuditConfig()

		broadRows, sliceRows, stabilityRows, err := lab.RunDetectorDurabilitySweep(candles, detectorCfg, durabilityCfg, lab.DefaultSplits())
		if err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_sweep.json"), broadRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilitySweepCSV(filepath.Join(*outDir, "detector_durability_sweep.csv"), broadRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_slices.json"), sliceRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilitySlicesCSV(filepath.Join(*outDir, "detector_durability_slices.csv"), sliceRows); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(*outDir, "detector_durability_stability.json"), stabilityRows); err != nil {
			return err
		}
		if err := writeDetectorDurabilityStabilityCSV(filepath.Join(*outDir, "detector_durability_stability.csv"), stabilityRows); err != nil {
			return err
		}
		fmt.Printf("detector_durability_sweep profiles=%d broad_rows=%d slice_rows=%d stability_rows=%d quick_invalidation_bars=%d horizons=%s\n",
			len(lab.DefaultDetectorSweepProfiles(*detectorLookbackDays)),
			len(broadRows),
			len(sliceRows),
			len(stabilityRows),
			durabilityCfg.QuickInvalidationBars,
			formatIntSlice(durabilityCfg.HorizonsBars),
		)
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

func writeDetectorDurabilitySweepCSV(path string, rows []lab.DetectorDurabilitySweepRow) error {
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
		"horizon_bars",
		"active_bars",
		"total_bars",
		"duty_cycle",
		"detector_episodes",
		"avg_detector_episode_length",
		"median_detector_episode_length",
		"longest_detector_episode_length",
		"durability_episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
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
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.ActiveBars),
			strconv.Itoa(row.TotalBars),
			formatFloat(row.DutyCycle),
			strconv.Itoa(row.DetectorEpisodes),
			formatFloat(row.AvgDetectorEpisodeLength),
			formatFloat(row.MedianDetectorEpisodeLength),
			strconv.Itoa(row.LongestDetectorEpisodeLength),
			strconv.Itoa(row.DurabilityEpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDurabilitySlicesCSV(path string, rows []lab.DetectorDurabilitySliceRow) error {
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
		"horizon_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_width_bucket",
		"width_to_atr_bucket",
		"episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
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
			strconv.Itoa(row.HorizonBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			row.EpisodeWidthBucket,
			row.WidthToATRBucket,
			strconv.Itoa(row.EpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeDetectorDurabilityStabilityCSV(path string, rows []lab.DetectorDurabilityStabilityRow) error {
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
		"horizon_bars",
		"period_splits",
		"period_episode_count",
		"episode_count_min",
		"episode_count_max",
		"episode_count_delta",
		"duty_cycle_min",
		"duty_cycle_max",
		"duty_cycle_delta",
		"label_reentered_range_rate_min",
		"label_reentered_range_rate_max",
		"label_reentered_range_rate_delta",
		"label_persisted_inside_range_rate_min",
		"label_persisted_inside_range_rate_max",
		"label_persisted_inside_range_rate_delta",
		"label_quick_invalidated_rate_min",
		"label_quick_invalidated_rate_max",
		"label_quick_invalidated_rate_delta",
		"label_chopped_rate_min",
		"label_chopped_rate_max",
		"label_chopped_rate_delta",
		"label_trended_rate_min",
		"label_trended_rate_max",
		"label_trended_rate_delta",
		"label_avg_close_drift_pct_min",
		"label_avg_close_drift_pct_max",
		"label_avg_close_drift_pct_delta",
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
			strconv.Itoa(row.HorizonBars),
			strconv.Itoa(row.PeriodSplits),
			strconv.Itoa(row.PeriodEpisodeCount),
			strconv.Itoa(row.EpisodeCountMin),
			strconv.Itoa(row.EpisodeCountMax),
			strconv.Itoa(row.EpisodeCountDelta),
			formatFloat(row.DutyCycleMin),
			formatFloat(row.DutyCycleMax),
			formatFloat(row.DutyCycleDelta),
			formatFloat(row.LabelReenteredRangeRateMin),
			formatFloat(row.LabelReenteredRangeRateMax),
			formatFloat(row.LabelReenteredRangeRateDelta),
			formatFloat(row.LabelPersistedInsideRangeRateMin),
			formatFloat(row.LabelPersistedInsideRangeRateMax),
			formatFloat(row.LabelPersistedInsideRangeRateDelta),
			formatFloat(row.LabelQuickInvalidatedRateMin),
			formatFloat(row.LabelQuickInvalidatedRateMax),
			formatFloat(row.LabelQuickInvalidatedRateDelta),
			formatFloat(row.LabelChoppedRateMin),
			formatFloat(row.LabelChoppedRateMax),
			formatFloat(row.LabelChoppedRateDelta),
			formatFloat(row.LabelTrendedRateMin),
			formatFloat(row.LabelTrendedRateMax),
			formatFloat(row.LabelTrendedRateDelta),
			formatFloat(row.LabelAvgCloseDriftPctMin),
			formatFloat(row.LabelAvgCloseDriftPctMax),
			formatFloat(row.LabelAvgCloseDriftPctDelta),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeCompressionBreakoutCandidatesCSV(path string, rows []lab.CompressionBreakoutCandidateRow) error {
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
		"breakout_delay_bars",
		"horizon_bars",
		"episode_raw_length_bucket",
		"episode_active_length_bucket",
		"episode_range_width_bucket",
		"breakout_move_bucket",
		"decision_true_range_expansion_bucket",
		"detector_profile_id",
		"candidate_count",
		"avg_episode_raw_length_bars",
		"avg_episode_active_length_bars",
		"avg_episode_range_width_pct",
		"avg_breakout_move_pct",
		"avg_decision_true_range_atr",
		"label_reentered_range_count",
		"label_opposite_close_break_count",
		"label_favorable_greater_than_adverse_count",
		"label_reentered_range_rate",
		"label_opposite_close_break_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.BreakoutDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.EpisodeRawLengthBucket,
			row.EpisodeActiveLengthBucket,
			row.EpisodeRangeWidthBucket,
			row.BreakoutMoveBucket,
			row.DecisionTrueRangeExpansionBucket,
			row.DetectorProfileID,
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgEpisodeRawLengthBars),
			formatFloat(row.AvgEpisodeActiveLengthBars),
			formatFloat(row.AvgEpisodeRangeWidthPct),
			formatFloat(row.AvgBreakoutMovePct),
			formatFloat(row.AvgDecisionTrueRangeATR),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelOppositeCloseBreakCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelOppositeCloseBreakRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeCompressionBreakoutSummaryCSV(path string, rows []lab.CompressionBreakoutSummaryRow) error {
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
		"horizon_bars",
		"detector_profile_id",
		"candidate_count",
		"avg_breakout_delay_bars",
		"avg_episode_raw_length_bars",
		"avg_episode_active_length_bars",
		"avg_episode_range_width_pct",
		"avg_breakout_move_pct",
		"avg_decision_true_range_atr",
		"label_reentered_range_rate",
		"label_opposite_close_break_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgBreakoutDelayBars),
			formatFloat(row.AvgEpisodeRawLengthBars),
			formatFloat(row.AvgEpisodeActiveLengthBars),
			formatFloat(row.AvgEpisodeRangeWidthPct),
			formatFloat(row.AvgBreakoutMovePct),
			formatFloat(row.AvgDecisionTrueRangeATR),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelOppositeCloseBreakRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeRegimeDurabilityEpisodesCSV(path string, rows []lab.RangeRegimeDurabilityEpisodeRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"episode_id",
		"start_index",
		"end_index",
		"start_time",
		"end_time",
		"horizon_bars",
		"detector_profile_id",
		"raw_length_bars",
		"active_length_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_high",
		"episode_low",
		"episode_end_close",
		"episode_width_pct",
		"episode_width_bucket",
		"avg_normalized_atr",
		"end_normalized_atr",
		"width_to_atr_ratio",
		"width_to_atr_bucket",
		"label_window_start_index",
		"label_window_end_index",
		"label_window_start_time",
		"label_window_end_time",
		"label_reentered_range",
		"label_persisted_inside_range",
		"label_quick_invalidated",
		"label_invalidated_up",
		"label_invalidated_down",
		"label_chopped",
		"label_trended_up",
		"label_trended_down",
		"label_close_drift_pct",
		"label_max_up_move_pct",
		"label_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.EpisodeID),
			strconv.Itoa(row.StartIndex),
			strconv.Itoa(row.EndIndex),
			row.StartTime,
			row.EndTime,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.Itoa(row.RawLengthBars),
			strconv.Itoa(row.ActiveLengthBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			formatFloat(row.EpisodeHigh),
			formatFloat(row.EpisodeLow),
			formatFloat(row.EpisodeEndClose),
			formatFloat(row.EpisodeWidthPct),
			row.EpisodeWidthBucket,
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.EndNormalizedATR),
			formatFloat(row.WidthToATRRatio),
			row.WidthToATRBucket,
			strconv.Itoa(row.LabelWindowStartIndex),
			strconv.Itoa(row.LabelWindowEndIndex),
			row.LabelWindowStartTime,
			row.LabelWindowEndTime,
			strconv.FormatBool(row.LabelReenteredRange),
			strconv.FormatBool(row.LabelPersistedInsideRange),
			strconv.FormatBool(row.LabelQuickInvalidated),
			strconv.FormatBool(row.LabelInvalidatedUp),
			strconv.FormatBool(row.LabelInvalidatedDown),
			strconv.FormatBool(row.LabelChopped),
			strconv.FormatBool(row.LabelTrendedUp),
			strconv.FormatBool(row.LabelTrendedDown),
			formatFloat(row.LabelCloseDriftPct),
			formatFloat(row.LabelMaxUpMovePct),
			formatFloat(row.LabelMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeRangeRegimeDurabilitySummaryCSV(path string, rows []lab.RangeRegimeDurabilitySummaryRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"split",
		"horizon_bars",
		"raw_length_bucket",
		"active_length_bucket",
		"episode_width_bucket",
		"width_to_atr_bucket",
		"detector_profile_id",
		"episode_count",
		"avg_raw_length_bars",
		"avg_active_length_bars",
		"avg_episode_width_pct",
		"avg_normalized_atr",
		"avg_end_normalized_atr",
		"avg_width_to_atr_ratio",
		"label_reentered_range_count",
		"label_persisted_inside_range_count",
		"label_quick_invalidated_count",
		"label_invalidated_up_count",
		"label_invalidated_down_count",
		"label_chopped_count",
		"label_trended_up_count",
		"label_trended_down_count",
		"label_reentered_range_rate",
		"label_persisted_inside_range_rate",
		"label_quick_invalidated_rate",
		"label_invalidated_up_rate",
		"label_invalidated_down_rate",
		"label_chopped_rate",
		"label_trended_up_rate",
		"label_trended_down_rate",
		"label_avg_close_drift_pct",
		"label_avg_max_up_move_pct",
		"label_avg_max_down_move_pct",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			strconv.Itoa(row.HorizonBars),
			row.RawLengthBucket,
			row.ActiveLengthBucket,
			row.EpisodeWidthBucket,
			row.WidthToATRBucket,
			row.DetectorProfileID,
			strconv.Itoa(row.EpisodeCount),
			formatFloat(row.AvgRawLengthBars),
			formatFloat(row.AvgActiveLengthBars),
			formatFloat(row.AvgEpisodeWidthPct),
			formatFloat(row.AvgNormalizedATR),
			formatFloat(row.AvgEndNormalizedATR),
			formatFloat(row.AvgWidthToATRRatio),
			strconv.Itoa(row.LabelReenteredRangeCount),
			strconv.Itoa(row.LabelPersistedInsideRangeCount),
			strconv.Itoa(row.LabelQuickInvalidatedCount),
			strconv.Itoa(row.LabelInvalidatedUpCount),
			strconv.Itoa(row.LabelInvalidatedDownCount),
			strconv.Itoa(row.LabelChoppedCount),
			strconv.Itoa(row.LabelTrendedUpCount),
			strconv.Itoa(row.LabelTrendedDownCount),
			formatFloat(row.LabelReenteredRangeRate),
			formatFloat(row.LabelPersistedInsideRangeRate),
			formatFloat(row.LabelQuickInvalidatedRate),
			formatFloat(row.LabelInvalidatedUpRate),
			formatFloat(row.LabelInvalidatedDownRate),
			formatFloat(row.LabelChoppedRate),
			formatFloat(row.LabelTrendedUpRate),
			formatFloat(row.LabelTrendedDownRate),
			formatFloat(row.LabelAvgCloseDriftPct),
			formatFloat(row.LabelAvgMaxUpMovePct),
			formatFloat(row.LabelAvgMaxDownMovePct),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRTouchAuditCSV(path string, rows []lab.SRAuditRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"index",
		"open_time",
		"close_time",
		"split",
		"close",
		"timeframe",
		"mode",
		"lookback_bars",
		"warmup_bars",
		"min_strength",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"qualified_zone_count",
		"raw_zone_count",
		"has_support",
		"near_support",
		"nearest_support",
		"nearest_support_distance",
		"nearest_support_distance_pct",
		"nearest_support_strength",
		"nearest_support_score",
		"nearest_support_top",
		"nearest_support_bottom",
		"nearest_support_last_touch_index",
		"nearest_support_source_pivots",
		"has_resistance",
		"near_resistance",
		"nearest_resistance",
		"nearest_resistance_distance",
		"nearest_resistance_distance_pct",
		"nearest_resistance_strength",
		"nearest_resistance_score",
		"nearest_resistance_top",
		"nearest_resistance_bottom",
		"nearest_resistance_last_touch_index",
		"nearest_resistance_source_pivots",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			strconv.Itoa(row.Index),
			row.OpenTime,
			row.CloseTime,
			row.Split,
			formatFloat(row.Close),
			row.Timeframe,
			row.Mode,
			strconv.Itoa(row.LookbackBars),
			strconv.Itoa(row.WarmupBars),
			strconv.Itoa(row.MinStrength),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.QualifiedZoneCount),
			strconv.Itoa(row.RawZoneCount),
			strconv.FormatBool(row.HasSupport),
			strconv.FormatBool(row.NearSupport),
			formatFloat(row.NearestSupport),
			formatFloat(row.NearestSupportDistance),
			formatFloat(row.NearestSupportDistancePct),
			strconv.Itoa(row.NearestSupportStrength),
			formatFloat(row.NearestSupportScore),
			formatFloat(row.NearestSupportTop),
			formatFloat(row.NearestSupportBottom),
			strconv.Itoa(row.NearestSupportLastTouchIndex),
			formatIntSlice(row.NearestSupportSourcePivots),
			strconv.FormatBool(row.HasResistance),
			strconv.FormatBool(row.NearResistance),
			formatFloat(row.NearestResistance),
			formatFloat(row.NearestResistanceDistance),
			formatFloat(row.NearestResistanceDistancePct),
			strconv.Itoa(row.NearestResistanceStrength),
			formatFloat(row.NearestResistanceScore),
			formatFloat(row.NearestResistanceTop),
			formatFloat(row.NearestResistanceBottom),
			strconv.Itoa(row.NearestResistanceLastTouchIndex),
			formatIntSlice(row.NearestResistanceSourcePivots),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryEventsCSV(path string, rows []lab.SRBoundaryEventRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"index",
		"open_time",
		"close_time",
		"split",
		"side",
		"close",
		"boundary_price",
		"zone_top",
		"zone_bottom",
		"zone_width",
		"rejection_threshold",
		"distance_pct",
		"strength",
		"strength_bucket",
		"score",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"horizon_bars",
		"future_max_high",
		"future_min_low",
		"future_close",
		"favorable_move",
		"adverse_move",
		"favorable_move_pct",
		"adverse_move_pct",
		"distance_bucket",
		"wick_break",
		"close_break",
		"reclaimed_after_break",
		"rejected",
		"favorable_greater_than_adverse",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			strconv.Itoa(row.Index),
			row.OpenTime,
			row.CloseTime,
			row.Split,
			row.Side,
			formatFloat(row.Close),
			formatFloat(row.BoundaryPrice),
			formatFloat(row.ZoneTop),
			formatFloat(row.ZoneBottom),
			formatFloat(row.ZoneWidth),
			formatFloat(row.RejectionThreshold),
			formatFloat(row.DistancePct),
			strconv.Itoa(row.Strength),
			row.StrengthBucket,
			formatFloat(row.Score),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.HorizonBars),
			formatFloat(row.FutureMaxHigh),
			formatFloat(row.FutureMinLow),
			formatFloat(row.FutureClose),
			formatFloat(row.FavorableMove),
			formatFloat(row.AdverseMove),
			formatFloat(row.FavorableMovePct),
			formatFloat(row.AdverseMovePct),
			row.DistanceBucket,
			strconv.FormatBool(row.WickBreak),
			strconv.FormatBool(row.CloseBreak),
			strconv.FormatBool(row.ReclaimedAfterBreak),
			strconv.FormatBool(row.Rejected),
			strconv.FormatBool(row.FavorableGreaterThanAdverse),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryQualityCSV(path string, rows []lab.SRBoundaryQualityRow) error {
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
		"horizon_bars",
		"strength_bucket",
		"distance_bucket",
		"event_count",
		"avg_score",
		"avg_distance_pct",
		"avg_favorable_pct",
		"median_favorable_pct",
		"avg_adverse_pct",
		"median_adverse_pct",
		"close_break_rate",
		"wick_break_rate",
		"reclaim_after_break_rate",
		"rejection_rate",
		"favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.StrengthBucket,
			row.DistanceBucket,
			strconv.Itoa(row.EventCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgFavorablePct),
			formatFloat(row.MedianFavorablePct),
			formatFloat(row.AvgAdversePct),
			formatFloat(row.MedianAdversePct),
			formatFloat(row.CloseBreakRate),
			formatFloat(row.WickBreakRate),
			formatFloat(row.ReclaimAfterBreakRate),
			formatFloat(row.RejectionRate),
			formatFloat(row.FavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRBoundaryCandidateComparisonCSV(path string, rows []lab.SRBoundaryCandidateComparisonRow) error {
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
		"horizon_bars",
		"strength_bucket",
		"distance_bucket",
		"event_count",
		"close_break_count",
		"rejected_count",
		"reclaimed_after_break_count",
		"close_break_rate",
		"rejection_rate",
		"reclaim_event_rate",
		"reclaim_given_close_break_rate",
		"all_avg_favorable_pct",
		"all_avg_adverse_pct",
		"all_favorable_minus_adverse_pct",
		"all_favorable_greater_than_adverse_rate",
		"rejected_avg_favorable_pct",
		"rejected_avg_adverse_pct",
		"rejected_favorable_minus_adverse_pct",
		"rejected_favorable_greater_than_adverse_rate",
		"reclaimed_avg_favorable_pct",
		"reclaimed_avg_adverse_pct",
		"reclaimed_favorable_minus_adverse_pct",
		"reclaimed_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.StrengthBucket,
			row.DistanceBucket,
			strconv.Itoa(row.EventCount),
			strconv.Itoa(row.CloseBreakCount),
			strconv.Itoa(row.RejectedCount),
			strconv.Itoa(row.ReclaimedAfterBreakCount),
			formatFloat(row.CloseBreakRate),
			formatFloat(row.RejectionRate),
			formatFloat(row.ReclaimEventRate),
			formatFloat(row.ReclaimGivenCloseBreakRate),
			formatFloat(row.AllAvgFavorablePct),
			formatFloat(row.AllAvgAdversePct),
			formatFloat(row.AllFavorableMinusAdversePct),
			formatFloat(row.AllFavorableGreaterThanAdverseRate),
			formatFloat(row.RejectedAvgFavorablePct),
			formatFloat(row.RejectedAvgAdversePct),
			formatFloat(row.RejectedFavorableMinusAdversePct),
			formatFloat(row.RejectedFavorableGreaterThanAdverseRate),
			formatFloat(row.ReclaimedAvgFavorablePct),
			formatFloat(row.ReclaimedAvgAdversePct),
			formatFloat(row.ReclaimedFavorableMinusAdversePct),
			formatFloat(row.ReclaimedFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRRejectionTimingCandidatesCSV(path string, rows []lab.SRRejectionTimingCandidateRow) error {
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
		"horizon_bars",
		"close_location",
		"touched_zone",
		"pierced_zone",
		"closed_back",
		"decision_rejection_candidate",
		"wick_beyond_bucket",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_wick_beyond_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.CloseLocation,
			strconv.FormatBool(row.TouchedZone),
			strconv.FormatBool(row.PiercedZone),
			strconv.FormatBool(row.ClosedBack),
			strconv.FormatBool(row.DecisionRejectionCandidate),
			row.WickBeyondBucket,
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgWickBeyondPct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRRejectionTimingSummaryCSV(path string, rows []lab.SRRejectionTimingSummaryRow) error {
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
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"touched_count",
		"pierced_count",
		"closed_back_count",
		"decision_rejection_candidate_count",
		"touched_rate",
		"pierced_rate",
		"closed_back_rate",
		"decision_rejection_candidate_rate",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"decision_candidate_label_close_break_rate",
		"decision_candidate_label_wick_break_rate",
		"decision_candidate_label_reclaim_event_rate",
		"decision_candidate_label_reclaim_given_close_break_rate",
		"decision_candidate_label_rejection_rate",
		"decision_candidate_label_avg_favorable_pct",
		"decision_candidate_label_avg_adverse_pct",
		"decision_candidate_label_favorable_minus_adverse_pct",
		"decision_candidate_label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.TouchedCount),
			strconv.Itoa(row.PiercedCount),
			strconv.Itoa(row.ClosedBackCount),
			strconv.Itoa(row.DecisionRejectionCandidateCount),
			formatFloat(row.TouchedRate),
			formatFloat(row.PiercedRate),
			formatFloat(row.ClosedBackRate),
			formatFloat(row.DecisionRejectionCandidateRate),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.DecisionCandidateLabelCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelWickBreakRate),
			formatFloat(row.DecisionCandidateLabelReclaimEventRate),
			formatFloat(row.DecisionCandidateLabelReclaimGivenCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelRejectionRate),
			formatFloat(row.DecisionCandidateLabelAvgFavorablePct),
			formatFloat(row.DecisionCandidateLabelAvgAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableMinusAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRConfirmationTimingCandidatesCSV(path string, rows []lab.SRConfirmationTimingCandidateRow) error {
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
		"confirmation_delay_bars",
		"horizon_bars",
		"seed_close_location",
		"seed_pierced_zone",
		"seed_wick_beyond_bucket",
		"confirmation_close_location",
		"confirmation_favorable_close",
		"confirmation_wrong_side_close",
		"decision_confirmation_candidate",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_seed_wick_beyond_pct",
		"avg_confirmation_move_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.ConfirmationDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.SeedCloseLocation,
			strconv.FormatBool(row.SeedPiercedZone),
			row.SeedWickBeyondBucket,
			row.ConfirmationCloseLocation,
			strconv.FormatBool(row.ConfirmationFavorableClose),
			strconv.FormatBool(row.ConfirmationWrongSideClose),
			strconv.FormatBool(row.DecisionConfirmationCandidate),
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgSeedWickBeyondPct),
			formatFloat(row.AvgConfirmationMovePct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRConfirmationTimingSummaryCSV(path string, rows []lab.SRConfirmationTimingSummaryRow) error {
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
		"confirmation_delay_bars",
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"confirmation_favorable_close_count",
		"confirmation_wrong_side_close_count",
		"decision_confirmation_candidate_count",
		"confirmation_favorable_close_rate",
		"confirmation_wrong_side_close_rate",
		"decision_confirmation_candidate_rate",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"decision_candidate_label_close_break_rate",
		"decision_candidate_label_wick_break_rate",
		"decision_candidate_label_reclaim_event_rate",
		"decision_candidate_label_reclaim_given_close_break_rate",
		"decision_candidate_label_rejection_rate",
		"decision_candidate_label_avg_favorable_pct",
		"decision_candidate_label_avg_adverse_pct",
		"decision_candidate_label_favorable_minus_adverse_pct",
		"decision_candidate_label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.ConfirmationDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.ConfirmationFavorableCloseCount),
			strconv.Itoa(row.ConfirmationWrongSideCloseCount),
			strconv.Itoa(row.DecisionConfirmationCandidateCount),
			formatFloat(row.ConfirmationFavorableCloseRate),
			formatFloat(row.ConfirmationWrongSideCloseRate),
			formatFloat(row.DecisionConfirmationCandidateRate),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.DecisionCandidateLabelCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelWickBreakRate),
			formatFloat(row.DecisionCandidateLabelReclaimEventRate),
			formatFloat(row.DecisionCandidateLabelReclaimGivenCloseBreakRate),
			formatFloat(row.DecisionCandidateLabelRejectionRate),
			formatFloat(row.DecisionCandidateLabelAvgFavorablePct),
			formatFloat(row.DecisionCandidateLabelAvgAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableMinusAdversePct),
			formatFloat(row.DecisionCandidateLabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRFalseBreakReclaimTimingCandidatesCSV(path string, rows []lab.SRFalseBreakReclaimTimingCandidateRow) error {
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
		"break_delay_bars",
		"reclaim_delay_bars",
		"total_delay_bars",
		"horizon_bars",
		"anchor_close_location",
		"break_close_location",
		"reclaim_close_location",
		"break_move_bucket",
		"reclaim_move_bucket",
		"decision_false_break_reclaim_candidate",
		"strength_bucket",
		"distance_bucket",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"avg_score",
		"avg_distance_pct",
		"avg_break_move_pct",
		"avg_reclaim_move_pct",
		"label_close_break_count",
		"label_wick_break_count",
		"label_reclaimed_after_break_count",
		"label_rejected_count",
		"label_favorable_greater_than_adverse_count",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.BreakDelayBars),
			strconv.Itoa(row.ReclaimDelayBars),
			strconv.Itoa(row.TotalDelayBars),
			strconv.Itoa(row.HorizonBars),
			row.AnchorCloseLocation,
			row.BreakCloseLocation,
			row.ReclaimCloseLocation,
			row.BreakMoveBucket,
			row.ReclaimMoveBucket,
			strconv.FormatBool(row.DecisionFalseBreakReclaimCandidate),
			row.StrengthBucket,
			row.DistanceBucket,
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			formatFloat(row.AvgScore),
			formatFloat(row.AvgDistancePct),
			formatFloat(row.AvgBreakMovePct),
			formatFloat(row.AvgReclaimMovePct),
			strconv.Itoa(row.LabelCloseBreakCount),
			strconv.Itoa(row.LabelWickBreakCount),
			strconv.Itoa(row.LabelReclaimedAfterBreakCount),
			strconv.Itoa(row.LabelRejectedCount),
			strconv.Itoa(row.LabelFavorableGreaterThanAdverseCount),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func writeSRFalseBreakReclaimTimingSummaryCSV(path string, rows []lab.SRFalseBreakReclaimTimingSummaryRow) error {
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
		"horizon_bars",
		"detector_profile_id",
		"detector_raw_active",
		"detector_active",
		"candidate_count",
		"decision_false_break_reclaim_candidate_count",
		"decision_false_break_reclaim_candidate_rate",
		"avg_break_delay_bars",
		"avg_reclaim_delay_bars",
		"avg_total_delay_bars",
		"avg_break_move_pct",
		"avg_reclaim_move_pct",
		"label_close_break_rate",
		"label_wick_break_rate",
		"label_reclaim_event_rate",
		"label_reclaim_given_close_break_rate",
		"label_rejection_rate",
		"label_avg_favorable_pct",
		"label_avg_adverse_pct",
		"label_favorable_minus_adverse_pct",
		"label_favorable_greater_than_adverse_rate",
		"label_decision_candidate_close_break_rate",
		"label_decision_candidate_wick_break_rate",
		"label_decision_candidate_reclaim_event_rate",
		"label_decision_candidate_reclaim_given_close_break_rate",
		"label_decision_candidate_rejection_rate",
		"label_decision_candidate_avg_favorable_pct",
		"label_decision_candidate_avg_adverse_pct",
		"label_decision_candidate_favorable_minus_adverse_pct",
		"label_decision_candidate_favorable_greater_than_adverse_rate",
	}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{
			row.Split,
			row.Side,
			strconv.Itoa(row.HorizonBars),
			row.DetectorProfileID,
			strconv.FormatBool(row.DetectorRawActive),
			strconv.FormatBool(row.DetectorActive),
			strconv.Itoa(row.CandidateCount),
			strconv.Itoa(row.DecisionFalseBreakReclaimCandidateCount),
			formatFloat(row.DecisionFalseBreakReclaimCandidateRate),
			formatFloat(row.AvgBreakDelayBars),
			formatFloat(row.AvgReclaimDelayBars),
			formatFloat(row.AvgTotalDelayBars),
			formatFloat(row.AvgBreakMovePct),
			formatFloat(row.AvgReclaimMovePct),
			formatFloat(row.LabelCloseBreakRate),
			formatFloat(row.LabelWickBreakRate),
			formatFloat(row.LabelReclaimEventRate),
			formatFloat(row.LabelReclaimGivenCloseBreakRate),
			formatFloat(row.LabelRejectionRate),
			formatFloat(row.LabelAvgFavorablePct),
			formatFloat(row.LabelAvgAdversePct),
			formatFloat(row.LabelFavorableMinusAdversePct),
			formatFloat(row.LabelFavorableGreaterThanAdverseRate),
			formatFloat(row.LabelDecisionCandidateCloseBreakRate),
			formatFloat(row.LabelDecisionCandidateWickBreakRate),
			formatFloat(row.LabelDecisionCandidateReclaimEventRate),
			formatFloat(row.LabelDecisionCandidateReclaimGivenCloseBreakRate),
			formatFloat(row.LabelDecisionCandidateRejectionRate),
			formatFloat(row.LabelDecisionCandidateAvgFavorablePct),
			formatFloat(row.LabelDecisionCandidateAvgAdversePct),
			formatFloat(row.LabelDecisionCandidateFavorableMinusAdversePct),
			formatFloat(row.LabelDecisionCandidateFavorableGreaterThanAdverseRate),
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

func formatIntSlice(values []int) string {
	if len(values) == 0 {
		return ""
	}
	out := strconv.Itoa(values[0])
	for _, value := range values[1:] {
		out += ";" + strconv.Itoa(value)
	}
	return out
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
