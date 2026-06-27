# Next Codex Brief: Futures Range Post-Rotation Premise Failure Pivot

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md.
- Read docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md only for the closed
  premise dependency details.
- Read docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md and
  docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md only if dependency facts
  need confirmation.
- Read docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md if a
  bounded post-failure direction is needed.
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
- Router facts:
  - rule_rows=58;
  - router_rows=29,784;
  - cohort_rows=84;
  - ranking_rows=12;
  - passing_cohorts=3;
  - no_trade=13,546;
  - tradable_rotation=1,299;
  - trend_continuation=0;
  - diagnostic_only=14,939;
  - conflicts=0.
- The router rotation premise spec named:
  router_gated_boundary_reclaim_rotation_v1.
- The implemented router rotation premise audit result directory is:
  results/futures-range-router-rotation-premise-audit/
- Premise audit facts:
  - context_segments=278;
  - events=97;
  - outcomes=97;
  - cohort_rows=12;
  - ranking_rows=3;
  - passing_cohorts=0;
  - lower_events=43;
  - upper_events=54;
  - full_midline_outcomes=71;
  - hard_adverse_outcomes=22;
  - chop_or_no_resolution_outcomes=4.
- Top failure reasons:
  inadequate_event_count,inadequate_split_event_count,single_split_contribution_above_gate,behavior_gate_failed
- Current effective stop state:
  range_router_rotation_premise_audit_failed_no_premise.

Goal:
Author only a documentation-only post-failure pivot review/spec that decides the
next bounded research direction after the failed router-gated boundary-reclaim
rotation premise.

The doc should:
- preserve the failed verdict in
  docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md;
- state that router_gated_boundary_reclaim_rotation_v1 is closed in reviewed
  form;
- not convert the 278 segments, 97 events, or 1,299 tradable_rotation router
  rows into trades;
- identify whether any materially different non-trading premise or context audit
  is worth specifying next, or stop with no next audit if not;
- distinguish any proposed direction from closed families:
  range_occupancy_rotation_v1, hold-inside/midline, breakout-retest/acceptance,
  clean breakout continuation, structured compression, impulse absorption,
  higher-timeframe nested range rotation, range quality/session/failure-mode
  triage by themselves, and router_gated_boundary_reclaim_rotation_v1;
- remain documentation-only unless the user explicitly asks for implementation.

Boundaries:
- Do not add Go code, strategy code, entries, exits, P&L strategy backtests,
  optimizer grids, replay, walk-forward logic, strategy packages,
  paper/testnet/live paths, exchange API, credentials, deploy files, source
  expansion, symbol expansion, broad mining, martingale, averaging down, or
  two-exchange logic.
- Do not retune, rename, relax gates for, or directly repackage the failed
  router-gated boundary-reclaim rotation premise.
- Do not import old binance-bot strategy/scoring/live code.
- Do not use future labels as premise inputs.

Verification for a documentation-only closeout:
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Add or update a focused docs/ review or spec only if a materially new bounded
  direction is identified.
- Update README.md docs index if a new doc is added.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the review creates a durable boundary,
  no-go rule, or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed docs/memory updates after checks pass unless explicitly told
  not to commit.
```
