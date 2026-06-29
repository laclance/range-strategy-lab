# Futures Derivatives Context Strategy-Premise Spec

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_strategy_premise_spec_ready_for_user_approval`.

This documentation-only spec converts the passed zero-trade derivatives context
audit into one selected later premise track:

```text
BTCUSDT 15m derivatives-context no-trade filter premise
```

The five BTCUSDT `15m` no-trade/toxic cohorts justify a later zero-trade filter
premise audit. The single BTCUSDT `15m` rotation candidate does not justify a
rotation-entry premise, and the evidence does not justify two parallel premise
tracks. The rotation row is preserved as diagnostic context only.

This spec does not authorize implementation. No Go code, CLI flag, generated
result directory, audit run, source download, network request, source
materialization, data write under `../binance-bot/data/derivatives/`, entry,
exit, P&L backtest, optimizer grid, replay, walk-forward, portfolio
construction, paper/testnet/live path, exchange API, credential, deploy file,
broad mining, martingale, averaging down, two-exchange logic, or closed-family
rescue is approved by this document alone.

## Authority Chain

The immediate authority chain is:

1. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`
   materialized durable Binance USDT-M futures `5m` mark/index/premium source
   files under `../binance-bot/data/derivatives/`.
2. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md` validated those
   files, SHA-256-bound provenance, and established the conservative source lag
   rule.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md`
   defined the separation-only derivatives context audit.
