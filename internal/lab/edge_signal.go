package lab

import "strconv"

func (s *BTC15MRangeEdgeExhaustionFadeStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	d := ctx.Index
	candles := ctx.Candles
	split := btc15MEdgeSplit(candles[d].CloseTime, s.splits)
	if d < s.cfg.LookbackBars || d < 3 {
		s.skips[split+"|missing_warmup"]++
		return Signal{}, false
	}
	if d-1 >= len(s.atr) || !validNumber(s.atr[d-1]) || s.atr[d-1] <= 0 {
		s.skips[split+"|missing_prior_atr"]++
		return Signal{}, false
	}
	r, ok := btc15MEdgeRange(candles, d, s.cfg.LookbackBars)
	if !ok {
		s.skips[split+"|missing_prior_range"]++
		return Signal{}, false
	}
	startPos := btc15MEdgeStartPosition(candles[d-3].Close, r)
	if startPos < s.cfg.InteriorLowPct || startPos > s.cfg.InteriorHighPct {
		return Signal{}, false
	}
	priorATR := s.atr[d-1]
	side, progressATR, ok := btc15MEdgeSide(candles, d, r, priorATR, s.cfg)
	if !ok {
		return Signal{}, false
	}
	stop := r.low - s.cfg.StopATRMultiple*priorATR
	if side == Short {
		stop = r.high + s.cfg.StopATRMultiple*priorATR
	}
	signalID := BTC15MRangeEdgeExhaustionFadeCandidateID + "_" + strconv.Itoa(len(s.rows)+1)
	row := BTC15MRangeEdgeExhaustionFadeSignalRow{SignalID: signalID, CandidateID: BTC15MRangeEdgeExhaustionFadeCandidateID, Split: split, DecisionIndex: d, DecisionCloseTime: candles[d].CloseTime.UTC().Format(timeLayout), DecisionClose: candles[d].Close, Side: side, TimingLabel: "next_15m_open", LookbackBars: s.cfg.LookbackBars, RangeHigh: r.high, RangeLow: r.low, RangeMidpoint: r.mid, RangeWidth: r.width, StartClosePosition: startPos, FinalProgressATR: progressATR, PriorATR14: priorATR, Stop: stop, Target: r.mid, MaxHoldBars: s.cfg.MaxHoldBars}
	s.rows = append(s.rows, row)
	return Signal{Side: side, Stop: stop, Target: r.mid, MaxHoldBars: s.cfg.MaxHoldBars, Reason: signalID}, true
}
