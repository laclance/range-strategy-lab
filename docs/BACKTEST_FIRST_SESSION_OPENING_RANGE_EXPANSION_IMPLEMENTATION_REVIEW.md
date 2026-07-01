# Backtest-First Session Opening-Range Expansion Implementation Review

Date: 2026-07-01

## Verdict

Stop state:

```text
btc_15m_session_opening_range_expansion_backtest_failed_no_usable_strategy
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
btc_15m_session_opening_range_expansion_v1
```

Implementation matched
`docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md`:

- accepted BTCUSDT Binance USDT-M futures `5m` source only;
- exact closed UTC `15m` resampling;
- fixed `13:30:00Z` UTC session anchor with no DST shifting or alternate anchor
  comparison;
- opening range from the four closed `15m` bars with opens at `13:30`, `13:45`,
  `14:00`, and `14:15` UTC;
- expansion window using closed decision bars with opens in `[14:30, 17:30) UTC`;
- closed-candle acceptance outside the opening range by `0.10 * ATR(14)`;
- next `15m` open entry;
- opposite-side opening-range stop with `0.10 * ATR(14)` buffer;
- `1.5R` target from the slipped next-open entry;
- `24` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side;
- one open position max and at most one qualified signal per UTC date;
- no derivatives-veto, optimizer, source expansion, session mining, volume
  filter, volatility filter, side selection, replay, or walk-forward
  interaction.

## Command

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -backtest-first-btc-15m-session-opening-range-expansion-v1 -out-dir results/backtest-first-btc-15m-session-opening-range-expansion-v1
```

Console result:

```text
backtest_first_btc_15m_session_opening_range_expansion session_range_rows=1993 signal_rows=1652 trades=1652 summary_rows=12 stop_state=btc_15m_session_opening_range_expansion_backtest_failed_no_usable_strategy
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
results/backtest-first-btc-15m-session-opening-range-expansion-v1/
```

CSV line counts:

| Artifact | Lines |
| --- | ---: |
| `btc_15m_session_opening_range_expansion_coverage.csv` | `2` |
| `btc_15m_session_opening_range_expansion_falsification.csv` | `2` |
| `btc_15m_session_opening_range_expansion_session_ranges.csv` | `1,994` |
| `btc_15m_session_opening_range_expansion_signals.csv` | `1,653` |
| `btc_15m_session_opening_range_expansion_skips.csv` | `5` |
| `btc_15m_session_opening_range_expansion_sources.csv` | `2` |
| `btc_15m_session_opening_range_expansion_summary.csv` | `13` |
| `btc_15m_session_opening_range_expansion_trades.csv` | `1,653` |
| `summary.csv` | `13` |
| Total CSV lines | `5,337` |

Common artifacts also written:

- `source_manifest.json`;
- `summary.json`;
- `summary.csv`;
- `trades.json`.

## Gate Results

| Gate | Result |
| --- | --- |
| Source and resample | pass |
| Fixed session spec | pass |
| Leakage | pass |
| Trade count | pass |
| Gross edge | fail |
| Net edge | fail |
| Profit factor | fail |
| Drawdown | fail |
| Side reporting | pass |
| Combined baseline selection | pass |
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

| Split | Trades | Win Rate | Gross P&L | Net P&L | PF | Gross PF | Max DD | Avg Hold Bars |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `604` | `0.420530` | `45.770589` | `-237.826052` | `0.889580` | `1.023023` | `0.294840` | `16.203642` |
| `2023_2024_oos` | `604` | `0.418874` | `122.959395` | `-139.289001` | `0.899614` | `1.099462` | `0.184699` | `15.821192` |
| `2025_2026_recent` | `444` | `0.421171` | `-16.997499` | `-169.125465` | `0.799265` | `0.977505` | `0.190476` | `16.650901` |
| `full_2021_2026` | `1,652` | `0.420097` | `151.732485` | `-546.240518` | `0.875398` | `1.038124` | `0.578097` | `16.184019` |

Full-period side split:

| Side | Trades | Win Rate | Gross P&L | Net P&L | PF | Gross PF | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Long | `825` | `0.431515` | `116.773979` | `-237.193826` | `0.890752` | `1.059232` | `0.295941` |
| Short | `827` | `0.408706` | `34.958506` | `-309.046692` | `0.860333` | `1.017406` | `0.361856` |

## Interpretation

Trade count is not the issue. The run produced `1,652` executed trades, with at
least `444` trades in every primary split. The fixed opening-range rule created
a mildly positive full-sample gross result, but the recent split was
gross-negative and the full sample could not survive the predeclared fees and
slippage.

Both sides are net-negative in the full sample. The long side has better gross
P&L than the short side, but selecting it post-result would violate the packet's
combined-baseline and no-side-selection boundary. The drawdown failure is also
large: full-sample max drawdown is `57.809706%`, versus the `25%` limit.

The failure is therefore not a source, resample, leakage, fixed-session, or
sample-size failure. It is a strategy-economics failure for this exact fixed
model.

## Commands And Outcomes

```bash
gofmt -w cmd/rangelab/session_opening_range.go cmd/rangelab/main_test.go internal/lab/backtest_first_btc_15m_session_opening_range_expansion.go internal/lab/backtest_first_btc_15m_session_opening_range_expansion_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -backtest-first-btc-15m-session-opening-range-expansion-v1 -out-dir results/backtest-first-btc-15m-session-opening-range-expansion-v1
wc -l results/backtest-first-btc-15m-session-opening-range-expansion-v1/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
```

Outcomes:

- pre-change and post-change Go tests passed;
- backtest wrote the expected artifacts;
- source/resample passed;
- stop state:
  `btc_15m_session_opening_range_expansion_backtest_failed_no_usable_strategy`.

## Closed Boundary

Do not rescue this fixed baseline with alternate UTC anchors, opening-range
lengths, expansion windows, acceptance buffers, ATR windows, stop buffers,
target R values, time stops, one-trade-per-day changes, side selection, weekday
filters, volume filters, volatility filters, derivatives-veto interaction,
source expansion, replay, walk-forward, or optimizer grids.

Do not rebrand closed range-reversion, midpoint, edge-fade, previous-day range,
value-area, range-optimization, post-compression, clean-breakout,
router/rotation, or trend-pullback branches as opening-range expansion.

The next strategy work, if any, needs a separately approved materially different
lane or an explicit docs-only strategy-class/premise decision. This result does
not create a promotion path.
