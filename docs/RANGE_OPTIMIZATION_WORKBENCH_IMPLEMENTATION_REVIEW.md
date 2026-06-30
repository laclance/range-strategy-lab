# Range Optimization Workbench Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
range_optimization_workbench_failed_no_candidate
```

The bounded offline workbench specified in `docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md` was implemented and locally verified.

The workbench failed to find a candidate that passed the declared minimum filters. It produced `112` trials, preserved every trial in an immutable run directory, and selected no candidate for fixed validation.

## Scope

Offline CLI flag:

```text
-range-optimization-workbench-v1
```

Run controls:

```text
-run-id <unique_run_id>
-out-dir results/range-optimization-workbench-v1/runs/<run_id>
```

The implementation validates the accepted BTCUSDT Binance USDT-M futures 5m source, exact-resamples closed UTC 15m candles, builds a bounded grid below the 2,500 trial cap, runs each trial through the existing backtest engine, emits every trial, ranks candidates by robustness, selects at most one candidate for a later locked validation lane, and refuses to overwrite an existing run directory.

Candidate counting/reporting rules:

- `top_candidates` contains the full sorted passing-candidate set, not a capped top-10 slice.
- `PassingCandidates` in robustness/falsification artifacts therefore reports the full passing-cell count.
- One-sided candidates are not rejected solely for having zero trades on the other side; side concentration remains a ranking/reporting penalty.

This implementation does not authorize paper/testnet/live trading, exchange API work, credentials, deployment files, production integration, martingale, averaging down, two-exchange execution, or promotion from optimizer output.

## Verified Run

Immutable run path:

```text
results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/
```

Command output:

```text
range_optimization_workbench run_id=20260630T200041Z-78f9a9e trials=112 passing_candidates=0 selected= stop_state=range_optimization_workbench_failed_no_candidate
```

CSV line counts:

```text
     2 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/coverage.csv
   113 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/rejected_candidates.csv
     2 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/source_contract.csv
     1 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/top_candidates.csv
   113 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/trial_results.csv
  1345 results/range-optimization-workbench-v1/runs/20260630T200041Z-78f9a9e/trial_summary.csv
  1576 total
```

Falsification artifact:

```json
{
  "backtest_name": "range_optimization_workbench_v1",
  "run_id": "20260630T200041Z-78f9a9e",
  "stop_state": "range_optimization_workbench_failed_no_candidate",
  "source_resample_pass": true,
  "total_trials": 112,
  "max_trials": 2500,
  "passing_candidates": 0,
  "failure_reasons": [
    "no_candidate_passed_minimum_workbench_filters"
  ]
}
```

Robustness summary:

```json
{
  "run_id": "20260630T200041Z-78f9a9e",
  "total_trials": 112,
  "max_trials": 2500,
  "passing_candidates": 0,
  "rejected_candidates": 112,
  "stop_state": "range_optimization_workbench_failed_no_candidate"
}
```

## Interpretation

The workbench result is clean but negative:

- Source and 15m resample gates passed.
- All `112` trials were emitted and preserved.
- `112` trials were rejected.
- `0` trials passed the minimum workbench filters.
- No selected candidate exists for fixed validation.

This means the bounded first-pass combination/tuning workbench did not recover a robust candidate from the failed range-family components under the declared filters.

## Closure Boundary

This result does not authorize paper/testnet/live, exchange API work, credentials, deployment, promotion, martingale, averaging down, two-exchange logic, or production changes.

Do not promote any workbench cell from this run. If further search is desired, it must be a separately approved spec revision or a materially different research lane, with explicit changes to the search space and guardrails.
