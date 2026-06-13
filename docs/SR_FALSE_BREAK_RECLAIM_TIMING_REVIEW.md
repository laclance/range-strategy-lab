# SR False-Break Reclaim Timing Review

Date: 2026-06-14

## Verdict

False-break reclaim timing audit is not entry-ready; keep
`lab.EmptyStrategy`.

The audit has a useful closed-candle framing: the reclaim candle is the
decision candle, and all `label_*` outcomes start after that reclaim candle.
But the evidence does not show a simple, split-stable first-entry template.
Broad side/horizon aggregates are small-positive, recent split support is
negative across all horizons, OOS resistance is negative at the longer
horizons, and candidate cohorts either become sparse or degrade when row
thresholds are stressed.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this audit.

## Inputs Reviewed

Generated audit paths:

- `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv`
- `results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv`

Audit size:

- candidate rows: `17,652`
- summary rows: `24`
- candidate CSV lines including header: `17,653`
- summary CSV lines including header: `25`
- max break delay: `3` bars after the anchor candle
- max reclaim delay: `12` bars after the break candle
- forward horizons: `1`, `3`, `6`, `12` bars after the reclaim candle
- support reclaim decisions across all splits: `4,120` per horizon
- resistance reclaim decisions across all splits: `4,150` per horizon

The `label_*` columns are forward outcomes only. They are not decision inputs.

Audit semantics:

- The anchor boundary comes from existing SR audit rows.
- Support false break closes below the frozen anchor support zone bottom, then
  reclaims with a close at or above the frozen anchor support level.
- Resistance false break closes above the frozen anchor resistance zone top,
  then reclaims with a close at or below the frozen anchor resistance level.
- The reclaim candle is the decision candle.

## Side And Horizon Summary

Rows below aggregate the three period splits for each side and horizon,
weighted by `candidate_count`. `diff` is label favorable-minus-adverse in
basis points.

| Side | Horizon | Candidates | Diff | FGTA | Rej |
| --- | ---: | ---: | ---: | ---: | ---: |
| support | 1 | 4,120 | +1.03bp | 49.85% | 22.18% |
| support | 3 | 4,120 | +1.76bp | 51.09% | 37.99% |
| support | 6 | 4,120 | +1.61bp | 49.27% | 40.07% |
| support | 12 | 4,120 | +0.75bp | 50.07% | 36.67% |
| resistance | 1 | 4,150 | +0.84bp | 49.49% | 22.36% |
| resistance | 3 | 4,150 | +1.28bp | 49.90% | 36.67% |
| resistance | 6 | 4,150 | +1.85bp | 50.72% | 39.90% |
| resistance | 12 | 4,150 | +2.41bp | 51.54% | 37.54% |

The broad grid is not enough for entry readiness. The best aggregate is only
`+2.41bp`, and FGTA is mostly near coin-flip.

## Split Stability

Cells are `diff / FGTA / n`.

| Split | Side | H1 | H3 | H6 | H12 |
| --- | --- | ---: | ---: | ---: | ---: |
| 2021_2022_stress | support | +2.07bp / 50.94% / n=1,710 | +3.04bp / 52.69% / n=1,710 | +3.24bp / 49.77% / n=1,710 | +0.65bp / 50.06% / n=1,710 |
| 2021_2022_stress | resistance | +1.21bp / 48.89% / n=1,661 | +2.52bp / 50.33% / n=1,661 | +3.42bp / 50.33% / n=1,661 | +3.99bp / 50.33% / n=1,661 |
| 2023_2024_oos | support | +0.57bp / 49.34% / n=1,445 | +1.70bp / 51.00% / n=1,445 | +1.46bp / 50.03% / n=1,445 | +2.16bp / 51.76% / n=1,445 |
| 2023_2024_oos | resistance | +0.24bp / 48.69% / n=1,485 | -0.38bp / 48.48% / n=1,485 | -0.72bp / 48.22% / n=1,485 | -1.00bp / 51.78% / n=1,485 |
| 2025_2026_recent | support | -0.12bp / 48.70% / n=965 | -0.42bp / 48.39% / n=965 | -1.03bp / 47.25% / n=965 | -1.20bp / 47.56% / n=965 |
| 2025_2026_recent | resistance | +1.12bp / 51.69% / n=1,004 | +1.68bp / 51.29% / n=1,004 | +3.04bp / 55.08% / n=1,004 | +4.86bp / 53.19% / n=1,004 |

Recent support does not hold up at any horizon. OOS resistance also turns
negative at `3`, `6`, and `12` bars. That side/split mismatch is enough to
block first-entry work.

## Candidate Cohorts

Thresholds below require every split to have at least the listed number of
candidates in the same cohort.

| Cohort Grouping | >=25 | >=50 | >=100 | >=250 | >=500 |
| --- | ---: | ---: | ---: | ---: | ---: |
| full feature: side, horizon, delays, close locations, move buckets, strength, distance | 0 | 0 | 0 | 0 | 0 |
| coarse: side, horizon, strength, distance | 48 | 48 | 36 | 8 | 0 |
| close shape: side, horizon, anchor/break/reclaim close locations | 32 | 32 | 24 | 8 | 8 |

Best stable examples:

| Screen | Min Split Rows | Best Cohort | Total | Min Count | Min Diff | Min FGTA | Note |
| --- | ---: | --- | ---: | ---: | ---: | ---: | --- |
| coarse | 100 | resistance h12, strength 4plus, distance 5-10bp | 472 | 120 | +4.22bp | 51.35% | positive, but sparse and side-specific |
| coarse | 250 | support h3, strength 4plus, distance 0-5bp | 1,023 | 292 | +1.00bp | 50.00% | larger, but too weak |
| close shape | 100 | resistance h12, below_zone to above_zone to in_zone_below_boundary | 897 | 209 | +2.98bp | 51.40% | does not survive the 250-row threshold |
| close shape | 500 | resistance h6, in_zone_below_boundary to above_zone to in_zone_below_boundary | 1,774 | 509 | +0.51bp | 46.69% | enough rows, but not enough edge |

The fully specified false-break reclaim shape has no stable cohort even at
`25` candidates per split. Coarser shapes can find positive rows, but the
stronger rows are sparse, and the larger rows are small or near coin-flip by
FGTA.

## Conclusion

This audit should be treated as a useful closed-candle false-break reclaim
inspection, not as an entry signal. The reclaim framing is cleaner than an
ex-post boundary label, but the observed edge is too small, too side-specific,
and too fragile by split for first trade logic.

Do not continue narrow SR timing slices unless the next hypothesis changes
materially. Next work should stay non-trading and move to a different idea
before any entries are added.
