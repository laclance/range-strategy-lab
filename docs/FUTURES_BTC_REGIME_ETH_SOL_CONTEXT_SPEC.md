# Futures BTC Regime Plus ETH/SOL Range Context Spec

Date: 2026-06-27

## Verdict

Stop state:
`btc_regime_eth_sol_context_spec_parked_pending_user_scope_approval`.

This is a documentation-only parked scope-expansion spec. It records a possible
future non-trading audit where BTCUSDT acts as market-regime context while
ETHUSDT and SOLUSDT may act as authority symbols only after explicit user
approval.

It does not authorize implementation, source expansion beyond the already local
BTC/ETH/SOL futures files, entries, exits, P&L backtests, optimizer grids,
replay, walk-forward, strategy packaging, live/paper/testnet work, exchange API,
credentials, deployment, martingale, averaging down, broad symbol mining, or
two-exchange logic.

## Motivation

Prior range-universe work showed that BTCUSDT was often weak as an authority
trade symbol while ETHUSDT and SOLUSDT carried some aggregate evidence in the
structured-compression branch. That branch later failed walk-forward robustness
and remains exclusion evidence.

The materially different research question is:

> Is BTCUSDT more useful as a regime input than as the traded authority symbol
> for ETH/SOL range contexts?

This is not structured-compression rescue. The first step must be a non-trading
context audit.

## Preconditions

This spec may become implementation-ready only if the user explicitly approves a
controlled BTC/ETH/SOL context audit or if a later range-state review recommends
this exact scope.

The initial data scope may use only the already validated local Binance USDT-M
futures `5m` sources:

| Symbol | Role In Parked Spec |
| --- | --- |
| `BTCUSDT` | market-regime context and diagnostic-only authority row |
| `ETHUSDT` | possible authority symbol after audit approval |
| `SOLUSDT` | possible authority symbol after audit approval |

No new symbols, broad symbol mining, data downloads, or exchange API use is
authorized.

## Proposed Non-Trading Audit

The future audit should combine two layers:

### BTC Market-Regime Layer

Closed-candle BTCUSDT context features:

- BTC range state from the range-state construction loop;
- BTC realized-volatility bucket;
- BTC trend-pressure bucket;
- BTC impulse bucket;
- BTC participation proxy;
- BTC close location relative to its own active range;
- BTC higher-timeframe slope or direction proxy from the same `5m` parent
  source.

### ETH/SOL Local Range Layer

Closed-candle ETHUSDT and SOLUSDT local features:

- local range geometry;
- local volatility state;
- local trend state;
- local impulse state;
- local participation proxy;
- relative strength versus BTC;
- local forward labels as labels only.

## State ID Construction

A future audit may build state IDs such as:

```text
btc_regime_eth_sol_v1::<symbol>::<timeframe>::<btc_regime_bucket>::<local_range_bucket>::<relative_strength_bucket>
```

All buckets must be known at decision-candle close. Future labels may appear only
in label/cohort artifacts.

## Exclusion Boundary

The audit must stop if it becomes any of these:

- structured-compression retune or rescue;
- ETH/SOL authority replay from the failed walk-forward branch;
- BTC promotion despite BTC authority weakness;
- broad symbol mining;
- portfolio construction;
- live/paper/testnet/deploy path;
- cross-exchange or two-exchange execution.

BTC may be diagnostic-only or context-only. It may not be promoted to authority
from this parked spec without a later passing audit and explicit review.

## Possible Route Questions

The non-trading audit may ask:

1. Do ETH/SOL local range states behave better when BTC is flat, compressed, or
   volatility-normal?
2. Are ETH/SOL range-continuation labels more stable when BTC trend pressure is
   aligned?
3. Are ETH/SOL rotation labels more stable when BTC is range-bound and not
   impulsive?
4. Are some BTC states toxic for all ETH/SOL range entries?

These are context questions, not trade entries.

## Required Later Artifacts

If later approved, use:

```text
results/futures-btc-regime-eth-sol-context-audit/
```

Expected artifacts:

- `futures_btc_regime_eth_sol_context_sources.csv/json`;
- `futures_btc_regime_eth_sol_context_coverage.csv/json`;
- `futures_btc_regime_eth_sol_context_btc_states.csv/json`;
- `futures_btc_regime_eth_sol_context_local_states.csv/json`;
- `futures_btc_regime_eth_sol_context_relative_strength.csv/json`;
- `futures_btc_regime_eth_sol_context_labels.csv/json`;
- `futures_btc_regime_eth_sol_context_cohorts.csv/json`;
- `futures_btc_regime_eth_sol_context_rankings.csv/json`;
- `futures_btc_regime_eth_sol_context_summary.csv/json`.

Common outputs must remain zero-trade compatible unless a much later strategy
spec explicitly changes that.

## Gates For Later Approval

A context cohort may pass only if:

- each authority symbol has adequate full and split counts;
- ETHUSDT and SOLUSDT do not rely on the same single fragile period;
- BTC context is known before decision time;
- BTC is not used as hidden future information;
- full and weakest-split useful/toxic separation beats the later declared gates;
- the result is not a structured-compression or breakout-retest reslice.

## Stop States

- `btc_regime_eth_sol_context_spec_parked_pending_user_scope_approval`;
- `btc_regime_eth_sol_context_spec_ready_for_audit_implementation`;
- `btc_regime_eth_sol_context_source_gap`;
- `btc_regime_eth_sol_context_failed_no_usable_context`;
- `btc_regime_eth_sol_context_passed_needs_strategy_premise_spec`;
- `btc_regime_eth_sol_context_rejected_closed_family_reslice`.

## Review Rule

Do not implement from current state. This spec exists so the future ETH/SOL/BTC
context idea has a controlled scope if the user chooses it.