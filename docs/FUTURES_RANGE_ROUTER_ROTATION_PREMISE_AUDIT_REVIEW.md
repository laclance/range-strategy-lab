# Futures Range Router Rotation Premise Audit Review

Date: 2026-06-27

## Verdict

Stop state:
`range_router_rotation_premise_audit_failed_no_premise`.

The zero-trade BTCUSDT Binance USDT-M futures router rotation premise audit was
implemented and run behind `-futures-range-router-rotation-premise-audit`.

The audit recomputed the reviewed range-state construction loop and context
router in-process, required router dependency stop state
`range_context_router_passed_needs_rotation_premise_spec`, selected only
`range_context_router_v1|15m|h24|tradable_rotation` context, collapsed
consecutive eligible router rows into context segments, froze range bounds at
segment start, searched the next `6` closed `15m` candles for boundary-reclaim
events, and labeled only `24`-bar forward outcomes after the event candle.

Source, `15m` resampling, and router dependency passed. The event premise did
not pass: it produced only `97` full-period valid events versus the required
`150`, only `23` events in the weakest split versus the required `40`, one
split contributed `47.4227%` of events versus the `45%` cap, and the adverse
rate behavior gates failed. This result does not authorize entries, exits, P&L
strategy backtests, optimizer grids, replay, walk-forward, packaging, source
expansion, symbol expansion, live-adjacent work, or closed-family retuning.

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

Closed UTC `15m` resampling also passed.

| Timeframe | Rows | First open | Last open | Status |
| --- | ---: | --- | --- | --- |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | `accepted` |

The `15m` resample had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and `coverage_facts_pass=true`.

## Router Dependency

The router dependency passed.

| Field | Value |
| --- | ---: |
| Required router stop state | `range_context_router_passed_needs_rotation_premise_spec` |
| Actual router stop state | `range_context_router_passed_needs_rotation_premise_spec` |
| Router rules | `58` |
| Router rows | `29,784` |
| Router cohorts | `84` |
| Router rankings | `12` |
| Router passing cohorts | `3` |
| Required cohort | `range_context_router_v1|15m|h24|tradable_rotation` |
| Required cohort passed | `true` |
| Required cohort full rows | `892` |
| Required cohort weakest split rows | `210` |
| Required cohort full expected hit rate | `0.599776` |
| Required cohort weakest split expected hit rate | `0.554667` |
| Required cohort full adverse rate | `0.181614` |
| Required cohort worst split adverse rate | `0.247619` |
| Required cohort dominant forward label | `contained_rotation` |
| Required cohort dominant forward label rate | `0.385650` |

No router row used forward labels as router input, and no forward-label columns
were present as premise inputs.

## Artifacts

Result directory:

```text
results/futures-range-router-rotation-premise-audit/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_range_router_rotation_premise_sources.csv` | `2` |
| `futures_range_router_rotation_premise_coverage.csv` | `2` |
| `futures_range_router_rotation_premise_router_dependency.csv` | `2` |
| `futures_range_router_rotation_premise_context_segments.csv` | `279` |
| `futures_range_router_rotation_premise_events.csv` | `98` |
| `futures_range_router_rotation_premise_outcomes.csv` | `98` |
| `futures_range_router_rotation_premise_cohorts.csv` | `13` |
| `futures_range_router_rotation_premise_rankings.csv` | `4` |
| `futures_range_router_rotation_premise_summary.csv` | `5` |
| `futures_range_router_rotation_premise_skips.csv` | `13` |
| common `summary.csv` | `13` |

The common `source_manifest.json`, `summary.csv/json`, and `trades.json`
outputs remained zero-trade compatible. `summary.csv` has `0` trades in every
split/side row.

## Premise Inventory

The audit found `278` context segments, `97` boundary-reclaim events, `97`
complete outcomes, and `0` missing-future outcomes.

