# Futures Derivatives Context Zero-Trade Context-Audit Brief

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval`.

This documentation-only brief converts the passed source audit in
`docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md` into a
decision-complete plan for a later zero-trade derivatives context audit. It
decides the exact later question, the derived context inputs, the anti-leakage
rules, the expected artifacts, the cohort gates, and the stop states.

This brief does not authorize implementation. The next step after this brief is
explicit user approval or rejection of the context-audit implementation scope.
No Go code, CLI flag, generated result directory, audit run, source download,
network request, data write under `../binance-bot/data/derivatives/`, context
feature, label, cohort, ranking, entry, exit, P&L strategy backtest, optimizer
grid, replay, walk-forward, portfolio construction, paper/testnet/live path,
exchange API, credential, deploy file, broad mining, martingale, averaging
down, two-exchange logic, or closed-family rescue is approved by this document
alone.

## Authority Chain

The immediate approval chain is:

1. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md` stopped
   automatic BTCUSDT-only price-range audit work at
   `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
2. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md` closed
   BTC regime plus ETH/SOL context at
   `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md` scoped
   the derivatives source/alignment question and stopped before code.
4. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`
   materialized the durable mark/index/premium `5m` source files and passed at
   `derivatives_context_source_materialization_passed_ready_for_source_audit_approval`.
5. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_AUDIT_REVIEW.md` validated those
   files, SHA-256-bound their provenance, and proved anti-lookahead alignment to
   the `5m` candle anchors, passing at
   `derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.

This branch is a context-separation question only. It is not a strategy premise,
not a context-gain promotion, and not a rescue path for any closed family.

## Why Basis/Premium Context Is Materially Different (Not A Closed-Family Rescue)

Every closed price-only family (structured compression, breakout-retest/
acceptance, clean breakout continuation, hold-inside/midline, impulse
absorption, higher-timeframe nested range rotation, `range_occupancy_rotation_v1`,
and `router_gated_boundary_reclaim_rotation_v1`) used only the OHLCV geometry of
the traded candle. The BTC regime plus ETH/SOL context family also failed, and
it conditioned on BTCUSDT candle price/volume geometry as cross-symbol context.

Mark-minus-index basis and premium-index level are a different, orthogonal
market-data source. They measure the dislocation between the perpetual mark and
the underlying index, i.e. leverage/funding/crowding pressure, which is not
derivable from the traded symbol's candle geometry or from BTC candle geometry.
The later audit therefore conditions local range states on a new orthogonal
source, not on a reslice, retune, rename, or gate relaxation of any closed
family.

This material-difference claim is a gate, not a guarantee. The later audit must
reject itself if the derived basis/premium buckets turn out to be dominated by,
or collinear with, the local price/volume state it already has (see the
orthogonality gate under Cohort Gates and the
`...rejected_closed_family_rescue` stop state). Conditioning on a new source is
only legitimate while that source adds information beyond closed price geometry.

## Minimum Audit Question

A later implementation, only after explicit user approval, may answer only:

```text
Do mark-minus-index basis, premium-index level, or basis-change context buckets,
known at closed decision-candle time, improve separation of BTCUSDT, ETHUSDT, and
SOLUSDT local range states into usable, toxic, rotation, continuation, or
no-trade contexts beyond the local price/volume range state alone?
```

This is a context-separation question. It is not a basis-tradability question,
not an entry question, not a strategy-P&L question, not a replay question, and
not a portfolio-construction question. The later audit must not ask whether a
basis signal can be traded and must not measure entry, exit, adverse/favorable
fill, or P&L.

## Source Contract

The only approved inputs for the later audit are the `9` validated derivatives
source CSVs plus the `3` candle anchors. No new symbol, interval, family, source
download, or source materialization is approved.

Derivatives context source files (durable, materialized, SHA-256 bound; hashes
match `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md` and the
source audit's recomputed provenance):

```text
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
```

Per-stream lag-aligned coverage proven by the source audit (lag-aligned decision
candles / `573,984`):

| Stream | Role | Required | Lag coverage |
| --- | --- | --- | ---: |
| mark BTCUSDT | basis input | yes | `0.995976` |
| index BTCUSDT | basis input | yes | `0.994472` |
| mark ETHUSDT | basis input | yes | `0.998984` |
| index ETHUSDT | basis input | yes | `0.998486` |
| mark SOLUSDT | basis input | yes | `0.996477` |
| index SOLUSDT | basis input | yes | `0.998486` |
| premium BTCUSDT | optional cross-check | no | `0.996470` |
| premium ETHUSDT | optional cross-check | no | `0.996472` |
| premium SOLUSDT | optional cross-check | no | `0.996973` |

Required-stream minimum lag coverage is `0.994472` (index BTCUSDT). Mark-minus-
index basis depends on the two required streams per symbol, so per-symbol basis
context coverage is bounded by the weaker of that symbol's mark/index streams.

Candle alignment anchors only (the local range-state and label source; not
derivatives context sources):

```text
../binance-bot/data/{btcusdt,ethusdt,solusdt}_futures_um_5m_2021_2026.csv
```

Anchor facts to preserve: each anchor has `573,984` loaded candles spanning
`2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`; sorted streams had
`gap_count=0` and `duplicate_count=0`; zero-volume counts were BTC `66`, ETH
`47`, SOL `47`; SOL had one physical non-monotonic row and is accepted only
after sorting by open time.

If any of these source files, provenance hashes, coverage facts, or anchor facts
are missing, contradicted, or impossible to reproduce in the later
implementation, stop at:
`derivatives_context_zero_trade_context_audit_source_gap`.

## Anti-Lookahead And Finality Contract (carried forward)

The later implementation must reuse the passed source audit's loader and
alignment facts. Publication lag for mark/index/premium klines is unproven, so
the conservative one-native-interval (`5m`) lag view is mandatory:

- a derivatives source interval may serve decision candle `D` only when
  `source_close_time + 5m <= decision_candle_close_time` (decision candle `D` is
  served only by the source interval at `D - 5m`);
- no forward fill (`MaxExtraStalenessIntervals=0`), no interpolation, no
  nearest-future joins, no future source revisions;
- every basis/premium context value must be computed from closed source rows
  only, and any basis-change lookback window must end at or before the lagged
  source close;
- missing context must produce missingness or skip rows, never silent defaults
  (no implicit `basis=0`, no carry-forward of the last basis, no zero-premium
  assumption);
- the per-symbol bounded missingness from the source audit stays bounded; a
  decision candle whose lagged basis/premium context is missing is recorded as a
  context skip and excluded from basis-conditioned cohorts, not filled.

## Allowed Derived Context Inputs

The later audit may build only high-level closed-candle derivatives context
buckets from the lagged source rows. All inputs below are conditioning variables
for local range states, never trades.

Primary basis context (from required mark + index streams):

- mark-minus-index basis level bucket, where basis is computed from the lagged
  mark and index closes for the same symbol and interval (raw and/or normalized
  to index, e.g. basis in basis points);
- basis sign/regime bucket (premium versus discount of perp to index);
- basis-change bucket over a closed lookback window ending at or before the
  lagged source close;
- basis dispersion/volatility bucket over a closed lookback window.

Optional corroborating context (from optional premium streams):

- premium-index level bucket;
- premium sign bucket.

Premium-index context is an optional cross-check only. It may corroborate or
refine basis buckets but may not be the sole required context, and a passing
result may not depend on premium streams alone.

Local range-state layer (from the candle anchors, reusing current range-state
construction infrastructure):

- per-symbol local range geometry, volatility state, trend state, impulse state,
  participation proxy, and close location relative to the active range, all from
  closed candles only and identical in spirit to the existing range-state
  construction loop.

The audit may define state IDs in this shape:

```text
derivatives_context_v1::<symbol>::<timeframe>::<local_range_bucket>::<basis_level_bucket>::<basis_change_or_premium_bucket>
```

State IDs must use only data known at closed decision-candle time. Basis and
premium buckets in a state ID must be built only from source rows satisfying
`source_close_time + 5m <= decision_candle_close_time`.

## Symbol Roles

All three symbols may be local range-state authority candidates in this
zero-trade audit, because the conditioning variable is a new orthogonal
derivatives source per symbol rather than the closed price-only or BTC-regime
geometry. BTCUSDT is no longer restricted to regime context here; it carries its
own basis/premium context.

This does not reopen any closed family. BTCUSDT price-only range premises remain
closed, the BTC regime plus ETH/SOL context family remains closed, and no symbol
may be promoted to a traded strategy from this audit. Any future passing context
result would still need a separate strategy premise spec before entries, exits,
P&L backtests, replay, walk-forward, or packaging.

## Anti-Leakage Rule

Forward labels may appear only in label, cohort, ranking, and summary artifacts.
They must never be premise inputs, state-ID inputs, context inputs, gating
inputs, or feature-bucket inputs.

The later implementation must make the time split explicit:

- local range-state and derivatives context features end at the closed decision
  candle, and derivatives context additionally respects the one-`5m`-interval
  source lag;
- labels begin strictly after that decision candle;
- any future return, adverse excursion, favorable excursion, or resolution label
  belongs only to label/cohort artifacts;
- missing derivatives context produces missingness/skip rows, not silent
  defaults, so absence of basis is never encoded as a value.

If hidden future-label input, future-row basis input, or any lookahead is
discovered, stop at:
`derivatives_context_zero_trade_context_audit_rejected_future_label_leak`.

## Expected Future Artifacts

If explicitly approved later, the implementation should write generated outputs
under:

```text
results/futures-derivatives-context-audit/
```

Expected audit-specific artifacts:

- `futures_derivatives_context_sources.csv/json`;
- `futures_derivatives_context_coverage.csv/json`;
- `futures_derivatives_context_basis_features.csv/json`;
- `futures_derivatives_context_local_states.csv/json`;
- `futures_derivatives_context_labels.csv/json`;
- `futures_derivatives_context_cohorts.csv/json`;
- `futures_derivatives_context_rankings.csv/json`;
- `futures_derivatives_context_missingness.csv/json`;
- `futures_derivatives_context_summary.csv/json`.

Common outputs must remain zero-trade compatible:

- `summary.json` and `summary.csv` must be valid with `0` trades;
- `trades.json` must contain no trades;
- common outputs must not be repurposed to hold context, basis, or label rows as
  pseudo-trades.

## Cohort Gates For Later Implementation

The future audit must rank context cohorts by separation quality, not P&L. The
later implementation may define exact numeric thresholds, but the gates must
preserve these rules:

- each candidate cohort must have adequate full-sample and split counts for each
  symbol and local range state it claims to describe;
- a useful cohort must survive the weakest split, not only the full sample;
- BTCUSDT, ETHUSDT, and SOLUSDT may not each rely on the same single fragile
  period;
- toxic/no-trade separation and useful/rotation/continuation separation must be
  reported separately;
- orthogonality gate: a derivatives context bucket must add separation versus the
  local price/volume range state alone; a bucket that is collinear with or
  dominated by the local state contributes no usable context and must fail;
- missing basis/premium context must produce skips, never a silent default
  bucket; cohorts that depend materially on missing-context rows, or that breach
  the source-audit-proven bounded coverage (required basis context coverage no
  worse than the `0.994472` floor), must fail;
- premium-index context is corroboration only; no cohort may pass on premium
  streams alone;
- no result may be treated as strategy evidence without a later strategy premise
  spec.

The later audit may stop with a passing context result only if it can show
closed-candle, lag-correct basis/premium buckets improve BTC/ETH/SOL local
range-state separation, beyond local price/volume state alone, without
future-label leakage and without rescuing any closed family.

## Rejection Criteria

Reject the later implementation or stop immediately if it becomes any of these:

- a price-only or BTC-regime reslice, retune, rename, or gate relaxation of any
  closed family (structured compression, breakout-retest/acceptance, clean
  breakout continuation, hold-inside/midline, impulse absorption, higher-
  timeframe nested range rotation, `range_occupancy_rotation_v1`,
  `router_gated_boundary_reclaim_rotation_v1`, or BTC regime plus ETH/SOL
  context);
- a basis-tradability test, entry, exit, fill model, or P&L measurement;
- use of derivatives context to rescue or revive any closed family;
- broad symbol, interval, or source-family mining;
- unapproved source expansion, source download, network request, or any data
  write under `../binance-bot/data/derivatives/`;
- forward fill, interpolation, nearest-future joins, future source revisions, or
  any basis/premium value built from rows newer than `D - 5m`;
- silent default context where basis/premium is missing;
- hidden future-label input or any lookahead in state/context/gating inputs;
- entries, exits, P&L strategy backtests, optimizer grids, replay, walk-forward,
  or portfolio construction;
- spread-range or pair-range engine work;
- paper, testnet, live, exchange API, credential, or deploy work;
- martingale, averaging down, or two-exchange logic;
- import of old `binance-bot` strategy, scoring, order-management, live,
  credential, deploy, or portfolio-coordinator logic.

## Stop States

This brief stops at:

- `derivatives_context_zero_trade_context_audit_brief_ready_for_user_approval`.

Allowed later implementation stop states, only after explicit user approval:

- `derivatives_context_zero_trade_context_audit_source_gap`;
- `derivatives_context_zero_trade_context_audit_rejected_future_label_leak`;
- `derivatives_context_zero_trade_context_audit_rejected_closed_family_rescue`;
- `derivatives_context_zero_trade_context_audit_failed_no_usable_context`;
- `derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.

The passing implementation stop state would still not authorize entries, exits,
P&L backtests, replay, walk-forward, packaging, paper/testnet/live paths, source
downloads, strategy promotion, or deployment.

## Later Implementation Verification

Only after explicit user approval, the implementation closeout should include:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-audit -out-dir results/futures-derivatives-context-audit
wc -l results/futures-derivatives-context-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

The future implementation review must record source paths, product, symbol,
source family, interval, timestamp semantics, finality/lag rule, per-stream and
per-symbol basis context coverage, recomputed provenance checksums, local
range-state counts, basis/premium feature counts, label counts, cohort counts,
ranking counts, passing cohort counts, missingness/skip counts, common
zero-trade output status, stop state, and exact command outcomes. The
implementation should reuse the passed source audit's loader and one-`5m`-
interval anti-lookahead alignment and stay zero-trade.

## Current Documentation Closeout Verification

For this brief-writing closeout, run only:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
