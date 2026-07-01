# Next Codex Brief: No Selected Next Implementation After Session OR Failure

```text
Current state:
- The current range-reversion / midpoint / edge-fade / previous-day range /
  bounded range-optimization path remains closed.
- The fixed trend-pullback baseline failed and is closed:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_IMPLEMENTATION_REVIEW.md.
  btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy.
- The session-based opening-range expansion lane was selected and packeted:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_LANE_SELECTION.md.
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md.
- The approved fixed session opening-range expansion baseline was implemented
  and backtested:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_IMPLEMENTATION_REVIEW.md.
- Stop state:
  btc_15m_session_opening_range_expansion_backtest_failed_no_usable_strategy.

Latest fixed-baseline result:
- Candidate id: btc_15m_session_opening_range_expansion_v1.
- CLI flag: -backtest-first-btc-15m-session-opening-range-expansion-v1.
- Output path:
  results/backtest-first-btc-15m-session-opening-range-expansion-v1/.
- Source/resample passed on the accepted BTCUSDT Binance USDT-M futures 5m
  source: 573,984 loaded candles; 2021-01-01T00:00:00Z through
  2026-06-16T23:55:00Z; gap_count=0; duplicate_count=0;
  zero_volume_count=66; comparison_only=false; validation_status=accepted;
  exact closed UTC 15m row_count=191,328.
- Result: 1,993 session-range rows; 1,652 signal rows; 1,652 executed trades;
  12 summary rows.
- Full all-side metrics: gross P&L 151.732485; net P&L -546.240518;
  PF 0.875398; max drawdown 0.578097.
- Primary split all-side metrics:
  2021_2022_stress gross 45.770589 / net -237.826052;
  2023_2024_oos gross 122.959395 / net -139.289001;
  2025_2026_recent gross -16.997499 / net -169.125465.
- Failed gates: gross edge, net edge, profit factor, drawdown.

Closed boundaries:
- Do not rescue failed trend-pullback with alternate EMA lengths, slope
  lookbacks, pullback windows, EMA-band definitions, continuation triggers, stop
  buffers, target R values, time stops, side selection, session filters, volume
  filters, volatility filters, derivatives-veto interaction, source expansion,
  replay, walk-forward, or optimizer grids.
- Do not rescue failed session opening-range expansion with alternate UTC
  anchors, opening-range lengths, expansion windows, acceptance buffers, ATR
  windows, stop buffers, target R values, time stops, one-trade-per-day changes,
  side selection, weekday filters, volume filters, volatility filters,
  derivatives-veto interaction, source expansion, replay, walk-forward, or
  optimizer grids.
- Do not rescue closed range-reversion, midpoint, edge-fade, previous-day range,
  value-area, range-optimization, post-compression, clean-breakout,
  router/rotation, trend-pullback, or opening-range branches by renaming them.
- Do not add paper/testnet/live flow, exchange APIs, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.

Next bounded gate:
- No implementation, backtest, optimizer, replay, walk-forward, derivatives-veto
  interaction, source expansion, or promotion is currently selected.
- The next valid action requires the operator to provide either:
  1. a materially different BTCUSDT Binance USDT-M futures offline strategy
     premise for a docs-only candidate packet or short lane selection; or
  2. an explicit docs-only progress/status/handoff task.
- If the operator supplies a materially different premise, use
  docs/BACKTEST_FIRST_RESEARCH_LANE.md and docs/STRATEGY_WORKFLOW.md to decide
  whether a short candidate packet is enough before any code.
- If no materially different premise is supplied, stay read-only or docs/memory
  only and do not invent a new implementation lane.

Required reads for any next strategy-continuation work:
- Read AGENTS.md, README.md, memory/README.md, memory/PROGRESS.md,
  memory/DECISIONS.md, memory/NEXT_CODEX_BRIEF.md,
  docs/BACKTEST_FIRST_RESEARCH_LANE.md, docs/STRATEGY_WORKFLOW.md,
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_IMPLEMENTATION_REVIEW.md,
  and only the additional docs directly relevant to the operator's new premise.

Hard exclusions:
- No rescue retuning of closed candidates or families.
- No generated results/backtests unless a specific materially different fixed
  candidate is explicitly approved.
- No source expansion, derivatives-veto interaction, optimizer, replay,
  walk-forward, exchange API, credentials, deployment, paper/testnet/live flow,
  martingale, averaging down, two-exchange logic, or promotion.
- If a docs-only milestone is completed, update memory/PROGRESS.md, update
  memory/DECISIONS.md only for durable decisions, and update
  memory/NEXT_CODEX_BRIEF.md to point to the next bounded gate.
```
