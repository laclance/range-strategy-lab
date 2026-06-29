# Futures Derivatives Context Zero-Trade Source Audit Review

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`.

The user explicitly approved implementing the zero-trade derivatives context
source audit described in
`docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_SOURCE_AUDIT_BRIEF.md`. The audit
validated the nine durable materialized Binance USDT-M futures mark/index/
premium `5m` source files, bound their provenance by SHA-256, and proved they
can be anti-lookahead-aligned to the existing `5m` candle anchors with bounded,
recorded missingness. All six required mark/index streams aligned above the
`0.99` coverage bar, so the source family is good enough to justify a separate
later zero-trade context-audit brief.

This is a source and alignment verdict only. It does not create context
features, labels, cohorts, rankings, entries, exits, P&L, replay, walk-forward,
or any promotion. The next derivatives step (a zero-trade context-audit brief,
then implementation) still requires separate explicit user approval.

## What Was Implemented

- CLI flag `-futures-derivatives-context-source-audit` (default out-dir
  `results/futures-derivatives-context-source-audit/`), wired in
  `cmd/rangelab/main.go` with the same source-product and flag-conflict guards
  as the other zero-trade audits.
- Audit engine `internal/lab/futures_derivatives_context_source_audit.go`:
  `RunFuturesDerivativesContextSourceAudit`. It is read-only over the durable
  source files plus the candle anchors and writes only diagnostic artifacts.
- Tests: `internal/lab/futures_derivatives_context_source_audit_test.go`
  (pass path plus source-gap, timestamp/finality, alignment-gap, closed-family
  rescue, and live/private-path rejections) and a CLI flag test in
  `cmd/rangelab/main_test.go` (artifact emission, default no-emission, zero
  trades, conflict rejection).

## Inputs

Candidate source files (durable, materialized, SHA-256 bound; hashes match
`docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_MATERIALIZATION_REVIEW.md`):

```text
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
```

Candle alignment anchors only (not derivatives context sources):

```text
../binance-bot/data/{btcusdt,ethusdt,solusdt}_futures_um_5m_2021_2026.csv
```

Each anchor reproduced `573,984` candles spanning `2021-01-01T00:00:00Z`
through `2026-06-16T23:55:00Z` with `gap_count=0` and `duplicate_count=0`.

## Anti-Lookahead Alignment Contract

Publication lag for mark/index/premium klines is not formally proven, so the
audit uses the conservative one-native-interval (`5m`) lag view: decision candle
`D` is served only by the source interval at `D - 5m`, i.e.
`source_close_time + 5m <= decision_candle_close_time`. No forward fill (max
extra staleness `0`), no interpolation, no nearest-future joins, and no
labels/future returns/future revisions are used as validation inputs. Every
`timestamp_alignment` row records `uses_future_rows=false` and
`exact_closed_interval_join=true`. The materialization gaps surface here as
bounded missingness with `forward_filled_rows=0`, never silent fills.

## Coverage And Missingness Outcomes

All nine streams validated with `0` duplicate-conflict, `0` close-time
violations, `0` non-monotonic rows, and `0` parse errors. Lag coverage =
lag-aligned decision candles / `573,984` (one warmup boundary candle at the era
start is unavoidable and reported as a skip).

| Stream | Required | Lag-aligned | Lag-missing | Lag coverage | Missing intervals |
| --- | --- | ---: | ---: | ---: | ---: |
| mark BTCUSDT | yes | 571,674 | 2,310 | 0.995976 | 2,309 |
| index BTCUSDT | yes | 570,811 | 3,173 | 0.994472 | 3,172 |
| mark ETHUSDT | yes | 573,401 | 583 | 0.998984 | 582 |
| index ETHUSDT | yes | 573,115 | 869 | 0.998486 | 868 |
| mark SOLUSDT | yes | 571,962 | 2,022 | 0.996477 | 2,021 |
| index SOLUSDT | yes | 573,115 | 869 | 0.998486 | 868 |
| premium BTCUSDT | no | 571,958 | 2,026 | 0.996470 | 2,025 |
| premium ETHUSDT | no | 571,959 | 2,025 | 0.996472 | 2,024 |
| premium SOLUSDT | no | 572,247 | 1,737 | 0.996973 | 1,736 |

Required-stream minimum lag coverage is `0.994472` (index BTCUSDT), above the
`0.99` bar, so `aligned_required_streams=6` of `6`. Required-stream validation
status is `accepted_with_recorded_gaps`; optional premium streams are
`accepted_optional_cross_check_with_recorded_gaps`.

## Artifacts

Output dir `results/futures-derivatives-context-source-audit/` (Git-ignored).
Nine artifact families, each CSV + JSON, exactly matching the brief's allowed
set, plus the zero-trade common outputs:

```text
futures_derivatives_context_source_audit_sources.{csv,json}            (9 rows)
futures_derivatives_context_source_audit_candle_anchors.{csv,json}     (3 rows)
futures_derivatives_context_source_audit_external_coverage.{csv,json}  (9 rows)
futures_derivatives_context_source_audit_timestamp_alignment.{csv,json}(9 rows)
futures_derivatives_context_source_audit_publication_lag.{csv,json}    (9 rows)
futures_derivatives_context_source_audit_missingness.{csv,json}        (9 rows)
futures_derivatives_context_source_audit_provenance.{csv,json}         (9 rows)
futures_derivatives_context_source_audit_skips.{csv,json}             (18 rows)
futures_derivatives_context_source_audit_summary.{csv,json}           (10 rows)
source_manifest.json, summary.json, summary.csv, trades.json (no trades)
```

No basis/feature/cohort/ranking/label artifact was written; mark-minus-index
basis derivation is deferred to a later context-audit stage. The `provenance`
artifact records source owner `Binance Data Vision`, archive family, durable
path, recomputed file SHA-256 (matching the materialization manifests), first
`source_object_id`, and the materialization manifest paths.

## Commands And Outcomes

```bash
gofmt -l internal/lab/futures_derivatives_context_source_audit*.go cmd/rangelab/main*.go   # clean
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...                # ok
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -futures-derivatives-context-source-audit -out-dir results/futures-derivatives-context-source-audit
wc -l results/futures-derivatives-context-source-audit/*.csv
git diff --check    # clean
git status --short
```

Outcomes: all package tests passed; the audit reproduced `source_rows=9`,
`anchor_rows=3`, `coverage_rows=9`, `alignment_rows=9`, `lag_rows=9`,
`missingness_rows=9`, `provenance_rows=9`, `skip_rows=18`, `required_streams=6`,
`aligned_required_streams=6`, and stop state
`derivatives_context_zero_trade_source_audit_passed_needs_context_audit_brief`;
`trades.json` contained no trades; recomputed provenance SHA-256 matched the
materialization manifests; `git diff --check` passed. (The brief's
`rg --files ../binance-bot/data | rg basis` check returns empty because the
adjacent `binance-bot` repo Git-ignores its `data/` tree; `find` confirms all
nine durable source CSVs are present and the audit read them successfully.)

## Boundaries

A passing source audit does not authorize context-gain implementation, labels,
cohorts, rankings, entries, exits, P&L backtests, optimizer grids, replay,
walk-forward, packaging, source downloads, paper/testnet/live paths, exchange
API work, credentials, deploy files, strategy promotion, martingale, averaging
down, or two-exchange logic. The next step is a separate, approval-gated
zero-trade derivatives context-audit brief.
