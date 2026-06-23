# Next Codex Brief: Futures Clean Breakout Baseline Backtest

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for BTCUSDT futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_REVIEW.md
  - docs/FUTURES_RANGE_CANDIDATE_DISCOVERY_SPEC.md
  - docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The active project objective is BTCUSDT futures range-strategy research.
- Active source is Binance USDT-M futures, not spot.
- Authoritative parent source:
  - path: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
  - required product: -source-product binance-usdm-futures
  - market type: Binance USDT-M futures BTCUSDT 5m
  - loaded candles: 573,984
  - CSV lines including header: 573,985
  - open-time coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z
  - accepted manifest facts: gap_count=0, duplicate_count=0, zero_volume_count=66, comparison_only=false, validation_status=accepted
- The futures range candidate discovery audit stop state was:
  range_discovery_audit_ready.
- Discovery result:
  - result dir: results/futures-range-candidate-discovery-audit/
  - complete resamples: 5m=573,984 rows; 15m=191,328; 1h=47,832; 4h=11,958
  - candidate rows: 144,636 across 48,212 distinct events
  - ranking rows: 120
  - passing rows: 24, all from clean_breakout_continuation
  - common summary/trades stayed zero-trade
- Top non-duplicative backtest candidates:
  1. 4h clean_breakout_continuation, up side, 12 bar horizon:
     - full count 496; weakest split count 112
     - full favorable rate 86.69%; weakest favorable rate 82.38%
     - worst adverse rate 16.58%; worst quick invalidation 19.64%
     - weakest cost buffer 4.2910%
  2. 1h clean_breakout_continuation, all sides, 12 bar horizon:
     - full count 4,352; weakest split count 1,213
     - full favorable rate 90.37%; weakest favorable rate 89.01%
     - worst adverse rate 10.99%; worst quick invalidation 17.51%
     - weakest cost buffer 2.5907%
- Do not backtest all 24 passing rows. The first baseline should cover only the two non-duplicative candidates above.
- Touch/rejection/re-entry/balance families did not pass the discovery gate and are not authorized for this baseline.
- Default runs still use lab.EmptyStrategy unless an explicit offline prototype/backtest flag is passed.
- No live orders, exchange keys, deploy scripts, grid, martingale, averaging down, two-exchange execution, paper, or testnet work is allowed.

Goal:
- Implement a fixed-rule, offline-only baseline backtest for the two clean breakout continuation candidates.
- Add a new explicit CLI flag:
  -futures-clean-breakout-baseline-backtest
- Keep default runs on lab.EmptyStrategy.
- Use only Binance USDT-M futures BTCUSDT data and reject spot/comparison sources.

Implementation boundaries:
- Do not optimize parameters.
- Do not add live, paper, testnet, exchange API, deploy, credential, grid, martingale, averaging down, two-exchange, data download, symbol expansion, or sibling repo mutation work.
- Do not add touch/rejection/re-entry/balance entries.
- Do not broaden to 15m in this first baseline; keep 15m clean breakout as later comparison only.

Backtest candidates:
- Candidate A:
  - timeframe: 4h closed UTC resample from accepted 5m parent
  - family: clean_breakout_continuation
  - side: up only
  - signal: closed 4h candle closes above a completed mature range high
  - entry: next 4h bar open, long
  - max hold: 12 closed 4h bars
- Candidate B:
  - timeframe: 1h closed UTC resample from accepted 5m parent
  - family: clean_breakout_continuation
  - side: up and down
  - signal long: closed 1h candle closes above a completed mature range high
  - signal short: closed 1h candle closes below a completed mature range low
  - entry: next 1h bar open
  - max hold: 12 closed 1h bars

Fixed trade template:
- Use completed mature range episodes from the same detector backbone as discovery:
  - percentile 0.30
  - min consecutive bars 12
  - Bollinger on
  - ADX off
  - BarsPerDay 24 for 1h and 6 for 4h
- Signal only on a closed breakout candle after a completed mature range.
- Entry on the next higher-timeframe bar open.
- Long stop: broken range high.
- Short stop: broken range low.
- Long target: entry + one completed range width.
- Short target: entry - one completed range width.
- Max hold: 12 higher-timeframe bars.
- Skip any signal whose next-bar open makes stop/target geometry invalid or non-positive.
- Use the existing engine fee/slippage/risk settings, one-position max, next-bar-open entry, and stop-first ambiguity.

Required outputs:
- Result dir:
  results/futures-clean-breakout-baseline-backtest/
- Write normal source_manifest.json, summary.csv/json, and trades.json.
- Write strategy-specific artifacts:
  - futures_clean_breakout_baseline_signals.csv/json
  - futures_clean_breakout_baseline_trades.csv/json
  - futures_clean_breakout_baseline_summary.csv/json
- Signal rows must include source timeframe, candidate id, episode bounds, breakout candle, side, entry bar, stop, target, max hold bars, and skip reason when skipped.
- Trade rows must join executed trades back to signal metadata and include fees, slippage, gross/net P&L, gross/net R, exit reason, close split, and candidate id.
- Summary rows must split by candidate id, timeframe, side, period split, and full sample.

Review gate:
- Stop with clean_breakout_source_or_resample_gap if source manifest or 1h/4h resampling fails.
- Stop with clean_breakout_codegen_or_test_blocked if implementation or verification cannot complete.
- Stop with clean_breakout_baseline_failed_no_promotion if both candidates fail full/sample/split P&L after costs.
- Stop with clean_breakout_baseline_passed_needs_optimization_brief only if at least one candidate has:
  - full gross and net P&L positive;
  - full profit factor at least 1.2;
  - 2023-2024 and 2025-2026 net P&L non-negative or near-flat with clear gross edge;
  - no single side or split carrying all positive performance;
  - adequate executed trade counts after skips.
- Stop with clean_breakout_review_only_no_strategy_change only if outputs are complete but do not change next direction.

Review and memory:
- Create docs/FUTURES_CLEAN_BREAKOUT_BASELINE_REVIEW.md after outputs exist.
- Add the review doc to README.
- Update memory/PROGRESS.md with commands, result paths, manifest/resample facts, row counts, trade counts, P&L facts, and stop state.
- Update memory/DECISIONS.md only if the review creates a durable promotion/no-promotion rule.
- Refresh memory/NEXT_CODEX_BRIEF.md from the chosen stop state.
- If the baseline passes, the next brief may be an offline optimization brief. If it fails, do not optimize it.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -csv ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv -source-product binance-usdm-futures -futures-clean-breakout-baseline-backtest -out-dir results/futures-clean-breakout-baseline-backtest
- wc -l results/futures-clean-breakout-baseline-backtest/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit completed implementation, generated review/doc memory updates, and refreshed next brief after verification unless explicitly told not to.
```
