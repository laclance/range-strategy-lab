# Backtest-First Trend-Pullback Continuation Candidate Packet

Date: 2026-07-01

## Verdict

Stop state:

```text
trend_pullback_candidate_packet_ready_for_implementation_approval
```

Selected fixed baseline:

```text
btc_15m_trend_pullback_continuation_v1
```

This is a docs-only backtest-first candidate packet for the selected
trend-pullback continuation lane. It defines exactly one fixed BTCUSDT Binance
USDT-M futures baseline for a later approval gate. It does not authorize Go code,
CLI flags, generated outputs, source downloads, optimizers, replay,
walk-forward, derivatives-veto interaction, paper/testnet/live paths, exchange
API work, credentials, deploy files, martingale, averaging down, two-exchange
logic, or promotion.

## Source Contract

The candidate uses the current accepted local source:

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

The implementation, if later approved, must resample the accepted `5m` source to
exact closed UTC `15m` bars. Decisions use confirmed closed `15m` candles only.
Entries happen on the next `15m` bar open. Same-bar stop/target ambiguity is
stop-first. Costs use `0.0004` fee per side and `0.000116` slippage per side.

## Why This Is A New Lane

This candidate follows an established directional trend after a shallow
pullback. It does not fade range edges, target range midpoints, trade
previous-day range reversion, select range-router states, reuse the range
optimization workbench, or revive the failed post-compression directional
expansion candidate.

It also does not use cross-asset context, derivatives source data, the
`btc_15m_basis_discount_no_trade_veto_v1` veto, volume filters, session filters,
or any post-result side selection. The first test is only a BTCUSDT OHLCV
trend-pullback continuation baseline.

## Candidate

Candidate id:

```text
btc_15m_trend_pullback_continuation_v1
```

### Hypothesis

When BTCUSDT is in a clear `15m` directional trend, a shallow pullback into the
fast/medium EMA band that holds the trend structure can resume far enough to pay
standard futures fees and slippage.

### Source And Timeframe

- Source: current accepted BTCUSDT Binance USDT-M futures `5m` CSV.
- Decision timeframe: exact closed UTC `15m` bars resampled from the `5m` source.
- Entry timing: next `15m` open after the closed decision candle.
- Position model: one open position max.

### Fixed Indicators

All indicators are computed on resampled closed UTC `15m` candles:

- `EMA(20)` of close;
- `EMA(50)` of close;
- `EMA(200)` of close;
- `ATR(14)` using true range with Wilder-style smoothing.

Skip all decision candles until every indicator and the slope reference below is
available. The trend slope reference is `16` closed `15m` bars.

### Fixed Entry Rule

For each closed `15m` decision candle `d`, skip if a position is already open.
If both long and short conditions somehow qualify on the same candle, skip the
candle as ambiguous.

Long setup:

1. At `d-1`, require `EMA(20) > EMA(50) > EMA(200)`.
2. At `d-1`, require `EMA(50)[d-1] > EMA(50)[d-17]`.
3. At `d-1`, require `close[d-1] > EMA(50)[d-1]`.
4. In the pullback window `[d-8, d-1]`, require at least one candle `p` where
   `low[p] <= EMA(20)[p]` and `close[p] >= EMA(50)[p]`.
5. In the same pullback window, reject the setup if any candle closes below
   `EMA(50)`.
6. On the decision candle `d`, require `close[d] > high[d-1]`.
7. On the decision candle `d`, require `close[d] > EMA(20)[d]`.
8. On the decision candle `d`, require `close[d] > open[d]`.
9. Enter long on `open[d+1]`.

Short setup:

1. At `d-1`, require `EMA(20) < EMA(50) < EMA(200)`.
2. At `d-1`, require `EMA(50)[d-1] < EMA(50)[d-17]`.
3. At `d-1`, require `close[d-1] < EMA(50)[d-1]`.
4. In the pullback window `[d-8, d-1]`, require at least one candle `p` where
   `high[p] >= EMA(20)[p]` and `close[p] <= EMA(50)[p]`.
5. In the same pullback window, reject the setup if any candle closes above
   `EMA(50)`.
6. On the decision candle `d`, require `close[d] < low[d-1]`.
7. On the decision candle `d`, require `close[d] < EMA(20)[d]`.
8. On the decision candle `d`, require `close[d] < open[d]`.
9. Enter short on `open[d+1]`.

### Fixed Risk And Exits

Use the same fixed risk model for both sides:

- Long stop: `min(low[d-8] ... low[d]) - 0.25 * ATR(14)[d]`.
- Short stop: `max(high[d-8] ... high[d]) + 0.25 * ATR(14)[d]`.
- Skip the signal if stop distance from the next-bar-open entry is non-positive.
- Target: `2.0R` from entry using the initial stop distance.
- Time stop: close after `32` closed `15m` bars after entry.
- No trailing stop, break-even move, scale-in, scale-out, re-entry while a
  position is open, martingale, or averaging down.
- Sizing: `1%` risk at stop, capped at `1x` notional.
- Costs: `0.0004` fee per side and `0.000116` slippage per side.

### Expected Output Path And Artifacts

If later approved for implementation, use this output path:

```text
results/backtest-first-btc-15m-trend-pullback-continuation-v1/
```

Expected artifacts:

- `source_manifest.json`;
- `summary.json`, `summary.csv`, and `trades.json`;
- candidate-specific coverage, sources, signals, skips, trades, summary, and
  falsification CSV/JSON files.

### Pass/Fail Gates

This fixed baseline may request review only if all gates below pass:

- source validation passes with the accepted BTCUSDT Binance USDT-M futures `5m`
  source contract;
- `15m` resampling uses only complete closed UTC bars and records coverage;
- no leakage is found in features, trend state, pullback state, trigger, stop,
  target, sizing, or selection;
- full-sample executed trades are at least `120`;
- `2021_2022_stress`, `2023_2024_oos`, and `2025_2026_recent` each have at
  least `25` executed trades;
- full-sample gross P&L is positive;
- `2023_2024_oos` and `2025_2026_recent` gross P&L are non-negative;
- full-sample net P&L after the fixed fees and slippage is positive;
- `2023_2024_oos` and `2025_2026_recent` net P&L after fixed fees and slippage
  are non-negative;
- full-sample profit factor is at least `1.10`;
- `2023_2024_oos` and `2025_2026_recent` profit factors are at least `1.00`;
- full-sample max drawdown is no worse than `25%`;
- each primary split max drawdown is no worse than `30%`;
- long and short side results are reported separately; and
- the combined fixed baseline, not post-result side selection, is what passes or
  fails.

Fail fast if gross edge is absent, costs kill the edge, OOS/recent splits fail,
trade count is too small, drawdown is unacceptable, same-bar optimism is
required, leakage is found, or the result needs retuning to survive.

### No-Rescue Boundaries

If this fixed baseline fails, do not rescue it with alternate EMA lengths, slope
lookbacks, pullback windows, EMA-band definitions, continuation triggers, stop
buffers, target R values, time stops, side selection, session filters, volume
filters, volatility filters, derivatives-veto interaction, source expansion,
replay, walk-forward, or optimizer grids.

Record the result and either close the candidate or move to a materially
different approved lane.

## Next Approval Gate

The next bounded gate is explicit operator approval to implement and run this
one fixed offline baseline. Without that approval, no Go code, CLI flag,
generated result, backtest, optimizer, source expansion, derivatives-veto
interaction, exchange API work, credential, deployment, paper/testnet/live flow,
martingale, averaging down, two-exchange logic, or promotion is authorized.
