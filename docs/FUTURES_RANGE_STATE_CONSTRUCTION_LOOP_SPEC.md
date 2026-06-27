# Futures Range State Construction Loop Spec

Date: 2026-06-27

## Verdict

Stop state:
`range_state_construction_loop_spec_ready_for_audit_implementation`.

This is a documentation-only implementation spec for the next non-trading
BTCUSDT range-state audit. It does not add a strategy, entry, exit, optimizer,
replay, walk-forward run, source download, symbol expansion, paper/testnet/live
path, exchange API, credential, deploy file, martingale, averaging down, or
two-exchange logic.

## Intent

The next research step should test a broader construction premise:

> A range is not a trade signal by itself. A range becomes useful only after its
> geometry is interpreted together with volatility, trend, impulse, and
> liquidity/participation state.

The audit must produce closed-candle state evidence before any strategy grammar
is written. The expected output is a state/routing map, not trades.

## Source Contract

Default source:

| Field | Contract |
| --- | --- |
| Market | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Parent interval | `5m` |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Loaded candles | `573,984` |
| Coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Accepted facts | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

The implementation must reject gaps, duplicates, irregular `5m` cadence,
non-positive or non-finite OHLC prices, negative volume, invalid high/low
containment, non-BTCUSDT identity, non-`5m` identity, spot paths, and
comparison-only paths.

## Timeframes

The audit may derive only these closed UTC bars from the accepted `5m` source:

| Timeframe | Expected rows | Expected first open | Expected last open |
| --- | ---: | --- | --- |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` |

Every resample row must have `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Exclusion Boundary

The audit may reuse infrastructure only:

- source guard and source manifest discipline;
- closed UTC resampling;
- confirmed-candle semantics;
- detector and range episode helpers as feature extraction only;
- forward event labels as labels only;
- split definitions and artifact patterns.

The audit must stop with
`range_state_construction_loop_rejected_closed_family_reslice` if the
implementation becomes any of these:

- `range_occupancy_rotation_v1` retune;
- structured-compression retune;
- breakout-retest/acceptance retune;
- clean breakout continuation retune;
- hold-inside/midline retune;
- impulse absorption entry;
- higher-timeframe nested range rotation;
- quality/session/failure-mode cohort scan without the new state dimensions;
- old spot-only SR timing or compression promotion evidence.

## Audit Unit

For each signal timeframe `T`, create state rows only from closed candles known
at decision time.

A row is eligible when all are true:

1. source and resample validation passed;
2. the row belongs to a mature active range episode under existing range episode
   helpers;
3. all feature windows can be computed using candles up to and including the
   decision candle;
4. forward-label windows exist for every declared horizon.

Rows that fail eligibility must be counted with skip reasons. They must not be
silently dropped.

## Feature Families

All features must be known at the decision candle close. Future candles may be
used only for label artifacts and cohort outcomes.

### 1. Range Geometry State

Required fields:

- `range_age_bars`;
- `range_width_pct`;
- `range_width_atr_ratio`;
- `close_location_pct` within the active range;
- `distance_to_low_pct`;
- `distance_to_high_pct`;
- `distance_to_mid_pct`;
- `boundary_touch_count_lookback`;
- `midline_cross_count_lookback`;
- `close_location_entropy_lookback`;
- `range_duty_cycle_lookback`.

Suggested buckets:

- `geometry_narrow_orderly`;
- `geometry_balanced_orderly`;
- `geometry_wide_volatile`;
- `geometry_boundary_crowded`;
- `geometry_midline_balanced`.

### 2. Volatility State

Required fields:

- ATR over a short closed-candle window;
- ATR percentile over a longer closed-candle lookback;
- realized close-to-close volatility percentile;
- `atr_to_range_width`;
- `true_range_expansion_ratio` versus prior window;
- `abnormal_range_bar_count_lookback`.

Suggested buckets:

- `vol_compressed`;
- `vol_normal`;
- `vol_expanding`;
- `vol_extreme`.

### 3. Trend State

Required fields:

