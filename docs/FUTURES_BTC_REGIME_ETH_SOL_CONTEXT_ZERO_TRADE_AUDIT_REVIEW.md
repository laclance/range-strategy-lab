# Futures BTC Regime Plus ETH/SOL Zero-Trade Audit Review

Date: 2026-06-28

## Verdict

Stop state:
`btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.

The explicitly approved zero-trade BTC regime plus ETH/SOL context audit was
implemented and run behind
`-futures-btc-regime-eth-sol-context-audit`.

The audit used BTCUSDT only as closed-candle market-regime context and
diagnostic authority. ETHUSDT and SOLUSDT were evaluated only as possible
zero-trade context authority rows. The audit produced no passing cohorts, wrote
no trade rows, and does not authorize entries, exits, P&L strategy backtests,
optimizer grids, replay, walk-forward, source downloads, strategy promotion,
paper/testnet/live paths, exchange API work, credentials, deploy files,
martingale, averaging down, or two-exchange logic.

The path is closed in this reviewed zero-trade form. Do not retune, reslice,
rename, gate-relax, or promote this BTC regime plus ETH/SOL context surface
from the current result.

## Source Scope

Only the approved local Binance USDT-M futures `5m` files were used.

| Symbol | Role | Path | Rows | First open | Last open | Gaps | Duplicates | Zero volume | Physical non-monotonic | Sorted for validation | Status |
| --- | --- | --- | ---: | --- | --- | ---: | ---: | ---: | ---: | --- | --- |
| `BTCUSDT` | BTC market-regime context, diagnostic only | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `0` | `66` | `0` | `false` | `accepted` |
| `ETHUSDT` | local range context authority candidate | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `0` | `47` | `0` | `false` | `accepted` |
| `SOLUSDT` | local range context authority candidate | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `0` | `47` | `1` | `true` | `accepted` |

No spot data, new symbols, source downloads, derivatives context sources,
private endpoints, exchange APIs, or broad mining were used.

## Coverage And Context Rows

Closed UTC resampling passed for all symbols and timeframes. Each resampled
stream had `gap_count=0`, `duplicate_count=0`, `missing_child_open_count=0`,
`complete=true`, and `validation_status=accepted`.

| Symbol | Timeframe | Resampled rows | State rows | Label rows | Context-matched rows | Missing BTC context rows |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `BTCUSDT` | `15m` | `191,328` | `24,067` | `0` | `24,067` | `0` |
| `BTCUSDT` | `1h` | `47,832` | `5,014` | `0` | `5,014` | `0` |
| `BTCUSDT` | `4h` | `11,958` | `703` | `0` | `703` | `0` |
| `ETHUSDT` | `15m` | `191,328` | `21,875` | `65,625` | `14,199` | `7,676` |
| `ETHUSDT` | `1h` | `47,832` | `4,849` | `14,547` | `3,051` | `1,798` |
| `ETHUSDT` | `4h` | `11,958` | `820` | `2,460` | `307` | `513` |
| `SOLUSDT` | `15m` | `191,328` | `20,862` | `62,586` | `10,579` | `10,283` |
| `SOLUSDT` | `1h` | `47,832` | `4,605` | `13,815` | `2,353` | `2,252` |
| `SOLUSDT` | `4h` | `11,958` | `845` | `2,535` | `228` | `617` |

Total audit rows:

| Artifact family | Rows |
| --- | ---: |
| Source rows | `3` |
| Coverage rows | `9` |
| BTC state rows | `29,784` |
| ETH/SOL local state rows | `30,717` |
| Relative-strength rows | `30,717` |
| Label rows | `92,151` |
| Cohort rows | `359,055` |
| Ranking rows | `160,983` |
| Passing cohorts | `0` |

## Anti-Leakage And Zero-Trade Status

The audit kept the time split explicit:

- state IDs used only data available at the closed decision candle;
- relative-strength buckets used closed-candle returns only;
- forward labels appeared only in label, cohort, ranking, and summary
  artifacts;
- no forward label was used as a source, state-ID, context, router, or gating
  input;
- common outputs remained zero-trade compatible.

The full-period summary row recorded:

| Field | Value |
| --- | --- |
| `source_scope_pass` | `true` |
| `coverage_pass` | `true` |
| `common_outputs_zero_trade` | `true` |
| `btc_context_diagnostic_only` | `true` |
| `eth_sol_authority_candidate_only` | `true` |
| `forward_labels_as_inputs` | `false` |

Common `summary.csv` has `0` trades in every split/side row, and `trades.json`
contains no trade rows.

## Ranking Outcome

The audit ranked `160,983` BTC-regime-plus-local context rows against their
local-only baselines. No row passed the declared context gates.

Failure reason inventory:

| Failure reason | Ranking rows |
| --- | ---: |
| `btc_context_improvement_gate_failed` | `142,615` |
| `single_split_contribution_above_gate` | `16,095` |
| `route_rate_gate_failed` | `1,627` |
| `missing_period_split` | `419` |
| `inadequate_cohort_count` | `227` |

The top-ranked rows showed high apparent separation only in tiny cohorts. The
highest-ranked row was a SOLUSDT `15m` `h48` tradable-rotation candidate with
`6` full-period rows and `1` weakest-split row, so it failed
`inadequate_cohort_count`. The dominant failure mode across the full ranking set
was that BTC regime context did not add durable separation beyond the ETH/SOL
local-only baseline.

## Artifacts

Result directory:

```text
results/futures-btc-regime-eth-sol-context-audit/
```

CSV line counts, including headers:

| Artifact | Lines |
| --- | ---: |
| `futures_btc_regime_eth_sol_context_sources.csv` | `4` |
| `futures_btc_regime_eth_sol_context_coverage.csv` | `10` |
| `futures_btc_regime_eth_sol_context_btc_states.csv` | `29,785` |
| `futures_btc_regime_eth_sol_context_local_states.csv` | `30,718` |
| `futures_btc_regime_eth_sol_context_relative_strength.csv` | `30,718` |
| `futures_btc_regime_eth_sol_context_labels.csv` | `92,152` |
| `futures_btc_regime_eth_sol_context_cohorts.csv` | `359,056` |
| `futures_btc_regime_eth_sol_context_rankings.csv` | `160,984` |
| `futures_btc_regime_eth_sol_context_summary.csv` | `74` |
| common `summary.csv` | `13` |

The audit also wrote JSON counterparts for the audit-specific artifacts plus
the common `source_manifest.json`, `summary.json`, and `trades.json` outputs.

## Review Gate Outcome

The implementation satisfied the approved zero-trade audit boundary:

- approved local BTC/ETH/SOL Binance USDT-M futures `5m` sources only;
- BTCUSDT remained diagnostic market-regime context only;
- ETHUSDT and SOLUSDT remained possible context authority rows only;
- source validation and closed UTC resampling passed;
- forward labels stayed out of premise and state inputs;
- common outputs stayed zero-trade compatible;
- no entry, exit, P&L strategy backtest, optimizer grid, replay, walk-forward,
  source download, strategy package, promotion, or live-adjacent path was added;
- no cohort passed the declared context gates.

Current stop state:
`btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.