4. `docs/FUTURES_DERIVATIVES_CONTEXT_AUDIT_REVIEW.md` implemented that
   zero-trade audit and passed at
   `derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.

The context audit is separation evidence only. Basis/premium context is a
closed-candle conditioning source, not a tradable signal by itself.

## Evidence Preserved

All passing context rows were BTCUSDT `15m`; ETHUSDT and SOLUSDT produced `0`
passing cohorts. The audit stayed zero-trade and produced `0` trades.

No-trade/toxic evidence:

| Rank | Horizon | Local state | Derivatives bucket | Rows | Weakest split | Toxic evidence |
| ---: | ---: | --- | --- | ---: | ---: | --- |
| 1 | 48 | `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none` | `basis_discount_small + premium_discount_small` | 515 | 110 | full toxic `0.732039`, worst split toxic `0.800000` |
| 2 | 48 | `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none` | `basis_discount_small` | 622 | 142 | full toxic `0.729904`, worst split toxic `0.802817` |
| 3 | 48 | `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` | `basis_discount_small + basis_change_flat` | 356 | 62 | full toxic `0.662921`, worst split toxic `0.699387` |
| 4 | 48 | `geometry_midline_balanced::vol_compressed::trend_up_pressure::impulse_none` | `basis_discount_small` | 613 | 124 | full toxic `0.654160`, worst split toxic `0.759358` |
| 5 | 48 | `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` | `basis_discount_small + premium_discount_small` | 538 | 115 | full toxic `0.659851`, worst split toxic `0.719212` |

The single rotation candidate:

| Rank | Horizon | Local state | Derivatives bucket | Rows | Weakest split | Useful evidence |
| ---: | ---: | --- | --- | ---: | ---: | --- |
| 6 | 24 | `geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale` | `basis_discount_small` | 313 | 71 | useful `0.632588` full / `0.521127` weakest split, margin improvement `0.072540` full / `0.025844` weakest split |

## Premise Decision

The selected premise is a no-trade filter premise, not a rotation-entry premise.

Why the filter track is justified:

- five independent passing rows share a toxic/no-trade interpretation;
- those rows survive weakest-split checks with at least `62` rows in the weakest
  split;
- full toxic rates range from `0.654160` to `0.732039`;
- worst-split toxic rates range from `0.699387` to `0.802817`;
- each row improved separation versus the local price/volume state alone;
- the common motif is interpretable: BTCUSDT `15m` midline-balanced,
  volume-compressed local range states become toxic when lagged basis context is
  in a small discount, sometimes corroborated by premium discount or flat basis
  change.

Why rotation is not selected:

- only one rotation row passed;
- its weakest-split useful rate was only `0.521127`;
- its weakest-split margin improvement was only `0.025844`;
- no entry trigger, exit model, fill model, P&L evidence, replay, or
  walk-forward evidence exists;
- converting one diagnostic row into an entry premise would be a direct
  overreach from a zero-trade context audit.

Why two tracks are not selected:

- the no-trade and rotation rows answer different questions;
- combining them would blur toxic filtering with entry selection;
- the rotation evidence is too thin to deserve equal implementation authority;
- a combined track would risk quietly rescuing the closed router-rotation and
  BTCUSDT price-only families.

The rejected alternative stop state,
`derivatives_context_strategy_premise_spec_rejected_no_strategy_premise`, is not
used because the five toxic cohorts are coherent enough to justify one later
zero-trade no-trade filter audit.

## Premise Definition

The later premise may test only this statement:

```text
When BTCUSDT 15m local range state is midline-balanced and volume-compressed,
lagged small-discount basis context identifies toxic/no-trade conditions strongly
enough to define a closed-candle no-trade veto candidate.
```

The premise is a veto candidate, not an entry signal. It may say "do not trade
this context" only after a later audit proves the filter definition remains
stable. It may not say "trade the opposite", "fade basis", "trade basis", "enter
rotation", "enter continuation", or "improve P&L".

## Initial Filter Candidates

The later audit must begin from the exact five passing toxic rows. It may report
both exact-row filters and a de-duplicated canonical union, but it may not
generalize beyond these contexts unless the zero-trade audit proves that the
broader candidate preserves the declared gates.

Initial exact candidates:

1. `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none`
   with `basis_discount_small + premium_discount_small`;
2. `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none`
   with `basis_discount_small`;
3. `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` with
   `basis_discount_small + basis_change_flat`;
4. `geometry_midline_balanced::vol_compressed::trend_up_pressure::impulse_none`
   with `basis_discount_small`;
5. `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` with
   `basis_discount_small + premium_discount_small`.

Canonicalization rules for the later audit:

- the `trend_down_pressure` exact premium-confirmed row is nested inside the
  broader `trend_down_pressure + basis_discount_small` row and must be reported
  as overlap, not double-counted;
- the flat-trend rows require their corroborating bucket
  (`basis_change_flat` or `premium_discount_small`) unless a later zero-trade
  result proves a broader flat-trend basis-discount filter still passes;
- no ETHUSDT or SOLUSDT filter may be inferred from this evidence;
- no `4h`, `1h`, `5m`, or other timeframe filter may be inferred from this
  evidence;
- no rotation-entry rule may be inferred from the no-trade filter candidates.

## Source And Finality Contract

The later audit may use only local Binance USDT-M futures BTCUSDT data and the
already validated derivatives context rows needed by this selected premise.

Allowed inputs:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv
```

Carried-forward source facts:

- BTCUSDT candle anchor: Binance USDT-M futures `5m`, `573,984` loaded candles,
  `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`, `gap_count=0`,
  `duplicate_count=0`, `zero_volume_count=66`;
- source-audit lag coverage: mark BTCUSDT `0.995976`, index BTCUSDT
  `0.994472`, premium BTCUSDT `0.996470`;
- required-stream floor for the derivatives source family remains `0.994472`;
- missing lagged context is bounded recorded missingness, never a fillable value.

Finality rules:

- confirmed closed-candle decisions only;
- derivatives source rows may serve decision candle `D` only when
  `source_close_time + 5m <= decision_candle_close_time`;
- no forward fill, interpolation, nearest-future joins, future source revisions,
  or silent default context;
- labels begin strictly after the decision candle and may appear only in label,
  cohort, ranking, and summary artifacts;
- basis/premium context must be a conditioning source only, not a traded signal.

## Allowed Later Implementation

Only after separate explicit user approval, a later implementation may add a
zero-trade no-trade filter premise audit. Suggested flag and output directory:

```text
-futures-derivatives-no-trade-filter-premise-audit
results/futures-derivatives-no-trade-filter-premise-audit/
```

Allowed audit question:

