# Futures Higher-Timeframe Nested Range Rotation Premise Spec

Date: 2026-06-26

## Verdict

Stop state:
`higher_tf_nested_range_rotation_premise_ready_for_audit`.

This is a review-only premise spec. It authorizes a future non-trading audit
brief only. It does not add code, CLI flags, generated results, entries, exits,
scoring, sizing, optimization, replay, walk-forward, paper/testnet/live wiring,
exchange API use, credentials, deploy files, data downloads, broad symbol
mining, martingale, averaging down, or two-exchange logic.

The next premise is materially different from the closed families because it
does not enter on a parent-range breakout, does not retest a broken parent
boundary, does not optimize structured-compression parameters, and does not
reuse the failed BTCUSDT `5m` hold-inside/midline template. It asks whether a
completed lower-timeframe range nested inside a higher-timeframe range can
identify internal rotation pressure before the parent range breaks.

## Exclusion Context

The current implementation chain is stopped:

- the frozen ETH/SOL `4h` structured-compression stream failed walk-forward
  robustness and is exclusion evidence for strategy packaging;
- the bounded breakout-retest/acceptance baseline failed after costs on
  BTCUSDT, ETHUSDT, and SOLUSDT;
- the BTCUSDT clean-breakout baseline failed after costs;
- `5m` hold-inside/midline, impulse absorption, SR timing, and legacy
  compression families are closed or diagnostic only.

Reusable infrastructure remains valuable: accepted futures source guards,
closed UTC resampling, mature range episode generation, audit-only artifact
patterns, split-stability review, and closed-candle finality.

## Source Contract

Use only the already specified BTCUSDT higher-timeframe source lane from
`docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.

Parent source:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol / interval | `BTCUSDT` / `5m` |
| Loaded candles | `573,984` |
| First parent open | `2021-01-01T00:00:00Z` |
| Last parent open | `2026-06-16T23:55:00Z` |
| Manifest status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

Required closed UTC resamples for the next audit:

| Interval | Rows | First open | Last open |
| --- | ---: | --- | --- |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` |

The `15m` resample is not part of the first audit. It may remain available as
future diagnostic context only after this premise is reviewed.

## Premise

Premise:

> When a mature `4h` BTCUSDT range contains a completed `1h` child range in one
> side of the parent range, the first closed `1h` displacement from the child
> range toward the parent interior may predict internal rotation to the parent
> midpoint or far quartile better than immediate invalidation back through the
> child range.

This is a non-trading hypothesis. The audit should measure event outcomes only.
It must not create orders, strategy signals, position sizing, or P&L rows.

## Candidate Event

Use the existing mature range detector backbone without tuning:

- detector profile: `p30_c12_bollinger_on_adx_off`;
- parent timeframe: `4h`;
- child timeframe: `1h`;
- parent and child candles must be complete closed UTC bars;
- every event must be knowable from the event candle close.

Parent range eligibility:

- a `4h` range is mature after at least `12` consecutive detector-qualified
  bars;
- freeze parent high, low, midpoint, upper quartile, and lower quartile at the
  first mature `4h` close;
- parent width must be positive;
- candidate child ranges are ignored after a closed `4h` candle closes outside
  the frozen parent range.

Child range eligibility:

- a `1h` range is mature after at least `12` consecutive detector-qualified
  bars;
- the child range high and low must be fully inside the frozen parent range;
- child width must be positive and no more than `40%` of parent width;
- the child midpoint must be in the parent lower half for upside rotation
  candidates or in the parent upper half for downside rotation candidates;
- a child range may emit at most one event.

Event definitions:

- `nested_rotation_up`: first closed `1h` candle after child maturity closes
  above the child high, remains inside the frozen parent range, closes below
  the parent midpoint, and starts from a child midpoint below the parent
  midpoint.
- `nested_rotation_down`: first closed `1h` candle after child maturity closes
  below the child low, remains inside the frozen parent range, closes above
  the parent midpoint, and starts from a child midpoint above the parent
  midpoint.

Skip and count, at minimum:

- missing parent or child resample coverage;
- non-positive parent or child width;
- child range not fully inside parent range;
- child width above `40%` of parent width;
- child midpoint not in the required parent half;
- event candle outside the parent range;
- event candle already beyond the parent midpoint;
- duplicate child event.

## Outcome Labels

