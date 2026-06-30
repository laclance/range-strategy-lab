# Next Codex Brief: Range-Edge Exhaustion Fade Candidate

```text
Current state:
- btc_5m_rolling_value_area_reversion_v1 failed and is closed.
- btc_15m_previous_day_range_reversion_v1 failed and is closed.
- Previous-day result: 1,956 trades, full gross P&L -506.1461042753044,
  full net P&L -935.263471879938, profit factor 0.5772171917112771, max
  drawdown 0.9388204368154983.
- Do not rescue either failed baseline by retuning.

Next candidate, only after explicit user approval:
- btc_15m_range_edge_exhaustion_fade_v1.
- Read docs/BACKTEST_FIRST_CANDIDATE_PACKET.md for the fixed candidate packet.
- Keep it fixed-baseline only: no optimizer, replay, walk-forward, derivatives
  veto interaction, paper/testnet/live path, exchange API, credentials, deploy
  files, martingale, averaging down, two-exchange logic, or promotion.
```
