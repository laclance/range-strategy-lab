package lab

import (
	"fmt"
	"math"
	"strings"
)

type BTC5MRollingValueAreaReversionStrategy struct {
	cfg    BacktestFirstBTC5MRollingValueAreaReversionConfig
	atr    []float64
	splits []Split
	rows   []BTC5MRollingValueAreaReversionSignalRow
	skips  map[string]int
}

func RunBacktestFirstBTC5MRollingValueAreaReversion(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC5MRollingValueAreaReversionResult, error) {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc5MValueAreaSourceRow(manifest, cfg)
	result := BacktestFirstBTC5MRollingValueAreaReversionResult{SourceRows: []BTC5MRollingValueAreaReversionSourceRow{source}}
	if !source.SourceFactsPass {
		result.Falsification = btc5MValueAreaFalsification(BTC5MRollingValueAreaReversionFalsification{SourcePass: false, LeakagePass: true, SideReportingPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}
	strategy := NewBTC5MRollingValueAreaReversionStrategy(candles, cfg, splits)
	bt := RunBacktest(candles, strategy, btCfg)
	strategy.markExecuted(bt.Trades)
	result.Trades = bt.Trades
	result.SignalRows = strategy.rows
	result.SkipRows = strategy.skipRows()
	result.TradeRows = btc5MValueAreaTradeRows(bt.Trades, splits)
	result.SummaryRows = SummarizeSplits(bt.Trades, btCfg.StartBalance, splits)
	result.Falsification = btc5MValueAreaFalsification(btc5MValueAreaEvaluate(source, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}

func NewBTC5MRollingValueAreaReversionStrategy(candles []Candle, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig, splits []Split) *BTC5MRollingValueAreaReversionStrategy {
	cfg = cfg.withDefaults()
	return &BTC5MRollingValueAreaReversionStrategy{cfg: cfg, atr: ATR(candles, cfg.ATRPeriod), splits: splits, skips: map[string]int{}}
}

func (s *BTC5MRollingValueAreaReversionStrategy) Name() string {
	return BacktestFirstBTC5MRollingValueAreaReversionName
}

func (s *BTC5MRollingValueAreaReversionStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	d := ctx.Index
	candles := ctx.Candles
	split := btc5MValueAreaSplit(candles[d].CloseTime, s.splits)
	if d < s.cfg.LookbackBars {
		s.skips[split+"|missing_warmup"]++
		return Signal{}, false
	}
	if d <= 0 || d-1 >= len(s.atr) || !validNumber(s.atr[d-1]) || s.atr[d-1] <= 0 {
		s.skips[split+"|missing_prior_atr"]++
		return Signal{}, false
	}
	facts, ok := btc5MValueAreaFacts(candles, d, s.cfg.LookbackBars)
	if !ok {
		s.skips[split+"|missing_prior_value_window"]++
		return Signal{}, false
	}
	priorATR := s.atr[d-1]
	if facts.width < s.cfg.MinRangeATRs*priorATR {
		return Signal{}, false
	}
	side, ok := btc5MValueAreaSide(candles[d].Close, facts, s.cfg)
	if !ok {
		return Signal{}, false
	}
	stop := math.Min(candles[d].Low, facts.low) - s.cfg.StopATRMultiple*priorATR
	if side == Short {
		stop = math.Max(candles[d].High, facts.high) + s.cfg.StopATRMultiple*priorATR
	}
	row := BTC5MRollingValueAreaReversionSignalRow{
		SignalID:                  fmt.Sprintf("%s_%06d", BTC5MRollingValueAreaReversionCandidateID, len(s.rows)+1),
		CandidateID:               BTC5MRollingValueAreaReversionCandidateID,
		Split:                     split,
		DecisionIndex:             d,
		DecisionCloseTime:         candles[d].CloseTime.UTC().Format(timeLayout),
		DecisionClose:             candles[d].Close,
		Side:                      side,
		TimingLabel:               "next_5m_open",
		LookbackBars:              s.cfg.LookbackBars,
		RangeHigh:                 facts.high,
		RangeLow:                  facts.low,
		RangeWidth:                facts.width,
		RangeWidthATRs:            facts.width / priorATR,
		RollingVWAP:               facts.vwap,
		PriorATR14:                priorATR,
		Stop:                      stop,
		Target:                    facts.vwap,
		MaxHoldBars:               s.cfg.MaxHoldBars,
		ForwardLabelsAsInput:      false,
		UsesFutureRowsForFeatures: false,
		DerivativesVetoUsed:       false,
		OptimizerSelectionUsed:    false,
	}
	s.rows = append(s.rows, row)
	return Signal{Side: side, Stop: stop, Target: facts.vwap, MaxHoldBars: s.cfg.MaxHoldBars, Reason: row.SignalID}, true
}

type btc5MValueAreaWindow struct{ high, low, width, vwap float64 }

func btc5MValueAreaFacts(candles []Candle, d, lookback int) (btc5MValueAreaWindow, bool) {
	high, low := candles[d-lookback].High, candles[d-lookback].Low
	weighted, volume := 0.0, 0.0
	for i := d - lookback; i <= d-1; i++ {
		c := candles[i]
		high = math.Max(high, c.High)
		low = math.Min(low, c.Low)
		weighted += ((c.High + c.Low + c.Close) / 3) * c.Volume
		volume += c.Volume
	}
	if high <= low || volume <= 0 {
		return btc5MValueAreaWindow{}, false
	}
	return btc5MValueAreaWindow{high: high, low: low, width: high - low, vwap: weighted / volume}, true
}

func btc5MValueAreaSide(close float64, facts btc5MValueAreaWindow, cfg BacktestFirstBTC5MRollingValueAreaReversionConfig) (Direction, bool) {
	if close <= facts.low+cfg.OuterZonePct*facts.width && close <= facts.vwap-cfg.VWAPDistanceRangePct*facts.width {
		return Long, true
	}
	if close >= facts.high-cfg.OuterZonePct*facts.width && close >= facts.vwap+cfg.VWAPDistanceRangePct*facts.width {
		return Short, true
	}
	return "", false
}

func (s *BTC5MRollingValueAreaReversionStrategy) markExecuted(trades []Trade) {
	executed := map[string]bool{}
	for _, tr := range trades {
		executed[tr.Signal] = true
	}
	for i := range s.rows {
		s.rows[i].Executed = executed[s.rows[i].SignalID]
	}
}

func (s *BTC5MRollingValueAreaReversionStrategy) skipRows() []BTC5MRollingValueAreaReversionSkipRow {
	rows := make([]BTC5MRollingValueAreaReversionSkipRow, 0, len(s.skips))
	for key, count := range s.skips {
		parts := strings.SplitN(key, "|", 2)
		rows = append(rows, BTC5MRollingValueAreaReversionSkipRow{CandidateID: BTC5MRollingValueAreaReversionCandidateID, Split: parts[0], Reason: parts[1], Count: count, MissingDataPolicy: "skip_explicit_missing_or_invalid_rows_no_fill"})
	}
	return rows
}
