# Next Codex Brief: Futures Range-Strategy No Automatic Implementation Stop

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to the user's request. For the current stop
  state, the most relevant docs are:
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research.
- Default CLI behavior remains BTCUSDT futures 5m with `lab.EmptyStrategy`
  unless an explicit offline audit/backtest flag is passed.
- Structured compression passed one fixed ETH/SOL replay, then failed
  walk-forward robustness.
- Structured-compression walk-forward stop state:
  structured_compression_walk_forward_fragile_needs_review.
- Breakout-retest/acceptance was the one automatic materially different
  follow-up from the post-compression pivot review, but its fixed-rule
  baseline failed after costs.
- Breakout-retest/acceptance stop state:
  breakout_retest_acceptance_baseline_failed_no_promotion.
- The user then approved a higher-timeframe nested range-rotation premise:
  closed UTC 1h child ranges inside frozen mature closed UTC 4h parent ranges,
  measuring internal rotation toward parent midpoint/far quartile before
  invalidation.
- That premise has now been implemented as a non-trading audit, run, and
  reviewed.
- Nested range-rotation audit stop state:
  higher_tf_nested_range_rotation_audit_failed_no_baseline.
- Source and resample validation passed for the accepted BTCUSDT Binance
  USDT-M futures data:
  - source rows: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - closed UTC 4h rows: 11,958 through 2026-06-16T20:00:00Z;
  - closed UTC 1h rows: 47,832 through 2026-06-16T23:00:00Z.
- The audit found 68 parent ranges, 282 child ranges, 11 eligible children,
  and only 3 valid full-sample events, all upside. No downside event passed.
- The failure is strategy-premise/count-gate evidence, not a source-gap stop.

Goal:
- Review-only stop unless the user explicitly supplies a materially different
  offline range-strategy premise.
- Do not implement a new strategy, optimizer, replay, walk-forward, grid,
  source expansion, symbol expansion, or result-producing run from the closed
  or fragile families without user input.
- If the user asks what is available, inventory the current docs and present a
  short choice set of materially different offline premises, clearly marking
  closed families as exclusion evidence.

Boundaries:
- Do not retune structured compression.
- Do not retune breakout_retest_acceptance target, stop, max hold, timeframe,
  side, symbol set, selection rules, or review gates around the failed result.
- Do not retune the nested range-rotation 40% child-width gate, 24 bar outcome
  horizon, 6 bar quick-invalidation horizon, parent/child timeframes, detector
  profile, or split gates around the failed result.
- Do not reopen the failed 1h structured-compression surface.
- Do not rerun failed clean-breakout, hold-inside/midline, impulse absorption,
  boundary touch rejection, single-candle wick rejection, failed breakout
  re-entry, mature balance persistence, breakout_retest_acceptance, or
  nested range-rotation as entries unless the user supplies a materially new
  data or structure premise.
- Do not add live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, data downloads, broad symbol mining, martingale, averaging down,
  or two-exchange logic.

Expected review-only shape if the user asks for next options:
- Start from current docs only.
- Treat structured compression walk-forward, breakout-retest/acceptance
  baseline, and nested range-rotation audit as exclusion evidence.
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
