package lab

import "time"

type RangeOptimizationWorkbenchTrialStrategy struct {
	spec   RangeOptimizationWorkbenchTrialSpec
	atr    []float64
	splits []Split
	skips  map[string]int
}

func NewRangeOptimizationWorkbenchTrialStrategy(candles []Candle, spec RangeOptimizationWorkbenchTrialSpec, splits []Split) *RangeOptimizationWorkbenchTrialStrategy {
	return &RangeOptimizationWorkbenchTrialStrategy{spec: spec, atr: ATR(candles, 14), splits: splits, skips: map[string]int{}}
}

func (s *RangeOptimizationWorkbenchTrialStrategy) Name() string {
	return RangeOptimizationWorkbenchName
}

func (s *RangeOptimizationWorkbenchTrialStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	d := ctx.Index
	candles := ctx.Candles
	split := rangeWorkbenchSplit(candles[d].CloseTime, s.splits)
	if d < s.spec.RangeLookback || d < 3 {
		s.skips[split+"|missing_warmup"]++
		return Signal{}, false
	}
	if d-1 >= len(s.atr) || !validNumber(s.atr[d-1]) || s.atr[d-1] <= 0 {
		s.skips[split+"|missing_prior_atr"]++
		return Signal{}, false
	}
	facts, ok := s.facts(candles, d)
	if !ok {
		s.skips[split+"|missing_range_facts"]++
		return Signal{}, false
	}
	priorATR := s.atr[d-1]
	if s.spec.MinRangeATR > 0 && facts.width < s.spec.MinRangeATR*priorATR {
		return Signal{}, false
	}
	if s.spec.MaxRangeATR > 0 && facts.width > s.spec.MaxRangeATR*priorATR {
		return Signal{}, false
	}
	side, ok := s.side(candles, d, facts, priorATR)
	if !ok {
		return Signal{}, false
	}
	stop := s.stop(side, candles[d].Close, facts, priorATR)
	target := s.target(side, facts)
	return Signal{Side: side, Stop: stop, Target: target, MaxHoldBars: s.spec.TimeStopBars, Reason: s.spec.TrialID}, true
}

func (s *RangeOptimizationWorkbenchTrialStrategy) facts(candles []Candle, d int) (rangeWorkbenchFacts, bool) {
	if s.spec.EntryArchetype == "previous_day_edge_reversion" {
		return rangeWorkbenchPreviousDayFacts(candles, d)
	}
	return rangeWorkbenchRollingFacts(candles, d, s.spec.RangeLookback)
}

func (s *RangeOptimizationWorkbenchTrialStrategy) side(candles []Candle, d int, facts rangeWorkbenchFacts, atr float64) (Direction, bool) {
	close := candles[d].Close
	lowerEdge := facts.low + s.spec.EdgeZonePct*facts.width
	upperEdge := facts.high - s.spec.EdgeZonePct*facts.width
	switch s.spec.EntryArchetype {
	case "midpoint_reversion", "previous_day_edge_reversion":
		if close >= facts.low && close <= lowerEdge {
			return Long, true
		}
		if close <= facts.high && close >= upperEdge {
			return Short, true
		}
	case "vwap_reversion":
		if close >= facts.low && close <= lowerEdge && facts.vwap-close >= s.spec.VWAPDistanceRangePct*facts.width {
			return Long, true
		}
		if close <= facts.high && close >= upperEdge && close-facts.vwap >= s.spec.VWAPDistanceRangePct*facts.width {
			return Short, true
		}
	case "interior_to_edge_exhaustion_fade":
		startPos := rangeWorkbenchPosition(candles[d-3].Close, facts)
		if startPos < s.spec.InteriorLowPct || startPos > s.spec.InteriorHighPct {
			return "", false
		}
		c2 := candles[d-2].Close
		c1 := candles[d-1].Close
		c0 := candles[d].Close
		longProgress := c1 - c0
		if c2 > c1 && c1 > c0 && c0 > facts.low && c0 <= lowerEdge && longProgress > 0 && longProgress < s.spec.ProgressATRMultiple*atr {
			return Long, true
		}
		shortProgress := c0 - c1
		if c2 < c1 && c1 < c0 && c0 < facts.high && c0 >= upperEdge && shortProgress > 0 && shortProgress < s.spec.ProgressATRMultiple*atr {
			return Short, true
		}
	}
	return "", false
}

func (s *RangeOptimizationWorkbenchTrialStrategy) stop(side Direction, entry float64, facts rangeWorkbenchFacts, atr float64) float64 {
	switch s.spec.StopMode {
	case "range_edge_plus_0.50_atr":
		if side == Short {
			return facts.high + 0.50*atr
		}
		return facts.low - 0.50*atr
	case "entry_minus_1.00_atr":
		if side == Short {
			return entry + atr
		}
		return entry - atr
	default:
		if side == Short {
			return facts.high + 0.25*atr
		}
		return facts.low - 0.25*atr
	}
}

func (s *RangeOptimizationWorkbenchTrialStrategy) target(side Direction, facts rangeWorkbenchFacts) float64 {
	switch s.spec.TargetMode {
	case "vwap":
		return facts.vwap
	case "opposite_inner_zone":
		if side == Short {
			return facts.low + s.spec.EdgeZonePct*facts.width
		}
		return facts.high - s.spec.EdgeZonePct*facts.width
	default:
		return facts.mid
	}
}

func rangeWorkbenchPosition(close float64, facts rangeWorkbenchFacts) float64 {
	if facts.width <= 0 {
		return 0
	}
	return (close - facts.low) / facts.width
}

func rangeWorkbenchSplit(t time.Time, splits []Split) string {
	for _, split := range splits {
		if split.Contains(t) {
			return split.Name
		}
	}
	return "unassigned"
}
