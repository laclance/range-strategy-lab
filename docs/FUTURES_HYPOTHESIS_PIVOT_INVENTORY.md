# Futures Hypothesis Pivot Inventory

Date: 2026-06-25

## Verdict

Stop state:
`pivot_inventory_needs_user_hypothesis`.

The current reviewed material does not expose a materially different
futures-authoritative hypothesis ready for an automatic new audit. The
hold-inside/midline reaction surface was revalidated on Binance USDT-M futures
data, but the first executable close-back prototype failed P&L across the full
sample, every period split, and both sides.

Do not keep mining that entry family by retuning, broadening, or slicing it
again. The next step needs a new user-supplied hypothesis or data premise before
any new non-trading audit is opened.

This inventory is review-only. It adds no entries, exits, scoring, sizing,
strategy replacement, paper/testnet/live wiring, exchange API use, deployment,
credentials, grid, martingale, averaging down, or two-exchange logic.

## Source Authority

Current research authority is Binance USDT-M futures BTCUSDT 5m data:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Market type | Binance USDT-M futures |
| CSV lines including header | `573,985` |
| Loaded candles | `573,984` |
| Open-time coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Manifest status | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

Legacy spot-only reviews are historical context. They can describe old
thinking, code shape, and audit plumbing, but they cannot promote futures
strategy work unless rerun and reviewed under the futures source contract.

## Inventory

| Family | Status | Evidence Authority | Reusable Infrastructure | No-Retry Boundary | Possible Future Use |
| --- | --- | --- | --- | --- | --- |
| SR rejection, confirmation, and false-break timing | Closed | Legacy spot-only; not futures-authoritative | SR feature extraction, closed-candle timing labels, split review format | Do not reslice rejection, delayed confirmation, or reclaim timing into entries from old spot outputs | Only rerun if the user supplies a materially new futures SR premise |
| Compression breakout | Closed | Legacy spot-only; not futures-authoritative | Compression episode framing, frozen high/low handling, breakout audit artifacts | Do not keep tuning thresholds, delays, side buckets, or old breakout cohorts | Only revisit with a new futures compression premise, not a narrower version of the closed review |
| Range regime durability | Diagnostic/infrastructure | Legacy spot-only review, useful as detector plumbing context | Episode durability labels, persistence/quick-invalidation/trend-leakage summaries | Do not treat durability labels as entries, exits, or scoring inputs | Reuse as background quality control for a different futures audit |
| Detector durability and context refinement | Reusable infrastructure | Futures data impact review revalidated the source contract and context role, not entry promotion | Detector profiles, delayed hold-inside context rules, source manifests, stability outputs | Do not promote detector profiles or `hold_3_inside`/`hold_6_inside` context directly into trades | Reuse as feature extraction and filtering infrastructure for a new futures premise |
| Hold-inside directional edge | Closed | Legacy spot-only conclusion; no futures promotion and no directional edge approval | Paper-side audit shape and side/bucket split summaries | Do not promote `paper_side=toward_high` or `paper_side=toward_low` | Only reopen if a new futures hypothesis changes the measured directional target |
| Hold-inside midline transition/reaction | Diagnostic only | Reaction surface was rerun on futures, then downstream prototype failed | Reindexed midline-event audit, event buckets, reaction labels, funnel gate | Do not broaden into `hold_6_inside`, `mid_close_across`, side cohorts, `hold_3_inside_mid_50`, or old spot authority | Reuse only as diagnostics if a different futures premise needs midline-event context |
| Futures midline touch prototype | Closed futures-authoritative failure | Binance USDT-M futures full-history prototype review | Explicit prototype flagging, joined signal/trade details, side/split summaries, entry-geometry guard | No promotion, retune, broadening, paper/testnet/live wiring, or strategy replacement | Carry forward code/reporting shape only |

## Pivot Read

The usable carry-forward is infrastructure, not a next strategy:

- futures source guard and manifest discipline
- audit-only CLI flag pattern
- closed-candle and next-bar-open semantics
- split-stability review format
- event-to-trade artifact joining
- entry-geometry guard

The reviewed hypothesis families now form an exclusion map. The strongest
futures-authoritative candidate, `hold_3_inside` plus first `mid_touch` within
`12` bars plus `mid_50`, did not survive conversion into a first executable
trade template. That failure blocks automatic follow-up work inside the same
hold-inside/midline family.

## Next Step

Ask for a new hypothesis or data premise before opening another audit. A valid
next brief may still be non-trading and futures-only, but it should start from
a materially different question than:

- SR timing reslices
- compression breakout reslices
- detector/context promotion
- hold-inside directional continuation
- hold-inside/midline close-back entries
- spot-based promotion evidence

Until that premise is explicit, keep `lab.EmptyStrategy` as the default and do
not add strategy behavior.
