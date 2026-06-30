package lab

import (
	"fmt"
	"math"
	"strings"
)

type BTC15MPreviousDayRangeReversionStrategy struct {
	cfg       BacktestFirstBTC15MPreviousDayRangeReversionConfig
	atr       []float64
	splits    []Split
	dayRanges map[string]btc15MDayRange
	rows      []BTC15MPreviousDayRangeReversionSignalRow
	skips     map[string]int
}

func RunBacktestFirstBTC15MPreviousDayRangeReversion(candles []Candle, manifest SourceManifest, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig, btCfg BacktestConfig, splits []Split) (BacktestFirstBTC15MPreviousDayRangeReversionResult, error) {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if btCfg.MaxHoldBars <= 0 {
		btCfg.MaxHoldBars = cfg.MaxHoldBars
	}
	source := btc15MPrevDaySourceRow(manifest, cfg)
	result := BacktestFirstBTC15MPreviousDayRangeReversionResult{SourceRows: []BTC15MPreviousDayRangeReversionSourceRow{source}}
	resampled, coverage := btc15MPrevDayResample(candles, cfg, source.SourceFactsPass)
	result.CoverageRows = []BTC15MPreviousDayRangeReversionCoverageRow{coverage}
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = btc15MPrevDayFalsification(BTC15MPreviousDayRangeReversionFalsification{SourceResamplePass: false, LeakagePass: true, SideReportingPass: true}, cfg)
		result.StopState = result.Falsification.StopState
		return result, nil
	}
	strategy := NewBTC15MPreviousDayRangeReversionStrategy(resampled, cfg, splits)
	bt := RunBacktest(resampled, strategy, btCfg)
	strategy.markExecuted(bt.Trades)
	result.Trades = bt.Trades
	result.SignalRows = strategy.rows
	result.SkipRows = strategy.skipRows()
	result.TradeRows = btc15MPrevDayTradeRows(bt.Trades, splits)
	result.SummaryRows = SummarizeSplits(bt.Trades, btCfg.StartBalance, splits)
	result.Falsification = btc15MPrevDayFalsification(btc15MPrevDayEvaluate(source, coverage, result.SummaryRows, cfg), cfg)
	result.StopState = result.Falsification.StopState
	return result, nil
}

func NewBTC15MPreviousDayRangeReversionStrategy(candles []Candle, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig, splits []Split) *BTC15MPreviousDayRangeReversionStrategy {
	cfg = cfg.withDefaults()
	return &BTC15MPreviousDayRangeReversionStrategy{cfg: cfg, atr: ATR(candles, cfg.ATRPeriod), splits: splits, dayRanges: btc15MPrevDayRanges(candles), skips: map[string]int{}}
}

func (s *BTC15MPreviousDayRangeReversionStrategy) Name() string {
	return BacktestFirstBTC15MPreviousDayRangeReversionName
}

func (s *BTC15MPreviousDayRangeReversionStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	d := ctx.Index
	candles := ctx.Candles
	decision := candles[d]
	split := btc15MPrevDaySplit(decision.CloseTime, s.splits)
	if d == 0 {
		s.skips[split+"|missing_prior_day"]++
		return Signal{}, false
	}
	if d-1 >= len(s.atr) || !validNumber(s.atr[d-1]) || s.atr[d-1] <= 0 {
		s.skips[split+"|missing_prior_atr"]++
		return Signal{}, false
	}
	priorKey := decision.OpenTime.UTC().AddDate(0, 0, -1).Format("2006-01-02")
	prior, ok := s.dayRanges[priorKey]
	if !ok || !prior.complete || prior.width <= 0 {
		s.skips[split+"|missing_complete_prior_day"]++
		return Signal{}, false
	}
	if !btc15MPrevDayCurrentDayStillInside(candles, d, prior) {
		return Signal{}, false
	}
	side, ok := btc15MPrevDaySide(decision.Close, prior, s.cfg.OuterZonePct)
	if !ok {
		return Signal{}, false
	}
	priorATR := s.atr[d-1]
	stop := prior.low - s.cfg.StopATRMultiple*priorATR
	if side == Short {
		stop = prior.high + s.cfg.StopATRMultiple*priorATR
	}
	row := BTC15MPreviousDayRangeReversionSignalRow{SignalID: fmt.Sprintf("%s_%06d", BTC15MPreviousDayRangeReversionCandidateID, len(s.rows)+1), CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Split: split, DecisionIndex: d, DecisionCloseTime: decision.CloseTime.UTC().Format(timeLayout), DecisionClose: decision.Close, Side: side, TimingLabel: "next_15m_open", PriorDay: priorKey, PriorDayHigh: prior.high, PriorDayLow: prior.low, PriorDayMidpoint: prior.mid, PriorDayRangeWidth: prior.width, OuterZonePct: s.cfg.OuterZonePct, PriorATR14: priorATR, Stop: stop, Target: prior.mid, MaxHoldBars: s.cfg.MaxHoldBars, ForwardLabelsAsInput: false, UsesFutureRowsForFeatures: false, DerivativesVetoUsed: false, OptimizerSelectionUsed: false}
	s.rows = append(s.rows, row)
	return Signal{Side: side, Stop: stop, Target: prior.mid, MaxHoldBars: s.cfg.MaxHoldBars, Reason: row.SignalID}, true
}

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

func btc15MPrevDayCurrentDayStillInside(candles []Candle, d int, prior btc15MDayRange) bool {
	day := candles[d].OpenTime.UTC().Format("2006-01-02")
	for i := d - 1; i >= 0 && candles[i].OpenTime.UTC().Format("2006-01-02") == day; i-- {
		if candles[i].Close < prior.low || candles[i].Close > prior.high {
			return false
		}
	}
	return true
}

func (s *BTC15MPreviousDayRangeReversionStrategy) markExecuted(trades []Trade) {
	executed := map[string]bool{}
	for _, tr := range trades {
		executed[tr.Signal] = true
	}
	for i := range s.rows {
		s.rows[i].Executed = executed[s.rows[i].SignalID]
	}
}

func (s *BTC15MPreviousDayRangeReversionStrategy) skipRows() []BTC15MPreviousDayRangeReversionSkipRow {
	rows := make([]BTC15MPreviousDayRangeReversionSkipRow, 0, len(s.skips))
	for key, count := range s.skips {
		parts := strings.SplitN(key, "|", 2)
		rows = append(rows, BTC15MPreviousDayRangeReversionSkipRow{CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Split: parts[0], Reason: parts[1], Count: count, MissingDataPolicy: "skip_explicit_missing_or_invalid_rows_no_fill"})
	}
	return rows
}
