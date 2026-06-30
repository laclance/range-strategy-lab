package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"range-strategy-lab/internal/lab"
)

func TestRunWithArgsAcceptsExplicitFuturesSourceAndWritesManifest(t *testing.T) {
	dir := t.TempDir()
	csvPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	outDir := filepath.Join(dir, "out")

	if err := runWithArgs([]string{
		"-csv", csvPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}

	manifest := readCLITestManifest(t, outDir)
	if manifest.Product != "Binance USDT-M futures" || manifest.RowCount != 2 || manifest.ComparisonOnly {
		t.Fatalf("unexpected futures manifest: %+v", manifest)
	}
}

func TestRunWithArgsRequiresSourceProductForNonDefaultCSV(t *testing.T) {
	dir := t.TempDir()
	csvPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")

	err := runWithArgs([]string{"-csv", csvPath, "-out-dir", filepath.Join(dir, "out")})
	if err == nil || !strings.Contains(err.Error(), "non-default -csv path requires explicit -source-product") {
		t.Fatalf("expected explicit source-product error, got %v", err)
	}
}

func TestRunWithArgsRejectsSpotByDefaultAndAllowsExplicitComparison(t *testing.T) {
	dir := t.TempDir()
	csvPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")

	err := runWithArgs([]string{
		"-csv", csvPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-out-dir", filepath.Join(dir, "rejected"),
	})
	if err == nil || !strings.Contains(err.Error(), "comparison-only") {
		t.Fatalf("expected spot comparison-only error, got %v", err)
	}

	outDir := filepath.Join(dir, "accepted")
	if err := runWithArgs([]string{
		"-csv", csvPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	manifest := readCLITestManifest(t, outDir)
	if manifest.Product != "Binance spot" || !manifest.ComparisonOnly {
		t.Fatalf("unexpected spot comparison manifest: %+v", manifest)
	}
}

func TestRunWithArgsRejectsSpotFilenameDuringFuturesRun(t *testing.T) {
	dir := t.TempDir()
	csvPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")

	err := runWithArgs([]string{
		"-csv", csvPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", filepath.Join(dir, "out"),
	})
	if err == nil || !strings.Contains(err.Error(), "spot-looking") {
		t.Fatalf("expected spot-looking futures error, got %v", err)
	}
}

func TestRunWithArgsPrototypeFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	outDir := filepath.Join(dir, "prototype")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-hold-inside-midline-touch-prototype",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"hold_inside_midline_touch_prototype_signals.csv",
		"hold_inside_midline_touch_prototype_trades.csv",
		"hold_inside_midline_touch_prototype_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected prototype artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-hold-inside-midline-touch-prototype",
		"-out-dir", filepath.Join(dir, "spot-prototype"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected prototype futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesImpulseAbsorptionFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_impulse_absorption_candidates.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write impulse absorption artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "impulse")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-impulse-absorption-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"futures_impulse_absorption_candidates.csv",
		"futures_impulse_absorption_summary.csv",
		"futures_impulse_absorption_stability.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected impulse absorption artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-impulse-absorption-audit",
		"-out-dir", filepath.Join(dir, "spot-impulse"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected impulse absorption futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeDiscoveryFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_discovery_candidates.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write range discovery artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "discovery")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-candidate-discovery-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"futures_range_discovery_coverage.csv",
		"futures_range_discovery_candidates.csv",
		"futures_range_discovery_summary.csv",
		"futures_range_discovery_rankings.csv",
		"futures_range_discovery_stability.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected range discovery artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("range discovery audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-candidate-discovery-audit",
		"-out-dir", filepath.Join(dir, "spot-discovery"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected range discovery futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseDiscoveryFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_universe.csv", 3)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_universe.csv", 3)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_universe.csv", 3)
	oldUniverseConfig := futuresRangeUniverseDiscoveryConfigForRun
	futuresRangeUniverseDiscoveryConfigForRun = func() lab.FuturesRangeUniverseDiscoveryAuditConfig {
		cfg := lab.DefaultFuturesRangeUniverseDiscoveryAuditConfig()
		cfg.Sources = []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.Discovery.DetectorLookbackBarsOverride = 2
		cfg.Discovery.DetectorMinConsecutiveBars = 1
		cfg.Discovery.MinCandidatesPerSplit = 1
		return cfg
	}
	defer func() { futuresRangeUniverseDiscoveryConfigForRun = oldUniverseConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_sources.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write universe artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "universe")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-discovery-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"futures_range_universe_sources.csv",
		"futures_range_universe_coverage.csv",
		"futures_range_universe_candidates.csv",
		"futures_range_universe_summary.csv",
		"futures_range_universe_rankings.csv",
		"futures_range_universe_stability.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected universe artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("range universe audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-discovery-audit",
		"-out-dir", filepath.Join(dir, "spot-universe"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected universe futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesCleanBreakoutBaselineFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_clean_breakout_baseline_signals.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write clean breakout artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "clean-breakout")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-clean-breakout-baseline-backtest",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_clean_breakout_baseline_signals.csv",
		"futures_clean_breakout_baseline_trades.csv",
		"futures_clean_breakout_baseline_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected clean breakout artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-clean-breakout-baseline-backtest",
		"-out-dir", filepath.Join(dir, "spot-clean-breakout"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected clean breakout futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseBreakoutRetestAcceptanceBaselineFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_breakout_retest.csv", 48)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_breakout_retest.csv", 48)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_breakout_retest.csv", 48)
	oldBreakoutRetestConfig := futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun
	futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun = func() lab.FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig {
		cfg := lab.DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
		cfg.DiscoveryConfig.Sources = []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.DiscoveryConfig.Discovery.DetectorLookbackBarsOverride = 2
		cfg.DiscoveryConfig.Discovery.DetectorMinConsecutiveBars = 1
		cfg.DiscoveryConfig.Discovery.MinCandidatesPerSplit = 1
		return cfg
	}
	defer func() { futuresRangeUniverseBreakoutRetestAcceptanceBaselineConfigForRun = oldBreakoutRetestConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_breakout_retest_acceptance_baseline_signals.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write breakout retest acceptance artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "breakout-retest")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-breakout-retest-acceptance-baseline-backtest",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_universe_breakout_retest_acceptance_baseline_sources.csv",
		"futures_range_universe_breakout_retest_acceptance_baseline_coverage.csv",
		"futures_range_universe_breakout_retest_acceptance_baseline_selection.csv",
		"futures_range_universe_breakout_retest_acceptance_baseline_signals.csv",
		"futures_range_universe_breakout_retest_acceptance_baseline_trades.csv",
		"futures_range_universe_breakout_retest_acceptance_baseline_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected breakout retest acceptance artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-breakout-retest-acceptance-baseline-backtest",
		"-out-dir", filepath.Join(dir, "spot-breakout-retest"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected breakout retest acceptance futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-breakout-retest-acceptance-baseline-backtest",
		"-futures-range-universe-structured-compression-baseline-backtest",
		"-out-dir", filepath.Join(dir, "combined-breakout-retest"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected breakout retest acceptance combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseStructuredCompressionBaselineFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_structured.csv", 48)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_structured.csv", 48)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_structured.csv", 48)
	oldStructuredConfig := futuresRangeUniverseStructuredCompressionBaselineConfigForRun
	futuresRangeUniverseStructuredCompressionBaselineConfigForRun = func() lab.FuturesRangeUniverseStructuredCompressionBaselineConfig {
		cfg := lab.DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
		cfg.Sources = []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.DetectorMinConsecutiveBars = 1
		return cfg
	}
	defer func() { futuresRangeUniverseStructuredCompressionBaselineConfigForRun = oldStructuredConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_structured_compression_baseline_signals.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write structured compression artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "structured")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-structured-compression-baseline-backtest",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_universe_structured_compression_baseline_sources.csv",
		"futures_range_universe_structured_compression_baseline_coverage.csv",
		"futures_range_universe_structured_compression_baseline_signals.csv",
		"futures_range_universe_structured_compression_baseline_trades.csv",
		"futures_range_universe_structured_compression_baseline_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected structured compression artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-structured-compression-baseline-backtest",
		"-out-dir", filepath.Join(dir, "spot-structured"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected structured compression futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseStructuredCompressionOptimizationFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_optimization.csv", 48)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_optimization.csv", 48)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_optimization.csv", 48)
	oldOptimizationConfig := futuresRangeUniverseStructuredCompressionOptimizationConfigForRun
	futuresRangeUniverseStructuredCompressionOptimizationConfigForRun = func() lab.FuturesRangeUniverseStructuredCompressionOptimizationConfig {
		cfg := lab.DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig()
		cfg.Sources = []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.ConfirmationWindowBars = []int{2}
		cfg.MaxHoldBars = []int{4}
		cfg.TargetRangeWidthMultiples = []float64{1}
		cfg.StopBoundaryBufferRangeWidths = []float64{0}
		cfg.SymbolSets = []string{lab.StructuredCompressionOptimizationSymbolSetAll}
		cfg.DetectorMinConsecutiveBars = 1
		return cfg
	}
	defer func() { futuresRangeUniverseStructuredCompressionOptimizationConfigForRun = oldOptimizationConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_structured_compression_optimization_grid.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write structured compression optimization artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "optimization")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-structured-compression-optimization",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_universe_structured_compression_optimization_sources.csv",
		"futures_range_universe_structured_compression_optimization_coverage.csv",
		"futures_range_universe_structured_compression_optimization_grid.csv",
		"futures_range_universe_structured_compression_optimization_trades.csv",
		"futures_range_universe_structured_compression_optimization_summary.csv",
		"futures_range_universe_structured_compression_optimization_rankings.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected structured compression optimization artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-structured-compression-optimization",
		"-out-dir", filepath.Join(dir, "spot-optimization"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected structured compression optimization futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseStructuredCompressionStrategyReplayFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_strategy.csv", 48)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_strategy.csv", 48)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_strategy.csv", 48)
	oldStrategyReplayConfig := futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun
	futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun = func() lab.FuturesRangeUniverseStructuredCompressionStrategyReplayConfig {
		cfg := lab.DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
		cfg.Sources = []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.DetectorMinConsecutiveBars = 1
		return cfg
	}
	defer func() { futuresRangeUniverseStructuredCompressionStrategyReplayConfigForRun = oldStrategyReplayConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_structured_compression_strategy_signals.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write structured compression strategy replay artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "strategy-replay")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-structured-compression-strategy-replay",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_universe_structured_compression_strategy_sources.csv",
		"futures_range_universe_structured_compression_strategy_coverage.csv",
		"futures_range_universe_structured_compression_strategy_signals.csv",
		"futures_range_universe_structured_compression_strategy_trades.csv",
		"futures_range_universe_structured_compression_strategy_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected structured compression strategy replay artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-structured-compression-strategy-replay",
		"-out-dir", filepath.Join(dir, "spot-strategy-replay"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected structured compression strategy replay futures-source error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeUniverseStructuredCompressionWalkForwardFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSV(t, dir, "btcusdt_futures_um_5m_test.csv")
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_walk_forward.csv", 48)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_walk_forward.csv", 48)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_walk_forward.csv", 48)
	oldWalkForwardConfig := futuresRangeUniverseStructuredCompressionWalkForwardConfigForRun
	futuresRangeUniverseStructuredCompressionWalkForwardConfigForRun = func() lab.FuturesRangeUniverseStructuredCompressionWalkForwardConfig {
		cfg := lab.DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig()
		sources := []lab.FuturesRangeUniverseSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSplitEligibilityCheck: true},
		}
		cfg.OptimizationConfig.Sources = sources
		cfg.OptimizationConfig.DetectorMinConsecutiveBars = 1
		cfg.FrozenReplayConfig.Sources = sources
		cfg.FrozenReplayConfig.DetectorMinConsecutiveBars = 1
		return cfg
	}
	defer func() { futuresRangeUniverseStructuredCompressionWalkForwardConfigForRun = oldWalkForwardConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_universe_structured_compression_walk_forward_grid.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write structured compression walk-forward artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "walk-forward")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-structured-compression-walk-forward-robustness",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_universe_structured_compression_walk_forward_sources.csv",
		"futures_range_universe_structured_compression_walk_forward_coverage.csv",
		"futures_range_universe_structured_compression_walk_forward_grid.csv",
		"futures_range_universe_structured_compression_walk_forward_folds.csv",
		"futures_range_universe_structured_compression_walk_forward_trades.csv",
		"futures_range_universe_structured_compression_walk_forward_summary.csv",
		"futures_range_universe_structured_compression_walk_forward_rankings.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected structured compression walk-forward artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-universe-structured-compression-walk-forward-robustness",
		"-out-dir", filepath.Join(dir, "spot-walk-forward"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected structured compression walk-forward futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-universe-structured-compression-optimization",
		"-futures-range-universe-structured-compression-walk-forward-robustness",
		"-out-dir", filepath.Join(dir, "combined-walk-forward"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected structured compression combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesHigherTFNestedRangeRotationAuditFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_nested_rotation.csv", 96)
	oldNestedConfig := futuresHigherTFNestedRangeRotationAuditConfigForRun
	futuresHigherTFNestedRangeRotationAuditConfigForRun = func() lab.FuturesHigherTFNestedRangeRotationAuditConfig {
		cfg := lab.DefaultFuturesHigherTFNestedRangeRotationAuditConfig()
		cfg.ApprovedSourcePath = futuresPath
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.DetectorLookbackBarsOverride = 1
		cfg.DetectorMinConsecutiveBars = 1
		return cfg
	}
	defer func() { futuresHigherTFNestedRangeRotationAuditConfigForRun = oldNestedConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_higher_tf_nested_range_rotation_events.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write nested range rotation artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "nested-rotation")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-higher-tf-nested-range-rotation-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_higher_tf_nested_range_rotation_sources.csv",
		"futures_higher_tf_nested_range_rotation_coverage.csv",
		"futures_higher_tf_nested_range_rotation_parent_ranges.csv",
		"futures_higher_tf_nested_range_rotation_child_ranges.csv",
		"futures_higher_tf_nested_range_rotation_events.csv",
		"futures_higher_tf_nested_range_rotation_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected nested range rotation artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_test.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-higher-tf-nested-range-rotation-audit",
		"-out-dir", filepath.Join(dir, "spot-nested-rotation"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected nested range rotation futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-higher-tf-nested-range-rotation-audit",
		"-futures-clean-breakout-baseline-backtest",
		"-out-dir", filepath.Join(dir, "combined-nested-rotation"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected nested range rotation combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeFirstOccupancyRotationV1OptimizationFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_occupancy_rotation.csv", 60)
	oldOccupancyConfig := futuresRangeFirstOccupancyRotationV1OptimizationConfigForRun
	futuresRangeFirstOccupancyRotationV1OptimizationConfigForRun = func() lab.FuturesRangeFirstOccupancyRotationV1OptimizationConfig {
		cfg := lab.DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig()
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.Timeframes = []string{lab.RangeDiscoveryTimeframe1h}
		cfg.LookbackHours = []int{1, 48}
		cfg.MaxWidthPcts = []float64{0.50}
		cfg.OccupancyWindows = []int{2, 12}
		cfg.OccupancyZoneLevels = []float64{0.25}
		cfg.OccupancyMinFractions = []float64{0.50}
		cfg.RecaptureLevels = []float64{0.33}
		cfg.TargetLevels = []float64{0.66}
		cfg.MaxHoldBars = []int{2}
		cfg.StopBufferWidths = []float64{0.05}
		cfg.SideModes = []string{lab.RangeDiscoverySideAll}
		cfg.MinTrainTrades = 1
		cfg.MinOOSTrades = 1
		cfg.MinRecentTrades = 1
		cfg.MinFullTrades = 1
		return cfg
	}
	defer func() { futuresRangeFirstOccupancyRotationV1OptimizationConfigForRun = oldOccupancyConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_first_occupancy_rotation_v1_grid.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write occupancy rotation artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "occupancy")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-first-occupancy-rotation-v1-optimization",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_first_occupancy_rotation_v1_sources.csv",
		"futures_range_first_occupancy_rotation_v1_coverage.csv",
		"futures_range_first_occupancy_rotation_v1_grid.csv",
		"futures_range_first_occupancy_rotation_v1_baseline.csv",
		"futures_range_first_occupancy_rotation_v1_signals.csv",
		"futures_range_first_occupancy_rotation_v1_trades.csv",
		"futures_range_first_occupancy_rotation_v1_summary.csv",
		"futures_range_first_occupancy_rotation_v1_rankings.csv",
		"futures_range_first_occupancy_rotation_v1_selection.csv",
		"futures_range_first_occupancy_rotation_v1_skips.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected occupancy rotation artifact %s: %v", name, err)
		}
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_occupancy.csv")
	err := runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-first-occupancy-rotation-v1-optimization",
		"-out-dir", filepath.Join(dir, "spot-occupancy"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected occupancy rotation futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-first-occupancy-rotation-v1-optimization",
		"-futures-range-universe-structured-compression-optimization",
		"-out-dir", filepath.Join(dir, "combined-occupancy"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected occupancy rotation combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeContextTriageAuditFlagWritesArtifactsAndRejectsSpotComparison(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_context_triage.csv", 96)
	oldTriageConfig := futuresRangeContextTriageAuditConfigForRun
	futuresRangeContextTriageAuditConfigForRun = func() lab.FuturesRangeContextTriageAuditConfig {
		cfg := lab.DefaultFuturesRangeContextTriageAuditConfig()
		cfg.ApprovedSourcePath = futuresPath
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		cfg.DetectorLookbackBarsOverride = 1
		cfg.DetectorMinConsecutiveBars = 1
		cfg.HorizonsBars = []int{2}
		cfg.MinFullCohortCount = 1
		cfg.MinSplitCohortCount = 1
		cfg.MinSessionSplitCohortCount = 1
		return cfg
	}
	defer func() { futuresRangeContextTriageAuditConfigForRun = oldTriageConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_context_triage_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write range context triage artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "triage")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-context-triage-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_context_triage_sources.csv",
		"futures_range_context_triage_coverage.csv",
		"futures_range_context_triage_episodes.csv",
		"futures_range_context_triage_quality.csv",
		"futures_range_context_triage_sessions.csv",
		"futures_range_context_triage_failure_modes.csv",
		"futures_range_context_triage_cohorts.csv",
		"futures_range_context_triage_rankings.csv",
		"futures_range_context_triage_summary.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected range context triage artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("range context triage audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_context_triage.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-context-triage-audit",
		"-out-dir", filepath.Join(dir, "spot-triage"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected range context triage futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-context-triage-audit",
		"-futures-range-first-occupancy-rotation-v1-optimization",
		"-out-dir", filepath.Join(dir, "combined-triage"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected range context triage combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeStateConstructionLoopAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_state_loop.csv", 96)
	oldStateConfig := futuresRangeStateConstructionLoopAuditConfigForRun
	futuresRangeStateConstructionLoopAuditConfigForRun = func() lab.FuturesRangeStateConstructionLoopAuditConfig {
		cfg := lab.DefaultFuturesRangeStateConstructionLoopAuditConfig()
		cfg.ApprovedSourcePath = futuresPath
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		cfg.DetectorLookbackBarsOverride = 1
		cfg.DetectorMinConsecutiveBars = 1
		cfg.ShortWindowBars = 2
		cfg.MediumWindowBars = 4
		cfg.FeatureLookbackBars = 4
		cfg.LongLookbackDays = 1
		cfg.Horizons15M = []int{2}
		cfg.MinFullCohortCount = 1
		cfg.MinSplitCohortCount = 1
		cfg.MinNoTradeFullCohortCount = 1
		cfg.MinNoTradeSplitCohortCount = 1
		return cfg
	}
	defer func() { futuresRangeStateConstructionLoopAuditConfigForRun = oldStateConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_state_construction_loop_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write range-state construction artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "state-loop")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-state-construction-loop-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_state_construction_loop_sources.csv",
		"futures_range_state_construction_loop_coverage.csv",
		"futures_range_state_construction_loop_feature_windows.csv",
		"futures_range_state_construction_loop_states.csv",
		"futures_range_state_construction_loop_labels.csv",
		"futures_range_state_construction_loop_cohorts.csv",
		"futures_range_state_construction_loop_rankings.csv",
		"futures_range_state_construction_loop_summary.csv",
		"futures_range_state_construction_loop_skips.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected range-state construction artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("range-state construction audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_state_loop.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-state-construction-loop-audit",
		"-out-dir", filepath.Join(dir, "spot-state-loop"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected range-state construction futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-state-construction-loop-audit",
		"-futures-range-context-triage-audit",
		"-out-dir", filepath.Join(dir, "combined-state-loop"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected range-state construction combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeContextRouterAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_context_router.csv", 96)
	oldRouterConfig := futuresRangeContextRouterAuditConfigForRun
	futuresRangeContextRouterAuditConfigForRun = func() lab.FuturesRangeContextRouterAuditConfig {
		cfg := lab.DefaultFuturesRangeContextRouterAuditConfig()
		cfg.StateAuditConfig.ApprovedSourcePath = futuresPath
		cfg.StateAuditConfig.SkipSourceFactCheck = true
		cfg.StateAuditConfig.SkipCoverageCountCheck = true
		cfg.StateAuditConfig.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		cfg.StateAuditConfig.DetectorLookbackBarsOverride = 1
		cfg.StateAuditConfig.DetectorMinConsecutiveBars = 1
		cfg.StateAuditConfig.ShortWindowBars = 2
		cfg.StateAuditConfig.MediumWindowBars = 4
		cfg.StateAuditConfig.FeatureLookbackBars = 4
		cfg.StateAuditConfig.LongLookbackDays = 1
		cfg.StateAuditConfig.Horizons15M = []int{2}
		cfg.StateAuditConfig.MinFullCohortCount = 1
		cfg.StateAuditConfig.MinSplitCohortCount = 1
		cfg.StateAuditConfig.MinNoTradeFullCohortCount = 1
		cfg.StateAuditConfig.MinNoTradeSplitCohortCount = 1
		cfg.MinFullRouterRows = 1
		cfg.MinSplitRouterRows = 1
		cfg.MinNoTradeFullRouterRows = 1
		cfg.MinNoTradeSplitRouterRows = 1
		return cfg
	}
	defer func() { futuresRangeContextRouterAuditConfigForRun = oldRouterConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_context_router_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write range context router artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "router")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-context-router-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_context_router_sources.csv",
		"futures_range_context_router_coverage.csv",
		"futures_range_context_router_rules.csv",
		"futures_range_context_router_rows.csv",
		"futures_range_context_router_cohorts.csv",
		"futures_range_context_router_rankings.csv",
		"futures_range_context_router_summary.csv",
		"futures_range_context_router_skips.csv",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected range context router artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("range context router audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_context_router.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-context-router-audit",
		"-out-dir", filepath.Join(dir, "spot-router"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected range context router futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-context-router-audit",
		"-futures-range-state-construction-loop-audit",
		"-out-dir", filepath.Join(dir, "combined-router"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected range context router combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesRangeRouterRotationPremiseAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_router_rotation_premise.csv", 180)
	oldPremiseConfig := futuresRangeRouterRotationPremiseAuditConfigForRun
	futuresRangeRouterRotationPremiseAuditConfigForRun = func() lab.FuturesRangeRouterRotationPremiseAuditConfig {
		cfg := lab.DefaultFuturesRangeRouterRotationPremiseAuditConfig()
		cfg.SkipCoverageCountCheck = true
		cfg.RouterAuditConfig.StateAuditConfig.ApprovedSourcePath = futuresPath
		cfg.RouterAuditConfig.StateAuditConfig.SkipSourceFactCheck = true
		cfg.RouterAuditConfig.StateAuditConfig.SkipCoverageCountCheck = true
		cfg.RouterAuditConfig.StateAuditConfig.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		cfg.RouterAuditConfig.StateAuditConfig.DetectorLookbackBarsOverride = 1
		cfg.RouterAuditConfig.StateAuditConfig.DetectorMinConsecutiveBars = 1
		cfg.RouterAuditConfig.StateAuditConfig.ShortWindowBars = 2
		cfg.RouterAuditConfig.StateAuditConfig.MediumWindowBars = 4
		cfg.RouterAuditConfig.StateAuditConfig.FeatureLookbackBars = 4
		cfg.RouterAuditConfig.StateAuditConfig.LongLookbackDays = 1
		cfg.RouterAuditConfig.StateAuditConfig.Horizons15M = []int{24}
		cfg.RouterAuditConfig.StateAuditConfig.MinFullCohortCount = 1
		cfg.RouterAuditConfig.StateAuditConfig.MinSplitCohortCount = 1
		cfg.RouterAuditConfig.StateAuditConfig.MinNoTradeFullCohortCount = 1
		cfg.RouterAuditConfig.StateAuditConfig.MinNoTradeSplitCohortCount = 1
		cfg.RouterAuditConfig.MinFullRouterRows = 1
		cfg.RouterAuditConfig.MinSplitRouterRows = 1
		cfg.RouterAuditConfig.MinNoTradeFullRouterRows = 1
		cfg.RouterAuditConfig.MinNoTradeSplitRouterRows = 1
		cfg.MinContextSegmentsFull = 1
		cfg.MinValidEventsFull = 1
		cfg.MinValidEventsPerSplit = 1
		cfg.MinValidEventsPerSide = 1
		return cfg
	}
	defer func() { futuresRangeRouterRotationPremiseAuditConfigForRun = oldPremiseConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_range_router_rotation_premise_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write router rotation premise artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "premise")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-router-rotation-premise-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_range_router_rotation_premise_sources.csv",
		"futures_range_router_rotation_premise_sources.json",
		"futures_range_router_rotation_premise_coverage.csv",
		"futures_range_router_rotation_premise_coverage.json",
		"futures_range_router_rotation_premise_router_dependency.csv",
		"futures_range_router_rotation_premise_router_dependency.json",
		"futures_range_router_rotation_premise_context_segments.csv",
		"futures_range_router_rotation_premise_context_segments.json",
		"futures_range_router_rotation_premise_events.csv",
		"futures_range_router_rotation_premise_events.json",
		"futures_range_router_rotation_premise_outcomes.csv",
		"futures_range_router_rotation_premise_outcomes.json",
		"futures_range_router_rotation_premise_cohorts.csv",
		"futures_range_router_rotation_premise_cohorts.json",
		"futures_range_router_rotation_premise_rankings.csv",
		"futures_range_router_rotation_premise_rankings.json",
		"futures_range_router_rotation_premise_summary.csv",
		"futures_range_router_rotation_premise_summary.json",
		"futures_range_router_rotation_premise_skips.csv",
		"futures_range_router_rotation_premise_skips.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected router rotation premise artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("router rotation premise audit must remain zero-trade, got %d", len(trades))
	}

	defaultDirCWD := t.TempDir()
	t.Chdir(defaultDirCWD)
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-router-rotation-premise-audit",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultDirCWD, "results", "futures-range-router-rotation-premise-audit", "futures_range_router_rotation_premise_summary.csv")); err != nil {
		t.Fatalf("expected default router rotation premise out dir artifact: %v", err)
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_router_rotation_premise.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-range-router-rotation-premise-audit",
		"-out-dir", filepath.Join(dir, "spot-premise"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected router rotation premise futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-router-rotation-premise-audit",
		"-futures-clean-breakout-baseline-backtest",
		"-out-dir", filepath.Join(dir, "combined-premise"),
	})
	if err == nil || !strings.Contains(err.Error(), lab.RangeRouterRotationPremiseStopStateRejectedStrategyBacktest) {
		t.Fatalf("expected router rotation premise strategy/backtest rejection, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-range-router-rotation-premise-audit",
		"-futures-range-context-router-audit",
		"-out-dir", filepath.Join(dir, "combined-router-premise"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected router rotation premise combination error, got %v", err)
	}
}

func TestRunWithArgsFuturesBTCRegimeETHSOLContextAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_btc_context_cli.csv", 96)
	btcPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_btc_context.csv", 96)
	ethPath := writeCLITestCSVN(t, dir, "ethusdt_futures_um_5m_btc_context.csv", 96)
	solPath := writeCLITestCSVN(t, dir, "solusdt_futures_um_5m_btc_context.csv", 96)
	oldContextConfig := futuresBTCRegimeETHSOLContextAuditConfigForRun
	futuresBTCRegimeETHSOLContextAuditConfigForRun = func() lab.FuturesBTCRegimeETHSOLContextAuditConfig {
		cfg := lab.DefaultFuturesBTCRegimeETHSOLContextAuditConfig()
		cfg.Sources = []lab.FuturesBTCRegimeETHSOLContextSourceConfig{
			{Symbol: lab.RangeUniverseSymbolBTCUSDT, Path: btcPath, ApprovedPath: btcPath, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolETHUSDT, Path: ethPath, ApprovedPath: ethPath, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
			{Symbol: lab.RangeUniverseSymbolSOLUSDT, Path: solPath, ApprovedPath: solPath, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
		}
		cfg.StateConfig.SkipSourceFactCheck = true
		cfg.StateConfig.SkipCoverageCountCheck = true
		cfg.StateConfig.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		cfg.StateConfig.DetectorLookbackBarsOverride = 1
		cfg.StateConfig.DetectorMinConsecutiveBars = 1
		cfg.StateConfig.ShortWindowBars = 2
		cfg.StateConfig.MediumWindowBars = 4
		cfg.StateConfig.FeatureLookbackBars = 4
		cfg.StateConfig.LongLookbackDays = 1
		cfg.StateConfig.Horizons15M = []int{2}
		cfg.StateConfig.MinFullCohortCount = 1
		cfg.StateConfig.MinSplitCohortCount = 1
		cfg.StateConfig.MinNoTradeFullCohortCount = 1
		cfg.StateConfig.MinNoTradeSplitCohortCount = 1
		cfg.MinFullCohortRows = 1
		cfg.MinSplitCohortRows = 1
		return cfg
	}
	defer func() { futuresBTCRegimeETHSOLContextAuditConfigForRun = oldContextConfig }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_btc_regime_eth_sol_context_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write BTC regime ETH/SOL context artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "btc-context")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-btc-regime-eth-sol-context-audit",
		"-out-dir", outDir,
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json",
		"summary.csv",
		"summary.json",
		"trades.json",
		"futures_btc_regime_eth_sol_context_sources.csv",
		"futures_btc_regime_eth_sol_context_sources.json",
		"futures_btc_regime_eth_sol_context_coverage.csv",
		"futures_btc_regime_eth_sol_context_coverage.json",
		"futures_btc_regime_eth_sol_context_btc_states.csv",
		"futures_btc_regime_eth_sol_context_btc_states.json",
		"futures_btc_regime_eth_sol_context_local_states.csv",
		"futures_btc_regime_eth_sol_context_local_states.json",
		"futures_btc_regime_eth_sol_context_relative_strength.csv",
		"futures_btc_regime_eth_sol_context_relative_strength.json",
		"futures_btc_regime_eth_sol_context_labels.csv",
		"futures_btc_regime_eth_sol_context_labels.json",
		"futures_btc_regime_eth_sol_context_cohorts.csv",
		"futures_btc_regime_eth_sol_context_cohorts.json",
		"futures_btc_regime_eth_sol_context_rankings.csv",
		"futures_btc_regime_eth_sol_context_rankings.json",
		"futures_btc_regime_eth_sol_context_summary.csv",
		"futures_btc_regime_eth_sol_context_summary.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected BTC regime ETH/SOL context artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("BTC regime ETH/SOL context audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_btc_context.csv")
	err = runWithArgs([]string{
		"-csv", spotPath,
		"-source-product", lab.SourceProductBinanceSpot,
		"-allow-spot-comparison",
		"-futures-btc-regime-eth-sol-context-audit",
		"-out-dir", filepath.Join(dir, "spot-btc-context"),
	})
	if err == nil || !strings.Contains(err.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected BTC regime ETH/SOL context futures-source error, got %v", err)
	}

	err = runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-btc-regime-eth-sol-context-audit",
		"-futures-clean-breakout-baseline-backtest",
		"-out-dir", filepath.Join(dir, "combined-btc-context"),
	})
	if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected BTC regime ETH/SOL context combination error, got %v", err)
	}
}

func writeCLIDerivSourceCSV(t *testing.T, dir, name string, rows int) string {
	t.Helper()
	path := filepath.Join(dir, name)
	data := "open_time,open,high,low,close,close_time,source_object_id\n"
	for i := 0; i < rows; i++ {
		open := int64(1609459200000 + i*300000)
		closeTime := open + 299999
		data += strconv.FormatInt(open, 10) + ",100,101,99,100," + strconv.FormatInt(closeTime, 10) + ",obj-" + strconv.Itoa(i) + "\n"
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunWithArgsFuturesDerivativesContextSourceAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_deriv_cli.csv", 96)

	// Derivative + anchor sources must live outside /tmp (the audit rejects /tmp).
	srcDir, err := os.MkdirTemp(".", "derivcli")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(srcDir) })
	srcAbs, err := filepath.Abs(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	const n = 24
	anchorPath := writeCLITestCSVN(t, srcAbs, "btcusdt_futures_um_5m_deriv_anchor.csv", n)
	markPath := writeCLIDerivSourceCSV(t, srcAbs, "binance_usdm_mark_price_klines_5m_BTCUSDT_cli.csv", n)
	indexPath := writeCLIDerivSourceCSV(t, srcAbs, "binance_usdm_index_price_klines_5m_BTCUSDT_cli.csv", n)
	premiumPath := writeCLIDerivSourceCSV(t, srcAbs, "binance_usdm_premium_index_klines_5m_BTCUSDT_cli.csv", n)

	oldCfg := futuresDerivativesContextSourceAuditConfigForRun
	futuresDerivativesContextSourceAuditConfigForRun = func() lab.FuturesDerivativesContextSourceAuditConfig {
		return lab.FuturesDerivativesContextSourceAuditConfig{
			DerivativeSources: []lab.FuturesDerivativesContextSourceFileConfig{
				{Symbol: "BTCUSDT", SourceFamily: "mark_price_klines", ArchiveFamily: "markPriceKlines", Path: markPath, Required: true},
				{Symbol: "BTCUSDT", SourceFamily: "index_price_klines", ArchiveFamily: "indexPriceKlines", Path: indexPath, Required: true},
				{Symbol: "BTCUSDT", SourceFamily: "premium_index_klines", ArchiveFamily: "premiumIndexKlines", Path: premiumPath, Required: false, AllowNonPositive: true},
			},
			Anchors: []lab.FuturesRangeUniverseSourceConfig{
				{Symbol: "BTCUSDT", Path: anchorPath, ApprovedPath: anchorPath, SkipSplitEligibilityCheck: true},
			},
			MinAlignedCoverage:       0.5,
			ConservativeLagIntervals: 1,
			EraStartMs:               1609459200000,
			EraEndMs:                 1609459200000 + int64(n-1)*300000,
		}
	}
	defer func() { futuresDerivativesContextSourceAuditConfigForRun = oldCfg }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-out-dir", defaultOutDir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_derivatives_context_source_audit_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write derivatives source audit artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "deriv-audit")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-context-source-audit", "-out-dir", outDir}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json", "summary.csv", "summary.json", "trades.json",
		"futures_derivatives_context_source_audit_sources.csv", "futures_derivatives_context_source_audit_sources.json",
		"futures_derivatives_context_source_audit_candle_anchors.csv", "futures_derivatives_context_source_audit_candle_anchors.json",
		"futures_derivatives_context_source_audit_external_coverage.csv", "futures_derivatives_context_source_audit_external_coverage.json",
		"futures_derivatives_context_source_audit_timestamp_alignment.csv", "futures_derivatives_context_source_audit_timestamp_alignment.json",
		"futures_derivatives_context_source_audit_publication_lag.csv", "futures_derivatives_context_source_audit_publication_lag.json",
		"futures_derivatives_context_source_audit_missingness.csv", "futures_derivatives_context_source_audit_missingness.json",
		"futures_derivatives_context_source_audit_provenance.csv", "futures_derivatives_context_source_audit_provenance.json",
		"futures_derivatives_context_source_audit_skips.csv", "futures_derivatives_context_source_audit_skips.json",
		"futures_derivatives_context_source_audit_summary.csv", "futures_derivatives_context_source_audit_summary.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected derivatives source audit artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("derivatives source audit must remain zero-trade, got %d", len(trades))
	}

	conflictErr := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-context-source-audit", "-futures-btc-regime-eth-sol-context-audit", "-out-dir", filepath.Join(dir, "conflict")})
	if conflictErr == nil {
		t.Fatalf("expected error combining derivatives source audit with another audit flag")
	}
}

func TestRunWithArgsFuturesDerivativesContextAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_deriv_context_cli.csv", 96)

	srcDir, err := os.MkdirTemp(".", "derivcontextcli")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(srcDir) })
	srcAbs, err := filepath.Abs(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	const n = 96
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}
	anchors := map[string]string{}
	for _, sym := range symbols {
		anchors[sym] = writeCLITestCSVN(t, srcAbs, strings.ToLower(sym)+"_futures_um_5m_deriv_context_anchor.csv", n)
	}
	sourceFamilies := []struct {
		filePart string
		family   string
		archive  string
		required bool
	}{
		{"mark_price_klines", "mark_price_klines", "markPriceKlines", true},
		{"index_price_klines", "index_price_klines", "indexPriceKlines", true},
		{"premium_index_klines", "premium_index_klines", "premiumIndexKlines", false},
	}
	var derivSources []lab.FuturesDerivativesContextSourceFileConfig
	for _, sym := range symbols {
		for _, family := range sourceFamilies {
			path := writeCLIDerivSourceCSV(t, srcAbs, "binance_usdm_"+family.filePart+"_5m_"+sym+"_context_cli.csv", n)
			derivSources = append(derivSources, lab.FuturesDerivativesContextSourceFileConfig{
				Symbol:           sym,
				SourceFamily:     family.family,
				ArchiveFamily:    family.archive,
				Path:             path,
				Required:         family.required,
				AllowNonPositive: !family.required,
			})
		}
	}

	oldCfg := futuresDerivativesContextAuditConfigForRun
	futuresDerivativesContextAuditConfigForRun = func() lab.FuturesDerivativesContextAuditConfig {
		stateCfg := lab.DefaultFuturesRangeStateConstructionLoopAuditConfig()
		stateCfg.SkipSourceFactCheck = true
		stateCfg.SkipCoverageCountCheck = true
		stateCfg.DetectorLookbackBarsOverride = 1
		stateCfg.DetectorMinConsecutiveBars = 1
		stateCfg.ShortWindowBars = 2
		stateCfg.MediumWindowBars = 4
		stateCfg.FeatureLookbackBars = 4
		stateCfg.LongLookbackDays = 1
		stateCfg.Timeframes = []string{lab.RangeDiscoveryTimeframe15m}
		stateCfg.Horizons15M = []int{2}
		return lab.FuturesDerivativesContextAuditConfig{
			SourceAuditConfig: lab.FuturesDerivativesContextSourceAuditConfig{
				DerivativeSources:        derivSources,
				Anchors:                  []lab.FuturesRangeUniverseSourceConfig{{Symbol: "BTCUSDT", Path: anchors["BTCUSDT"], ApprovedPath: anchors["BTCUSDT"], SkipSplitEligibilityCheck: true}, {Symbol: "ETHUSDT", Path: anchors["ETHUSDT"], ApprovedPath: anchors["ETHUSDT"], SkipSplitEligibilityCheck: true}, {Symbol: "SOLUSDT", Path: anchors["SOLUSDT"], ApprovedPath: anchors["SOLUSDT"], SkipSplitEligibilityCheck: true}},
				MinAlignedCoverage:       0.5,
				ConservativeLagIntervals: 1,
				EraStartMs:               1609459200000,
				EraEndMs:                 1609459200000 + int64(n-1)*300000,
			},
			LocalSources: []lab.FuturesBTCRegimeETHSOLContextSourceConfig{
				{Symbol: "BTCUSDT", Path: anchors["BTCUSDT"], ApprovedPath: anchors["BTCUSDT"], SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
				{Symbol: "ETHUSDT", Path: anchors["ETHUSDT"], ApprovedPath: anchors["ETHUSDT"], SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
				{Symbol: "SOLUSDT", Path: anchors["SOLUSDT"], ApprovedPath: anchors["SOLUSDT"], SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
			},
			StateConfig:                        stateCfg,
			MinBasisContextCoverage:            0.5,
			BasisChangeLookbackIntervals:       2,
			BasisVolatilityLookbackIntervals:   4,
			MinFullCohortRows:                  1,
			MinSplitCohortRows:                 1,
			MaxSplitContributionRate:           1,
			MinUsefulRateFull:                  0.1,
			MinUsefulRateSplit:                 0.1,
			MaxToxicRateFull:                   1,
			MaxToxicRateSplit:                  1,
			MinUsefulMinusToxicMarginFull:      -1,
			MinUsefulMinusToxicMarginSplit:     -1,
			MinToxicRateFull:                   0.1,
			MinToxicRateSplit:                  0.1,
			MinContextUsefulImprovementFull:    -1,
			MinContextUsefulImprovementSplit:   -1,
			MinContextMarginImprovementFull:    -1,
			MinContextMarginImprovementSplit:   -1,
			MinContextToxicImprovementFull:     -1,
			MinContextToxicImprovementSplit:    -1,
			MaxBucketShareOfLocalBaselineFull:  1,
			MaxBucketShareOfLocalBaselineSplit: 1,
		}
	}
	defer func() { futuresDerivativesContextAuditConfigForRun = oldCfg }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-out-dir", defaultOutDir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_derivatives_context_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write derivatives context audit artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "deriv-context-audit")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-context-audit", "-out-dir", outDir}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json", "summary.csv", "summary.json", "trades.json",
		"futures_derivatives_context_sources.csv", "futures_derivatives_context_sources.json",
		"futures_derivatives_context_coverage.csv", "futures_derivatives_context_coverage.json",
		"futures_derivatives_context_basis_features.csv", "futures_derivatives_context_basis_features.json",
		"futures_derivatives_context_local_states.csv", "futures_derivatives_context_local_states.json",
		"futures_derivatives_context_labels.csv", "futures_derivatives_context_labels.json",
		"futures_derivatives_context_cohorts.csv", "futures_derivatives_context_cohorts.json",
		"futures_derivatives_context_rankings.csv", "futures_derivatives_context_rankings.json",
		"futures_derivatives_context_missingness.csv", "futures_derivatives_context_missingness.json",
		"futures_derivatives_context_summary.csv", "futures_derivatives_context_summary.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected derivatives context audit artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("derivatives context audit must remain zero-trade, got %d", len(trades))
	}
	conflictErr := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-context-audit", "-futures-derivatives-context-source-audit", "-out-dir", filepath.Join(dir, "conflict")})
	if conflictErr == nil {
		t.Fatalf("expected error combining derivatives context audit with another audit flag")
	}
}

func TestRunWithArgsFuturesDerivativesNoTradeFilterPremiseAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_no_trade_filter_cli.csv", 96)
	oldCfg := futuresDerivativesNoTradeFilterPremiseAuditConfigForRun
	futuresDerivativesNoTradeFilterPremiseAuditConfigForRun = func() lab.FuturesDerivativesNoTradeFilterPremiseAuditConfig {
		cfg := lab.DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
		cfg.RescueClosedFamily = true
		return cfg
	}
	defer func() { futuresDerivativesNoTradeFilterPremiseAuditConfigForRun = oldCfg }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-out-dir", defaultOutDir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "futures_derivatives_no_trade_filter_premise_summary.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write derivatives no-trade filter premise artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "deriv-no-trade-filter")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-no-trade-filter-premise-audit", "-out-dir", outDir}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json", "summary.csv", "summary.json", "trades.json",
		"futures_derivatives_no_trade_filter_premise_sources.csv", "futures_derivatives_no_trade_filter_premise_sources.json",
		"futures_derivatives_no_trade_filter_premise_coverage.csv", "futures_derivatives_no_trade_filter_premise_coverage.json",
		"futures_derivatives_no_trade_filter_premise_filter_definitions.csv", "futures_derivatives_no_trade_filter_premise_filter_definitions.json",
		"futures_derivatives_no_trade_filter_premise_exact_candidates.csv", "futures_derivatives_no_trade_filter_premise_exact_candidates.json",
		"futures_derivatives_no_trade_filter_premise_canonical_union.csv", "futures_derivatives_no_trade_filter_premise_canonical_union.json",
		"futures_derivatives_no_trade_filter_premise_overlap.csv", "futures_derivatives_no_trade_filter_premise_overlap.json",
		"futures_derivatives_no_trade_filter_premise_veto_candidates.csv", "futures_derivatives_no_trade_filter_premise_veto_candidates.json",
		"futures_derivatives_no_trade_filter_premise_collateral_damage.csv", "futures_derivatives_no_trade_filter_premise_collateral_damage.json",
		"futures_derivatives_no_trade_filter_premise_missingness.csv", "futures_derivatives_no_trade_filter_premise_missingness.json",
		"futures_derivatives_no_trade_filter_premise_summary.csv", "futures_derivatives_no_trade_filter_premise_summary.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected derivatives no-trade filter premise artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("derivatives no-trade filter premise audit must remain zero-trade, got %d", len(trades))
	}
	conflictErr := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-derivatives-no-trade-filter-premise-audit", "-futures-derivatives-context-audit", "-out-dir", filepath.Join(dir, "conflict")})
	if conflictErr == nil {
		t.Fatalf("expected error combining derivatives no-trade filter premise audit with another audit flag")
	}
}

func TestRunWithArgsFuturesBTC15MPostCompressionDirectionalExpansionAuditFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_post_compression_cli.csv", 96)
	oldCfg := futuresBTC15MPostCompressionDirectionalExpansionAuditConfigForRun
	futuresBTC15MPostCompressionDirectionalExpansionAuditConfigForRun = func() lab.FuturesBTC15MPostCompressionDirectionalExpansionAuditConfig {
		cfg := lab.DefaultFuturesBTC15MPostCompressionDirectionalExpansionAuditConfig()
		cfg.ApprovedSourcePath = futuresPath
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.CompressionLookbacks = []int{2}
		cfg.CompressionPercentiles = []float64{0.8}
		cfg.PercentileReferenceBars = 2
		cfg.BreakoutATRMultiples = []float64{0.1}
		cfg.VolumeModes = []string{lab.BTC15MPostCompressionVolumeNone}
		cfg.VolumeLookbackBars = 2
		cfg.ATRPeriod = 2
		cfg.HorizonsBars = []int{1}
		cfg.MinFullDedupCandidates = 1
		cfg.MinSplitDedupCandidates = 1
		return cfg
	}
	defer func() { futuresBTC15MPostCompressionDirectionalExpansionAuditConfigForRun = oldCfg }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-out-dir", defaultOutDir}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "btc_15m_post_compression_directional_expansion_sources.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write post-compression audit artifacts, stat err=%v", err)
	}

	outDir := filepath.Join(dir, "post-compression")
	if err := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-btc-15m-post-compression-directional-expansion-audit", "-out-dir", outDir}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"source_manifest.json", "summary.csv", "summary.json", "trades.json",
		"btc_15m_post_compression_directional_expansion_sources.csv", "btc_15m_post_compression_directional_expansion_sources.json",
		"btc_15m_post_compression_directional_expansion_resample_coverage.csv", "btc_15m_post_compression_directional_expansion_resample_coverage.json",
		"btc_15m_post_compression_directional_expansion_parameter_cells.csv", "btc_15m_post_compression_directional_expansion_parameter_cells.json",
		"btc_15m_post_compression_directional_expansion_candidates.csv", "btc_15m_post_compression_directional_expansion_candidates.json",
		"btc_15m_post_compression_directional_expansion_dedup_events.csv", "btc_15m_post_compression_directional_expansion_dedup_events.json",
		"btc_15m_post_compression_directional_expansion_baseline.csv", "btc_15m_post_compression_directional_expansion_baseline.json",
		"btc_15m_post_compression_directional_expansion_split_summary.csv", "btc_15m_post_compression_directional_expansion_split_summary.json",
		"btc_15m_post_compression_directional_expansion_adjacency.csv", "btc_15m_post_compression_directional_expansion_adjacency.json",
		"btc_15m_post_compression_directional_expansion_missingness.csv", "btc_15m_post_compression_directional_expansion_missingness.json",
		"btc_15m_post_compression_directional_expansion_falsification.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected post-compression audit artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}
	if len(trades) != 0 {
		t.Fatalf("post-compression audit must remain zero-trade, got %d", len(trades))
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_post_compression.csv")
	spotErr := runWithArgs([]string{"-csv", spotPath, "-source-product", lab.SourceProductBinanceSpot, "-allow-spot-comparison", "-futures-btc-15m-post-compression-directional-expansion-audit", "-out-dir", filepath.Join(dir, "spot")})
	if spotErr == nil || !strings.Contains(spotErr.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected post-compression futures-source error, got %v", spotErr)
	}

	conflictErr := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-btc-15m-post-compression-directional-expansion-audit", "-futures-derivatives-no-trade-filter-premise-audit", "-out-dir", filepath.Join(dir, "conflict")})
	if conflictErr == nil || !strings.Contains(conflictErr.Error(), "cannot be combined") {
		t.Fatalf("expected error combining post-compression audit with another audit flag, got %v", conflictErr)
	}
}

func TestRunWithArgsFuturesBTC15MPostCompressionFixedBacktestFlagWritesArtifactsAndRejectsConflicts(t *testing.T) {
	dir := t.TempDir()
	futuresPath := writeCLITestCSVN(t, dir, "btcusdt_futures_um_5m_post_compression_backtest_cli.csv", 180)
	oldCfg := futuresBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfigForRun
	futuresBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfigForRun = func() lab.FuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig {
		cfg := lab.DefaultFuturesBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfig()
		cfg.ApprovedSourcePath = futuresPath
		cfg.SkipSourceFactCheck = true
		cfg.SkipCoverageCountCheck = true
		cfg.LookbackBars = 2
		cfg.CompressionPercentile = 0.8
		cfg.PercentileReferenceBars = 2
		cfg.BreakoutATRMultiple = 0.1
		cfg.ATRPeriod = 2
		cfg.MaxHoldBars = 2
		cfg.SkipCandidateIdentityGate = true
		cfg.MinFullTrades = 1
		cfg.MinSplitTrades = 1
		return cfg
	}
	defer func() { futuresBTC15MPostCompressionL192Q20M020NoneLongH48BacktestConfigForRun = oldCfg }()

	defaultOutDir := filepath.Join(dir, "default")
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-out-dir", defaultOutDir,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(defaultOutDir, "btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.csv")); !os.IsNotExist(err) {
		t.Fatalf("default run should not write fixed post-compression backtest artifacts, stat err=%v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := runWithArgs([]string{
		"-csv", futuresPath,
		"-source-product", lab.SourceProductBinanceUSDMFutures,
		"-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest",
	}); err != nil {
		_ = os.Chdir(cwd)
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	outDir := filepath.Join(dir, "results", "futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest")
	for _, name := range []string{
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_sources.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_sources.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_resample_coverage.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_resample_coverage.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_skips.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_skips.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_trades.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_trades.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_summary.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_summary.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_cost_stress.csv", "btc_15m_post_compression_l192_q20_m020_none_long_h48_cost_stress.json",
		"btc_15m_post_compression_l192_q20_m020_none_long_h48_falsification.json",
		"source_manifest.json", "summary.json", "summary.csv", "trades.json",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected fixed post-compression backtest artifact %s: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "trades.json"))
	if err != nil {
		t.Fatal(err)
	}
	var trades []lab.Trade
	if err := json.Unmarshal(data, &trades); err != nil {
		t.Fatal(err)
	}

	spotPath := writeCLITestCSV(t, dir, "btcusdt_spot_5m_post_compression_backtest.csv")
	spotErr := runWithArgs([]string{"-csv", spotPath, "-source-product", lab.SourceProductBinanceSpot, "-allow-spot-comparison", "-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest", "-out-dir", filepath.Join(dir, "spot")})
	if spotErr == nil || !strings.Contains(spotErr.Error(), "requires Binance USDT-M futures source") {
		t.Fatalf("expected fixed post-compression backtest futures-source error, got %v", spotErr)
	}

	conflictErr := runWithArgs([]string{"-csv", futuresPath, "-source-product", lab.SourceProductBinanceUSDMFutures, "-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest", "-futures-btc-15m-post-compression-directional-expansion-audit", "-out-dir", filepath.Join(dir, "conflict")})
	if conflictErr == nil || !strings.Contains(conflictErr.Error(), "cannot be combined") {
		t.Fatalf("expected error combining fixed post-compression backtest with another audit flag, got %v", conflictErr)
	}
}

func writeCLITestCSV(t *testing.T, dir string, name string) string {
	return writeCLITestCSVN(t, dir, name, 2)
}

func writeCLITestCSVN(t *testing.T, dir string, name string, rows int) string {
	t.Helper()
	path := filepath.Join(dir, name)
	data := "open_time,open,high,low,close,volume,close_time\n"
	for i := 0; i < rows; i++ {
		open := int64(1609459200000 + i*300000)
		closeTime := open + 299999
		openPrice := 100 + i
		data += strconv.FormatInt(open, 10) + "," +
			strconv.Itoa(openPrice) + "," +
			strconv.Itoa(openPrice+10) + "," +
			strconv.Itoa(openPrice-10) + "," +
			strconv.Itoa(openPrice+5) + "," +
			strconv.Itoa(i+1) + "," +
			strconv.FormatInt(closeTime, 10) + "\n"
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readCLITestManifest(t *testing.T, outDir string) lab.SourceManifest {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(outDir, "source_manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest lab.SourceManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}
