package lab

func btc15MPrevDayEvaluate(source BTC15MPreviousDayRangeReversionSourceRow, coverage BTC15MPreviousDayRangeReversionCoverageRow, rows []SummaryRow, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) BTC15MPreviousDayRangeReversionFalsification {
	byKey := map[string]SummaryRow{}
	for _, row := range rows {
		byKey[row.Split+"|"+row.Side] = row
	}
	full := byKey["full_2021_2026|all"]
	stress := byKey["2021_2022_stress|all"]
	oos := byKey["2023_2024_oos|all"]
	recent := byKey["2025_2026_recent|all"]
	minPrimary := minInt(stress.TotalTrades, minInt(oos.TotalTrades, recent.TotalTrades))
	dominantShare := 0.0
	if full.TotalTrades > 0 {
		dominantShare = float64(maxInt(stress.TotalTrades, maxInt(oos.TotalTrades, recent.TotalTrades))) / float64(full.TotalTrades)
	}
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit && stress.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && oos.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && recent.MaxDrawdown <= cfg.SplitMaxDrawdownLimit
	tradeCountPass := full.TotalTrades >= cfg.MinFullTrades && minPrimary >= cfg.MinSplitTrades
	return BTC15MPreviousDayRangeReversionFalsification{SourceResamplePass: source.SourceFactsPass && coverage.SourceResamplePass, LeakagePass: true, TradeCountPass: tradeCountPass, GrossEdgePass: full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0, NetEdgePass: full.NetPnL > 0, DrawdownPass: drawdownPass, RobustnessPass: tradeCountPass && dominantShare <= 0.90, SideReportingPass: true, FullExecutedTrades: full.TotalTrades, RequiredFullExecutedTrades: cfg.MinFullTrades, MinimumPrimarySplitExecutedTrades: minPrimary, RequiredPrimarySplitTrades: cfg.MinSplitTrades, FullGrossPnL: full.GrossPnL, FullNetPnL: full.NetPnL, FullProfitFactor: full.ProfitFactor, FullMaxDrawdown: full.MaxDrawdown, DominantPrimarySplitTradeShare: dominantShare}
}
