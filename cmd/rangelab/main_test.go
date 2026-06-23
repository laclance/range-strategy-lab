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
