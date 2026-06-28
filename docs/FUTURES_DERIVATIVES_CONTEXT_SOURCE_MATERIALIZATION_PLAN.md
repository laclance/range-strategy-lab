# Futures Derivatives Context Source Materialization Plan

Date: 2026-06-28

## Verdict

Stop state:
`derivatives_context_source_materialization_plan_ready_for_execution_approval`.

This docs-only plan records the user-approved offline materialization scope for
Binance public Data Vision derivatives context sources. It does not download,
parse, normalize, or audit source rows today.

The next executable step still requires an explicit future approval to run a
materialization command, use network access, and write durable files outside
this repository under:

```text
../binance-bot/data/derivatives/
```

This plan does not authorize context-gain features, labels, cohorts, rankings,
entries, exits, P&L backtests, replay, walk-forward, paper/testnet/live paths,
private endpoints, exchange API keys, credentials, deploy files, broad mining,
strategy promotion, martingale, averaging down, or two-exchange logic.

## Authority Chain

The immediate boundary is:

1. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_SCOPE_REVIEW.md` approved only a
   separate zero-trade source-audit brief-writing task.
2. `docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md` selected
   Binance USDT-M mark/index/premium basis klines as the first source family,
   but required durable local/offline files or an explicitly approved offline
   materialization plan before implementation.
3. The user explicitly approved an offline materialization plan for Binance
   public Data Vision mark-price, index-price, and optional premium-index `5m`
   klines for `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` under
   `../binance-bot/data/derivatives/`.

This branch is source materialization only. It is not a source-audit result, not
a context-audit result, and not a strategy premise.

## Materialization Question

A later execution, only after explicit approval, may answer:

```text
Can the approved Binance public Data Vision USDT-M futures mark-price,
index-price, and optional premium-index 5m archive objects for BTCUSDT,
ETHUSDT, and SOLUSDT be downloaded once, checksum-bound, normalized, and stored
as durable local/offline source files under ../binance-bot/data/derivatives/?
```

The execution must stop after source materialization and provenance. It must
not run the derivatives source audit, context features, labels, cohorts,
rankings, strategy tests, or any P&L path.

## Approved Public Archive Scope

Required source families:

| Family | Binance Data Vision archive family | Required |
| --- | --- | --- |
| mark price klines | `markPriceKlines` | yes |
| index price klines | `indexPriceKlines` | yes |

Optional source family:

| Family | Binance Data Vision archive family | Required |
| --- | --- | --- |
| premium-index klines | `premiumIndexKlines` | no, cross-check only |

Allowed symbols:

```text
BTCUSDT
ETHUSDT
SOLUSDT
```

Allowed interval:

```text
5m
```

Allowed era is bounded by the existing candle anchors:

```text
2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
```

Archive object set:

- monthly objects from `2021-01` through `2026-05`;
- daily tail objects from `2026-06-01` through `2026-06-16`;
- no current-month monthly inference beyond the approved daily tail;
- no symbols, intervals, source families, or date ranges outside this list.

Path template:

```text
https://data.binance.vision/data/futures/um/{monthly|daily}/{archive_family}/{symbol}/5m/{symbol}-5m-{key}.zip
```

Example objects:

```text
https://data.binance.vision/data/futures/um/monthly/markPriceKlines/BTCUSDT/5m/BTCUSDT-5m-2021-01.zip
https://data.binance.vision/data/futures/um/monthly/indexPriceKlines/ETHUSDT/5m/ETHUSDT-5m-2024-12.zip
https://data.binance.vision/data/futures/um/daily/premiumIndexKlines/SOLUSDT/5m/SOLUSDT-5m-2026-06-16.zip
```

Do not use REST tail endpoints, WebSockets, signed/private endpoints, archive
listing as broad discovery, `/futures/data/basis`, funding, aggregate trades,
taker flow, open interest, long/short ratios, liquidations, force orders,
order-book/depth, spot sources, or cross-exchange sources in this plan.

## Expected Object And Download Size

The local adjacent source-proof inventory at
`../binance-bot/research/2026-06-18_futures_perp_basis_reversion/event_study/source_inventory.csv`
recorded the required mark/index object set as:

| Scope | Objects | Compressed bytes |
| --- | ---: | ---: |
| `markPriceKlines` | `243` | `42,833,309` |
| `indexPriceKlines` | `243` | `47,233,885` |
| required mark plus index | `486` | `90,067,194` |

That is about `90.1` MB decimal, or about `85.9` MiB, before durable manifests.

The optional `premiumIndexKlines` family is expected to add another same-shaped
family, approximately `243` archive objects and roughly `45` MB compressed if
available for the full approved era. The planning estimate for required plus
optional families is therefore roughly `135` MB, or about `140` MB with
metadata, retries, and small manifests. Actual execution must record measured
bytes from downloaded files and must not treat this estimate as proof.

## Durable Target Layout

Future execution may write only under:

```text
../binance-bot/data/derivatives/
```

Preferred raw archive layout:

```text
../binance-bot/data/derivatives/raw/binance_usdm/monthly/markPriceKlines/BTCUSDT/5m/BTCUSDT-5m-2021-01.zip
../binance-bot/data/derivatives/raw/binance_usdm/monthly/indexPriceKlines/BTCUSDT/5m/BTCUSDT-5m-2021-01.zip
../binance-bot/data/derivatives/raw/binance_usdm/daily/markPriceKlines/BTCUSDT/5m/BTCUSDT-5m-2026-06-16.zip
```

Preferred normalized CSV layout:

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

The optional premium-index CSVs may be absent only if the execution records a
family-level skip with URL keys, HTTP status or local error, and explicit
`optional_cross_check_missing` status. Missing required mark/index objects must
reject the materialization.

Preferred manifest layout:

```text
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_objects.csv
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_objects.json
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_files.csv
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_files.json
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_summary.json
```

## Normalized Schema

Normalized CSVs should keep only source facts needed for the later source audit:

```text
open_time,open,high,low,close,close_time,source_object_id
```

Rules:

- timestamps are UTC Unix milliseconds in the raw archive semantics;
- rows are sorted by `open_time`;
- duplicate identical rows are counted and collapsed;
- duplicate conflicting rows reject the file;
- non-monotonic physical order is counted in the manifest before sorting;
- the normalized file covers only the approved candle-anchor era;
- no imputation, interpolation, future fill, symbol substitution, or
  candle-volume proxy is allowed.

The later source audit may derive mark-minus-index basis only after source
validation. Materialization itself should not compute context rows, basis
buckets, labels, or cohort outputs.

## Required Provenance Fields

Every archive object manifest row must record:

- object id;
- source owner: `Binance Data Vision`;
- product: `USDT-M futures`;
- archive family;
- source family;
- symbol;
- interval;
- scope: `monthly` or `daily`;
- key: `YYYY-MM` or `YYYY-MM-DD`;
- URL;
- durable raw zip path;
- captured timestamp;
- file bytes;
- SHA-256;
- ETag if returned by the archive response;
- Last-Modified if returned by the archive response;
- Content-Length if returned by the archive response;
- row count in the zip CSV;
- used row count after era clipping;
- out-of-range row count;
- parse error count;
- duplicate same count;
- duplicate conflict count;
- first and last open time;
- validation status.

Every normalized file manifest row must record:

- durable normalized path;
- source family and archive family;
- symbol;
- interval;
- first and last open time;
- first and last close time;
- row count;
- expected row count for the approved era;
- missing interval count;
- duplicate count;
- gap count;
- non-monotonic physical row count;
- source object count;
- normalized file SHA-256;
- timestamp semantics;
- finality rule;
- missing-data policy;
- comparison-only status;
- validation status.

## Finality And Alignment Contract

Materialized rows are not yet approved context inputs. They become candidate
inputs only after the later zero-trade source audit validates them.

For kline-like rows:

- `open_time` identifies the interval start;
- `close_time` identifies when the interval is complete;
- a row may be known only after its close time;
- future alignment must use
  `source_close_time + publication_lag <= decision_candle_close_time`;
- if publication lag cannot be proven, the later source audit must either
  reject the source or use a conservative one-native-interval-lag view;
- nearest-future joins and interpolation from future rows are forbidden.

## Future Execution Stop States

Allowed future execution stop states are:

- `derivatives_context_source_materialization_rejected_public_archive_gap`;
- `derivatives_context_source_materialization_rejected_checksum_or_schema_gap`;
- `derivatives_context_source_materialization_rejected_unapproved_source_path`;
- `derivatives_context_source_materialization_rejected_unapproved_live_or_private_path`;
- `derivatives_context_source_materialization_passed_ready_for_source_audit_approval`;
- `derivatives_context_source_materialization_passed_with_optional_premium_gap_ready_for_source_audit_approval`.

Passing materialization means durable files and manifests exist. It does not
approve the source audit implementation or any context/strategy work.

## Future Execution Verification

If a later user request explicitly approves execution, the materialization
closeout should run at least:

```bash
find ../binance-bot/data/derivatives -type f | sort
wc -l ../binance-bot/data/derivatives/*.csv
wc -l ../binance-bot/data/derivatives/manifests/*.csv
sha256sum ../binance-bot/data/derivatives/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

The future execution may also run source-materialization-specific validation
commands, but it must not run the derivatives context source audit unless that
separate audit implementation is explicitly approved after materialization.

## Docs-Only Verification

This plan requires no Go tests, generated results, source downloads, or audit
runs. It changes only documentation and tracked project memory.

Closeout commands:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
