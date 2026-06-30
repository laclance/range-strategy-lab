package lab

func btc15MPrevDayFalsification(report BTC15MPreviousDayRangeReversionFalsification, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) BTC15MPreviousDayRangeReversionFalsification {
	report.BacktestName = BacktestFirstBTC15MPreviousDayRangeReversionName
	report.CandidateID = BTC15MPreviousDayRangeReversionCandidateID
	if !report.SourceResamplePass || !report.LeakagePass {
		report.StopState = BTC15MPreviousDayRangeReversionStopStateFailedSourceOrResample
	} else if report.TradeCountPass && report.GrossEdgePass && report.NetEdgePass && report.DrawdownPass && report.RobustnessPass {
		report.StopState = BTC15MPreviousDayRangeReversionStopStatePassedNeedsReview
	} else {
		report.StopState = BTC15MPreviousDayRangeReversionStopStateFailedNoUsableStrategy
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
