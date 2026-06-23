# Futures Range Universe Discovery Spec

Date: 2026-06-26

## Verdict

Stop state:
`range_universe_spec_ready_for_audit_implementation`.

The BTCUSDT clean-breakout baseline had enough trades to be meaningful and
still failed after costs. The next useful path is not another BTCUSDT-only
slice of the failed clean-breakout template. The project should broaden to a
small local Binance USDT-M futures range universe and run one source-validated,
non-trading discovery audit that can route quickly into a fixed-rule baseline
backtest if a surface earns it.

This is a scope change from BTCUSDT-only to range-first local-universe
discovery. It is not optimization approval, portfolio-stream approval, symbol
download approval, live wiring, or strategy promotion.

## Current Baseline Failure

The last baseline was:

- `clean_breakout_4h_up_h12`: `121` trades, gross P&L `2.55`, net P&L
  `-69.29`, PF `0.8386`;
- `clean_breakout_1h_all_h12`: `1,064` trades, gross P&L `272.63`, net P&L
  `-320.25`, PF `0.8846`;
- aggregate compatibility summary: `1,185` trades, gross P&L `275.18`,
  net P&L `-389.54`, PF `0.8784`.

Both candidates were negative after costs in `2023_2024_oos` and
`2025_2026_recent`. The user-approved portfolio-stream routing rule did not
trigger because neither stream was near-viable after costs.

The next audit should therefore look for better range behavior across a small
source-validated futures universe, not retune the failed immediate-breakout
entry.

## V1 Source Universe

Use only local Binance USDT-M futures `5m` CSVs that already exist under
`../binance-bot/data/`:

| Symbol | Path | CSV Lines Including Header | Quick Physical-Order Check |
| --- | --- | ---: | --- |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | `573,985` | min open `1609459200000`, max open `1781654100000`, non-monotonic count `0` |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | `573,985` | min open `1609459200000`, max open `1781654100000`, non-monotonic count `0` |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | `573,985` | min open `1609459200000`, max open `1781654100000`, non-monotonic count `1` |

The SOL file has full quick min/max coverage but at least one non-monotonic
physical row in a simple local check. The future implementation must validate
the source explicitly and either sort accepted candles before downstream use
or fail closed with a source-gap stop state. It must not silently trust filename
coverage.

BNB, XRP, larger Binance USD-M universes, spot sources, external downloads,
official archive listing, REST tails, websockets, authenticated exchange access,
and sibling-repo data mutation are not part of this v1 scope.

## Source And Resample Gates

The future audit must validate every selected symbol before computing range
labels:

- filename and manifest identity must be Binance USDT-M futures, not spot or
  comparison-only;
- timestamps are UTC open times from closed `5m` candles;
- duplicate opens, missing `5m` opens, invalid high/low containment,
  non-positive OHLC prices, non-finite values, and negative volume reject the
  source;
- physical non-monotonic rows must be reported; accepted downstream candle
  arrays must be strictly monotonic after validation;
- split eligibility must cover `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent` or mark the symbol/split unusable;
- `15m`, `1h`, and `4h` bars must be closed UTC resamples with complete child
  `5m` buckets only.

For resampling: `15m` uses `3` child bars, `1h` uses `12`, and `4h` uses `48`.
OHLCV aggregation is first child open, max high, min low, last child close, and
summed volume. Partial final buckets, missing child opens, duplicate child
opens, forward fills, and synthetic candles are rejected.

## Candidate Families

The future audit should stay range-first and closed-candle knowable. It should
rank surfaces by symbol, timeframe, family, side, and horizon.

| Family | Timeframes | Candidate Event | Favorable Label | Adverse Label |
| --- | --- | --- | --- | --- |
| Breakout retest / acceptance | `15m`, `1h`, `4h` | Completed mature range, closed break outside the boundary, then a closed retest/acceptance from the new side | continuation away from retested boundary | re-entry through the range boundary |
| Boundary touch / rejection | `5m`, `15m`, `1h` | Closed candle touches mature range boundary without decisive outside close | inward rejection | outside acceptance |
| Failed breakout re-entry | `5m`, `15m`, `1h` | Closed break outside mature range followed by closed re-entry within a fixed window | re-entry continuation | second break / renewed outside continuation |
| Mature balance rotation / persistence | `1h`, `4h` | Mature range remains compact with repeated inside closes | internal rotation / persistence | expansion failure |
| Compression-to-expansion | `15m`, `1h`, `4h` | Reframed compression release after completed mature range, not immediate clean-breakout entry | post-release continuation with structure | failed release / re-entry |

The immediate clean-breakout baseline is closed as a failed entry template. The
future audit may keep compression-to-expansion only when the event is materially
reframed, such as requiring retest/acceptance or post-break structure before
counting a candidate.

## Discovery-To-Backtest Funnel

The next implementation should be a discovery audit, not a strategy:

1. Validate the local universe and write source/coverage artifacts.
2. Build non-trading candidate rows across symbols, timeframes, families,
   sides, and fixed horizons.
3. Rank surfaces by split-stable favorable-versus-adverse behavior, adverse
   excursion, quick invalidation, rough cost buffer, and transfer evidence.
4. If one or two surfaces clear the gate, refresh the next brief as a
   fixed-rule baseline backtest for those surfaces only.
5. Optimize only after a baseline backtest is positive or near-positive after
   costs with adequate split stability.

## Review Gate

A surface can move to a baseline backtest only if:

- source validation and resample coverage are accepted for every symbol used;
- BTC has adequate counts in every key split and at least one transfer symbol
  confirms the same family/timeframe/side behavior, or the surface is clearly
  marked symbol-specific with enough split evidence to justify that exception;
- favorable rate beats adverse rate in every key split;
- quick invalidation is below favorable rate;
- adverse excursion does not dominate favorable excursion;
- rough cost buffer survives the weakest split;
- rankings prefer cross-symbol confirmation over isolated full-sample strength.

If no surface passes, the correct next step is to stop with a no-candidate
review, not to optimize, retune, or keep slicing.

## Required Future Outputs

The next audit should write compact artifacts under
`results/futures-range-universe-discovery-audit/`:

- `futures_range_universe_sources.csv/json`
- `futures_range_universe_coverage.csv/json`
- `futures_range_universe_candidates.csv/json`
- `futures_range_universe_summary.csv/json`
- `futures_range_universe_rankings.csv/json`
- `futures_range_universe_stability.csv/json`
- normal zero-trade `summary.*` and `trades.json`

Default `cmd/rangelab` behavior must remain `lab.EmptyStrategy`. The audit must
not create entries, exits, scoring, sizing, optimization, paper/testnet/live
wiring, exchange API use, deploy files, credentials, grid, martingale,
averaging down, two-exchange logic, data downloads, or sibling repo mutations.

## Stop States

- `range_universe_source_gap`: any selected source or resample fails the
  acceptance gate.
- `range_universe_no_backtest_candidate`: implementation and outputs complete,
  but no surface passes the review gate.
- `range_universe_audit_ready_for_baseline_backtest`: at least one and at most
  two surfaces pass and the next brief is a fixed-rule baseline backtest.
- `range_universe_codegen_or_test_blocked`: implementation or verification
  cannot complete.
- `range_universe_review_only_no_strategy_change`: outputs document context
  but do not change implementation direction.
