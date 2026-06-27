# Progress

This file is the always-read project snapshot. Keep it compact: current state,
latest verification, result paths, and a milestone index. Detailed evidence
belongs in focused `docs/` reviews, generated artifacts under `results/`, and
git history.

## Current State

- Scope is offline Binance USDT-M futures range-strategy research. The default
  CLI remains BTCUSDT `5m` with `lab.EmptyStrategy`; trades remain `0` unless an
  explicit offline research flag is passed.
- Active source contract remains:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed-candle decision semantics, next-bar-open entries, one-position max,
  stop-first ambiguity, cost accounting, split metrics, and source manifests
  remain core infrastructure.
- Research is not stopped, but automatic reuse of failed premises is stopped.
  Reviewed exclusion evidence includes structured compression, breakout-retest
  acceptance, clean breakout continuation, hold-inside/midline, impulse
  absorption, higher-timeframe nested range rotation, `range_occupancy_rotation_v1`,
  and range quality/session/failure-mode triage cohorts in their reviewed forms.
- The latest user-approved direction is a broader construction loop: combine
  range geometry with volatility, trend, impulse, and OHLCV liquidity or
  participation state before any new entry is specified.
- The next implementation-ready doc is
  `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`. It authorizes only a
  non-trading BTCUSDT state audit behind a future explicit flag. It does not
  authorize a strategy, entry, exit, optimizer, replay, walk-forward, symbol
  expansion, source expansion, live/paper/testnet path, exchange API, deploy
  file, martingale, averaging down, or two-exchange logic.
- Parked future directions are documented but not implementation-ready:
  range context router, volatility-aware exits, BTC regime plus ETH/SOL context,
  spread-range/pair-range work, and derivatives context source expansion.
- `memory/NEXT_CODEX_BRIEF.md` is the canonical next-session prompt.

## 2026-06-27

Future range-strategy directions and specs:

- Added research map:
  `docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md`.
- Added next implementation-ready spec:
  `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
- Added parked follow-up specs:
  - `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md`;
  - `docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md`;
  - `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md`;
  - `docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md`;
  - `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md`.
- Updated `README.md` docs index so future sessions can discover the new docs.
- Updated memory and next Codex brief to make the non-trading range-state audit
  the next concrete task.
- Stop state:
  `range_state_construction_loop_spec_ready_for_audit_implementation`.
- This was a connector-only documentation/memory update. No local Go tests,
  `rg`, `git diff --check`, or `git status` were run in this chat. The next
  Codex brief requires those checks during implementation closeout.

Recent failed premise evidence to preserve:

- `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`: source/resampling passed,
  but no range-quality, UTC-session, or quality-plus-session cohort passed the
  declared gates; stop state `range_context_triage_failed_no_strategy_premise`.
- `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`: fixed
  baseline lost after costs and `0` of `1,152` grid rows passed selection gates;
  stop state `range_first_strategy_v1_optimizer_failed_no_replay`.
- `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`:
  ETH/SOL authority stream stayed positive in the frozen replay but failed
  walk-forward robustness; no package approved.
- `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`:
  selected `15m` and `1h` breakout-retest/acceptance baselines failed after
  costs across BTCUSDT, ETHUSDT, and SOLUSDT.
- `docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md`: only
  `3` valid events across the full BTCUSDT sample, so no baseline was approved.

## Milestone Index

Use `README.md` as the full docs index. The most relevant current docs are:

1. `docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md`.
2. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
3. `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`.
4. `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`.
5. `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`.
6. `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`.
7. `memory/NEXT_CODEX_BRIEF.md`.

Historical details remain in the focused docs and git history rather than this
always-read memory file.