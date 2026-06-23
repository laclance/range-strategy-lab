# Futures Range Universe Structured Compression Baseline Review

Date: 2026-06-26

## Verdict

Stop state:
`structured_compression_baseline_passed_needs_optimization_brief`.

The fixed-rule structured-compression universe baseline was implemented and
run across the approved local Binance USDT-M futures `5m` sources for
`BTCUSDT`, `ETHUSDT`, and `SOLUSDT`.

The `4h structured_compression_expansion all h6` surface passed as a
cross-symbol aggregate after costs. The result is not clean enough to treat as
strategy approval: BTCUSDT was negative, and the full stream depends on ETH and
SOL strength. It is strong enough to justify a bounded offline optimization
and robustness brief.

The `1h structured_compression_expansion all h12` surface failed after costs
and should not be optimized from this result.

Default `cmd/rangelab` still uses `lab.EmptyStrategy`. This backtest runs only
behind the explicit
`-futures-range-universe-structured-compression-baseline-backtest` flag.

## Sources And Coverage

Result directory:
`results/futures-range-universe-structured-compression-baseline-backtest/`.

The normal `source_manifest.json` records the default BTCUSDT futures parent
source:

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

The baseline source artifact validated all three approved universe sources:

| Symbol | Rows | First Open | Last Open | Physical Non-Monotonic | Sorted | Gaps | Duplicates | Zero Volume | Status |
| --- | ---: | --- | --- | ---: | --- | ---: | ---: | ---: | --- |
| `BTCUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `66` | accepted |
| `ETHUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `47` | accepted |
| `SOLUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `1` | `true` | `0` | `0` | `47` | accepted |

Only closed UTC `1h` and `4h` resamples were used:

| Symbol | Timeframe | Rows | First Open | Last Open | Missing Children | Status |
| --- | --- | ---: | --- | --- | ---: | --- |
| `BTCUSDT` | `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `0` | accepted |
| `BTCUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `ETHUSDT` | `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `0` | accepted |
| `ETHUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `SOLUSDT` | `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `0` | accepted |
| `SOLUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_range_universe_structured_compression_baseline_coverage.csv` | `7` |
| `futures_range_universe_structured_compression_baseline_signals.csv` | `926` |
| `futures_range_universe_structured_compression_baseline_sources.csv` | `4` |
| `futures_range_universe_structured_compression_baseline_summary.csv` | `97` |
| `futures_range_universe_structured_compression_baseline_trades.csv` | `913` |
| `summary.csv` | `13` |

## Baseline Template

Both selected surfaces used the same fixed non-optimized template:

- detector: percentile `0.30`, min consecutive bars `12`, Bollinger on, ADX
  off;
- completed mature range episodes only;
- first closed breakout candle within `24` higher-timeframe bars after episode
  end;
- closed confirmation candle within the next `3` higher-timeframe bars;
- up confirmation: close above range high and high above breakout high;
- down confirmation: close below range low and low below breakout low;
- signal on the closed confirmation candle;
- entry on the next higher-timeframe bar open;
- stop at the broken range boundary;
- target one completed range width from slipped entry;
- existing engine fees, slippage, risk sizing, one-position max, and
  stop-first ambiguity.

The evaluated surfaces were:

| Candidate | Timeframe | Side | Max Hold | Signals | Skipped | Trades |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| `structured_compression_4h_all_h6` | `4h` | long + short | `6` | `186` | `1` | `185` |
| `structured_compression_1h_all_h12` | `1h` | long + short | `12` | `739` | `5` | `727` |

Exit reasons across both candidates:

| Exit Reason | Trades |
| --- | ---: |
| `stop_loss` | `510` |
| `take_profit` | `163` |
| `time_stop` | `239` |

## Results

Full-sample all-side rows:

