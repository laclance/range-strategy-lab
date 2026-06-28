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
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md`. It approved only
  a later docs-only zero-trade derivatives source-audit brief and stopped
  before implementation at
  `derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.
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
- The post-rotation premise failure pivot review stopped with
  `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
  No automatic BTCUSDT-only price-only audit is selected.
- User explicitly approved the BTC regime plus ETH/SOL zero-trade context audit
  described in
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`.
  Implementation used only the approved local Binance USDT-M futures `5m`
  BTCUSDT, ETHUSDT, and SOLUSDT files, produced `0` passing context cohorts,
  and closed the path in reviewed zero-trade form. BTCUSDT remains
  market-regime context and diagnostic-only authority; ETHUSDT/SOLUSDT remain
  failed context authority candidates only. This does not authorize entries,
  exits, P&L backtests, optimizer grids, replay, walk-forward,
  paper/testnet/live paths, exchange API, credentials, deploy files, broad
  mining, martingale, averaging down, or two-exchange logic.
- Parked future directions remain documented but not implementation-ready:
  derivatives context is approved only for a later source/alignment
  brief-writing task, spread-range/pair-range is parked behind source/engine
  complexity, and volatility-aware exits remain rejected until a new independent
  entry premise first shows gross edge before costs.
- `memory/NEXT_CODEX_BRIEF.md` is the canonical next-session prompt.

## 2026-06-28

Derivatives market-data context source scope review:

- Added docs-only scope review:
  `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md`.
- Stop state:
  `derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.
- Decision: a later zero-trade derivatives source-audit brief is justified only
  for source and alignment approval, not context-gain implementation.
- The review found no approved durable derivatives market-data rows in the
  lab's current local data scope. `../binance-bot/data/` contained durable
  candle CSVs only, and `../binance-bot/data/raw/` had no files in the reviewed
  search depth.
- Approved local/offline inputs for the later brief are only alignment anchors:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Adjacent `../binance-bot/research/` source-proof artifacts may be referenced
  only as process/source evidence, not as lab input data or strategy evidence.
  Candidate source families for the later brief are mark/index/premium basis
  first, funding second, and aggregate trades only as a high-volume secondary
  source-proof candidate.
- Rejected or blocked from current scope: `/tmp` caches as durable inputs,
  source downloads, live probes, private endpoints, exchange API keys, open
  interest from the current evidence set, long/short ratios without full-era
  archive proof, liquidation/force-order history from current evidence,
  order-book/depth without separate historical archive proof, entries, exits,
  P&L backtests, replay, walk-forward, optimizer grids, and promotion.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a separate zero-trade source-audit
  brief-writing task that stops before code at
  `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.
- Commands run:
  - `rg --files ../binance-bot/data | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade|trade)"`
  - `rg --files ../binance-bot/data | rg -i "(btcusdt|ethusdt|solusdt).*futures.*(5m|1h|15m|4h|um)|futures_um"`
  - `rg --files . | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade)"`
  - `rg -n "funding|open interest|basis|premium|taker|long/short|order-book|order book|derivatives" docs memory README.md AGENTS.md`
  - `rg --files ../binance-bot | rg -i "(funding|open[_-]?interest|\\boi\\b|oi_|basis|premium[_-]?index|premiumIndex|mark[_-]?price|index[_-]?price|taker|long[_-]?short|longshort|depth|order[_-]?book|book[_-]?ticker|liquidat|agg[_-]?trade)"`
  - `rg --files ../binance-bot/data`
  - `ls -lah ../binance-bot/data`
  - `find ../binance-bot/data/raw -maxdepth 3 -type f`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: inventory checks found durable local candle CSVs but
  no approved durable derivatives market-data source rows under
  `../binance-bot/data/`; adjacent source-proof artifacts exist only as
  references; final brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.

BTC regime plus ETH/SOL zero-trade audit implementation:

- Added zero-trade audit implementation behind
  `-futures-btc-regime-eth-sol-context-audit`.
- Review doc:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md`.
- Result directory:
  `results/futures-btc-regime-eth-sol-context-audit/`.
- Stop state:
  `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
- The audit used only the approved local Binance USDT-M futures `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts reproduced: each file had `573,984` loaded candles from
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`; all accepted sorted
  streams had `gap_count=0` and `duplicate_count=0`; zero-volume counts were
  BTC `66`, ETH `47`, SOL `47`; SOL had one physical non-monotonic row and was
  accepted only after sorting.
