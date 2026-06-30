package lab

type BTC15MRangeEdgeExhaustionFadeStrategy struct {
	cfg    BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig
	atr    []float64
	splits []Split
	rows   []BTC15MRangeEdgeExhaustionFadeSignalRow
	skips  map[string]int
}

func NewBTC15MRangeEdgeExhaustionFadeStrategy(candles []Candle, cfg BacktestFirstBTC15MRangeEdgeExhaustionFadeConfig, splits []Split) *BTC15MRangeEdgeExhaustionFadeStrategy {
	cfg = cfg.withDefaults()
	return &BTC15MRangeEdgeExhaustionFadeStrategy{cfg: cfg, atr: ATR(candles, cfg.ATRPeriod), splits: splits, skips: map[string]int{}}
}

func (s *BTC15MRangeEdgeExhaustionFadeStrategy) Name() string {
	return BacktestFirstBTC15MRangeEdgeExhaustionFadeName
}
