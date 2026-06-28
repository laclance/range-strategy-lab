# Futures Derivatives Context Source Scope Review

Date: 2026-06-28

## Verdict

Stop state:
`derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.

This docs-only review approves a later separate zero-trade source-audit brief
for derivatives market-data context, but only as a source and alignment brief.
It does not approve data downloads, Go code, CLI flags, generated result
directories, context feature implementation, entries, exits, P&L backtests,
optimizer grids, replay, walk-forward, source mining, source promotion, paper,
testnet, live execution, exchange API keys, private endpoints, credentials,
deploy files, martingale, averaging down, or two-exchange logic.

The later brief is justified because derivatives market data is materially
different from the closed candle-only and BTC regime plus ETH/SOL context
surfaces. The later brief must still prove a durable local/offline source
contract before any implementation. If durable local files or an explicitly
approved offline materialization plan are missing, the later source audit must
stop at source gap.

## Authority Chain

The immediate boundary is:

1. `docs/FUTURES_RANGE_POST_ROTATION_PREMISE_FAILURE_PIVOT_REVIEW.md` stopped
   automatic BTCUSDT-only price-range audit work.
2. `docs/FUTURES_BTC_REGIME_ETH_SOL_CONTEXT_ZERO_TRADE_AUDIT_REVIEW.md` closed
   BTC regime plus ETH/SOL context with `0` passing cohorts at
   `btc_regime_eth_sol_context_zero_trade_audit_failed_no_usable_context`.
3. `docs/FUTURES_DERIVATIVES_CONTEXT_SOURCE_EXPANSION_SPEC.md` parked
   derivatives context as market-data source expansion only, pending explicit
   source approval.
4. The current user request approves this docs-only source scope review, not
   implementation.

This branch is a source premise, not a strategy premise.

## Current Local Inventory

No derivative market-data rows are currently approved as direct inputs for this
lab.

The current durable local source directory checked from this repo is:

```text
../binance-bot/data/
```

That directory currently contains durable candle CSVs only: BTCUSDT spot and
futures candles, BTCUSD external spot candles, BTCUSDT sample candles, plus the
ETHUSDT and SOLUSDT Binance USDT-M futures `5m` files used in the prior
context audit. It does not contain durable funding, open-interest, mark-price,
index-price, premium-index, taker-flow, aggregate-trade, order-book,
long/short, or liquidation source files.

The checked raw directory:

```text
../binance-bot/data/raw/
```

had no files in the reviewed search depth.

The only currently approved durable candle anchors for a later source-audit
brief are the already validated Binance USDT-M futures `5m` files:

| Symbol | Path | Role in later source audit |
| --- | --- | --- |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor only if ETH is explicitly kept in scope |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | candle alignment anchor only if SOL is explicitly kept in scope |

Known source facts from prior accepted validation:

- each file has `573,984` loaded candles;
- coverage is `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z`;
- sorted streams had `gap_count=0` and `duplicate_count=0`;
- zero-volume counts were BTC `66`, ETH `47`, SOL `47`;
- SOL had one physical non-monotonic row and was accepted only after sorting.

These candle files are not derivatives context sources. They are alignment
anchors only.

## Adjacent Source-Proof References

Adjacent `../binance-bot/research/` artifacts contain useful source-proof
references, but they are not approved as direct data inputs for this lab.

Usable only as reference evidence for a later brief:

| Family | Adjacent evidence | Scope decision |
| --- | --- | --- |
| Funding rate history | `../binance-bot/research/2026-06-18_futures_native_non_ohlcv_source_proof/README.md`; `../binance-bot/research/2026-06-18_futures_funding_state_transition/event_study/source_inventory.csv` | Candidate for later source-audit brief only. Must be materialized as durable approved local files or rejected. |
| Mark/index/premium basis | `../binance-bot/research/2026-06-18_futures_native_non_ohlcv_source_proof/README.md`; `../binance-bot/research/2026-06-18_futures_perp_basis_reversion/event_study/source_inventory.csv` | Candidate for later source-audit brief only. Best first family if durable local files are approved because it aligns naturally to `5m`. |
| Aggregate trades / taker flow | `../binance-bot/research/2026-06-18_futures_aggtrades_taker_flow_event_study/source_inventory.csv`; `../binance-bot/research/memory/futures_source_proof_task_ranking_2026_06.md` | Candidate source-proof reference only. High data-volume risk; not first implementation input. |
| Liquidations / force orders | `../binance-bot/research/2026-06-20_futures_liquidation_source_proof/source_inventory.csv` | Rejected for full-era public historical source from current evidence. |
| Open interest / metrics | `../binance-bot/research/memory/futures_native_non_ohlcv_source_proof.md` | Rejected for current full-era BTC/ETH/SOL source scope because ETH/SOL metrics start later and schema/finality are not sufficiently proven. |

The adjacent source inventories reference cache paths under `/tmp`, such as
funding, perp-basis, and aggregate-trade caches. Those paths are not durable
repo-local source paths and are not approved as inputs for this lab.

## Allowed Later Brief Scope

A later zero-trade source-audit brief may be written with this narrow question:

```text
Can one or more durable local/offline derivatives market-data sources be
approved and aligned to the existing Binance USDT-M futures 5m candle anchors
without look-ahead, private endpoints, live streams, or source downloads?
```

The later brief may consider only these source families:

1. Funding rate history.
2. Mark-price, index-price, or premium-index klines for basis context.
3. Aggregate trades only as an explicitly secondary, high-volume source-proof
   candidate.

The later brief must choose at most one first implementation candidate family.
The recommended first family is mark/index/premium basis because it is a public
archive kline family with `5m` timestamp structure and direct candle alignment.
Funding is the second candidate because it has low row volume but lower native
cadence and tail/finality caveats. Aggregate trades should stay parked until a
lighter source family either fails source proof or fails context value.

## Required Later Source Contract

The later brief must define a manifest before implementation. For every
candidate source it must specify:

- durable local file path or explicitly approved offline materialization path;
- source family and owner;
- symbol scope;
- product scope;
- timestamp semantics;
- first and last timestamp;
- row count;
- duplicate count;
- gap or missing-interval count;
- native cadence;
- publication/finality rule;
- timezone assumption;
- max staleness if forward-filled;
- alignment target timeframe;
- missing-data policy;
- checksum or provenance identifier when available;
- comparison-only status;
- validation status.

The later brief must require source rejection when timestamp finality,
publication lag, schema, symbol identity, product identity, or missing-data
behavior cannot be explained.

## Forbidden Source Scope

The following remain forbidden from this scope review:

- source downloads or live source probes;
- private account endpoints;
- exchange API keys or signed requests;
- paper, testnet, or live exchange work;
- WebSocket live streams as historical substitutes;
- `/tmp` caches as durable source inputs;
- adjacent research result CSVs as strategy evidence;
- spot or cross-exchange sources;
- broad symbol mining;
- open-interest metrics from the current evidence set;
- long/short ratio endpoints without full-era public archive proof;
- liquidation or force-order endpoints from the current evidence set;
- order-book/depth sources without a separate historical archive proof;
- aggregate-trade implementation before a narrower funding or basis source
  brief explicitly rejects the lighter family;
- any direct strategy, entry, exit, P&L backtest, replay, walk-forward,
  optimizer grid, or promotion path.

## Decision On Later Brief

A later zero-trade audit brief is justified only for source and alignment
approval. It is not yet justified for context-gain implementation because this
review found no approved durable derivatives market-data source rows in the
lab's current local data scope.

The later brief should stop before code at:

```text
derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval
```

If the later brief cannot name durable local/offline source files or an
explicitly approved offline materialization plan, it must instead stop at:

```text
derivatives_context_zero_trade_source_audit_brief_rejected_source_gap
```

## Stop States

This review stops at:

- `derivatives_context_source_scope_review_approved_needs_zero_trade_source_audit_brief`.

Allowed next brief-writing stop states:

- `derivatives_context_zero_trade_source_audit_brief_ready_for_user_approval`;
- `derivatives_context_zero_trade_source_audit_brief_rejected_source_gap`;
- `derivatives_context_zero_trade_source_audit_brief_rejected_live_or_private_api_path`;
- `derivatives_context_zero_trade_source_audit_brief_rejected_closed_family_rescue`.

No implementation stop state is approved by this review.

## Verification

Commands run for this docs-only closeout:

```bash
rg --files ../binance-bot/data | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade|trade)"
rg --files ../binance-bot/data | rg -i "(btcusdt|ethusdt|solusdt).*futures.*(5m|1h|15m|4h|um)|futures_um"
rg --files . | rg -i "(funding|fund|open.?interest|basis|premium|mark|index|taker|long.?short|book|depth|liquid|oi|agg.?trade)"
rg -n "funding|open interest|basis|premium|taker|long/short|order-book|order book|derivatives" docs memory README.md AGENTS.md
rg --files ../binance-bot | rg -i "(funding|open[_-]?interest|\boi\b|oi_|basis|premium[_-]?index|premiumIndex|mark[_-]?price|index[_-]?price|taker|long[_-]?short|longshort|depth|order[_-]?book|book[_-]?ticker|liquidat|agg[_-]?trade)"
rg --files ../binance-bot/data
ls -lah ../binance-bot/data
find ../binance-bot/data/raw -maxdepth 3 -type f
```

Final closeout also requires:

```bash
rg -n "CODEX_BRIEF|NEXT_CODEX_BRIEF" README.md docs memory AGENTS.md
git diff --check
git status --short
```
