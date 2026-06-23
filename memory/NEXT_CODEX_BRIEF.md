# Next Codex Brief: Futures Range-Universe Breakout Retest Acceptance Baseline

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_RANGE_UNIVERSE_POST_STRUCTURED_COMPRESSION_PIVOT_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_DISCOVERY_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_BASELINE_REVIEW.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- The frozen ETH/SOL 4h structured-compression stream
  sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00 passed fixed replay but
  failed walk-forward robustness.
- Walk-forward stop state:
  structured_compression_walk_forward_fragile_needs_review.
- Post-compression pivot stop state:
  post_structured_compression_pivot_ready_for_breakout_retest_acceptance_baseline.
- The structured-compression walk-forward result is exclusion evidence:
  do not retune confirmation, max hold, target, stop buffer, symbol set, or
  training gates; do not reopen the failed 1h structured-compression surface;
  BTCUSDT remains diagnostic-only for the completed structured-compression
  stream.
- The only next automatic range-universe premise authorized by the pivot review
  is a bounded offline fixed-rule baseline for breakout retest/acceptance after
  a completed mature range, selected from existing
  results/futures-range-universe-discovery-audit/ evidence.

Goal:
- Implement a bounded offline baseline backtest behind a new explicit flag,
  suggested name:
  -futures-range-universe-breakout-retest-acceptance-baseline-backtest
- Before coding the baseline rules, identify the top one or two
  non-duplicative passing breakout_retest_acceptance rows from the existing
  range-universe discovery artifacts.
- Stop with a review-only no-candidate result if those artifacts do not contain
  a passing non-duplicative breakout_retest_acceptance row.
- Reuse existing BTC/ETH/SOL Binance USDT-M futures source validation, SOL
  sorting acceptance, closed UTC resampling, one-position max, engine costs,
  next-bar-open entries, max-hold handling, and stop-first behavior.

Boundaries:
- Do not rerun or retune structured compression.
- Do not add optimization dimensions, new symbols, data downloads, broad symbol
  mining, paper/testnet/live wiring, exchange APIs, credentials, deploy files,
  martingale, averaging down, or two-exchange logic.
- Do not promote BTCUSDT or any symbol from this baseline before a review doc
  passes after-cost, split, side, drawdown, and symbol-transfer gates.
- Default cmd/rangelab runs must continue to use lab.EmptyStrategy and must not
  write breakout-retest artifacts unless the new flag is passed.

Expected implementation shape:
- Add focused lab code for:
  - selecting the top non-duplicative breakout_retest_acceptance candidate rows
    from the existing discovery ranking/candidate artifacts;
  - replaying those fixed candidates as baseline rules without inspecting
    result P&L to change rules;
  - emitting source, coverage, signal, trade, summary, and review rows.
- Add CLI wiring for the new explicit flag, futures-only source guard,
  incompatible flag checks, artifact writers, and console summary.
- Write artifacts under:
  results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/
  with CSV/JSON files named consistently with the existing futures range
  universe baseline patterns, plus normal source_manifest.json, summary.csv,
  summary.json, and trades.json.
- Evaluate BTCUSDT, ETHUSDT, and SOLUSDT independently and in aggregate, but do
  not use the failed structured-compression BTC diagnostic result as promotion
  evidence for this new premise.

Review gates:
- Source and resample validation must pass for all required local futures
  sources.
- Existing discovery artifacts must contain at least one passing,
  non-duplicative breakout_retest_acceptance row.
- A passing baseline must have adequate trade count, positive net P&L after
  costs, PF at least 1.2, positive 2023_2024_oos and 2025_2026_recent splits,
  no severe side weakness, no single-symbol domination that hides transfer
  failure, and explainable drawdown.
- Stop states:
  - breakout_retest_acceptance_baseline_source_gap
  - breakout_retest_acceptance_baseline_no_ranked_candidate
  - breakout_retest_acceptance_baseline_failed_no_promotion
  - breakout_retest_acceptance_baseline_passed_needs_robustness_review

Docs and memory:
- Add a concise review doc after the run:
  docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
- Add the review doc to README.md.
- Update memory/PROGRESS.md with exact commands, result paths, source/resample
  facts, CSV line counts, candidate selection, baseline outcomes, and stop
  state.
- Update memory/DECISIONS.md only if the baseline creates a durable promotion,
  no-promotion, or no-strategy-change rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the final stop state.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-universe-breakout-retest-acceptance-baseline-backtest -out-dir results/futures-range-universe-breakout-retest-acceptance-baseline-backtest
- wc -l results/futures-range-universe-breakout-retest-acceptance-baseline-backtest/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Commit the completed implementation, generated review doc, README/memory
  updates, refreshed next brief, and verification evidence after checks pass
  unless explicitly told not to commit.
```
