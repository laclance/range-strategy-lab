# Futures Range Context Router Audit Review

Date: 2026-06-27

## Verdict

Stop state:
`range_context_router_passed_needs_rotation_premise_spec`.

The non-trading BTCUSDT futures range context router audit was implemented and
run behind `-futures-range-context-router-audit`.

The router reused the accepted futures range-state construction loop. It turned
only passing state-audit ranking rows into deterministic closed-candle router
rules, assigned exactly one router label per eligible state row, and kept
forward labels out of router inputs. Forward labels were used only after router
assignment to audit cohorts.

This result authorizes only a later materially new rotation premise spec. It
does not authorize entries, exits, P&L backtests, optimizer grids, replay,
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

| Timeframe | Resample rows | Router rows | First open | Last open | Status |
| --- | ---: | ---: | --- | --- | --- |
| `15m` | `191,328` | `24,067` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | `accepted` |
| `1h` | `47,832` | `5,014` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `accepted` |
| `4h` | `11,958` | `703` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `accepted` |

Every resample had `gap_count=0`, `duplicate_count=0`,
`missing_child_open_count=0`, `complete=true`, and
`validation_status=accepted`.

## Artifacts

Result directory:

```text
results/futures-range-context-router-audit/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_range_context_router_sources.csv` | `2` |
| `futures_range_context_router_coverage.csv` | `4` |
| `futures_range_context_router_rules.csv` | `59` |
| `futures_range_context_router_rows.csv` | `29,785` |
| `futures_range_context_router_cohorts.csv` | `85` |
| `futures_range_context_router_rankings.csv` | `13` |
| `futures_range_context_router_summary.csv` | `38` |
| `futures_range_context_router_skips.csv` | `25` |
| common `summary.csv` | `13` |

The common `source_manifest.json`, `summary.csv/json`, and `trades.json`
outputs remained zero-trade compatible. `trades.json` contains no trades.

## Router Inventory

The router produced `58` active rules from state-audit passing rankings:

| Router label | Rules |
| --- | ---: |
| `no_trade` | `52` |
| `tradable_rotation` | `6` |
| `trend_continuation` | `0` |

It assigned `29,784` eligible state rows:

| Router label | Rows |
| --- | ---: |
| `no_trade` | `13,546` |
| `tradable_rotation` | `1,299` |
| `trend_continuation` | `0` |
| `diagnostic_only` | `14,939` |

No row had conflicting router-rule matches. All `14,939` diagnostic-only rows
had `no_passing_rule_match`. State-audit skip preservation added `24,889`
`state_audit_not_mature_active` rows to the skip summary.

## Cohort Rankings

The router summarized `84` cohorts and ranked `12` full-period router cohorts.
`3` cohorts passed:

| Router label | Passing cohorts |
| --- | ---: |
| `tradable_rotation` | `1` |
| `no_trade` | `2` |
| `trend_continuation` | `0` |

Top ranked row:

```text
range_context_router_v1|15m|h24|tradable_rotation
```

It had `892` full-period rows, weakest split rows `210`, full expected-route hit
rate `0.599776`, weakest split expected hit rate `0.554667`, full adverse rate
`0.181614`, worst split adverse rate `0.247619`, and dominant forward label
`contained_rotation` at `0.385650`.

Passing no-trade rows:

- `range_context_router_v1|15m|h48|no_trade`: full toxic hit rate `0.623421`,
  weakest split toxic hit rate `0.606882`, dominant label `boundary_chop` at
  `0.613265`.
- `range_context_router_v1|1h|h48|no_trade`: full toxic hit rate `0.612686`,
  weakest split toxic hit rate `0.538636`, dominant label `boundary_chop` at
  `0.610886`.

The first failed ranking row was `15m` `h12` `tradable_rotation`; it failed
`route_rate_gate_failed` with full expected-route hit rate `0.464126`, below
the declared positive-route full-rate gate.

## Review Gate Outcome

The router passed as a route-selection milestone, not as a strategy approval:

- source and resample validation passed;
- state-audit stop state was
  `range_state_construction_loop_audit_passed_needs_router_spec`;
- router rules were created only from the `58` passing state-audit rankings;
- router rows contained no forward-label input columns;
- no `trend_continuation` route survived;
- both no-trade filtering and one rotation route survived router review.

Next allowed state:
`range_rotation_premise_spec_ready_for_authoring`.

The next document may specify only a materially new rotation premise. It must
explain why it is not `range_occupancy_rotation_v1`, hold-inside/midline,
breakout-retest/acceptance, clean breakout continuation, structured
compression, impulse absorption, or higher-timeframe nested range rotation under
a new name.

## Verification

Commands run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-router-audit -out-dir results/futures-range-context-router-audit
wc -l results/futures-range-context-router-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

Observed audit summary:

```text
futures_range_context_router_audit source_rows=1 coverage_rows=3 rule_rows=58 router_rows=29784 cohort_rows=84 ranking_rows=12 passing_cohorts=3 stop_state=range_context_router_passed_needs_rotation_premise_spec
strategy=empty trades=0 output=results/futures-range-context-router-audit
```

`wc -l` over router CSV artifacts totaled `30,024` lines. The brief-reference
scan found the canonical `memory/NEXT_CODEX_BRIEF.md` references and checklist
mentions, with no duplicate root `CODEX_BRIEF.md`. `git diff --check` passed.
The pre-commit status contained only intended code, docs, and memory changes
for this milestone.
