# Next Codex Brief: Previous-Day Range Reversion Baseline

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/STRATEGY_WORKFLOW.md.
- Read docs/BACKTEST_FIRST_RESEARCH_LANE.md.
- Read docs/BACKTEST_FIRST_CANDIDATE_PACKET.md.
- Read docs/BACKTEST_FIRST_BTC_5M_ROLLING_VALUE_AREA_REVERSION_IMPLEMENTATION_REVIEW.md.
- Inspect git status before editing.

Current state:
- The first selected backtest-first fixed baseline,
  btc_5m_rolling_value_area_reversion_v1, was implemented and locally verified.
- Stop state:
  btc_5m_rolling_value_area_reversion_backtest_failed_no_usable_strategy.
- It produced 20,144 trades, but failed gross edge, net edge, and drawdown gates.
- Full gross P&L was -489.49542250301454.
- Full net P&L was -999.9999999684289.
- Full profit factor was 0.6949768102636272.
- Full max drawdown was 0.9999999999689176.
- Do not rescue this candidate with alternate VWAP windows, outer-zone
  percentages, target changes, time-stop changes, side selection, volume filters,
  derivatives-veto interaction, replay, walk-forward, or optimizer grids.

Next candidate from the existing docs-only packet:
- Candidate id: btc_15m_previous_day_range_reversion_v1.
- Candidate packet section: docs/BACKTEST_FIRST_CANDIDATE_PACKET.md.
- Source: current accepted BTCUSDT Binance USDT-M futures 5m CSV resampled to
  exact closed UTC 15m bars.
- Entry timing: next 15m open.
- Hypothesis: when the current UTC day remains inside the previous UTC day's
  high-low range, outer-decile tests of that previous-day range may revert
  toward the previous-day midpoint before a true daily expansion starts.

Allowed next task, only after explicit user approval:
- Implement exactly one fixed offline baseline backtest for
  btc_15m_previous_day_range_reversion_v1.

Fixed candidate rule from packet:
- Build the prior UTC day's high, low, and midpoint from complete 15m candles.
- Skip the current day if any current-day candle before d has closed outside the
  prior day's high-low range.
- Long candidate: close[d] is inside the prior-day range and in its lower 10%.
- Short candidate: close[d] is inside the prior-day range and in its upper 10%.
- Do not enter if a position is already open.
- Stop: beyond the prior-day range extreme by 0.25 * ATR(14)[d-1].
- Target: prior-day midpoint.
- Time stop: 24 closed 15m bars.
- Sizing/costs: 1% risk, 1x notional cap, 0.0004 fee per side, 0.000116 slippage
  per side.

Boundaries:
- No optimizer grid.
- No adjacent-cell P&L selection.
- No post-result filters.
- No derivatives-veto interaction.
- No replay or walk-forward.
- No paper/testnet/live path, exchange API work, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.
- Do not reopen or rescue the failed rolling value-area candidate.
```
