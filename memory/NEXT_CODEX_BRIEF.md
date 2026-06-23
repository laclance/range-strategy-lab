# Next Codex Brief: Futures Structured Compression ETH/SOL Strategy Replay

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The implemented CLI default remains BTCUSDT 5m, but the active approved strategy replay path is the local ETH/SOL futures 4h structured-compression authority stream.
- The structured-compression strategy spec stop state was:
  structured_compression_strategy_spec_ready_for_offline_replay.
- Frozen config:
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00
- Authority symbols:
  ETHUSDT,SOLUSDT
- Diagnostic-only symbol:
  BTCUSDT
- Approved local universe sources are:
  - BTCUSDT: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - ETHUSDT: ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - SOLUSDT: ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
  - each source loaded 573,984 candles from 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - gaps / duplicates were 0 / 0 for every symbol
  - zero-volume counts were BTCUSDT=66, ETHUSDT=47, SOLUSDT=47
  - SOLUSDT had physical_non_monotonic_count=1 and may be sorted before accepted downstream validation
- Frozen strategy rules:
  - closed UTC 4h resampled from accepted 5m parents;
  - detector p30_c12_bollinger_on_adx_off;
  - completed mature range episodes only;
  - first closed breakout within 24 closed 4h bars after episode end;
  - closed confirmation within 2 closed 4h bars after breakout;
  - up confirmation closes above range high and makes a high above breakout high;
  - down confirmation closes below range low and makes a low below breakout low;
  - signal on the closed confirmation candle;
  - entry on next 4h bar open;
  - long stop at broken range high, short stop at broken range low;
  - target 1.0 completed range width from slipped entry;
  - stop buffer 0.0 range width;
  - max hold 12 closed 4h bars;
  - existing engine fees, slippage, risk sizing, one-position max per symbol replay, and stop-first ambiguity unchanged.
- Selected authority optimization result:
  - 129 trades, gross P&L 641.05, net P&L 573.87, PF 1.8089, max DD 9.82%, average net R 0.3465
  - 2021_2022_stress: 54 trades, net P&L 151.79, PF 1.4867
  - 2023_2024_oos: 43 trades, net P&L 229.02, PF 2.2318
  - 2025_2026_recent: 32 trades, net P&L 193.06, PF 1.9121
  - long side: 69 trades, net P&L 385.79, PF 2.0488
  - short side: 60 trades, net P&L 188.08, PF 1.5506
- Caveats to preserve:
  - BTCUSDT is not promoted: diagnostic-only BTCUSDT had 55 trades, net P&L -100.67, PF 0.6507.
  - ETHUSDT was positive full-sample but its recent split was negative after costs.
  - SOLUSDT was positive full-sample but its stress split was negative after costs.
  - A material replay mismatch must stop for review instead of being tuned around.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype/optimization/replay flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement the frozen offline replay/backtest behind:
  -futures-range-universe-structured-compression-strategy-replay
- This is a fixed candidate strategy replay, not a new grid search.
- Write a review doc and route the next step from the replay result.

Implementation boundaries:
- Reuse the existing universe source validation, SOL sorting rule, closed UTC 4h resampling, structured-compression detector/event helpers, engine costs/sizing, one-position max, max hold, and stop-first behavior.
- Do not add optimization grids, new symbols, new timeframes, new candidate families, scoring search, live wiring, paper/testnet, exchange API use, credentials, deployment files, data downloads, sibling repo mutation, martingale, averaging down, or two-exchange logic.
- Do not promote BTCUSDT. BTCUSDT may be emitted only as diagnostic rows.
- Do not reopen the failed 1h structured-compression surface or other range families.

Required implementation:
- Add a Go-native fixed replay module, likely:
  internal/lab/futures_range_universe_structured_compression_strategy_replay.go
- Add CLI flag:
  -futures-range-universe-structured-compression-strategy-replay
- The flag must:
  - reject spot/comparison sources and incompatible manifests;
  - validate only the approved local BTC/ETH/SOL Binance USDT-M futures files;
  - run ETHUSDT and SOLUSDT as authority;
  - optionally emit BTCUSDT diagnostic rows in the same fixed config;
  - write common trades.json and summary.* using ETH/SOL authority trades only;
  - keep default runs on lab.EmptyStrategy.
