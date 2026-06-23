# Futures Range-Universe Breakout Retest Acceptance Baseline Review

Date: 2026-06-26

## Scope

This review covers the bounded offline baseline behind:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-breakout-retest-acceptance-baseline-backtest -out-dir results/futures-range-universe-breakout-retest-acceptance-baseline-backtest
```

The run recomputes range-universe discovery rows from the accepted local
BTCUSDT, ETHUSDT, and SOLUSDT Binance USDT-M futures sources, then selects the
top passing all-side `breakout_retest_acceptance` rows without using ignored
`results/` files as runtime input.

No structured-compression replay, optimization, or walk-forward result was
rerun or retuned.

## Source And Resample Validation

All required sources passed the existing Binance USDT-M futures source guard:

| Symbol | 5m rows | First open | Last open | Gaps | Duplicates | Zero volume | Physical non-monotonic | Status |
| --- | ---: | --- | --- | ---: | ---: | ---: | ---: | --- |
| BTCUSDT | 573,984 | 2021-01-01T00:00:00Z | 2026-06-16T23:55:00Z | 0 | 0 | 66 | 0 | accepted |
| ETHUSDT | 573,984 | 2021-01-01T00:00:00Z | 2026-06-16T23:55:00Z | 0 | 0 | 47 | 0 | accepted |
| SOLUSDT | 573,984 | 2021-01-01T00:00:00Z | 2026-06-16T23:55:00Z | 0 | 0 | 47 | 1 | accepted after validation sort |

Selected closed UTC resamples also passed:

| Symbol | Timeframe | Rows | First open | Last open | Complete | Status |
| --- | --- | ---: | --- | --- | --- | --- |
| BTCUSDT | 15m | 191,328 | 2021-01-01T00:00:00Z | 2026-06-16T23:45:00Z | true | accepted |
| BTCUSDT | 1h | 47,832 | 2021-01-01T00:00:00Z | 2026-06-16T23:00:00Z | true | accepted |
| ETHUSDT | 15m | 191,328 | 2021-01-01T00:00:00Z | 2026-06-16T23:45:00Z | true | accepted |
| ETHUSDT | 1h | 47,832 | 2021-01-01T00:00:00Z | 2026-06-16T23:00:00Z | true | accepted |
| SOLUSDT | 15m | 191,328 | 2021-01-01T00:00:00Z | 2026-06-16T23:45:00Z | true | accepted |
| SOLUSDT | 1h | 47,832 | 2021-01-01T00:00:00Z | 2026-06-16T23:00:00Z | true | accepted |

Every selected coverage row had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Candidate Selection

The runtime selection used passing discovery ranking rows where
`family=breakout_retest_acceptance`, `side=all`, then de-duplicated by
`(family, timeframe, side)`.

| Selected | Candidate | Discovery rank | Timeframe | Side | Horizon | Rank score |
| ---: | --- | ---: | --- | --- | ---: | ---: |
| 1 | `breakout_retest_acceptance_15m_all_h12` | 22 | 15m | all | 12 | 5.778311 |
| 2 | `breakout_retest_acceptance_1h_all_h12` | 28 | 1h | all | 12 | 5.239488 |

Directional breakout-retest rows and shorter horizon duplicates were not
selected for this default baseline.

## Baseline Rules

For each selected candidate, signals fire on the closed retest/acceptance
candle emitted by the existing `rangeUniverseBreakoutRetestEvents` logic.

Execution is unchanged from the offline engine contract:

- enter on the next higher-timeframe bar open;
- one open position max;
- long stop at the broken range high;
- short stop at the broken range low;
- target one completed range width from slipped entry;
- max hold is the selected horizon, `12` bars;
- stop-first ambiguity handling;
- existing fees, slippage, risk sizing, and 1x notional cap.

Invalid geometry, missing entry candles, non-positive range width, and
non-positive trade prices are skipped and recorded in signal artifacts.

## Results

The run wrote:

- result directory:
  `results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/`
- specific CSV line counts including headers:
  - coverage: `7`
  - selection: `3`
  - signals: `10,827`
  - sources: `4`
  - summary: `97`
  - trades: `7,183`
  - common summary: `13`
- common `trades.json` contains `7,182` trades.

Aggregate outcomes:

| Candidate | Split | Signals | Skipped | Trades | Net P&L | PF | Max DD | Pass |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `breakout_retest_acceptance_15m_all_h12` | full_2021_2026 | 8,701 | 2,849 | 5,825 | -2,329.18 | 0.6778 | 2.2464 | false |
| `breakout_retest_acceptance_15m_all_h12` | 2021_2022_stress | 3,314 | 1,113 | 2,189 | -1,256.08 | 0.7310 | 1.2433 | false |
| `breakout_retest_acceptance_15m_all_h12` | 2023_2024_oos | 3,173 | 985 | 2,180 | -860.87 | 0.5274 | 0.8677 | false |
| `breakout_retest_acceptance_15m_all_h12` | 2025_2026_recent | 2,214 | 751 | 1,456 | -212.23 | 0.7128 | 0.2195 | false |
| `breakout_retest_acceptance_1h_all_h12` | full_2021_2026 | 2,125 | 751 | 1,357 | -604.03 | 0.8652 | 0.7372 | false |
| `breakout_retest_acceptance_1h_all_h12` | 2021_2022_stress | 835 | 313 | 508 | -98.70 | 0.9540 | 0.3145 | false |
| `breakout_retest_acceptance_1h_all_h12` | 2023_2024_oos | 793 | 266 | 525 | -285.08 | 0.8066 | 0.3260 | false |
| `breakout_retest_acceptance_1h_all_h12` | 2025_2026_recent | 497 | 172 | 324 | -220.25 | 0.7445 | 0.3291 | false |

Full-period symbol transfer outcomes:

| Candidate | BTCUSDT net | ETHUSDT net | SOLUSDT net | Transfer read |
| --- | ---: | ---: | ---: | --- |
| `breakout_retest_acceptance_15m_all_h12` | -828.60 | -639.03 | -861.54 | failed on every symbol |
| `breakout_retest_acceptance_1h_all_h12` | -361.79 | -5.18 | -237.06 | failed on every symbol |

Both candidates had enough trades, but neither produced positive full-period
net P&L after costs, neither reached PF `1.2`, both lost in the `2023_2024_oos`
and `2025_2026_recent` splits, and neither showed symbol transfer strength.

## Verdict

Stop state:

```text
breakout_retest_acceptance_baseline_failed_no_promotion
```

The breakout-retest/acceptance premise is not promoted. This result does not
authorize optimization, robustness review, parameter retuning, BTCUSDT
promotion, a portfolio stream, paper/testnet/live wiring, exchange API use,
deployment, data downloads, broad symbol mining, martingale, averaging down,
or two-exchange work.

The reviewed baseline is exclusion evidence for this fixed-rule
`breakout_retest_acceptance` implementation. A new implementation brief should
be review-only unless the user chooses a materially different offline
range-strategy premise.

## Verification

Passed:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-breakout-retest-acceptance-baseline-backtest -out-dir results/futures-range-universe-breakout-retest-acceptance-baseline-backtest
wc -l results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/*.csv
```
