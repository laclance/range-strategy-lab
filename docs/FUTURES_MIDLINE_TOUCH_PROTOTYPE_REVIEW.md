# Futures Midline Touch Prototype Review

Date: 2026-06-25

## Verdict

Stop state:
`prototype_failed_no_promotion`.

The first minimal offline entry prototype for the futures-revalidated
`hold_3_inside + mid_touch + mid_50` surface failed. The reaction audit
survived on Binance USDT-M futures data, but the first executable close-back
trade template did not convert that diagnostic behavior into P&L.

No strategy promotion, paper/testnet/live work, sizing change, deployment, API
key, exchange wiring, grid, martingale, averaging down, or two-exchange work is
approved from this result. The default `cmd/rangelab` path still uses
`lab.EmptyStrategy`; prototype trades appear only with the explicit
`-hold-inside-midline-touch-prototype` flag.

## Prototype Tested

Source:

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

Trade template:

- detector profile `p30_c12_bollinger_on_adx_off`
- context rule `hold_3_inside`
- first closed-candle `mid_touch` within `12` bars after the hold decision
- event close-position bucket `mid_50`
- close-back side model:
  - event close below midline -> long
  - event close above midline -> short
  - event close exactly at midline -> skipped
- stop at same-side frozen range boundary
- target at opposite frozen range boundary
- entry on next candle open
- time stop after `6` bars
- existing engine costs, risk-at-stop sizing, one-position limit, and stop-first
  ambiguity

Outputs:

- `results/futures-hold-inside-midline-touch-prototype/source_manifest.json`
- `results/futures-hold-inside-midline-touch-prototype/summary.csv`
- `results/futures-hold-inside-midline-touch-prototype/trades.json`
- `results/futures-hold-inside-midline-touch-prototype/hold_inside_midline_touch_prototype_signals.csv`
- `results/futures-hold-inside-midline-touch-prototype/hold_inside_midline_touch_prototype_trades.csv`
- `results/futures-hold-inside-midline-touch-prototype/hold_inside_midline_touch_prototype_summary.csv`

CSV line counts including headers:

| File | Lines |
| --- | ---: |
| `hold_inside_midline_touch_prototype_signals.csv` | `533` |
| `hold_inside_midline_touch_prototype_trades.csv` | `532` |
| `hold_inside_midline_touch_prototype_summary.csv` | `13` |
| `summary.csv` | `13` |

The prototype emitted `532` signal rows: `531` executed trades and one
documented exact-mid skip. Exit reasons were `140` stop losses, `82` take
profits, and `309` time stops.

## P&L Review

The gate fails at the broadest level: full-sample gross P&L, net P&L, profit
factor, every period split, and both sides are negative after costs.

| Split | Trades | Win Rate | Gross P&L | Net P&L | PF | Avg Net R | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `2021_2022_stress` | `220` | `32.73%` | `-71.44` | `-226.03` | `0.3824` | `-0.3655` | `22.88%` |
| `2023_2024_oos` | `198` | `29.80%` | `-3.80` | `-116.58` | `0.3165` | `-0.4323` | `11.79%` |
| `2025_2026_recent` | `113` | `24.78%` | `-20.31` | `-76.38` | `0.2298` | `-0.5402` | `7.64%` |
| `full_2021_2026` | `531` | `29.94%` | `-95.54` | `-418.99` | `0.3409` | `-0.4276` | `42.11%` |

Side split:

| Side | Trades | Win Rate | Gross P&L | Net P&L | PF | Avg Net R | Max DD |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| long | `248` | `29.03%` | `-83.51` | `-234.50` | `0.2805` | `-0.4822` | `23.73%` |
| short | `283` | `30.74%` | `-12.04` | `-184.49` | `0.4045` | `-0.3797` | `18.49%` |

The failure is not a cost-only problem. Full-sample gross P&L is already
negative, and costs widen the loss materially. The recent split is the weakest
by win rate and average net R, while the stress split is the worst by net P&L.

## Conclusion

The futures reaction gate was useful as a diagnostic observation point, but the
first executable close-back entry template is not viable. Do not promote,
parameter-tune, or live-wire this prototype. Do not broaden it to
`hold_6_inside`, `mid_close_across`, side-specific cohorts, or the old
spot-based approval.

Carry forward only the code infrastructure:

- explicit prototype flagging
- futures source guard
- signal/trade/summary artifact shape
- joined event-to-trade detail
- entry-geometry guard in the engine

The next task should stop mining this same hold-inside/midline entry family and
produce a compact futures hypothesis inventory before any new non-trading audit
is chosen.
