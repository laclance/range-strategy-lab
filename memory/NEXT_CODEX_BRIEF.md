# Next Codex Brief: Futures Higher-Timeframe Nested Range Rotation Audit

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read only the docs relevant to this task:
  - docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md
  - docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md
  - docs/FUTURES_RANGE_UNIVERSE_BREAKOUT_RETEST_ACCEPTANCE_BASELINE_REVIEW.md
  - docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_WALK_FORWARD_REVIEW.md
  - docs/FUTURES_HYPOTHESIS_PIVOT_INVENTORY.md
  - docs/VERIFICATION.md
- Check git status before editing.

Current state:
- Structured compression passed a fixed ETH/SOL replay but failed
  walk-forward robustness; stop state:
  structured_compression_walk_forward_fragile_needs_review.
- Breakout-retest/acceptance was the one automatic materially different
  follow-up, but its fixed-rule baseline failed after costs; stop state:
  breakout_retest_acceptance_baseline_failed_no_promotion.
- The user approved option 2: a review-only higher-timeframe premise spec
  rather than another automatic backtest.
- The new premise spec is:
  docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_PREMISE_SPEC.md.
- Premise stop state:
  higher_tf_nested_range_rotation_premise_ready_for_audit.
- The premise is BTCUSDT-only, range-only, and non-trading:
  closed UTC 1h child ranges nested inside frozen mature closed UTC 4h parent
  ranges, measuring internal rotation to parent midpoint/far quartile before
  child-range invalidation.
- This is not a rescue retune of structured compression, clean breakout,
  breakout retest, hold-inside/midline, impulse absorption, mature balance
  persistence, or the failed range-universe families.

Goal:
- Implement a non-trading audit behind:
  -futures-higher-tf-nested-range-rotation-audit
- Reuse existing Binance USDT-M futures source validation, closed UTC resample
  behavior, mature range episode/detector plumbing, artifact writer patterns,
  split labels, and zero-trade default compatibility outputs.
- Do not add entries, exits, P&L, strategy replacement, optimizer, replay,
  walk-forward, grid, live/paper/testnet wiring, exchange API, credentials,
  deploy scripts, data downloads, symbol expansion, martingale, averaging down,
  or two-exchange logic.

Source requirements:
- Parent source must be the accepted BTCUSDT Binance USDT-M futures 5m CSV:
  ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv
- Required source facts:
  - loaded candles: 573,984;
  - first open: 2021-01-01T00:00:00Z;
  - last open: 2026-06-16T23:55:00Z;
  - gap_count=0, duplicate_count=0, zero_volume_count=66,
    comparison_only=false, validation_status=accepted.
- Required closed UTC resamples:
  - 1h rows: 47,832, last open 2026-06-16T23:00:00Z;
  - 4h rows: 11,958, last open 2026-06-16T20:00:00Z.
- Reject spot/comparison sources.

Audit rules:
- Detector profile: p30_c12_bollinger_on_adx_off.
- Parent timeframe: 4h.
- Child timeframe: 1h.
- A parent range is mature after at least 12 consecutive detector-qualified
  4h bars; freeze parent high, low, midpoint, upper quartile, and lower
  quartile at the first mature 4h close.
- Ignore child candidates after a closed 4h candle closes outside the frozen
  parent range.
- A child range is mature after at least 12 consecutive detector-qualified
  1h bars.
- Child range eligibility:
  - high/low fully inside frozen parent range;
  - positive width;
  - width no more than 40% of parent width;
  - midpoint in parent lower half for upside rotation candidates;
  - midpoint in parent upper half for downside rotation candidates;
  - at most one event per child range.
- Event definitions:
  - nested_rotation_up: first closed 1h candle after child maturity closes
    above child high, remains inside parent range, closes below parent
    midpoint, and child midpoint is below parent midpoint;
  - nested_rotation_down: first closed 1h candle after child maturity closes
    below child low, remains inside parent range, closes above parent
    midpoint, and child midpoint is above parent midpoint.
