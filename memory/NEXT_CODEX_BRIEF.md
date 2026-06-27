# Next Codex Brief: Futures Range Router Rotation Premise Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md.
- Read docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md for dependency facts.
- Read docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md only for state-audit
  dependency details if needed.
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
- Router result directory:
  results/futures-range-context-router-audit/
- Router audit facts:
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
- The passed rotation cohort is:
  range_context_router_v1|15m|h24|tradable_rotation
- That cohort had:
  - full_period_rows=892;
  - weakest_split_rows=210;
  - full_expected_route_hit_rate=0.599776;
  - weakest_split_expected_hit_rate=0.554667;
  - full_adverse_route_hit_rate=0.181614;
  - worst_split_adverse_hit_rate=0.247619;
  - dominant_forward_label=contained_rotation at 0.385650.
- The router rotation premise spec was authored at:
  docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md
- Current effective stop state:
  range_router_rotation_premise_spec_ready_for_non_trading_audit.

Goal:
Implement only the explicit, zero-trade, non-trading BTCUSDT futures audit
behind:
  -futures-range-router-rotation-premise-audit

Default output directory when the flag is used and -out-dir is not explicitly
set:
  results/futures-range-router-rotation-premise-audit

The audit must:
- validate the accepted Binance USDT-M futures BTCUSDT 5m source and closed UTC
  15m resample;
- recompute or reuse the reviewed range-state construction loop and context
  router deterministically, not depend on ignored results/ files as runtime
  input;
- require router dependency stop state
  range_context_router_passed_needs_rotation_premise_spec;
- start only from qualifying rows in the passed router cohort:
  range_context_router_v1|15m|h24|tradable_rotation;
- use the router only as closed-candle context, not as an entry signal;
- collapse consecutive eligible router rows in the same range episode into one
  context segment;
- freeze range high, low, midpoint, and quartiles at the segment start using
  only candles known through that closed segment-start candle;
- search only the first 6 closed 15m candles after the segment start for the
  first in-range boundary-reclaim event described in the spec;
- label only forward outcomes over 24 closed 15m bars after the event candle;
- keep forward labels out of context selection, event formation, grouping, and
  skip decisions;
- rank the event premise through the exact count and behavior gates in the
  spec;
- preserve zero-trade common outputs via lab.EmptyStrategy.

Required audit-specific outputs:
- futures_range_router_rotation_premise_sources.csv/json
- futures_range_router_rotation_premise_coverage.csv/json
- futures_range_router_rotation_premise_router_dependency.csv/json
- futures_range_router_rotation_premise_context_segments.csv/json
- futures_range_router_rotation_premise_events.csv/json
- futures_range_router_rotation_premise_outcomes.csv/json
- futures_range_router_rotation_premise_cohorts.csv/json
- futures_range_router_rotation_premise_rankings.csv/json
- futures_range_router_rotation_premise_summary.csv/json
- futures_range_router_rotation_premise_skips.csv/json

Common outputs must stay:
- source_manifest.json
- summary.csv/json
- trades.json with no trades

Stop states:
- range_router_rotation_premise_audit_source_router_gap
- range_router_rotation_premise_audit_rejected_closed_family_reslice
- range_router_rotation_premise_audit_no_eligible_events
- range_router_rotation_premise_audit_failed_no_premise
- range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit
- range_router_rotation_premise_audit_rejected_as_strategy_backtest_request

Boundaries:
- Do not add Go strategy code, entries, exits, P&L strategy backtests, optimizer
  grids, replay, walk-forward logic, strategy packages, paper/testnet/live
  paths, exchange API, credentials, deploy files, source expansion, symbol
  expansion, broad mining, martingale, averaging down, or two-exchange logic.
- Do not turn the 1,299 rotation-routed rows directly into trades.
- Do not retune or rename failed reviewed families under new labels.
- Do not use future labels as premise inputs.
- Do not import old binance-bot strategy/scoring/live code.
- If a requested path tries to build a strategy or backtest from this premise,
  stop with
  range_router_rotation_premise_audit_rejected_as_strategy_backtest_request.

Test plan:
- Add focused unit tests for accepted futures source/router dependency and
  source/router gap behavior.
- Add tests proving spot/path/count/comparison-only rejection where applicable.
- Add event-construction tests proving context segments collapse duplicate
  router rows, frozen bounds use only candles through segment start, events use
  only closed candles through the event candle, and future labels do not enter
  inputs.
- Add label/gate tests for lower and upper boundary reclaim precedence, no
  event, no eligible events, count gates, behavior gates, closed-family reslice
  protection, and stop-state selection.
- Add CLI tests proving the default run writes no new audit artifacts, the
  explicit flag writes all required artifacts to a temp out dir with zero
  trades, spot comparison is rejected, and conflicting flags are rejected.

Closeout verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-router-rotation-premise-audit -out-dir results/futures-range-router-rotation-premise-audit
- wc -l results/futures-range-router-rotation-premise-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Add docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md after the run
  with source facts, router dependency facts, CSV line counts, segment/event
  counts, cohort/ranking counts, failed reasons, final stop state, and
  verification evidence.
- Update README.md docs index if a new review doc is added.
- Update memory/PROGRESS.md with exact commands and factual outcomes.
- Update memory/DECISIONS.md only if the audit creates a durable new boundary
  or permission rule.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next bounded step.
- Commit completed code/docs/memory updates after checks pass unless explicitly
  told not to commit.
```
