# Futures Scope Pivot Review

Date: 2026-06-26

## Verdict

Stop state:
`range_scope_pivot_ready_for_higher_timeframe_source_spec`.

BTCUSDT 5m range research is not close enough to keep mining automatically.
The next range-only route is a BTCUSDT higher-timeframe futures source and
premise specification, not an audit, prototype, strategy replacement, or symbol
expansion.

This review pauses BTCUSDT 5m range continuation work until a materially new
premise appears. It selects a docs/source-spec task for closed UTC `15m`, `1h`,
and `4h` BTCUSDT futures bars derived from the accepted 5m source. That keeps
the project range-only and BTCUSDT-only while changing the source shape enough
to avoid another reslice of failed 5m families.

No code, CLI flag, audit run, generated result, entry, exit, scoring, sizing,
paper/testnet/live wiring, exchange API use, deploy file, credential, grid,
martingale, averaging down, or two-exchange logic is approved here.

## Current Source Authority

The current lab authority remains Binance USDT-M futures BTCUSDT 5m:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Market type | Binance USDT-M futures |
| Symbol / interval | `BTCUSDT` / `5m` |
| CSV lines including header | `573,985` |
| Loaded candles | `573,984` |
| Open-time coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Manifest status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

This source is valid for BTCUSDT 5m work and can be the parent source for a
higher-timeframe source-spec review. It is not, by itself, authority for any
higher-timeframe range audit until resampling rules, coverage, finality, and
premise gates are specified.

## Closed 5m Range Families

| Family | Current Status | Evidence Authority | Reusable Piece | Boundary |
| --- | --- | --- | --- | --- |
| SR rejection, confirmation, and false-break timing | Closed | Legacy spot-only; not futures-authoritative | SR feature extraction and timing-label report shape | Do not reslice SR timing into entries from old spot outputs |
| Compression breakout | Closed | Legacy spot-only; not futures-authoritative | Compression episode framing and breakout artifact pattern | Do not tune thresholds, delays, or side buckets without a materially new futures premise |
| Range regime durability | Diagnostic/infrastructure | Legacy spot review; useful only as detector context | Episode durability, quick-invalidation, and trend-leakage summaries | Do not treat durability labels as entry, exit, or scoring inputs |
| Detector durability/context refinement | Reusable infrastructure | Futures review revalidated source discipline and context role, not entries | Detector profiles, delayed hold-inside context, stability outputs | Do not promote detector profiles or hold-inside context directly into trades |
| Hold-inside directional edge | Closed | Legacy spot-only conclusion with no futures promotion | Paper-side and bucket split report shape | Do not promote `toward_high` or `toward_low` |
| Hold-inside midline transition/reaction | Diagnostic only | Reaction was futures-rerun, but downstream prototype failed | Midline-event labels and funnel gate | Do not broaden into `hold_6_inside`, `mid_close_across`, side cohorts, or old spot authority |
| Futures midline touch prototype | Closed futures-authoritative failure | Binance USDT-M futures full-history prototype review | Explicit prototype flagging and joined signal/trade artifacts | No promotion, retune, broadening, paper/testnet/live wiring, or strategy replacement |
| Futures impulse absorption | Closed futures-authoritative failure | Binance USDT-M futures full-history audit review | Prior-window percentile ranks and horizon labeling | Continuation dominated midpoint reclaim; do not convert into a prototype |

The 5m lane now acts as an exclusion map plus infrastructure. The useful code
and review patterns can support a future audit, but the reviewed 5m hypotheses
do not justify more automatic mining.

## Sibling Repo Context

Sibling repos are process or exclusion context only. They are not promotion
authority for this lab.

| Source | Label | Range-Relevant Fact | Use Here | Boundary |
| --- | --- | --- | --- | --- |
| `~/binance-bot` Binance-only range seed | Exclusion-only | BTCUSDT 5m range sleeve on legacy spot data failed even gross and collapsed after costs. | Caution against simple single-venue SR fades and static parameter sweeps. | Legacy spot-only; do not import as futures evidence. |
| `~/binance-bot` cross-exchange spread seed | Exclusion-only | Venue-spread mean reversion during range regimes failed local data gates and lacked modern overlap. | Reinforces blocking cross-exchange range work here. | Two-exchange execution remains forbidden. |
| `~/binance-bot` multi-pair selector work | Exclusion-only | Broad alt selection and BTC+SOL overlays did not justify widening. | Caution against broad symbol selectors. | Not range-only authority for this lab. |
| `~/binance-bot` daily ATR contraction | Outside-scope context | A higher-timeframe daily futures trend/volatility path had positive evidence. | Shows non-range futures work may exist elsewhere. | Not a range strategy; cannot justify this pivot. |
| `~/crypto-trading-bot` BTC/ETH USD-M source contract | Process-only | BTC/ETH `1h` candles, complete `4h`/`1d` resamples, and funding-aware cutoff are documented there. | Useful source-contract discipline for future higher-timeframe specs. | Do not import strategy results, parameters, or Python implementation. |
| `~/crypto-trading-bot` BTC/ETH 4h range reversal/re-entry | Exclusion-only | The frozen BTC/ETH USD-M 4h range-reversal/re-entry workbench found no benchmark-beating edge. | Warning against porting that family. | Do not rerun, retune, or port it here. |
| `~/crypto-trading-bot` broader USD-M screens | Exclusion-only | Broader return-targeted and volatility work emphasizes source discipline and stop states. | Caution against broad automatic screening. | Broad multi-pair search is not authorized. |

## Scope Lane Decision

| Lane | Decision | Reason | Next Action |
| --- | --- | --- | --- |
| BTCUSDT 5m range continuation | Not selected | Current futures-authoritative 5m prototypes/audits failed or stayed diagnostic. | Pause until a materially new 5m premise appears. |
| BTCUSDT higher-timeframe range | Selected | Keeps BTCUSDT-only scope while changing the closed-candle source shape to `15m`/`1h`/`4h`. | Write a source/premise spec before any audit. |
| BTC/ETH range-only source expansion | Deferred | Could be useful later, but it changes symbol scope and sibling BTC/ETH range work already has exclusion warnings. | Reconsider only after BTCUSDT higher-timeframe source/premise review or explicit user scope change. |
| Larger multi-pair range universe | Blocked | Broad selectors and symbol mining are not justified by current evidence. | No generic universe screen. |
| Cross-exchange range spread | Out of scope | This lab forbids two-exchange execution and modern multi-venue proof is absent. | Requires explicit scope change. |
| Non-range higher-timeframe trend/volatility | Out of scope | Positive-looking sibling evidence is trend/volatility, not range strategy. | Requires explicit user decision to stop being range-only. |

## Recommended Next Route

Create `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.

The next task should be source-spec and premise-spec only. It should:

- derive closed UTC `15m`, `1h`, and `4h` candidate bars from the accepted
  BTCUSDT futures 5m CSV;
- define resampling finality, coverage, no-gap/no-duplicate acceptance gates,
  and the parent-source manifest facts that must carry through;
- require a higher-timeframe range premise that is materially different from
  closed 5m SR timing, compression, hold-inside/midline, and impulse surfaces;
- define what would falsify the premise before any entry prototype;
- continue to block data downloads, audit outputs, entries, exits, scoring,
  sizing, live-adjacent work, and symbol expansion.

Until that spec exists and is reviewed, keep `lab.EmptyStrategy` as the default
and do not add strategy behavior.
