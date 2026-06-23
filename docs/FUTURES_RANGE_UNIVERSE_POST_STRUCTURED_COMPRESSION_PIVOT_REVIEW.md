# Futures Range Universe Post Structured Compression Pivot Review

Date: 2026-06-26

## Verdict

Stop state:
`post_structured_compression_pivot_ready_for_breakout_retest_acceptance_baseline`.

This is a review-only pivot pass. No strategy code, optimizer, replay, data
download, live/paper/testnet path, exchange API, credential, deploy file,
martingale, averaging down, or two-exchange logic was added.

The structured-compression walk-forward result is now exclusion evidence for
candidate packaging. The fixed ETH/SOL authority replay was positive, but the
walk-forward selection test did not reproduce enough robustness to package the
stream or tune around it.

The only still-open range-first premise visible from the current docs is the
secondary `breakout_retest_acceptance` family from the range-universe discovery
review. It is materially different from the failed structured-compression path
because the next test would require a breakout retest/acceptance event after a
completed mature range, selected from existing discovery evidence, rather than
another compression-expansion confirmation, target, stop, hold, symbol-set, or
training-gate retune.

## Compression Exclusion Evidence

Walk-forward result directory:
`results/futures-range-universe-structured-compression-walk-forward-robustness/`.

The already-generated walk-forward artifacts are present with the documented
CSV line counts:

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

Source and closed UTC `4h` resample validation passed for `BTCUSDT`,
`ETHUSDT`, and `SOLUSDT`: each local Binance USDT-M futures `5m` source had
`573,984` loaded candles from `2021-01-01T00:00:00Z` through
`2026-06-16T23:55:00Z`, and each accepted `4h` resample had `11,958` rows
through `2026-06-16T20:00:00Z`.

The common `trades.json` and `summary.*` outputs stayed frozen to the ETH/SOL
authority replay. BTCUSDT remained diagnostic-only for this completed
structured-compression stream.

Package readiness failed for these reasons:

- fold `wf_2021_2022_train__2023_2024_test` selected
  `sc4h_btc_diagnostic_eth_sol_cw2_h12_t0_75_sb0_10`, but its test net P&L
  was `92.68` versus frozen test net P&L `229.02`;
- fold `wf_2021_2024_train__2025_2026_test` selected no config because the
  combined train period had `97` authority trades, below the `100`
  multi-split train gate;
- fold `wf_2023_2024_train__2025_2026_test` selected the exact frozen config
  and passed, but this was only one of three folds;
- passing the failed folds would require relaxing gates or changing the frozen
  strategy shape after seeing the walk-forward result.

Therefore the frozen config
`sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` is not package-ready. Do not
retune confirmation, max hold, target, stop buffer, symbol set, training gates,
or BTC authority around this failure.

## Premise Inventory

| Premise | Current Status | Pivot Decision |
| --- | --- | --- |
| `4h` ETH/SOL structured compression | Fragile exclusion evidence | No package, retune, gate relaxation, BTC promotion, or live-adjacent work. |
| `1h` structured compression | Failed baseline | Do not reopen from this branch. |
| Immediate clean breakout continuation | Failed baseline | Do not optimize or combine into a stream from the failed result. |
| `breakout_retest_acceptance` | Secondary universe-discovery evidence; `14` passing rows | Only open next premise. Select top non-duplicative rows from existing discovery evidence before any baseline. |
| Boundary touch rejection | Did not pass universe gate | No baseline from current docs. |
| Single-candle wick rejection | Did not pass universe gate | No baseline from current docs. |
| Failed breakout re-entry | Did not pass universe gate | No baseline from current docs. |
| Mature balance persistence | Did not pass universe gate | No baseline from current docs. |
| Hold-inside, midline, SR timing, compression breakout legacy, and impulse absorption | Closed or diagnostic | Exclusion/infrastructure only. |
| BTCUSDT higher-timeframe source lane | Valid source/premise discipline | Available only if the user chooses a BTCUSDT-only source pivot; it is not the next automatic implementation while an existing universe premise remains open. |

## Next Authorized Premise

The next bounded offline premise is:

> After a completed mature range, a breakout that retests and then accepts the
> broken range boundary may carry a cleaner range-continuation edge than the
> failed compression-expansion stream.

The next implementation may only use the existing local BTC/ETH/SOL Binance
USDT-M futures sources and the existing range-universe discovery evidence. It
must first identify the top one or two non-duplicative
`breakout_retest_acceptance` rows from
`results/futures-range-universe-discovery-audit/` and stop if no passing row is
present.

The baseline must be fixed-rule and offline:

- no new symbols, data downloads, broad mining, or optimization dimensions;
- no reuse of structured-compression target, stop, hold, confirmation, or
  training-gate tuning as a rescue path;
- closed UTC resampling from accepted `5m` parents only;
- next-bar-open entries and stop-first ambiguity only if a baseline is built;
- independent BTCUSDT, ETHUSDT, and SOLUSDT rows plus aggregate review rows;
- no strategy-promotion language until the baseline review passes after costs
  and across splits.

## Stop Conditions

The next task should stop with one of these states:

- `breakout_retest_acceptance_baseline_source_gap` if any required local
  source or resample validation fails;
- `breakout_retest_acceptance_baseline_no_ranked_candidate` if existing
  discovery artifacts do not contain a passing, non-duplicative
  `breakout_retest_acceptance` row;
- `breakout_retest_acceptance_baseline_failed_no_promotion` if the fixed
  baseline fails after costs, split stability, side balance, or symbol-transfer
  review;
- `breakout_retest_acceptance_baseline_passed_needs_robustness_review` only if
  the fixed baseline passes without changing rules after result inspection.
