# Next Codex Brief: Implement Session Opening-Range Expansion Baseline

```text
Current state:
- The current range-reversion / midpoint / edge-fade / previous-day range /
  bounded range-optimization path remains closed.
- The approved fixed trend-pullback baseline also failed and is closed:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_IMPLEMENTATION_REVIEW.md.
  btc_15m_trend_pullback_continuation_backtest_failed_no_usable_strategy.
- The docs-only lane selection selected session-based opening-range expansion:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_LANE_SELECTION.md.
- The docs-only candidate packet is now ready:
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md.
- Stop state:
  session_opening_range_expansion_candidate_packet_ready_for_implementation_approval.

Locked fixed baseline:
- Candidate id: btc_15m_session_opening_range_expansion_v1.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv.
- Product/symbol/interval: Binance USDT-M futures BTCUSDT 5m.
- Accepted source facts: 573,984 loaded candles; 2021-01-01T00:00:00Z through
  2026-06-16T23:55:00Z; gap_count=0; duplicate_count=0;
  zero_volume_count=66; comparison_only=false; validation_status=accepted.
- Decision timeframe: exact closed UTC 15m bars resampled from the accepted 5m
  source.
- Session anchor: every UTC calendar date at 13:30:00Z, fixed in UTC with no DST
  shifting and no alternate anchor comparison.
- Opening range: four closed 15m bars with open times 13:30, 13:45, 14:00, and
  14:15 UTC.
- Expansion window: closed decision bars with open times in [14:30, 17:30) UTC
  on the same UTC date.
- Trigger: first same-date closed 15m candle that closes outside the opening
  range by 0.10 * ATR(14), long above the opening-range high or short below the
  opening-range low.
- Entry: next 15m bar open after the closed decision candle.
- Stop: opposite side of the opening range with a 0.10 * ATR(14) buffer.
- Target: 1.5R from entry using initial stop distance.
- Time stop: close after 24 closed 15m bars after entry.
- Sizing/costs: 1% risk at stop, capped at 1x notional; 0.0004 fee per side;
  0.000116 slippage per side.
- Output path:
  results/backtest-first-btc-15m-session-opening-range-expansion-v1/.

Closed boundaries:
- Do not rescue failed trend-pullback with alternate EMA lengths, slope
  lookbacks, pullback windows, EMA-band definitions, continuation triggers, stop
  buffers, target R values, time stops, side selection, session filters, volume
  filters, volatility filters, derivatives-veto interaction, source expansion,
  replay, walk-forward, or optimizer grids.
- Do not rescue closed range-reversion, midpoint, edge-fade, previous-day range,
  value-area, range-optimization, post-compression, clean-breakout, or
  router/rotation branches by renaming them as opening-range expansion.
- Do not retune this opening-range candidate while implementing it. Use the
  candidate packet exactly; no alternate UTC anchor, opening-range length,
  expansion window, acceptance buffer, ATR window, stop, target, time stop,
  one-trade-per-day rule, side selection, weekday filter, volume filter,
  volatility filter, derivatives-veto interaction, source expansion, replay,
  walk-forward, or optimizer grid.
- Do not add paper/testnet/live flow, exchange APIs, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.

Next bounded gate:
- Only if explicitly approved, implement and run exactly one fixed offline
  baseline:
  btc_15m_session_opening_range_expansion_v1.
- Add the minimal Go code and CLI flag needed for that one baseline.
- Write generated artifacts only under:
  results/backtest-first-btc-15m-session-opening-range-expansion-v1/.
- Record source/resample facts, artifact counts, full and split metrics,
  side-separated metrics, pass/fail gates, and the factual verdict in a review
  doc after the run.
- If the baseline fails, close it instead of rescuing it. If it passes all gates,
  stop at a later review approval gate; do not promote or add live flow.

Required reads:
- Read AGENTS.md, README.md, memory/README.md, memory/PROGRESS.md,
  memory/DECISIONS.md, memory/NEXT_CODEX_BRIEF.md,
  docs/BACKTEST_FIRST_RESEARCH_LANE.md, docs/STRATEGY_WORKFLOW.md,
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_LANE_SELECTION.md, and
  docs/BACKTEST_FIRST_SESSION_OPENING_RANGE_EXPANSION_CANDIDATE_PACKET.md.

Verification:
- Run the relevant Go tests before and after the implementation.
- Run the fixed offline backtest command for the approved CLI flag.
- Run git diff --check.

Hard exclusions:
- No source expansion, derivatives-veto interaction, optimizer, replay,
  walk-forward, exchange API, credentials, deployment, paper/testnet/live flow,
  martingale, averaging down, two-exchange logic, or promotion.
- Update memory/PROGRESS.md after the implementation/backtest milestone, update
  memory/DECISIONS.md only if a durable decision is created, and update
  memory/NEXT_CODEX_BRIEF.md to point to the next bounded gate.
```
