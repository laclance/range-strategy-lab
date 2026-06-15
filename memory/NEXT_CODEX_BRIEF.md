# Next Codex Brief: Refine Detector Context Before Entries

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/RANGE_REGIME_DURABILITY_REVIEW.md, docs/COMPRESSION_BREAKOUT_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Compression breakout audit was not entry-ready.
- Range regime durability review says the current balanced detector regimes are not durable enough to use as context for future entry hypotheses.
- The current detector is split-stable mainly because weakness repeats across splits, not because regime quality is high.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest range regime durability review:
- Review doc:
  - docs/RANGE_REGIME_DURABILITY_REVIEW.md
- Reviewed outputs:
  - results/range-regime-durability-audit/range_regime_durability_episodes.csv/json
  - results/range-regime-durability-audit/range_regime_durability_summary.csv/json
- Audit size:
  - episode_rows=11984
  - unique episodes=2996
  - summary_rows=452
  - episode CSV lines including header: 11,985
  - summary CSV lines including header: 453
- Detector profile:
  - p30_c12_bollinger_on_adx_off
- Compact evidence:
  - full-sample persistence fell from 44.66% at 1 bar to 14.95% at 12 bars
  - full-sample quick invalidation reached 70.06% by the 3 bar horizon and stayed there
  - 12 bar persistence by period split was 15.68%, 15.07%, and 13.61%
  - no fully specified bucket slice had at least 100 episodes in every period split
  - the best 12 bar fully specified slice with at least 25 episodes in every split still had only 23.93% minimum persistence and 58.12% maximum quick invalidation
  - the broad 12 bar width-only gt_50bp slice still had only 16.80% minimum persistence and 68.07% maximum quick invalidation

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
Refine or reframe detector/context first, before testing another entry trigger. The next implementation should stay non-trading and produce inspectable detector/context durability outputs that answer whether the refined context persists, avoids quick invalidation, and remains stable across 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent. Do not add entries, exits, scoring, sizing, or strategy replacement.

Suggested verification for docs/memory-only closeouts:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check
```
