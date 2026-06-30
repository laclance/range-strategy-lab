# Futures Independent Entry Premise And Hypothesis Map

Date: 2026-06-30

## Verdict

Stop state:
`independent_entry_premise_and_hypothesis_map_needs_user_scope_choice`.

This documentation-only milestone inventories the remaining materially different
research directions and tests whether the current repo state exposes exactly one
independent BTCUSDT `15m` entry-premise candidate for a later zero-trade audit.
It does not. The available BTCUSDT `15m` local-source candidates either collapse
back into reviewed closed families, remain no-trade/context evidence only, or
need a fresh user-supplied candidate event before they can become an entry
premise.

The canonical derivatives veto remains parked as future filter evidence only:

```text
btc_15m_basis_discount_no_trade_veto_v1
```

It is not an entry premise, is not used to select an entry premise, and is not
tested here. A later veto interaction audit still requires a separate,
independently approved candidate-entry stream.

This milestone authorizes no Go code, CLI flag, generated result directory,
audit run, source download, source materialization, data write, entry, exit,
P&L backtest, optimizer grid, replay, walk-forward, portfolio construction,
paper/testnet/live path, exchange API, credential, deploy file, martingale,
averaging down, two-exchange logic, closed-family rescue, or strategy
promotion.

## Fixed Selection Gates

An independent entry-premise candidate had to clear every gate below:

| Gate | Required outcome |
| --- | --- |
| Symbol | BTCUSDT only |
| Decision stream | closed-candle `15m` candidate rows |
| Sources | local accepted futures candles and deterministic resamples only |
| Derivatives veto | not used in entry inputs, candidate selection, side, ranking, or scoring |
| Family separation | materially different from reviewed closed families |
| Audit path | zero-trade-auditable before P&L or strategy work |
| Evidence shape | defined candidate event, side/timing, decision candle, and falsification rule |

These gates deliberately keep the first entry-premise lane compatible with a
future `15m` veto interaction while preventing the veto from inventing the
entry stream it would later filter.

## Inventory

| Direction | Current status | Reusable part | Decision here |
| --- | --- | --- | --- |
| BTCUSDT `15m` router rotation | Closed after `router_gated_boundary_reclaim_rotation_v1` failed with `97` valid events and `0` passing rankings | Router rows, range-state features, split/report shape | Not eligible; direct reslice risk |
| BTCUSDT `15m` no-trade/router toxic states | Useful filter/context evidence | No-trade labels and audit plumbing | Not an entry premise |
| Derivatives basis/premium veto | Passed as no-trade premise; `1,823` veto rows, `1,241` toxic rows | Future skip/retain filter only | Parked; cannot shape entry |
| Higher-timeframe nested range rotation | Reviewed and failed with only `3` valid events | Closed UTC resampling and parent/child range tooling | Not eligible for this `15m` lane |
| BTC regime plus ETH/SOL context | Reviewed zero-trade audit found `0` passing cohorts | Source-validation/context-audit pattern | Closed; violates BTCUSDT-only lane |
| Spread-range / pair-range | Parked behind source/engine complexity | Future source/engine spec shape | Needs separate scope choice |
| Funding/open interest/taker-flow/source expansion | Parked source families | Source-scope and source-audit process | Needs separate source approval |
| Volatility-aware exits | Parked until independent gross entry evidence exists | Future exit-review framing | Not eligible before entry premise |

## Candidate Evaluation

| Candidate | Gate result | Reason |
| --- | --- | --- |
| Reuse `range_context_router_v1|15m|h24|tradable_rotation` as entry stream | Fail | The only tested event premise from this route, `router_gated_boundary_reclaim_rotation_v1`, failed and is closed. Converting router rows directly into entries would bypass the failed premise gate. |
| Convert `15m` no-trade/toxic states into opposite-side entries | Fail | No-trade evidence may protect a future strategy, but "trade the opposite" is explicitly disallowed and would manufacture an entry from a filter. |
| Use `btc_15m_basis_discount_no_trade_veto_v1` as entry context | Fail | The veto is derivatives filter evidence only. It cannot define candidate rows, side, timing, or P&L expectation. |
| Reopen higher-timeframe nested rotation through `15m` projection | Fail | The reviewed `4h`/`1h` nested premise failed with sparse events. Projecting it to `15m` would be a closed-family rescue without a new premise. |
| Start spread-range or new derivatives source work | Defer | Potentially materially different, but outside the fixed BTCUSDT `15m` local-source lane and requires explicit source/engine scope approval. |
| User supplies a new BTCUSDT `15m` local-source event | Open | This could become eligible, but it is not present in current reviewed evidence. It must define the event, side/timing, and falsification rule before any audit. |

No candidate clears all fixed gates. The correct outcome is not rejection of all
future research; it is a user-choice stop. A future step must either supply a
new BTCUSDT `15m` local-source entry premise or explicitly change scope.

## Required User Choice

The next approved step must choose exactly one route:

| Route | What must be supplied | First allowed task |
| --- | --- | --- |
| New BTCUSDT `15m` local-source premise | Closed-candle candidate event, intended side/timing, why it is not a closed-family reslice, and falsification rule | Docs-only premise spec |
| Higher-timeframe premise | Interval set and materially different range behavior, not the failed nested rotation premise | Docs-only premise spec or zero-trade audit brief |
| Spread-range / pair-range | Source and engine scope acceptance | Docs-only source/engine scope spec |
| New derivatives/source family | Source family, provenance plan, and anti-leakage alignment premise | Docs-only source-scope or source-audit brief |
| No further audit | Explicit stop | Memory/brief closeout only |

Until one route is chosen, do not implement an audit.

## Conditional Future Entry-Premise Requirements

If the user supplies a new BTCUSDT `15m` local-source premise, the later spec
must define:

- source facts from `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv`;
- resampling/finality rules for closed UTC `15m` bars;
- candidate-event stream before any derivatives veto is applied;
- side and decision-candle timing, without forward labels in inputs;
- skip/missingness accounting;
- outcome labels and falsification rule for a zero-trade audit;
- closed-family separation versus router rotation, occupancy rotation,
  hold-inside/midline, breakout-retest/acceptance, clean breakout continuation,
  structured compression, impulse absorption, higher-timeframe nested rotation,
  BTC regime plus ETH/SOL context, and derivatives no-trade veto work.

That later audit may count and label candidate rows. It must not simulate fills,
model stops/targets, score P&L, optimize, replay, walk forward, promote a
strategy, or apply the derivatives veto unless a separate interaction audit is
explicitly approved afterward.

## Stop States

Used:

```text
independent_entry_premise_and_hypothesis_map_needs_user_scope_choice
```

Not used:

```text
independent_entry_premise_and_hypothesis_map_ready_for_user_approval
independent_entry_premise_and_hypothesis_map_rejected_closed_family_reslice
```

The ready state is not used because no single independent candidate exists in
the current reviewed evidence. The rejected state is not used because there are
still materially different possible routes, but they require a fresh user
premise or an explicit scope change before implementation.
