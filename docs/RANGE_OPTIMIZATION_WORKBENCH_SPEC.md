# Range Optimization Workbench Spec

Date: 2026-06-30

## Verdict

Stop state:

```text
range_optimization_workbench_spec_ready_for_implementation_approval
```

This is a docs-only specification for a bounded offline optimization/workbench
lane. It does **not** authorize paper trading, testnet, live trading, exchange API
work, credentials, deployment, promotion, or production changes.

The purpose is to answer a different question from the fixed-baseline lane:

```text
Fixed-baseline lane:
Does this exact simple idea work as specified?

Optimization/workbench lane:
Can a bounded family of related range features produce a robust candidate worth
locking and validating as a new fixed strategy?
```

The three completed backtest-first candidates failed as fixed baselines:

- `btc_5m_rolling_value_area_reversion_v1`;
- `btc_15m_previous_day_range_reversion_v1`;
- `btc_15m_range_edge_exhaustion_fade_v1`.

Those failures do not prove that every combination of range context, regime
state, entry timing, and exit construction is useless. They do mean that any
further combination/tuning work must be treated as an explicit search problem,
with all trials logged and no post-hoc promotion.

## Core Principle

Optimization is allowed only as **controlled discovery**.

It must not be used to quietly rescue one of the failed baselines, retune a few
thresholds until a curve looks good, or present a selected parameter cell as if
it were an independent first-pass test.

The optimizer may search for a candidate. A selected candidate must then be
locked into a later fixed validation spec before any stronger claim is made.

## Scope

Allowed:

- offline only;
- current accepted BTCUSDT Binance USDT-M futures source;
- 5m native candles and exact closed 15m resamples;
- bounded combinations of existing range-family features;
- bounded regime/context features available at decision time;
- trial enumeration with every trial recorded;
- robustness ranking;
- candidate selection for later fixed validation.

Forbidden:

- paper/testnet/live trading;
- exchange API calls;
- credentials;
- deployment files;
- production bot integration;
- martingale;
- averaging down;
- two-exchange execution;
- hidden trial deletion;
- hand-picking a lucky cell;
- post-result filter additions;
- changing the data source after seeing results;
- final promotion from optimizer output alone.

## Source Contract

Use the current accepted BTCUSDT Binance USDT-M futures CSV:

```text
../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
```

Required source facts:

```text
product: Binance USDT-M futures
symbol: BTCUSDT
base interval: 5m
loaded candles: 573,984
first open: 2021-01-01T00:00:00Z
last open: 2026-06-16T23:55:00Z
gap_count: 0
duplicate_count: 0
zero_volume_count: 66
comparison_only: false
validation_status: accepted
```

If 15m features are used, the resample must be exact closed UTC 15m bars from
complete three-child 5m buckets.

Expected 15m source facts:

```text
row_count: 191,328
last_open: 2026-06-16T23:45:00Z
expected_child_bars: 3
missing_child_buckets: 0
```

No forward filling, interpolation, synthetic candles, spot substitution, or
comparison-only source is allowed.

## Families Eligible For Combination

The workbench may combine components from the three failed families, but must
label them as components rather than as rescued strategies.

### Rolling value-area components

Allowed components:

- rolling high/low range from closed candles;
- rolling midpoint;
- rolling VWAP from typical price weighted by volume;
- range width relative to ATR;
- close position inside range;
- distance from VWAP/midpoint.

### Previous-day range components

Allowed components:

- prior complete UTC day high/low/midpoint;
- current day inside/outside previous-day range state;
- distance to prior-day edge/midpoint;
- previous-day range width relative to ATR.

### Range-edge exhaustion components

Allowed components:

- movement from range interior toward an edge;
- final-progress weakness relative to ATR;
- edge-zone position;
- fade target to midpoint.

## Additional Context Features

The workbench may add bounded decision-time context features.

Allowed context groups:

- volatility regime:
  - ATR percentile;
  - range width divided by ATR;
  - realized volatility percentile;
  - compression/expansion state.
- trend/impulse regime:
  - EMA slope direction;
  - close relative to EMA;
  - recent impulse size relative to ATR;
  - number of same-direction closes.
