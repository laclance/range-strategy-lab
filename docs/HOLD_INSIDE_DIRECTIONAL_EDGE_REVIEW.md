# Hold-Inside Directional Edge Review

Date: 2026-06-15

## Verdict

Hold-inside directional edge audit is not entry-ready; keep
`lab.EmptyStrategy`.

The delayed hold-inside rules remain useful decision-candle context for range
survival, but this directional paper-side audit does not show a split-stable
edge toward the frozen range high or low. The strongest broad rows either have
negative worst-split favorable-minus-adverse, favorable-greater-than-adverse
below `50%`, or both. Close-position buckets do not rescue the result: no
non-`all` bucket reaches `100` candidates in every period split, and the best
relaxed `50` candidate rows are too sparse, too low by FGTA, or too exposed to
quick invalidation.

No entries, exits, scoring, sizing, or strategy replacement should be added from
this audit. Do not convert `paper_side=toward_high` or `paper_side=toward_low`
labels into trades.

## Inputs Reviewed

Generated hold-inside directional edge audit paths:

- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv`
- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.json`
- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv`
- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.json`
- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv`
- `results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.json`

Audit size:

- profiles: `1`
- context rules: `3`
- paper sides: `2`
- candidate rows: `15,976`
- summary rows: `624`
- stability rows: `168`
- candidate CSV lines including header: `15,977`
- summary CSV lines including header: `625`
- stability CSV lines including header: `169`
- horizons: `1`, `3`, `6`, `12`
- quick invalidation window: `3` bars after the decision candle

Review semantics:

- Detector profile is only `p30_c12_bollinger_on_adx_off`.
- Context rules are `hold_3_inside`, `hold_6_inside`, and
  `hold_3_inside_mid_50`.
- Candidate rows are one row per passed source episode, context rule, horizon,
  and paper side.
- `toward_high` treats upward excursion as favorable and downward excursion as
  adverse.
- `toward_low` treats downward excursion as favorable and upward excursion as
  adverse.
- Decision fields use only data known at the decision candle: frozen episode
  high/low/mid, decision close position, distance to high/low/mid, raw/active
  length, width, ATR, and width/ATR context.
- All `label_*` fields start at `decision_index + 1` and are forward outcomes
  only. They are not decision inputs.
- Stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent`.

## Promotion Gate

For this review, a row could only be carried forward as directional context if
`hold_3_inside` or `hold_6_inside` passed all of these:

- all-bucket row, or a decision-close-position bucket with at least `100`
  candidates in every period split
- positive worst-split average favorable-minus-adverse
- worst-split favorable-greater-than-adverse above `50%`
- stable support in all three period splits
- no obvious contradiction from quick invalidation or opposite close-break rates

`hold_3_inside_mid_50` is diagnostic only because it is a stricter comparison
with smaller samples. Bucket rows below `100` candidates per split are
diagnostic only.

## Primary Stability Rows

Worst-split rows for the two primary context rules:

| Rule | Horizon | Side | Min Split Cand | Min FMA | Min FGTA | Min Side Touch | Max Opp Break | Max Quick Inv |
| --- | ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | 1 | `toward_high` | 222 | -0.34bp | 47.44% | 8.33% | 6.73% | 11.86% |
| `hold_3_inside` | 1 | `toward_low` | 222 | -1.11bp | 49.86% | 8.82% | 5.41% | 11.86% |
| `hold_3_inside` | 3 | `toward_high` | 222 | +0.78bp | 46.56% | 16.99% | 12.18% | 25.32% |
| `hold_3_inside` | 3 | `toward_low` | 222 | -2.66bp | 46.85% | 17.95% | 13.14% | 25.32% |
| `hold_3_inside` | 6 | `toward_high` | 222 | +0.62bp | 45.51% | 24.77% | 21.62% | 25.32% |
| `hold_3_inside` | 6 | `toward_low` | 222 | -1.96bp | 46.85% | 26.28% | 20.83% | 25.32% |
| `hold_3_inside` | 12 | `toward_high` | 222 | -0.98bp | 46.79% | 33.97% | 31.98% | 25.32% |
| `hold_3_inside` | 12 | `toward_low` | 222 | -5.37bp | 49.10% | 37.50% | 29.73% | 25.32% |
| `hold_6_inside` | 1 | `toward_high` | 170 | -0.28bp | 47.06% | 4.24% | 5.88% | 10.73% |
| `hold_6_inside` | 1 | `toward_low` | 170 | -0.61bp | 46.78% | 6.44% | 6.01% | 10.73% |
| `hold_6_inside` | 3 | `toward_high` | 170 | -0.80bp | 47.65% | 11.76% | 12.35% | 21.46% |
| `hold_6_inside` | 3 | `toward_low` | 170 | -0.39bp | 49.36% | 15.02% | 10.30% | 21.46% |
| `hold_6_inside` | 6 | `toward_high` | 170 | -2.17bp | 45.88% | 18.88% | 22.35% | 21.46% |
| `hold_6_inside` | 6 | `toward_low` | 170 | -0.03bp | 50.18% | 21.89% | 18.24% | 21.46% |
| `hold_6_inside` | 12 | `toward_high` | 170 | -1.10bp | 46.35% | 31.33% | 30.00% | 21.46% |
| `hold_6_inside` | 12 | `toward_low` | 170 | -5.31bp | 45.23% | 33.05% | 27.06% | 21.46% |

