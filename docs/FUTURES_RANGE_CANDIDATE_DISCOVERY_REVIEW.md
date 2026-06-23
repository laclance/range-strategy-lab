# Futures Range Candidate Discovery Review

Date: 2026-06-26

## Verdict

Stop state:
`range_discovery_audit_ready`.

The broad futures range discovery audit found one candidate family worth
moving to a fixed-rule baseline backtest: clean breakout continuation after a
completed mature range. The strongest non-duplicative next candidates are:

1. `4h` clean breakout continuation, `up` side, `12` bar horizon.
2. `1h` clean breakout continuation, `all` sides, `12` bar horizon.

This is not strategy promotion and not optimization approval. It is permission
to build a first offline baseline backtest for those candidate definitions only.
Default `cmd/rangelab` still uses `lab.EmptyStrategy`, and this audit produced
zero trades.

No live orders, exchange API keys, deploy scripts, paper/testnet wiring, grid,
martingale, averaging down, two-exchange execution, symbol expansion, or data
download is approved.

## Source And Coverage

Parent source:

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
`results/futures-range-candidate-discovery-audit/`.

Generated timeframe coverage:

| Timeframe | Rows | First Open | Last Open | Gaps | Duplicates | Missing Child Opens | Status |
| --- | ---: | --- | --- | ---: | ---: | ---: | --- |
| `5m` | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` | `0` | `0` | accepted |
| `15m` | `191,328` | `2021-01-01T00:00:00Z` | `2026-06-16T23:45:00Z` | `0` | `0` | `0` | accepted |
| `1h` | `47,832` | `2021-01-01T00:00:00Z` | `2026-06-16T23:00:00Z` | `0` | `0` | `0` | accepted |
| `4h` | `11,958` | `2021-01-01T00:00:00Z` | `2026-06-16T20:00:00Z` | `0` | `0` | `0` | accepted |

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `futures_range_discovery_candidates.csv` | `144,637` |
| `futures_range_discovery_coverage.csv` | `5` |
| `futures_range_discovery_rankings.csv` | `121` |
| `futures_range_discovery_stability.csv` | `121` |
| `futures_range_discovery_summary.csv` | `481` |
| `summary.csv` | `13` |

The audit emitted `144,636` event/horizon rows representing `48,212` distinct
events across the three horizons. Common `summary.*` and `trades.json` remained
zero-trade.

## Discovery Results

The audit ranked `120` candidate surfaces. `24` passed the balanced discovery
gate, and every passing row came from `clean_breakout_continuation`.

Passing rows by timeframe:

| Timeframe | Passing Rows |
| --- | ---: |
| `15m` | `9` |
| `1h` | `9` |
| `4h` | `6` |

Passing rows by side:

| Side | Passing Rows |
| --- | ---: |
| `all` | `9` |
| `up` | `9` |
| `down` | `6` |

Top rankings:

| Rank | Candidate | Full Count | Weakest Split Count | Full Favorable | Weakest Favorable | Worst Adverse | Worst Quick Invalid. | Weakest Cost Buffer |
| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `1` | `4h clean_breakout_continuation up h12` | `496` | `112` | `86.69%` | `82.38%` | `16.58%` | `19.64%` | `4.2910%` |
| `2` | `4h clean_breakout_continuation up h6` | `496` | `112` | `86.90%` | `82.38%` | `15.18%` | `19.64%` | `4.0041%` |
| `3` | `1h clean_breakout_continuation all h12` | `4,352` | `1,213` | `90.37%` | `89.01%` | `10.99%` | `17.51%` | `2.5907%` |

The top two ranking rows are the same `4h` up-breakout candidate at different
horizons, so the next backtest should avoid duplicating that surface. Use the
best `4h` up candidate and the best sample-rich `1h` all-side candidate.

Split detail for the primary `4h` up `h12` candidate:

| Split | Count | Favorable | Adverse | Quick Invalid. | Avg Fav. Move | Avg Adv. Move | Cost Buffer |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `191` | `92.67%` | `7.33%` | `13.61%` | `6.8153%` | `0.5282%` | `6.1839%` |
| `2023_2024_oos` | `193` | `82.38%` | `16.58%` | `13.47%` | `5.1852%` | `0.7909%` | `4.2910%` |
| `2025_2026_recent` | `112` | `83.93%` | `16.07%` | `19.64%` | `6.1052%` | `0.2963%` | `5.7057%` |
| `full_2021_2026` | `496` | `86.69%` | `12.90%` | `14.92%` | `6.0206%` | `0.5781%` | `5.3394%` |

Split detail for the secondary `1h` all-side `h12` candidate:

| Split | Count | Favorable | Adverse | Quick Invalid. | Avg Fav. Move | Avg Adv. Move | Cost Buffer |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `1,683` | `91.74%` | `8.26%` | `13.01%` | `4.0804%` | `0.4342%` | `3.5429%` |
| `2023_2024_oos` | `1,456` | `89.01%` | `10.99%` | `17.51%` | `2.9812%` | `0.2873%` | `2.5907%` |
| `2025_2026_recent` | `1,213` | `90.11%` | `9.89%` | `16.49%` | `3.0622%` | `0.2383%` | `2.7207%` |
| `full_2021_2026` | `4,352` | `90.37%` | `9.63%` | `15.49%` | `3.4288%` | `0.3305%` | `2.9952%` |

The `15m` clean breakout rows also pass, but they rank behind the `4h` and
`1h` shapes. They should remain a later comparison, not the first baseline
backtest target.

## Rejected Families

The other range-first families did not pass the discovery gate:

| Family | Main Failure |
| --- | --- |
| `boundary_touch_rejection` | Adverse acceptance outside the range beat inward rejection across full-sample rows. |
| `single_candle_wick_rejection` | Boundary breaks and quick invalidation dominated wick-rejection follow-through. |
| `failed_breakout_reentry` | Second breaks/adverse continuation beat re-entry continuation. |
| `mature_balance_persistence` | Expansion failure dominated persistence despite some positive excursion buffers. |

These failures are non-trading discovery results, not permanent bans on every
possible future touch or reaction idea. They do mean the next implementation
should not backtest touch/rejection/re-entry/balance variants from this audit.

## Next Step

Create a baseline offline backtest for:

- `4h clean_breakout_continuation up h12`;
- `1h clean_breakout_continuation all h12`.

The next baseline must be fixed-rule and non-optimized:

- derive `1h` and `4h` bars from the accepted `5m` futures parent source;
- use completed mature range episodes only;
- signal on a closed break beyond the completed range boundary;
- enter on the next higher-timeframe bar open;
- skip invalid entry geometry;
- use one simple stop/target/time-stop template per candidate;
- keep the default run on `lab.EmptyStrategy`;
- write explicit signal, trade, and summary artifacts;
- review gross/net P&L before any optimization.

Do not optimize, broaden symbols, add live-adjacent wiring, or convert the
discovery labels directly into a deployed strategy.
