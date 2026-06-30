# Backtest-First BTC 15m Range-Edge Exhaustion Fade Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy
```

The fixed offline baseline candidate was implemented and locally verified:

```text
btc_15m_range_edge_exhaustion_fade_v1
```

The fixed baseline failed. It produced enough trades to pass the first-baseline
trade-count gates, but failed gross edge, net edge, and drawdown gates. It is
closed as no usable strategy in this form.

## Fixed Baseline Scope

The implementation follows the selected packet:

- current accepted BTCUSDT Binance USDT-M futures `5m` CSV;
- exact closed UTC `15m` resample from complete three-child `5m` buckets;
- prior `96` closed `15m` bars `[d-96,d-1]` as the local range;
- require `close[d-3]` to start between the lower and upper `40%` of the range;
- long fade: last three closes move downward into the lower `15%` of the range,
  `close[d]` remains above range low, and final downward progress is less than
  `0.35 * ATR(14)[d-1]`;
- short fade: symmetric upward move into the upper `15%` of the range, still
  below range high, with weak final progress;
- next-`15m`-open execution through the existing backtest engine;
- stop beyond the local range edge by `0.25 * ATR(14)[d-1]`;
- target at the local range midpoint;
- `16` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side.

No optimizer, replay, walk-forward, derivatives-veto interaction, source
expansion, paper/testnet/live path, exchange API work, credentials, deploy file,
martingale, averaging down, two-exchange logic, or promotion was added.

## Local Verification

The user ran the required local verification from branch `edge-fade-baseline`.

Build/test:

```text
ok      range-strategy-lab/cmd/rangelab 0.366s
ok      range-strategy-lab/internal/lab 0.065s
```

Backtest command output:

```text
backtest_first_btc_15m_range_edge_exhaustion_fade signal_rows=156 trades=156 summary_rows=12 stop_state=btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy
```

CSV line counts:

```text
    2 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_coverage.csv
    2 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_falsification.csv
  157 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_signals.csv
    2 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_skips.csv
    2 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_sources.csv
   13 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_summary.csv
  157 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_trades.csv
   13 results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/summary.csv
  348 total
```

## Falsification Result

```json
{
  "backtest_name": "backtest_first_btc_15m_range_edge_exhaustion_fade",
  "candidate_id": "btc_15m_range_edge_exhaustion_fade_v1",
  "stop_state": "btc_15m_range_edge_exhaustion_fade_backtest_failed_no_usable_strategy",
  "source_resample_pass": true,
  "leakage_pass": true,
  "trade_count_pass": true,
  "gross_edge_pass": false,
  "net_edge_pass": false,
  "drawdown_pass": false,
  "robustness_pass": true,
  "side_reporting_pass": true,
  "full_executed_trades": 156,
  "required_full_executed_trades": 120,
  "minimum_primary_split_executed_trades": 38,
  "required_primary_split_trades": 25,
  "full_gross_pnl": -154.40528599997904,
  "full_net_pnl": -261.59525647142874,
  "full_profit_factor": 0.48125879295748447,
  "full_max_drawdown": 0.28473381700333156,
  "dominant_primary_split_trade_share": 0.40384615384615385,
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
| `2021_2022_stress` | all | 63 | 0.238095 | -34.287918 | -81.312641 | 0.656318 | 0.126355 | 3.682540 |
| `2023_2024_oos` | all | 55 | 0.181818 | -75.755449 | -112.884803 | 0.350124 | 0.116621 | 3.654545 |
| `2025_2026_recent` | all | 38 | 0.184211 | -44.361918 | -67.397812 | 0.282957 | 0.072759 | 3.947368 |
| `full_2021_2026` | all | 156 | 0.205128 | -154.405286 | -261.595256 | 0.481259 | 0.284734 | 3.737179 |
| `full_2021_2026` | long | 88 | 0.181818 | -106.970632 | -167.547236 | 0.437620 | 0.196456 | 3.352273 |
| `full_2021_2026` | short | 68 | 0.235294 | -47.434654 | -94.048021 | 0.544260 | 0.094048 | 4.235294 |

## Interpretation

The baseline failed before any question of tuning or confirmation:

- trade count passed with `156` trades;
- minimum primary split trade count passed with `38` trades;
- source/resample, leakage, side reporting, and robustness gates passed;
- gross edge failed across full, stress, OOS, and recent splits;
- net edge failed across all primary splits;
- drawdown failed, with full max drawdown of `0.284734` versus the fixed `0.25`
  gate;
- profit factor was far below a usable threshold.

This candidate should not be rescued by retuning.

## Closure Boundary

Do not rescue `btc_15m_range_edge_exhaustion_fade_v1` with alternate range
windows, progress thresholds, edge zones, midpoint variants, added volume
filters, derivatives context, replay, walk-forward, or optimizer grids.

All candidates in the current backtest-first candidate packet have now failed as
fixed baselines. The next research step should be a new materially different
candidate packet, not a retune of these three baselines.
