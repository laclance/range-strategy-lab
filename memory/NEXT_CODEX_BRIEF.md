# Next Codex Brief: Review False-Break Reclaim Timing Audit

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/ENTRY_READINESS_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/SR_CONFIRMATION_TIMING_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit was not entry-ready.
- Delayed confirmation after SR rejection was not entry-ready.
- False-break reclaim timing audit has been implemented, but it has not been reviewed for entry-readiness yet.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Current false-break reclaim evidence:
- -sr-false-break-reclaim-timing-audit wrote:
  - results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_candidates.csv/json
  - results/sr-false-break-reclaim-timing-audit/sr_false_break_reclaim_timing_summary.csv/json
  - candidate_rows=17652
  - summary_rows=24
  - max_break_delay=3
  - max_reclaim_delay=12
  - horizons=1;3;6;12
  - strategy=empty trades=0
- CSV line counts including header:
  - sr_false_break_reclaim_timing_candidates.csv: 17,653
  - sr_false_break_reclaim_timing_summary.csv: 25
- Broad summary read:
  - support reclaim decisions across all splits: 4,120 per horizon
  - resistance reclaim decisions across all splits: 4,150 per horizon
  - broad side/horizon favorable-minus-adverse was small-positive:
    - support h1/h3/h6/h12 about +1.03bp/+1.76bp/+1.61bp/+0.75bp
    - resistance h1/h3/h6/h12 about +0.84bp/+1.28bp/+1.85bp/+2.41bp
  - broad favorable-greater-than-adverse rates were about 49.3% to 51.5%

Audit semantics to preserve:
- The anchor boundary comes from existing SR audit rows.
- Support false break closes below the frozen anchor support zone bottom, then reclaims with a close at or above the frozen anchor support level.
- Resistance false break closes above the frozen anchor resistance zone top, then reclaims with a close at or below the frozen anchor resistance level.
- The reclaim candle is the decision candle.
- All forward outcome metrics are label_* fields and start after the reclaim candle.

Goal:
Review the false-break reclaim timing outputs for split and cohort stability before any entries. Produce a concise verdict grounded in the generated CSVs. If adding a durable review doc, keep it factual and make clear whether the audit is entry-ready or not.

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- No entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.
- Keep generated CSV/JSON outputs under results/.
- Treat label_* columns as forward outcomes, not decision inputs.
- Update memory/PROGRESS.md with commands, result paths, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.

Suggested verification:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```
```
