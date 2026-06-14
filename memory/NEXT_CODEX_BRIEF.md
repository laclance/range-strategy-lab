# Next Codex Brief: Review Compression Breakout Audit Before Any Entries

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
- Compression breakout audit outputs now exist, but they have not yet received a durable entry-readiness review.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Latest compression breakout audit:
- CLI flag:
  - -compression-breakout-audit
- Outputs:
  - results/compression-breakout-audit/compression_breakout_candidates.csv/json
  - results/compression-breakout-audit/compression_breakout_summary.csv/json
- Audit size:
  - candidate_rows=5096
  - summary_rows=24
  - candidate CSV lines including header: 5,097
  - summary CSV lines including header: 25
  - breakout decisions across all splits: 2,548 per horizon
  - one-bar horizon side counts: 1,290 up breakouts and 1,258 down breakouts
- Defaults:
  - detector_profile_id=p30_c12_bollinger_on_adx_off
  - max_breakout_delay=12
  - horizons=1;3;6;12
- Decision semantics:
  - episodes are contiguous RawActive detector runs that eventually become Active
  - episode high/low are frozen using only closed candles through the episode end
  - the first close above the frozen episode high or below the frozen episode low within 12 bars is the decision breakout candle
  - all forward outcome metrics are label_* fields and start after the breakout candle
- Latest smoke:
  - loaded 569451 candles from 2021-01-01T00:00:00Z to 2026-06-01T23:59:59Z
  - strategy=empty trades=0

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
Review the compression breakout audit outputs for split and cohort stability before any entries. Add a durable review document only if the review reaches a clear verdict. Keep lab.EmptyStrategy unless the user explicitly changes scope after the review.

Suggested verification for docs/memory-only closeouts:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- git diff --check
```
