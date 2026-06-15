# Next Codex Brief: Review Hold-Inside Midline Transition Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this review:
  - docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md
  - docs/RANGE_REGIME_DURABILITY_REVIEW.md
  - docs/RESEARCH_HELPERS.md and docs/VERIFICATION.md only if implementation or verification questions arise
- Do not read every historical review doc by default; open older SR/compression docs only if the current review needs that comparison.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the balanced detector regimes are not durable enough as context for future entry hypotheses.
- Detector durability sweep review says no current DefaultDetectorSweepProfiles profile is approved as future entry context; p30_c12_bollinger_on_adx_on is diagnostic only.
- Detector context refinement review says delayed hold_3_inside and hold_6_inside are the first context refinement that materially and split-stably reduces quick invalidation and trend leakage with adequate candidates. They are leading decision-candle context, but NOT promoted to entry context.
- Hold-inside directional edge review says hold_3_inside/hold_6_inside do not show a split-stable directional edge toward the frozen range high or low. No all-bucket row passed the review gate, and no non-all decision-close-position bucket reached 100 candidates in every period split.
- Hold-inside midline transition audit is built and generated, but not reviewed yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest hold-inside midline transition audit:
- CLI flag:
  - -hold-inside-midline-transition-audit
- Result directory:
  - results/hold-inside-midline-transition-audit/
- Outputs:
  - hold_inside_midline_transition_candidates.csv/json
  - hold_inside_midline_transition_summary.csv/json
  - hold_inside_midline_transition_stability.csv/json
- Audit size:
  - profiles: 1
  - context rules: 3
  - candidate rows: 7,988
  - summary rows: 720
  - stability rows: 192
  - candidate CSV lines including header: 7,989
  - summary CSV lines including header: 721
  - stability CSV lines including header: 193
  - base summary CSV lines including header: 13
  - horizons: 1, 3, 6, 12
  - quick invalidation window: 3 bars after the decision candle
- Scope:
  - profile p30_c12_bollinger_on_adx_off
  - primary context rules hold_3_inside and hold_6_inside
  - hold_3_inside_mid_50 included as diagnostic output only
  - no paper_side labels and no favorable/adverse fields
  - missing first-delay labels use -1
  - same-bar midline and boundary events do not satisfy before-ordering labels
  - stability rows compare only 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent
  - stability rows include min/max/delta rates for reentry, persistence, quick invalidation, invalidation up/down, and trend up/down labels
- Last run printed:
  - hold_inside_midline_transition_audit profiles=1 rules=3 candidate_rows=7988 summary_rows=720 stability_rows=192 quick_invalidation_bars=3 horizons=1;3;6;12
  - loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z
  - strategy=empty trades=0

Recommended next task:
Stay non-trading. Review existing results/hold-inside-midline-transition-audit/ outputs only; do not add Go code, strategy logic, entries, exits, scoring, sizing, or live wiring.

Review order:
- Start with hold_inside_midline_transition_stability.csv.
- Use summary/candidate rows only to explain notable slices.
- Prioritize hold_3_inside and hold_6_inside.
- Treat hold_3_inside_mid_50 as diagnostic unless the review later justifies promoting it as context.
- Focus on split-stable behavior across 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent.

Suggested review questions:
- Do all-bucket hold_3_inside or hold_6_inside rows show stable rates for label_touched_mid and label_closed_across_mid?
- Do label_mid_touch_before_boundary_touch and label_mid_cross_before_boundary_close_break hold up in worst split with adequate candidate_count_min?
- Are persisted-inside, quick-invalidated, invalidated-up/down, and trended-up/down rates consistent with a useful non-trading context?
- Are any positive rows dependent on sparse decision_mid_side or close-position buckets?
- Do first-delay averages suggest a practical closed-candle observation window, or are they too slow/noisy?

Expected review deliverables if the evidence supports a clear verdict:
- Add docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md.
- Update README.md docs order if the review doc is added.
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint or no-promotion rule changes.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next materially different non-trading hypothesis or review follow-up.

Verification:
- wc -l results/hold-inside-midline-transition-audit/hold_inside_midline_transition_candidates.csv results/hold-inside-midline-transition-audit/hold_inside_midline_transition_summary.csv results/hold-inside-midline-transition-audit/hold_inside_midline_transition_stability.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check

Closeout:
- Commit completed repo changes after closeout unless explicitly told not to.
```
