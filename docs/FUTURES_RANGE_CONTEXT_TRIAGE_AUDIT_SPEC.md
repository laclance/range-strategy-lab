# Futures Range Context Triage Audit Spec

Date: 2026-06-27

## Verdict

Stop state:
`range_context_triage_spec_ready_for_audit_implementation`.

This is a documentation-only spec for the next bounded offline audit. It
authorizes a non-trading BTCUSDT futures range-context triage pass that assesses
range quality, UTC session behavior, and failure-mode taxonomy in parallel.

It does not add a strategy, entry, exit, optimizer, replay, walk-forward run,
CLI flag, generated result directory, paper/testnet/live path, exchange API
use, credentials, deploy files, data downloads, broad symbol mining,
martingale, averaging down, or two-exchange logic.

## Intent

The latest executable grammar, `range_occupancy_rotation_v1`, failed its
bounded optimizer review with no selectable config. The next useful move is not
another entry grammar. It is a non-trading context audit that asks which range
episodes are structurally tradable, whether their behavior depends on UTC
session, and how weak range contexts fail.

The audit combines three lenses in one pass:

- range quality and tradability;
- session-conditioned behavior;
- failure-mode taxonomy.

These lenses should be evaluated together because a cohort may be unusable in
aggregate but useful after quality and session context are combined, or a
session edge may be explained by a concentrated failure mode.

## Scope

The first implementation must remain BTCUSDT-first and futures-only.

| Field | Contract |
| --- | --- |
| Market | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Parent source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Loaded rows | `573,984` |
| Open-time coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Source status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

The audit may derive only closed UTC `15m`, `1h`, and `4h` bars from the
accepted `5m` parent source.

Expected resample facts:

