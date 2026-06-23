# Futures Scope Pivot Review Spec

Date: 2026-06-26

## Purpose

This is a planning specification for a future
`docs/FUTURES_SCOPE_PIVOT_REVIEW.md`.

The review should decide whether this lab should keep pressing BTCUSDT 5m
range research or pivot the range-strategy scope to a different futures
timeframe or tightly controlled symbol set. The pivot remains range-only. It
does not authorize trend-following, volatility expansion trend systems,
cross-sectional momentum, carry, grid, martingale, averaging down,
cross-exchange execution, paper/testnet/live work, exchange API use, deploy
work, or strategy replacement.

## Current Authority

The current authoritative source for this lab remains Binance USDT-M futures
BTCUSDT 5m:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Market type | Binance USDT-M futures |
| CSV lines including header | `573,985` |
| Loaded candles | `573,984` |
| Open-time coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Manifest status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

That source is enough for BTCUSDT 5m audits. It is not enough by itself for a
higher-timeframe or multi-symbol pivot until the review specifies the new
source contract, resampling rules, coverage, and acceptance gates.

## Why Pivot

The BTCUSDT 5m range path has been tested seriously enough that the current
state is not "almost there." It is an exclusion map plus reusable
infrastructure:

- hold-inside/midline reaction looked promising as diagnostics but the
  first futures prototype failed full sample, every period split, and both
  sides after costs;
- abnormal OHLCV impulse absorption had ample candidates but was
  continuation-dominant, not midpoint-reclaim-dominant;
- earlier SR timing, compression breakout, range durability, detector/context,
  and hold-inside directional families are closed, diagnostic-only, or
  reusable infrastructure rather than near-entry surfaces.

The future review should therefore ask a scope question before another audit:
is the next honest range strategy attempt still BTCUSDT 5m, or should the
range premise move to a different closed-candle source shape?

## Sibling Repo Evidence Boundaries

Sibling repos are useful as process context and exclusion evidence, not as
authority for this lab.

| Source | Range-Relevant Fact | Use In Review | Boundary |
| --- | --- | --- | --- |
| `~/binance-bot` range seed | A Binance-only BTCUSDT 5m range sleeve on spot data failed even gross and collapsed after costs. | Caution against simple single-venue support/resistance fades and static parameter sweeps. | Legacy spot-only and not futures-authoritative. Do not import results as lab evidence. |
| `~/binance-bot` cross-exchange spread seed | Venue-spread mean reversion during range regimes failed local data gates and lacked modern overlap. | Exclusion evidence for cross-exchange range execution. | Two-exchange work remains out of scope here. |
| `~/binance-bot` multi-pair selector work | Broad alt selection damaged the BTC-only edge; BTC+SOL variants failed quality gates. | Caution against broad symbol selectors and alt overlays. | Not a range-only futures premise and not authority for this lab. |
| `~/binance-bot` daily ATR contraction | A higher-timeframe daily futures trend/volatility path had positive evidence. | Mention only as outside-scope evidence that non-range futures work may exist elsewhere. | Not a range strategy; do not use it to justify this range-only pivot. |
| `~/crypto-trading-bot` BTC/ETH source contract | BTC/ETH Binance USD-M 1h candles with complete 4h/1d resamples and funding-aware cutoff exist as a source-contract pattern. | Possible template if the review recommends a higher-timeframe BTC/ETH range source spec. | Process/source-contract reuse only. |
| `~/crypto-trading-bot` 4h range reversal | The bounded BTC/ETH USD-M 4h range-reversal/re-entry workbench found no benchmark-beating edge. | Exclusion evidence for the frozen 4h range-reversal/re-entry family. | Do not rerun, retune, or port the closed family here. |
| `~/crypto-trading-bot` broad USD-M screens | Broader return-targeted screens and volatility breakout lanes failed gates. | Caution against broad automatic screening. | Broad multi-pair search is not the next range-only move. |

## Candidate Scope Lanes

The future review should classify these lanes before any implementation:

| Lane | Review Default | What Could Make It Viable | What Is Blocked |
| --- | --- | --- | --- |
| Continue BTCUSDT 5m range | Not recommended automatically | A materially new closed-candle range observable that is not a reslice of closed hold-inside/midline, SR timing, compression, or impulse surfaces | More 5m reslicing, retuning, side buckets, old spot authority |
| BTCUSDT higher-timeframe range | Plausible review lane | A source spec for closed 15m/1h/4h bars derived from accepted futures candles, plus a range premise materially different from the failed 5m families | Immediate audit/prototype before source and premise review |
| BTC/ETH range-only source expansion | Plausible review lane | A narrow BTC/ETH futures source contract and benchmark plan for range behavior only | Broad selectors, alt overlays, 25-symbol screens, imported ETH/SOL results |
| Larger multi-pair range universe | Not first | Only after BTC/ETH range evidence justifies widening | Generic symbol mining and return-targeted ranked screens |
| Cross-exchange range spread | Out of scope | Requires explicit project scope change and modern multi-venue source proof | Two-exchange execution, venue-spread live architecture |
| Higher-timeframe trend/volatility strategy | Out of scope for this review | Requires explicit user decision to stop being range-strategy-only | Using daily ATR contraction or trend-core evidence as range promotion |

## Review Questions

The future `FUTURES_SCOPE_PIVOT_REVIEW` should answer:

1. Is BTCUSDT 5m range research still a good primary scope after the futures
   failures?
2. If not, is the next range-only scope a higher timeframe for BTCUSDT, a
   narrow BTC/ETH source contract, or no viable range lane?
3. Which sibling-repo artifacts are useful only as process templates or
   exclusion evidence?
4. What exact source contract would the next non-trading range audit need?
5. What closed-candle observable would define the next candidate event?
6. What would falsify the premise before any entry prototype?

## Required Future Review Output

The future review should create `docs/FUTURES_SCOPE_PIVOT_REVIEW.md` with:

- current BTCUSDT 5m futures source facts;
- a table of closed 5m range families and why they do not justify more
  automatic mining;
- a sibling-repo evidence table with process-only and exclusion-only labels;
- a lane table for BTCUSDT 5m, BTCUSDT higher timeframe, BTC/ETH range-only,
  larger multi-pair, cross-exchange, and non-range trend work;
- one recommended next route, or a stop state saying no range-only lane is
  viable without a new premise;
- a refreshed `memory/NEXT_CODEX_BRIEF.md` based on that route.

## Stop States

- `range_scope_pivot_ready_for_higher_timeframe_source_spec`: the review
  selects a higher-timeframe BTCUSDT range source/spec task.
- `range_scope_pivot_ready_for_btc_eth_range_source_spec`: the review selects
  a narrow BTC/ETH range-only source/spec task.
- `range_scope_pivot_no_viable_range_lane`: no range-only lane is justified
  from current evidence.
- `range_scope_pivot_needs_user_scope_change`: the only attractive direction
  is non-range, cross-exchange, live-adjacent, or otherwise outside the current
  project scope.
- `range_scope_pivot_review_only_no_strategy_change`: the review documents
  context but does not change implementation direction.

## Verification For This Spec

This spec is docs/memory-only. It authorizes no audit run, strategy code, data
mutation, source expansion, entry, exit, scoring, sizing, paper/testnet/live
work, exchange API use, deploy work, grid, martingale, averaging down, or
two-exchange logic.

