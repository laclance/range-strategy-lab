package main

import (
	"path/filepath"

	"range-strategy-lab/internal/lab"
)

func writeEdgeFadeArtifacts(outDir string, result lab.BacktestFirstBTC15MRangeEdgeExhaustionFadeResult) error {
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_sources.json"), result.SourceRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_coverage.json"), result.CoverageRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_signals.json"), result.SignalRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_skips.json"), result.SkipRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_trades.json"), result.TradeRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_summary.json"), result.SummaryRows); err != nil { return err }
	if err := writeJSON(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_falsification.json"), []lab.BTC15MRangeEdgeExhaustionFadeFalsification{result.Falsification}); err != nil { return err }
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_sources.csv"), result.SourceRows); err != nil { return err }
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_coverage.csv"), result.CoverageRows); err != nil { return err }
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_signals.csv"), result.SignalRows); err != nil { return err }
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_skips.csv"), result.SkipRows); err != nil { return err }
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_trades.csv"), result.TradeRows); err != nil { return err }
	if err := writeSummaryCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_summary.csv"), result.SummaryRows); err != nil { return err }
	return writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_range_edge_exhaustion_fade_falsification.csv"), []lab.BTC15MRangeEdgeExhaustionFadeFalsification{result.Falsification})
}
