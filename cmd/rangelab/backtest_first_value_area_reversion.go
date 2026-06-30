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

const backtestFirstValueAreaFlagName = "backtest-first-btc-5m-rolling-value-area-reversion-v1"

func init() {
	if !backtestFirstValueAreaFlagPresent(os.Args[1:]) {
		return
	}
	if err := runBacktestFirstValueAreaReversionWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func backtestFirstValueAreaFlagPresent(args []string) bool {
	prefix := "-" + backtestFirstValueAreaFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}

func runBacktestFirstValueAreaReversionWithArgs(args []string) error {
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	csvPath := fs.String("csv", defaultCSVPath, "5m BTCUSDT candle CSV")
	sourceProduct := fs.String("source-product", "", "source product: binance-usdm-futures or binance-spot")
	allowSpotComparison := fs.Bool("allow-spot-comparison", false, "allow an explicitly labeled Binance spot comparison run")
	outDir := fs.String("out-dir", "results/backtest-first-btc-5m-rolling-value-area-reversion-v1", "output directory")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	runFlag := fs.Bool(backtestFirstValueAreaFlagName, false, "run offline BTCUSDT 5m rolling value-area reversion baseline")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*runFlag {
		return nil
	}

	sourceProductWasSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "source-product" {
			sourceProductWasSet = true
		}
	})
	product := *sourceProduct
	if *csvPath == defaultCSVPath && !sourceProductWasSet {
		product = lab.SourceProductBinanceUSDMFutures
	}
	if *csvPath != defaultCSVPath && !sourceProductWasSet {
		return fmt.Errorf("non-default -csv path requires explicit -source-product=%s or -source-product=%s", lab.SourceProductBinanceUSDMFutures, lab.SourceProductBinanceSpot)
	}

	candles, sourceManifest, err := lab.LoadResearchSourceCSV(*csvPath, lab.SourceValidationOptions{
		Product:             product,
		AllowSpotComparison: *allowSpotComparison,
	})
	if err != nil {
		return err
	}
	if sourceManifest.ComparisonOnly || sourceManifest.Product != "Binance USDT-M futures" {
		return fmt.Errorf("-%s requires Binance USDT-M futures source; got product=%q comparison_only=%t", backtestFirstValueAreaFlagName, sourceManifest.Product, sourceManifest.ComparisonOnly)
	}

	strategyCfg := lab.DefaultBacktestFirstBTC5MRollingValueAreaReversionConfig()
	btCfg := lab.BacktestConfig{
		StartBalance:   *startBalance,
		RiskPct:        *riskPct,
		MaxNotionalPct: *maxNotionalPct,
		FeePct:         *feePct,
		SlippagePct:    *slippagePct,
		MaxHoldBars:    strategyCfg.MaxHoldBars,
	}
	result, err := lab.RunBacktestFirstBTC5MRollingValueAreaReversion(candles, sourceManifest, strategyCfg, btCfg, lab.DefaultSplits())
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return err
	}

	commonSummary := lab.SummarizeSplits(result.Trades, *startBalance, lab.DefaultSplits())
	if err := writeJSON(filepath.Join(*outDir, "source_manifest.json"), sourceManifest); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "summary.json"), commonSummary); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(*outDir, "summary.csv"), commonSummary); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(*outDir, "trades.json"), result.Trades); err != nil {
		return err
	}

	artifacts := []struct {
		name string
		rows any
	}{
		{"btc_5m_rolling_value_area_reversion_sources", result.SourceRows},
		{"btc_5m_rolling_value_area_reversion_signals", result.SignalRows},
		{"btc_5m_rolling_value_area_reversion_skips", result.SkipRows},
		{"btc_5m_rolling_value_area_reversion_trades", result.TradeRows},
		{"btc_5m_rolling_value_area_reversion_summary", result.SummaryRows},
		{"btc_5m_rolling_value_area_reversion_falsification", []lab.BTC5MRollingValueAreaReversionFalsification{result.Falsification}},
	}
	for _, artifact := range artifacts {
		if err := writeJSON(filepath.Join(*outDir, artifact.name+".json"), artifact.rows); err != nil {
			return err
		}
	}
	if err := writeJSONTaggedCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_sources.csv"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_signals.csv"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_skips.csv"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_trades.csv"), result.TradeRows); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_summary.csv"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(*outDir, "btc_5m_rolling_value_area_reversion_falsification.csv"), []lab.BTC5MRollingValueAreaReversionFalsification{result.Falsification}); err != nil {
		return err
	}

	fmt.Printf("backtest_first_btc_5m_rolling_value_area_reversion signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n", len(result.SignalRows), len(result.Trades), len(result.SummaryRows), result.StopState)
	return nil
}
