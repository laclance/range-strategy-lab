# Compression Breakout Review

Date: 2026-06-14

## Verdict

Compression breakout audit is not entry-ready; keep `lab.EmptyStrategy`.

The audit has clean closed-candle decision semantics: the frozen compression
episode high/low are known through the episode end, the first close outside
that frozen range is the breakout decision candle, and all `label_*` outcomes
start after that breakout candle. The broad side/horizon averages are positive
by favorable-minus-adverse, but the evidence is not stable enough for first
trade logic. Recent up breakouts turn negative from `3` through `12` bars, down
breakouts are negative in the stress split at `12` bars and in the recent split
at `1` bar, favorable-greater-than-adverse rates are mostly below `50%`, and
candidate cohorts collapse under per-split row thresholds.

No entries, exits, scoring, sizing, or strategy replacement should be added
from this audit.

## Inputs Reviewed

Generated audit paths:

- `results/compression-breakout-audit/compression_breakout_summary.csv`
- `results/compression-breakout-audit/compression_breakout_summary.json`
- `results/compression-breakout-audit/compression_breakout_candidates.csv`
- `results/compression-breakout-audit/compression_breakout_candidates.json`

Audit size:

- candidate rows: `5,096`
- summary rows: `24`
- candidate CSV lines including header: `5,097`
- summary CSV lines including header: `25`
- breakout decisions across all splits: `2,548` per horizon
- one-bar side counts: `1,290` up breakouts and `1,258` down breakouts
- detector profile: `p30_c12_bollinger_on_adx_off`
- max breakout delay: `12` bars after the raw-active compression episode ends
- forward horizons: `1`, `3`, `6`, `12` bars after the breakout candle

The `label_*` columns are forward outcomes only. They are not decision inputs.

Audit semantics:

- Episodes are contiguous `RawActive` detector runs that eventually become
  `Active`.
- Episode high/low are frozen using only closed candles through the episode
  end.
- The first close above the frozen episode high or below the frozen episode low
  within `12` bars is the breakout decision candle.
- All forward outcome metrics start after the breakout candle.

## Side And Horizon Summary

Rows below aggregate the three period splits for each side and horizon,
weighted by `candidate_count`. `diff` is label favorable-minus-adverse in
basis points.

| Side | Horizon | Candidates | Diff | FGTA | Reentered Range | Opposite Close Break |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| up | 1 | 1,290 | +2.23bp | 46.59% | 22.09% | 0.31% |
| up | 3 | 1,290 | +2.09bp | 47.44% | 43.26% | 1.32% |
| up | 6 | 1,290 | +2.25bp | 46.12% | 56.98% | 2.64% |
| up | 12 | 1,290 | +2.73bp | 47.21% | 65.97% | 7.36% |
| down | 1 | 1,258 | +2.10bp | 44.28% | 26.63% | 0.00% |
| down | 3 | 1,258 | +2.78bp | 47.30% | 46.34% | 0.95% |
| down | 6 | 1,258 | +2.97bp | 46.66% | 59.46% | 2.78% |
| down | 12 | 1,258 | +2.31bp | 46.50% | 70.11% | 6.76% |

The broad grid is positive by average excursion difference, but FGTA stays
mostly below `50%` and range re-entry rises sharply by the `12` bar horizon.
That is not enough for entry readiness.

## Split Stability

Cells are `diff / FGTA / n`.

| Split | Side | H1 | H3 | H6 | H12 |
| --- | --- | ---: | ---: | ---: | ---: |
| 2021_2022_stress | up | +5.02bp / 48.50% / n=501 | +4.30bp / 47.31% / n=501 | +4.21bp / 45.31% / n=501 | +5.40bp / 47.50% / n=501 |
| 2021_2022_stress | down | +4.13bp / 45.88% / n=510 | +4.22bp / 49.02% / n=510 | +2.07bp / 46.86% / n=510 | -1.87bp / 44.51% / n=510 |
| 2023_2024_oos | up | +0.75bp / 46.93% / n=473 | +1.28bp / 47.36% / n=473 | +2.70bp / 47.15% / n=473 | +3.24bp / 48.41% / n=473 |
| 2023_2024_oos | down | +1.47bp / 44.21% / n=423 | +2.73bp / 46.57% / n=423 | +5.21bp / 47.52% / n=423 | +5.14bp / 48.46% / n=423 |
| 2025_2026_recent | up | +0.04bp / 43.04% / n=316 | -0.21bp / 47.78% / n=316 | -1.53bp / 45.89% / n=316 | -2.27bp / 44.94% / n=316 |
| 2025_2026_recent | down | -0.27bp / 41.85% / n=325 | +0.57bp / 45.54% / n=325 | +1.48bp / 45.23% / n=325 | +5.20bp / 47.08% / n=325 |

The split mismatch blocks first-entry work. The recent period does not support
up breakouts beyond one bar, while down breakouts fail in different cells.

## Candidate Cohorts

Thresholds below require every split to have at least the listed number of
candidates in the same cohort and a positive minimum split `diff`.

| Cohort Grouping | >=25 | >=50 | >=100 | >=250 | >=500 |
| --- | ---: | ---: | ---: | ---: | ---: |
| full feature: side, horizon, delay, raw length, active length, width, breakout move, true-range expansion | 0 | 0 | 0 | 0 | 0 |
| no delay: side, horizon, raw length, active length, width, breakout move, true-range expansion | 0 | 0 | 0 | 0 | 0 |
| breakout shape: side, horizon, delay, width, breakout move, true-range expansion | 0 | 0 | 0 | 0 | 0 |
| coarse episode: side, horizon, raw length, active length, width | 7 | 0 | 0 | 0 | 0 |
| coarse breakout: side, horizon, width, breakout move, true-range expansion | 0 | 0 | 0 | 0 | 0 |
| width and move: side, horizon, width, breakout move | 5 | 0 | 0 | 0 | 0 |

Best stable examples were still too sparse or too weak:

| Screen | Min Split Rows | Best Cohort | Total | Min Count | Min Diff | Min FGTA | Note |
| --- | ---: | --- | ---: | ---: | ---: | ---: | --- |
| coarse episode | 25 | down h3, raw 48plus, active 48plus, width gt_50bp | 105 | 26 | +6.75bp | 42.00% | positive diff, but tiny and poor FGTA |
| width and move | 25 | down h6, width gt_50bp, breakout move 10_20bp | 151 | 25 | +3.94bp | 44.00% | sparse and below coin-flip by FGTA |

No positive stable cohort survived the `50` candidates-per-split threshold.

## Conclusion

This audit should be treated as a useful closed-candle compression breakout
inspection, not as an entry signal. The breakout decision framing is clean, but
the observed edge is too split-sensitive and cohort-fragile for first trade
logic.

Do not continue narrower compression-breakout slicing unless the next
hypothesis changes materially. Next work should stay non-trading and move to a
different idea before any entries are added.
