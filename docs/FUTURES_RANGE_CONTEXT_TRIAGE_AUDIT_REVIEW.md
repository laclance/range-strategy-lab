# Futures Range Context Triage Audit Review

Date: 2026-06-27

## Scope

This review covers the non-trading BTCUSDT futures range-context triage audit
run behind `-futures-range-context-triage-audit`.

The audit evaluated range quality, UTC session behavior, and failure-mode
taxonomy in parallel from the accepted Binance USDT-M futures `5m` BTCUSDT
source resampled to closed UTC `15m`, `1h`, and `4h` bars. It did not add
entries, exits, P&L, optimizer grids, replay, walk-forward logic, source
expansion, symbol expansion, or live-adjacent behavior.

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

## Artifacts

Result directory:

```text
results/futures-range-context-triage-audit/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_range_context_triage_sources.csv` | `2` |
| `futures_range_context_triage_coverage.csv` | `4` |
| `futures_range_context_triage_episodes.csv` | `3,590` |
| `futures_range_context_triage_quality.csv` | `3,590` |
| `futures_range_context_triage_sessions.csv` | `533` |
| `futures_range_context_triage_failure_modes.csv` | `4,261` |
| `futures_range_context_triage_cohorts.csv` | `883` |
| `futures_range_context_triage_rankings.csv` | `160` |
| `futures_range_context_triage_summary.csv` | `38` |
| common `summary.csv` | `13` |

The common `source_manifest.json`, `summary.csv/json`, and `trades.json`
outputs remained zero-trade compatible. `trades.json` contains no trades.

## Episode Inventory

The audit produced `3,589` episode rows. Of those, `1,420` reached a mature
active range and were eligible for context labeling.

| Timeframe | Eligible episodes | Skipped `no_mature_active_bar` |
| --- | ---: | ---: |
| `15m` | `1,070` | `1,499` |
| `1h` | `282` | `465` |
| `4h` | `68` | `205` |

Eligible quality buckets were concentrated in `balanced_orderly` on `15m` and
`1h`, while the small `4h` set was mostly `wide_volatile`.

| Timeframe | Quality bucket | Eligible episodes |
| --- | --- | ---: |
| `15m` | `balanced_orderly` | `955` |
| `15m` | `wide_volatile` | `104` |
| `15m` | `narrow_orderly` | `7` |
| `15m` | `too_narrow_noise` | `4` |
| `1h` | `balanced_orderly` | `162` |
| `1h` | `wide_volatile` | `119` |
| `1h` | `narrow_orderly` | `1` |
| `4h` | `wide_volatile` | `66` |
| `4h` | `balanced_orderly` | `2` |

## Failure-Mode Taxonomy

The audit wrote `4,260` failure-mode rows across the declared horizons
`12`, `24`, and `48`. The most common primary label was `boundary_chop`, and
the next largest groups were clean expansions and false-break reentries.

| Primary label | Rows |
| --- | ---: |
| `boundary_chop` | `1,706` |
| `clean_expansion_up` | `556` |
| `clean_expansion_down` | `505` |
| `false_break_reentry_down` | `421` |
| `false_break_reentry_up` | `415` |
| `contained_rotation` | `317` |
| `no_resolution` | `174` |
| `drift_through_up` | `84` |
| `drift_through_down` | `70` |
| `low_width_noise` | `12` |

The `48` bar `15m` horizon was especially chop-heavy, with `704`
`boundary_chop` labels. Shorter `1h` rows were less noisy, but still failed
the usable/toxic gates.

## Cohort Rankings

The audit summarized `882` cohorts and ranked `159` decision-context cohorts.
No cohort passed the review gates.

The top ranked row was:

```text
range_context_1h_h12_quality_bucket_mature_session_balanced_orderly_us_late_utc_17_23
```

It had `66` full-period candidates, full usable context rate `0.560606`,
weakest split usable rate `0.391304`, full toxic context rate `0.439394`, and
worst split toxic rate `0.608696`. It failed for inadequate full count,
inadequate split count, usable rate below gate, and toxic rate above gate.

The broad `1h h12 all` cohort ranked ninth. It had `282` full-period
candidates, full usable context rate `0.478723`, weakest split usable rate
`0.376623`, full toxic context rate `0.521277`, and worst split toxic rate
`0.623377`. It also failed count, usable-rate, and toxic-rate gates.

Outcome-label cohorts are useful taxonomy rows, but they are not rankable
decision context because they depend on future-known labels. They are therefore
excluded from strategy-premise selection.

## Verdict

Final stop state:

```text
range_context_triage_failed_no_strategy_premise
```

The source and closed UTC resampling passed, so this is not a source-gap stop.
The audit failed because no range-quality, session, or quality-plus-session
context cohort passed the declared gates with adequate count, split stability,
usable/toxic balance, and closed-family protection.

This result does not authorize a strategy spec, baseline backtest, optimizer,
fixed replay, walk-forward review, strategy package, retune, gate relaxation,
source expansion, symbol expansion, or live-adjacent path. Future strategy work
requires a materially different user-approved offline range-first premise.

## Verification

Commands run:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-range-context-triage-audit -out-dir results/futures-range-context-triage-audit
wc -l results/futures-range-context-triage-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
