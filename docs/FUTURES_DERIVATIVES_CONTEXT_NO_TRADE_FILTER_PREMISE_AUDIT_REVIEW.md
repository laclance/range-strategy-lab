# Futures Derivatives Context No-Trade Filter Premise Audit Review

Date: 2026-06-29

## Verdict

Stop state:
`derivatives_context_no_trade_filter_premise_audit_passed_needs_filter_integration_spec`.

The user explicitly approved implementing the zero-trade derivatives no-trade
filter premise audit selected by
`docs/FUTURES_DERIVATIVES_CONTEXT_STRATEGY_PREMISE_SPEC.md`. The audit used only
the approved BTCUSDT Binance USDT-M futures candle anchor and BTCUSDT
mark/index/premium `5m` derivatives context files, kept the conservative
one-`5m` source lag, and tested whether the five BTCUSDT `15m` `h48`
toxic/no-trade cohorts reproduce and combine into a stable closed-candle veto
candidate.

The premise passed as a filter-integration candidate, not as a strategy. All
five exact toxic rows reproduced, the de-duplicated canonical union remained
toxic/no-trade dominated in the full sample and every period split, overlap and
collateral damage were reported, and common outputs stayed zero-trade. No entry,
exit, fill model, P&L, optimizer, replay, walk-forward, paper/testnet/live path,
exchange API, credential, deployment, source download, or promotion was
performed.

The only authorized next step is a separate docs-only filter integration spec.

## What Was Implemented

- CLI flag `-futures-derivatives-no-trade-filter-premise-audit` with default
  output directory
  `results/futures-derivatives-no-trade-filter-premise-audit/`, wired in
  `cmd/rangelab/main.go` with futures source-product, zero-trade, and
  audit-conflict guards.
- Audit engine
  `internal/lab/futures_derivatives_no_trade_filter_premise_audit.go`:
  `RunFuturesDerivativesNoTradeFilterPremiseAudit`. It validates BTCUSDT
  source facts, rebuilds closed lagged derivatives context, reproduces the five
  selected exact toxic rows, constructs a de-duplicated canonical veto union,
  reports overlap and collateral labels, and writes zero-trade diagnostics.
- Tests in
  `internal/lab/futures_derivatives_no_trade_filter_premise_audit_test.go` and
  `cmd/rangelab/main_test.go` cover de-duplication, boundary stop states,
  no-trade-only definitions, full local-state construction context, artifact
  emission, zero-trade common outputs, and flag-conflict rejection.

## Inputs

