# Next Codex Brief: Futures Structured Compression Universe Baseline Backtest

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md
  - docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The implemented CLI default remains BTCUSDT 5m, but the active approved research path is now the local BTC/ETH/SOL futures range-universe funnel.
- The local universe discovery audit stop state was:
  range_universe_audit_ready_for_baseline_backtest.
- The v1 local universe sources were accepted:
  - BTCUSDT: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ETHUSDT: ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - SOLUSDT: ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
  - each source loaded 573,984 candles from 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - gaps / duplicates were 0 / 0 for every symbol
  - zero-volume counts were BTCUSDT=66, ETHUSDT=47, SOLUSDT=47
  - SOLUSDT had physical_non_monotonic_count=1 and was sorted before accepted downstream validation
- Closed UTC resamples were accepted for every symbol:
  - 5m: 573,984 rows
  - 15m: 191,328 rows
  - 1h: 47,832 rows
  - 4h: 11,958 rows
- The universe audit emitted 415,953 candidate rows and 147 ranked surfaces.
- 35 rows passed the universe gate:
  - 21 structured_compression_expansion rows
  - 14 breakout_retest_acceptance rows
- The authorized next baseline surfaces are:
  - structured_compression_4h_all_h6: 4h structured_compression_expansion, all sides, horizon/max hold 6 higher-timeframe bars
  - structured_compression_1h_all_h12: 1h structured_compression_expansion, all sides, horizon/max hold 12 higher-timeframe bars
- These surfaces are not the failed immediate clean-breakout baseline. They require post-break closed-candle structure after a completed mature range.
- Breakout retest/acceptance is secondary evidence only for now.
- Boundary touch rejection, single-candle wick rejection, failed breakout re-entry, and mature balance persistence are not approved for baseline backtest from this audit.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement -futures-range-universe-structured-compression-baseline-backtest as an offline-only fixed-rule baseline for the two approved structured-compression surfaces.
- Evaluate BTCUSDT, ETHUSDT, and SOLUSDT independently first, then report aggregate compatibility.
- Move toward an actual strategy only through fixed baseline evidence: no optimization until the baseline is positive or near-positive after costs with stable splits.

Implementation boundaries:
- No optimization, scoring search, parameter sweep, live wiring, paper/testnet, exchange API use, credentials, deployment files, grid, martingale, averaging down, two-exchange logic, data downloads, broad symbol mining, spot comparison, or sibling repo mutation.
- Do not backtest breakout_retest_acceptance, boundary_touch_rejection, single_candle_wick_rejection, failed_breakout_reentry, or mature_balance_persistence in this milestone.
- Do not import sibling repo strategy results, parameters, Python implementation, or runtime code as evidence.
- Keep common default runs on lab.EmptyStrategy.
- The new baseline flag may produce trades; default runs must remain zero-trade.
- Reject spot/comparison sources and reject any universe source outside the approved local BTC/ETH/SOL Binance USDT-M futures files.

Source and resampling:
- Reuse the universe source-validation contract from -futures-range-universe-discovery-audit.
- Accepted source paths:
  - ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
- Use closed UTC resamples from the accepted 5m parent:
  - 1h = 12 complete child bars
  - 4h = 48 complete child bars
- Reject partial final bars, missing child opens, duplicate child opens, forward fills, synthetic candles, or non-UTC/future-looking semantics.

Detector and event definition:
- Detector backbone:
  - compression detector percentile 0.30
  - min consecutive bars 12
  - Bollinger on
  - ADX off
  - interval-specific BarsPerDay: 1h=24, 4h=6
- Use completed mature range episodes only; candidate boundaries must be known at the completed episode end.
- For each episode, scan for the first closed breakout candle within 24 higher-timeframe bars after episode end.
- Up structured-compression event:
  - breakout candle closes above the completed range high;
  - within the next 3 higher-timeframe bars, a closed confirmation candle closes above the completed range high and makes a high above the breakout candle high;
  - signal on the closed confirmation candle.
- Down structured-compression event:
  - breakout candle closes below the completed range low;
  - within the next 3 higher-timeframe bars, a closed confirmation candle closes below the completed range low and makes a low below the breakout candle low;
  - signal on the closed confirmation candle.
