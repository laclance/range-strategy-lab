# Next Codex Brief: Futures Range Context Router Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md.
- Read docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md.
- Skim docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md only if implementation
  details from the state audit are needed.
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
- The futures range-state construction loop audit was implemented and run.
- Result directory:
  results/futures-range-state-construction-loop-audit/
- Generated state audit facts:
  - state_rows=29,784;
  - label_rows=89,352;
  - cohort_rows=68,796;
  - ranking_rows=16,335;
  - passing_cohorts=58.
- Passing route counts:
  - no_trade_toxic=52;
  - tradable_rotation_candidate=6;
  - trend_continuation_candidate=0.
- Current effective stop state:
  range_state_construction_loop_audit_passed_needs_router_spec.

Goal:
Implement only the non-trading BTCUSDT futures range context router audit
authorized by docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md and bounded by
docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md.

Expected flag:
-futures-range-context-router-audit

Default result directory when the flag is used:
results/futures-range-context-router-audit/

The router must:
- validate the accepted BTCUSDT Binance USDT-M futures 5m source;
- derive or reuse closed UTC 15m, 1h, and 4h state construction inputs;
- apply deterministic closed-candle route rules only from the passing state
  audit evidence;
- assign exactly one router label per eligible state row:
  no_trade, tradable_rotation, trend_continuation, or diagnostic_only;
- preserve skipped, ambiguous, missing, and conflicting rows with explicit
  reasons;
- keep forward labels out of router inputs;
- keep common source_manifest.json, summary.csv/json, and trades.json
  zero-trade compatible;
- write router-specific sources, coverage, rules, rows, cohorts, rankings,
  summary, and skips CSV/JSON artifacts as specified by the router spec;
- write a review doc with source facts, artifact counts, router label counts,
  passing/failing router cohorts, final stop state, and verification evidence.

Boundaries:
- Do not add entries, exits, P&L strategy backtests, optimizer grids, replay,
  walk-forward logic, strategy package, paper/testnet/live path, exchange API,
  credentials, deploy files, source expansion, symbol expansion, broad mining,
  martingale, averaging down, or two-exchange logic.
- Do not turn the 6 rotation candidates directly into a strategy.
- Do not retune or rename failed reviewed families under new labels.
- Do not use future labels as router features.
- Do not import old binance-bot strategy/scoring/live code.
- Default cmd/rangelab behavior must remain lab.EmptyStrategy unless the explicit
  router audit flag is passed.

Relevant closed or failed premises:
- structured compression;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation;
- range_occupancy_rotation_v1;
- range quality, UTC session, or failure-mode cohorts by themselves;
- legacy spot-only SR timing/compression evidence.

Expected stop states:
- range_context_router_source_gap
- range_context_router_failed_no_actionable_route
- range_context_router_passed_no_trade_filter_only
- range_context_router_passed_needs_rotation_premise_spec
- range_context_router_passed_needs_continuation_premise_spec
- range_context_router_rejected_closed_family_reslice

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-router-audit -out-dir results/futures-range-context-router-audit
- wc -l results/futures-range-context-router-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if a durable rule changes or the review result
  creates a durable no-promotion/no-strategy-change decision.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed code/docs/memory updates and verification evidence after
  checks pass unless explicitly told not to commit.
```
