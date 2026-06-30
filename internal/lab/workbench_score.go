package lab

import "math"

func rangeWorkbenchFailureReasons(full SummaryRow, oos SummaryRow, recent SummaryRow, longFull SummaryRow, shortFull SummaryRow, minPrimary int, dominantShare float64) []string {
	reasons := []string{}
	if full.TotalTrades < 200 {
		reasons = append(reasons, "full_trades_lt_200")
	}
	if minPrimary < 40 {
		reasons = append(reasons, "primary_split_trades_lt_40")
	}
	if full.GrossPnL <= 0 {
		reasons = append(reasons, "full_gross_not_positive")
	}
	if full.NetPnL <= 0 {
		reasons = append(reasons, "full_net_not_positive")
	}
	if oos.GrossPnL < 0 {
		reasons = append(reasons, "oos_gross_negative")
	}
	if recent.GrossPnL < 0 {
		reasons = append(reasons, "recent_gross_negative")
	}
	if full.MaxDrawdown > 0.30 {
		reasons = append(reasons, "full_drawdown_gt_0_30")
	}
	if dominantShare > 0.75 {
		reasons = append(reasons, "dominant_split_trade_share_gt_0_75")
	}
	if longFull.TotalTrades == 0 || shortFull.TotalTrades == 0 {
		reasons = append(reasons, "missing_long_or_short_side")
	}
	return reasons
}

func rangeWorkbenchRobustnessScore(full SummaryRow, stress SummaryRow, oos SummaryRow, recent SummaryRow, longFull SummaryRow, shortFull SummaryRow, dominantShare float64) float64 {
	netConsistency := 0.0
	grossConsistency := 0.0
	for _, row := range []SummaryRow{stress, oos, recent} {
		if row.NetPnL > 0 {
			netConsistency++
		}
		if row.GrossPnL >= 0 {
			grossConsistency++
		}
	}
	pfStability := math.Min(full.ProfitFactor, math.Min(stress.ProfitFactor, math.Min(oos.ProfitFactor, recent.ProfitFactor)))
	drawdownPenalty := math.Max(0, full.MaxDrawdown-0.20) * 5
	sidePenalty := 0.0
	if full.TotalTrades > 0 {
		sideShare := math.Max(float64(longFull.TotalTrades), float64(shortFull.TotalTrades)) / float64(full.TotalTrades)
		sidePenalty = sideShare
	}
	return netConsistency + grossConsistency + pfStability - drawdownPenalty - dominantShare - sidePenalty
}
