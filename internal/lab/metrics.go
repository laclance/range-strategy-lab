package lab

type SummaryRow struct {
	Split             string  `json:"split"`
	Side              string  `json:"side"`
	TotalTrades       int     `json:"total_trades"`
	Wins              int     `json:"wins"`
	Losses            int     `json:"losses"`
	WinRate           float64 `json:"win_rate"`
	GrossPnL          float64 `json:"gross_pnl"`
	NetPnL            float64 `json:"net_pnl"`
	TotalCosts        float64 `json:"total_costs"`
	ProfitFactor      float64 `json:"profit_factor"`
	GrossProfitFactor float64 `json:"gross_profit_factor"`
	MaxDrawdown       float64 `json:"max_drawdown"`
	Expectancy        float64 `json:"expectancy"`
	AvgHoldBars       float64 `json:"avg_hold_bars"`
}

func SummarizeSplits(trades []Trade, startBalance float64, splits []Split) []SummaryRow {
	rows := make([]SummaryRow, 0, len(splits)*3)
	for _, split := range splits {
		for _, side := range []string{"all", "long", "short"} {
			var filtered []Trade
			for _, tr := range trades {
				closeTime, err := parseTime(tr.ExitTime)
				if err != nil || !split.Contains(closeTime) {
					continue
				}
				if side != "all" && string(tr.Side) != side {
					continue
				}
				filtered = append(filtered, tr)
			}
			row := summarizeTrades(filtered, startBalance)
			row.Split = split.Name
			row.Side = side
			rows = append(rows, row)
		}
	}
	return rows
}

func summarizeTrades(trades []Trade, startBalance float64) SummaryRow {
	row := SummaryRow{TotalTrades: len(trades)}
	balance := startBalance
	equity := []float64{startBalance}
	netProfit, netLoss := 0.0, 0.0
	grossProfit, grossLoss := 0.0, 0.0
	holdBars := 0
	for _, tr := range trades {
		row.GrossPnL += tr.GrossPnL
		row.NetPnL += tr.NetPnL
		row.TotalCosts += tr.Fees + tr.Slippage
		holdBars += tr.HoldBars
		if tr.NetPnL > 0 {
			row.Wins++
			netProfit += tr.NetPnL
		} else if tr.NetPnL < 0 {
			row.Losses++
			netLoss += -tr.NetPnL
		}
		if tr.GrossPnL > 0 {
			grossProfit += tr.GrossPnL
		} else if tr.GrossPnL < 0 {
			grossLoss += -tr.GrossPnL
		}
		balance += tr.NetPnL
		equity = append(equity, balance)
	}
	if row.TotalTrades > 0 {
		row.WinRate = float64(row.Wins) / float64(row.TotalTrades)
		row.Expectancy = row.NetPnL / float64(row.TotalTrades)
		row.AvgHoldBars = float64(holdBars) / float64(row.TotalTrades)
	}
	if netLoss > 0 {
		row.ProfitFactor = netProfit / netLoss
	} else if netProfit > 0 {
		row.ProfitFactor = 999.99
	}
	if grossLoss > 0 {
		row.GrossProfitFactor = grossProfit / grossLoss
	} else if grossProfit > 0 {
		row.GrossProfitFactor = 999.99
	}
	row.MaxDrawdown = MaxDrawdown(equity)
	return row
}

func MaxDrawdown(equity []float64) float64 {
	if len(equity) == 0 {
		return 0
	}
	peak := equity[0]
	maxDD := 0.0
	for _, value := range equity {
		if value > peak {
			peak = value
		}
		if peak > 0 {
			dd := (peak - value) / peak
			if dd > maxDD {
				maxDD = dd
			}
		}
	}
	return maxDD
}
