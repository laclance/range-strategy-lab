# Next Codex Brief: Futures Range Occupancy Rotation V1 Optimizer

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md
  - docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_PROTOCOL.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The project is offline Binance USDT-M futures range-strategy construction.
- Default CLI behavior remains BTCUSDT futures 5m with lab.EmptyStrategy unless
  an explicit offline audit/backtest flag is passed.
- The latest spec selected the first v1 grammar:
  range_occupancy_rotation_v1.
- V1 spec stop state:
  range_first_strategy_v1_spec_ready_for_optimizer_implementation.
- Scope remains range-first and BTCUSDT-first.
- The first source remains the accepted Binance USDT-M futures BTCUSDT 5m file:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv.
- Source facts:
  - loaded candles: 573,984;
  - coverage: 2021-01-01T00:00:00Z through 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66;
  - comparison_only=false, validation_status=accepted.
- Expected closed UTC resamples:
  - 15m: 191,328 rows, 2021-01-01T00:00:00Z through
    2026-06-16T23:45:00Z;
  - 1h: 47,832 rows, 2021-01-01T00:00:00Z through
    2026-06-16T23:00:00Z.
- Structured compression, breakout-retest/acceptance, clean breakout,
  hold-inside/midline, impulse absorption, and higher-timeframe nested
  range-rotation are exclusion evidence in their reviewed forms.

Goal:
- Implement the bounded offline BTCUSDT optimizer/backtester behind:
  -futures-range-first-occupancy-rotation-v1-optimization
- Implement exactly the range_occupancy_rotation_v1 grammar, grid, fixed
  baseline, gates, ranking score, artifacts, and stop states from
  docs/FUTURES_RANGE_FIRST_STRATEGY_CONSTRUCTION_V1_SPEC.md.
- Do not invent research policy during implementation.

Key implementation requirements:
- Add lab runner/config/result/row types for range_occupancy_rotation_v1.
- Reuse the existing futures source guard, closed UTC resampling, split labels,
  one-position engine, next-bar-open entries, engine costs, max-hold behavior,
  and stop-first ambiguity.
- Derive signal bars only from closed UTC 15m and 1h BTCUSDT resamples.
- Build range envelopes from the previous lookback bars ending at i-1.
- Emit long/short signals only from close-location occupancy imbalance plus
  interior recapture, exactly as specified.
- Reject and record invalid geometry/skipped signals with the spec reasons.
- Evaluate the fixed baseline row:
  range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005.
- Declare and evaluate the exact 1,152 config grid from the spec.
- Common source_manifest.json, summary.csv/json, and trades.json must describe
  the fixed baseline only.
- Grid, ranking, selected-config, comparison, skipped-signal, and selected
  trade details belong only in V1-specific artifacts.

Expected result directory:
- results/futures-range-first-occupancy-rotation-v1-optimization/

Required artifacts:
- futures_range_first_occupancy_rotation_v1_sources.csv/json
- futures_range_first_occupancy_rotation_v1_coverage.csv/json
- futures_range_first_occupancy_rotation_v1_grid.csv/json
- futures_range_first_occupancy_rotation_v1_baseline.csv/json
- futures_range_first_occupancy_rotation_v1_signals.csv/json
- futures_range_first_occupancy_rotation_v1_trades.csv/json
- futures_range_first_occupancy_rotation_v1_summary.csv/json
- futures_range_first_occupancy_rotation_v1_rankings.csv/json
- futures_range_first_occupancy_rotation_v1_selection.csv/json
- futures_range_first_occupancy_rotation_v1_skips.csv/json
- common source_manifest.json, summary.csv/json, and trades.json

Stop states:
- range_first_strategy_v1_source_gap
- range_first_strategy_v1_no_valid_signals
- range_first_strategy_v1_baseline_failed_grid_still_reviewed
- range_first_strategy_v1_optimizer_failed_no_replay
- range_first_strategy_v1_passed_needs_fixed_replay_spec
- range_first_strategy_v1_rejected_closed_family_reslice

Boundaries:
- BTCUSDT only.
- No ETH/SOL or broad symbol expansion.
- No new data downloads.
- No structured-compression, breakout-retest/acceptance, clean-breakout,
  hold-inside/midline, impulse-absorption, or nested-rotation entries.
- No old binance-bot strategy/scoring/live code.
- No live orders, paper/testnet, exchange API keys, deploy scripts,
  credentials, martingale, averaging down, or two-exchange logic.
- Passing the optimizer gate authorizes only a later fixed replay spec, not a
  strategy package, walk-forward, or live-adjacent work.

Tests:
- Lab tests for source/resample acceptance, range-envelope lookback excluding
  the signal candle, long/short occupancy recapture signals, dual-side skips,
  invalid geometry skips, fixed baseline ID/parameters, grid size 1,152,
  summary gate pass/fail, ranking score/tie-breaks, common-output baseline
  selection, and stop-state selection.
- CLI tests proving default runs do not write V1 artifacts, the new flag writes
  all required artifacts, spot/comparison sources are rejected, and
  combinations with other trade-producing prototype/baseline/optimization/
  replay/walk-forward flags are rejected.

Review and memory after the run:
- Add docs/FUTURES_RANGE_FIRST_OCCUPANCY_ROTATION_V1_OPTIMIZATION_REVIEW.md.
- Add the review doc to README.md.
- Update memory/PROGRESS.md with exact commands, result path, source/resample
  facts, CSV line counts, baseline outcome, selected config/ranking outcome,
  and stop state.
- Update memory/DECISIONS.md only if the optimizer creates a durable
  promotion, no-promotion, or no-strategy-change rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the final stop state.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-first-occupancy-rotation-v1-optimization -out-dir results/futures-range-first-occupancy-rotation-v1-optimization
- wc -l results/futures-range-first-occupancy-rotation-v1-optimization/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the completed implementation, generated review doc, README/memory
  updates, refreshed next brief, and verification evidence after checks pass
  unless explicitly told not to commit.
```
