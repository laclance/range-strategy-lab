package lab

import "time"

func btc15MPrevDayTradeRows(trades []Trade, splits []Split) []BTC15MPreviousDayRangeReversionTradeRow {
	rows := make([]BTC15MPreviousDayRangeReversionTradeRow, 0, len(trades))
	for _, tr := range trades {
		exitTime, _ := time.Parse(timeLayout, tr.ExitTime)
		row := BTC15MPreviousDayRangeReversionTradeRow{SignalID: tr.Signal, CloseSplit: btc15MPrevDaySplit(exitTime, splits), Side: tr.Side, EntryTime: tr.EntryTime, ExitTime: tr.ExitTime, GrossPnL: tr.GrossPnL, NetPnL: tr.NetPnL, Fees: tr.Fees, Slippage: tr.Slippage, ExitReason: tr.Reason, HoldBars: tr.HoldBars}
		rows = append(rows, row)
	}
	return rows
}
