# Backtest-First Trend-Pullback Continuation Implementation Review

Date: 2026-07-01

## Verdict

Stop state:

```text
btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy
```

The approved fixed baseline failed. Source validation and exact closed UTC
`15m` resampling passed, and the run produced enough trades to satisfy the
first-baseline trade-count gates. The candidate failed gross edge, net edge,
profit factor, and drawdown gates. It is closed as no usable strategy in this
form.

No optimizer, source expansion, derivatives-veto interaction, replay,
walk-forward, paper/testnet/live flow, exchange API work, credentials,
deployment, martingale, averaging down, two-exchange logic, or promotion is
authorized by this result.

## Implemented Candidate

Candidate:

```text
btc_15m_trend_pullback_continuation_v1
```

Implementation matched
`docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md`:

- accepted BTCUSDT Binance USDT-M futures `5m` source only;
- exact closed UTC `15m` resampling;
- `EMA(20)`, `EMA(50)`, `EMA(200)`, and `ATR(14)` on closed `15m` candles;
- fixed trend stack, `16`-bar `EMA(50)` slope, `8`-bar pullback window, and
  continuation close trigger;
- next `15m` open entry;
- stop at pullback-window invalidation plus `0.25 * ATR(14)`;
- `2.0R` target from the slipped next-open entry;
- `32` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side;
- one open position max;
- no derivatives-veto, optimizer, source expansion, session filter, volume
  filter, side selection, replay, or walk-forward interaction.

## Command

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -backtest-first-btc-15m-trend-pullback-continuation-v1 -out-dir results/backtest-first-btc-15m-trend-pullback-continuation-v1
```

Console result:

```text
backtest_first_btc_15m_trend_pullback_continuation signal_rows=13190 trades=3816 summary_rows=12 stop_state=btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy
```

## Source And Resample

Source contract:

| Field | Value |
| --- | --- |
| Source path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Base interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |
| Source facts pass | `true` |

Resample contract:

| Field | Value |
| --- | --- |
| Timeframe | exact closed UTC `15m` |
| Row count | `191,328` |
| Expected row count | `191,328` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:45:00Z` |
| Expected child bars | `3` |
| Missing child buckets | `0` |
| Source resample pass | `true` |
| Validation status | `accepted` |

## Artifacts

Output path:

```text
results/backtest-first-btc-15m-trend-pullback-continuation-v1/
```

CSV line counts:

| Artifact | Lines |
| --- | ---: |
| `btc_15m_trend_pullback_continuation_coverage.csv` | `2` |
| `btc_15m_trend_pullback_continuation_falsification.csv` | `2` |
| `btc_15m_trend_pullback_continuation_signals.csv` | `13,191` |
| `btc_15m_trend_pullback_continuation_skips.csv` | `13` |
| `btc_15m_trend_pullback_continuation_sources.csv` | `2` |
| `btc_15m_trend_pullback_continuation_summary.csv` | `13` |
| `btc_15m_trend_pullback_continuation_trades.csv` | `3,817` |
| `summary.csv` | `13` |
| Total CSV lines | `17,053` |

Common artifacts also written:

- `source_manifest.json`;
- `summary.json`;
- `summary.csv`;
- `trades.json`.

## Gate Results

| Gate | Result |
| --- | --- |
| Source and resample | pass |
| Leakage | pass |
| Trade count | pass |
| Gross edge | fail |
| Net edge | fail |
| Profit factor | fail |
| Drawdown | fail |
| Robustness | pass |
| Side reporting | pass |
| Optimizer contamination | pass |
| Derivatives-veto contamination | pass |

Falsification failure reasons:

```text
gross_edge_gate_failed
net_edge_gate_failed
profit_factor_gate_failed
drawdown_gate_failed
```

## Split Results

All-side results:

| Split | Trades | Win Rate | Gross P&L | Net P&L | PF | Max DD | Avg Hold Bars |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `1,390` | `0.375540` | `-67.949612` | `-608.659531` | `0.860832` | `0.640290` | `16.627338` |
| `2023_2024_oos` | `1,395` | `0.356989` | `-36.293536` | `-268.423414` | `0.779940` | `0.275926` | `16.318280` |
| `2025_2026_recent` | `1,031` | `0.350145` | `-19.184695` | `-79.483644` | `0.743475` | `0.085763` | `16.688652` |
| `full_2021_2026` | `3,816` | `0.361897` | `-123.427844` | `-956.566589` | `0.837958` | `0.958438` | `16.530922` |

Full-period side split:

| Side | Trades | Win Rate | Gross P&L | Net P&L | PF | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Long | `2,005` | `0.368579` | `-122.716000` | `-578.903589` | `0.816389` | `0.597338` |
| Short | `1,811` | `0.354500` | `-0.711844` | `-377.663000` | `0.862684` | `0.431704` |

## Interpretation

The candidate fails before costs: full-sample gross P&L is negative and every
primary split is gross-negative. Costs and slippage deepen the loss, but they
are not the root cause. The short side is closer to flat gross over the full
sample, but it is still net-negative and cannot be selected post-result without
violating the fixed packet's no side-selection boundary.

Trade count is not the issue. The run produced `3,816` trades, with at least
`1,031` trades in every primary split. The failure is signal quality and loss
shape, not insufficient sample size.

## Closed Boundary

Do not rescue this fixed baseline with alternate EMA lengths, slope lookbacks,
pullback windows, EMA-band definitions, continuation triggers, stop buffers,
target R values, time stops, side selection, session filters, volume filters,
volatility filters, derivatives-veto interaction, source expansion, replay,
walk-forward, or optimizer grids.

The next strategy work, if any, needs a separately approved materially different
lane or a docs-only strategy-class/premise decision. This result does not create
a promotion path.
