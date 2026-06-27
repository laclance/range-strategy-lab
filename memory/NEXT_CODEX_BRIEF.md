# Next Codex Brief: Futures Range Scope Choice After BTCUSDT Price-Only Stop

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md.
- Read docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md only for the
  failed premise evidence.
- Read docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md only if
  the user asks which parked direction exists.
- Read parked direction specs only for a user-selected lane:
  - docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md;
  - docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md;
  - docs/FUTURES_SPREAD_RANGE_STRATEGY_SPEC.md;
  - docs/FUTURES_VOLATILITY_AWARE_EXIT_MODEL_SPEC.md.
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- The accepted BTCUSDT futures 5m source remains:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- Source facts:
  - loaded candles: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - comparison_only=false, validation_status=accepted.
- The futures range state construction loop audit stopped at:
  range_state_construction_loop_audit_passed_needs_router_spec.
- The futures range context router audit stopped at:
  range_context_router_passed_needs_rotation_premise_spec.
- The router rotation premise audit stopped at:
  range_router_rotation_premise_audit_failed_no_premise.
- The post-rotation premise failure pivot review stopped at:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- The reviewed premise router_gated_boundary_reclaim_rotation_v1 is closed in
  reviewed form.
- Do not convert the 278 premise context segments, 97 boundary-reclaim events,
  or 1,299 tradable_rotation router rows into trades.
- No materially different BTCUSDT-only, candle-price-only range-premise audit is
  selected from current evidence.

Goal:
Do not implement anything automatically. First get an explicit user scope
choice, then act only within that selected scope.

If the user has not already chosen a lane, respond with a concise scope-choice
question and do not edit files. The valid choices are:
- stop with no further range audit for now;
- documentation-only approval review for BTC/ETH/SOL context scope;
- documentation-only approval review for derivatives market-data context scope;
- documentation-only approval review for spread-range source/engine scope;
- documentation-only exit-model revisit, only if the user also supplies a new
  independent entry premise with gross-edge evidence.

Boundaries:
- Do not add Go code, strategy code, entries, exits, P&L strategy backtests,
  optimizer grids, replay, walk-forward logic, strategy packages,
  paper/testnet/live paths, exchange API, credentials, deploy files, source
  expansion, symbol expansion, broad mining, martingale, averaging down, or
  two-exchange logic without a new explicit user-approved scope brief.
- Do not retune, rename, relax gates for, or directly repackage the failed
  router-gated boundary-reclaim rotation premise.
- Do not import old binance-bot strategy/scoring/live code.
- Do not use future labels as premise inputs.

If the user chooses a documentation-only scope approval review:
- Keep it docs/memory only.
- Preserve the failed BTCUSDT price-only verdict.
- Explain why the chosen direction is materially different from closed families:
  range_occupancy_rotation_v1, hold-inside/midline, breakout-retest/acceptance,
  clean breakout continuation, structured compression, impulse absorption,
  higher-timeframe nested range rotation, range quality/session/failure-mode
  triage by themselves, and router_gated_boundary_reclaim_rotation_v1.
- End with either no next audit or a zero-trade audit brief that still requires
  explicit user approval before implementation.

Verification for a documentation-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update README.md docs index if a new doc is added.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the review creates a durable boundary,
  no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed docs/memory updates after checks pass unless explicitly told
  not to commit.
```
