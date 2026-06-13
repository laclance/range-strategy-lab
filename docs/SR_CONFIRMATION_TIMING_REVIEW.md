# SR Confirmation Timing Review

Date: 2026-06-13

## Verdict

Delayed confirmation after an SR rejection candle is not entry-ready; keep
`lab.EmptyStrategy`.

The audit re-indexes the confirmation candle as the decision candle and starts
all `label_*` outcomes after that confirmation candle. That makes it useful as
a closed-candle timing inspection, but it does not produce a robust first-entry
template. Broad side/delay/horizon aggregates show only tiny
favorable-minus-adverse lift, favorable-greater-than-adverse rates are mostly
near or below `50%`, and the best candidate cohorts collapse as the per-split
row threshold rises.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this audit.

## Inputs Reviewed

Generated audit paths:

- `results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv`
- `results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv`

Audit size:

- candidate rows: `9,692`
- summary rows: `72`
- candidate CSV lines including header: `9,693`
- summary CSV lines including header: `73`
- confirmation delays: `1`, `2`, `3` bars
- forward horizons: `1`, `3`, `6`, `12` bars

The `label_*` columns are forward outcomes only. They are not decision inputs.

## Side, Delay, And Horizon Summary

Rows below aggregate the three period splits for each side, confirmation delay,
and horizon. `diff` is decision-candidate label favorable-minus-adverse in
basis points.

| Side | Delay | Horizon | Candidates | Decision Candidates | Decision Rate | Diff | FGTA | Rej |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| support | 1 | 1 | 27,474 | 13,478 | 49.06% | +0.18bp | 49.32% | 13.13% |
| support | 1 | 3 | 27,474 | 13,478 | 49.06% | +0.16bp | 48.42% | 27.60% |
| support | 1 | 6 | 27,474 | 13,478 | 49.06% | +0.17bp | 48.56% | 36.79% |
| support | 1 | 12 | 27,474 | 13,478 | 49.06% | +0.62bp | 49.15% | 41.72% |
| support | 2 | 1 | 27,474 | 13,554 | 49.33% | +0.20bp | 48.24% | 13.49% |
| support | 2 | 3 | 27,474 | 13,554 | 49.33% | +0.19bp | 47.68% | 27.83% |
| support | 2 | 6 | 27,474 | 13,554 | 49.33% | +0.34bp | 48.04% | 37.35% |
| support | 2 | 12 | 27,474 | 13,554 | 49.33% | +0.80bp | 48.61% | 43.03% |
| support | 3 | 1 | 27,474 | 13,461 | 49.00% | +0.15bp | 48.05% | 13.40% |
| support | 3 | 3 | 27,474 | 13,461 | 49.00% | +0.28bp | 47.54% | 28.33% |
| support | 3 | 6 | 27,474 | 13,461 | 49.00% | +0.46bp | 47.97% | 38.38% |
| support | 3 | 12 | 27,474 | 13,461 | 49.00% | +0.77bp | 48.60% | 44.22% |
| resistance | 1 | 1 | 27,056 | 13,221 | 48.87% | +0.40bp | 49.65% | 12.72% |
| resistance | 1 | 3 | 27,056 | 13,221 | 48.87% | +0.60bp | 49.03% | 28.10% |
| resistance | 1 | 6 | 27,056 | 13,221 | 48.87% | +0.53bp | 48.33% | 37.00% |
| resistance | 1 | 12 | 27,056 | 13,221 | 48.87% | +0.80bp | 48.91% | 41.46% |
| resistance | 2 | 1 | 27,056 | 13,408 | 49.56% | +0.34bp | 49.40% | 13.38% |
| resistance | 2 | 3 | 27,056 | 13,408 | 49.56% | +0.39bp | 48.11% | 28.51% |
| resistance | 2 | 6 | 27,056 | 13,408 | 49.56% | +0.46bp | 48.52% | 37.84% |
| resistance | 2 | 12 | 27,056 | 13,408 | 49.56% | +0.76bp | 48.64% | 43.29% |
| resistance | 3 | 1 | 27,056 | 13,414 | 49.58% | +0.29bp | 48.93% | 13.75% |
| resistance | 3 | 3 | 27,056 | 13,414 | 49.58% | +0.50bp | 48.02% | 28.93% |
| resistance | 3 | 6 | 27,056 | 13,414 | 49.58% | +0.79bp | 48.58% | 38.85% |
| resistance | 3 | 12 | 27,056 | 13,414 | 49.58% | +0.96bp | 48.90% | 44.80% |

The broad grid is small-positive by diff, but not convincing: the best
aggregate is still under `+1bp`, and FGTA stays mostly below `50%`.

## Split Stability

Recent split support was flat or negative in several rows even when the broad
aggregate was positive:

| Split | Side | Delay | Horizon | Diff | FGTA | Rej |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| 2025_2026_recent | support | 1 | 12 | -0.02bp | 49.34% | 41.12% |
| 2025_2026_recent | support | 3 | 1 | -0.09bp | 48.01% | 7.69% |
| 2025_2026_recent | support | 3 | 3 | -0.06bp | 47.58% | 21.74% |
| 2025_2026_recent | support | 3 | 12 | -0.07bp | 47.96% | 42.60% |
| 2025_2026_recent | resistance | 3 | 3 | -0.02bp | 47.48% | 22.15% |
| 2025_2026_recent | resistance | 3 | 6 | -0.09bp | 47.35% | 33.94% |

This does not satisfy the project's split-stability gate for a first entry.

## Candidate Cohorts

The best split-stable cohorts with at least `250` decision-confirmation rows in
each split were still narrow and fragile:

| Side | Delay | Horizon | Seed Close | Confirm Close | Strength | Distance | Min Split Count | Min Diff | Min FGTA |
| --- | ---: | ---: | --- | --- | --- | --- | ---: | ---: | ---: |
| resistance | 1 | 12 | in zone below boundary | below zone | 3 | 0-5bp | 300 | +2.06bp | 49.9% |
| support | 3 | 12 | above zone | above zone | 2 | 5-10bp | 341 | +1.62bp | 49.5% |
| resistance | 2 | 12 | in zone below boundary | below zone | 3 | 0-5bp | 335 | +1.52bp | 48.0% |
| resistance | 3 | 12 | below zone | below zone | 3 | 5-10bp | 276 | +1.26bp | 51.1% |
| resistance | 3 | 12 | in zone below boundary | below zone | 3 | 0-5bp | 337 | +1.25bp | 50.3% |

The threshold stress is worse:

- with at least `250` rows in every split, there were `152` stable cohorts, but
  top rows were small and often near coin-flip by FGTA
- with at least `500` rows in every split, only `16` stable cohorts remained,
  and the best minimum split diff fell to about `+0.16bp`
- with at least `750` rows in every split, no stable cohorts remained

## Conclusion

This audit should be treated as a useful closed-candle confirmation inspection,
not as an entry signal. It improves the timing framing versus the earlier
rejection audit, but the observed edge is too small, too threshold-sensitive,
and too fragile by split for first trade logic.

Do not add a narrower confirmation follow-up unless the next hypothesis changes
materially. A next non-trading branch should test a different idea rather than
continue broad rejection-confirmation slicing.
