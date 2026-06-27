# Futures Spread-Range Strategy Spec

Date: 2026-06-27

## Verdict

Stop state:
`spread_range_strategy_spec_parked_pending_engine_scope_approval`.

This is a documentation-only parked spec. It records a possible future branch
where ranges are defined on relative-price or beta-adjusted spread series rather
than on a single instrument price box.

It does not authorize implementation, engine changes, multi-leg accounting,
entries, exits, P&L backtests, optimizer grids, replay, walk-forward, source
expansion, symbol expansion beyond controlled BTC/ETH/SOL inputs,
paper/testnet/live work, exchange API, credentials, deployment, martingale,
averaging down, broad pair mining, or two-exchange execution.

## Motivation

Single-symbol range events have repeatedly failed or proved fragile. A materially
different range-first idea is to define the range on a spread:

```text
relative value moves into a range, dislocates, and reverts or expands
```

The first step must be a non-trading spread-state audit. Trading a spread
requires engine capabilities that are not part of the current one-position
single-instrument backtester.

## Candidate Spread Series

A future non-trading audit may define only controlled BTC/ETH/SOL spread series
unless the user approves a broader universe.

Allowed initial candidates:

- `ETHUSDT / BTCUSDT` relative price;
- `SOLUSDT / BTCUSDT` relative price;
- `SOLUSDT / ETHUSDT` relative price;
- log-price differences for the same pairs;
- beta-adjusted ETH residual versus BTC;
- beta-adjusted SOL residual versus BTC.

All spreads must be constructed from synchronized closed candles. No missing or
look-ahead alignment is allowed.

## Required Engine Gap Before Trading

The current lab engine is single-position and single-instrument. Before any
trade-producing spread strategy can exist, a separate engine spec must define:

- synchronized multi-series candle iteration;
- multi-leg signal representation;
- leg-level entry prices;
- leg-level fees and slippage;
- notional allocation per leg;
- hedge ratio or beta calculation freeze rules;
- margin/notional caps;
- leg-level stop/target or spread-level invalidation;
- force-close behavior when one leg is unavailable;
- trade and summary artifact schema for multi-leg positions.

This parked spec does not authorize those changes.

## Proposed Non-Trading Audit

A future implementation-ready audit should answer:

1. Which controlled spreads have stable range episodes?
2. Are spread ranges more durable than individual symbol ranges?
3. Do spread boundary states separate reversion, expansion, and chop labels?
4. Are spread states stable across `2021_2022_stress`, `2023_2024_oos`, and
   `2025_2026_recent`?
5. Does the result remain useful without broad pair mining?

## Feature Families

Closed-candle spread features may include:

- spread value;
- spread z-score from a closed rolling window;
- spread range width;
- spread range age;
- spread volatility percentile;
- spread trend/slope;
- residual volatility;
- hedge-ratio stability;
- relative-strength impulse;
- correlation stability;
- cointegration-style diagnostics as diagnostics only, if implemented without
  future leakage.

## Forward Labels

Labels must be spread outcomes, not single-leg outcomes:

- spread contained rotation;
- spread boundary chop;
- spread expansion upward;
- spread expansion downward;
- spread false-break reentry;
- spread drift-through;
- spread no-resolution.

Future labels cannot be features.

## Exclusion Boundary

Stop if the branch becomes:

- broad pair mining;
- hidden portfolio construction;
- two-exchange execution;
- exchange API or live-adjacent logic;
- a disguised structured-compression or breakout-retest retry on a synthetic
  series;
- an entry strategy before non-trading spread-state evidence exists.

## Required Later Artifacts

If later approved, use:

```text
results/futures-spread-range-state-audit/
```

Expected artifacts:

- `futures_spread_range_sources.csv/json`;
- `futures_spread_range_coverage.csv/json`;
- `futures_spread_range_series.csv/json`;
- `futures_spread_range_states.csv/json`;
- `futures_spread_range_labels.csv/json`;
- `futures_spread_range_cohorts.csv/json`;
- `futures_spread_range_rankings.csv/json`;
- `futures_spread_range_summary.csv/json`;
- `futures_spread_range_skips.csv/json`.

Common outputs must remain zero-trade compatible for the non-trading audit.

## Stop States

- `spread_range_strategy_spec_parked_pending_engine_scope_approval`;
- `spread_range_state_audit_spec_ready_for_implementation`;
- `spread_range_state_audit_source_gap`;
- `spread_range_state_audit_failed_no_usable_spread_state`;
- `spread_range_state_audit_passed_needs_engine_spec`;
- `spread_range_state_audit_passed_needs_strategy_premise_spec`;
- `spread_range_rejected_broad_pair_mining`;
- `spread_range_rejected_closed_family_reslice`.

## Review Rule

Do not implement from current state. This is a materially different future branch
that requires explicit engine/source approval before any code.