package main

import (
	"path/filepath"

	"range-strategy-lab/internal/lab"
)

type rangeWorkbenchLatestRun struct {
	RunID     string `json:"run_id"`
	OutDir    string `json:"out_dir"`
	StopState string `json:"stop_state"`
}

func writeRangeWorkbenchOutputs(outDir string, manifest lab.SourceManifest, result lab.RangeOptimizationWorkbenchResult) error {
	parentDir := filepath.Dir(filepath.Dir(outDir))
	if err := writeJSON(filepath.Join(outDir, "source_manifest.json"), manifest); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "source_contract.json"), result.SourceRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "coverage.json"), result.CoverageRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "optimization_grid.json"), result.GridRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "trial_results.json"), result.TrialRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "trial_summary.json"), result.TrialSummaryRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "top_candidates.json"), result.TopCandidates); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "rejected_candidates.json"), result.RejectedCandidates); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "robustness_summary.json"), result.RobustnessSummary); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "falsification.json"), result.Falsification); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "source_contract.csv"), result.SourceRows); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "coverage.csv"), result.CoverageRows); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "trial_results.csv"), result.TrialRows); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "trial_summary.csv"), result.TrialSummaryRows); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "top_candidates.csv"), result.TopCandidates); err != nil { return err }
	if err := writeRangeWorkbenchRowsCSV(filepath.Join(outDir, "rejected_candidates.csv"), result.RejectedCandidates); err != nil { return err }
	latest := rangeWorkbenchLatestRun{RunID: result.RunID, OutDir: outDir, StopState: result.StopState}
	return writeJSON(filepath.Join(parentDir, "latest_run.json"), latest)
}
