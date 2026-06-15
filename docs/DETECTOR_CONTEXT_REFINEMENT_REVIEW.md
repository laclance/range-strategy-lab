# Detector Context Refinement Review

Date: 2026-06-15

## Verdict

The delayed "hold inside" context rules are the first context refinement that
**materially and split-stably** reduces quick invalidation and trend leakage,
with adequate candidate counts. Lead rules: `hold_3_inside` and `hold_6_inside`.
Carry them forward as the leading decision-candle context.

They are **not** promoted to entry context, and no detector profile or context
rule is approved for entries. Keep `lab.EmptyStrategy`. The improvement is a
heavy survivorship/conditioning effect, residual trend leakage is still material
at the `12` bar horizon, and these are regime-durability labels, not P&L.

No entries, exits, scoring, sizing, or strategy replacement should be added from
this review.

## Inputs Reviewed

Generated detector context refinement audit paths:

- `results/detector-context-refinement-audit/detector_context_refinement_candidates.csv`
- `results/detector-context-refinement-audit/detector_context_refinement_candidates.json`
- `results/detector-context-refinement-audit/detector_context_refinement_summary.csv`
- `results/detector-context-refinement-audit/detector_context_refinement_summary.json`
- `results/detector-context-refinement-audit/detector_context_refinement_stability.csv`
- `results/detector-context-refinement-audit/detector_context_refinement_stability.json`

Audit size:

- profiles: `8`
- context rules: `5`
- candidate rows: `113,824`
- summary rows: `640`
- stability rows: `160`
- candidate CSV lines including header: `113,825`
- summary CSV lines including header: `641`
- stability CSV lines including header: `161`
- horizons: `1`, `3`, `6`, `12`
- quick invalidation window: `3` bars after the decision candle

Review semantics:

- Summary rows are one row per profile, context rule, split, and horizon.
- Stability rows compare `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent` only (`min`, `max`, `delta` across the three).
- `episode_end` decides at the raw-active episode end; `hold_N_inside` rules
  delay the decision to `episode_end_index + N` and require every close in
  `episode_end_index+1 .. decision_index` to be inside the frozen episode range;
  `hold_3_inside_mid_50` additionally requires the decision close in the middle
  `50%` of the range.
- The hold condition is closed-candle knowable at the decision candle. All
  `label_*` fields start at `decision_index + 1` and are forward outcomes only.
- Source episode counts stay in the summary denominators even when a context
  rule rejects the episode, so `candidate_rate` measures survivorship.

## Context Rule Effect (balanced baseline)

For `p30_c12_bollinger_on_adx_off`, the rules sharply improve the durability
surface. Full-sample (`full_2021_2026`) at the `12` bar horizon:

| Rule | Cand | Cand Rate | H12 Persisted | H12 Quick Invalidated | H12 Trended | H12 Chopped |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `episode_end` | 2,996 | 100.0% | 15.0% | 70.1% | 51.6% | 33.4% |
| `hold_1_inside` | 1,338 | 44.7% | 32.0% | 39.9% | 38.6% | 29.4% |
| `hold_3_inside` | 897 | 29.9% | 43.0% | 23.5% | 35.2% | 21.7% |
| `hold_6_inside` | 686 | 22.9% | 48.4% | 19.7% | 32.5% | 19.1% |
| `hold_3_inside_mid_50` | 414 | 13.8% | 49.0% | 14.5% | 31.6% | 19.3% |

The same monotonic improvement holds across the period splits. Worst-split
(`min`/`max` across the three splits) at the `12` bar horizon:

| Rule | Min Split Cand | H12 Min Persisted | H12 Max Quick Invalidated | H12 Max Trended |
| --- | ---: | ---: | ---: | ---: |
| `episode_end` | 742 | 13.6% | 70.4% | 51.8% |
| `hold_1_inside` | 346 | 27.7% | 42.2% | 40.8% |
| `hold_3_inside` | 222 | 40.5% | 25.3% | 40.1% |
| `hold_6_inside` | 170 | 44.7% | 21.5% | 35.9% |
| `hold_3_inside_mid_50` | 94 | 47.8% | 16.3% | 33.1% |

