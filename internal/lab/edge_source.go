package lab

import (
	"path/filepath"
	"strings"
)

func btc15MEdgeSourceRow(manifest SourceManifest, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig) BTC15MRangeEdgeExhaustionFadeSourceRow {
	row := BTC15MRangeEdgeExhaustionFadeSourceRow{BacktestName: BacktestFirstBTC15MRangeEdgeExhaustionFadeName, CandidateID: BTC15MRangeEdgeExhaustionFadeCandidateID, Path: manifest.Path, ApprovedPath: cfg.ApprovedSourcePath, Product: manifest.Product, Symbol: manifest.Symbol, Interval: manifest.Interval, RowCount: manifest.RowCount, ExpectedRowCount: cfg.ExpectedSourceRows, FirstOpenTime: manifest.FirstOpenTime, ExpectedFirstOpenTime: cfg.ExpectedFirstOpenTime, LastOpenTime: manifest.LastOpenTime, ExpectedLastOpenTime: cfg.ExpectedLastOpenTime, GapCount: manifest.GapCount, DuplicateCount: manifest.DuplicateCount, ZeroVolumeCount: manifest.ZeroVolumeCount, ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount, ComparisonOnly: manifest.ComparisonOnly, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	failures := []string{}
	if manifest.ValidationStatus != "accepted" || manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly || manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" || filepath.Clean(manifest.Path) != filepath.Clean(cfg.ApprovedSourcePath) || manifest.RowCount != cfg.ExpectedSourceRows || manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime || manifest.LastOpenTime != cfg.ExpectedLastOpenTime || manifest.GapCount != 0 || manifest.DuplicateCount != 0 || manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
		failures = append(failures, "source contract mismatch")
	}
	row.SourceFactsPass = len(failures) == 0
	if !row.SourceFactsPass {
		row.ValidationStatus = "rejected"
		row.ValidationError = strings.Join(failures, "; ")
	}
	return row
}
