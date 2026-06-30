package lab

func RunBacktestFirstBTC15MPreviousDayRangeReversion(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC15MPreviousDayRangeReversionResult, error) {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc15MPrevDaySourceRow(manifest, cfg)
	result := BacktestFirstBTC15MPreviousDayRangeReversionResult{SourceRows: []BTC15MPreviousDayRangeReversionSourceRow{source}}
	resampled, coverage := btc15MPrevDayResample(candles, cfg, source.SourceFactsPass)
	result.CoverageRows = []BTC15MPreviousDayRangeReversionCoverageRow{coverage}
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = btc15MPrevDayFalsification(BTC15MPreviousDayRangeReversionFalsification{SourceResamplePass: false, LeakagePass: true, SideReportingPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}
	strategy := NewBTC15MPreviousDayRangeReversionStrategy(resampled, cfg, splits)
	bt := RunBacktest(resampled, strategy, btCfg)
	strategy.markExecuted(bt.Trades)
	result.Trades = bt.Trades
	result.SignalRows = strategy.rows
	result.SkipRows = strategy.skipRows()
	result.TradeRows = btc15MPrevDayTradeRows(bt.Trades, splits)
	result.SummaryRows = SummarizeSplits(bt.Trades, btCfg.StartBalance, splits)
	result.Falsification = btc15MPrevDayFalsification(btc15MPrevDayEvaluate(source, coverage, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}
