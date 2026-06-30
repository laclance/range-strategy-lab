package lab

import "math"

func btc15MEdgeRange(candles []Candle, d int, lookback int) (btc15MEdgeRange, bool) {
	if lookback <= 0 || d < lookback {
		return btc15MEdgeRange{}, false
	}
	high := candles[d-lookback].High
	low := candles[d-lookback].Low
	for i := d - lookback; i <= d-1; i++ {
		high = math.Max(high, candles[i].High)
		low = math.Min(low, candles[i].Low)
	}
	width := high - low
	if width <= 0 {
		return btc15MEdgeRange{}, false
	}
	return btc15MEdgeRange{high: high, low: low, mid: (high + low) / 2, width: width}, true
}

func btc15MEdgeStartPosition(close float64, r btc15MEdgeRange) float64 {
	if r.width <= 0 {
		return 0
	}
	return (close - r.low) / r.width
}