This is the first reviewed refinement that moves worst-split quick invalidation
from roughly `70%` down to roughly `16%-25%` while lifting worst-split
persistence from `13.6%` to over `40%`, without collapsing the sample below
useful counts for `hold_3_inside` and `hold_6_inside`.

## Horizon Shape

The gain decays with horizon, as expected, but stays meaningful out to `12`
bars. Full-sample `hold_3_inside`:

| Horizon | Persisted | Quick Invalidated | Trended |
| ---: | ---: | ---: | ---: |
| 1 | 89.6% | 10.4% | 9.9% |
| 3 | 76.5% | 23.5% | 16.5% |
| 6 | 61.4% | 23.5% | 24.9% |
| 12 | 43.0% | 23.5% | 35.2% |

Quick invalidation is fixed at the `3` bar post-decision window, so it is flat
from the `3` bar horizon onward; persistence erodes and trend leakage grows as
the horizon extends.

## Cross-Profile Picture

The effect is not specific to the baseline. The strongest worst-split `12` bar
rows with at least `150` candidates in every split are all delayed-hold rules:

| Profile And Rule | Min Split Cand | H12 Min Persisted | H12 Max Quick Invalidated | H12 Max Trended |
| --- | ---: | ---: | ---: | ---: |
| `p40_c24_bollinger_on_adx_off`, `hold_3_inside` | 181 | 46.8% | 27.0% | 31.5% |
| `p30_c12_bollinger_on_adx_off`, `hold_6_inside` | 170 | 44.7% | 21.5% | 35.9% |
| `p40_c24_bollinger_on_adx_on`, `hold_3_inside` | 164 | 42.9% | 24.4% | 32.4% |
| `p30_c12_bollinger_on_adx_off`, `hold_3_inside` | 222 | 40.5% | 25.3% | 40.1% |

`mid_50` pushes persistence and quick invalidation slightly further but drops to
roughly `94`-`111` candidates in the weakest split, so it is a promising tighten
rather than a primary rule at current sample sizes.

## Why This Is Not Yet Entry Context

1. Survivorship/conditioning. The rule keeps only the episodes that already held
   inside for `N` bars: `hold_3_inside` keeps `~30%` of episodes,
   `hold_6_inside` `~22%-24%`, `hold_3_inside_mid_50` `~14%`. Much of the gain
   is persistence-of-persistence, and the rule both delays the decision and
   shrinks the sample.
2. Residual trend leakage. Even after the filter, worst-split trended rate at
   `12` bars is still `~32%-40%` and persistence is only `~40%-48%` — close to a
   coin flip on `12` bar range survival.
3. Durability, not edge. The `label_*` fields measure whether the range
   persists, quick-invalidates, or trends. They are not directional P&L.
   Persistence of a range is not yet a profitable mean-reversion entry.

## Project State

Current project state remains audit-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The detector context refinement audit smoke reported `strategy=empty trades=0`.
- No Go API, CLI flag, result schema, strategy code, entries, exits, scoring,
  sizing, or strategy replacement changed for this review.

## Conclusion

This audit produced the first genuine, split-stable context improvement in the
project: requiring a few closed candles of hold inside the frozen episode range
materially reduces quick invalidation and trend leakage. Because the hold
condition is closed-candle knowable at the decision candle, `hold_3_inside` and
`hold_6_inside` are legitimate candidate decision-candle context, not lookahead.

The next implementation should stay non-trading and re-express the leading
hold-inside context as a forward directional/mean-reversion edge audit: condition
on `hold_3_inside`/`hold_6_inside` at the decision candle and measure forward,
split-stable outcomes (for example reversion toward range mid versus continued
trend), before any entry trigger, scoring, or sizing is added.
