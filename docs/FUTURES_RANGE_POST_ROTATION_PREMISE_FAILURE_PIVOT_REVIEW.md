# Futures Range Post-Rotation Premise Failure Pivot Review

Date: 2026-06-28

## Verdict

Stop state:
`range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.

This is a documentation-only post-failure pivot review. It preserves the failed
verdict in `docs/FUTURES_RANGE_ROUTER_ROTATION_PREMISE_AUDIT_REVIEW.md` and
closes `router_gated_boundary_reclaim_rotation_v1` in reviewed form.

No new BTCUSDT-only, candle-price-only, range-premise audit is worth specifying
from the current evidence. The range-state construction loop and context router
were useful non-trading infrastructure, but the only passed positive router
branch collapsed to a failed boundary-reclaim premise. Continuing by retuning,
renaming, relaxing gates, or slicing the same router/rotation surface would be a
closed-family reslice.

The next bounded action is not implementation. If research continues, it must
start with an explicit user scope choice for a materially different parked
direction, such as controlled BTC/ETH/SOL context, derivatives market-data
context, or spread-range source/engine work. Without that approval, stop with no
next audit.

This review does not authorize Go code, strategy code, entries, exits, P&L
strategy backtests, optimizer grids, replay, walk-forward logic, strategy
packages, paper/testnet/live paths, exchange API use, credentials, deploy files,
source expansion, symbol expansion, broad mining, martingale, averaging down, or
two-exchange logic.

## Source And Dependency Facts

Accepted source:

| Field | Value |
| --- | --- |
| Source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Reviewed dependency chain:

| Milestone | Key facts | Stop state |
| --- | --- | --- |
| Range-state construction loop audit | `29,784` states, `89,352` labels, `68,796` cohorts, `16,335` rankings, `58` passing cohorts | `range_state_construction_loop_audit_passed_needs_router_spec` |
| Range context router audit | `58` rules, `29,784` router rows, `3` passing cohorts, `1,299` `tradable_rotation`, `0` `trend_continuation`, `0` conflicts | `range_context_router_passed_needs_rotation_premise_spec` |
| Router rotation premise spec | Named `router_gated_boundary_reclaim_rotation_v1`; allowed only a zero-trade non-trading audit | `range_router_rotation_premise_spec_ready_for_non_trading_audit` |
| Router rotation premise audit | `278` segments, `97` events, `97` outcomes, `3` rankings, `0` passing cohorts | `range_router_rotation_premise_audit_failed_no_premise` |

The failed premise audit produced `43` lower events, `54` upper events, `71`
midline outcomes, `22` hard-adverse outcomes, and `4` chop or no-resolution
outcomes. Its top failure reasons were:

```text
inadequate_event_count,inadequate_split_event_count,single_split_contribution_above_gate,behavior_gate_failed
```

Do not convert the `278` context segments, `97` boundary-reclaim events, or
`1,299` `tradable_rotation` router rows into trades.

## Pivot Decision

The BTCUSDT-only price-range ladder has now passed through the intended
non-trading sequence:

1. construct closed-candle range states;
2. route states into no-trade, rotation, continuation, or diagnostic context;
3. test the only passed positive router family through a materially new event
   premise.

That sequence did not produce a usable premise. The router found no passing
`trend_continuation` branch, and the only rotation branch failed its premise
audit on event count, split stability, event concentration, and adverse behavior.
The remaining passed router evidence is primarily no-trade filtering, which can
be useful protection for a future strategy but is not itself a strategy premise.

Therefore:

- no trigger audit should be built from `router_gated_boundary_reclaim_rotation_v1`;
- no BTCUSDT-only price-only range audit should be invented automatically from
  the same state/router artifacts;
- no volatility-aware exit work is valid because no materially new entry has
  shown gross edge before costs;
- no no-trade filter implementation is useful without a later independent
  entry premise to protect;
- no source, symbol, or engine expansion may start without explicit user
  approval.

## Closed-Family Separation

| Closed family | Why this pivot does not reopen it |
| --- | --- |
| `range_occupancy_rotation_v1` | The reviewed occupancy grammar and optimizer grid already failed; this pivot does not retune its occupancy buckets, targets, stops, or replay route. |
| Hold-inside / midline | The failed premise already used midline rotation as an outcome label and still failed; this pivot does not promote hold-inside or midpoint transition events. |
| Breakout-retest / acceptance | Broken range boundaries and outside acceptance remain closed; this pivot does not trade acceptance after a break. |
| Clean breakout continuation | The router found `0` passing continuation cohorts, and clean expansion against rotation was adverse in the failed premise. |
| Structured compression | The fragile ETH/SOL authority stream and its walk-forward failure remain exclusion evidence; this pivot does not rescue it through router language. |
| Impulse absorption | Impulse buckets remain context only; abnormal impulse candles are not promoted as entry events. |
| Higher-timeframe nested range rotation | The nested-range audit had insufficient valid events; this pivot does not repackage it through the `15m` router branch. |
| Range quality/session/failure-mode triage | Those cohorts failed by themselves and remain diagnostic context only. |
| `router_gated_boundary_reclaim_rotation_v1` | This exact premise is closed in reviewed form after `97` valid events and `0` passing rankings. |
| Legacy spot-only SR timing/compression evidence | Spot-only outputs remain historical context only and are not futures promotion authority. |

## Candidate Direction Disposition

| Candidate direction | Decision | Reason |
| --- | --- | --- |
| More BTCUSDT-only price-range slicing | Stop | The state/router/premise ladder already tested the available positive route and failed; more slicing risks gate relaxation. |
| No-trade filter audit | Do not start now | A no-trade filter can only protect a separate strategy. The repo has no approved independent entry premise to protect. |
| Volatility-aware exit model | Parked | Exit research requires a materially new entry template with gross edge before costs. |
| BTC regime plus ETH/SOL context | Materially different, but parked | It changes symbol scope and can start only with explicit user approval for a zero-trade context audit. |
| Derivatives market-data context | Materially different, but parked | It changes source scope and can start only with explicit user source approval and source-alignment gates. |
| Spread-range / pair-range work | Materially different, but parked | It changes source and engine scope, and any trading version would need separate multi-leg accounting. |

No automatic next audit is selected.

## Stop And Next-Step Rules

Current stop state:
`range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.

Allowed next interactive states:

- `range_post_rotation_premise_failure_pivot_user_chooses_no_further_audit`;
- `range_post_rotation_premise_failure_pivot_user_approves_btc_eth_sol_context_scope`;
- `range_post_rotation_premise_failure_pivot_user_approves_derivatives_context_source_scope`;
- `range_post_rotation_premise_failure_pivot_user_approves_spread_range_scope`;
- `range_post_rotation_premise_failure_pivot_rejected_strategy_backtest_request`;
- `range_post_rotation_premise_failure_pivot_rejected_closed_family_reslice`.

If the user chooses a parked direction, the next task should still be
documentation or a zero-trade source/context audit brief first. It must not jump
to entries, exits, P&L backtests, optimizers, replay, walk-forward, paper,
testnet, live paths, exchange APIs, credentials, deployment, source downloads,
broad symbol mining, martingale, averaging down, or two-exchange logic.

## Verification

Documentation-only closeout verification required:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
