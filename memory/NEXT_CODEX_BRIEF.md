# Next Codex Brief: Derivatives Context No-Trade Filter Integration Spec Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do NOT implement anything in this task without explicit user approval. The
approved zero-trade derivatives no-trade filter premise audit passed, but the
next step is a docs-only filter integration spec, not code. If approval is not
explicitly given in the session, make no docs edits, no Go code changes, no CLI
flag, no generated result directory, no audit run, no source download, no
network request, no data write, and no strategy/P&L work; report the waiting
state and stop.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md
  for the passed no-trade filter premise audit.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md for the
  selected premise and forbidden alternatives.
- Read docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md only if exact upstream
  context-audit evidence is needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline audit/backtest flag is passed.
- BTCUSDT-only price-range audits remain stopped by:
  range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit.
- BTC regime plus ETH/SOL context audit failed with 0 cohorts and remains
  closed.
- Derivatives source materialization, source audit, and context audit all passed
  in zero-trade form.
- The docs-only derivatives context strategy-premise spec selected exactly one
  later track: a BTCUSDT 15m derivatives-context no-trade filter premise audit.
  It rejected rotation-entry and two-track alternatives. The single rotation
  candidate remains diagnostic only.
- The zero-trade derivatives no-trade filter premise audit passed at:
  derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec.
- Result path:
  results/futures-derivatives-no-trade-filter-premise-audit/.
- Result counts:
  source_rows=4, coverage_rows=7, filter_definition_rows=5,
  exact_candidate_rows=20, canonical_union_rows=4, overlap_rows=40,
  veto_candidate_rows=1823, collateral_rows=37, missingness_rows=4,
  exact_candidates_passed=5, canonical_union_passed=true, trades=0.
- The canonical filter ID is:
  btc_15m_basis_discount_no_trade_veto_v1.
- Canonical union full-sample facts: 1,823 deduplicated veto rows; 821 overlap
  rows; 515 nested trend-down premium overlap rows; 1,241 no-trade toxic rows;
  toxic rate 0.680746; min split toxic rate 0.665485; weakest split rows 387;
  full toxic improvement versus local-only baseline 0.046269; 311
  rotation-useful and 271 continuation-useful full-sample collateral rows were
  reported.
- Exact candidate facts reproduced:
  1. trend_down_pressure + basis_discount_small + premium_discount_small:
     rows=515, weakest_split_rows=110, full_toxic=0.732039,
     worst_split_toxic=0.800000, full_toxic_improvement=0.049738.
  2. trend_down_pressure + basis_discount_small:
     rows=622, weakest_split_rows=142, full_toxic=0.729904,
     worst_split_toxic=0.802817, full_toxic_improvement=0.047603.
  3. trend_flat + basis_discount_small + basis_change_flat:
     rows=356, weakest_split_rows=62, full_toxic=0.662921,
     worst_split_toxic=0.699387, full_toxic_improvement=0.045126.
  4. trend_up_pressure + basis_discount_small:
     rows=613, weakest_split_rows=124, full_toxic=0.654160,
     worst_split_toxic=0.759358, full_toxic_improvement=0.050840.
  5. trend_flat + basis_discount_small + premium_discount_small:
     rows=538, weakest_split_rows=115, full_toxic=0.659851,
     worst_split_toxic=0.719212, full_toxic_improvement=0.042056.

Task (only after explicit user approval):
- Write a docs-only filter integration spec, suggested path:
  docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_INTEGRATION_SPEC.md
- The spec may answer only:
  Given the passed zero-trade no-trade filter premise audit, what exactly would a
  later integration be allowed to test, and under what gates, without converting
  the veto into an entry signal or P&L claim?
- The spec must decide whether the passed veto can be held as a future
  integration candidate, should be deferred until an independent entry premise
  exists, or should be rejected as not integration-ready despite passing the
  premise audit.
- Preserve the filter as a veto only. It may say "skip this context" only inside
  a later explicitly approved integration audit. It may not say "trade the
  opposite", "trade basis", "fade basis", "enter rotation", "enter
  continuation", or "improve P&L".
- Define exact allowed later integration inputs and outputs if, and only if, the
  docs-only spec selects a later implementation gate. The later implementation
  still requires separate explicit user approval.
- Require any later implementation to stay offline and zero-trade until a
  separate independent entry premise exists and is explicitly approved for
  interaction testing.
- Record that the diagnostic rotation row remains diagnostic only.
- Record that no ETHUSDT/SOLUSDT, other timeframe, source expansion, or closed
  family rescue is authorized.

Forbidden in this docs-only task:
- Go code, CLI flags, generated result directories, audit runs, source downloads,
  network requests, source materialization, data writes under
  ../binance-bot/data/derivatives/, entries, exits, P&L backtests, optimizer
  grids, replay, walk-forward, portfolio construction, paper/testnet/live paths,
  exchange API, credentials, deploy files, martingale, averaging down,
  two-exchange logic, or strategy promotion.
- Treating the veto as an entry signal, a basis-tradability claim, a direct
  strategy, or an excuse to reopen closed BTCUSDT price-only, BTC-regime,
  router-rotation, spread-range, or volatility-aware-exit families.

Suggested docs-only stop states:
- derivatives_context_no_trade_filter_integration_spec_ready_for_user_approval;
- derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise;
- derivatives_context_no_trade_filter_integration_spec_rejected_no_integration_premise.

Closeout:
- Add/update only the docs-only integration spec and memory/index files needed
  to point at it.
- Update README.md docs index and intro if the current next gate changes.
- Update memory/PROGRESS.md with the spec stop state, commands, and short
  factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh memory/NEXT_CODEX_BRIEF.md to the next gate.
- Run:
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- Commit completed docs/memory changes after checks pass unless told not to.
```
