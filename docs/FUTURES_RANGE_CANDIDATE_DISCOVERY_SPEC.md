# Futures Range Candidate Discovery Spec

Date: 2026-06-26

## Verdict

Stop state:
`range_discovery_spec_ready_for_audit_implementation`.

The project should move from one-premise-at-a-time rejection into a
range-first broad discovery funnel. The next milestone should implement a
non-trading discovery audit that compares multiple BTCUSDT futures
range-adjacent candidate families across `5m`, `15m`, `1h`, and `4h`, ranks
them, and sends only the strongest candidates to a baseline backtest brief.

This is not a strategy promotion. It is a faster path from idea discovery to
baseline backtesting while keeping the current offline, futures-only,
BTCUSDT-only boundaries.

## Scope Correction

The recent scope pivot was useful because it stopped automatic mining of the
failed BTCUSDT `5m` lane. It should not be interpreted as a permanent ban on:

- `5m` research;
- buy/sell touch hypotheses;
- single-candle reactions;
- boundary rejection;
- breakout continuation;
- failed-breakout re-entry;
- other range-first strategy shapes.

The durable boundary is narrower: do not rerun the exact failed idea under a
new label. A reviewed family can be revisited only when it is materially
reframed, compared against other families, and tied to a fast backtest gate
rather than another isolated review loop.

## Source Universe

The parent source remains Binance USDT-M futures BTCUSDT `5m`:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Required product | `-source-product binance-usdm-futures` |
| Market type | Binance USDT-M futures BTCUSDT `5m` |
| CSV lines including header | `573,985` |
| Loaded candles | `573,984` |
| Open-time coverage | `2021-01-01T00:00:00Z` through `2026-06-16T23:55:00Z` |
| Accepted manifest facts | `gap_count=0`, `duplicate_count=0`, `zero_volume_count=66`, `comparison_only=false`, `validation_status=accepted` |

The discovery audit should use:

- native `5m` parent candles;
- closed UTC `15m`, `1h`, and `4h` resamples from
  `docs/FUTURES_HIGHER_TIMEFRAME_RANGE_SOURCE_SPEC.md`.

No spot source, comparison-only manifest, alternate symbol, exchange API,
download, sibling repo mutation, paper/testnet/live wiring, or deploy work is
allowed.

## Candidate Families

The next audit should evaluate a compact menu instead of a single premise.
Each candidate must be closed-candle knowable and range-first.

| Family | Timeframes | Candidate Event | Non-Trading Outcome Labels |
| --- | --- | --- | --- |
| Mature balance persistence | `1h`, `4h`, optional `15m` cross-check | A compact multi-bar range with adequate inside closes and no fresh expansion | inside persistence, internal rotation, expansion failure |
| Boundary touch rejection | `5m`, `15m`, `1h` | Closed candle touches a mature range boundary without closing decisively beyond it | reject inward, accept outside, stall/none |
| Single-candle wick rejection | `5m`, `15m`, `1h` | Large wick at a mature range boundary with close back inside the range | inward follow-through, boundary break, no follow-through |
| Failed breakout re-entry | `5m`, `15m`, `1h` | Closed break outside a mature range followed by closed re-entry within a fixed window | re-entry continuation, second break, no follow-through |
| Clean breakout continuation | `15m`, `1h`, `4h` | Mature range compression followed by a decisive closed break | continuation beyond break, failed break/re-entry, no extension |

The audit may use existing detector/SR/report infrastructure, but it must not
reuse old spot evidence as authority or copy old strategy scoring. Candidate
families should compete in one ranked report.

## Discovery-To-Backtest Funnel

The next implementation should produce non-trading audit artifacts first, then
route quickly:

1. Discovery audit across the candidate menu.
2. Rank candidates by split-stable counts, directional follow-through,
   adverse move, excursion profile, and rough cost buffer.
3. If one or two candidates clear the gate, write the next brief as a baseline
   offline backtest prototype.
4. Optimize only after the baseline backtest is positive or near-positive on
   gross and net metrics with adequate split stability.

Do not optimize, retune, or add entry variants before a baseline candidate has
earned it.

## Required Audit Outputs

The next audit should write compact CSV/JSON artifacts under
`results/futures-range-candidate-discovery-audit/`:

- `futures_range_discovery_candidates.csv/json`
- `futures_range_discovery_summary.csv/json`
- `futures_range_discovery_rankings.csv/json`
- `futures_range_discovery_stability.csv/json`
- normal `source_manifest.json`, `summary.*`, and `trades.json`

Common `summary.*` and `trades.json` should remain zero-trade because this is
still non-trading discovery.

## Review Gate

A candidate family can move to a baseline backtest brief only if:

- source manifest is accepted Binance USDT-M futures BTCUSDT, not comparison;
- generated resample coverage is documented for any `15m`, `1h`, or `4h`
  candidate;
- each period split has adequate candidate counts;
- favorable outcome rate beats adverse outcome rate in every period split;
- adverse move and quick invalidation do not dominate the shape;
- rough cost buffer is plausible before sizing or optimization;
- both the full sample and worst split are coherent enough for a fixed-rule
  prototype.

If no family clears the gate, the next step should be a short review explaining
why and asking whether to broaden beyond BTCUSDT/range-first. If a family
clears, the next brief should be a baseline backtest, not another review-only
inventory.

## Stop States

- `range_discovery_audit_ready`: implementation and outputs complete; at least
  one family clears the review gate and the next brief is a baseline backtest.
- `range_discovery_no_backtest_candidate`: implementation and outputs
  complete; no family clears the review gate.
- `range_discovery_source_or_resample_gap`: source validation or resampling
  completeness fails.
- `range_discovery_codegen_or_test_blocked`: implementation or verification
  cannot complete.
- `range_discovery_review_only_no_strategy_change`: outputs document context
  but do not change next implementation direction.

## Next Brief

The canonical next brief should implement
`-futures-range-candidate-discovery-audit` as a non-trading audit. It should
not add entry, exit, scoring, sizing, optimization, paper/testnet/live wiring,
exchange API use, deploy files, credentials, grid, martingale, averaging down,
two-exchange logic, data downloads, spot comparison, symbol expansion, or
sibling repo mutation.
