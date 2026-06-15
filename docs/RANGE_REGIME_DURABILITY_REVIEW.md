# Range Regime Durability Review

Date: 2026-06-15

## Verdict

The current balanced detector regimes are not durable enough to use as context
for future entry hypotheses; keep `lab.EmptyStrategy`.

The regimes are reasonably split-stable, but mainly because the weakness repeats
across splits. The current detector marks many episodes that stop behaving like
durable ranges almost immediately after the episode ends. This is a
detector/context quality problem, not a reason to test a new entry trigger.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this audit.

## Inputs Reviewed

Generated audit paths:

- `results/range-regime-durability-audit/range_regime_durability_summary.csv`
- `results/range-regime-durability-audit/range_regime_durability_episodes.csv`

Audit size:

- episode rows: `11,984`
- unique episodes: `2,996`
- summary rows: `452`
- episode CSV lines including header: `11,985`
- summary CSV lines including header: `453`
- detector profile: `p30_c12_bollinger_on_adx_off`
- quick invalidation window: `3` bars after the episode end
- forward horizons: `1`, `3`, `6`, `12` bars after the episode end

Audit semantics:

- Episodes are contiguous `RawActive` detector runs that eventually become
  `Active`.
- Episode high/low, width, length, and ATR context use only closed candles
  through the episode end.
- All `label_*` fields are forward outcomes starting at
  `episode_end_index + 1`; they are labels only, not decision inputs.

## Broad Durability

Rows below aggregate the reviewed episode rows by split and horizon.

| Split | Horizon | Episodes | Reentered | Persisted Inside | Quick Invalidated | Chopped | Trended Up | Trended Down |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 2021_2022_stress | 1 | 1,199 | 44.79% | 44.79% | 55.21% | 0.00% | 26.44% | 27.61% |
| 2021_2022_stress | 3 | 1,199 | 62.05% | 30.28% | 69.72% | 17.93% | 27.11% | 24.60% |
| 2021_2022_stress | 6 | 1,199 | 71.89% | 23.60% | 69.72% | 27.02% | 24.52% | 24.85% |
| 2021_2022_stress | 12 | 1,199 | 78.57% | 15.68% | 69.72% | 32.53% | 26.61% | 25.19% |
| 2023_2024_oos | 1 | 1,055 | 43.13% | 43.13% | 56.87% | 0.00% | 29.76% | 25.59% |
| 2023_2024_oos | 3 | 1,055 | 61.33% | 29.57% | 70.43% | 19.43% | 28.15% | 22.65% |
| 2023_2024_oos | 6 | 1,055 | 71.28% | 22.09% | 70.43% | 27.96% | 27.39% | 22.37% |
| 2023_2024_oos | 12 | 1,055 | 78.67% | 15.07% | 70.43% | 33.18% | 27.01% | 24.74% |
| 2025_2026_recent | 1 | 742 | 46.63% | 46.63% | 53.37% | 0.00% | 26.01% | 25.34% |
| 2025_2026_recent | 3 | 742 | 65.09% | 29.92% | 70.08% | 19.81% | 25.20% | 25.07% |
| 2025_2026_recent | 6 | 742 | 74.53% | 22.91% | 70.08% | 28.71% | 24.53% | 23.85% |
| 2025_2026_recent | 12 | 742 | 81.27% | 13.61% | 70.08% | 35.31% | 24.39% | 26.55% |
| full_2021_2026 | 1 | 2,996 | 44.66% | 44.66% | 55.34% | 0.00% | 27.50% | 26.34% |
| full_2021_2026 | 3 | 2,996 | 62.55% | 29.94% | 70.06% | 18.93% | 27.00% | 24.03% |
| full_2021_2026 | 6 | 2,996 | 72.33% | 22.90% | 70.06% | 27.77% | 25.53% | 23.73% |
| full_2021_2026 | 12 | 2,996 | 79.27% | 14.95% | 70.06% | 33.44% | 26.20% | 25.37% |

The full sample is the clearest problem statement: persistence falls from
`44.66%` at `1` bar to `14.95%` at `12` bars, while quick invalidation reaches
`70.06%` by the `3` bar horizon and remains there. Re-entry rises with horizon,
but that is mostly re-entry after an invalidation, not clean range persistence.

## Slice Stability

The bucketed summary does not reveal a durable slice that is strong enough to
serve as future entry context:

| Minimum episodes in every period split | Fully specified bucket rows |
| ---: | ---: |
| 10 | 52 |
| 25 | 20 |
| 50 | 12 |
| 100 | 0 |
| 250 | 0 |
| 500 | 0 |

The best `12` bar fully specified slice with at least `25` episodes in every
period split was still not good enough:

| Slice | Total | Min Split Count | Min Persisted Inside | Max Quick Invalidated | Max Trended |
| --- | ---: | ---: | ---: | ---: | ---: |
| raw `48plus`, active `48plus`, width `gt_50bp`, width/ATR `gt_4x` | 283 | 75 | 23.93% | 58.12% | 42.73% |

Coarser slices were also too broad. For example, at `12` bars the width-only
`gt_50bp` slice had at least `258` episodes in every split, but only `16.80%`
minimum persistence, `68.07%` maximum quick invalidation, and `51.37%` maximum
trend behavior.

## Project State

Current project state remains audit-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The latest durability audit smoke reported `strategy=empty trades=0`.
- No strategy code, public APIs, result files, entries, exits, scoring, sizing,
  or strategy replacement were changed for this review.

## Conclusion

This audit should be treated as a detector/context quality review, not an entry
readiness review. The current balanced compression detector is too broad and
too quick to invalidate after episode end to be trusted as context for a new
entry trigger.

The next implementation should refine or reframe the detector/context first,
then review the new durability outputs before any entry hypothesis is tested.
