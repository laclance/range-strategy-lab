# Next Codex Brief: Futures Structured Compression Walk-Forward Robustness

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_SPEC.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_OPTIMIZATION_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The active candidate stream is the local ETH/SOL futures 4h structured-compression authority stream.
- The strategy replay stop state was:
  structured_compression_strategy_replay_passed_needs_walk_forward_robustness_brief.
- Frozen replay config:
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00
- Authority symbols:
  ETHUSDT,SOLUSDT
- Diagnostic-only symbol:
  BTCUSDT
- Replay source facts:
  - BTCUSDT, ETHUSDT, and SOLUSDT approved local Binance USDT-M futures 5m files each loaded 573,984 candles from 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z.
  - gaps / duplicates were 0 / 0 for every symbol.
  - zero-volume counts were BTCUSDT=66, ETHUSDT=47, SOLUSDT=47.
  - SOLUSDT had physical_non_monotonic_count=1, was sorted, and was accepted.
  - closed UTC 4h resamples had 11,958 rows per symbol from 2021-01-01T00:00:00Z through 2026-06-16T20:00:00Z.
- Replay result:
  - strategy-specific rows: 185 signals and 184 trades;
  - common authority outputs: 129 ETH/SOL trades only;
  - full authority gross P&L 641.05, net P&L 573.87, PF 1.8089, max DD 9.82%;
  - 2021_2022_stress: 54 trades, net P&L 151.79, PF 1.4867;
  - 2023_2024_oos: 43 trades, net P&L 229.02, PF 2.2318;
  - 2025_2026_recent: 32 trades, net P&L 193.06, PF 1.9121;
  - long side: 69 trades, net P&L 385.79, PF 2.0488;
  - short side: 60 trades, net P&L 188.08, PF 1.5506.
- Caveats to preserve:
  - BTCUSDT is not promoted: diagnostic-only BTCUSDT had 55 trades, net P&L -100.67, PF 0.6507.
  - ETHUSDT is positive full-sample but negative in 2025_2026_recent.
  - SOLUSDT is positive full-sample but negative in 2021_2022_stress.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype/optimization/replay/robustness flag is passed.
- No live orders, exchange keys, deploy scripts, grid/martingale/averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement a bounded offline walk-forward robustness pass behind:
  -futures-range-universe-structured-compression-walk-forward-robustness
- This is not a new open-ended optimization. It may reuse only the already declared 4h structured-compression optimization grid to test forward-selection robustness and compare it with the frozen ETH/SOL replay config.

Implementation boundaries:
- Reuse existing universe source validation, SOL sorting rule, closed UTC 4h resampling, structured-compression strategy helpers, existing optimization grid definitions, engine costs/sizing, one-position max, max hold, and stop-first behavior.
- Do not add new symbols, new timeframes, new candidate families, new grid dimensions, scoring search outside the existing ranking logic, live wiring, paper/testnet, exchange API use, credentials, deployment files, data downloads, sibling repo mutation, martingale, averaging down, or two-exchange logic.
- Do not promote BTCUSDT. BTCUSDT may be emitted only as diagnostic rows or as part of historical grid comparison rows clearly marked non-authority.
- Do not reopen the failed 1h structured-compression surface or other range families.

Required implementation:
- Add a Go-native walk-forward robustness module, likely:
  internal/lab/futures_range_universe_structured_compression_walk_forward.go
- Add CLI flag:
  -futures-range-universe-structured-compression-walk-forward-robustness
- The flag must:
  - reject spot/comparison sources and incompatible manifests;
  - validate only the approved local BTC/ETH/SOL Binance USDT-M futures files;
  - run the existing 4h structured-compression grid exactly as declared in the optimization review;
  - evaluate forward-selection folds without adding grid dimensions;
  - always include a row for the frozen replay config;
  - keep default runs on lab.EmptyStrategy;
  - keep normal trades.json and summary.* frozen ETH/SOL authority-only, matching the replay milestone; BTCUSDT rows may appear only in strategy-specific diagnostic artifacts.

Walk-forward folds:
- wf_2021_2022_train__2023_2024_test:
  - select on 2021_2022_stress;
  - test on 2023_2024_oos.
- wf_2021_2024_train__2025_2026_test:
  - select on 2021_2022_stress plus 2023_2024_oos;
  - test on 2025_2026_recent.
