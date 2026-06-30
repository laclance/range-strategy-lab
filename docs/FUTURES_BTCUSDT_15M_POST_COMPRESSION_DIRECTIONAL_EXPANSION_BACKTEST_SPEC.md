# Futures BTCUSDT 15m Post-Compression Directional Expansion Backtest Spec

Date: 2026-06-30

## Verdict

Stop state:
`post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval`.

The user explicitly approved this docs-only offline backtest spec for the
selected representative candidate:

```text
btc_15m_post_compression_l192_q20_m020_none_long_h48_v1
```

This spec is ready for a separate implementation approval gate. It defines one
candidate stream, one fixed risk/exit model, expected outputs, and hard
falsification gates for a later offline backtest implementation. It does not
authorize Go code, CLI flags, generated outputs, a backtest run, replay,
walk-forward, optimizer selection, derivatives-veto interaction,
paper/testnet/live paths, or promotion.

## Evidence Status

Current evidence is zero-trade label separation only, not P&L.

The approved zero-trade audit passed at:

```text
btc_15m_post_compression_directional_expansion_zero_trade_audit_passed_needs_review
```

The docs-only strategy-premise spec then selected exactly one representative
candidate from the passing pocket and stopped at:

```text
post_compression_directional_expansion_strategy_premise_spec_ready_for_backtest_approval
```

The only eligible zero-trade evidence pocket was:

| Field | Eligible evidence |
| --- | --- |
| Side | `long` only |
| Horizon | `48` closed `15m` bars only |
| Compression lookback | `192` prior closed `15m` bars |
| Compression threshold | bottom `20%` of prior range width |
| Breakout threshold | `0.1`, `0.2`, `0.3` prior-bar `ATR(14)` |
| Volume mode | `none`, `above_prior_96_median`, `above_prior_96_p60` |

The selected representative cell is the center/no-extra-filter cell:

| Dimension | Value |
| --- | --- |
| Candidate id | `btc_15m_post_compression_l192_q20_m020_none_long_h48_v1` |
| Compression lookback | `192` prior closed `15m` bars |
| Compression threshold | bottom `20%` |
| Breakout threshold | `0.2` prior-bar `ATR(14)` |
| Volume confirmation | `none` |
| Side | `long` |
| Evidence horizon | `48` closed `15m` bars |

The adjacent passing cells remain supporting robustness evidence only. They may
be referenced in the later review, but they may not be converted into a P&L
grid, selected after result inspection, or used to rescue a failed
representative-cell backtest.

## Candidate Construction

The later implementation may build only this candidate stream.

Use exact closed UTC `15m` candles resampled from complete local BTCUSDT
Binance USDT-M futures `5m` children. For each closed `15m` decision candle
`d`:

1. Compute the prior range over closed `15m` bars `[d-192, d-1]`.
2. `range_high(d)` is the maximum high in that prior window.
3. `range_low(d)` is the minimum low in that prior window.
4. `range_width_pct(d)` is
   `(range_high(d) - range_low(d)) / close[d-1]`.
5. Build the compression reference from the previous `1,920` valid
   `range_width_pct` observations ending at `d-1`.
6. Require `range_width_pct(d)` at or below the `20%` percentile of that prior
   reference set.
7. Compute Wilder-style `ATR(14)` on the resampled `15m` candles and use only
   `ATR(14)[d-1]`.
8. Require `close[d] >= range_high(d) + 0.2 * ATR(14)[d-1]`.
9. Apply no volume confirmation filter.
10. Emit a long-only signal candidate for the next `15m` bar open.

The implementation must reproduce the zero-trade audit's representative-cell
candidate count before position-overlap filtering:

```text
l192_q20_m020_none long h48 full rows = 468
```

Any mismatch must be explained in an implementation review and fail closed if
it comes from source drift, resample drift, lookahead, missingness policy
change, parameter drift, or label/candidate confusion.

## Source And Resample Contract

The later implementation must carry forward the accepted source contract:

