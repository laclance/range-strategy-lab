# Futures Range Router Rotation Premise Spec

Date: 2026-06-27

## Verdict

Stop state:
`range_router_rotation_premise_spec_ready_for_non_trading_audit`.

This is a documentation-only premise spec. It authorizes a future non-trading
audit only. It does not add code, CLI flags, generated results, entries, exits,
P&L strategy backtests, optimizer grids, replay, walk-forward logic, strategy
packages, paper/testnet/live paths, exchange API use, credentials, deploy
files, source expansion, symbol expansion, broad mining, martingale, averaging
down, or two-exchange logic.

The premise starts from the passed router cohort:

```text
range_context_router_v1|15m|h24|tradable_rotation
```

The router is closed-candle context only. It is not an entry signal, and the
`1,299` `tradable_rotation` router rows must not be converted directly into
trades.

## Dependency Evidence

The dependency audit is
`docs/FUTURES_RANGE_CONTEXT_ROUTER_AUDIT_REVIEW.md`.

Accepted source:

| Field | Value |
| --- | --- |
| Source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Router audit facts:

| Field | Value |
| --- | ---: |
| Router rules | `58` |
| Router rows | `29,784` |
| `no_trade` rows | `13,546` |
| `tradable_rotation` rows | `1,299` |
| `trend_continuation` rows | `0` |
| `diagnostic_only` rows | `14,939` |
| Conflicts | `0` |
| Passing router cohorts | `3` |
| Passing `tradable_rotation` cohorts | `1` |

The passed rotation cohort had `892` full-period rows, weakest split rows
`210`, full expected-route hit rate `0.599776`, weakest split expected hit rate
`0.554667`, full adverse route hit rate `0.181614`, worst split adverse route
hit rate `0.247619`, and dominant forward label `contained_rotation` at
`0.385650`.

That is enough to authorize a premise audit, but not enough to define an entry.
The dominant forward label is not a trade rule, and it is not a majority of the
cohort.

## Premise

Premise name:
`router_gated_boundary_reclaim_rotation_v1`.

Premise:

> When a closed `15m` BTCUSDT futures range state belongs to the passed
> `h24` `tradable_rotation` router cohort, the next closed in-range boundary
> probe that reclaims toward the frozen range interior may isolate contained
> rotation to the range midpoint or opposite inner quartile better than using
> the router row alone.

This is an event-quality hypothesis, not a strategy. The later audit may ask
whether a specific closed-candle event improves the router's contained-rotation
evidence. It must not define order direction, entries, stops, targets, sizing,
P&L, or portfolio behavior.

## Source And Context Contract

The next audit must use only:

- Binance USDT-M futures;
- BTCUSDT;
- the accepted `5m` parent source listed above;
- closed UTC `15m` resampling;
- the deterministic range-state construction loop and router rules already
  reviewed;
- closed candles through the router context row or event candle.

No ETH/SOL authority stream, derivatives context source, spot comparison,
external strategy score, live code, or old `binance-bot` strategy logic is part
of this premise.

## Qualifying Router Context

The later audit must not treat every rotation-routed row as an event.

A router context row is eligible only when all of these are true:

- `router_label=tradable_rotation`;
- `timeframe=15m`;
- at least one matched rule id starts with
  `range_context_router_v1|tradable_rotation|15m|h24|`;
- `conflicting_rule_match=false`;
- `missing_rule_match=false`;
- `closed_candle_only=true`;
- `forward_labels_as_router_input=false`;
- `forward_label_columns_present=false`.

Consecutive eligible rows in the same `range_episode_id` must be collapsed into
a single router context segment. The segment start is the first eligible row in
that contiguous run. Later rows in the segment are duplicate context and must be
counted as skipped, not as additional opportunities.

Range high, range low, midpoint, and quartiles are frozen at the segment start
using only candles from the state episode start through that closed segment
start candle. Later candles may form an event, but they must not move the
frozen range bounds used for labeling this premise.

## Candidate Event

The next audit may search only the first `6` closed `15m` candles after the
router context segment start. If no event forms in that window, the segment is
skipped as `no_boundary_reclaim_event`.

Use `20%` of frozen range width as the boundary zone. The range width must be
positive.

Lower-boundary reclaim event:

- the event candle low enters the lower boundary zone;
- the event candle closes back above the lower boundary zone;
- the event candle closes below the frozen midpoint;
- the event candle closes inside the frozen range;
- the event candle has an interior-pointing body, with close above open and
  close position at least `0.60` of that candle's high-low range.

Upper-boundary reclaim event:

- the event candle high enters the upper boundary zone;
- the event candle closes back below the upper boundary zone;
- the event candle closes above the frozen midpoint;
- the event candle closes inside the frozen range;
- the event candle has an interior-pointing body, with close below open and
  close position no more than `0.40` of that candle's high-low range.

An event is rejected if it closes outside the frozen range, already reaches or
crosses the frozen midpoint, has zero candle range, forms after the first valid
event in the segment, or requires a future label to identify it.

## Outcome Labels

Evaluate outcomes over `24` closed `15m` candles after the event candle. This
is an audit horizon, not a max hold.

For a lower-boundary reclaim event, label:

- `midline_rotation_first`: high reaches the frozen midpoint before hard
  lower-boundary failure;
- `opposite_inner_quartile_first`: high reaches the frozen upper quartile
  before hard lower-boundary failure;
- `boundary_failure_first`: a later candle closes below the frozen range low
  before midpoint rotation;
