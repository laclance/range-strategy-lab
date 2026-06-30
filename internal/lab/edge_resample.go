package lab

import (
	"fmt"
	"math"
	"time"
)

func btc15MEdgeResample(candles []Candle, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig, sourcePass bool) ([]Candle, BTC15MRangeEdgeExhaustionFadeCoverageRow) {
	out := []Candle{}
	missingBuckets := 0
	for i := 0; i+2 < len(candles); i += 3 {
		first, second, third := candles[i], candles[i+1], candles[i+2]
		start := first.OpenTime.UTC()
		if start.Minute()%15 != 0 || !second.OpenTime.UTC().Equal(start.Add(5*time.Minute)) || !third.OpenTime.UTC().Equal(start.Add(10*time.Minute)) {
			missingBuckets++
			continue
		}
		high := math.Max(first.High, math.Max(second.High, third.High))
		low := math.Min(first.Low, math.Min(second.Low, third.Low))
		out = append(out, Candle{OpenTime: start, CloseTime: third.CloseTime.UTC(), Open: first.Open, High: high, Low: low, Close: third.Close, Volume: first.Volume + second.Volume + third.Volume})
	}
	row := BTC15MRangeEdgeExhaustionFadeCoverageRow{BacktestName: BacktestFirstBTC15MRangeEdgeExhaustionFadeName, CandidateID: BTC15MRangeEdgeExhaustionFadeCandidateID, Timeframe: "15m", RowCount: len(out), ExpectedRowCount: cfg.Expected15MRows, ExpectedLastOpenTime: cfg.Expected15MLastOpenTime, ExpectedChildBars: 3, MissingChildBuckets: missingBuckets, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	if len(out) > 0 {
		row.FirstOpenTime = out[0].OpenTime.UTC().Format(timeLayout)
		row.LastOpenTime = out[len(out)-1].OpenTime.UTC().Format(timeLayout)
	}
	row.SourceResamplePass = sourcePass && missingBuckets == 0 && row.RowCount == cfg.Expected15MRows && row.LastOpenTime == cfg.Expected15MLastOpenTime
	if !row.SourceResamplePass {
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf("15m resample mismatch rows=%d expected=%d missing_child_buckets=%d last_open=%s expected=%s", row.RowCount, cfg.Expected15MRows, missingBuckets, row.LastOpenTime, cfg.Expected15MLastOpenTime)
	}
	return out, row
}
