# Next Codex Brief: Trend-Pullback Candidate Packet

```text
Current state:
- The bounded offline range optimization workbench has been implemented and run:
  docs/RANGE_OPTIMIZATION_WORKBENCH_IMPLEMENTATION_REVIEW.md.
- Run id: 20260630T200041Z-78f9a9e.
- Total trials: 112.
- Passing candidates: 0.
- Rejected candidates: 112.
- Selected candidate: none.
- Stop state: range_optimization_workbench_failed_no_candidate.
- The docs-only strategy-class pivot assessment has been added:
  docs/STRATEGY_CLASS_PIVOT_ASSESSMENT.md.
- Stop state:
  strategy_class_pivot_assessment_recommends_trend_pullback.
- Recommended next research lane:
  trend-pullback continuation.

Boundaries:
- The current range-reversion / midpoint / edge-fade / previous-day range /
  bounded range-optimization path is closed for now.
- Do not rescue failed range work with retuned thresholds, windows, sessions,
  filters, side selection, derivatives-veto interaction, replay, walk-forward,
  or optimizer grids.
- Do not implement a strategy or run a backtest from this brief.

Next allowed artifact after operator approval:
- Create exactly one docs-only backtest-first candidate packet:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md.
- The packet should define one fixed BTCUSDT futures trend-pullback continuation
  baseline, likely on exact closed UTC 15m bars from the accepted 5m source.
- It must lock the source, timeframe, trend definition, pullback definition,
  continuation trigger, stop, target, time stop, sizing, fees, slippage,
  pass/fail gates, side reporting, split reporting, output path, and no-rescue
  boundaries.
- No Go code, CLI flags, generated results, optimizer, source expansion,
  derivatives-veto interaction, paper/testnet/live path, exchange API,
  credentials, deploy files, martingale, averaging down, two-exchange logic, or
  promotion is authorized.
```