- Use a fixed config, not configurable grid values:
  - config_id=sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00
  - candidate_id=structured_compression_4h_all_h12
  - timeframe=4h
  - confirmation_window_bars=2
  - max_hold_bars=12
  - target_range_width_multiple=1.0
  - stop_boundary_buffer_range_width=0.0
  - event_delay_bars=24
  - detector_lookback_days=20
  - detector_percentile=0.30
  - detector_min_consecutive_bars=12
  - Bollinger on, ADX off

Outputs:
- Write under:
  results/futures-range-universe-structured-compression-strategy-replay/
- Write:
  - futures_range_universe_structured_compression_strategy_sources.csv/json
  - futures_range_universe_structured_compression_strategy_coverage.csv/json
  - futures_range_universe_structured_compression_strategy_signals.csv/json
  - futures_range_universe_structured_compression_strategy_trades.csv/json
  - futures_range_universe_structured_compression_strategy_summary.csv/json
  - normal source_manifest.json, summary.csv/json, and trades.json
- Strategy-specific rows should include config_id, symbol role, authority/diagnostic booleans, signal metadata, entry/exit metadata, fees, slippage, gross/net P&L, gross/net R, exit reason, hold bars, split, side, and skipped-signal reason.

Review gate:
- Source and 4h resample validation must be accepted for BTCUSDT, ETHUSDT, and SOLUSDT.
- ETH/SOL authority full-sample trades must be at least 100.
- ETH/SOL authority 2023_2024_oos and 2025_2026_recent trades must each be at least 25.
- Full authority net P&L must be positive after costs and PF at least 1.2.
- Stress, OOS, and recent authority splits must be positive or no worse than the optimization review within a documented replay-tolerance explanation.
- Long and short authority sides must remain positive or no side loss may dominate.
- ETHUSDT and SOLUSDT must each retain positive full-sample net P&L and PF at least 1.0.
- BTCUSDT diagnostics must be reported but cannot make the replay pass.
- If trade count, selected authority rows, or key metrics materially diverge from the optimization review, stop with replay mismatch unless the review proves a harmless artifact-schema difference.

Stop states:
- structured_compression_strategy_replay_source_gap
- structured_compression_strategy_replay_codegen_or_test_blocked
- structured_compression_strategy_replay_regression_or_mismatch
- structured_compression_strategy_replay_failed_no_promotion
- structured_compression_strategy_replay_passed_needs_walk_forward_robustness_brief
- structured_compression_strategy_replay_review_only_no_strategy_change

Review and memory:
- After the full run, create:
  docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md
- Add the review doc to README.md.
- Update memory/PROGRESS.md with commands, result paths, source/resample facts, row counts, trade counts, authority/diagnostic results, and stop state.
- Update memory/DECISIONS.md only if the replay creates a durable promotion/no-promotion rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the stop state.

Test plan:
- Add unit tests for:
  - fixed config values and no grid expansion;
  - approved-path source reuse and spot/comparison rejection;
  - closed UTC 4h resampling;
  - up/down confirmation signals;
  - next-bar-open entry;
  - long/short stop and target geometry;
  - skipped missing entry / non-positive width / invalid geometry signals;
  - max-hold, target, stop-loss, and stop-first behavior;
  - authority vs diagnostic row marking;
  - common trades.json and summary.* excluding BTC diagnostics;
  - replay review gate and stop-state selection.
- Add CLI tests proving:
  - default runs do not write strategy replay artifacts;
  - the new flag writes all required artifacts;
  - spot/comparison sources are rejected;
  - normal outputs remain present and reflect authority-only trades.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-strategy-replay -out-dir results/futures-range-universe-structured-compression-strategy-replay
- wc -l results/futures-range-universe-structured-compression-strategy-replay/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the implementation, generated review doc, memory updates, refreshed next brief, and verification evidence after checks pass unless explicitly told not to.
```
