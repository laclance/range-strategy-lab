# Futures Clean Breakout Baseline Review

Date: 2026-06-26

## Verdict

Stop state:
`clean_breakout_baseline_failed_no_promotion`.

The fixed-rule clean breakout continuation baseline was implemented and run on
Binance USDT-M futures BTCUSDT data. Both approved candidates failed after
costs. No optimization, portfolio-stream combination, 15m expansion, live
wiring, paper/testnet path, or strategy replacement is approved from this
baseline.

The user-approved portfolio-stream routing rule was evaluated. It is not
triggered here: neither candidate is near-viable after costs, and no single
candidate is strong enough to optimize.

Default `cmd/rangelab` still uses `lab.EmptyStrategy` unless an explicit
offline prototype or backtest flag is passed.

## Source And Coverage

Parent source:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Venue / product | `Binance` / `Binance USDT-M futures` |
| Symbol / interval | `BTCUSDT` / `5m` |
| Row count | `573,984` loaded candles |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps / duplicates | `0` / `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Baseline resamples:

| Timeframe | Rows | First Open | Last Open |
| --- | ---: | --- | --- |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` |

Result directory:
`results/futures-clean-breakout-baseline-backtest/`.

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_clean_breakout_baseline_signals.csv` | `4,849` |
| `futures_clean_breakout_baseline_trades.csv` | `1,186` |
| `futures_clean_breakout_baseline_summary.csv` | `25` |
| `summary.csv` | `13` |

The run emitted `4,848` signal rows and `1,185` executed trade rows. Common
`summary.*` and `trades.json` are compatibility views across the independent
candidate runs; the strategy-specific summary is the review authority.

## Baseline Template

Both candidates used the same fixed, non-optimized template:

- completed mature range episodes from the discovery detector backbone:
  percentile `0.30`, min consecutive bars `12`, Bollinger on, ADX off;
- signal only on a closed higher-timeframe breakout candle after a completed
  mature range;
- entry on the next higher-timeframe bar open;
- long stop at the broken range high and target one completed range width
  above slipped entry;
- short stop at the broken range low and target one completed range width
  below slipped entry;
- max hold `12` higher-timeframe bars;
- existing engine fees, slippage, 1% risk-at-stop sizing, 1x notional cap,
  one-position max, and stop-first ambiguity.

The two candidates were evaluated independently:

| Candidate | Timeframe | Side | Signals | Skipped | Trades |
| --- | --- | --- | ---: | ---: | ---: |
| `clean_breakout_4h_up_h12` | `4h` | long only | `496` | `0` | `121` |
| `clean_breakout_1h_all_h12` | `1h` | long + short | `4,352` | `47` | `1,064` |

## Results

Full-sample candidate results:

| Candidate | Trades | Win Rate | Gross P&L | Net P&L | Costs | PF | Gross PF | Max DD | Avg Net R |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `clean_breakout_4h_up_h12` | `121` | `28.93%` | `2.55` | `-69.29` | `92.67` | `0.8386` | `1.0069` | `13.31%` | `-0.7057` |
| `clean_breakout_1h_all_h12` | `1,064` | `27.44%` | `272.63` | `-320.25` | `764.82` | `0.8846` | `1.1177` | `37.91%` | `-0.3479` |

Period-split all-side results:

| Candidate | Split | Trades | Gross P&L | Net P&L | PF |
| --- | --- | ---: | ---: | ---: | ---: |
| `clean_breakout_4h_up_h12` | `2021_2022_stress` | `44` | `32.57` | `7.19` | `1.0405` |
| `clean_breakout_4h_up_h12` | `2023_2024_oos` | `44` | `-22.94` | `-49.44` | `0.7198` |
| `clean_breakout_4h_up_h12` | `2025_2026_recent` | `33` | `-7.08` | `-27.04` | `0.6416` |
| `clean_breakout_1h_all_h12` | `2021_2022_stress` | `373` | `153.10` | `-75.06` | `0.9415` |
| `clean_breakout_1h_all_h12` | `2023_2024_oos` | `392` | `56.61` | `-164.38` | `0.8167` |
| `clean_breakout_1h_all_h12` | `2025_2026_recent` | `299` | `62.93` | `-80.81` | `0.8644` |

The `1h` candidate was negative on both sides:

| Side | Trades | Gross P&L | Net P&L | PF |
| --- | ---: | ---: | ---: | ---: |
| long | `532` | `199.66` | `-99.25` | `0.9252` |
| short | `532` | `72.97` | `-221.00` | `0.8474` |

The aggregate compatibility summary was also negative:

| Trades | Gross P&L | Net P&L | Costs | PF |
| ---: | ---: | ---: | ---: | ---: |
| `1,185` | `275.18` | `-389.54` | `857.49` | `0.8784` |

## Review Gate

The baseline fails the promotion gate:

- full net P&L is negative for both candidates;
- full profit factor is below `1.2` for both candidates;
- `2023_2024_oos` and `2025_2026_recent` are negative after costs for both
  candidates;
- the `1h` all-side candidate is negative for both long and short;
- executed trade counts are adequate, so the failure is behavioral, not just
  a sample-size gap.

The mixed portfolio-stream route is not triggered. The routing rule requires
two or more viable or near-viable candidates with no single excellent
candidate. Here, both candidates are below near-flat thresholds after costs.

## Next Step

Do not optimize this clean breakout baseline. Do not combine these two
candidates into a portfolio-style stream from this result. Do not broaden the
same failed template to `15m` automatically.

The next task should ask for a new futures premise or an explicit scope choice
before adding another backtest. A later `15m` clean-breakout comparison or a
portfolio-style stream can still be proposed, but only from a new user-approved
premise, not as a continuation of this failed baseline.
