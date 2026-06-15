# Hold-Inside Midline Reaction Review

Date: 2026-06-15

## Verdict

The hold-inside midline reaction audit clears the gate for one first minimal
non-live entry prototype, but not for strategy promotion or live use.

The only primary surface worth carrying forward is `hold_3_inside` +
`mid_touch`, with strongest support when the midline event close remains in the
frozen range `mid_50` bucket. This surface has enough split-stable candidates,
mid-touch events occur in every split more often than they are missing, and the
post-event labels consistently favor close-back/mid-rejection behavior before a
boundary-first continuation.

This is still a prototype-only finding. No entry code, exits, scoring, sizing,
or strategy replacement was added in this review. The next step may build a
minimal offline prototype to test this exact surface with side-split P&L and
stress evidence. Keep `lab.EmptyStrategy` until that prototype exists and is
reviewed.

## Inputs Reviewed

Generated hold-inside midline reaction audit paths:

- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_candidates.csv`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_candidates.json`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_funnel_summary.csv`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_funnel_summary.json`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_summary.csv`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_summary.json`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_stability.csv`
- `results/hold-inside-midline-reaction-audit/hold_inside_midline_reaction_stability.json`

Audit size:

- profiles: `1`
- context rules: `3`
- event types: `2`
- candidate rows: `9,080`
- funnel rows: `24`
- summary rows: `1,296`
- stability rows: `352`
- CSV lines including headers: `9,081` / `25` / `1,297` / `353`
- horizons: `1`, `3`, `6`, `12`
- max midline event delay: `12` bars after the hold decision
- quick invalidation window: `3` bars after the midline event

Review semantics:

- Detector profile is only `p30_c12_bollinger_on_adx_off`.
- Primary context rules are `hold_3_inside` and `hold_6_inside`.
- `hold_3_inside_mid_50` is diagnostic.
- Event types are `mid_touch` and `mid_close_across`.
- The event candle is the reindexed decision candle.
- All `label_*` fields start at `event_index + 1` and are forward outcomes
  only. They are not decision inputs.
- Stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent`.

## Funnel Gate

Primary-rule funnel rows by weakest split:

| Rule | Event | Min Source Holds | Min Event Count | Min Event Rate | Max Missing Rate | Avg Delay Range | Gate |
| --- | --- | ---: | ---: | ---: | ---: | ---: | --- |
| `hold_3_inside` | `mid_touch` | 222 | 116 | 52.25% | 47.75% | 4.04-4.44 | pass |
| `hold_3_inside` | `mid_close_across` | 222 | 100 | 45.05% | 54.95% | 4.73-5.11 | fail: missing dominates |
| `hold_6_inside` | `mid_touch` | 170 | 89 | 52.35% | 47.65% | 3.50-4.04 | fail: weak split below 100 events |
| `hold_6_inside` | `mid_close_across` | 170 | 79 | 46.47% | 53.53% | 4.20-5.01 | fail: missing dominates and sparse |

`hold_3_inside` + `mid_touch` is the only primary event surface where the event
occurs more often than it is missing in every split while preserving at least
`100` event candidates in the weakest split.

## Reaction Gate

Primary all-bucket stability rows:

| Rule | Event | Horizon | Min Cand | Min Persist | Max Quick Inv | Max Trend | Min Close Back | Min Mid Reject Before Boundary | Max Boundary Before Reject |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | `mid_touch` | 3 | 116 | 81.61% | 18.39% | 12.11% | 47.09% | 42.60% | 23.28% |
| `hold_3_inside` | `mid_touch` | 6 | 116 | 69.06% | 18.39% | 22.22% | 55.17% | 45.69% | 28.25% |
| `hold_3_inside` | `mid_touch` | 12 | 116 | 47.41% | 18.39% | 35.34% | 70.69% | 52.59% | 36.21% |
| `hold_3_inside` | `mid_close_across` | 6 | 100 | 69.49% | 17.51% | 21.19% | 53.00% | 46.00% | 30.00% |
| `hold_6_inside` | `mid_touch` | 6 | 89 | 75.18% | 12.36% | 22.63% | 63.50% | 56.93% | 19.35% |

The reaction shape is strongest as a rejection/close-back surface, not a
boundary-continuation surface:

- For `hold_3_inside` + `mid_touch`, close-back and mid-rejection rates are
  consistently above boundary-first rates in every split and horizon.
- The `6` bar horizon is the cleanest prototype horizon: enough time for the
  reaction to develop, while trend leakage is still roughly `20%-22%` rather
  than the `31%-35%` seen by `12` bars.
- `mid_close_across` has acceptable all-bucket candidate count at `hold_3_inside`
  but fails the funnel because missing close-across events dominate.
- `hold_6_inside` has cleaner reaction rates but fails the weakest-split sample
  threshold after reindexing.

## Cohort Check

No `below_mid` or `above_mid` event-side cohort reaches `100` candidates in the
weakest period split. The largest primary side cohorts are only `58` for
`hold_3_inside` + `mid_touch` and `55` for `hold_3_inside` +
`mid_close_across`, so side-specific conclusions remain diagnostic.

The strongest qualifying narrower row is the event close-position `mid_50`
bucket for `hold_3_inside` + `mid_touch`:

| Horizon | Min Cand | Min Persist | Max Quick Inv | Max Trend | Min Close Back | Min Mid Reject Before Boundary | Max Boundary Before Reject |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 3 | 104 | 88.46% | 11.54% | 8.38% | 50.96% | 48.08% | 18.27% |
| 6 | 104 | 74.04% | 11.54% | 19.23% | 58.65% | 50.96% | 22.12% |
| 12 | 104 | 49.04% | 11.54% | 36.54% | 71.15% | 57.69% | 30.77% |

This bucket improves quick invalidation and boundary-first risk without
collapsing the weakest split below `100` candidates. It should be the first
filter tested in the prototype.

`hold_3_inside_mid_50` is still diagnostic after reindexing: its weakest split
falls to `62` `mid_touch` candidates and `56` `mid_close_across` candidates.

## Project State

Current project state remains review-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The generated audit run reported `strategy=empty trades=0`.
- No Go API, result schema, strategy code, entries, exits, scoring, sizing, or
  strategy replacement changed for this review.

## Conclusion

Build the first minimal offline entry prototype only around:

- `hold_3_inside`
- first `mid_touch` within `12` bars after the hold decision
- event close-position bucket `mid_50` as the primary filter
- closed-candle event decision with any trade entered no earlier than the next
  bar open

The prototype must report side splits and stress splits before any promotion
claim. If that prototype fails P&L or split-stability checks, stop mining this
hold-inside/midline detector family and pivot to a materially different
non-trading hypothesis.
