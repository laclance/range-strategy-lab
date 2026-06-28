# Futures BTC Regime Plus ETH/SOL Context Scope Review

Date: 2026-06-28

## Verdict

Stop state:
`btc_regime_eth_sol_context_scope_review_approved_needs_zero_trade_audit_brief`.

This documentation-only scope review approves
`docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` as materially different enough
and tightly scoped enough to become the next separate zero-trade audit brief.

This approval does not authorize Go code, CLI flags, generated result
directories, source downloads, strategy code, entries, exits, P&L strategy
backtests, optimizer grids, replay, walk-forward logic, strategy packaging,
paper/testnet/live paths, exchange API use, credentials, deploy files, broad
mining, martingale, averaging down, or two-exchange logic.

The next permitted step is only a docs-only brief or implementation plan for a
future zero-trade context audit. That later audit still requires explicit
approval before code.

## Approval Rationale

BTC regime plus ETH/SOL local range context is materially different from another
BTCUSDT-only price-range slice because it changes the research question and
symbol roles:

- BTCUSDT is market-regime context, not the authority symbol being promoted.
- ETHUSDT and SOLUSDT are the only possible later authority rows.
- The future question is whether BTC regime buckets improve separation of
  ETH/SOL usable, toxic, rotation, continuation, or no-trade range states.
- BTCUSDT price-only range premises remain stopped by the post-rotation premise
  failure pivot review.

This is a context-scope approval, not a strategy approval. A future audit may
only ask whether BTC context known at the decision candle improves ETH/SOL range
state separation. It may not turn that separation directly into entries,
exits, P&L, replay, walk-forward, or strategy packaging.

## Allowed Source Scope

The only allowed source scope for a future brief is the already local Binance
USDT-M futures `5m` files:

| Symbol | Path | Loaded candles | Coverage | Status |
| --- | --- | ---: | --- | --- |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted after sort |

Prior source validation recorded `gap_count=0` and `duplicate_count=0` for all
three sorted streams. Zero-volume counts were `66` for BTCUSDT, `47` for
ETHUSDT, and `47` for SOLUSDT. The SOLUSDT file had one physical
non-monotonic row and was accepted only after sorting by open time.

No new symbols, source downloads, spot comparisons, derivatives context files,
private endpoints, exchange APIs, or broad symbol mining are approved.

## Role Boundaries

BTCUSDT may provide only:

- market-regime buckets;
- diagnostic-only authority rows;
- context features known at closed decision-candle time.

ETHUSDT and SOLUSDT may provide only:

- local range-state rows in a later zero-trade context audit;
- possible authority rows for context separation only;
- forward labels as labels, never as premise inputs.

BTCUSDT may not be promoted to authority by this review. ETHUSDT and SOLUSDT may
not be promoted to strategy authority by this review. Any future passing context
audit would still need a separate strategy premise spec before trades.

## Exclusion Boundary

This review is not:

- a structured-compression rescue;
- an ETH/SOL authority replay from the fragile walk-forward branch;
- BTCUSDT promotion after BTC authority weakness;
- a direct repackaging of `router_gated_boundary_reclaim_rotation_v1`;
- a retune, rename, or gate relaxation of any closed family;
- broad symbol mining;
- portfolio construction;
- spread-range or pair-range engine work;
- volatility-aware exit work;
- a paper/testnet/live path.

Structured compression remains exclusion evidence. The prior ETH/SOL authority
stream cannot be reused as promotion evidence, and its rules may not be copied
into the BTC-regime context audit under new names.

## Minimum Next Audit Question

The next zero-trade audit brief must answer only:

```text
Do BTCUSDT regime buckets, known at closed decision-candle time, improve
separation of ETHUSDT and SOLUSDT local range states into usable, toxic,
rotation, continuation, or no-trade contexts?
```

The future brief must reject any design that uses hidden future-label inputs,
converts context labels into trades, or starts from P&L strategy objectives.

## Rejection Criteria

Reject the later brief or implementation plan if it becomes any of these:

- closed-family reslice;
- broad mining or new symbol search;
- source gap or unapproved source expansion;
- hidden future-label input;
- BTCUSDT authority promotion;
- structured-compression rescue;
- ETH/SOL replay;
- entries, exits, P&L backtests, optimizer grids, replay, or walk-forward.

## Next Step

Refresh `memory/NEXT_CODEX_BRIEF.md` to a separate brief-writing task for the
future zero-trade BTC regime plus ETH/SOL context audit. The next task should
stop before code at:

```text
btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval
```

## Verification

Documentation-only closeout verification required:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
