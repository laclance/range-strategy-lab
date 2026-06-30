package main

import (
	"fmt"
	"os"

	"range-strategy-lab/internal/lab"
)

func runRangeWorkbench(csvPath string, outDir string, runID string, startBalance float64, riskPct float64, maxNotionalPct float64, feePct float64, slippagePct float64) error {
	if _, err := os.Stat(outDir); err == nil {
		return fmt.Errorf("range workbench output directory already exists: %s", outDir)
	} else if !os.IsNotExist(err) {
		return err
	}
	candles, manifest, err := lab.LoadResearchSourceCSV(csvPath, lab.SourceValidationOptions{Product: lab.SourceProductBinanceUSDMFutures})
	if err != nil {
		return err
	}
	strategyCfg := lab.DefaultRangeOptimizationWorkbenchConfig()
	btCfg := lab.BacktestConfig{StartBalance: startBalance, RiskPct: riskPct, MaxNotionalPct: maxNotionalPct, FeePct: feePct, SlippagePct: slippagePct}
	result, err := lab.RunRangeOptimizationWorkbench(candles, manifest, strategyCfg, btCfg, lab.DefaultSplits(), runID)
	if err != nil {
		return err
	}
	if err := writeRangeWorkbenchOutputs(outDir, manifest, result); err != nil {
		return err
	}
	fmt.Printf("range_optimization_workbench run_id=%s trials=%d passing_candidates=%d selected=%s stop_state=%s\n", result.RunID, len(result.TrialRows), len(result.TopCandidates), result.Falsification.SelectedTrialID, result.StopState)
	return nil
}
