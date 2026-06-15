# Next Codex Brief: Build Hold-Inside Midline Reaction Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
  - docs/HOLD_INSIDE_DIRECTIONAL_EDGE_REVIEW.md
  - docs/DETECTOR_CONTEXT_REFINEMENT_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the balanced detector regimes are not durable enough as context for future entry hypotheses.
- Detector durability sweep review says no current DefaultDetectorSweepProfiles profile is approved as future entry context; p30_c12_bollinger_on_adx_on is diagnostic only.
- Detector context refinement review says delayed hold_3_inside and hold_6_inside materially improve range-survival context, but are not entry context.
- Hold-inside directional edge review says hold_3_inside/hold_6_inside do not show a split-stable edge toward frozen high or low.
- Hold-inside midline transition review says hold_3_inside/hold_6_inside do show split-stable midline touch and close-across behavior by 12 bars, but current midline labels are not entry context or strategy inputs.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest reviewed midline-transition facts:
- Result directory:
  - results/hold-inside-midline-transition-audit/
- Review doc:
  - docs/HOLD_INSIDE_MIDLINE_TRANSITION_REVIEW.md
- Broad h12 all-bucket facts:
  - hold_3_inside: 222 minimum split candidates, 52.25% minimum mid touch, 45.05% minimum close-across, 37.84% minimum cross-before-boundary-break, 25.32% maximum quick invalidation, 40.09% maximum trend leakage.
  - hold_6_inside: 170 minimum split candidates, 52.35% minimum mid touch, 46.47% minimum close-across, 40.00% minimum cross-before-boundary-break, 21.46% maximum quick invalidation, 35.88% maximum trend leakage.
  - mid-position h12 rows are cleaner but diagnostic because weakest split counts are 94 for hold_3_inside and 98 for hold_6_inside.

Recommended next task:
Stay non-trading. Add a hold-inside midline reaction audit that re-indexes the first closed-candle midline event as the decision candle, then labels what happens after that event.

Implementation outline:
- Add CLI flag:
  - -hold-inside-midline-reaction-audit
- Default output directory for validation:
  - results/hold-inside-midline-reaction-audit/
- Emit compact CSV/JSON outputs:
  - hold_inside_midline_reaction_candidates.csv/json
  - hold_inside_midline_reaction_summary.csv/json
  - hold_inside_midline_reaction_stability.csv/json
- Reuse the balanced detector profile p30_c12_bollinger_on_adx_off.
- Use context rules:
  - hold_3_inside
  - hold_6_inside
  - hold_3_inside_mid_50 as diagnostic only
- For each source episode/context rule, search after the hold decision candle for the first midline event within max_midline_event_delay_bars=12.
- Include two event types as separate candidate families:
  - mid_touch: first candle whose high/low touches the frozen episode mid
  - mid_close_across: first candle whose close crosses to the opposite side of the frozen episode mid
- The midline event candle is the new decision candle. All forward label_* fields must start at event_index + 1.
- Record decision-candle-known fields only: source hold decision time, event delay bars, event type, frozen high/low/mid, event close position and bucket, event mid side, distances to high/low/mid, episode width, ATR, and width/ATR context.
- Add forward labels over horizons 1, 3, 6, and 12 bars:
  - label_reentered_range
  - label_persisted_inside_range
  - label_quick_invalidated
  - label_invalidated_up/down
  - label_trended_up/down
  - label_touched_high/low
  - label_touched_opposite_half
  - label_closed_back_across_mid
  - label_mid_rejection_before_boundary_touch
  - label_boundary_touch_before_mid_rejection
- Keep labels outcome-only; do not add paper_side, favorable/adverse, entries, exits, scoring, sizing, or strategy replacement.
- Summary rows should aggregate by profile, context rule, event type, split, horizon, event_mid_side, and event close-position bucket, including all.
- Stability rows should compare only 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent.

Expected deliverables:
- Add the audit implementation and focused tests for no-lookahead event selection, touch/cross event separation, label-window start, missing-event skipping, missing-future skipping, deterministic sorting, summary denominators, and stability aggregation.
- Run the audit on BTCUSDT 5m and record result paths and row counts.
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.
- Refresh memory/NEXT_CODEX_BRIEF.md with the next review task for the generated outputs.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build GOPATH=/tmp/range-strategy-lab-go GOMODCACHE=/tmp/range-strategy-lab-go/pkg/mod /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
    -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
    -out-dir results/hold-inside-midline-reaction-audit \
    -hold-inside-midline-reaction-audit
- wc -l results/hold-inside-midline-reaction-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Commit completed repo changes after closeout unless explicitly told not to.
```
