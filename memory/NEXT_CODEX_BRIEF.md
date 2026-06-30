# Next Codex Brief: New Candidate Packet

```text
Current state:
- All three candidates from docs/BACKTEST_FIRST_CANDIDATE_PACKET.md failed.
- btc_5m_rolling_value_area_reversion_v1 failed and is closed.
- btc_15m_previous_day_range_reversion_v1 failed and is closed.
- btc_15m_range_edge_exhaustion_fade_v1 failed and is closed.

Latest result:
- btc_15m_range_edge_exhaustion_fade_v1 produced 156 trades.
- Full gross P&L: -154.40528599997904.
- Full net P&L: -261.59525647142874.
- Full profit factor: 0.48125879295748447.
- Full max drawdown: 0.28473381700333156.
- Failed gates: gross edge, net edge, drawdown.

Do not rescue any failed baseline by retuning.

Next allowed research task:
- Create a new backtest-first candidate packet with materially different BTCUSDT
  range-entry ideas.
- Select exactly one new fixed baseline candidate only after explicit user
  approval.
- Offline only; no paper/testnet/live, exchange API, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.
```