Allowed source files:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_mark_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_index_price_klines_5m_BTCUSDT_2021_2026.csv
../binance-bot/data/derivatives/binance_usdm_premium_index_klines_5m_BTCUSDT_2021_2026.csv
```

Source facts reproduced:

| Role | Rows | Gaps | Duplicates | Zero volume | SHA-256 |
| --- | ---: | ---: | ---: | ---: | --- |
| BTCUSDT mark price `5m` | 571,675 | 6 | 0 | 0 | `424c05ca880a31270eea1286d6cdd96ac1132d848e8f5e9d6b3b7177bb7c2858` |
| BTCUSDT index price `5m` | 570,812 | 8 | 0 | 0 | `7ba5a375311e0324dab38f18c2a7137376b619ce63cfb01a17e5684c58390aca` |
| BTCUSDT premium index `5m` | 571,959 | 7 | 0 | 0 | `094e610617f812f032e6a68b3ae6186b20359592415cdb678845c5d287ec298c` |
| BTCUSDT trade candle anchor `5m` | 573,984 | 0 | 0 | 66 | recorded in source manifest |

All four sources span `2021-01-01T00:00:00Z` through
`2026-06-16T23:55:00Z`.

## Anti-Lookahead And Missingness

The audit keeps the source-audit finality rule:
`source_close_time + 5m <= decision_candle_close_time`. No forward fill,
interpolation, nearest-future join, future source revision, or silent default
context is used. Missing lagged context is skipped and counted. Forward labels
are evaluation metadata only and never feature, state, context, or gate inputs.

Required lag coverage reproduced:

| Stream | Lag coverage |
| --- | ---: |
| mark BTCUSDT | 0.995975 |
| index BTCUSDT | 0.994472 |
| premium BTCUSDT | 0.996470 |

The required floor remains the rounded source-audit floor `0.994472`. Derived
BTCUSDT `15m` local-state context recorded `24,067` state rows, `23,756` feature
rows, and `311` missing lagged basis-context rows; those rows were skipped, not
filled. The full local-state construction context (`15m`, `1h`, and `4h`) is
kept internally because the accepted local bucket grammar uses the
higher-timeframe proxy. The audited filter question and generated veto rows
remain BTCUSDT `15m/h48` only.

## Artifacts

Output dir:
`results/futures-derivatives-no-trade-filter-premise-audit/`.

Audit artifacts are written as CSV and JSON:

```text
futures_derivatives_no_trade_filter_premise_sources             4 rows
futures_derivatives_no_trade_filter_premise_coverage            7 rows
futures_derivatives_no_trade_filter_premise_filter_definitions  5 rows
futures_derivatives_no_trade_filter_premise_exact_candidates    20 rows
futures_derivatives_no_trade_filter_premise_canonical_union     4 rows
futures_derivatives_no_trade_filter_premise_overlap             40 rows
futures_derivatives_no_trade_filter_premise_veto_candidates     1,823 rows
futures_derivatives_no_trade_filter_premise_collateral_damage   37 rows
futures_derivatives_no_trade_filter_premise_missingness         4 rows
futures_derivatives_no_trade_filter_premise_summary             4 rows
source_manifest.json, summary.json, summary.csv, trades.json (0 trades)
```

CSV line count total including headers: `1,971`.

## Exact Candidate Reproduction

All five selected exact rows reproduced before the canonical union was
evaluated:

| Candidate | Rows | Weakest split | Full toxic | Worst split toxic | Full toxic improvement | Gate |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| trend down + basis discount + premium discount | 515 | 110 | 0.732039 | 0.800000 | 0.049738 | pass |
| trend down + basis discount | 622 | 142 | 0.729904 | 0.802817 | 0.047603 | pass |
| trend flat + basis discount + basis change flat | 356 | 62 | 0.662921 | 0.699387 | 0.045126 | pass |
| trend up + basis discount | 613 | 124 | 0.654160 | 0.759358 | 0.050840 | pass |
| trend flat + basis discount + premium discount | 538 | 115 | 0.659851 | 0.719212 | 0.042056 | pass |

The diagnostic BTCUSDT `15m/h24` rotation row from the context audit was not
selected, not converted into an entry premise, and not used by the veto.

## Canonical Union

The canonical filter ID is
`btc_15m_basis_discount_no_trade_veto_v1`.

Full sample:

- exact candidate row sum: `2,644`;
- de-duplicated veto rows: `1,823`;
- overlap rows: `821`;
- nested trend-down premium overlap rows: `515`;
- no-trade toxic rows: `1,241`;
- toxic rate: `0.680746`;
- weakest split rows: `387`;
- min split toxic rate: `0.665485`;
- worst split toxic rate: `0.708475`;
- local-only baseline toxic rate: `0.634477`;
- full toxic improvement: `0.046269`;
- rotation-useful collateral rows: `311`;
- continuation-useful collateral rows: `271`.

Split toxic rates:

| Split | Veto rows | Toxic rows | Toxic rate |
| --- | ---: | ---: | ---: |
| `2021_2022_stress` | 387 | 260 | 0.671835 |
| `2023_2024_oos` | 590 | 418 | 0.708475 |
| `2025_2026_recent` | 846 | 563 | 0.665485 |
| `full_2021_2026` | 1,823 | 1,241 | 0.680746 |

The union passed the toxic-dominance gate in full sample and every split. Its
double-counting protection passed because
`sum_exact_candidate_rows - deduplicated_rows = overlap_rows`.

Collateral damage was reported rather than hidden. In the full sample, the
largest blocked non-toxic labels were `clean_expansion_up` (`164`),
`false_break_reentry_up` (`130`), `false_break_reentry_down` (`99`),
`clean_expansion_down` (`94`), and `contained_rotation` (`82`).

## Commands And Outcomes

```bash
gofmt -w internal/lab/futures_derivatives_no_trade_filter_premise_audit.go internal/lab/futures_derivatives_no_trade_filter_premise_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-derivatives-no-trade-filter-premise-audit -out-dir results/futures-derivatives-no-trade-filter-premise-audit
wc -l results/futures-derivatives-no-trade-filter-premise-audit/*.csv
```

Outcomes: all package tests passed; the audit reproduced
`source_rows=4`, `coverage_rows=7`, `filter_definition_rows=5`,
`exact_candidate_rows=20`, `canonical_union_rows=4`, `overlap_rows=40`,
`veto_candidate_rows=1823`, `collateral_rows=37`, `missingness_rows=4`,
`exact_candidates_passed=5`, `canonical_union_passed=true`, `trades=0`, and the
passing stop state above. Common `summary.csv` stayed at `0` trades for every
split and `trades.json` contained no trades.

## Boundaries

This pass authorizes no trading implementation. The veto is a premise for a
later filter-integration spec only. It may not be interpreted as "trade the
opposite", basis tradability, a rotation entry, a continuation entry, P&L
evidence, replay evidence, walk-forward evidence, or promotion.

No closed family is reopened. BTCUSDT price-only range families, BTC regime plus
ETH/SOL context, the router rotation premise, spread-range work, and
volatility-aware exits retain their prior boundaries. ETHUSDT, SOLUSDT, other
timeframes, broad source expansion, source downloads, live/paper/testnet paths,
exchange API work, credentials, deploy files, martingale, averaging down, and
two-exchange logic remain outside scope.
