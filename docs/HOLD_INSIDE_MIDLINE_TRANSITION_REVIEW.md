# Hold-Inside Midline Transition Review

Date: 2026-06-15

## Verdict

Hold-inside midline transition audit is not entry-ready; keep
`lab.EmptyStrategy`.

The leading `hold_3_inside` and `hold_6_inside` context rules do show a
split-stable midline-transition surface. By the `12` bar horizon, at least half
of every period split touches the frozen range mid, and roughly `38%-43%` of
every split touches mid before any boundary touch. That is materially better
than the failed high/low paper-side direction test.

This is still not a trade trigger. The broad rows need several bars for the
midline event to emerge, `12` bar trend leakage remains material, and the
cleanest mid-position slice is below the primary sample threshold in the
weakest split. Treat the midline as a promising non-trading observation point,
not as entry context, scoring, sizing, or strategy logic.

## Inputs Reviewed

Generated hold-inside midline transition audit paths:

- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.csv`
- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.json`
- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.csv`
- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.json`
- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.csv`
- `results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.json`

Audit size:

- profiles: `1`
- context rules: `3`
- candidate rows: `7,988`
- summary rows: `720`
- stability rows: `192`
- candidate CSV lines including header: `7,989`
- summary CSV lines including header: `721`
- stability CSV lines including header: `193`
- horizons: `1`, `3`, `6`, `12`
- quick invalidation window: `3` bars after the decision candle

Review semantics:

- Detector profile is only `p30_c12_bollinger_on_adx_off`.
- Primary context rules are `hold_3_inside` and `hold_6_inside`.
- `hold_3_inside_mid_50` is diagnostic.
- Candidate rows are one row per passed source episode, context rule, horizon,
  decision mid side, and decision close-position bucket.
- Decision fields use only data known at the hold decision candle.
- All `label_*` fields start at `decision_index + 1` and are forward outcomes
  only. They are not decision inputs.
- Stability rows compare only `2021_2022_stress`, `2023_2024_oos`, and
  `2025_2026_recent`.

## Primary Stability Rows

Worst-split rows for the two primary context rules:

| Rule | Horizon | Min Cand | Min Touch Mid | Min Close Across | Min Touch Before Boundary | Min Cross Before Break | Min Persist | Max Quick Inv | Max Trend | Max Cross Delay |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | 3 | 222 | 27.03% | 19.37% | 22.97% | 18.47% | 74.68% | 25.32% | 17.57% | 1.96 |
| `hold_3_inside` | 6 | 222 | 38.29% | 30.63% | 32.43% | 27.93% | 58.65% | 25.32% | 26.58% | 3.14 |
| `hold_3_inside` | 12 | 222 | 52.25% | 45.05% | 38.74% | 37.84% | 40.54% | 25.32% | 40.09% | 5.11 |
| `hold_6_inside` | 3 | 170 | 30.00% | 20.59% | 28.24% | 19.41% | 78.54% | 21.46% | 15.88% | 1.84 |
| `hold_6_inside` | 6 | 170 | 38.82% | 30.59% | 35.29% | 28.24% | 59.41% | 21.46% | 27.06% | 2.85 |
| `hold_6_inside` | 12 | 170 | 52.35% | 46.47% | 43.53% | 40.00% | 44.71% | 21.46% | 35.88% | 5.01 |

The all-bucket result is directionally useful as research context:

- Midline touch and close-across rates rise monotonically with horizon.
- `hold_6_inside` is cleaner than `hold_3_inside` on quick invalidation and
  most ordering labels.
- The `12` bar rows are the first broad hold-inside outputs where every split
  has a midline touch rate above `50%`.

The same rows also explain why this is not entry-ready:

- Average first close-across delay reaches roughly `5` bars in the weakest
  split by the `12` bar horizon.
- Broad `12` bar persistence is only `40.54%` to `44.71%` in the weakest split.
- Broad `12` bar trend leakage reaches `35.88%` to `40.09%` in the worst split.

## Split Detail

Primary `12` bar all-bucket split rows:

| Rule | Split | Cand | Touch Mid | Close Across | Touch Before Boundary | Cross Before Break | Persist | Quick Inv | Trend | Avg Cross Delay |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | `2021_2022_stress` | 363 | 61.43% | 48.76% | 45.45% | 42.42% | 44.08% | 22.04% | 32.51% | 4.74 |
| `hold_3_inside` | `2023_2024_oos` | 312 | 57.69% | 48.40% | 43.91% | 41.67% | 43.59% | 25.32% | 34.94% | 4.73 |
| `hold_3_inside` | `2025_2026_recent` | 222 | 52.25% | 45.05% | 38.74% | 37.84% | 40.54% | 23.42% | 40.09% | 5.11 |
| `hold_6_inside` | `2021_2022_stress` | 283 | 65.72% | 53.71% | 52.30% | 47.35% | 49.12% | 17.31% | 32.16% | 4.76 |
| `hold_6_inside` | `2023_2024_oos` | 233 | 58.80% | 52.79% | 51.50% | 48.07% | 50.21% | 21.46% | 30.47% | 4.20 |
| `hold_6_inside` | `2025_2026_recent` | 170 | 52.35% | 46.47% | 43.53% | 40.00% | 44.71% | 21.18% | 35.88% | 5.01 |

The recent split is usually the floor, but it does not collapse. That makes the
midline event more credible as an observation target than the earlier high/low
paper-side direction labels.

## Bucket Dependence

Narrower rows do not justify promotion:

| Group | Rows |
| --- | ---: |
| Non-`all` rows for `hold_3_inside`/`hold_6_inside` with `>=100` candidates in every split | 4 |
| Non-`all` rows for `hold_3_inside`/`hold_6_inside` with `>=50` candidates in every split | 44 |

The only non-`all` rows at the `>=100` threshold are `hold_3_inside`,
`below_mid`, `all` bucket rows. They do not improve the broad result enough to
change the verdict; at `12` bars they have `38.89%` minimum cross-before-break,
`35.71%` minimum persistence, and `44.44%` maximum trend leakage.

The strongest narrower shape is the mid-position bucket, but it is diagnostic:

| Row | Min Cand | Min Touch Mid | Min Close Across | Min Cross Before Break | Min Persist | Max Quick Inv | Max Trend |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside`, h12, `all`/`mid_50` | 94 | 65.96% | 59.57% | 53.93% | 47.75% | 16.29% | 33.10% |
| `hold_6_inside`, h12, `all`/`mid_50` | 98 | 68.37% | 61.22% | 55.10% | 56.89% | 10.46% | 29.59% |

Those rows are promising, but both fall below `100` candidates in the weakest
split. `hold_3_inside_mid_50` is therefore still diagnostic rather than a
promoted context rule.

## Project State

Current project state remains audit-only:

- `cmd/rangelab` still uses `lab.EmptyStrategy`.
- The latest hold-inside midline transition audit smoke reported
  `strategy=empty trades=0`.
- No Go API, CLI flag, result schema, strategy code, entries, exits, scoring,
  sizing, or strategy replacement changed for this review.

## Conclusion

The midline is a better next research object than frozen high/low direction.
`hold_3_inside` and `hold_6_inside` create a stable enough decision-candle
context to ask what happens after the first midline touch or close-across.

Do not add entries from the current labels. The next step should stay
non-trading and re-index the first closed-candle midline event as the decision
candle, then label post-midline acceptance, rejection, persistence,
invalidation, and trend behavior before considering any entry trigger.
