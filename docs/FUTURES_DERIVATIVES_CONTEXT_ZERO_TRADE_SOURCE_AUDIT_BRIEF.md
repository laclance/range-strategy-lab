# Futures Derivatives Context Zero-Trade Source Audit Brief

Date: 2026-06-28

## Verdict

Stop state:
`derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.

This documentation-only brief converts
`docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md` into a
decision-complete plan for a later zero-trade derivatives source audit.

This brief does not authorize implementation. The next step after this brief is
explicit user approval or rejection of the source-audit implementation scope. No
Go code, CLI flag, generated result directory, source download, source
materialization, source parser, context feature implementation, entry, exit,
P&L strategy backtest, optimizer grid, replay, walk-forward, paper/testnet/live
path, exchange API, credential, deploy file, broad mining, martingale,
averaging down, or two-exchange logic is approved by this document alone.

## Authority Chain

The immediate approval chain is:

1. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md` stopped
   automatic BTCUSDT-only price-range audit work at
   `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
2. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md` closed
   BTC regime plus ETH/SOL context at
   `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md` parked
   derivatives context as market-data source expansion only, pending explicit
   source approval.
4. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md` approved only this
   separate zero-trade source-audit brief-writing step.

This branch is a source and alignment question. It is not a strategy premise,
not a context-gain result, and not a rescue path for a closed family.

## Exact Later Source-Audit Question

A later implementation, only after explicit user approval, may answer:

```text
Can durable local/offline Binance USDT-M mark-price, index-price, or
premium-index kline source rows for BTCUSDT, ETHUSDT, and SOLUSDT be validated,
provenance-bound, and anti-lookahead-aligned to the existing Binance USDT-M
futures 5m candle anchors well enough to justify a separate later zero-trade
context-audit brief?
```

The later source audit must not ask whether a basis signal can be traded. It
must not measure entry/exit performance, P&L, strategy returns, replay,
walk-forward behavior, or portfolio construction.

## Proposed Source Family

The first source family is:

```text
Binance USDT-M futures mark/index/premium basis klines
```

The family may include only these public market-data concepts:

- mark-price klines;
- index-price klines;
- premium-index klines;
- derived mark minus index basis, premium-index level, or basis change rows
  computed only for source validation and alignment reporting.

The later source audit may not include funding rate history, aggregate trades,
taker flow, open interest, long/short ratios, order book/depth, liquidations, or
force-order rows. Funding remains the second candidate only if this basis source
family is rejected or closed by a separate later review. Aggregate trades remain
parked as a high-volume secondary candidate.

## Durable Local/Offline Source Requirements

No derivatives market-data source rows are approved as direct lab inputs today.
The later implementation must first prove one of these conditions:

- durable local/offline mark/index/premium files already exist under an
  approved local data directory; or
- the user has explicitly approved an offline materialization plan before the
  implementation runs.

Acceptable durable local file scope is limited to `../binance-bot/data/` or a
documented subdirectory under it. Preferred path shape:

```text
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_ETHUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_SOLUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_ETHUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_SOLUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_ETHUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_SOLUSDT_2021_2026.csv
```

Equivalent names may be accepted only if a source manifest maps the actual
paths to the same product, symbol, interval, source family, and coverage.

For every candidate file, the later audit must record:

- durable path;
- source family;
- source owner or archive family;
- symbol;
- product scope;
- native interval;
- timestamp semantics;
- first and last timestamp;
- row count;
- duplicate count;
- gap or missing-interval count;
- non-monotonic physical row count;
- timezone assumption;
- publication/finality rule;
- max staleness;
- missing-data policy;
- checksum or provenance identifier;
- comparison-only status;
- validation status.

Reject `/tmp` caches, adjacent research result CSVs, private endpoints, signed
requests, WebSocket streams, spot sources, cross-exchange sources, or any source
whose product, symbol, timestamp, finality, provenance, or missing-data behavior
cannot be explained.

## Candle Alignment Anchor Contract

The later source audit may use only these approved Binance USDT-M futures `5m`
candle files as alignment anchors:

| Symbol | Path | Role |
| --- | --- | --- |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor |

Known accepted source facts:

- each candle anchor has `573,984` loaded candles;
- coverage is `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
- sorted streams had `gap_count=0` and `duplicate_count=0`;
- zero-volume counts were BTC `66`, ETH `47`, and SOL `47`;
- SOL had one physical non-monotonic row and was accepted only after sorting.

The candle files are not derivatives context sources. They are anchors for
coverage, timestamp, and anti-lookahead alignment only.

## Timestamp, Finality, And Staleness Requirements

The later audit must treat time as UTC and record whether the source timestamp
is an open time, close time, event time, or publication time.

For kline-like mark/index/premium rows:

- a row may be aligned to a decision candle only when the source observation is
  known at or before that decision candle close;
- if a source row describes a completed interval, the implementation must use
  `source_close_time + publication_lag <= decision_candle_close_time`;
- if publication lag cannot be proven, the source must either be rejected or
  reported only in a conservative one-native-interval-lag alignment view;
