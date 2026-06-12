package lab

import "math"

func ATR14(candles []Candle) []float64 {
	return NormalizedATR(candles, 14)
}

func NormalizedATR(candles []Candle, period int) []float64 {
	atr := ATR(candles, period)
	out := nanSlice(len(candles))
	for i, value := range atr {
		if !validNumber(value) || candles[i].Close <= 0 {
			continue
		}
		out[i] = value / candles[i].Close
	}
	return out
}

func ATR(candles []Candle, period int) []float64 {
	out := nanSlice(len(candles))
	if period <= 0 || len(candles) < period {
		return out
	}
	tr := make([]float64, len(candles))
	for i, c := range candles {
		tr[i] = c.High - c.Low
		if i == 0 {
			continue
		}
		highGap := math.Abs(c.High - candles[i-1].Close)
		lowGap := math.Abs(c.Low - candles[i-1].Close)
		tr[i] = math.Max(tr[i], math.Max(highGap, lowGap))
	}
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += tr[i]
	}
	out[period-1] = sum / float64(period)
	for i := period; i < len(candles); i++ {
		out[i] = (out[i-1]*float64(period-1) + tr[i]) / float64(period)
	}
	return out
}

func Donchian20Width(candles []Candle) []float64 {
	return DonchianWidth(candles, 20)
}

func DonchianWidth(candles []Candle, period int) []float64 {
	out := nanSlice(len(candles))
	if period <= 0 || len(candles) < period {
		return out
	}
	for i := period - 1; i < len(candles); i++ {
		highest := candles[i-period+1].High
		lowest := candles[i-period+1].Low
		for j := i - period + 2; j <= i; j++ {
			if candles[j].High > highest {
				highest = candles[j].High
			}
			if candles[j].Low < lowest {
				lowest = candles[j].Low
			}
		}
		if candles[i].Close > 0 {
			out[i] = (highest - lowest) / candles[i].Close
		}
	}
	return out
}

func Bollinger20Width(candles []Candle) []float64 {
	return BollingerWidth(candles, 20)
}

func BollingerWidth(candles []Candle, period int) []float64 {
	out := nanSlice(len(candles))
	if period <= 0 || len(candles) < period {
		return out
	}
	for i := period - 1; i < len(candles); i++ {
		sum := 0.0
		for j := i - period + 1; j <= i; j++ {
			sum += candles[j].Close
		}
		mean := sum / float64(period)
		if mean <= 0 {
			continue
		}
		variance := 0.0
		for j := i - period + 1; j <= i; j++ {
			diff := candles[j].Close - mean
			variance += diff * diff
		}
		stddev := math.Sqrt(variance / float64(period))
		out[i] = (4 * stddev) / mean
	}
	return out
}

func ADX14(candles []Candle) []float64 {
	return ADX(candles, 14)
}

func ADX(candles []Candle, period int) []float64 {
	out := nanSlice(len(candles))
	if period <= 0 || len(candles) <= period*2-1 {
		return out
	}

	tr := make([]float64, len(candles))
	plusDM := make([]float64, len(candles))
	minusDM := make([]float64, len(candles))
	for i := 1; i < len(candles); i++ {
		c := candles[i]
		prev := candles[i-1]
		tr[i] = math.Max(c.High-c.Low, math.Max(math.Abs(c.High-prev.Close), math.Abs(c.Low-prev.Close)))

		upMove := c.High - prev.High
		downMove := prev.Low - c.Low
		if upMove > downMove && upMove > 0 {
			plusDM[i] = upMove
		}
		if downMove > upMove && downMove > 0 {
			minusDM[i] = downMove
		}
	}

	smoothedTR := 0.0
	smoothedPlusDM := 0.0
	smoothedMinusDM := 0.0
	for i := 1; i <= period; i++ {
		smoothedTR += tr[i]
		smoothedPlusDM += plusDM[i]
		smoothedMinusDM += minusDM[i]
	}

	dx := nanSlice(len(candles))
	for i := period; i < len(candles); i++ {
		if i > period {
			smoothedTR = smoothedTR - smoothedTR/float64(period) + tr[i]
			smoothedPlusDM = smoothedPlusDM - smoothedPlusDM/float64(period) + plusDM[i]
			smoothedMinusDM = smoothedMinusDM - smoothedMinusDM/float64(period) + minusDM[i]
		}
		if smoothedTR <= 0 {
			dx[i] = 0
			continue
		}
		plusDI := 100 * smoothedPlusDM / smoothedTR
		minusDI := 100 * smoothedMinusDM / smoothedTR
		denom := plusDI + minusDI
		if denom <= 0 {
			dx[i] = 0
			continue
		}
		dx[i] = 100 * math.Abs(plusDI-minusDI) / denom
	}

	firstADX := period*2 - 1
	sumDX := 0.0
	for i := period; i <= firstADX; i++ {
		sumDX += dx[i]
	}
	out[firstADX] = sumDX / float64(period)
	for i := firstADX + 1; i < len(candles); i++ {
		out[i] = (out[i-1]*float64(period-1) + dx[i]) / float64(period)
	}
	return out
}

func nanSlice(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = math.NaN()
	}
	return out
}

func validNumber(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
