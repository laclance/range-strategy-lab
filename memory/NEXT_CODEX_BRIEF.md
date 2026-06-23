# Next Codex Brief: Futures Range-Universe No Automatic Implementation Stop

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this stop state:
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The frozen ETH/SOL 4h structured-compression stream
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00 passed fixed replay but
  failed walk-forward robustness.
- Structured-compression walk-forward stop state:
  structured_compression_walk_forward_fragile_needs_review.
- The post-compression pivot authorized only one automatic materially
  different implementation: a bounded offline fixed-rule
  breakout_retest_acceptance baseline selected from existing range-universe
  discovery evidence.
- That breakout-retest/acceptance baseline has now been implemented, run, and
  reviewed.
- Breakout-retest/acceptance baseline stop state:
  breakout_retest_acceptance_baseline_failed_no_promotion.
- Selected candidates and outcomes:
  - breakout_retest_acceptance_15m_all_h12, discovery rank 22, full aggregate
    5,825 trades, net P&L -2329.18, PF 0.6778;
  - breakout_retest_acceptance_1h_all_h12, discovery rank 28, full aggregate
    1,357 trades, net P&L -604.03, PF 0.8652;
  - both candidates lost in 2023_2024_oos and 2025_2026_recent and had no
    positive full-period symbol transfer.
- Source and resample validation passed for BTCUSDT, ETHUSDT, and SOLUSDT
  Binance USDT-M futures data, so the failure is strategy-premise evidence,
  not a source-gap stop.

Goal:
- Review-only stop unless the user explicitly supplies a materially different
  offline range-strategy premise.
- Do not implement a new strategy, optimizer, replay, walk-forward, grid,
  source expansion, or result-producing run from the existing closed families
  without user input.
- If the user asks what is available, inventory the current docs and present a
  short choice set of materially different offline premises, clearly marking
  closed families as exclusion evidence.

Boundaries:
- Do not retune structured compression.
- Do not retune breakout_retest_acceptance target, stop, max hold, timeframe,
  side, symbol set, selection rules, or review gates around the failed result.
- Do not reopen the failed 1h structured-compression surface.
- Do not rerun failed clean-breakout, hold-inside/midline, impulse absorption,
  boundary touch rejection, single-candle wick rejection, failed breakout
  re-entry, or mature balance persistence as entries unless the user supplies
  a materially new data or structure premise.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Expected review-only shape if the user asks for next options:
- Start from current docs only.
- Treat structured compression walk-forward and breakout-retest baseline as
  exclusion evidence.
- Separate reusable infrastructure from rejected strategy premises.
- End with either:
  - a user-approved bounded offline implementation brief for a materially
    different premise, or
  - a no-implementation stop that says no automatic strategy work is
    authorized.

Verification for any documentation-only update:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- If files are changed, update memory/PROGRESS.md with exact commands and
  factual outcomes.
- Update memory/DECISIONS.md only for a durable user-approved rule.
- Commit completed docs/memory updates and verification evidence after checks
  pass unless explicitly told not to commit.
```