- nearest-future joins are forbidden;
- interpolation from future rows is forbidden;
- forward fill requires an explicit max staleness;
- missing context must produce missingness or skip rows, not silent defaults.

For `5m` basis klines, preferred max staleness is one closed `5m` interval. Any
longer staleness must be reported as a rejection candidate and cannot justify a
passing source decision without a separate review.

## Allowed Future Implementation Artifacts

Only after explicit user approval, the source audit may write under:

```text
results/futures-derivatives-context-source-audit/
```

Allowed source-audit artifacts:

- `futures_derivatives_context_source_audit_sources.csv/json`;
- `futures_derivatives_context_source_audit_candle_anchors.csv/json`;
- `futures_derivatives_context_source_audit_external_coverage.csv/json`;
- `futures_derivatives_context_source_audit_timestamp_alignment.csv/json`;
- `futures_derivatives_context_source_audit_publication_lag.csv/json`;
- `futures_derivatives_context_source_audit_missingness.csv/json`;
- `futures_derivatives_context_source_audit_provenance.csv/json`;
- `futures_derivatives_context_source_audit_skips.csv/json`;
- `futures_derivatives_context_source_audit_summary.csv/json`.

Common outputs may be written only in zero-trade-compatible form:

- `summary.json` and `summary.csv` must be valid with `0` trades;
- `trades.json` must contain no trades;
- common outputs must not store source rows as pseudo-trades.

Forbidden implementation artifacts in this first source audit:

- context feature/cohort/ranking artifacts that claim context gain;
- labels, forward returns, adverse/favorable excursion rows, or trade outcome
  rows;
- entry, exit, P&L, replay, walk-forward, optimizer, portfolio, paper/testnet,
  live, deploy, credential, or exchange-order artifacts.

## Anti-Lookahead Join Rules

The later implementation must make the source timeline auditable:

1. Validate source rows before alignment.
2. Align by symbol, product, interval, and UTC timestamp.
3. Use only source observations whose known/publication time is less than or
   equal to the decision candle close.
4. If source rows are native `5m`, require exact closed-interval alignment or
   report the row as missing.
5. If source rows are lower cadence or irregular, reject this first basis audit
   unless a separate review approves the cadence.
6. Do not fill from future rows.
7. Do not use labels, future returns, or future source revisions as source
   validation inputs.
8. Report skipped rows, missing rows, stale rows, duplicate rows, and rejected
   rows as first-class artifacts.

## Rejection Criteria

The later implementation must stop or reject if any of these occur:

- no durable local/offline source files or explicitly approved offline
  materialization plan exists;
- implementation requires source downloads, live probes, private endpoints,
  signed requests, API keys, credentials, or WebSocket streams;
- files are under `/tmp` or exist only as adjacent research results;
- product, symbol, interval, schema, timestamp semantics, or source ownership is
  ambiguous;
- source finality or publication lag cannot be explained;
- coverage cannot be aligned to enough of the existing candle anchor era to
  make a source decision;
- duplicate, gap, stale, or missing rows are unbounded or silently filled;
- source rows are spot, cross-exchange, open-interest, long/short, liquidation,
  force-order, order-book/depth, funding, aggregate-trade, or taker-flow rows;
- source alignment is used to rescue structured compression, BTC regime plus
  ETH/SOL context, router-gated boundary reclaim, breakout-retest/acceptance,
  clean breakout, hold-inside/midline, impulse absorption, nested range
  rotation, or another closed family;
- entries, exits, P&L backtests, optimizer grids, replay, walk-forward, paper,
  testnet, live paths, exchange API work, credentials, deploy files,
  martingale, averaging down, or two-exchange logic appear.

## Stop States

This brief stops at:

- `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`.

Allowed later implementation stop states, only after explicit user approval:

- `derivatives_context_zero_trade_source_audit_rejected_source_gap`;
- `derivatives_context_zero_trade_source_audit_rejected_live_or_private_api_path`;
- `derivatives_context_zero_trade_source_audit_rejected_timestamp_or_finality_gap`;
- `derivatives_context_zero_trade_source_audit_rejected_alignment_gap`;
- `derivatives_context_zero_trade_source_audit_rejected_closed_family_rescue`;
- `derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.

A passing source-audit stop state would still not authorize context-gain
implementation, entries, exits, P&L strategy backtests, optimizer grids,
replay, walk-forward, packaging, source downloads, paper/testnet/live paths,
exchange API work, credentials, deploy files, strategy promotion, martingale,
averaging down, or two-exchange logic.

## Future Implementation Verification

Only after explicit user approval, the implementation closeout should include:

```bash
rg --files ../binance-bot/data | rg -i "(mark[_-]?price|index[_-]?price|premium[_-]?index|premiumIndex|basis)"
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-source-audit -out-dir results/futures-derivatives-context-source-audit
wc -l results/futures-derivatives-context-source-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

The future implementation review must record source paths, product, symbol,
source family, interval, timestamp semantics, finality rule, coverage, row
counts, duplicate counts, gap or missing counts, stale counts, checksums or
provenance identifiers, generated artifact paths, common zero-trade output
status, stop state, and exact command outcomes.

## Current Documentation Closeout Verification

For this brief-writing closeout, run only:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
