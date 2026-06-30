# Next Codex Brief: No Selected Implementation After Post-Compression Backtest Failure

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Current state:
- The approved offline backtest implementation for
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1 is complete.
- It stopped at:
  post_compression_directional_expansion_backtest_failed_no_usable_strategy.
- Review doc:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md.
- Result path:
  results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/.
- Source reproduced:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv;
  Binance USDT-M futures BTCUSDT 5m; 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; comparison_only=false;
  validation_status=accepted.
- Exact 15m resample reproduced:
  191,328 rows; first open 2021-01-01T00:00:00Z; last open
  2026-06-16T23:45:00Z; last close 2026-06-16T23:59:59Z;
  3 expected child bars; 0 missing child opens; validation accepted.
- Candidate identity passed:
  expected 468 raw representative-cell rows before one-position filtering,
  got 468.
- Execution produced 421 trades after one-position filtering.
- Trade-count gate passed: full 421, primary splits 152/146/123.
- Economics failed:
  full gross P&L 208.560999; engine net -129.258571;
  extra slippage-stress net -227.226250; stress PF 0.799666;
  stress max drawdown 0.289326.
- Recent split failed before costs:
  2025_2026_recent gross P&L -15.799742, engine net -106.272938,
  extra stress net -132.510165, stress PF 0.526275.
- Falsification failed gross edge, extra slippage-stress edge, stress PF, and
  drawdown gates. Source/resample, candidate identity, leakage, trade count,
  robustness, optimizer-contamination, closed-family, and derivatives-veto
  gates passed.

Important boundary:
- This exact fixed candidate is closed as no usable strategy in this form.
- Do not rescue it with adjacent-cell P&L selection, stop/target/hold retuning,
  volume filters, side changes, the full 81-cell grid, derivatives-veto
  interaction, replay, walk-forward, source expansion, paper/testnet/live paths,
  exchange API work, credentials, deploy files, or promotion.
- The canonical derivatives veto
  btc_15m_basis_discount_no_trade_veto_v1 remains parked because there is no
  passed independent entry stream for it to annotate.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md.
- If the user asks about how we got here, also read:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md,
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md,
  and docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md.
- Inspect git status before editing.

No selected next implementation:
- If the user asks to continue without supplying a materially different premise
  or explicit docs-only scope task, make no code changes, no generated results,
  no backtest/audit run, no source download, and no strategy/veto work.
- Report that the latest fixed post-compression backtest failed and that there
  is no selected next implementation.
- Valid next work requires explicit user direction, such as:
  1. a docs-only hypothesis map for materially different independent entry
     premise candidates;
  2. a new BTCUSDT closed-candle local-source premise spec that is not a
     retuned post-compression/closed-family rescue;
  3. a docs-only stop/parking decision for this research lane.

If a later user supplies and approves a materially different docs-only premise
task:
- Keep it docs-only unless the user explicitly approves implementation.
- Preserve BTCUSDT Binance USDT-M futures as the active source unless a
  separate source-scope review changes it.
- Keep derivatives veto facts parked unless a future independent entry stream
  first passes its own fixed backtest and the user separately approves an
  interaction audit.
- Update README.md and memory only for the completed bounded milestone.
- Run the appropriate docs/memory or code closeout checks.
- Commit completed changes after checks pass unless explicitly told not to.
```
