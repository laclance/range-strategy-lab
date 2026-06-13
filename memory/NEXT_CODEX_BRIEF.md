# Next Codex Brief: Boundary-Rejection Timing Audit Before Entries

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and docs/*.md, especially docs/ENTRY_READINESS_REVIEW.md, docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Goal:
Add one compact non-trading timing audit that tests whether SR boundary rejection can be identified using only the decision candle and prior data. Do not add trade entries in this milestone.

Current state:
- Latest review gate added docs/ENTRY_READINESS_REVIEW.md.
- Strategy remains lab.EmptyStrategy.
- Trades remain 0.
- Internal lab coverage is 99.8%; remaining uncovered lines are justified defensive SR audit error propagation branches.
- The project is offline BTCUSDT 5m research only.
- Existing audit modes:
  - -detector
  - -detector-sweep
  - -sr-audit
  - -sr-boundary-audit
  - -sr-boundary-inspect
- SR boundary inspection conclusion: boundary rejection has better current evidence than false-break reclaim, but rejection is still an ex-post label and not yet an entry rule.

Non-negotiables:
- Offline research only.
- No live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only.
- 5m candles first.
- Confirmed closed-candle decisions only.
- Enter on next bar open when entries are eventually added.
- One open position max.
- Stop-first ambiguity.
- Keep every result explainable and reproducible.
- Do not reuse strategy, scoring, or live-execution logic from the old binance-bot project.
- Helper modules may provide feature extraction/audit outputs only; strategy hypotheses, entries, exits, scoring, sizing, and backtest behavior stay inside this lab.

Timing audit task:
1. Add a compact non-trading audit mode, preferably named -sr-rejection-timing-audit.
2. Reuse existing SR audit rows and balanced detector context.
3. Build decision-candle features only from the current closed candle and previously computed SR context:
   - support/resistance side
   - touched or pierced zone on the decision candle
   - closed back above support or below resistance
   - wick size beyond boundary
   - close location relative to zone
   - zone strength bucket
   - distance bucket
   - detector_active context
4. Reuse 1, 3, 6, and 12 bar forward outcomes only as labels, not decision inputs.
5. Keep outputs compact under results/, for example:
   - sr_rejection_timing_candidates.csv/json
   - sr_rejection_timing_summary.csv/json
6. Preserve existing audit modes and output schemas.
7. Do not add entries, exits, scoring, sizing changes, or strategy replacement.

Acceptance criteria:
- The new timing audit can distinguish decision-candle candidate features from future outcome labels.
- Outputs are compact enough to inspect quickly.
- Existing detector, SR audit, boundary audit, and boundary inspect modes still work.
- Tests cover no-lookahead boundaries, support/resistance symmetry, bucket behavior, and summary denominators.
- Strategy remains empty and every smoke run prints trades=0.
- memory/PROGRESS.md records commands, result paths, row counts, and concise outcome.

Required verification commands:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-rejection-timing-audit \
  -sr-rejection-timing-audit
```
```
