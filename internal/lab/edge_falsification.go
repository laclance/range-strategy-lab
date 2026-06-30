package lab

func btc15MEdgeFalsification(report BTC15MRangeEdgeExhaustionFadeFalsification, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig) BTC15MRangeEdgeExhaustionFadeFalsification {
	report.BacktestName = BacktestFirstBTC15MRangeEdgeExhaustionFadeName
	report.CandidateID = BTC15MRangeEdgeExhaustionFadeCandidateID
	if !report.SourceResamplePass || !report.LeakagePass {
		report.StopState = BTC15MRangeEdgeExhaustionFadeStopStateFailedSourceOrResample
	} else if report.TradeCountPass && report.GrossEdgePass && report.NetEdgePass && report.DrawdownPass && report.RobustnessPass {
		report.StopState = BTC15MRangeEdgeExhaustionFadeStopStatePassedNeedsReview
	} else {
		report.StopState = BTC15MRangeEdgeExhaustionFadeStopStateFailedNoUsableStrategy
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
