# Backtest-First BTC 15m Previous-Day Range Reversion Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_15m_previous_day_range_reversion_backtest_implementation_added_needs_local_verification
```

This implementation adds the selected fixed offline baseline candidate:

```text
btc_15m_previous_day_range_reversion_v1
```

The implementation is not promoted and no result verdict is claimed here because
this connector session could not run the local Go test suite or the BTCUSDT CSV
backtest command.

## What Was Added

- Lab implementation:
  - `internal/lab/backtest_first_btc_15m_day_types.go`
  - `internal/lab/backtest_first_btc_15m_day_rows.go`
  - `internal/lab/backtest_first_btc_15m_day_runner.go`
  - `internal/lab/backtest_first_btc_15m_day_support.go`
  - `internal/lab/day_state.go`
  - `internal/lab/pdr.go`
  - `internal/lab/pdr_row.go`
  - `internal/lab/pdrh.go`
  - `internal/lab/pdr_inside.go`
  - `internal/lab/pdr_exec.go`
- Offline CLI entrypoint:
  - `cmd/rangelab/day.go`
  - `cmd/rangelab/day_outputs.go`
- Fixed flag:

```text
-backtest-first-btc-15m-previous-day-range-reversion-v1
```

- Default output path:

```text
results/backtest-first-btc-15m-previous-day-range-reversion-v1/
```

## Fixed Baseline Scope

The implementation follows the selected packet:

- current accepted BTCUSDT Binance USDT-M futures `5m` CSV;
- exact closed UTC `15m` resample from complete three-child `5m` buckets;
- prior complete UTC day's high, low, and midpoint;
- skip the current day after any prior current-day close outside the previous
  day's high-low range;
- long entry when a closed `15m` candle is inside the previous-day range and in
  its lower `10%`;
- short entry when a closed `15m` candle is inside the previous-day range and in
  its upper `10%`;
- next-`15m`-open execution through the existing backtest engine;
- stop beyond the prior-day range extreme by `0.25 * ATR(14)[d-1]`;
- target at the prior-day midpoint;
- `24` closed `15m` bar time stop;
- `1%` risk at stop, `1x` notional cap, `0.0004` fee per side, and `0.000116`
  slippage per side.

No optimizer, replay, walk-forward, derivatives-veto interaction, source
expansion, paper/testnet/live path, exchange API work, credentials, deploy file,
martingale, averaging down, two-exchange logic, or promotion was added.

## Expected Artifacts

The fixed run should write:

- `source_manifest.json`
- `summary.json`
- `summary.csv`
- `trades.json`
- `btc_15m_previous_day_range_reversion_sources.json/csv`
- `btc_15m_previous_day_range_reversion_coverage.json/csv`
- `btc_15m_previous_day_range_reversion_signals.json/csv`
- `btc_15m_previous_day_range_reversion_skips.json/csv`
- `btc_15m_previous_day_range_reversion_trades.json/csv`
- `btc_15m_previous_day_range_reversion_summary.json/csv`
- `btc_15m_previous_day_range_reversion_falsification.json/csv`

## Required Local Verification

Run from repo root:

```bash
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
```

## Next Gate

After local verification, record the result with source facts, coverage rows,
signal rows, executed trades, split metrics, long/short behavior, gross and net
P&L, drawdown, pass/fail gates, and final stop state.

If this fixed baseline fails, do not rescue it with alternate UTC sessions,
previous `2`/`3` day windows, changed outer-decile thresholds, derivatives
context, calendar/session mining, replay, walk-forward, or optimizer grids.
