# Backtest-First Candidate Packet

Date: 2026-06-30

## Verdict

Stop state:
`backtest_first_candidate_packet_selected_value_reversion_baseline_needs_implementation_approval`.

Select exactly one first baseline candidate for the next implementation task:

```text
btc_5m_rolling_value_area_reversion_v1
```

This is a docs-only candidate packet. It lists simple, materially different
BTCUSDT range-entry hypotheses, rejects closed-family rescues, and selects the
simplest fixed baseline to test first. It does not authorize Go code, CLI flags,
generated outputs, source downloads, optimizer grids, replay, walk-forward,
derivatives-veto interaction, paper/testnet/live paths, exchange API work,
credentials, deploy files, martingale, averaging down, two-exchange logic, or
promotion.

## Common Source Contract

Unless a later source-scope review changes scope, every candidate in this packet
uses the current default source:

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

All rules must use confirmed closed candles only. Entries happen on the next bar
open. Same-bar stop/target ambiguity is stop-first. Costs use `0.0004` fee per
side and `0.000116` slippage per side unless a later fixed baseline packet
explicitly changes them before implementation.

## Closed-Family Rejection Filter

The candidates below may not reopen, rename, retune, or rescue these reviewed
families:

- failed post-compression directional expansion;
- structured compression baseline/optimization/replay/walk-forward;
- breakout-retest acceptance;
- clean breakout continuation;
- hold-inside/midline;
- impulse absorption;
- higher-timeframe nested range rotation;
- range occupancy rotation;
- range router rotation;
- BTC-regime plus ETH/SOL context;
- derivatives no-trade veto or derivatives-context rotation.

No candidate may use adjacent-cell P&L selection, post-result filter additions,
exit retuning after a failed backtest, derivatives-veto interaction, optimizer
selection, replay, or walk-forward as part of its first baseline.

## Candidate 1 — Rolling Value-Area Reversion

Candidate id:

```text
btc_5m_rolling_value_area_reversion_v1
```

### Hypothesis

When BTCUSDT is inside a broad rolling intraday value area, excursions into the
outer value zones often mean-revert toward the rolling volume-weighted value
anchor before a true directional break develops.

### Material Difference From Closed Failures

This is not a compression breakout, post-compression expansion, breakout-retest,
router rotation, occupancy rotation, midline/hold-inside replay, or derivatives
veto. It uses a rolling value anchor inside a local `5m` range and asks whether
outer-zone value reversion has enough gross edge before any optimizer or veto is
considered.

### Source And Timeframe

- Source: current accepted BTCUSDT Binance USDT-M futures `5m` CSV.
- Decision timeframe: native closed `5m` candles.
- Entry timing: next `5m` open.

### Fixed Entry Rule

For each closed `5m` decision candle `d`:

1. Use the prior `288` closed `5m` bars `[d-288, d-1]` as the rolling value
   window.
2. Compute `range_high`, `range_low`, and `range_width = range_high - range_low`.
3. Compute rolling VWAP from typical price `(high + low + close) / 3` weighted by
   volume over `[d-288, d-1]`; skip if total volume is zero.
4. Compute `ATR(14)[d-1]` on native `5m` candles; skip if missing or non-positive.
5. Require `range_width >= 6 * ATR(14)[d-1]` so the target distance is not too
   small relative to noise.
6. Long candidate: `close[d]` is in the lower `20%` of the prior range and at
   least `0.15 * range_width` below rolling VWAP.
7. Short candidate: `close[d]` is in the upper `20%` of the prior range and at
   least `0.15 * range_width` above rolling VWAP.
8. Do not enter if a position is already open.

### Fixed Risk And Exits

- Long entry: next `5m` open after signal.
- Long stop: `min(low[d], range_low) - 0.25 * ATR(14)[d-1]`.
- Long target: rolling VWAP from the decision window.
- Short entry/stop/target are symmetric.
- Time stop: close after `36` closed `5m` bars after entry.
- Sizing: `1%` risk at stop, capped at `1x` notional.
- Costs: `0.0004` fee per side and `0.000116` slippage per side.

### Expected Output Path And Artifacts

If later approved for implementation, use:

```text
results/backtest-first-btc-5m-rolling-value-area-reversion-v1/
```

Expected artifacts:

- `source_manifest.json`;
- `summary.json`, `summary.csv`, `trades.json`;
- strategy-specific sources, signals, skips, trades, summary, and falsification
  CSV/JSON files.

### Pass/Fail Gates

A first fixed baseline may request review only if:

- source validation passes with the accepted BTCUSDT futures `5m` contract;
- no leakage is found;
- full-sample executed trades are at least `120`;
- every primary split has at least `25` executed trades;
- full-sample gross P&L is positive;
- `2023_2024_oos` and `2025_2026_recent` gross P&L are not clearly negative;
- full-sample net after fees/slippage is positive;
- results do not depend only on `2021_2022_stress`;
- drawdown is acceptable versus return;
- long and short sides are reported separately.

Fail fast if gross edge is absent, costs kill the edge, OOS/recent splits fail,
trade count is too small, drawdown is unacceptable, or the result requires
retuning to survive.

### No-Rescue Boundaries

If this fixed baseline fails, do not rescue it with alternate VWAP windows,
outer-zone percentages, target changes, time-stop changes, side selection,
volume filters, derivatives-veto interaction, replay, walk-forward, or optimizer
grids. Record the result and move to a materially different candidate.

## Candidate 2 — Previous-Day Range Reversion

Candidate id:

```text
btc_15m_previous_day_range_reversion_v1
```

### Hypothesis

When the current UTC day remains inside the previous UTC day's high-low range,
outer-decile tests of that previous-day range may revert toward the previous-day
midpoint before a true daily expansion starts.