Evaluate outcomes over `24` closed `1h` bars after the event candle. This is an
audit horizon, not a max hold for trading.

For `nested_rotation_up`:

- favorable midpoint: high reaches the parent midpoint before child-low
  invalidation;
- favorable far quartile: high reaches the parent upper quartile before
  child-low invalidation;
- adverse child invalidation: low reaches the child low before favorable
  midpoint;
- adverse parent invalidation: close below the frozen parent low before
  favorable midpoint;
- no resolution: neither favorable midpoint nor adverse child invalidation
  occurs inside `24` bars.

For `nested_rotation_down`, mirror the labels:

- favorable midpoint: low reaches the parent midpoint before child-high
  invalidation;
- favorable far quartile: low reaches the parent lower quartile before
  child-high invalidation;
- adverse child invalidation: high reaches the child high before favorable
  midpoint;
- adverse parent invalidation: close above the frozen parent high before
  favorable midpoint;
- no resolution: neither favorable midpoint nor adverse child invalidation
  occurs inside `24` bars.

Also record quick invalidation when adverse child invalidation occurs within
`6` closed `1h` bars.

## Review Gate

The next audit may pass only if all of these are true:

- source and `1h`/`4h` resample validation pass;
- at least `100` full-sample events exist;
- every period split has at least `25` events;
- both sides have at least `25` full-sample events, or the weaker side is
  explicitly marked as a caveat and excluded from any later baseline brief;
- favorable midpoint rate is greater than adverse child invalidation rate in
  every period split;
- favorable midpoint rate is greater than quick invalidation rate in every
  period split;
- favorable far-quartile rate is non-trivial in every split, with no split
  below `20%`;
- average favorable excursion to midpoint is greater than average adverse
  excursion to child invalidation in every split;
- the result is not carried entirely by one short historical regime.

The audit must fail if the only way to pass is to change the `40%` child-width
gate, the `24` bar outcome horizon, the `6` bar quick-invalidation horizon, or
the split gates after seeing results.

## Not A Closed-Family Retry

This premise is allowed only because it is structurally different from the
closed result set:

- not clean breakout continuation: parent-range boundary breaks are excluded
  from event formation;
- not breakout retest/acceptance: no broken parent boundary is retested or
  accepted;
- not structured compression: there is no confirmation-window, target, stop,
  max-hold, symbol-set, or training-gate grid;
- not `5m` hold-inside/midline: the event is a nested `1h` range inside a
  frozen `4h` parent range; the parent midpoint is an outcome marker, not an
  entry trigger;
- not mature balance persistence: the child range must displace toward the
  parent interior, so the audit is about nested rotation, not passive
  persistence.

If implementation work drifts into any closed family, stop with
`higher_tf_nested_range_rotation_audit_rejected_as_closed_family_reslice`.

## Next Audit Brief

The next implementation may add a non-trading audit behind:

```text
-futures-higher-tf-nested-range-rotation-audit
```

Expected artifacts under
`results/futures-higher-tf-nested-range-rotation-audit/`:

- `futures_higher_tf_nested_range_rotation_sources.csv/json`
- `futures_higher_tf_nested_range_rotation_coverage.csv/json`
- `futures_higher_tf_nested_range_rotation_parent_ranges.csv/json`
- `futures_higher_tf_nested_range_rotation_child_ranges.csv/json`
- `futures_higher_tf_nested_range_rotation_events.csv/json`
- `futures_higher_tf_nested_range_rotation_summary.csv/json`
- common `source_manifest.json`, `summary.csv/json`, and `trades.json`

Common `summary.*` and `trades.json` must remain zero-trade compatibility
outputs.

Audit stop states:

- `higher_tf_nested_range_rotation_audit_source_gap`
- `higher_tf_nested_range_rotation_audit_no_candidate_events`
- `higher_tf_nested_range_rotation_audit_rejected_as_closed_family_reslice`
- `higher_tf_nested_range_rotation_audit_failed_no_baseline`
- `higher_tf_nested_range_rotation_audit_ready_for_baseline_brief`

Only `higher_tf_nested_range_rotation_audit_ready_for_baseline_brief` may
authorize a later fixed-rule baseline brief. It still would not authorize
optimization, live/paper/testnet, exchange API use, deployment, data downloads,
symbol expansion, martingale, averaging down, or two-exchange work.
