# Futures Range Context Router Spec

Date: 2026-06-27

## Verdict

Stop state:
`range_context_router_spec_ready_for_audit_implementation`.

This spec is now implementation-ready because
`docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md` stopped at
`range_state_construction_loop_audit_passed_needs_router_spec`.

It authorizes only a non-trading router audit. No strategy, entry, exit,
optimizer, replay, walk-forward run, source expansion, symbol expansion, data
download, paper/testnet/live path, exchange API, credential, deploy file,
martingale, averaging down, or two-exchange logic is approved by this document.

## Dependency

This spec became implementation-ready because
`docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md` stopped with:

- `range_state_construction_loop_audit_passed_needs_router_spec`;

Do not implement this router by retuning feature buckets or relaxing gates.

## Intent

The router converts pre-entry range states into a small set of closed-candle
research routes:

```text
no_trade | tradable_rotation | trend_continuation | diagnostic_only
```

The router is not a trade signal. It is a decision layer that says which later
strategy premise, if any, is allowed to be specified.

## Source Contract

Initial router scope is inherited from the range-state audit:

- Binance USDT-M futures;
- BTCUSDT;
- accepted `5m` parent source;
- closed UTC `15m`, `1h`, and `4h` resamples;
- closed-candle features only;
- no ETH/SOL authority;
- no external derivatives market data.

Any ETH/SOL or derivatives-source router requires a separate user-approved scope
brief.

## Router Labels

### `no_trade`

Use when the state audit identifies stable toxic states, especially states with
high boundary chop, low-width noise, unstable label mixes, or adverse split
behavior.

A `no_trade` route may be the only useful outcome. It can later become a filter
for other strategy research, but it does not authorize entries.

### `tradable_rotation`

Use only when the state audit identifies a stable rotation-style state without
becoming the failed `range_occupancy_rotation_v1` grammar.

This route may later authorize a new mean-reversion premise spec, but only with
fresh entry/exit rules and review gates.

### `trend_continuation`

Use only when the state audit identifies stable continuation-style states. This
must not reopen structured compression, clean breakout, or breakout-retest under
new names.

A later continuation strategy premise must be materially new and must explain
why it is not a retune of the closed families.

### `diagnostic_only`

Use when the state is interesting but fails count, split-stability, useful-rate,
toxic-rate, or closed-family protection gates.

Diagnostic states must not be converted into strategies without a later passing
audit or explicit user scope change.

## Router Construction Rules

The later router implementation, if approved, must be deterministic:

1. Load the state audit outputs or recompute state rows from source.
2. Apply only predeclared route rules from the passing review.
3. Assign exactly one router label per eligible state row.
4. Record all ambiguous, missing, or conflicting rows.
5. Produce zero trades.

The router must not train a classifier, optimize thresholds, infer hidden labels,
or use future outcomes as inputs.

## Required Artifacts For Later Implementation

Default result directory, if later approved:

```text
results/futures-range-context-router-audit/
```

Common outputs must remain zero-trade compatible:

- `source_manifest.json`;
- `summary.csv/json`;
- `trades.json` with no trades.

Router-specific outputs:

- `futures_range_context_router_sources.csv/json`;
- `futures_range_context_router_coverage.csv/json`;
- `futures_range_context_router_rules.csv/json`;
- `futures_range_context_router_rows.csv/json`;
- `futures_range_context_router_cohorts.csv/json`;
- `futures_range_context_router_rankings.csv/json`;
- `futures_range_context_router_summary.csv/json`;
- `futures_range_context_router_skips.csv/json`.

## Promotion Rules

A passing router may authorize only one of the following later docs:

- a no-trade filter integration spec;
- a materially new rotation premise spec;
- a materially new continuation premise spec;
- a review-only decision that no entry should be built.

It does not authorize a backtest, optimizer, replay, walk-forward, packaging,
live-adjacent work, symbol expansion, source expansion, or gate relaxation.

## Stop States

- `range_context_router_spec_parked_pending_range_state_audit`;
- `range_context_router_spec_ready_for_audit_implementation`;
- `range_context_router_source_gap`;
- `range_context_router_failed_no_actionable_route`;
- `range_context_router_passed_no_trade_filter_only`;
- `range_context_router_passed_needs_rotation_premise_spec`;
- `range_context_router_passed_needs_continuation_premise_spec`;
- `range_context_router_rejected_closed_family_reslice`.

## Review Rule

The range-state construction loop review now exists and explicitly authorizes
the router as the next step. Keep the implementation non-trading and
zero-trade compatible.
