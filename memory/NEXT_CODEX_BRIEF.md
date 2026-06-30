# Next Codex Brief: Backtest-First Research Lane

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/STRATEGY_WORKFLOW.md.
- Read docs/BACKTEST_FIRST_RESEARCH_LANE.md.
- Inspect git status before editing.

Current state:
- The latest fixed offline backtest implementation was for
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1.
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
- Economics failed:
  full gross P&L 208.560999; engine net -129.258571;
  extra slippage-stress net -227.226250; stress PF 0.799666;
  stress max drawdown 0.289326.
- Recent split failed before costs:
  2025_2026_recent gross P&L -15.799742, engine net -106.272938,
  extra stress net -132.510165, stress PF 0.526275.
- The exact fixed candidate is closed as no usable strategy in this form.

Workflow decision:
- The lab is now using a backtest-first research lane for materially different
  simple offline BTCUSDT range-entry ideas.
- Do not create another long premise-spec -> zero-trade audit -> strategy-premise
  spec -> backtest spec chain unless the user explicitly asks for that or the
  idea changes source, timeframe, source family, engine mechanics, promotion
  state, or paper/live readiness.
- Prefer a short candidate packet plus one fixed baseline backtest.
- Become stricter only after a fixed baseline passes.

Non-negotiable safeguards:
- Offline research only.
- Binance USDT-M futures BTCUSDT 5m remains the active source unless a reviewed
  source-scope decision changes it.
- Use confirmed closed-candle decisions only.
- Use next-bar-open entries.
- Preserve accepted source validation, source manifests, fees/slippage,
  stop-first ambiguity, one-position max, split metrics, and reproducible
  artifacts under results/.
- No paper/testnet/live path, exchange API work, credentials, deploy files,
  martingale, averaging down, or two-exchange logic.
- No lookahead in features, thresholds, entry filters, stops, targets, sizing,
  ranking, or selection.

Closed-family boundaries:
- Do not rescue the failed post-compression fixed candidate with adjacent-cell
  P&L selection, stop/target/hold retuning, side changes, volume filters,
  derivatives-veto interaction, replay, walk-forward, or promotion.
- Do not reopen structured compression, breakout-retest acceptance, clean
  breakout continuation, hold-inside/midline, impulse absorption,
  higher-timeframe nested range rotation, range_occupancy_rotation_v1, router
  rotation, BTC-regime/ETH/SOL context, or derivatives-veto paths merely by
  renaming or retuning them.

Valid next work:
- If the user asks to continue strategy research without supplying a specific
  idea, create a compact candidate packet with 3 to 5 materially different
  BTCUSDT range-entry hypotheses.
- Reject candidates that are only retuned closed-family rescues.
- Select the simplest first fixed baseline candidate and state why.
- If the user explicitly asks for implementation, implement exactly one fixed
  offline baseline backtest for that selected candidate.
- The first baseline may include one entry template, one fixed exit model, one
  sizing model, fees/slippage, source manifest, split metrics, and artifacts.
- It must not include optimizer grids, adjacent-cell P&L selection, post-result
  filters, derivatives-veto interaction, replay, walk-forward, paper/testnet/live
  paths, exchange API work, credentials, deploy files, or promotion.

Candidate packet contents:
1. hypothesis;
2. material difference from closed failures;
3. source and timeframe;
4. closed-candle entry rule;
5. fixed stop, target, time stop, sizing, fee, and slippage assumptions;
6. expected output path and artifacts;
7. pass/fail gates;
8. no-rescue boundaries.

Fast-fail rule:
- If the fixed baseline fails gross edge, OOS/recent split quality, net/stress
  edge, trade count, drawdown, leakage, source validation, or explainability,
  record the result and close it. Move to a materially different candidate
  instead of rescuing it by retuning.

Memory and closeout:
- Update README.md only when adding a durable doc that future users need to find
  from the index.
- Update memory/PROGRESS.md after completed milestones.
- Update memory/DECISIONS.md only for durable decisions or constraints.
- Run relevant closeout checks. For docs-only changes, at minimum run markdown or
  diff hygiene checks available in the local environment plus git diff checks.
- Commit completed changes after checks pass unless the user explicitly says not
  to commit.
```
