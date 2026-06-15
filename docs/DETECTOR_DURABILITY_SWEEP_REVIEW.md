# Detector Durability Sweep Review

Date: 2026-06-15

## Verdict

No profile in the current `DefaultDetectorSweepProfiles` detector durability
sweep is approved as future entry context; keep `lab.EmptyStrategy`.

The ADX comparison profile, `p30_c12_bollinger_on_adx_on`, is a useful
diagnostic clue because it improves short-horizon persistence and lowers quick
invalidation versus the balanced baseline. It is not a promoted detector. At
the `12` bar horizon, persistence is still too low, trend behavior remains too
high, and the best fully specified slices are either too sparse or still
quick-invalidate too often.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this review.

## Inputs Reviewed

Generated detector durability sweep paths:

- `results/detector-durability-sweep/detector_durability_sweep.csv`
- `results/detector-durability-sweep/detector_durability_sweep.json`
- `results/detector-durability-sweep/detector_durability_slices.csv`
- `results/detector-durability-sweep/detector_durability_slices.json`
- `results/detector-durability-sweep/detector_durability_stability.csv`
- `results/detector-durability-sweep/detector_durability_stability.json`

Audit size:

- profiles: `19`
- broad rows: `304`
- slice rows: `9,088`
- stability rows: `76`
- broad CSV lines including header: `305`
- slice CSV lines including header: `9,089`
- stability CSV lines including header: `77`
- horizons: `1`, `3`, `6`, `12`
- quick invalidation window: `3` bars after the episode end

Review semantics:

- Broad rows are one row per detector profile, split, and horizon.
- Slice rows use raw length, active length, width, and width/ATR buckets.
- Stability rows compare `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent`.
- All `label_*` fields are forward outcomes only. They are not decision inputs.

## Broad Profile Stability

The balanced baseline remains weak across splits:

| Profile | Episodes | Min Split Count | Duty Cycle Range | H12 Min Persisted Inside | H12 Max Quick Invalidated | H12 Max Trended |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `p30_c12_bollinger_on_adx_off` | 2,996 | 742 | 12.70%-14.59% | 13.61% | 70.43% | 51.79% |
| `p30_c12_bollinger_on_adx_on` | 1,881 | 494 | 4.25%-4.50% | 16.47% | 56.89% | 50.75% |

The ADX comparison is cleaner at short horizons:

| Horizon | Episodes | Min Split Count | Min Persisted Inside | Max Quick Invalidated | Max Trended |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 1 | 1,881 | 494 | 59.28% | 40.72% | 39.67% |
| 3 | 1,881 | 494 | 43.11% | 56.89% | 40.61% |
| 6 | 1,881 | 494 | 28.59% | 56.89% | 42.51% |
| 12 | 1,881 | 494 | 16.47% | 56.89% | 50.75% |

That shape is not enough to approve it as context. The `12` bar profile still
has only `16.47%` minimum persistence and more than half of episodes
quick-invalidated in the weakest split.

The best `12` bar persistence floors among broad profiles also fail promotion:

| Profile | Episodes | Min Split Count | Duty Cycle Range | H12 Min Persisted Inside | H12 Max Quick Invalidated | H12 Max Trended |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `p20_c24_bollinger_on_adx_off` | 1,013 | 282 | 3.57%-5.71% | 17.57% | 70.57% | 52.13% |
| `p30_c24_bollinger_on_adx_off` | 1,616 | 413 | 7.60%-10.05% | 17.43% | 69.57% | 50.56% |
| `p40_c24_bollinger_on_adx_off` | 2,237 | 545 | 12.43%-15.16% | 17.30% | 66.79% | 48.86% |

Increasing the consecutive-bar requirement improves the persistence floor only
slightly. It does not solve the quick-invalidation or trend-behavior problem.

## Slice Stability

Fully specified `12` bar slices were checked across all three period splits.
Eligible row counts by minimum episodes in every period split were:

| Minimum episodes in every period split | Fully specified H12 slices |
| ---: | ---: |
| 10 | 203 |
| 25 | 80 |
| 50 | 45 |
| 100 | 13 |
| 250 | 0 |
| 500 | 0 |

The best fully specified slices with at least `100` episodes in every period
split still do not define durable context:

| Profile And Slice | Total | Min Split Count | Min Persisted Inside | Max Quick Invalidated | Max Trended |
| --- | ---: | ---: | ---: | ---: | ---: |
| `p40_c24_bollinger_on_adx_off`, raw `48plus`, active `48plus`, width `gt_50bp`, width/ATR `gt_4x` | 391 | 117 | 26.95% | 60.28% | 45.30% |
| `p40_c06_bollinger_on_adx_off`, raw `48plus`, active `48plus`, width `gt_50bp`, width/ATR `gt_4x` | 645 | 167 | 25.65% | 63.57% | 44.98% |
| `p40_c12_bollinger_on_adx_off`, raw `48plus`, active `48plus`, width `gt_50bp`, width/ATR `gt_4x` | 547 | 147 | 25.17% | 62.21% | 46.26% |

The slices are better than the broad detector rows, but they remain weak as
entry context: persistence is still low, quick invalidation is still high, and
no fully specified slice survives a `250` episodes-per-split threshold.

## Project State

Current project state remains audit-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The latest detector durability sweep smoke reported `strategy=empty trades=0`.
- No Go API, CLI flag, result schema, strategy code, entries, exits, scoring,
  sizing, or strategy replacement changed for this review.

## Conclusion

The sweep is useful because it shows that the detector problem is not limited
to one balanced baseline profile. ADX gating and stricter consecutive-bar
profiles can improve parts of the durability surface, but the reviewed grid
still does not produce durable enough range context for first-entry work.

The next implementation should stay non-trading and refine or reframe detector
context before any entry trigger is tested.
