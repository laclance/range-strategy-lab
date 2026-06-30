# Backtest-First BTC 5m Rolling Value-Area Reversion Implementation Review

Date: 2026-06-30

## Verdict

Stop state:

```text
btc_5m_rolling_value_area_reversion_backtest_implementation_added_needs_local_verification
```

This implementation adds the selected fixed offline baseline candidate:

```text
btc_5m_rolling_value_area_reversion_v1
```

The implementation is not promoted and no result verdict is claimed here because
this connector session could not run the local Go test suite or the BTCUSDT CSV
backtest command.

## What Was Added

- Lab implementation:
  - `internal/lab/backtest_first_btc_5m_value_area_types.go`
  - `internal/lab/backtest_first_btc_5m_value_area_runner.go`
  - `internal/lab/backtest_first_btc_5m_value_area_support.go`
- Offline CLI entrypoint:
  - `cmd/rangelab/backtest_first_value_area_reversion.go`
- Fixed flag:

```text
-backtest-first-btc-5m-rolling-value-area-reversion-v1
```

- Default output path:

```text
results/backtest-first-btc-5m-rolling-value-area-reversion-v1/
```

## Fixed Baseline Scope

The implementation follows the selected packet:

- native closed BTCUSDT `5m` candles;
- prior `288` closed `5m` bars as the rolling value-area window;
- rolling VWAP from typical price weighted by volume;
- `ATR(14)[d-1]` known at decision time;
- minimum range width of `6 * ATR(14)[d-1]`;
- lower/upper `20%` outer-zone entries;
- `0.15 * range_width` minimum distance from VWAP;
- next-`5m`-open execution through the existing backtest engine;
- `36` closed `5m` bar time stop;
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
- `btc_5m_rolling_value_area_reversion_sources.json/csv`
- `btc_5m_rolling_value_area_reversion_signals.json/csv`
- `btc_5m_rolling_value_area_reversion_skips.json/csv`
- `btc_5m_rolling_value_area_reversion_trades.json/csv`
- `btc_5m_rolling_value_area_reversion_summary.json/csv`
- `btc_5m_rolling_value_area_reversion_falsification.json/csv`

## Required Local Verification

Run from repo root:

```bash
/usr/local/go/bin/gofmt -w \
  cmd/rangelab/backtest_first_value_area_reversion.go \
  internal/lab/backtest_first_btc_5m_value_area_types.go \
  internal/lab/backtest_first_btc_5m_value_area_runner.go \
  internal/lab/backtest_first_btc_5m_value_area_support.go

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...

env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab \
  -backtest-first-btc-5m-rolling-value-area-reversion-v1

wc -l results/backtest-first-btc-5m-rolling-value-area-reversion-v1/*.csv

git diff --check
git status --short
```

## Next Gate

After local verification, create a result review recording:

- source manifest facts;
- signal rows;
- executed trades;
- split metrics;
- long/short behavior;
- gross and net P&L;
- drawdown;
- pass/fail gate outcomes;
- stop state.

If the fixed baseline fails, do not rescue it with alternate VWAP windows,
outer-zone percentages, target changes, time-stop changes, side selection,
volume filters, derivatives-veto interaction, replay, walk-forward, or optimizer
grids.
