# Next Codex Brief: Futures Range-First Strategy Construction Review Stop

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this stop state:
  - docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project remains offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- Research is not stopped, but automatic reuse or retuning of failed premises
  is stopped.
- The user-approved range-first, BTCUSDT-first construction protocol produced
  the first V1 grammar:
  range_occupancy_rotation_v1.
- That V1 optimizer/backtester has now been implemented, run, and reviewed.
- V1 stop state:
  range_first_strategy_v1_optimizer_failed_no_replay.
- Source and closed UTC resample validation passed:
  - source:
    ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv;
  - loaded candles: 573,984;
  - open-time coverage:
    2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - 15m resample: 191,328 rows through 2026-06-16T23:45:00Z;
  - 1h resample: 47,832 rows through 2026-06-16T23:00:00Z.
- The fixed baseline
  range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005
  lost after costs:
  - full_2021_2026: 43 trades, net P&L -94.208745, PF 0.466232;
  - 2021_2022_stress: 8 trades, net P&L -22.653070, PF 0.464536;
  - 2023_2024_oos: 16 trades, net P&L -38.248119, PF 0.428362;
  - 2025_2026_recent: 19 trades, net P&L -33.307557, PF 0.504958.
- The declared 1,152-row grid had 0 passing configs and no selected config.
- A fixed replay spec, walk-forward, package review, retune, gate relaxation,
  symbol expansion, or live-adjacent path is not authorized from V1.

Goal:
- Review-only stop unless the user explicitly supplies a materially different
  offline range-first premise.
- Do not implement a new strategy, optimizer, replay, walk-forward, grid,
  source expansion, or result-producing run from the existing closed families
  without user input.
- If the user asks what is available, inventory the current docs and present a
  short choice set of materially different offline premises, clearly marking
  closed families as exclusion evidence.

Closed or failed premises in their reviewed forms:
- structured compression;
- breakout-retest/acceptance;
- clean breakout;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation;
- range_occupancy_rotation_v1.

Boundaries:
- Do not retune range_occupancy_rotation_v1 target, stop, max hold, timeframe,
  lookback, occupancy window, occupancy threshold, recapture level, symbol set,
  grid, ranking score, or review gates around the failed result.
- Do not retune structured compression.
- Do not retune breakout_retest_acceptance target, stop, max hold, timeframe,
  side, symbol set, selection rules, or review gates around the failed result.
- Do not reopen the failed 1h structured-compression surface.
- Do not rerun failed clean-breakout, hold-inside/midline, impulse absorption,
  boundary touch rejection, single-candle wick rejection, failed breakout
  re-entry, mature balance persistence, nested range-rotation, or occupancy
  rotation as entries unless the user supplies a materially new data or
  structure premise.
- Do not import old binance-bot strategy/scoring/live code.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Expected review-only shape if the user asks for next options:
- Start from current docs only.
- Treat structured compression walk-forward, breakout-retest baseline,
  higher-timeframe nested range-rotation audit, and occupancy rotation V1
  optimizer as exclusion evidence.
- Separate reusable infrastructure from rejected strategy premises.
- End with either:
  - a user-approved bounded offline spec brief for a materially different
    premise, or
  - a no-implementation stop that says no automatic strategy work is
    authorized.

Verification for any documentation-only update:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- If files are changed, update memory/PROGRESS.md with exact commands and
  factual outcomes.
- Update memory/DECISIONS.md only for a durable user-approved rule.
- Commit completed docs/memory updates and verification evidence after checks
  pass unless explicitly told not to commit.
```
