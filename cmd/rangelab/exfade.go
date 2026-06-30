package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"range-strategy-lab/internal/lab"
)

const edgeFadeFlagName = "backtest-first-btc-15m-range-edge-exhaustion-fade-v1"

func init() {
	if !edgeFadeFlagPresent(os.Args[1:]) {
		return
	}
	if err := runEdgeFadeWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func edgeFadeFlagPresent(args []string) bool {
	prefix := "-" + edgeFadeFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}

func runEdgeFadeWithArgs(args []string) error {
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	csvPath := fs.String("csv", defaultCSVPath, "5m BTCUSDT candle CSV")
	outDir := fs.String("out-dir", "results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1", "output directory")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	runFlag := fs.Bool(edgeFadeFlagName, false, "run fixed range-edge exhaustion fade baseline")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*runFlag {
		return nil
	}
	if *csvPath != defaultCSVPath {
		return fmt.Errorf("-%s uses only the approved default BTCUSDT futures CSV", edgeFadeFlagName)
	}
	candles, manifest, err := lab.LoadResearchSourceCSV(*csvPath, lab.SourceValidationOptions{Product: lab.SourceProductBinanceUSDMFutures})
	if err != nil {
		return err
	}
	strategyCfg := lab.DefaultBacktestFirstBTC15MRangeEdgeExhaustionFadeConfig()
	btCfg := lab.BacktestConfig{StartBalance: *startBalance, RiskPct: *riskPct, MaxNotionalPct: *maxNotionalPct, FeePct: *feePct, SlippagePct: *slippagePct, MaxHoldBars: strategyCfg.MaxHoldBars}
	result, err := lab.RunBacktestFirstBTC15MRangeEdgeExhaustionFade(candles, manifest, strategyCfg, btCfg, lab.DefaultSplits())
	if err != nil {
		return err
	}
	if err := writeEdgeFadeOutputs(*outDir, manifest, result, *startBalance); err != nil {
		return err
	}
	fmt.Printf("backtest_first_btc_15m_range_edge_exhaustion_fade signal_rows=%d trades=%d summary_rows=%d stop_state=%s\n", len(result.SignalRows), len(result.Trades), len(result.SummaryRows), result.StopState)
	return nil
}
