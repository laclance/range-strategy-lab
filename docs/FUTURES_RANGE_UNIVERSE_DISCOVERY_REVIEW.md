# Futures Range Universe Discovery Review

Date: 2026-06-26

## Verdict

Stop state:
`range_universe_audit_ready_for_baseline_backtest`.

The local BTC/ETH/SOL Binance USDT-M futures range-universe discovery audit
found transfer-confirmed range surfaces worth a fixed-rule offline baseline
backtest. The strongest non-duplicative next candidates are:

1. `4h structured_compression_expansion all h6`.
2. `1h structured_compression_expansion all h12`.

This is not optimization approval and not live, paper, testnet, portfolio, or
strategy promotion. It is permission to build the first fixed-rule baseline
backtest for those two structured-compression surfaces only.

Default `cmd/rangelab` still uses `lab.EmptyStrategy`, and this audit produced
zero trades.

## Sources And Coverage

Result directory:
`results/futures-range-universe-discovery-audit/`.

All three approved local sources validated as Binance USDT-M futures `5m`
closed candles:

| Symbol | Rows | First Open | Last Open | Physical Non-Monotonic | Sorted | Gaps | Duplicates | Zero Volume | Status |
| --- | ---: | --- | --- | ---: | --- | ---: | ---: | ---: | --- |
| `BTCUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `66` | accepted |
| `ETHUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `false` | `0` | `0` | `47` | accepted |
| `SOLUSDT` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `1` | `true` | `0` | `0` | `47` | accepted |

The SOL file caveat was handled explicitly: the audit reported one physical
non-monotonic row, sorted by open time, then accepted only after the sorted
stream had no gaps or duplicates.

Every symbol produced complete closed UTC resamples:

| Timeframe | Rows Per Symbol | First Open | Last Open | Gaps | Duplicates | Missing Children | Status |
| --- | ---: | --- | --- | ---: | ---: | ---: | --- |
| `5m` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `0` | `0` | accepted |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | `0` | `0` | `0` | accepted |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `0` | `0` | `0` | accepted |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | `0` | `0` | accepted |

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_range_universe_candidates.csv` | `415,954` |
| `futures_range_universe_coverage.csv` | `13` |
| `futures_range_universe_rankings.csv` | `148` |
| `futures_range_universe_sources.csv` | `4` |
| `futures_range_universe_stability.csv` | `442` |
| `futures_range_universe_summary.csv` | `1,765` |
| `summary.csv` | `13` |

The audit emitted `415,953` candidate rows:

| Symbol | Candidate Rows |
| --- | ---: |
| `BTCUSDT` | `134,379` |
| `ETHUSDT` | `138,324` |
| `SOLUSDT` | `143,250` |

Common `summary.*` and `trades.json` remained zero-trade.

## Discovery Results

The universe audit ranked `147` surfaces. `35` passed the universe gate:

| Family | Passing Rows |
| --- | ---: |
| `structured_compression_expansion` | `21` |
| `breakout_retest_acceptance` | `14` |

Passing rows by timeframe:

| Timeframe | Passing Rows |
| --- | ---: |
| `15m` | `18` |
| `1h` | `14` |
| `4h` | `3` |

Passing rows by side:

| Side | Passing Rows |
| --- | ---: |
| `all` | `15` |
| `up` | `11` |
| `down` | `9` |

No `boundary_touch_rejection`, `single_candle_wick_rejection`,
`failed_breakout_reentry`, or `mature_balance_persistence` surface passed the
universe gate.

Top-ranked rows:

| Rank | Surface | BTC Weak Count | Transfer Symbols | BTC Weak Fav. | BTC Worst Adv. | BTC Worst Quick Invalid. | BTC Weak Cost Buffer | Best Transfer |
| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `1` | `4h structured_compression_expansion all h6` | `100` | `2` | `85.50%` | `10.00%` | `14.00%` | `4.4982%` | `SOLUSDT` |
| `2` | `4h structured_compression_expansion all h12` | `100` | `2` | `85.50%` | `13.00%` | `14.00%` | `4.4546%` | `SOLUSDT` |
| `3` | `4h structured_compression_expansion all h3` | `100` | `2` | `83.50%` | `7.50%` | `14.00%` | `4.2620%` | `SOLUSDT` |
| `4` | `1h structured_compression_expansion all h12` | `799` | `2` | `93.24%` | `6.76%` | `11.76%` | `3.1030%` | `SOLUSDT` |

The top three rows are the same `4h` all-side structured-compression surface at
different horizons. To avoid duplicating the same hypothesis in the first
baseline, use the top `4h h6` row and the strongest non-duplicative `1h h12`
row.

## Selected Baseline Surfaces

Primary candidate: `4h structured_compression_expansion all h6`.

| Symbol | Split | Count | Favorable | Adverse | Quick Invalid. | Cost Buffer |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `BTCUSDT` | `2021_2022_stress` | `224` | `92.41%` | `7.14%` | `12.95%` | `5.3879%` |
| `BTCUSDT` | `2023_2024_oos` | `200` | `85.50%` | `10.00%` | `13.50%` | `4.4982%` |
| `BTCUSDT` | `2025_2026_recent` | `100` | `94.00%` | `6.00%` | `14.00%` | `5.1342%` |
| `ETHUSDT` | weakest split | `132` | `91.67%` | `8.33%` | `10.00%` | `6.2767%` |
| `SOLUSDT` | weakest split | `205` | `91.71%` | `8.29%` | `15.89%` | `9.1636%` |

Secondary candidate: `1h structured_compression_expansion all h12`.

| Symbol | Split | Count | Favorable | Adverse | Quick Invalid. | Cost Buffer |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `BTCUSDT` | `2021_2022_stress` | `1,092` | `95.24%` | `4.76%` | `8.24%` | `4.0239%` |
| `BTCUSDT` | `2023_2024_oos` | `928` | `95.04%` | `4.96%` | `10.88%` | `3.1030%` |
| `BTCUSDT` | `2025_2026_recent` | `799` | `93.24%` | `6.76%` | `11.76%` | `3.1323%` |
| `ETHUSDT` | weakest split | `613` | `91.52%` | `8.48%` | `13.15%` | `3.1778%` |
| `SOLUSDT` | weakest split | `673` | `93.94%` | `6.06%` | `12.24%` | `4.5270%` |

Both selected surfaces have BTC split eligibility, both transfer symbols
passing, favorable rate above adverse rate in every key split, quick
invalidation below favorable rate, and cost buffers above the rough round-trip
fee plus slippage proxy.

## Secondary Evidence

`breakout_retest_acceptance` passed in `14` rows but ranked below structured
compression and carried much higher quick invalidation, especially in `15m`
all-side rows. It should stay as later comparison evidence, not the first
baseline target.

The failed immediate clean-breakout baseline remains closed. The selected
structured-compression surfaces differ because the candidate event requires a
post-break confirmation candle with extension after a completed mature range,
not entry on the first clean breakout candle.

## Next Step

Build a fixed-rule offline baseline backtest for:

- `4h structured_compression_expansion all h6`;
- `1h structured_compression_expansion all h12`.

The baseline should evaluate BTCUSDT, ETHUSDT, and SOLUSDT independently and
in aggregate, use the accepted local futures sources and closed UTC resamples,
enter only on the next higher-timeframe bar after the closed confirmation
candle, and keep the template fixed before result inspection.

Do not optimize, add live-adjacent wiring, download data, broaden symbols, or
convert the discovery labels directly into a strategy without the baseline
review.
