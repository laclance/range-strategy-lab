# Next Codex Brief: Build Hold-Inside Midline Transition Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md, docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md, docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md, docs/RANGE_REGIME_DURABILITY_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
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
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest hold-inside directional edge review:
- Review doc:
  - docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md
- Inputs reviewed:
  - results/hold-inside-directional-edge-audit/hold_inside_directional_edge_candidates.csv/json
  - results/hold-inside-directional-edge-audit/hold_inside_directional_edge_summary.csv/json
  - results/hold-inside-directional-edge-audit/hold_inside_directional_edge_stability.csv/json
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
- Compact evidence:
  - hold_3_inside h3 toward_high all-bucket had 222 minimum split candidates and +0.78bp minimum favorable-minus-adverse, but only 46.56% minimum favorable-greater-than-adverse.
  - hold_3_inside h6 toward_high all-bucket had 222 minimum split candidates and +0.62bp minimum favorable-minus-adverse, but only 45.51% minimum favorable-greater-than-adverse.
  - hold_6_inside h6 toward_low all-bucket had 170 minimum split candidates and 50.18% minimum favorable-greater-than-adverse, but -0.03bp minimum favorable-minus-adverse.
  - all other primary all-bucket rows had negative minimum favorable-minus-adverse, minimum favorable-greater-than-adverse below 50%, or both.
  - non-all bucket rows for hold_3_inside/hold_6_inside: 0 rows with >=100 candidates in every period split; 32 rows with >=50 candidates in every period split, but the best relaxed positives were sparse or weak by FGTA.

Recommended next task:
Stay non-trading. Build a compact hold-inside midline transition audit that asks a different question from the failed toward-high/toward-low paper-side audit:
- condition only on p30_c12_bollinger_on_adx_off with hold_3_inside and hold_6_inside first
- keep hold_3_inside_mid_50 as diagnostic comparison only if useful
- use only confirmed closed-candle decision fields known at the decision candle
- label what happens after the decision candle around the frozen range midline, not a paper trade side
- useful label examples:
  - touched_mid within horizon
  - closed_across_mid within horizon
  - first_mid_touch_delay_bars
  - mid_touch_before_boundary_touch
  - mid_cross_before_boundary_close_break
  - persisted_inside_range
  - quick_invalidated within 3 bars
  - invalidated_up/down and trended_up/down for context
- aggregate by profile, context rule, split, horizon, decision_mid_side, and decision close position bucket including all
- emit candidate, summary, and stability CSV/JSON under results/hold-inside-midline-transition-audit/
- prioritize stability rows over full-sample averages and require adequate per-split candidates before any follow-up

Do not add entries, exits, scoring, sizing, or strategy replacement in this implementation. Do not convert midline labels into trades. Do not add live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.

Suggested verification for the next audit:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
    -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
    -out-dir results/hold-inside-midline-transition-audit \
    -hold-inside-midline-transition-audit
- wc -l results/hold-inside-midline-transition-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next session prompt.
- Commit completed repo changes after closeout unless explicitly told not to.
```
