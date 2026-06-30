package main

import (
	"os"
	"path/filepath"

	"range-strategy-lab/internal/lab"
)

func writeEdgeFadeOutputs(outDir string, manifest lab.SourceManifest, result lab.BacktestFirstBTC15MRangeEdgeExhaustionFadeResult, startBalance float64) error {
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
	if err := writeEdgeFadeArtifacts(outDir, result); err != nil {
		return err
	}
	return nil
}
