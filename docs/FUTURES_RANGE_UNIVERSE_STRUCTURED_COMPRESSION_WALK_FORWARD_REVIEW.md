# Futures Range Universe Structured Compression Walk-Forward Review

Date: 2026-06-26

## Verdict

Stop state:
`structured_compression_walk_forward_fragile_needs_review`.

The bounded offline walk-forward robustness pass was implemented and run behind
`-futures-range-universe-structured-compression-walk-forward-robustness`.
It reused the already declared `4h` structured-compression grid and compared
fold-selected configurations against the frozen replay config:
`sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`.

The walk-forward result is not robust enough for a candidate strategy package.
Only one of three folds selected the frozen config and passed the fold gate.
One fold selected a same-shape ETH/SOL authority config that tested worse than
the frozen config, and one fold had no selectable training config because the
combined train period had only `97` authority trades against the `100` trade
gate.

This review does not approve live, paper, testnet, exchange API, deploy,
credential, data download, broad symbol mining, new grid dimensions, BTCUSDT
promotion, martingale, averaging down, or two-exchange work.

## Sources And Coverage

Result directory:
`results/futures-range-universe-structured-compression-walk-forward-robustness/`.

All three approved local Binance USDT-M futures `5m` sources were accepted:

| Symbol | Rows | First Open | Last Open | Physical Non-Monotonic | Sorted | Gaps | Duplicates | Zero Volume | Status |
| --- | ---: | --- | --- | ---: | --- | ---: | ---: | ---: | --- |
| `BTCUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | false | `0` | `0` | `66` | accepted |
| `ETHUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | false | `0` | `0` | `47` | accepted |
| `SOLUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `1` | true | `0` | `0` | `47` | accepted |

Closed UTC `4h` resamples were accepted for every symbol:

| Symbol | Rows | First Open | Last Open | Missing Children | Status |
| --- | ---: | --- | --- | ---: | --- |
| `BTCUSDT` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `ETHUSDT` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |
| `SOLUSDT` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | accepted |

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_range_universe_structured_compression_walk_forward_coverage.csv` | `4` |
| `futures_range_universe_structured_compression_walk_forward_folds.csv` | `4` |
| `futures_range_universe_structured_compression_walk_forward_grid.csv` | `217` |
| `futures_range_universe_structured_compression_walk_forward_rankings.csv` | `649` |
| `futures_range_universe_structured_compression_walk_forward_sources.csv` | `4` |
| `futures_range_universe_structured_compression_walk_forward_summary.csv` | `9,505` |
| `futures_range_universe_structured_compression_walk_forward_trades.csv` | `35,881` |
| `summary.csv` | `13` |

Normal `trades.json` and `summary.csv/json` stayed fixed to the frozen ETH/SOL
authority replay. `trades.json` contains `129` authority trades. BTCUSDT
appears only in walk-forward-specific diagnostic/comparison artifacts.

## Walk-Forward Folds

The grid stayed at the declared `216` configurations:

- confirmation window `2,3,4`;
- max hold `4,6,8,12`;
- target multiple `0.75,1.0,1.25`;
- stop buffer `0.0,0.10`;
- symbol sets `BTC_ETH_SOL`, `ETH_SOL`, and `BTC_DIAGNOSTIC_ETH_SOL`.

Fold results:

| Fold | Train | Test | Selected Config | Selected Test Net | Frozen Test Net | Pass | Reason |
| --- | --- | --- | --- | ---: | ---: | --- | --- |
| `wf_2021_2022_train__2023_2024_test` | `2021_2022_stress` | `2023_2024_oos` | `sc4h_btc_diagnostic_eth_sol_cw2_h12_t0_75_sb0_10` | `92.68` | `229.02` | false | `selected_test_net_worse_than_frozen` |
| `wf_2021_2024_train__2025_2026_test` | `2021_2022_stress+2023_2024_oos` | `2025_2026_recent` | none | `0.00` | `193.06` | false | `no_training_config_selected` |
| `wf_2023_2024_train__2025_2026_test` | `2023_2024_oos` | `2025_2026_recent` | `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` | `193.06` | `193.06` | true |  |

The frozen config was present in every fold:

| Fold | Frozen Rank | Frozen Train Trades | Frozen Train Net | Frozen Train PF | Frozen Test Trades | Frozen Test Net | Frozen Test PF |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `wf_2021_2022_train__2023_2024_test` | `7` | `54` | `151.79` | `1.4867` | `43` | `229.02` | `2.2318` |
| `wf_2021_2024_train__2025_2026_test` | `1` | `97` | `380.81` | `1.7650` | `32` | `193.06` | `1.9121` |
| `wf_2023_2024_train__2025_2026_test` | `1` | `43` | `229.02` | `2.2318` | `32` | `193.06` | `1.9121` |

## Review Gate

The source and resample gate passed, and BTCUSDT was never required as
authority for a passing fold.

The robustness gate failed:

- fold 1 selected a BTC-diagnostic ETH/SOL authority row with the same
  confirmation window `2` and max hold `12`, but its test net P&L was
  materially worse than frozen: `92.68` versus `229.02`;
- fold 2 had no selectable training config because the train aggregate had
  `97` authority trades, below the required `100` for multi-split training;
- fold 3 selected the exact frozen config and passed;
- only one of three folds selected the exact frozen config;
- the result would require changing strategy rules or relaxing gates to pass,
  so it must stop for review.

ETH/SOL transfer evidence remains visible, but fragile. The frozen stream is
positive after costs in every period split, yet the walk-forward selection test
does not reproduce enough robustness to justify packaging it as a strategy.
BTCUSDT remains diagnostic-only and cannot offset this fragility.

## Next Step

Do not tune around this walk-forward failure. The structured-compression
ETH/SOL authority stream should remain review-only unless the user explicitly
approves a materially new premise or a separate review-only decision task.

The next safe action is a post-fragility hypothesis pivot brief: keep the
walk-forward result as exclusion evidence, avoid reopening the failed `1h`
surface or widening this grid, and define a materially different offline
range-strategy premise before any new implementation.
