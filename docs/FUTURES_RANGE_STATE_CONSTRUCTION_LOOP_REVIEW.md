# Futures Range-State Construction Loop Audit Review

Date: 2026-06-27

## Verdict

Stop state:
`range_state_construction_loop_audit_passed_needs_router_spec`.

The non-trading BTCUSDT futures range-state construction loop audit was
implemented and run behind `-futures-range-state-construction-loop-audit`.

Source and closed UTC resample validation passed. The audit produced
closed-candle state rows that combine range geometry, volatility, trend,
impulse, and OHLCV participation proxy buckets. Forward outcomes were written
only as labels. The run found both stable `no_trade_toxic` filter cohorts and
some `tradable_rotation_candidate` cohorts, so the next allowed step is a
documentation-only router spec or router audit implementation. This result does
not authorize entries, exits, P&L backtests, optimizer grids, replay,
walk-forward, packaging, source expansion, symbol expansion, live-adjacent work,
or closed-family retuning.

## Source And Resampling

Source validation passed.

| Field | Value |
| --- | --- |
| Source | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
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

Closed UTC resampling also passed.

| Timeframe | Rows | First open | Last open | Status |
| --- | ---: | --- | --- | --- |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | `accepted` |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `accepted` |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `accepted` |

Every resample had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Artifacts

Result directory:

```text
results/futures-range-state-construction-loop-audit/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_range_state_construction_loop_sources.csv` | `2` |
| `futures_range_state_construction_loop_coverage.csv` | `4` |
| `futures_range_state_construction_loop_feature_windows.csv` | `4` |
| `futures_range_state_construction_loop_states.csv` | `29,785` |
| `futures_range_state_construction_loop_labels.csv` | `89,353` |
| `futures_range_state_construction_loop_cohorts.csv` | `68,797` |
| `futures_range_state_construction_loop_rankings.csv` | `16,336` |
| `futures_range_state_construction_loop_summary.csv` | `38` |
| `futures_range_state_construction_loop_skips.csv` | `13` |
| common `summary.csv` | `13` |

The common `source_manifest.json`, `summary.csv/json`, and `trades.json`
outputs remained zero-trade compatible. `trades.json` contains no trades.

## State And Label Inventory

The audit produced `29,784` eligible state rows and `89,352` forward-label rows.

| Timeframe | State rows |
| --- | ---: |
| `15m` | `24,067` |
| `1h` | `5,014` |
| `4h` | `703` |

Forward labels:

| Label | Rows |
| --- | ---: |
| `boundary_chop` | `25,911` |
| `contained_rotation` | `13,505` |
| `no_resolution` | `9,791` |
| `false_break_reentry_up` | `9,353` |
| `clean_expansion_up` | `9,165` |
| `false_break_reentry_down` | `8,386` |
| `clean_expansion_down` | `8,197` |
| `drift_through_up` | `2,838` |
| `drift_through_down` | `2,161` |
| `low_width_noise` | `45` |

Skip rows were aggregate counts only. All skipped rows were
`not_mature_active`: `18,214` on `15m`, `5,020` on `1h`, and `1,655` on `4h`.

## Cohort Rankings

The audit summarized `68,796` route cohorts and ranked `16,335` decision
cohorts. `58` cohorts passed:

| Route | Passing cohorts |
| --- | ---: |
| `no_trade_toxic` | `52` |
| `tradable_rotation_candidate` | `6` |
| `trend_continuation_candidate` | `0` |

Top ranked row:

```text
range_state_v1|15m|h48|no_trade_toxic|all_families|range_state_v1::15m::geometry_midline_balanced::vol_compressed::trend_down_pressure::impulse_none::participation_normal
```

It had `432` full-period rows, weakest split rows `134`, full toxic rate
`0.731481`, worst split toxic rate `0.768657`, and dominant label
`boundary_chop` at `0.729167`.

Top passing rotation row:

```text
range_state_v1|15m|h24|tradable_rotation_candidate|all_families|range_state_v1::15m::geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale::participation_low
```

It had `431` full-period rows, weakest split rows `92`, full useful rate
`0.617169`, weakest split useful rate `0.530612`, full toxic rate `0.169374`,
worst split toxic rate `0.239130`, and dominant label `contained_rotation` at
`0.364269`.

The first failed ranking row was a `tradable_rotation_candidate` using
`geometry+vol+trend+impulse`; it failed for `inadequate_cohort_count`.

## Review Gate Outcome

The audit passed as a state/router discovery milestone, not as a strategy
approval:

- source and resample validation passed;
- deterministic closed-candle feature buckets and state IDs were written;
- forward labels were not used as feature inputs;
- `no_trade_toxic` and `tradable_rotation_candidate` both had passing cohorts;
- mixed route evidence requires a router layer before any later strategy
  premise is specified.

Next allowed state:
`range_context_router_spec_ready_for_audit_implementation`.

The router must remain non-trading and zero-trade compatible. It may only map
closed-candle state rows into `no_trade`, `tradable_rotation`,
`trend_continuation`, or `diagnostic_only`. It must not choose trade direction
from forward labels, reopen `range_occupancy_rotation_v1`, structured
compression, breakout-retest/acceptance, clean breakout continuation,
hold-inside/midline, impulse absorption, or higher-timeframe nested range
rotation under new names.

## Verification

Commands run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-state-construction-loop-audit -out-dir results/futures-range-state-construction-loop-audit
wc -l results/futures-range-state-construction-loop-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

Observed audit summary:

```text
futures_range_state_construction_loop_audit source_rows=1 coverage_rows=3 feature_window_rows=3 state_rows=29784 label_rows=89352 cohort_rows=68796 ranking_rows=16335 passing_cohorts=58 stop_state=range_state_construction_loop_audit_passed_needs_router_spec
strategy=empty trades=0 output=results/futures-range-state-construction-loop-audit
```

`git diff --check` passed. The final pre-commit status contained only intended
code, docs, and memory changes for this milestone.