| Field | Value |
| --- | ---: |
| Source path | `../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv` |
| Product | Binance USDT-M futures |
| Symbol | BTCUSDT |
| Interval | `5m` |
| Loaded candles | `573,984` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:55:00Z` |
| Gaps | `0` |
| Duplicates | `0` |
| Zero-volume rows | `66` |
| Comparison-only | `false` |
| Validation status | `accepted` |

Accepted closed UTC `15m` resample facts:

| Field | Value |
| --- | ---: |
| Row count | `191,328` |
| First open | `2021-01-01T00:00:00Z` |
| Last open | `2026-06-16T23:45:00Z` |
| Last close | `2026-06-16T23:59:59Z` |
| Expected child bars | `3` |
| Missing child opens | `0` |
| Validation status | `accepted` |

The later implementation must reject spot, comparison-only sources,
non-BTCUSDT sources, non-`5m` sources, missing child `15m` buckets, gaps,
duplicates, irregular cadence, invalid OHLC, negative volume, and source
coverage drift unless a separate approved source-impact review changes scope.

## Fixed Backtest Model

The later implementation may use only this fixed model:

| Component | Decision |
| --- | --- |
| Direction | long only |
| Entry timing | next `15m` bar open after decision candle `d` |
| Position limit | one open position max |
| Stop | `entry_price - 1.0 * ATR(14)[d-1]` |
| Target | `entry_price + 2.0 * ATR(14)[d-1]` |
| Max hold | `48` closed `15m` bars after entry |
| Exit ambiguity | stop-first when stop and target are touched in the same bar |
| Time exit | close of the bar where `hold_bars >= 48`, using normal exit slippage |
| Sizing | `1%` risk-at-stop sizing |
| Notional cap | `1x` equity |
| Starting balance | `1000` |
| Fee | `0.0004` per side |
| Slippage | `0.000116` per side |

For a long signal, `entry_price` is the slipped next-open execution price under
the repo's existing slippage convention. Stop and target are fixed from that
entry price and `ATR(14)[d-1]`. The implementation must skip a signal if the
entry candle is missing, `ATR(14)[d-1]` is missing or non-positive, or stop and
target geometry is invalid.

This fixed model is deliberately simple. It is not an exit optimizer and it is
not a claim that `1R`/`2R` is optimal. It is the first executable stress test of
whether the zero-trade label separation survives a conservative, bounded
offline backtest.

## No-Lookahead Rules

The later implementation must preserve these rules:

- prior range uses `[d-192, d-1]`;
- compression percentile references exclude `d` and require exactly `1,920`
  prior valid observations;
- ATR uses `ATR(14)[d-1]`;
- decision close is known only at `d` close;
- entry uses only `open[d+1]`;
- stop and target may use `open[d+1]` only as the actual execution anchor, not
  as an entry filter;
- future highs/lows, forward returns, exit outcomes, split membership after
  close, and P&L may never shape candidate construction;
- missing required values produce explicit skip/missingness rows, never silent
  defaults.

## Expected CLI And Output Scope

If separately approved, the implementation should add exactly one offline flag:

```text
-futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest
```

Default output directory:

```text
results/futures-btc-15m-post-compression-l192-q20-m020-none-long-h48-backtest/
```

Expected artifacts:

- `source_manifest.json`;
- `summary.json`, `summary.csv`, and `trades.json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_sources.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_resample_coverage.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_signals.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_skips.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_trades.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_summary.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_cost_stress.csv/json`;
- `btc_15m_post_compression_l192_q20_m020_none_long_h48_falsification.json`.

The strategy-specific summary is the review authority. Common `summary.*` and
`trades.json` must remain compatibility views of the same executed trades.

## Metrics

The implementation review must report these metrics by `full_2021_2026`,
`2021_2022_stress`, `2023_2024_oos`, and `2025_2026_recent`:

- raw candidate rows before one-position filtering;
- skipped rows by reason;
- executed trades;
- win rate;
- gross P&L;
- engine net P&L using slipped execution prices and fees;
- extra slippage-stress net P&L;
- total fees;
- total slippage;
- profit factor using extra slippage-stress net;
- gross profit factor;
- max drawdown;
- average gross R;
- average engine net R;
- average extra slippage-stress net R;
- average initial risk;
- exit-reason counts for `stop_loss`, `take_profit`, `time_stop`, and
  `force_close`.

The current engine applies slipped execution prices and then subtracts fees in
`NetPnL`; it also records slippage as a separate cost proxy. Therefore the
later implementation must explicitly report both the current engine net and an
extra conservative stress view that subtracts the recorded slippage proxy once
more. Promotion or failure decisions must use that extra slippage-stress view.

## Pass And Fail Gates

The later backtest implementation may pass only if every gate below clears.
Failure is a valid result and must not be rescued by adding filters, changing
parameters, using adjacent cells, or applying the derivatives veto.

| Gate | Required result |
| --- | --- |
| Source/resample | accepted source contract and exact closed UTC `15m` resample reproduced |
| Candidate identity | `468` representative-cell raw candidate rows reproduced before one-position filtering, or a fail-closed documented source/resample/candidate mismatch |
| Leakage | no feature, threshold, entry filter, stop, target, or sizing value uses data after its allowed timestamp |
| Trade count | at least `120` full-sample executed trades and at least `25` executed trades in each primary split |
| Gross edge | full-sample gross P&L positive, with `2023_2024_oos` and `2025_2026_recent` gross P&L non-negative |
| Costed edge | full-sample extra slippage-stress net P&L positive, with `2023_2024_oos` and `2025_2026_recent` stress net P&L non-negative |
| Profit factor | full-sample extra slippage-stress PF at least `1.20`; `2023_2024_oos` and `2025_2026_recent` stress PF at least `1.05` |
| Drawdown | full-sample max drawdown no worse than `25%`; each primary split no worse than `30%` |
| Robustness | pass does not depend on a single split, a single exit reason, a handful of trades, or positive results only in `2021_2022_stress` |
| Optimizer protection | no stop/target/hold/cell/volume/side/veto grid or post-result parameter selection |
| Closed-family protection | no detector episodes, structured compression fields, clean breakout continuation, router rotation, occupancy rotation, midline/hold-inside, HTF nested rotation, BTC/ETH/SOL context, or old compression-breakout rescue |
| Veto protection | `btc_15m_basis_discount_no_trade_veto_v1` is absent from entries, exits, side, ranking, P&L, and pass/fail decisions |

If any required gate fails, the implementation review must stop at a failed
backtest state. It must not request a replay, walk-forward, veto interaction,
or strategy promotion from a failing fixed backtest.

## Implementation Stop States

Allowed later implementation stop states:

```text
post_compression_directional_expansion_backtest_passed_needs_review
post_compression_directional_expansion_backtest_failed_no_usable_strategy
post_compression_directional_expansion_backtest_failed_source_or_resample
post_compression_directional_expansion_backtest_rejected_optimizer_contamination
post_compression_directional_expansion_backtest_rejected_closed_family_reslice
post_compression_directional_expansion_backtest_rejected_veto_contamination
```

A passing backtest stop state would still not authorize replay, walk-forward,
veto interaction, strategy promotion, or paper/testnet/live work. It would only
authorize a separate docs-only review of whether the fixed backtest result is
strong enough to request the next stage.

## Rejections And Boundaries

This backtest spec explicitly rejects:

- short-side post-compression entries;
- `16`-bar or `32`-bar label evidence as strategy authority;
- lookback `48` or `96`;
- compression thresholds `30%` or `40%`;
- breakout thresholds `0.1` or `0.3` as implementation candidates;
- volume confirmation modes as implementation candidates;
- adjacent-cell P&L selection or rescue;
- the full `81`-cell parameter grid;
- derivatives veto interaction;
- old detector episodes;
- router rows;
- occupancy rotation;
- midline/hold-inside states;
- BTC/ETH/SOL context;
- source or symbol expansion.

The canonical derivatives veto remains parked as future skip/retain evidence
only:

```text
btc_15m_basis_discount_no_trade_veto_v1
```

It may not shape candidate rows, create entries, choose side, rank, score P&L,
replay, walk forward, optimize, promote a strategy, or reopen closed families.
Any future veto interaction audit requires a separate approval after an
independent backtest candidate exists and after the fixed backtest result is
reviewed.

## Next Gate

The next allowed step is explicit user approval for implementation of the
offline backtest defined here:

```text
btc_15m_post_compression_l192_q20_m020_none_long_h48_v1
```

Until that approval is given, do not add Go code, CLI flags, generated outputs,
audit/backtest runs, source downloads, P&L artifacts, replay, walk-forward,
derivatives-veto interaction, or promotion.

## Stop States

Used:

```text
post_compression_directional_expansion_backtest_spec_ready_for_implementation_approval
```

Reserved for implementation:

```text
post_compression_directional_expansion_backtest_passed_needs_review
post_compression_directional_expansion_backtest_failed_no_usable_strategy
post_compression_directional_expansion_backtest_failed_source_or_resample
post_compression_directional_expansion_backtest_rejected_optimizer_contamination
post_compression_directional_expansion_backtest_rejected_closed_family_reslice
post_compression_directional_expansion_backtest_rejected_veto_contamination
```
