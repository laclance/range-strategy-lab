# Next Codex Brief: Verify Range Optimization Workbench

```text
Current state:
- The bounded offline range optimization workbench implementation has been added.
- Implementation review doc:
  docs/RANGE_OPTIMIZATION_WORKBENCH_IMPLEMENTATION_REVIEW.md.
- Stop state:
  range_optimization_workbench_implementation_added_needs_local_run.

Required next task:
- Verify locally/CI and run exactly one immutable workbench run.
- Do not delete any existing workbench run directories.
- Use a unique RUN_ID and OUT_DIR under:
  results/range-optimization-workbench-v1/runs/<run_id>/.

Required commands:
/usr/local/go/bin/gofmt -w \
  cmd/rangelab/workbench.go \
  cmd/rangelab/workbench_run.go \
  cmd/rangelab/workbench_outputs.go \
  cmd/rangelab/workbench_csv.go \
  internal/lab/range_workbench_types.go \
  internal/lab/range_workbench_source.go \
  internal/lab/range_workbench_grid.go \
  internal/lab/range_workbench_strategy.go \
  internal/lab/range_workbench_run.go \
  internal/lab/range_workbench_evaluate.go \
  internal/lab/workbench_score.go \
  internal/lab/workbench_row.go \
  internal/lab/workbench_rank.go

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-$(git rev-parse --short HEAD)"
OUT_DIR="results/range-optimization-workbench-v1/runs/${RUN_ID}"

test ! -e "${OUT_DIR}"

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -range-optimization-workbench-v1 \
  -out-dir "${OUT_DIR}" \
  -run-id "${RUN_ID}"

wc -l "${OUT_DIR}"/*.csv
cat "${OUT_DIR}"/falsification.json
cat "${OUT_DIR}"/robustness_summary.json

git diff --check
git status --short

After verification:
- Record trial count, source/resample facts, passing candidate count, selected
  candidate if any, final stop state, and immutable run path in the review doc.
- Preserve all generated run artifacts locally; do not delete failed/ugly trials.
- Optimizer output alone cannot authorize paper/testnet/live or promotion.
```
