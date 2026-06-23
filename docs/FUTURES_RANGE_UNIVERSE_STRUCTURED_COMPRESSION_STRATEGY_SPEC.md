# Futures Range Universe Structured Compression Strategy Spec

Date: 2026-06-26

## Verdict

Stop state:
`structured_compression_strategy_spec_ready_for_offline_replay`.

The selected `4h` structured-compression optimization result is
implementation-ready for one offline candidate strategy replay/backtest. The
strategy authority is ETHUSDT and SOLUSDT only. BTCUSDT remains diagnostic-only
because the selected BTC diagnostic stream was negative in every period split.

This spec freezes the selected configuration:
`sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00`.

The next milestone may implement a fixed offline replay/backtest for this exact
configuration. It must not run a new grid, add new symbols, add new timeframes,
promote BTCUSDT, or add live, paper, testnet, exchange API, deployment,
credential, grid, martingale, averaging down, or two-exchange logic.

## Source Contract

Authority sources:

| Symbol | Path | Role | Loaded Candles | First Open | Last Open | Gaps / Duplicates | Zero Volume | Physical Non-Monotonic |
| --- | --- | --- | ---: | --- | --- | --- | ---: | ---: |
| `ETHUSDT` | `../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv` | authority | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` / `0` | `47` | `0` |
| `SOLUSDT` | `../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv` | authority | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` / `0` | `47` | `1` |

Diagnostic source:

