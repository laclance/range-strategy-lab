package lab

import "strconv"

func (s *BTC15MPreviousDayRangeReversionStrategy) btc15MPrevDaySignalRow(candles []Candle, d int, split string, priorKey string, prior btc15MDayRange, side Direction) BTC15MPreviousDayRangeReversionSignalRow {
	priorATR := s.atr[d-1]
	stop := prior.low - s.cfg.StopATRMultiple*priorATR
	if side == Short {
		stop = prior.high + s.cfg.StopATRMultiple*priorATR
	}
	signalID := BTC15MPreviousDayRangeReversionCandidateID + "_" + strconv.Itoa(len(s.rows)+1)
	return BTC15MPreviousDayRangeReversionSignalRow{SignalID: signalID, CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Split: split, DecisionIndex: d, DecisionCloseTime: candles[d].CloseTime.UTC().Format(timeLayout), DecisionClose: candles[d].Close, Side: side, TimingLabel: "next_15m_open", PriorDay: priorKey, PriorDayHigh: prior.high, PriorDayLow: prior.low, PriorDayMidpoint: prior.mid, PriorDayRangeWidth: prior.width, OuterZonePct: s.cfg.OuterZonePct, PriorATR14: priorATR, Stop: stop, Target: prior.mid, MaxHoldBars: s.cfg.MaxHoldBars}
}