- Resampled coverage passed for BTC/ETH/SOL `15m`, `1h`, and `4h`: rows were
  `191,328`, `47,832`, and `11,958` respectively for each symbol, with
  `gap_count=0`, `duplicate_count=0`, `missing_child_open_count=0`, and
  `complete=true`.
- Generated artifact counts: BTC state rows `29,784`, ETH/SOL local state rows
  `30,717`, relative-strength rows `30,717`, label rows `92,151`, cohort rows
  `359,055`, ranking rows `160,983`, and passing cohorts `0`.
- Main ranking failure inventory: `142,615` rows failed
  `btc_context_improvement_gate_failed`, `16,095` failed
  `single_split_contribution_above_gate`, `1,627` failed
  `route_rate_gate_failed`, `419` failed `missing_period_split`, and `227`
  failed `inadequate_cohort_count`.
- Anti-leakage and boundary outcomes: BTCUSDT stayed diagnostic market-regime
  context only; ETHUSDT/SOLUSDT stayed zero-trade context authority candidates
  only; forward labels were not used as state/context/gating inputs; common
  outputs remained zero-trade compatible; no entries, exits, P&L backtests,
  replay, walk-forward, source downloads, or strategy promotion were added.
- Commands run:
  - `gofmt -w internal/lab/futures_btc_regime_eth_sol_context_audit.go internal/lab/futures_btc_regime_eth_sol_context_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...`
  - `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-regime-eth-sol-context-audit -out-dir results/futures-btc-regime-eth-sol-context-audit`
  - `wc -l results/futures-btc-regime-eth-sol-context-audit/*.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: tests passed; audit rerun reproduced `source_rows=3`,
  `coverage_rows=9`, `btc_state_rows=29784`, `local_state_rows=30717`,
  `relative_strength_rows=30717`, `label_rows=92151`,
  `cohort_rows=359055`, `ranking_rows=160983`, `passing_cohorts=0`, and stop
  state `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`;
  `wc -l` over CSV artifacts plus common `summary.csv` totaled `703,514`
  lines; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed; pre-commit
  status showed only intended code, docs, and memory changes.

BTC regime plus ETH/SOL zero-trade audit brief:

- Added docs-only brief:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_BRIEF.md`.
- Stop state:
  `btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.
- The brief defines the later audit question: whether BTCUSDT regime buckets
  known at closed decision-candle time improve separation of ETHUSDT/SOLUSDT
  usable, toxic, rotation, continuation, or no-trade local range states.
- BTCUSDT role remains market-regime context and diagnostic-only authority row;
  ETHUSDT/SOLUSDT may become possible authority rows only inside a zero-trade
  context audit, not strategy promotion.
- Allowed source scope remains only the already local Binance USDT-M futures
  `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts preserved from prior validation: each file has `573,984` loaded
  candles from `2021-01-01T00:00:00Z` through
  `2026-06-16T23:55:00Z`; all sorted streams had `gap_count=0` and
  `duplicate_count=0`; zero-volume counts were BTC `66`, ETH `47`, SOL `47`;
  SOL had one physical non-monotonic row and was accepted only after sorting.
- Required anti-leakage rule: forward labels may appear only in label, cohort,
  ranking, and summary artifacts, never premise, state-ID, router, gating, or
  feature-bucket inputs.
- Required common-output rule: `summary.json`, `summary.csv`, and `trades.json`
  must remain zero-trade compatible for the later audit.
- Rejection criteria include closed-family reslice, broad mining, source gap,
  hidden future-label input, structured-compression rescue, ETH/SOL replay,
  BTC promotion, entries, exits, P&L backtests, optimizer grids, replay, or
  walk-forward.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to wait for explicit user approval
  before any implementation.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit `git status --short` showed only
  intended docs and memory changes.

BTC regime plus ETH/SOL context scope review:

- Added docs-only review:
  `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md`.
- Approved `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` as materially
  different enough and tightly scoped enough for a separate zero-trade audit
  brief-writing task, not audit implementation.
- Stop state:
  `btc_regime_eth_sol_context_scope_review_approved_needs_zero_trade_audit_brief`.
