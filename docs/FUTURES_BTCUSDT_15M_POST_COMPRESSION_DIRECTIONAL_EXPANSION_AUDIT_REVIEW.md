# Futures BTCUSDT 15m Post-Compression Directional Expansion Audit Review

Date: 2026-06-30

## Verdict

Stop state:
`btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review`.

The user explicitly approved implementing the zero-trade audit for:

```text
btc_15m_post_compression_directional_expansion_v1
```

The audit used only the accepted local BTCUSDT Binance USDT-M futures `5m`
source, resampled exact closed UTC `15m` candles from complete `5m` children,
and evaluated the predeclared `81`-cell post-compression directional expansion
grid. It produced diagnostic labels only and kept common outputs at
`trades=0`.

The premise passed as zero-trade entry-premise evidence, not as a strategy. The
passing evidence is narrow: long-side `48`-bar forward labels only, at the
`192`-bar compression lookback and bottom `20%` compression threshold, with
adjacent passing cells across all predeclared breakout thresholds and volume
modes. No short-side, `16`-bar, or `32`-bar side/horizon passed the full gate.

The only authorized next step is a separate docs-only strategy-premise spec that
decides whether this zero-trade evidence is strong enough to request a later
offline backtest spec. This review does not authorize entries, exits, P&L,
optimizer selection, replay, walk-forward, derivatives veto interaction, or
strategy promotion.

## What Was Implemented

- CLI flag
  `-futures-btc-15m-post-compression-directional-expansion-audit` with default
  output directory
  `results/futures-btc-15m-post-compression-directional-expansion-audit/`,
  wired in `cmd/rangelab/main.go` with futures source-product, zero-trade, and
  audit-conflict guards.
- Audit engine
  `internal/lab/futures_btc_15m_post_compression_directional_expansion_audit.go`:
  `RunFuturesBTC15MPostCompressionDirectionalExpansionAudit`. It validates the
  accepted BTCUSDT source manifest, resamples closed `15m` candles, builds the
  predeclared grid, labels candidates from `open[d+1]`, compares every
  cell/side/horizon against the unconditional eligible `15m` baseline, and
  applies the de-duplicated candidate, split-stability, and adjacent-cell gates.
- Tests in
  `internal/lab/futures_btc_15m_post_compression_directional_expansion_audit_test.go`
  and `cmd/rangelab/main_test.go` cover prior-window exclusion, exact prior
  valid percentile windows, ATR `d-1`, volume-threshold exclusion, side
  assignment, forward-label anchoring and symmetry, de-duplication, baseline
  comparison, adjacency, stop-state precedence, artifact emission, zero-trade
  common outputs, spot rejection, and flag-conflict rejection.

## Inputs

Allowed source file:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

Source facts reproduced:

| Field | Value |
| --- | ---: |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |

Exact `15m` resample facts:

