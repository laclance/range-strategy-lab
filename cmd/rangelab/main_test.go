package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func writeCLITestCSV(t *testing.T, dir string, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	data := "open_time,open,high,low,close,volume,close_time\n" +
		"1609459200000,100,110,90,105,1,1609459499999\n" +
		"1609459500000,105,112,100,108,2,1609459799999\n"
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
