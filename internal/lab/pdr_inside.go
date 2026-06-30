package lab

func btc15MPrevDayCurrentDayStillInside(candles []Candle, d int, prior btc15MDayRange) bool {
	day := candles[d].OpenTime.UTC().Format("2006-01-02")
	for i := d - 1; i >= 0 && candles[i].OpenTime.UTC().Format("2006-01-02") == day; i-- {
		if candles[i].Close < prior.low || candles[i].Close > prior.high {
			return false
		}
	}
	return true
}
