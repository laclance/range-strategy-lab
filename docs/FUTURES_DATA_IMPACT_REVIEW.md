# Futures Data Impact Review

Date: 2026-06-25

## Verdict

Stop state:
`futures_reaction_gate_passed_needs_minimal_entry_brief`.

The old spot-based hold-inside midline touch prototype idea survives the move
to Binance USDT-M futures data, but the authority is now the futures rerun, not
the old spot review. The only surface revalidated for a first minimal offline
entry prototype is:

- `hold_3_inside`
- first `mid_touch` within `12` bars after the hold decision
- event close-position bucket `mid_50`
- closed-candle event decision, with any later prototype entering no earlier
  than the next bar open

This is not strategy promotion or live approval. No entry, exit, scoring,
sizing, strategy replacement, paper, testnet, live, exchange-key, deploy, grid,
martingale, averaging-down, or two-exchange work was added in this review.
`cmd/rangelab` still uses `lab.EmptyStrategy`.

## Inputs Reviewed

Generated futures audit paths:

- `results/futures-detector-context-refinement-audit/`
- `results/futures-hold-inside-midline-transition-audit/`
- `results/futures-hold-inside-midline-reaction-audit/`

Every run wrote `source_manifest.json` with the same accepted source contract:

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Venue / product | `Binance` / `Binance USDT-M futures` |
| Symbol / interval | `BTCUSDT` / `5m` |
| Row count | `573,984` loaded candles |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Schema | Binance archive shape, `open_time` through `ignore` |
| Timestamp semantics | `open_time` |
| Finality rule | closed 5m candles; `close_time = open_time + 5m - 1ms` |
| Gaps / duplicates | `0` / `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Audit sizes:

| Audit | Candidate Rows | Summary Rows | Stability Rows | Other |
| --- | ---: | ---: | ---: | --- |
| Detector context refinement | `117,848` | `640` | `160` | profiles `8`, rules `5` |
| Hold-inside midline transition | `8,600` | `672` | `168` | profiles `1`, rules `3` |
| Hold-inside midline reaction | `10,172` | `1,240` | `336` | funnel rows `24` |

All three runs loaded `573,984` candles through close time
`2026-06-16T23:59:59Z` and reported `strategy=empty trades=0`.

## Supporting Futures Context

The futures detector context rerun preserved the delayed hold-inside context:

| Context | H12 Min Cand | H12 Min Persist | H12 Max Quick Inv | H12 Max Trend |
| --- | ---: | ---: | ---: | ---: |
| `hold_3_inside` | `232` | `39.22%` | `25.00%` | `37.93%` |
| `hold_6_inside` | `174` | `45.98%` | `24.14%` | `33.95%` |
| `hold_3_inside_mid_50` | `106` | `45.96%` | `17.03%` | `33.96%` |

The futures transition rerun also preserved the midline as a useful observation
point, while still not making the transition labels entry inputs:

| Rule | Horizon | Min Cand | Min Touch Mid | Min Close Across | Min Touch Before Boundary | Min Cross Before Break | Min Persist | Max Quick Inv | Max Trend |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside` | 3 | `232` | `28.88%` | `18.97%` | `25.43%` | `18.10%` | `75.00%` | `25.00%` | `19.83%` |
| `hold_3_inside` | 6 | `232` | `40.09%` | `30.17%` | `34.48%` | `28.45%` | `56.90%` | `25.00%` | `27.16%` |
| `hold_3_inside` | 12 | `232` | `55.60%` | `46.98%` | `41.81%` | `39.22%` | `39.22%` | `25.00%` | `37.93%` |
| `hold_6_inside` | 3 | `174` | `30.46%` | `22.41%` | `28.16%` | `21.84%` | `75.86%` | `24.14%` | `19.54%` |
| `hold_6_inside` | 6 | `174` | `42.53%` | `33.91%` | `37.36%` | `32.18%` | `62.07%` | `24.14%` | `25.29%` |
| `hold_6_inside` | 12 | `174` | `56.32%` | `50.57%` | `45.40%` | `45.98%` | `45.98%` | `24.14%` | `33.95%` |

These rows support continuing to the reaction surface, but they still do not
authorize entries by themselves.

## Reaction Gate

Futures funnel rows by weakest split:

| Rule | Event | Min Source Holds | Min Event Count | Min Event Rate | Max Missing Rate | Missing Future | Avg Delay Range | Gate |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| `hold_3_inside` | `mid_touch` | `232` | `129` | `55.60%` | `44.40%` | `0` | `3.92-4.43` | pass |
| `hold_3_inside` | `mid_close_across` | `232` | `109` | `46.98%` | `53.02%` | `0` | `4.77-5.22` | fail: missing dominates |
| `hold_6_inside` | `mid_touch` | `174` | `98` | `56.32%` | `43.68%` | `0` | `3.26-4.07` | fail: weak split below `100` events |
| `hold_6_inside` | `mid_close_across` | `174` | `88` | `50.57%` | `49.43%` | `0` | `3.89-4.88` | fail: weak split below `100` events |

