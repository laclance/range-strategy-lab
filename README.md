# Range Strategy Lab

Standalone Binance USDT-M futures range-strategy starter project. The current
implemented CLI default remains BTCUSDT futures 5m data. Research is not
stopped, but automatic reuse of failed premises is stopped. The user-approved
range-first, BTCUSDT-first construction protocol produced and tested the first
V1 grammar, BTCUSDT range occupancy rotation, and that optimizer review failed
with no selected config for fixed replay. The next authorized path is a
non-trading futures range-context triage audit that assesses range quality,
session behavior, and failure modes before any new strategy grammar. Older
spot-data outputs are historical context unless a futures rerun explicitly
revalidates a specific conclusion.

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
17. [docs/FUTURES_DATA_IMPACT_REVIEW.md](docs/FUTURES_DATA_IMPACT_REVIEW.md)
18. [docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md](docs/FUTURES_MIDLINE_TOUCH_PROTOTYPE_REVIEW.md)
19. [docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md](docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md)
20. [docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md](docs/FUTURES_IMPULSE_ABSORPTION_AUDIT_REVIEW.md)
21. [docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md](docs/FUTURES_SCOPE_PIVOT_REVIEW_SPEC.md)
22. [docs/FUTURES_SCOPE_PIVOT_REVIEW.md](docs/FUTURES_SCOPE_PIVOT_REVIEW.md)
23. [docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md](docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md)
24. [docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md](docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md)
25. [docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md](docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md)
26. [docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md](docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md)
27. [docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md](docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md)
28. [docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md](docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md)
29. [docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md](docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md)
30. [docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md)
31. [docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md)
32. [docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md)
33. [docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md](docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md)
34. [docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md)
35. [docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md)
36. [docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md)
37. [docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md](docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md)
38. [docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md](docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md)
39. [docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md](docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md)
40. [docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md](docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md)
41. [docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_SPEC.md](docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_SPEC.md)
42. [memory/NEXT_CODEX_BRIEF.md](memory/NEXT_CODEX_BRIEF.md)

## Quick Start

From inside this folder:

```bash
go test ./...

go run ./cmd/rangelab \
  -source-product binance-usdm-futures \
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

Use Binance USDT-M futures data for current research. The implemented CLI still
defaults to BTCUSDT 5m; `docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md`
records the first local BTC/ETH/SOL futures source-validation and discovery
audit. A source change between spot and futures is a research break: record the
CSV path, coverage, and row count, then rerun/review affected audits before
trusting a verdict for entries.

The default CLI source is:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

That file is the current full-history Binance USDT-M futures 5m source visible
from this checkout: `573,985` CSV lines including header, `573,984` loaded
candles, spanning open times `2021-01-01T00:00:00Z` through
`2026-06-16T23:55:00Z`. The accepted manifest has `gap_count=0`,
`duplicate_count=0`, and `zero_volume_count=66`.

When passing any non-default CSV path, also pass `-source-product`. Valid
values are:

- `binance-usdm-futures`
- `binance-spot`

Spot CSVs are rejected unless the run is explicitly marked as a comparison with
both `-source-product binance-spot` and `-allow-spot-comparison`. A spot
comparison cannot satisfy futures promotion or entry gates.

Source validation rejects non-BTCUSDT/non-5m paths, spot-looking paths in
futures runs, gaps, duplicates, irregular 5m cadence, non-positive prices,
non-finite values, negative volume, and invalid high/low containment. Zero
volume is allowed and counted because official exchange archives can contain
closed flat candles with no trades.

The CLI accepts either:

- Binance archive shape:
  `open_time,open,high,low,close,volume,close_time,...`
- normalized shape:
  `open_time,open,high,low,close,volume`

Timestamps may be Unix milliseconds or RFC3339.

## Outputs

Each run writes:

- `source_manifest.json`
- `summary.json`
- `summary.csv`
- `trades.json`

`source_manifest.json` records the source path, venue, product, symbol,
interval, row count, first/last open time, schema, timestamp semantics, finality
rule, gap/duplicate/zero-volume counts, comparison-only status, and validation
status.

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
