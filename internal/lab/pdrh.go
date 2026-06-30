package lab

func btc15MPrevDaySide(close float64, prior btc15MDayRange, outerPct float64) (Direction, bool) {
	if close < prior.low || close > prior.high {
		return "", false
	}
	if close <= prior.low+outerPct*prior.width {
		return Long, true
	}
	if close >= prior.high-outerPct*prior.width {
		return Short, true
	}
	return "", false
}
