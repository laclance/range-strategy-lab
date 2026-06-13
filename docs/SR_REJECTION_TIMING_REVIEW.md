# SR Rejection Timing Review

Date: 2026-06-13

## Verdict

Boundary-rejection timing audit is not entry-ready; keep
`lab.EmptyStrategy`.

The timing audit does show that closed-candle rejection features can slightly
lift forward rejection rates, but it does not show a simple, support/resistance
symmetric edge that survives all splits. Favorable-minus-adverse differences
are small, often mixed by side, and strongest cohorts depend on split-specific
or pierced-zone details.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this audit alone.

## Inputs Reviewed

Generated audit paths:

- `results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv`
- `results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv`

Audit size:

- candidate rows: `968`
- summary rows: `24`
- candidate CSV lines including header: `969`
- summary CSV lines including header: `25`

The `label_*` columns are forward outcomes only. They are not decision inputs.

## Side And Horizon Summary

Rows below aggregate the three period splits for each side and horizon.
`diff` is label favorable-minus-adverse in basis points.

| Side | Horizon | Candidates | Decision Candidates | All Diff | Decision Diff | All Rej | Decision Rej |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| support | 1 | 35,204 | 27,474 | -0.04bp | +0.02bp | 10.68% | 11.78% |
| support | 3 | 35,204 | 27,474 | -0.13bp | -0.10bp | 24.40% | 25.49% |
| support | 6 | 35,204 | 27,474 | -0.15bp | -0.20bp | 32.79% | 32.88% |
| support | 12 | 35,204 | 27,474 | +0.04bp | -0.06bp | 36.77% | 35.32% |
| resistance | 1 | 35,066 | 27,056 | +0.01bp | +0.06bp | 10.52% | 11.69% |
| resistance | 3 | 35,066 | 27,056 | +0.13bp | +0.24bp | 24.89% | 26.26% |
| resistance | 6 | 35,066 | 27,056 | +0.09bp | +0.21bp | 33.24% | 33.65% |
| resistance | 12 | 35,066 | 27,056 | +0.25bp | +0.45bp | 36.31% | 35.10% |

The broadest support-side decision candidates are negative or flat across most
horizons. Resistance is better, but the aggregate lift is still tiny.

## Candidate Cohorts

Top positive OOS/recent decision-candidate cohorts with at least `250` rows:

| Split | Side | Horizon | Close Location | Pierced | Strength | Distance | Count | Diff | FGTA | Rej |
| --- | --- | ---: | --- | --- | --- | --- | ---: | ---: | ---: | ---: |
| 2023_2024_oos | support | 12 | in zone above boundary | true | 2 | 0-5bp | 250 | +5.05bp | 55.6% | 36.8% |
| 2025_2026_recent | resistance | 12 | below zone | false | 3 | 5-10bp | 529 | +3.44bp | 53.3% | 46.3% |
| 2025_2026_recent | resistance | 12 | in zone below boundary | false | 3 | 0-5bp | 1,150 | +2.99bp | 53.8% | 38.9% |
| 2025_2026_recent | support | 12 | in zone above boundary | false | 2 | 0-5bp | 1,533 | +2.15bp | 54.6% | 40.1% |
| 2025_2026_recent | resistance | 6 | below zone | false | 3 | 5-10bp | 529 | +2.13bp | 54.1% | 40.5% |
| 2023_2024_oos | support | 6 | in zone above boundary | true | 2 | 0-5bp | 250 | +1.95bp | 50.8% | 40.0% |
| 2023_2024_oos | support | 12 | in zone above boundary | false | 4plus | 0-5bp | 3,456 | +1.93bp | 51.3% | 31.7% |
| 2023_2024_oos | support | 12 | above zone | false | 3 | 5-10bp | 744 | +1.86bp | 50.7% | 44.2% |

These rows do not define one clean first-entry template. The largest OOS
support cohort is pierced and barely meets the row threshold; the best recent
resistance cohorts use different close-location and strength settings. The
common in-zone/strength-2 shape is not stable enough:

| Split | Side | Pierced | Count | Diff | FGTA | Rej |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| 2021_2022_stress | support | false | 978 | -4.34bp | 47.14% | 29.65% |
| 2021_2022_stress | support | true | 311 | +1.36bp | 52.73% | 31.51% |
| 2021_2022_stress | resistance | false | 906 | +1.53bp | 52.98% | 32.67% |
| 2021_2022_stress | resistance | true | 328 | +5.81bp | 51.83% | 31.10% |
| 2023_2024_oos | support | false | 1,918 | -0.08bp | 49.22% | 33.16% |
| 2023_2024_oos | support | true | 250 | +5.05bp | 55.60% | 36.80% |
| 2023_2024_oos | resistance | false | 1,920 | +0.14bp | 48.70% | 35.42% |
| 2023_2024_oos | resistance | true | 241 | -1.03bp | 49.38% | 31.12% |
| 2025_2026_recent | support | false | 1,533 | +2.15bp | 54.60% | 40.05% |
| 2025_2026_recent | support | true | 160 | +0.39bp | 53.75% | 41.25% |
| 2025_2026_recent | resistance | false | 1,444 | +0.08bp | 49.38% | 35.53% |
| 2025_2026_recent | resistance | true | 189 | +2.49bp | 52.38% | 29.63% |

## Conclusion

This audit should be treated as a useful rejection-timing inspection, not as an
entry signal. The simplest closed-candle rejection shape is either too weak,
too side-specific, or too dependent on a narrow pierced-zone slice.

Next work should stay non-trading unless the user explicitly asks otherwise.
The best follow-up is another compact audit that tests a materially different
confirmation hypothesis, such as delayed confirmation after a rejection candle
or a false-break reclaim timing audit.
