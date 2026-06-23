# Futures Range Universe Structured Compression Strategy Replay Review

Date: 2026-06-26

## Verdict

Stop state:
`structured_compression_strategy_replay_passed_needs_walk_forward_robustness_brief`.

The fixed offline replay for
`sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` was implemented and run.
It reproduces the selected optimization authority stream closely enough to
stay promoted to the next offline robustness step.

Authority remains ETHUSDT and SOLUSDT only. BTCUSDT remains diagnostic-only
and is still negative. This review does not approve live, paper, testnet,
exchange API, deploy, credential, data download, broad symbol mining, grid,
martingale, averaging down, or two-exchange work.

## Sources And Coverage

Result directory:
`results/futures-range-universe-structured-compression-strategy-replay/`.

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

The strategy source artifact validated all three approved universe sources:

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
| `futures_range_universe_structured_compression_strategy_coverage.csv` | `4` |
| `futures_range_universe_structured_compression_strategy_signals.csv` | `186` |
| `futures_range_universe_structured_compression_strategy_sources.csv` | `4` |
| `futures_range_universe_structured_compression_strategy_summary.csv` | `49` |
| `futures_range_universe_structured_compression_strategy_trades.csv` | `185` |
| `summary.csv` | `13` |

## Replay Scope

Frozen replay config:

| Field | Value |
| --- | --- |
| Config | `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` |
| Symbol set | `BTC_DIAGNOSTIC_ETH_SOL` |
| Authority | `ETHUSDT,SOLUSDT` |
| Diagnostic | `BTCUSDT` |
| Candidate | `structured_compression_4h_all_h12` |
| Timeframe | closed UTC `4h` |
| Detector | `p30_c12_bollinger_on_adx_off` |
| Event delay | first closed breakout within `24` closed `4h` bars |
| Confirmation window | `2` closed `4h` bars |
| Target | `1.0` completed range width |
| Stop buffer | `0.0` range width |
| Max hold | `12` closed `4h` bars |

The replay emitted `185` signal rows and `184` strategy-specific trade rows.
Common `trades.json` and common `summary.csv/json` include only the `129`
ETH/SOL authority trades.

## Authority Result

Common authority split rows:

| Split | Trades | Gross P&L | Net P&L | PF | Max DD | Avg Hold |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `54` | `174.57` | `151.79` | `1.4867` | `9.82%` | `5.44` |
| `2023_2024_oos` | `43` | `252.13` | `229.02` | `2.2318` | `5.05%` | `5.86` |
| `2025_2026_recent` | `32` | `214.34` | `193.06` | `1.9121` | `9.75%` | `3.59` |
| `full_2021_2026` | `129` | `641.05` | `573.87` | `1.8089` | `9.82%` | `5.12` |

Authority side rows:

| Side | Trades | Gross P&L | Net P&L | PF | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: |
| long | `69` | `423.19` | `385.79` | `2.0488` | `4.97%` |
| short | `60` | `217.86` | `188.08` | `1.5506` | `10.89%` |

Selected full-sample symbol rows:

| Symbol | Authority | Diagnostic | Trades | Gross P&L | Net P&L | PF | Max DD |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| `BTCUSDT` | false | true | `55` | `-71.32` | `-100.67` | `0.6507` | `13.22%` |
| `ETHUSDT` | true | false | `70` | `394.76` | `351.81` | `1.9044` | `5.37%` |
| `SOLUSDT` | true | false | `59` | `246.28` | `222.06` | `1.6930` | `11.25%` |
| `all` | true | false | `129` | `641.05` | `573.87` | `1.8089` | `9.82%` |

The replay preserves the known caveats:

- BTCUSDT is negative and diagnostic-only.
- ETHUSDT is positive full-sample but negative in `2025_2026_recent`.
- SOLUSDT is positive full-sample but negative in `2021_2022_stress`.
- The authority short side is positive full-sample, but weaker than the long
  side and has the worst drawdown row.

## Review Gate

The replay passes the review gate:

- all source and `4h` resample rows are accepted;
- common outputs include ETH/SOL authority trades only;
- BTCUSDT appears only as diagnostic rows;
- full authority trades are above `100`;
- `2023_2024_oos` and `2025_2026_recent` authority trades are above `25`;
- full authority net P&L is positive after costs;
- full authority PF is above `1.2`;
- stress, OOS, and recent authority splits are positive after costs;
- long and short authority sides are both positive after costs;
- ETHUSDT and SOLUSDT each remain positive full-sample with PF above `1.0`;
- replay metrics match the selected optimization review within tolerance.

## Next Step

Run a bounded offline walk-forward robustness pass before any further
strategy-promotion language. The next pass may reuse the already declared
structured-compression `4h` grid only to test forward-selection robustness and
compare against the frozen ETH/SOL authority replay. It must not add new grid
dimensions, new symbols, live/paper/testnet wiring, exchange API use,
credentials, deployment, martingale, averaging down, or two-exchange logic.
