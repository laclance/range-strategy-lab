# Next Codex Brief: Futures Range Rotation Premise Spec

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md.
- Read docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md only for the state
  audit dependency details.
- Read docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md only if router construction
  boundaries need confirmation.
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
- The futures range context router audit was implemented and run.
- Result directory:
  results/futures-range-context-router-audit/
- Generated router audit facts:
  - rule_rows=58;
  - router_rows=29,784;
  - cohort_rows=84;
  - ranking_rows=12;
  - passing_cohorts=3.
- Router row counts:
  - no_trade=13,546;
  - tradable_rotation=1,299;
  - trend_continuation=0;
  - diagnostic_only=14,939;
  - conflicts=0.
- Passing router cohort counts:
  - no_trade=2;
  - tradable_rotation=1;
  - trend_continuation=0.
- Current effective stop state:
  range_context_router_passed_needs_rotation_premise_spec.

Goal:
Author only a materially new futures range rotation premise spec, likely at:
docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md

The spec must:
- stay documentation-only unless the user explicitly asks for implementation;
- start from the passed router cohort:
  range_context_router_v1|15m|h24|tradable_rotation;
- use the router only as closed-candle context, not as an entry signal;
- define what later non-trading audit would have to prove before any entry,
  exit, P&L backtest, optimizer, replay, or walk-forward work exists;
- explain why the premise is materially different from closed/failed families:
  range_occupancy_rotation_v1, hold-inside/midline, breakout-retest/acceptance,
  clean breakout continuation, structured compression, impulse absorption,
  higher-timeframe nested range rotation, range quality/session/failure-mode
  triage by themselves, and legacy spot-only SR timing/compression evidence;
- include explicit stop states for the next audit, including at minimum:
  source/router gap, closed-family reslice, no eligible events, failed no premise,
  passed needs non-trading audit, and rejected as strategy/backtest request.

Boundaries:
- Do not add Go strategy code, entries, exits, P&L strategy backtests, optimizer
  grids, replay, walk-forward logic, strategy packages, paper/testnet/live paths,
  exchange API, credentials, deploy files, source expansion, symbol expansion,
  broad mining, martingale, averaging down, or two-exchange logic.
- Do not turn the 1,299 rotation-routed rows directly into trades.
- Do not retune or rename failed reviewed families under new labels.
- Do not use future labels as premise inputs.
- Do not import old binance-bot strategy/scoring/live code.

Verification for a documentation-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update README.md docs index if a new doc is added.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the spec creates a durable new boundary or
  permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed docs/memory updates after checks pass unless explicitly told
  not to commit.
```