- Entry is next higher-timeframe bar open.
- If multiple events compete, keep deterministic ordering and one open position max through the existing engine.

Trade template:
- Long:
  - entry at next higher-timeframe bar open with existing slippage model;
  - stop at the broken range high;
  - target at slipped entry plus one completed range width.
- Short:
  - entry at next higher-timeframe bar open with existing slippage model;
  - stop at the broken range low;
  - target at slipped entry minus one completed range width.
- Skip missing entry bars, non-positive range width, non-positive prices, non-finite values, and invalid stop/target geometry.
- Use the existing engine for fees, slippage, risk sizing, one-position max, max hold, and stop-first ambiguity.
- Max hold:
  - 6 higher-timeframe bars for structured_compression_4h_all_h6.
  - 12 higher-timeframe bars for structured_compression_1h_all_h12.

Required outputs:
- Write results under results/futures-range-universe-structured-compression-baseline-backtest/:
  - futures_range_universe_structured_compression_baseline_sources.csv/json
  - futures_range_universe_structured_compression_baseline_coverage.csv/json
  - futures_range_universe_structured_compression_baseline_signals.csv/json
  - futures_range_universe_structured_compression_baseline_trades.csv/json
  - futures_range_universe_structured_compression_baseline_summary.csv/json
  - normal summary.csv/json and trades.json
- Create docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md after outputs exist.
- Add the review doc to the README docs index.
- Update memory/PROGRESS.md with commands, result paths, source facts, resample facts, row counts, signal/trade counts, candidate-level P&L, aggregate compatibility facts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the computed stop state.

Review gate:
- Adequacy:
  - each selected surface should have at least 100 full-sample executed trades;
  - each selected surface should have at least 25 executed trades in 2023_2024_oos and 2025_2026_recent, or the review must document a sample-size gap.
- Pass requires:
  - full net P&L positive after costs;
  - full PF at least 1.2;
  - 2023_2024_oos and 2025_2026_recent not negative after costs for BTCUSDT and at least one transfer symbol, or a clearly positive aggregate with documented symbol weakness;
  - no side, symbol, or timeframe weakness dominates the result;
  - drawdown and average net R are not obviously incompatible with a later optimization pass.
- Mixed routing:
  - use structured_compression_baseline_mixed_needs_portfolio_stream_review only if at least two streams are near-flat or positive after costs, no single stream is strong enough alone, and the combination has a plausible risk-diversification reason beyond averaging losers.
- Optimization is allowed in the next brief only if the fixed baseline passes the review gate.

Stop states:
- structured_compression_baseline_source_gap
- structured_compression_baseline_codegen_or_test_blocked
- structured_compression_baseline_failed_no_promotion
- structured_compression_baseline_passed_needs_optimization_brief
- structured_compression_baseline_mixed_needs_portfolio_stream_review
- structured_compression_baseline_review_only_no_strategy_change

Test plan:
- Add unit tests for:
  - multi-symbol source validation reuse and approved-path enforcement;
  - closed UTC 1h/4h resampling;
  - structured-compression signal detection for up/down events;
  - closed confirmation candle and next-bar-open entry;
  - long/short stop/target geometry;
  - skipped signals for missing entry and invalid geometry;
  - max-hold exits, stop/target exits, and stop-first ambiguity;
  - metadata joins from signal to trade rows;
  - summary rows by candidate, symbol, split, side, and aggregate;
  - stop-state selection for pass, mixed, failed, and source-gap outcomes.
- Add CLI tests that:
  - default runs do not write structured-compression baseline artifacts;
  - the new flag writes all required artifacts;
  - spot/comparison sources are rejected;
  - normal summary.* and trades.json remain present.
- Verify:
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-baseline-backtest -out-dir results/futures-range-universe-structured-compression-baseline-backtest
  - wc -l results/futures-range-universe-structured-compression-baseline-backtest/*.csv
  - rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  - git diff --check
  - git status --short

Closeout:
- Commit completed implementation, generated review/doc memory updates, refreshed next brief, and verification evidence unless explicitly told not to.
```
