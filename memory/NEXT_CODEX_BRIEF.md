# Next Codex Brief: Verify Previous-Day Range Reversion Baseline

```text
We are in /home/lance/range-strategy-lab, a standalone offline Go project for
Binance USDT-M futures range-strategy research.

Before work:
- Read AGENTS.md.
- Read memory/README.md, memory/PROGRESS.md, and memory/DECISIONS.md.
- Read README.md as the docs index.
- Read docs/STRATEGY_WORKFLOW.md.
- Read docs/BACKTEST_FIRST_RESEARCH_LANE.md.
- Read docs/BACKTEST_FIRST_CANDIDATE_PACKET.md.
- Read docs/BACKTEST_FIRST_BTC_5M_ROLLING_VALUE_AREA_REVERSION_IMPLEMENTATION_REVIEW.md.
- Read docs/BACKTEST_FIRST_BTC_15M_PREVIOUS_DAY_RANGE_REVERSION_IMPLEMENTATION_REVIEW.md.
- Inspect git status before editing.

Current state:
- The first selected backtest-first fixed baseline,
  btc_5m_rolling_value_area_reversion_v1, failed and is closed as no usable
  strategy.
- The next fixed baseline implementation for
  btc_15m_previous_day_range_reversion_v1 has been added.
- Stop state:
  btc_15m_previous_day_range_reversion_backtest_implementation_added_needs_local_verification.
- Review doc:
  docs/BACKTEST_FIRST_BTC_15M_PREVIOUS_DAY_RANGE_REVERSION_IMPLEMENTATION_REVIEW.md.

Selected baseline:
- Candidate id: btc_15m_previous_day_range_reversion_v1.
- Flag:
  -backtest-first-btc-15m-previous-day-range-reversion-v1
- Output path:
  results/backtest-first-btc-15m-previous-day-range-reversion-v1/.
- Source: ../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv.
- Source facts to reproduce: 573,984 loaded candles; 2021-01-01T00:00:00Z
  through 2026-06-16T23:55:00Z; gap_count=0; duplicate_count=0;
  zero_volume_count=66; comparison_only=false; validation_status=accepted.
- Exact 15m resample expected: 191,328 rows, last open 2026-06-16T23:45:00Z.

Required next task:
- Verify the implementation locally/CI and run the fixed backtest.
- Record the actual result review.

Required commands:
/usr/local/go/bin/gofmt -w \
  cmd/rangelab/day.go \
  cmd/rangelab/day_outputs.go \
  internal/lab/backtest_first_btc_15m_day_types.go \
  internal/lab/backtest_first_btc_15m_day_rows.go \
  internal/lab/backtest_first_btc_15m_day_runner.go \
  internal/lab/backtest_first_btc_15m_day_support.go \
  internal/lab/day_state.go \
  internal/lab/pdr.go \
  internal/lab/pdr_row.go \
  internal/lab/pdrh.go \
  internal/lab/pdr_inside.go \
  internal/lab/pdr_exec.go

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

rm -rf results/backtest-first-btc-15m-previous-day-range-reversion-v1

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -backtest-first-btc-15m-previous-day-range-reversion-v1

wc -l results/backtest-first-btc-15m-previous-day-range-reversion-v1/*.csv

cat results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_falsification.json

cat results/backtest-first-btc-15m-previous-day-range-reversion-v1/btc_15m_previous_day_range_reversion_summary.csv

git diff --check
git status --short

After verification:
- Add a result review that records source facts, coverage rows, signal rows,
  executed trades, split metrics, long/short behavior, gross and net P&L,
  drawdown, pass/fail gates, and final stop state.
- Update memory/PROGRESS.md with the completed verification/result milestone.
- Update memory/NEXT_CODEX_BRIEF.md to the next gate.
- Update memory/DECISIONS.md only for durable decisions or constraints.

No-rescue boundaries:
- If this fixed baseline fails, do not rescue it with alternate UTC sessions,
  previous 2/3 day windows, changed outer-decile thresholds, derivatives context,
  calendar/session mining, replay, walk-forward, or optimizer grids.
- No paper/testnet/live path, exchange API work, credentials, deploy files,
  martingale, averaging down, two-exchange logic, or promotion.
```
