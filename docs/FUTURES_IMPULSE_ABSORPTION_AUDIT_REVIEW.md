# Futures Impulse Absorption Audit Review

Date: 2026-06-25

## Verdict

Stop state:
`impulse_absorption_no_viable_edge`.

The futures impulse absorption premise does not pass the non-trading review
gate. The audit found plenty of abnormal OHLCV impulse candles, but the
post-event behavior is continuation-dominant rather than absorption-dominant.
Across every period split and every tested horizon, continuation beyond the
event extreme happens before midpoint reclaim far more often than midpoint
reclaim happens first.

No entry prototype, exit rule, sizing change, strategy replacement,
paper/testnet/live wiring, API key, deploy script, grid, martingale, averaging
down, or two-exchange work is approved from this result. `cmd/rangelab` still
uses `lab.EmptyStrategy` by default, and this audit produced zero trades.

## Source

| Field | Value |
| --- | --- |
| Path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Venue / product | `Binance` / `Binance USDT-M futures` |
| Symbol / interval | `BTCUSDT` / `5m` |
| Row count | `573,984` loaded candles |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps / duplicates | `0` / `0` |
| Zero-volume candles | `66` |
| Comparison only | `false` |
| Validation status | `accepted` |

Result directory:
`results/futures-impulse-absorption-audit/`.

Output files:

- `source_manifest.json`
- `summary.csv/json`
- `trades.json`
- `futures_impulse_absorption_candidates.csv/json`
- `futures_impulse_absorption_summary.csv/json`
- `futures_impulse_absorption_stability.csv/json`

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_impulse_absorption_candidates.csv` | `13,737` |
| `futures_impulse_absorption_summary.csv` | `949` |
| `futures_impulse_absorption_stability.csv` | `241` |
| `summary.csv` | `13` |

The audit emitted `13,736` candidate horizon rows, representing `3,434`
distinct impulse events across the four horizons. It found `1,838` up impulse
events and `1,596` down impulse events. No horizon rows had missing future
windows.

## Candidate Definition

Each candidate is a closed BTCUSDT 5m futures candle after a `30` day prior
rolling warmup where:

- true-range percentile rank against the prior `8,640` candles is at least
  `p99`;
- volume percentile rank against the prior `8,640` candles is at least `p95`;
- close position is `>=0.75` for an up impulse or `<=0.25` for a down impulse;
- zero-range candles are skipped.

The event candle is excluded from its own percentile reference window. Labels
are measured over `3`, `6`, `12`, and `24` closed bars after the event.

For an up impulse, midpoint reclaim means a later candle low touches or crosses
the event midpoint; continuation means a later candle high exceeds the event
high. For a down impulse, midpoint reclaim means a later candle high touches or
crosses the event midpoint; continuation means a later candle low falls below
the event low.

## Gate Results

The count criterion passes: every period split has more than `100` all-direction
all-bucket candidates at every horizon. The behavior criteria fail decisively.

All-direction, all-bucket rows:

| Split | Horizon | Events | Reclaim First | Continuation First | Same-Bar Ambiguous | Quick Continuation |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `3` | `1,106` | `24.32%` | `58.41%` | `14.92%` | `77.76%` |
| `2023_2024_oos` | `3` | `1,325` | `24.68%` | `57.74%` | `16.00%` | `79.55%` |
| `2025_2026_recent` | `3` | `1,003` | `24.13%` | `58.42%` | `16.55%` | `80.76%` |
| `2021_2022_stress` | `6` | `1,106` | `25.14%` | `59.40%` | `15.01%` | `77.76%` |
| `2023_2024_oos` | `6` | `1,325` | `25.21%` | `58.34%` | `16.00%` | `79.55%` |
| `2025_2026_recent` | `6` | `1,003` | `24.63%` | `58.72%` | `16.55%` | `80.76%` |
| `2021_2022_stress` | `12` | `1,106` | `25.32%` | `59.58%` | `15.01%` | `77.76%` |
| `2023_2024_oos` | `12` | `1,325` | `25.36%` | `58.42%` | `16.00%` | `79.55%` |
| `2025_2026_recent` | `12` | `1,003` | `24.63%` | `58.72%` | `16.55%` | `80.76%` |
| `2021_2022_stress` | `24` | `1,106` | `25.32%` | `59.58%` | `15.01%` | `77.76%` |
| `2023_2024_oos` | `24` | `1,325` | `25.43%` | `58.57%` | `16.00%` | `79.55%` |
| `2025_2026_recent` | `24` | `1,003` | `24.63%` | `58.72%` | `16.55%` | `80.76%` |

Full-sample direction rows are not a rescue path:

| Direction | Horizon | Events | Reclaim First | Continuation First | Same-Bar Ambiguous | Quick Continuation |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `up` | `3` | `1,838` | `24.27%` | `60.50%` | `13.28%` | `79.11%` |
| `down` | `3` | `1,596` | `24.56%` | `55.45%` | `18.73%` | `79.57%` |
| `up` | `24` | `1,838` | `25.08%` | `61.59%` | `13.33%` | `79.11%` |
| `down` | `24` | `1,596` | `25.25%` | `55.89%` | `18.73%` | `79.57%` |

The stability rows confirm the same failure. The all-direction/all-bucket
minimum reclaim-minus-continuation margin is negative at every horizon, from
`-34.30` percentage points at `h3` to `-34.27` percentage points at `h24`.
The maximum quick-continuation rate is `80.76%` in every horizon family.

## Conclusion

The abnormal OHLCV impulse event is futures-authoritative and useful as an
audit artifact, but the tested absorption premise is not viable. The actual
shape is fast continuation through the event extreme, not stable midpoint
reclaim before continuation.

Do not convert this impulse absorption hypothesis into an entry prototype. Do
not retune it into an entry, paper/testnet/live path, or strategy replacement
without a materially new futures hypothesis or data premise.

Carry forward only the reusable infrastructure:

- futures-only source guard enforcement
- prior-window percentile-rank event detection
- closed-candle horizon labeling
- summary and stability artifact pattern
- explicit no-trade audit flag shape

The next step should ask for a new materially different futures hypothesis or
data premise before opening another audit.