No primary all-bucket row passes the gate. `hold_3_inside` toward high at the
`3` and `6` bar horizons has positive worst-split FMA, but worst-split FGTA is
only `46.56%` and `45.51%`. `hold_6_inside` toward low at the `6` bar horizon
is the opposite shape: worst-split FGTA is `50.18%`, but worst-split FMA is
slightly negative at `-0.03bp`.

## Split Detail

The closest broad candidates fail for different reasons:

| Row | Stress | OOS | Recent | Failure |
| --- | ---: | ---: | ---: | --- |
| `hold_3_inside`, h3, `toward_high` | +2.66bp / 46.56% / n=363 | +0.78bp / 48.08% / n=312 | +0.90bp / 53.15% / n=222 | FMA positive, but two splits are below `50%` FGTA |
| `hold_3_inside`, h6, `toward_high` | +1.96bp / 47.93% / n=363 | +0.62bp / 45.51% / n=312 | +1.31bp / 53.15% / n=222 | FMA positive, but two splits are below `50%` FGTA |
| `hold_6_inside`, h6, `toward_low` | +0.08bp / 50.18% / n=283 | +2.17bp / 52.79% / n=233 | -0.03bp / 54.12% / n=170 | FGTA clears `50%`, but recent FMA is negative |
| `hold_3_inside_mid_50`, h3, `toward_low` | -0.09bp / 53.93% / n=178 | -0.16bp / 55.63% / n=142 | +1.15bp / 52.13% / n=94 | diagnostic-only rule; two splits have negative FMA |

The full-sample averages can look better than the split-stability result. For
example, full-sample `hold_3_inside` toward high is `+1.57bp` at `3` bars and
`+1.34bp` at `6` bars, but the worst-split FGTA remains below the gate. This is
why broad positivity is insufficient here.

## Decision-Close Buckets

Decision-close-position buckets do not create a promotable row:

| Threshold | Non-`all` bucket rows for `hold_3_inside`/`hold_6_inside` |
| ---: | ---: |
| `>=100` candidates in every split | 0 |
| `>=50` candidates in every split | 32 |

Best relaxed rows at `>=50` candidates per split:

| Rule | Horizon | Side | Bucket | Min Split Cand | Min FMA | Min FGTA | Min Side Touch | Max Quick Inv |
| --- | ---: | --- | --- | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | 3 | `toward_high` | `high_25` | 54 | +2.76bp | 39.08% | 43.18% | 37.04% |
| `hold_3_inside` | 6 | `toward_high` | `high_25` | 54 | +1.14bp | 43.68% | 54.55% | 37.04% |
| `hold_3_inside` | 1 | `toward_high` | `high_25` | 54 | +0.81bp | 44.32% | 25.00% | 22.22% |
| `hold_3_inside` | 1 | `toward_low` | `mid_50` | 94 | +0.40bp | 50.56% | 2.13% | 5.62% |
| `hold_6_inside` | 6 | `toward_low` | `mid_50` | 98 | +0.04bp | 49.67% | 15.69% | 10.46% |

These rows are diagnostic. The `high_25` rows have positive FMA but poor FGTA,
small samples, and high quick invalidation. The `mid_50` rows are cleaner on
quick invalidation but either too small, too close to flat, or below `50%`
FGTA.

## Project State

Current project state remains audit-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The latest hold-inside directional edge audit smoke reported
  `strategy=empty trades=0`.
- No Go API, CLI flag, result schema, strategy code, entries, exits, scoring,
  sizing, or strategy replacement changed for this review.

## Conclusion

The delayed hold-inside context still matters: it reduces quick invalidation and
creates a cleaner decision candle than raw episode end. But "after holding
inside, paper-side toward high or low has directional edge" is not supported by
the current evidence.

Do not continue narrower high/low paper-side slicing unless the hypothesis
changes materially. The next implementation should stay non-trading and test a
different hold-inside question, such as whether closed-candle midline transition
after `hold_3_inside`/`hold_6_inside` identifies reversion, chop, or trend
continuation better than raw toward-high/toward-low labels.
