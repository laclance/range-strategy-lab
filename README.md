# Range Strategy Lab

Standalone BTCUSDT range-strategy starter project. The current trading target is
Binance USDT-M futures 5m data; older spot-data outputs are historical context
until rerun and reviewed on futures data.

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

Use the docs list as an index. For routine work, read only the docs relevant to
the current task plus `memory/NEXT_CODEX_BRIEF.md`; for onboarding or broad
review, read the docs in this order:

1. [docs/PURPOSE.md](docs/PURPOSE.md)
2. [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
3. [docs/STRATEGY_WORKFLOW.md](docs/STRATEGY_WORKFLOW.md)
4. [docs/RESEARCH_HELPERS.md](docs/RESEARCH_HELPERS.md)
5. [docs/VERIFICATION.md](docs/VERIFICATION.md)
6. [docs/ENTRY_READINESS_REVIEW.md](docs/ENTRY_READINESS_REVIEW.md)
7. [docs/SR_REJECTION_TIMING_REVIEW.md](docs/SR_REJECTION_TIMING_REVIEW.md)
8. [docs/SR_CONFIRMATION_TIMING_REVIEW.md](docs/SR_CONFIRMATION_TIMING_REVIEW.md)
9. [docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md](docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md)
10. [docs/COMPRESSION_BREAKOUT_REVIEW.md](docs/COMPRESSION_BREAKOUT_REVIEW.md)
11. [docs/RANGE_REGIME_DURABILITY_REVIEW.md](docs/RANGE_REGIME_DURABILITY_REVIEW.md)
12. [docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md](docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md)
13. [docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md](docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md)
14. [docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md](docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md)
15. [docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md](docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md)
16. [docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md](docs/HOLD_INSIDE_MIDLINE_REACTION_REVIEW.md)
17. [memory/NEXT_CODEX_BRIEF.md](memory/NEXT_CODEX_BRIEF.md)

## Quick Start

From inside this folder:

```bash
go test ./...

go run ./cmd/rangelab \
  -csv /absolute/path/to/btcusdt_futures_um_5m_2021_2026.csv \
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

Use Binance USDT-M futures BTCUSDT 5m data for current research. A source
change between spot and futures is a research break: record the CSV path,
coverage, and row count, then rerun/review affected audits before trusting a
verdict for entries.

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