- Allowed source scope for the later brief is only the already local Binance
  USDT-M futures `5m` files:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`,
  `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv`, and
  `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`.
- Source facts preserved from prior validation: each file has `573,984` loaded
  candles from `2021-01-01T00:00:00Z` through
  `2026-06-16T23:55:00Z`; all sorted streams had `gap_count=0` and
  `duplicate_count=0`; zero-volume counts were BTC `66`, ETH `47`, SOL `47`;
  SOL had one physical non-monotonic row and was accepted only after sorting.
- BTCUSDT role is market-regime context and diagnostic-only authority row.
  ETHUSDT/SOLUSDT role is possible authority rows only in a later zero-trade
  context audit, not strategy promotion.
- Minimum next audit question: whether BTC regime buckets improve separation of
  ETH/SOL usable, toxic, rotation, continuation, or no-trade range states.
- Rejection criteria preserved: closed-family reslice, broad mining, source
  gap, hidden future-label input, or any move toward entries/backtests.
- Refreshed `memory/NEXT_CODEX_BRIEF.md` to a separate zero-trade audit brief
  task that stops before code at
  `btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.
- Commands run:
  - `wc -l ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv`
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: source files existed with `573,985` CSV lines each
  including headers; reference scan found canonical `memory/NEXT_CODEX_BRIEF.md`
  references and checklist mentions only; `git diff --check` passed;
  pre-commit `git status --short` showed only intended docs and memory changes.

BTC/ETH/SOL context review handoff selection:

- Reviewed `memory/NEXT_CODEX_BRIEF.md` valid choices and selected BTC regime
  plus ETH/SOL context as the next bounded lane to review.
- The next task is a documentation-only scope approval review, not audit
  implementation.
- Recommended ranking preserved in memory:
  1. BTC regime plus ETH/SOL context first;
  2. derivatives market-data context second, pending source/alignment approval;
  3. spread-range source/engine work third, pending engine/source scope;
  4. volatility-aware exits only after a future independent entry premise shows
     gross edge before costs.
- The next brief should create or update a docs-only review that decides whether
  the existing parked BTC/ETH/SOL context spec can become a zero-trade audit
  brief. It must preserve the failed BTCUSDT price-only verdict and avoid
  structured-compression rescue.
- Stop state:
  `btc_regime_eth_sol_context_scope_review_selected_for_next_brief`.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; `git status --short` showed only intended changes
  to `memory/DECISIONS.md`, `memory/NEXT_CODEX_BRIEF.md`, and
  `memory/PROGRESS.md`.

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

## 2026-06-28

Futures range post-rotation premise failure pivot review:

- Added docs-only review:
  `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md`.
- Preserved the failed audit verdict in
  `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md` and kept
  `router_gated_boundary_reclaim_rotation_v1` closed in reviewed form.
- Source facts remain:
  `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
  Binance USDT-M futures `BTCUSDT` `5m`; `573,984` loaded candles;
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
  `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`,
  `comparison_only=false`, `validation_status=accepted`.
- Dependency facts preserved: state audit stopped at
  `range_state_construction_loop_audit_passed_needs_router_spec`; router audit
  stopped at `range_context_router_passed_needs_rotation_premise_spec`; premise
  audit stopped at `range_router_rotation_premise_audit_failed_no_premise`.
- Premise audit facts preserved: context segments `278`, events `97`, outcomes
  `97`, cohort rows `12`, ranking rows `3`, passing cohorts `0`, lower events
  `43`, upper events `54`, midline outcomes `71`, hard-adverse outcomes `22`,
  and chop/no-resolution outcomes `4`.
- Pivot decision: no materially different BTCUSDT-only, candle-price-only
  range-premise audit is worth specifying next from current evidence. Do not
  convert the `278` segments, `97` events, or `1,299` `tradable_rotation`
  router rows into trades.
- Stop state:
  `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
- Current next step: no automatic implementation. A future session should first
  obtain an explicit user scope choice before BTC/ETH/SOL context, derivatives
  context, spread-range source/engine work, or a no-further-audit stop is
  pursued.
- Commands run:
  - `rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md`
  - `git diff --check`
  - `git status --short`
- Verification outcomes: brief-reference scan found canonical
  `memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only;
  `git diff --check` passed; pre-commit status showed only intended docs and
  memory changes.

## Milestone Index

Use `README.md` as the full docs index. The most relevant current docs are:

1. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md`.
2. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md`.
3. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md`.
4. `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md`.
5. `docs/FUTURES_RANGE_STRATEGY_FUTURE_DIRECTIONS_RESEARCH_MAP.md`.
6. `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_SPEC.md`.
7. `docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.
8. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md`.
9. `docs/FUTURES_RANGE_CONTEXT_ROUTER_SPEC.md`.
10. `docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_SPEC.md`.
11. `docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`.
12. `docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md`.
13. `docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md`.
14. `docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md`.
15. `memory/NEXT_CODEX_BRIEF.md`.

Historical details remain in the focused docs and git history rather than this
always-read memory file.
