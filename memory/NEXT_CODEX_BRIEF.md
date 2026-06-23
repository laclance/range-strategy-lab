# Next Codex Brief: New Futures Premise After Clean Breakout Failure

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is BTCUSDT futures strategy research.
- Active source is Binance USDT-M futures, not spot.
- Authoritative parent source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - required product: -source-product binance-usdm-futures
  - market type: Binance USDT-M futures BTCUSDT 5m
  - loaded candles: 573,984
  - CSV lines including header: 573,985
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- The futures clean breakout baseline stop state was:
  clean_breakout_baseline_failed_no_promotion.
- Baseline result:
  - result dir: results/futures-clean-breakout-baseline-backtest/
  - 1h resample rows: 47,832
  - 4h resample rows: 11,958
  - signal rows: 4,848
  - executed trades: 1,185
  - clean_breakout_4h_up_h12: 121 trades, gross P&L 2.55, net P&L -69.29, PF 0.8386
  - clean_breakout_1h_all_h12: 1,064 trades, gross P&L 272.63, net P&L -320.25, PF 0.8846
  - aggregate compatibility summary: 1,185 trades, gross P&L 275.18, net P&L -389.54, PF 0.8784
- The user-approved portfolio-stream routing rule was checked and did not trigger: neither clean breakout candidate was near-viable after costs.
- Default runs still use lab.EmptyStrategy unless an explicit offline prototype/backtest/audit flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Do not optimize the failed clean breakout baseline.
- Do not automatically combine the failed clean breakout candidates into a portfolio-style stream.
- Do not automatically broaden the failed clean breakout template to 15m.
- Ask the user for a new futures premise or explicit scope choice before adding another audit or backtest.

Acceptable next directions only after explicit user choice:
- A materially new BTCUSDT futures premise with a closed-candle event and falsification rule.
- An explicitly authorized 15m clean-breakout comparison, acknowledging that the first 4h/1h baseline failed after costs.
- An explicitly authorized portfolio-style stream review, but only with a new premise for why combination should overcome the independent failures.
- A scope change beyond the current BTCUSDT futures range-first lane.

Implementation boundaries:
- Keep this task planning/review-only unless the user supplies a clear new premise or scope change in the same session.
- Do not add entries, exits, scoring, sizing changes, optimization, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, two-exchange logic, source downloads, or sibling repo mutation without a new approved brief.
- Do not use spot outputs as authority for futures promotion.
- Do not retune the failed clean breakout baseline.

Required next interaction:
- Ask for the next premise or scope choice in concrete terms:
  - What market behavior should be tested next?
  - Which timeframe(s) should be included?
  - Why should it differ from the failed clean breakout baseline and prior failed families?
  - Which closed-candle event defines the candidate?
  - What would falsify it before optimization or live-adjacent work?
- If the user supplies a clear premise, create a scoped next brief. Prefer a fast fixed-rule baseline backtest when prior discovery evidence already supports the premise; otherwise scope a compact non-trading discovery audit first.
- If no premise is supplied, stop with:
  clean_breakout_baseline_failed_no_promotion.

Suggested verification after any docs/memory refresh:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check

Closeout:
- Update memory/PROGRESS.md with the chosen premise or lack of one, exact commands, and stop state.
- Update memory/DECISIONS.md only if a durable rule changes.
- Commit completed repo changes after verification unless explicitly told not to.
```
