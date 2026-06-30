package lab

func (s *BTC15MPreviousDayRangeReversionStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	d := ctx.Index
	candles := ctx.Candles
	decision := candles[d]
	split := btc15MPrevDaySplit(decision.CloseTime, s.splits)
	if d == 0 {
		s.skips[split+"|missing_prior_day"]++
		return Signal{}, false
	}
	if d-1 >= len(s.atr) || !validNumber(s.atr[d-1]) || s.atr[d-1] <= 0 {
		s.skips[split+"|missing_prior_atr"]++
		return Signal{}, false
	}
	priorKey := decision.OpenTime.UTC().AddDate(0, 0, -1).Format("2006-01-02")
	prior, ok := s.dayRanges[priorKey]
	if !ok || !prior.complete || prior.width <= 0 {
		s.skips[split+"|missing_complete_prior_day"]++
		return Signal{}, false
	}
	if !btc15MPrevDayCurrentDayStillInside(candles, d, prior) {
		return Signal{}, false
	}
	side, ok := btc15MPrevDaySide(decision.Close, prior, s.cfg.OuterZonePct)
	if !ok {
		return Signal{}, false
	}
	row := s.btc15MPrevDaySignalRow(candles, d, split, priorKey, prior, side)
	s.rows = append(s.rows, row)
	return Signal{Side: side, Stop: row.Stop, Target: row.Target, MaxHoldBars: s.cfg.MaxHoldBars, Reason: row.SignalID}, true
}
