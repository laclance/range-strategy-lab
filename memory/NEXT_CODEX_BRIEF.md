# Next Codex Brief: Review Detector Durability Sweep

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
- Detector durability sweep has been implemented but not reviewed for profile stability yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest detector durability sweep:
- CLI flag:
  - -detector-durability-sweep
- Outputs:
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
- Defaults:
  - profile grid: existing DefaultDetectorSweepProfiles
  - horizons=1;3;6;12
  - quick_invalidation_bars=3
- Semantics:
  - broad rows are one row per detector profile, split, and horizon
  - slice rows use existing raw length, active length, width, and width/ATR buckets
  - stability rows compare 2021_2022_stress, 2023_2024_oos, and 2025_2026_recent
  - label_* fields are forward outcomes only, not decision inputs
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
Review the detector durability sweep outputs for profile-level split stability and regime quality before any detector promotion or entry-trigger work. Decide whether any detector profile is durable enough to become future entry context, or whether detector/context refinement must continue. Do not add entries, exits, scoring, sizing, or strategy replacement in that review. Do not add a durable verdict doc unless the sweep outputs are actually reviewed in that same session.

Suggested verification for docs/memory-only closeouts:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check
```
