# Next Codex Brief: Verify Range-Edge Exhaustion Fade Baseline

```text
Current state:
- btc_5m_rolling_value_area_reversion_v1 failed and is closed.
- btc_15m_previous_day_range_reversion_v1 failed and is closed.
- btc_15m_range_edge_exhaustion_fade_v1 implementation has been added.
- Stop state:
  btc_15m_range_edge_exhaustion_fade_backtest_implementation_added_needs_local_verification.

Selected baseline:
- Candidate id: btc_15m_range_edge_exhaustion_fade_v1.
- Flag:
  -backtest-first-btc-15m-range-edge-exhaustion-fade-v1
- Output path:
  results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv.
- Exact 15m resample expected: 191,328 rows, last open 2026-06-16T23:45:00Z.

Required next task:
- Verify locally/CI and run the fixed backtest.
- Record actual result review in
  docs/BACKTEST_FIRST_BTC_15M_RANGE_EDGE_EXHAUSTION_FADE_IMPLEMENTATION_REVIEW.md.

No-rescue boundaries:
- If this fixed baseline fails, do not rescue it with alternate range windows,
  progress thresholds, edge zones, midpoint variants, added volume filters,
  derivatives context, replay, walk-forward, or optimizer grids.
- No paper/testnet/live path, exchange API, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.
```