- `clean_expansion_against_rotation`: a later candle closes below the frozen
  range low by at least `15%` of range width, or two consecutive candles close
  below the frozen range low, before midpoint rotation;
- `boundary_chop_no_rotation`: at least three later candles probe the same
  lower boundary zone without midpoint rotation or hard boundary failure;
- `no_resolution`: none of the above resolves inside `24` bars.

For an upper-boundary reclaim event, mirror the labels around the frozen range
midpoint and upper/lower boundaries.

Outcome labels must start strictly after the event candle. They must never be
inputs to router context selection, event formation, grouping, or skip
decisions.

## Review Gates For The Next Audit

The next audit may pass the premise gate only if all source, router, event, and
behavior gates pass.

Source and dependency gates:

- accepted BTCUSDT Binance USDT-M futures source facts match this spec;
- closed UTC `15m` resample validation passes;
- router dependency stop state is
  `range_context_router_passed_needs_rotation_premise_spec`;
- no spot comparison, source expansion, or symbol expansion is used.

Event integrity gates:

- every candidate context and event row is closed-candle only;
- event inputs use only data known at or before the event candle close;
- frozen range bounds are taken from the router context segment start;
- forward labels are stored only as labels;
- duplicate router rows in one context segment are skipped.

Count gates:

- at least `250` full-period router context segments;
- at least `150` full-period valid events;
- at least `40` valid events in every period split;
- both event sides have at least `40` full-period valid events, or the weaker
  side is explicitly marked as diagnostic and no side-symmetric later premise
  is allowed;
- no single period split contributes more than `45%` of valid events.

Behavior gates:

- full-period `midline_rotation_first` rate is at least `3` percentage points
  above the router cohort's `0.599776` expected-route hit rate;
- weakest-split `midline_rotation_first` rate is at least `2` percentage
  points above the router cohort's `0.554667` weakest-split expected hit rate;
- full-period hard adverse rate, defined as
  `boundary_failure_first + clean_expansion_against_rotation`, is no higher
  than `0.22`;
- worst-split hard adverse rate is no higher than `0.28`;
- full-period `boundary_chop_no_rotation + no_resolution` is no higher than
  `0.30`;
- the event result is not carried entirely by one state-id rollup or one short
  historical regime.

Passing this audit would still not authorize a strategy backtest. A pass may
only authorize another non-trading audit that studies whether a later
closed-candle trigger can be specified without reopening a failed family.

## Not A Closed-Family Retry

This premise is allowed only as a materially different, router-gated audit:

- not `range_occupancy_rotation_v1`: it does not use the V1 occupancy grammar,
  target/stop template, optimizer grid, or fixed replay route;
- not hold-inside/midline: a router context segment and boundary-reclaim event
  must exist before the midpoint becomes an outcome label;
- not breakout-retest/acceptance: broken range boundaries and outside
  acceptance are rejected, not traded;
- not clean breakout continuation: clean expansion against the rotation is an
  adverse label;
- not structured compression: there is no symbol-set authority stream,
  confirmation-window grid, target/stop/max-hold grid, or walk-forward
  selection;
- not impulse absorption: abnormal OHLCV impulse candles are not the event
  source, and impulse buckets remain context only;
- not higher-timeframe nested range rotation: this is the passed `15m` router
  cohort, not a `4h` parent / `1h` child nested-range event;
- not range quality, session, or failure-mode triage by themselves: those
  features are not sufficient without the passed router context and new event;
- not legacy spot-only SR timing or compression evidence: the source is the
  accepted BTCUSDT Binance USDT-M futures file, and no spot-only SR conclusion
  is imported.

If implementation work drifts into any closed family, stop with
`range_router_rotation_premise_audit_rejected_closed_family_reslice`.

## Next Audit Brief

The next implementation may add a non-trading audit behind:

```text
-futures-range-router-rotation-premise-audit
```

Default result directory:

```text
results/futures-range-router-rotation-premise-audit/
```

Common outputs must remain zero-trade compatible:

- `source_manifest.json`;
- `summary.csv/json`;
- `trades.json` with no trades.

Audit-specific outputs should include:

- `futures_range_router_rotation_premise_sources.csv/json`;
- `futures_range_router_rotation_premise_coverage.csv/json`;
- `futures_range_router_rotation_premise_router_dependency.csv/json`;
- `futures_range_router_rotation_premise_context_segments.csv/json`;
- `futures_range_router_rotation_premise_events.csv/json`;
- `futures_range_router_rotation_premise_outcomes.csv/json`;
- `futures_range_router_rotation_premise_cohorts.csv/json`;
- `futures_range_router_rotation_premise_rankings.csv/json`;
- `futures_range_router_rotation_premise_summary.csv/json`;
- `futures_range_router_rotation_premise_skips.csv/json`.

Required audit stop states:

- `range_router_rotation_premise_audit_source_router_gap`;
- `range_router_rotation_premise_audit_rejected_closed_family_reslice`;
- `range_router_rotation_premise_audit_no_eligible_events`;
- `range_router_rotation_premise_audit_failed_no_premise`;
- `range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit`;
- `range_router_rotation_premise_audit_rejected_as_strategy_backtest_request`.

Only `range_router_rotation_premise_audit_passed_needs_non_trading_trigger_audit`
may authorize a later non-trading trigger audit. It still would not authorize
entries, exits, P&L backtests, optimizer grids, replay, walk-forward,
paper/testnet/live paths, exchange API use, deployment, source expansion,
symbol expansion, martingale, averaging down, or two-exchange logic.