- close-to-close return over short and medium lookbacks;
- moving-average slope or equivalent closed-candle slope proxy;
- higher-timeframe direction proxy when available from the same BTCUSDT parent
  source;
- trend-strength proxy such as ADX or an in-lab equivalent;
- distance from moving average or range midpoint.

Suggested buckets:

- `trend_flat`;
- `trend_up_pressure`;
- `trend_down_pressure`;
- `trend_strong_up`;
- `trend_strong_down`.

### 4. Impulse State

Required fields:

- last abnormal candle side;
- bars since last abnormal candle;
- count of large-body candles in the lookback;
- count of large-range candles in the lookback;
- impulse continuation pressure based only on already-closed candles;
- impulse exhaustion proxy based only on already-closed candles.

This feature family must not recreate the failed impulse-absorption entry
premise. It is context only.

Suggested buckets:

- `impulse_none`;
- `impulse_up_recent`;
- `impulse_down_recent`;
- `impulse_clustered`;
- `impulse_stale`.

### 5. OHLCV Liquidity / Participation Proxy

Required fields from existing candles only:

- volume percentile;
- volume change ratio versus prior window;
- volume per range-width proxy;
- candle spread proxy: `(high-low)/close`;
- wick/body structure summary;
- zero-volume row flag.

Suggested buckets:

- `participation_low`;
- `participation_normal`;
- `participation_high`;
- `participation_dislocated`.

No order-book, funding, open-interest, taker-flow, or other external market data
is part of this V1 audit.

## State ID Construction

Every eligible decision row must receive a deterministic `state_id` built from
bucketed feature families:

```text
range_state_v1::<timeframe>::<geometry_bucket>::<vol_bucket>::<trend_bucket>::<impulse_bucket>::<participation_bucket>
```

Also write less granular rollups:

```text
geometry+vol
geometry+trend
geometry+impulse
geometry+participation
geometry+vol+trend
geometry+vol+trend+impulse
all_families
```

The implementation must not add a learned classifier in V1. Buckets must be
predeclared, deterministic, closed-candle knowable, and inspectable.

## Forward Labels

Forward labels are outcome labels only. They may be used for cohort summaries,
rankings, and review, but never as features.

Use horizons:

| Timeframe | Horizons |
| --- | --- |
| `15m` | `12`, `24`, `48` bars |
| `1h` | `12`, `24`, `48` bars |
| `4h` | `6`, `12`, `24` bars |

Required labels:

- `contained_rotation`;
- `boundary_chop`;
- `clean_expansion_up`;
- `clean_expansion_down`;
- `false_break_reentry_up`;
- `false_break_reentry_down`;
- `drift_through_up`;
- `drift_through_down`;
- `no_resolution`;
- `low_width_noise`.

The existing failure-mode taxonomy from the range-context triage audit may be
reused, but the cohort decision must come from the new state dimensions.

## Route Interpretation

The audit should summarize each cohort into route candidates:

| Route | Useful Labels | Toxic Labels | Meaning |
| --- | --- | --- | --- |
| `tradable_rotation_candidate` | `contained_rotation`, false-break reentry back inside range | `boundary_chop`, clean expansion against the candidate route | Possible later mean-reversion spec, not an entry now. |
| `trend_continuation_candidate` | clean expansion, drift-through with low chop | false-break reentry, boundary chop | Possible later continuation spec, not structured-compression rescue. |
| `no_trade_toxic` | none required | high `boundary_chop`, `low_width_noise`, unstable outcomes | Useful as a filter even if no strategy is authorized. |
| `diagnostic_only` | mixed or unstable | mixed or unstable | Record but do not promote. |

A route candidate can pass only if it is stable across period splits and has
adequate counts. Outcome labels cannot directly select trade direction.

## Cohort Gates

A cohort is reviewable only when all are true:

- full-period eligible rows >= `150`;
- each period split rows >= `30`;
- no single calendar split contributes more than `60%` of full-period rows;
- no single route label contributes more than `80%` of rows unless it is the
  declared `no_trade_toxic` route;
- all required feature buckets are non-empty and known at decision time.

