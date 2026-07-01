# Backtest-First Session Opening-Range Expansion Candidate Packet

Date: 2026-07-01

## Verdict

Stop state:

```text
session_opening_range_expansion_candidate_packet_ready_for_implementation_approval
```

Selected fixed baseline:

```text
btc_15m_session_opening_range_expansion_v1
```

This is a docs-only backtest-first candidate packet for the selected
session-based opening-range expansion lane. It defines exactly one fixed BTCUSDT
Binance USDT-M futures baseline for a later approval gate. It does not authorize
Go code, CLI flags, generated outputs, source downloads, optimizers, replay,
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

This candidate uses a predeclared UTC session time box and tests continuation
away from that box after a closed-candle expansion. It does not fade range edges,
target a midpoint, target VWAP, target value area, trade previous-day range
reversion, select range-router states, reuse the range optimization workbench,
or revive failed post-compression or clean-breakout paths.

It also does not use EMA trend stacks, EMA slopes, EMA pullback bands, the failed
trend-pullback continuation trigger, cross-asset context, derivatives source
data, the `btc_15m_basis_discount_no_trade_veto_v1` veto, volume filters,
volatility filters, weekday filters, or post-result side selection. The first
test is only a BTCUSDT OHLCV opening-range expansion baseline.

## Candidate

Candidate id:

```text
btc_15m_session_opening_range_expansion_v1
```

### Hypothesis

After the fixed `13:30 UTC` liquidity handoff window opens, a closed `15m`
acceptance outside the first-hour opening range can continue far enough away
from that range to pay standard futures fees and slippage.

### Source And Timeframe

- Source: current accepted BTCUSDT Binance USDT-M futures `5m` CSV.
- Decision timeframe: exact closed UTC `15m` bars resampled from the `5m` source.
- Session anchor: every UTC calendar date at `13:30:00Z`.
- Opening-range length: `60` minutes, using the four closed `15m` bars with open
  times `13:30`, `13:45`, `14:00`, and `14:15` UTC.
- Expansion window: closed decision bars with open times in `[14:30, 17:30) UTC`
  on the same UTC date.
- Entry timing: next `15m` open after the closed decision candle.
- Position model: one open position max, and at most one signal per UTC date.

The `13:30 UTC` anchor is fixed in UTC for every date. Do not daylight-saving
shift it, compare it against alternate anchors, or mine a replacement session
after seeing results.

### Fixed Calculations

All calculations use resampled closed UTC `15m` candles:

- Opening-range high: maximum high of the four opening-range bars.
- Opening-range low: minimum low of the four opening-range bars.
- Opening-range width: opening-range high minus opening-range low.
- `ATR(14)` using true range with Wilder-style smoothing.

Skip a UTC date until all four opening-range bars and `ATR(14)` are available.
Skip a date if the opening-range width is non-positive.

### Fixed Entry Rule

For each UTC date, evaluate the expansion-window decision candles in timestamp
order after the opening range is complete. If a position is already open, skip
signals until the position is closed. If both long and short conditions somehow
qualify on the same candle, skip the candle as ambiguous. Emit only the first
qualified signal for the UTC date.

Long setup:

1. On the closed decision candle `d`, require
   `close[d] >= opening_range_high + 0.10 * ATR(14)[d]`.
2. Enter long on `open[d+1]`.

Short setup:

1. On the closed decision candle `d`, require
   `close[d] <= opening_range_low - 0.10 * ATR(14)[d]`.
2. Enter short on `open[d+1]`.

### Fixed Risk And Exits

Use the same fixed risk model for both sides:

- Long stop: `opening_range_low - 0.10 * ATR(14)[d]`.
- Short stop: `opening_range_high + 0.10 * ATR(14)[d]`.
- Skip the signal if stop distance from the next-bar-open entry is non-positive.
- Target: `1.5R` from entry using the initial stop distance.
- Time stop: close after `24` closed `15m` bars after entry.
- No trailing stop, break-even move, scale-in, scale-out, re-entry while a
  position is open, martingale, or averaging down.
- Sizing: `1%` risk at stop, capped at `1x` notional.
- Costs: `0.0004` fee per side and `0.000116` slippage per side.

### Expected Output Path And Artifacts

If later approved for implementation, use this output path:

```text
results/backtest-first-btc-15m-session-opening-range-expansion-v1/
```

Expected artifacts:

- `source_manifest.json`;
- `summary.json`, `summary.csv`, and `trades.json`;
- candidate-specific coverage, sources, session-ranges, signals, skips, trades,
  summary, and falsification CSV/JSON files.

### Pass/Fail Gates

This fixed baseline may request review only if all gates below pass:

- source validation passes with the accepted BTCUSDT Binance USDT-M futures `5m`
  source contract;
- `15m` resampling uses only complete closed UTC bars and records coverage;
- the `13:30 UTC` anchor, `60` minute opening range, `[14:30, 17:30) UTC`
  expansion window, `0.10 * ATR(14)` acceptance buffer, stop, target, and time
  stop are used exactly as predeclared;
- no leakage is found in session construction, opening-range state, ATR,
  trigger, stop, target, sizing, or selection;
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
required, leakage is found, the result depends on a mined session anchor, or the
result needs retuning to survive.

### No-Rescue Boundaries

If this fixed baseline fails, do not rescue it with alternate UTC anchors,
opening-range lengths, expansion windows, acceptance buffers, ATR windows, stop
buffers, target R values, time stops, one-trade-per-day changes, side selection,
weekday filters, volume filters, volatility filters, derivatives-veto
interaction, source expansion, replay, walk-forward, or optimizer grids.

Do not rebrand closed range-reversion, midpoint, edge-fade, previous-day range,
value-area, range-optimization, post-compression, clean-breakout,
router/rotation, or trend-pullback branches as opening-range expansion.

Record the result and either close the candidate or move to a materially
different approved lane.

## Next Approval Gate

The next bounded gate is explicit operator approval to implement and run this
one fixed offline baseline. Without that approval, no Go code, CLI flag,
generated result, backtest, optimizer, source expansion, derivatives-veto
interaction, exchange API work, credential, deployment, paper/testnet/live flow,
martingale, averaging down, two-exchange logic, or promotion is authorized.
