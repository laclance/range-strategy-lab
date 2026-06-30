package lab

import (
	"path/filepath"
	"strings"
	"time"
)

func (cfg BacktestFirstBTC5MRollingValueAreaReversionConfig) withDefaults() BacktestFirstBTC5MRollingValueAreaReversionConfig {
	defaults := DefaultBacktestFirstBTC5MRollingValueAreaReversionConfig()
	if cfg.ApprovedSourcePath == "" {
		cfg.ApprovedSourcePath = defaults.ApprovedSourcePath
	}
	if cfg.ExpectedSourceRows == 0 {
		cfg.ExpectedSourceRows = defaults.ExpectedSourceRows
	}
	if cfg.ExpectedFirstOpenTime == "" {
		cfg.ExpectedFirstOpenTime = defaults.ExpectedFirstOpenTime
	}
	if cfg.ExpectedLastOpenTime == "" {
		cfg.ExpectedLastOpenTime = defaults.ExpectedLastOpenTime
	}
	if cfg.ExpectedZeroVolumeCount == 0 {
		cfg.ExpectedZeroVolumeCount = defaults.ExpectedZeroVolumeCount
	}
	if cfg.LookbackBars == 0 {
		cfg.LookbackBars = defaults.LookbackBars
	}
	if cfg.ATRPeriod == 0 {
		cfg.ATRPeriod = defaults.ATRPeriod
	}
	if cfg.MinRangeATRs == 0 {
		cfg.MinRangeATRs = defaults.MinRangeATRs
	}
	if cfg.OuterZonePct == 0 {
		cfg.OuterZonePct = defaults.OuterZonePct
	}
	if cfg.VWAPDistanceRangePct == 0 {
		cfg.VWAPDistanceRangePct = defaults.VWAPDistanceRangePct
	}
	if cfg.StopATRMultiple == 0 {
		cfg.StopATRMultiple = defaults.StopATRMultiple
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = defaults.MaxHoldBars
	}
	if cfg.MinFullTrades == 0 {
		cfg.MinFullTrades = defaults.MinFullTrades
	}
	if cfg.MinSplitTrades == 0 {
		cfg.MinSplitTrades = defaults.MinSplitTrades
	}
	if cfg.FullMaxDrawdownLimit == 0 {
		cfg.FullMaxDrawdownLimit = defaults.FullMaxDrawdownLimit
	}
	if cfg.SplitMaxDrawdownLimit == 0 {
		cfg.SplitMaxDrawdownLimit = defaults.SplitMaxDrawdownLimit
	}
	return cfg
}