| Candidate | Symbol | Trades | Gross P&L | Net P&L | Costs | PF | Max DD | Avg Net R |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `structured_compression_4h_all_h6` | `BTCUSDT` | `56` | `-62.55` | `-92.27` | `38.34` | `0.6404` | `9.39%` | `-0.3725` |
| `structured_compression_4h_all_h6` | `ETHUSDT` | `70` | `283.09` | `243.61` | `50.93` | `1.7811` | `5.56%` | `0.2345` |
| `structured_compression_4h_all_h6` | `SOLUSDT` | `59` | `175.14` | `151.58` | `30.39` | `1.5543` | `7.95%` | `0.2412` |
| `structured_compression_4h_all_h6` | `all` | `185` | `395.68` | `302.92` | `119.66` | `1.3598` | `12.92%` | `0.0529` |
| `structured_compression_1h_all_h12` | `BTCUSDT` | `255` | `-26.09` | `-184.28` | `204.06` | `0.7902` | `19.84%` | `-0.2929` |
| `structured_compression_1h_all_h12` | `ETHUSDT` | `235` | `236.00` | `76.44` | `205.84` | `1.0678` | `15.24%` | `-0.1016` |
| `structured_compression_1h_all_h12` | `SOLUSDT` | `237` | `29.11` | `-92.29` | `156.60` | `0.9183` | `25.93%` | `-0.1076` |
| `structured_compression_1h_all_h12` | `all` | `727` | `239.01` | `-200.13` | `566.50` | `0.9362` | `35.46%` | `-0.1706` |

Period splits for the passing `4h` aggregate:

| Split | Trades | Gross P&L | Net P&L | PF | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `77` | `27.31` | `-4.82` | `0.9878` | `12.73%` |
| `2023_2024_oos` | `63` | `213.20` | `180.30` | `1.7867` | `6.70%` |
| `2025_2026_recent` | `45` | `155.16` | `127.44` | `1.5896` | `10.19%` |
| `full_2021_2026` | `185` | `395.68` | `302.92` | `1.3598` | `12.92%` |

Side split for the passing `4h` aggregate:

| Side | Trades | Gross P&L | Net P&L | PF | Max DD | Avg Net R |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| long | `98` | `246.01` | `195.96` | `1.4762` | `7.11%` | `0.1221` |
| short | `87` | `149.67` | `106.96` | `1.2485` | `15.31%` | `-0.0250` |

The common combined compatibility view across both candidates had `912` trades,
gross P&L `634.69`, net P&L `102.79`, and PF `1.0258`. It is not the promotion
authority because the `1h` candidate failed.

## Review Gate

The `4h` surface clears the baseline gate as an aggregate universe stream:

- full net P&L is positive after costs;
- full PF is above `1.2`;
- `2023_2024_oos` and `2025_2026_recent` are positive after costs;
- both long and short aggregate sides are positive after costs;
- drawdown is not obviously incompatible with a bounded optimization pass.

Important weaknesses:

- BTCUSDT is negative full-sample and has only `56` full trades, below the
  `100` full-sample adequacy preference for a standalone symbol;
- the aggregate `2021_2022_stress` split is slightly negative after costs;
- the `4h` aggregate result is carried by ETHUSDT and SOLUSDT;
- the `1h` surface fails and should not be optimized.

This means the next step should not be broad strategy promotion. It should be a
bounded offline optimization and robustness pass for the `4h` structured
compression universe stream only, with explicit treatment of BTC weakness,
symbol inclusion, side behavior, and stress-split fragility.

## Next Step

Create a bounded offline optimization brief for:

- `structured_compression_4h_all_h6` only;
- approved local `BTCUSDT`, `ETHUSDT`, and `SOLUSDT` Binance USDT-M futures
  sources only;
- fixed source validation and closed UTC `4h` resampling;
- no live, paper, testnet, exchange API, data download, symbol expansion, grid,
  martingale, averaging down, or two-exchange logic.

Do not optimize the `1h` structured-compression surface from this result.