- liquidity/activity proxy:
  - volume percentile;
  - zero-volume guard;
  - volume expansion/contraction relative to recent median.
- time/session context:
  - UTC hour bucket;
  - weekday;
  - optional coarse session bucket.

Forbidden context groups unless separately approved:

- derivatives mark/index/premium veto integration;
- ETH/SOL cross-market context;
- funding-based selection;
- order-book or exchange-adapter data;
- any future label or forward return as an input feature.

## Search Space Boundaries

The first implementation must be deliberately small enough to inspect and rerun.

Maximum first-pass trial budget:

```text
max_trials: 2,500
```

A smaller trial count is preferred if the grid can still cover the key families.

### Timeframes

Allowed:

```text
5m native
15m exact resample
```

### Range lookbacks

Allowed 5m lookbacks:

```text
144, 288, 576
```

Allowed 15m lookbacks:

```text
48, 96, 192
```

### Entry archetypes

Allowed:

```text
edge_fade
midpoint_reversion
vwap_reversion
interior_to_edge_exhaustion_fade
previous_day_edge_reversion
```

### Edge/zone thresholds

Allowed:

```text
edge_zone_pct: 0.10, 0.15, 0.20
interior_low_pct: 0.35, 0.40, 0.45
interior_high_pct: 0.55, 0.60, 0.65
vwap_distance_range_pct: 0.10, 0.15, 0.20
```

### Volatility/range gates

Allowed:

```text
min_range_atr: 3, 4, 6
max_range_atr: none, 10, 14
atr_percentile_min: none, 20, 40
atr_percentile_max: none, 80, 95
```

### Trend/impulse gates

Allowed:

```text
trend_mode: none, with_trend_only, against_trend_only, flat_only
ema_length: 48, 96, 192
impulse_atr_min: none, 1.0, 1.5
impulse_atr_max: none, 2.5, 4.0
```

### Volume/activity gates

Allowed:

```text
volume_mode: none, above_median, below_80pct, expansion_only
volume_lookback: 96, 192
```

### Exits

Allowed target modes:

```text
midpoint
vwap
opposite_inner_zone
```

Allowed stop modes:

```text
range_edge_plus_0.25_atr
range_edge_plus_0.50_atr
entry_minus_1.00_atr
```

Allowed time stops:

```text
5m: 24, 36, 72 closed bars
15m: 8, 16, 24 closed bars
```

No trailing stop, partial take-profit, pyramiding, martingale, or averaging down
is allowed in the first-pass workbench.

## Fixed Cost And Risk Contract

Every trial must use the same cost/risk assumptions:

```text
start_balance: 1000
risk_pct: 0.01
max_notional_pct: 1.0
fee_pct_per_side: 0.0004
slippage_pct_per_side: 0.000116
one_position_only: true
entry_timing: next_bar_open
stop_first_ambiguity: true
```

No trial may improve costs, increase leverage, widen notional cap, or change
risk sizing to rescue results.

## Validation Design

The workbench must separate **search ranking** from **candidate validation**.

Required split labels:

```text
2021_2022_stress
2023_2024_oos
2025_2026_recent
full_2021_2026
```

Because these split metrics have already been viewed in prior work, the recent
split should not be called a pristine hidden holdout. It is still required as a
confirmation split and must be reported separately.

The first-pass workbench may rank candidates using a robustness score, but it
must report all split metrics for every trial.

## Trial Logging Requirements

Every trial must be recorded. No failed or ugly trial may be omitted.

The implementation must write each run into an immutable run directory. The
canonical parent directory is only an index/container and must not be deleted to
rerun verification.

Required run directory shape:

```text
results/range-optimization-workbench-v1/
  runs/
    <run_id>/
      source_manifest.json
      optimization_grid.json
      trial_results.json
      trial_results.csv
      trial_summary.csv
      top_candidates.csv
      rejected_candidates.csv
      robustness_summary.json
      falsification.json
  latest_run.json
```

`run_id` must be unique and stable for the run, for example a UTC timestamp plus
short git SHA. If a run directory already exists, the implementation must fail or
choose a new run ID; it must not overwrite prior trial outputs.

Every trial row must include:

