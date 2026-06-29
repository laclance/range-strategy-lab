# Next Codex Brief: Derivatives Context No-Trade Filter Premise Audit Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do NOT implement anything in this task without explicit user approval. The
docs-only derivatives context strategy-premise spec has selected a BTCUSDT 15m
no-trade filter premise, but the next step is the user approving or rejecting a
zero-trade implementation audit. If approval is not explicitly given in the
session, make no Go code changes, no CLI flag, no generated result directory, no
audit run, no source download, no network request, no data write, and no
strategy/P&L work; report the waiting state and stop.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md for the
  selected no-trade filter premise and later-audit gates.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md for the passed context
  audit result and exact cohort evidence.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md only
  for the carried-forward anti-leakage and source-finality contract if needed.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md only for
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
- The derivatives zero-trade context audit passed at:
  derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec.
- The docs-only strategy-premise spec stopped at:
  derivatives_context_strategy_premise_spec_ready_for_user_approval.
- The spec selected exactly one later track: a BTCUSDT 15m derivatives-context
  no-trade filter premise audit. It rejected the rotation-entry premise and the
  two-track alternative for now. The single rotation candidate remains
  diagnostic only.
- The context audit result path remains:
  results/futures-derivatives-context-audit/.
- Context audit counts: source_rows=12, coverage_rows=18,
  basis_feature_rows=83,004, local_state_rows=83,640, label_rows=249,012,
  cohort_rows=512,190, ranking_rows=181,827, missingness_rows=36,
  passing_cohorts=6, trades=0.
- The five selected toxic/no-trade cohorts are all BTCUSDT 15m h48, all
  midline-balanced and volume-compressed, and all use lagged small-discount
  basis context, sometimes with premium-discount or basis-change-flat
  corroboration.

Selected no-trade filter candidates to preserve:
1. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none,
   bucket=basis_discount_small + premium_discount_small,
   rows=515, weakest_split_rows=110, full_toxic=0.732039,
   worst_split_toxic=0.800000.
2. BTCUSDT 15m h48 no_trade_toxic:
   local=geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none,
   bucket=basis_discount_small,
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

Diagnostic-only rotation row, not selected for implementation:
- BTCUSDT 15m h24 tradable_rotation_candidate:
  local=geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale,
  bucket=basis_discount_small,
  rows=313, weakest_split_rows=71, full_useful=0.632588,
  weakest_split_useful=0.521127, full_margin_improvement=0.072540,
  weakest_split_margin_improvement=0.025844.

Task (only after explicit user approval):
- Implement a zero-trade derivatives no-trade filter premise audit, suggested
  flag and output directory:
  - -futures-derivatives-no-trade-filter-premise-audit
  - results/futures-derivatives-no-trade-filter-premise-audit/
- The implementation may answer only:
  Do the five BTCUSDT 15m toxic derivatives-context cohorts define a stable
  closed-candle no-trade veto candidate after exact-row reproduction,
  de-duplication, split checks, and missingness accounting?
- Start from the exact five toxic rows above. Report both exact-row filters and
  a de-duplicated canonical union. Do not generalize beyond these rows unless
  the zero-trade audit proves the broader candidate preserves the declared gates.
- Preserve the nested-row rule: the trend_down_pressure premium-confirmed row is
  nested inside the broader trend_down_pressure + basis_discount_small row and
  must be reported as overlap, not double-counted.
- Keep the flat-trend rows corroborator-bound (basis_change_flat or
  premium_discount_small) unless the audit proves a broader flat-trend
  basis-discount filter passes.
- Use only local BTCUSDT Binance USDT-M futures candles and the already
  validated BTCUSDT derivatives context rows needed by the selected premise:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  ../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv
  ../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv
  ../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv
- Carry forward finality constraints: confirmed closed-candle decisions only;
  source_close_time + 5m <= decision_candle_close_time; no forward fill,
  interpolation, nearest-future joins, future source revisions, or silent
  default context; missing lagged context becomes skip/missingness rows.
- Forward labels may appear only as evaluation metadata in label/cohort/ranking/
  summary artifacts; they must never be feature, state, context, or gate inputs.
- Common summary.json/summary.csv/trades.json must stay zero-trade compatible.

Required later audit gates:
- BTCUSDT source/provenance/coverage facts reproduce, including one-5m lag and
  no-fill policy.
- All five exact toxic context rows reproduce as candidates before canonical
  union evaluation.
- Each exact candidate keeps full rows >= 300 and weakest split rows >= 60.
- Each exact candidate keeps full toxic rate >= 0.65 and worst split toxic rate
  >= 0.69.
- Each exact candidate keeps full toxic-rate improvement versus local-only
  baseline >= 0.04.
- The canonical de-duplicated filter union reports overlap counts and does not
  rely on double counting.
- The filter union remains toxic/no-trade dominated in full sample and every
  split.
- Useful/rotation labels blocked by the veto are reported as collateral damage,
  not hidden.
- The audit proves it is not a price-only or BTC-regime reslice and does not
  reopen, retune, rename, gate-relax, or promote any closed family.

Forbidden in this task:
- Any implementation before explicit user approval.
- Rotation-entry implementation or treating the diagnostic rotation row as a
  trade premise.
- ETHUSDT/SOLUSDT promotion, other timeframes, broad symbol/interval/source
  mining, source downloads, network requests, source materialization, data writes
  under ../binance-bot/data/derivatives/, entries, exits, P&L backtests,
  optimizer grids, replay, walk-forward, portfolio construction, paper/testnet/
  live paths, exchange API, credentials, deploy files, martingale, averaging
  down, two-exchange logic, or closed-family rescue.
- Importing old binance-bot strategy/scoring/order-management/live/deploy/
  credential/portfolio-coordinator logic.

Allowed later implementation stop states:
- derivatives_context_no_trade_filter_premise_audit_source_gap;
- derivatives_context_no_trade_filter_premise_audit_rejected_future_label_leak;
- derivatives_context_no_trade_filter_premise_audit_rejected_closed_family_rescue;
- derivatives_context_no_trade_filter_premise_audit_rejected_rotation_entry_rescue;
- derivatives_context_no_trade_filter_premise_audit_failed_no_usable_filter;
- derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec.

Closeout (only after an approved implementation):
- Add/update the no-trade filter premise audit code/tests and review doc.
- Keep generated CSV/JSON outputs under
  results/futures-derivatives-no-trade-filter-premise-audit/.
- Update README.md docs index and intro if the current next gate changes.
- Update memory/PROGRESS.md with stop state, result path, row counts, source
  facts, verification commands, and short factual outcome.
- Update memory/DECISIONS.md only if the audit creates a durable boundary.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next gate.
- Run:
  gofmt -w <touched go files>
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-no-trade-filter-premise-audit -out-dir results/futures-derivatives-no-trade-filter-premise-audit
  wc -l results/futures-derivatives-no-trade-filter-premise-audit/*.csv
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- Commit completed code/docs/memory changes after checks pass unless told not to.
```
