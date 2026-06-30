package lab

func btc15MEdgeSide(candles []Candle, d int, r btc15MEdgeRange, atr float64, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig) (Direction, float64, bool) {
	c2 := candles[d-2].Close
	c1 := candles[d-1].Close
	c0 := candles[d].Close
	longProgress := c1 - c0
	if c2 > c1 && c1 > c0 && c0 > r.low && c0 <= r.low+cfg.EdgeZonePct*r.width && longProgress > 0 && longProgress < cfg.ProgressATRMultiple*atr {
		return Long, longProgress / atr, true
	}
	shortProgress := c0 - c1
	if c2 < c1 && c1 < c0 && c0 < r.high && c0 >= r.high-cfg.EdgeZonePct*r.width && shortProgress > 0 && shortProgress < cfg.ProgressATRMultiple*atr {
		return Short, shortProgress / atr, true
	}
	return "", 0, false
}
