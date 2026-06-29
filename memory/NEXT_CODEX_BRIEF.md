# Next Codex Brief: Derivatives Context Strategy-Premise Spec Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do NOT implement anything in this task without explicit user approval. The
zero-trade derivatives context audit has passed only as context-separation
evidence; the next step is the user approving or rejecting a docs-only
strategy-premise spec. If approval is not explicitly given in the session, make
no Go code changes, no CLI flag, no generated result directory, no audit run, no
source download, no network request, no data write, and no strategy/P&L work;
report the waiting state and stop.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md for the passed context
  audit result.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md for
  the audit scope and anti-leakage contract.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md only for the
  source/alignment facts if needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed (0 cohorts), closed at:
  btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context.
- Derivatives source materialization passed; 729 checksum-verified zips, 9
  normalized CSVs, 5 manifests under ../binance-bot/data/derivatives/.
- The derivatives zero-trade source audit passed at:
  derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief.
- The derivatives zero-trade context-audit brief was approved and implemented.
  The implementation passed at:
  derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec.
- Result path: results/futures-derivatives-context-audit/.
- Result counts: source_rows=12, coverage_rows=18,
  basis_feature_rows=83,004, local_state_rows=83,640, label_rows=249,012,
  cohort_rows=512,190, ranking_rows=181,827, missingness_rows=36,
  passing_cohorts=6, trades=0.
- The pass is narrow: all 6 passing cohorts are BTCUSDT 15m. Five are
  no-trade/toxic context cohorts and one is a rotation candidate. ETHUSDT and
  SOLUSDT produced 0 passing cohorts. Rows ranked after the passing set mostly
  failed inadequate-count gates.
- The context audit was zero-trade only. It did not test entries, exits, fills,
  P&L, replay, walk-forward, or basis tradability.

Passing context cohorts to preserve:
1. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none,
   bucket=basis_discount_small + premium_discount_small,
   rows=515, weakest_split_rows=110, full_toxic=0.732039,
   worst_split_toxic=0.800000.
2. BTCUSDT 15m h48 no_trade_toxic:
   same local state, bucket=basis_discount_small,
   rows=622, weakest_split_rows=142, full_toxic=0.729904,
   worst_split_toxic=0.802817.
3. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_flat::impulse_none,
   bucket=basis_discount_small + basis_change_flat,
   rows=356, weakest_split_rows=62, full_toxic=0.662921,
   worst_split_toxic=0.699387.
4. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_up_pressure::impulse_none,
   bucket=basis_discount_small,
   rows=613, weakest_split_rows=124, full_toxic=0.654160,
   worst_split_toxic=0.759358.
5. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_flat::impulse_none,
   bucket=basis_discount_small + premium_discount_small,
   rows=538, weakest_split_rows=115, full_toxic=0.659851,
   worst_split_toxic=0.719212.
6. BTCUSDT 15m h24 tradable_rotation_candidate:
   local=geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale,
   bucket=basis_discount_small,
   rows=313, weakest_split_rows=71, full_useful=0.632588,
   weakest_split_useful=0.521127, full_margin_improvement=0.072540,
   weakest_split_margin_improvement=0.025844.

Task (only after explicit user approval):
- Write a docs-only strategy-premise spec, suggested path:
  docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md.
- The spec must decide whether the 6 passing BTCUSDT 15m context cohorts justify:
  a no-trade filter premise, a rotation-entry premise, two separate premise
  tracks, or no further strategy work.
- Keep no-trade/toxic filtering and rotation-entry premises separate unless the
  spec gives a clear reason to combine them.
- Preserve the audit's meaning: basis/premium context is a closed-candle
  conditioning source, not a tradable signal by itself. The context pass is
  separation evidence only.
- Define what a later implementation would be allowed to test, but do not
  implement it in this task. If the spec authorizes a later audit/backtest, it
  must require separate explicit user approval and must name exact stop states.
- Carry forward source and finality constraints: only local Binance USDT-M
  futures BTCUSDT data and the already validated derivatives context rows
  needed by the selected premise; confirmed closed-candle decisions; source
  lag rule source_close_time + 5m <= decision_candle_close_time; no forward
  fill/interpolation/nearest-future joins; no new source downloads.
- Carry forward exclusion constraints: do not reopen, retune, rename, gate-relax,
  or promote closed families. Do not import old binance-bot strategy/scoring/
  order-management/live/deploy/credential logic.

Forbidden in this task:
- Any implementation before explicit user approval.
- Go code, CLI flags, generated result directories, audit runs, source downloads,
  network requests, source materialization, data writes under
  ../binance-bot/data/derivatives/, entries, exits, P&L backtests, optimizer
  grids, replay, walk-forward, portfolio construction, paper/testnet/live paths,
  exchange API, credentials, deploy files, martingale, averaging down,
  two-exchange logic, broad symbol/interval/source mining, or closed-family
  rescue.

Suggested docs-only stop states:
- derivatives_context_strategy_premise_spec_ready_for_user_approval;
- derivatives_context_strategy_premise_spec_rejected_no_strategy_premise.

Closeout (only after an approved docs-only spec):
- Add the new spec doc.
- Update README.md docs index and intro if the current next gate changes.
- Update memory/PROGRESS.md with the docs-only spec result and stop state.
- Update memory/DECISIONS.md only if the spec creates a durable boundary.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next gate.
- Run:
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- Commit completed docs and memory changes after checks pass unless told not to.
```