### Material Difference From Closed Failures

This is a fixed calendar-range value test, not a compression pocket, router
rotation, clean breakout, breakout-retest, nested higher-timeframe rotation, or
derivatives-context filter. It uses the previous day's range as the only range
anchor and does not mine session cohorts.

### Source And Timeframe

- Source: current accepted BTCUSDT Binance USDT-M futures `5m` CSV resampled to
  exact closed UTC `15m` bars.
- Decision timeframe: closed UTC `15m` candles.
- Entry timing: next `15m` open.

### Fixed Entry Rule

For each closed `15m` decision candle `d`:

1. Build the prior UTC day's `high`, `low`, and midpoint from complete `15m`
   candles.
2. Skip the current day if any current-day candle before `d` has closed outside
   the prior day's high-low range.
3. Long candidate: `close[d]` is inside the prior-day range and in its lower
   `10%`.
4. Short candidate: `close[d]` is inside the prior-day range and in its upper
   `10%`.
5. Do not enter if a position is already open.

### Fixed Risk And Exits

- Stop: beyond the prior-day range extreme by `0.25 * ATR(14)[d-1]`.
- Target: prior-day midpoint.
- Time stop: `24` closed `15m` bars.
- Sizing/costs: `1%` risk, `1x` notional cap, `0.0004` fee per side,
  `0.000116` slippage per side.

### Expected Output Path And Artifacts

```text
results/backtest-first-btc-15m-previous-day-range-reversion-v1/
```

Artifacts should mirror the common backtest outputs plus strategy-specific
sources, signals, skips, trades, summary, and falsification files.

### Pass/Fail Gates

Use the same first-baseline gates as Candidate 1, with exact UTC daily source
coverage and resample validation added.

### No-Rescue Boundaries

If this fails, do not rescue with alternative UTC session definitions, previous
`2`/`3` day windows, changed outer-decile thresholds, derivatives context, or
calendar/session mining.

## Candidate 3 — Range-Edge Exhaustion Fade

Candidate id:

```text
btc_15m_range_edge_exhaustion_fade_v1
```

### Hypothesis

A fast move from the middle of a local range into a range edge can become
exhausted before breakout. If the final approach candle shows weak incremental
progress while price remains inside the prior range, fading back toward the range
midpoint may have a tradable gross edge.

### Material Difference From Closed Failures

This is not a clean breakout, breakout-retest, post-compression expansion,
occupancy rotation, or router-gated boundary reclaim. The first baseline uses a
single deceleration condition at a local range edge, not a route-selection state,
zero-trade cohort ranking, or parameter grid.

### Source And Timeframe

- Source: current accepted BTCUSDT Binance USDT-M futures `5m` CSV resampled to
  exact closed UTC `15m` bars.
- Decision timeframe: closed UTC `15m` candles.
- Entry timing: next `15m` open.

### Fixed Entry Rule

For each closed `15m` decision candle `d`:

1. Use prior `96` closed `15m` bars `[d-96, d-1]` as the local range.
2. Compute range high, range low, midpoint, and `ATR(14)[d-1]`.
3. Require `close[d-3]` to be between the lower and upper `40%` of the range,
   so the move begins near the interior.
4. Long fade candidate: the last three closes move downward toward the lower
   `15%` of the range, `close[d]` remains above `range_low`, and the final close
   progress from `d-1` to `d` is less than `0.35 * ATR(14)[d-1]` in the trend
   direction.
5. Short fade candidate: symmetric upward move into the upper `15%` of the
   range, still below `range_high`, with weak final progress.
6. Do not enter if a position is already open.

### Fixed Risk And Exits

- Stop: beyond the local range edge by `0.25 * ATR(14)[d-1]`.
- Target: local range midpoint.
- Time stop: `16` closed `15m` bars.
- Sizing/costs: `1%` risk, `1x` notional cap, `0.0004` fee per side,
  `0.000116` slippage per side.

### Expected Output Path And Artifacts

```text
results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/
```

Artifacts should mirror the common backtest outputs plus strategy-specific
sources, signals, skips, trades, summary, and falsification files.

### Pass/Fail Gates

Use the same first-baseline gates as Candidate 1, with exact `15m` resample
validation added.

### No-Rescue Boundaries

If this fails, do not rescue with alternate range windows, progress thresholds,
edge zones, midpoint variants, added volume filters, derivatives context, or
optimizer grids.

## Rejected Lookalikes

These ideas are intentionally not selected for the first backtest-first packet:

- post-compression breakout variants;
- clean breakout continuation variants;
- breakout-retest acceptance variants;
- router-gated boundary reclaim variants;
- occupancy-rotation variants;
- midline/hold-inside variants;
- derivatives-veto skip/retain interaction.

They are too close to reviewed closed paths or require a passed independent entry
baseline before interaction testing.

## Selected First Baseline

Select Candidate 1:

```text
btc_5m_rolling_value_area_reversion_v1
```

Reasons:

- It is the simplest native-`5m` implementation.
- It uses only the already accepted BTCUSDT futures source.
- It does not require resampling before the first baseline.
- It is materially different from the latest failed post-compression expansion
  path.
- It tests a direct range value-reversion premise before adding context, vetoes,
  optimizers, replay, or walk-forward.

## Next Gate

The next allowed task is explicit user approval to implement exactly one fixed
offline baseline backtest for:

```text
btc_5m_rolling_value_area_reversion_v1
```

Until that approval is given, do not add Go code, CLI flags, generated outputs,
source downloads, optimizer grids, replay, walk-forward, derivatives-veto
interaction, paper/testnet/live paths, exchange API work, credentials, deploy
files, martingale, averaging down, two-exchange logic, or promotion.
