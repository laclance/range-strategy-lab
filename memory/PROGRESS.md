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
- The latest completed research doc is
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`. It implemented
  the zero-trade router rotation premise audit and stopped at
  `range_router_rotation_premise_audit_failed_no_premise`.
- The prior dependency docs are
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md` and
  `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
- The router produced `58` rules from passing state-audit rankings, assigned
  `13,546` `no_trade`, `1,299` `tradable_rotation`, `0`
  `trend_continuation`, and `14,939` `diagnostic_only` rows, and found `3`
  passing router cohorts: `2` no-trade and `1` rotation. The follow-up router
  rotation premise audit collapsed that rotation context to `278` segments and
  `97` valid boundary-reclaim events, found `0` passing premise cohorts, and
  closed `router_gated_boundary_reclaim_rotation_v1` in reviewed form. It does
  not authorize a strategy, entry, exit, optimizer, replay, walk-forward, symbol
  expansion, source expansion, live/paper/testnet path, exchange API, deploy
  file, martingale, averaging down, or two-exchange logic.
- Parked future directions are documented but not implementation-ready:
  volatility-aware exits, BTC regime plus ETH/SOL context, spread-range/pair-range
  work, and derivatives context source expansion.
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

Futures range-state construction loop audit implementation:

- Added `-futures-range-state-construction-loop-audit`.
- Added audit implementation and tests for source/coverage gates, closed-candle
  state construction, forward-label separation, route cohort gates, stop-state
  selection, CLI artifacts, spot rejection, and conflicting flag rejection.
- Result directory:
  `results/futures-range-state-construction-loop-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample rows: `15m=191,328`, `1h=47,832`, `4h=11,958`; all
  accepted with `gap_count=0`, `duplicate_count=0`,
  `missing_child_open_count=0`, and complete buckets.
- Generated artifact counts: states `29,784`, labels `89,352`, cohorts
  `68,796`, rankings `16,335`, passing cohorts `58`.
- Passing route counts: `no_trade_toxic=52`,
  `tradable_rotation_candidate=6`, `trend_continuation_candidate=0`.
- Common outputs stayed zero-trade compatible; `trades.json` contains no trades.
- Review doc:
  `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md`.
- Stop state:
  `range_state_construction_loop_audit_passed_needs_router_spec`.
- Commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-state-construction-loop-audit -out-dir results/futures-range-state-construction-loop-audit`
  - `wc -l results/futures-range-state-construction-loop-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `state_rows=29784`,
  `label_rows=89352`, `cohort_rows=68796`, `ranking_rows=16335`,
  `passing_cohorts=58`, and stop state
  `range_state_construction_loop_audit_passed_needs_router_spec`;
  CSV line-count check totaled `204,345`; `rg` found only canonical
  `memory/NEXT_CODEX_BRIEF.md` references; `git diff --check` passed; pre-commit
  status showed only intended code, docs, and memory changes.
- That next step was completed by the router audit below; the state audit still
  does not authorize an entry strategy by itself.

Futures range context router audit implementation:

- Added `-futures-range-context-router-audit`.
- Added audit implementation and tests for source/coverage inheritance, state
  rollup rule matching, forward-label separation, conflicting-rule preservation,
  router cohort gates, stop-state selection, CLI artifacts, spot rejection, and
  conflicting flag rejection.
- Result directory:
  `results/futures-range-context-router-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC resample rows remained `15m=191,328`, `1h=47,832`,
  `4h=11,958`; all accepted with `gap_count=0`, `duplicate_count=0`,
  `missing_child_open_count=0`, and complete buckets.
- Generated artifact counts: rules `58`, router rows `29,784`, cohorts `84`,
  rankings `12`, passing cohorts `3`.
- Router row counts: `no_trade=13,546`, `tradable_rotation=1,299`,
  `trend_continuation=0`, `diagnostic_only=14,939`, conflicts `0`.
- Passing router cohorts: `no_trade=2`, `tradable_rotation=1`,
  `trend_continuation=0`.
- Common outputs stayed zero-trade compatible; `trades.json` contains no trades.
- Review doc:
  `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
- Stop state:
  `range_context_router_passed_needs_rotation_premise_spec`.
- Commands run:
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-router-audit -out-dir results/futures-range-context-router-audit`
  - `wc -l results/futures-range-context-router-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `rule_rows=58`,
  `router_rows=29784`, `cohort_rows=84`, `ranking_rows=12`,
  `passing_cohorts=3`, and stop state
  `range_context_router_passed_needs_rotation_premise_spec`; CSV line-count
  check totaled `30,024`; `rg` found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed;
  pre-commit status showed only intended code, docs, and memory changes.
