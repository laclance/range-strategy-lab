# Backtest-First BTC 15m Range-Edge Exhaustion Fade Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_15m_range_edge_exhaustion_fade_backtest_implementation_added_needs_local_verification
```

This implementation adds the fixed offline baseline candidate:

```text
btc_15m_range_edge_exhaustion_fade_v1
```

The implementation is not promoted and no result verdict is claimed here because
this connector session could not run the local Go test suite or BTCUSDT CSV
backtest command.

## Fixed Baseline Scope

The implementation follows the selected packet:

- current accepted BTCUSDT Binance USDT-M futures `5m` CSV;
- exact closed UTC `15m` resample from complete three-child `5m` buckets;
- prior `96` closed `15m` bars `[d-96,d-1]` as the local range;
- require `close[d-3]` to start between the lower and upper `40%` of the range;
- long fade: last three closes move downward into the lower `15%` of the range,
  `close[d]` remains above range low, and final downward progress is less than
  `0.35 * ATR(14)[d-1]`;
- short fade: symmetric upward move into the upper `15%` of the range, still
  below range high, with weak final progress;
- next-`15m`-open execution through the existing backtest engine;
- stop beyond the local range edge by `0.25 * ATR(14)[d-1]`;
- target at the local range midpoint;
- `16` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side.

No optimizer, replay, walk-forward, derivatives-veto interaction, source
expansion, paper/testnet/live path, exchange API work, credentials, deploy file,
martingale, averaging down, two-exchange logic, or promotion was added.

## Fixed Flag And Output Path

```text
-backtest-first-btc-15m-range-edge-exhaustion-fade-v1
```

```text
results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/
```

## Expected Artifacts

The fixed run should write:

- `source_manifest.json`
- `summary.json`
- `summary.csv`
- `trades.json`
- `btc_15m_range_edge_exhaustion_fade_sources.json/csv`
- `btc_15m_range_edge_exhaustion_fade_coverage.json/csv`
- `btc_15m_range_edge_exhaustion_fade_signals.json/csv`
- `btc_15m_range_edge_exhaustion_fade_skips.json/csv`
- `btc_15m_range_edge_exhaustion_fade_trades.json/csv`
- `btc_15m_range_edge_exhaustion_fade_summary.json/csv`
- `btc_15m_range_edge_exhaustion_fade_falsification.json/csv`

## Required Local Verification

Run from repo root:

```bash
/usr/local/go/bin/gofmt -w \
  cmd/rangelab/exfade.go \
  cmd/rangelab/exfade_outputs.go \
  cmd/rangelab/exfade_artifacts.go \
  internal/lab/edge_types.go \
  internal/lab/edge_rows.go \
  internal/lab/edge_defaults.go \
  internal/lab/edge_source.go \
  internal/lab/edge_resample.go \
  internal/lab/edge_runner.go \
  internal/lab/edge_state.go \
  internal/lab/edge_range.go \
  internal/lab/edge_side.go \
  internal/lab/edge_signal.go \
  internal/lab/edge_exec.go \
  internal/lab/edge_split.go \
  internal/lab/edge_trade_rows.go \
  internal/lab/edge_eval.go \
  internal/lab/edge_falsification.go

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

rm -rf results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -backtest-first-btc-15m-range-edge-exhaustion-fade-v1

wc -l results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/*.csv

cat results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_falsification.json

cat results/backtest-first-btc-15m-range-edge-exhaustion-fade-v1/btc_15m_range_edge_exhaustion_fade_summary.csv

git diff --check
git status --short
```

## Next Gate

After local verification, record the result with source facts, coverage rows,
signal rows, executed trades, split metrics, long/short behavior, gross and net
P&L, drawdown, pass/fail gates, and final stop state.

If this fixed baseline fails, do not rescue it with alternate range windows,
progress thresholds, edge zones, midpoint variants, added volume filters,
derivatives context, replay, walk-forward, or optimizer grids.
