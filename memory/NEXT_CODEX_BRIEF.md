# Next Codex Brief: Entry Readiness Review Before Trades

```text
We are in /home/lance/range-strategy-lab, a standalone Go project named range-strategy-lab.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md and all docs/*.md, especially docs/RESEARCH_HELPERS.md, docs/STRATEGY_WORKFLOW.md, docs/ARCHITECTURE.md, and docs/VERIFICATION.md.
- Check git status before editing.

Goal:
Perform a full entry-readiness review of the project before adding the first trade entries. Treat this as an important quality gate. Do not add entries in the review pass unless the user explicitly asks after the review is complete.

Current state:
- Latest completed implementation commit before this brief: 1e26695 Add SR boundary inspection mode.
- Strategy remains lab.EmptyStrategy.
- Trades remain 0.
- The project is offline BTCUSDT 5m research only.
- CSV loading, one-position backtest engine, costs, risk sizing, split metrics, detector diagnostics, detector sweep, SR audit, SR boundary-quality audit, and compact SR boundary candidate comparison modes exist.
- Pinned helper modules may only provide feature extraction or audit outputs:
  - github.com/laclance/go-sr v1.0.0
  - github.com/markcheno/go-talib v0.0.0-20250114000313-ec55a20c902f
  - nproject.io/gitlab/libraries/talib-cdl-go v0.0.0-20211217160304-2ed8176448cc

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

Known research facts:
- Balanced detector baseline:
  - profile: p30_c12_bollinger_on_adx_off
  - full split: 77,231 active bars / 569,451 total bars
  - duty cycle: 13.56%
  - episodes: 2,996
- SR audit:
  - flag: -sr-audit
  - defaults: 5m zone mode, 120-bar lookback, min strength 2, warmup 138 bars
  - latest row count: 569,313
  - near-support rows: 126,861
  - near-resistance rows: 128,085
- SR boundary-quality audit:
  - flag: -sr-boundary-audit
  - horizons: 1, 3, 6, 12 bars
  - filter: detector_active=true
  - event rows: 281,080
  - quality rows: 192
- SR boundary inspection:
  - flag: -sr-boundary-inspect
  - output: results/sr-boundary-inspection/sr_boundary_candidate_comparison.csv/json
  - comparison rows: 192
  - conclusion: boundary rejection has better current evidence than false-break reclaim
  - rejection rates by side/horizon across all splits: 10.52% to 36.77%
  - rejected-cohort favorable-minus-adverse: 14.75bp to 27.73bp
  - reclaimed-cohort favorable-minus-adverse was negative at 3, 6, and 12 bars for both support and resistance in this event definition
  - 10_20bp distance bucket remains sparse: 208 resistance events and 96 support events across all splits/horizons

Quality bar before entries:
- Code should be correct, bug free as far as review and tests can establish, clean, readable, and not over-complicated.
- Test coverage should be pushed to 100% for core lab logic before entries, or every remaining uncovered line must be explicitly justified in the review output before the user accepts it.
- Documentation and memory should be current enough that a fresh Codex session can verify the same facts without guessing.
- No audit result should rely on future bars for a field that would later be used as an entry decision.

Current coverage note:
- `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab` currently reports 84.5% statement coverage.
- `env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./...` currently hits the known snap-confine issue for cmd/rangelab even though `/usr/local/go/bin/go test ./cmd/rangelab` passes.
- The review should investigate coverage gaps, add missing focused tests where appropriate, and document any tooling limitation separately from real coverage gaps.

Review task:
1. Review the entire project before entries:
   - CLI behavior in cmd/rangelab
   - CSV loading and timestamp assumptions
   - backtest engine, stop-first ambiguity, fees, slippage, sizing, max-notional cap, and force-close behavior
   - split metrics and summary outputs
   - detector and detector sweep logic
   - SR adapter no-lookahead behavior
   - SR boundary audit and compact inspection summaries
   - tests, docs, memory, and generated-output hygiene
2. Look for bugs, lookahead leaks, incorrect denominators, unstable sort/order behavior, missing tests, unclear names, over-complication, stale docs, and any path that could accidentally add live/exchange behavior.
3. Produce a concise review report with:
   - critical findings first, with file/line references
   - coverage gaps and the exact tests needed
   - documentation/memory gaps
   - whether the project is ready for entry implementation
4. If the review finds small, unambiguous fixes or missing tests, implement them and update memory/PROGRESS.md with commands and outcomes.
5. If the review finds a design ambiguity that affects entries, stop and ask the user before implementing entries.

Next research task after the review is clean:
- Add one more compact non-trading timing audit before entries.
- The audit should ask whether boundary rejection can be identified using only the decision candle and prior data, not future bars.
- Candidate decision-candle fields may include:
  - support/resistance side
  - touched or pierced zone
  - closed back above support or below resistance
  - wick size beyond boundary
  - close location relative to zone
  - zone strength bucket
  - distance bucket
  - detector_active context
- Reuse 1, 3, 6, and 12 bar forward outcomes only as labels, not decision inputs.
- Keep outputs compact under results/.
- Do not add trade entries until this review and timing audit support a closed-candle entry rule.

Required verification commands:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
git diff --check

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -csv ../binance-bot/data/btcusdt_spot_5m_2021_2026.csv \
  -out-dir results/sr-boundary-inspection-check \
  -sr-boundary-inspect
```

Coverage commands to investigate:
```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./internal/lab
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test -cover ./...
```

Acceptance criteria:
- The review report is specific enough to act on, with file/line references for findings.
- Tests pass.
- Coverage gaps are either fixed or explicitly listed with a reason and a proposed test.
- Documentation and memory are updated if any project facts or next steps change.
- Strategy remains empty and trades remain 0.
- No entries, exits, scoring, sizing changes, live code, or exchange logic are added during the review unless the user explicitly asks after reviewing the report.
```
