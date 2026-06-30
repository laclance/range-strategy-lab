package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const rangeWorkbenchFlagName = "range-optimization-workbench-v1"

func init() {
	if !rangeWorkbenchFlagPresent(os.Args[1:]) {
		return
	}
	if err := runRangeWorkbenchWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func rangeWorkbenchFlagPresent(args []string) bool {
	prefix := "-" + rangeWorkbenchFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}

func runRangeWorkbenchWithArgs(args []string) error {
	fs := flag.NewFlagSet("rangelab", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	csvPath := fs.String("csv", defaultCSVPath, "5m BTCUSDT candle CSV")
	outDir := fs.String("out-dir", "", "immutable range workbench run output directory")
	runID := fs.String("run-id", "", "unique range workbench run id")
	startBalance := fs.Float64("start-balance", 1000, "starting balance")
	riskPct := fs.Float64("risk-pct", 0.01, "fraction of equity risked at stop")
	maxNotionalPct := fs.Float64("max-notional-pct", 1.0, "maximum entry notional as fraction of equity")
	feePct := fs.Float64("fee-pct", 0.0004, "fee fraction per side")
	slippagePct := fs.Float64("slippage-pct", 0.000116, "slippage fraction per side")
	runFlag := fs.Bool(rangeWorkbenchFlagName, false, "run bounded offline range optimization workbench")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*runFlag {
		return nil
	}
	if *csvPath != defaultCSVPath {
		return fmt.Errorf("-%s uses only the approved default BTCUSDT futures CSV", rangeWorkbenchFlagName)
	}
	if *runID == "" {
		*runID = time.Now().UTC().Format("20060102T150405Z")
	}
	if *outDir == "" {
		*outDir = filepath.Join("results", "range-optimization-workbench-v1", "runs", *runID)
	}
	return runRangeWorkbench(*csvPath, *outDir, *runID, *startBalance, *riskPct, *maxNotionalPct, *feePct, *slippagePct)
}
