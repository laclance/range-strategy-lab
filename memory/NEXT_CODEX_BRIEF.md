# Next Codex Brief: Futures Range Universe Discovery Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_SPEC.md
  - docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md
  - docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy discovery.
- The implemented CLI default remains BTCUSDT 5m, but the explicit current research scope is a local BTC/ETH/SOL futures range universe.
- The BTCUSDT clean-breakout baseline stop state was:
  clean_breakout_baseline_failed_no_promotion.
- Clean-breakout baseline facts:
  - clean_breakout_4h_up_h12: 121 trades, gross P&L 2.55, net P&L -69.29, PF 0.8386
  - clean_breakout_1h_all_h12: 1,064 trades, gross P&L 272.63, net P&L -320.25, PF 0.8846
  - aggregate compatibility summary: 1,185 trades, gross P&L 275.18, net P&L -389.54, PF 0.8784
  - both candidates were negative after costs in 2023_2024_oos and 2025_2026_recent
- The user-approved portfolio-stream routing rule did not trigger because neither clean-breakout stream was near-viable after costs.
- The next direction is a non-trading local universe discovery audit, not optimization, not a portfolio stream, and not a direct 15m clean-breakout expansion.
- Default runs still use lab.EmptyStrategy unless an explicit offline audit/backtest/prototype flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

V1 local universe:
- BTCUSDT:
  - ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - 573,985 CSV lines including header
  - quick check: 573,984 data rows, min open 1609459200000, max open 1781654100000, non-monotonic physical row count 0
- ETHUSDT:
  - ../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv
  - 573,985 CSV lines including header
  - quick check: 573,984 data rows, min open 1609459200000, max open 1781654100000, non-monotonic physical row count 0
- SOLUSDT:
  - ../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv
  - 573,985 CSV lines including header
  - quick check: 573,984 data rows, min open 1609459200000, max open 1781654100000, non-monotonic physical row count 1
  - required behavior: validate explicitly and either sort accepted candles before downstream use or fail closed with range_universe_source_gap; do not silently trust filename coverage

Goal:
- Implement -futures-range-universe-discovery-audit as a futures-only, offline-only, non-trading audit.
- Validate the BTC/ETH/SOL local source universe, generate closed UTC 15m/1h/4h resamples, evaluate range-first candidate families, rank symbol/timeframe/family surfaces, and write a review doc.
- If one or two surfaces pass, refresh the next brief as a fixed-rule baseline backtest for those surfaces only.
- If no surface passes, stop with a no-candidate review instead of optimizing or reslicing.

Implementation boundaries:
- No entries, exits, scoring, sizing, optimization, strategy replacement, paper/testnet/live wiring, exchange API use, deployment files, credentials, grid, martingale, averaging down, two-exchange logic, data downloads, broad symbol mining, spot comparison, or sibling repo mutation.
- Do not import sibling repo strategy results, parameters, Python implementation, or runtime code as evidence.
- Keep common summary.* and trades.json zero-trade.
- Keep lab.EmptyStrategy as the default run.
- Treat spot or comparison sources as invalid for this audit.

Audit design:
- Add a Go-native audit module for the local range universe.
- Add CLI flag:
  - -futures-range-universe-discovery-audit
- Source validation:
  - accept only Binance USDT-M futures 5m CSVs for BTCUSDT, ETHUSDT, and SOLUSDT from the approved local paths;
  - require UTC open-time semantics, complete 5m cadence, no duplicate opens, no missing opens, finite positive OHLC, non-negative finite volume, and valid high/low containment;
  - report physical non-monotonic rows separately;
  - accepted downstream candle arrays must be strictly monotonic;
  - reject or mark unusable any symbol/split that cannot satisfy stress 2021_2022, OOS 2023_2024, and recent 2025_2026 eligibility.
- Resampling:
  - 5m native;
  - 15m = 3 complete 5m children;
  - 1h = 12 complete 5m children;
  - 4h = 48 complete 5m children;
  - open first child open, high max child high, low min child low, close last child close, volume sum;
  - reject partial final bars, missing child opens, duplicate child opens, forward fills, and synthetic bars.
- Candidate families:
  - breakout retest / acceptance after completed mature range on 15m, 1h, 4h;
  - boundary touch / rejection on 5m, 15m, 1h;
  - failed breakout re-entry on 5m, 15m, 1h;
  - mature balance rotation / persistence on 1h, 4h;
  - compression-to-expansion on 15m, 1h, 4h only if materially reframed from immediate clean breakout, such as by requiring post-break structure.
- Use closed-candle candidate events only.
- Use fixed discovery horizons and quick-invalidation windows chosen in the implementation from existing local discovery conventions, keeping them predeclared before result inspection.

Required outputs:
- Write results under results/futures-range-universe-discovery-audit/:
  - futures_range_universe_sources.csv/json
  - futures_range_universe_coverage.csv/json
  - futures_range_universe_candidates.csv/json
  - futures_range_universe_summary.csv/json
  - futures_range_universe_rankings.csv/json
  - futures_range_universe_stability.csv/json
  - normal source_manifest.json when applicable, summary.csv/json, and trades.json
- Create docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md after outputs exist.
- Add the review doc to the README docs index.
- Update memory/PROGRESS.md with commands, result paths, source facts, coverage/resample facts, row counts, candidate counts, ranking counts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md based on the stop state.

Review gate:
- Source validation and resample coverage must be accepted for every symbol/surface used.
- A candidate must have adequate BTC counts in every key split and at least one transfer symbol confirming the same family/timeframe/side behavior, unless it is explicitly marked symbol-specific with enough split evidence to justify a baseline exception.
- Favorable rate must beat adverse rate in every key split.
- Quick invalidation must stay below favorable rate.
- Adverse excursion must not dominate favorable excursion.
- Rough cost buffer must survive the weakest split.
- Rankings must prefer cross-symbol confirmation over isolated full-sample strength.

Stop states:
- range_universe_source_gap
- range_universe_no_backtest_candidate
- range_universe_audit_ready_for_baseline_backtest
- range_universe_codegen_or_test_blocked
- range_universe_review_only_no_strategy_change

Test plan:
- Add unit tests for:
  - multi-symbol source validation;
  - non-monotonic physical-row detection and accepted monotonic downstream arrays;
  - gap/duplicate rejection;
  - closed UTC resampling coverage and partial-final handling;
  - all candidate-family labels;
  - quick invalidation, favorable/adverse labels, rough cost buffer, rankings, and stability excluding full split;
  - stop-state selection.
- Add CLI tests that:
  - default runs do not write universe artifacts;
  - the new flag writes source/coverage/candidate/summary/ranking/stability artifacts;
  - spot/comparison sources are rejected;
  - common summary.* and trades.json remain zero-trade.
- Verify:
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
  - env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-discovery-audit -out-dir results/futures-range-universe-discovery-audit
  - wc -l results/futures-range-universe-discovery-audit/*.csv
  - rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
  - git diff --check
  - git status --short

Closeout:
- Commit completed implementation, generated review/doc memory updates, and refreshed next brief after verification unless explicitly told not to.
```