- That next step was completed by the router rotation premise spec below; the
  router result still does not authorize entries, exits, backtests, optimizers,
  replay, or walk-forward.

Futures range router rotation premise spec:

- Added docs-only spec:
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md`.
- The spec starts from the passed router cohort
  `range_context_router_v1|15m|h24|tradable_rotation` and names the non-trading
  event premise `router_gated_boundary_reclaim_rotation_v1`.
- The router remains closed-candle context only, not an entry signal. The
  `1,299` `tradable_rotation` rows must not be converted directly into trades.
- The next permitted implementation is only a zero-trade non-trading audit
  behind `-futures-range-router-rotation-premise-audit`, with default output
  directory `results/futures-range-router-rotation-premise-audit/`.
- The spec requires context segments from `15m` `h24` `tradable_rotation`
  matched rules, frozen range bounds at the segment start, a subsequent closed
  in-range boundary-reclaim event, and labels beginning strictly after the event
  candle.
- Required next-audit stop states include
  `range_router_rotation_premise_audit_source_router_gap`,
  `range_router_rotation_premise_audit_rejected_closed_family_reslice`,
  `range_router_rotation_premise_audit_no_eligible_events`,
  `range_router_rotation_premise_audit_failed_no_premise`,
  `range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit`,
  and
  `range_router_rotation_premise_audit_rejected_as_strategy_backtest_request`.
- Stop state:
  `range_router_rotation_premise_spec_ready_for_non_trading_audit`.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.
- Current next step: implement only the non-trading audit from the premise spec;
  do not add entries, exits, P&L backtests, optimizers, replay, walk-forward,
  paper/testnet/live paths, source expansion, or symbol expansion.

Futures range router rotation premise audit implementation:

- Added `-futures-range-router-rotation-premise-audit`.
- Added audit implementation and tests for source/router dependency checks,
  router context collapse, closed-candle event construction, outcome labeling,
  gate failures, stop-state selection, CLI artifacts, spot rejection, default
  output behavior, and conflicting flag rejection.
- Result directory:
  `results/futures-range-router-rotation-premise-audit/`.
- Source facts: `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Closed UTC `15m` resample rows: `191,328`, from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:45:00Z`, accepted with
  `gap_count=0`, `duplicate_count=0`, `missing_child_open_count=0`, and
  `complete=true`.
- Router dependency passed with stop state
  `range_context_router_passed_needs_rotation_premise_spec`; the required
  cohort was `range_context_router_v1|15m|h24|tradable_rotation`.
- Generated artifact counts: context segments `278`, events `97`, outcomes
  `97`, cohorts `12`, rankings `3`, passing cohorts `0`.
- Full-period event inventory: lower events `43`, upper events `54`, midline
  outcomes `71`, hard adverse outcomes `22`, chop/no-resolution outcomes `4`.
- Top failure reasons:
  `inadequate_event_count,inadequate_split_event_count,single_split_contribution_above_gate,behavior_gate_failed`.
- Common outputs stayed zero-trade compatible; `summary.csv` has `0` trades in
  every split/side row.
- Review doc:
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`.
- Stop state:
  `range_router_rotation_premise_audit_failed_no_premise`.
- Commands run:
  - `gofmt -w internal/lab/futures_range_context_router_audit.go internal/lab/futures_range_router_rotation_premise_audit.go internal/lab/futures_range_router_rotation_premise_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-router-rotation-premise-audit -out-dir results/futures-range-router-rotation-premise-audit`
  - `wc -l results/futures-range-router-rotation-premise-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced
  `context_segments=278`, `events=97`, `outcomes=97`, `cohort_rows=12`,
  `ranking_rows=3`, `passing_cohorts=0`, and stop state
  `range_router_rotation_premise_audit_failed_no_premise`; CSV line-count check
  totaled `529`; `rg` found canonical `memory/NEXT_CODEX_BRIEF.md` references
  and checklist mentions only; `git diff --check` passed; pre-commit status
  showed only intended code, docs, and memory changes.
- Current next step: do not build entries, exits, P&L backtests, optimizers,
  replay, walk-forward, or trigger audits from this premise. Any follow-up must
  be a materially different non-trading premise or context audit.

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
2. `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`.
3. `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md`.
4. `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
5. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md`.
6. `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md`.
7. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
8. `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`.
9. `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`.
10. `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`.
11. `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`.
12. `memory/NEXT_CODEX_BRIEF.md`.

Historical details remain in the focused docs and git history rather than this
always-read memory file.
