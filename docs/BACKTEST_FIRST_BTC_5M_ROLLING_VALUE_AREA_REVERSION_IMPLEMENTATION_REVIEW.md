# Backtest-First BTC 5m Rolling Value-Area Reversion Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy
```

The selected fixed offline baseline candidate was implemented and locally
verified:

```text
btc_5m_rolling_value_area_reversion_v1
```

The fixed baseline failed decisively. It produced more than enough trades, but
failed gross edge, net edge, and drawdown gates. It is closed as no usable
strategy in this form.

## What Was Added

- Lab implementation:
  - `internal/lab/backtest_first_btc_5m_value_area_types.go`
  - `internal/lab/backtest_first_btc_5m_value_area_runner.go`
  - `internal/lab/backtest_first_btc_5m_value_area_support.go`
- Offline CLI entrypoint:
  - `cmd/rangelab/backtest_first_value_area_reversion.go`
- Fixed flag:

```text
-backtest-first-btc-5m-rolling-value-area-reversion-v1
```

- Default output path:

```text
results/backtest-first-btc-5m-rolling-value-area-reversion-v1/
```

## Fixed Baseline Scope

The implementation follows the selected packet:

- native closed BTCUSDT `5m` candles;
- prior `288` closed `5m` bars as the rolling value-area window;
- rolling VWAP from typical price weighted by volume;
- `ATR(14)[d-1]` known at decision time;
- minimum range width of `6 * ATR(14)[d-1]`;
- lower/upper `20%` outer-zone entries;
- `0.15 * range_width` minimum distance from VWAP;
- next-`5m`-open execution through the existing backtest engine;
- `36` closed `5m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side.

No optimizer, replay, walk-forward, derivatives-veto interaction, source
expansion, paper/testnet/live path, exchange API work, credentials, deploy file,
martingale, averaging down, two-exchange logic, or promotion was added.

## Local Verification

The user ran the required local verification from branch
`value-area-reversion-baseline-2`.

Source file check:

```text
573985 ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

Build/test:

```text
ok      range-strategy-lab/cmd/rangelab 0.310s
ok      range-strategy-lab/internal/lab 0.044s
```

Backtest command output:

```text
backtest_first_btc_5m_rolling_value_area_reversion signal_rows=20144 trades=20144 summary_rows=12 stop_state=btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy
```

CSV line counts:

```text
      2 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_falsification.csv
  20145 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_signals.csv
      2 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_skips.csv
      2 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_sources.csv
     13 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_summary.csv
  20145 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/btc_5m_rolling_value_area_reversion_trades.csv
     13 results/backtest-first-btc-5m-rolling-value-area-reversion-v1/summary.csv
  40322 total
```

## Falsification Result

```json
{
  "backtest_name": "backtest_first_btc_5m_rolling_value_area_reversion",
  "candidate_id": "btc_5m_rolling_value_area_reversion_v1",
  "stop_state": "btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy",
  "source_pass": true,
  "leakage_pass": true,
  "trade_count_pass": true,
  "gross_edge_pass": false,
  "net_edge_pass": false,
  "drawdown_pass": false,
  "robustness_pass": true,
  "side_reporting_pass": true,
  "full_executed_trades": 20144,
  "required_full_executed_trades": 120,
  "minimum_primary_split_executed_trades": 6070,
  "required_primary_split_trades": 25,
  "full_gross_pnl": -489.49542250301454,
  "full_net_pnl": -999.9999999684289,
  "full_profit_factor": 0.6949768102636272,
  "full_max_drawdown": 0.9999999999689176,
  "dominant_primary_split_trade_share": 0.3622418586179508,
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
| `2021_2022_stress` | all | 7297 | 0.286556 | -489.466883 | -999.933313 | 0.694983 | 0.999933 | 15.207208 |
| `2023_2024_oos` | all | 6777 | 0.278589 | -0.028533 | -0.066666 | 0.558752 | 0.000067 | 14.642467 |
| `2025_2026_recent` | all | 6070 | 0.246623 | -0.000007 | -0.000021 | 0.565299 | 0.000000 | 13.625535 |
| `full_2021_2026` | all | 20144 | 0.271843 | -489.495423 | -1000.000000 | 0.694977 | 1.000000 | 14.540608 |
| `full_2021_2026` | long | 9509 | 0.273110 | -170.280394 | -356.450109 | 0.751346 | 0.371184 | 14.406983 |
| `full_2021_2026` | short | 10635 | 0.270710 | -319.215028 | -643.549891 | 0.651178 | 0.643550 | 14.660085 |

## Interpretation

The baseline failed before any question of tuning or confirmation:

- trade count passed with `20,144` trades;
- split trade count passed with minimum primary split trades of `6,070`;
- both long and short sides were reported;
- source and leakage gates passed;
- gross edge failed across the full sample and primary splits;
- net edge failed, with full net P&L effectively depleting starting equity;
- drawdown failed at effectively `100%` full-sample drawdown;
- profit factor was far below a usable threshold.

This is not a close failure. The candidate should not be rescued by retuning.

## Closure Boundary

Do not rescue `btc_5m_rolling_value_area_reversion_v1` with alternate VWAP
windows, outer-zone percentages, target changes, time-stop changes, side
selection, volume filters, derivatives-veto interaction, replay, walk-forward,
or optimizer grids.

The next research step, if continuing the backtest-first lane, should move to a
materially different parked candidate from `docs/BACKTEST_FIRST_CANDIDATE_PACKET.md`,
not a retuned value-area variant.
