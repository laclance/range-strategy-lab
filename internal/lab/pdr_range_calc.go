package lab

import "math"

func btc15MPrevDayRanges(candles []Candle) map[string]btc15MDayRange {
	ranges := map[string]btc15MDayRange{}
	for _, c := range candles {
		key := c.OpenTime.UTC().Format("2006-01-02")
		r := ranges[key]
		if r.count == 0 {
			r = btc15MDayRange{key: key, high: c.High, low: c.Low}
		}
		r.high = math.Max(r.high, c.High)
		r.low = math.Min(r.low, c.Low)
		r.count++
		ranges[key] = r
	}
	for key, r := range ranges {
		r.width = r.high - r.low
		r.mid = (r.high + r.low) / 2
		r.complete = r.count == 96 && r.width > 0
		ranges[key] = r
	}
	return ranges
}