| Field | Value |
| --- | ---: |
| Row count | `191,328` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:45:00Z` |
| First close | `2021-01-01T00:14:59Z` |
| Last close | `2026-06-16T23:59:59Z` |
| Expected child bars | `3` |
| Complete bucket count | `191,328` |
| Partial final child bars | `0` |
| Missing child opens | `0` |
| Validation status | `accepted` |

## Audit Definition

The grid is the predeclared family from the premise spec:

| Dimension | Values |
| --- | --- |
| Compression lookback | `48`, `96`, `192` prior closed `15m` bars |
| Compression threshold | bottom `20%`, `30%`, `40%` of prior range width |
| Percentile reference | prior `1,920` valid closed `15m` range-width observations |
| Breakout threshold | `0.1`, `0.2`, `0.3` prior-bar `ATR(14)` beyond prior range |
| Volume confirmation | `none`, `above_prior_96_median`, `above_prior_96_p60` |

No-lookahead rules enforced:

- prior range uses `[d-L, d-1]`;
- percentile references exclude `d` and require exact prior valid observations;
- ATR uses `ATR(14)[d-1]`;
- volume thresholds use `[d-96, d-1]`;
- breakout and decision volume are known only at decision candle close;
- labels are diagnostics from `open[d+1]` through horizons `16`, `32`, and `48`;
- candidate viability uses de-duplicated `(decision_close, side)` events, not
  repeated grid hits.

## Artifacts

Output dir:
`results/futures-btc-15m-post-compression-directional-expansion-audit/`.

Audit artifacts were written as CSV and JSON sidecars:

```text
btc_15m_post_compression_directional_expansion_sources             1 row
btc_15m_post_compression_directional_expansion_resample_coverage   1 row
btc_15m_post_compression_directional_expansion_parameter_cells      81 rows
btc_15m_post_compression_directional_expansion_candidates           386,694 rows
btc_15m_post_compression_directional_expansion_dedup_events         4,677 rows
btc_15m_post_compression_directional_expansion_baseline             24 rows
btc_15m_post_compression_directional_expansion_split_summary        1,944 rows
btc_15m_post_compression_directional_expansion_adjacency            486 rows
btc_15m_post_compression_directional_expansion_missingness          4 rows
btc_15m_post_compression_directional_expansion_falsification.json
source_manifest.json, summary.json, summary.csv, trades.json (0 trades)
```

CSV line count total including headers: `393,934`.

## Falsification Result

The falsification report passed every declared gate:

| Gate | Result |
| --- | --- |
| Source and resample | pass |
| Leakage protection | pass |
| Full de-duplicated candidate size | `4,677` events, pass versus required `300` |
| Minimum primary split size | `1,484` events, pass versus required `50` |
| Baseline separation | pass |
| Adjacent-cell cluster | pass |
| Split stability | pass |
| Closed-family protection | pass |
| Derivatives-veto contamination | pass |
| Common outputs zero-trade | pass, `trades=0` |

Missingness was skipped and counted, not filled:

| Reason | Count | Rate |
| --- | ---: | ---: |
| audit warmup | `2,112` | `0.011039` |
| missing forward label | `27` | `0.000141` |
| missing max-horizon future | `48` | `0.000251` |
| missing volume reference | `780` | `0.004077` |

## Passing Pocket

All `9` adjacent passing cell/side/horizon rows share the same surface:

- side: `long`;
- horizon: `48` closed `15m` bars;
- lookback: `192` prior closed `15m` bars;
- compression threshold: bottom `20%`;
- breakout thresholds: `0.1`, `0.2`, and `0.3` prior-bar `ATR(14)`;
- volume modes: `none`, `above_prior_96_median`, and `above_prior_96_p60`.

Every one of those cells passed the same-side/same-horizon baseline gate in the
full sample and in all primary splits. The cells are adjacent by the declared
one-ordered-step rule, so the pass is not an isolated parameter-cell artifact.

Representative full-sample rows:

| Cell | Rows | Mean close-return delta bp | FMA delta bp | FGTA delta | Adjacent passing cells |
| --- | ---: | ---: | ---: | ---: | ---: |
| `l192_q20_m010_none` | `514` | `10.250396` | `36.813284` | `0.042029` | `2` |
| `l192_q20_m010_above_prior_96_median` | `514` | `10.250396` | `36.813284` | `0.042029` | `3` |
| `l192_q20_m010_above_prior_96_p60` | `511` | `10.343778` | `37.013819` | `0.043225` | `2` |
| `l192_q20_m020_none` | `468` | `12.986691` | `39.661224` | `0.049936` | `3` |
| `l192_q20_m020_above_prior_96_median` | `468` | `12.986691` | `39.661224` | `0.049936` | `4` |
| `l192_q20_m020_above_prior_96_p60` | `466` | `13.371315` | `40.193430` | `0.052275` | `3` |
| `l192_q20_m030_none` | `410` | `14.352689` | `42.978777` | `0.051406` | `2` |
| `l192_q20_m030_above_prior_96_median` | `410` | `14.352689` | `42.978777` | `0.051406` | `3` |
| `l192_q20_m030_above_prior_96_p60` | `410` | `14.352689` | `42.978777` | `0.051406` | `2` |

The weakest confirming split is `2025_2026_recent`, where the passing cells
still clear all three metrics but with much smaller margins. For example,
`l192_q20_m010_none` at `h48` has `162` rows, a close-return delta of
`3.093987` bp, favorable-minus-adverse delta of `9.262758` bp, and
favorable-greater-than-adverse delta of `0.003822`.

## What Did Not Pass

- No short-side cell passed the full gate.
- No `16`-bar or `32`-bar horizon passed the full gate.
- No `48`-bar or `96`-bar compression lookback produced adjacent passing
  evidence.
- No `30%` or `40%` compression threshold produced adjacent passing evidence.
- The result is not evidence for the derivatives no-trade veto, basis
  tradability, a basis-fade rule, a router-rotation rescue, or a reviewed
  compression-breakout family.

## Commands And Outcomes

```bash
gofmt -w cmd/rangelab/main.go cmd/rangelab/main_test.go internal/lab/futures_btc_15m_post_compression_directional_expansion_audit.go internal/lab/futures_btc_15m_post_compression_directional_expansion_audit_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-15m-post-compression-directional-expansion-audit
wc -l results/futures-btc-15m-post-compression-directional-expansion-audit/*.csv
```

Outcomes: all package tests passed; the approved audit reproduced
`source_rows=1`, `coverage_rows=1`, `parameter_cells=81`,
`candidate_rows=386694`, `dedup_events=4677`, `baseline_rows=24`,
`split_summary_rows=1944`, `adjacency_rows=486`, `missingness_rows=4`,
`passing_cells=9`, `adjacent_pass_clusters=9`, `trades=0`, and the passing stop
state above. Common `summary.csv` stayed at `0` trades for every split and
`trades.json` contained no trades.

## Boundaries

This pass authorizes no trading implementation. It is a zero-trade diagnostic
entry-premise audit only. It does not prove P&L, cost tolerance, executable
fill quality, stop/target behavior, replay robustness, walk-forward durability,
or promotion readiness.

No closed family is reopened. Compression breakout, structured compression,
clean breakout continuation, router rotation, occupancy rotation,
midline/hold-inside, higher-timeframe nested rotation, BTC/ETH/SOL context, and
the derivatives veto retain their prior boundaries. The canonical derivatives
veto `btc_15m_basis_discount_no_trade_veto_v1` remains parked as future
skip/retain evidence only and may not be applied until a separate approved
interaction audit exists.

The next allowed decision is a docs-only strategy-premise spec. That spec must
either define a tightly bounded later offline backtest request from the passing
long `48`-bar pocket, or stop the line if the zero-trade evidence is too narrow
to justify strategy construction.