- wf_2023_2024_train__2025_2026_test:
  - select on 2023_2024_oos;
  - test on 2025_2026_recent.

Selection rule:
- Candidate configs are the existing bounded optimization grid only:
  - confirmation window 2,3,4;
  - max hold 4,6,8,12;
  - target multiple 0.75,1.0,1.25;
  - stop buffer 0.0,0.10;
  - symbol sets BTC_ETH_SOL, ETH_SOL, and BTC_DIAGNOSTIC_ETH_SOL.
- A fold-selected config must pass the same training adequacy basics used by optimization:
  - at least 100 aggregate authority training trades when the train set contains multiple period splits;
  - at least 25 authority training trades in each individual train split segment used by the fold;
  - positive training net after costs;
  - training PF at least 1.2;
  - no side or authority symbol weakness dominating.
- Ranking may reuse the existing optimization rank score over the training rows.
- BTC_DIAGNOSTIC_ETH_SOL means ETH/SOL determine authority; BTC rows are diagnostic.

Outputs:
- Write under:
  results/futures-range-universe-structured-compression-walk-forward-robustness/
- Write:
  - futures_range_universe_structured_compression_walk_forward_sources.csv/json
  - futures_range_universe_structured_compression_walk_forward_coverage.csv/json
  - futures_range_universe_structured_compression_walk_forward_grid.csv/json
  - futures_range_universe_structured_compression_walk_forward_folds.csv/json
  - futures_range_universe_structured_compression_walk_forward_trades.csv/json
  - futures_range_universe_structured_compression_walk_forward_summary.csv/json
  - futures_range_universe_structured_compression_walk_forward_rankings.csv/json
  - normal source_manifest.json, summary.csv/json, and trades.json
- Fold rows must include fold id, train split set, test split, selected config id, frozen config id, selected-vs-frozen rank, authority symbols, diagnostic symbols, train/test trade counts, train/test net P&L, train/test PF, drawdown, side/symbol caveat flags, and pass/fail reason.

Review gate:
- Source and 4h resample validation must be accepted for BTCUSDT, ETHUSDT, and SOLUSDT.
- The frozen replay config must be present in every fold.
- At least two of three folds must select the frozen config or a config with the same ETH/SOL authority, same confirmation window 2, max hold 12, and non-worse test result.
- Every fold test must have at least 25 authority trades for the selected config or the frozen config comparison.
- Selected test net P&L must be positive after costs in every fold.
- Frozen-config test net P&L must be positive in every fold where it has at least 25 trades.
- No fold may require BTCUSDT authority to pass.
- ETH/SOL transfer evidence must remain visible; BTC diagnostic weakness cannot be used to offset or justify promotion.
- If the walk-forward result depends on changing the selected config away from the frozen replay config, stop for review instead of changing strategy rules.

Stop states:
- structured_compression_walk_forward_source_gap
- structured_compression_walk_forward_codegen_or_test_blocked
- structured_compression_walk_forward_failed_no_promotion
- structured_compression_walk_forward_fragile_needs_review
- structured_compression_walk_forward_passed_needs_candidate_strategy_package
- structured_compression_walk_forward_review_only_no_strategy_change

Review and memory:
- After the full run, create:
  docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
- Add the review doc to README.md.
- Update memory/PROGRESS.md with commands, result paths, source/resample facts, CSV counts, fold counts, selected/frozen config results, and stop state.
- Update memory/DECISIONS.md only if the walk-forward review creates a durable promotion/no-promotion rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the stop state.

Test plan:
- Add unit tests for:
  - fold definitions and split membership;
  - grid reuse without new dimensions;
  - BTC diagnostic-only symbol-set handling;
  - source validation and 4h resampling reuse;
  - train/test summary extraction;
  - fold selection ranking;
  - frozen-config comparison rows;
  - stop-state selection;
  - CLI default writes no walk-forward artifacts;
  - CLI flag writes all required artifacts and rejects spot comparison.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-structured-compression-walk-forward-robustness -out-dir results/futures-range-universe-structured-compression-walk-forward-robustness
- wc -l results/futures-range-universe-structured-compression-walk-forward-robustness/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the implementation, generated review doc, memory updates, refreshed next brief, and verification evidence after checks pass unless explicitly told not to.
```
