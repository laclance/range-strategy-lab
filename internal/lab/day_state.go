package lab

type BTC15MPreviousDayRangeReversionStrategy struct {
	cfg       BacktestFirstBTC15MPreviousDayRangeReversionConfig
	atr       []float64
	splits    []Split
	dayRanges map[string]btc15MDayRange
	rows      []BTC15MPreviousDayRangeReversionSignalRow
	skips     map[string]int
}

func NewBTC15MPreviousDayRangeReversionStrategy(candles []Candle, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig, splits []Split) *BTC15MPreviousDayRangeReversionStrategy {
	cfg = cfg.withDefaults()
	return &BTC15MPreviousDayRangeReversionStrategy{cfg: cfg, atr: ATR(candles, cfg.ATRPeriod), splits: splits, dayRanges: btc15MPrevDayRanges(candles), skips: map[string]int{}}
}

func (s *BTC15MPreviousDayRangeReversionStrategy) Name() string {
	return BacktestFirstBTC15MPreviousDayRangeReversionName
}
