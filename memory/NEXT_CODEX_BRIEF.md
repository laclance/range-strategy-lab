# Next Codex Brief: BTCUSDT 15m Post-Compression Offline Backtest Implementation Approval Gate

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Do not implement anything unless the user explicitly approves the offline
backtest implementation for the selected candidate:

btc_15m_post_compression_l192_q20_m020_none_long_h48_v1

If the user has not approved that exact implementation, make no docs edits, no
Go code changes, no CLI flag, no generated result directory, no audit/backtest
run, no source download, no network request, no data write, no P&L artifact,
and no strategy/veto work. Report that the project is waiting for explicit
backtest implementation approval and stop.

Before any nontrivial work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_SPEC.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_STRATEGY_PREMISE_SPEC.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_AUDIT_REVIEW.md.
- Read docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_PREMISE_SPEC.md only if exact event-construction details are needed.
- Inspect git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy research. The
  default CLI remains BTCUSDT 5m with lab.EmptyStrategy; trades remain 0 unless
  an explicit offline research flag is passed.
- The selected independent entry-premise family is:
  btc_15m_post_compression_directional_expansion_v1.
- The approved zero-trade audit passed at:
  btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review.
- The docs-only strategy-premise spec stopped at:
  post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval.
- The docs-only backtest spec stopped at:
  post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval.
- The only selected later implementation candidate is:
  btc_15m_post_compression_l192_q20_m020_none_long_h48_v1.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  Binance USDT-M futures BTCUSDT 5m; 573,984 loaded candles;
  2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z; gap_count=0;
  duplicate_count=0; zero_volume_count=66; comparison_only=false;
  validation_status=accepted.
- Exact closed UTC 15m resample facts from the zero-trade audit: 191,328 rows;
  first open 2021-01-01T00:00:00Z; last open 2026-06-16T23:45:00Z; last close
  2026-06-16T23:59:59Z; 3 expected child bars; 0 missing child opens;
  validation accepted.
- The implementation must reproduce the representative zero-trade audit cell
  count before one-position filtering:
  l192_q20_m020_none long h48 full rows = 468.
- The fixed model is long only; next 15m open entry; stop at
  entry_price - 1.0 * ATR(14)[d-1]; target at
  entry_price + 2.0 * ATR(14)[d-1]; max hold 48 closed 15m bars; one-position
  max; stop-first ambiguity; 1% risk-at-stop sizing; 1x notional cap; start
  balance 1000; fee 0.0004 per side; slippage 0.000116 per side.
- The current engine applies slipped execution prices and then subtracts fees in
  NetPnL; it also records slippage as a separate cost proxy. The implementation
  must report both current-engine net and an extra conservative stress view that
  subtracts the recorded slippage proxy once more. Pass/fail decisions must use
  that stress view.
- The derivatives veto candidate remains parked as future skip/retain evidence
  only:
  btc_15m_basis_discount_no_trade_veto_v1.
- The veto may not shape candidate rows, create entries, choose side, rank,
  score P&L, replay, walk forward, optimize, promote a strategy, or reopen
  closed families. Any future veto interaction audit requires a separate
  approval after this fixed backtest result is reviewed.

If the user explicitly approves implementation:
- Add exactly one offline CLI flag:
  -futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest
- Default its output directory to:
  results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/
- Add a focused lab module and tests for this candidate only.
- Reuse existing source validation, exact 15m resampling, ATR, splits,
  backtest engine, CSV/JSON writing, and source-manifest helpers where
  practical.
- Do not add a parameter grid, optimizer, replay, walk-forward, source
  download, source expansion, symbol expansion, derivatives veto interaction,
  paper/testnet/live path, exchange API, credentials, deploy files, martingale,
  averaging down, two-exchange logic, or promotion.

Implementation requirements:
- Candidate construction:
  closed UTC BTCUSDT 15m candles resampled from exact local 5m children; prior
  range [d-192,d-1]; q20 compression against prior 1,920 valid range-width
  observations; decision close above prior range high by 0.2 * ATR(14)[d-1];
  volume mode none; long side only; next 15m open entry timing.
- Fixed exits:
  stop at entry_price - 1.0 * ATR(14)[d-1]; target at
  entry_price + 2.0 * ATR(14)[d-1]; max hold 48 closed 15m bars.
- No-lookahead:
  percentile references exclude d; ATR uses d-1; entry uses open[d+1] only;
  future highs/lows, returns, exits, split close-time outcomes, and P&L never
  shape candidate rows.
- Missingness/skips:
  skip rows explicitly for missing prior range, missing percentile reference,
  missing ATR, missing entry candle, invalid geometry, open position already
  active, missing future path, source/resample mismatch, and any other
  non-silent exclusion.
- Costs:
  emit current-engine net and extra slippage-stress net; use stress net for
  falsification/pass-fail.
- Outputs:
  source_manifest.json; summary.json; summary.csv; trades.json; and the
  strategy-specific sources, resample coverage, signals, skips, trades,
  summary, cost_stress, and falsification CSV/JSON artifacts named in the spec.

Required tests:
- Unit-test event construction: prior range exclusion, exact prior-valid
  percentile window, ATR d-1, q20 threshold, 0.2 ATR breakout, no volume filter,
  long-only side, missing warmup skips, and expected representative candidate
  identity.
- Unit-test signal/backtest geometry: next-open entry anchoring, stop/target
  from slipped entry and ATR d-1, max hold 48, stop-first ambiguity, invalid
  geometry skips, and one-position overlap skips.
- Unit-test summaries: current-engine net versus extra slippage-stress net,
  split assignment by trade close time, exit-reason counts, trade-count gates,
  PF/drawdown gates, source/resample mismatch failure, optimizer contamination
  failure, closed-family reslice failure, and veto contamination failure.
- CLI-test flag wiring, default output directory, exact artifact creation,
  source-product enforcement, spot rejection, and conflict rejection with other
  audit/backtest flags.

Allowed implementation stop states:
- post_compression_directional_expansion_backtest_passed_needs_review
- post_compression_directional_expansion_backtest_failed_no_usable_strategy
- post_compression_directional_expansion_backtest_failed_source_or_resample
- post_compression_directional_expansion_backtest_rejected_optimizer_contamination
- post_compression_directional_expansion_backtest_rejected_closed_family_reslice
- post_compression_directional_expansion_backtest_rejected_veto_contamination

Closeout if implementation is completed:
- Add a review doc:
  docs/FUTURES_BTCUSDT_15M_POST_COMPRESSION_DIRECTIONAL_EXPANSION_BACKTEST_REVIEW.md
- Update README.md intro and docs index.
- Update memory/PROGRESS.md with the stop state, verification commands, result
  path, source facts, trade count, P&L facts, and short factual outcome.
- Update memory/DECISIONS.md only for durable boundaries.
- Refresh this memory/NEXT_CODEX_BRIEF.md to the new canonical next state:
  review/spec gate, veto-interaction gate only if legitimately supported, or no
  selected next implementation.
- Run:
  gofmt -w on changed Go files
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest
  wc -l results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/*.csv
  rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  git diff --check
  git status --short
- After staging, run:
  git diff --cached --check
- Commit completed code/docs/memory changes after checks pass unless told not
  to.
```