| Symbol | Path | Role | Loaded Candles | First Open | Last Open | Gaps / Duplicates | Zero Volume | Physical Non-Monotonic |
| --- | --- | --- | ---: | --- | --- | --- | ---: | ---: |
| `BTCUSDT` | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` | diagnostic only | `573,984` | `2021-01-01T00:00:00Z` | `2026-06-16T23:55:00Z` | `0` / `0` | `66` | `0` |

All sources must be Binance USDT-M futures `5m` closed candles. Spot sources,
comparison-only manifests, non-approved symbols, source downloads, sibling-repo
mutation, and exchange API access are out of scope.

SOLUSDT may be sorted by open time during validation because the accepted local
file has one physical non-monotonic row. It may be used only if the sorted
stream has no duplicates, gaps, invalid OHLCV, or bad finality.

## Resampling Contract

The replay uses closed UTC `4h` bars derived from the accepted `5m` parent
sources.

| Rule | Value |
| --- | --- |
| Child bars per `4h` bar | `48` |
| Expected rows per symbol | `11,958` |
| First `4h` open | `2021-01-01T00:00:00Z` |
| Last `4h` open | `2026-06-16T20:00:00Z` |
| Open | first child open |
| High | maximum child high |
| Low | minimum child low |
| Close | last child close |
| Volume | sum of child volume |

Every accepted `4h` bucket must contain each expected child open exactly once.
Partial final buckets, duplicate child opens, missing child opens, forward
fills, synthetic candles, local-time buckets, and future-looking semantics are
rejections.

## Frozen Candidate

| Field | Value |
| --- | --- |
| Candidate | `structured_compression_4h_all_h12` for replay naming |
| Source optimization config | `sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00` |
| Family | `structured_compression_expansion` |
| Timeframe | closed UTC `4h` |
| Detector | `p30_c12_bollinger_on_adx_off` |
| Detector lookback | `20` days |
| Detector percentile | `0.30` |
| Detector min consecutive bars | `12` |
| Bollinger | on |
| ADX | off |
| Sides | long and short |
| Event delay | first closed breakout within `24` closed `4h` bars after completed mature range end |
| Confirmation window | `2` closed `4h` bars after breakout |
| Entry | next `4h` bar open after the closed confirmation candle |
| Target | `1.0` completed range width from slipped entry |
| Stop buffer | `0.0` completed range width |
| Max hold | `12` closed `4h` bars |

The replay should name the implementation surface clearly enough to distinguish
it from the old baseline candidate `structured_compression_4h_all_h6`. The
source candidate came from that family, but the optimized max hold is now
`12`.

## Event And Signal Rules

Use completed mature range episodes from the detector only. Candidate
boundaries come from completed range episodes and may not inspect future
candles.

For each completed episode:

1. Scan forward from the first closed `4h` bar after episode end.
2. Stop scanning after `24` closed `4h` bars.
3. For an up breakout, require a closed candle with close above the completed
   range high.
4. For a down breakout, require a closed candle with close below the completed
   range low.
5. After the first breakout, search only the next `2` closed `4h` bars for
   confirmation.
6. Up confirmation requires close above range high and high above the breakout
   candle high.
7. Down confirmation requires close below range low and low below the breakout
   candle low.
8. Signal on the closed confirmation candle.
9. Use only the first valid structured-compression signal per completed
   episode.

Skip signals with:

- missing next-bar entry candle;
- non-positive range width;
- non-positive or non-finite entry, stop, or target;
- invalid stop/target geometry after entry slippage;
- duplicate signal index for the same symbol and configuration.

## Trade Template

The replay inherits the current lab engine assumptions:

- starting balance: CLI default unless explicitly changed by the run;
- risk at stop: CLI default `1%`;
- max notional: CLI default `1x` equity;
- fee: CLI default `0.0004` per side;
- slippage: CLI default `0.000116` per side;
- one open position max inside each symbol replay;
- stop-first same-bar ambiguity;
- next-bar-open entries only.

Long trade:

- entry at next `4h` open with entry slippage;
- stop at completed range high because the high boundary is the broken
  same-side boundary;
- target at slipped entry plus one completed range width;
- time stop after `12` closed `4h` bars.

Short trade:

- entry at next `4h` open with entry slippage;
- stop at completed range low because the low boundary is the broken
  same-side boundary;
- target at slipped entry minus one completed range width;
- time stop after `12` closed `4h` bars.

## Portfolio-Stream Semantics

The authority strategy stream is the combined ETHUSDT plus SOLUSDT replay.
Authority aggregate rows, common `trades.json`, and common `summary.*` should
include only ETHUSDT and SOLUSDT trades.

BTCUSDT may be replayed in the same run for diagnostics, with rows marked as
diagnostic. BTCUSDT must not:

- determine pass/fail authority;
- enter common strategy outputs;
- be promoted as a BTC strategy from this spec;
- offset ETH/SOL weakness in review math.

This is still an offline research stream, not a live portfolio coordinator.
Do not add cross-symbol capital allocation, exchange routing, position
netting, or overlap optimization in the replay milestone.

## Expected Baseline Facts

The optimization selected authority result was:

| Scope | Trades | Gross P&L | Net P&L | PF | Max DD | Avg Net R |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| ETH/SOL authority full sample | `129` | `641.05` | `573.87` | `1.8089` | `9.82%` | `0.3465` |
| `2021_2022_stress` | `54` | `174.57` | `151.79` | `1.4867` | `9.82%` | `0.1673` |
| `2023_2024_oos` | `43` | `252.13` | `229.02` | `2.2318` | `5.05%` | `0.4585` |
| `2025_2026_recent` | `32` | `214.34` | `193.06` | `1.9121` | `9.75%` | `0.4983` |

Authority symbol caveats:

- ETHUSDT full sample was positive, but its `2025_2026_recent` split was
  negative after costs.
- SOLUSDT full sample was positive, but its `2021_2022_stress` split was
  negative after costs.
- BTCUSDT diagnostic-only was negative full sample: `55` trades, net P&L
  `-100.67`, PF `0.6507`.

The replay does not need bit-identical floating point formatting, but material
differences in trade count, source coverage, selected rows, or pass/fail
metrics should stop as a replay mismatch until explained.

## Replay Review Gate

The next replay/backtest can stay promoted only if all of these hold:

- every source and `4h` resample validates as accepted;
- authority outputs include ETHUSDT and SOLUSDT only;
- BTCUSDT remains diagnostic-only even if emitted;
- full authority trades are at least `100`;
- `2023_2024_oos` and `2025_2026_recent` authority trades are at least `25`;
- full authority net P&L is positive after costs;
- full authority PF is at least `1.2`;
- stress, OOS, and recent authority splits are positive or no worse than the
  optimization review within a documented replay-tolerance explanation;
- long and short authority sides are both positive or no side loss dominates;
- ETHUSDT and SOLUSDT each retain positive full-sample net P&L and PF at least
  `1.0`;
- BTCUSDT diagnostic weakness is reported and not used as authority;
- max drawdown is not materially worse than the optimization selected stream;
- skipped signals are counted and explained.

Failing these gates should stop the replay as a regression or no-promotion
review. Do not tune around a failure inside the replay milestone.

## Next Implementation Brief

The next brief should implement:
`-futures-range-universe-structured-compression-strategy-replay`.

It should write fixed replay artifacts under
`results/futures-range-universe-structured-compression-strategy-replay/`,
document the review in
`docs/FUTURES_RANGE_UNIVERSE_STRUCTURED_COMPRESSION_STRATEGY_REPLAY_REVIEW.md`,
and keep default `cmd/rangelab` behavior on `lab.EmptyStrategy`.