| Timeframe | Expected rows | Expected first open | Expected last open |
| --- | ---: | --- | --- |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` |

Every resample must use closed UTC buckets, have no gaps, no duplicates, no
missing child opens, complete final buckets after dropping any partial final
bucket, and `validation_status=accepted`.

## Exclusion Boundary

This audit is not a retune or reslice of the closed families. The implementation
must reject drift into:

- structured compression;
- breakout-retest/acceptance;
- clean breakout continuation;
- hold-inside/midline entries;
- impulse absorption;
- higher-timeframe nested range rotation;
- `range_occupancy_rotation_v1`;
- old `binance-bot` strategy, scoring, or live execution code.

The audit may reuse source guards, closed-candle resampling, detector and range
episode helpers, split labels, CSV/JSON artifact writers, and zero-trade common
outputs. It must not emit trade signals, run the backtest engine for strategy
P&L, optimize parameters, or rank executable trade configs.

## Default Audit Contract

Future implementation flag:
`-futures-range-context-triage-audit`.

Future result directory:
`results/futures-range-context-triage-audit/`.

Default detector and range candidate settings:

| Field | Value |
| --- | --- |
| Detector profile | `p30_c12_bollinger_on_adx_off` |
| Detector percentile | `0.30` |
| Detector min consecutive bars | `12` |
| Detector lookback days | `20` |
| Detector ADX filter | off |
| Detector Bollinger filter | on |
| Timeframes | `15m`, `1h`, `4h` |
| Outcome horizons | `12`, `24`, `48` closed bars per timeframe |
| Quick failure horizon | `6` closed bars |
| Reentry window | `6` closed bars |
| Clean expansion threshold | `0.75` frozen range widths beyond the broken boundary |
| Drift threshold | final close outside the range with less than `0.50` range-width extension |
| Boundary chop transition threshold | `3` inside/outside state transitions after first outside close |
| Low width threshold | `width_pct < 0.0015` or `width_to_atr_ratio < 0.75` |
| Wide volatile threshold | `width_pct >= 0.0150` or `width_to_atr_ratio >= 4.00` |
| Choppy pre-mature threshold | at least `3` midpoint crosses and at least `4` boundary touches |
| Minimum full cohort count | `300` |
| Minimum split cohort count | `50` |
| Minimum session split cohort count | `30` |
| Minimum usable context rate | `0.55` full, `0.45` in every period split |
| Maximum toxic context rate | `0.45` full, `0.55` in every period split |
| Maximum missing future rate | `0.02` full |

Test hooks may relax source paths and expected row counts for small fixtures,
but CLI defaults must remain locked to this spec.

## Range Candidate Construction

For each selected timeframe:

1. Resample the accepted `5m` parent source into closed UTC bars.
2. Run the balanced detector profile on the resampled bars.
3. Build one audit range candidate from each raw-active detector run that
   reaches a mature `Active` candle.
4. Freeze candidate geometry at the first mature close, using only bars from
   raw-active start through that mature close.
5. Record raw-active end metadata separately as descriptive context only; it
   must not be used for future strategy eligibility or cohort ranking.

The frozen geometry must include high, low, midpoint, upper quartile, lower
quartile, width, width percent, close position, width-to-ATR ratio, raw bars to
maturity, active bars to maturity, and UTC session at mature close.

Invalid candidates are skipped with explicit reasons:

- `no_mature_active_bar`;
- `non_positive_width`;
- `non_positive_price`;
- `missing_atr`;
- `missing_future`;
- `source_or_resample_gap`;
- `closed_family_reslice`.

## Quality Lens

The quality lens scores the candidate range before any future outcome label is
read. It may use only bars from raw-active start through the mature close.

Required quality fields:

- `width_pct`;
- `width_to_atr_ratio`;
- `duration_bars_to_maturity`;
- `active_bars_to_maturity`;
- `mature_close_position`;
- `pre_mature_mid_cross_count`;
- `pre_mature_boundary_touch_count`;
- `pre_mature_close_inside_rate`;
- `pre_mature_wick_overshoot_count`;
- `quality_bucket`.

Required quality buckets:

- `too_narrow_noise`;
- `narrow_orderly`;
- `balanced_orderly`;
- `wide_volatile`;
- `choppy`;
- `unknown`.

Quality bucket precedence must be deterministic:

1. `too_narrow_noise` when `width_pct < 0.0015` or
   `width_to_atr_ratio < 0.75`;
2. `wide_volatile` when `width_pct >= 0.0150` or
   `width_to_atr_ratio >= 4.00`;
3. `choppy` when `pre_mature_mid_cross_count >= 3` and
   `pre_mature_boundary_touch_count >= 4`;
4. `narrow_orderly` when `width_pct < 0.0030`, `width_to_atr_ratio < 1.50`,
   and `pre_mature_close_inside_rate >= 0.90`;
5. `balanced_orderly` when `pre_mature_close_inside_rate >= 0.90`;
6. `unknown`.

The bucket is context evidence only and must not be used as an entry signal.

## Session Lens

Every candidate and every labeled outcome must carry UTC session labels.

Default mature-close sessions:

| Session | UTC hours |
| --- | --- |
| `asia_utc_00_07` | `00:00` through `07:59` |
| `europe_utc_08_12` | `08:00` through `12:59` |
| `us_overlap_utc_13_16` | `13:00` through `16:59` |
| `us_late_utc_17_23` | `17:00` through `23:59` |

The audit must summarize by split, timeframe, horizon, mature-close session,
and outcome session when present. It must record whether any session improves
usable context rate by at least `0.15` versus the worst adequately populated
session in the same timeframe and horizon. Session results are context evidence
only; they do not authorize a session filter strategy until a later spec.

## Failure-Mode Lens

For each valid candidate and horizon, label the future window beginning at the
next closed bar after mature close.

Required future-window fields:

- first outside close side: `up`, `down`, or `none`;
- bars to first outside close;
- first outside session;
- reentry within `6` bars after first outside close;
- maximum excursion above frozen high in range widths;
- maximum excursion below frozen low in range widths;
- final close position relative to the frozen range;
- inside close rate;
- midpoint cross count;
- outside state transition count;
- missing future flag.

Each candidate-horizon row must receive exactly one primary context label:

- `contained_rotation`: no outside close, at least two midpoint crosses, and
  at least one upper-half and one lower-half close inside the frozen range;
- `clean_expansion_up`;
- `clean_expansion_down`;
- `false_break_reentry_up`;
- `false_break_reentry_down`;
- `boundary_chop`;
- `drift_through_up`;
- `drift_through_down`;
- `low_width_noise`;
- `no_resolution`;
- `missing_future`.

Label precedence must be deterministic:

1. `missing_future`;
2. `low_width_noise`;
3. `contained_rotation`;
4. `boundary_chop`;
5. `false_break_reentry_up/down`;
6. `clean_expansion_up/down`;
7. `drift_through_up/down`;
8. `no_resolution`.

Constructive context labels are `contained_rotation`,
`clean_expansion_up`, and `clean_expansion_down`. Toxic context labels are
`false_break_reentry_up`, `false_break_reentry_down`, `boundary_chop`,
`drift_through_up`, `drift_through_down`, `low_width_noise`, and
`no_resolution`.

## Cohorts And Review Gates

The audit must summarize at these cohort levels:

- split, timeframe, horizon;
- split, timeframe, horizon, quality bucket;
- split, timeframe, horizon, mature session;
- split, timeframe, horizon, primary context label;
- split, timeframe, horizon, quality bucket, mature session.

The audit may rank cohorts as future strategy-spec candidates only when the
cohort is not an executable trade config. A ranked cohort is a context premise,
not a strategy.

A cohort passes the review gate only if all conditions hold:

- source and resample validation pass;
- full-period count is at least `300`;
- every period split count is at least `50`;
- if the cohort includes mature session, every period split count is at least
  `30`;
- full usable context rate is at least `0.55`;
- every period split usable context rate is at least `0.45`;
- full toxic context rate is no more than `0.45`;
- every period split toxic context rate is no more than `0.55`;
- missing future rate is no more than `0.02`;
- no single toxic label dominates every period split above `0.50`;
- the cohort is not a closed-family reslice.

Ranking score for passing cohorts:

```text
usable_rate_full
+ weakest_split_usable_rate
- toxic_rate_full
- worst_split_toxic_rate
+ min(0.20, session_edge_rate)
+ min(0.20, log10(full_count) / 20)
```

Tie-breaks:

1. higher weakest-split usable context rate;
2. lower worst-split toxic context rate;
3. higher full candidate count;
4. timeframe order `1h`, `15m`, `4h`;
5. simpler cohort: no session before session, no quality bucket before quality
   bucket;
6. lexical cohort ID.

## Stop States

Spec stop states:

- `range_context_triage_spec_ready_for_audit_implementation`;
- `range_context_triage_spec_needs_user_premise_or_scope_input`;
- `range_context_triage_spec_rejected_closed_family_reslice`.

Later implementation stop states:

- `range_context_triage_source_gap`;
- `range_context_triage_no_range_episodes`;
- `range_context_triage_no_usable_cohorts`;
- `range_context_triage_failed_no_strategy_premise`;
- `range_context_triage_ready_for_strategy_spec`;
- `range_context_triage_rejected_closed_family_reslice`.

`range_context_triage_ready_for_strategy_spec` is allowed only when source and
resample validation pass and at least one non-closed-family cohort passes the
review gate. It authorizes only a later documentation-only strategy spec, not
a baseline backtest, optimizer, replay, walk-forward, or strategy package.

## Artifacts For Later Implementation

The later audit must keep common outputs zero-trade compatible:

- `source_manifest.json`;
- `summary.csv`;
- `summary.json`;
- `trades.json`.

Common `trades.json` must contain zero trades. Common summary rows must
represent the no-trade compatibility run only.

Triage-specific artifacts must be written under
`results/futures-range-context-triage-audit/` as CSV and JSON:

- `futures_range_context_triage_sources`;
- `futures_range_context_triage_coverage`;
- `futures_range_context_triage_episodes`;
- `futures_range_context_triage_quality`;
- `futures_range_context_triage_sessions`;
- `futures_range_context_triage_failure_modes`;
- `futures_range_context_triage_cohorts`;
- `futures_range_context_triage_rankings`;
- `futures_range_context_triage_summary`.

The implementation review doc should be:
`docs/FUTURES_RANGE_CONTEXT_TRIAGE_AUDIT_REVIEW.md`.

## Verification For This Spec

This spec milestone should close with documentation checks only:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
