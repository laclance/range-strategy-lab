# Futures BTC Regime Plus ETH/SOL Zero-Trade Audit Brief

Date: 2026-06-28

## Verdict

Stop state:
`btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.

This documentation-only brief converts the approved scope in
`docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md` into a
decision-complete plan for a future zero-trade context audit.

This brief does not authorize implementation. The next step after this brief is
user approval or rejection of the audit implementation scope. No Go code, CLI
flag, generated result directory, source download, entry, exit, P&L strategy
backtest, optimizer grid, replay, walk-forward logic, strategy package,
paper/testnet/live path, exchange API, credential, deploy file, broad mining,
martingale, averaging down, or two-exchange logic is approved by this document
alone.

## Authority Chain

The immediate approval chain is:

1. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md` stopped
   automatic BTCUSDT-only price-range audit work at
   `range_post_rotation_premise_failure_pivot_stopped_no_next_btcusdt_price_only_audit`.
2. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SPEC.md` parked a materially
   different context idea where BTCUSDT is regime context rather than the
   promoted authority symbol.
3. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_SCOPE_REVIEW.md` approved only this
   separate zero-trade audit brief-writing step.

The reviewed premise `router_gated_boundary_reclaim_rotation_v1` remains closed.
Do not convert its `278` context segments, `97` boundary-reclaim events, or the
router's `1,299` `tradable_rotation` rows into trades.

## Minimum Audit Question

A later implementation may answer only:

```text
Do BTCUSDT regime buckets, known at closed decision-candle time, improve
separation of ETHUSDT and SOLUSDT local range states into usable, toxic,
rotation, continuation, or no-trade contexts?
```

This is a context-separation question. It is not an entry question, not a
strategy-P&L question, not a replay question, and not a portfolio-construction
question.

## Source Contract

The only approved inputs for the future audit are the already local Binance
USDT-M futures `5m` files:

