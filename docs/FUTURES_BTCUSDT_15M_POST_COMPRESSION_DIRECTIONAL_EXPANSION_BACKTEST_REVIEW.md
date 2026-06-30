# Futures BTCUSDT 15m Post-Compression Directional Expansion Backtest Review

Date: 2026-06-30

## Verdict

Stop state:
`post_compression_directional_expansion_backtest_failed_no_usable_strategy`.

The user explicitly approved implementing the offline backtest for:

```text
btc_15m_post_compression_l192_q20_m020_none_long_h48_v1
```

The implementation reproduced the approved source contract, exact closed UTC
`15m` resample, and representative raw candidate identity (`468` signal rows
before one-position filtering). It produced enough trades for review (`421`
executed trades; minimum primary split `123`), but the fixed model failed the
declared strategy gates after costs and under the extra slippage-stress view.

This review closes the fixed representative candidate as not usable in this
form. It does not authorize parameter rescue, adjacent-cell P&L selection,
exit retuning, derivatives-veto interaction, replay, walk-forward, source
expansion, paper/testnet/live work, exchange API use, or promotion.

## What Was Implemented

- CLI flag:
  `-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest`.
- Default output directory:
  `results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/`.
- Focused fixed-cell backtest module:
  `internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest.go`.
- Tests:
  `internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest_test.go`
  and `cmd/rangelab/main_test.go`.

The model is exactly the docs-only spec model:

| Component | Value |
| --- | --- |
| Direction | long only |
| Decision stream | closed UTC BTCUSDT `15m` candles |
| Source | local Binance USDT-M futures `5m` BTCUSDT CSV |
| Event | `L=192`, `q20`, `M=0.2`, volume `none` |
| Entry | next `15m` open |
| Stop | `entry_price - 1.0 * ATR(14)[d-1]` |
| Target | `entry_price + 2.0 * ATR(14)[d-1]` |
| Max hold | `48` bars |
| Positioning | one position max |
| Ambiguity | stop-first |
| Sizing/costs | `1%` risk, `1x` notional cap, `0.0004` fee, `0.000116` slippage |

The implementation reports both engine net and extra slippage-stress net. The
falsification decision uses the extra stress view.

## Source And Resample

Source facts reproduced:

| Field | Value |
| --- | ---: |
| Source path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |

Exact `15m` resample facts reproduced:

| Field | Value |
| --- | ---: |
| Row count | `191,328` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:45:00Z` |
| Last close | `2026-06-16T23:59:59Z` |
| Expected child bars | `3` |
| Missing child opens | `0` |
| Validation status | `accepted` |

## Artifacts

Output dir:
`results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/`.

Artifacts written:

```text
source_manifest.json
summary.json
summary.csv
trades.json
btc_15m_post_compression_l192_q20_m020_none_long_h48_sources.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_resample_coverage.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_skips.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_trades.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_summary.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_cost_stress.csv/json
btc_15m_post_compression_l192_q20_m020_none_long_h48_falsification.json
```

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.csv` | `469` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_trades.csv` | `422` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_summary.csv` | `13` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_cost_stress.csv` | `13` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_skips.csv` | `7` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_sources.csv` | `2` |
| `btc_15m_post_compression_l192_q20_m020_none_long_h48_resample_coverage.csv` | `2` |
| `summary.csv` | `13` |
| Total CSV lines | `941` |

## Falsification

| Gate | Result |
| --- | --- |
| Source/resample | pass |
| Candidate identity | pass, `468 / 468` raw rows |
| Leakage | pass |
| Trade count | pass, `421` full and min primary split `123` |
| Gross edge | fail |
| Extra slippage-stress edge | fail |
| Stress profit factor | fail |
| Drawdown | fail |
| Robustness | pass |
| Optimizer contamination | pass |
| Closed-family protection | pass |
| Derivatives-veto contamination | pass |

Failure reasons:

```text
gross_edge_gate_failed
extra_slippage_stress_edge_gate_failed
stress_profit_factor_gate_failed
drawdown_gate_failed
```

## Split Results

Strategy-specific pass/fail uses extra slippage-stress net, not common engine
net.

| Split | Raw rows | Trades | Win rate | Gross PnL | Engine net | Extra stress net | Stress PF | Stress max DD |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `167` | `152` | `0.427632` | `185.552840` | `57.866623` | `20.837617` | `1.042668` | `0.077273` |
| `2023_2024_oos` | `159` | `146` | `0.342466` | `38.807900` | `-80.852256` | `-115.553702` | `0.684405` | `0.117669` |
| `2025_2026_recent` | `142` | `123` | `0.325203` | `-15.799742` | `-106.272938` | `-132.510165` | `0.526275` | `0.141135` |
| `full_2021_2026` | `468` | `421` | `0.368171` | `208.560999` | `-129.258571` | `-227.226250` | `0.799666` | `0.289326` |

Common `summary.csv` is retained as the engine-net compatibility view. It
shows full-sample gross PnL positive (`208.560999`) but engine net negative
(`-129.258571`) after slipped prices and fees. The required extra stress view
subtracts the recorded slippage proxy once more, producing full stress net
`-227.226250`.

## Interpretation

The zero-trade label-separation pocket was real enough to produce a coherent
candidate stream and sufficient trade count. It did not survive the first fixed
offline execution model after costs. The recent split was gross-negative before
the extra stress adjustment, and both OOS/recent stress profit factors were far
below the required thresholds.

The failure is therefore not a source, resample, leakage, or sample-size
failure. It is a strategy-economics failure for this exact fixed model.

## Boundaries

Closed by this review:

- `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1` as a usable fixed
  offline strategy in this form;
- adjacent-cell P&L rescue from the zero-trade passing pocket;
- stop/target/hold retuning from this result;
- derivatives-veto interaction on this failed candidate stream.

Still not authorized:

- full `81`-cell grid P&L;
- shorts, `h16`/`h32`, `L48`/`L96`, `q30`/`q40`, other breakout multiples, or
  volume filters as rescue candidates;
- derivatives veto application;
- replay, walk-forward, optimizer selection, strategy promotion,
  paper/testnet/live paths, exchange API work, credentials, deploy files,
  martingale, averaging down, or two-exchange logic.

## Commands And Outcomes

```bash
/usr/local/go/bin/gofmt -w cmd/rangelab/main.go cmd/rangelab/main_test.go internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest.go internal/lab/futures_btc_15m_post_compression_l192_q20_m020_none_long_h48_backtest_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest
wc -l results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/*.csv
```

Outcomes:

- tests passed;
- backtest wrote the expected artifacts;
- source/resample and candidate identity passed;
- stop state:
  `post_compression_directional_expansion_backtest_failed_no_usable_strategy`.
