# Futures Range Universe Structured Compression Optimization Review

Date: 2026-06-26

## Verdict

Stop state:
`structured_compression_optimization_passed_needs_strategy_spec`.

The bounded offline optimization/robustness pass was implemented and run for
the passing `4h` structured-compression universe stream only. The failed `1h`
surface was not optimized.

The selected configuration is:
`sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`.

This is an ETH/SOL-authoritative universe candidate with BTCUSDT emitted as a
diagnostic-only stream. It is not a BTC strategy. It is strong enough to move
to a first offline candidate strategy spec, with the weak BTC diagnostic result
and symbol-specific split fragility carried forward as constraints.

No live, paper, testnet, exchange API, deployment, data download, broad symbol
mining, grid, martingale, averaging down, or two-exchange path is approved.

## Sources And Coverage

Result directory:
`results/futures-range-universe-structured-compression-optimization/`.

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

The optimization source artifact validated all three approved universe sources:

| Symbol | Rows | First Open | Last Open | Physical Non-Monotonic | Sorted | Gaps | Duplicates | Zero Volume | Status |
| --- | ---: | --- | --- | ---: | --- | ---: | ---: | ---: | --- |
| `BTCUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `66` | accepted |
| `ETHUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `47` | accepted |
| `SOLUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `1` | `true` | `0` | `0` | `47` | accepted |

Only closed UTC `4h` resamples were used:

| Symbol | Timeframe | Rows | First Open | Last Open | Missing Children | Status |
| --- | --- | ---: | --- | --- | ---: | --- |
| `BTCUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `ETHUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `SOLUSDT` | `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_range_universe_structured_compression_optimization_coverage.csv` | `4` |
| `futures_range_universe_structured_compression_optimization_grid.csv` | `217` |
| `futures_range_universe_structured_compression_optimization_rankings.csv` | `217` |
| `futures_range_universe_structured_compression_optimization_sources.csv` | `4` |
| `futures_range_universe_structured_compression_optimization_summary.csv` | `9,505` |
| `futures_range_universe_structured_compression_optimization_trades.csv` | `35,881` |
| `summary.csv` | `13` |

## Optimization Scope

The grid was predeclared and bounded:

| Dimension | Values |
| --- | --- |
| Timeframe | `4h` only |
| Family | `structured_compression_expansion` only |
| Detector | `p30_c12_bollinger_on_adx_off` |
| Event delay | first closed breakout within `24` closed `4h` bars |
| Confirmation window | `2`, `3`, `4` bars |
| Max hold | `4`, `6`, `8`, `12` bars |
| Target multiple | `0.75`, `1.0`, `1.25` range widths |
| Stop buffer | `0.0`, `0.10` range widths |
| Symbol sets | `BTC_ETH_SOL`, `ETH_SOL`, `BTC_DIAGNOSTIC_ETH_SOL` |

The grid produced `216` configurations. `115` passed the optimization gate.
The ranking selected the BTC-diagnostic ETH/SOL configuration because it keeps
the BTC diagnostic rows visible while giving authority only to ETHUSDT and
SOLUSDT.

The top-ranked rows were:

| Rank | Config | Symbol Set | Authority | Trades | Net P&L | PF | Max DD | Stress Net | OOS Net | Recent Net |
| ---: | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `1` | `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` | `BTC_DIAGNOSTIC_ETH_SOL` | `ETHUSDT,SOLUSDT` | `129` | `573.87` | `1.8089` | `9.82%` | `151.79` | `229.02` | `193.06` |
| `2` | `sc4h_eth_sol_cw2_h12_t1_00_sb0_00` | `ETH_SOL` | `ETHUSDT,SOLUSDT` | `129` | `573.87` | `1.8089` | `9.82%` | `151.79` | `229.02` | `193.06` |
| `21` | `sc4h_btc_eth_sol_cw2_h12_t1_00_sb0_00` | `BTC_ETH_SOL` | `BTCUSDT,ETHUSDT,SOLUSDT` | `184` | `473.20` | `1.4743` | `13.22%` | `94.27` | `221.84` | `157.09` |

## Selected Result

The selected authority stream uses:

- confirmation window: `2` closed `4h` bars;
- max hold: `12` closed `4h` bars;
- target: `1.0` completed range width from slipped entry;
- stop buffer: `0.0` range width;
- authority symbols: `ETHUSDT`, `SOLUSDT`;
- diagnostic symbol: `BTCUSDT`.

Authority split rows:

| Split | Trades | Gross P&L | Net P&L | PF | Max DD | Avg Net R |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `54` | `174.57` | `151.79` | `1.4867` | `9.82%` | `0.1673` |
| `2023_2024_oos` | `43` | `252.13` | `229.02` | `2.2318` | `5.05%` | `0.4585` |
| `2025_2026_recent` | `32` | `214.34` | `193.06` | `1.9121` | `9.75%` | `0.4983` |
| `full_2021_2026` | `129` | `641.05` | `573.87` | `1.8089` | `9.82%` | `0.3465` |

Authority side rows:

| Side | Trades | Gross P&L | Net P&L | PF | Max DD | Avg Net R |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| long | `69` | `423.19` | `385.79` | `2.0488` | `4.97%` | `0.4760` |
| short | `60` | `217.86` | `188.08` | `1.5506` | `10.89%` | `0.1975` |

Selected full-sample symbol rows:

| Symbol | Authority | Diagnostic | Trades | Gross P&L | Net P&L | PF | Max DD |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| `BTCUSDT` | false | true | `55` | `-71.32` | `-100.67` | `0.6507` | `13.22%` |
| `ETHUSDT` | true | false | `70` | `394.76` | `351.81` | `1.9044` | `5.37%` |
| `SOLUSDT` | true | false | `59` | `246.28` | `222.06` | `1.6930` | `11.25%` |
| `all` | true | false | `129` | `641.05` | `573.87` | `1.8089` | `9.82%` |

Selected symbol split caveats:

- `BTCUSDT` remained negative in every period split and is diagnostic-only.
- `ETHUSDT` was positive full-sample and in stress/OOS, but its recent split
  was negative after costs: `15` trades, net P&L `-17.76`, PF `0.8542`.
- `SOLUSDT` was positive full-sample and in OOS/recent, but its stress split
  was negative after costs: `24` trades, net P&L `-101.84`, PF `0.3720`.

The selected common compatibility outputs contain only authority trades:
`trades.json` has `129` trades and `summary.csv/json` matches the selected
ETH/SOL authority stream.

## Review Gate

The selected configuration passes the optimization gate:

- full authority trades are above `100`;
- `2023_2024_oos` and `2025_2026_recent` authority trades are above `25`;
- full net P&L is positive after costs;
- full PF is above `1.2`;
- stress, OOS, and recent aggregate authority splits are positive after costs;
- long and short authority sides are both positive after costs;
- drawdown improved versus the prior `4h` baseline aggregate;
- ETHUSDT and SOLUSDT each retain adequate full-sample evidence.

The result is still constrained:

- it is not BTCUSDT-positive;
- the selected stream is cross-symbol ETH/SOL, not a single-symbol BTC strategy;
- each transfer symbol has one weak period split;
- the high pass count means the next step must freeze a strategy spec and
  robustness checks rather than widen optimization.

## Next Step

Write a first offline candidate strategy spec for the selected ETH/SOL
structured-compression configuration:

- fixed local ETHUSDT and SOLUSDT Binance USDT-M futures sources;
- BTCUSDT diagnostic-only replay allowed, not authority;
- closed UTC `4h` resampling from `5m` parents;
- detector and event logic frozen to this review;
- confirmation window `2`, max hold `12`, target multiple `1.0`, stop buffer
  `0.0`;
- no additional grid search until a strategy spec defines walk-forward,
  robustness, and rejection rules.

The next step is still offline. It should not add live, paper, testnet,
exchange API, deploy, credential, data-download, broad symbol-mining, grid,
martingale, averaging down, or two-exchange logic.
