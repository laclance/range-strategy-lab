# Next Codex Brief: Rolling Value-Area Reversion Baseline

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
- Inspect git status before editing.

Current state:
- The lab uses the backtest-first research lane for materially different simple
  offline BTCUSDT range-entry ideas.
- The docs-only candidate packet selected exactly one first fixed baseline:
  btc_5m_rolling_value_area_reversion_v1.
- Candidate packet doc:
  docs/BACKTEST_FIRST_CANDIDATE_PACKET.md.
- Stop state:
  backtest_first_candidate_packet_selected_value_reversion_baseline_needs_implementation_approval.

Prior failed backtest context:
- The latest fixed offline backtest implementation was for
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1.
- It stopped at:
  post_compression_directional_expansion_backtest_failed_no_usable_strategy.
- Review doc:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md.
- It reproduced 468 raw representative-cell rows and executed 421 trades, but
  failed economics: full gross P&L 208.560999; engine net -129.258571; extra
  slippage-stress net -227.226250; stress PF 0.799666; recent split gross P&L
  -15.799742 before stress costs.
- The exact fixed post-compression candidate is closed as no usable strategy in
  this form.

Selected implementation candidate:
- Candidate id: btc_5m_rolling_value_area_reversion_v1.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv.
- Product: Binance USDT-M futures.
- Symbol/interval: BTCUSDT 5m.
- Loaded candles: 573,984.
- Coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z.
- Source validation facts: gap_count=0, duplicate_count=0, zero_volume_count=66,
  comparison_only=false, validation_status=accepted.
- Decision timeframe: native closed 5m candles.
- Entry timing: next 5m open.

Fixed baseline rule to implement only after explicit user approval:
- For each closed 5m decision candle d, use prior 288 closed 5m bars [d-288,d-1]
  as the rolling value window.
- Compute range_high, range_low, range_width, rolling VWAP from typical price
  weighted by volume, and ATR(14)[d-1].
- Skip if rolling volume is zero, ATR is missing/non-positive, or range_width is
  less than 6 * ATR(14)[d-1].
- Long candidate: close[d] is in the lower 20% of the prior range and at least
  0.15 * range_width below rolling VWAP.
- Short candidate: close[d] is in the upper 20% of the prior range and at least
  0.15 * range_width above rolling VWAP.
- Do not enter if a position is already open.
- Long stop: min(low[d], range_low) - 0.25 * ATR(14)[d-1].
- Long target: rolling VWAP from the decision window.
- Short stop/target are symmetric.
- Time stop: close after 36 closed 5m bars after entry.
- Sizing/costs: 1% risk at stop, 1x notional cap, 0.0004 fee per side,
  0.000116 slippage per side.

Expected implementation scope, only if approved:
- Add exactly one offline CLI flag for the fixed baseline.
- Suggested output path:
  results/backtest-first-btc-5m-rolling-value-area-reversion-v1/.
- Write source_manifest.json, summary.json, summary.csv, trades.json, and
  strategy-specific sources/signals/skips/trades/summary/falsification CSV/JSON
  files.
- Report split metrics for 2021_2022_stress, 2023_2024_oos,
  2025_2026_recent, and full_2021_2026.
- Report long and short sides separately.

Pass/fail gates:
- Accepted BTCUSDT futures 5m source contract reproduced.
- No leakage in features, thresholds, entry, stops, targets, sizing, ranking, or
  selection.
- At least 120 full-sample executed trades and at least 25 in each primary split.
- Full-sample gross P&L positive.
- 2023_2024_oos and 2025_2026_recent gross P&L not clearly negative.
- Full-sample net after fees/slippage positive.
- Result does not depend only on 2021_2022_stress.
- Drawdown acceptable versus return.
- Long and short side behavior reported separately.

No-rescue boundaries:
- If the fixed baseline fails, do not rescue it with alternate VWAP windows,
  outer-zone percentages, target changes, time-stop changes, side selection,
  volume filters, derivatives-veto interaction, replay, walk-forward, or
  optimizer grids.
- Record the result and move to a materially different candidate.

Hard boundaries:
- Offline research only.
- No paper/testnet/live path.
- No exchange API work, credentials, deploy files, martingale, averaging down, or
  two-exchange logic.
- No derivatives-veto interaction until an independent entry baseline passes and
  the user separately approves interaction testing.
- Do not reopen or rescue failed post-compression, structured compression,
  breakout-retest acceptance, clean breakout continuation, hold-inside/midline,
  impulse absorption, higher-timeframe nested range rotation, range occupancy
  rotation, range router rotation, BTC-regime/ETH/SOL context, or derivatives
  context paths.

If the user asks for the next implementation:
- Implement exactly btc_5m_rolling_value_area_reversion_v1 as fixed above.
- Do not add extra filters or parameter choices.
- Run gofmt, go test ./..., the fixed backtest command, CSV line count checks,
  git diff --check, and git status --short.
- Update README.md only if a new durable doc needs indexing.
- Update memory/PROGRESS.md after the completed milestone.
- Update memory/NEXT_CODEX_BRIEF.md to the resulting review/next gate.
- Update memory/DECISIONS.md only for durable decisions or constraints.
- Commit completed changes and open a PR unless the user explicitly says not to.
```
