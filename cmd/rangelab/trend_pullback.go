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

const trendPullbackFlagName = "backtest-first-btc-15m-trend-pullback-continuation-v1"

func init() {
	if !trendPullbackFlagPresent(os.Args[1:]) {
		return
	}
	if err := runTrendPullbackWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func trendPullbackFlagPresent(args []string) bool {
	prefix := "-" + trendPullbackFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}

func runTrendPullbackWithArgs(args []string) error {
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	csvPath := fs.String("csv", defaultCSVPath, "5m BTCUSDT candle CSV")
	outDir := fs.String("out-dir", "results/backtest-first-btc-15m-trend-pullback-continuation-v1", "output directory")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	runFlag := fs.Bool(trendPullbackFlagName, false, "run fixed BTCUSDT 15m trend-pullback continuation baseline")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*runFlag {
		return nil
	}
	if *csvPath != defaultCSVPath {
		return fmt.Errorf("-%s uses only the approved default BTCUSDT futures CSV", trendPullbackFlagName)
	}
	candles, manifest, err := lab.LoadResearchSourceCSV(*csvPath, lab.SourceValidationOptions{Product: lab.SourceProductBinanceUSDMFutures})
	if err != nil {
		return err
	}
	strategyCfg := lab.DefaultBacktestFirstBTC15MTrendPullbackContinuationConfig()
	btCfg := lab.BacktestConfig{StartBalance: *startBalance, RiskPct: *riskPct, MaxNotionalPct: *maxNotionalPct, FeePct: *feePct, SlippagePct: *slippagePct, MaxHoldBars: strategyCfg.MaxHoldBars}
	result, err := lab.RunBacktestFirstBTC15MTrendPullbackContinuation(candles, manifest, strategyCfg, btCfg, lab.DefaultSplits())
	if err != nil {
		return err
	}
	if err := writeTrendPullbackOutputs(*outDir, manifest, result, *startBalance); err != nil {
		return err
	}
	fmt.Printf("backtest_first_btc_15m_trend_pullback_continuation signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n", len(result.SignalRows), len(result.Trades), len(result.SummaryRows), result.StopState)
	return nil
}

func writeTrendPullbackOutputs(outDir string, manifest lab.SourceManifest, result lab.BacktestFirstBTC15MTrendPullbackContinuationResult, startBalance float64) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	summary := lab.SummarizeSplits(result.Trades, startBalance, lab.DefaultSplits())
	if err := writeJSON(filepath.Join(outDir, "source_manifest.json"), manifest); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "summary.json"), summary); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(outDir, "summary.csv"), summary); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "trades.json"), result.Trades); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_sources.json"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_coverage.json"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_signals.json"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_skips.json"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_trades.json"), result.TradeRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_summary.json"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_falsification.json"), []lab.BTC15MTrendPullbackContinuationFalsification{result.Falsification}); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_sources.csv"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_coverage.csv"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_signals.csv"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_skips.csv"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_trades.csv"), result.TradeRows); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_summary.csv"), result.SummaryRows); err != nil {
		return err
	}
	return writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_trend_pullback_continuation_falsification.csv"), []lab.BTC15MTrendPullbackContinuationFalsification{result.Falsification})
}
