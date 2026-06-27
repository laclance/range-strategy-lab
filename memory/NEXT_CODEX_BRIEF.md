# Next Codex Brief: Futures Range-State Construction Loop Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md.
- Read docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md.
- Skim only these exclusion reviews as needed:
  - docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md
  - docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
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
- The latest range-context triage audit passed source/resample validation but
  failed with no gated strategy premise.
- The latest strategy grammar, range_occupancy_rotation_v1, failed optimizer
  review with zero selectable configs.
- Current effective stop state:
  range_state_construction_loop_spec_ready_for_audit_implementation.

Closed or failed premises:
- structured compression;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation;
- range_occupancy_rotation_v1;
- range quality, UTC session, or failure-mode cohorts by themselves;
- legacy spot-only SR timing/compression evidence.

Goal:
Implement the non-trading BTCUSDT futures range-state construction loop audit
specified in docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md.

Add an explicit CLI flag:
-futures-range-state-construction-loop-audit

Default result directory when the flag is used:
results/futures-range-state-construction-loop-audit/

The implementation must:
- validate the accepted BTCUSDT Binance USDT-M futures 5m source;
- derive closed UTC 15m, 1h, and 4h resamples from the 5m parent;
- validate expected resample coverage and complete child buckets;
- create eligible mature active range state rows using closed-candle data only;
- compute range geometry, volatility, trend, impulse, and OHLCV participation
  proxy features;
- bucket those features into deterministic inspectable state IDs;
- create forward labels only as labels, never as features;
- summarize route candidates:
  tradable_rotation_candidate, trend_continuation_candidate, no_trade_toxic,
  diagnostic_only;
- apply the predeclared count, split-stability, useful-rate, toxic-rate, and
  closed-family protection gates;
- keep common source_manifest.json, summary.csv/json, and trades.json
  zero-trade compatible;
- write the spec-required CSV/JSON artifacts:
  - futures_range_state_construction_loop_sources.csv/json
  - futures_range_state_construction_loop_coverage.csv/json
  - futures_range_state_construction_loop_feature_windows.csv/json
  - futures_range_state_construction_loop_states.csv/json
  - futures_range_state_construction_loop_labels.csv/json
  - futures_range_state_construction_loop_cohorts.csv/json
  - futures_range_state_construction_loop_rankings.csv/json
  - futures_range_state_construction_loop_summary.csv/json
  - futures_range_state_construction_loop_skips.csv/json
- write docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md with source facts,
  resample facts, artifact line counts, state counts, top route cohorts, failed
  cohort reasons, final stop state, and verification evidence.

Boundaries:
- Do not add entries, exits, P&L strategy backtests, optimizer grids, replay,
  walk-forward logic, strategy package, paper/testnet/live path, exchange API,
  credentials, deploy files, source expansion, symbol expansion, broad mining,
  martingale, averaging down, or two-exchange logic.
- Do not retune or rename failed reviewed families under new labels.
- Do not use future labels as features.
- Do not import old binance-bot strategy/scoring/live code.
- Default cmd/rangelab behavior must remain lab.EmptyStrategy unless the explicit
  audit flag is passed.

Expected stop states:
- range_state_construction_loop_source_gap
- range_state_construction_loop_no_eligible_states
- range_state_construction_loop_audit_failed_no_usable_state
- range_state_construction_loop_audit_passed_no_trade_filter_only
- range_state_construction_loop_audit_passed_needs_router_spec
- range_state_construction_loop_audit_passed_needs_strategy_premise_spec
- range_state_construction_loop_rejected_closed_family_reslice

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-state-construction-loop-audit -out-dir results/futures-range-state-construction-loop-audit
- wc -l results/futures-range-state-construction-loop-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if a durable rule changes or the review result
  creates a durable no-promotion/no-strategy-change decision.
- Commit completed code/docs/memory updates and verification evidence after
  checks pass unless explicitly told not to commit.
```
