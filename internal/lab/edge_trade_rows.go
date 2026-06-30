package lab

import "time"

func btc15MEdgeTradeRows(trades []Trade, splits []Split) []BTC15MRangeEdgeExhaustionFadeTradeRow {
	rows := make([]BTC15MRangeEdgeExhaustionFadeTradeRow, 0, len(trades))
	for _, tr := range trades {
		exitTime, _ := time.Parse(timeLayout, tr.ExitTime)
		rows = append(rows, BTC15MRangeEdgeExhaustionFadeTradeRow{SignalID: tr.Signal, CloseSplit: btc15MEdgeSplit(exitTime, splits), Side: tr.Side, EntryTime: tr.EntryTime, ExitTime: tr.ExitTime, GrossPnL: tr.GrossPnL, NetPnL: tr.NetPnL, Fees: tr.Fees, Slippage: tr.Slippage, ExitReason: tr.Reason, HoldBars: tr.HoldBars})
	}
	return rows
}
