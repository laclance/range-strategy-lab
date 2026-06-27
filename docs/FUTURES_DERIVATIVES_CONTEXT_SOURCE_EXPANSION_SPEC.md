# Futures Derivatives Context Source Expansion Spec

Date: 2026-06-27

## Verdict

Stop state:
`derivatives_context_source_expansion_spec_parked_pending_user_source_approval`.

This is a documentation-only parked source-expansion spec. It records a possible
future branch that evaluates whether derivatives market context improves range
state usefulness.

It does not authorize data downloads, exchange API keys, live trading,
paper/testnet trading, order placement, deployment, entries, exits, P&L
backtests, optimizer grids, replay, walk-forward, strategy packaging,
martingale, averaging down, broad symbol mining, or two-exchange logic.

## Intent

Some failed range trades may be adverse-selection or crowding problems rather
than pure geometry problems. A later source-expansion audit could test whether
funding, basis, open interest, taker-flow, long/short, or order-book context
separates tradable range states from toxic ones.

This branch is market-data context only. It must start as a non-trading source
and alignment audit.

## Allowed Future Data Families

A later user-approved brief may select one or more market-data families:

- funding rate history;
- mark price and index price basis;
- open interest and open-interest statistics;
- global or account-ratio long/short statistics where available;
- taker buy/sell volume;
- top-of-book or order-book depth snapshots;
- aggregate trades only if timestamp alignment and storage rules are defined.

No order placement, account access, private endpoints, API keys, or live
execution path is allowed by this spec.

## Source Manifest Requirements

Every external context source must have a manifest that records:

- source name;
- endpoint or file path;
- symbol;
- market/product;
- timestamp semantics;
- first and last timestamp;
- row count;
- duplicate count;
- gap or missing-interval count;
- timezone assumptions;
- expected update cadence;
- alignment target timeframe;
- missing-data policy;
- comparison-only status;
- validation status.

The implementation must reject unknown timestamp semantics, non-monotonic data
that cannot be sorted without duplicates, missing required fields, and rows that
cannot be aligned without look-ahead.

## Alignment Rules

External context must be joined to candle decisions using closed-candle
semantics:

1. A decision candle may use only context observations whose timestamps are known
   at or before that decision candle close.
2. If a metric is published after the period it describes, the implementation
   must model that publication lag or mark the source rejected.
3. Forward-filled context must declare a max staleness.
4. Missing context must produce skip or missingness flags, not silent defaults.
5. Alignment artifacts must make the join auditable.

## Candidate Context Features

A later non-trading audit may include:

### Funding / Basis

- funding rate level;
- funding percentile;
- funding sign and persistence;
- mark-index basis;
- basis percentile;
- basis expansion/contraction.

### Open Interest

- open-interest change over closed windows;
- open-interest percentile;
- price up plus OI up / down combinations;
- OI impulse flags.

### Taker Flow / Long-Short

- taker buy/sell imbalance;
- taker imbalance percentile;
- long/short ratio level;
- long/short ratio change;
- crowding flags.

### Order Book / Liquidity

- spread;
- top-depth imbalance;
- depth slope;
- depth withdrawal or replenishment proxy;
- book dislocation flags.

These features may be used only as state context until a separate strategy spec
is approved.

## Non-Trading Audit Question

The first implementation, if later approved, should ask:

> Do derivatives context features improve the separation between usable range
> states, toxic range states, continuation states, and no-trade states after
> source alignment and split-stability gates?

It should not ask whether a funding or order-book signal can be traded directly.

## Exclusion Boundary

Stop if the branch becomes:

- private-account or order endpoint usage;
- exchange-key handling;
- live, paper, or testnet execution;
- broad symbol mining;
- a strategy before source validation and non-trading context gates pass;
- a rescue filter for failed occupancy rotation, structured compression,
  breakout-retest, clean breakout, hold-inside/midline, impulse absorption, or
  nested range rotation;
- a two-exchange arbitrage path.

## Required Later Artifacts

If later approved, use:

```text
results/futures-derivatives-context-source-audit/
```

Expected artifacts:

- `futures_derivatives_context_sources.csv/json`;
- `futures_derivatives_context_candle_coverage.csv/json`;
- `futures_derivatives_context_external_coverage.csv/json`;
- `futures_derivatives_context_alignment.csv/json`;
- `futures_derivatives_context_features.csv/json`;
- `futures_derivatives_context_states.csv/json`;
- `futures_derivatives_context_labels.csv/json`;
- `futures_derivatives_context_cohorts.csv/json`;
- `futures_derivatives_context_rankings.csv/json`;
- `futures_derivatives_context_summary.csv/json`;
- `futures_derivatives_context_skips.csv/json`.

Common outputs must remain zero-trade compatible for the first audit.

## Gates

A context source may be admitted only when:

- timestamp semantics are explicit;
- no look-ahead join is required;
- aligned rows cover enough of every period split;
- missingness is reported and bounded;
- the source is reproducible from local files or a later approved fetch process;
- the source does not require private keys or account access.

A context cohort may pass only when:

- source coverage gates pass;
- full and split counts are adequate;
- useful/toxic separation improves versus candle-only state rows;
- improvement is not concentrated in one split;
- the result is not a closed-family rescue.

## Stop States

- `derivatives_context_source_expansion_spec_parked_pending_user_source_approval`;
- `derivatives_context_source_expansion_spec_ready_for_source_audit`;
- `derivatives_context_source_expansion_source_gap`;
- `derivatives_context_source_expansion_alignment_failed`;
- `derivatives_context_source_expansion_failed_no_context_gain`;
- `derivatives_context_source_expansion_passed_needs_strategy_premise_spec`;
- `derivatives_context_source_expansion_rejected_live_or_private_api_path`;
- `derivatives_context_source_expansion_rejected_closed_family_rescue`.

## Review Rule

Do not implement from current state. This spec exists only to prevent future
source expansion from being improvised or confused with live/exchange work.