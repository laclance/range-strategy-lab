# Next Codex Brief: Choose Next Non-Trading Hypothesis After SR Timing Reviews

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/ENTRY_READINESS_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit was not entry-ready.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest false-break reclaim review:
- Durable report:
  - docs/SR_FALSE_BREAK_RECLAIM_TIMING_REVIEW.md
- Inputs reviewed:
  - results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv/json
  - results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv/json
- Audit size:
  - candidate_rows=17652
  - summary_rows=24
  - candidate CSV lines including header: 17,653
  - summary CSV lines including header: 25
  - support reclaim decisions across all splits: 4,120 per horizon
  - resistance reclaim decisions across all splits: 4,150 per horizon
- Compact verdict evidence:
  - broad side/horizon favorable-minus-adverse was small-positive, topping at about +2.41bp
  - broad FGTA was mostly near coin-flip, about 49.27% to 51.54%
  - 2025_2026_recent support was negative across all horizons
  - 2023_2024_oos resistance was negative at 3, 6, and 12 bars
  - fully sliced cohorts had 0 stable rows at every threshold from 25 to 500 candidates per split
  - coarse cohorts had no stable rows at 500 candidates per split
  - close-shape cohorts survived at 500 candidates per split, but the best row had only about +0.51bp min split diff and 46.69% min FGTA

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- No entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.
- Keep generated CSV/JSON outputs under results/.
- Treat label_* columns as forward outcomes, not decision inputs.
- Update memory/PROGRESS.md with commands, result paths, and concise factual outcome after a completed milestone.
- Update memory/DECISIONS.md only if a durable constraint changes.
- After completing a brief or milestone, run closeout checks and commit the completed repo changes unless the user explicitly says not to commit.

Recommended next task:
Do not continue narrow SR timing slices by default. Choose or design a materially different non-trading audit before any entries. A reasonable default branch is a compact compression-breakout or volatility-expansion audit built from the existing detector outputs, but first keep the work to an inspectable audit plan or audit outputs, not trade logic.

Suggested verification for docs/memory-only closeouts:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```
```