| Symbol | Path | Loaded candles | Coverage | Validation status |
| --- | --- | ---: | --- | --- |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | `573,984` | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` | accepted after sorting |

Prior validation facts to preserve:

- sorted streams had `gap_count=0` and `duplicate_count=0`;
- zero-volume counts were BTCUSDT `66`, ETHUSDT `47`, and SOLUSDT `47`;
- SOLUSDT had one physical non-monotonic row and may be accepted only after
  sorting by open time;
- no spot data, source downloads, new symbols, derivatives files, private
  endpoints, exchange APIs, or broad mining are approved.

If any of those local source facts are missing, contradicted, or impossible to
reproduce in the later implementation, stop at:
`btc_regime_eth_sol_context_zero_trade_audit_source_gap`.

## Symbol Roles

BTCUSDT role:

- market-regime context only;
- diagnostic-only authority row if a common authority table needs a BTC row;
- closed-candle context features known at or before the decision candle close;
- no promotion to traded authority from this audit.

ETHUSDT and SOLUSDT roles:

- local range-state rows;
- possible authority rows only for context separation in the zero-trade audit;
- forward labels as labels only;
- no strategy promotion from this audit.

Any future passing context audit would still need a separate strategy premise
spec before entries, exits, P&L backtests, replay, walk-forward, or packaging.

## Allowed Feature Layer

The future audit may build high-level closed-candle features only.

Allowed BTCUSDT regime buckets:

- BTC range state from current range-state construction infrastructure;
- BTC realized-volatility bucket;
- BTC trend-pressure bucket;
- BTC impulse bucket;
- BTC participation proxy;
- BTC close location relative to BTC's active range;
- BTC higher-timeframe slope or direction proxy derived from the same local
  `5m` futures source.

Allowed ETHUSDT/SOLUSDT local buckets:

- local range geometry;
- local volatility state;
- local trend state;
- local impulse state;
- local participation proxy;
- close location relative to the local active range;
- relative strength versus BTC, computed only from closed candles.

The audit may define state IDs in this shape:

```text
btc_regime_eth_sol_v1::<symbol>::<timeframe>::<btc_regime_bucket>::<local_range_bucket>::<relative_strength_bucket>
```

State IDs must use only data known at closed decision-candle time.

## Anti-Leakage Rule

Forward labels may appear only in label, cohort, ranking, and summary artifacts.
They must never be premise inputs, state-ID inputs, router inputs, gating inputs,
or feature-bucket inputs.

The later implementation must make the time split explicit:

- premise/context features end at the closed decision candle;
- labels begin strictly after that decision candle;
- any relative-strength input must be computed from closed candles only;
- any future return, adverse excursion, favorable excursion, or resolution label
  belongs only to label/cohort artifacts.

If hidden future-label input is discovered, stop at:
`btc_regime_eth_sol_context_zero_trade_audit_rejected_future_label_leak`.

## Expected Future Artifacts

If explicitly approved later, the implementation should write generated outputs
under:

```text
results/futures-btc-regime-eth-sol-context-audit/
```

Expected audit-specific artifacts:

- `futures_btc_regime_eth_sol_context_sources.csv/json`;
- `futures_btc_regime_eth_sol_context_coverage.csv/json`;
- `futures_btc_regime_eth_sol_context_btc_states.csv/json`;
- `futures_btc_regime_eth_sol_context_local_states.csv/json`;
- `futures_btc_regime_eth_sol_context_relative_strength.csv/json`;
- `futures_btc_regime_eth_sol_context_labels.csv/json`;
- `futures_btc_regime_eth_sol_context_cohorts.csv/json`;
- `futures_btc_regime_eth_sol_context_rankings.csv/json`;
- `futures_btc_regime_eth_sol_context_summary.csv/json`.

Common outputs must remain zero-trade compatible:

- `summary.json` and `summary.csv` must be valid with `0` trades;
- `trades.json` must contain no trades;
- common outputs must not be repurposed to hold context rows as pseudo-trades.

## Cohort Gates For Later Implementation

The future audit should rank context cohorts by separation quality, not P&L.
The later implementation may define exact numeric thresholds, but the gates must
preserve these rules:

- each candidate cohort must have adequate full-sample and split counts for
  each authority symbol it claims to describe;
- a useful cohort must survive the weakest split, not only the full sample;
- ETHUSDT and SOLUSDT may not rely on the same single fragile period;
- toxic/no-trade separation and useful/rotation/continuation separation must be
  reported separately;
- BTCUSDT context must add information versus ETH/SOL local state alone;
- diagnostic BTC rows may be reported, but BTC cannot become authority;
- no result may be treated as strategy evidence without a later strategy
  premise spec.

The later audit may stop with a passing context result only if it can show
closed-candle BTC regime buckets improve ETH/SOL context separation without
future-label leakage or closed-family reuse.

## Rejection Criteria

Reject the later implementation or stop immediately if it becomes any of these:

- BTCUSDT-only price-range reslice;
- retune, rename, or gate relaxation of `router_gated_boundary_reclaim_rotation_v1`;
- structured-compression rescue;
- ETH/SOL authority replay from the fragile walk-forward branch;
- BTCUSDT promotion;
- broad symbol mining;
- unapproved source expansion or source download;
- hidden future-label input;
- entries, exits, P&L strategy backtests, optimizer grids, replay, or
  walk-forward;
- portfolio construction;
- spread-range or pair-range engine work;
- paper, testnet, live, exchange API, credential, or deploy work;
- martingale, averaging down, or two-exchange logic;
- import of old `binance-bot` strategy, scoring, order-management, live,
  credential, deploy, or portfolio-coordinator logic.

## Stop States

This brief stops at:

- `btc_regime_eth_sol_context_zero_trade_audit_brief_ready_for_user_approval`.

Allowed later implementation stop states, only after explicit user approval:

- `btc_regime_eth_sol_context_zero_trade_audit_source_gap`;
- `btc_regime_eth_sol_context_zero_trade_audit_rejected_closed_family_reslice`;
- `btc_regime_eth_sol_context_zero_trade_audit_rejected_future_label_leak`;
- `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`;
- `btc_regime_eth_sol_context_zero_trade_audit_passed_needs_strategy_premise_spec`.

The passing implementation stop state would still not authorize entries, exits,
P&L backtests, replay, walk-forward, packaging, paper/testnet/live paths, source
downloads, strategy promotion, or deployment.

## Later Implementation Verification

Only after explicit user approval, the implementation closeout should include:

```bash
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/range-strategy-lab-go-build /usr/local/go/bin/go run ./cmd/rangelab -futures-btc-regime-eth-sol-context-audit -out-dir results/futures-btc-regime-eth-sol-context-audit
wc -l results/futures-btc-regime-eth-sol-context-audit/*.csv
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```

The future implementation review must record source path, product, symbol,
coverage, row counts, gap count, duplicate count, zero-volume count, generated
artifact paths, common zero-trade output status, stop state, and exact command
outcomes.

## Current Documentation Closeout Verification

For this brief-writing closeout, run only:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
