# Futures Range-First Strategy Construction V1 Spec

Date: 2026-06-26

## Verdict

Stop state:
`range_first_strategy_v1_spec_ready_for_optimizer_implementation`.

This is a documentation-only strategy construction spec. It selects the first
fresh range-derived BTCUSDT grammar for a later bounded offline
optimizer/backtester, but it does not implement a strategy, optimizer, replay,
walk-forward run, CLI flag, tests, generated result directory, artifact writer,
data download, paper/testnet/live path, exchange API, credentials, deploy file,
martingale, averaging down, or two-exchange logic.

The selected V1 grammar is:
`range_occupancy_rotation_v1`.

In plain language: build a rolling closed-candle range envelope, require recent
closes to cluster persistently in one outer part of the range, then enter only
after a closed candle recaptures an interior line back toward range value. The
trade is a rotation away from close-location imbalance inside a still-contained
range, not a breakout, retest, compression expansion, impulse absorption, or
nested parent/child rotation.

## Source Contract

The first construction source remains BTCUSDT Binance USDT-M futures `5m`:

| Field | Contract |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Market | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Parent interval | `5m` |
| Loaded candles | `573,984` |
| Coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Accepted facts | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

The later implementation must reject source gaps, duplicates, irregular `5m`
cadence, non-positive or non-finite OHLC prices, negative volume, invalid
high/low containment, non-BTCUSDT identity, non-`5m` identity, and spot or
comparison-only paths. Spot comparison rows cannot satisfy any V1 gate.

## Timeframes

V1 may derive only these closed UTC bars from the accepted `5m` parent source:

| Role | Timeframes | Use |
| --- | --- | --- |
| Parent source | `5m` | Source validation and optional execution alignment only |
| Signal bars | `15m`, `1h` | Range envelope, signal, entry, stop, target, max-hold accounting |

No `4h` parent/child context is part of V1. That avoids re-slicing the failed
higher-timeframe nested range-rotation audit. ETH/SOL and broader symbols are
not part of V1; they can be considered only by a later optional transfer review
after BTCUSDT evidence exists.

Expected closed UTC resample validation:

| Timeframe | Expected rows | Expected first open | Expected last open |
| --- | ---: | --- | --- |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |

