# Next Codex Brief: Non-Trading Confirmation Audit After Rejection Review

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/ENTRY_READINESS_REVIEW.md, docs/SR_REJECTION_TIMING_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Current verdict:
- Boundary-rejection timing audit is not entry-ready.
- Keep lab.EmptyStrategy.
- Trades remain 0.
- Do not add entries, exits, scoring, sizing, strategy replacement, live code, deploy scripts, API keys, grid, martingale, averaging down, or two-exchange execution.

Current evidence:
- -sr-rejection-timing-audit wrote:
  - results/sr-rejection-timing-audit/sr_rejection_timing_candidates.csv/json
  - results/sr-rejection-timing-audit/sr_rejection_timing_summary.csv/json
  - candidate_rows=968
  - summary_rows=24
- Review note:
  - docs/SR_REJECTION_TIMING_REVIEW.md
- Review conclusion:
  - closed-candle rejection features lift some rejection rates, but favorable-minus-adverse is too small, side-specific, and split-fragile for a first entry.

Goal:
Add or design one more compact non-trading audit that tests a materially different confirmation hypothesis before entries.

Preferred next hypothesis:
- delayed confirmation after an SR rejection candle, re-indexed so the confirmation candle is the decision candle and all label_* outcomes start after that confirmation candle.

Acceptable alternative:
- false-break reclaim timing audit, if you first document why it is now more promising than delayed confirmation.

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
- Keep strategy empty and every smoke run at trades=0.
- Separate decision features from forward label_* outcome fields.
- Update memory/PROGRESS.md with commands, result paths, row counts, and concise factual outcome.
- Update memory/DECISIONS.md only if a durable constraint changes.

Suggested verification:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit
```
```
