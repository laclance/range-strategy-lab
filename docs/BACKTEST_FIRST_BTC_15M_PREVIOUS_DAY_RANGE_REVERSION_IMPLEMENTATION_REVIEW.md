# Backtest-First BTC 15m Previous-Day Range Reversion Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy
```

The selected fixed offline baseline candidate was implemented and locally
verified:

```text
btc_15m_previous_day_range_reversion_v1
```

The fixed baseline failed. It produced enough trades, but failed gross edge, net
edge, and drawdown gates. It is closed as no usable strategy in this form.

## What Was Added

- Lab implementation for the previous-day range reversion fixed baseline.
- Offline CLI flag:

```text
-backtest-first-btc-15m-previous-day-range-reversion-v1
```

- Default output path:

```text
results/backtest-first-btc-15m-previous-day-range-reversion-v1/
```

## Fixed Baseline Scope

The implementation follows the selected packet:

- current accepted BTCUSDT Binance USDT-M futures `5m` CSV;
- exact closed UTC `15m` resample from complete three-child `5m` buckets;
- prior complete UTC day's high, low, and midpoint;
- skip the current day after any prior current-day close outside the previous
  day's high-low range;
- long entry when a closed `15m` candle is inside the previous-day range and in
  its lower `10%`;
- short entry when a closed `15m` candle is inside the previous-day range and in
  its upper `10%`;
- next-`15m`-open execution through the existing backtest engine;
- stop beyond the prior-day range extreme by `0.25 * ATR(14)[d-1]`;
- target at the prior-day midpoint;
- `24` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side.

No optimizer, replay, walk-forward, derivatives-veto interaction, source
expansion, paper/testnet/live path, exchange API work, credentials, deploy file,
martingale, averaging down, two-exchange logic, or promotion was added.

## Local Verification

The user ran the required local verification from branch `day-baseline`.

Build/test:

```text
ok      range-strategy-lab/cmd/rangelab 0.389s
ok      range-strategy-lab/internal/lab 0.050s
```

Backtest command output:

```text
backtest_first_btc_15m_previous_day_range_reversion signal_rows=1956 trades=1956 summary_rows=12 stop_state=btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy
```

CSV line counts:

```text
     2 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_coverage.csv
     2 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_falsification.csv
  1957 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_signals.csv
     4 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_skips.csv
     2 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_sources.csv
    13 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_summary.csv
  1957 results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_trades.csv
    13 results/backtest-first-btc-15m-previous-day-range-reversion-v1/summary.csv
  3950 total
```

## Falsification Result

```json
{
  "backtest_name": "backtest_first_btc_15m_previous_day_range_reversion",
  "candidate_id": "btc_15m_previous_day_range_reversion_v1",
  "stop_state": "btc_15m_previous_day_range_reversion_backtest_failed_no_usable_strategy",
  "source_resample_pass": true,
  "leakage_pass": true,
  "trade_count_pass": true,
  "gross_edge_pass": false,
  "net_edge_pass": false,
  "drawdown_pass": false,
  "robustness_pass": true,
  "side_reporting_pass": true,
  "full_executed_trades": 1956,
  "required_full_executed_trades": 120,
  "minimum_primary_split_executed_trades": 577,
  "required_primary_split_trades": 25,
  "full_gross_pnl": -506.1461042753044,
  "full_net_pnl": -935.263471879938,
  "full_profit_factor": 0.5772171917112771,
  "full_max_drawdown": 0.9388204368154983,
  "dominant_primary_split_trade_share": 0.3537832310838446,
  "failure_reasons": [
    "gross_edge_gate_failed",
    "net_edge_gate_failed",
    "drawdown_gate_failed"
  ]
}
```

## Key Split Summary

| Split | Side | Trades | Win rate | Gross P&L | Net P&L | Profit factor | Max drawdown | Avg hold bars |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | all | 687 | 0.227074 | -442.487425 | -726.424204 | 0.563333 | 0.728566 | 7.113537 |
| `2023_2024_oos` | all | 692 | 0.229769 | -50.692741 | -153.571356 | 0.611159 | 0.153571 | 7.523121 |
| `2025_2026_recent` | all | 577 | 0.220104 | -12.965939 | -55.267912 | 0.640300 | 0.060864 | 7.759099 |
| `full_2021_2026` | all | 1956 | 0.225971 | -506.146104 | -935.263472 | 0.577217 | 0.938820 | 7.448875 |
| `full_2021_2026` | long | 915 | 0.220765 | -328.516719 | -525.656640 | 0.513433 | 0.540102 | 7.311475 |
| `full_2021_2026` | short | 1041 | 0.230548 | -177.629386 | -409.606832 | 0.638100 | 0.409887 | 7.569645 |

## Interpretation

The baseline failed before any question of tuning or confirmation:

- trade count passed with `1,956` trades;
- minimum primary split trade count passed with `577` trades;
- source/resample, leakage, side reporting, and robustness gates passed;
- gross edge failed across full, stress, OOS, and recent splits;
- net edge failed across all primary splits;
- drawdown failed, with full max drawdown of `0.938820`;
- profit factor was far below a usable threshold.

This is not a close failure. The candidate should not be rescued by retuning.

## Closure Boundary

Do not rescue `btc_15m_previous_day_range_reversion_v1` with alternate UTC
sessions, previous `2`/`3` day windows, changed outer-decile thresholds,
derivatives context, calendar/session mining, replay, walk-forward, or optimizer
grids.

The next research step, if continuing the backtest-first lane, should move to a
materially different candidate from `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md`,
not a retuned previous-day range variant.
