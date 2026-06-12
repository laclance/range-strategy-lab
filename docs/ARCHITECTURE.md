# Architecture

The project is intentionally small.

## Layout

```text
range-strategy-lab/
  cmd/rangelab/          CLI entrypoint
  internal/lab/          reusable loading, engine, strategy, and metrics code
  docs/                  project guidance
  results/               generated outputs, ignored by git
```

## Data Flow

1. `cmd/rangelab` parses CLI flags.
2. `internal/lab.LoadCSV` loads candles from a Binance-style or normalized CSV.
3. `lab.RunBacktest` walks candles in order.
4. The configured `lab.Strategy` receives a confirmed candle index.
5. If the strategy returns a signal, the engine enters on the next candle open.
6. The engine exits by stop, target, time stop, or force close at end of data.
7. `lab.SummarizeSplits` writes split metrics by trade close time.

## Strategy Boundary

All strategy work should live behind this interface:

```go
type Strategy interface {
    Name() string
    OnCandle(ctx StrategyContext) (Signal, bool)
}
```

The starter uses `lab.EmptyStrategy`, which returns no trades. Replace that in
`cmd/rangelab/main.go` after adding your own strategy.

## Execution Model

- The strategy sees only candles up to the current confirmed bar.
- Entries are scheduled on the next candle open.
- Only one position may be open at a time.
- Stop checks happen before target checks.
- Costs are applied through `FeePct` and `SlippagePct`.
- Sizing risks `RiskPct` of current equity at the stop, capped by
  `MaxNotionalPct`.

## Extension Points

Good first additions:

- Indicator helpers, such as ATR, Donchian width, Bollinger width, or ADX.
- Detector-only outputs, such as duty cycle and regime duration.
- A simple range detector strategy that returns no trades but writes diagnostics.
- Candidate strategy structs under `internal/lab/strategies/`.
- CLI flags to select strategy and strategy parameters.

Avoid adding live execution or exchange clients until offline evidence is
strong. The current purpose is research, not trading.