func btc5MValueAreaSourceRow(manifest SourceManifest, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig) BTC5MRollingValueAreaReversionSourceRow {
	row := BTC5MRollingValueAreaReversionSourceRow{BacktestName: BacktestFirstBTC5MRollingValueAreaReversionName, CandidateID: BTC5MRollingValueAreaReversionCandidateID, Path: manifest.Path, ApprovedPath: cfg.ApprovedSourcePath, Product: manifest.Product, Symbol: manifest.Symbol, Interval: manifest.Interval, RowCount: manifest.RowCount, ExpectedRowCount: cfg.ExpectedSourceRows, FirstOpenTime: manifest.FirstOpenTime, ExpectedFirstOpenTime: cfg.ExpectedFirstOpenTime, LastOpenTime: manifest.LastOpenTime, ExpectedLastOpenTime: cfg.ExpectedLastOpenTime, GapCount: manifest.GapCount, DuplicateCount: manifest.DuplicateCount, ZeroVolumeCount: manifest.ZeroVolumeCount, ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount, ComparisonOnly: manifest.ComparisonOnly, ClosedCandleOnly: true, Native5MOnly: true, ValidationStatus: "accepted"}
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

func btc5MValueAreaTradeRows(trades []Trade, splits []Split) []BTC5MRollingValueAreaReversionTradeRow {
	rows := make([]BTC5MRollingValueAreaReversionTradeRow, 0, len(trades))
	for _, tr := range trades {
		exitTime, _ := time.Parse(timeLayout, tr.ExitTime)
		rows = append(rows, BTC5MRollingValueAreaReversionTradeRow{SignalID: tr.Signal, CloseSplit: btc5MValueAreaSplit(exitTime, splits), Side: tr.Side, EntryTime: tr.EntryTime, ExitTime: tr.ExitTime, GrossPnL: tr.GrossPnL, NetPnL: tr.NetPnL, Fees: tr.Fees, Slippage: tr.Slippage, ExitReason: tr.Reason, HoldBars: tr.HoldBars})
	}
	return rows
}

func btc5MValueAreaEvaluate(source BTC5MRollingValueAreaReversionSourceRow, rows []SummaryRow, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig) BTC5MRollingValueAreaReversionFalsification {
	byKey := map[string]SummaryRow{}
	for _, row := range rows {
		byKey[row.Split+"|"+row.Side] = row
	}
	full, oos, recent := byKey["full_2021_2026|all"], byKey["2023_2024_oos|all"], byKey["2025_2026_recent|all"]
	minPrimary := minInt(byKey["2021_2022_stress|all"].TotalTrades, minInt(oos.TotalTrades, recent.TotalTrades))
	dominantShare := 0.0
	if full.TotalTrades > 0 {
		dominantShare = float64(maxInt(byKey["2021_2022_stress|all"].TotalTrades, maxInt(oos.TotalTrades, recent.TotalTrades))) / float64(full.TotalTrades)
	}
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit && byKey["2021_2022_stress|all"].MaxDrawdown <= cfg.SplitMaxDrawdownLimit && oos.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && recent.MaxDrawdown <= cfg.SplitMaxDrawdownLimit
	tradeCountPass := full.TotalTrades >= cfg.MinFullTrades && minPrimary >= cfg.MinSplitTrades
	return BTC5MRollingValueAreaReversionFalsification{SourcePass: source.SourceFactsPass, LeakagePass: true, TradeCountPass: tradeCountPass, GrossEdgePass: full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0, NetEdgePass: full.NetPnL > 0, DrawdownPass: drawdownPass, RobustnessPass: tradeCountPass && dominantShare <= 0.90, SideReportingPass: true, FullExecutedTrades: full.TotalTrades, RequiredFullExecutedTrades: cfg.MinFullTrades, MinimumPrimarySplitExecutedTrades: minPrimary, RequiredPrimarySplitTrades: cfg.MinSplitTrades, FullGrossPnL: full.GrossPnL, FullNetPnL: full.NetPnL, FullProfitFactor: full.ProfitFactor, FullMaxDrawdown: full.MaxDrawdown, DominantPrimarySplitTradeShare: dominantShare}
}

func btc5MValueAreaFalsification(report BTC5MRollingValueAreaReversionFalsification, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig) BTC5MRollingValueAreaReversionFalsification {
	report.BacktestName = BacktestFirstBTC5MRollingValueAreaReversionName
	report.CandidateID = BTC5MRollingValueAreaReversionCandidateID
	if !report.SourcePass || !report.LeakagePass {
		report.StopState = BTC5MRollingValueAreaReversionStopStateFailedSourceOrLeakage
	} else if report.TradeCountPass && report.GrossEdgePass && report.NetEdgePass && report.DrawdownPass && report.RobustnessPass {
		report.StopState = BTC5MRollingValueAreaReversionStopStatePassedNeedsReview
	} else {
		report.StopState = BTC5MRollingValueAreaReversionStopStateFailedNoUsableStrategy
	}
	if !report.TradeCountPass {
		report.FailureReasons = append(report.FailureReasons, "insufficient_executed_trades")
	}
	if !report.GrossEdgePass {
		report.FailureReasons = append(report.FailureReasons, "gross_edge_gate_failed")
	}
	if !report.NetEdgePass {
		report.FailureReasons = append(report.FailureReasons, "net_edge_gate_failed")
	}
	if !report.DrawdownPass {
		report.FailureReasons = append(report.FailureReasons, "drawdown_gate_failed")
	}
	if !report.RobustnessPass {
		report.FailureReasons = append(report.FailureReasons, "robustness_gate_failed")
	}
	_ = cfg
	return report
}

func btc5MValueAreaSplit(t time.Time, splits []Split) string {
	for _, split := range splits {
		if split.Contains(t) {
			return split.Name
		}
	}
	return "unassigned"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
