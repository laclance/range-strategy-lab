# Futures Derivatives Context No-Trade Filter Integration Spec

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise`.

This documentation-only spec preserves
`btc_15m_basis_discount_no_trade_veto_v1` as a future veto candidate, but it
does not select an implementation gate now. The passed zero-trade no-trade
filter premise audit proved that the declared BTCUSDT `15m` basis-discount
contexts are toxic/no-trade dominated. It did not prove there is an independent
entry stream worth filtering.

Integration is therefore deferred until a separate, independently approved
entry premise exists. A veto cannot be integrated in isolation without becoming
either a repeated label audit or an accidental entry/P&L claim. The filter may
only ever say "skip this already-approved candidate context"; it may not say
"trade the opposite", "trade basis", "fade basis", "enter rotation", "enter
continuation", or "improve P&L".

This spec authorizes no Go code, CLI flag, generated result directory, audit
run, source download, network request, source materialization, data write under
`../binance-bot/data/derivatives/`, entry, exit, P&L backtest, optimizer grid,
replay, walk-forward, portfolio construction, paper/testnet/live path, exchange
API, credential, deploy file, martingale, averaging down, two-exchange logic,
closed-family rescue, or strategy promotion.

## Authority Chain

The immediate authority chain is:

1. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`
   materialized the local Binance USDT-M futures mark/index/premium sources.
2. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md` validated
   provenance and established the conservative one-`5m` source lag.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md` proved that lagged
   derivatives context separated BTCUSDT `15m` local range states beyond the
   local price/volume state alone.
4. `docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md` selected only a
   BTCUSDT `15m` no-trade filter premise and rejected rotation-entry and
   two-track alternatives.
5. `docs/FUTURES_DERIVATIVES_CONTEXT_NO_TRADE_FILTER_PREMISE_AUDIT_REVIEW.md`
   reproduced the five exact toxic rows and passed the canonical veto union as
   a zero-trade filter premise.

That chain proves a veto premise, not an entry premise.

## Decision

The passed veto is held as a future integration candidate only:

```text
btc_15m_basis_discount_no_trade_veto_v1
```

It is not rejected, because the premise audit found stable toxic/no-trade
dominance in full sample and every split. It is not implementation-ready,
because there is no separately approved candidate-entry stream for it to
suppress. The next selected state is deferred, not ready-for-implementation.

The diagnostic BTCUSDT `15m/h24` rotation row remains diagnostic only. ETHUSDT,
SOLUSDT, other timeframes, broad derivatives source expansion, and closed
BTCUSDT price-only or BTC-regime families receive no authority from this spec.

## Preserved Evidence

The no-trade filter premise audit passed with these facts:

- canonical filter ID: `btc_15m_basis_discount_no_trade_veto_v1`;
- exact candidates passed: `5`;
- de-duplicated veto rows: `1,823`;
- overlap rows: `821`;
- no-trade toxic rows: `1,241`;
- full toxic rate: `0.680746`;
- minimum split toxic rate: `0.665485`;
- weakest split rows: `387`;
- full toxic improvement versus local-only baseline: `0.046269`;
- reported full-sample collateral: `311` rotation-useful rows and `271`
  continuation-useful rows blocked;
- common outputs stayed zero-trade: `summary.json`, `summary.csv`, and
  `trades.json` reported `0` trades.

The five exact rows remain the only definition source:

| Local state motif | Required derivatives context | Horizon | Rows | Weakest split | Full toxic |
| --- | --- | ---: | ---: | ---: | ---: |
| trend down pressure | `basis_discount_small + premium_discount_small` | 48 | 515 | 110 | 0.732039 |
| trend down pressure | `basis_discount_small` | 48 | 622 | 142 | 0.729904 |
| trend flat | `basis_discount_small + basis_change_flat` | 48 | 356 | 62 | 0.662921 |
| trend up pressure | `basis_discount_small` | 48 | 613 | 124 | 0.654160 |
| trend flat | `basis_discount_small + premium_discount_small` | 48 | 538 | 115 | 0.659851 |

The flat-trend rows keep their corroborating bucket requirements. The
trend-down premium-confirmed row remains nested inside the broader
trend-down/basis-discount row and must be counted as overlap, never as a second
independent veto row.

## Why Integration Is Deferred

A no-trade filter needs an independent thing to filter. The current repo state
has a passed veto premise, but no live, approved entry premise with:

