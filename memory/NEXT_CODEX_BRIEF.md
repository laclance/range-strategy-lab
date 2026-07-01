# Next Codex Brief: Session Opening-Range Expansion Candidate Packet

```text
Current state:
- The current range-reversion / midpoint / edge-fade / previous-day range /
  bounded range-optimization path remains closed.
- The approved fixed trend-pullback baseline also failed and is closed:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_IMPLEMENTATION_REVIEW.md.
  btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy.
- Trend-pullback result: source/resample passed; 13,190 signal rows; 3,816
  executed trades; full gross P&L -123.427844; full net P&L -956.566589;
  full PF 0.837958; full max drawdown 0.958438.
- The docs-only next-lane selection selected session-based opening-range
  expansion as the next materially different lane:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_LANE_SELECTION.md.
- Stop state:
  session_opening_range_expansion_lane_selected_for_candidate_packet.

Closed boundaries:
- Do not rescue failed trend-pullback with alternate EMA lengths, slope
  lookbacks, pullback windows, EMA-band definitions, continuation triggers, stop
  buffers, target R values, time stops, side selection, session filters, volume
  filters, volatility filters, derivatives-veto interaction, source expansion,
  replay, walk-forward, or optimizer grids.
- Do not rescue closed range-reversion, midpoint, edge-fade, previous-day range,
  value-area, range-optimization, post-compression, clean-breakout, or
  router/rotation branches by renaming them as opening-range expansion.
- Session-based opening-range expansion must use a predeclared UTC session time
  box and continuation away from that box. It must not fade range edges, target
  midpoints, use EMA trend-pullback rules, or mine session anchors after seeing
  results.
- Do not add paper/testnet/live flow, exchange APIs, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.

Next bounded gate:
- Create the docs-only backtest-first candidate packet for the selected
  session-based opening-range expansion lane:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md.
- Define exactly one fixed BTCUSDT Binance USDT-M futures baseline, likely named
  btc_15m_session_opening_range_expansion_v1.
- The packet must lock one accepted source, one decision timeframe, one fixed
  UTC session anchor, one fixed opening-range length, one closed-candle
  expansion/acceptance trigger, next-bar-open entry, one stop, one target, one
  time stop, fixed sizing, fee, slippage, output path, artifacts, pass/fail
  gates, and no-rescue boundaries.
- Stop after the packet is ready for operator approval to implement/backtest.

Required reads:
- Read AGENTS.md, README.md, memory/README.md, memory/PROGRESS.md,
  memory/DECISIONS.md, memory/NEXT_CODEX_BRIEF.md,
  docs/BACKTEST_FIRST_RESEARCH_LANE.md, docs/STRATEGY_WORKFLOW.md, and
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_LANE_SELECTION.md.

Hard exclusions:
- Do not add Go code, CLI flags, generated results, backtests, optimizers,
  source expansion, derivatives-veto interaction, exchange APIs, credentials,
  deployment, paper/testnet/live flow, martingale, averaging down,
  two-exchange logic, or promotion.
- Update memory/PROGRESS.md after the docs-only milestone, update
  memory/DECISIONS.md only if a durable decision is created, and update
  memory/NEXT_CODEX_BRIEF.md to point to the next bounded approval gate.
```