```text
Do the five BTCUSDT 15m toxic derivatives-context cohorts define a stable
closed-candle no-trade veto candidate after exact-row reproduction,
de-duplication, split checks, and missingness accounting?
```

Allowed generated artifacts:

- source/provenance and coverage rows for the BTCUSDT candle, mark, index, and
  premium inputs;
- exact filter definition rows for the five toxic cohorts;
- canonical union and overlap rows that prevent double counting;
- lagged basis/premium feature rows needed by those filters;
- local-state rows and veto-candidate rows;
- forward labels only as evaluation metadata;
- split stability and toxic-rate lift summaries versus local-only baselines;
- missingness and skip rows;
- zero-trade common outputs (`summary.json`, `summary.csv`, `trades.json`).

The later audit must remain zero-trade. It may count veto candidates and label
outcomes, but it may not simulate entries, exits, fills, adverse/favorable
execution, P&L, replay, walk-forward, or portfolio behavior.

## Later Audit Gates

A later no-trade filter premise audit may pass only if all of these hold:

- the BTCUSDT source/provenance/coverage facts reproduce, including the
  one-`5m` source lag and no-fill policy;
- all five exact toxic context rows are reproduced as candidates before any
  canonical union is evaluated;
- each exact candidate preserves full-sample rows at or above `300` and weakest
  split rows at or above `60`;
- each exact candidate preserves full toxic rate at or above `0.65` and worst
  split toxic rate at or above `0.69`;
- each exact candidate preserves toxic-rate improvement versus its local-only
  baseline, with full improvement at or above `0.04`;
- the canonical de-duplicated filter union is reported with overlap counts and
  does not rely on double-counting the nested `trend_down_pressure` premium row;
- the filter union remains toxic/no-trade dominated in full sample and every
  split;
- useful/rotation labels blocked by the veto are reported as collateral damage,
  not hidden;
- missing lagged context is skipped and counted, never encoded as a default
  bucket;
- the audit proves it is not a price-only or BTC-regime reslice and does not
  reopen, retune, rename, gate-relax, or promote any closed family.

Failure to preserve the exact toxic rows, failure of the canonical union, hidden
future-label input, closed-family rescue, or any attempt to convert the rotation
candidate into an entry premise must fail the audit.

## Stop States

This spec stops at:

- `derivatives_context_strategy_premise_spec_ready_for_user_approval`.

Rejected docs-only alternative:

- `derivatives_context_strategy_premise_spec_rejected_no_strategy_premise`.

Allowed later implementation stop states, only after separate explicit user
approval:

- `derivatives_context_no_trade_filter_premise_audit_source_gap`;
- `derivatives_context_no_trade_filter_premise_audit_rejected_future_label_leak`;
- `derivatives_context_no_trade_filter_premise_audit_rejected_closed_family_rescue`;
- `derivatives_context_no_trade_filter_premise_audit_rejected_rotation_entry_rescue`;
- `derivatives_context_no_trade_filter_premise_audit_failed_no_usable_filter`;
- `derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec`.

The passing implementation stop state would still not authorize entries, exits,
P&L backtests, optimizer grids, replay, walk-forward, packaging, source
downloads, source expansion, symbol expansion, paper/testnet/live paths,
exchange API work, credentials, deploy files, martingale, averaging down, or
two-exchange logic.

## Boundaries

The selected no-trade filter premise does not reopen any closed family:

- BTCUSDT price-only range families remain closed in their reviewed forms;
- BTC regime plus ETH/SOL context remains closed at `0` passing cohorts;
- `router_gated_boundary_reclaim_rotation_v1` remains closed and may not be
  rescued by the single derivatives rotation candidate;
- spread-range/pair-range remains parked;
- volatility-aware exits remain unavailable until an independent entry premise
  first shows gross edge before costs.

No old `binance-bot` strategy, scoring, order-management, portfolio coordinator,
live, credential, deploy, grid, martingale, averaging-down, or two-exchange
logic may be imported.

## Current Documentation Closeout Verification

For this docs-only spec closeout, run:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
