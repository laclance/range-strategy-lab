# Next Codex Brief: Futures Structured Compression 4h Universe Optimization

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The implemented CLI default remains BTCUSDT 5m, but the active approved research path is the local BTC/ETH/SOL futures range-universe funnel.
- The local universe sources are:
  - BTCUSDT: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ETHUSDT: ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - SOLUSDT: ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
  - each source loaded 573,984 candles from 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - gaps / duplicates were 0 / 0 for every symbol
  - zero-volume counts were BTCUSDT=66, ETHUSDT=47, SOLUSDT=47
  - SOLUSDT had physical_non_monotonic_count=1 and was sorted before accepted downstream validation
- Closed UTC resamples used by the structured-compression baseline:
  - 1h: 47,832 rows per symbol, last open 2026-06-16T23:00:00Z
  - 4h: 11,958 rows per symbol, last open 2026-06-16T20:00:00Z
- The structured-compression universe baseline stop state was:
  structured_compression_baseline_passed_needs_optimization_brief.
- Baseline outputs are under:
  results/futures-range-universe-structured-compression-baseline-backtest/
- Baseline facts:
  - 925 signal rows, 912 executed trades, 6 skipped signals.
  - structured_compression_4h_all_h6: 186 signals, 185 trades, gross P&L 395.68, net P&L 302.92, PF 1.3598, max DD 12.92%.
  - structured_compression_1h_all_h12: 739 signals, 727 trades, gross P&L 239.01, net P&L -200.13, PF 0.9362, max DD 35.46%.
  - common combined compatibility view: 912 trades, gross P&L 634.69, net P&L 102.79, PF 1.0258.
- Important constraints from the review:
  - The 4h surface passed as an aggregate universe stream, not as a BTC-only strategy.
  - BTCUSDT was weak on the 4h surface: 56 trades, net P&L -92.27, PF 0.6404.
  - ETHUSDT and SOLUSDT carried the 4h result: net P&L 243.61 and 151.58.
  - The 4h aggregate 2021_2022_stress split was slightly negative after costs at net P&L -4.82.
  - The 1h surface failed and is not approved for optimization.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype/optimization flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement a bounded offline optimization and robustness pass for the passing 4h structured-compression universe stream only.
- Keep the work fixed-source, closed-candle, explainable, and still far from live deployment.
- Decide whether the 4h stream can move to a first candidate strategy spec or must stop as overfit/fragile.

Implementation boundaries:
- No live wiring, paper/testnet, exchange API use, credentials, deployment files, data downloads, broad symbol mining, spot comparison, sibling repo mutation, grid, martingale, averaging down, or two-exchange logic.
- Do not optimize the failed 1h structured-compression surface.
- Do not introduce unrelated candidate families.
- Do not use sibling repo strategy results, parameters, Python implementation, or runtime code as evidence.
- Keep default runs on lab.EmptyStrategy.

Optimization scope:
- Add an explicit offline flag:
  -futures-range-universe-structured-compression-optimization
- Source contract:
  - reuse the approved BTCUSDT, ETHUSDT, and SOLUSDT local Binance USDT-M futures 5m files;
  - validate sources with the existing universe source contract;
  - derive only closed UTC 4h bars from the 5m parents.
- Base candidate:
  - detector p30_c12_bollinger_on_adx_off
  - completed mature range only
  - first closed breakout within 24 4h bars after episode end
  - closed confirmation candle after breakout
  - all sides
  - next 4h bar open entry
- Bound the optimization grid before looking at results:
  - confirmation window bars: 2, 3, 4
  - max hold bars: 4, 6, 8, 12
  - target range-width multiple: 0.75, 1.0, 1.25
  - stop boundary buffer: 0.0 and 0.10 of range width, where long stops move below broken high and short stops move above broken low
  - symbol set: all three symbols, ETH+SOL only, BTC+ETH+SOL with BTC capped to diagnostic-only output
- Keep fees, slippage, risk sizing, one-position max per stream, and stop-first ambiguity unchanged unless the review explicitly routes to a later strategy spec.
- Do not optimize detector percentile, detector min bars, Bollinger/ADX, broad timeframes, or alternate families in this pass.

Required outputs:
- Write results under results/futures-range-universe-structured-compression-optimization/:
  - futures_range_universe_structured_compression_optimization_sources.csv/json
  - futures_range_universe_structured_compression_optimization_coverage.csv/json
  - futures_range_universe_structured_compression_optimization_grid.csv/json
  - futures_range_universe_structured_compression_optimization_trades.csv/json
  - futures_range_universe_structured_compression_optimization_summary.csv/json
  - futures_range_universe_structured_compression_optimization_rankings.csv/json
  - normal summary.csv/json and trades.json for the top selected configuration only, or zero trades if no configuration is selected
- Create docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md after outputs exist.
- Add the review doc to the README docs index.
- Update memory/PROGRESS.md with commands, result paths, source facts, row counts, grid size, top configuration, candidate-level P&L, symbol/split/side facts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the computed stop state.

Review gate:
- Adequacy:
  - selected configuration should have at least 100 full-sample aggregate trades;
  - selected configuration should have at least 25 trades in 2023_2024_oos and 2025_2026_recent;
  - ETH and SOL must each retain adequate full-sample evidence if BTC is excluded or diagnostic-only.
- Pass requires:
  - full net P&L positive after costs;
  - full PF at least 1.2;
  - 2023_2024_oos and 2025_2026_recent net P&L positive after costs;
  - stress split not worse than near-flat after costs;
  - both sides are not dominated by one losing side;
  - drawdown improves or remains acceptable versus the baseline;
  - performance is not produced only by removing BTC while leaving too little transfer evidence.
- If optimization finds a stronger ETH+SOL stream, document it as an ETH/SOL universe candidate, not a BTC strategy.

Stop states:
- structured_compression_optimization_source_gap
- structured_compression_optimization_codegen_or_test_blocked
- structured_compression_optimization_failed_no_promotion
- structured_compression_optimization_passed_needs_strategy_spec
- structured_compression_optimization_mixed_needs_portfolio_stream_review
- structured_compression_optimization_review_only_no_strategy_change

Test plan:
- Add unit tests for:
  - grid generation and fixed bounds;
  - source validation reuse and approved-path enforcement;
  - closed UTC 4h resampling;
  - target multiple and stop buffer geometry for long/short trades;
  - selected top configuration output;
  - symbol-set handling including ETH+SOL and BTC diagnostic-only routing;
  - ranking and stop-state selection for pass, mixed, failed, and source-gap outcomes.
- Add CLI tests that:
  - default runs do not write optimization artifacts;
  - the new flag writes all required artifacts;
  - spot/comparison sources are rejected;
  - normal summary.* and trades.json remain present.
- Verify:
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-optimization -out-dir results/futures-range-universe-structured-compression-optimization
  - wc -l results/futures-range-universe-structured-compression-optimization/*.csv
  - rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  - git diff --check
  - git status --short

Closeout:
- Commit completed implementation, generated review/doc memory updates, refreshed next brief, and verification evidence unless explicitly told not to.
```
