# Futures Derivatives Context Source Materialization Review

Date: 2026-06-28

## Verdict

Stop state:
`derivatives_context_source_materialization_passed_ready_for_source_audit_approval`.

The user explicitly approved executing the offline materialization in
`docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_PLAN.md`. Execution
downloaded, checksum-bound, normalized, and stored durable Binance public Data
Vision USDT-M futures `markPriceKlines`, `indexPriceKlines`, and optional
`premiumIndexKlines` `5m` archive objects for `BTCUSDT`, `ETHUSDT`, and
`SOLUSDT` under `../binance-bot/data/derivatives/`.

This is source materialization and provenance only. It does not run the
derivatives source audit, context features, labels, cohorts, rankings, entries,
exits, P&L, replay, walk-forward, or any strategy/promotion path. The next
derivatives step still requires a separate explicit approval for the zero-trade
source-audit implementation.

## Approved Scope Executed

- Source owner: `Binance Data Vision`. Product: `USDT-M futures`.
- Archive families: `markPriceKlines` and `indexPriceKlines` (required);
  `premiumIndexKlines` (optional cross-check).
- Symbols: `BTCUSDT`, `ETHUSDT`, `SOLUSDT`. Interval: `5m`.
- Era: `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`
  (`era_start_ms=1609459200000`, `era_end_ms=1781654100000`,
  `expected_rows_per_stream=573984`).
- Object set: monthly `2021-01`..`2026-05` (65) plus daily `2026-06-01`..
  `2026-06-16` (16) = `81` objects per family/symbol, `729` total.
- Object URL template (deterministic, public archive only):
  `https://data.binance.vision/data/futures/um/{monthly|daily}/{archive_family}/{symbol}/5m/{symbol}-5m-{key}.zip`
- Each `.zip` verified against its published `.zip.CHECKSUM` SHA-256 before use.

No REST tail, WebSocket, signed/private, funding, aggregate-trade, open-interest,
long/short, liquidation, order-book, spot, or cross-exchange source was touched.

## Measured Object And Byte Outcomes

- `objects_ok=729`, `objects_missing=0`, `objects_error=0`. All `729` object
  manifest rows have `validation_status=accepted` (checksum-verified).
- Total compressed bytes downloaded: `125,508,895`.
- Per archive family (`81` objects per symbol, `243` per family):
  - `markPriceKlines`: `243` objects, `42,833,309` bytes.
  - `indexPriceKlines`: `243` objects, `47,233,885` bytes.
  - `premiumIndexKlines`: `243` objects, `35,441,701` bytes.
- The required `markPriceKlines` + `indexPriceKlines` total of `486` objects and
  `90,067,194` bytes reproduces the adjacent planning inventory exactly. Optional
  premium was unknown in the plan and measured here at `35,441,701` bytes.

## Normalized Stream Outcomes

Schema (source facts only):
`open_time,open,high,low,close,close_time,source_object_id`. Rows are sorted by
`open_time`, duplicate-identical rows collapsed, no imputation/interpolation/
future-fill/symbol substitution. All `9` streams span the full era
(`first_open=1609459200000`, `last_open=1781654100000`; no front or tail
truncation).

| Stream | Rows | Expected | Missing | Gaps | DupSame | DupConflict | NonMono | Status |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| mark BTCUSDT | 571,675 | 573,984 | 2,309 | 6 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| mark ETHUSDT | 573,402 | 573,984 | 582 | 4 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| mark SOLUSDT | 571,963 | 573,984 | 2,021 | 5 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| index BTCUSDT | 570,812 | 573,984 | 3,172 | 8 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| index ETHUSDT | 573,116 | 573,984 | 868 | 4 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| index SOLUSDT | 573,116 | 573,984 | 868 | 4 | 0 | 0 | 0 | accepted_with_recorded_gaps |
| premium BTCUSDT | 571,959 | 573,984 | 2,025 | 7 | 0 | 0 | 0 | accepted_optional_cross_check_with_recorded_gaps |
| premium ETHUSDT | 571,960 | 573,984 | 2,024 | 6 | 0 | 0 | 0 | accepted_optional_cross_check_with_recorded_gaps |
| premium SOLUSDT | 572,248 | 573,984 | 1,736 | 5 | 0 | 0 | 0 | accepted_optional_cross_check_with_recorded_gaps |

Required mark/index missing intervals total `9,820` of `3,443,904`
(`0.285%`). Across all streams there are `0` duplicate-conflicting rows, `0`
parse errors, `0` out-of-range rows, and `0` physical non-monotonic rows.

## Gap Decision

The required mark/index `5m` public archives are **not** gap-contiguous over the
approved era. Missing intervals are real published-archive holes, not parse
artifacts: they are whole-day aligned (multiples of `288` `5m` intervals) and
recur on the same calendar windows across symbols and families, indicating
Binance mark/index publication outages. Observed windows include
`2021-06-30→07-02`, `2021-07-23→07-28` (4 days), several `2022-07` windows,
`2022-10-01→10-03`, `2023-02-23→02-25`, `2023-04` windows, and a brief
`2023-11-10 ~04:00Z` incident. The trade-candle anchors had `gap_count=0` over
the same era, so these derivatives series are genuinely holed relative to the
decision candles.

