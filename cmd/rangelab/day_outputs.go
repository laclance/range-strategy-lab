package main

import (
	"path/filepath"

	"range-strategy-lab/internal/lab"
)

func writeBTC15MDayArtifacts(outDir string, result lab.BacktestFirstBTC15MPreviousDayRangeReversionResult) error {
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_sources.json"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_coverage.json"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_signals.json"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_skips.json"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_trades.json"), result.TradeRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_summary.json"), result.SummaryRows); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_falsification.json"), []lab.BTC15MPreviousDayRangeReversionFalsification{result.Falsification}); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_sources.csv"), result.SourceRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_coverage.csv"), result.CoverageRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_signals.csv"), result.SignalRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_skips.csv"), result.SkipRows); err != nil {
		return err
	}
	if err := writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_trades.csv"), result.TradeRows); err != nil {
		return err
	}
	if err := writeSummaryCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_summary.csv"), result.SummaryRows); err != nil {
		return err
	}
	return writeJSONTaggedCSV(filepath.Join(outDir, "btc_15m_previous_day_range_reversion_falsification.csv"), []lab.BTC15MPreviousDayRangeReversionFalsification{result.Falsification})
}