- an independently defined candidate-event stream;
- a declared side, timing, and decision candle;
- pre-filter evidence that does not depend on this veto;
- a reason to test skipped versus retained candidate contexts;
- approval to run an interaction audit.

Testing the veto again against all BTCUSDT `15m` range-state rows would only
repeat the premise audit. Testing it through P&L would manufacture strategy
evidence from a filter. Converting the toxic context into a short, a basis fade,
or a rotation/continuation entry would violate the premise chain.

## Conditional Future Integration Question

No future implementation is selected by this spec. If a separate independent
entry premise is later approved, the first possible interaction question may be
only:

```text
Within an independently approved BTCUSDT candidate-entry context stream, does
btc_15m_basis_discount_no_trade_veto_v1 identify candidate rows that should be
skipped as toxic/no-trade, without changing the entry definition and without
making a P&L claim?
```

That question is conditional. It requires a later explicit user approval and a
fresh brief or audit spec that names the independent entry premise. This spec by
itself does not authorize that implementation.

## Conditional Future Inputs

If an independent entry premise is later approved for interaction testing, the
veto side of the interaction may use only the already validated BTCUSDT inputs:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv
```

The carried-forward source facts remain:

- BTCUSDT candle anchor: Binance USDT-M futures `5m`, `573,984` loaded candles,
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`, `gap_count=0`,
  `duplicate_count=0`, `zero_volume_count=66`;
- BTCUSDT mark rows `571,675`, gaps `6`, SHA-256
  `424c05ca880a31270eea1286d6cdd96ac1132d848e8f5e9d6b3b7177bb7c2858`;
- BTCUSDT index rows `570,812`, gaps `8`, SHA-256
  `7ba5a375311e0324dab38f18c2a7137376b619ce63cfb01a17e5684c58390aca`;
- BTCUSDT premium rows `571,959`, gaps `7`, SHA-256
  `094e610617f812f032e6a68b3ae6186b20359592415cdb678845c5d287ec298c`;
- required lag coverage floor `0.994472`.

The finality rule remains:
`source_close_time + 5m <= decision_candle_close_time`. No forward fill,
interpolation, nearest-future join, future source revision, or silent default
context is allowed. Missing lagged context must be skipped and counted.
Forward labels may be used only as evaluation metadata.

Any later interaction must stay in the BTCUSDT `15m` decision-candle family
unless a separate approved spec proves a non-leaking projection. This spec does
not authorize projection to `5m`, `1h`, `4h`, ETHUSDT, SOLUSDT, or any new
source family.

## Conditional Future Outputs

If a later approved interaction audit exists, it must report at least:

- reproduction of the canonical veto definition and source/coverage facts;
- the independent entry premise's candidate-event rows before the veto is
  applied;
- veto match rows, retained candidate rows, and skipped candidate rows;
- missingness and skipped-context counts;
- overlap between exact veto components and the canonical union;
- forward-label summaries for skipped versus retained candidate contexts;
- split stability for skipped versus retained candidate contexts;
- collateral labels blocked by the veto, including rotation-useful and
  continuation-useful outcomes;
- zero-trade common outputs.

The interaction may annotate candidate rows as skipped or retained. It may not
create entries, alter the entry premise, simulate fills, model stops/targets,
score P&L, optimize, replay, walk forward, or promote a strategy.

## Conditional Future Fail Gates

A later interaction audit must fail or defer if any of these occur:

- no independent entry premise exists;
- the candidate-entry stream is created from the veto itself;
- the veto is used as an entry, side, exit, sizing, ranking, or P&L rule;
- the entry premise changes after seeing veto results;
- forward labels enter feature, state, context, veto, or entry inputs;
- missing source context is filled, interpolated, nearest-future joined, or
  encoded as a default bucket;
- the canonical veto definition cannot be reproduced;
- overlap or collateral damage is hidden;
- ETHUSDT, SOLUSDT, other timeframes, new source families, source downloads, or
  source writes are introduced;
- the work reopens, retunes, renames, gate-relaxes, or promotes a closed family.

## Stop States

This spec uses:

```text
derivatives_context_no_trade_filter_integration_spec_deferred_until_entry_premise
```

The ready-for-implementation stop state is intentionally not used:

```text
derivatives_context_no_trade_filter_integration_spec_ready_for_user_approval
```

The rejection stop state is also not used:

```text
derivatives_context_no_trade_filter_integration_spec_rejected_no_integration_premise
```

The veto evidence is coherent enough to keep. It is simply waiting for an
independent entry premise before any integration audit can be meaningful.