The materialization plan's normalization rules said "fail closed on gaps", but
the same plan defines `gap_count`, `missing_interval_count`, and a no-imputation
`missing_data_policy`, and the planned later source audit is specified to handle
"bounded missingness/staleness". Because that downstream stage owns the
missingness tolerance decision, the user chose to **record gaps and pass**
rather than fail closed: write all `9` normalized CSVs unfilled, record gaps in
the manifests, and leave bounded-missingness gating to the source audit.
Integrity faults remain fail-closed (a duplicate-conflicting row, schema
ambiguity, checksum mismatch, or any missing required object would have rejected
materialization); none occurred.

## Durable Paths

Target root: `../binance-bot/data/derivatives/` (outside this repo, not tracked
by Git). `743` files total: `729` raw zips, `9` normalized CSVs, `5` manifests.

Normalized CSVs:

```text
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
```

Raw archive zips:

```text
../binance-bot/data/derivatives/raw/binance_usdm/{monthly|daily}/{archive_family}/{symbol}/5m/{symbol}-5m-{key}.zip
```

Manifests:

```text
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_objects.csv
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_objects.json
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_files.csv
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_files.json
../binance-bot/data/derivatives/manifests/binance_usdm_derivatives_context_source_materialization_summary.json
```

The objects manifest records, per object: object id, source owner, product,
archive family, source family, symbol, interval, scope, key, URL, durable raw
path, captured timestamp, file bytes, SHA-256, ETag, Last-Modified,
Content-Length, rows_total, rows_used, rows_out_of_range, parse_errors,
duplicate_same, duplicate_conflict, first_open, last_open, and validation
status. The files manifest records, per normalized stream: durable path, source
and archive family, symbol, interval, first/last open, first/last close, row
count, expected row count, missing interval count, duplicate count, gap count,
non-monotonic physical count, source object count, SHA-256, timestamp
semantics, finality rule, missing-data policy, comparison-only flag, and
validation status.

Normalized CSV SHA-256:

```text
424c05ca880a31270eea1286d6cdd96ac1132d848e8f5e9d6b3b7177bb7c2858  mark_price_klines_5m_BTCUSDT
53b1af03c69e610b35332c3f8d06324fd8fff75dbd9211f27526c7c0b1b5604f  mark_price_klines_5m_ETHUSDT
30a2659ab9e36bad5203d4f8bf4606848b502e22ac1a577a8dfe5928bad31361  mark_price_klines_5m_SOLUSDT
7ba5a375311e0324dab38f18c2a7137376b619ce63cfb01a17e5684c58390aca  index_price_klines_5m_BTCUSDT
26fcefa1d66a1e3e61eb9a31cbfdd75f87d9d2b740221480f4d857c6b5478a92  index_price_klines_5m_ETHUSDT
e378e229bebb5206ef4dfa15282e2b43ecc59150623dadafe77d8519248f0c7e  index_price_klines_5m_SOLUSDT
094e610617f812f032e6a68b3ae6186b20359592415cdb678845c5d287ec298c  premium_index_klines_5m_BTCUSDT
5e39c4a00935b3c925037d32554aa9f661a16a0e95017ae0bd22593d4ad35c2c  premium_index_klines_5m_ETHUSDT
df2450cb3ba4a324a4e4a373037824e5c57663747c692090dac24754232301c1  premium_index_klines_5m_SOLUSDT
```

## Finality And Alignment Contract (carried forward)

Materialized rows are not yet approved context inputs. They become candidate
inputs only after the later zero-trade source audit validates them. `open_time`
identifies the interval start, `close_time` identifies completion, and any
future alignment must use
`source_close_time + publication_lag <= decision_candle_close_time`. If
publication lag cannot be proven, the source audit must reject the source or use
a conservative one-native-interval-lag view. Nearest-future joins and
interpolation from future rows remain forbidden.

## Commands And Outcomes

Materialization was performed by a one-shot offline Go generator (run from the
session scratchpad, not tracked in this repo) that downloaded each public
object, verified it against the published `.CHECKSUM`, wrote the raw zip,
normalized the CSV, and wrote the manifests.

Verification:

```bash
find ../binance-bot/data/derivatives -type f | sort        # 743 files
find ../binance-bot/data/derivatives/raw -name '*.zip' | wc -l   # 729
wc -l ../binance-bot/data/derivatives/*.csv                 # 5,150,260 total (rows + 9 headers)
wc -l ../binance-bot/data/derivatives/manifests/*.csv       # files.csv 10, objects.csv 730
sha256sum ../binance-bot/data/derivatives/*.csv             # 9 hashes above
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check                                            # clean
git status --short                                          # empty (data is outside the repo)
```

Outcomes: all `729` objects downloaded with `0` missing and `0` errors and all
checksums matched; the `9` normalized streams cover the full era with `0`
duplicate-conflict, `0` parse-error, `0` out-of-range, and `0` non-monotonic
rows and recorded gaps as above; the reference scan found only canonical
`memory/NEXT_CODEX_BRIEF.md` mentions; `git diff --check` passed; pre-commit
`git status --short` showed only intended `docs/` and `memory/` changes
(generated derivatives data lives outside this repository and is not committed).

## Next Step

The next derivatives step is a separate, approval-gated zero-trade derivatives
context source-audit implementation. It must not start without explicit user
approval and must not add context features, labels, cohorts, rankings, entries,
exits, P&L, replay, walk-forward, or promotion. Passing materialization here
does not authorize that audit by itself.