- Outcome horizon: 24 closed 1h bars after the event candle.
- Quick invalidation horizon: 6 closed 1h bars.
- Record favorable midpoint, favorable far quartile, adverse child
  invalidation, adverse parent invalidation, no resolution, quick invalidation,
  and excursion fields exactly enough to support the review gate in the spec.
- Skip and count missing coverage, invalid widths, child not inside parent,
  child width above 40%, child midpoint in the wrong parent half, event outside
  parent, event already beyond parent midpoint, and duplicate child events.

Outputs:
- Write artifacts under:
  results/futures-higher-tf-nested-range-rotation-audit/
- Required specific artifacts:
  - futures_higher_tf_nested_range_rotation_sources.csv/json
  - futures_higher_tf_nested_range_rotation_coverage.csv/json
  - futures_higher_tf_nested_range_rotation_parent_ranges.csv/json
  - futures_higher_tf_nested_range_rotation_child_ranges.csv/json
  - futures_higher_tf_nested_range_rotation_events.csv/json
  - futures_higher_tf_nested_range_rotation_summary.csv/json
- Also write common source_manifest.json, summary.csv/json, and trades.json.
- Common summary/trades must remain zero-trade compatibility outputs.
- Add a review doc after the run:
  docs/FUTURES_HIGHER_TIMEFRAME_NESTED_RANGE_ROTATION_AUDIT_REVIEW.md
- Add the review doc to README.md.

Stop states:
- higher_tf_nested_range_rotation_audit_source_gap
- higher_tf_nested_range_rotation_audit_no_candidate_events
- higher_tf_nested_range_rotation_audit_rejected_as_closed_family_reslice
- higher_tf_nested_range_rotation_audit_failed_no_baseline
- higher_tf_nested_range_rotation_audit_ready_for_baseline_brief

Review gate:
- Pass only if source/resample validation passes; at least 100 full-sample
  events exist; every period split has at least 25 events; both sides have at
  least 25 full-sample events or the weaker side is explicitly excluded from
  any future baseline; favorable midpoint beats adverse child invalidation and
  quick invalidation in every period split; favorable far-quartile rate is at
  least 20% in every split; average favorable excursion exceeds average adverse
  excursion in every split; and the result is not carried entirely by one
  historical regime.
- Fail if the only way to pass is to change the 40% child-width gate, 24 bar
  outcome horizon, 6 bar quick-invalidation horizon, or split gates after
  seeing results.

Test plan:
- Add focused lab tests for:
  - 1h/4h closed UTC resample source acceptance;
  - parent range maturity and frozen parent geometry;
  - child range eligibility and skip reasons;
  - nested_rotation_up and nested_rotation_down event emission;
  - duplicate child event handling;
  - outcome labeling order for favorable midpoint, far quartile, adverse child
    invalidation, adverse parent invalidation, no resolution, and quick
    invalidation;
  - summary gate pass/fail stop-state selection.
- Add CLI tests proving:
  - default runs do not write nested range rotation artifacts;
  - the new flag writes all required artifacts;
  - spot comparison is rejected;
  - combinations with trade-producing baseline/replay/optimization/
    walk-forward flags are rejected.

Verification:
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
- env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-higher-tf-nested-range-rotation-audit -out-dir results/futures-higher-tf-nested-range-rotation-audit
- wc -l results/futures-higher-tf-nested-range-rotation-audit/*.csv
- rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
- git diff --check
- git status --short

Closeout:
- Update memory/PROGRESS.md with exact commands, result path, source/resample
  facts, CSV line counts, event counts, summary outcome, and stop state.
- Update memory/DECISIONS.md only if the audit creates a durable
  promotion/no-baseline/no-strategy-change rule.
- Replace memory/NEXT_CODEX_BRIEF.md based on the final stop state.
- Commit the completed audit, review doc, README/memory updates, refreshed
  next brief, and verification evidence after checks pass unless explicitly
  told not to commit.
```
