# Next Codex Brief: New Non-Trading Hypothesis After SR Confirmation Review

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
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution unless the user explicitly changes scope.

Current evidence:
- -sr-confirmation-timing-audit wrote:
  - results/sr-confirmation-timing-audit/sr_confirmation_timing_candidates.csv/json
  - results/sr-confirmation-timing-audit/sr_confirmation_timing_summary.csv/json
  - candidate_rows=9692
  - summary_rows=72
  - delays=1;2;3
  - horizons=1;3;6;12
  - strategy=empty trades=0
- Review note:
  - docs/SR_CONFIRMATION_TIMING_REVIEW.md
- Review conclusion:
  - broad decision-confirmation cohorts were only about 48.9%-49.6% of seed rejection candidates
  - aggregate favorable-minus-adverse topped near +0.96bp
  - FGTA was mostly below 50%
  - with minimum 500 rows in every split, best minimum split diff fell to about +0.16bp
  - with minimum 750 rows in every split, no stable cohorts remained
  - no entries or narrower broad-confirmation follow-up should be added from this audit

Goal:
Choose or design one materially different compact non-trading audit before entries. Do not continue broad SR rejection-confirmation slicing unless the hypothesis changes materially.

Preferred next directions:
- false-break reclaim timing audit, but only if it uses a closed-candle reclaim as the decision candle and keeps all future outcomes as label_* fields
- midpoint reclaim after failed continuation, if defined as a simple closed-candle audit over existing SR/range context
- compression breakout inspection, if it stays detector/audit-only and does not add trade signals

Non-negotiables:
- Offline BTCUSDT 5m research only.
- Confirmed closed-candle decisions only.
- Enter on next bar open only when entries eventually exist.
- One open position max.
- Stop-first ambiguity.
- Helper modules may provide feature extraction/audit outputs only.
- Strategy hypotheses, entries, exits, scoring, sizing, and backtest behavior stay inside this lab.
- Generated CSV/JSON outputs stay under results/.

Acceptance criteria:
- Preserve existing audit modes and output schemas.
- Keep strategy empty and every smoke/review run at trades=0 unless the user explicitly asks for entries.
- Treat label_* columns as forward outcomes, not decision inputs.
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.

Suggested verification:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
```
```