Every resample row must have `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Exclusion Boundary

This spec would have stopped with
`range_first_strategy_v1_spec_rejected_closed_family_reslice` if the selected
grammar drifted into any reviewed closed family. A later implementation must
stop with `range_first_strategy_v1_rejected_closed_family_reslice` if the code
or config drifts into any reviewed closed family:

- structured-compression expansion or ETH/SOL authority replay;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline touch, close-back, or first-midline-event entries;
- impulse absorption after abnormal candles;
- higher-timeframe nested range rotation;
- old spot-only SR rejection, confirmation, false-break, or compression
  promotion evidence.

V1 may reuse range envelopes, split metrics, closed-candle semantics,
next-bar-open execution, costs, one-position max, and stop-first handling. It
must not rescue failed entry rules by renaming them.

## Feature Primitives

For a signal timeframe `T` and closed signal candle index `i`, V1 may use only
features known at the close of candle `i`:

- rolling range high and low from the previous `lookback_hours` of closed bars,
  ending at `i-1`;
- range width, midpoint, quartiles, and thirds derived from that rolling high
  and low;
- range width as a fraction of the signal candle close;
- recent close-location occupancy counts over the last `occupancy_window`
  closed signal bars ending at `i`;
- inside-range close containment for that recent occupancy window;
- signal candle close recapture of a declared interior line;
- side label: `long` from lower-range occupancy imbalance, `short` from
  upper-range occupancy imbalance;
- split label by trade close time;
- standard engine trade metrics, side metrics, and drawdown metrics.

V1 must not use future candles, wick-only rejection labels, first-break labels,
midline-touch labels, parent/child nesting labels, external strategy scores, or
old `binance-bot` strategy/scoring/live code.

## Range Envelope

For each signal candle `i`:

1. Convert `lookback_hours` to signal bars for timeframe `T`.
2. Build the envelope from bars `[i-lookback_bars, i-1]`.
3. Freeze:
   - `range_high`;
   - `range_low`;
   - `range_width = range_high - range_low`;
   - `range_mid = range_low + 0.50 * range_width`;
   - `range_q1 = range_low + 0.25 * range_width`;
   - `range_q3 = range_low + 0.75 * range_width`;
   - `lower_recapture = range_low + recapture_level * range_width`;
   - `upper_recapture = range_high - recapture_level * range_width`.
4. Reject the envelope when width is non-positive, prices are non-positive,
   there is insufficient history, or `range_width / close_i` exceeds
   `max_width_pct`.

The signal candle itself is not allowed to define the high/low envelope. This
keeps the range known before the decision candle closes.

## Entry Grammar

All entries are next signal-timeframe bar open. The engine remains one-position
max; overlapping signals are skipped or recorded as already-in-position.

### Long Setup

A long signal is valid when all are true:

1. The envelope is valid.
2. Every close in the last `occupancy_window` signal bars is inside
   `[range_low, range_high]`.
3. At least `occupancy_min_fraction` of those closes are in the lower
   occupancy zone:
   `close <= range_low + occupancy_zone_level * range_width`.
4. The signal candle closes at or above `lower_recapture`.
5. The signal candle close remains below `range_mid`.
6. The next signal bar open exists and produces valid entry geometry.

### Short Setup

A short signal mirrors the long setup:

1. The envelope is valid.
2. Every close in the last `occupancy_window` signal bars is inside
   `[range_low, range_high]`.
3. At least `occupancy_min_fraction` of those closes are in the upper
   occupancy zone:
   `close >= range_high - occupancy_zone_level * range_width`.
4. The signal candle closes at or below `upper_recapture`.
5. The signal candle close remains above `range_mid`.
6. The next signal bar open exists and produces valid entry geometry.

If both sides trigger on the same signal candle, skip the signal as
`ambiguous_dual_side_signal`. Do not choose a side from later data.

## Stop, Target, And Max Hold

Execution uses the existing offline engine contract:

- next signal-timeframe bar open entry;
- one open position max;
- existing fees, slippage, 1% risk sizing, and 1x notional cap;
- stop-first handling when stop and target are both touched in the same bar.

Long trade geometry:

- stop: `range_low - stop_buffer_width * range_width`;
- target: `range_low + target_level * range_width`;
- target must be above slipped entry;
- stop must be below slipped entry.

Short trade geometry:

- stop: `range_high + stop_buffer_width * range_width`;
- target: `range_high - target_level * range_width`;
- target must be below slipped entry;
- stop must be above slipped entry.

Max hold is `max_hold_bars` signal bars after entry. If neither stop nor target
is hit before then, exit at the close of the max-hold bar using the existing
time-exit behavior.

## Invalid Geometry And Skip Reasons

The later implementation must record skipped signal reasons at minimum:

- `insufficient_history`;
- `source_or_resample_rejected`;
- `non_positive_price`;
- `non_positive_range_width`;
- `range_too_wide`;
- `recent_close_outside_range`;
- `occupancy_threshold_not_met`;
- `recapture_not_confirmed`;
- `signal_crossed_midpoint`;
- `missing_entry_bar`;
- `ambiguous_dual_side_signal`;
- `entry_stop_target_invalid`;
- `already_in_position`.

Skipped rows belong in V1-specific artifacts, not in the normal common
`trades.json`.

## Parameter Grid

The default V1 optimization grid is exactly `1,152` declared configs:

| Parameter | Values |
| --- | --- |
| `signal_timeframe` | `15m`, `1h` |
| `lookback_hours` | `24`, `48`, `72` |
| `max_width_pct` | `0.020`, `0.035` |
| `occupancy_window` | `8`, `12` signal bars |
| `occupancy_zone_level` | `0.25` |
| `occupancy_min_fraction` | `0.60`, `0.70` |
| `recapture_level` | `0.25`, `0.33` |
| `target_level` | `0.50`, `0.66` |
| `max_hold_bars` | `8`, `12`, `24` signal bars |
| `stop_buffer_width` | `0.00`, `0.05` |
| `side_mode` | `all` |

No implementation may add grid dimensions, broaden values, switch symbols, or
change gates to rescue a weak result without a new reviewed spec.

## Fixed Baseline

The later implementation must always evaluate this fixed baseline row before
ranking the grid:

```text
range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005
```

Baseline parameters:

| Parameter | Value |
| --- | --- |
| `signal_timeframe` | `1h` |
| `lookback_hours` | `48` |
| `max_width_pct` | `0.035` |
| `occupancy_window` | `12` |
| `occupancy_zone_level` | `0.25` |
| `occupancy_min_fraction` | `0.60` |
| `recapture_level` | `0.33` |
| `target_level` | `0.66` |
| `max_hold_bars` | `12` |
| `stop_buffer_width` | `0.05` |
| `side_mode` | `all` |

This fixed baseline is not selected from prior results. It exists so the review
can separate a human-declared first template from optimized configs.

Common `summary.csv/json` and `trades.json` should describe the fixed baseline
only. Grid, ranking, selected-config, and comparison rows belong only in
V1-specific artifacts until a later fixed replay is approved.

## Optimization Behavior

The optimizer may rank only declared grid rows. It must compute all split and
side summaries for every config, but it must select candidates using training
data only.

Initial optimization split:

| Role | Split |
| --- | --- |
| Train/rank | `2021_2022_stress` |
| OOS gate | `2023_2024_oos` |
| Recent gate | `2025_2026_recent` |
| Full review | `full_2021_2026` |

A config is selectable only when all are true:

- train trades >= `100`;
- OOS trades >= `50`;
- recent trades >= `25`;
- full trades >= `200`;
- train net P&L after costs > `0`;
- OOS net P&L after costs > `0`;
- recent net P&L after costs > `0`;
- full net P&L after costs > `0`;
- train PF >= `1.20`;
- OOS PF >= `1.05`;
- recent PF >= `1.05`;
- full PF >= `1.15`;
- train net-to-drawdown ratio >= `1.00` when drawdown is positive;
- full net-to-drawdown ratio >= `1.00` when drawdown is positive;
- no side with at least `25` full-period trades loses money in both OOS and
  recent splits;
- no side accounts for more than `75%` of full-period net P&L unless the
  weaker side is still positive after costs.

If no config is selectable, stop with
`range_first_strategy_v1_optimizer_failed_no_replay`.

## Ranking Score

Rank selectable configs by training rows only:

```text
pf_component = min(train_profit_factor, 3.0) - 1.0
trade_component = min(train_trades, 400) / 400.0
drawdown_component = train_net_pnl / max(train_max_drawdown, 1.0)
rank_score = drawdown_component + pf_component + trade_component - caveat_penalty
```

`caveat_penalty` is:

- `0.50` for side concentration above `60%` of full-period net P&L;
- `0.50` for either side having fewer than `50` full-period trades;
- additive when both caveats are present.

Tie-breaks, in order:

1. higher train net P&L;
2. higher train PF;
3. lower train max drawdown;
4. higher full-period trade count;
5. `1h` before `15m`;
6. shorter max hold;
7. lexicographic config ID.

OOS and recent rows are gates and review evidence, not ranking inputs.

## Walk-Forward Expectations

A later walk-forward brief, if earned by the first optimizer review, should use
these folds:

| Fold | Train | Test |
| --- | --- | --- |
| `wf_2021_2022_train__2023_2024_test` | `2021_2022_stress` | `2023_2024_oos` |
| `wf_2021_2024_train__2025_2026_test` | `2021_2022_stress+2023_2024_oos` | `2025_2026_recent` |
| `wf_2023_2024_train__2025_2026_test` | `2023_2024_oos` | `2025_2026_recent` |

Walk-forward selection must use train rows only, require at least `100` train
trades for single-split trains and `150` aggregate train trades for multi-split
trains, and require positive net P&L plus PF >= `1.20` in each training
segment. A package review is not authorized unless at least two of three fold
tests pass and no fold requires changing the declared grid.

## Expected Implementation Flag

The later bounded offline optimizer/backtester should be wired behind:

```text
-futures-range-first-occupancy-rotation-v1-optimization
```

The flag must reject spot or comparison-only sources and must reject
combinations with other trade-producing prototype, baseline, optimization,
replay, or walk-forward flags.

## Expected Result Directory And Artifacts

Default result directory:

```text
results/futures-range-first-occupancy-rotation-v1-optimization/
```

Common outputs:

- `source_manifest.json`;
- `summary.csv`;
- `summary.json`;
- `trades.json`.

V1-specific outputs:

- `futures_range_first_occupancy_rotation_v1_sources.csv/json`;
- `futures_range_first_occupancy_rotation_v1_coverage.csv/json`;
- `futures_range_first_occupancy_rotation_v1_grid.csv/json`;
- `futures_range_first_occupancy_rotation_v1_baseline.csv/json`;
- `futures_range_first_occupancy_rotation_v1_signals.csv/json`;
- `futures_range_first_occupancy_rotation_v1_trades.csv/json`;
- `futures_range_first_occupancy_rotation_v1_summary.csv/json`;
- `futures_range_first_occupancy_rotation_v1_rankings.csv/json`;
- `futures_range_first_occupancy_rotation_v1_selection.csv/json`;
- `futures_range_first_occupancy_rotation_v1_skips.csv/json`.

The V1 `signals` and `trades` artifacts should include the fixed baseline and
the selected top-ranked config when one exists. The `summary` artifact should
include all grid config split/side rows. The `selection` artifact should compare
the fixed baseline to the selected optimized config and include pass/fail
reasons.

## Stop States

Spec stop states:

- `range_first_strategy_v1_spec_ready_for_optimizer_implementation`;
- `range_first_strategy_v1_spec_needs_user_premise_or_scope_input`;
- `range_first_strategy_v1_spec_rejected_closed_family_reslice`.

Later implementation stop states:

- `range_first_strategy_v1_source_gap`;
- `range_first_strategy_v1_no_valid_signals`;
- `range_first_strategy_v1_baseline_failed_grid_still_reviewed`;
- `range_first_strategy_v1_optimizer_failed_no_replay`;
- `range_first_strategy_v1_passed_needs_fixed_replay_spec`;
- `range_first_strategy_v1_rejected_closed_family_reslice`.

Passing the optimizer review would authorize only a fixed replay spec for the
selected config. It would not authorize strategy packaging, walk-forward,
paper/testnet/live wiring, exchange API use, deployment, source expansion,
symbol expansion, martingale, averaging down, or two-exchange logic.
