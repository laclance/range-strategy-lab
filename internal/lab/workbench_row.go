package lab

import "strings"

func rangeWorkbenchBuildTrialRow(runID string, sourceRef string, spec RangeOptimizationWorkbenchTrialSpec, summaries []SummaryRow) (RangeOptimizationWorkbenchTrialRow, []RangeOptimizationWorkbenchSummaryRow, bool) {
	byKey := map[string]SummaryRow{}
	for _, row := range summaries {
		byKey[row.Split+"|"+row.Side] = row
	}
	full := byKey["full_2021_2026|all"]
	stress := byKey["2021_2022_stress|all"]
	oos := byKey["2023_2024_oos|all"]
	recent := byKey["2025_2026_recent|all"]
	longFull := byKey["full_2021_2026|long"]
	shortFull := byKey["full_2021_2026|short"]
	minPrimary := minInt(stress.TotalTrades, minInt(oos.TotalTrades, recent.TotalTrades))
	dominantShare := 0.0
	if full.TotalTrades > 0 {
		dominantShare = float64(maxInt(stress.TotalTrades, maxInt(oos.TotalTrades, recent.TotalTrades))) / float64(full.TotalTrades)
	}
	reasons := rangeWorkbenchFailureReasons(full, oos, recent, longFull, shortFull, minPrimary, dominantShare)
	robustness := rangeWorkbenchRobustnessScore(full, stress, oos, recent, longFull, shortFull, dominantShare)
	trialRow := RangeOptimizationWorkbenchTrialRow{RangeOptimizationWorkbenchTrialSpec: spec, RunID: runID, SourceRef: sourceRef, FullTrades: full.TotalTrades, StressTrades: stress.TotalTrades, OOSTrades: oos.TotalTrades, RecentTrades: recent.TotalTrades, FullGrossPnL: full.GrossPnL, FullNetPnL: full.NetPnL, FullProfitFactor: full.ProfitFactor, FullMaxDrawdown: full.MaxDrawdown, StressGrossPnL: stress.GrossPnL, StressNetPnL: stress.NetPnL, OOSGrossPnL: oos.GrossPnL, OOSNetPnL: oos.NetPnL, RecentGrossPnL: recent.GrossPnL, RecentNetPnL: recent.NetPnL, LongTrades: longFull.TotalTrades, ShortTrades: shortFull.TotalTrades, LongNetPnL: longFull.NetPnL, ShortNetPnL: shortFull.NetPnL, MinimumPrimarySplitTrades: minPrimary, DominantPrimarySplitTradeShare: dominantShare, RobustnessScore: robustness, FailureReasons: strings.Join(reasons, ";")}
	return trialRow, rangeWorkbenchSummaryRows(runID, spec.TrialID, summaries), len(reasons) == 0
}

func rangeWorkbenchSummaryRows(runID string, trialID string, summaries []SummaryRow) []RangeOptimizationWorkbenchSummaryRow {
	rows := []RangeOptimizationWorkbenchSummaryRow{}
	for _, row := range summaries {
		rows = append(rows, RangeOptimizationWorkbenchSummaryRow{RunID: runID, TrialID: trialID, Split: row.Split, Side: row.Side, TotalTrades: row.TotalTrades, Wins: row.Wins, Losses: row.Losses, WinRate: row.WinRate, GrossPnL: row.GrossPnL, NetPnL: row.NetPnL, TotalCosts: row.TotalCosts, ProfitFactor: row.ProfitFactor, GrossProfitFactor: row.GrossProfitFactor, MaxDrawdown: row.MaxDrawdown, Expectancy: row.Expectancy, AvgHoldBars: row.AvgHoldBars})
	}
	return rows
}
