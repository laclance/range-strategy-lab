# Next Codex Brief: Review Hold-Inside Directional Edge Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md, docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md, docs/RANGE_REGIME_DURABILITY_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the balanced detector regimes are not durable enough as context for future entry hypotheses.
- Detector durability sweep review says no current DefaultDetectorSweepProfiles profile is approved as future entry context; p30_c12_bollinger_on_adx_on is diagnostic only.
- Detector context refinement review says delayed hold_3_inside and hold_6_inside are the first context refinement that materially and split-stably reduces quick invalidation and trend leakage with adequate candidates. They are leading decision-candle context, but NOT promoted to entry context.
- A hold-inside directional edge audit has now been built. It is output-only and has no verdict yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest hold-inside directional edge audit build:
- CLI flag: -hold-inside-directional-edge-audit
- Result directory:
  - results/hold-inside-directional-edge-audit/
- Outputs:
  - hold_inside_directional_edge_candidates.csv/json
  - hold_inside_directional_edge_summary.csv/json
  - hold_inside_directional_edge_stability.csv/json
- Audit size:
  - profiles: 1
  - context rules: 3
  - paper sides: 2
  - candidate rows: 15,976
  - summary rows: 624
  - stability rows: 168
  - candidate CSV lines including header: 15,977
  - summary CSV lines including header: 625
  - stability CSV lines including header: 169
  - horizons: 1, 3, 6, 12
  - quick invalidation window: 3 bars after the decision candle
- Run printed:
  - loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z
  - strategy=empty trades=0

Audit semantics:
- Detector profile is only p30_c12_bollinger_on_adx_off.
- Context rules are hold_3_inside, hold_6_inside, and hold_3_inside_mid_50.
- Candidate rows are one row per passed source episode, context rule, horizon, and paper side:
  - paper_side=toward_high
  - paper_side=toward_low
- Decision fields use only data known at the decision candle:
  - frozen episode high/low/mid
  - decision close position and bucket
  - distance to high/low/mid
  - raw/active length, width, ATR, and width/ATR context
- Forward label_* fields start at decision_index + 1.
- Summary rows aggregate by profile, context rule, split, horizon, paper side, and decision close position bucket, including all.
- Stability rows compare only 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent.

Recommended next task:
Stay non-trading. Review the hold-inside directional edge outputs for split-stable directional edge before any entry trigger. Focus on:
- hold_3_inside and hold_6_inside first; use hold_3_inside_mid_50 as a stricter comparison, not the primary sample.
- Compare paper_side=toward_high vs paper_side=toward_low by horizon and decision close position bucket.
- Prioritize stability rows over full-sample averages.
- Look for worst-split evidence in label_avg_favorable_minus_adverse_pct, label_favorable_greater_than_adverse_rate, label_touched_mid_rate, label_closed_across_mid_rate, label_side_boundary_touch_rate, label_opposite_close_break_rate, and label_quick_invalidated_rate.
- Treat broad positivity as insufficient unless it survives the three period splits and has adequate candidate counts.
- If there is a clear verdict, add a durable review doc such as docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md and update README.md docs order.
- If there is no clear split-stable edge, record a no-promotion conclusion and pivot the next brief to a materially different non-trading hypothesis.

Do not add entries, exits, scoring, sizing, or strategy replacement in the review session. Do not convert paper-side labels into trades. Do not add live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.

Suggested verification for the review:
- wc -l results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check

Closeout:
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next session prompt.
- Commit completed repo changes after closeout unless explicitly told not to.
```
