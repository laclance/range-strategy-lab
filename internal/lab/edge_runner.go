package lab

func RunBacktestFirstBTC15MRangeEdgeExhaustionFade(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC15MRangeEdgeExhaustionFadeResult, error) {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc15MEdgeSourceRow(manifest, cfg)
	result := BacktestFirstBTC15MRangeEdgeExhaustionFadeResult{SourceRows: []BTC15MRangeEdgeExhaustionFadeSourceRow{source}}
	resampled, coverage := btc15MEdgeResample(candles, cfg, source.SourceFactsPass)
	result.CoverageRows = []BTC15MRangeEdgeExhaustionFadeCoverageRow{coverage}
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = btc15MEdgeFalsification(BTC15MRangeEdgeExhaustionFadeFalsification{SourceResamplePass: false, LeakagePass: true, SideReportingPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}
	strategy := NewBTC15MRangeEdgeExhaustionFadeStrategy(resampled, cfg, splits)
	bt := RunBacktest(resampled, strategy, btCfg)
	strategy.markExecuted(bt.Trades)
	result.Trades = bt.Trades
	result.SignalRows = strategy.rows
	result.SkipRows = strategy.skipRows()
	result.TradeRows = btc15MEdgeTradeRows(bt.Trades, splits)
	result.SummaryRows = SummarizeSplits(bt.Trades, btCfg.StartBalance, splits)
	result.Falsification = btc15MEdgeFalsification(btc15MEdgeEvaluate(source, coverage, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}
