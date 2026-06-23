# Next Codex Brief: Futures Structured Compression ETH/SOL Strategy Spec

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The implemented CLI default remains BTCUSDT 5m, but the active approved strategy-spec path is the local ETH/SOL futures 4h structured-compression universe stream.
- The approved local universe sources are:
  - BTCUSDT: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ETHUSDT: ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - SOLUSDT: ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
  - each source loaded 573,984 candles from 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - gaps / duplicates were 0 / 0 for every symbol
  - zero-volume counts were BTCUSDT=66, ETHUSDT=47, SOLUSDT=47
  - SOLUSDT had physical_non_monotonic_count=1 and was sorted before accepted downstream validation
- The structured-compression optimization stop state was:
  structured_compression_optimization_passed_needs_strategy_spec.
- Optimization outputs are under:
  results/futures-range-universe-structured-compression-optimization/
- Selected config:
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00
- Selected authority symbols:
  ETHUSDT,SOLUSDT
- Diagnostic-only symbol:
  BTCUSDT
- Selected fixed parameters:
  - timeframe: closed UTC 4h resampled from accepted 5m parents
  - detector: p30_c12_bollinger_on_adx_off
  - event: completed mature range, first closed breakout within 24 closed 4h bars, then closed confirmation
  - confirmation window: 2 closed 4h bars
  - max hold: 12 closed 4h bars
  - target: 1.0 completed range width from slipped entry
  - stop buffer: 0.0 range width
  - entry: next 4h bar open
  - fees, slippage, risk sizing, one-position max, and stop-first ambiguity unchanged
- Selected authority result:
  - 129 trades, gross P&L 641.05, net P&L 573.87, PF 1.8089, max DD 9.82%, average net R 0.3465
  - 2021_2022_stress: 54 trades, net P&L 151.79, PF 1.4867
  - 2023_2024_oos: 43 trades, net P&L 229.02, PF 2.2318
  - 2025_2026_recent: 32 trades, net P&L 193.06, PF 1.9121
  - long side: 69 trades, net P&L 385.79, PF 2.0488
  - short side: 60 trades, net P&L 188.08, PF 1.5506
- Important caveats:
  - BTCUSDT is not promoted: diagnostic-only BTCUSDT had 55 trades, net P&L -100.67, PF 0.6507.
  - ETHUSDT was positive full-sample but its recent split was negative after costs.
  - SOLUSDT was positive full-sample but its stress split was negative after costs.
  - The next step must freeze a strategy spec and robustness/rejection rules, not run another grid search.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype/optimization flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Create docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md.
- Freeze the selected ETH/SOL 4h structured-compression candidate into an implementation-ready offline strategy spec.
- Define the next implementation milestone after the spec as an offline candidate strategy replay/backtest only, not live/paper/testnet.

Implementation boundaries:
- This milestone is docs/memory-only unless the user explicitly requests code.
- Do not add CLI flags, entries, exits, scoring changes, sizing changes, new optimization grids, new results directories, live wiring, paper/testnet, exchange API use, credentials, deployment files, data downloads, broad symbol mining, sibling repo mutation, martingale, averaging down, or two-exchange logic.
- Do not promote BTCUSDT. BTC may remain diagnostic-only in the spec.
- Do not reopen the failed 1h structured-compression surface or other range families.
- Do not change source files or generated results.

Required work:
- Create docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md with:
  - source contract for ETHUSDT and SOLUSDT authority sources;
  - BTCUSDT diagnostic-only replay contract;
  - closed UTC 4h resampling contract from the accepted 5m parents;
  - frozen detector/event/signal/entry/stop/target/time-stop rules;
  - portfolio-stream semantics for ETH+SOL authority trades;
  - risk/sizing assumptions inherited from the lab engine;
  - robustness and rejection rules for the next offline replay/backtest milestone;
  - explicit no-live/no-paper/no-testnet boundaries.
- Add the spec doc to the README docs index.
- Update memory/PROGRESS.md with doc path, selected config facts, no-code/no-results boundary, verification commands, and stop state.
- Update memory/DECISIONS.md only if the spec adds a durable rule not already recorded.
- Replace memory/NEXT_CODEX_BRIEF.md with the next implementation brief for the offline candidate strategy replay/backtest, or with a blocker if the spec concludes the candidate is not implementation-ready.

Stop states:
- structured_compression_strategy_spec_ready_for_offline_replay
- structured_compression_strategy_spec_needs_user_scope_decision
- structured_compression_strategy_spec_rejected_no_strategy_change
- structured_compression_strategy_spec_review_only_no_strategy_change

Suggested verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed doc/memory updates after verification unless explicitly told not to.
```
