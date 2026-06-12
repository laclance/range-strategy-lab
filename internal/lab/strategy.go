package lab

type Direction string

const (
	Long  Direction = "long"
	Short Direction = "short"
)

type Signal struct {
	Side        Direction `json:"side"`
	Stop        float64   `json:"stop"`
	Target      float64   `json:"target"`
	MaxHoldBars int       `json:"max_hold_bars"`
	Reason      string    `json:"reason"`
}

type StrategyContext struct {
	Candles []Candle
	Index   int
}

func (c StrategyContext) Candle() Candle {
	return c.Candles[c.Index]
}

type Strategy interface {
	Name() string
	OnCandle(ctx StrategyContext) (Signal, bool)
}

type EmptyStrategy struct{}

func (EmptyStrategy) Name() string {
	return "empty"
}

func (EmptyStrategy) OnCandle(StrategyContext) (Signal, bool) {
	return Signal{}, false
}
