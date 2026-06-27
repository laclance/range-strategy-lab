# Futures Range-First Occupancy Rotation V1 Optimization Review

Date: 2026-06-27

## Verdict

Stop state:
`range_first_strategy_v1_optimizer_failed_no_replay`.

The bounded offline BTCUSDT optimizer/backtester for
`range_occupancy_rotation_v1` was implemented and run behind
`-futures-range-first-occupancy-rotation-v1-optimization`.

Source and closed UTC resample validation passed, the declared `1,152` config
grid was evaluated, and the fixed baseline common outputs were written. No grid
config passed the predeclared train, OOS, recent, full-period, PF,
net-to-drawdown, and side-caveat gates. The fixed replay spec is not
authorized.

This is no-promotion evidence for the reviewed V1 occupancy-rotation grammar.
It does not authorize retuning the grid, relaxing gates, changing the target,
stop, max hold, timeframe, symbol scope, or renaming this grammar into another
implementation pass.

## Source And Coverage

Accepted source:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | `BTCUSDT` |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Closed UTC resamples:

| Timeframe | Rows | First open | Last open | Status |
| --- | ---: | --- | --- | --- |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | accepted |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | accepted |

Both resamples had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, and complete closed buckets.

## Artifacts

Result directory:

```text
results/futures-range-first-occupancy-rotation-v1-optimization/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_range_first_occupancy_rotation_v1_sources.csv` | `2` |
| `futures_range_first_occupancy_rotation_v1_coverage.csv` | `3` |
| `futures_range_first_occupancy_rotation_v1_grid.csv` | `1,153` |
| `futures_range_first_occupancy_rotation_v1_baseline.csv` | `2` |
| `futures_range_first_occupancy_rotation_v1_signals.csv` | `84` |
| `futures_range_first_occupancy_rotation_v1_trades.csv` | `44` |
| `futures_range_first_occupancy_rotation_v1_summary.csv` | `13,825` |
| `futures_range_first_occupancy_rotation_v1_rankings.csv` | `1,153` |
| `futures_range_first_occupancy_rotation_v1_selection.csv` | `3` |
| `futures_range_first_occupancy_rotation_v1_skips.csv` | `42,543` |
| common `summary.csv` | `13` |

Common outputs stayed fixed-baseline only:

- `source_manifest.json`;
- `summary.csv`;
- `summary.json`;
- `trades.json`.

Grid, ranking, selection, skip, signal, and V1 trade details are confined to
the V1-specific artifacts.

## Fixed Baseline

Fixed baseline:

```text
range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005
```

Baseline result:

| Split | Trades | Net P&L | Profit Factor |
| --- | ---: | ---: | ---: |
| `2021_2022_stress` | `8` | `-22.653070` | `0.464536` |
| `2023_2024_oos` | `16` | `-38.248119` | `0.428362` |
| `2025_2026_recent` | `19` | `-33.307557` | `0.504958` |
| `full_2021_2026` | `43` | `-94.208745` | `0.466232` |

Baseline rank was `769` of `1,152`. It failed the minimum train, OOS, recent,
and full trade-count gates; all net P&L gates; all PF gates; and both
train/full net-to-drawdown gates.

## Optimization Result

Declared grid rows: `1,152`.

Passing grid rows: `0`.

Top-ranked non-passing row:

```text
range_occupancy_rotation_v1_1h_l24_w020_ow8_occ070_rec25_t66_h12_sb000
```

Its rank score was high only because the train split had `2` trades, positive
train net P&L, and no train drawdown. It still failed selection because it had
too few train, OOS, recent, and full trades; negative OOS and full net P&L; OOS
and full PF below gate; full net-to-drawdown below gate; and side weakness.

Selection artifact outcome:

| Role | Config | Rank | Passes Gate | Failure |
| --- | --- | ---: | --- | --- |
| baseline | `range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005` | `769` | `false` | baseline failed gates |
| selected | none | `0` | `false` | `no_selectable_config` |

The stop-state precedence was applied as specified: a passing optimized config
would have advanced the run even if the fixed baseline failed, but no optimized
config passed.

## Review Gate Outcome

The V1 occupancy-rotation grammar failed the optimizer review:

- source validation passed;
- `15m` and `1h` closed UTC resample validation passed;
- fixed baseline failed across every split after costs;
- no declared grid config passed all selection gates;
- no selected config exists for fixed replay.

Next allowed state is review-only unless the user supplies a materially
different offline range-first premise and a new spec. Do not retune this V1
grid, relax gates, rerun the failed closed families, expand symbols, or move to
paper/testnet/live/deploy from this result.

## Verification

Commands run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-first-occupancy-rotation-v1-optimization -out-dir results/futures-range-first-occupancy-rotation-v1-optimization
wc -l results/futures-range-first-occupancy-rotation-v1-optimization/*.csv
```

Additional closeout checks are recorded in `memory/PROGRESS.md`.