| Split | Segments | Events | Lower events | Upper events | Midline outcomes | Hard adverse | Chop/no-resolution |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `full_2021_2026` | `278` | `97` | `43` | `54` | `71` | `22` | `4` |
| `2021_2022_stress` | `76` | `28` | `16` | `12` | `23` | `4` | `1` |
| `2023_2024_oos` | `119` | `46` | `20` | `26` | `33` | `11` | `2` |
| `2025_2026_recent` | `83` | `23` | `7` | `16` | `15` | `7` | `1` |

Skip inventory:

| Reason | Full count |
| --- | ---: |
| `ineligible_router_context` | `23,175` |
| `duplicate_router_context_row` | `614` |
| `no_boundary_reclaim_event` | `181` |

The `278` context segments were enough to pass the full context-segment count
gate, and both sides had at least `40` full-period events (`43` lower, `54`
upper). The premise failed because the event and behavior gates did not survive
the declared review thresholds.

## Rankings

The audit wrote `12` cohort rows and `3` ranking rows. No ranking passed.

Top ranking row:

```text
range_router_rotation_premise_v1|15m|h24|router_gated_boundary_reclaim_rotation|full_2021_2026|all
```

It had:

| Metric | Value |
| --- | ---: |
| Full context segments | `278` |
| Full valid events | `97` |
| Weakest split events | `23` |
| Lower side full events | `43` |
| Upper side full events | `54` |
| Max split contribution | `0.474227` |
| Full midline rotation rate | `0.731959` |
| Weakest split midline rotation rate | `0.652174` |
| Full hard adverse rate | `0.226804` |
| Worst split hard adverse rate | `0.304348` |
| Full chop/no-resolution rate | `0.041237` |
| Dominant outcome | `midline_rotation_first` at `0.371134` |
| Dominant state ID rate | `0.319588` |

Failure reason:

```text
inadequate_event_count,inadequate_split_event_count,single_split_contribution_above_gate,behavior_gate_failed
```

The full and weakest-split midline rates cleared the improvement checks versus
the router baseline, and chop/no-resolution stayed below the cap. The premise
still failed because full valid events were below `150`, the weakest split had
only `23` events, max split contribution was above `0.45`, full hard adverse
rate was above `0.22`, and worst-split hard adverse rate was above `0.28`.

## Review Gate Outcome

The audit failed as a premise gate:

- source validation passed;
- closed UTC `15m` resampling passed;
- router dependency passed;
- context construction stayed closed-candle and used frozen segment-start range
  bounds;
- event selection did not use forward labels;
- outcome labels started strictly after the event candle;
- common outputs remained zero-trade;
- no cohort passed the declared event premise gates.

Current stop state:
`range_router_rotation_premise_audit_failed_no_premise`.

This closes `router_gated_boundary_reclaim_rotation_v1` in its reviewed form.
Do not convert the `97` events, `278` segments, or `1,299` rotation-routed
router rows into trades. Any later work must be a materially different
non-trading premise or context audit with its own review gate.

## Verification

Commands run:

```bash
gofmt -w internal/lab/futures_range_context_router_audit.go internal/lab/futures_range_router_rotation_premise_audit.go internal/lab/futures_range_router_rotation_premise_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-router-rotation-premise-audit -out-dir results/futures-range-router-rotation-premise-audit
wc -l results/futures-range-router-rotation-premise-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

Observed audit summary:

```text
futures_range_router_rotation_premise_audit source_rows=1 coverage_rows=1 router_dependency_rows=1 context_segments=278 events=97 outcomes=97 cohort_rows=12 ranking_rows=3 passing_cohorts=0 stop_state=range_router_rotation_premise_audit_failed_no_premise
strategy=empty trades=0 output=results/futures-range-router-rotation-premise-audit
```

`go test ./...` passed. `wc -l` over premise-audit CSV artifacts plus common
`summary.csv` totaled `529` lines. The brief-reference scan found canonical
`memory/NEXT_CODEX_BRIEF.md` references and checklist mentions only.
`git diff --check` passed. Pre-commit status showed only intended code, docs,
and memory changes.
