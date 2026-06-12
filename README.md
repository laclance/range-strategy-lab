# Range Strategy Lab

Standalone BTCUSDT range-strategy starter project.

## Purpose

This folder is designed to be moved out of `binance-bot` and developed as a
fresh project. It keeps only the basic research plumbing:

- Binance-style 5m CSV loading
- confirmed-candle signal flow with next-bar-open entries
- one-position backtest engine
- stop-first ambiguity handling
- fees/slippage, 1% risk sizing, and 1x notional cap
- split metrics by trade close time
- JSON/CSV outputs

It deliberately does not copy the bot's existing strategy scoring, live
execution, deploy scripts, presets, portfolio coordinator, credentials, or
order code.

Read the docs in this order:

1. [docs/PURPOSE.md](docs/PURPOSE.md)
2. [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
3. [docs/STRATEGY_WORKFLOW.md](docs/STRATEGY_WORKFLOW.md)
4. [docs/RESEARCH_HELPERS.md](docs/RESEARCH_HELPERS.md)
5. [docs/VERIFICATION.md](docs/VERIFICATION.md)
6. [CODEX_BRIEF.md](CODEX_BRIEF.md)

## Quick Start

From inside this folder:

```bash
go test ./...

go run ./cmd/rangelab \
  -csv ../data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/smoke
```

In this current machine, use `/usr/local/go/bin/go` if the snap-wrapped `go`
binary complains.

The included strategy is a no-op placeholder, so the smoke run should produce
zero trades. That is intentional: the folder is a clean workbench for building
range logic from scratch.

## Where To Add A Strategy

Implement the `lab.Strategy` interface in a new file, then swap it into
`cmd/rangelab/main.go`:

```go
strategy := lab.EmptyStrategy{}
```

Replace that with your own strategy type once it exists.

The strategy receives a confirmed candle index and returns an optional signal.
The engine enters on the next candle open and manages stop/target exits.

## Data

The CLI accepts either:

- Binance archive shape:
  `open_time,open,high,low,close,volume,close_time,...`
- normalized shape:
  `open_time,open,high,low,close,volume`

Timestamps may be Unix milliseconds or RFC3339.

## Outputs

Each run writes:

- `summary.json`
- `summary.csv`
- `trades.json`

Metrics are split by trade close time:

- `2021_2022_stress`
- `2023_2024_oos`
- `2025_2026_recent`
- `full_2021_2026`

## Suggested First Experiments

- Build a range detector only, and output duty cycle before adding entries.
- Add one simple entry rule at a time.
- Require gross profitability before studying cost-sensitive variants.
- Keep rules closed-candle and next-bar-open until you have a reason not to.
- Track long and short separately from the first run.
