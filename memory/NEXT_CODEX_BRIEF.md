# Next Codex Brief: Refine Detector Context

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
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest detector durability sweep review:
- Review doc:
  - docs/DETECTOR_DURABILITY_SWEEP_REVIEW.md
- Reviewed outputs:
  - results/detector-durability-sweep/detector_durability_sweep.csv/json
  - results/detector-durability-sweep/detector_durability_slices.csv/json
  - results/detector-durability-sweep/detector_durability_stability.csv/json
- Audit size:
  - profiles=19
  - broad_rows=304
  - slice_rows=9088
  - stability_rows=76
  - broad CSV lines including header: 305
  - slice CSV lines including header: 9,089
  - stability CSV lines including header: 77
- Compact evidence:
  - balanced baseline p30_c12_bollinger_on_adx_off had only 13.61% minimum 12-bar persistence, 70.43% maximum quick invalidation, and 51.79% maximum trended rate across period splits
  - ADX comparison p30_c12_bollinger_on_adx_on improved 1-bar minimum persistence to 59.28% and maximum quick invalidation to 40.72%, but still had only 16.47% minimum 12-bar persistence and 50.75% maximum trended rate
  - best broad 12-bar persistence floors were 17.57%, 17.43%, and 17.30%, all still blocked by quick invalidation or trend behavior
  - fully specified 12-bar slice counts by minimum episodes in every period split were 203, 80, 45, 13, 0, and 0 at thresholds 10, 25, 50, 100, 250, and 500
  - the best fully specified 12-bar slice at the 100 episodes-per-split threshold had 26.95% minimum persistence and 60.28% maximum quick invalidation

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
Implement a compact, non-trading detector/context refinement audit that tries to reduce quick invalidation and trend leakage before any entry-trigger work. Use the detector durability review as the evidence base: treat ADX gating as a diagnostic clue, not a promoted detector. Keep outputs inspectable under results/, keep all label_* fields as forward outcomes only, and do not add entries, exits, scoring, sizing, or strategy replacement.

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
    -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
    -out-dir results/<new-detector-context-refinement-dir> \
    <new-audit-flag>
- wc -l results/<new-detector-context-refinement-dir>/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
```