A `tradable_rotation_candidate` or `trend_continuation_candidate` passes only
when all are true:

- full useful-rate >= `0.58`;
- weakest split useful-rate >= `0.52`;
- full toxic-rate <= `0.42`;
- worst split toxic-rate <= `0.48`;
- full useful-minus-toxic margin >= `0.12`;
- weakest split useful-minus-toxic margin >= `0.04`;
- the route does not depend on a future-known label as an input;
- the route is not a closed-family reslice.

A `no_trade_toxic` cohort passes when all are true:

- full toxic-rate >= `0.58`;
- weakest split toxic-rate >= `0.52`;
- full rows >= `200`;
- each split rows >= `40`;
- the toxic state is closed-candle knowable and not just a future label.

Passing a cohort authorizes only a later documentation-only router or strategy
premise spec. It does not authorize entries, exits, P&L, optimization, replay,
walk-forward, packaging, symbol expansion, source expansion, or live-adjacent
work.

## Required Artifacts

Default result directory:

```text
results/futures-range-state-construction-loop-audit/
```

Common outputs must remain zero-trade compatible:

- `source_manifest.json`;
- `summary.csv`;
- `summary.json`;
- `trades.json` with no trades.

Audit-specific outputs:

- `futures_range_state_construction_loop_sources.csv/json`;
- `futures_range_state_construction_loop_coverage.csv/json`;
- `futures_range_state_construction_loop_feature_windows.csv/json`;
- `futures_range_state_construction_loop_states.csv/json`;
- `futures_range_state_construction_loop_labels.csv/json`;
- `futures_range_state_construction_loop_cohorts.csv/json`;
- `futures_range_state_construction_loop_rankings.csv/json`;
- `futures_range_state_construction_loop_summary.csv/json`;
- `futures_range_state_construction_loop_skips.csv/json`.

At minimum, `states` rows must include:

- timestamp;
- timeframe;
- split;
- range episode id;
- all raw feature values;
- all feature buckets;
- `state_id`;
- route rollup id;
- skip status if applicable.

At minimum, `labels` rows must include:

- timestamp;
- timeframe;
- horizon;
- forward label;
- favorable/toxic flags by route interpretation;
- future-window metadata used only for labels.

At minimum, `cohorts` and `rankings` rows must include:

- cohort id;
- route candidate;
- timeframe;
- horizon;
- full and split counts;
- useful/toxic rates;
- margins;
- pass/fail flags;
- failure reasons;
- closed-family protection flag.

## Expected CLI Flag

The later implementation should run only behind:

```text
-futures-range-state-construction-loop-audit
```

The flag must reject combinations with trade-producing strategy, baseline,
optimization, replay, walk-forward, source-expansion, symbol-expansion, or
external-data flags.

## Review Document

A later implementation run must produce:

```text
docs/FUTURES_RANGE_STATE_CONSTRUCTION_LOOP_REVIEW.md
```

The review must include source facts, resample facts, artifact line counts,
state counts, top route cohorts, failed cohort reasons, and final stop state.

## Stop States

Spec stop states:

- `range_state_construction_loop_spec_ready_for_audit_implementation`;
- `range_state_construction_loop_spec_needs_user_scope_input`;
- `range_state_construction_loop_spec_rejected_closed_family_reslice`.

Implementation/review stop states:

- `range_state_construction_loop_source_gap`;
- `range_state_construction_loop_no_eligible_states`;
- `range_state_construction_loop_audit_failed_no_usable_state`;
- `range_state_construction_loop_audit_passed_no_trade_filter_only`;
- `range_state_construction_loop_audit_passed_needs_router_spec`;
- `range_state_construction_loop_audit_passed_needs_strategy_premise_spec`;
- `range_state_construction_loop_rejected_closed_family_reslice`.

## Verification Expectations

For the documentation-only spec update:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

For the later implementation review, also run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -futures-range-state-construction-loop-audit \
  -out-dir results/futures-range-state-construction-loop-audit
wc -l results/futures-range-state-construction-loop-audit/*.csv
```

Do not claim a passing state unless the review gates pass from generated
artifacts.