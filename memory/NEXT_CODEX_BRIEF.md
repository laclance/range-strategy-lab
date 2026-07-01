# Next Codex Brief: Trend-Pullback Implementation Approval Gate

```text
Current state:
- The bounded offline range optimization workbench failed with no candidate:
  range_optimization_workbench_failed_no_candidate.
- The docs-only strategy-class pivot assessment recommended trend-pullback
  continuation:
  docs/STRATEGY_CLASS_PIVOT_ASSESSMENT.md.
- Assessment stop state:
  strategy_class_pivot_assessment_recommends_trend_pullback.
- The docs-only backtest-first trend-pullback candidate packet has been added:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md.
- Packet stop state:
  trend_pullback_candidate_packet_ready_for_implementation_approval.
- Selected fixed baseline:
  btc_15m_trend_pullback_continuation_v1.
- Source contract remains:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv;
  Binance USDT-M futures BTCUSDT 5m; 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; validation_status=accepted.

Boundaries:
- The current range-reversion / midpoint / edge-fade / previous-day range /
  bounded range-optimization path is closed for now.
- Do not rescue failed range work with retuned thresholds, windows, sessions,
  filters, side selection, derivatives-veto interaction, replay, walk-forward,
  or optimizer grids.
- The trend-pullback packet is not implementation approval by itself.
- Do not change the fixed packet rules after seeing results.
- Do not add optimizer, source expansion, derivatives-veto interaction,
  paper/testnet/live flow, exchange APIs, credentials, deploy files, martingale,
  averaging down, two-exchange logic, or promotion.

Next bounded gate:
- Stop until the operator explicitly approves implementation/backtest of
  btc_15m_trend_pullback_continuation_v1.

If implementation/backtest is explicitly approved:
- Read AGENTS.md, README.md, memory/README.md, memory/PROGRESS.md,
  memory/DECISIONS.md, memory/NEXT_CODEX_BRIEF.md,
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_CANDIDATE_PACKET.md.
- Also read docs/BACKTEST_FIRST_RESEARCH_LANE.md and docs/STRATEGY_WORKFLOW.md.
- Implement exactly one offline fixed baseline matching the packet:
  btc_15m_trend_pullback_continuation_v1.
- Run exactly that fixed backtest against the accepted BTCUSDT Binance USDT-M
  futures source.
- Write generated outputs only under:
  results/backtest-first-btc-15m-trend-pullback-continuation-v1/.
- Add a concise implementation review doc only if the implementation/backtest is
  approved:
  docs/BACKTEST_FIRST_TREND_PULLBACK_CONTINUATION_IMPLEMENTATION_REVIEW.md.
- Update memory with the exact command, output path, stop state, and factual
  pass/fail result.
```
