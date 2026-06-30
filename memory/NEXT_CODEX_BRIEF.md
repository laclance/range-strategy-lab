# Next Codex Brief: Range Optimization Workbench Implementation

```text
Current state:
- All three candidates from docs/BACKTEST_FIRST_CANDIDATE_PACKET.md failed as
  fixed baselines.
- btc_5m_rolling_value_area_reversion_v1 failed and is closed.
- btc_15m_previous_day_range_reversion_v1 failed and is closed.
- btc_15m_range_edge_exhaustion_fade_v1 failed and is closed.
- A docs-only optimization/workbench spec has been added:
  docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md.
- Stop state:
  range_optimization_workbench_spec_ready_for_implementation_approval.

User intent:
- Combining, tweaking, and optimizing range-family components is allowed only as
  a controlled offline discovery workbench.
- It must not be treated as a quiet rescue of failed baselines.
- It must not authorize paper/testnet/live, exchange API work, credentials,
  deploy files, martingale, averaging down, two-exchange logic, or promotion.

Allowed next task only after explicit user approval:
- Implement the bounded offline `-range-optimization-workbench-v1` harness from
  docs/RANGE_OPTIMIZATION_WORKBENCH_SPEC.md.
- Use the fixed source contract and trial logging rules from the spec.
- Emit every trial into a unique immutable run directory under
  `results/range-optimization-workbench-v1/runs/<run_id>/`.
- The implementation must support `-out-dir` and `-run-id`.
- Do not use `rm -rf` on the canonical workbench results parent during
  verification; create a new run id for reruns instead.
- Emit every trial; do not delete failed trials.
- Select at most one candidate for later locked validation.
- Optimizer output alone must stop at either failed/no-candidate, rejected
  overfit risk, or candidate-selected-needs-fixed-validation.
```