```text
trial_id
family_id
timeframe
all parameter values
source fingerprint or manifest reference
trade counts by split
gross P&L by split
net P&L by split
profit factor by split
max drawdown by split
long/short trade counts
long/short net P&L
robustness score
failure reasons
selected_for_locking
```

## Robustness Ranking

A candidate may be ranked, but not promoted, if it passes minimum workbench
filters.

Minimum candidate filters:

```text
full_trades >= 200
minimum_primary_split_trades >= 40
full_gross_pnl > 0
full_net_pnl > 0
2023_2024_oos_gross_pnl >= 0
2025_2026_recent_gross_pnl >= 0
full_max_drawdown <= 0.30
no single primary split has more than 75% of full trades
both side reporting rows exist
```

Preferred ranking score:

```text
robustness_score =
  split_net_consistency
  + gross_edge_consistency
  + profit_factor_stability
  + drawdown_penalty
  + trade_distribution_penalty
  + side_concentration_penalty
```

The exact score implementation must be documented in the implementation PR.

A candidate with the highest net P&L must not automatically win. A lower-P&L
candidate with better split stability may rank higher.

## Selection Rules

At most one candidate may be selected by the first workbench run.

A selected candidate must produce a locked candidate packet containing:

```text
candidate_id
fixed source contract
complete fixed parameter set
entry rules
exit rules
risk/cost rules
all split metrics from the workbench
why it was selected
known weaknesses
no-rescue boundary
next fixed-validation command
```

The selected candidate must then be rerun through a later fixed-validation lane.
The optimizer's selected result is not itself enough for paper/shadow/live or
promotion.

## Falsification Rules

The workbench must stop with one of these states:

```text
range_optimization_workbench_spec_ready_for_implementation_approval
range_optimization_workbench_implementation_added_needs_local_run
range_optimization_workbench_failed_no_candidate
range_optimization_workbench_candidate_selected_needs_fixed_validation
range_optimization_workbench_rejected_overfit_risk
```

Use `range_optimization_workbench_failed_no_candidate` if no candidate passes the
minimum candidate filters.

Use `range_optimization_workbench_rejected_overfit_risk` if candidates pass only
through narrow, unstable, or obviously data-mined parameter islands.

Use `range_optimization_workbench_candidate_selected_needs_fixed_validation` only
if a candidate passes the workbench filters and is selected for a separate locked
validation run.

## Anti-Overfitting Guardrails

The implementation must include the following guardrails:

- the full grid is defined before the run;
- every trial is emitted to disk;
- trial count is reported;
- top candidates and rejected candidates are both reported;
- no manual deletion of failed trials;
- no deleting a prior workbench run to rerun verification;
- no changing the grid after seeing the first result within the same run;
- no adding a new filter because it improves the current best cell;
- no optimizer output may authorize paper/testnet/live;
- if selected, one candidate is locked and tested later as a fixed rule.

## Next Implementation Boundary

The next PR may implement only this offline workbench harness and, if locally
run, record its first result review.

It may not:

- introduce exchange API code;
- touch bot runtime code;
- introduce deployment files;
- add credentials;
- enable paper/testnet/live;
- promote any strategy;
- add unbounded optimizers;
- add genetic algorithms, Bayesian optimization, or reinforcement learning;
- mutate source files;
- delete failed result artifacts.

## Required Next Local Verification For Implementation PR

When implemented, the local verification command should preserve prior runs by
writing into a unique run directory. Do not use `rm -rf` on the canonical
workbench results parent once any workbench output may exist.

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-$(git rev-parse --short HEAD)"
OUT_DIR="results/range-optimization-workbench-v1/runs/${RUN_ID}"

test ! -e "${OUT_DIR}"

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -range-optimization-workbench-v1 \
  -out-dir "${OUT_DIR}" \
  -run-id "${RUN_ID}"

wc -l "${OUT_DIR}"/*.csv

cat "${OUT_DIR}"/falsification.json

git diff --check
git status --short
```

If a rerun is needed, create a new `RUN_ID`. Existing run directories may be
archived or indexed, but must not be deleted as part of normal verification.

## Operator Decision Required

This spec is ready for user/operator review.

If accepted, authorize exactly one bounded offline implementation PR for:

```text
-range-optimization-workbench-v1
```

No paper/testnet/live or promotion is authorized by this spec.
