package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"range-strategy-lab/internal/lab"
)

const sessionOpeningRangeExpansionFlagName = "backtest-first-btc-15m-session-opening-range-expansion-v1"

var sessionOpeningRangeExpansionConfigForRun = lab.DefaultBacktestFirstBTC15MSessionOpeningRangeExpansionConfig

func init() {
	if !sessionOpeningRangeExpansionFlagPresent(os.Args[1:]) {
		return
	}
	if err := runSessionOpeningRangeExpansionWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func sessionOpeningRangeExpansionFlagPresent(args []string) bool {
	prefix := "-" + sessionOpeningRangeExpansionFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}

func runSessionOpeningRangeExpansionWithArgs(args []string) error {
	cfg := sessionOpeningRangeExpansionConfigForRun()
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	csvPath := fs.String("csv", cfg.ApprovedSourcePath, "5m BTCUSDT candle CSV")
	outDir := fs.String("out-dir", "results/backtest-first-btc-15m-session-opening-range-expansion-v1", "output directory")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	runFlag := fs.Bool(sessionOpeningRangeExpansionFlagName, false, "run fixed BTCUSDT 15m session opening-range expansion baseline")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*runFlag {
		return nil
	}
	if !samePathForCLI(*csvPath, cfg.ApprovedSourcePath) {
		return fmt.Errorf("-%s uses only the approved BTCUSDT futures CSV %s", sessionOpeningRangeExpansionFlagName, cfg.ApprovedSourcePath)
	}
	candles, manifest, err := lab.LoadResearchSourceCSV(*csvPath, lab.SourceValidationOptions{Product: lab.SourceProductBinanceUSDMFutures})
	if err != nil {
		return err
	}
	btCfg := lab.BacktestConfig{
		StartBalance:   *startBalance,
		RiskPct:        *riskPct,
		MaxNotionalPct: *maxNotionalPct,
		FeePct:         *feePct,
		SlippagePct:    *slippagePct,
		MaxHoldBars:    cfg.MaxHoldBars,
	}
	result, err := lab.RunBacktestFirstBTC15MSessionOpeningRangeExpansion(candles, manifest, cfg, btCfg, lab.DefaultSplits())
	if err != nil {
		return err
	}
	if err := writeSessionOpeningRangeExpansionOutputs(*outDir, manifest, result); err != nil {
		return err
	}
	fmt.Printf("backtest_first_btc_15m_session_opening_range_expansion session_range_rows=%d signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n", len(result.SessionRangeRows), len(result.SignalRows), len(result.Trades), len(result.SummaryRows), result.StopState)
	return nil
}

func samePathForCLI(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func writeSessionOpeningRangeExpansionOutputs(outDir string, manifest lab.SourceManifest, result lab.BacktestFirstBTC15MSessionOpeningRangeExpansionResult) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "source_manifest.json"), manifest); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "summary.json"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(outDir, "summary.csv"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "trades.json"), result.Trades); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_sources.json"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_coverage.json"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_session_ranges.json"), result.SessionRangeRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_signals.json"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_skips.json"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_trades.json"), result.TradeRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_summary.json"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_falsification.json"), []lab.BTC15MSessionOpeningRangeExpansionFalsification{result.Falsification}); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_sources.csv"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_coverage.csv"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_session_ranges.csv"), result.SessionRangeRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_signals.csv"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_skips.csv"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_trades.csv"), result.TradeRows); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_summary.csv"), result.SummaryRows); err != nil {
		return err
	}
	return writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_session_opening_range_expansion_falsification.csv"), []lab.BTC15MSessionOpeningRangeExpansionFalsification{result.Falsification})
}