This result closes the BTC regime plus ETH/SOL context surface in reviewed
zero-trade form. Future work must choose a materially different approved scope;
it must not retune this audit, convert its labels into entries, treat ETH/SOL
context rows as strategy authority, or promote BTC regime rows.

## Verification

Commands run:

```bash
gofmt -w internal/lab/futures_btc_regime_eth_sol_context_audit.go internal/lab/futures_btc_regime_eth_sol_context_audit_test.go cmd/rangelab/main.go cmd/rangelab/main_test.go
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-regime-eth-sol-context-audit -out-dir results/futures-btc-regime-eth-sol-context-audit
wc -l results/futures-btc-regime-eth-sol-context-audit/*.csv
```

Observed audit summary:

```text
futures_btc_regime_eth_sol_context_audit source_rows=3 coverage_rows=9 btc_state_rows=29784 local_state_rows=30717 relative_strength_rows=30717 label_rows=92151 cohort_rows=359055 ranking_rows=160983 passing_cohorts=0 stop_state=btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context
loaded 573984 candles from 2021-01-01T00:00:00Z to 2026-06-16T23:59:59Z
strategy=empty trades=0 output=results/futures-btc-regime-eth-sol-context-audit
```

`go test ./...` passed. `wc -l` over audit CSV artifacts plus common
`summary.csv` totaled `703,514` lines.

Final closeout also requires:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
