# Next Codex Brief: Review Detector Context Refinement Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md, docs/RANGE_REGIME_DURABILITY_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the current balanced detector regimes are not durable enough to use as context for future entry hypotheses.
- Detector durability sweep review says no current DefaultDetectorSweepProfiles profile is approved as future entry context.
- p30_c12_bollinger_on_adx_on improves short-horizon durability and quick invalidation, but it is diagnostic only and not promoted.
- Detector context refinement audit has been implemented but not reviewed for profile/rule stability yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest detector context refinement audit:
- CLI flag:
  - -detector-context-refinement-audit
- Outputs:
  - results/detector-context-refinement-audit/detector_context_refinement_candidates.csv/json
  - results/detector-context-refinement-audit/detector_context_refinement_summary.csv/json
  - results/detector-context-refinement-audit/detector_context_refinement_stability.csv/json
- Audit size:
  - profiles=8
  - rules=5
  - candidate_rows=113824
  - summary_rows=640
  - stability_rows=160
  - candidate CSV lines including header: 113,825
  - summary CSV lines including header: 641
  - stability CSV lines including header: 161
- Profiles:
  - p30_c12_bollinger_on_adx_off
  - p30_c12_bollinger_on_adx_on
  - p20_c24_bollinger_on_adx_off
  - p30_c24_bollinger_on_adx_off
  - p40_c24_bollinger_on_adx_off
  - p20_c24_bollinger_on_adx_on
  - p30_c24_bollinger_on_adx_on
  - p40_c24_bollinger_on_adx_on
- Context rules:
  - episode_end
  - hold_1_inside
  - hold_3_inside
  - hold_6_inside
  - hold_3_inside_mid_50
- Semantics:
  - episode high/low are frozen at the original raw-active episode end
  - delayed rules set decision_index = episode_end_index + hold_bars
  - source episode counts remain in summary denominators even when a context rule rejects the episode
  - all label_* fields start at decision_index + 1 and are forward outcomes only, not decision inputs
  - stability rows compare 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent
- Latest smoke:
  - loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z
  - strategy=empty trades=0

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- No entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.
- Keep generated CSV/JSON outputs under results/.
- Treat label_* columns as forward outcomes, not decision inputs.
- Keep helper modules behind adapters and limited to feature extraction or audit outputs.
- Update memory/PROGRESS.md with commands, result paths, and concise factual outcome after a completed milestone.
- Update memory/DECISIONS.md only if a durable constraint changes.
- After completing a brief or milestone, run closeout checks and commit the completed repo changes unless the user explicitly says not to commit.

Recommended next task:
Review the detector context refinement audit outputs for profile/rule-level split stability and regime quality before any detector promotion or entry-trigger work. Decide whether any context rule materially reduces quick invalidation and trend leakage with enough split-stable candidates, or whether detector/context refinement must continue. Do not add entries, exits, scoring, sizing, or strategy replacement in that review. Do not add a durable verdict doc unless the outputs are actually reviewed in that same session.

Suggested verification for docs/memory-only closeouts:
- wc -l results/detector-context-refinement-audit/detector_context_refinement_candidates.csv results/detector-context-refinement-audit/detector_context_refinement_summary.csv results/detector-context-refinement-audit/detector_context_refinement_stability.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check
```
