# Futures Derivatives Context Zero-Trade Audit Review

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_zero_trade_context_audit_passed_needs_strategy_premise_spec`.

The user explicitly approved implementing the zero-trade derivatives context
audit described in
`docs/FUTURES_DERIVATIVES_CONTEXT_ZERO_TRADE_CONTEXT_AUDIT_BRIEF.md`. The audit
used only the nine validated Binance USDT-M futures mark/index/premium `5m`
source CSVs plus the three Binance USDT-M futures candle anchors, computed
lagged mark-minus-index basis and optional premium context from closed source
rows, and tested whether those buckets separate local range states beyond the
local price/volume state alone.

The result is a narrow context pass, not a strategy pass: `6` BTCUSDT `15m`
derivatives-context cohorts passed the declared separation gates (`5`
no-trade/toxic and `1` rotation candidate). ETHUSDT and SOLUSDT produced no
passing cohorts. The run remained zero-trade throughout: no entries, exits, P&L,
optimizer, replay, walk-forward, source download, API, deployment, or promotion
was performed.

The only authorized next step is a separate strategy-premise spec that decides
whether and how to test the passing context cohorts. This audit does not by
itself authorize a strategy implementation.

## What Was Implemented

- CLI flag `-futures-derivatives-context-audit` with default output directory
  `results/futures-derivatives-context-audit/`, wired in `cmd/rangelab/main.go`
  with futures source-product and zero-trade audit flag-conflict guards.
- Audit engine `internal/lab/futures_derivatives_context_audit.go`:
  `RunFuturesDerivativesContextAudit`. It reuses the passed source audit's
  source/provenance/alignment contract, derives closed lagged basis/premium
  buckets, joins them to local range states, labels forward outcomes only as
  metadata, and ranks separation cohorts versus the local-only baseline.
- Tests: `internal/lab/futures_derivatives_context_audit_test.go`
  (lagged closed source-row use, derivatives-improvement requirement,
  future-label rejection, closed-family rescue rejection) and a CLI flag test in
  `cmd/rangelab/main_test.go` (artifact emission, default no-emission,
  zero-trade common output, conflict rejection).

## Inputs

Derivatives context source files (durable, materialized, SHA-256 bound; hashes
match the materialization and source-audit records):

```text
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_{BTCUSDT,ETHUSDT,SOLUSDT}_2021_2026.csv
```

Candle anchors:

```text
../binance-bot/data/{btcusdt,ethusdt,solusdt}_futures_um_5m_2021_2026.csv
```

Each candle anchor reproduced `573,984` candles spanning
`2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`, with `gap_count=0`,
`duplicate_count=0`, and zero-volume counts BTC `66`, ETH `47`, SOL `47`.

## Anti-Lookahead And Missingness

The audit keeps the source-audit finality rule: a decision candle may use only a
source interval satisfying
`source_close_time + 5m <= decision_candle_close_time`. With the repo's
second-precision candle close timestamps, a decision close such as
`00:14:59Z` maps to the lagged source open `00:05:00Z`, not to the future
`00:10:00Z` row.

No forward fill, interpolation, nearest-future join, or silent default bucket is
used. Missing lagged basis context produces skipped local-state rows and
missingness artifacts. Forward labels appear only in label/cohort/ranking/
summary artifacts; `ForwardLabelUsedAsFeature=false` throughout.

The fail/pass coverage gate is the carried-forward source-audit lag coverage
floor. All required mark/index streams remain at or above the rounded
`0.994472` floor under the one-`5m` lag. Local range-state subsets have their
own diagnostic basis-context coverage, with minimum `0.982930` for BTCUSDT `4h`;
those missing rows are recorded and skipped rather than filled or treated as
context values.

## Artifacts

Output dir `results/futures-derivatives-context-audit/` (Git-ignored). Audit
artifacts are written as CSV and JSON:

```text
futures_derivatives_context_sources          12 rows
futures_derivatives_context_coverage         18 rows
futures_derivatives_context_basis_features   83,004 rows
futures_derivatives_context_local_states     83,640 rows
futures_derivatives_context_labels           249,012 rows
futures_derivatives_context_cohorts          512,190 rows
futures_derivatives_context_rankings         181,827 rows
futures_derivatives_context_missingness      36 rows
futures_derivatives_context_summary          109 rows
source_manifest.json, summary.json, summary.csv, trades.json (0 trades)
```

The audit-specific summary row reports `source_rows=12`, `coverage_rows=18`,
`basis_feature_rows=83004`, `local_state_rows=83640`, `label_rows=249012`,
`cohort_rows=512190`, `ranking_rows=181827`, `passing_cohorts=6`,
`missingness_rows=36`, `common_outputs_zero_trade=true`,
`forward_labels_as_inputs=false`, `orthogonality_gate_applied=true`, and
`trades=0`.

## Passing Context Cohorts

All passing rows are BTCUSDT `15m` and all pass the orthogonality/improvement
gates versus the local-only baseline:

| Rank | Route | Horizon | Local state | Derivatives bucket | Rows | Weakest split | Main separation |
| ---: | --- | ---: | --- | --- | ---: | ---: | --- |
| 1 | no-trade/toxic | 48 | `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none` | `basis_discount_small + premium_discount_small` | 515 | 110 | full toxic `0.732039`, worst split toxic `0.800000`, toxic improvement `0.049738` full / `0.109942` split |
| 2 | no-trade/toxic | 48 | `geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none` | `basis_discount_small` | 622 | 142 | full toxic `0.729904`, worst split toxic `0.802817`, toxic improvement `0.047603` full / `0.112758` split |
| 3 | no-trade/toxic | 48 | `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` | `basis_discount_small + basis_change_flat` | 356 | 62 | full toxic `0.662921`, worst split toxic `0.699387`, toxic improvement `0.045126` full / `0.058361` split |
| 4 | no-trade/toxic | 48 | `geometry_midline_balanced::vol_compressed::trend_up_pressure::impulse_none` | `basis_discount_small` | 613 | 124 | full toxic `0.654160`, worst split toxic `0.759358`, toxic improvement `0.050840` full / `0.137293` split |
| 5 | no-trade/toxic | 48 | `geometry_midline_balanced::vol_compressed::trend_flat::impulse_none` | `basis_discount_small + premium_discount_small` | 538 | 115 | full toxic `0.659851`, worst split toxic `0.719212`, toxic improvement `0.042056` full / `0.078186` split |
| 6 | rotation candidate | 24 | `geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale` | `basis_discount_small` | 313 | 71 | useful `0.632588` full / `0.521127` weakest split, margin improvement `0.072540` full / `0.025844` split |

Rows ranked after these failed the review gates. The most common high-rank
failure was inadequate cohort count, not source leakage or closed-family rescue.

## Commands And Outcomes

```bash
gofmt -w internal/lab/futures_derivatives_context_audit.go internal/lab/futures_derivatives_context_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-context-audit -out-dir results/futures-derivatives-context-audit
wc -l results/futures-derivatives-context-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

Outcomes: all package tests passed; the approved audit reproduced the passing
stop state and the row counts above; common `summary.csv` stayed at `0` trades
for every split and `trades.json` contained no trades; generated CSV line count
total was `1,109,870` including headers; `git diff --check` passed.

## Boundaries

This passing context audit does not authorize entries, exits, P&L backtests,
optimizer grids, replay, walk-forward, packaging, source downloads, source
expansion, symbol expansion, paper/testnet/live paths, exchange API work,
credentials, deploy files, strategy promotion, martingale, averaging down, or
two-exchange logic.

No closed family is reopened. The BTCUSDT price-only range families, BTC regime
plus ETH/SOL context family, router rotation premise, spread-range strategy
work, and volatility-aware exits retain their prior boundaries. A later
strategy-premise spec must decide whether these six BTCUSDT `15m` context
cohorts justify a no-trade filter premise, a rotation premise, both as separate
tracks, or no further strategy work.
