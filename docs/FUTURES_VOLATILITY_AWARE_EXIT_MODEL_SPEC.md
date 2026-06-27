# Futures Volatility-Aware Exit Model Spec

Date: 2026-06-27

## Verdict

Stop state:
`volatility_aware_exit_model_spec_parked_pending_gross_entry_edge`.

This is a documentation-only parked spec. It records a future exit-model research
branch for the case where a materially new range-derived entry template first
shows gross edge before costs.

It does not authorize implementation, entries, exits, P&L backtests,
optimization, replay, walk-forward, source expansion, symbol expansion,
paper/testnet/live work, exchange API use, credentials, deployment, martingale,
averaging down, or two-exchange logic.

## Dependency

This branch may start only after a future, materially new strategy premise
passes a fixed baseline gross-edge review. It must not be used to rescue:

- `range_occupancy_rotation_v1`;
- structured compression;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation.

The project workflow remains: gross edge first, then cost-sensitive and exit
variants.

## Intent

The current lab has mostly tested simple structural exits: boundary stops,
range-width targets, fixed max hold, and stop-first ambiguity. If a future entry
has gross edge, the next exit question is whether volatility state improves
survival without overfitting.

This branch treats exits as risk-shape research, not signal generation.

## Candidate Exit Families

A later implementation-ready brief may compare only predeclared exit families.

### Baseline Fixed Exit

- structural invalidation stop;
- fixed range-width target or fixed R target;
- fixed max hold;
- stop-first ambiguity.

This is the control row. A volatility-aware exit is not useful if it cannot beat
or at least stabilize the control without adding fragile complexity.

### ATR Stop / Target Variant

- stop distance as a multiple of ATR;
- target as fixed R from ATR stop distance or range-derived objective;
- ATR window and multiple must be declared before the run;
- no trailing update unless declared.

### Volatility-State Max Hold

- shorter max hold during volatility expansion;
- longer max hold during compressed or normal volatility;
- no dynamic change after entry unless the rule is closed-candle and predeclared.

### ATR Trailing Variant

- trail only after unrealized progress threshold;
- update on closed candles only;
- stop can only move toward reduced risk;
- no optimistic intrabar trail updates.

### Range-Mid / Boundary Staged Exit

- optional first objective at midpoint or prior internal line;
- final objective at opposite boundary or range-width target;
- requires explicit one-position accounting rules before implementation.

This is parked until the engine can represent partial exits or until the spec
chooses a single-exit approximation.

## Evaluation Rules

A later exit study must compare:

- gross P&L before costs;
- normal net P&L after existing costs;
- fee/slippage stress;
- train/OOS/recent/full splits;
- long and short sides;
- trade-count adequacy;
- drawdown versus net return;
- sensitivity to volatility regime.

It must rank by training data only and treat OOS/recent as gates, not ranking
inputs.

## Anti-Rescue Rule

If the entry fails gross P&L before costs, the exit branch stops immediately.
Do not try volatility-aware exits to make a negative signal look tradable.

If the entry is a closed-family reslice, stop with:

```text
volatility_aware_exit_model_rejected_closed_family_rescue
```

## Required Later Artifacts

If later approved, use a dedicated result directory:

```text
results/futures-volatility-aware-exit-model-review/
```

Expected artifacts:

- `futures_volatility_aware_exit_sources.csv/json`;
- `futures_volatility_aware_exit_coverage.csv/json`;
- `futures_volatility_aware_exit_entry_source.csv/json`;
- `futures_volatility_aware_exit_grid.csv/json`;
- `futures_volatility_aware_exit_trades.csv/json`;
- `futures_volatility_aware_exit_summary.csv/json`;
- `futures_volatility_aware_exit_rankings.csv/json`;
- `futures_volatility_aware_exit_selection.csv/json`.

Common `summary.csv/json` and `trades.json` may describe only the fixed selected
exit row after a later replay is approved.

## Stop States

- `volatility_aware_exit_model_spec_parked_pending_gross_entry_edge`;
- `volatility_aware_exit_model_spec_ready_for_implementation`;
- `volatility_aware_exit_model_source_gap`;
- `volatility_aware_exit_model_entry_failed_gross_edge`;
- `volatility_aware_exit_model_failed_no_selected_exit`;
- `volatility_aware_exit_model_passed_needs_fixed_replay_spec`;
- `volatility_aware_exit_model_rejected_closed_family_rescue`.

## Review Rule

Do not implement this branch from current project state. It is a future option
only after a materially new entry template earns gross-edge evidence.