The primary `hold_3_inside + mid_touch` surface clears the event-count and
event-rate gate in every period split.

Primary reaction rows:

| Row | Horizon | Min Cand | Min Persist | Max Quick Inv | Max Trend | Min Close Back | Min Mid Reject Before Boundary | Max Boundary Before Reject |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `hold_3_inside + mid_touch`, all bucket | 3 | `129` | `81.40%` | `18.60%` | `13.11%` | `47.54%` | `44.67%` | `21.71%` |
| `hold_3_inside + mid_touch`, all bucket | 6 | `129` | `66.67%` | `18.60%` | `22.48%` | `56.59%` | `48.06%` | `29.46%` |
| `hold_3_inside + mid_touch`, all bucket | 12 | `129` | `46.95%` | `18.60%` | `34.27%` | `70.49%` | `55.04%` | `35.66%` |
| `hold_3_inside + mid_touch`, `mid_50` bucket | 3 | `114` | `86.84%` | `13.16%` | `10.45%` | `50.00%` | `47.73%` | `15.79%` |
| `hold_3_inside + mid_touch`, `mid_50` bucket | 6 | `114` | `71.21%` | `13.16%` | `21.93%` | `57.89%` | `52.63%` | `23.68%` |
| `hold_3_inside + mid_touch`, `mid_50` bucket | 12 | `114` | `50.51%` | `13.16%` | `32.83%` | `69.30%` | `59.65%` | `30.70%` |

The `6` bar `mid_50` row is the cleanest first prototype target. It keeps the
weakest split above `100` candidates and preserves the intended behavior:
close-back and mid-rejection-before-boundary remain above
boundary-before-rejection, while quick invalidation stays far below the
reaction rates.

Split-level `h6` futures detail:

| Split | Bucket | Cand | Close Back | Mid Reject Before Boundary | Boundary Before Reject | Quick Inv |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | all | `244` | `61.07%` | `55.74%` | `27.05%` | `15.98%` |
| `2021_2022_stress` | `mid_50` | `220` | `62.73%` | `59.09%` | `21.82%` | `11.36%` |
| `2023_2024_oos` | all | `213` | `67.14%` | `60.09%` | `24.41%` | `16.90%` |
| `2023_2024_oos` | `mid_50` | `198` | `70.20%` | `64.65%` | `19.19%` | `11.62%` |
| `2025_2026_recent` | all | `129` | `56.59%` | `48.06%` | `29.46%` | `18.60%` |
| `2025_2026_recent` | `mid_50` | `114` | `57.89%` | `52.63%` | `23.68%` | `13.16%` |

## Spot Comparison

The old spot review is now comparison context only. The futures rerun is close
enough to the spot shape to revalidate the prototype idea:

| Source | Bucket | H6 Min Cand | H6 Min Persist | H6 Max Quick Inv | H6 Max Trend | H6 Min Close Back | H6 Min Mid Reject | H6 Max Boundary Before Reject |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| spot | all | `116` | `69.06%` | `18.39%` | `22.22%` | `55.17%` | `45.69%` | `28.25%` |
| futures | all | `129` | `66.67%` | `18.60%` | `22.48%` | `56.59%` | `48.06%` | `29.46%` |
| spot | `mid_50` | `104` | `74.04%` | `11.54%` | `19.23%` | `58.65%` | `50.96%` | `22.12%` |
| futures | `mid_50` | `114` | `71.21%` | `13.16%` | `21.93%` | `57.89%` | `52.63%` | `23.68%` |

Futures is slightly weaker on persistence, quick invalidation, trend, and
boundary-before-rejection in the `mid_50` h6 row, but not enough to break the
gate. Candidate counts are stronger, event occurrence remains split-stable,
and the post-event behavior still points to the same rejection/close-back
prototype shape.

Side-specific cohorts remain diagnostic only. At h6, `below_mid` has weakest
split `60` candidates and `above_mid` has weakest split `68` candidates, so no
side-specific conclusion reaches the `100` candidate threshold yet.

## Carry Forward

Revalidated on futures:

- first minimal offline prototype around `hold_3_inside + mid_touch + mid_50`

Reusable as code only:

- the detector/context, transition, and reaction audit implementations
- old spot CSV extraction and review structure

Historical or diagnostic only:

- old spot-generated approval as an evidence source
- `hold_6_inside`
- `mid_close_across`
- side-specific `below_mid` / `above_mid` cohorts
- `hold_3_inside_mid_50`
- any live, paper, testnet, scoring, sizing, or strategy-promotion claim

## Next Step

Build the first minimal offline entry prototype only around the revalidated
futures surface. The prototype must report side splits, period splits, trade
counts, P&L after costs, drawdown, and stress behavior before any further
promotion claim. If that prototype fails split-stable P&L or exposes a weak
side, stop this detector family and pivot to a materially different
non-trading hypothesis.
